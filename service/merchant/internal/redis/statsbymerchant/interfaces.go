package merchantstatsbymerchant

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type MerchantStatsMethodByMerchantCache interface {
	// GetMonthlyPaymentMethodByMerchantsCache retrieves the monthly payment method statistics per merchant from cache.
	// If the cache is found and contains a valid response, it will return the cached
	// response. Otherwise, it will return nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing month, year, and merchant ID.
	//
	// Returns:
	//   - []*response.MerchantResponseMonthlyPaymentMethod: The cached data.
	//   - bool: Whether the cache is found and valid.
	GetMonthlyPaymentMethodByMerchantsCache(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*response.MerchantResponseMonthlyPaymentMethod, bool)

	// SetMonthlyPaymentMethodByMerchantsCache stores the monthly payment method statistics per merchant into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object used as cache key.
	//   - data: The data to cache.
	SetMonthlyPaymentMethodByMerchantsCache(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant, data []*response.MerchantResponseMonthlyPaymentMethod)

	// GetYearlyPaymentMethodByMerchantsCache retrieves the yearly payment method statistics per merchant from cache.
	// If the cache is found and contains a valid response, it will return the cached
	// response. Otherwise, it will return nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing year and merchant ID.
	//
	// Returns:
	//   - []*response.MerchantResponseYearlyPaymentMethod: The cached data.
	//   - bool: Whether the cache is found and valid.
	GetYearlyPaymentMethodByMerchantsCache(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*response.MerchantResponseYearlyPaymentMethod, bool)

	// SetYearlyPaymentMethodByMerchantsCache stores the yearly payment method statistics per merchant into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object used as cache key.
	//   - data: The data to cache.
	SetYearlyPaymentMethodByMerchantsCache(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant, data []*response.MerchantResponseYearlyPaymentMethod)
}

type MerchantStatsAmountByMerchantCache interface {
	// GetMonthlyAmountByMerchantsCache retrieves the monthly amount statistics per merchant from cache.
	// If the cache is found and contains a valid response, it will return the cached
	// response. Otherwise, it will return nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing month, year, and merchant ID.
	//
	// Returns:
	//   - []*response.MerchantResponseMonthlyAmount: The cached data.
	//   - bool: Whether the cache is found and valid.
	GetMonthlyAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*response.MerchantResponseMonthlyAmount, bool)

	// SetMonthlyAmountByMerchantsCache stores the monthly amount statistics per merchant into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object used as cache key.
	//   - data: The data to cache.
	SetMonthlyAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearAmountMerchant, data []*response.MerchantResponseMonthlyAmount)

	// GetYearlyAmountByMerchantsCache retrieves the yearly amount statistics per merchant from cache.
	// If the cache is found and contains a valid response, it will return the cached
	// response. Otherwise, it will return nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing year and merchant ID.
	//
	// Returns:
	//   - []*response.MerchantResponseYearlyAmount: The cached data.
	//   - bool: Whether the cache is found and valid.
	GetYearlyAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*response.MerchantResponseYearlyAmount, bool)

	// SetYearlyAmountByMerchantsCache stores the yearly amount statistics per merchant into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object used as cache key.
	//   - data: The data to cache.
	SetYearlyAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearAmountMerchant, data []*response.MerchantResponseYearlyAmount)
}

type MerchantStatsTotalAmountByMerchantCache interface {
	// GetMonthlyTotalAmountByMerchantsCache retrieves the monthly total amount statistics per merchant from cache.
	// If the cache is found and contains a valid response, it will return the cached
	// response. Otherwise, it will return nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing month, year, and merchant ID.
	//
	// Returns:
	//   - []*response.MerchantResponseMonthlyTotalAmount: The cached data.
	//   - bool: Whether the cache is found and valid.
	GetMonthlyTotalAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) ([]*response.MerchantResponseMonthlyTotalAmount, bool)

	// SetMonthlyTotalAmountByMerchantsCache stores the monthly total amount statistics per merchant into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object used as cache key.
	//   - data: The data to cache.
	SetMonthlyTotalAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearTotalAmountMerchant, data []*response.MerchantResponseMonthlyTotalAmount)

	// GetYearlyTotalAmountByMerchantsCache retrieves the yearly total amount statistics per merchant from cache.
	// If the cache is found and contains a valid response, it will return the cached
	// response. Otherwise, it will return nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing year and merchant ID.
	//
	// Returns:
	//   - []*response.MerchantResponseYearlyTotalAmount: The cached data.
	//   - bool: Whether the cache is found and valid.
	GetYearlyTotalAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) ([]*response.MerchantResponseYearlyTotalAmount, bool)

	// SetYearlyTotalAmountByMerchantsCache stores the yearly total amount statistics per merchant into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object used as cache key.
	//   - data: The data to cache.
	SetYearlyTotalAmountByMerchantsCache(ctx context.Context, req *requests.MonthYearTotalAmountMerchant, data []*response.MerchantResponseYearlyTotalAmount)
}
