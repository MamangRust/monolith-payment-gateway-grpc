package transferstatsrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/transfer/stats"
)

type transferStatsAmountRepository struct {
	db     *db.Queries
	mapper recordmapper.TransferStatisticAmountRecordMapper
}

func NewTransferStatsAmountRepository(db *db.Queries, mapper recordmapper.TransferStatisticAmountRecordMapper) TransferStatsAmountRepository {
	return &transferStatsAmountRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyTransferAmounts retrieves transfer amount statistics grouped by month.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which the monthly amount data is requested.
//
// Returns:
//   - []*record.TransferMonthAmount: List of monthly transfer amount records.
//   - error: Any error encountered during the operation.
func (r *transferStatsAmountRepository) GetMonthlyTransferAmounts(ctx context.Context, year int) ([]*record.TransferMonthAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransferAmounts(ctx, yearStart)

	if err != nil {
		return nil, transfer_errors.ErrGetMonthlyTransferAmountsFailed
	}

	so := r.mapper.ToTransferMonthAmounts(res)

	return so, nil
}

// GetYearlyTransferAmounts retrieves transfer amount statistics grouped by year.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which the yearly amount data is requested.
//
// Returns:
//   - []*record.TransferYearAmount: List of yearly transfer amount records.
//   - error: Any error encountered during the operation.
func (r *transferStatsAmountRepository) GetYearlyTransferAmounts(ctx context.Context, year int) ([]*record.TransferYearAmount, error) {
	res, err := r.db.GetYearlyTransferAmounts(ctx, year)

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferAmountsFailed
	}

	so := r.mapper.ToTransferYearAmounts(res)

	return so, nil
}
