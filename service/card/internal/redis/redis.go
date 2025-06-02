package mencache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Mencache struct {
	CardCommandCache           CardCommandCache
	CardDashboardCache         CardDashboardCache
	CardQueryCache             CardQueryCache
	CardStatisticCache         CardStatisticCache
	CardStatisticByNumberCache CardStatisticByNumberCache
}

type Deps struct {
	Ctx    context.Context
	Redis  *redis.Client
	Logger logger.LoggerInterface
}

func NewMencache(deps *Deps) *Mencache {
	cacheStore := NewCacheStore(deps.Ctx, deps.Redis, deps.Logger)

	return &Mencache{
		CardCommandCache:           NewCardCommandCache(cacheStore),
		CardDashboardCache:         NewCardDashboardCache(cacheStore),
		CardQueryCache:             NewCardQueryCache(cacheStore),
		CardStatisticCache:         NewCardStatisticCache(cacheStore),
		CardStatisticByNumberCache: NewCardStatisticByNumberCache(cacheStore),
	}
}
