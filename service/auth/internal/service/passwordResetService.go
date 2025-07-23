package service

import (
	"context"
	"encoding/json"
	"errors"
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
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// PasswordResetServiceDeps holds the dependencies for creating a PasswordResetService.
type PasswordResetServiceDeps struct {
	// ErrorHandler handles general password reset errors.
	ErrorHandler errorhandler.PasswordResetErrorHandler

	// ErrorRandomString handles errors related to random string generation.
	ErrorRandomString errorhandler.RandomStringErrorHandler

	// ErrorMarshal handles JSON or data marshalling errors.
	ErrorMarshal errorhandler.MarshalErrorHandler

	// ErrorPassword handles password hashing or validation errors.
	ErrorPassword errorhandler.PasswordErrorHandler

	// ErrorKafka handles Kafka publishing errors.
	ErrorKafka errorhandler.KafkaErrorHandler

	// Cache provides caching for password reset tokens or sessions.
	Cache mencache.PasswordResetCache

	// Kafka is the Kafka client used to publish password reset events.
	Kafka *kafka.Kafka

	// Logger logs service-level events and errors.
	Logger logger.LoggerInterface

	// User provides access to user account data.
	User repository.UserRepository

	// ResetToken manages password reset token storage and retrieval.
	ResetToken repository.ResetTokenRepository
}

// passwordResetService implements the logic for handling password reset requests,
// including token generation, validation, and email notification via Kafka.
type passwordResetService struct {
	// Handles general password reset errors.
	errorhandler errorhandler.PasswordResetErrorHandler

	// Handles random string generation errors.
	errorRandomString errorhandler.RandomStringErrorHandler

	// Handles marshalling errors (e.g., JSON encoding).
	errorMarshal errorhandler.MarshalErrorHandler

	// Handles password validation and hashing errors.
	errorPassword errorhandler.PasswordErrorHandler

	// Handles Kafka-related publishing errors.
	errorKafka errorhandler.KafkaErrorHandler

	// Caching layer for password reset-related data.
	mencache mencache.PasswordResetCache

	// Kafka publisher for sending password reset events.
	kafka *kafka.Kafka

	// Logger for internal logs and errors.
	logger logger.LoggerInterface

	// Access to user account data.
	user repository.UserRepository

	// Repository for storing and validating reset tokens.
	resetToken repository.ResetTokenRepository

	observability observability.TraceLoggerObservability
}

// NewPasswordResetService initializes and returns a new instance of passwordResetService.
// It sets up Prometheus metrics for tracking request counts and durations, and registers these metrics.
// The function takes PasswordResetServiceDeps which includes context, error handlers, cache, logger,
// user repository, refresh token repository, token manager, and token service.
// Returns a pointer to the initialized passwordResetService.
func NewPasswordResetService(params *PasswordResetServiceDeps) *passwordResetService {
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

	observability := observability.NewTraceLoggerObservability(otel.Tracer("password-reset-service"), params.Logger, requestCounter, requestDuration)

	return &passwordResetService{
		errorhandler:      params.ErrorHandler,
		errorRandomString: params.ErrorRandomString,
		errorMarshal:      params.ErrorMarshal,
		errorPassword:     params.ErrorPassword,
		errorKafka:        params.ErrorKafka,
		mencache:          params.Cache,
		kafka:             params.Kafka,
		logger:            params.Logger,
		user:              params.User,
		resetToken:        params.ResetToken,
		observability:     observability,
	}
}

// ForgotPassword initiates the password reset process by sending a verification code to the user's email.
//
// Parameters:
//   - ctx: the context for the operation
//   - email: the user's email address
//
// Returns:
//   - true if the code was sent successfully, or an ErrorResponse if it fails.
func (s *passwordResetService) ForgotPassword(ctx context.Context, email string) (bool, *response.ErrorResponse) {
	const method = "ForgotPassword"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("email", email))

	defer func() {
		end(status)
	}()

	res, err := s.user.FindByEmail(ctx, email)
	if err != nil {
		return s.errorhandler.HandleFindEmailError(err, method, "FORGOT_PASSWORD_ERR", span, &status, zap.String("email", email))
	}

	span.SetAttributes(attribute.Int("user.id", res.ID))

	random, err := randomstring.GenerateRandomString(10)
	if err != nil {
		return s.errorRandomString.HandleRandomStringErrorForgotPassword(err, method, "FORGOT_PASSWORD_ERR", span, &status, zap.String("email", email), zap.Error(err))
	}

	_, err = s.resetToken.CreateResetToken(ctx, &requests.CreateResetTokenRequest{
		UserID:     res.ID,
		ResetToken: random,
		ExpiredAt:  time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		return s.errorhandler.HandleCreateResetTokenError(err, method, "FORGOT_PASSWORD_ERR", span, &status, zap.String("email", email), zap.Error(err))
	}

	s.mencache.SetResetTokenCache(ctx, random, res.ID, 5*time.Minute)

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
		return s.errorMarshal.HandleMarsalForgotPassword(err, method, "FORGOT_PASSWORD_ERR", span, &status, zap.Error(err))
	}

	err = s.kafka.SendMessage("email-service-topic-auth-forgot-password", strconv.Itoa(res.ID), payloadBytes)
	if err != nil {
		return s.errorKafka.HandleSendEmailForgotPassword(err, method, "FORGOT_PASSWORD_ERR", span, &status, zap.Error(err))
	}

	logSuccess("Successfully sent password reset email", zap.String("email", email))

	return true, nil
}

