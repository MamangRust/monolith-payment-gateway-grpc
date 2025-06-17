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
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transferCommandService struct {
	kafka                     *kafka.Kafka
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
	kafka *kafka.Kafka,
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
	const method = "CreateTransacton"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	card, err := s.cardRepository.FindUserCardByCardNumber(request.TransferFrom)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_CARD_BY_CARD_NUMBER", span, &status, card_errors.ErrFailedFindByCardNumber, zap.Error(err))
	}

	_, err = s.cardRepository.FindCardByCardNumber(request.TransferTo)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_CARD_BY_CARD_NUMBER", span, &status, card_errors.ErrFailedFindByCardNumber, zap.Error(err))
	}

	senderSaldo, err := s.saldoRepository.FindByCardNumber(request.TransferFrom)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_SALDO_BY_CARD_NUMBER", span, &status, saldo_errors.ErrFailedSaldoNotFound, zap.Error(err))
	}

	receiverSaldo, err := s.saldoRepository.FindByCardNumber(request.TransferTo)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_SALDO_BY_CARD_NUMBER", span, &status, saldo_errors.ErrFailedSaldoNotFound, zap.Error(err))
	}

	if senderSaldo.TotalBalance < request.TransferAmount {
		return s.errorhandler.HandleSenderInsufficientBalanceError(err, method, "FAILED_FIND_SALDO_BY_CARD_NUMBER", span, &status, request.TransferFrom, zap.Error(err))
	}

	senderSaldo.TotalBalance -= request.TransferAmount
	_, err = s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   senderSaldo.CardNumber,
		TotalBalance: senderSaldo.TotalBalance,
	})
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
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
			return s.errorhandler.HandleRepositorySingleError(rollbackErr, method, "FAILED_ROLLBACK_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
		}

		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
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
			return s.errorhandler.HandleRepositorySingleError(rollbackErr, method, "FAILED_ROLLBACK_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
		}

		_, rollbackErr = s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
			CardNumber:   receiverSaldo.CardNumber,
			TotalBalance: receiverSaldo.TotalBalance,
		})
		if rollbackErr != nil {
			return s.errorhandler.HandleRepositorySingleError(rollbackErr, method, "FAILED_ROLLBACK_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
		}

		if _, err := s.transferCommandRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
			TransferID: transfer.ID,
			Status:     "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_TRANSFER_STATUS", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Error(err))
		}

		return s.errorhandler.HandleCreateTransferError(err, method, "FAILED_CREATE_TRANSFER", span, &status, zap.Error(err))
	}

	res, err := s.transferCommandRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
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

	so := s.mapping.ToTransferResponse(transfer)

	logSuccess("Transfer created successfully", zap.Bool("success", true))

	return so, nil
}

func (s *transferCommandService) UpdateTransaction(request *requests.UpdateTransferRequest) (*response.TransferResponse, *response.ErrorResponse) {
	const method = "UpdateTransaction"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	transfer, err := s.transferQueryRepository.FindById(*request.TransferID)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_TRANSFER", span, &status, transfer_errors.ErrTransferNotFound, zap.Error(err))
	}

	amountDifference := request.TransferAmount - transfer.TransferAmount

	senderSaldo, err := s.saldoRepository.FindByCardNumber(transfer.TransferFrom)
	if err != nil {
		if _, err := s.transferCommandRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
			TransferID: *request.TransferID,
			Status:     "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_TRANSFER_STATUS", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Error(err))
		}

		return nil, saldo_errors.ErrFailedSaldoNotFound
	}

	newSenderBalance := senderSaldo.TotalBalance - amountDifference
	if newSenderBalance < 0 {

		if _, err := s.transferCommandRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
			TransferID: *request.TransferID,
			Status:     "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_TRANSFER_STATUS", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Error(err))
		}

		return s.errorhandler.HandleSenderInsufficientBalanceError(err, method, "FAILED_UPDATE_TRANSFER_STATUS", span, &status, senderSaldo.CardNumber, zap.Error(err))
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
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_TRANSFER_STATUS", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Error(err))
		}

		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
	}

	receiverSaldo, err := s.saldoRepository.FindByCardNumber(transfer.TransferTo)
	if err != nil {
		rollbackSenderBalance := &requests.UpdateSaldoBalance{
			CardNumber:   transfer.TransferFrom,
			TotalBalance: senderSaldo.TotalBalance + amountDifference,
		}
		_, rollbackErr := s.saldoRepository.UpdateSaldoBalance(rollbackSenderBalance)
		if rollbackErr != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
		}

		if _, err := s.transferCommandRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
			TransferID: *request.TransferID,
			Status:     "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_TRANSFER_STATUS", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Error(err))
		}

		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
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
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
		}
		if _, err := s.saldoRepository.UpdateSaldoBalance(rollbackReceiverBalance); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
		}

		if _, err := s.transferCommandRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
			TransferID: *request.TransferID,
			Status:     "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_TRANSFER_STATUS", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Error(err))
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
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
		}
		if _, err := s.saldoRepository.UpdateSaldoBalance(rollbackReceiverBalance); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
		}

		if _, err := s.transferCommandRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
			TransferID: *request.TransferID,
			Status:     "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_TRANSFER_STATUS", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Error(err))
		}

		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_TRANSFER", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Error(err))
	}

	if _, err := s.transferCommandRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
		TransferID: *request.TransferID,
		Status:     "success",
	}); err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_TRANSFER_STATUS", span, &status, transfer_errors.ErrFailedUpdateTransfer, zap.Error(err))
	}

	so := s.mapping.ToTransferResponse(updatedTransfer)

	s.mencache.DeleteTransferCache(*request.TransferID)

	logSuccess("Successfully update transfer", zap.Int("transfer.id", *request.TransferID))

	return so, nil
}

