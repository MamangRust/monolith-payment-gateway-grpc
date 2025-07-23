package rolehandler

import (
	mencache "github.com/MamangRust/monolith-payment-gateway-apigateway/internal/redis"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/role"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/role"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsRole struct {
	Client *grpc.ClientConn

	// Kafka is the Kafka instance used for producing and consuming messages.
	Kafka *kafka.Kafka

	E *echo.Echo

	Logger logger.LoggerInterface

	Cache mencache.RoleCache
}

func RegisterRoleHandler(deps *DepsRole) {
	mapper := apimapper.NewRoleResponseMapper()

	handlers := []func(){
		setupRoleQueryHandler(deps, deps.Cache, mapper.QueryMapper()),
		setupRoleCommandHandler(deps, deps.Cache, mapper.CommandMapper()),
	}

	for _, h := range handlers {
		h()
	}
}

func setupRoleQueryHandler(deps *DepsRole, cache mencache.RoleCache, mapper apimapper.RoleQueryResponseMapper) func() {
	return func() {
		NewRoleQueryHandleApi(&roleQueryHandleDeps{
			client: pb.NewRoleQueryServiceClient(deps.Client),
			router: deps.E,
			logger: deps.Logger,
			mapper: mapper,
			cache:  cache,
			kafka:  deps.Kafka,
		})
	}
}

func setupRoleCommandHandler(deps *DepsRole, cache mencache.RoleCache, mapper apimapper.RoleCommandResponseMapper) func() {
	return func() {
		NewRoleCommandHandleApi(&roleCommandHandleDeps{
			client: pb.NewRoleCommandServiceClient(deps.Client),
			router: deps.E,
			logger: deps.Logger,
			mapper: mapper,
			kafka:  deps.Kafka,
			cache:  cache,
		})
	}
}
