package mencache

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/redis/go-redis/v9"
)

// Mencache is a struct that holds various cache interfaces.
type Mencache struct {
	// IdentityCache is responsible for identity-related caching operations.
	IdentityCache IdentityCache
	// LoginCache handles caching operations related to user login.
	LoginCache LoginCache
	// PasswordResetCache manages caching for password reset tokens and codes.
	PasswordResetCache PasswordResetCache
	// RegisterCache deals with caching operations during user registration.
	RegisterCache RegisterCache
}

// Deps is a struct containing dependencies required for creating a Mencache instance.
type Deps struct {
	// Redis is the Redis client used for accessing the cache.
	Redis *redis.Client
	// Logger is used for logging cache operations and errors.
	Logger logger.LoggerInterface
}

// NewMencache creates a new mencache instance using the given dependencies.
//
// It creates a new cache store using the given context, redis client, and logger.
// Then it creates a new mencache instance with the given cache store for identity,
// login, password reset, and register caches.
func NewMencache(deps *Deps) *Mencache {
	cacheStore := sharedcachehelpers.NewCacheStore(deps.Redis, deps.Logger)

	return &Mencache{
		IdentityCache:      NewidentityCache(cacheStore),
		LoginCache:         NewLoginCache(cacheStore),
		PasswordResetCache: NewPasswordResetCache(cacheStore),
		RegisterCache:      NewRegisterCache(cacheStore),
	}
}
