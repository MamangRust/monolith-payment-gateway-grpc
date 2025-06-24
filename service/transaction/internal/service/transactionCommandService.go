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
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-transaction/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transactionCommandService struct {
	kafka                        *kafka.Kafka
	ctx                          context.Context
	errorhandler                 errorhandler.TransactionCommandErrorHandler
	mencache                     mencache.TransactionCommandCache
	trace                        trace.Tracer
	merchantRepository           repository.MerchantRepository
	cardRepository               repository.CardRepository
	saldoRepository              repository.SaldoRepository
	transactionQueryRepository   repository.TransactionQueryRepository
	transactionCommandRepository repository.TransactionCommandRepository
	logger                       logger.LoggerInterface
	mapping                      responseservice.TransactionResponseMapper
	requestCounter               *prometheus.CounterVec
	requestDuration              *prometheus.HistogramVec
}

func NewTransactionCommandService(
	kafka *kafka.Kafka,
	ctx context.Context,
	errorhandler errorhandler.TransactionCommandErrorHandler,
	mencache mencache.TransactionCommandCache,
	merchantRepository repository.MerchantRepository,
	cardRepository repository.CardRepository,
	saldoRepository repository.SaldoRepository,
	transactionCommandRepository repository.TransactionCommandRepository,
	transactionQueryRepository repository.TransactionQueryRepository,
	logger logger.LoggerInterface,
	mapping responseservice.TransactionResponseMapper,
) *transactionCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transaction_command_service_request_total",
			Help: "Total number of requests to the TransactionCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transaction_command_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TransactionCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &transactionCommandService{
		kafka:                        kafka,
		ctx:                          ctx,
		errorhandler:                 errorhandler,
		mencache:                     mencache,
		trace:                        otel.Tracer("transaction-command-service"),
		merchantRepository:           merchantRepository,
		cardRepository:               cardRepository,
		saldoRepository:              saldoRepository,
		transactionCommandRepository: transactionCommandRepository,
		transactionQueryRepository:   transactionQueryRepository,
		logger:                       logger,
		mapping:                      mapping,
		requestCounter:               requestCounter,
		requestDuration:              requestDuration,
	}
}

func (s *transactionCommandService) Create(apiKey string, request *requests.CreateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse) {
	const method = "CreateTransacton"

	s.logger.Debug("CreateTransaction called", zap.String("card_number", request.CardNumber), zap.String("api_key", apiKey))

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.String("apikey", apiKey))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantRepository.FindByApiKey(apiKey)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_MERCHANT_BY_API_KEY", span, &status, merchant_errors.ErrFailedFindByApiKey, zap.Error(err))
	}

	card, err := s.cardRepository.FindUserCardByCardNumber(request.CardNumber)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_CARD_BY_CARD_NUMBER", span, &status, card_errors.ErrFailedFindByCardNumber, zap.Error(err))
	}

	saldo, err := s.saldoRepository.FindByCardNumber(card.CardNumber)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_SALDO_BY_CARD_NUMBER", span, &status, saldo_errors.ErrFailedFindSaldoByCardNumber, zap.Error(err))
	}

	if saldo.TotalBalance < request.Amount {
		return s.errorhandler.HandleInsufficientBalanceError(err, method, "FAILED_INSUFFICIENT_BALANCE", span, &status, card.CardNumber, zap.Error(err))
	}

	saldo.TotalBalance -= request.Amount
	if _, err := s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   card.CardNumber,
		TotalBalance: saldo.TotalBalance,
	}); err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_SALDO", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
	}

	request.MerchantID = &merchant.ID

	transaction, err := s.transactionCommandRepository.CreateTransaction(request)
	if err != nil {
		saldo.TotalBalance += request.Amount
		_, err := s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
			CardNumber:   card.CardNumber,
			TotalBalance: saldo.TotalBalance,
		})
		if err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_SALDO", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
		}

		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: transaction.ID,
			Status:        "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, transaction_errors.ErrFailedUpdateTransaction, zap.Error(err))
		}

		return s.errorhandler.HandleCreateTransactionError(err, method, "FAILED_CREATE_TRANSACTION", span, &status, zap.Error(err))
	}

	if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
		TransactionID: transaction.ID,
		Status:        "success",
	}); err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, transaction_errors.ErrFailedUpdateTransaction, zap.Error(err))
	}

	merchantCard, err := s.cardRepository.FindCardByUserId(merchant.UserID)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_CARD_BY_USER_ID_MERCHANT", span, &status, card_errors.ErrFailedFindByCardNumber, zap.Error(err))
	}

	merchantSaldo, err := s.saldoRepository.FindByCardNumber(merchantCard.CardNumber)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_SALDO_BY_CARD_NUMBER_MERCHANT", span, &status, merchant_errors.ErrFailedFindByApiKey, zap.Error(err))
	}

	merchantSaldo.TotalBalance += request.Amount

	if _, err := s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   merchantCard.CardNumber,
		TotalBalance: merchantSaldo.TotalBalance,
	}); err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_SALDO_MERCHANT", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
	}

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   "Transaction Successful",
		"Message": fmt.Sprintf("Your transaction of %d has been processed successfully.", request.Amount),
		"Button":  "View History",
		"Link":    "https://sanedge.example.com/transaction/history",
	})

	emailPayload := map[string]any{
		"email":   card.Email,
		"subject": "Transaction Successful - SanEdge",
		"body":    htmlBody,
	}

	payloadBytes, err := json.Marshal(emailPayload)
	if err != nil {
		return errorhandler.HandleErrorMarshal[*response.TransactionResponse](s.logger, err, method, "FAILED_MARSHAL_EMAIL_PAYLOAD", span, &status, transaction_errors.ErrFailedCreateTransaction, zap.Error(err))
	}

	err = s.kafka.SendMessage("email-service-topic-transaction-create", strconv.Itoa(transaction.ID), payloadBytes)
	if err != nil {
		return errorhandler.HandleErrorKafkaSend[*response.TransactionResponse](s.logger, err, method, "FAILED_SEND_EMAIL", span, &status, transaction_errors.ErrFailedCreateTransaction, zap.Error(err))
	}

	so := s.mapping.ToTransactionResponse(transaction)

	logSuccess("Successfully created transaction", zap.Int("transaction.id", transaction.ID))

	return so, nil
}

