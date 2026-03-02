package service

import (
	"context"
	"time"

	mencache "github.com/MamangRust/monolith-payment-gateway-auth/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-auth/internal/repository"

	"github.com/MamangRust/monolith-payment-gateway-pkg/auth"
	"github.com/MamangRust/monolith-payment-gateway-pkg/hash"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	sharederrorhandler "github.com/MamangRust/monolith-payment-gateway-shared/errorhandler"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// LoginServiceDeps defines all dependencies required by LoginService.
type LoginServiceDeps struct {
	Cache  mencache.LoginCache
	Logger logger.LoggerInterface
	Hash   hash.HashPassword

	UserRepository repository.UserRepository
	RefreshToken   repository.RefreshTokenRepository

	Token        auth.TokenManager
	TokenService *tokenService

	Observability observability.TraceLoggerObservability
}

type loginService struct {
	mencache mencache.LoginCache
	logger   logger.LoggerInterface
	hash     hash.HashPassword

	user         repository.UserRepository
	refreshToken repository.RefreshTokenRepository

	token        auth.TokenManager
	tokenService *tokenService

	observability observability.TraceLoggerObservability
}

func NewLoginService(params *LoginServiceDeps) *loginService {
	return &loginService{
		mencache:      params.Cache,
		logger:        params.Logger,
		hash:          params.Hash,
		user:          params.UserRepository,
		refreshToken:  params.RefreshToken,
		token:         params.Token,
		tokenService:  params.TokenService,
		observability: params.Observability,
	}
}

func (s *loginService) Login(ctx context.Context, request *requests.AuthRequest) (*response.TokenResponse, error) {
	const method = "Login"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("email", request.Email))

	defer func() {
		end(status)
	}()

	if cachedToken, found := s.mencache.GetCachedLogin(ctx, request.Email); found {
		logSuccess("Successfully logged in from cache", zap.String("email", request.Email))
		return cachedToken, nil
	}

	res, err := s.user.FindByEmailAndVerify(ctx, request.Email)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*response.TokenResponse](s.logger, err, method, span, zap.String("email", request.Email))
	}

	err = s.hash.ComparePassword(res.Password, request.Password)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*response.TokenResponse](s.logger, err, method, span, zap.String("email", request.Email))
	}

	token, err := s.tokenService.createAccessToken(int(res.UserID))
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*response.TokenResponse](s.logger, err, method, span, zap.Int("user.id", int(res.UserID)))
	}

	refreshToken, err := s.tokenService.createRefreshToken(ctx, int(res.UserID))
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*response.TokenResponse](s.logger, err, method, span, zap.Int("user.id", int(res.UserID)))
	}

	tokenResp := &response.TokenResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}

	s.mencache.SetCachedLogin(ctx, request.Email, tokenResp, time.Minute)

	logSuccess("Successfully logged in", zap.String("email", request.Email))

	return tokenResp, nil
}
