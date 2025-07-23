package transactionstatsrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/transaction/stats"
)

type transactionStatsMethodRepository struct {
	db     *db.Queries
	mapper recordmapper.TransactionStatisticMethodRecordMapper
}

func NewTransactionStatsMethodRepository(db *db.Queries, mapper recordmapper.TransactionStatisticMethodRecordMapper) TransactionStatsMethodRepository {
	return &transactionStatsMethodRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyPaymentMethods retrieves monthly statistics grouped by payment method.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*record.TransactionMonthMethod: List of monthly payment method usage statistics.
//   - error: Error if any occurs during query.
func (r *transactionStatsMethodRepository) GetMonthlyPaymentMethods(ctx context.Context, year int) ([]*record.TransactionMonthMethod, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyPaymentMethods(ctx, yearStart)

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyPaymentMethodsFailed
	}

	return r.mapper.ToTransactionMonthlyMethods(res), nil
}

// GetYearlyPaymentMethods retrieves yearly statistics grouped by payment method.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*record.TransactionYearMethod: List of yearly payment method usage statistics.
//   - error: Error if any occurs during query.
func (r *transactionStatsMethodRepository) GetYearlyPaymentMethods(ctx context.Context, year int) ([]*record.TransactionYearMethod, error) {
	res, err := r.db.GetYearlyPaymentMethods(ctx, year)

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyPaymentMethodsFailed
	}

	return r.mapper.ToTransactionYearlyMethods(res), nil
}
