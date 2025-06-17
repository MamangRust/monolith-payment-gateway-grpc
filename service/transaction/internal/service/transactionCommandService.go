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
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("CreateTransaction", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "CreateTransaction")
	defer span.End()

	span.SetAttributes(
		attribute.String("apiKey", apiKey),
	)

	s.logger.Debug("Starting CreateTransaction process",
		zap.String("apiKey", apiKey),
		zap.Any("request", request),
	)

	merchant, err := s.merchantRepository.FindByApiKey(apiKey)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "FindByApiKey", "FAILED_FIND_MERCHANT_BY_API_KEY", span, &status, merchant_errors.ErrFailedFindByApiKey, zap.String("api_key", apiKey))
	}

	card, err := s.cardRepository.FindUserCardByCardNumber(request.CardNumber)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "FindUserCardByCardNumber", "FAILED_FIND_CARD_BY_CARD_NUMBER", span, &status, merchant_errors.ErrFailedFindByApiKey, zap.String("card_number", request.CardNumber))
	}

	saldo, err := s.saldoRepository.FindByCardNumber(card.CardNumber)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "FindByCardNumber", "FAILED_FIND_SALDO_BY_CARD_NUMBER", span, &status, merchant_errors.ErrFailedFindByApiKey, zap.String("card_number", card.CardNumber))
	}

	if saldo.TotalBalance < request.Amount {
		return s.errorhandler.HandleInsufficientBalanceError(err, "CreateTransaction", "FAILED_INSUFFICIENT_BALANCE", span, &status, card.CardNumber, zap.String("card_number", card.CardNumber))
	}

	saldo.TotalBalance -= request.Amount
	if _, err := s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   card.CardNumber,
		TotalBalance: saldo.TotalBalance,
	}); err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "UpdateSaldoBalance", "FAILED_UPDATE_SALDO", span, &status, merchant_errors.ErrFailedFindByApiKey, zap.String("card_number", card.CardNumber))
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
			return s.errorhandler.HandleRepositorySingleError(err, "UpdateSaldoBalance", "FAILED_UPDATE_SALDO", span, &status, merchant_errors.ErrFailedFindByApiKey, zap.String("card_number", card.CardNumber))
		}

		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: transaction.ID,
			Status:        "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, "UpdateTransactionStatus", "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, merchant_errors.ErrFailedFindByApiKey, zap.String("card_number", card.CardNumber))
		}

		return s.errorhandler.HandleCreateTransactionError(err, "CreateTransaction", "FAILED_CREATE_TRANSACTION", span, &status, zap.String("card_number", card.CardNumber))
	}

	if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
		TransactionID: transaction.ID,
		Status:        "success",
	}); err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "UpdateTransactionStatus", "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, merchant_errors.ErrFailedFindByApiKey, zap.String("card_number", card.CardNumber))
	}

	merchantCard, err := s.cardRepository.FindCardByUserId(merchant.UserID)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "FindCardByUserId", "FAILED_FIND_CARD_BY_USER_ID", span, &status, merchant_errors.ErrFailedFindByApiKey, zap.Int("user_id", merchant.UserID))
	}

	merchantSaldo, err := s.saldoRepository.FindByCardNumber(merchantCard.CardNumber)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "FindByCardNumber", "FAILED_FIND_SALDO_BY_CARD_NUMBER", span, &status, merchant_errors.ErrFailedFindByApiKey, zap.String("card_number", merchantCard.CardNumber))
	}

	merchantSaldo.TotalBalance += request.Amount

	if _, err := s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   merchantCard.CardNumber,
		TotalBalance: merchantSaldo.TotalBalance,
	}); err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "UpdateSaldoBalance", "FAILED_UPDATE_SALDO", span, &status, merchant_errors.ErrFailedFindByApiKey, zap.String("card_number", merchantCard.CardNumber))
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
		return errorhandler.HandleErrorMarshal[*response.TransactionResponse](s.logger, err, "CreateTransaction", "FAILED_MARSHAL_EMAIL_PAYLOAD", span, &status, transaction_errors.ErrFailedCreateTransaction, zap.String("card_number", card.CardNumber))
	}

	err = s.kafka.SendMessage("email-service-topic-transaction-create", strconv.Itoa(transaction.ID), payloadBytes)
	if err != nil {
		return errorhandler.HandleErrorKafkaSend[*response.TransactionResponse](s.logger, err, "CreateTransaction", "FAILED_SEND_EMAIL", span, &status, transaction_errors.ErrFailedCreateTransaction, zap.String("card_number", card.CardNumber))
	}

	so := s.mapping.ToTransactionResponse(transaction)

	s.logger.Debug("CreateTransaction process completed",
		zap.String("apiKey", apiKey),
		zap.Int("transactionID", transaction.ID),
	)

	return so, nil
}

