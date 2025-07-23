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
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/service"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/topup"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-topup/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

// topupCommandDeps holds the dependencies required to construct a topupCommandService.
type topupCommandDeps struct {
	// Kafka is the Kafka client used for publishing top-up related events.
	Kafka *kafka.Kafka

	// ErrorHandler handles domain-specific errors during top-up operations.
	ErrorHandler errorhandler.TopupCommandErrorHandler

	// Cache is used to manage cache invalidation or updates for top-up command operations.
	Cache mencache.TopupCommandCache

	// CardRepository provides access to card-related data used in validation or lookup.
	CardRepository repository.CardRepository

	// TopupQueryRepository provides access to read operations for top-up data.
	TopupQueryRepository repository.TopupQueryRepository

	// TopupCommandRepository provides access to write operations for top-up data.
	TopupCommandRepository repository.TopupCommandRepository

	// SaldoRepository handles saldo updates during top-up operations.
	SaldoRepository repository.SaldoRepository

	// Logger provides structured logging capabilities for debugging and observability.
	Logger logger.LoggerInterface

	// Mapper converts internal domain data to top-up response DTOs.
	Mapper responseservice.TopupCommandResponseMapper
}

// topupCommandService handles the logic for creating and managing top-up operations,
// including saldo updates, Kafka publishing, caching, and metrics tracking.
type topupCommandService struct {
	// kafka is the Kafka client used for publishing top-up related events to topics.
	kafka *kafka.Kafka

	// errorhandler handles domain-level errors that may occur during top-up operations.
	errorhandler errorhandler.TopupCommandErrorHandler

	// mencache provides access to the cache layer for invalidating or managing top-up command cache entries.
	mencache mencache.TopupCommandCache

	// topupQueryRepository provides read-only access to top-up data for validation or pre-checks.
	topupQueryRepository repository.TopupQueryRepository

	// cardRepository provides access to card data, used to validate ownership or card existence.
	cardRepository repository.CardRepository

	// topupCommandRepository provides write operations for creating or updating top-up records.
	topupCommandRepository repository.TopupCommandRepository

	// saldoRepository provides access to update saldo (balance) after top-up operations.
	saldoRepository repository.SaldoRepository

	// logger is used for structured and leveled logging throughout the service logic.
	logger logger.LoggerInterface

	// mapper transforms internal top-up domain models to response DTOs for external use.
	mapper responseservice.TopupCommandResponseMapper

	observability observability.TraceLoggerObservability
}

// NewTopupCommandService initializes a new instance of topupCommandService with the provided parameters.
// It sets up the prometheus metrics for counting and measuring the duration of top-up command requests.
//
// Parameters:
// - params: A pointer to topupCommandDeps containing the necessary dependencies.
//
// Returns:
// - A pointer to a newly created topupCommandService.
func NewTopupCommandService(
	params *topupCommandDeps,
) TopupCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "topup_command_service_request_total",
			Help: "Total number of requests to the TopupCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "topup_command_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TopupCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("topup-command-service"), params.Logger, requestCounter, requestDuration)

	return &topupCommandService{
		kafka:                  params.Kafka,
		errorhandler:           params.ErrorHandler,
		mencache:               params.Cache,
		topupQueryRepository:   params.TopupQueryRepository,
		topupCommandRepository: params.TopupCommandRepository,
		saldoRepository:        params.SaldoRepository,
		cardRepository:         params.CardRepository,
		logger:                 params.Logger,
		mapper:                 params.Mapper,
		observability:          observability,
	}
}

