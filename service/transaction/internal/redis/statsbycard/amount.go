package transactionstatsbycarcache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type transactionStatsByCardAmountCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTransactionStatsByCardAmountCache(store *sharedcachehelpers.CacheStore) TransactionStatsByCardAmountCache {
	return &transactionStatsByCardAmountCache{store: store}
}

// GetMonthlyAmountsByCardCache retrieves cached monthly transaction amount statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Month, year, card number filter.
//
// Returns:
//   - []*response.TransactionMonthAmountResponse: Monthly amounts.
//   - bool: Whether the cache was found.
func (t *transactionStatsByCardAmountCache) GetMonthlyAmountsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*response.TransactionMonthAmountResponse, bool) {
	key := fmt.Sprintf(monthTopupAmountByCardCacheKey, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.TransactionMonthAmountResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true

}

// SetMonthlyAmountsByCardCache stores monthly transaction amount statistics by card number into cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request details.
//   - data: Amounts to cache.
func (t *transactionStatsByCardAmountCache) SetMonthlyAmountsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod, data []*response.TransactionMonthAmountResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupAmountByCardCacheKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

// GetYearlyAmountsByCardCache retrieves cached yearly transaction amount statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Card number and year info.
//
// Returns:
//   - []*response.TransactionYearlyAmountResponse: Yearly amounts.
//   - bool: Whether the cache was found.
func (t *transactionStatsByCardAmountCache) GetYearlyAmountsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*response.TransactionYearlyAmountResponse, bool) {
	key := fmt.Sprintf(yearTopupAmountByCardCacheKey, req.CardNumber, req.Year)

	result, found := sharedcachehelpers.GetFromCache[[]*response.TransactionYearlyAmountResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

// SetYearlyAmountsByCardCache stores yearly transaction amount statistics by card number into cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request filter.
//   - data: Yearly amount data.
func (t *transactionStatsByCardAmountCache) SetYearlyAmountsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod, data []*response.TransactionYearlyAmountResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupAmountByCardCacheKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}
