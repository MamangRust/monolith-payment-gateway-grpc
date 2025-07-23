package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-merchant/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-pkg/email"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/service"
	user_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors/service"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/merchant"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

// MerchantCommandServiceDeps contains the dependencies required to initialize a new instance
// of merchantCommandService.
type merchantCommandServiceDeps struct {
	// Kafka is the Kafka producer/consumer instance used to publish or consume merchant-related events.
	Kafka *kafka.Kafka

	// Ctx is the base context for controlling cancellations, timeouts, and tracing across the service.
	Ctx context.Context

	// ErrorHandler handles and wraps errors returned during merchant command operations.
	ErrorHandler errorhandler.MerchantCommandErrorHandler

	// Cache provides caching functionality for merchant command-related data.
	Cache mencache.MerchantCommandCache

	// UserRepository provides access to user data from the database.
	UserRepository repository.UserRepository

	// MerchantQueryRepository is used to fetch merchant data in a read-only manner.
	MerchantQueryRepository repository.MerchantQueryRepository

	// MerchantCommandRepository is responsible for creating, updating, and deleting merchant records.
	MerchantCommandRepository repository.MerchantCommandRepository

	// Logger provides structured logging functionality for observability and debugging.
	Logger logger.LoggerInterface

	// Mapper maps internal data to response formats.
	Mapper responseservice.MerchantCommandResponseMapper
}

// merchantCommandService provides an interface for interacting with the merchant command service,
// handling operations such as create, update, delete, and business logic for merchants.
type merchantCommandService struct {
	// ctx is the base context used for controlling cancellation, timeouts,
	// and carrying deadlines or trace metadata across service operations.
	ctx context.Context

	// kafka is the Kafka client used to publish merchant-related events
	// (e.g., merchant created, updated, or deleted) to Kafka topics.
	kafka *kafka.Kafka

	// errorHandler handles application-specific errors that occur in the merchant command service,
	// converting them into consistent responses or logs.
	errorHandler errorhandler.MerchantCommandErrorHandler

	// mencache provides caching functionality for merchant data to reduce repeated database access,
	// typically backed by Redis or in-memory cache.
	mencache mencache.MerchantCommandCache

	// userRepository provides access to user-related data required during merchant operations,
	// such as owner lookups or permission checks.
	userRepository repository.UserRepository

	// merchantQueryRepository is responsible for retrieving merchant data in a read-only manner,
	// often used to validate or enrich command operations.
	merchantQueryRepository repository.MerchantQueryRepository

	// merchantCommandRepository handles the actual persistence of merchant entities,
	// including creation, update, and soft/hard deletion in the database.
	merchantCommandRepository repository.MerchantCommandRepository

	// logger is the logging interface used to record structured logs
	// for observability and debugging during merchant command operations.
	logger logger.LoggerInterface

	// mapper maps internal domain models or records into response models or DTOs
	// suitable for gRPC/HTTP API layer.
	mapper responseservice.MerchantCommandResponseMapper

	observability observability.TraceLoggerObservability
}

// NewMerchantCommandService initializes a new instance of merchantCommandService with the provided parameters.
// It sets up Prometheus metrics for tracking request counts and durations and returns a configured
// merchantCommandService ready for handling merchant-related commands.
//
// Parameters:
// - params: A pointer to merchantCommandServiceDeps containing the necessary dependencies.
//
// Returns:
// - A pointer to an initialized merchantCommandService.
func NewMerchantCommandService(params *merchantCommandServiceDeps) MerchantCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_command_service_requests_total",
			Help: "Total number of requests to the MerchantCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_command_service_request_duration_seconds",
			Help:    "Histogram of request durations for the MerchantCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("merchant-command-service"), params.Logger, requestCounter, requestDuration)

	return &merchantCommandService{
		kafka:                     params.Kafka,
		ctx:                       params.Ctx,
		errorHandler:              params.ErrorHandler,
		mencache:                  params.Cache,
		merchantCommandRepository: params.MerchantCommandRepository,
		userRepository:            params.UserRepository,
		merchantQueryRepository:   params.MerchantQueryRepository,
		logger:                    params.Logger,
		mapper:                    params.Mapper,
		observability:             observability,
	}
}

