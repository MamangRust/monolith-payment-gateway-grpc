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
	withdraw_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors/service"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/withdraw"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-withdraw/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

// WithdrawCommandServiceDeps holds the dependencies for WithdrawCommandService.
type withdrawCommandServiceDeps struct {

	// ErrorHandler handles command-related errors for withdraw operations.
	ErrorHandler errorhandler.WithdrawCommandErrorHandler

	// Cache provides methods to manage withdraw command cache.
	Cache mencache.WithdrawCommandCache

	// Kafka is the Kafka producer for event publishing (e.g., withdraw events).
	Kafka *kafka.Kafka

	// CardRepository accesses card data for validation or enrichment.
	CardRepository repository.CardRepository

	// SaldoRepository handles saldo (balance) logic such as deducting balance.
	SaldoRepository repository.SaldoRepository

	// CommandRepository handles writing withdraw data to persistent storage.
	CommandRepository repository.WithdrawCommandRepository

	// QueryRepository reads withdraw data, e.g., to confirm success or get history.
	QueryRepository repository.WithdrawQueryRepository

	// Logger is used to log operations and errors.
	Logger logger.LoggerInterface

	// Mapper maps withdraw domain models to response formats.
	Mapper responseservice.WithdrawCommandResponseMapper
}

// withdrawCommandService implements the WithdrawCommandService interface.
type withdrawCommandService struct {

	// ErrorHandler handles command-related errors for withdraw operations.
	errorhandler errorhandler.WithdrawCommandErrorHandler

	// Cache provides methods to manage withdraw command cache.
	mencache mencache.WithdrawCommandCache

	// Kafka is the Kafka producer for event publishing (e.g., withdraw events).
	kafka *kafka.Kafka

	// CardRepository accesses card data for validation or enrichment.
	cardRepository repository.CardRepository

	// SaldoRepository handles saldo (balance) logic such as deducting balance.
	saldoRepository repository.SaldoRepository

	// CommandRepository handles writing withdraw data to persistent storage.
	withdrawCommandRepository repository.WithdrawCommandRepository

	// QueryRepository reads withdraw data, e.g., to confirm success or get history.
	withdrawQueryRepository repository.WithdrawQueryRepository

	// Logger is used to log operations and errors.
	logger logger.LoggerInterface

	// Mapper maps withdraw domain models to response formats.
	mapper responseservice.WithdrawCommandResponseMapper

	observability observability.TraceLoggerObservability
}

// NewWithdrawCommandService initializes a new instance of withdrawCommandService with the provided parameters.
// It sets up the prometheus metrics for counting and measuring the duration of withdraw command requests.
//
// Parameters:
// - deps: A pointer to withdrawCommandServiceDeps containing the necessary dependencies.
//
// Returns:
// - A pointer to a newly created withdrawCommandService.
func NewWithdrawCommandService(
	deps *withdrawCommandServiceDeps,
) WithdrawCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "withdraw_command_service_request_total",
			Help: "Total number of requests to the WithdrawCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "withdraw_command_service_request_duration_seconds",
			Help:    "Histogram of request durations for the WithdrawCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("withdraw-command-service"), deps.Logger, requestCounter, requestDuration)

	return &withdrawCommandService{
		kafka:                     deps.Kafka,
		errorhandler:              deps.ErrorHandler,
		mencache:                  deps.Cache,
		cardRepository:            deps.CardRepository,
		saldoRepository:           deps.SaldoRepository,
		withdrawCommandRepository: deps.CommandRepository,
		withdrawQueryRepository:   deps.QueryRepository,
		logger:                    deps.Logger,
		mapper:                    deps.Mapper,
		observability:             observability,
	}
}

