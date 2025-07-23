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
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/service"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/service"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/service"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/service"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/transaction"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-transaction/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// TransactionCommandServiceDeps defines the dependencies required to initialize a TransactionCommandService.
type transactionCommandServiceDeps struct {
	// Kafka producer for publishing transaction events.
	Kafka *kafka.Kafka

	// Context for base operations (usually request context).
	Ctx context.Context

	// Error handler to process and log errors consistently.
	ErrorHandler errorhandler.TransactionCommandErrorHandler

	// Redis cache layer used for storing command-side data temporarily.
	Mencache mencache.TransactionCommandCache

	// OpenTelemetry tracer for distributed tracing.
	Tracer trace.Tracer

	// Repository to access merchant-related data.
	MerchantRepository repository.MerchantRepository

	// Repository to access card-related data.
	CardRepository repository.CardRepository

	// Repository for saldo (balance) operations.
	SaldoRepository repository.SaldoRepository

	// Repository for reading/querying transaction data.
	TransactionQueryRepository repository.TransactionQueryRepository

	// Repository for writing/updating transaction data.
	TransactionCommandRepository repository.TransactionCommandRepository

	// Structured logger interface.
	Logger logger.LoggerInterface

	// Mapper for converting transaction records to response DTOs.
	Mapping responseservice.TransactionCommandResponseMapper
}

// transactionCommandService provides command-side business logic related to transactions,
// such as creating new transactions, updating statuses, publishing events to Kafka,
// and interacting with Redis cache and repositories.
type transactionCommandService struct {
	// kafka is the Kafka producer used to publish transaction-related events.
	kafka *kafka.Kafka

	// ctx is the base context used for all operations in the service.
	ctx context.Context

	// errorhandler handles standardized error responses and telemetry for command operations.
	errorhandler errorhandler.TransactionCommandErrorHandler

	// mencache is the Redis-based cache layer used for transaction command data.
	mencache mencache.TransactionCommandCache

	// merchantRepository handles database operations related to merchants.
	merchantRepository repository.MerchantRepository

	// cardRepository handles database operations related to cards.
	cardRepository repository.CardRepository

	// saldoRepository handles database operations related to saldo/balance.
	saldoRepository repository.SaldoRepository

	// transactionQueryRepository is used to query historical transaction data when needed.
	transactionQueryRepository repository.TransactionQueryRepository

	// transactionCommandRepository handles writes and updates to transaction data.
	transactionCommandRepository repository.TransactionCommandRepository

	// logger is the structured logger interface used to log service activities and errors.
	logger logger.LoggerInterface

	// mapper provides functionality to map internal transaction data to response models.
	mapper responseservice.TransactionCommandResponseMapper

	observability observability.TraceLoggerObservability
}

// NewTransactionCommandService initializes a new instance of transactionCommandService with the provided parameters.
// It sets up Prometheus metrics for tracking request counts and durations and returns a configured
// transactionCommandService ready for handling transaction-related commands.
//
// Parameters:
// - params: A pointer to transactionCommandServiceDeps containing the necessary dependencies.
//
// Returns:
// - A pointer to an initialized transactionCommandService.
func NewTransactionCommandService(
	params *transactionCommandServiceDeps,
) TransactionCommandService {
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

	observability := observability.NewTraceLoggerObservability(
		otel.Tracer("transaction-command-service"), params.Logger, requestCounter, requestDuration)

	return &transactionCommandService{
		kafka:                        params.Kafka,
		ctx:                          params.Ctx,
		errorhandler:                 params.ErrorHandler,
		mencache:                     params.Mencache,
		merchantRepository:           params.MerchantRepository,
		cardRepository:               params.CardRepository,
		saldoRepository:              params.SaldoRepository,
		transactionCommandRepository: params.TransactionCommandRepository,
		transactionQueryRepository:   params.TransactionQueryRepository,
		logger:                       params.Logger,
		mapper:                       params.Mapping,
		observability:                observability,
	}
}

