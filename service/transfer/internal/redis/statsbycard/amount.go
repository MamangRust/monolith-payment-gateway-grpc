package transferstatsbycardcache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type transferStatsByCardAmountCache struct {
	store *sharedcachehelpers.CacheStore
}

func NewTransferStatsByCardAmountCache(store *sharedcachehelpers.CacheStore) TransferStatsByCardAmountCache {
	return &transferStatsByCardAmountCache{store: store}
}

// GetMonthlyTransferAmountsBySenderCard retrieves cached monthly transfer amounts for a sender card.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Contains sender card number and year/month filter.
//
// Returns:
//   - []*response.TransferMonthAmountResponse: Monthly transfer amount statistics.
//   - bool: Whether the cache was found.
func (t *transferStatsByCardAmountCache) GetMonthlyTransferAmountsBySenderCard(ctx context.Context, req *requests.MonthYearCardNumber) ([]*response.TransferMonthAmountResponse, bool) {
	key := fmt.Sprintf(transferMonthTransferAmountByCardKey, req.CardNumber, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.TransferMonthAmountResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthlyTransferAmountsBySenderCard stores monthly transfer amounts for a sender card into cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Contains sender card number and year/month.
//   - data: List of amounts to cache.
func (t *transferStatsByCardAmountCache) SetMonthlyTransferAmountsBySenderCard(ctx context.Context, req *requests.MonthYearCardNumber, data []*response.TransferMonthAmountResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferMonthTransferAmountByCardKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

// GetMonthlyTransferAmountsByReceiverCard retrieves cached monthly transfer amounts for a receiver card.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Contains receiver card number and year/month filter.
//
// Returns:
//   - []*response.TransferMonthAmountResponse: Monthly transfer amount statistics.
//   - bool: Whether the cache was found.
func (t *transferStatsByCardAmountCache) GetMonthlyTransferAmountsByReceiverCard(ctx context.Context, req *requests.MonthYearCardNumber) ([]*response.TransferMonthAmountResponse, bool) {
	key := fmt.Sprintf(transferMonthTransferAmountByCardKey, req.CardNumber, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.TransferMonthAmountResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetMonthlyTransferAmountsByReceiverCard stores monthly transfer amounts for a receiver card into cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Contains receiver card number and year/month.
//   - data: List of amounts to cache.
func (t *transferStatsByCardAmountCache) SetMonthlyTransferAmountsByReceiverCard(ctx context.Context, req *requests.MonthYearCardNumber, data []*response.TransferMonthAmountResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferMonthTransferAmountByCardKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

// GetYearlyTransferAmountsBySenderCard retrieves cached yearly transfer amounts for a sender card.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Contains sender card number and year.
//
// Returns:
//   - []*response.TransferYearAmountResponse: Yearly transfer amount statistics.
//   - bool: Whether the cache was found.
func (t *transferStatsByCardAmountCache) GetYearlyTransferAmountsBySenderCard(ctx context.Context, req *requests.MonthYearCardNumber) ([]*response.TransferYearAmountResponse, bool) {
	key := fmt.Sprintf(transferYearTransferAmountByCardKey, req.CardNumber, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.TransferYearAmountResponse](ctx, t.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetYearlyTransferAmountsBySenderCard stores yearly transfer amounts for a sender card into cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Contains sender card number and year.
//   - data: List of yearly transfer amounts to cache.
func (t *transferStatsByCardAmountCache) SetYearlyTransferAmountsBySenderCard(ctx context.Context, req *requests.MonthYearCardNumber, data []*response.TransferYearAmountResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferYearTransferAmountByCardKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}

// GetYearlyTransferAmountsByReceiverCard retrieves cached yearly transfer amounts for a receiver card.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Contains receiver card number and year.
//
// Returns:
//   - []*response.TransferYearAmountResponse: Yearly transfer amount statistics.
//   - bool: Whether the cache was found.
func (t *transferStatsByCardAmountCache) GetYearlyTransferAmountsByReceiverCard(ctx context.Context, req *requests.MonthYearCardNumber) ([]*response.TransferYearAmountResponse, bool) {
	key := fmt.Sprintf(transferYearTransferAmountByCardKey, req.CardNumber, req.Year)
	result, found := sharedcachehelpers.GetFromCache[[]*response.TransferYearAmountResponse](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

// SetYearlyTransferAmountsByReceiverCard stores yearly transfer amounts for a receiver card into cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: Contains receiver card number and year.
//   - data: List of yearly transfer amounts to cache.
func (t *transferStatsByCardAmountCache) SetYearlyTransferAmountsByReceiverCard(ctx context.Context, req *requests.MonthYearCardNumber, data []*response.TransferYearAmountResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transferYearTransferAmountByCardKey, req.CardNumber, req.Year)
	sharedcachehelpers.SetToCache(ctx, t.store, key, &data, ttlDefault)
}
