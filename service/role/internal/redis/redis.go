package mencache

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/redis/go-redis/v9"
)

// Mencache is a struct that holds various cache interfaces
type Mencache struct {
	RoleCommandCache RoleCommandCache
	RoleQueryCache   RoleQueryCache
}

// Deps is a struct that holds dependencies
type Deps struct {
	Redis  *redis.Client
	Logger logger.LoggerInterface
}

func NewMencache(deps *Deps) *Mencache {
	cacheStore := sharedcachehelpers.NewCacheStore(deps.Redis, deps.Logger)

	return &Mencache{
		RoleCommandCache: NewRoleCommandCache(cacheStore),
		RoleQueryCache:   NewRoleQueryCache(cacheStore),
	}
}
