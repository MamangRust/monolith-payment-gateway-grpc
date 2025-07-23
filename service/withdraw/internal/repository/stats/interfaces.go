package withdrawstatsrepository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type WithdrawStatsStatusRepository interface {
	// GetMonthWithdrawStatusSuccess retrieves monthly withdraw statistics with status "success".
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing month, year, and additional filters.
	//
	// Returns:
	//   - []*record.WithdrawRecordMonthStatusSuccess: List of successful monthly withdraw records.
	//   - error: An error if the operation fails.
	GetMonthWithdrawStatusSuccess(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*record.WithdrawRecordMonthStatusSuccess, error)

	// GetYearlyWithdrawStatusSuccess retrieves yearly withdraw statistics with status "success".
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*record.WithdrawRecordYearStatusSuccess: List of successful yearly withdraw records.
	//   - error: An error if the operation fails.
	GetYearlyWithdrawStatusSuccess(ctx context.Context, year int) ([]*record.WithdrawRecordYearStatusSuccess, error)

	// GetMonthWithdrawStatusFailed retrieves monthly withdraw statistics with status "failed".
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing month, year, and additional filters.
	//
	// Returns:
	//   - []*record.WithdrawRecordMonthStatusFailed: List of failed monthly withdraw records.
	//   - error: An error if the operation fails.
	GetMonthWithdrawStatusFailed(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*record.WithdrawRecordMonthStatusFailed, error)

	// GetYearlyWithdrawStatusFailed retrieves yearly withdraw statistics with status "failed".
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*record.WithdrawRecordYearStatusFailed: List of failed yearly withdraw records.
	//   - error: An error if the operation fails.
	GetYearlyWithdrawStatusFailed(ctx context.Context, year int) ([]*record.WithdrawRecordYearStatusFailed, error)
}

type WithdrawStatsAmountRepository interface {
	// GetMonthlyWithdraws retrieves the total amount of withdraws grouped by month for the given year.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which monthly amounts are requested.
	//
	// Returns:
	//   - []*record.WithdrawMonthlyAmount: List of monthly withdraw amounts.
	//   - error: An error if the operation fails.
	GetMonthlyWithdraws(ctx context.Context, year int) ([]*record.WithdrawMonthlyAmount, error)

	// GetYearlyWithdraws retrieves the total amount of withdraws grouped by year.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which yearly amounts are requested.
	//
	// Returns:
	//   - []*record.WithdrawYearlyAmount: List of yearly withdraw amounts.
	//   - error: An error if the operation fails.
	GetYearlyWithdraws(ctx context.Context, year int) ([]*record.WithdrawYearlyAmount, error)
}
