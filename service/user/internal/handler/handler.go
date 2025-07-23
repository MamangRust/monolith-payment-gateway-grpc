package handler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/user"
	"github.com/MamangRust/monolith-payment-gateway-user/internal/service"
)

// Deps holds dependencies required to initialize HTTP API handlers.
type Deps struct {
	// Service is a service instance providing user query, command, and statistic functionalities.
	Service service.Service
	// Logger is a logger interface for structured logging.
	Logger logger.LoggerInterface
}

// Handler holds gRPC handlers for user operations.
type Handler interface {
	UserQueryHandleGrpc
	UserCommandHandleGrpc
}

type handler struct {
	UserQueryHandleGrpc
	UserCommandHandleGrpc
}

// NewHandler creates a new Handler instance.
//
// It takes a pointer to a Deps struct as argument, which contains the dependencies
// required to set up the handler.
//
// The returned handler is ready to be used.
func NewHandler(deps *Deps) Handler {
	mapper := protomapper.NewUserProtoMapper()

	return &handler{
		UserQueryHandleGrpc:   NewUserQueryHandleGrpc(deps.Service, deps.Logger, mapper.UserQueryProtoMapper),
		UserCommandHandleGrpc: NewUserCommandHandleGrpc(deps.Service, deps.Logger, mapper.UserCommandProtoMapper),
	}
}
