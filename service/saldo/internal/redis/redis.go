package mencache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	saldostatscache "github.com/MamangRust/monolith-payment-gateway-saldo/internal/redis/stats"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/redis/go-redis/v9"
)

type Mencache interface {
	SaldoQueryCache
	SaldoCommandCache
	saldostatscache.SaldoStatsCache
}

type mencache struct {
	SaldoQueryCache
	SaldoCommandCache
	saldostatscache.SaldoStatsCache
}

// Deps is a struct that holds the dependencies for creating a new Mencache instance.
type Deps struct {
	Ctx    context.Context
	Redis  *redis.Client
	Logger logger.LoggerInterface
}

// NewMencache creates a new Mencache instance using the given dependencies.
// It creates a new cache store using the given context, redis client, and logger.
// Then it creates a new Mencache instance with the given cache store for saldo query,
// saldo command, and saldo statistic caches.
func NewMencache(deps *Deps) Mencache {
	cacheStore := sharedcachehelpers.NewCacheStore(deps.Redis, deps.Logger)

	return &mencache{
		SaldoQueryCache:   NewSaldoQueryCache(cacheStore),
		SaldoCommandCache: NewSaldoCommandCache(cacheStore),
		SaldoStatsCache:   saldostatscache.NewSaldoStatsCache(cacheStore),
	}
}
