package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
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

func (r *cardRepository) FindCardByCardNumber(card_number string) (*record.CardRecord, error) {
	res, err := r.db.GetCardByCardNumber(r.ctx, card_number)

	if err != nil {
		return nil, card_errors.ErrFindCardByCardNumberFailed
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

func (r *cardRepository) UpdateCard(request *requests.UpdateCardRequest) (*record.CardRecord, error) {
	req := db.UpdateCardParams{
		CardID:       int32(request.CardID),
		CardType:     request.CardType,
		ExpireDate:   request.ExpireDate,
		Cvv:          request.CVV,
		CardProvider: request.CardProvider,
	}

	res, err := r.db.UpdateCard(r.ctx, req)

	if err != nil {
		return nil, card_errors.ErrUpdateCardFailed
	}

	return r.mapping.ToCardRecord(res), nil
}

func (r *cardRepository) TrashedCard(card_id int) (*record.CardRecord, error) {
	res, err := r.db.TrashCard(r.ctx, int32(card_id))

	if err != nil {
		return nil, card_errors.ErrTrashCardFailed
	}

	return r.mapping.ToCardRecord(res), nil
}

func (r *cardRepository) RestoreCard(card_id int) (*record.CardRecord, error) {
	res, err := r.db.RestoreCard(r.ctx, int32(card_id))

	if err != nil {
		return nil, card_errors.ErrRestoreCardFailed
	}

	return r.mapping.ToCardRecord(res), nil
}

func (r *cardRepository) DeleteCardPermanent(card_id int) (bool, error) {
	err := r.db.DeleteCardPermanently(r.ctx, int32(card_id))

	if err != nil {
		return false, card_errors.ErrDeleteCardPermanentFailed
	}

	return true, nil
}

func (r *cardRepository) RestoreAllCard() (bool, error) {
	err := r.db.RestoreAllCards(r.ctx)

	if err != nil {
		return false, card_errors.ErrRestoreAllCardsFailed
	}

	return true, nil
}

func (r *cardRepository) DeleteAllCardPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentCards(r.ctx)

	if err != nil {
		return false, card_errors.ErrDeleteAllCardsPermanentFailed
	}

	return true, nil
}
