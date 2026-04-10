package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/repository"
)

type topupQueryRepository struct {
	db *db.Queries
}

func NewTopupQueryRepository(db *db.Queries) TopupQueryRepository {
	return &topupQueryRepository{
		db: db,
	}
}

func (r *topupQueryRepository) FindAllTopups(ctx context.Context, req *requests.FindAllTopups) ([]*db.GetTopupsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTopupsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTopups(ctx, reqDb)

	if err != nil {
		return nil, topup_errors.ErrFindAllTopupsFailed
	}

	return res, nil
}

func (r *topupQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllTopups) ([]*db.GetActiveTopupsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveTopupsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveTopups(ctx, reqDb)

	if err != nil {
		return nil, topup_errors.ErrFindTopupsByActiveFailed
	}

	return res, nil
}

func (r *topupQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllTopups) ([]*db.GetTrashedTopupsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedTopupsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedTopups(ctx, reqDb)

	if err != nil {
		return nil, topup_errors.ErrFindTopupsByTrashedFailed
	}

	return res, nil
}

func (r *topupQueryRepository) FindAllTopupByCardNumber(ctx context.Context, req *requests.FindAllTopupsByCardNumber) ([]*db.GetTopupsByCardNumberRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTopupsByCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    req.Search,
		Limit:      int32(req.PageSize),
		Offset:     int32(offset),
	}

	res, err := r.db.GetTopupsByCardNumber(ctx, reqDb)

	if err != nil {
		return nil, topup_errors.ErrFindTopupsByCardNumberFailed
	}

	return res, nil
}

func (r *topupQueryRepository) FindById(ctx context.Context, topup_id int) (*db.GetTopupByIDRow, error) {
	res, err := r.db.GetTopupByID(ctx, int32(topup_id))
	if err != nil {
		return nil, topup_errors.ErrFindTopupByIdFailed
	}
	return res, nil
}
