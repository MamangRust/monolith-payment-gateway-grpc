package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-merchant/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-pkg/email"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantDocumentCommandService struct {
	kafka                             *kafka.Kafka
	ctx                               context.Context
	mencache                          mencache.MerchantDocumentCommandCache
	errorMerchantDocumentCommand      errorhandler.MerchantDocumentCommandErrorHandler
	trace                             trace.Tracer
	merchantQueryRepository           repository.MerchantQueryRepository
	merchantDocumentCommandRepository repository.MerchantDocumentCommandRepository
	userRepository                    repository.UserRepository
	logger                            logger.LoggerInterface
	mapping                           responseservice.MerchantDocumentResponseMapper
	requestCounter                    *prometheus.CounterVec
	requestDuration                   *prometheus.HistogramVec
}

func NewMerchantDocumentCommandService(
	kafka *kafka.Kafka,
	ctx context.Context,
	mencache mencache.MerchantDocumentCommandCache,
	errorMerchantDocumentCommand errorhandler.MerchantDocumentCommandErrorHandler,
	merchantDocumentCommandRepository repository.MerchantDocumentCommandRepository,
	merchantQueryRepository repository.MerchantQueryRepository,
	userRepository repository.UserRepository,
	logger logger.LoggerInterface,
	mapping responseservice.MerchantDocumentResponseMapper,
) *merchantDocumentCommandService {
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

	return &merchantDocumentCommandService{
		kafka:                             kafka,
		ctx:                               ctx,
		mencache:                          mencache,
		trace:                             otel.Tracer("merchant-document-command-service"),
		errorMerchantDocumentCommand:      errorMerchantDocumentCommand,
		merchantQueryRepository:           merchantQueryRepository,
		merchantDocumentCommandRepository: merchantDocumentCommandRepository,
		userRepository:                    userRepository,
		logger:                            logger,
		mapping:                           mapping,
		requestCounter:                    requestCounter,
		requestDuration:                   requestDuration,
	}
}

func (s *merchantDocumentCommandService) CreateMerchantDocument(request *requests.CreateMerchantDocumentRequest) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	const method = "CreateMerchantDocument"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantQueryRepository.FindById(request.MerchantID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantDocumentResponse](s.logger, err, method, "FAILED_FIND_MERCHANT_BY_ID", span, &status, merchant_errors.ErrMerchantNotFoundRes, zap.Error(err))
	}

	user, err := s.userRepository.FindById(merchant.UserID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantDocumentResponse](s.logger, err, method, "FAILED_FIND_USER_BY_ID", span, &status, user_errors.ErrUserNotFoundRes, zap.Error(err))
	}

	merchantDocument, err := s.merchantDocumentCommandRepository.CreateMerchantDocument(request)

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

	so := s.mapping.ToMerchantDocumentResponse(merchantDocument)

	logSuccess("Successfully created merchant document", zap.Bool("success", true))

	return so, nil
}

func (s *merchantDocumentCommandService) UpdateMerchantDocument(request *requests.UpdateMerchantDocumentRequest) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	const method = "UpdateMerchantDocument"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	merchantDocument, err := s.merchantDocumentCommandRepository.UpdateMerchantDocument(request)
	if err != nil {
		return s.errorMerchantDocumentCommand.HandleUpdateMerchantDocumentError(err, method, "FAILED_UPDATE_MERCHANT_DOCUMENT", span, &status, zap.Error(err))
	}

	so := s.mapping.ToMerchantDocumentResponse(merchantDocument)

	s.mencache.DeleteCachedMerchantDocuments(merchantDocument.ID)

	logSuccess("Successfully updated merchant document", zap.Bool("success", true))

	return so, nil
}

