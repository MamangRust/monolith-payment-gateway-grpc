package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/email"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-transfer/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transferCommandService struct {
	kafka                     kafka.Kafka
	errorhandler              errorhandler.TransferCommandErrorHandler
	mencache                  mencache.TransferCommandCache
	ctx                       context.Context
	trace                     trace.Tracer
	cardRepository            repository.CardRepository
	saldoRepository           repository.SaldoRepository
	transferQueryRepository   repository.TransferQueryRepository
	transferCommandRepository repository.TransferCommandRepository
	logger                    logger.LoggerInterface
	mapping                   responseservice.TransferResponseMapper
	requestCounter            *prometheus.CounterVec
	requestDuration           *prometheus.HistogramVec
}

func NewTransferCommandService(
	kafka kafka.Kafka,
	ctx context.Context,
	errorhandler errorhandler.TransferCommandErrorHandler,
	mencache mencache.TransferCommandCache,
	cardRepository repository.CardRepository,
	saldoRepository repository.SaldoRepository,
	transferQueryRepository repository.TransferQueryRepository,
	transferCommandRepository repository.TransferCommandRepository,
	logger logger.LoggerInterface,
	mapping responseservice.TransferResponseMapper,
) *transferCommandService {
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

	return &transferCommandService{
		kafka:                     kafka,
		ctx:                       ctx,
		mencache:                  mencache,
		errorhandler:              errorhandler,
		trace:                     otel.Tracer("transfer-command-service"),
		cardRepository:            cardRepository,
		saldoRepository:           saldoRepository,
		transferQueryRepository:   transferQueryRepository,
		transferCommandRepository: transferCommandRepository,
		logger:                    logger,
		mapping:                   mapping,
		requestCounter:            requestCounter,
		requestDuration:           requestDuration,
	}
}

