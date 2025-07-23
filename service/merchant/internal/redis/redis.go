package mencache

import (
	merchantstatscache "github.com/MamangRust/monolith-payment-gateway-merchant/internal/redis/stats"
	merchantstatsapikey "github.com/MamangRust/monolith-payment-gateway-merchant/internal/redis/statsbyapikey"
	merchantstatsbymerchant "github.com/MamangRust/monolith-payment-gateway-merchant/internal/redis/statsbymerchant"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/redis/go-redis/v9"
)

// Mencache is a struct that represents the cache store
type mencache struct {
	MerchantQueryCache
	MerchantCommandCache
	MerchantDocumentQueryCache
	MerchantDocumentCommandCache
	MerchantTransactionCache
	merchantstatscache.MerchantStatsCache
	merchantstatsapikey.MerchantStatsByApiKeyCache
	merchantstatsbymerchant.MerchantStatsByMerchantCache
}

type Mencache interface {
	MerchantQueryCache
	MerchantCommandCache
	MerchantDocumentQueryCache
	MerchantDocumentCommandCache
	MerchantTransactionCache
	merchantstatscache.MerchantStatsCache
	merchantstatsapikey.MerchantStatsByApiKeyCache
	merchantstatsbymerchant.MerchantStatsByMerchantCache
}

// Deps is a struct that represents the dependencies required for creating a Mencache
type Deps struct {
	Redis  *redis.Client
	Logger logger.LoggerInterface
}

// NewMencache creates a new instance of Mencache using the provided dependencies.
// It creates a new cache store using the given Redis client and logger,
// and returns a Mencache struct with initialized caches for merchant query,
// merchant command, merchant document query, merchant document command,
// merchant transaction, merchant statistic, merchant statistic by API key,
// and merchant statistic by merchant.
func NewMencache(deps *Deps) Mencache {
	cacheStore := sharedcachehelpers.NewCacheStore(deps.Redis, deps.Logger)

	return &mencache{
		MerchantQueryCache:           NewMerchantQueryCache(cacheStore),
		MerchantCommandCache:         NewMerchantCommandCache(cacheStore),
		MerchantDocumentQueryCache:   NewMerchantDocumentQueryCache(cacheStore),
		MerchantDocumentCommandCache: NewMerchantDocumentCommandCache(cacheStore),
		MerchantTransactionCache:     NewMerchantTransactionCachhe(cacheStore),
		MerchantStatsCache:           merchantstatscache.NewMerchantStatsCache(cacheStore),
		MerchantStatsByApiKeyCache:   merchantstatsapikey.NewMerchantStatsByApiKeyCache(cacheStore),
		MerchantStatsByMerchantCache: merchantstatsbymerchant.NewMerchantStatsByMerchantCache(cacheStore),
	}
}
