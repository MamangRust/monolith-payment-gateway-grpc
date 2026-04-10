package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/repository"
	"github.com/jackc/pgx/v5/pgtype"
)

type saldoCommandRepository struct {
	db *db.Queries
}

func NewSaldoCommandRepository(db *db.Queries) SaldoCommandRepository {
	return &saldoCommandRepository{
		db: db,
	}
}

func (r *saldoCommandRepository) CreateSaldo(ctx context.Context, request *requests.CreateSaldoRequest) (*db.CreateSaldoRow, error) {
	req := db.CreateSaldoParams{
		CardNumber:   request.CardNumber,
		TotalBalance: int32(request.TotalBalance),
	}
	res, err := r.db.CreateSaldo(ctx, req)

	if err != nil {
		return nil, saldo_errors.ErrCreateSaldoFailed
	}

	return res, nil
}

func (r *saldoCommandRepository) UpdateSaldo(ctx context.Context, request *requests.UpdateSaldoRequest) (*db.UpdateSaldoRow, error) {
	req := db.UpdateSaldoParams{
		SaldoID:      int32(*request.SaldoID),
		CardNumber:   request.CardNumber,
		TotalBalance: int32(request.TotalBalance),
	}

	res, err := r.db.UpdateSaldo(ctx, req)

	if err != nil {
		return nil, saldo_errors.ErrUpdateSaldoFailed
	}

	return res, nil
}

func (r *saldoCommandRepository) UpdateSaldoBalance(ctx context.Context, request *requests.UpdateSaldoBalance) (*db.UpdateSaldoBalanceRow, error) {
	req := db.UpdateSaldoBalanceParams{
		CardNumber:   request.CardNumber,
		TotalBalance: int32(request.TotalBalance),
	}

	res, err := r.db.UpdateSaldoBalance(ctx, req)

	if err != nil {
		return nil, saldo_errors.ErrUpdateSaldoBalanceFailed
	}

	return res, nil
}

func (r *saldoCommandRepository) UpdateSaldoWithdraw(ctx context.Context, request *requests.UpdateSaldoWithdraw) (*db.UpdateSaldoWithdrawRow, error) {
	var withdrawAmount pgtype.Int4
	if request.WithdrawAmount != nil {
		withdrawAmount = pgtype.Int4{
			Int32: int32(*request.WithdrawAmount),
			Valid: true,
		}
	}

	var withdrawTime pgtype.Timestamp
	if request.WithdrawTime != nil {
		withdrawTime = pgtype.Timestamp{
			Time:  *request.WithdrawTime,
			Valid: true,
		}
	}

	req := db.UpdateSaldoWithdrawParams{
		CardNumber:     request.CardNumber,
		WithdrawAmount: &withdrawAmount.Int32,
		WithdrawTime:   withdrawTime,
	}

	res, err := r.db.UpdateSaldoWithdraw(ctx, req)
	if err != nil {
		return nil, saldo_errors.ErrUpdateSaldoWithdrawFailed
	}

	return res, nil
}

func (r *saldoCommandRepository) TrashedSaldo(ctx context.Context, saldo_id int) (*db.Saldo, error) {
	res, err := r.db.TrashSaldo(ctx, int32(saldo_id))
	if err != nil {
		return nil, saldo_errors.ErrTrashSaldoFailed
	}
	return res, nil
}

func (r *saldoCommandRepository) RestoreSaldo(ctx context.Context, saldo_id int) (*db.Saldo, error) {
	res, err := r.db.RestoreSaldo(ctx, int32(saldo_id))
	if err != nil {
		return nil, saldo_errors.ErrRestoreSaldoFailed
	}
	return res, nil
}

func (r *saldoCommandRepository) DeleteSaldoPermanent(ctx context.Context, saldo_id int) (bool, error) {
	err := r.db.DeleteSaldoPermanently(ctx, int32(saldo_id))
	if err != nil {
		return false, saldo_errors.ErrDeleteSaldoPermanentFailed
	}
	return true, nil
}

func (r *saldoCommandRepository) RestoreAllSaldo(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllSaldos(ctx)

	if err != nil {
		return false, saldo_errors.ErrRestoreAllSaldosFailed
	}

	return true, nil
}

func (r *saldoCommandRepository) DeleteAllSaldoPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentSaldos(ctx)

	if err != nil {
		return false, saldo_errors.ErrDeleteAllSaldosPermanentFailed
	}

	return true, nil
}