func (s *transferCommandService) CreateTransaction(request *requests.CreateTransferRequest) (*response.TransferResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("CreateTransaction", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "CreateTransaction")
	defer span.End()

	span.SetAttributes(
		attribute.String("transfer_from", request.TransferFrom),
		attribute.String("transfer_to", request.TransferTo),
		attribute.Float64("transfer_amount", float64(request.TransferAmount)),
	)

	s.logger.Debug("Starting create transaction process",
		zap.Any("request", request),
	)

	card, err := s.cardRepository.FindUserCardByCardNumber(request.TransferFrom)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "FindUserCardByCardNumber", "FAILED_FIND_CARD_BY_CARD_NUMBER", span, &status, card_errors.ErrFailedFindByCardNumber, zap.String("card_number", request.TransferFrom))
	}

	_, err = s.cardRepository.FindCardByCardNumber(request.TransferTo)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "FindCardByCardNumber", "FAILED_FIND_CARD_BY_CARD_NUMBER", span, &status, card_errors.ErrFailedFindByCardNumber, zap.String("card_number", request.TransferTo))
	}

	senderSaldo, err := s.saldoRepository.FindByCardNumber(request.TransferFrom)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "FindByCardNumber", "FAILED_FIND_SALDO_BY_CARD_NUMBER", span, &status, saldo_errors.ErrFailedSaldoNotFound, zap.String("card_number", request.TransferFrom))
	}

	receiverSaldo, err := s.saldoRepository.FindByCardNumber(request.TransferTo)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "FindByCardNumber", "FAILED_FIND_SALDO_BY_CARD_NUMBER", span, &status, saldo_errors.ErrFailedSaldoNotFound, zap.String("card_number", request.TransferTo))
	}

	if senderSaldo.TotalBalance < request.TransferAmount {
		return s.errorhandler.HandleSenderInsufficientBalanceError(err, "FindByCardNumber", "FAILED_FIND_SALDO_BY_CARD_NUMBER", span, &status, request.TransferFrom, zap.String("card_number", request.TransferFrom))
	}

	senderSaldo.TotalBalance -= request.TransferAmount
	_, err = s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   senderSaldo.CardNumber,
		TotalBalance: senderSaldo.TotalBalance,
	})
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "UpdateSaldoBalance", "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.String("card_number", request.TransferFrom))
	}

	receiverSaldo.TotalBalance += request.TransferAmount
	_, err = s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   receiverSaldo.CardNumber,
		TotalBalance: receiverSaldo.TotalBalance,
	})
	if err != nil {
		senderSaldo.TotalBalance += request.TransferAmount
		_, rollbackErr := s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
			CardNumber:   senderSaldo.CardNumber,
			TotalBalance: senderSaldo.TotalBalance,
		})
		if rollbackErr != nil {
			return s.errorhandler.HandleRepositorySingleError(rollbackErr, "UpdateSaldoBalance", "FAILED_ROLLBACK_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.String("card_number", request.TransferFrom))
		}

		return s.errorhandler.HandleRepositorySingleError(err, "UpdateSaldoBalance", "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.String("card_number", request.TransferTo))
	}

	transfer, err := s.transferCommandRepository.CreateTransfer(request)
	if err != nil {
		senderSaldo.TotalBalance += request.TransferAmount
		receiverSaldo.TotalBalance -= request.TransferAmount

		_, rollbackErr := s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
			CardNumber:   senderSaldo.CardNumber,
			TotalBalance: senderSaldo.TotalBalance,
		})
		if rollbackErr != nil {
			return s.errorhandler.HandleRepositorySingleError(rollbackErr, "UpdateSaldoBalance", "FAILED_ROLLBACK_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.String("card_number", request.TransferFrom))
		}

		_, rollbackErr = s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
			CardNumber:   receiverSaldo.CardNumber,
			TotalBalance: receiverSaldo.TotalBalance,
		})
		if rollbackErr != nil {
			return s.errorhandler.HandleRepositorySingleError(rollbackErr, "UpdateSaldoBalance", "FAILED_ROLLBACK_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.String("card_number", request.TransferTo))
		}

		if _, err := s.transferCommandRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
			TransferID: transfer.ID,
			Status:     "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, "UpdateTransferStatus", "FAILED_UPDATE_TRANSFER_STATUS", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Int32("transfer_id", int32(transfer.ID)))
		}

		return s.errorhandler.HandleCreateTransferError(err, "CreateTransfer", "FAILED_CREATE_TRANSFER", span, &status, zap.Int32("transfer_id", int32(transfer.ID)))
	}

	res, err := s.transferCommandRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
		TransferID: transfer.ID,
		Status:     "success",
	})

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "UpdateTransferStatus", "FAILED_UPDATE_TRANSFER_STATUS", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Int32("transfer_id", int32(transfer.ID)))
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
		return errorhandler.HandleErrorJSONMarshal[*response.TransferResponse](s.logger, err, "UpdateTransferStatus", "FAILED_UPDATE_TRANSFER_STATUS", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Int32("transfer_id", int32(transfer.ID)))
	}

	err = s.kafka.SendMessage("email-service-topic-transfer-create", strconv.Itoa(res.ID), payloadBytes)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "SendMessage", "FAILED_SEND_EMAIL", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Int32("transfer_id", int32(transfer.ID)))
	}

	so := s.mapping.ToTransferResponse(transfer)

	s.logger.Debug("successfully create transaction",
		zap.Int("transfer_id", transfer.ID),
	)

	return so, nil
}