// Create creates a new transaction based on the provided request and API key.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - apiKey: The API key for merchant authorization.
//   - request: The transaction creation request payload.
//
// Returns:
//   - *response.TransactionResponse: The created transaction response.
//   - *response.ErrorResponse: Error detail if the operation fails.
func (s *transactionCommandService) Create(ctx context.Context, apiKey string, request *requests.CreateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse) {
	const method = "CreateTransacton"

	s.logger.Debug("CreateTransaction called", zap.String("card_number", request.CardNumber), zap.String("api_key", apiKey))

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("apikey", apiKey))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantRepository.FindByApiKey(ctx, apiKey)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_MERCHANT_BY_API_KEY", span, &status, merchant_errors.ErrFailedFindByApiKey, zap.Error(err))
	}

	card, err := s.cardRepository.FindUserCardByCardNumber(ctx, request.CardNumber)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_CARD_BY_CARD_NUMBER", span, &status, card_errors.ErrFailedFindByCardNumber, zap.Error(err))
	}

	saldo, err := s.saldoRepository.FindByCardNumber(ctx, card.CardNumber)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_SALDO_BY_CARD_NUMBER", span, &status, saldo_errors.ErrFailedFindSaldoByCardNumber, zap.Error(err))
	}

	if saldo.TotalBalance < request.Amount {
		return s.errorhandler.HandleInsufficientBalanceError(err, method, "FAILED_INSUFFICIENT_BALANCE", span, &status, card.CardNumber, zap.Error(err))
	}

	saldo.TotalBalance -= request.Amount
	if _, err := s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
		CardNumber:   card.CardNumber,
		TotalBalance: saldo.TotalBalance,
	}); err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_SALDO", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
	}

	request.MerchantID = &merchant.ID

	transaction, err := s.transactionCommandRepository.CreateTransaction(ctx, request)
	if err != nil {
		saldo.TotalBalance += request.Amount
		_, err := s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
			CardNumber:   card.CardNumber,
			TotalBalance: saldo.TotalBalance,
		})
		if err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_SALDO", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
		}

		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(ctx, &requests.UpdateTransactionStatus{
			TransactionID: transaction.ID,
			Status:        "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, transaction_errors.ErrFailedUpdateTransaction, zap.Error(err))
		}

		return s.errorhandler.HandleCreateTransactionError(err, method, "FAILED_CREATE_TRANSACTION", span, &status, zap.Error(err))
	}

	if _, err := s.transactionCommandRepository.UpdateTransactionStatus(ctx, &requests.UpdateTransactionStatus{
		TransactionID: transaction.ID,
		Status:        "success",
	}); err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, transaction_errors.ErrFailedUpdateTransaction, zap.Error(err))
	}

	merchantCard, err := s.cardRepository.FindCardByUserId(ctx, merchant.UserID)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_CARD_BY_USER_ID_MERCHANT", span, &status, card_errors.ErrFailedFindByCardNumber, zap.Error(err))
	}

	merchantSaldo, err := s.saldoRepository.FindByCardNumber(ctx, merchantCard.CardNumber)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_SALDO_BY_CARD_NUMBER_MERCHANT", span, &status, merchant_errors.ErrFailedFindByApiKey, zap.Error(err))
	}

	merchantSaldo.TotalBalance += request.Amount

	if _, err := s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
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

	so := s.mapper.ToTransactionResponse(transaction)

	logSuccess("Successfully created transaction", zap.Int("transaction.id", transaction.ID))

	return so, nil
}

