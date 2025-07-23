package authhandler

import (
	pb "github.com/MamangRust/monolith-payment-gateway-pb/auth"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/auth"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsAuth struct {
	Client *grpc.ClientConn
	E      *echo.Echo
	Logger logger.LoggerInterface
}

func RegisterAuthHandler(deps *DepsAuth) {
	mapper := apimapper.NewAuthResponseMapper()

	NewHandlerAuth(&authHandleParams{
		client: pb.NewAuthServiceClient(deps.Client),
		router: deps.E,
		logger: deps.Logger,
		mapper: mapper,
	})
}
