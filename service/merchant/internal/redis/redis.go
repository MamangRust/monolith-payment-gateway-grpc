package mencache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Mencache struct {
	MerchantQueryCache               MerchantQueryCache
	MerchantCommandCache             MerchantCommandCache
	MerchantDocumentQueryCache       MerchantDocumentQueryCache
	MerchantDocumentCommandCache     MerchantDocumentCommandCache
	MerchantTransactionCache         MerchantTransactionCache
	MerchantStatisticCache           MerchantStatisticCache
	MerchantStatisticByMerchantCache MerchantStatisticByMerchantCache
	MerchantStatisticByApiCache      MerchantStatisticByApikeyCache
}

type Deps struct {
	Ctx    context.Context
	Redis  *redis.Client
	Logger logger.LoggerInterface
}

func NewMencache(deps *Deps) *Mencache {
	cacheStore := NewCacheStore(deps.Ctx, deps.Redis, deps.Logger)

	return &Mencache{
		MerchantQueryCache:               NewMerchantQueryCache(cacheStore),
		MerchantCommandCache:             NewMerchantCommandCache(cacheStore),
		MerchantDocumentQueryCache:       NewMerchantDocumentQueryCache(cacheStore),
		MerchantDocumentCommandCache:     NewMerchantDocumentCommandCache(cacheStore),
		MerchantTransactionCache:         NewMerchantTransactionCachhe(cacheStore),
		MerchantStatisticCache:           NewMerchantStatisticCache(cacheStore),
		MerchantStatisticByMerchantCache: NewMerchantStatisticByMerchantCache(cacheStore),
		MerchantStatisticByApiCache:      NewMerchantStatisticByApiKeyCache(cacheStore),
	}
}
