package transactionbycardrepository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type TransactonStatsByCardStatusRepository interface {
	// GetMonthTransactionStatusSuccessByCardNumber retrieves monthly successful transaction statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing month, year, and card number.
	//
	// Returns:
	//   - []*record.TransactionRecordMonthStatusSuccess: List of monthly success transaction stats.
	//   - error: Error if any occurs.
	GetMonthTransactionStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*record.TransactionRecordMonthStatusSuccess, error)

	// GetYearlyTransactionStatusSuccessByCardNumber retrieves yearly successful transaction statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//
	// Returns:
	//   - []*record.TransactionRecordYearStatusSuccess: List of yearly success transaction stats.
	//   - error: Error if any occurs.
	GetYearlyTransactionStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*record.TransactionRecordYearStatusSuccess, error)

	// GetMonthTransactionStatusFailedByCardNumber retrieves monthly failed transaction statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing month, year, and card number.
	//
	// Returns:
	//   - []*record.TransactionRecordMonthStatusFailed: List of monthly failed transaction stats.
	//   - error: Error if any occurs.
	GetMonthTransactionStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*record.TransactionRecordMonthStatusFailed, error)

	// GetYearlyTransactionStatusFailedByCardNumber retrieves yearly failed transaction statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//
	// Returns:
	//   - []*record.TransactionRecordYearStatusFailed: List of yearly failed transaction stats.
	//   - error: Error if any occurs.
	GetYearlyTransactionStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*record.TransactionRecordYearStatusFailed, error)
}

type TransactionStatsByCardMethodRepository interface {
	// GetMonthlyPaymentMethodsByCardNumber retrieves monthly transaction method statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//
	// Returns:
	//   - []*record.TransactionMonthMethod: List of monthly payment method usage.
	//   - error: Error if any occurs.
	GetMonthlyPaymentMethodsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*record.TransactionMonthMethod, error)

	// GetYearlyPaymentMethodsByCardNumber retrieves yearly transaction method statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//
	// Returns:
	//   - []*record.TransactionYearMethod: List of yearly payment method usage.
	//   - error: Error if any occurs.
	GetYearlyPaymentMethodsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*record.TransactionYearMethod, error)
}

type TransactonStatsByCardAmountRepository interface {
	// GetMonthlyAmountsByCardNumber retrieves monthly transaction amount statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//
	// Returns:
	//   - []*record.TransactionMonthAmount: List of monthly transaction amounts.
	//   - error: Error if any occurs.
	GetMonthlyAmountsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*record.TransactionMonthAmount, error)

	// GetYearlyAmountsByCardNumber retrieves yearly transaction amount statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//
	// Returns:
	//   - []*record.TransactionYearlyAmount: List of yearly transaction amounts.
	//   - error: Error if any occurs.
	GetYearlyAmountsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*record.TransactionYearlyAmount, error)
}
