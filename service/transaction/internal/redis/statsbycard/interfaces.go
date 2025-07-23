package transactionstatsbycarcache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type TransactionStatsByCardAmountCache interface {
	// GetMonthlyAmountsByCardCache retrieves cached monthly transaction amount statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Month, year, card number filter.
	//
	// Returns:
	//   - []*response.TransactionMonthAmountResponse: Monthly amounts.
	//   - bool: Whether the cache was found.
	GetMonthlyAmountsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*response.TransactionMonthAmountResponse, bool)

	// SetMonthlyAmountsByCardCache stores monthly transaction amount statistics by card number into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request details.
	//   - data: Amounts to cache.
	SetMonthlyAmountsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod, data []*response.TransactionMonthAmountResponse)

	// GetYearlyAmountsByCardCache retrieves cached yearly transaction amount statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Card number and year info.
	//
	// Returns:
	//   - []*response.TransactionYearlyAmountResponse: Yearly amounts.
	//   - bool: Whether the cache was found.
	GetYearlyAmountsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*response.TransactionYearlyAmountResponse, bool)

	// SetYearlyAmountsByCardCache stores yearly transaction amount statistics by card number into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request filter.
	//   - data: Yearly amount data.
	SetYearlyAmountsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod, data []*response.TransactionYearlyAmountResponse)
}

type TransactionStatsByCardMethodCache interface {
	// GetMonthlyPaymentMethodsByCardCache retrieves cached monthly payment method statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Month, year, and card number request.
	//
	// Returns:
	//   - []*response.TransactionMonthMethodResponse: Payment methods per month.
	//   - bool: Whether the cache was found.
	GetMonthlyPaymentMethodsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*response.TransactionMonthMethodResponse, bool)

	// SetMonthlyPaymentMethodsByCardCache stores monthly payment method statistics by card number into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request parameters.
	//   - data: Monthly method stats.
	SetMonthlyPaymentMethodsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod, data []*response.TransactionMonthMethodResponse)

	// GetYearlyPaymentMethodsByCardCache retrieves cached yearly payment method statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request object with card and year info.
	//
	// Returns:
	//   - []*response.TransactionYearMethodResponse: Yearly method stats.
	//   - bool: Whether the cache was found.
	GetYearlyPaymentMethodsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*response.TransactionYearMethodResponse, bool)

	// SetYearlyPaymentMethodsByCardCache stores yearly payment method statistics by card number into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request details.
	//   - data: The method stats to cache.
	SetYearlyPaymentMethodsByCardCache(ctx context.Context, req *requests.MonthYearPaymentMethod, data []*response.TransactionYearMethodResponse)
}

type TransactionStatsByCardStatusCache interface {
	// GetMonthTransactionStatusSuccessByCardCache retrieves cached monthly successful transaction statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request containing card number, month, and year.
	//
	// Returns:
	//   - []*response.TransactionResponseMonthStatusSuccess: List of successful transactions.
	//   - bool: Whether the cache was found.
	GetMonthTransactionStatusSuccessByCardCache(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*response.TransactionResponseMonthStatusSuccess, bool)

	// SetMonthTransactionStatusSuccessByCardCache stores monthly successful transaction statistics by card number into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request with filtering info.
	//   - data: Data to be cached.
	SetMonthTransactionStatusSuccessByCardCache(ctx context.Context, req *requests.MonthStatusTransactionCardNumber, data []*response.TransactionResponseMonthStatusSuccess)

	// GetYearTransactionStatusSuccessByCardCache retrieves cached yearly successful transaction statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request containing card number and year.
	//
	// Returns:
	//   - []*response.TransactionResponseYearStatusSuccess: Yearly successful transactions.
	//   - bool: Whether the cache was found.
	GetYearTransactionStatusSuccessByCardCache(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*response.TransactionResponseYearStatusSuccess, bool)

	// SetYearTransactionStatusSuccessByCardCache stores yearly successful transaction statistics by card number into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request object with filtering details.
	//   - data: The data to cache.
	SetYearTransactionStatusSuccessByCardCache(ctx context.Context, req *requests.YearStatusTransactionCardNumber, data []*response.TransactionResponseYearStatusSuccess)

	// GetMonthTransactionStatusFailedByCardCache retrieves cached monthly failed transaction statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request with month, year, and card number.
	//
	// Returns:
	//   - []*response.TransactionResponseMonthStatusFailed: Monthly failed transactions.
	//   - bool: Whether the cache was found.
	GetMonthTransactionStatusFailedByCardCache(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*response.TransactionResponseMonthStatusFailed, bool)

	// SetMonthTransactionStatusFailedByCardCache stores monthly failed transaction statistics by card number into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request filter.
	//   - data: Failed transaction data.
	SetMonthTransactionStatusFailedByCardCache(ctx context.Context, req *requests.MonthStatusTransactionCardNumber, data []*response.TransactionResponseMonthStatusFailed)

	// GetYearTransactionStatusFailedByCardCache retrieves cached yearly failed transaction statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request with card number and year.
	//
	// Returns:
	//   - []*response.TransactionResponseYearStatusFailed: Failed transactions.
	//   - bool: Whether the cache was found.
	GetYearTransactionStatusFailedByCardCache(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*response.TransactionResponseYearStatusFailed, bool)

	// SetYearTransactionStatusFailedByCardCache stores yearly failed transaction statistics by card number into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request with filter info.
	//   - data: Yearly failed transaction data.
	SetYearTransactionStatusFailedByCardCache(ctx context.Context, req *requests.YearStatusTransactionCardNumber, data []*response.TransactionResponseYearStatusFailed)
}
