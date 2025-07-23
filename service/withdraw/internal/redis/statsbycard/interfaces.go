package withdrawstatsbycardcache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type WithdrawStatsByCardStatusCache interface {
	// GetCachedMonthWithdrawStatusSuccessByCardNumber retrieves cached monthly successful withdraw statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing month and card number.
	//
	// Returns:
	//   - []*response.WithdrawResponseMonthStatusSuccess: List of monthly successful withdraw statistics.
	//   - bool: Whether the cache was found.
	GetCachedMonthWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*response.WithdrawResponseMonthStatusSuccess, bool)

	// SetCachedMonthWithdrawStatusSuccessByCardNumber stores monthly successful withdraw statistics in the cache by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing month and card number.
	//   - data: The data to cache.
	SetCachedMonthWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber, data []*response.WithdrawResponseMonthStatusSuccess)

	// GetCachedYearlyWithdrawStatusSuccessByCardNumber retrieves cached yearly successful withdraw statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//
	// Returns:
	//   - []*response.WithdrawResponseYearStatusSuccess: List of yearly successful withdraw statistics.
	//   - bool: Whether the cache was found.
	GetCachedYearlyWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*response.WithdrawResponseYearStatusSuccess, bool)

	// SetCachedYearlyWithdrawStatusSuccessByCardNumber stores yearly successful withdraw statistics in the cache by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//   - data: The data to cache.
	SetCachedYearlyWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber, data []*response.WithdrawResponseYearStatusSuccess)

	// GetCachedMonthWithdrawStatusFailedByCardNumber retrieves cached monthly failed withdraw statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing month and card number.
	//
	// Returns:
	//   - []*response.WithdrawResponseMonthStatusFailed: List of monthly failed withdraw statistics.
	//   - bool: Whether the cache was found.
	GetCachedMonthWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*response.WithdrawResponseMonthStatusFailed, bool)

	// SetCachedMonthWithdrawStatusFailedByCardNumber stores monthly failed withdraw statistics in the cache by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing month and card number.
	//   - data: The data to cache.
	SetCachedMonthWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber, data []*response.WithdrawResponseMonthStatusFailed)

	// GetCachedYearlyWithdrawStatusFailedByCardNumber retrieves cached yearly failed withdraw statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//
	// Returns:
	//   - []*response.WithdrawResponseYearStatusFailed: List of yearly failed withdraw statistics.
	//   - bool: Whether the cache was found.
	GetCachedYearlyWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*response.WithdrawResponseYearStatusFailed, bool)

	// SetCachedYearlyWithdrawStatusFailedByCardNumber stores yearly failed withdraw statistics in the cache by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//   - data: The data to cache.
	SetCachedYearlyWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber, data []*response.WithdrawResponseYearStatusFailed)
}

type WithdrawStatsByCardAmountCache interface {
	// GetCachedMonthlyWithdrawsByCardNumber retrieves cached monthly withdraw amounts by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year, month, and card number.
	//
	// Returns:
	//   - []*response.WithdrawMonthlyAmountResponse: List of monthly withdraw amounts.
	//   - bool: Whether the cache was found.
	GetCachedMonthlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*response.WithdrawMonthlyAmountResponse, bool)

	// SetCachedMonthlyWithdrawsByCardNumber stores monthly withdraw amounts in the cache by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year, month, and card number.
	//   - data: The data to cache.
	SetCachedMonthlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber, data []*response.WithdrawMonthlyAmountResponse)

	// GetCachedYearlyWithdrawsByCardNumber retrieves cached yearly withdraw amounts by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//
	// Returns:
	//   - []*response.WithdrawYearlyAmountResponse: List of yearly withdraw amounts.
	//   - bool: Whether the cache was found.
	GetCachedYearlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*response.WithdrawYearlyAmountResponse, bool)

	// SetCachedYearlyWithdrawsByCardNumber stores yearly withdraw amounts in the cache by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//   - data: The data to cache.
	SetCachedYearlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber, data []*response.WithdrawYearlyAmountResponse)
}
