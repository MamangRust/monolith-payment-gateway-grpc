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
	traceunic "github.com/MamangRust/monolith-payment-gateway-pkg/trace_unic"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transactionCommandService struct {
	kafka                        kafka.Kafka
	ctx                          context.Context
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
	kafka kafka.Kafka,
	ctx context.Context,
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
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &transactionCommandService{
		kafka:                        kafka,
		ctx:                          ctx,
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
		traceID := traceunic.GenerateTraceID("FAILED_FIND_BY_API_KEY")

		s.logger.Error("Failed to find merchant", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find merchant")
		status = "failed_to_find_merchant"

		return nil, merchant_errors.ErrFailedFindByApiKey
	}

	card, err := s.cardRepository.FindUserCardByCardNumber(request.CardNumber)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_CARD_BY_CARD_NUMBER")

		s.logger.Error("Card not found for card number",
			zap.String("card_number", request.CardNumber),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Card not found for card number")
		status = "card_not_found_for_card_number"

		return nil, card_errors.ErrFailedFindByCardNumber
	}

	saldo, err := s.saldoRepository.FindByCardNumber(card.CardNumber)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_SALDO_BY_CARD_NUMBER")

		s.logger.Error("Failed to retrieve saldo details",
			zap.String("card_number", card.CardNumber),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve saldo details")
		status = "failed_to_retrieve_saldo_details"

		return nil, saldo_errors.ErrFailedSaldoNotFound
	}

	if saldo.TotalBalance < request.Amount {
		traceID := traceunic.GenerateTraceID("INSUFFICIENT_BALANCE")

		s.logger.Error("Insufficient balance",
			zap.String("card_number", card.CardNumber),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Insufficient balance")
		status = "insufficient_balance"

		return nil, saldo_errors.ErrFailedInsuffientBalance
	}

	saldo.TotalBalance -= request.Amount
	if _, err := s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   card.CardNumber,
		TotalBalance: saldo.TotalBalance,
	}); err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_SALDO")

		s.logger.Error("failed to update saldo", zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update saldo")
		status = "failed_to_update_saldo"

		return nil, saldo_errors.ErrFailedUpdateSaldo
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
			traceID := traceunic.GenerateTraceID("FAILED_UPDATE_SALDO")

			s.logger.Error("failed to update saldo", zap.Error(err))
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to update saldo")
			status = "failed_to_update_saldo"

			return nil, saldo_errors.ErrFailedUpdateSaldo
		}

		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: transaction.ID,
			Status:        "failed",
		}); err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_UPDATE_TRANSACTION_STATUS")

			s.logger.Error("failed to update transaction status", zap.Error(err))
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to update transaction status")
			status = "failed_to_update_transaction_status"
		}

		return nil, transaction_errors.ErrFailedCreateTransaction
	}

	if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
		TransactionID: transaction.ID,
		Status:        "success",
	}); err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_TRANSACTION_STATUS")

		s.logger.Error("failed to update transaction status", zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update transaction status")
		status = "failed_to_update_transaction_status"

		return nil, transaction_errors.ErrFailedUpdateTransaction
	}

	merchantCard, err := s.cardRepository.FindCardByUserId(merchant.UserID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_CARD_BY_USER_ID")

		s.logger.Error("failed to find merchant card", zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find merchant card")
		status = "failed_to_find_merchant_card"

		return nil, card_errors.ErrCardNotFoundRes
	}

	merchantSaldo, err := s.saldoRepository.FindByCardNumber(merchantCard.CardNumber)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_SALDO_BY_CARD_NUMBER")

		s.logger.Error("failed to find merchant saldo", zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find merchant saldo")
		status = "failed_to_find_merchant_saldo"

		return nil, saldo_errors.ErrFailedSaldoNotFound
	}

	merchantSaldo.TotalBalance += request.Amount

	if _, err := s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   merchantCard.CardNumber,
		TotalBalance: merchantSaldo.TotalBalance,
	}); err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_SALDO")

		s.logger.Error("failed to update merchant saldo", zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update merchant saldo")
		status = "failed_to_update_merchant_saldo"

		return nil, saldo_errors.ErrFailedUpdateSaldo
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
		traceID := traceunic.GenerateTraceID("TransactionErr")
		s.logger.Error("Failed to marshal transaction email payload", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal transaction email payload")
		return nil, withdraw_errors.ErrFailedSendEmail
	}

	err = s.kafka.SendMessage("email-service-topic-transaction-create", strconv.Itoa(transaction.ID), payloadBytes)
	if err != nil {
		traceID := traceunic.GenerateTraceID("TransactionErr")
		s.logger.Error("Failed to send transaction email message", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to send transaction email")
		return nil, withdraw_errors.ErrFailedSendEmail
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
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TRANSACTION_BY_ID")

		s.logger.Error("failed to find transaction by id", zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find transaction by id")
		status = "failed_to_find_transaction_by_id"

		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_UPDATE_TRANSACTION_STATUS")

			s.logger.Error("failed to update transaction status", zap.Error(err))
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to update transaction status")
			status = "failed_to_update_transaction_status"
		}

		return nil, transaction_errors.ErrFailedUpdateTransaction
	}

	merchant, err := s.merchantRepository.FindByApiKey(apiKey)
	if err != nil || transaction.MerchantID != merchant.ID {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MERCHANT_BY_API_KEY")

		s.logger.Error("failed to find merchant by api key", zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find merchant by api key")
		status = "failed_to_find_merchant_by_api_key"

		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_UPDATE_TRANSACTION_STATUS")

			s.logger.Error("failed to update transaction status", zap.String("trace.id", traceID), zap.Error(err))
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to update transaction status")
			status = "failed_to_update_transaction_status"
		}

		return nil, transaction_errors.ErrFailedUpdateTransaction
	}

	card, err := s.cardRepository.FindCardByCardNumber(transaction.CardNumber)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_CARD_BY_CARD_NUMBER")

		s.logger.Error("failed to find card by card number", zap.String("trace.id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find card by card number")
		status = "failed_to_find_card_by_card_number"

		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_UPDATE_TRANSACTION_STATUS")

			s.logger.Error("failed to update transaction status", zap.String("trace.id", traceID), zap.Error(err))
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to update transaction status")
			status = "failed_to_update_transaction_status"
		}

		return nil, card_errors.ErrCardNotFoundRes
	}

	saldo, err := s.saldoRepository.FindByCardNumber(card.CardNumber)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_SALDO_BY_CARD_NUMBER")

		s.logger.Error("failed to find saldo by card number", zap.String("trace.id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find saldo by card number")
		status = "failed_to_find_saldo_by_card_number"

		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_UPDATE_TRANSACTION_STATUS")

			s.logger.Error("failed to update transaction status", zap.String("trace.id", traceID), zap.Error(err))
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to update transaction status")
			status = "failed_to_update_transaction_status"
		}

		return nil, saldo_errors.ErrFailedSaldoNotFound
	}

	saldo.TotalBalance += transaction.Amount
	s.logger.Debug("Restoring balance for old transaction amount", zap.Int("RestoredBalance", saldo.TotalBalance))

	if _, err := s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   card.CardNumber,
		TotalBalance: saldo.TotalBalance,
	}); err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_SALDO_BALANCE")

		s.logger.Error("failed to update saldo balance", zap.String("trace.id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update saldo balance")
		status = "failed_to_update_saldo_balance"

		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_UPDATE_TRANSACTION_STATUS")

			s.logger.Error("failed to update transaction status", zap.String("trace.id", traceID), zap.Error(err))
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to update transaction status")
			status = "failed_to_update_transaction_status"
		}

		return nil, saldo_errors.ErrFailedUpdateSaldo
	}

	if saldo.TotalBalance < request.Amount {
		traceID := traceunic.GenerateTraceID("FAILED_INSUFFICIENT_FUNDS")

		s.logger.Error("failed to update transaction status", zap.String("trace.id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update transaction status")
		status = "failed_to_update_transaction_status"

		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_UPDATE_TRANSACTION_STATUS")

			s.logger.Error("failed to update transaction status", zap.String("trace.id", traceID), zap.Error(err))
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to update transaction status")
			status = "failed_to_update_transaction_status"
		}

		return nil, transaction_errors.ErrFailedUpdateTransaction
	}

	saldo.TotalBalance -= request.Amount
	s.logger.Info("Updating balance for updated transaction amount")

	if _, err := s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   card.CardNumber,
		TotalBalance: saldo.TotalBalance,
	}); err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_SALDO_BALANCE")

		s.logger.Error("failed to update saldo balance", zap.String("trace.id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update saldo balance")
		status = "failed_to_update_saldo_balance"

		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_UPDATE_TRANSACTION_STATUS")

			s.logger.Error("failed to update transaction status", zap.String("trace.id", traceID), zap.Error(err))
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to update transaction status")
			status = "failed_to_update_transaction_status"
		}

		return nil, saldo_errors.ErrFailedUpdateSaldo
	}

	transaction.Amount = request.Amount
	transaction.PaymentMethod = request.PaymentMethod

	layout := "2006-01-02 15:04:05"
	parsedTime, err := time.Parse(layout, transaction.TransactionTime)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_PARSE_TRANSACTION_TIME")

		s.logger.Error("failed to parse transaction time", zap.String("trace.id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to parse transaction time")
		status = "failed_to_parse_transaction_time"

		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_UPDATE_TRANSACTION_STATUS")

			s.logger.Error("failed to update transaction status", zap.String("trace.id", traceID), zap.Error(err))
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to update transaction status")
			status = "failed_to_update_transaction_status"
		}

		return nil, transaction_errors.ErrFailedUpdateTransaction
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
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_TRANSACTION")

		s.logger.Error("failed to update transaction", zap.String("trace.id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update transaction")
		status = "failed_to_update_transaction"

		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_UPDATE_TRANSACTION_STATUS")

			s.logger.Error("failed to update transaction status", zap.String("trace.id", traceID), zap.Error(err))
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to update transaction status")
			status = "failed_to_update_transaction_status"
		}

		return nil, transaction_errors.ErrFailedUpdateTransaction
	}

	if _, err := s.transactionCommandRepository.UpdateTransactionStatus(&requests.UpdateTransactionStatus{
		TransactionID: transaction.ID,
		Status:        "success",
	}); err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_TRANSACTION_STATUS")

		s.logger.Error("failed to update transaction status", zap.String("trace.id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update transaction status")
		status = "failed_to_update_transaction_status"

		return nil, transaction_errors.ErrFailedUpdateTransaction
	}

	so := s.mapping.ToTransactionResponse(res)

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
		traceID := traceunic.GenerateTraceID("FAILED_TRASHED_TRANSACTION")

		s.logger.Error("failed to trashed transaction", zap.String("trace.id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to trashed transaction")
		status = "failed_to_trashed_transaction"

		return nil, transaction_errors.ErrFailedTrashedTransaction
	}

	so := s.mapping.ToTransactionResponse(res)

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
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_TRANSACTION")

		s.logger.Error("failed to restore transaction", zap.String("trace.id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore transaction")
		status = "failed_to_restore_transaction"

		return nil, transaction_errors.ErrFailedRestoreTransaction
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
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_TRANSACTION_PERMANENT")

		s.logger.Error("failed to permanently delete transaction", zap.String("trace.id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete transaction")
		status = "failed_to_permanently_delete_transaction"

		return false, transaction_errors.ErrFailedDeleteTransactionPermanent
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
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ALL_TRANSACTIONS")

		s.logger.Error("failed to restore all transactions", zap.String("trace.id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore all transactions")
		status = "failed_to_restore_all_transactions"

		return false, transaction_errors.ErrFailedRestoreAllTransactions
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
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_ALL_TRANSACTIONS_PERMANENT")

		s.logger.Error("failed to permanently delete all transactions", zap.String("trace.id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete all transactions")
		status = "failed_to_permanently_delete_all_transactions"

		return false, transaction_errors.ErrFailedDeleteAllTransactionsPermanent
	}

	s.logger.Debug("Successfully deleted all transactions permanently")

	return true, nil
}

func (s *transactionCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
