package repository

import (
	"context"
	"database/sql"

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
//   - card_number: The card number to search for.
//
// Returns:
//   - *record.SaldoRecord: The found saldo record.
//   - error: An error if the operation fails.
func (r *saldoRepository) FindByCardNumber(ctx context.Context, card_number string) (*record.SaldoRecord, error) {
	res, err := r.db.GetSaldoByCardNumber(ctx, card_number)

	if err != nil {
		return nil, saldo_errors.ErrFindSaldoByCardNumberFailed
	}

	so := r.mapper.ToSaldoRecord(res)

	return so, nil
}

// UpdateSaldoBalance updates the saldo balance for a specific card.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The update request containing balance data.
//
// Returns:
//   - *record.SaldoRecord: The updated saldo record.
//   - error: An error if the update fails.
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

// UpdateSaldoWithdraw updates the saldo balance after a withdrawal.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The update request related to withdrawal.
//
// Returns:
//   - *record.SaldoRecord: The updated saldo record.
//   - error: An error if the update fails.
func (r *saldoRepository) UpdateSaldoWithdraw(ctx context.Context, request *requests.UpdateSaldoWithdraw) (*record.SaldoRecord, error) {
	withdrawAmount := sql.NullInt32{
		Int32: int32(*request.WithdrawAmount),
		Valid: request.WithdrawAmount != nil,
	}
	var withdrawTime sql.NullTime
	if request.WithdrawTime != nil {
		withdrawTime = sql.NullTime{
			Time:  *request.WithdrawTime,
			Valid: true,
		}
	}

	req := db.UpdateSaldoWithdrawParams{
		CardNumber:     request.CardNumber,
		WithdrawAmount: withdrawAmount,
		WithdrawTime:   withdrawTime,
	}

	res, err := r.db.UpdateSaldoWithdraw(ctx, req)

	if err != nil {
		return nil, saldo_errors.ErrUpdateSaldoWithdrawFailed
	}

	so := r.mapper.ToSaldoRecord(res)

	return so, nil
}