// Create creates a new withdraw record.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The request data to create a withdraw.
//
// Returns:
//   - *response.WithdrawResponse: The created withdraw response.
//   - *response.ErrorResponse: Error information if any occurred.
func (s *withdrawCommandService) Create(ctx context.Context, request *requests.CreateWithdrawRequest) (*response.WithdrawResponse, *response.ErrorResponse) {
	const method = "CreateWithdraw"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	card, err := s.cardRepository.FindUserCardByCardNumber(ctx, request.CardNumber)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_CARD_BY_CARD_NUMBER", span, &status, card_errors.ErrFailedFindByCardNumber, zap.Error(err))
	}

	saldo, err := s.saldoRepository.FindByCardNumber(ctx, request.CardNumber)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_SALDO_BY_CARD_NUMBER", span, &status, saldo_errors.ErrFailedFindSaldoByCardNumber, zap.Error(err))
	}

	if saldo.TotalBalance < request.WithdrawAmount {
		return s.errorhandler.HandleInsufficientBalanceError(err, method, "INSUFFICIENT_BALANCE", span, &status, request.CardNumber, zap.Error(err))
	}
	newTotalBalance := saldo.TotalBalance - request.WithdrawAmount
	updateData := &requests.UpdateSaldoWithdraw{
		CardNumber:     request.CardNumber,
		TotalBalance:   newTotalBalance,
		WithdrawAmount: &request.WithdrawAmount,
		WithdrawTime:   &request.WithdrawTime,
	}
	_, err = s.saldoRepository.UpdateSaldoWithdraw(ctx, updateData)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_SALDO", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
	}
	withdrawRecord, err := s.withdrawCommandRepository.CreateWithdraw(ctx, request)
	if err != nil {
		rollbackData := &requests.UpdateSaldoWithdraw{
			CardNumber:     request.CardNumber,
			TotalBalance:   saldo.TotalBalance,
			WithdrawAmount: &request.WithdrawAmount,
			WithdrawTime:   &request.WithdrawTime,
		}
		if _, rollbackErr := s.saldoRepository.UpdateSaldoWithdraw(ctx, rollbackData); rollbackErr != nil {
			return s.errorhandler.HandleRepositorySingleError(rollbackErr, method, "FAILED_ROLLBACK_SALDO", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
		}
		if _, err := s.withdrawCommandRepository.UpdateWithdrawStatus(ctx, &requests.UpdateWithdrawStatus{
			WithdrawID: withdrawRecord.ID,
			Status:     "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_WITHDRAW_STATUS", span, &status, withdraw_errors.ErrFailedUpdateWithdraw, zap.Error(err))
		}

		return s.errorhandler.HandleCreateWithdrawError(err, method, "FAILED_CREATE_WITHDRAW", span, &status, zap.Error(err))
	}
	if _, err := s.withdrawCommandRepository.UpdateWithdrawStatus(ctx, &requests.UpdateWithdrawStatus{
		WithdrawID: withdrawRecord.ID,
		Status:     "success",
	}); err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_WITHDRAW_STATUS", span, &status, withdraw_errors.ErrFailedUpdateWithdraw, zap.Error(err))
	}

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   "Withdraw Successful",
		"Message": fmt.Sprintf("Your withdrawal of %d has been processed successfully.", request.WithdrawAmount),
		"Button":  "View History",
		"Link":    "https://sanedge.example.com/withdraw/history",
	})

	emailPayload := map[string]any{
		"email":   card.Email,
		"subject": "Withdraw Successful - SanEdge",
		"body":    htmlBody,
	}

	payloadBytes, err := json.Marshal(emailPayload)
	if err != nil {
		return errorhandler.HandleErrorMarshal[*response.WithdrawResponse](s.logger, err, method, "FAILED_MARSHAL_EMAIL_PAYLOAD", span, &status, withdraw_errors.ErrFailedSendEmail, zap.Error(err))
	}

	err = s.kafka.SendMessage("email-service-topic-withdraw-create", strconv.Itoa(withdrawRecord.ID), payloadBytes)
	if err != nil {
		return errorhandler.HandleErrorKafkaSend[*response.WithdrawResponse](s.logger, err, method, "FAILED_SEND_EMAIL", span, &status, withdraw_errors.ErrFailedSendEmail, zap.Error(err))
	}

	so := s.mapper.ToWithdrawResponse(withdrawRecord)

	logSuccess("Successfully created withdraw", zap.Int("withdraw.id", withdrawRecord.ID))

	return so, nil
}

