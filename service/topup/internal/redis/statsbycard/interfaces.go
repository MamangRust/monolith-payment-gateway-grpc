package topupstatsbycardcache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type TopupStatsStatusByCardCache interface {
	// GetMonthTopupStatusSuccessByCardNumberCache retrieves cached monthly successful topup statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing month, year, and card number.
	//
	// Returns:
	//   - []*response.TopupResponseMonthStatusSuccess: List of successful topup statistics.
	//   - bool: Whether the cache was found.
	GetMonthTopupStatusSuccessByCardNumberCache(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*response.TopupResponseMonthStatusSuccess, bool)

	// SetMonthTopupStatusSuccessByCardNumberCache stores monthly successful topup statistics by card number in cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing month, year, and card number.
	//   - data: The data to store in cache.
	SetMonthTopupStatusSuccessByCardNumberCache(ctx context.Context, req *requests.MonthTopupStatusCardNumber, data []*response.TopupResponseMonthStatusSuccess)

	// GetYearlyTopupStatusSuccessByCardNumberCache retrieves cached yearly successful topup statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//
	// Returns:
	//   - []*response.TopupResponseYearStatusSuccess: List of yearly topup statistics.
	//   - bool: Whether the cache was found.
	GetYearlyTopupStatusSuccessByCardNumberCache(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*response.TopupResponseYearStatusSuccess, bool)

	// SetYearlyTopupStatusSuccessByCardNumberCache stores yearly successful topup statistics by card number in cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//   - data: The data to store in cache.
	SetYearlyTopupStatusSuccessByCardNumberCache(ctx context.Context, req *requests.YearTopupStatusCardNumber, data []*response.TopupResponseYearStatusSuccess)

	// GetMonthTopupStatusFailedByCardNumberCache retrieves cached monthly failed topup statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing month, year, and card number.
	//
	// Returns:
	//   - []*response.TopupResponseMonthStatusFailed: List of failed topup statistics.
	//   - bool: Whether the cache was found.
	GetMonthTopupStatusFailedByCardNumberCache(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*response.TopupResponseMonthStatusFailed, bool)

	// SetMonthTopupStatusFailedByCardNumberCache stores monthly failed topup statistics by card number in cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing month, year, and card number.
	//   - data: The data to store in cache.
	SetMonthTopupStatusFailedByCardNumberCache(ctx context.Context, req *requests.MonthTopupStatusCardNumber, data []*response.TopupResponseMonthStatusFailed)

	// GetYearlyTopupStatusFailedByCardNumberCache retrieves cached yearly failed topup statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//
	// Returns:
	//   - []*response.TopupResponseYearStatusFailed: List of failed topup statistics.
	//   - bool: Whether the cache was found.
	GetYearlyTopupStatusFailedByCardNumberCache(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*response.TopupResponseYearStatusFailed, bool)

	// SetYearlyTopupStatusFailedByCardNumberCache stores yearly failed topup statistics by card number in cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//   - data: The data to store in cache.
	SetYearlyTopupStatusFailedByCardNumberCache(ctx context.Context, req *requests.YearTopupStatusCardNumber, data []*response.TopupResponseYearStatusFailed)
}

type TopupStatsMethodByCardCache interface {
	// GetMonthlyTopupMethodsByCardNumberCache retrieves cached monthly topup methods by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year, month, and card number.
	//
	// Returns:
	//   - []*response.TopupMonthMethodResponse: List of topup method statistics.
	//   - bool: Whether the cache was found.
	GetMonthlyTopupMethodsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod) ([]*response.TopupMonthMethodResponse, bool)

	// SetMonthlyTopupMethodsByCardNumberCache stores monthly topup methods by card number in cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year, month, and card number.
	//   - data: The data to store in cache.
	SetMonthlyTopupMethodsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod, data []*response.TopupMonthMethodResponse)

	// GetYearlyTopupMethodsByCardNumberCache retrieves cached yearly topup methods by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//
	// Returns:
	//   - []*response.TopupYearlyMethodResponse: List of topup method statistics.
	//   - bool: Whether the cache was found.
	GetYearlyTopupMethodsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod) ([]*response.TopupYearlyMethodResponse, bool)

	// SetYearlyTopupMethodsByCardNumberCache stores yearly topup methods by card number in cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//   - data: The data to store in cache.
	SetYearlyTopupMethodsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod, data []*response.TopupYearlyMethodResponse)
}

type TopupStatsAmountByCardCache interface {
	// GetMonthlyTopupAmountsByCardNumberCache retrieves cached monthly topup amounts by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year, month, and card number.
	//
	// Returns:
	//   - []*response.TopupMonthAmountResponse: List of topup amount statistics.
	//   - bool: Whether the cache was found.
	GetMonthlyTopupAmountsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod) ([]*response.TopupMonthAmountResponse, bool)

	// SetMonthlyTopupAmountsByCardNumberCache stores monthly topup amounts by card number in cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year, month, and card number.
	//   - data: The data to store in cache.
	SetMonthlyTopupAmountsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod, data []*response.TopupMonthAmountResponse)

	// GetYearlyTopupAmountsByCardNumberCache retrieves cached yearly topup amounts by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//
	// Returns:
	//   - []*response.TopupYearlyAmountResponse: List of yearly topup amount statistics.
	//   - bool: Whether the cache was found.
	GetYearlyTopupAmountsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod) ([]*response.TopupYearlyAmountResponse, bool)

	// SetYearlyTopupAmountsByCardNumberCache stores yearly topup amounts by card number in cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//   - data: The data to store in cache.
	SetYearlyTopupAmountsByCardNumberCache(ctx context.Context, req *requests.YearMonthMethod, data []*response.TopupYearlyAmountResponse)
}