// CreateTopup creates a new topup record and performs the associated operations, such as updating
// the user's saldo balance and sending an email notification.
//
// Parameters:
//   - request: A pointer to a requests.CreateTopupRequest containing the user's card number,
//     topup amount, and topup method.
//
// Returns:
//   - A pointer to a response.TopupResponse containing the created topup record, if successful.
//   - A pointer to a response.ErrorResponse containing error information, if any.
func (s *topupCommandService) CreateTopup(ctx context.Context, request *requests.CreateTopupRequest) (*response.TopupResponse, *response.ErrorResponse) {
	const method = "CreateTopup"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	card, err := s.cardRepository.FindUserCardByCardNumber(ctx, request.CardNumber)
	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.TopupResponse](s.logger, err, method, "FAILED_FIND_CARD_BY_CARD_NUMBER", span, &status, card_errors.ErrCardNotFoundRes, zap.Error(err))
	}

	topup, err := s.topupCommandRepository.CreateTopup(ctx, request)
	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.TopupResponse](s.logger, err, method, "FAILED_CREATE_TOPUP", span, &status, topup_errors.ErrFailedCreateTopup, zap.Error(err))
	}

	saldo, err := s.saldoRepository.FindByCardNumber(ctx, request.CardNumber)
	if err != nil {
		req := requests.UpdateTopupStatus{
			TopupID: topup.ID,
			Status:  "failed",
		}
		s.topupCommandRepository.UpdateTopupStatus(ctx, &req)

		return errorhandler.HandleRepositorySingleError[*response.TopupResponse](s.logger, err, method, "FAILED_FIND_SALDO_BY_CARD_NUMBER", span, &status, topup_errors.ErrFailedCreateTopup, zap.Error(err))
	}

	newBalance := saldo.TotalBalance + request.TopupAmount
	_, err = s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
		CardNumber:   request.CardNumber,
		TotalBalance: newBalance,
	})
	if err != nil {

		req := requests.UpdateTopupStatus{
			TopupID: topup.ID,
			Status:  "failed",
		}
		s.topupCommandRepository.UpdateTopupStatus(ctx, &req)

		return errorhandler.HandleRepositorySingleError[*response.TopupResponse](s.logger, err, method, "FAILED_UPDATE_SALDO_BALANCE", span, &status, topup_errors.ErrFailedCreateTopup, zap.Error(err))
	}

	expireDate, err := time.Parse("2006-01-02", card.ExpireDate)
	if err != nil {
		req := requests.UpdateTopupStatus{
			TopupID: topup.ID,
			Status:  "failed",
		}
		s.topupCommandRepository.UpdateTopupStatus(ctx, &req)

		return s.errorhandler.HandleInvalidParseTimeError(err, "CreateTopup", "FAILED_PARSE_EXPIRE_DATE", span, &status, card.ExpireDate, zap.Error(err))
	}

	_, err = s.cardRepository.UpdateCard(ctx, &requests.UpdateCardRequest{
		CardID:       card.ID,
		UserID:       card.UserID,
		CardType:     card.CardType,
		ExpireDate:   expireDate,
		CVV:          card.CVV,
		CardProvider: card.CardProvider,
	})
	if err != nil {
		req := requests.UpdateTopupStatus{
			TopupID: topup.ID,
			Status:  "failed",
		}
		s.topupCommandRepository.UpdateTopupStatus(ctx, &req)

		return errorhandler.HandleRepositorySingleError[*response.TopupResponse](s.logger, err, method, "FAILED_UPDATE_CARD", span, &status, topup_errors.ErrFailedCreateTopup, zap.Error(err))
	}

	req := requests.UpdateTopupStatus{
		TopupID: topup.ID,
		Status:  "success",
	}

	res, err := s.topupCommandRepository.UpdateTopupStatus(ctx, &req)
	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.TopupResponse](s.logger, err, method, "FAILED_UPDATE_TOPUP_STATUS", span, &status, topup_errors.ErrFailedUpdateTopup, zap.Error(err))
	}

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   "Topup Successful",
		"Message": fmt.Sprintf("Your topup of %d has been processed successfully.", request.TopupAmount),
		"Button":  "View History",
		"Link":    "https://sanedge.example.com/topup/history",
	})

	emailPayload := map[string]any{
		"email":   card.Email,
		"subject": "Topup Successful - SanEdge",
		"body":    htmlBody,
	}

	payloadBytes, err := json.Marshal(emailPayload)
	if err != nil {
		return errorhandler.HandleErrorJSONMarshal[*response.TopupResponse](s.logger, err, "CreateTopup", "FAILED_JSON_MARSHAL", span, &status, topup_errors.ErrFailedUpdateTopup)
	}

	err = s.kafka.SendMessage("email-service-topic-topup-create", strconv.Itoa(res.ID), payloadBytes)
	if err != nil {
		return errorhandler.HandleErrorKafkaSend[*response.TopupResponse](s.logger, err, "CreateTopup", "FAILED_KAFKA_SEND", span, &status, topup_errors.ErrFailedUpdateTopup)
	}

	so := s.mapper.ToTopupResponse(topup)

	logSuccess("Topup created successfully",
		zap.String("cardNumber", request.CardNumber),
		zap.Int("topupID", topup.ID),
		zap.Float64("topupAmount", float64(request.TopupAmount)),
	)

	return so, nil
}

