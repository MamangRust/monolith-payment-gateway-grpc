package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
)

type cardQueryRepository struct {
	db *db.Queries
}

func NewCardQueryRepository(db *db.Queries) CardQueryRepository {
	return &cardQueryRepository{
		db: db,
	}
}

func (r *cardQueryRepository) FindAllCards(ctx context.Context, req *requests.FindAllCards) ([]*db.GetCardsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCardsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	cards, err := r.db.GetCards(ctx, reqDb)

	if err != nil {
		return nil, card_errors.ErrFindAllCardsFailed
	}

	return cards, nil
}

func (r *cardQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllCards) ([]*db.GetActiveCardsWithCountRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveCardsWithCountParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveCardsWithCount(ctx, reqDb)

	if err != nil {
		return nil, card_errors.ErrFindActiveCardsFailed
	}

	return res, nil
}

func (r *cardQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllCards) ([]*db.GetTrashedCardsWithCountRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedCardsWithCountParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedCardsWithCount(ctx, reqDb)

	if err != nil {
		return nil, card_errors.ErrFindTrashedCardsFailed
	}

	return res, nil
}

func (r *cardQueryRepository) FindById(ctx context.Context, card_id int) (*db.GetCardByIDRow, error) {
	res, err := r.db.GetCardByID(ctx, int32(card_id))

	if err != nil {
		return nil, card_errors.ErrFindCardByIdFailed
	}

	return res, nil
}

func (r *cardQueryRepository) FindCardByUserId(ctx context.Context, user_id int) (*db.GetCardByUserIDRow, error) {
	res, err := r.db.GetCardByUserID(ctx, int32(user_id))

	if err != nil {
		return nil, card_errors.ErrFindCardByUserIdFailed
	}

	return res, nil
}

func (r *cardQueryRepository) FindCardByCardNumber(ctx context.Context, card_number string) (*db.GetCardByCardNumberRow, error) {
	res, err := r.db.GetCardByCardNumber(ctx, card_number)

	if err != nil {
		return nil, card_errors.ErrFindCardByCardNumberFailed
	}

	return res, nil
}
