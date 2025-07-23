package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/hash"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/user"
	"github.com/MamangRust/monolith-payment-gateway-user/internal/errorhandler"
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
	Ctx          context.Context
	ErrorHandler *errorhandler.ErrorHandler
	Mencache     mencache.Mencache
	Repositories repository.Repositories
	Hash         hash.HashPassword
	Logger       logger.LoggerInterface
}

// NewService initializes and returns a new instance of the Service struct,
// which provides a suite of user services including query and command services.
// It sets up these services using the provided dependencies and response mapper.
func NewService(deps *Deps) Service {
	userMapper := responseservice.NewUserResponseMapper()

	return &service{
		UserQueryService:   newUserQueryService(deps, userMapper.QueryMapper()),
		UserCommandService: newUserCommandService(deps, userMapper.CommandMapper()),
	}
}

func newUserQueryService(
	deps *Deps,
	mapper responseservice.UserQueryResponseMapper,
) UserQueryService {
	return NewUserQueryService(
		&userQueryDeps{
			Ctx:          deps.Ctx,
			ErrorHandler: deps.ErrorHandler.UserQueryError,
			Cache:        deps.Mencache,
			Repository:   deps.Repositories.UserQuery(),
			Logger:       deps.Logger,
			Mapper:       mapper,
		},
	)
}

func newUserCommandService(
	deps *Deps,
	mapper responseservice.UserCommandResponseMapper,
) UserCommandService {
	return NewUserCommandService(
		&userCommandDeps{
			Ctx:                   deps.Ctx,
			ErrorHandler:          deps.ErrorHandler.UserCommandError,
			Cache:                 deps.Mencache,
			UserQueryRepository:   deps.Repositories.UserQuery(),
			UserCommandRepository: deps.Repositories.UserCommand(),
			RoleRepository:        deps.Repositories.Role(),
			Logger:                deps.Logger,
			Mapper:                mapper,
			Hashing:               deps.Hash,
		},
	)
}