// UpdateTopup updates an existing topup record with the provided details.
//
// Parameters:
//   - request: A pointer to a requests.UpdateTopupRequest containing the topup ID,
//     card number, topup amount, and topup method.
//
// Returns:
//   - A pointer to a response.TopupResponse containing the updated topup record, if successful.
//   - A pointer to a response.ErrorResponse containing error information, if any.
func (s *topupCommandService) UpdateTopup(ctx context.Context, request *requests.UpdateTopupRequest) (*response.TopupResponse, *response.ErrorResponse) {
	const method = "UpdateTopup"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.cardRepository.FindCardByCardNumber(ctx, request.CardNumber)
	if err != nil {
		req := requests.UpdateTopupStatus{
			TopupID: *request.TopupID,
			Status:  "failed",
		}

		s.topupCommandRepository.UpdateTopupStatus(ctx, &req)

		return errorhandler.HandleRepositorySingleError[*response.TopupResponse](s.logger, err, method, "FAILED_FIND_CARD_BY_CARD_NUMBER", span, &status, topup_errors.ErrFailedUpdateTopup, zap.Error(err))
	}

	existingTopup, err := s.topupQueryRepository.FindById(ctx, *request.TopupID)
	if err != nil || existingTopup == nil {
		req := requests.UpdateTopupStatus{
			TopupID: *request.TopupID,
			Status:  "failed",
		}

		s.topupCommandRepository.UpdateTopupStatus(ctx, &req)

		return errorhandler.HandleRepositorySingleError[*response.TopupResponse](s.logger, err, method, "FAILED_FIND_TOPUP_BY_ID", span, &status, topup_errors.ErrFailedUpdateTopup, zap.Error(err))
	}

	topupDifference := request.TopupAmount - existingTopup.TopupAmount

	_, err = s.topupCommandRepository.UpdateTopup(ctx, request)
	if err != nil {
		req := requests.UpdateTopupStatus{
			TopupID: *request.TopupID,
			Status:  "failed",
		}

		s.topupCommandRepository.UpdateTopupStatus(ctx, &req)

		return errorhandler.HandleRepositorySingleError[*response.TopupResponse](s.logger, err, method, "FAILED_UPDATE_TOPUP", span, &status, topup_errors.ErrFailedUpdateTopup, zap.Error(err))
	}

	currentSaldo, err := s.saldoRepository.FindByCardNumber(ctx, request.CardNumber)
	if err != nil {
		req := requests.UpdateTopupStatus{
			TopupID: *request.TopupID,
			Status:  "failed",
		}

		s.topupCommandRepository.UpdateTopupStatus(ctx, &req)

		return errorhandler.HandleRepositorySingleError[*response.TopupResponse](s.logger, err, method, "FAILED_FIND_SALDO_BY_CARD_NUMBER", span, &status, topup_errors.ErrFailedUpdateTopup, zap.Error(err))
	}

	if currentSaldo == nil {
		req := requests.UpdateTopupStatus{
			TopupID: *request.TopupID,
			Status:  "failed",
		}

		s.topupCommandRepository.UpdateTopupStatus(ctx, &req)

		return errorhandler.HandleRepositorySingleError[*response.TopupResponse](s.logger, err, method, "FAILED_FIND_SALDO_BY_CARD_NUMBER", span, &status, topup_errors.ErrFailedUpdateTopup, zap.Error(err))
	}

	newBalance := currentSaldo.TotalBalance + topupDifference
	_, err = s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
		CardNumber:   request.CardNumber,
		TotalBalance: newBalance,
	})
	if err != nil {
		_, rollbackErr := s.topupCommandRepository.UpdateTopupAmount(ctx, &requests.UpdateTopupAmount{
			TopupID:     *request.TopupID,
			TopupAmount: existingTopup.TopupAmount,
		})
		if rollbackErr != nil {
			return errorhandler.HandleRepositorySingleError[*response.TopupResponse](s.logger, err, method, "FAILED_ROLLBACK_TOPUP_AMOUNT", span, &status, topup_errors.ErrFailedUpdateTopup, zap.Error(err))
		}

		req := requests.UpdateTopupStatus{
			TopupID: *request.TopupID,
			Status:  "failed",
		}

		s.topupCommandRepository.UpdateTopupStatus(ctx, &req)

		return errorhandler.HandleRepositorySingleError[*response.TopupResponse](s.logger, err, method, "FAILED_UPDATE_SALDO_BALANCE", span, &status, topup_errors.ErrFailedUpdateTopup, zap.Error(err))
	}

	updatedTopup, err := s.topupQueryRepository.FindById(ctx, *request.TopupID)
	if err != nil || updatedTopup == nil {
		req := requests.UpdateTopupStatus{
			TopupID: *request.TopupID,
			Status:  "failed",
		}

		s.topupCommandRepository.UpdateTopupStatus(ctx, &req)

		return errorhandler.HandleRepositorySingleError[*response.TopupResponse](s.logger, err, method, "FAILED_FIND_TOPUP_BY_ID", span, &status, topup_errors.ErrFailedUpdateTopup, zap.Error(err))
	}

	req := requests.UpdateTopupStatus{
		TopupID: *request.TopupID,
		Status:  "success",
	}

	_, err = s.topupCommandRepository.UpdateTopupStatus(ctx, &req)
	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.TopupResponse](s.logger, err, method, "FAILED_UPDATE_TOPUP_STATUS", span, &status, topup_errors.ErrFailedUpdateTopup, zap.Error(err))
	}

	so := s.mapper.ToTopupResponse(updatedTopup)

	s.mencache.DeleteCachedTopupCache(ctx, *request.TopupID)

	logSuccess("UpdateTopup process completed", zap.Bool("success", true))

	return so, nil
}

