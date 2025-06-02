package mencache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Mencache struct {
	TransferQueryCache           TransferQueryCache
	TransferCommandCache         TransferCommandCache
	TransferStatisticCache       TransferStatisticCache
	TransferStatisticByCardCache TransferStatisticByCardCache
}

type Deps struct {
	Ctx    context.Context
	Redis  *redis.Client
	Logger logger.LoggerInterface
}

func NewMencache(deps *Deps) *Mencache {
	cacheStore := NewCacheStore(deps.Ctx, deps.Redis, deps.Logger)

	return &Mencache{
		TransferQueryCache:           NewTransferQueryCache(cacheStore),
		TransferCommandCache:         NewTransferCommandCache(cacheStore),
		TransferStatisticCache:       NewTransferStatisticCache(cacheStore),
		TransferStatisticByCardCache: NewTransferStatisticByCardCache(cacheStore),
	}
}