// Update updates an existing withdraw record.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The request data to update the withdraw.
//
// Returns:
//   - *response.WithdrawResponse: The updated withdraw response.
//   - *response.ErrorResponse: Error information if any occurred.
func (s *withdrawCommandService) Update(ctx context.Context, request *requests.UpdateWithdrawRequest) (*response.WithdrawResponse, *response.ErrorResponse) {
	const method = "UpdateWithdraw"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.withdrawQueryRepository.FindById(ctx, *request.WithdrawID)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_WITHDRAW", span, &status, withdraw_errors.ErrWithdrawNotFound)
	}

	saldo, err := s.saldoRepository.FindByCardNumber(ctx, request.CardNumber)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_SALDO", span, &status, saldo_errors.ErrFailedSaldoNotFound, zap.Error(err))
	}

	if saldo.TotalBalance < request.WithdrawAmount {
		return s.errorhandler.HandleInsufficientBalanceError(err, method, "FAILED_INSUFFICIENT_BALANCE", span, &status, request.CardNumber, zap.Error(err))
	}

	newTotalBalance := saldo.TotalBalance - request.WithdrawAmount
	updateSaldoData := &requests.UpdateSaldoWithdraw{
		CardNumber:     saldo.CardNumber,
		TotalBalance:   newTotalBalance,
		WithdrawAmount: &request.WithdrawAmount,
		WithdrawTime:   &request.WithdrawTime,
	}

	_, err = s.saldoRepository.UpdateSaldoWithdraw(ctx, updateSaldoData)
	if err != nil {
		if _, err := s.withdrawCommandRepository.UpdateWithdrawStatus(ctx, &requests.UpdateWithdrawStatus{
			WithdrawID: *request.WithdrawID,
			Status:     "failed",
		}); err != nil {
			return s.errorhandler.HandleUpdateWithdrawError(err, method, "FAILED_UPDATE_WITHDRAW_STATUS", span, &status, zap.Error(err))
		}

		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_SALDO", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
	}

	updatedWithdraw, err := s.withdrawCommandRepository.UpdateWithdraw(ctx, request)
	if err != nil {
		rollbackData := &requests.UpdateSaldoBalance{
			CardNumber:   saldo.CardNumber,
			TotalBalance: saldo.TotalBalance,
		}
		_, rollbackErr := s.saldoRepository.UpdateSaldoBalance(ctx, rollbackData)
		if rollbackErr != nil {
			return s.errorhandler.HandleUpdateWithdrawError(err, method, "FAILED_ROLLBACK_SALDO", span, &status, zap.Error(err))
		}
		if _, err := s.withdrawCommandRepository.UpdateWithdrawStatus(ctx, &requests.UpdateWithdrawStatus{
			WithdrawID: *request.WithdrawID,
			Status:     "failed",
		}); err != nil {
			return s.errorhandler.HandleUpdateWithdrawError(err, method, "FAILED_UPDATE_WITHDRAW_STATUS", span, &status, zap.Error(err))
		}

		return s.errorhandler.HandleUpdateWithdrawError(err, method, "FAILED_UPDATE_WITHDRAW", span, &status, zap.Error(err))
	}

	if _, err := s.withdrawCommandRepository.UpdateWithdrawStatus(ctx, &requests.UpdateWithdrawStatus{
		WithdrawID: updatedWithdraw.ID,
		Status:     "success",
	}); err != nil {
		return s.errorhandler.HandleUpdateWithdrawError(err, method, "FAILED_UPDATE_WITHDRAW_STATUS", span, &status, zap.Error(err))
	}

	so := s.mapper.ToWithdrawResponse(updatedWithdraw)

	s.mencache.DeleteCachedWithdrawCache(ctx, *request.WithdrawID)

	logSuccess("Successfully updated withdraw", zap.Int("withdraw.id", updatedWithdraw.ID))

	return so, nil
}

