package repository

import (
	"context"
	"fmt"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/randomvcc"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/card"
)

// cardCommandRepository is a struct that implements the CardCommandRepository interface
type cardCommandRepository struct {
	db     *db.Queries
	mapper recordmapper.CardCommandRecordMapper
}

// NewCardCommandRepository initializes a new instance of cardCommandRepository with the provided
// database queries, context, and card record mapper. This repository is responsible for executing
// command operations related to card records in the database.
//
// Parameters:
//   - db: A pointer to the db.Queries object for executing database queries.
//   - mapper: A CardRecordMapping that provides methods to map database rows to Card domain models.
//
// Returns:
//   - A pointer to the newly created cardCommandRepository instance.
func NewCardCommandRepository(db *db.Queries, mapper recordmapper.CardCommandRecordMapper) CardCommandRepository {
	return &cardCommandRepository{
		db:     db,
		mapper: mapper,
	}
}

// CreateCard generates a new card number, constructs a CreateCardParams object
// from the provided CreateCardRequest, and inserts a new card record into the database.
// It returns the created CardRecord or an error if the operation fails.
//
// Parameters:
//   - ctx: the context for the database operation
//   - request: A CreateCardRequest object containing the details of the card to be created.
//
// Returns:
//   - A pointer to the created CardRecord, or an error if the operation fails.
func (r *cardCommandRepository) CreateCard(ctx context.Context, request *requests.CreateCardRequest) (*record.CardRecord, error) {
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

	res, err := r.db.CreateCard(ctx, req)

	if err != nil {
		return nil, card_errors.ErrCreateCardFailed
	}

	return r.mapper.ToCardRecord(res), nil
}

// UpdateCard updates a card record in the database.
//
// Parameters:
//   - ctx: the context for the database operation
//   - request: An UpdateCardRequest object containing the details of the card to be updated.
//
// Returns:
//   - A pointer to the updated CardRecord, or an error if the operation fails.
func (r *cardCommandRepository) UpdateCard(ctx context.Context, request *requests.UpdateCardRequest) (*record.CardRecord, error) {
	req := db.UpdateCardParams{
		CardID:       int32(request.CardID),
		CardType:     request.CardType,
		ExpireDate:   request.ExpireDate,
		Cvv:          request.CVV,
		CardProvider: request.CardProvider,
	}

	res, err := r.db.UpdateCard(ctx, req)

	if err != nil {
		return nil, card_errors.ErrUpdateCardFailed
	}

	return r.mapper.ToCardRecord(res), nil
}

// TrashedCard permanently deletes a card record from the database.
//
// Parameters:
//   - ctx: the context for the database operation
//   - card_id: The ID of the card to be trashed.
//
// Returns:
//   - A pointer to the trashed CardRecord, or an error if the operation fails.
func (r *cardCommandRepository) TrashedCard(ctx context.Context, card_id int) (*record.CardRecord, error) {
	res, err := r.db.TrashCard(ctx, int32(card_id))

	if err != nil {
		return nil, card_errors.ErrTrashCardFailed
	}

	return r.mapper.ToCardRecord(res), nil
}

// RestoreCard restores a previously trashed card by setting its deleted_at field to NULL.
//
// Parameters:
//   - ctx: the context for the database operation
//   - card_id: The ID of the card to be restored.
//
// Returns:
//   - A pointer to the restored CardRecord, or an error if the operation fails.
func (r *cardCommandRepository) RestoreCard(ctx context.Context, card_id int) (*record.CardRecord, error) {
	res, err := r.db.RestoreCard(ctx, int32(card_id))

	if err != nil {
		return nil, card_errors.ErrRestoreCardFailed
	}

	return r.mapper.ToCardRecord(res), nil
}

// DeleteCardPermanent permanently deletes a card record from the database.
//
// Parameters:
//   - ctx: the context for the database operation
//   - card_id: The ID of the card to be deleted permanently.
//
// Returns:
//   - A boolean indicating if the operation was successful, and an error if the operation fails.
func (r *cardCommandRepository) DeleteCardPermanent(ctx context.Context, card_id int) (bool, error) {
	err := r.db.DeleteCardPermanently(ctx, int32(card_id))

	if err != nil {
		return false, card_errors.ErrDeleteCardPermanentFailed
	}

	return true, nil
}

// RestoreAllCard restores all previously trashed card records by setting their deleted_at fields to NULL.
//
// Parameters:
//   - ctx: the context for the database operation
//
// Returns:
//   - A boolean indicating if the operation was successful.
//   - An error if the operation fails.
func (r *cardCommandRepository) RestoreAllCard(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllCards(ctx)

	if err != nil {
		return false, card_errors.ErrRestoreAllCardsFailed
	}

	return true, nil
}

// DeleteAllCardPermanent permanently deletes all card records from the database.
// Parameters:
//   - ctx: the context for the database operation
//
// Returns:
//   - A boolean indicating if the operation was successful, and an error if the operation fails.
func (r *cardCommandRepository) DeleteAllCardPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentCards(ctx)

	if err != nil {
		return false, card_errors.ErrDeleteAllCardsPermanentFailed
	}

	return true, nil
}
