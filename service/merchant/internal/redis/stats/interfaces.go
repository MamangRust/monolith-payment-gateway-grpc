package merchantstatscache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type MerchantStatsMethodCache interface {
	// GetMonthlyPaymentMethodsMerchantCache retrieves the monthly payment method statistics of merchants from cache.
	// If the cache is found and contains a valid response, it will return the cached
	// response. Otherwise, it will return nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the data is requested.
	//
	// Returns:
	//   - []*response.MerchantResponseMonthlyPaymentMethod: The cached monthly payment method data.
	//   - bool: Whether the cache is found and valid.
	GetMonthlyPaymentMethodsMerchantCache(ctx context.Context, year int) ([]*response.MerchantResponseMonthlyPaymentMethod, bool)

	// SetMonthlyPaymentMethodsMerchantCache stores the monthly payment method statistics of merchants into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the data is cached.
	//   - data: The data to cache.
	SetMonthlyPaymentMethodsMerchantCache(ctx context.Context, year int, data []*response.MerchantResponseMonthlyPaymentMethod)

	// GetYearlyPaymentMethodMerchantCache retrieves the yearly payment method statistics of merchants from cache.
	// If the cache is found and contains a valid response, it will return the cached
	// response. Otherwise, it will return nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the data is requested.
	//
	// Returns:
	//   - []*response.MerchantResponseYearlyPaymentMethod: The cached yearly payment method data.
	//   - bool: Whether the cache is found and valid.
	GetYearlyPaymentMethodMerchantCache(ctx context.Context, year int) ([]*response.MerchantResponseYearlyPaymentMethod, bool)

	// SetYearlyPaymentMethodMerchantCache stores the yearly payment method statistics of merchants into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the data is cached.
	//   - data: The data to cache.
	SetYearlyPaymentMethodMerchantCache(ctx context.Context, year int, data []*response.MerchantResponseYearlyPaymentMethod)
}

type MerchantStatsAmountCache interface {
	// GetMonthlyAmountMerchantCache retrieves the monthly amount statistics of merchants from cache.
	// If the cache is found and contains a valid response, it will return the cached
	// response. Otherwise, it will return nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the data is requested.
	//
	// Returns:
	//   - []*response.MerchantResponseMonthlyAmount: The cached monthly amount data.
	//   - bool: Whether the cache is found and valid.
	GetMonthlyAmountMerchantCache(ctx context.Context, year int) ([]*response.MerchantResponseMonthlyAmount, bool)

	// SetMonthlyAmountMerchantCache stores the monthly amount statistics of merchants into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the data is cached.
	//   - data: The data to cache.
	SetMonthlyAmountMerchantCache(ctx context.Context, year int, data []*response.MerchantResponseMonthlyAmount)

	// GetYearlyAmountMerchantCache retrieves the yearly amount statistics of merchants from cache.
	// If the cache is found and contains a valid response, it will return the cached
	// response. Otherwise, it will return nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the data is requested.
	//
	// Returns:
	//   - []*response.MerchantResponseYearlyAmount: The cached yearly amount data.
	//   - bool: Whether the cache is found and valid.
	GetYearlyAmountMerchantCache(ctx context.Context, year int) ([]*response.MerchantResponseYearlyAmount, bool)

	// SetYearlyAmountMerchantCache stores the yearly amount statistics of merchants into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the data is cached.
	//   - data: The data to cache.
	SetYearlyAmountMerchantCache(ctx context.Context, year int, data []*response.MerchantResponseYearlyAmount)
}

type MerchantStatsTotalAmountCache interface {
	// GetMonthlyTotalAmountMerchantCache retrieves the monthly total amount statistics of merchants from cache.
	// If the cache is found and contains a valid response, it will return the cached
	// response. Otherwise, it will return nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the data is requested.
	//
	// Returns:
	//   - []*response.MerchantResponseMonthlyTotalAmount: The cached monthly total amount data.
	//   - bool: Whether the cache is found and valid.
	GetMonthlyTotalAmountMerchantCache(ctx context.Context, year int) ([]*response.MerchantResponseMonthlyTotalAmount, bool)

	// SetMonthlyTotalAmountMerchantCache stores the monthly total amount statistics of merchants into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the data is cached.
	//   - data: The data to cache.
	SetMonthlyTotalAmountMerchantCache(ctx context.Context, year int, data []*response.MerchantResponseMonthlyTotalAmount)

	// GetYearlyTotalAmountMerchantCache retrieves the yearly total amount statistics of merchants from cache.
	// If the cache is found and contains a valid response, it will return the cached
	// response. Otherwise, it will return nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the data is requested.
	//
	// Returns:
	//   - []*response.MerchantResponseYearlyTotalAmount: The cached yearly total amount data.
	//   - bool: Whether the cache is found and valid.
	GetYearlyTotalAmountMerchantCache(ctx context.Context, year int) ([]*response.MerchantResponseYearlyTotalAmount, bool)

	// SetYearlyTotalAmountMerchantCache stores the yearly total amount statistics of merchants into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the data is cached.
	//   - data: The data to cache.
	SetYearlyTotalAmountMerchantCache(ctx context.Context, year int, data []*response.MerchantResponseYearlyTotalAmount)
}
