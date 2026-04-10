package repository

import (
	"context"
	"fmt"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/randomvcc"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	"github.com/jackc/pgx/v5/pgtype"
)

type cardCommandRepository struct {
	db *db.Queries
}

func NewCardCommandRepository(db *db.Queries) CardCommandRepository {
	return &cardCommandRepository{
		db: db,
	}
}

func (r *cardCommandRepository) CreateCard(ctx context.Context, request *requests.CreateCardRequest) (*db.CreateCardRow, error) {
	number, err := randomvcc.RandomCardNumber()
	if err != nil {
		return nil, fmt.Errorf("failed to generate card number: %w", err)
	}

	expireDate := pgtype.Date{
		Time:  request.ExpireDate,
		Valid: true,
	}

	req := db.CreateCardParams{
		UserID:       int32(request.UserID),
		CardNumber:   number,
		CardType:     request.CardType,
		ExpireDate:   expireDate,
		Cvv:          request.CVV,
		CardProvider: request.CardProvider,
	}

	res, err := r.db.CreateCard(ctx, req)
	if err != nil {
		return nil, card_errors.ErrCreateCardFailed
	}

	return res, nil
}

func (r *cardCommandRepository) UpdateCard(ctx context.Context, request *requests.UpdateCardRequest) (*db.UpdateCardRow, error) {
	expireDate := pgtype.Date{
		Time:  request.ExpireDate,
		Valid: true,
	}

	req := db.UpdateCardParams{
		CardID:       int32(request.CardID),
		CardType:     request.CardType,
		ExpireDate:   expireDate,
		Cvv:          request.CVV,
		CardProvider: request.CardProvider,
	}

	res, err := r.db.UpdateCard(ctx, req)
	if err != nil {
		return nil, card_errors.ErrUpdateCardFailed
	}

	return res, nil
}

func (r *cardCommandRepository) TrashedCard(ctx context.Context, card_id int) (*db.Card, error) {
	res, err := r.db.TrashCard(ctx, int32(card_id))

	if err != nil {
		return nil, card_errors.ErrTrashCardFailed
	}

	return res, nil
}

func (r *cardCommandRepository) RestoreCard(ctx context.Context, card_id int) (*db.Card, error) {
	res, err := r.db.RestoreCard(ctx, int32(card_id))

	if err != nil {
		return nil, card_errors.ErrRestoreCardFailed
	}

	return res, nil
}

func (r *cardCommandRepository) DeleteCardPermanent(ctx context.Context, card_id int) (bool, error) {
	err := r.db.DeleteCardPermanently(ctx, int32(card_id))

	if err != nil {
		return false, card_errors.ErrDeleteCardPermanentFailed
	}

	return true, nil
}

func (r *cardCommandRepository) RestoreAllCard(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllCards(ctx)

	if err != nil {
		return false, card_errors.ErrRestoreAllCardsFailed
	}

	return true, nil
}

func (r *cardCommandRepository) DeleteAllCardPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentCards(ctx)

	if err != nil {
		return false, card_errors.ErrDeleteAllCardsPermanentFailed
	}

	return true, nil
}
