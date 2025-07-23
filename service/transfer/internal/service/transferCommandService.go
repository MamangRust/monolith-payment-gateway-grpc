package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/MamangRust/monolith-payment-gateway-pkg/email"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/service"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/service"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/service"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/transfer"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-transfer/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

// transferCommandDeps contains the dependencies required to construct a
// transferCommandService.
type transferCommandDeps struct {
	// Kafka is a pointer to the Kafka client for producing and consuming messages.
	Kafka *kafka.Kafka

	// Ctx is the context for the service.
	Ctx context.Context

	// ErrorHandler is an error handler for the service.
	ErrorHandler errorhandler.TransferCommandErrorHandler

	// Cache is a pointer to the Redis cache for storing data.
	Cache mencache.TransferCommandCache

	// CardRepository is a pointer to the card repository for performing database
	// operations.
	CardRepository repository.CardRepository

	// SaldoRepository is a pointer to the saldo repository for performing database
	// operations.
	SaldoRepository repository.SaldoRepository

	// TransferQueryRepository is a pointer to the transfer query repository for
	// performing database operations.
	TransferQueryRepository repository.TransferQueryRepository

	// TransferCommandRepository is a pointer to the transfer command repository for
	// performing database operations.
	TransferCommandRepository repository.TransferCommandRepository

	// Logger is a pointer to the logger for logging information.
	Logger logger.LoggerInterface

	// Mapper is a pointer to the mapper for mapper the response.
	Mapper responseservice.TransferCommandResponseMapper
}

// transferCommandService provides command-side business logic related to transfers.
// It interfaces with various repositories and services to manage transfer operations.
type transferCommandService struct {
	// kafka is the Kafka client used for producing and consuming messages.
	kafka *kafka.Kafka

	// errorhandler handles errors specific to the transfer command service.
	errorhandler errorhandler.TransferCommandErrorHandler

	// mencache is the Redis cache for storing transfer command data.
	mencache mencache.TransferCommandCache

	// ctx is the context used for managing request lifecycle and cancellation.
	ctx context.Context

	// cardRepository interfaces with the card repository for database operations.
	cardRepository repository.CardRepository

	// saldoRepository interfaces with the saldo repository for database operations.
	saldoRepository repository.SaldoRepository

	// transferQueryRepository interfaces with the transfer query repository for database operations.
	transferQueryRepository repository.TransferQueryRepository

	// transferCommandRepository interfaces with the transfer command repository for database operations.
	transferCommandRepository repository.TransferCommandRepository

	// logger is the logger used for logging information.
	logger logger.LoggerInterface

	// mapper maps responses from domain models to response objects.
	mapper responseservice.TransferCommandResponseMapper

	observability observability.TraceLoggerObservability
}

func NewTransferCommandService(
	params *transferCommandDeps,
) TransferCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transfer_command_service_request_total",
			Help: "Total number of requests to the TransferCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transfer_command_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TransferCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(
		otel.Tracer("transfer-command-service"), params.Logger, requestCounter, requestDuration)

	return &transferCommandService{
		kafka:                     params.Kafka,
		ctx:                       params.Ctx,
		mencache:                  params.Cache,
		errorhandler:              params.ErrorHandler,
		cardRepository:            params.CardRepository,
		saldoRepository:           params.SaldoRepository,
		transferQueryRepository:   params.TransferQueryRepository,
		transferCommandRepository: params.TransferCommandRepository,
		logger:                    params.Logger,
		mapper:                    params.Mapper,
		observability:             observability,
	}
}