func (s *transferCommandService) TrashedTransfer(transfer_id int) (*response.TransferResponse, *response.ErrorResponse) {
	const method = "TrashedTransfer"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	res, err := s.transferCommandRepository.TrashedTransfer(transfer_id)

	if err != nil {
		return s.errorhandler.HandleTrashedTransferError(err, method, "FAILED_TRASHED_TRANSFER", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransferResponse(res)

	s.mencache.DeleteTransferCache(transfer_id)

	logSuccess("Successfully trashed transfer",
		zap.Int("transfer.id", transfer_id),
	)

	return so, nil
}

func (s *transferCommandService) RestoreTransfer(transfer_id int) (*response.TransferResponse, *response.ErrorResponse) {
	const method = "RestoreTransfer"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	res, err := s.transferCommandRepository.RestoreTransfer(transfer_id)

	if err != nil {
		return s.errorhandler.HandleRestoreTransferError(err, method, "FAILED_RESTORE_TRANSFER", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransferResponse(res)

	logSuccess("RestoreTransfer process completed", zap.Int("transfer.id", transfer_id))

	return so, nil
}

func (s *transferCommandService) DeleteTransferPermanent(transfer_id int) (bool, *response.ErrorResponse) {
	const method = "DeleteTransferPermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	_, err := s.transferCommandRepository.DeleteTransferPermanent(transfer_id)

	if err != nil {
		return s.errorhandler.HandleDeleteTransferPermanentError(err, method, "FAILED_DELETE_TRANSFER_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("DeleteTransferPermanent process completed", zap.Int("transfer.id", transfer_id))

	return true, nil
}

func (s *transferCommandService) RestoreAllTransfer() (bool, *response.ErrorResponse) {
	const method = "RestoreAllTransfer"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	_, err := s.transferCommandRepository.RestoreAllTransfer()

	if err != nil {
		return s.errorhandler.HandleRestoreAllTransferError(err, method, "FAILED_RESTORE_ALL_TRANSFERS", span, &status, zap.Error(err))
	}

	logSuccess("RestoreAllTransfer process completed", zap.Bool("success", true))

	return true, nil
}

func (s *transferCommandService) DeleteAllTransferPermanent() (bool, *response.ErrorResponse) {
	const method = "DeleteAllTransferPermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	_, err := s.transferCommandRepository.DeleteAllTransferPermanent()

	if err != nil {
		return s.errorhandler.HandleDeleteAllTransferPermanentError(err, method, "FAILED_DELETE_ALL_TRANSFER_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("DeleteAllTransferPermanent process completed", zap.Bool("success", true))

	return true, nil
}

func (s *transferCommandService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
	trace.Span,
	func(string),
	string,
	func(string, ...zap.Field),
) {
	start := time.Now()
	status := "success"

	_, span := s.trace.Start(s.ctx, method)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.AddEvent("Start: " + method)

	s.logger.Info("Start: " + method)

	end := func(status string) {
		s.recordMetrics(method, status, start)
		code := codes.Ok
		if status != "success" {
			code = codes.Error
		}
		span.SetStatus(code, status)
		span.End()
	}

	logSuccess := func(msg string, fields ...zap.Field) {
		span.AddEvent(msg)
		s.logger.Info(msg, fields...)
	}

	return span, end, status, logSuccess
}

func (s *transferCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