// ResetPassword sets a new password for the user using the provided reset token and new password.
//
// Parameters:
//   - ctx: the context for the operation
//   - request: the payload containing reset token and new password
//
// Returns:
//   - true if the password reset is successful, or an ErrorResponse if it fails.
func (s *passwordResetService) ResetPassword(ctx context.Context, req *requests.CreateResetPasswordRequest) (bool, *response.ErrorResponse) {
	const method = "ResetPassword"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("reset_token", req.ResetToken))

	defer func() {
		end(status)
	}()

	var userID int
	var found bool

	userID, found = s.mencache.GetResetTokenCache(ctx, req.ResetToken)
	if !found {
		res, err := s.resetToken.FindByToken(ctx, req.ResetToken)
		if err != nil {
			return s.errorhandler.HandleFindTokenError(err, method, "RESET_PASSWORD_ERR", span, &status, zap.String("reset_token", req.ResetToken))
		}
		userID = int(res.UserID)

		s.mencache.SetResetTokenCache(ctx, req.ResetToken, userID, 5*time.Minute)
	}

	if req.Password != req.ConfirmPassword {
		err := errors.New("password and confirm password do not match")
		return s.errorPassword.HandlePasswordNotMatchError(err, method, "RESET_PASSWORD_ERR", span, &status, zap.String("reset_token", req.ResetToken))
	}

	_, err := s.user.UpdateUserPassword(ctx, userID, req.Password)
	if err != nil {
		return s.errorhandler.HandleUpdatePasswordError(err, method, "RESET_PASSWORD_ERR", span, &status, zap.String("reset_token", req.ResetToken))
	}

	_ = s.resetToken.DeleteResetToken(ctx, userID)
	s.mencache.DeleteResetTokenCache(ctx, req.ResetToken)

	logSuccess("Successfully reset password", zap.String("reset_token", req.ResetToken))

	return true, nil
}

// VerifyCode validates the verification code sent to the user's email.
//
// Parameters:
//   - ctx: the context for the operation
//   - code: the verification code to validate
//
// Returns:
//   - true if the code is valid, or an ErrorResponse if the code is invalid or expired.
func (s *passwordResetService) VerifyCode(ctx context.Context, code string) (bool, *response.ErrorResponse) {
	const method = "VerifyCode"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("code", code))

	defer func() {
		end(status)
	}()

	res, err := s.user.FindByVerificationCode(ctx, code)
	if err != nil {
		return s.errorhandler.HandleVerifyCodeError(err, method, "VERIFY_CODE_ERR", span, &status, zap.String("code", code))
	}

	_, err = s.user.UpdateUserIsVerified(ctx, res.ID, true)
	if err != nil {
		return s.errorhandler.HandleUpdateVerifiedError(err, method, "VERIFY_CODE_ERR", span, &status, zap.Int("user.id", res.ID))
	}

	s.mencache.DeleteVerificationCodeCache(ctx, res.Email)

	htmlBody := emails.GenerateEmailHTML(map[string]string{
		"Title":   "Verification Success",
		"Message": "Your account has been successfully verified. Click the button below to view or manage your card.",
		"Button":  "Go to Dashboard",
		"Link":    "https://sanedge.example.com/card/create",
	})

	payloadBytes, err := json.Marshal(htmlBody)
	if err != nil {
		return s.errorMarshal.HandleMarshalVerifyCode(err, method, "SEND_EMAIL_VERIFY_CODE_ERR", span, &status, zap.Error(err))
	}

	err = s.kafka.SendMessage("email-service-topic-auth-verify-code-success", strconv.Itoa(res.ID), payloadBytes)
	if err != nil {
		return s.errorKafka.HandleSendEmailVerifyCode(err, method, "SEND_EMAIL_VERIFY_CODE_ERR", span, &status, zap.Error(err))
	}

	logSuccess("Successfully verify code", zap.String("code", code))

	return true, nil
}