func (s *transferCommandService) UpdateTransaction(request *requests.UpdateTransferRequest) (*response.TransferResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("UpdateTransaction", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateTransaction")
	defer span.End()

	span.SetAttributes(
		attribute.Int("transfer_id", *request.TransferID),
	)

	s.logger.Debug("Starting update transaction process",
		zap.Int("transfer_id", *request.TransferID),
	)

	transfer, err := s.transferQueryRepository.FindById(*request.TransferID)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "FindById", "FAILED_FIND_TRANSFER", span, &status, transfer_errors.ErrTransferNotFound, zap.Int32("transfer_id", int32(*request.TransferID)))
	}

	amountDifference := request.TransferAmount - transfer.TransferAmount

	senderSaldo, err := s.saldoRepository.FindByCardNumber(transfer.TransferFrom)
	if err != nil {
		if _, err := s.transferCommandRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
			TransferID: *request.TransferID,
			Status:     "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, "UpdateTransferStatus", "FAILED_UPDATE_TRANSFER_STATUS", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Int32("transfer_id", int32(*request.TransferID)))
		}

		return nil, saldo_errors.ErrFailedSaldoNotFound
	}

	newSenderBalance := senderSaldo.TotalBalance - amountDifference
	if newSenderBalance < 0 {

		if _, err := s.transferCommandRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
			TransferID: *request.TransferID,
			Status:     "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, "UpdateTransferStatus", "FAILED_UPDATE_TRANSFER_STATUS", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Int32("transfer_id", int32(*request.TransferID)))
		}

		return s.errorhandler.HandleSenderInsufficientBalanceError(err, "UpdateTransferStatus", "FAILED_UPDATE_TRANSFER_STATUS", span, &status, senderSaldo.CardNumber, zap.Int32("transfer_id", int32(*request.TransferID)))
	}

	senderSaldo.TotalBalance = newSenderBalance
	_, err = s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   senderSaldo.CardNumber,
		TotalBalance: senderSaldo.TotalBalance,
	})
	if err != nil {
		if _, err := s.transferCommandRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
			TransferID: *request.TransferID,
			Status:     "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, "UpdateTransferStatus", "FAILED_UPDATE_TRANSFER_STATUS", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Int32("transfer_id", int32(*request.TransferID)))
		}

		return s.errorhandler.HandleRepositorySingleError(err, "UpdateSaldoBalance", "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Int32("transfer_id", int32(*request.TransferID)))
	}

	receiverSaldo, err := s.saldoRepository.FindByCardNumber(transfer.TransferTo)
	if err != nil {
		rollbackSenderBalance := &requests.UpdateSaldoBalance{
			CardNumber:   transfer.TransferFrom,
			TotalBalance: senderSaldo.TotalBalance + amountDifference,
		}
		_, rollbackErr := s.saldoRepository.UpdateSaldoBalance(rollbackSenderBalance)
		if rollbackErr != nil {
			return s.errorhandler.HandleRepositorySingleError(err, "UpdateSaldoBalance", "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Int32("transfer_id", int32(*request.TransferID)))
		}

		if _, err := s.transferCommandRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
			TransferID: *request.TransferID,
			Status:     "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, "UpdateTransferStatus", "FAILED_UPDATE_TRANSFER_STATUS", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Int32("transfer_id", int32(*request.TransferID)))
		}

		return s.errorhandler.HandleRepositorySingleError(err, "UpdateSaldoBalance", "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Int32("transfer_id", int32(*request.TransferID)))
	}

	newReceiverBalance := receiverSaldo.TotalBalance + amountDifference
	receiverSaldo.TotalBalance = newReceiverBalance
	_, err = s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
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

		if _, err := s.saldoRepository.UpdateSaldoBalance(rollbackSenderBalance); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, "UpdateSaldoBalance", "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Int32("transfer_id", int32(*request.TransferID)))
		}
		if _, err := s.saldoRepository.UpdateSaldoBalance(rollbackReceiverBalance); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, "UpdateSaldoBalance", "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Int32("transfer_id", int32(*request.TransferID)))
		}

		if _, err := s.transferCommandRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
			TransferID: *request.TransferID,
			Status:     "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, "UpdateTransferStatus", "FAILED_UPDATE_TRANSFER_STATUS", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Int32("transfer_id", int32(*request.TransferID)))
		}

		return nil, saldo_errors.ErrFailedUpdateSaldo
	}

	updatedTransfer, err := s.transferCommandRepository.UpdateTransfer(request)
	if err != nil {
		rollbackSenderBalance := &requests.UpdateSaldoBalance{
			CardNumber:   transfer.TransferFrom,
			TotalBalance: senderSaldo.TotalBalance + amountDifference,
		}
		rollbackReceiverBalance := &requests.UpdateSaldoBalance{
			CardNumber:   transfer.TransferTo,
			TotalBalance: receiverSaldo.TotalBalance - amountDifference,
		}

		if _, err := s.saldoRepository.UpdateSaldoBalance(rollbackSenderBalance); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, "UpdateSaldoBalance", "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Int32("transfer_id", int32(*request.TransferID)))
		}
		if _, err := s.saldoRepository.UpdateSaldoBalance(rollbackReceiverBalance); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, "UpdateSaldoBalance", "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Int32("transfer_id", int32(*request.TransferID)))
		}

		if _, err := s.transferCommandRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
			TransferID: *request.TransferID,
			Status:     "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, "UpdateTransferStatus", "FAILED_UPDATE_TRANSFER_STATUS", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Int32("transfer_id", int32(*request.TransferID)))
		}

		return s.errorhandler.HandleRepositorySingleError(err, "UpdateTransfer", "FAILED_UPDATE_TRANSFER", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Int32("transfer_id", int32(*request.TransferID)))
	}

	if _, err := s.transferCommandRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
		TransferID: *request.TransferID,
		Status:     "success",
	}); err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "UpdateTransferStatus", "FAILED_UPDATE_TRANSFER_STATUS", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Int32("transfer_id", int32(*request.TransferID)))
	}

	so := s.mapping.ToTransferResponse(updatedTransfer)

	s.mencache.DeleteTransferCache(*request.TransferID)

	s.logger.Debug("successfully update transaction",
		zap.Int("transfer_id", *request.TransferID),
	)

	return so, nil
}

