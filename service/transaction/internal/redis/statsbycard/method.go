package transactionstatsbycarcache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type transactionStatsByCardMethodCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTransactionStatsByCardMethodCache(store *sharedcachehelpers.CacheStore) TransactionStatsByCardMethodCache {
	return &transactionStatsByCardMethodCache{store: store}
}

// GetMonthlyPaymentMethodsByCardCache retrieves cached monthly payment method statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Month, year, and card number request.
//
// Returns:
//   - []*response.TransactionMonthMethodResponse: Payment methods per month.
//   - bool: Whether the cache was found.
func (t *transactionStatsByCardMethodCache) GetMonthlyPaymentMethodsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*response.TransactionMonthMethodResponse, bool) {
	key := fmt.Sprintf(monthTopupMethodByCardCacheKey, req.CardNumber, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.TransactionMonthMethodResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

// SetMonthlyPaymentMethodsByCardCache stores monthly payment method statistics by card number into cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request parameters.
//   - data: Monthly method stats.
func (t *transactionStatsByCardMethodCache) SetMonthlyPaymentMethodsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod, data []*response.TransactionMonthMethodResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupMethodByCardCacheKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

// GetYearlyPaymentMethodsByCardCache retrieves cached yearly payment method statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request object with card and year info.
//
// Returns:
//   - []*response.TransactionYearMethodResponse: Yearly method stats.
//   - bool: Whether the cache was found.
func (t *transactionStatsByCardMethodCache) GetYearlyPaymentMethodsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*response.TransactionYearMethodResponse, bool) {
	key := fmt.Sprintf(yearTopupMethodByCardCacheKey, req.CardNumber, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.TransactionYearMethodResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyPaymentMethodsByCardCache stores yearly payment method statistics by card number into cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request details.
//   - data: The method stats to cache.
func (t *transactionStatsByCardMethodCache) SetYearlyPaymentMethodsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod, data []*response.TransactionYearMethodResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupMethodByCardCacheKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}