// CreateTransaction creates a new transfer transaction.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The request containing transfer details.
//
// Returns:
//   - *response.TransferResponse: The created transfer data.
//   - *response.ErrorResponse: Error details if operation fails.
func (s *transferCommandService) CreateTransaction(ctx context.Context, request *requests.CreateTransferRequest) (*response.TransferResponse, *response.ErrorResponse) {
	const method = "CreateTransacton"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	card, err := s.cardRepository.FindUserCardByCardNumber(ctx, request.TransferFrom)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_CARD_BY_CARD_NUMBER", span, &status, card_errors.ErrFailedFindByCardNumber, zap.Error(err))
	}

	_, err = s.cardRepository.FindCardByCardNumber(ctx, request.TransferTo)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_CARD_BY_CARD_NUMBER", span, &status, card_errors.ErrFailedFindByCardNumber, zap.Error(err))
	}

	senderSaldo, err := s.saldoRepository.FindByCardNumber(ctx, request.TransferFrom)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_SALDO_BY_CARD_NUMBER", span, &status, saldo_errors.ErrFailedSaldoNotFound, zap.Error(err))
	}

	receiverSaldo, err := s.saldoRepository.FindByCardNumber(ctx, request.TransferTo)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_SALDO_BY_CARD_NUMBER", span, &status, saldo_errors.ErrFailedSaldoNotFound, zap.Error(err))
	}

	if senderSaldo.TotalBalance < request.TransferAmount {
		return s.errorhandler.HandleSenderInsufficientBalanceError(err, method, "FAILED_FIND_SALDO_BY_CARD_NUMBER", span, &status, request.TransferFrom, zap.Error(err))
	}

	senderSaldo.TotalBalance -= request.TransferAmount
	_, err = s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
		CardNumber:   senderSaldo.CardNumber,
		TotalBalance: senderSaldo.TotalBalance,
	})
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
	}

	receiverSaldo.TotalBalance += request.TransferAmount
	_, err = s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
		CardNumber:   receiverSaldo.CardNumber,
		TotalBalance: receiverSaldo.TotalBalance,
	})
	if err != nil {
		senderSaldo.TotalBalance += request.TransferAmount
		_, rollbackErr := s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
			CardNumber:   senderSaldo.CardNumber,
			TotalBalance: senderSaldo.TotalBalance,
		})
		if rollbackErr != nil {
			return s.errorhandler.HandleRepositorySingleError(rollbackErr, method, "FAILED_ROLLBACK_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
		}

		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
	}

	transfer, err := s.transferCommandRepository.CreateTransfer(ctx, request)
	if err != nil {
		senderSaldo.TotalBalance += request.TransferAmount
		receiverSaldo.TotalBalance -= request.TransferAmount

		_, rollbackErr := s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
			CardNumber:   senderSaldo.CardNumber,
			TotalBalance: senderSaldo.TotalBalance,
		})
		if rollbackErr != nil {
			return s.errorhandler.HandleRepositorySingleError(rollbackErr, method, "FAILED_ROLLBACK_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
		}

		_, rollbackErr = s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
			CardNumber:   receiverSaldo.CardNumber,
			TotalBalance: receiverSaldo.TotalBalance,
		})
		if rollbackErr != nil {
			return s.errorhandler.HandleRepositorySingleError(rollbackErr, method, "FAILED_ROLLBACK_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
		}

		if _, err := s.transferCommandRepository.UpdateTransferStatus(ctx, &requests.UpdateTransferStatus{
			TransferID: transfer.ID,
			Status:     "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_TRANSFER_STATUS", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Error(err))
		}

		return s.errorhandler.HandleCreateTransferError(err, method, "FAILED_CREATE_TRANSFER", span, &status, zap.Error(err))
	}

	res, err := s.transferCommandRepository.UpdateTransferStatus(ctx, &requests.UpdateTransferStatus{
		TransferID: transfer.ID,
		Status:     "success",
	})

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_TRANSFER_STATUS", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Error(err))
	}

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   "Transfer Successful",
		"Message": fmt.Sprintf("Your Transfer of %d has been processed successfully.", request.TransferAmount),
		"Button":  "View History",
		"Link":    "https://sanedge.example.com/withdraw/history",
	})

	emailPayload := map[string]any{
		"email":   card.Email,
		"subject": "Transfer Successful - SanEdge",
		"body":    htmlBody,
	}

	payloadBytes, err := json.Marshal(emailPayload)
	if err != nil {
		return errorhandler.HandleErrorMarshal[*response.TransferResponse](s.logger, err, method, "FAILED_UPDATE_TRANSFER_STATUS", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Error(err))
	}

	err = s.kafka.SendMessage("email-service-topic-transfer-create", strconv.Itoa(res.ID), payloadBytes)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_SEND_EMAIL", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Error(err))
	}

	so := s.mapper.ToTransferResponse(transfer)

	logSuccess("Transfer created successfully", zap.Bool("success", true))

	return so, nil
}

