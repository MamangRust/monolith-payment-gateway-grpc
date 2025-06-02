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
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantCommandService struct {
	ctx                       context.Context
	redis                     *redis.Client
	kafka                     kafka.Kafka
	trace                     trace.Tracer
	errorHandler              errorhandler.MerchantCommandErrorHandler
	userRepository            repository.UserRepository
	merchantQueryRepository   repository.MerchantQueryRepository
	merchantCommandRepository repository.MerchantCommandRepository
	logger                    logger.LoggerInterface
	mapping                   responseservice.MerchantResponseMapper
	mencache                  mencache.MerchantCommandCache
	requestCounter            *prometheus.CounterVec
	requestDuration           *prometheus.HistogramVec
}

func NewMerchantCommandService(kafka kafka.Kafka, ctx context.Context,
	errorHandler errorhandler.MerchantCommandErrorHandler,
	mencache mencache.MerchantCommandCache,
	userRepository repository.UserRepository,
	merchantQueryRepository repository.MerchantQueryRepository,
	merchantCommandRepository repository.MerchantCommandRepository, logger logger.LoggerInterface, mapping responseservice.MerchantResponseMapper) *merchantCommandService {
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

	return &merchantCommandService{
		kafka:                     kafka,
		ctx:                       ctx,
		errorHandler:              errorHandler,
		mencache:                  mencache,
		trace:                     otel.Tracer("merchant-command-service"),
		merchantCommandRepository: merchantCommandRepository,
		userRepository:            userRepository,
		merchantQueryRepository:   merchantQueryRepository,
		logger:                    logger,
		mapping:                   mapping,
		requestCounter:            requestCounter,
		requestDuration:           requestDuration,
	}
}

func (s *merchantCommandService) CreateMerchant(request *requests.CreateMerchantRequest) (*response.MerchantResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("CreateMerchant", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "CreateMerchant")
	defer span.End()

	span.SetAttributes(
		attribute.String("name", request.Name),
	)

	s.logger.Debug("Creating new merchant", zap.String("merchant_name", request.Name))

	user, err := s.userRepository.FindById(request.UserID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantResponse](s.logger, err, "CreateMerchant", "FAILED_FIND_USER_BY_ID", span, &status, user_errors.ErrUserNotFoundRes, zap.Int("user_id", request.UserID))
	}

	res, err := s.merchantCommandRepository.CreateMerchant(request)

	if err != nil {
		return s.errorHandler.HandleCreateMerchantError(err, "CreateMerchant", "FAILED_CREATE_MERCHANT", span, &status, zap.Int("user_id", request.UserID))
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
		return errorhandler.HandleMarshalError[*response.MerchantResponse](s.logger, err, "CreateMerchant", "FAILED_MARSHAL_EMAIL_PAYLOAD", span, &status, merchant_errors.ErrFailedSendEmail, zap.Int("user_id", user.ID))
	}

	err = s.kafka.SendMessage("email-service-topic-merchant-created", strconv.Itoa(res.ID), payloadBytes)
	if err != nil {
		return errorhandler.HandleSendEmailError[*response.MerchantResponse](s.logger, err, "CreateMerchant", "FAILED_SEND_EMAIL", span, &status, merchant_errors.ErrFailedSendEmail, zap.Int("user_id", user.ID))
	}

	so := s.mapping.ToMerchantResponse(res)

	s.logger.Debug("Successfully created merchant", zap.Int("merchant_id", res.ID))

	return so, nil
}

func (s *merchantCommandService) UpdateMerchant(request *requests.UpdateMerchantRequest) (*response.MerchantResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("UpdateMerchant", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateMerchant")
	defer span.End()

	span.SetAttributes(
		attribute.Int("merchant_id", *request.MerchantID),
	)

	s.logger.Debug("Updating merchant", zap.Int("merchant_id", *request.MerchantID))

	res, err := s.merchantCommandRepository.UpdateMerchant(request)

	if err != nil {
		return s.errorHandler.HandleUpdateMerchantError(err, "UpdateMerchant", "FAILED_UPDATE_MERCHANT", span, &status, zap.Int("merchant_id", *request.MerchantID))
	}

	so := s.mapping.ToMerchantResponse(res)

	cacheKey := fmt.Sprintf("merchant:id:%d", *request.MerchantID)
	if err := s.redis.Del(s.ctx, cacheKey).Err(); err != nil {
		s.logger.Error("Failed to delete merchant from cache", zap.Error(err))
	}

	s.logger.Debug("Successfully updated merchant", zap.Int("merchant_id", res.ID))

	return so, nil
}

func (s *merchantCommandService) UpdateMerchantStatus(request *requests.UpdateMerchantStatusRequest) (*response.MerchantResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("UpdateMerchantStatus", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateMerchantStatus")
	defer span.End()

	span.SetAttributes(
		attribute.Int("merchant_id", *request.MerchantID),
	)

	s.logger.Debug("Updating merchant status", zap.Int("merchant_id", *request.MerchantID))

	merchant, err := s.merchantQueryRepository.FindById(*request.MerchantID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantResponse](s.logger, err, "UpdateMerchantStatus", "FAILED_FIND_MERCHANT_BY_ID", span, &status, merchant_errors.ErrFailedFindMerchantById, zap.Int("merchant_id", *request.MerchantID))
	}

	user, err := s.userRepository.FindById(merchant.UserID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantResponse](s.logger, err, "UpdateMerchantStatus", "FAILED_FIND_USER_BY_ID", span, &status, user_errors.ErrUserNotFoundRes, zap.Int("user_id", merchant.UserID))
	}

	res, err := s.merchantCommandRepository.UpdateMerchantStatus(request)

	if err != nil {
		return s.errorHandler.HandleUpdateMerchantStatusError(err, "UpdateMerchantStatus", "FAILED_UPDATE_MERCHANT_STATUS", span, &status, zap.Int("merchant_id", *request.MerchantID))
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
		return errorhandler.HandleMarshalError[*response.MerchantResponse](s.logger, err, "UpdateMerchantStatus", "FAILED_MARSHAL_EMAIL_PAYLOAD", span, &status, merchant_errors.ErrFailedSendEmail, zap.Int("merchant_id", *request.MerchantID))
	}

	err = s.kafka.SendMessage("email-service-topic-merchant-update-status", strconv.Itoa(*request.MerchantID), payloadBytes)
	if err != nil {
		return errorhandler.HandleSendEmailError[*response.MerchantResponse](s.logger, err, "UpdateMerchantStatus", "FAILED_SEND_EMAIL", span, &status, merchant_errors.ErrFailedSendEmail, zap.Int("merchant_id", *request.MerchantID))
	}

	so := s.mapping.ToMerchantResponse(res)
	s.mencache.DeleteCachedMerchant(res.ID)

	s.logger.Debug("Successfully updated merchant status", zap.Int("merchant_id", res.ID))

	return so, nil
}

