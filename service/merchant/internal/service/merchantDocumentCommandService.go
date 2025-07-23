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
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/merchantdocument"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

// merchantDocumentCommandDeps contains the dependencies required to
// construct a merchantDocumentCommandService.
type merchantDocumentCommandDeps struct {
	// Kafka is the Kafka client used to produce events.
	Kafka *kafka.Kafka

	// Ctx is the base context used in the service.
	Ctx context.Context

	// Cache is the cache layer for merchant document commands.
	Cache mencache.MerchantDocumentCommandCache

	// ErrorHandler handles errors for merchant document commands.
	ErrorHandler errorhandler.MerchantDocumentCommandErrorHandler

	// CommandRepository handles write operations on merchant documents.
	CommandRepository repository.MerchantDocumentCommandRepository

	// MerchantQueryRepository provides access to merchant query data.
	MerchantQueryRepository repository.MerchantQueryRepository

	// UserRepository provides access to user-related data.
	UserRepository repository.UserRepository

	// Logger is used for structured logging.
	Logger logger.LoggerInterface

	// Mapper maps internal data to response formats.
	Mapper responseservice.MerchantDocumentCommandResponseMapper
}

// merchantDocumentCommandService handles merchant document command operations,
// including creation, update, deletion, and restoration.
type merchantDocumentCommandService struct {
	// kafka is the Kafka client used to produce events.
	kafka *kafka.Kafka

	// ctx is the base context used throughout the service.
	ctx context.Context

	// mencache is the cache layer for merchant document commands.
	mencache mencache.MerchantDocumentCommandCache

	// errorMerchantDocumentCommand handles errors for merchant document commands.
	errorMerchantDocumentCommand errorhandler.MerchantDocumentCommandErrorHandler

	// merchantQueryRepository provides access to merchant query data.
	merchantQueryRepository repository.MerchantQueryRepository

	// merchantDocumentCommandRepository handles write operations on merchant documents.
	merchantDocumentCommandRepository repository.MerchantDocumentCommandRepository

	// userRepository provides access to user-related data.
	userRepository repository.UserRepository

	// logger is used for logging within the service.
	logger logger.LoggerInterface

	// mapper maps internal data to response formats.
	mapper responseservice.MerchantDocumentCommandResponseMapper

	observability observability.TraceLoggerObservability
}

// NewMerchantDocumentCommandService initializes a new instance of merchantDocumentCommandService with the provided parameters.
// It sets up Prometheus metrics for tracking request counts and durations and returns a configured
// merchantDocumentCommandService ready for handling merchant document-related commands.
//
// Parameters:
// - params: A pointer to merchantDocumentCommandDeps containing the necessary dependencies.
//
// Returns:
// - A pointer to an initialized merchantDocumentCommandService.
func NewMerchantDocumentCommandService(
	params *merchantDocumentCommandDeps,
) MerchantDocumentCommandService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "merchant_document_command_request_count",
		Help: "Number of merchant document command requests MerchantDocumentCommandService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "merchant_document_command_request_duration_seconds",
		Help:    "The duration of requests MerchantDocumentCommandService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("merchant-document-command-service"), params.Logger, requestCounter, requestDuration)

	return &merchantDocumentCommandService{
		kafka:                             params.Kafka,
		ctx:                               params.Ctx,
		mencache:                          params.Cache,
		errorMerchantDocumentCommand:      params.ErrorHandler,
		merchantQueryRepository:           params.MerchantQueryRepository,
		merchantDocumentCommandRepository: params.CommandRepository,
		userRepository:                    params.UserRepository,
		logger:                            params.Logger,
		mapper:                            params.Mapper,
		observability:                     observability,
	}
}

