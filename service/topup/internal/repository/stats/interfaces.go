package topupstatsrepository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type TopupStatsAmountRepository interface {
	// GetMonthlyTopupAmounts retrieves monthly statistics of topup amounts.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the data is requested.
	//
	// Returns:
	//   - []*record.TopupMonthAmount: List of monthly topup amount statistics.
	//   - error: Error if the query fails.
	GetMonthlyTopupAmounts(ctx context.Context, year int) ([]*record.TopupMonthAmount, error)

	// GetYearlyTopupAmounts retrieves yearly statistics of topup amounts.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the data is requested.
	//
	// Returns:
	//   - []*record.TopupYearlyAmount: List of yearly topup amount statistics.
	//   - error: Error if the query fails.
	GetYearlyTopupAmounts(ctx context.Context, year int) ([]*record.TopupYearlyAmount, error)
}

type TopupStatsStatusRepository interface {
	// GetMonthTopupStatusSuccess retrieves monthly statistics of successful topups.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing month, year, and method for filtering.
	//
	// Returns:
	//   - []*record.TopupRecordMonthStatusSuccess: List of monthly successful topup records.
	//   - error: Error if the query fails.
	GetMonthTopupStatusSuccess(ctx context.Context, req *requests.MonthTopupStatus) ([]*record.TopupRecordMonthStatusSuccess, error)

	// GetYearlyTopupStatusSuccess retrieves yearly statistics of successful topups.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the statistics are requested.
	//
	// Returns:
	//   - []*record.TopupRecordYearStatusSuccess: List of yearly successful topup records.
	//   - error: Error if the query fails.
	GetYearlyTopupStatusSuccess(ctx context.Context, year int) ([]*record.TopupRecordYearStatusSuccess, error)

	// GetMonthTopupStatusFailed retrieves monthly statistics of failed topups.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing month, year, and method for filtering.
	//
	// Returns:
	//   - []*record.TopupRecordMonthStatusFailed: List of monthly failed topup records.
	//   - error: Error if the query fails.
	GetMonthTopupStatusFailed(ctx context.Context, req *requests.MonthTopupStatus) ([]*record.TopupRecordMonthStatusFailed, error)

	// GetYearlyTopupStatusFailed retrieves yearly statistics of failed topups.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the statistics are requested.
	//
	// Returns:
	//   - []*record.TopupRecordYearStatusFailed: List of yearly failed topup records.
	//   - error: Error if the query fails.
	GetYearlyTopupStatusFailed(ctx context.Context, year int) ([]*record.TopupRecordYearStatusFailed, error)
}

type TOpupStatsMethodRepository interface {
	// GetMonthlyTopupMethods retrieves monthly statistics of topup methods used.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which monthly method statistics are requested.
	//
	// Returns:
	//   - []*record.TopupMonthMethod: List of monthly topup method usage.
	//   - error: Error if the query fails.
	GetMonthlyTopupMethods(ctx context.Context, year int) ([]*record.TopupMonthMethod, error)

	// GetYearlyTopupMethods retrieves yearly statistics of topup methods used.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which yearly method statistics are requested.
	//
	// Returns:
	//   - []*record.TopupYearlyMethod: List of yearly topup method usage.
	//   - error: Error if the query fails.
	GetYearlyTopupMethods(ctx context.Context, year int) ([]*record.TopupYearlyMethod, error)
}