func (s *transactionCommandService) Update(apiKey string, request *requests.UpdateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse) {
	const method = "UpdateTransaction"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	transaction, err := s.transactionQueryRepository.FindById(*request.TransactionID)
	if err != nil {
		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, transaction_errors.ErrFailedUpdateTransaction, zap.Error(err))
		}

		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_TRANSACTION_BY_ID", span, &status, transaction_errors.ErrFailedUpdateTransaction, zap.Error(err))
	}

	merchant, err := s.merchantRepository.FindByApiKey(apiKey)
	if err != nil || transaction.MerchantID != merchant.ID {
		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, merchant_errors.ErrFailedFindByApiKey, zap.Error(err))
		}

		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_MERCHANT_BY_API_KEY", span, &status, merchant_errors.ErrFailedFindByApiKey, zap.Error(err))
	}

	card, err := s.cardRepository.FindCardByCardNumber(transaction.CardNumber)
	if err != nil {
		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			return s.errorhandler.HandleUpdateTransactionError(err, method, "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, zap.Error(err))
		}

		return nil, card_errors.ErrCardNotFoundRes
	}

	saldo, err := s.saldoRepository.FindByCardNumber(card.CardNumber)
	if err != nil {
		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			return s.errorhandler.HandleUpdateTransactionError(err, method, "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, zap.Error(err))
		}

		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_SALDO_BY_CARD_NUMBER", span, &status, saldo_errors.ErrFailedSaldoNotFound, zap.Error(err))
	}

	saldo.TotalBalance += transaction.Amount
	s.logger.Info("Restoring balance for old transaction amount", zap.Int("RestoredBalance", saldo.TotalBalance))

	if _, err := s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   card.CardNumber,
		TotalBalance: saldo.TotalBalance,
	}); err != nil {
		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			return s.errorhandler.HandleUpdateTransactionError(err, method, "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, zap.Error(err))
		}

		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
	}

	if saldo.TotalBalance < request.Amount {
		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			return s.errorhandler.HandleUpdateTransactionError(err, method, "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, zap.Error(err))

		}
		return s.errorhandler.HandleInsufficientBalanceError(err, method, "INSUFFICIENT_BALANCE", span, &status, card.CardNumber, zap.Error(err))
	}

	saldo.TotalBalance -= request.Amount
	s.logger.Info("Updating balance for updated transaction amount")

	if _, err := s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   card.CardNumber,
		TotalBalance: saldo.TotalBalance,
	}); err != nil {
		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			return s.errorhandler.HandleUpdateTransactionError(err, method, "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, zap.Error(err))
		}

		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
	}

	transaction.Amount = request.Amount
	transaction.PaymentMethod = request.PaymentMethod

	layout := "2006-01-02 15:04:05"
	parsedTime, err := time.Parse(layout, transaction.TransactionTime)
	if err != nil {
		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			return s.errorhandler.HandleUpdateTransactionError(err, method, "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, zap.Error(err))
		}

		return s.errorhandler.HandleInvalidParseTimeError(err, method, "INVALID_TRANSACTION_TIME", span, &status, transaction.TransactionTime, zap.Error(err))
	}

	res, err := s.transactionCommandRepository.UpdateTransaction(&requests.UpdateTransactionRequest{
		TransactionID:   &transaction.ID,
		CardNumber:      transaction.CardNumber,
		Amount:          transaction.Amount,
		PaymentMethod:   transaction.PaymentMethod,
		MerchantID:      &transaction.MerchantID,
		TransactionTime: parsedTime,
	})
	if err != nil {
		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			return s.errorhandler.HandleUpdateTransactionError(err, method, "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, zap.Error(err))
		}

		return s.errorhandler.HandleUpdateTransactionError(err, method, "FAILED_UPDATE_TRANSACTION", span, &status, zap.Error(err))
	}

	if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
		TransactionID: transaction.ID,
		Status:        "success",
	}); err != nil {
		return s.errorhandler.HandleUpdateTransactionError(err, method, "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionResponse(res)

	s.mencache.DeleteTransactionCache(*request.TransactionID)

	logSuccess("Successfully updated transaction", zap.Bool("success", true))

	return so, nil
}

func (s *transactionCommandService) TrashedTransaction(transaction_id int) (*response.TransactionResponse, *response.ErrorResponse) {
	const method = "TrashedTransaction"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	res, err := s.transactionCommandRepository.TrashedTransaction(transaction_id)

	if err != nil {
		return s.errorhandler.HandleTrashedTransactionError(err, method, "FAILED_TRASHED_TRANSACTION", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionResponse(res)

	s.mencache.DeleteTransactionCache(transaction_id)

	logSuccess("Successfully trashed transaction", zap.Bool("success", true))

	return so, nil
}

func (s *transactionCommandService) RestoreTransaction(transaction_id int) (*response.TransactionResponse, *response.ErrorResponse) {
	const method = "RestoreTransaction"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	res, err := s.transactionCommandRepository.RestoreTransaction(transaction_id)

	if err != nil {
		return s.errorhandler.HandleRestoreTransactionError(err, method, "FAILED_RESTORE_TRANSACTION", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionResponse(res)

	logSuccess("Successfully restored transaction", zap.Bool("success", true))

	return so, nil
}

func (s *transactionCommandService) DeleteTransactionPermanent(transaction_id int) (bool, *response.ErrorResponse) {
	const method = "DeleteTransactionPermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	_, err := s.transactionCommandRepository.DeleteTransactionPermanent(transaction_id)

	if err != nil {
		return s.errorhandler.HandleDeleteTransactionPermanentError(err, method, "FAILED_DELETE_TRANSACTION_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted transaction permanent", zap.Bool("success", true))

	return true, nil
}

func (s *transactionCommandService) RestoreAllTransaction() (bool, *response.ErrorResponse) {
	const method = "RestoreAllTransaction"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	_, err := s.transactionCommandRepository.RestoreAllTransaction()
	if err != nil {
		return s.errorhandler.HandleRestoreAllTransactionError(err, method, "FAILED_RESTORE_ALL_TRANSACTIONS", span, &status, zap.Error(err))
	}

	logSuccess("Successfully restored all transactions", zap.Bool("success", true))

	return true, nil
}

func (s *transactionCommandService) DeleteAllTransactionPermanent() (bool, *response.ErrorResponse) {
	const method = "DeleteAllTransactionPermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	_, err := s.transactionCommandRepository.DeleteAllTransactionPermanent()

	if err != nil {
		return s.errorhandler.HandleDeleteAllTransactionPermanentError(err, method, "FAILED_DELETE_ALL_TRANSACTION_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted all transactions permanent", zap.Bool("success", true))

	return true, nil
}

func (s *transactionCommandService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *transactionCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
