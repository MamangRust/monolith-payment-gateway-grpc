package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type saldoStatisticsRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.SaldoRecordMapping
}

func NewSaldoStatisticsRepository(db *db.Queries, ctx context.Context, mapping recordmapper.SaldoRecordMapping) *saldoStatisticsRepository {
	return &saldoStatisticsRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *saldoStatisticsRepository) GetMonthlyTotalSaldoBalance(req *requests.MonthTotalSaldoBalance) ([]*record.SaldoMonthTotalBalance, error) {
	year := req.Year
	month := req.Month

	currentDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTotalSaldoBalance(r.ctx, db.GetMonthlyTotalSaldoBalanceParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, saldo_errors.ErrGetMonthlyTotalSaldoBalanceFailed
	}

	so := r.mapping.ToSaldoMonthTotalBalances(res)
	return so, nil
}

func (r *saldoStatisticsRepository) GetYearTotalSaldoBalance(year int) ([]*record.SaldoYearTotalBalance, error) {
	res, err := r.db.GetYearlyTotalSaldoBalances(r.ctx, int32(year))

	if err != nil {
		return nil, saldo_errors.ErrGetYearTotalSaldoBalanceFailed
	}

	so := r.mapping.ToSaldoYearTotalBalances(res)

	return so, nil
}

func (r *saldoStatisticsRepository) GetMonthlySaldoBalances(year int) ([]*record.SaldoMonthSaldoBalance, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlySaldoBalances(r.ctx, yearStart)

	if err != nil {
		return nil, saldo_errors.ErrGetMonthlySaldoBalancesFailed
	}

	so := r.mapping.ToSaldoMonthBalances(res)

	return so, nil
}

func (r *saldoStatisticsRepository) GetYearlySaldoBalances(year int) ([]*record.SaldoYearSaldoBalance, error) {
	res, err := r.db.GetYearlySaldoBalances(r.ctx, year)

	if err != nil {
		return nil, saldo_errors.ErrGetYearlySaldoBalancesFailed
	}

	so := r.mapping.ToSaldoYearSaldoBalances(res)

	return so, nil
}
