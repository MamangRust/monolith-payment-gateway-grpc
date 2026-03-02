package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/repository"
	"github.com/jackc/pgx/v5/pgtype"
)

type saldoRepository struct {
	db *db.Queries
}

func NewSaldoRepository(db *db.Queries) SaldoRepository {
	return &saldoRepository{
		db: db,
	}
}

func (r *saldoRepository) FindByCardNumber(ctx context.Context, card_number string) (*db.Saldo, error) {
	res, err := r.db.GetSaldoByCardNumber(ctx, card_number)

	if err != nil {
		return nil, saldo_errors.ErrFindSaldoByCardNumberFailed
	}

	return res, nil
}

func (r *saldoRepository) UpdateSaldoBalance(ctx context.Context, request *requests.UpdateSaldoBalance) (*db.UpdateSaldoBalanceRow, error) {
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

func (r *saldoRepository) UpdateSaldoWithdraw(ctx context.Context, request *requests.UpdateSaldoWithdraw) (*db.UpdateSaldoWithdrawRow, error) {
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
