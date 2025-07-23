package mencache

import (
	carddashboardmencache "github.com/MamangRust/monolith-payment-gateway-card/internal/redis/dashboard"
	cardstatsmencache "github.com/MamangRust/monolith-payment-gateway-card/internal/redis/stats"
	cardstatsbycardmencache "github.com/MamangRust/monolith-payment-gateway-card/internal/redis/statsbycard"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/redis/go-redis/v9"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
)

type Mencache interface {
	CardQueryCache
	CardCommandCache
	cardstatsmencache.CardStatsCache
	cardstatsbycardmencache.CardStatsByCardCache
	carddashboardmencache.CardDashboardCache
}

// Mencache is a struct that represents the cache store
type mencache struct {
	CardQueryCache
	CardCommandCache
	cardstatsmencache.CardStatsCache
	cardstatsbycardmencache.CardStatsByCardCache
	carddashboardmencache.CardDashboardCache
}

// Deps is a struct that represents the dependencies needed to create a Mencache
type Deps struct {
	Redis  *redis.Client
	Logger logger.LoggerInterface
}

// NewMencache creates a new Mencache instance using the provided dependencies.
// It initializes a cache store with the given context, Redis client, and logger,
// and returns a Mencache struct with initialized caches for card command, dashboard,
// query, statistic, and statistic by number.
func NewMencache(deps *Deps) Mencache {
	cacheStore := sharedcachehelpers.NewCacheStore(deps.Redis, deps.Logger)

	return &mencache{
		CardCommandCache:     NewCardCommandCache(cacheStore),
		CardQueryCache:       NewCardQueryCache(cacheStore),
		CardStatsCache:       cardstatsmencache.NewMencacheStats(cacheStore),
		CardStatsByCardCache: cardstatsbycardmencache.NewMencacheStatsByCard(cacheStore),
		CardDashboardCache:   carddashboardmencache.NewMencacheDashboard(cacheStore),
	}
}
