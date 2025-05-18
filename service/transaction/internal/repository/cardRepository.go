package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type cardRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.CardRecordMapping
}

func NewCardRepository(db *db.Queries, ctx context.Context, mapping recordmapper.CardRecordMapping) *cardRepository {
	return &cardRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *cardRepository) FindCardByUserId(user_id int) (*record.CardRecord, error) {
	res, err := r.db.GetCardByUserID(r.ctx, int32(user_id))

	if err != nil {
		return nil, card_errors.ErrFindCardByUserIdFailed
	}

	return r.mapping.ToCardRecord(res), nil
}

func (r *cardRepository) FindUserCardByCardNumber(card_number string) (*record.CardEmailRecord, error) {
	res, err := r.db.GetUserEmailByCardNumber(r.ctx, card_number)

	if err != nil {
		return nil, card_errors.ErrFindCardByCardNumberFailed
	}

	return r.mapping.ToCardEmailRecord(res), nil
}

func (r *cardRepository) FindCardByCardNumber(card_number string) (*record.CardRecord, error) {
	res, err := r.db.GetCardByCardNumber(r.ctx, card_number)

	if err != nil {
		return nil, card_errors.ErrFindCardByCardNumberFailed
	}

	return r.mapping.ToCardRecord(res), nil
}
