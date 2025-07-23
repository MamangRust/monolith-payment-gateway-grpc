package handler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-role/internal/service"
	roleprotomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/role"
)

// Deps is a struct that holds the dependencies for the role service.
type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

// Handler is a struct that holds the dependencies for the role service.
type Handler struct {
	RoleQuery   RoleQueryHandlerGrpc
	RoleCommand RoleCommandHandlerGrpc
}

// NewHandler sets up the handler for the role service.
//
// It takes a pointer to a Deps struct as argument, which contains the dependencies
// required to set up the handler.
//
// The returned handler is ready to be used.
func NewHandler(deps *Deps) *Handler {
	mapper := roleprotomapper.NewRoleProtoMapper()

	return &Handler{
		RoleQuery:   NewRoleQueryHandleGrpc(deps.Service.RoleQuery, deps.Logger, mapper.RoleQueryProtoMapper),
		RoleCommand: NewRoleCommandHandleGrpc(deps.Service.RoleCommand, deps.Logger, mapper.RoleCommandProtoMapper),
	}
}
