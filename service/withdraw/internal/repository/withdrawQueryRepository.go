package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	withdraw_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors/repository"
)

type withdrawQueryRepository struct {
	db *db.Queries
}

func NewWithdrawQueryRepository(db *db.Queries) WithdrawQueryRepository {
	return &withdrawQueryRepository{
		db: db,
	}
}

func (r *withdrawQueryRepository) FindAll(ctx context.Context, req *requests.FindAllWithdraws) ([]*db.GetWithdrawsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetWithdrawsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	withdraw, err := r.db.GetWithdraws(ctx, reqDb)

	if err != nil {
		return nil, withdraw_errors.ErrFindAllWithdrawsFailed
	}

	return withdraw, nil

}

func (r *withdrawQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllWithdraws) ([]*db.GetActiveWithdrawsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveWithdrawsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveWithdraws(ctx, reqDb)

	if err != nil {
		return nil, withdraw_errors.ErrFindActiveWithdrawsFailed
	}

	return res, nil
}

func (r *withdrawQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllWithdraws) ([]*db.GetTrashedWithdrawsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedWithdrawsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedWithdraws(ctx, reqDb)

	if err != nil {
		return nil, withdraw_errors.ErrFindTrashedWithdrawsFailed
	}

	return res, nil
}

func (r *withdrawQueryRepository) FindAllByCardNumber(ctx context.Context, req *requests.FindAllWithdrawCardNumber) ([]*db.GetWithdrawsByCardNumberRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetWithdrawsByCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    req.Search,
		Limit:      int32(req.PageSize),
		Offset:     int32(offset),
	}

	withdraw, err := r.db.GetWithdrawsByCardNumber(ctx, reqDb)

	if err != nil {
		return nil, withdraw_errors.ErrFindWithdrawsByCardNumberFailed
	}

	return withdraw, nil

}

func (r *withdrawQueryRepository) FindById(ctx context.Context, id int) (*db.GetWithdrawByIDRow, error) {
	withdraw, err := r.db.GetWithdrawByID(ctx, int32(id))

	if err != nil {
		return nil, withdraw_errors.ErrFindWithdrawByIdFailed
	}

	return withdraw, nil
}
