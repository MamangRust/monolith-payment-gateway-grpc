package mencache

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	transferstatscache "github.com/MamangRust/monolith-payment-gateway-transfer/internal/redis/stats"
	transferstatsbycardcache "github.com/MamangRust/monolith-payment-gateway-transfer/internal/redis/statsbycard"
	"github.com/redis/go-redis/v9"
)

type Mencache interface {
	TransferQueryCache
	TransferCommandCache
	transferstatscache.TransferStatsCache
	transferstatsbycardcache.TransferStatsByCardCache
}

type mencache struct {
	TransferQueryCache
	TransferCommandCache
	transferstatscache.TransferStatsCache
	transferstatsbycardcache.TransferStatsByCardCache
}

// Deps represents the dependencies required to create a new Mencache instance.
type Deps struct {
	Redis  *redis.Client
	Logger logger.LoggerInterface
}

// NewMencache creates a new instance of Mencache using the provided dependencies.
// It creates a new cache store using the given context, Redis client, and logger,
// and returns a Mencache struct with initialized caches for transfer query,
// transfer command, transfer statistic, and transfer statistic by card.
func NewMencache(deps *Deps) Mencache {
	cacheStore := sharedcachehelpers.NewCacheStore(deps.Redis, deps.Logger)

	return &mencache{
		TransferQueryCache:       NewTransferQueryCache(cacheStore),
		TransferCommandCache:     NewTransferCommandCache(cacheStore),
		TransferStatsCache:       transferstatscache.NewTransferStatsCache(cacheStore),
		TransferStatsByCardCache: transferstatsbycardcache.NewTransferStatsByCardCache(cacheStore),
	}
}