// Update updates an existing transaction with the given request and API key.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - apiKey: The API key for merchant authorization.
//   - request: The transaction update request payload.
//
// Returns:
//   - *response.TransactionResponse: The updated transaction response.
//   - *response.ErrorResponse: Error detail if the operation fails.
func (s *transactionCommandService) Update(ctx context.Context, apiKey string, request *requests.UpdateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse) {
	const method = "UpdateTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	transaction, err := s.transactionQueryRepository.FindById(ctx, *request.TransactionID)
	if err != nil {
		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(ctx, &requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, transaction_errors.ErrFailedUpdateTransaction, zap.Error(err))
		}

		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_TRANSACTION_BY_ID", span, &status, transaction_errors.ErrFailedUpdateTransaction, zap.Error(err))
	}

	merchant, err := s.merchantRepository.FindByApiKey(ctx, apiKey)
	if err != nil || transaction.MerchantID != merchant.ID {
		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(ctx, &requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, merchant_errors.ErrFailedFindByApiKey, zap.Error(err))
		}

		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_MERCHANT_BY_API_KEY", span, &status, merchant_errors.ErrFailedFindByApiKey, zap.Error(err))
	}

	card, err := s.cardRepository.FindCardByCardNumber(ctx, transaction.CardNumber)
	if err != nil {
		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(ctx, &requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			return s.errorhandler.HandleUpdateTransactionError(err, method, "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, zap.Error(err))
		}

		return nil, card_errors.ErrCardNotFoundRes
	}

	saldo, err := s.saldoRepository.FindByCardNumber(ctx, card.CardNumber)
	if err != nil {
		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(ctx, &requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			return s.errorhandler.HandleUpdateTransactionError(err, method, "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, zap.Error(err))
		}

		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_SALDO_BY_CARD_NUMBER", span, &status, saldo_errors.ErrFailedSaldoNotFound, zap.Error(err))
	}

	saldo.TotalBalance += transaction.Amount
	s.logger.Info("Restoring balance for old transaction amount", zap.Int("RestoredBalance", saldo.TotalBalance))

	if _, err := s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
		CardNumber:   card.CardNumber,
		TotalBalance: saldo.TotalBalance,
	}); err != nil {
		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(ctx, &requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			return s.errorhandler.HandleUpdateTransactionError(err, method, "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, zap.Error(err))
		}

		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_SALDO_BALANCE", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
	}

	if saldo.TotalBalance < request.Amount {
		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(ctx, &requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			return s.errorhandler.HandleUpdateTransactionError(err, method, "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, zap.Error(err))

		}
		return s.errorhandler.HandleInsufficientBalanceError(err, method, "INSUFFICIENT_BALANCE", span, &status, card.CardNumber, zap.Error(err))
	}

	saldo.TotalBalance -= request.Amount
	s.logger.Info("Updating balance for updated transaction amount")

	if _, err := s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
		CardNumber:   card.CardNumber,
		TotalBalance: saldo.TotalBalance,
	}); err != nil {
		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(ctx, &requests.UpdateTransactionStatus{
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
		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(ctx, &requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			return s.errorhandler.HandleUpdateTransactionError(err, method, "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, zap.Error(err))
		}

		return s.errorhandler.HandleInvalidParseTimeError(err, method, "INVALID_TRANSACTION_TIME", span, &status, transaction.TransactionTime, zap.Error(err))
	}

	res, err := s.transactionCommandRepository.UpdateTransaction(ctx, &requests.UpdateTransactionRequest{
		TransactionID:   &transaction.ID,
		CardNumber:      transaction.CardNumber,
		Amount:          transaction.Amount,
		PaymentMethod:   transaction.PaymentMethod,
		MerchantID:      &transaction.MerchantID,
		TransactionTime: parsedTime,
	})
	if err != nil {
		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(ctx, &requests.UpdateTransactionStatus{
			TransactionID: *request.TransactionID,
			Status:        "failed",
		}); err != nil {
			return s.errorhandler.HandleUpdateTransactionError(err, method, "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, zap.Error(err))
		}

		return s.errorhandler.HandleUpdateTransactionError(err, method, "FAILED_UPDATE_TRANSACTION", span, &status, zap.Error(err))
	}

	if _, err := s.transactionCommandRepository.UpdateTransactionStatus(ctx, &requests.UpdateTransactionStatus{
		TransactionID: transaction.ID,
		Status:        "success",
	}); err != nil {
		return s.errorhandler.HandleUpdateTransactionError(err, method, "FAILED_UPDATE_TRANSACTION_STATUS", span, &status, zap.Error(err))
	}

	so := s.mapper.ToTransactionResponse(res)

	s.mencache.DeleteTransactionCache(ctx, *request.TransactionID)

	logSuccess("Successfully updated transaction", zap.Bool("success", true))

	return so, nil
}

// TrashedTransaction moves the transaction to the trash (soft delete).
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - transaction_id: The ID of the transaction to be trashed.
//
// Returns:
//   - *response.TransactionResponse: The trashed transaction response.
//   - *response.ErrorResponse: Error detail if the operation fails.
func (s *transactionCommandService) TrashedTransaction(ctx context.Context, transaction_id int) (*response.TransactionResponseDeleteAt, *response.ErrorResponse) {
	const method = "TrashedTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.transactionCommandRepository.TrashedTransaction(ctx, transaction_id)

	if err != nil {
		return s.errorhandler.HandleTrashedTransactionError(err, method, "FAILED_TRASHED_TRANSACTION", span, &status, zap.Error(err))
	}

	so := s.mapper.ToTransactionResponseDeleteAt(res)

	s.mencache.DeleteTransactionCache(ctx, transaction_id)

	logSuccess("Successfully trashed transaction", zap.Bool("success", true))

	return so, nil
}

// RestoreTransaction restores a previously trashed transaction.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - transaction_id: The ID of the transaction to be restored.
//
// Returns:
//   - *response.TransactionResponse: The restored transaction response.
//   - *response.ErrorResponse: Error detail if the operation fails.
func (s *transactionCommandService) RestoreTransaction(ctx context.Context, transaction_id int) (*response.TransactionResponse, *response.ErrorResponse) {
	const method = "RestoreTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.transactionCommandRepository.RestoreTransaction(ctx, transaction_id)

	if err != nil {
		return s.errorhandler.HandleRestoreTransactionError(err, method, "FAILED_RESTORE_TRANSACTION", span, &status, zap.Error(err))
	}

	so := s.mapper.ToTransactionResponse(res)

	logSuccess("Successfully restored transaction", zap.Bool("success", true))

	return so, nil
}

// DeleteTransactionPermanent permanently deletes a transaction from the database.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - transaction_id: The ID of the transaction to delete permanently.
//
// Returns:
//   - bool: Whether the operation was successful.
//   - *response.ErrorResponse: Error detail if the operation fails.
func (s *transactionCommandService) DeleteTransactionPermanent(ctx context.Context, transaction_id int) (bool, *response.ErrorResponse) {
	const method = "DeleteTransactionPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.transactionCommandRepository.DeleteTransactionPermanent(ctx, transaction_id)

	if err != nil {
		return s.errorhandler.HandleDeleteTransactionPermanentError(err, method, "FAILED_DELETE_TRANSACTION_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted transaction permanent", zap.Bool("success", true))

	return true, nil
}

// RestoreAllTransaction restores all trashed transactions.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//
// Returns:
//   - bool: Whether the operation was successful.
//   - *response.ErrorResponse: Error detail if the operation fails.
func (s *transactionCommandService) RestoreAllTransaction(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "RestoreAllTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.transactionCommandRepository.RestoreAllTransaction(ctx)
	if err != nil {
		return s.errorhandler.HandleRestoreAllTransactionError(err, method, "FAILED_RESTORE_ALL_TRANSACTIONS", span, &status, zap.Error(err))
	}

	logSuccess("Successfully restored all transactions", zap.Bool("success", true))

	return true, nil
}

// DeleteAllTransactionPermanent permanently deletes all trashed transactions.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//
// Returns:
//   - bool: Whether the operation was successful.
//   - *response.ErrorResponse: Error detail if the operation fails.
func (s *transactionCommandService) DeleteAllTransactionPermanent(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "DeleteAllTransactionPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.transactionCommandRepository.DeleteAllTransactionPermanent(ctx)

	if err != nil {
		return s.errorhandler.HandleDeleteAllTransactionPermanentError(err, method, "FAILED_DELETE_ALL_TRANSACTION_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted all transactions permanent", zap.Bool("success", true))

	return true, nil
}
