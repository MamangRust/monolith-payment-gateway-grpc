package transactionstatsrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/transaction/stats"
)

type transactionStatsAmountRepository struct {
	db     *db.Queries
	mapper recordmapper.TransactionStatisticAmountRecordMapper
}

func NewTransactionStatsAmountRepository(db *db.Queries, mapper recordmapper.TransactionStatisticAmountRecordMapper) TransactionStatsAmountRepository {
	return &transactionStatsAmountRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyAmounts retrieves monthly transaction amount statistics.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*record.TransactionMonthAmount: List of monthly transaction amount statistics.
//   - error: Error if any occurs during query.
func (r *transactionStatsAmountRepository) GetMonthlyAmounts(ctx context.Context, year int) ([]*record.TransactionMonthAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyAmounts(ctx, yearStart)

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountsFailed
	}

	return r.mapper.ToTransactionMonthlyAmounts(res), nil
}

// GetYearlyAmounts retrieves yearly transaction amount statistics.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*record.TransactionYearlyAmount: List of yearly transaction amount statistics.
//   - error: Error if any occurs during query.
func (r *transactionStatsAmountRepository) GetYearlyAmounts(ctx context.Context, year int) ([]*record.TransactionYearlyAmount, error) {
	res, err := r.db.GetYearlyAmounts(ctx, year)

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountsFailed
	}

	return r.mapper.ToTransactionYearlyAmounts(res), nil
}