// CreateMerchant creates a new merchant with the provided request data.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The merchant creation request payload.
//
// Returns:
//   - *response.MerchantResponse: The created merchant's data.
//   - *response.ErrorResponse: An error if creation fails.
func (s *merchantCommandService) CreateMerchant(ctx context.Context, request *requests.CreateMerchantRequest) (*response.MerchantResponse, *response.ErrorResponse) {
	const method = "CreateMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	user, err := s.userRepository.FindByUserId(ctx, request.UserID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantResponse](s.logger, err, method, "FAILED_FIND_USER_BY_ID", span, &status, user_errors.ErrUserNotFoundRes, zap.Error(err))
	}

	res, err := s.merchantCommandRepository.CreateMerchant(ctx, request)

	if err != nil {
		return s.errorHandler.HandleCreateMerchantError(err, method, "FAILED_CREATE_MERCHANT", span, &status, zap.Error(err))
	}

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   "Welcome to SanEdge Merchant Portal",
		"Message": "Your merchant account has been created successfully. To continue, please upload the required documents for verification. Once completed, our team will review and activate your account.",
		"Button":  "Upload Documents",
		"Link":    fmt.Sprintf("https://sanedge.example.com/merchant/%d/documents", user.ID),
	})

	emailPayload := map[string]any{
		"email":   user.Email,
		"subject": "Initial Verification - SanEdge",
		"body":    htmlBody,
	}

	payloadBytes, err := json.Marshal(emailPayload)
	if err != nil {
		return errorhandler.HandleMarshalError[*response.MerchantResponse](s.logger, err, method, "FAILED_MARSHAL_EMAIL_PAYLOAD", span, &status, merchant_errors.ErrFailedSendEmail, zap.Error(err))
	}

	err = s.kafka.SendMessage("email-service-topic-merchant-created", strconv.Itoa(res.ID), payloadBytes)
	if err != nil {
		return errorhandler.HandleSendEmailError[*response.MerchantResponse](s.logger, err, method, "FAILED_SEND_EMAIL", span, &status, merchant_errors.ErrFailedSendEmail, zap.Error(err))
	}

	so := s.mapper.ToMerchantResponse(res)

	logSuccess("Successfully created merchant", zap.Bool("success", true))

	return so, nil
}

// UpdateMerchant updates an existing merchant's data.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The request containing updated merchant data.
//
// Returns:
//   - *response.MerchantResponse: The updated merchant's data.
//   - *response.ErrorResponse: An error if update fails.
func (s *merchantCommandService) UpdateMerchant(ctx context.Context, request *requests.UpdateMerchantRequest) (*response.MerchantResponse, *response.ErrorResponse) {
	const method = "UpdateMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.merchantCommandRepository.UpdateMerchant(ctx, request)

	if err != nil {
		return s.errorHandler.HandleUpdateMerchantError(err, method, "FAILED_UPDATE_MERCHANT", span, &status, zap.Error(err))
	}

	so := s.mapper.ToMerchantResponse(res)

	s.mencache.DeleteCachedMerchant(ctx, res.ID)

	logSuccess("Successfully updated merchant", zap.Bool("success", true))

	return so, nil
}

// UpdateMerchantStatus updates the status (active/inactive) of a merchant.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The request containing status update information.
//
// Returns:
//   - *response.MerchantResponse: The updated merchant data.
//   - *response.ErrorResponse: An error if status update fails.
func (s *merchantCommandService) UpdateMerchantStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*response.MerchantResponse, *response.ErrorResponse) {
	const method = "UpdateMerchantStatus"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantQueryRepository.FindById(ctx, *request.MerchantID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantResponse](s.logger, err, method, "FAILED_FIND_MERCHANT_BY_ID", span, &status, merchant_errors.ErrFailedFindMerchantById, zap.Error(err))
	}

	user, err := s.userRepository.FindByUserId(ctx, merchant.UserID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantResponse](s.logger, err, method, "FAILED_FIND_USER_BY_ID", span, &status, user_errors.ErrUserNotFoundRes, zap.Error(err))
	}

	res, err := s.merchantCommandRepository.UpdateMerchantStatus(ctx, request)

	if err != nil {
		return s.errorHandler.HandleUpdateMerchantStatusError(err, method, "FAILED_UPDATE_MERCHANT_STATUS", span, &status, zap.Error(err))
	}

	statusReq := request.Status
	subject := ""
	message := ""
	buttonLabel := "Go to Portal"
	link := fmt.Sprintf("https://sanedge.example.com/merchant/%d/dashboard", *request.MerchantID)

	switch statusReq {
	case "active":
		subject = "Your Merchant Account is Now Active"
		message = "Congratulations! Your merchant account has been verified and is now <b>active</b>. You can now fully access all features in the SanEdge Merchant Portal."
	case "inactive":
		subject = "Merchant Account Set to Inactive"
		message = "Your merchant account status has been set to <b>inactive</b>. Please contact support if you believe this is a mistake."
	case "rejected":
		subject = "Merchant Account Rejected"
		message = "We're sorry to inform you that your merchant account has been <b>rejected</b>. Please contact support or review your submissions."
	default:
		return nil, nil
	}

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   subject,
		"Message": message,
		"Button":  buttonLabel,
		"Link":    link,
	})

	emailPayload := map[string]any{
		"email":   user.Email,
		"subject": subject,
		"body":    htmlBody,
	}

	payloadBytes, err := json.Marshal(emailPayload)
	if err != nil {
		return errorhandler.HandleMarshalError[*response.MerchantResponse](s.logger, err, method, "FAILED_MARSHAL_EMAIL_PAYLOAD", span, &status, merchant_errors.ErrFailedSendEmail, zap.Error(err))
	}

	err = s.kafka.SendMessage("email-service-topic-merchant-update-status", strconv.Itoa(*request.MerchantID), payloadBytes)
	if err != nil {
		return errorhandler.HandleSendEmailError[*response.MerchantResponse](s.logger, err, method, "FAILED_SEND_EMAIL", span, &status, merchant_errors.ErrFailedSendEmail, zap.Error(err))
	}

	so := s.mapper.ToMerchantResponse(res)
	s.mencache.DeleteCachedMerchant(ctx, res.ID)

	logSuccess("Successfully updated merchant status", zap.Bool("success", true))

	return so, nil
}

