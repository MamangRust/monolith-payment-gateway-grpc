package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/saldo"
)

// saldoRepository is a struct that implements the SaldoRepository interface
type saldoRepository struct {
	db     *db.Queries
	mapper recordmapper.SaldoQueryRecordMapping
}

// NewSaldoRepository creates a new instance of saldoRepository, which provides
// methods for performing operations related to saldo records in the database.
//
// Parameters:
//   - db: A pointer to the db.Queries object used to execute database queries.
//   - ctx: The context for database operations, supporting cancellation and timeout.
//   - mapper: A SaldoRecordMapping for mapper database rows to domain models.
//
// Returns:
//   - A pointer to the initialized saldoRepository instance.
func NewSaldoRepository(db *db.Queries, mapper recordmapper.SaldoQueryRecordMapping) SaldoRepository {
	return &saldoRepository{
		db:     db,
		mapper: mapper,
	}
}

// FindByCardNumber retrieves a saldo record by the given card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - card_number: The card number to search the saldo for.
//
// Returns:
//   - *record.SaldoRecord: The retrieved saldo record.
//   - error: Any error encountered during the operation.
func (r *saldoRepository) FindByCardNumber(ctx context.Context, card_number string) (*record.SaldoRecord, error) {
	res, err := r.db.GetSaldoByCardNumber(ctx, card_number)

	if err != nil {
		return nil, saldo_errors.ErrFindSaldoByCardNumberFailed
	}

	so := r.mapper.ToSaldoRecord(res)

	return so, nil
}

// UpdateSaldoBalance updates the saldo balance based on the provided request.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The request containing card number and the new balance.
//
// Returns:
//   - *record.SaldoRecord: The updated saldo record.
//   - error: Any error encountered during the operation.
func (r *saldoRepository) UpdateSaldoBalance(ctx context.Context, request *requests.UpdateSaldoBalance) (*record.SaldoRecord, error) {
	req := db.UpdateSaldoBalanceParams{
		CardNumber:   request.CardNumber,
		TotalBalance: int32(request.TotalBalance),
	}

	res, err := r.db.UpdateSaldoBalance(ctx, req)

	if err != nil {
		return nil, saldo_errors.ErrUpdateSaldoBalanceFailed
	}

	so := r.mapper.ToSaldoRecord(res)

	return so, nil
}
