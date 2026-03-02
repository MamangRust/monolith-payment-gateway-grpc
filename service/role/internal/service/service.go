package service

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	mencache "github.com/MamangRust/monolith-payment-gateway-role/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-role/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
)

// Service aggregates role-related services.
type Service struct {
	RoleQuery   RoleQueryService
	RoleCommand RoleCommandService
}

// Deps defines dependencies for role services.
type Deps struct {
	Cache        *cache.CacheStore
	Repositories repository.Repositories
	Logger       logger.LoggerInterface
}

// NewService creates a new role Service.
func NewService(deps *Deps) *Service {
	cache := mencache.NewMencache(deps.Cache)

	return &Service{
		RoleQuery:   newRoleQueryService(deps, cache.RoleQueryCache),
		RoleCommand: newRoleCommandService(deps, cache.RoleCommandCache),
	}
}

// newRoleCommandService creates a RoleCommandService.
func newRoleCommandService(
	deps *Deps,
	cache mencache.RoleCommandCache,
) RoleCommandService {
	return NewRoleCommandService(&roleCommandDeps{
		Cache:      cache,
		Repository: deps.Repositories,
		Logger:     deps.Logger,
	})
}

// newRoleQueryService creates a RoleQueryService.
func newRoleQueryService(
	deps *Deps,
	cache mencache.RoleQueryCache,
) RoleQueryService {
	return NewRoleQueryService(&roleQueryDeps{
		Cache:      cache,
		Repository: deps.Repositories,
		Logger:     deps.Logger,
	})
}
