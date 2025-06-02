package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-auth/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-auth/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-auth/internal/repository"

	"github.com/MamangRust/monolith-payment-gateway-pkg/auth"
	"github.com/MamangRust/monolith-payment-gateway-pkg/hash"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type loginService struct {
	ctx             context.Context
	errorPassword   errorhandler.PasswordErrorHandler
	errorToken      errorhandler.TokenErrorHandler
	errorHandler    errorhandler.LoginErrorHandler
	mencache        mencache.LoginCache
	logger          logger.LoggerInterface
	hash            hash.HashPassword
	user            repository.UserRepository
	refreshToken    repository.RefreshTokenRepository
	token           auth.TokenManager
	trace           trace.Tracer
	tokenService    tokenService
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewLoginService(
	ctx context.Context,
	errorPassword errorhandler.PasswordErrorHandler,
	errorToken errorhandler.TokenErrorHandler,
	errorHandler errorhandler.LoginErrorHandler,
	mencache mencache.LoginCache,
	logger logger.LoggerInterface,
	hash hash.HashPassword,
	userRepository repository.UserRepository,
	refreshToken repository.RefreshTokenRepository,
	token auth.TokenManager,
	tokenService tokenService,
) *loginService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "login_service_requests_total",
			Help: "Total number of auth requests",
		},
		[]string{"method", "status"},
	)
	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "login_service_request_duration_seconds",
			Help:    "Duration of auth requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &loginService{
		ctx:             ctx,
		errorPassword:   errorPassword,
		errorToken:      errorToken,
		errorHandler:    errorHandler,
		mencache:        mencache,
		logger:          logger,
		hash:            hash,
		user:            userRepository,
		refreshToken:    refreshToken,
		token:           token,
		trace:           otel.Tracer("login-service"),
		tokenService:    tokenService,
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}
}

func (s *loginService) Login(request *requests.AuthRequest) (*response.TokenResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("Login", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "LoginService.Login")
	defer span.End()

	span.SetAttributes(
		attribute.String("user.email", request.Email),
	)

	s.logger.Debug("Starting login process",
		zap.String("email", request.Email),
	)

	cachedToken := s.mencache.GetCachedLogin(request.Email)
	if cachedToken != nil {
		s.logger.Debug("Returning cached login token", zap.String("email", request.Email))
		span.SetStatus(codes.Ok, "Login from cache")
		return cachedToken, nil
	}

	res, err := s.user.FindByEmail(request.Email)
	if err != nil {
		return s.errorHandler.HandleFindEmailError(err, "Login", "LOGIN_ERR", span, &status, zap.Error(err))
	}

	span.SetAttributes(
		attribute.Int("user.id", res.ID),
	)

	err = s.hash.ComparePassword(res.Password, request.Password)
	if err != nil {
		return s.errorPassword.HandleComparePasswordError(err, "Login", "COMPARE_PASSWORD_ERR", span, &status, zap.Error(err))
	}

	token, err := s.tokenService.createAccessToken(s.ctx, res.ID)
	if err != nil {
		return s.errorToken.HandleCreateAccessTokenError(err, "Login", "CREATE_ACCESS_TOKEN_ERR", span, &status, zap.Error(err))
	}

	refreshToken, err := s.tokenService.createRefreshToken(s.ctx, res.ID)
	if err != nil {
		return s.errorToken.HandleCreateRefreshTokenError(err, "Login", "CREATE_REFRESH_TOKEN_ERR", span, &status, zap.Error(err))
	}

	tokenResp := &response.TokenResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}

	s.mencache.SetCachedLogin(request.Email, tokenResp, time.Minute)

	s.logger.Debug("User logged in successfully",
		zap.String("email", request.Email),
		zap.Int("userID", res.ID),
	)
	span.SetStatus(codes.Ok, "Login successful")

	return tokenResp, nil
}

func (s *loginService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