// TrashedTopup marks a topup record as trashed by its ID.
//
// It returns a response.TopupResponseDeleteAt pointer containing the trashed topup record, and a response.ErrorResponse pointer if something goes wrong. The error response is nil if no error occurred.
//
// Parameters:
//   - topup_id: The ID of the topup record to trash.
//
// Returns:
//   - A response.TopupResponseDeleteAt pointer containing the trashed topup record.
//   - A response.ErrorResponse pointer containing error information.
func (s *topupCommandService) TrashedTopup(ctx context.Context, topup_id int) (*response.TopupResponseDeleteAt, *response.ErrorResponse) {
	const method = "TrashedTopup"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.topupCommandRepository.TrashedTopup(ctx, topup_id)

	if err != nil {
		return s.errorhandler.HandleTrashedTopupError(err, method, "FAILED_TRASH_TOPUP", span, &status, zap.Error(err))
	}

	so := s.mapper.ToTopupResponseDeleteAt(res)

	s.mencache.DeleteCachedTopupCache(ctx, topup_id)

	logSuccess("TrashedTopup process completed", zap.Bool("success", true))

	return so, nil
}

// RestoreTopup marks a trashed topup record as not trashed by its ID.
//
// It returns a response.TopupResponseDeleteAt pointer containing the restored topup record, and a response.ErrorResponse pointer if something goes wrong. The error response is nil if no error occurred.
//
// Parameters:
//   - topup_id: The ID of the topup record to restore.
//
// Returns:
//   - A response.TokenResponse pointer containing the restored topup record.
//   - A response.ErrorResponse pointer containing error information.
func (s *topupCommandService) RestoreTopup(ctx context.Context, topup_id int) (*response.TopupResponse, *response.ErrorResponse) {
	const method = "RestoreTopup"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.topupCommandRepository.RestoreTopup(ctx, topup_id)

	if err != nil {
		return s.errorhandler.HandleRestoreTopupError(err, method, "FAILED_RESTORE_TOPUP", span, &status, zap.Error(err))
	}

	so := s.mapper.ToTopupResponse(res)

	logSuccess("RestoreTopup process completed", zap.Bool("success", true))

	return so, nil
}

// DeleteTopupPermanent permanently deletes a topup record by its ID.
//
// It returns a boolean indicating success, and a response.ErrorResponse pointer if something goes wrong.
// The error response is nil if no error occurred.
//
// Parameters:
//   - topup_id: The ID of the topup record to permanently delete.
//
// Returns:
//   - A boolean indicating whether the deletion was successful.
//   - A response.ErrorResponse pointer containing error information, if any.
func (s *topupCommandService) DeleteTopupPermanent(ctx context.Context, topup_id int) (bool, *response.ErrorResponse) {
	const method = "DeleteTopupPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.topupCommandRepository.DeleteTopupPermanent(ctx, topup_id)

	if err != nil {
		return s.errorhandler.HandleDeleteTopupPermanentError(err, method, "FAILED_DELETE_TOPUP_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("DeleteTopupPermanent process completed", zap.Bool("success", true))

	return true, nil
}

// RestoreAllTopup restores all trashed topup records in the database.
//
// It returns a boolean indicating whether the operation was successful,
// and a response.ErrorResponse pointer if something goes wrong. The error
// response is nil if no error occurred.
//
// Returns:
//   - A boolean indicating whether the operation was successful.
//   - A response.ErrorResponse pointer containing error information, if any.
func (s *topupCommandService) RestoreAllTopup(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "RestoreAllTopup"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.topupCommandRepository.RestoreAllTopup(ctx)

	if err != nil {
		return s.errorhandler.HandleRestoreAllTopupError(err, method, "FAILED_RESTORE_ALL_TOPUP", span, &status, zap.Error(err))
	}

	logSuccess("RestoreAllTopup process completed", zap.Bool("success", true))

	return true, nil
}

// DeleteAllTopupPermanent permanently deletes all topup records from the database.
//
// Returns:
//   - A boolean indicating whether the operation was successful.
//   - A response.ErrorResponse pointer containing error information, if any.
func (s *topupCommandService) DeleteAllTopupPermanent(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "DeleteAllTopupPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.topupCommandRepository.DeleteAllTopupPermanent(ctx)

	if err != nil {
		return s.errorhandler.HandleDeleteAllTopupPermanentError(err, method, "FAILED_DELETE_ALL_TOPUP_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("DeleteAllTopupPermanent process completed", zap.Bool("success", true))

	return true, nil
}
