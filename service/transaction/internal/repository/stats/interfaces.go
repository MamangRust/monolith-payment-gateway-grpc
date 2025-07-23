package transactionstatsrepository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type TransactionStatsStatusRepository interface {
	// GetMonthTransactionStatusSuccess retrieves monthly statistics of successful transactions.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing month and status filter.
	//
	// Returns:
	//   - []*record.TransactionRecordMonthStatusSuccess: List of monthly successful transaction statistics.
	//   - error: Error if any occurs during query.
	GetMonthTransactionStatusSuccess(ctx context.Context, req *requests.MonthStatusTransaction) ([]*record.TransactionRecordMonthStatusSuccess, error)

	// GetYearlyTransactionStatusSuccess retrieves yearly statistics of successful transactions.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*record.TransactionRecordYearStatusSuccess: List of yearly successful transaction statistics.
	//   - error: Error if any occurs during query.
	GetYearlyTransactionStatusSuccess(ctx context.Context, year int) ([]*record.TransactionRecordYearStatusSuccess, error)

	// GetMonthTransactionStatusFailed retrieves monthly statistics of failed transactions.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing month and status filter.
	//
	// Returns:
	//   - []*record.TransactionRecordMonthStatusFailed: List of monthly failed transaction statistics.
	//   - error: Error if any occurs during query.
	GetMonthTransactionStatusFailed(ctx context.Context, req *requests.MonthStatusTransaction) ([]*record.TransactionRecordMonthStatusFailed, error)

	// GetYearlyTransactionStatusFailed retrieves yearly statistics of failed transactions.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*record.TransactionRecordYearStatusFailed: List of yearly failed transaction statistics.
	//   - error: Error if any occurs during query.
	GetYearlyTransactionStatusFailed(ctx context.Context, year int) ([]*record.TransactionRecordYearStatusFailed, error)
}

type TransactionStatsMethodRepository interface {
	// GetMonthlyPaymentMethods retrieves monthly statistics grouped by payment method.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*record.TransactionMonthMethod: List of monthly payment method usage statistics.
	//   - error: Error if any occurs during query.
	GetMonthlyPaymentMethods(ctx context.Context, year int) ([]*record.TransactionMonthMethod, error)

	// GetYearlyPaymentMethods retrieves yearly statistics grouped by payment method.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*record.TransactionYearMethod: List of yearly payment method usage statistics.
	//   - error: Error if any occurs during query.
	GetYearlyPaymentMethods(ctx context.Context, year int) ([]*record.TransactionYearMethod, error)
}

type TransactionStatsAmountRepository interface {
	// GetMonthlyAmounts retrieves monthly transaction amount statistics.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*record.TransactionMonthAmount: List of monthly transaction amount statistics.
	//   - error: Error if any occurs during query.
	GetMonthlyAmounts(ctx context.Context, year int) ([]*record.TransactionMonthAmount, error)

	// GetYearlyAmounts retrieves yearly transaction amount statistics.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*record.TransactionYearlyAmount: List of yearly transaction amount statistics.
	//   - error: Error if any occurs during query.
	GetYearlyAmounts(ctx context.Context, year int) ([]*record.TransactionYearlyAmount, error)
}
