package transactionstatscache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type TransactionStatsAmountCache interface {
	// GetMonthlyAmountsCache retrieves cached monthly transaction amount statistics.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*response.TransactionMonthAmountResponse: Monthly amount statistics.
	//   - bool: Whether the cache was found.
	GetMonthlyAmountsCache(ctx context.Context, year int) ([]*response.TransactionMonthAmountResponse, bool)

	// SetMonthlyAmountsCache caches monthly transaction amount statistics.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year used as key.
	//   - data: Monthly transaction amount data.
	SetMonthlyAmountsCache(ctx context.Context, year int, data []*response.TransactionMonthAmountResponse)

	// GetYearlyAmountsCache retrieves cached yearly transaction amount statistics.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*response.TransactionYearlyAmountResponse: Yearly amount statistics.
	//   - bool: Whether the cache was found.
	GetYearlyAmountsCache(ctx context.Context, year int) ([]*response.TransactionYearlyAmountResponse, bool)

	// SetYearlyAmountsCache caches yearly transaction amount statistics.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year used as cache key.
	//   - data: Yearly amount statistics to cache.
	SetYearlyAmountsCache(ctx context.Context, year int, data []*response.TransactionYearlyAmountResponse)
}

type TransactionStatsMethodCache interface {
	// GetMonthlyPaymentMethodsCache retrieves cached monthly payment method statistics.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*response.TransactionMonthMethodResponse: List of monthly payment method statistics.
	//   - bool: Whether the cache was found.
	GetMonthlyPaymentMethodsCache(ctx context.Context, year int) ([]*response.TransactionMonthMethodResponse, bool)

	// SetMonthlyPaymentMethodsCache caches monthly payment method statistics.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year used as key.
	//   - data: Monthly payment method data.
	SetMonthlyPaymentMethodsCache(ctx context.Context, year int, data []*response.TransactionMonthMethodResponse)

	// GetYearlyPaymentMethodsCache retrieves cached yearly payment method statistics.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*response.TransactionYearMethodResponse: Yearly payment method stats.
	//   - bool: Whether the cache was found.
	GetYearlyPaymentMethodsCache(ctx context.Context, year int) ([]*response.TransactionYearMethodResponse, bool)

	// SetYearlyPaymentMethodsCache caches yearly payment method statistics.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year used as cache key.
	//   - data: Yearly method statistics.
	SetYearlyPaymentMethodsCache(ctx context.Context, year int, data []*response.TransactionYearMethodResponse)
}

type TransactionStatsStatusCache interface {
	// GetMonthTransactionStatusSuccessCache retrieves cached monthly successful transactions.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request containing month and year filter.
	//
	// Returns:
	//   - []*response.TransactionResponseMonthStatusSuccess: List of successful monthly transactions.
	//   - bool: Whether the cache was found.
	GetMonthTransactionStatusSuccessCache(ctx context.Context, req *requests.MonthStatusTransaction) ([]*response.TransactionResponseMonthStatusSuccess, bool)

	// SetMonthTransactionStatusSuccessCache stores successful monthly transactions in cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Original request object.
	//   - data: Transactions to cache.
	SetMonthTransactionStatusSuccessCache(ctx context.Context, req *requests.MonthStatusTransaction, data []*response.TransactionResponseMonthStatusSuccess)

	// GetYearTransactionStatusSuccessCache retrieves cached yearly successful transactions.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which statistics are requested.
	//
	// Returns:
	//   - []*response.TransactionResponseYearStatusSuccess: List of successful yearly transactions.
	//   - bool: Whether the cache was found.
	GetYearTransactionStatusSuccessCache(ctx context.Context, year int) ([]*response.TransactionResponseYearStatusSuccess, bool)

	// SetYearTransactionStatusSuccessCache caches yearly successful transaction statistics.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year used as key.
	//   - data: Yearly successful transaction stats.
	SetYearTransactionStatusSuccessCache(ctx context.Context, year int, data []*response.TransactionResponseYearStatusSuccess)

	// GetMonthTransactionStatusFailedCache retrieves cached monthly failed transactions.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request containing month and year filter.
	//
	// Returns:
	//   - []*response.TransactionResponseMonthStatusFailed: List of failed monthly transactions.
	//   - bool: Whether the cache was found.
	GetMonthTransactionStatusFailedCache(ctx context.Context, req *requests.MonthStatusTransaction) ([]*response.TransactionResponseMonthStatusFailed, bool)

	// SetMonthTransactionStatusFailedCache caches monthly failed transaction statistics.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Original request object.
	//   - data: List of failed monthly transactions.
	SetMonthTransactionStatusFailedCache(ctx context.Context, req *requests.MonthStatusTransaction, data []*response.TransactionResponseMonthStatusFailed)

	// GetYearTransactionStatusFailedCache retrieves cached yearly failed transactions.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*response.TransactionResponseYearStatusFailed: List of failed transactions.
	//   - bool: Whether the cache was found.
	GetYearTransactionStatusFailedCache(ctx context.Context, year int) ([]*response.TransactionResponseYearStatusFailed, bool)

	// SetYearTransactionStatusFailedCache caches yearly failed transaction statistics.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year used as key.
	//   - data: List of failed yearly transactions.
	SetYearTransactionStatusFailedCache(ctx context.Context, year int, data []*response.TransactionResponseYearStatusFailed)
}