func (s *merchantDocumentCommandService) UpdateMerchantDocumentStatus(request *requests.UpdateMerchantDocumentStatusRequest) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	const method = "UpdateMerchantDocumentStatus"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantQueryRepository.FindById(request.MerchantID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantDocumentResponse](s.logger, err, method, "FAILED_FIND_MERCHANT", span, &status, merchant_errors.ErrMerchantNotFoundRes, zap.Error(err))
	}

	user, err := s.userRepository.FindById(merchant.UserID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantDocumentResponse](s.logger, err, method, "FAILED_FIND_USER", span, &status, user_errors.ErrUserNotFoundRes, zap.Error(err))
	}

	merchantDocument, err := s.merchantDocumentCommandRepository.UpdateMerchantDocumentStatus(request)
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

	so := s.mapping.ToMerchantDocumentResponse(merchantDocument)

	s.mencache.DeleteCachedMerchantDocuments(merchantDocument.ID)

	logSuccess("Successfully updated merchant document status", zap.Bool("success", true))

	return so, nil
}

func (s *merchantDocumentCommandService) TrashedMerchantDocument(documentID int) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	const method = "TrashedMerchantDocument"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	res, err := s.merchantDocumentCommandRepository.TrashedMerchantDocument(documentID)
	if err != nil {
		return s.errorMerchantDocumentCommand.HandleTrashedMerchantDocumentError(err, method, "FAILED_TRASH_DOCUMENT", span, &status, zap.Error(err))
	}

	s.mencache.DeleteCachedMerchantDocuments(documentID)

	so := s.mapping.ToMerchantDocumentResponse(res)

	logSuccess("Successfully trashed document", zap.Bool("success", true))

	return so, nil
}

func (s *merchantDocumentCommandService) RestoreMerchantDocument(documentID int) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	const method = "RestoreMerchantDocument"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	res, err := s.merchantDocumentCommandRepository.RestoreMerchantDocument(documentID)
	if err != nil {
		return s.errorMerchantDocumentCommand.HandleRestoreMerchantDocumentError(err, method, "FAILED_RESTORE_DOCUMENT", span, &status, zap.Int("document_id", documentID))
	}

	so := s.mapping.ToMerchantDocumentResponse(res)

	logSuccess("Successfully restored document", zap.Bool("success", true))

	return so, nil
}

func (s *merchantDocumentCommandService) DeleteMerchantDocumentPermanent(documentID int) (bool, *response.ErrorResponse) {
	const method = "DeleteMerchantDocumentPermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	_, err := s.merchantDocumentCommandRepository.DeleteMerchantDocumentPermanent(documentID)
	if err != nil {
		return s.errorMerchantDocumentCommand.HandleDeleteMerchantDocumentPermanentError(err, method, "FAILED_DELETE_DOCUMENT_PERMANENT", span, &status, zap.Int("document_id", documentID))
	}

	logSuccess("Successfully deleted document permanently", zap.Bool("success", true))

	return true, nil
}

func (s *merchantDocumentCommandService) RestoreAllMerchantDocument() (bool, *response.ErrorResponse) {
	const method = "RestoreAllMerchantDocument"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	_, err := s.merchantDocumentCommandRepository.RestoreAllMerchantDocument()
	if err != nil {
		return s.errorMerchantDocumentCommand.HandleRestoreAllMerchantDocumentError(err, method, "FAILED_RESTORE_ALL_DOCUMENTS", span, &status, zap.Error(err))
	}

	logSuccess("Successfully restored all documents", zap.Bool("success", true))

	return true, nil
}

func (s *merchantDocumentCommandService) DeleteAllMerchantDocumentPermanent() (bool, *response.ErrorResponse) {
	const method = "DeleteAllMerchantDocumentPermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	_, err := s.merchantDocumentCommandRepository.DeleteAllMerchantDocumentPermanent()

	if err != nil {
		return s.errorMerchantDocumentCommand.HandleDeleteAllMerchantDocumentPermanentError(err, method, "FAILED_DELETE_ALL_DOCUMENTS_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted all documents", zap.Bool("success", true))

	return true, nil
}

func (s *merchantDocumentCommandService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *merchantDocumentCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
