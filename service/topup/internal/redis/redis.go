package mencache

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	topupstatscache "github.com/MamangRust/monolith-payment-gateway-topup/internal/redis/stats"
	topupstatsbycardcache "github.com/MamangRust/monolith-payment-gateway-topup/internal/redis/statsbycard"
	"github.com/redis/go-redis/v9"
)

type Mencache interface {
	TopupQueryCache
	TopupCommandCache
	topupstatscache.TopupStatsCache
	topupstatsbycardcache.TopupStatsByCardCache
}

type mencache struct {
	TopupQueryCache
	TopupCommandCache
	topupstatscache.TopupStatsCache
	topupstatsbycardcache.TopupStatsByCardCache
}

// Deps is a struct that holds the dependencies required by the mencache
type Deps struct {
	Redis  *redis.Client
	Logger logger.LoggerInterface
}

// NewMencache creates a new Mencache instance using the given dependencies.
// It creates a new cache store using the given context, Redis client, and logger,
// and returns a Mencache struct with initialized caches for topup query, topup command,
// topup statistic, and topup statistic by card.
func NewMencache(deps *Deps) Mencache {
	cacheStore := sharedcachehelpers.NewCacheStore(deps.Redis, deps.Logger)

	return &mencache{
		TopupQueryCache:       NewTopupQueryCache(cacheStore),
		TopupCommandCache:     NewTopupCommandCache(cacheStore),
		TopupStatsCache:       topupstatscache.NewTopupStatsCache(cacheStore),
		TopupStatsByCardCache: topupstatsbycardcache.NewTopupStatsByCardCache(cacheStore),
	}
}
