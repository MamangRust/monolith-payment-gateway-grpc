package rolehandler

import (
	mencache "github.com/MamangRust/monolith-payment-gateway-apigateway/internal/redis"
	role_cache "github.com/MamangRust/monolith-payment-gateway-apigateway/internal/redis/api/role"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/role"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/role"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsRole struct {
	Client *grpc.ClientConn

	Kafka *kafka.Kafka

	E *echo.Echo

	Logger logger.LoggerInterface

	Cache mencache.RoleCache

	CacheStore *cache.CacheStore

	ApiHandler errors.ApiHandler
}

func RegisterRoleHandler(deps *DepsRole) {
	mapper := apimapper.NewRoleResponseMapper()
	cache := role_cache.NewRoleMencache(deps.CacheStore)

	handlers := []func(){
		setupRoleQueryHandler(deps, deps.Cache, mapper.QueryMapper(), cache),
		setupRoleCommandHandler(deps, deps.Cache, mapper.CommandMapper(), cache),
	}

	for _, h := range handlers {
		h()
	}
}

func setupRoleQueryHandler(deps *DepsRole, cache_role mencache.RoleCache, mapper apimapper.RoleQueryResponseMapper, cache role_cache.RoleMencache) func() {
	return func() {
		NewRoleQueryHandleApi(&roleQueryHandleDeps{
			client:     pb.NewRoleServiceClient(deps.Client),
			router:     deps.E,
			logger:     deps.Logger,
			mapper:     mapper,
			cache_role: cache_role,
			kafka:      deps.Kafka,
			cache:      cache,
			apiHandler: deps.ApiHandler,
		})
	}
}

func setupRoleCommandHandler(deps *DepsRole, cache_role mencache.RoleCache, mapper apimapper.RoleCommandResponseMapper, cache role_cache.RoleMencache) func() {
	return func() {
		NewRoleCommandHandleApi(&roleCommandHandleDeps{
			client:     pb.NewRoleCommandServiceClient(deps.Client),
			router:     deps.E,
			logger:     deps.Logger,
			mapper:     mapper,
			kafka:      deps.Kafka,
			cache_role: cache_role,
			cache:      cache,
			apiHandler: deps.ApiHandler,
		})
	}
}