// TrashedMerchant soft-deletes a merchant by its ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - merchant_id: The ID of the merchant to be soft-deleted.
//
// Returns:
//   - *response.MerchantResponse: The trashed merchant data.
//   - *response.ErrorResponse: An error if the operation fails.
func (s *merchantCommandService) TrashedMerchant(ctx context.Context, merchant_id int) (*response.MerchantResponseDeleteAt, *response.ErrorResponse) {
	const method = "TrashedMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.merchantCommandRepository.TrashedMerchant(ctx, merchant_id)

	if err != nil {
		return s.errorHandler.HandleTrashedMerchantError(err, method, "FAILED_TRASHED_MERCHANT", span, &status, zap.Error(err))
	}
	so := s.mapper.ToMerchantResponseDeleteAt(res)

	s.mencache.DeleteCachedMerchant(ctx, res.ID)

	logSuccess("Successfully trashed merchant", zap.Bool("success", true))

	return so, nil
}

// RestoreMerchant restores a soft-deleted merchant by its ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - merchant_id: The ID of the merchant to restore.
//
// Returns:
//   - *response.MerchantResponse: The restored merchant data.
//   - *response.ErrorResponse: An error if restoration fails.
func (s *merchantCommandService) RestoreMerchant(ctx context.Context, merchant_id int) (*response.MerchantResponse, *response.ErrorResponse) {
	const method = "RestoreMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.merchantCommandRepository.RestoreMerchant(ctx, merchant_id)

	if err != nil {
		return s.errorHandler.HandleRestoreMerchantError(err, method, "FAILED_RESTORE_MERCHANT", span, &status, zap.Error(err))
	}

	so := s.mapper.ToMerchantResponse(res)

	s.mencache.DeleteCachedMerchant(ctx, res.ID)

	logSuccess("Successfully restored merchant", zap.Bool("success", true))

	return so, nil
}

// DeleteMerchantPermanent permanently deletes a merchant by its ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - merchant_id: The ID of the merchant to delete permanently.
//
// Returns:
//   - bool: Whether the deletion was successful.
//   - *response.ErrorResponse: An error if deletion fails.
func (s *merchantCommandService) DeleteMerchantPermanent(ctx context.Context, merchant_id int) (bool, *response.ErrorResponse) {
	const method = "DeleteMerchantPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.merchantCommandRepository.DeleteMerchantPermanent(ctx, merchant_id)

	if err != nil {
		return s.errorHandler.HandleDeleteMerchantPermanentError(err, method, "FAILED_DELETE_MERCHANT_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted merchant permanently", zap.Bool("success", true))

	return true, nil
}

// RestoreAllMerchant restores all soft-deleted merchants.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//
// Returns:
//   - bool: Whether the restoration was successful.
//   - *response.ErrorResponse: An error if restoration fails.
func (s *merchantCommandService) RestoreAllMerchant(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "RestoreAllMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.merchantCommandRepository.RestoreAllMerchant(ctx)

	if err != nil {
		return s.errorHandler.HandleRestoreAllMerchantError(err, method, "FAILED_RESTORE_ALL_MERCHANT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully restored all merchants", zap.Bool("success", true))

	return true, nil
}

// DeleteAllMerchantPermanent permanently deletes all soft-deleted merchants.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//
// Returns:
//   - bool: Whether the deletion was successful.
//   - *response.ErrorResponse: An error if deletion fails.
func (s *merchantCommandService) DeleteAllMerchantPermanent(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "DeleteAllMerchantPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.merchantCommandRepository.DeleteAllMerchantPermanent(ctx)

	if err != nil {
		return s.errorHandler.HandleDeleteAllMerchantPermanentError(err, method, "FAILED_DELETE_ALL_MERCHANT_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted all merchants permanently", zap.Bool("success", true))

	return true, nil
}
