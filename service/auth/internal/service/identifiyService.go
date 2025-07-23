package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-auth/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-auth/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-auth/internal/repository"

	"github.com/MamangRust/monolith-payment-gateway-pkg/auth"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/user"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

// IdentityServiceDeps holds dependencies for the IdentityService.
type IdentityServiceDeps struct {
	// ErrorHandler handles identity-related errors.
	ErrorHandler errorhandler.IdentityErrorHandler

	// ErrorToken handles token-related errors.
	ErrorToken errorhandler.TokenErrorHandler

	// Cache provides caching for identity data.
	Cache mencache.IdentityCache

	// Token manages token generation and validation.
	Token auth.TokenManager

	// RefreshToken provides access to refresh token data.
	RefreshToken repository.RefreshTokenRepository

	// User provides access to user data.
	User repository.UserRepository

	// Logger logs system events and errors.
	Logger logger.LoggerInterface

	// Mapping maps user data to response models.
	Mapping responseservice.UserQueryResponseMapper

	// TokenService manages advanced token-related logic.
	TokenService *tokenService
}

// identityService is the implementation of the identity service.
type identityService struct {
	// errorhandler is the error handler for identity-related errors.
	errorhandler errorhandler.IdentityErrorHandler

	// errorToken is the error handler for token-related errors.
	errorToken errorhandler.TokenErrorHandler

	// mencache is the cache for identity-related data.
	mencache mencache.IdentityCache

	// logger is the logger for logging events and errors.
	logger logger.LoggerInterface

	// token is the token manager for generating and validating tokens.
	token auth.TokenManager

	// refreshToken is the repository for managing refresh tokens.
	refreshToken repository.RefreshTokenRepository

	// user is the repository for managing user data.
	user repository.UserRepository

	// mapper is the mapper for converting user data to a response format.
	mapper responseservice.UserQueryResponseMapper

	// tokenService is the token service for managing tokens.
	tokenService *tokenService

	observability observability.TraceLoggerObservability
}

// NewIdentityService initializes and returns the IdentityService with the given parameters.
//
// It sets up the prometheus metrics for request counters and durations, and registers them.
//
// It returns the initialized IdentityService.
func NewIdentityService(param *IdentityServiceDeps) *identityService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "identity_service_requests_total",
			Help: "Total number of auth requests",
		},
		[]string{"method", "status"},
	)
	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "identity_service_request_duration_seconds",
			Help:    "Duration of auth requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("identity-service"), param.Logger, requestCounter, requestDuration)

	return &identityService{
		errorhandler:  param.ErrorHandler,
		errorToken:    param.ErrorToken,
		mencache:      param.Cache,
		logger:        param.Logger,
		token:         param.Token,
		refreshToken:  param.RefreshToken,
		user:          param.User,
		mapper:        param.Mapping,
		tokenService:  param.TokenService,
		observability: observability,
	}
}

