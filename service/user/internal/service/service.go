package service

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/hash"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	mencache "github.com/MamangRust/monolith-payment-gateway-user/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-user/internal/repository"
)

type Service interface {
	UserQueryService
	UserCommandService
}

type service struct {
	UserQueryService
	UserCommandService
}

// Deps represents the dependencies required by the Service struct.
type Deps struct {
	Cache        *cache.CacheStore
	Repositories repository.Repositories
	Hash         hash.HashPassword
	Logger       logger.LoggerInterface
}

// NewService initializes and returns a new instance of the Service struct,
// which provides a suite of user services including query and command services.
// It sets up these services using the provided dependencies and response mapper.
func NewService(deps *Deps) Service {
	cache := mencache.NewMencache(deps.Cache)

	return &service{
		UserQueryService:   newUserQueryService(deps, cache),
		UserCommandService: newUserCommandService(deps, cache),
	}
}

func newUserQueryService(
	deps *Deps,
	cache mencache.UserQueryCache,
) UserQueryService {
	return NewUserQueryService(
		&userQueryDeps{
			Cache:      cache,
			Repository: deps.Repositories.UserQuery(),
			Logger:     deps.Logger,
		},
	)
}

func newUserCommandService(
	deps *Deps,
	cache mencache.UserCommandCache,
) UserCommandService {
	return NewUserCommandService(
		&userCommandDeps{
			Cache:                 cache,
			UserQueryRepository:   deps.Repositories.UserQuery(),
			UserCommandRepository: deps.Repositories.UserCommand(),
			RoleRepository:        deps.Repositories.Role(),
			Logger:                deps.Logger,
			Hashing:               deps.Hash,
		},
	)
}
