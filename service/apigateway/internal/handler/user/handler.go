package userhandler

import (
	pb "github.com/MamangRust/monolith-payment-gateway-pb/user"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/user"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsUser struct {
	Client *grpc.ClientConn

	E *echo.Echo

	Logger logger.LoggerInterface
}

func RegisterUserHandler(deps *DepsUser) {
	mapper := apimapper.NewUserResponseMapper()

	handlers := []func(){
		setupUserQueryHandler(deps, mapper.QueryMapper()),
		setupUserCommandHandler(deps, mapper.CommandMapper()),
	}

	for _, h := range handlers {
		h()
	}
}

func setupUserQueryHandler(deps *DepsUser, mapper apimapper.UserQueryResponseMapper) func() {
	return func() {
		NewUserQueryHandleApi(&userQueryHandleDeps{
			client: pb.NewUserQueryServiceClient(deps.Client),
			router: deps.E,
			logger: deps.Logger,
			mapper: mapper,
		})
	}
}

func setupUserCommandHandler(deps *DepsUser, mapper apimapper.UserCommandResponseMapper) func() {
	return func() {
		NewUserCommandHandleApi(&userCommandHandleDeps{
			client: pb.NewUserCommandServiceClient(deps.Client),
			router: deps.E,
			logger: deps.Logger,
			mapper: mapper,
		})
	}
}
