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
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantDocumentCommandService struct {
	kafka                             kafka.Kafka
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
	kafka kafka.Kafka,
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
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("CreateMerchantDocument", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "CreateMerchantDocument")
	defer span.End()

	span.SetAttributes(
		attribute.String("document-type", request.DocumentType),
	)

	s.logger.Debug("Creating new merchant document", zap.String("document-type", request.DocumentType))

	merchant, err := s.merchantQueryRepository.FindById(request.MerchantID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantDocumentResponse](s.logger, err, "CreateMerchantDocument", "FAILED_FIND_MERCHANT_BY_ID", span, &status, merchant_errors.ErrMerchantNotFoundRes, zap.Int("merchant_id", request.MerchantID))
	}

	user, err := s.userRepository.FindById(merchant.UserID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantDocumentResponse](s.logger, err, "CreateMerchantDocument", "FAILED_FIND_USER_BY_ID", span, &status, user_errors.ErrUserNotFoundRes, zap.Int("user_id", merchant.UserID))
	}

	merchantDocument, err := s.merchantDocumentCommandRepository.CreateMerchantDocument(request)

	if err != nil {
		return s.errorMerchantDocumentCommand.HandleCreateMerchantDocumentError(err, "CreateMerchantDocument", "FAILED_CREATE_MERCHANT_DOCUMENT", span, &status, zap.Int("merchant_id", request.MerchantID))
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
		return errorhandler.HandleMarshalError[*response.MerchantDocumentResponse](s.logger, err, "CreateMerchantDocument", "FAILED_MARSHAL_EMAIL_PAYLOAD", span, &status, merchant_errors.ErrFailedSendEmail, zap.Int("merchant_id", request.MerchantID))
	}

	err = s.kafka.SendMessage("email-service-topic-merchant-document-created", strconv.Itoa(merchantDocument.ID), payloadBytes)
	if err != nil {
		return errorhandler.HandleSendEmailError[*response.MerchantDocumentResponse](s.logger, err, "CreateMerchantDocument", "FAILED_SEND_EMAIL", span, &status, merchant_errors.ErrFailedSendEmail, zap.Int("merchant_id", request.MerchantID))
	}

	so := s.mapping.ToMerchantDocumentResponse(merchantDocument)

	s.logger.Debug("Successfully created merchant document", zap.String("document-type", request.DocumentType))

	return so, nil
}

func (s *merchantDocumentCommandService) UpdateMerchantDocument(request *requests.UpdateMerchantDocumentRequest) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("UpdateMerchantDocument", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateMerchantDocument")
	defer span.End()

	span.SetAttributes(
		attribute.Int("document_id", *request.DocumentID),
	)

	s.logger.Debug("Updating merchant document", zap.Int("document_id", *request.DocumentID))

	merchantDocument, err := s.merchantDocumentCommandRepository.UpdateMerchantDocument(request)
	if err != nil {
		return s.errorMerchantDocumentCommand.HandleUpdateMerchantDocumentError(err, "UpdateMerchantDocument", "FAILED_UPDATE_MERCHANT_DOCUMENT", span, &status, zap.Int("document_id", *request.DocumentID))
	}

	so := s.mapping.ToMerchantDocumentResponse(merchantDocument)

	s.mencache.DeleteCachedMerchantDocuments(merchantDocument.ID)

	s.logger.Debug("Successfully updated merchant document", zap.Int("document_id", *request.DocumentID))

	return so, nil
}

func (s *merchantDocumentCommandService) UpdateMerchantDocumentStatus(request *requests.UpdateMerchantDocumentStatusRequest) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("UpdateMerchantDocumentStatus", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateMerchantDocumentStatus")
	defer span.End()

	span.SetAttributes(
		attribute.Int("document_id", *request.DocumentID),
	)

	s.logger.Debug("Updating merchant document status", zap.Int("document_id", *request.DocumentID))

	merchant, err := s.merchantQueryRepository.FindById(request.MerchantID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantDocumentResponse](s.logger, err, "UpdateMerchantDocumentStatus", "FAILED_FIND_MERCHANT", span, &status, merchant_errors.ErrMerchantNotFoundRes, zap.Int("merchant_id", request.MerchantID))
	}

	user, err := s.userRepository.FindById(merchant.UserID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantDocumentResponse](s.logger, err, "UpdateMerchantDocumentStatus", "FAILED_FIND_USER", span, &status, user_errors.ErrUserNotFoundRes, zap.Int("user_id", merchant.UserID))
	}

	merchantDocument, err := s.merchantDocumentCommandRepository.UpdateMerchantDocumentStatus(request)
	if err != nil {
		return s.errorMerchantDocumentCommand.HandleUpdateMerchantDocumentStatusError(err, "UpdateMerchantDocumentStatus", "FAILED_UPDATE_MERCHANT_DOCUMENT_STATUS", span, &status, zap.Int("document_id", *request.DocumentID))
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
		return errorhandler.HandleMarshalError[*response.MerchantDocumentResponse](s.logger, err, "UpdateMerchantDocumentStatus", "FAILED_MARSHAL_EMAIL_PAYLOAD", span, &status, merchant_errors.ErrFailedSendEmail, zap.Int("merchant_id", request.MerchantID))
	}

	err = s.kafka.SendMessage("email-service-topic-merchant-document-update-status", strconv.Itoa(request.MerchantID), payloadBytes)
	if err != nil {
		return errorhandler.HandleSendEmailError[*response.MerchantDocumentResponse](s.logger, err, "UpdateMerchantDocumentStatus", "FAILED_SEND_EMAIL", span, &status, merchant_errors.ErrFailedSendEmail, zap.Int("merchant_id", request.MerchantID))
	}

	so := s.mapping.ToMerchantDocumentResponse(merchantDocument)

	s.mencache.DeleteCachedMerchantDocuments(merchantDocument.ID)

	s.logger.Debug("Successfully updated merchant document status", zap.Int("document_id", *request.DocumentID))

	return so, nil
}