// TrashedWithdraw soft-deletes a withdraw by ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - withdraw_id: The ID of the withdraw to soft-delete.
//
// Returns:
//   - *response.WithdrawResponse: The soft-deleted withdraw response.
//   - *response.ErrorResponse: Error information if any occurred.
func (s *withdrawCommandService) TrashedWithdraw(ctx context.Context, withdraw_id int) (*response.WithdrawResponseDeleteAt, *response.ErrorResponse) {
	const method = "TrashedWithdraw"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.withdrawCommandRepository.TrashedWithdraw(ctx, withdraw_id)

	if err != nil {
		return s.errorhandler.HandleTrashedWithdrawError(err, method, "FAILED_TRASHED_WITHDRAW", span, &status, zap.Error(err))
	}

	withdrawResponse := s.mapper.ToWithdrawResponseDeleteAt(res)

	s.mencache.DeleteCachedWithdrawCache(ctx, withdraw_id)

	logSuccess("Successfully trashed withdraw", zap.Int("withdraw.id", withdraw_id))

	return withdrawResponse, nil
}

// RestoreWithdraw restores a soft-deleted withdraw by ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - withdraw_id: The ID of the withdraw to restore.
//
// Returns:
//   - *response.WithdrawResponse: The restored withdraw response.
//   - *response.ErrorResponse: Error information if any occurred.
func (s *withdrawCommandService) RestoreWithdraw(ctx context.Context, withdraw_id int) (*response.WithdrawResponse, *response.ErrorResponse) {
	const method = "RestoreWithdraw"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.withdrawCommandRepository.RestoreWithdraw(ctx, withdraw_id)

	if err != nil {
		return s.errorhandler.HandleRestoreWithdrawError(err, method, "FAILED_RESTORE_WITHDRAW", span, &status, zap.Error(err))
	}

	withdrawResponse := s.mapper.ToWithdrawResponse(res)

	logSuccess("Successfully restored withdraw", zap.Int("withdraw.id", withdraw_id))

	return withdrawResponse, nil
}

// DeleteWithdrawPermanent permanently deletes a withdraw by ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - withdraw_id: The ID of the withdraw to delete permanently.
//
// Returns:
//   - bool: True if deletion was successful.
//   - *response.ErrorResponse: Error information if any occurred.
func (s *withdrawCommandService) DeleteWithdrawPermanent(ctx context.Context, withdraw_id int) (bool, *response.ErrorResponse) {
	const method = "DeleteWithdrawPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.withdrawCommandRepository.DeleteWithdrawPermanent(ctx, withdraw_id)

	if err != nil {
		return s.errorhandler.HandleDeleteWithdrawPermanentError(err, method, "FAILED_DELETE_WITHDRAW_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted withdraw permanent", zap.Int("withdraw.id", withdraw_id))

	return true, nil
}

// RestoreAllWithdraw restores all soft-deleted withdraws.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//
// Returns:
//   - bool: True if all records were successfully restored.
//   - *response.ErrorResponse: Error information if any occurred.
func (s *withdrawCommandService) RestoreAllWithdraw(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "RestoreAllWithdraw"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.withdrawCommandRepository.RestoreAllWithdraw(ctx)

	if err != nil {
		return s.errorhandler.HandleRestoreAllWithdrawError(err, method, "FAILED_RESTORE_ALL_WITHDRAW", span, &status, zap.Error(err))
	}

	logSuccess("Successfully restored all withdraws", zap.Bool("success", true))

	return true, nil
}

// DeleteAllWithdrawPermanent permanently deletes all soft-deleted withdraws.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//
// Returns:
//   - bool: True if all records were successfully deleted.
//   - *response.ErrorResponse: Error information if any occurred.
func (s *withdrawCommandService) DeleteAllWithdrawPermanent(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "DeleteAllWithdrawPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.withdrawCommandRepository.DeleteAllWithdrawPermanent(ctx)

	if err != nil {
		return s.errorhandler.HandleDeleteAllWithdrawPermanentError(err, method, "FAILED_DELETE_ALL_WITHDRAW_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted all withdraws permanent", zap.Bool("success", true))

	return true, nil
}