// UpdateTransaction updates an existing transfer transaction.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The request containing updated transfer details.
//
// Returns:
//   - *response.TransferResponse: The updated transfer data.
//   - *response.ErrorResponse: Error details if operation fails.
func (s *transferCommandService) UpdateTransaction(ctx context.Context, request *requests.UpdateTransferRequest) (*response.TransferResponse, *response.ErrorResponse) {
	const method = "UpdateTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	transfer, err := s.transferQueryRepository.FindById(ctx, *request.TransferID)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_TRANSFER", span, &status, transfer_errors.ErrTransferNotFound, zap.Error(err))
	}

	amountDifference := request.TransferAmount - transfer.TransferAmount

	senderSaldo, err := s.saldoRepository.FindByCardNumber(ctx, transfer.TransferFrom)
	if err != nil {
		if _, err := s.transferCommandRepository.UpdateTransferStatus(ctx, &requests.UpdateTransferStatus{
			TransferID: *request.TransferID,
			Status:     "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_TRANSFER_STATUS", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Error(err))
		}

		return nil, saldo_errors.ErrFailedSaldoNotFound
	}

	newSenderBalance := senderSaldo.TotalBalance - amountDifference
	if newSenderBalance < 0 {

		if _, err := s.transferCommandRepository.UpdateTransferStatus(ctx, &requests.UpdateTransferStatus{
			TransferID: *request.TransferID,
			Status:     "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_TRANSFER_STATUS", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Error(err))
		}

		return s.errorhandler.HandleSenderInsufficientBalanceError(err, method, "FAILED_UPDATE_TRANSFER_STATUS", span, &status, senderSaldo.CardNumber, zap.Error(err))
	}

	senderSaldo.TotalBalance = newSenderBalance
	_, err = s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
		CardNumber:   senderSaldo.CardNumber,
		TotalBalance: senderSaldo.TotalBalance,
	})
	if err != nil {
		if _, err := s.transferCommandRepository.UpdateTransferStatus(ctx, &requests.UpdateTransferStatus{
			TransferID: *request.TransferID,
			Status:     "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_TRANSFER_STATUS", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Error(err))
		}

		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
	}

	receiverSaldo, err := s.saldoRepository.FindByCardNumber(ctx, transfer.TransferTo)
	if err != nil {
		rollbackSenderBalance := &requests.UpdateSaldoBalance{
			CardNumber:   transfer.TransferFrom,
			TotalBalance: senderSaldo.TotalBalance + amountDifference,
		}
		_, rollbackErr := s.saldoRepository.UpdateSaldoBalance(ctx, rollbackSenderBalance)
		if rollbackErr != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
		}

		if _, err := s.transferCommandRepository.UpdateTransferStatus(ctx, &requests.UpdateTransferStatus{
			TransferID: *request.TransferID,
			Status:     "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_TRANSFER_STATUS", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Error(err))
		}

		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
	}

	newReceiverBalance := receiverSaldo.TotalBalance + amountDifference
	receiverSaldo.TotalBalance = newReceiverBalance
	_, err = s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
		CardNumber:   receiverSaldo.CardNumber,
		TotalBalance: receiverSaldo.TotalBalance,
	})
	if err != nil {
		rollbackSenderBalance := &requests.UpdateSaldoBalance{
			CardNumber:   transfer.TransferFrom,
			TotalBalance: senderSaldo.TotalBalance + amountDifference,
		}
		rollbackReceiverBalance := &requests.UpdateSaldoBalance{
			CardNumber:   transfer.TransferTo,
			TotalBalance: receiverSaldo.TotalBalance - amountDifference,
		}

		if _, err := s.saldoRepository.UpdateSaldoBalance(ctx, rollbackSenderBalance); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
		}
		if _, err := s.saldoRepository.UpdateSaldoBalance(ctx, rollbackReceiverBalance); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
		}

		if _, err := s.transferCommandRepository.UpdateTransferStatus(ctx, &requests.UpdateTransferStatus{
			TransferID: *request.TransferID,
			Status:     "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_TRANSFER_STATUS", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Error(err))
		}

		return nil, saldo_errors.ErrFailedUpdateSaldo
	}

	updatedTransfer, err := s.transferCommandRepository.UpdateTransfer(ctx, request)
	if err != nil {
		rollbackSenderBalance := &requests.UpdateSaldoBalance{
			CardNumber:   transfer.TransferFrom,
			TotalBalance: senderSaldo.TotalBalance + amountDifference,
		}
		rollbackReceiverBalance := &requests.UpdateSaldoBalance{
			CardNumber:   transfer.TransferTo,
			TotalBalance: receiverSaldo.TotalBalance - amountDifference,
		}

		if _, err := s.saldoRepository.UpdateSaldoBalance(ctx, rollbackSenderBalance); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
		}
		if _, err := s.saldoRepository.UpdateSaldoBalance(ctx, rollbackReceiverBalance); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
		}

		if _, err := s.transferCommandRepository.UpdateTransferStatus(ctx, &requests.UpdateTransferStatus{
			TransferID: *request.TransferID,
			Status:     "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_TRANSFER_STATUS", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Error(err))
		}

		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_TRANSFER", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Error(err))
	}

	if _, err := s.transferCommandRepository.UpdateTransferStatus(ctx, &requests.UpdateTransferStatus{
		TransferID: *request.TransferID,
		Status:     "success",
	}); err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_TRANSFER_STATUS", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Error(err))
	}

	so := s.mapper.ToTransferResponse(updatedTransfer)

	s.mencache.DeleteTransferCache(ctx, *request.TransferID)

	logSuccess("Successfully update transfer", zap.Int("transfer.id", *request.TransferID))

	return so, nil
}

