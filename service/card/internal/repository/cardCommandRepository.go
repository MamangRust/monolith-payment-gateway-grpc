package repository

import (
	"context"
	"fmt"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/randomvcc"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type cardCommandRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.CardRecordMapping
}

func NewCardCommandRepository(db *db.Queries, ctx context.Context, mapping recordmapper.CardRecordMapping) *cardCommandRepository {
	return &cardCommandRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *cardCommandRepository) CreateCard(request *requests.CreateCardRequest) (*record.CardRecord, error) {
	number, err := randomvcc.RandomCardNumber()

	if err != nil {
		return nil, fmt.Errorf("failed to generate card number: %w", err)
	}

	req := db.CreateCardParams{
		UserID:       int32(request.UserID),
		CardNumber:   number,
		CardType:     request.CardType,
		ExpireDate:   request.ExpireDate,
		Cvv:          request.CVV,
		CardProvider: request.CardProvider,
	}

	res, err := r.db.CreateCard(r.ctx, req)

	if err != nil {
		return nil, card_errors.ErrCreateCardFailed
	}

	return r.mapping.ToCardRecord(res), nil
}
func (r *cardCommandRepository) UpdateCard(request *requests.UpdateCardRequest) (*record.CardRecord, error) {
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

func (r *cardCommandRepository) TrashedCard(card_id int) (*record.CardRecord, error) {
	res, err := r.db.TrashCard(r.ctx, int32(card_id))

	if err != nil {
		return nil, card_errors.ErrTrashCardFailed
	}

	return r.mapping.ToCardRecord(res), nil
}

func (r *cardCommandRepository) RestoreCard(card_id int) (*record.CardRecord, error) {
	res, err := r.db.RestoreCard(r.ctx, int32(card_id))

	if err != nil {
		return nil, card_errors.ErrRestoreCardFailed
	}

	return r.mapping.ToCardRecord(res), nil
}

func (r *cardCommandRepository) DeleteCardPermanent(card_id int) (bool, error) {
	err := r.db.DeleteCardPermanently(r.ctx, int32(card_id))

	if err != nil {
		return false, card_errors.ErrDeleteCardPermanentFailed
	}

	return true, nil
}

func (r *cardCommandRepository) RestoreAllCard() (bool, error) {
	err := r.db.RestoreAllCards(r.ctx)

	if err != nil {
		return false, card_errors.ErrRestoreAllCardsFailed
	}

	return true, nil
}

func (r *cardCommandRepository) DeleteAllCardPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentCards(r.ctx)

	if err != nil {
		return false, card_errors.ErrDeleteAllCardsPermanentFailed
	}

	return true, nil
}
