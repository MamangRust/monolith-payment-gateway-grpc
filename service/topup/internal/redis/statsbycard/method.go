package topupstatsbycardcache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type topupStatsMethodByCardCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTopupStatsMethodByCardCache(store *sharedcachehelpers.CacheStore) TopupStatsMethodByCardCache {
	return &topupStatsMethodByCardCache{store: store}
}

// GetMonthlyTopupMethodsByCardNumberCache retrieves the monthly topup methods data
// for a specific card number from the cache. It takes as an argument a request containing
// the card number, year, and month, and returns a slice of TopupMonthMethodResponse and a
// boolean indicating whether the data was found in the cache.
func (s *topupStatsMethodByCardCache) GetMonthlyTopupMethodsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod) ([]*response.TopupMonthMethodResponse, bool) {
	key := fmt.Sprintf(monthTopupMethodByCardCacheKey, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.TopupMonthMethodResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthlyTopupMethodsByCardNumberCache stores the monthly topup methods data
// for a specific card number in the cache. It takes a request containing the
// card number, year, and month, and a slice of TopupMonthMethodResponse as arguments.
//
// If the provided data is nil, it returns immediately.
//
// It constructs a cache key using the request's parameters and stores the data in the
// cache with a default TTL.
func (s *topupStatsMethodByCardCache) SetMonthlyTopupMethodsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod, data []*response.TopupMonthMethodResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupMethodByCardCacheKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

// GetYearlyTopupMethodsByCardNumberCache retrieves the yearly topup methods data
// for a specific card number from the cache. It takes as an argument a request containing
// the card number, year, and month, and returns a slice of TopupYearlyMethodResponse and a
// boolean indicating whether the data was found in the cache.
func (s *topupStatsMethodByCardCache) GetYearlyTopupMethodsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod) ([]*response.TopupYearlyMethodResponse, bool) {
	key := fmt.Sprintf(yearTopupMethodByCardCacheKey, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.TopupYearlyMethodResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyTopupMethodsByCardNumberCache stores the yearly topup methods data
// for a specific card number in the cache. It takes a request containing the
// card number and year, and a slice of TopupYearlyMethodResponse as arguments.
//
// If the provided data is nil, it returns immediately.
//
// It constructs a cache key using the request's parameters and stores the data
// in the cache with a default TTL.
func (s *topupStatsMethodByCardCache) SetYearlyTopupMethodsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod, data []*response.TopupYearlyMethodResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupMethodByCardCacheKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, s.store, key, &data, ttlDefault)
}
