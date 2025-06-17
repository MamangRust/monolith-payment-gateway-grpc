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
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type identityService struct {
	ctx             context.Context
	errorhandler    errorhandler.IdentityErrorHandler
	errorToken      errorhandler.TokenErrorHandler
	mencache        mencache.IdentityCache
	trace           trace.Tracer
	logger          logger.LoggerInterface
	token           auth.TokenManager
	refreshToken    repository.RefreshTokenRepository
	user            repository.UserRepository
	mapping         responseservice.UserResponseMapper
	tokenService    tokenService
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewIdentityService(ctx context.Context, errohandler errorhandler.IdentityErrorHandler, errorToken errorhandler.TokenErrorHandler, mencache mencache.IdentityCache, token auth.TokenManager, refreshToken repository.RefreshTokenRepository, user repository.UserRepository, logger logger.LoggerInterface, mapping responseservice.UserResponseMapper, tokenService tokenService) *identityService {
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

	return &identityService{
		ctx:             ctx,
		errorhandler:    errohandler,
		errorToken:      errorToken,
		mencache:        mencache,
		trace:           otel.Tracer("identity-service"),
		logger:          logger,
		token:           token,
		refreshToken:    refreshToken,
		user:            user,
		mapping:         mapping,
		tokenService:    tokenService,
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}
}

func (s *identityService) RefreshToken(token string) (*response.TokenResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("RefreshToken", status, startTime)
	}()

	ctx, span := s.trace.Start(s.ctx, "IdentityService.RefreshToken")
	defer span.End()

	span.SetAttributes(attribute.String("token", token))
	s.logger.Debug("Refreshing token", zap.String("token", token))

	if cachedUserID, found := s.mencache.GetRefreshToken(token); found {
		userId, err := strconv.Atoi(cachedUserID)
		if err == nil {
			s.mencache.DeleteRefreshToken(token)
			s.logger.Debug("Invalidated old refresh token from cache", zap.String("token", token))

			accessToken, err := s.tokenService.createAccessToken(ctx, userId)
			if err != nil {
				return s.errorToken.HandleCreateAccessTokenError(err, "RefreshToken", "CREATE_ACCESS_TOKEN_FAILED", span, &status, zap.Int("user_id", userId))
			}

			refreshToken, err := s.tokenService.createRefreshToken(ctx, userId)
			if err != nil {
				return s.errorToken.HandleCreateRefreshTokenError(err, "RefreshToken", "CREATE_REFRESH_TOKEN_FAILED", span, &status, zap.Int("user_id", userId))
			}

			expiryTime := time.Now().Add(24 * time.Hour)
			expirationDuration := time.Until(expiryTime)

			s.mencache.SetRefreshToken(refreshToken, expirationDuration)

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
			s.mencache.DeleteRefreshToken(token)
			if err := s.refreshToken.DeleteRefreshToken(token); err != nil {
				return s.errorhandler.HandleDeleteRefreshTokenError(err, "RefreshToken", "DELETE_REFRESH_TOKEN", span, &status, zap.String("token", token))
			}
			return s.errorhandler.HandleExpiredRefreshTokenError(err, "RefreshToken", "TOKEN_EXPIRED", span, &status, zap.String("token", token))
		}

		return s.errorhandler.HandleInvalidTokenError(err, "RefreshToken", "INVALID_TOKEN", span, &status, zap.String("token", token))
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return errorhandler.HandleInvalidFormatUserIDError[*response.TokenResponse](s.logger, err, "RefreshToken", "INVALID_USER_ID", span, &status, zap.Int("user_id", userId))
	}

	span.SetAttributes(attribute.Int("user.id", userId))

	s.mencache.DeleteRefreshToken(token)
	if err := s.refreshToken.DeleteRefreshToken(token); err != nil {
		return s.errorhandler.HandleDeleteRefreshTokenError(err, "RefreshToken", "DELETE_REFRESH_TOKEN", span, &status, zap.String("token", token))
	}

	accessToken, err := s.tokenService.createAccessToken(ctx, userId)
	if err != nil {
		return s.errorToken.HandleCreateAccessTokenError(err, "RefreshToken", "CREATE_ACCESS_TOKEN_FAILED", span, &status, zap.Int("user_id", userId))
	}

	refreshToken, err := s.tokenService.createRefreshToken(ctx, userId)
	if err != nil {
		return s.errorToken.HandleCreateRefreshTokenError(err, "RefreshToken", "CREATE_REFRESH_TOKEN_FAILED", span, &status, zap.Int("user_id", userId))
	}

	expiryTime := time.Now().Add(24 * time.Hour)
	updateRequest := &requests.UpdateRefreshToken{
		UserId:    userId,
		Token:     refreshToken,
		ExpiresAt: expiryTime.Format("2006-01-02 15:04:05"),
	}

	if _, err = s.refreshToken.UpdateRefreshToken(updateRequest); err != nil {
		s.mencache.DeleteRefreshToken(refreshToken)
		return s.errorhandler.HandleUpdateRefreshTokenError(err, "RefreshToken", "UPDATE_REFRESH_TOKEN_FAILED", span, &status, zap.Int("user_id", userId))
	}

	expirationDuration := time.Until(expiryTime)
	s.mencache.SetRefreshToken(refreshToken, expirationDuration)
	s.logger.Debug("Stored new refresh token in cache after DB update",
		zap.String("new_token", refreshToken),
		zap.Duration("expiration", expirationDuration))

	s.logger.Debug("Refresh token refreshed successfully", zap.Int("user_id", userId))
	span.SetStatus(codes.Ok, "Token refreshed successfully")

	return &response.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
func (s *identityService) GetMe(token string) (*response.UserResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("GetMe", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "IdentityService.GetMe")
	defer span.End()

	span.SetAttributes(attribute.String("token", token))
	s.logger.Debug("Fetching user details", zap.String("token", token))

	userIdStr, err := s.token.ValidateToken(token)
	if err != nil {
		return s.errorhandler.HandleValidateTokenError(err, "GetMe", "INVALID_TOKEN", span, &status, zap.String("token", token))
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return errorhandler.HandleInvalidFormatUserIDError[*response.UserResponse](s.logger, err, "GetMe", "INVALID_USER_ID", span, &status, zap.Int("user_id", userId))
	}

	if cachedUser, found := s.mencache.GetCachedUserInfo(userIdStr); found {
		s.logger.Debug("User info retrieved from cache", zap.Int("user_id", userId))
		span.SetStatus(codes.Ok, "User details fetched from cache")
		return cachedUser, nil
	}

	span.SetAttributes(attribute.Int("user.id", userId))

	user, err := s.user.FindById(userId)
	if err != nil {
		return s.errorhandler.HandleFindByIdError(err, "GetMe", "FAILED_FETCH_USER", span, &status, zap.Int("user_id", userId))
	}

	userResponse := s.mapping.ToUserResponse(user)

	s.mencache.SetCachedUserInfo(userResponse, time.Minute*5)

	s.logger.Debug("User details fetched successfully", zap.Int("user_id", userId))
	span.SetStatus(codes.Ok, "User details fetched")

	return userResponse, nil
}

func (s *identityService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
