package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type cardQueryRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.CardRecordMapping
}

func NewCardQueryRepository(db *db.Queries, ctx context.Context, mapping recordmapper.CardRecordMapping) *cardQueryRepository {
	return &cardQueryRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *cardQueryRepository) FindAllCards(req *requests.FindAllCards) ([]*record.CardRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCardsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	cards, err := r.db.GetCards(r.ctx, reqDb)

	if err != nil {
		return nil, nil, card_errors.ErrFindAllCardsFailed
	}

	var totalCount int

	if len(cards) > 0 {
		totalCount = int(cards[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToCardsRecord(cards), &totalCount, nil
}

func (r *cardQueryRepository) FindByActive(req *requests.FindAllCards) ([]*record.CardRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveCardsWithCountParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveCardsWithCount(r.ctx, reqDb)

	if err != nil {
		return nil, nil, card_errors.ErrFindActiveCardsFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToCardRecordsActive(res), &totalCount, nil

}

func (r *cardQueryRepository) FindByTrashed(req *requests.FindAllCards) ([]*record.CardRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedCardsWithCountParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedCardsWithCount(r.ctx, reqDb)

	if err != nil {
		return nil, nil, card_errors.ErrFindTrashedCardsFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToCardRecordsTrashed(res), &totalCount, nil
}

func (r *cardQueryRepository) FindById(card_id int) (*record.CardRecord, error) {
	res, err := r.db.GetCardByID(r.ctx, int32(card_id))

	if err != nil {
		return nil, card_errors.ErrFindCardByIdFailed
	}

	return r.mapping.ToCardRecord(res), nil
}

func (r *cardQueryRepository) FindCardByUserId(user_id int) (*record.CardRecord, error) {
	res, err := r.db.GetCardByUserID(r.ctx, int32(user_id))

	if err != nil {
		return nil, card_errors.ErrFindCardByUserIdFailed
	}

	return r.mapping.ToCardRecord(res), nil
}

func (r *cardQueryRepository) FindCardByCardNumber(card_number string) (*record.CardRecord, error) {
	res, err := r.db.GetCardByCardNumber(r.ctx, card_number)

	if err != nil {
		return nil, card_errors.ErrFindCardByCardNumberFailed
	}

	return r.mapping.ToCardRecord(res), nil
}