// TrashedTransfer marks a transfer transaction as trashed (soft delete).
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - transfer_id: The ID of the transfer to be trashed.
//
// Returns:
//   - *response.TransferResponse: The trashed transfer data.
//   - *response.ErrorResponse: Error details if operation fails.
func (s *transferCommandService) TrashedTransfer(ctx context.Context, transfer_id int) (*response.TransferResponseDeleteAt, *response.ErrorResponse) {
	const method = "TrashedTransfer"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.transferCommandRepository.TrashedTransfer(ctx, transfer_id)

	if err != nil {
		return s.errorhandler.HandleTrashedTransferError(err, method, "FAILED_TRASHED_TRANSFER", span, &status, zap.Error(err))
	}

	so := s.mapper.ToTransferResponseDeleteAt(res)

	s.mencache.DeleteTransferCache(ctx, transfer_id)

	logSuccess("Successfully trashed transfer",
		zap.Int("transfer.id", transfer_id),
	)

	return so, nil
}

// RestoreTransfer restores a previously trashed transfer transaction.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - transfer_id: The ID of the transfer to be restored.
//
// Returns:
//   - *response.TransferResponse: The restored transfer data.
//   - *response.ErrorResponse: Error details if operation fails.
func (s *transferCommandService) RestoreTransfer(ctx context.Context, transfer_id int) (*response.TransferResponse, *response.ErrorResponse) {
	const method = "RestoreTransfer"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.transferCommandRepository.RestoreTransfer(ctx, transfer_id)

	if err != nil {
		return s.errorhandler.HandleRestoreTransferError(err, method, "FAILED_RESTORE_TRANSFER", span, &status, zap.Error(err))
	}

	so := s.mapper.ToTransferResponse(res)

	logSuccess("RestoreTransfer process completed", zap.Int("transfer.id", transfer_id))

	return so, nil
}

// DeleteTransferPermanent permanently deletes a transfer transaction.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - transfer_id: The ID of the transfer to be permanently deleted.
//
// Returns:
//   - bool: Whether the deletion was successful.
//   - *response.ErrorResponse: Error details if operation fails.
func (s *transferCommandService) DeleteTransferPermanent(ctx context.Context, transfer_id int) (bool, *response.ErrorResponse) {
	const method = "DeleteTransferPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.transferCommandRepository.DeleteTransferPermanent(ctx, transfer_id)

	if err != nil {
		return s.errorhandler.HandleDeleteTransferPermanentError(err, method, "FAILED_DELETE_TRANSFER_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("DeleteTransferPermanent process completed", zap.Int("transfer.id", transfer_id))

	return true, nil
}

// RestoreAllTransfer restores all trashed transfer transactions.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//
// Returns:
//   - bool: Whether the restore operation was successful.
//   - *response.ErrorResponse: Error details if operation fails.
func (s *transferCommandService) RestoreAllTransfer(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "RestoreAllTransfer"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.transferCommandRepository.RestoreAllTransfer(ctx)

	if err != nil {
		return s.errorhandler.HandleRestoreAllTransferError(err, method, "FAILED_RESTORE_ALL_TRANSFERS", span, &status, zap.Error(err))
	}

	logSuccess("RestoreAllTransfer process completed", zap.Bool("success", true))

	return true, nil
}

// DeleteAllTransferPermanent permanently deletes all trashed transfer transactions.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//
// Returns:
//   - bool: Whether the deletion was successful.
//   - *response.ErrorResponse: Error details if operation fails.
func (s *transferCommandService) DeleteAllTransferPermanent(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "DeleteAllTransferPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.transferCommandRepository.DeleteAllTransferPermanent(ctx)

	if err != nil {
		return s.errorhandler.HandleDeleteAllTransferPermanentError(err, method, "FAILED_DELETE_ALL_TRANSFER_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("DeleteAllTransferPermanent process completed", zap.Bool("success", true))

	return true, nil
}
