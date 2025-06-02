package mencache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Mencache struct {
	WithdrawQueryCache           WithdrawQueryCache
	WithdrawCommand              WithdrawCommandCache
	WithdrawStatisticCache       WithdrawStatisticCache
	WithdrawStatisticByCardCache WithdrawStasticByCardCache
}

type Deps struct {
	Ctx    context.Context
	Redis  *redis.Client
	Logger logger.LoggerInterface
}

func NewMencache(deps *Deps) *Mencache {
	cacheStore := NewCacheStore(deps.Ctx, deps.Redis, deps.Logger)

	return &Mencache{
		WithdrawQueryCache:           NewWithdrawQueryCache(cacheStore),
		WithdrawCommand:              NewWithdrawCommandCache(cacheStore),
		WithdrawStatisticCache:       NewWithdrawStatisticCache(cacheStore),
		WithdrawStatisticByCardCache: NewWithdrawStatisticByCardCache(cacheStore),
	}
}
