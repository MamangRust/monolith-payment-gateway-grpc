package merchantstatsapikey

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type MerchantStatsMethodByApiKeyCache interface {
	// GetMonthlyPaymentMethodByApikeysCache retrieves the monthly payment method statistics by API key from cache.
	// If the cache is found and contains a valid response, it will return the cached response. Otherwise, it will return nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing month, year, and API key.
	//
	// Returns:
	//   - []*response.MerchantResponseMonthlyPaymentMethod: The cached data.
	//   - bool: Whether the cache is found and valid.
	GetMonthlyPaymentMethodByApikeysCache(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*response.MerchantResponseMonthlyPaymentMethod, bool)

	// SetMonthlyPaymentMethodByApikeysCache stores the monthly payment method statistics by API key into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object used as the cache key.
	//   - data: The data to cache.
	SetMonthlyPaymentMethodByApikeysCache(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey, data []*response.MerchantResponseMonthlyPaymentMethod)

	// GetYearlyPaymentMethodByApikeysCache retrieves the yearly payment method statistics by API key from cache.
	// If the cache is found and contains a valid response, it will return the cached response. Otherwise, it will return nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing year and API key.
	//
	// Returns:
	//   - []*response.MerchantResponseYearlyPaymentMethod: The cached data.
	//   - bool: Whether the cache is found and valid.
	GetYearlyPaymentMethodByApikeysCache(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*response.MerchantResponseYearlyPaymentMethod, bool)

	// SetYearlyPaymentMethodByApikeysCache stores the yearly payment method statistics by API key into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object used as the cache key.
	//   - data: The data to cache.
	SetYearlyPaymentMethodByApikeysCache(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey, data []*response.MerchantResponseYearlyPaymentMethod)
}

type MerchantStatsAmountByApiKeyCache interface {
	// GetMonthlyAmountByApikeysCache retrieves the monthly amount statistics by API key from cache.
	// If the cache is found and contains a valid response, it will return the cached response. Otherwise, it will return nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing month, year, and API key.
	//
	// Returns:
	//   - []*response.MerchantResponseMonthlyAmount: The cached data.
	//   - bool: Whether the cache is found and valid.
	GetMonthlyAmountByApikeysCache(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*response.MerchantResponseMonthlyAmount, bool)

	// SetMonthlyAmountByApikeysCache stores the monthly amount statistics by API key into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object used as the cache key.
	//   - data: The data to cache.
	SetMonthlyAmountByApikeysCache(ctx context.Context, req *requests.MonthYearAmountApiKey, data []*response.MerchantResponseMonthlyAmount)

	// GetYearlyAmountByApikeysCache retrieves the yearly amount statistics by API key from cache.
	// If the cache is found and contains a valid response, it will return the cached response. Otherwise, it will return nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing year and API key.
	//
	// Returns:
	//   - []*response.MerchantResponseYearlyAmount: The cached data.
	//   - bool: Whether the cache is found and valid.
	GetYearlyAmountByApikeysCache(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*response.MerchantResponseYearlyAmount, bool)

	// SetYearlyAmountByApikeysCache stores the yearly amount statistics by API key into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object used as the cache key.
	//   - data: The data to cache.
	SetYearlyAmountByApikeysCache(ctx context.Context, req *requests.MonthYearAmountApiKey, data []*response.MerchantResponseYearlyAmount)
}

type MerchantStatsTotalAmountByApiKeyCache interface {
	// GetMonthlyTotalAmountByApikeysCache retrieves the monthly total amount statistics by API key from cache.
	// If the cache is found and contains a valid response, it will return the cached response. Otherwise, it will return nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing month, year, and API key.
	//
	// Returns:
	//   - []*response.MerchantResponseMonthlyTotalAmount: The cached data.
	//   - bool: Whether the cache is found and valid.
	GetMonthlyTotalAmountByApikeysCache(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*response.MerchantResponseMonthlyTotalAmount, bool)

	// SetMonthlyTotalAmountByApikeysCache stores the monthly total amount statistics by API key into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object used as the cache key.
	//   - data: The data to cache.
	SetMonthlyTotalAmountByApikeysCache(ctx context.Context, req *requests.MonthYearTotalAmountApiKey, data []*response.MerchantResponseMonthlyTotalAmount)

	// GetYearlyTotalAmountByApikeysCache retrieves the yearly total amount statistics by API key from cache.
	// If the cache is found and contains a valid response, it will return the cached response. Otherwise, it will return nil, false.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing year and API key.
	//
	// Returns:
	//   - []*response.MerchantResponseYearlyTotalAmount: The cached data.
	//   - bool: Whether the cache is found and valid.
	GetYearlyTotalAmountByApikeysCache(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*response.MerchantResponseYearlyTotalAmount, bool)

	// SetYearlyTotalAmountByApikeysCache stores the yearly total amount statistics by API key into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object used as the cache key.
	//   - data: The data to cache.
	SetYearlyTotalAmountByApikeysCache(ctx context.Context, req *requests.MonthYearTotalAmountApiKey, data []*response.MerchantResponseYearlyTotalAmount)
}
