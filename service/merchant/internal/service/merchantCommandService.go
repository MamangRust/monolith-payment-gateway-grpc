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

type merchantCommandService struct {
	ctx                       context.Context
	kafka                     *kafka.Kafka
	trace                     trace.Tracer
	errorHandler              errorhandler.MerchantCommandErrorHandler
	mencache                  mencache.MerchantCommandCache
	userRepository            repository.UserRepository
	merchantQueryRepository   repository.MerchantQueryRepository
	merchantCommandRepository repository.MerchantCommandRepository
	logger                    logger.LoggerInterface
	mapping                   responseservice.MerchantResponseMapper
	requestCounter            *prometheus.CounterVec
	requestDuration           *prometheus.HistogramVec
}

func NewMerchantCommandService(kafka *kafka.Kafka, ctx context.Context,
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
	const method = "CreateMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	user, err := s.userRepository.FindById(request.UserID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantResponse](s.logger, err, method, "FAILED_FIND_USER_BY_ID", span, &status, user_errors.ErrUserNotFoundRes, zap.Error(err))
	}

	res, err := s.merchantCommandRepository.CreateMerchant(request)

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

	so := s.mapping.ToMerchantResponse(res)

	logSuccess("Successfully created merchant", zap.Bool("success", true))

	return so, nil
}

func (s *merchantCommandService) UpdateMerchant(request *requests.UpdateMerchantRequest) (*response.MerchantResponse, *response.ErrorResponse) {
	const method = "UpdateMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	res, err := s.merchantCommandRepository.UpdateMerchant(request)

	if err != nil {
		return s.errorHandler.HandleUpdateMerchantError(err, method, "FAILED_UPDATE_MERCHANT", span, &status, zap.Error(err))
	}

	so := s.mapping.ToMerchantResponse(res)

	s.mencache.DeleteCachedMerchant(res.ID)

	logSuccess("Successfully updated merchant", zap.Bool("success", true))

	return so, nil
}

func (s *merchantCommandService) UpdateMerchantStatus(request *requests.UpdateMerchantStatusRequest) (*response.MerchantResponse, *response.ErrorResponse) {
	const method = "UpdateMerchantStatus"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantQueryRepository.FindById(*request.MerchantID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantResponse](s.logger, err, method, "FAILED_FIND_MERCHANT_BY_ID", span, &status, merchant_errors.ErrFailedFindMerchantById, zap.Error(err))
	}

	user, err := s.userRepository.FindById(merchant.UserID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantResponse](s.logger, err, method, "FAILED_FIND_USER_BY_ID", span, &status, user_errors.ErrUserNotFoundRes, zap.Error(err))
	}

	res, err := s.merchantCommandRepository.UpdateMerchantStatus(request)

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

	so := s.mapping.ToMerchantResponse(res)
	s.mencache.DeleteCachedMerchant(res.ID)

	logSuccess("Successfully updated merchant status", zap.Bool("success", true))

	return so, nil
}

func (s *merchantCommandService) TrashedMerchant(merchant_id int) (*response.MerchantResponse, *response.ErrorResponse) {
	const method = "TrashedMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	res, err := s.merchantCommandRepository.TrashedMerchant(merchant_id)

	if err != nil {
		return s.errorHandler.HandleTrashedMerchantError(err, method, "FAILED_TRASHED_MERCHANT", span, &status, zap.Error(err))
	}
	so := s.mapping.ToMerchantResponse(res)

	s.mencache.DeleteCachedMerchant(res.ID)

	logSuccess("Successfully trashed merchant", zap.Bool("success", true))

	return so, nil
}

func (s *merchantCommandService) RestoreMerchant(merchant_id int) (*response.MerchantResponse, *response.ErrorResponse) {
	const method = "RestoreMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	res, err := s.merchantCommandRepository.RestoreMerchant(merchant_id)

	if err != nil {
		return s.errorHandler.HandleRestoreMerchantError(err, method, "FAILED_RESTORE_MERCHANT", span, &status, zap.Error(err))
	}

	so := s.mapping.ToMerchantResponse(res)

	s.mencache.DeleteCachedMerchant(res.ID)

	logSuccess("Successfully restored merchant", zap.Bool("success", true))

	return so, nil
}

func (s *merchantCommandService) DeleteMerchantPermanent(merchant_id int) (bool, *response.ErrorResponse) {
	const method = "DeleteMerchantPermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	_, err := s.merchantCommandRepository.DeleteMerchantPermanent(merchant_id)

	if err != nil {
		return s.errorHandler.HandleDeleteMerchantPermanentError(err, method, "FAILED_DELETE_MERCHANT_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted merchant permanently", zap.Bool("success", true))

	return true, nil
}

func (s *merchantCommandService) RestoreAllMerchant() (bool, *response.ErrorResponse) {
	const method = "RestoreAllMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	_, err := s.merchantCommandRepository.RestoreAllMerchant()

	if err != nil {
		return s.errorHandler.HandleRestoreAllMerchantError(err, method, "FAILED_RESTORE_ALL_MERCHANT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully restored all merchants", zap.Bool("success", true))

	return true, nil
}

func (s *merchantCommandService) DeleteAllMerchantPermanent() (bool, *response.ErrorResponse) {
	const method = "DeleteAllMerchantPermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	_, err := s.merchantCommandRepository.DeleteAllMerchantPermanent()

	if err != nil {
		return s.errorHandler.HandleDeleteAllMerchantPermanentError(err, method, "FAILED_DELETE_ALL_MERCHANT_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted all merchants permanently", zap.Bool("success", true))

	return true, nil
}

func (s *merchantCommandService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *merchantCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
