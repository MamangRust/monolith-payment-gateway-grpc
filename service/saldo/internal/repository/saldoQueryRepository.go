package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type saldoQueryRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.SaldoRecordMapping
}

func NewSaldoQueryRepository(db *db.Queries, ctx context.Context, mapping recordmapper.SaldoRecordMapping) *saldoQueryRepository {
	return &saldoQueryRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *saldoQueryRepository) FindAllSaldos(req *requests.FindAllSaldos) ([]*record.SaldoRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetSaldosParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	saldos, err := r.db.GetSaldos(r.ctx, reqDb)

	if err != nil {
		return nil, nil, saldo_errors.ErrFindAllSaldosFailed
	}

	var totalCount int
	if len(saldos) > 0 {
		totalCount = int(saldos[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToSaldosRecordAll(saldos), &totalCount, nil
}

func (r *saldoQueryRepository) FindByActive(req *requests.FindAllSaldos) ([]*record.SaldoRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveSaldosParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveSaldos(r.ctx, reqDb)

	if err != nil {
		return nil, nil, saldo_errors.ErrFindActiveSaldosFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToSaldosRecordActive(res), &totalCount, nil

}

func (r *saldoQueryRepository) FindByTrashed(req *requests.FindAllSaldos) ([]*record.SaldoRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedSaldosParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	saldos, err := r.db.GetTrashedSaldos(r.ctx, reqDb)

	if err != nil {
		return nil, nil, saldo_errors.ErrFindTrashedSaldosFailed
	}

	var totalCount int
	if len(saldos) > 0 {
		totalCount = int(saldos[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToSaldosRecordTrashed(saldos), &totalCount, nil
}

func (r *saldoQueryRepository) FindByCardNumber(card_number string) (*record.SaldoRecord, error) {
	res, err := r.db.GetSaldoByCardNumber(r.ctx, card_number)

	if err != nil {
		return nil, saldo_errors.ErrFindSaldoByCardNumberFailed
	}

	return r.mapping.ToSaldoRecord(res), nil
}

func (r *saldoQueryRepository) FindById(saldo_id int) (*record.SaldoRecord, error) {
	res, err := r.db.GetSaldoByID(r.ctx, int32(saldo_id))

	if err != nil {
	}

	return r.mapping.ToSaldoRecord(res), nil
}