func (s *transactionCommandService) Update(apiKey string, request *requests.UpdateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("UpdateTransaction", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateTransaction")
	defer span.End()

	span.SetAttributes(
		attribute.String("apiKey", apiKey),
		attribute.Int("transaction_id", *request.TransactionID),
	)

	s.logger.Debug("Starting UpdateTransaction process",
		zap.String("apiKey", apiKey),
		zap.Int("transaction_id", *request.TransactionID),
	)

	transaction, err := s.transactionQueryRepository.FindById(*request.TransactionID)
	if err != nil {
		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, "UpdateTransactionStatus", "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, transaction_errors.ErrFailedUpdateTransaction, zap.Int("transaction_id", *request.TransactionID))
		}

		return s.errorhandler.HandleRepositorySingleError(err, "FindById", "FAILED_FIND_TRANSACTION_BY_ID", span, &status, transaction_errors.ErrFailedUpdateTransaction, zap.Int("transaction_id", *request.TransactionID))
	}

	merchant, err := s.merchantRepository.FindByApiKey(apiKey)
	if err != nil || transaction.MerchantID != merchant.ID {
		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, "UpdateTransactionStatus", "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, merchant_errors.ErrFailedFindByApiKey, zap.String("api_key", apiKey))
		}

		return s.errorhandler.HandleRepositorySingleError(err, "FindByApiKey", "FAILED_FIND_MERCHANT_BY_API_KEY", span, &status, merchant_errors.ErrFailedFindByApiKey, zap.String("api_key", apiKey))
	}

	card, err := s.cardRepository.FindCardByCardNumber(transaction.CardNumber)
	if err != nil {
		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			return s.errorhandler.HandleUpdateTransactionError(err, "UpdateTransactionStatus", "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, zap.String("card_number", transaction.CardNumber))
		}

		return nil, card_errors.ErrCardNotFoundRes
	}

	saldo, err := s.saldoRepository.FindByCardNumber(card.CardNumber)
	if err != nil {
		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			return s.errorhandler.HandleUpdateTransactionError(err, "UpdateTransactionStatus", "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, zap.String("card_number", card.CardNumber))
		}

		return s.errorhandler.HandleRepositorySingleError(err, "FindByCardNumber", "FAILED_FIND_SALDO_BY_CARD_NUMBER", span, &status, saldo_errors.ErrFailedSaldoNotFound, zap.String("card_number", card.CardNumber))
	}

	saldo.TotalBalance += transaction.Amount
	s.logger.Debug("Restoring balance for old transaction amount", zap.Int("RestoredBalance", saldo.TotalBalance))

	if _, err := s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   card.CardNumber,
		TotalBalance: saldo.TotalBalance,
	}); err != nil {
		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			return s.errorhandler.HandleUpdateTransactionError(err, "UpdateTransactionStatus", "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, zap.String("card_number", card.CardNumber))
		}

		return s.errorhandler.HandleRepositorySingleError(err, "UpdateSaldoBalance", "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.String("card_number", card.CardNumber))
	}

	if saldo.TotalBalance < request.Amount {
		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			return s.errorhandler.HandleUpdateTransactionError(err, "UpdateTransactionStatus", "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, zap.String("card_number", card.CardNumber))

		}
		return s.errorhandler.HandleInsufficientBalanceError(err, "UpdateTransaction", "INSUFFICIENT_BALANCE", span, &status, card.CardNumber, zap.String("card_number", card.CardNumber))
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
			return s.errorhandler.HandleUpdateTransactionError(err, "UpdateTransactionStatus", "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, zap.String("card_number", card.CardNumber))
		}

		return s.errorhandler.HandleRepositorySingleError(err, "UpdateSaldoBalance", "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.String("card_number", card.CardNumber))
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
			return s.errorhandler.HandleUpdateTransactionError(err, "UpdateTransactionStatus", "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, zap.String("card_number", card.CardNumber))
		}

		return s.errorhandler.HandleInvalidParseTimeError(err, "UpdateTransaction", "INVALID_TRANSACTION_TIME", span, &status, transaction.TransactionTime, zap.String("card_number", card.CardNumber))
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
			return s.errorhandler.HandleUpdateTransactionError(err, "UpdateTransactionStatus", "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, zap.String("card_number", card.CardNumber))
		}

		return s.errorhandler.HandleUpdateTransactionError(err, "UpdateTransaction", "FAILED_UPDATE_TRANSACTION", span, &status, zap.String("card_number", card.CardNumber))
	}

	if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
		TransactionID: transaction.ID,
		Status:        "success",
	}); err != nil {
		return s.errorhandler.HandleUpdateTransactionError(err, "UpdateTransactionStatus", "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, zap.String("card_number", card.CardNumber))
	}

	so := s.mapping.ToTransactionResponse(res)

	s.mencache.DeleteTransactionCache(*request.TransactionID)

	s.logger.Debug("UpdateTransaction process completed",
		zap.String("apiKey", apiKey),
		zap.Int("transaction_id", *request.TransactionID),
	)

	return so, nil
}

