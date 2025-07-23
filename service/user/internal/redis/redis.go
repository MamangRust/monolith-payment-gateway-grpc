package mencache

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/redis/go-redis/v9"
)

type Mencache interface {
	UserQueryCache
	UserCommandCache
}

// Mencache is a struct that holds user query and user command caches.
type mencache struct {
	UserQueryCache
	UserCommandCache
}

// Deps is a struct that holds dependencies for creating a new Mencache instance.
type Deps struct {
	Redis  *redis.Client
	Logger logger.LoggerInterface
}

// NewMencache creates a new instance of Mencache using the given dependencies.
// It creates a new cache store using the given context, redis client, and logger,
// and returns a Mencache struct with initialized caches for user query and user command.
func NewMencache(deps *Deps) Mencache {
	cacheStore := sharedcachehelpers.NewCacheStore(deps.Redis, deps.Logger)

	return &mencache{
		UserQueryCache:   NewUserQueryCache(cacheStore),
		UserCommandCache: NewUserCommandCache(cacheStore),
	}
}
