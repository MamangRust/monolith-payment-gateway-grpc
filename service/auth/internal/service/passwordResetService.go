package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-auth/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-auth/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-auth/internal/repository"

	emails "github.com/MamangRust/monolith-payment-gateway-pkg/email"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	randomstring "github.com/MamangRust/monolith-payment-gateway-pkg/random_string"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type passwordResetService struct {
	ctx               context.Context
	errorhandler      errorhandler.PasswordResetErrorHandler
	errorRandomString errorhandler.RandomStringErrorHandler
	errorMarshal      errorhandler.MarshalErrorHandler
	errorPassword     errorhandler.PasswordErrorHandler
	errorKafka        errorhandler.KafkaErrorHandler
	mencache          mencache.PasswordResetCache
	trace             trace.Tracer
	kafka             *kafka.Kafka
	logger            logger.LoggerInterface
	user              repository.UserRepository
	resetToken        repository.ResetTokenRepository
	requestCounter    *prometheus.CounterVec
	requestDuration   *prometheus.HistogramVec
}

func NewPasswordResetService(ctx context.Context,
	errorhandler errorhandler.PasswordResetErrorHandler,
	errorRandomString errorhandler.RandomStringErrorHandler,
	errorMarshal errorhandler.MarshalErrorHandler,
	errorPassword errorhandler.PasswordErrorHandler,
	errorKafka errorhandler.KafkaErrorHandler,
	mencache mencache.PasswordResetCache,
	kafka *kafka.Kafka, logger logger.LoggerInterface, user repository.UserRepository, resetToken repository.ResetTokenRepository) *passwordResetService {

	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "password_reset_service_requests_total",
			Help: "Total number of requests to the PasswordResetService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "password_reset_service_request_duration_seconds",
			Help:    "Histogram of request durations for the PasswordResetService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &passwordResetService{
		ctx:               ctx,
		errorhandler:      errorhandler,
		errorRandomString: errorRandomString,
		errorPassword:     errorPassword,
		errorMarshal:      errorMarshal,
		errorKafka:        errorKafka,
		mencache:          mencache,
		kafka:             kafka,
		trace:             otel.Tracer("password-reset-service"),
		user:              user,
		logger:            logger,
		resetToken:        resetToken,
		requestCounter:    requestCounter,
		requestDuration:   requestDuration,
	}
}

func (s *passwordResetService) ForgotPassword(email string) (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("ForgotPassword", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "PasswordService.ForgotPassword")
	defer span.End()

	span.SetAttributes(attribute.String("user.email", email))
	s.logger.Debug("Starting forgot password process", zap.String("email", email))

	res, err := s.user.FindByEmail(email)
	if err != nil {
		return s.errorhandler.HandleFindEmailError(err, "ForgotPassword", "FORGOT_PASSWORD_ERR", span, &status, zap.String("email", email))
	}

	span.SetAttributes(attribute.Int("user.id", res.ID))

	random, err := randomstring.GenerateRandomString(10)
	if err != nil {
		return s.errorRandomString.HandleRandomStringErrorForgotPassword(err, "ForgotPassword", "FORGOT_PASSWORD_ERR", span, &status, zap.String("email", email))
	}

	_, err = s.resetToken.CreateResetToken(&requests.CreateResetTokenRequest{
		UserID:     res.ID,
		ResetToken: random,
		ExpiredAt:  time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		return s.errorhandler.HandleCreateResetTokenError(err, "ForgotPassword", "FORGOT_PASSWORD_ERR", span, &status, zap.String("email", email))
	}

	s.mencache.SetResetTokenCache(random, res.ID, 5*time.Minute)

	htmlBody := emails.GenerateEmailHTML(map[string]string{
		"Title":   "Reset Your Password",
		"Message": "Click the button below to reset your password.",
		"Button":  "Reset Password",
		"Link":    "https://sanedge.example.com/reset-password?token=" + random,
	})

	emailPayload := map[string]any{
		"email":   res.Email,
		"subject": "Password Reset Request",
		"body":    htmlBody,
	}

	payloadBytes, err := json.Marshal(emailPayload)
	if err != nil {
		return s.errorMarshal.HandleMarsalForgotPassword(err, "ForgotPassword", "FORGOT_PASSWORD_ERR", span, &status)
	}

	err = s.kafka.SendMessage("email-service-topic-auth-forgot-password", strconv.Itoa(res.ID), payloadBytes)
	if err != nil {
		return s.errorKafka.HandleSendEmailForgotPassword(err, "ForgotPassword", "FORGOT_PASSWORD_ERR", span, &status)
	}

	s.logger.Debug("Password reset email sent successfully",
		zap.Int("user_id", res.ID),
		zap.String("email", res.Email),
	)
	span.SetStatus(codes.Ok, "Password reset initiated")
	return true, nil
}

func (s *passwordResetService) ResetPassword(req *requests.CreateResetPasswordRequest) (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("ResetPassword", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "PasswordService.ResetPassword")
	defer span.End()

	span.SetAttributes(attribute.String("reset_token", req.ResetToken))
	s.logger.Debug("Starting password reset process", zap.String("reset_token", req.ResetToken))

	var userID int
	var found bool

	userID, found = s.mencache.GetResetTokenCache(req.ResetToken)
	if !found {
		res, err := s.resetToken.FindByToken(req.ResetToken)
		if err != nil {
			return s.errorhandler.HandleFindTokenError(err, "ResetPassword", "RESET_PASSWORD_ERR", span, &status, zap.String("reset_token", req.ResetToken))
		}
		userID = int(res.UserID)

		s.mencache.SetResetTokenCache(req.ResetToken, userID, 5*time.Minute)
	}

	if req.Password != req.ConfirmPassword {
		return s.errorPassword.HandlePasswordNotMatchError(nil, "ResetPassword", "RESET_PASSWORD_ERR", span, &status, zap.String("reset_token", req.ResetToken))
	}

	_, err := s.user.UpdateUserPassword(userID, req.Password)
	if err != nil {
		return s.errorhandler.HandleUpdatePasswordError(err, "ResetPassword", "RESET_PASSWORD_ERR", span, &status, zap.String("reset_token", req.ResetToken))
	}

	_ = s.resetToken.DeleteResetToken(userID)
	s.mencache.DeleteResetTokenCache(req.ResetToken)

	s.logger.Debug("Password reset successfully", zap.Int("user_id", userID))
	span.SetStatus(codes.Ok, "Password reset successful")
	return true, nil
}

func (s *passwordResetService) VerifyCode(code string) (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("VerifyCode", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "PasswordService.VerifyCode")
	defer span.End()

	span.SetAttributes(
		attribute.String("verification_code", code),
	)

	s.logger.Debug("Starting verification code process",
		zap.String("code", code),
	)

	res, err := s.user.FindByVerificationCode(code)
	if err != nil {
		return s.errorhandler.HandleVerifyCodeError(err, "VerifyCode", "VERIFY_CODE_ERR", span, &status, zap.String("code", code))
	}

	_, err = s.user.UpdateUserIsVerified(res.ID, true)
	if err != nil {
		return s.errorhandler.HandleUpdateVerifiedError(err, "VerifyCode", "VERIFY_CODE_ERR", span, &status, zap.Int("user_id", res.ID))
	}

	s.mencache.DeleteVerificationCodeCache("register:verify:" + res.Email)

	s.logger.Debug("User verified successfully",
		zap.Int("user_id", res.ID),
	)
	return true, nil
}

func (s *passwordResetService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