func (s *transactionCommandService) TrashedTransaction(transaction_id int) (*response.TransactionResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("TrashedTransaction", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "TrashedTransaction")
	defer span.End()

	span.SetAttributes(
		attribute.Int("transaction_id", transaction_id),
	)

	s.logger.Debug("Starting TrashedTransaction process",
		zap.Int("transaction_id", transaction_id),
	)

	res, err := s.transactionCommandRepository.TrashedTransaction(transaction_id)

	if err != nil {
		return s.errorhandler.HandleTrashedTransactionError(err, "TrashedTransaction", "FAILED_TRASHED_TRANSACTION", span, &status, zap.Int("transaction_id", transaction_id))
	}

	so := s.mapping.ToTransactionResponse(res)

	s.mencache.DeleteTransactionCache(transaction_id)

	s.logger.Debug("Successfully trashed transaction", zap.Int("transaction_id", transaction_id))

	return so, nil
}

func (s *transactionCommandService) RestoreTransaction(transaction_id int) (*response.TransactionResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreTransaction", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreTransaction")
	defer span.End()

	span.SetAttributes(
		attribute.Int("transaction_id", transaction_id),
	)

	s.logger.Debug("Starting RestoreTransaction process",
		zap.Int("transaction_id", transaction_id),
	)

	res, err := s.transactionCommandRepository.RestoreTransaction(transaction_id)

	if err != nil {
		return s.errorhandler.HandleRestoreTransactionError(err, "RestoreTransaction", "FAILED_RESTORE_TRANSACTION", span, &status, zap.Int("transaction_id", transaction_id))
	}

	so := s.mapping.ToTransactionResponse(res)

	s.logger.Debug("Successfully restored transaction", zap.Int("transaction_id", transaction_id))

	return so, nil
}

func (s *transactionCommandService) DeleteTransactionPermanent(transaction_id int) (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteTransactionPermanent", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteTransactionPermanent")
	defer span.End()

	span.SetAttributes(
		attribute.Int("transaction_id", transaction_id),
	)

	_, err := s.transactionCommandRepository.DeleteTransactionPermanent(transaction_id)

	if err != nil {
		return s.errorhandler.HandleDeleteTransactionPermanentError(err, "DeleteTransactionPermanent", "FAILED_DELETE_TRANSACTION_PERMANENT", span, &status, zap.Int("transaction_id", transaction_id))
	}

	s.logger.Debug("Successfully permanently deleted transaction", zap.Int("transaction_id", transaction_id))

	return true, nil
}

func (s *transactionCommandService) RestoreAllTransaction() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreAllTransaction", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreAllTransaction")
	defer span.End()

	s.logger.Debug("Restoring all transactions")

	_, err := s.transactionCommandRepository.RestoreAllTransaction()
	if err != nil {
		return s.errorhandler.HandleRestoreAllTransactionError(err, "RestoreAllTransaction", "FAILED_RESTORE_ALL_TRANSACTIONS", span, &status)
	}

	s.logger.Debug("Successfully restored all transactions")
	return true, nil
}

func (s *transactionCommandService) DeleteAllTransactionPermanent() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteAllTransactionPermanent", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteAllTransactionPermanent")
	defer span.End()

	_, err := s.transactionCommandRepository.DeleteAllTransactionPermanent()

	if err != nil {
		return s.errorhandler.HandleDeleteAllTransactionPermanentError(err, "DeleteAllTransactionPermanent", "FAILED_DELETE_ALL_TRANSACTION_PERMANENT", span, &status)
	}

	s.logger.Debug("Successfully deleted all transactions permanently")

	return true, nil
}

func (s *transactionCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
