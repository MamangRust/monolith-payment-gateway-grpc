package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-auth/internal/repository"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/MamangRust/monolith-payment-gateway-pkg/auth"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	traceunic "github.com/MamangRust/monolith-payment-gateway-pkg/trace_unic"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// tokenServiceDeps holds the dependencies required to construct a tokenService.
type tokenServiceDeps struct {
	// Ctx carries deadlines, cancelation signals, and other request-scoped values.
	Ctx context.Context

	// Token manages JWT access and refresh token generation and validation.
	Token auth.TokenManager

	// RefreshToken provides access to the refresh token repository.
	RefreshToken repository.RefreshTokenRepository

	// Logger handles logging of service events and errors.
	Logger logger.LoggerInterface
}

// tokenService provides operations for issuing, refreshing, and revoking tokens.
// It includes observability instrumentation and logging.
type tokenService struct {
	// ctx carries deadlines, cancelation signals, and other request-scoped values.
	ctx context.Context

	// refreshToken accesses the persistence layer for refresh tokens.
	refreshToken repository.RefreshTokenRepository

	// token handles creation and validation of JWTs.
	token auth.TokenManager

	// logger records logs related to token operations.
	logger logger.LoggerInterface

	// trace enables distributed tracing with OpenTelemetry.
	trace trace.Tracer

	// requestCounter counts the number of token-related requests (Prometheus metric).
	requestCounter *prometheus.CounterVec

	// requestDuration measures the latency of token-related requests (Prometheus metric).
	requestDuration *prometheus.HistogramVec
}

// NewTokenService initializes and returns a new instance of tokenService.
// It sets up Prometheus metrics for tracking request counts and durations,
// and registers these metrics. The function takes tokenServiceDeps which
// includes context, token manager, refresh token repository, and logger.
// Returns a pointer to the initialized tokenService with tracing enabled.
func NewTokenService(
	params *tokenServiceDeps,
) *tokenService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "token_service_requests_total",
			Help: "Total number of auth requests",
		},
		[]string{"method", "status"},
	)
	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "token_service_request_duration_seconds",
			Help:    "Duration of auth requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	return &tokenService{
		refreshToken:    params.RefreshToken,
		token:           params.Token,
		logger:          params.Logger,
		trace:           otel.Tracer("token-service"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}
}

// createAccessToken generates an access token for a given user ID.
// It initiates tracing and logging for the token creation process.
// The function returns the generated token as a string if successful,
// or an error if the token generation fails. Tracing and logging are
// used to record the success or failure of the operation.
func (s *tokenService) createAccessToken(id int) (string, error) {
	const method = "createAccessToken"

	end, logSuccess, status, logError := s.startTracingAndLogging(method, attribute.Int("user.id", id))

	defer func() {
		end(status)
	}()

	res, err := s.token.GenerateToken(id, "access")
	if err != nil {
		status = "error"
		traceId := traceunic.GenerateTraceID("ACCESS_TOKEN_FAILED")

		logError(traceId, "Failed to create access token", err,
			zap.Int("userID", id),
			zap.Error(err),
		)

		return "", err
	}

	logSuccess("Created access token",
		zap.Int("userID", id),
	)

	return res, nil
}

// createRefreshToken generates a new refresh token for a given user ID.
// It initiates tracing and logging for the token creation process.
// The function deletes any existing refresh tokens for the user before
// creating a new one.
// The function returns the generated token as a string if successful,
// or an error if the token generation fails. Tracing and logging are
// used to record the success or failure of the operation.
func (s *tokenService) createRefreshToken(ctx context.Context, id int) (string, error) {
	const method = "createRefreshToken"

	end, logSuccess, status, logError := s.startTracingAndLogging(method, attribute.Int("user.id", id))

	defer func() {
		end(status)
	}()

	res, err := s.token.GenerateToken(id, "refresh")
	if err != nil {
		status = "error"
		traceId := traceunic.GenerateTraceID("REFRESH_TOKEN_FAILED")

		logError(traceId, "Failed to create refresh token", err, zap.Int("user.id", id), zap.Error(err))
		return "", err
	}

	if err := s.refreshToken.DeleteRefreshTokenByUserId(ctx, id); err != nil && !errors.Is(err, sql.ErrNoRows) {
		status = "error"

		traceId := traceunic.GenerateTraceID("DELETE_REFRESH_TOKEN_ERR")

		logError(traceId, "Failed to delete existing refresh token", err, zap.Int("userID", id), zap.Error(err))

		return "", err
	}

	_, err = s.refreshToken.CreateRefreshToken(ctx, &requests.CreateRefreshToken{
		Token:     res,
		UserId:    id,
		ExpiresAt: time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		status = "error"

		traceId := traceunic.GenerateTraceID("CREATE_REFRESH_TOKEN_ERR")

		logError(traceId, "Failed to create refresh token", err, zap.Int("userID", id), zap.Error(err))

		return "", err
	}

	logSuccess("Created refresh token",
		zap.Int("userID", id),
	)

	return res, nil
}

// startTracingAndLogging initializes tracing and logging for a given method.
// It starts a span with optional attributes and logs the method start.
// It returns the span, a function to end the span and record metrics, the initial
// status of the operation, and a function to log success messages.
//
// Parameters:
//   - method: The name of the method to trace and log.
//   - attrs: Optional attributes to add to the span.
//
// Returns:
//   - trace.Span: The OpenTelemetry span for the traced method.
//   - func(string): Function to end the span with a given status, recording metrics.
//   - string: Initial status of the operation, defaulting to "success".
//   - func(string, ...zap.Field): Function to log success messages with optional fields.
func (s *tokenService) startTracingAndLogging(
	method string,
	attrs ...attribute.KeyValue,
) (
	end func(string),
	logSuccess func(string, ...zap.Field),
	status string,
	logError func(traceID string, msg string, err error, fields ...zap.Field),
) {
	start := time.Now()
	status = "success"

	_, span := s.trace.Start(s.ctx, method)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.AddEvent("Start: " + method)
	s.logger.Debug("Start: " + method)

	end = func(status string) {
		s.recordMetrics(method, status, start)

		code := codes.Ok
		if status != "success" {
			code = codes.Error
		}

		span.SetStatus(code, status)
		span.End()
	}

	logSuccess = func(msg string, fields ...zap.Field) {
		span.AddEvent(msg)
		s.logger.Debug(msg, fields...)
	}

	logError = func(traceID string, msg string, err error, fields ...zap.Field) {
		span.RecordError(err)
		span.SetStatus(codes.Error, msg)
		span.AddEvent(msg)

		allFields := append([]zap.Field{
			zap.String("trace.id", traceID),
			zap.Error(err),
		}, fields...)

		s.logger.Error(msg, allFields...)
	}

	return end, logSuccess, status, logError
}

// recordMetrics records a Prometheus metric for the given method and status.
// It increments a counter and records the duration since the provided start time.
func (s *tokenService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
