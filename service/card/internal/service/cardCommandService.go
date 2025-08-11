package service

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/MamangRust/monolith-payment-gateway-card/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-card/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-card/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/service"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/card"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// cardCommandServiceDeps holds the dependencies required to create a new cardCommandService.
type cardCommandServiceDeps struct {
	// errorHandler handles domain-specific errors for the card command service.
	ErrorHandler errorhandler.CardCommandErrorHandler

	// cache is the in-memory cache layer for card commands.
	Cache mencache.CardCommandCache

	// kafka is the Kafka client used to publish card-related events.
	Kafka *kafka.Kafka

	// userRepository provides access to user data for validation or enrichment.
	UserRepository repository.UserRepository

	// cardCommandRepository provides access to persistent storage for card commands.
	CardCommandRepository repository.CardCommandRepository

	// logger is used to log service activity and errors.
	Logger logger.LoggerInterface

	// mapper converts internal data models to API response formats.
	Mapper responseservice.CardCommandResponseMapper
}

// cardCommandService implements business logic for card-related command operations.
type cardCommandService struct {

	// errorhandler handles domain-level errors for the card command service.
	errorhandler errorhandler.CardCommandErrorHandler

	// mencache provides cache functionality for reducing load and improving performance.
	mencache mencache.CardCommandCache

	// kafka is the Kafka producer used for emitting card-related events.
	kafka *kafka.Kafka

	// userRepository is used to fetch user data related to card operations.
	userRepository repository.UserRepository

	// cardCommentRepository handles database interactions for card commands.
	cardCommentRepository repository.CardCommandRepository

	// logger enables structured logging within the service.
	logger logger.LoggerInterface

	// mapper is responsible for converting internal entities to response DTOs.
	mapper responseservice.CardCommandResponseMapper

	observability observability.TraceLoggerObservability
}

// NewCardCommandService initializes a new instance of cardCommandService with the provided parameters.
// It sets up Prometheus metrics for tracking request counts and durations and returns a configured
// cardCommandService ready for handling card-related commands.
//
// Parameters:
// - params: A pointer to cardCommandServiceDeps containing the necessary dependencies.
//
// Returns:
// - A pointer to an initialized cardCommandService.
func NewCardCommandService(params *cardCommandServiceDeps) CardCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "card_command_service_requests_total",
			Help: "Total number of requests to the CardCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "card_command_service_request_duration_seconds",
			Help:    "Histogram of request durations for the CardCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("card-command-service"), params.Logger, requestCounter, requestDuration)

	return &cardCommandService{
		errorhandler:          params.ErrorHandler,
		mencache:              params.Cache,
		kafka:                 params.Kafka,
		userRepository:        params.UserRepository,
		cardCommentRepository: params.CardCommandRepository,
		logger:                params.Logger,
		mapper:                params.Mapper,
		observability:         observability,
	}
}

// CreateCard creates a new card for a user.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - request: A requests.CreateCardRequest object containing the details of the card to be created.
//
// Returns:
//   - A pointer to a responses.CardResponse object containing the created card's info.
//   - A pointer to a responses.ErrorResponse object describing the error if the operation fails.
func (s *cardCommandService) CreateCard(ctx context.Context, request *requests.CreateCardRequest) (*response.CardResponse, *response.ErrorResponse) {
	const method = "CreateCard"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.userRepository.FindById(ctx, request.UserID)

	if err != nil {
		return s.errorhandler.HandleFindByIdUserError(err, "CreateCard", "FAILED_USER_NOT_FOUND", span, &status, zap.Error(err))
	}

	res, err := s.cardCommentRepository.CreateCard(ctx, request)

	if err != nil {
		return s.errorhandler.HandleCreateCardError(err, "CreateCard", "FAILED_CREATE_CARD", span, &status, zap.Error(err))
	}

	saldoPayload := map[string]any{
		"card_number":   res.CardNumber,
		"total_balance": 0,
	}

	payloadBytes, err := json.Marshal(saldoPayload)

	s.logger.Info("hello world", zap.Any("payloadBytes", payloadBytes))

	if err != nil {
		return errorhandler.HandleMarshalError[*response.CardResponse](s.logger, err, "CreateCard", "FAILED_CREATE_CARD", span, &status, card_errors.ErrFailedCreateCard, zap.Error(err))
	}

	err = s.kafka.SendMessage("saldo-service-topic-create-saldo", strconv.Itoa(res.ID), payloadBytes)

	if err != nil {
		return errorhandler.HandleSendEmailError[*response.CardResponse](s.logger, err, "CreateCard", "FAILED_CREATE_CARD", span, &status, card_errors.ErrFailedCreateCard, zap.Error(err))
	}

	so := s.mapper.ToCardResponse(res)

	logSuccess("Successfully created card", zap.Int("card.id", so.ID), zap.String("card.card_number", so.CardNumber))

	return so, nil
}

