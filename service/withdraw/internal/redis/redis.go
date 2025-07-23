package mencache

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	withdrawstatscache "github.com/MamangRust/monolith-payment-gateway-withdraw/internal/redis/stats"
	withdrawstatsbycardcache "github.com/MamangRust/monolith-payment-gateway-withdraw/internal/redis/statsbycard"
	"github.com/redis/go-redis/v9"
)

// Mencache represents a cache store for withdraw operations.
type Mencache interface {
	WithdrawQueryCache
	WithdrawCommandCache
	withdrawstatscache.WithdrawStatsCache
	withdrawstatsbycardcache.WithdrawStatsByCardCache
}

type mencache struct {
	WithdrawQueryCache
	WithdrawCommandCache
	withdrawstatscache.WithdrawStatsCache
	withdrawstatsbycardcache.WithdrawStatsByCardCache
}

// Deps represents the dependencies required for creating a Mencache instance.
type Deps struct {
	Redis  *redis.Client
	Logger logger.LoggerInterface
}

// NewMencache creates a new Mencache instance using the given dependencies.
// It creates a new cache store using the given context, Redis client, and logger,
// and returns a Mencache struct with initialized caches for withdraw query, withdraw command,
// withdraw statistic, and withdraw statistic by card.
func NewMencache(deps *Deps) Mencache {
	cacheStore := sharedcachehelpers.NewCacheStore(deps.Redis, deps.Logger)

	return &mencache{
		WithdrawQueryCache:       NewWithdrawQueryCache(cacheStore),
		WithdrawCommandCache:     NewWithdrawCommandCache(cacheStore),
		WithdrawStatsCache:       withdrawstatscache.NewWithdrawStatsCache(cacheStore),
		WithdrawStatsByCardCache: withdrawstatsbycardcache.NewWithdrawStatsByCardCache(cacheStore),
	}
}
