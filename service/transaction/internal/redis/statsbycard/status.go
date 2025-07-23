package transactionstatsbycarcache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type transactionStatsByCardStatusCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTransactionStatsByCardStatusCache(store *sharedcachehelpers.CacheStore) TransactionStatsByCardStatusCache {
	return &transactionStatsByCardStatusCache{store: store}
}

// GetMonthTransactionStatusSuccessByCardCache retrieves cached monthly successful transaction statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request containing card number, month, and year.
//
// Returns:
//   - []*response.TransactionResponseMonthStatusSuccess: List of successful transactions.
//   - bool: Whether the cache was found.
func (t *transactionStatsByCardStatusCache) GetMonthTransactionStatusSuccessByCardCache(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*response.TransactionResponseMonthStatusSuccess, bool) {
	key := fmt.Sprintf(monthTopupStatusSuccessByCardCacheKey, req.CardNumber, req.Month, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.TransactionResponseMonthStatusSuccess](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthTransactionStatusSuccessByCardCache stores monthly successful transaction statistics by card number into cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request with filtering info.
//   - data: Data to be cached.
func (t *transactionStatsByCardStatusCache) SetMonthTransactionStatusSuccessByCardCache(ctx context.Context, req *requests.MonthStatusTransactionCardNumber, data []*response.TransactionResponseMonthStatusSuccess) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupStatusSuccessByCardCacheKey, req.CardNumber, req.Month, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

// GetYearTransactionStatusSuccessByCardCache retrieves cached yearly successful transaction statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request containing card number and year.
//
// Returns:
//   - []*response.TransactionResponseYearStatusSuccess: Yearly successful transactions.
//   - bool: Whether the cache was found.
func (t *transactionStatsByCardStatusCache) GetYearTransactionStatusSuccessByCardCache(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*response.TransactionResponseYearStatusSuccess, bool) {
	key := fmt.Sprintf(yearTopupStatusSuccessByCardCacheKey, req.CardNumber, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.TransactionResponseYearStatusSuccess](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearTransactionStatusSuccessByCardCache stores yearly successful transaction statistics by card number into cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request object with filtering details.
//   - data: The data to cache.
func (t *transactionStatsByCardStatusCache) SetYearTransactionStatusSuccessByCardCache(ctx context.Context, req *requests.YearStatusTransactionCardNumber, data []*response.TransactionResponseYearStatusSuccess) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupStatusSuccessByCardCacheKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

// GetMonthTransactionStatusFailedByCardCache retrieves cached monthly failed transaction statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request with month, year, and card number.
//
// Returns:
//   - []*response.TransactionResponseMonthStatusFailed: Monthly failed transactions.
//   - bool: Whether the cache was found.
func (t *transactionStatsByCardStatusCache) GetMonthTransactionStatusFailedByCardCache(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*response.TransactionResponseMonthStatusFailed, bool) {
	key := fmt.Sprintf(monthTopupStatusFailedByCardCacheKey, req.CardNumber, req.Month, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.TransactionResponseMonthStatusFailed](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

// SetMonthTransactionStatusFailedByCardCache stores monthly failed transaction statistics by card number into cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request filter.
//   - data: Failed transaction data.
func (t *transactionStatsByCardStatusCache) SetMonthTransactionStatusFailedByCardCache(ctx context.Context, req *requests.MonthStatusTransactionCardNumber, data []*response.TransactionResponseMonthStatusFailed) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(monthTopupStatusFailedByCardCacheKey, req.CardNumber, req.Month, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

// GetYearTransactionStatusFailedByCardCache retrieves cached yearly failed transaction statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request with card number and year.
//
// Returns:
//   - []*response.TransactionResponseYearStatusFailed: Failed transactions.
//   - bool: Whether the cache was found.
func (t *transactionStatsByCardStatusCache) GetYearTransactionStatusFailedByCardCache(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*response.TransactionResponseYearStatusFailed, bool) {
	key := fmt.Sprintf(yearTopupStatusFailedByCardCacheKey, req.CardNumber, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.TransactionResponseYearStatusFailed](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

// SetYearTransactionStatusFailedByCardCache stores yearly failed transaction statistics by card number into cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Request with filter info.
//   - data: Yearly failed transaction data.
func (t *transactionStatsByCardStatusCache) SetYearTransactionStatusFailedByCardCache(ctx context.Context, req *requests.YearStatusTransactionCardNumber, data []*response.TransactionResponseYearStatusFailed) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(yearTopupStatusFailedByCardCacheKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}
