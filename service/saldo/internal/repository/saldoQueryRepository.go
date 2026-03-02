package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/repository"
)

type saldoQueryRepository struct {
	db *db.Queries
}

func NewSaldoQueryRepository(db *db.Queries) SaldoQueryRepository {
	return &saldoQueryRepository{
		db: db,
	}
}

func (r *saldoQueryRepository) FindAllSaldos(ctx context.Context, req *requests.FindAllSaldos) ([]*db.GetSaldosRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetSaldosParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	saldos, err := r.db.GetSaldos(ctx, reqDb)

	if err != nil {
		return nil, saldo_errors.ErrFindAllSaldosFailed
	}

	return saldos, nil
}

func (r *saldoQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllSaldos) ([]*db.GetActiveSaldosRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveSaldosParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveSaldos(ctx, reqDb)

	if err != nil {
		return nil, saldo_errors.ErrFindActiveSaldosFailed
	}

	return res, nil
}

func (r *saldoQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllSaldos) ([]*db.GetTrashedSaldosRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedSaldosParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	saldos, err := r.db.GetTrashedSaldos(ctx, reqDb)

	if err != nil {
		return nil, saldo_errors.ErrFindTrashedSaldosFailed
	}

	return saldos, nil
}

func (r *saldoQueryRepository) FindByCardNumber(ctx context.Context, card_number string) (*db.Saldo, error) {
	res, err := r.db.GetSaldoByCardNumber(ctx, card_number)

	if err != nil {
		return nil, saldo_errors.ErrFindSaldoByCardNumberFailed
	}

	return res, nil
}

func (r *saldoQueryRepository) FindById(ctx context.Context, saldo_id int) (*db.GetSaldoByIDRow, error) {
	res, err := r.db.GetSaldoByID(ctx, int32(saldo_id))

	if err != nil {
		return nil, saldo_errors.ErrFindSaldoByIdFailed
	}

	return res, nil
}
