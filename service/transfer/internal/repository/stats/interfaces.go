package transferstatsrepository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type TransferStatsAmountRepository interface {
	// GetMonthlyTransferAmounts retrieves transfer amount statistics grouped by month.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the monthly amount data is requested.
	//
	// Returns:
	//   - []*record.TransferMonthAmount: List of monthly transfer amount records.
	//   - error: Any error encountered during the operation.
	GetMonthlyTransferAmounts(ctx context.Context, year int) ([]*record.TransferMonthAmount, error)

	// GetYearlyTransferAmounts retrieves transfer amount statistics grouped by year.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the yearly amount data is requested.
	//
	// Returns:
	//   - []*record.TransferYearAmount: List of yearly transfer amount records.
	//   - error: Any error encountered during the operation.
	GetYearlyTransferAmounts(ctx context.Context, year int) ([]*record.TransferYearAmount, error)
}

type TransferStatsStatusRepository interface {
	// GetMonthTransferStatusSuccess retrieves successful transfer statistics per month.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The month and year filter for the statistics.
	//
	// Returns:
	//   - []*record.TransferRecordMonthStatusSuccess: List of monthly successful transfer records.
	//   - error: Any error encountered during the operation.
	GetMonthTransferStatusSuccess(ctx context.Context, req *requests.MonthStatusTransfer) ([]*record.TransferRecordMonthStatusSuccess, error)

	// GetYearlyTransferStatusSuccess retrieves successful transfer statistics per year.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the statistics are requested.
	//
	// Returns:
	//   - []*record.TransferRecordYearStatusSuccess: List of yearly successful transfer records.
	//   - error: Any error encountered during the operation.
	GetYearlyTransferStatusSuccess(ctx context.Context, year int) ([]*record.TransferRecordYearStatusSuccess, error)

	// GetMonthTransferStatusFailed retrieves failed transfer statistics per month.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The month and year filter for the statistics.
	//
	// Returns:
	//   - []*record.TransferRecordMonthStatusFailed: List of monthly failed transfer records.
	//   - error: Any error encountered during the operation.
	GetMonthTransferStatusFailed(ctx context.Context, req *requests.MonthStatusTransfer) ([]*record.TransferRecordMonthStatusFailed, error)

	// GetYearlyTransferStatusFailed retrieves failed transfer statistics per year.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the statistics are requested.
	//
	// Returns:
	//   - []*record.TransferRecordYearStatusFailed: List of yearly failed transfer records.
	//   - error: Any error encountered during the operation.
	GetYearlyTransferStatusFailed(ctx context.Context, year int) ([]*record.TransferRecordYearStatusFailed, error)
}
