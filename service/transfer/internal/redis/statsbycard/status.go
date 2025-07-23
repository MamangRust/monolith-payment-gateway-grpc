package transferstatsbycardcache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type transferStatsByCardStatusCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTransferStatsByCardStatusCache(store *sharedcachehelpers.CacheStore) TransferStatsByCardStatusCache {
	return &transferStatsByCardStatusCache{store: store}
}

// GetMonthTransferStatusSuccessByCard retrieves cached monthly successful transfers for a specific card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Contains card number and month.
//
// Returns:
//   - []*response.TransferResponseMonthStatusSuccess: List of monthly successful transfers.
//   - bool: Whether the cache was found.
func (t *transferStatsByCardStatusCache) GetMonthTransferStatusSuccessByCard(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*response.TransferResponseMonthStatusSuccess, bool) {
	key := fmt.Sprintf(transferMonthTransferStatusSuccessByCardKey, req.CardNumber, req.Month, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.TransferResponseMonthStatusSuccess](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true

}

// SetMonthTransferStatusSuccessByCard stores monthly successful transfers for a specific card number into cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Contains card number and month.
//   - data: List of monthly successful transfers to cache.
func (t *transferStatsByCardStatusCache) SetMonthTransferStatusSuccessByCard(ctx context.Context, req *requests.MonthStatusTransferCardNumber, data []*response.TransferResponseMonthStatusSuccess) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferMonthTransferStatusSuccessByCardKey, req.CardNumber, req.Month, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

// GetYearlyTransferStatusSuccessByCard retrieves cached yearly successful transfers for a specific card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Contains card number and year.
//
// Returns:
//   - []*response.TransferResponseYearStatusSuccess: List of yearly successful transfers.
//   - bool: Whether the cache was found.
func (t *transferStatsByCardStatusCache) GetYearlyTransferStatusSuccessByCard(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*response.TransferResponseYearStatusSuccess, bool) {
	key := fmt.Sprintf(transferYearTransferStatusSuccessByCardKey, req.CardNumber, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.TransferResponseYearStatusSuccess](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyTransferStatusSuccessByCard stores yearly successful transfers for a specific card number into cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Contains card number and year.
//   - data: List of yearly successful transfers to cache.
func (t *transferStatsByCardStatusCache) SetYearlyTransferStatusSuccessByCard(ctx context.Context, req *requests.YearStatusTransferCardNumber, data []*response.TransferResponseYearStatusSuccess) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferYearTransferStatusSuccessByCardKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

// GetMonthTransferStatusFailedByCard retrieves cached monthly failed transfers for a specific card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Contains card number and month.
//
// Returns:
//   - []*response.TransferResponseMonthStatusFailed: List of monthly failed transfers.
//   - bool: Whether the cache was found.
func (t *transferStatsByCardStatusCache) GetMonthTransferStatusFailedByCard(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*response.TransferResponseMonthStatusFailed, bool) {
	key := fmt.Sprintf(transferMonthTransferStatusFailedByCardKey, req.CardNumber, req.Month, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.TransferResponseMonthStatusFailed](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthTransferStatusFailedByCard stores monthly failed transfers for a specific card number into cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Contains card number and month.
//   - data: List of monthly failed transfers to cache.
func (t *transferStatsByCardStatusCache) SetMonthTransferStatusFailedByCard(ctx context.Context, req *requests.MonthStatusTransferCardNumber, data []*response.TransferResponseMonthStatusFailed) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferMonthTransferStatusFailedByCardKey, req.CardNumber, req.Month, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

// GetYearlyTransferStatusFailedByCard retrieves cached yearly failed transfers for a specific card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Contains card number and year.
//
// Returns:
//   - []*response.TransferResponseYearStatusFailed: List of yearly failed transfers.
//   - bool: Whether the cache was found.
func (t *transferStatsByCardStatusCache) GetYearlyTransferStatusFailedByCard(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*response.TransferResponseYearStatusFailed, bool) {
	key := fmt.Sprintf(transferYearTransferStatusFailedByCardKey, req.CardNumber, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.TransferResponseYearStatusFailed](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

// SetYearlyTransferStatusFailedByCard stores yearly failed transfers for a specific card number into cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Contains card number and year.
//   - data: List of yearly failed transfers to cache.
func (t *transferStatsByCardStatusCache) SetYearlyTransferStatusFailedByCard(ctx context.Context, req *requests.YearStatusTransferCardNumber, data []*response.TransferResponseYearStatusFailed) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferYearTransferStatusFailedByCardKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}
