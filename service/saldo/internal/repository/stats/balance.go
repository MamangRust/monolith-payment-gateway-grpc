package saldostatsrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/repository"
)

type saldoStatsBalanceRepository struct {
	db *db.Queries
}

func NewSaldoStatsBalanceRepository(db *db.Queries) SaldoStatsBalanceRepository {
	return &saldoStatsBalanceRepository{
		db: db,
	}
}

func (r *saldoStatsBalanceRepository) GetMonthlySaldoBalances(ctx context.Context, year int) ([]*db.GetMonthlySaldoBalancesRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlySaldoBalances(ctx, yearStart)

	if err != nil {
		return nil, saldo_errors.ErrGetMonthlySaldoBalancesFailed
	}

	return res, nil
}

func (r *saldoStatsBalanceRepository) GetYearlySaldoBalances(ctx context.Context, year int) ([]*db.GetYearlySaldoBalancesRow, error) {
	res, err := r.db.GetYearlySaldoBalances(ctx, year)

	if err != nil {
		return nil, saldo_errors.ErrGetYearlySaldoBalancesFailed
	}

	return res, nil
}
