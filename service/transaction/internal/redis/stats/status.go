package transactionstatscache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type transactionStatsStatusCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTransactionStatsStatusCache(store *sharedcachehelpers.CacheStore) TransactionStatsStatusCache {
	return &transactionStatsStatusCache{store: store}
}

// GetMonthTransactionStatusSuccessCache retrieves cached monthly successful transactions.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request containing month and year filter.
//
// Returns:
//   - []*response.TransactionResponseMonthStatusSuccess: List of successful monthly transactions.
//   - bool: Whether the cache was found.
func (t *transactionStatsStatusCache) GetMonthTransactionStatusSuccessCache(ctx context.Context, req *requests.MonthStatusTransaction) ([]*response.TransactionResponseMonthStatusSuccess, bool) {
	key := fmt.Sprintf(monthTopupStatusSuccessCacheKey, req.Month, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.TransactionResponseMonthStatusSuccess](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthTransactionStatusSuccessCache stores successful monthly transactions in cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Original request object.
//   - data: Transactions to cache.
func (t *transactionStatsStatusCache) SetMonthTransactionStatusSuccessCache(ctx context.Context, req *requests.MonthStatusTransaction, data []*response.TransactionResponseMonthStatusSuccess) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupStatusSuccessCacheKey, req.Month, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

// GetYearTransactionStatusSuccessCache retrieves cached yearly successful transactions.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which statistics are requested.
//
// Returns:
//   - []*response.TransactionResponseYearStatusSuccess: List of successful yearly transactions.
//   - bool: Whether the cache was found.
func (t *transactionStatsStatusCache) GetYearTransactionStatusSuccessCache(ctx context.Context, year int) ([]*response.TransactionResponseYearStatusSuccess, bool) {
	key := fmt.Sprintf(yearTopupStatusSuccessCacheKey, year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.TransactionResponseYearStatusSuccess](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearTransactionStatusSuccessCache caches yearly successful transaction statistics.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year used as key.
//   - data: Yearly successful transaction stats.
func (t *transactionStatsStatusCache) SetYearTransactionStatusSuccessCache(ctx context.Context, year int, data []*response.TransactionResponseYearStatusSuccess) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupStatusSuccessCacheKey, year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

// GetMonthTransactionStatusFailedCache retrieves cached monthly failed transactions.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request containing month and year filter.
//
// Returns:
//   - []*response.TransactionResponseMonthStatusFailed: List of failed monthly transactions.
//   - bool: Whether the cache was found.
func (t *transactionStatsStatusCache) GetMonthTransactionStatusFailedCache(ctx context.Context, req *requests.MonthStatusTransaction) ([]*response.TransactionResponseMonthStatusFailed, bool) {
	key := fmt.Sprintf(monthTopupStatusFailedCacheKey, req.Month, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.TransactionResponseMonthStatusFailed](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthTransactionStatusFailedCache caches monthly failed transaction statistics.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Original request object.
//   - data: List of failed monthly transactions.
func (t *transactionStatsStatusCache) SetMonthTransactionStatusFailedCache(ctx context.Context, req *requests.MonthStatusTransaction, data []*response.TransactionResponseMonthStatusFailed) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupStatusFailedCacheKey, req.Month, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

// GetYearTransactionStatusFailedCache retrieves cached yearly failed transactions.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*response.TransactionResponseYearStatusFailed: List of failed transactions.
//   - bool: Whether the cache was found.
func (t *transactionStatsStatusCache) GetYearTransactionStatusFailedCache(ctx context.Context, year int) ([]*response.TransactionResponseYearStatusFailed, bool) {
	key := fmt.Sprintf(yearTopupStatusFailedCacheKey, year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.TransactionResponseYearStatusFailed](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearTransactionStatusFailedCache caches yearly failed transaction statistics.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year used as key.
//   - data: List of failed yearly transactions.
func (t *transactionStatsStatusCache) SetYearTransactionStatusFailedCache(ctx context.Context, year int, data []*response.TransactionResponseYearStatusFailed) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupStatusFailedCacheKey, year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}
