package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/hash"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-user/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-user/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-user/internal/repository"
)

type Service struct {
	UserQuery   UserQueryService
	UserCommand UserCommandService
}

type Deps struct {
	Ctx          context.Context
	ErrorHandler *errorhandler.ErrorHandler
	Mencache     *mencache.Mencache
	Repositories *repository.Repositories
	Hash         hash.HashPassword
	Logger       logger.LoggerInterface
}

func NewService(deps *Deps) *Service {
	userMapper := responseservice.NewUserResponseMapper()

	return &Service{
		UserQuery:   NewUserQueryService(deps.Ctx, deps.ErrorHandler.UserQueryError, deps.Mencache.UserQueryCache, deps.Repositories.UserQuery, deps.Logger, userMapper),
		UserCommand: NewUserCommandService(deps.Ctx, deps.ErrorHandler.UserCommandError, deps.Mencache.UserCommandCache, deps.Repositories.UserQuery, deps.Repositories.UserCommand, deps.Repositories.Role, deps.Logger, userMapper, deps.Hash),
	}
}
