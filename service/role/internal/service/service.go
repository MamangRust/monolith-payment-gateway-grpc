package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-role/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-role/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-role/internal/repository"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	RoleQuery   RoleQueryService
	RoleCommand RoleCommandService
}

type Deps struct {
	Ctx          context.Context
	ErrorHandler *errorhandler.ErrorHandler
	Mencache     *mencache.Mencache
	Redis        *redis.Client
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps *Deps) *Service {
	roleMapper := responseservice.NewRoleResponseMapper()

	return &Service{
		RoleQuery:   NewRoleQueryService(deps.Ctx, deps.ErrorHandler.RoleQueryError, deps.Mencache.RoleQueryCache, deps.Repositories.RoleQuery, deps.Logger, roleMapper),
		RoleCommand: NewRoleCommandService(deps.Ctx, deps.ErrorHandler.RoleCommandError, deps.Mencache.RoleCommandCache, deps.Repositories.RoleCommand, deps.Logger, roleMapper),
	}
}
