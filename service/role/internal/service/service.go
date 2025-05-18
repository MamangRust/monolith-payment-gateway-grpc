package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-role/internal/repository"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
)

type Service struct {
	RoleQuery   RoleQueryService
	RoleCommand RoleCommandService
}

type Deps struct {
	Ctx          context.Context
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps Deps) *Service {
	roleMapper := responseservice.NewRoleResponseMapper()

	return &Service{
		RoleQuery:   NewRoleQueryService(deps.Ctx, deps.Repositories.RoleQuery, deps.Logger, roleMapper),
		RoleCommand: NewRoleCommandService(deps.Ctx, deps.Repositories.RoleCommand, deps.Logger, roleMapper),
	}
}
