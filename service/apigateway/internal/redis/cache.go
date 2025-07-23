package mencache

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/redis/go-redis/v9"
)

type CacheApiGateway interface {
	MerchantCache
	RoleCache
}

type mencacheApiGateay struct {
	MerchantCache
	RoleCache
}

type Deps struct {
	Redis  *redis.Client
	Logger logger.LoggerInterface
}

func NewCacheApiGateway(deps *Deps) CacheApiGateway {
	store := sharedcachehelpers.NewCacheStore(deps.Redis, deps.Logger)

	return &mencacheApiGateay{
		MerchantCache: NewMerchantCache(store),
		RoleCache:     NewRoleCache(store),
	}
}
