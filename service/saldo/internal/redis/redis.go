package mencache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Mencache struct {
	SaldoQueryCache     SaldoQueryCache
	SaldoCommandCache   SaldoCommandCache
	SaldoStatisticCache SaldoStatisticCache
}

type Deps struct {
	Ctx    context.Context
	Redis  *redis.Client
	Logger logger.LoggerInterface
}

func NewMencache(deps *Deps) *Mencache {
	cacheStore := NewCacheStore(deps.Ctx, deps.Redis, deps.Logger)

	return &Mencache{
		SaldoQueryCache:     NewSaldoQueryCache(cacheStore),
		SaldoCommandCache:   NewSaldoCommandCache(cacheStore),
		SaldoStatisticCache: NewSaldoStatisticCache(cacheStore),
	}
}