func (s *merchantDocumentCommandService) TrashedMerchantDocument(documentID int) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("TrashedDocument", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "TrashedDocument")
	defer span.End()

	span.SetAttributes(attribute.Int("document_id", documentID))

	s.logger.Debug("Trashing merchant document", zap.Int("document_id", documentID))

	res, err := s.merchantDocumentCommandRepository.TrashedMerchantDocument(documentID)
	if err != nil {
		return s.errorMerchantDocumentCommand.HandleTrashedMerchantDocumentError(err, "TrashedDocument", "FAILED_TRASH_DOCUMENT", span, &status, zap.Int("document_id", documentID))
	}

	s.mencache.DeleteCachedMerchantDocuments(documentID)

	s.logger.Debug("Successfully trashed document", zap.Int("document_id", documentID))

	return s.mapping.ToMerchantDocumentResponse(res), nil
}

func (s *merchantDocumentCommandService) RestoreMerchantDocument(documentID int) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreDocument", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreDocument")
	defer span.End()

	span.SetAttributes(attribute.Int("document_id", documentID))

	s.logger.Debug("Restoring merchant document", zap.Int("document_id", documentID))

	res, err := s.merchantDocumentCommandRepository.RestoreMerchantDocument(documentID)
	if err != nil {
		return s.errorMerchantDocumentCommand.HandleRestoreMerchantDocumentError(err, "RestoreDocument", "FAILED_RESTORE_DOCUMENT", span, &status, zap.Int("document_id", documentID))
	}

	s.logger.Debug("Successfully restored document", zap.Int("document_id", documentID))

	return s.mapping.ToMerchantDocumentResponse(res), nil
}

func (s *merchantDocumentCommandService) DeleteMerchantDocumentPermanent(documentID int) (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteDocumentPermanent", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteDocumentPermanent")
	defer span.End()

	s.logger.Debug("Permanently deleting merchant document", zap.Int("document_id", documentID))

	_, err := s.merchantDocumentCommandRepository.DeleteMerchantDocumentPermanent(documentID)
	if err != nil {
		return s.errorMerchantDocumentCommand.HandleDeleteMerchantDocumentPermanentError(err, "DeleteDocumentPermanent", "FAILED_DELETE_DOCUMENT_PERMANENT", span, &status, zap.Int("document_id", documentID))
	}

	s.logger.Debug("Successfully deleted document permanently", zap.Int("document_id", documentID))

	return true, nil
}

func (s *merchantDocumentCommandService) RestoreAllMerchantDocument() (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreAllDocuments", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreAllDocuments")
	defer span.End()

	s.logger.Debug("Restoring all merchant documents")

	_, err := s.merchantDocumentCommandRepository.RestoreAllMerchantDocument()
	if err != nil {
		return s.errorMerchantDocumentCommand.HandleRestoreAllMerchantDocumentError(err, "RestoreAllDocuments", "FAILED_RESTORE_ALL_DOCUMENTS", span, &status)
	}

	s.logger.Debug("Successfully restored all merchant documents")

	return true, nil
}

func (s *merchantDocumentCommandService) DeleteAllMerchantDocumentPermanent() (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteAllDocumentsPermanent", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteAllDocumentsPermanent")
	defer span.End()

	s.logger.Debug("Deleting all merchant documents permanently")

	_, err := s.merchantDocumentCommandRepository.DeleteAllMerchantDocumentPermanent()
	if err != nil {
		return s.errorMerchantDocumentCommand.HandleDeleteAllMerchantDocumentPermanentError(err, "DeleteAllDocumentsPermanent", "FAILED_DELETE_ALL_DOCUMENTS_PERMANENT", span, &status)
	}

	s.logger.Debug("Successfully deleted all merchant documents permanently")

	return true, nil
}

func (s *merchantDocumentCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