func (s *transferCommandService) TrashedTransfer(transfer_id int) (*response.TransferResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("TrashedTransfer", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "TrashedTransfer")
	defer span.End()

	span.SetAttributes(
		attribute.Int("transfer_id", transfer_id),
	)

	s.logger.Debug("Starting trashed transfer process",
		zap.Int("transfer_id", transfer_id),
	)

	res, err := s.transferCommandRepository.TrashedTransfer(transfer_id)

	if err != nil {
		return s.errorhandler.HandleTrashedTransferError(err, "TrashedTransfer", "FAILED_TRASHED_TRANSFER", span, &status, zap.Int("transfer_id", transfer_id))
	}

	so := s.mapping.ToTransferResponse(res)

	s.mencache.DeleteTransferCache(transfer_id)

	s.logger.Debug("successfully trashed transfer",
		zap.Int("transfer_id", transfer_id),
	)

	return so, nil
}

func (s *transferCommandService) RestoreTransfer(transfer_id int) (*response.TransferResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreTransfer", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreTransfer")
	defer span.End()

	span.SetAttributes(
		attribute.Int("transfer_id", transfer_id),
	)

	s.logger.Debug("Starting restore transfer process",
		zap.Int("transfer_id", transfer_id),
	)

	res, err := s.transferCommandRepository.RestoreTransfer(transfer_id)

	if err != nil {
		return s.errorhandler.HandleRestoreTransferError(err, "RestoreTransfer", "FAILED_RESTORE_TRANSFER", span, &status, zap.Int("transfer_id", transfer_id))
	}

	so := s.mapping.ToTransferResponse(res)

	s.logger.Debug("successfully restore transfer",
		zap.Int("transfer_id", transfer_id),
	)

	return so, nil
}

func (s *transferCommandService) DeleteTransferPermanent(transfer_id int) (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteTransferPermanent", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteTransferPermanent")
	defer span.End()

	span.SetAttributes(
		attribute.Int("transfer_id", transfer_id),
	)

	s.logger.Debug("Starting delete transfer permanent process",
		zap.Int("transfer_id", transfer_id),
	)

	_, err := s.transferCommandRepository.DeleteTransferPermanent(transfer_id)

	if err != nil {
		return s.errorhandler.HandleDeleteTransferPermanentError(err, "DeleteTransferPermanent", "FAILED_DELETE_TRANSFER_PERMANENT", span, &status, zap.Int("transfer_id", transfer_id))
	}

	s.logger.Debug("successfully delete permanent transfer",
		zap.Int("transfer_id", transfer_id),
	)

	return true, nil
}

func (s *transferCommandService) RestoreAllTransfer() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreAllTransfer", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreAllTransfer")
	defer span.End()

	s.logger.Debug("Restoring all transfers")

	_, err := s.transferCommandRepository.RestoreAllTransfer()

	if err != nil {
		return s.errorhandler.HandleRestoreAllTransferError(err, "RestoreAllTransfer", "FAILED_RESTORE_ALL_TRANSFERS", span, &status, zap.Error(err))
	}

	s.logger.Debug("Successfully restored all transfers")

	return true, nil
}

func (s *transferCommandService) DeleteAllTransferPermanent() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteAllTransferPermanent", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteAllTransferPermanent")
	defer span.End()

	s.logger.Debug("Permanently deleting all transfers")

	_, err := s.transferCommandRepository.DeleteAllTransferPermanent()

	if err != nil {
		return s.errorhandler.HandleDeleteAllTransferPermanentError(err, "DeleteAllTransferPermanent", "FAILED_DELETE_ALL_TRANSFER_PERMANENT", span, &status, zap.Error(err))
	}

	s.logger.Debug("Successfully deleted all transfers permanently")
	return true, nil
}

func (s *transferCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