// RefreshToken generates a new access token using a valid refresh token.
//
// Parameters:
//   - ctx: the context for the operation (used for timeout, tracing, etc.)
//   - token: the refresh token string
//
// Returns:
//   - A new TokenResponse if the token is valid, or an ErrorResponse if the refresh fails.
func (s *identityService) RefreshToken(ctx context.Context, token string) (*response.TokenResponse, *response.ErrorResponse) {
	const method = "RefreshToken"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("token", token))

	defer func() {
		end(status)
	}()

	if cachedUserID, found := s.mencache.GetRefreshToken(ctx, token); found {
		userId, err := strconv.Atoi(cachedUserID)
		if err == nil {
			s.mencache.DeleteRefreshToken(ctx, token)
			s.logger.Debug("Invalidated old refresh token from cache", zap.String("token", token))

			accessToken, err := s.tokenService.createAccessToken(userId)
			if err != nil {
				return s.errorToken.HandleCreateAccessTokenError(err, method, "CREATE_ACCESS_TOKEN_FAILED", span, &status, zap.Int("user.id", userId))
			}

			refreshToken, err := s.tokenService.createRefreshToken(ctx, userId)
			if err != nil {
				return s.errorToken.HandleCreateRefreshTokenError(err, method, "CREATE_REFRESH_TOKEN_FAILED", span, &status, zap.Int("user.id", userId))
			}

			expiryTime := time.Now().Add(24 * time.Hour)
			expirationDuration := time.Until(expiryTime)

			s.mencache.SetRefreshToken(ctx, refreshToken, expirationDuration)
			s.logger.Debug("Stored new refresh token in cache",
				zap.String("new_token", refreshToken),
				zap.Duration("expiration", expirationDuration))

			s.logger.Debug("Refresh token refreshed successfully (cached)", zap.Int("user_id", userId))
			span.SetStatus(codes.Ok, "Token refreshed successfully from cache")

			return &response.TokenResponse{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			}, nil
		}
	}

	userIdStr, err := s.token.ValidateToken(token)
	if err != nil {
		if errors.Is(err, auth.ErrTokenExpired) {
			s.mencache.DeleteRefreshToken(ctx, token)
			if err := s.refreshToken.DeleteRefreshToken(ctx, token); err != nil {

				return s.errorhandler.HandleDeleteRefreshTokenError(err, method, "DELETE_REFRESH_TOKEN", span, &status, zap.String("token", token))
			}

			return s.errorhandler.HandleExpiredRefreshTokenError(err, method, "TOKEN_EXPIRED", span, &status, zap.String("token", token))
		}

		return s.errorhandler.HandleInvalidTokenError(err, method, "INVALID_TOKEN", span, &status, zap.String("token", token))
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {

		return errorhandler.HandleInvalidFormatUserIDError[*response.TokenResponse](s.logger, err, method, "INVALID_USER_ID", span, &status, zap.Int("user.id", userId))
	}

	span.SetAttributes(attribute.Int("user.id", userId))

	s.mencache.DeleteRefreshToken(ctx, token)
	if err := s.refreshToken.DeleteRefreshToken(ctx, token); err != nil {
		s.logger.Debug("Failed to delete old refresh token", zap.Error(err))

		return s.errorhandler.HandleDeleteRefreshTokenError(err, method, "DELETE_REFRESH_TOKEN", span, &status, zap.String("token", token))
	}

	accessToken, err := s.tokenService.createAccessToken(userId)
	if err != nil {

		return s.errorToken.HandleCreateAccessTokenError(err, method, "CREATE_ACCESS_TOKEN_FAILED", span, &status, zap.Int("user.id", userId))
	}

	refreshToken, err := s.tokenService.createRefreshToken(ctx, userId)
	if err != nil {

		return s.errorToken.HandleCreateRefreshTokenError(err, method, "CREATE_REFRESH_TOKEN_FAILED", span, &status, zap.Int("user.id", userId))
	}

	expiryTime := time.Now().Add(24 * time.Hour)
	updateRequest := &requests.UpdateRefreshToken{
		UserId:    userId,
		Token:     refreshToken,
		ExpiresAt: expiryTime.Format("2006-01-02 15:04:05"),
	}

	if _, err = s.refreshToken.UpdateRefreshToken(ctx, updateRequest); err != nil {
		s.mencache.DeleteRefreshToken(ctx, refreshToken)

		return s.errorhandler.HandleUpdateRefreshTokenError(err, method, "UPDATE_REFRESH_TOKEN_FAILED", span, &status, zap.Int("user.id", userId))
	}

	expirationDuration := time.Until(expiryTime)
	s.mencache.SetRefreshToken(ctx, refreshToken, expirationDuration)

	logSuccess("Refresh token refreshed successfully", zap.Int("user.id", userId))

	return &response.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// GetMe retrieves the current user's profile information based on the access token.
//
// Parameters:
//   - ctx: the context for the operation
//   - token: the access token string
//
// Returns:
//   - A UserResponse representing the authenticated user, or an ErrorResponse if unauthorized or failed.
func (s *identityService) GetMe(ctx context.Context, token string) (*response.UserResponse, *response.ErrorResponse) {
	const method = "GetMe"
	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("token", token))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Fetching user details", zap.String("token", token))

	userIdStr, err := s.token.ValidateToken(token)
	if err != nil {
		status = "error"

		return s.errorhandler.HandleValidateTokenError(err, method, "INVALID_TOKEN", span, &status, zap.String("token", token))
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		status = "error"

		return errorhandler.HandleInvalidFormatUserIDError[*response.UserResponse](
			s.logger, err, method, "INVALID_USER_ID", span, &status, zap.String("user_id_str", userIdStr),
		)
	}

	if cachedUser, found := s.mencache.GetCachedUserInfo(ctx, userIdStr); found {
		logSuccess("User info retrieved from cache", zap.Int("user.id", userId))
		return cachedUser, nil
	}

	user, err := s.user.FindById(ctx, userId)
	if err != nil {
		status = "error"

		return s.errorhandler.HandleFindByIdError(err, method, "FAILED_FETCH_USER", span, &status, zap.Int("user.id", userId))
	}

	userResponse := s.mapper.ToUserResponse(user)

	s.mencache.SetCachedUserInfo(ctx, userResponse, time.Minute*5)

	logSuccess("User details fetched successfully", zap.Int("user.id", userId))

	return userResponse, nil
}
