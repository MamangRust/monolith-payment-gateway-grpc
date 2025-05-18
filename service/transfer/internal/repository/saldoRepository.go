package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type saldoRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.SaldoRecordMapping
}

func NewSaldoRepository(db *db.Queries, ctx context.Context, mapping recordmapper.SaldoRecordMapping) *saldoRepository {
	return &saldoRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *saldoRepository) FindByCardNumber(card_number string) (*record.SaldoRecord, error) {
	res, err := r.db.GetSaldoByCardNumber(r.ctx, card_number)

	if err != nil {
		return nil, saldo_errors.ErrFindSaldoByCardNumberFailed
	}

	return r.mapping.ToSaldoRecord(res), nil
}

func (r *saldoRepository) UpdateSaldoBalance(request *requests.UpdateSaldoBalance) (*record.SaldoRecord, error) {
	req := db.UpdateSaldoBalanceParams{
		CardNumber:   request.CardNumber,
		TotalBalance: int32(request.TotalBalance),
	}

	res, err := r.db.UpdateSaldoBalance(r.ctx, req)

	if err != nil {
		return nil, saldo_errors.ErrUpdateSaldoBalanceFailed
	}

	return r.mapping.ToSaldoRecord(res), nil
}