// CreateMerchantDocument creates a new merchant document.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The merchant document creation request payload.
//
// Returns:
//   - *response.MerchantDocumentResponse: The created merchant document.
//   - *response.ErrorResponse: An error if creation fails.
func (s *merchantDocumentCommandService) CreateMerchantDocument(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	const method = "CreateMerchantDocument"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantQueryRepository.FindById(ctx, request.MerchantID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantDocumentResponse](s.logger, err, method, "FAILED_FIND_MERCHANT_BY_ID", span, &status, merchant_errors.ErrMerchantNotFoundRes, zap.Error(err))
	}

	user, err := s.userRepository.FindByUserId(ctx, merchant.UserID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantDocumentResponse](s.logger, err, method, "FAILED_FIND_USER_BY_ID", span, &status, user_errors.ErrUserNotFoundRes, zap.Error(err))
	}

	merchantDocument, err := s.merchantDocumentCommandRepository.CreateMerchantDocument(ctx, request)

	if err != nil {
		return s.errorMerchantDocumentCommand.HandleCreateMerchantDocumentError(err, method, "FAILED_CREATE_MERCHANT_DOCUMENT", span, &status, zap.Error(err))
	}

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   "Welcome to SanEdge Merchant Portal",
		"Message": "Thank you for registering your merchant account. Your account is currently <b>inactive</b> and under initial review. To proceed, please upload all required documents for verification. Once your documents are submitted, our team will review them and activate your account accordingly.",
		"Button":  "Upload Documents",
		"Link":    fmt.Sprintf("https://sanedge.example.com/merchant/%d/documents", user.ID),
	})

	emailPayload := map[string]any{
		"email":   user.Email,
		"subject": "Merchant Verification Pending - Action Required",
		"body":    htmlBody,
	}

	payloadBytes, err := json.Marshal(emailPayload)
	if err != nil {
		return errorhandler.HandleMarshalError[*response.MerchantDocumentResponse](s.logger, err, method, "FAILED_MARSHAL_EMAIL_PAYLOAD", span, &status, merchant_errors.ErrFailedSendEmail, zap.Error(err))
	}

	err = s.kafka.SendMessage("email-service-topic-merchant-document-created", strconv.Itoa(merchantDocument.ID), payloadBytes)
	if err != nil {
		return errorhandler.HandleSendEmailError[*response.MerchantDocumentResponse](s.logger, err, method, "FAILED_SEND_EMAIL", span, &status, merchant_errors.ErrFailedSendEmail, zap.Error(err))
	}

	so := s.mapper.ToMerchantDocumentResponse(merchantDocument)

	logSuccess("Successfully created merchant document", zap.Bool("success", true))

	return so, nil
}

// UpdateMerchantDocument updates an existing merchant document.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The update request containing new document data.
//
// Returns:
//   - *response.MerchantDocumentResponse: The updated merchant document.
//   - *response.ErrorResponse: An error if update fails.
func (s *merchantDocumentCommandService) UpdateMerchantDocument(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	const method = "UpdateMerchantDocument"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	merchantDocument, err := s.merchantDocumentCommandRepository.UpdateMerchantDocument(ctx, request)
	if err != nil {
		return s.errorMerchantDocumentCommand.HandleUpdateMerchantDocumentError(err, method, "FAILED_UPDATE_MERCHANT_DOCUMENT", span, &status, zap.Error(err))
	}

	so := s.mapper.ToMerchantDocumentResponse(merchantDocument)

	s.mencache.DeleteCachedMerchantDocuments(ctx, merchantDocument.ID)

	logSuccess("Successfully updated merchant document", zap.Bool("success", true))

	return so, nil
}

// UpdateMerchantDocumentStatus updates the status (e.g., verified, rejected) of a merchant document.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The request containing status update data.
//
// Returns:
//   - *response.MerchantDocumentResponse: The updated merchant document.
//   - *response.ErrorResponse: An error if status update fails.
func (s *merchantDocumentCommandService) UpdateMerchantDocumentStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	const method = "UpdateMerchantDocumentStatus"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantQueryRepository.FindById(ctx, request.MerchantID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantDocumentResponse](s.logger, err, method, "FAILED_FIND_MERCHANT", span, &status, merchant_errors.ErrMerchantNotFoundRes, zap.Error(err))
	}

	user, err := s.userRepository.FindByUserId(ctx, merchant.UserID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantDocumentResponse](s.logger, err, method, "FAILED_FIND_USER", span, &status, user_errors.ErrUserNotFoundRes, zap.Error(err))
	}

	merchantDocument, err := s.merchantDocumentCommandRepository.UpdateMerchantDocumentStatus(ctx, request)
	if err != nil {
		return s.errorMerchantDocumentCommand.HandleUpdateMerchantDocumentStatusError(err, method, "FAILED_UPDATE_MERCHANT_DOCUMENT_STATUS", span, &status, zap.Error(err))
	}

	statusReq := request.Status
	note := request.Note
	subject := ""
	message := ""
	buttonLabel := ""
	link := fmt.Sprintf("https://sanedge.example.com/merchant/%d/documents", request.MerchantID)

	switch statusReq {
	case "pending":
		subject = "Merchant Document Status: Pending Review"
		message = "Your merchant documents have been submitted and are currently pending review."
		buttonLabel = "View Documents"
	case "approved":
		subject = "Merchant Document Status: Approved"
		message = "Congratulations! Your merchant documents have been approved. Your account is now active and fully functional."
		buttonLabel = "Go to Dashboard"
		link = fmt.Sprintf("https://sanedge.example.com/merchant/%d/dashboard", request.MerchantID)
	case "rejected":
		subject = "Merchant Document Status: Rejected"
		message = "Unfortunately, your merchant documents were rejected. Please review the feedback below and re-upload the necessary documents."
		buttonLabel = "Re-upload Documents"
	default:
		return nil, nil
	}

	if note != "" {
		message += fmt.Sprintf(`<br><br><b>Reviewer Note:</b><br><i>%s</i>`, note)
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
		return errorhandler.HandleMarshalError[*response.MerchantDocumentResponse](s.logger, err, method, "FAILED_MARSHAL_EMAIL_PAYLOAD", span, &status, merchant_errors.ErrFailedSendEmail, zap.Error(err))
	}

	err = s.kafka.SendMessage("email-service-topic-merchant-document-update-status", strconv.Itoa(request.MerchantID), payloadBytes)
	if err != nil {
		return errorhandler.HandleSendEmailError[*response.MerchantDocumentResponse](s.logger, err, method, "FAILED_SEND_EMAIL", span, &status, merchant_errors.ErrFailedSendEmail, zap.Error(err))
	}

	so := s.mapper.ToMerchantDocumentResponse(merchantDocument)

	s.mencache.DeleteCachedMerchantDocuments(ctx, merchantDocument.ID)

	logSuccess("Successfully updated merchant document status", zap.Bool("success", true))

	return so, nil
}