func (s *merchantCommandService) TrashedMerchant(merchant_id int) (*response.MerchantResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("TrashedMerchant", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "TrashedMerchant")
	defer span.End()

	span.SetAttributes(
		attribute.Int("merchant_id", merchant_id),
	)

	s.logger.Debug("Trashing merchant", zap.Int("merchant_id", merchant_id))

	res, err := s.merchantCommandRepository.TrashedMerchant(merchant_id)

	if err != nil {
		return s.errorHandler.HandleTrashedMerchantError(err, "TrashedMerchant", "FAILED_TRASHED_MERCHANT", span, &status, zap.Int("merchant_id", merchant_id))
	}
	so := s.mapping.ToMerchantResponse(res)

	s.mencache.DeleteCachedMerchant(res.ID)

	s.logger.Debug("Successfully trashed merchant", zap.Int("merchant_id", merchant_id))

	return so, nil
}

func (s *merchantCommandService) RestoreMerchant(merchant_id int) (*response.MerchantResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreMerchant", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreMerchant")
	defer span.End()

	span.SetAttributes(
		attribute.Int("merchant_id", merchant_id),
	)

	s.logger.Debug("Restoring merchant", zap.Int("merchant_id", merchant_id))

	res, err := s.merchantCommandRepository.RestoreMerchant(merchant_id)

	if err != nil {
		return s.errorHandler.HandleRestoreMerchantError(err, "RestoreMerchant", "FAILED_RESTORE_MERCHANT", span, &status, zap.Int("merchant_id", merchant_id))
	}
	s.logger.Debug("Successfully restored merchant", zap.Int("merchant_id", merchant_id))

	so := s.mapping.ToMerchantResponse(res)

	return so, nil
}

func (s *merchantCommandService) DeleteMerchantPermanent(merchant_id int) (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteMerchantPermanent", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteMerchantPermanent")
	defer span.End()

	s.logger.Debug("Deleting merchant permanently", zap.Int("merchant_id", merchant_id))

	_, err := s.merchantCommandRepository.DeleteMerchantPermanent(merchant_id)

	if err != nil {
		return s.errorHandler.HandleDeleteMerchantPermanentError(err, "DeleteMerchantPermanent", "FAILED_DELETE_MERCHANT_PERMANENT", span, &status, zap.Int("merchant_id", merchant_id))
	}

	s.logger.Debug("Successfully deleted merchant permanently", zap.Int("merchant_id", merchant_id))

	return true, nil
}

func (s *merchantCommandService) RestoreAllMerchant() (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreAllMerchant", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreAllMerchant")
	defer span.End()

	s.logger.Debug("Restoring all merchants")

	_, err := s.merchantCommandRepository.RestoreAllMerchant()

	if err != nil {
		return s.errorHandler.HandleRestoreAllMerchantError(err, "RestoreAllMerchant", "FAILED_RESTORE_ALL_MERCHANT", span, &status)
	}

	s.logger.Debug("Successfully restored all merchants")
	return true, nil
}

func (s *merchantCommandService) DeleteAllMerchantPermanent() (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteAllMerchantPermanent", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteAllMerchantPermanent")
	defer span.End()

	s.logger.Debug("Permanently deleting all merchants")

	_, err := s.merchantCommandRepository.DeleteAllMerchantPermanent()

	if err != nil {
		return s.errorHandler.HandleDeleteAllMerchantPermanentError(err, "DeleteAllMerchantPermanent", "FAILED_DELETE_ALL_MERCHANT_PERMANENT", span, &status)
	}

	s.logger.Debug("Successfully deleted all merchants permanently")
	return true, nil
}

func (s *merchantCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
