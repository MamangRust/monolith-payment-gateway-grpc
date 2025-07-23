package mencache

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	transactionstatscache "github.com/MamangRust/monolith-payment-gateway-transaction/internal/redis/stats"
	transactionstatsbycarcache "github.com/MamangRust/monolith-payment-gateway-transaction/internal/redis/statsbycard"
	"github.com/redis/go-redis/v9"
)

type Mencache interface {
	TransactionQueryCache
	TransactionCommandCache
	transactionstatscache.TransactionStatsCache
	transactionstatsbycarcache.TransactionStatsByCardCache
}

// Mencache represents a cache store for transaction queries, commands, statistics, and statistics by card.
type mencache struct {
	TransactionQueryCache
	TransactionCommandCache
	transactionstatscache.TransactionStatsCache
	transactionstatsbycarcache.TransactionStatsByCardCache
}

// Deps represents the dependencies required to initialize a Mencache.
type Deps struct {
	Redis  *redis.Client
	Logger logger.LoggerInterface
}

// NewMencache creates a new Mencache instance using the given dependencies.
// It creates a new cache store using the given context, Redis client, and logger,
// and returns a Mencache struct with initialized caches for transaction query, transaction command,
// transaction statistic, and transaction statistic by card.
func NewMencache(deps *Deps) Mencache {
	cacheStore := sharedcachehelpers.NewCacheStore(deps.Redis, deps.Logger)

	return &mencache{
		TransactionQueryCache:       NewTransactionQueryCache(cacheStore),
		TransactionCommandCache:     NewTransactionCommandCache(cacheStore),
		TransactionStatsCache:       transactionstatscache.NewTransactionStatsCache(cacheStore),
		TransactionStatsByCardCache: transactionstatsbycarcache.NewTransactionStatsByCardCache(cacheStore),
	}
}