// TrashedMerchantDocument soft-deletes a merchant document by its ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - document_id: The ID of the document to be soft-deleted.
//
// Returns:
//   - *response.MerchantDocumentResponse: The trashed document.
//   - *response.ErrorResponse: An error if the operation fails.
func (s *merchantDocumentCommandService) TrashedMerchantDocument(ctx context.Context, documentID int) (*response.MerchantDocumentResponseDeleteAt, *response.ErrorResponse) {
	const method = "TrashedMerchantDocument"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.merchantDocumentCommandRepository.TrashedMerchantDocument(ctx, documentID)
	if err != nil {
		return s.errorMerchantDocumentCommand.HandleTrashedMerchantDocumentError(err, method, "FAILED_TRASH_DOCUMENT", span, &status, zap.Error(err))
	}

	s.mencache.DeleteCachedMerchantDocuments(ctx, documentID)

	so := s.mapper.ToMerchantDocumentResponseDeleteAt(res)

	logSuccess("Successfully trashed document", zap.Bool("success", true))

	return so, nil
}

// RestoreMerchantDocument restores a soft-deleted merchant document by its ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - document_id: The ID of the document to restore.
//
// Returns:
//   - *response.MerchantDocumentResponse: The restored document.
//   - *response.ErrorResponse: An error if restoration fails.
func (s *merchantDocumentCommandService) RestoreMerchantDocument(ctx context.Context, documentID int) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	const method = "RestoreMerchantDocument"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.merchantDocumentCommandRepository.RestoreMerchantDocument(ctx, documentID)
	if err != nil {
		return s.errorMerchantDocumentCommand.HandleRestoreMerchantDocumentError(err, method, "FAILED_RESTORE_DOCUMENT", span, &status, zap.Int("document_id", documentID))
	}

	so := s.mapper.ToMerchantDocumentResponse(res)

	logSuccess("Successfully restored document", zap.Bool("success", true))

	return so, nil
}

// DeleteMerchantDocumentPermanent permanently deletes a merchant document by its ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - document_id: The ID of the document to delete.
//
// Returns:
//   - bool: True if the deletion was successful.
//   - *response.ErrorResponse: An error if the deletion fails.
func (s *merchantDocumentCommandService) DeleteMerchantDocumentPermanent(ctx context.Context, documentID int) (bool, *response.ErrorResponse) {
	const method = "DeleteMerchantDocumentPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.merchantDocumentCommandRepository.DeleteMerchantDocumentPermanent(ctx, documentID)
	if err != nil {
		return s.errorMerchantDocumentCommand.HandleDeleteMerchantDocumentPermanentError(err, method, "FAILED_DELETE_DOCUMENT_PERMANENT", span, &status, zap.Int("document_id", documentID))
	}

	logSuccess("Successfully deleted document permanently", zap.Bool("success", true))

	return true, nil
}

// RestoreAllMerchantDocument restores all soft-deleted merchant documents.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//
// Returns:
//   - bool: True if all documents were restored successfully.
//   - *response.ErrorResponse: An error if restoration fails.
func (s *merchantDocumentCommandService) RestoreAllMerchantDocument(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "RestoreAllMerchantDocument"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.merchantDocumentCommandRepository.RestoreAllMerchantDocument(ctx)
	if err != nil {
		return s.errorMerchantDocumentCommand.HandleRestoreAllMerchantDocumentError(err, method, "FAILED_RESTORE_ALL_DOCUMENTS", span, &status, zap.Error(err))
	}

	logSuccess("Successfully restored all documents", zap.Bool("success", true))

	return true, nil
}

// DeleteAllMerchantDocumentPermanent permanently deletes all soft-deleted merchant documents.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//
// Returns:
//   - bool: True if all documents were deleted successfully.
//   - *response.ErrorResponse: An error if deletion fails.
func (s *merchantDocumentCommandService) DeleteAllMerchantDocumentPermanent(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "DeleteAllMerchantDocumentPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.merchantDocumentCommandRepository.DeleteAllMerchantDocumentPermanent(ctx)

	if err != nil {
		return s.errorMerchantDocumentCommand.HandleDeleteAllMerchantDocumentPermanentError(err, method, "FAILED_DELETE_ALL_DOCUMENTS_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted all documents", zap.Bool("success", true))

	return true, nil
}
