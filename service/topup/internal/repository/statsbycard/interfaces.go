package topupstatsbycardrepository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type TopupStatsByCardAmountRepository interface {
	// GetMonthlyTopupAmountsByCardNumber retrieves monthly topup amount statistics for a specific card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//
	// Returns:
	//   - []*record.TopupMonthAmount: List of monthly topup amount data.
	//   - error: Error if the query fails.
	GetMonthlyTopupAmountsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*record.TopupMonthAmount, error)

	// GetYearlyTopupAmountsByCardNumber retrieves yearly topup amount statistics for a specific card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//
	// Returns:
	//   - []*record.TopupYearlyAmount: List of yearly topup amount data.
	//   - error: Error if the query fails.
	GetYearlyTopupAmountsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*record.TopupYearlyAmount, error)
}

type TopupStatsByCardStatusRepository interface {
	// GetMonthTopupStatusSuccessByCardNumber retrieves monthly statistics of successful topups for a specific card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing month, year, and card number.
	//
	// Returns:
	//   - []*record.TopupRecordMonthStatusSuccess: List of monthly successful topup records.
	//   - error: Error if the query fails.
	GetMonthTopupStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*record.TopupRecordMonthStatusSuccess, error)

	// GetYearlyTopupStatusSuccessByCardNumber retrieves yearly statistics of successful topups for a specific card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//
	// Returns:
	//   - []*record.TopupRecordYearStatusSuccess: List of yearly successful topup records.
	//   - error: Error if the query fails.
	GetYearlyTopupStatusSuccessByCardNumber(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*record.TopupRecordYearStatusSuccess, error)

	// GetMonthTopupStatusFailedByCardNumber retrieves monthly statistics of failed topups for a specific card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing month, year, and card number.
	//
	// Returns:
	//   - []*record.TopupRecordMonthStatusFailed: List of monthly failed topup records.
	//   - error: Error if the query fails.
	GetMonthTopupStatusFailedByCardNumber(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*record.TopupRecordMonthStatusFailed, error)

	// GetYearlyTopupStatusFailedByCardNumber retrieves yearly statistics of failed topups for a specific card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//
	// Returns:
	//   - []*record.TopupRecordYearStatusFailed: List of yearly failed topup records.
	//   - error: Error if the query fails.
	GetYearlyTopupStatusFailedByCardNumber(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*record.TopupRecordYearStatusFailed, error)
}

type TopupStatsByCardMethodRepository interface {
	// GetMonthlyTopupMethodsByCardNumber retrieves monthly topup method statistics for a specific card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//
	// Returns:
	//   - []*record.TopupMonthMethod: List of monthly topup method usage.
	//   - error: Error if the query fails.
	GetMonthlyTopupMethodsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*record.TopupMonthMethod, error)

	// GetYearlyTopupMethodsByCardNumber retrieves yearly topup method statistics for a specific card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//
	// Returns:
	//   - []*record.TopupYearlyMethod: List of yearly topup method usage.
	//   - error: Error if the query fails.
	GetYearlyTopupMethodsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*record.TopupYearlyMethod, error)
}
