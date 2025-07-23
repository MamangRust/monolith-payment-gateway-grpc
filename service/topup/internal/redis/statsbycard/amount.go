package topupstatsbycardcache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type topupStatsAmountByCardCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTopupStatsAmountByCardCache(store *sharedcachehelpers.CacheStore) TopupStatsAmountByCardCache {
	return &topupStatsAmountByCardCache{store: store}
}

// GetMonthlyTopupAmountsByCardNumberCache retrieves the monthly topup amount data
// for a specific card number from the cache. It takes as an argument a request containing
// the card number and year, and returns a slice of TopupMonthAmountResponse and a boolean
// indicating whether the data was found in the cache.
func (s *topupStatsAmountByCardCache) GetMonthlyTopupAmountsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod) ([]*response.TopupMonthAmountResponse, bool) {
	key := fmt.Sprintf(monthTopupAmountByCardCacheKey, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.TopupMonthAmountResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthlyTopupAmountsByCardNumberCache stores the monthly topup amounts for a specific card number
// in the cache. It takes a request containing the card number, year, and month, and a slice of
// TopupMonthAmountResponse as arguments.
//
// If the provided data is nil, it returns immediately.
//
// It constructs a cache key using the request's parameters and stores the data in the cache
// with a default TTL.
func (s *topupStatsAmountByCardCache) SetMonthlyTopupAmountsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod, data []*response.TopupMonthAmountResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupAmountByCardCacheKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

// GetYearlyTopupAmountsByCardNumberCache retrieves the yearly topup amount data
// for a specific card number from the cache. It takes as an argument a request containing
// the card number and year, and returns a slice of TopupYearlyAmountResponse and a boolean
// indicating whether the data was found in the cache.
func (s *topupStatsAmountByCardCache) GetYearlyTopupAmountsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod) ([]*response.TopupYearlyAmountResponse, bool) {
	key := fmt.Sprintf(yearTopupAmountByCardCacheKey, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.TopupYearlyAmountResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyTopupAmountsByCardNumberCache stores the yearly topup amounts for a specific card number
// in the cache. It takes a request containing the card number, year, and month, and a slice of
// TopupYearlyAmountResponse as arguments.
//
// If the provided data is nil, it returns immediately.
//
// It constructs a cache key using the request's parameters and stores the data in the cache
// with a default TTL.
func (s *topupStatsAmountByCardCache) SetYearlyTopupAmountsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod, data []*response.TopupYearlyAmountResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupAmountByCardCacheKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, s.store, key, &data, ttlDefault)
}
