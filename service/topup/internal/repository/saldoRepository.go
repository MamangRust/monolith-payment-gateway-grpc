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

// NewSaldoRepository initializes a new instance of saldoRepository with the provided
// database queries, context, and saldo record mapper. This repository is responsible for
// executing operations related to saldo records in the database.
//
// Parameters:
//   - db: A pointer to the db.Queries object for executing database queries.
//   - mapper: A SaldoRecordMapping that provides methods to map database rows to SaldoRecord domain models.
//
// Returns:
//   - A pointer to the newly created saldoRepository instance.
func NewSaldoRepository(db *db.Queries, mapper recordmapper.SaldoQueryRecordMapping) SaldoRepository {
	return &saldoRepository{
		db:     db,
		mapper: mapper,
	}
}

// FindByCardNumber retrieves a saldo record by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - card_number: The unique card number to find the saldo.
//
// Returns:
//   - *record.SaldoRecord: The saldo record found.
//   - error: Error if the query fails or saldo is not found.
func (r *saldoRepository) FindByCardNumber(ctx context.Context, card_number string) (*record.SaldoRecord, error) {
	res, err := r.db.GetSaldoByCardNumber(ctx, card_number)

	if err != nil {
		return nil, saldo_errors.ErrFindSaldoByCardNumberFailed
	}

	return r.mapper.ToSaldoRecord(res), nil
}

// UpdateSaldoBalance updates the balance of a saldo record.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The request payload containing saldo ID and new balance.
//
// Returns:
//   - *record.SaldoRecord: The updated saldo record.
//   - error: Error if the update fails.
func (r *saldoRepository) UpdateSaldoBalance(ctx context.Context, request *requests.UpdateSaldoBalance) (*record.SaldoRecord, error) {
	req := db.UpdateSaldoBalanceParams{
		CardNumber:   request.CardNumber,
		TotalBalance: int32(request.TotalBalance),
	}

	res, err := r.db.UpdateSaldoBalance(ctx, req)

	if err != nil {
		return nil, saldo_errors.ErrUpdateSaldoBalanceFailed
	}

	return r.mapper.ToSaldoRecord(res), nil
}
