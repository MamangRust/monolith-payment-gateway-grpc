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
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// LoginServiceDeps holds the parameters for the login service.
type LoginServiceDeps struct {
	// ErrorPassword handles password-related errors.
	ErrorPassword errorhandler.PasswordErrorHandler

	// ErrorToken handles token-related errors.
	ErrorToken errorhandler.TokenErrorHandler

	// ErrorHandler handles login-related errors.
	ErrorHandler errorhandler.LoginErrorHandler

	// Cache is the cache for login/session data.
	Cache mencache.LoginCache

	// Logger logs events and errors.
	Logger logger.LoggerInterface

	// Hash verifies and compares password hashes.
	Hash hash.HashPassword

	// UserRepository accesses user data from the database.
	UserRepository repository.UserRepository

	// RefreshToken manages refresh token storage.
	RefreshToken repository.RefreshTokenRepository

	// Token manages access and refresh token generation.
	Token auth.TokenManager

	// TokenService handles session token management.
	TokenService *tokenService
}

// loginService is the implementation of the LoginService interface.
// It handles the logic for authenticating users, validating passwords,
// issuing tokens, and caching login sessions.
type loginService struct {
	// Handles password-related errors.
	errorPassword errorhandler.PasswordErrorHandler

	// Handles token-related errors.
	errorToken errorhandler.TokenErrorHandler

	// Handles general login-related errors.
	errorHandler errorhandler.LoginErrorHandler

	// Cache for login attempts or user sessions.
	mencache mencache.LoginCache

	// Logging utility.
	logger logger.LoggerInterface

	// Responsible for comparing hashed passwords.
	hash hash.HashPassword

	// Repository for fetching user data.
	user repository.UserRepository

	// Repository for managing refresh tokens.
	refreshToken repository.RefreshTokenRepository

	// Token generator and validator.
	token auth.TokenManager

	// Token service used for managing session tokens.
	tokenService *tokenService

	observability observability.TraceLoggerObservability
}

// NewLoginService initializes and returns a new instance of loginService.
// It sets up Prometheus metrics for tracking request counts and durations, and registers these metrics.
// The function takes LoginServiceDeps which includes context, error handlers, cache, logger,
// user repository, refresh token repository, token manager, and token service.
// Returns a pointer to the initialized loginService.
func NewLoginService(params *LoginServiceDeps) *loginService {
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

	observability := observability.NewTraceLoggerObservability(otel.Tracer("login-service"), params.Logger, requestCounter, requestDuration)

	return &loginService{
		errorPassword: params.ErrorPassword,
		errorToken:    params.ErrorToken,
		errorHandler:  params.ErrorHandler,
		mencache:      params.Cache,
		logger:        params.Logger,
		hash:          params.Hash,
		user:          params.UserRepository,
		refreshToken:  params.RefreshToken,
		token:         params.Token,
		tokenService:  params.TokenService,
		observability: observability,
	}
}

// Login authenticates a user using their credentials and returns a token upon success.
//
// Parameters:
//   - ctx: the context for the operation
//   - request: the authentication request payload containing email and password
//
// Returns:
//   - A TokenResponse if authentication is successful, or an ErrorResponse if it fails.
func (s *loginService) Login(ctx context.Context, request *requests.AuthRequest) (*response.TokenResponse, *response.ErrorResponse) {
	const method = "Login"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("email", request.Email))

	defer func() {
		end(status)
	}()

	if cachedToken, found := s.mencache.GetCachedLogin(ctx, request.Email); found {
		logSuccess("Successfully logged in", zap.String("email", request.Email))
		return cachedToken, nil
	}

	res, err := s.user.FindByEmailAndVerify(ctx, request.Email)
	if err != nil {
		return s.errorHandler.HandleFindEmailError(err, method, "LOGIN_ERR", span, &status, zap.Error(err))
	}

	err = s.hash.ComparePassword(res.Password, request.Password)
	if err != nil {
		return s.errorPassword.HandleComparePasswordError(err, method, "COMPARE_PASSWORD_ERR", span, &status, zap.Error(err))
	}

	token, err := s.tokenService.createAccessToken(res.ID)
	if err != nil {
		return s.errorToken.HandleCreateAccessTokenError(err, method, "CREATE_ACCESS_TOKEN_ERR", span, &status, zap.Error(err))
	}

	refreshToken, err := s.tokenService.createRefreshToken(ctx, res.ID)
	if err != nil {
		return s.errorToken.HandleCreateRefreshTokenError(err, method, "CREATE_REFRESH_TOKEN_ERR", span, &status, zap.Error(err))
	}

	tokenResp := &response.TokenResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}

	s.mencache.SetCachedLogin(ctx, request.Email, tokenResp, time.Minute)

	logSuccess("Successfully logged in", zap.String("email", request.Email))

	return tokenResp, nil
}
