package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/saldo"
)

// saldoRepository is a repository for handling saldo operations.
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

// FindByCardNumber retrieves saldo information by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - card_number: The card number to lookup.
//
// Returns:
//   - *record.SaldoRecord: The saldo record if found.
//   - error: Error if something went wrong during the query.
func (r *saldoRepository) FindByCardNumber(ctx context.Context, card_number string) (*record.SaldoRecord, error) {
	res, err := r.db.GetSaldoByCardNumber(ctx, card_number)

	if err != nil {
		return nil, saldo_errors.ErrFindSaldoByCardNumberFailed
	}

	return r.mapper.ToSaldoRecord(res), nil
}

// UpdateSaldoBalance updates the saldo balance based on the given request.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The update request containing new saldo balance data.
//
// Returns:
//   - *record.SaldoRecord: The updated saldo record.
//   - error: Error if something went wrong during the update.
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