// UpdateCard updates a card record in the database.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - request: An UpdateCardRequest object containing the details of the card to be updated.
//
// Returns:
//   - A pointer to the updated CardResponse, or an error if the operation fails.
func (s *cardCommandService) UpdateCard(ctx context.Context, request *requests.UpdateCardRequest) (*response.CardResponse, *response.ErrorResponse) {
	const method = "UpdateCard"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.userRepository.FindById(ctx, request.UserID)

	if err != nil {
		status = "error"

		return s.errorhandler.HandleFindByIdUserError(err, "UpdateCard", "FAILED_USER_NOT_FOUND", span, &status, zap.Error(err))
	}

	res, err := s.cardCommentRepository.UpdateCard(ctx, request)

	if err != nil {
		return s.errorhandler.HandleUpdateCardError(err, "UpdateCard", "FAILED_UPDATE_CARD", span, &status, zap.Error(err))
	}

	so := s.mapper.ToCardResponse(res)

	s.mencache.DeleteCardCommandCache(ctx, request.CardID)

	logSuccess("Successfully updated card", zap.Int("card.id", so.ID))

	return so, nil
}

// TrashedCard marks a card as trashed by updating its deleted_at field in the database
// and removes its cache entry.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - card_id: The ID of the card to be trashed.
//
// Returns:
//   - A pointer to a CardResponse containing the trashed card's info, or an error
//     if the operation fails.
func (s *cardCommandService) TrashedCard(ctx context.Context, card_id int) (*response.CardResponseDeleteAt, *response.ErrorResponse) {
	const method = "TrashedCard"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("card.id", card_id))

	defer func() {
		end(status)
	}()

	res, err := s.cardCommentRepository.TrashedCard(ctx, card_id)

	if err != nil {
		return s.errorhandler.HandleTrashedCardError(err, "TrashedCard", "FAILED_TO_TRASH_CARD", span, &status, zap.Error(err))
	}

	so := s.mapper.ToCardResponseDeleteAt(res)

	s.mencache.DeleteCardCommandCache(ctx, card_id)

	logSuccess("Successfully trashed card", zap.Int("card.id", so.ID))

	return so, nil
}

// RestoreCard restores a previously trashed card by setting its deleted_at field to NULL.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - card_id: The ID of the card to be restored.
//
// Returns:
//   - A pointer to the restored CardResponse, or an error if the operation fails.
func (s *cardCommandService) RestoreCard(ctx context.Context, card_id int) (*response.CardResponse, *response.ErrorResponse) {
	const method = "RestoreCard"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("card.id", card_id))

	defer func() {
		end(status)
	}()

	res, err := s.cardCommentRepository.RestoreCard(ctx, card_id)

	if err != nil {
		return s.errorhandler.HandleRestoreCardError(err, "RestoreCard", "FAILED_TO_RESTORE_CARD", span, &status, zap.Error(err))
	}

	so := s.mapper.ToCardResponse(res)

	logSuccess("Successfully restored card", zap.Int("card.id", so.ID))

	return so, nil
}

// DeleteCardPermanent permanently deletes a card by its ID from the database
// and removes its cache entry.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - card_id: The ID of the card to be deleted permanently.
//
// Returns:
//   - A boolean indicating if the operation was successful, and an error
//     if the operation fails.
func (s *cardCommandService) DeleteCardPermanent(ctx context.Context, card_id int) (bool, *response.ErrorResponse) {
	const method = "DeleteCardPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("card.id", card_id))

	defer func() {
		end(status)
	}()

	_, err := s.cardCommentRepository.DeleteCardPermanent(ctx, card_id)

	if err != nil {
		return s.errorhandler.HandleDeleteCardPermanentError(err, "DeleteCardPermanent", "FAILED_TO_DELETE_CARD_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted card permanently", zap.Int("card.id", card_id))

	return true, nil
}

// RestoreAllCard restores all previously trashed card records by setting their deleted_at fields to NULL.
//
// Parameters:
//   - ctx: The context for the database operation
//
// Returns:
//   - A boolean indicating if the operation was successful, and an error
//     if the operation fails.
func (s *cardCommandService) RestoreAllCard(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "RestoreAllCard"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.cardCommentRepository.RestoreAllCard(ctx)

	if err != nil {
		return s.errorhandler.HandleRestoreAllCardError(err, "RestoreAllCard", "FAILED_TO_RESTORE_ALL_CARDS", span, &status, zap.Error(err))
	}

	logSuccess("Successfully restored all cards", zap.Bool("success", true))

	return true, nil
}

// DeleteAllCardPermanent permanently deletes all card records from the database.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//
// Returns:
//   - A boolean indicating if the operation was successful, and an error if the operation fails.
func (s *cardCommandService) DeleteAllCardPermanent(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "DeleteAllCardPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.cardCommentRepository.DeleteAllCardPermanent(ctx)

	if err != nil {
		return s.errorhandler.HandleDeleteAllCardPermanentError(err, "DeleteAllCardPermanent", "FAILED_TO_DELETE_ALL_CARDS_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted all cards permanently", zap.Bool("success", true))

	return true, nil
}
