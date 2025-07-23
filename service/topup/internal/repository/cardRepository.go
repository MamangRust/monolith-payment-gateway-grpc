package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/card"
)

// cardRepository is a struct that implements the CardRepository interface
type cardRepository struct {
	db     *db.Queries
	mapper recordmapper.CardQueryRecordMapper
}

// NewCardRepository creates a new instance of cardRepository.
//
// Parameters:
//   - db: A pointer to the db.Queries object for executing database queries.
//   - ctx: The context to be used for database operations, allowing for cancellation and timeout.
//   - mapper: A CardRecordMapping that provides methods to map database rows to Card domain models.
//
// Returns:
//   - A new instance of cardRepository
func NewCardRepository(db *db.Queries, mapper recordmapper.CardQueryRecordMapper) CardRepository {
	return &cardRepository{
		db:     db,
		mapper: mapper,
	}
}

// FindUserCardByCardNumber retrieves card data with user email by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - card_number: The card number to search.
//
// Returns:
//   - *record.CardEmailRecord: The card record including associated user email.
//   - error: Error if retrieval fails.
func (r *cardRepository) FindCardByCardNumber(ctx context.Context, card_number string) (*record.CardRecord, error) {
	res, err := r.db.GetCardByCardNumber(ctx, card_number)

	if err != nil {
		return nil, card_errors.ErrFindCardByCardNumberFailed
	}

	return r.mapper.ToCardRecord(res), nil
}

// FindCardByCardNumber retrieves a card record by its card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - card_number: The card number to search.
//
// Returns:
//   - *record.CardRecord: The card record data.
//   - error: Error if retrieval fails.
func (r *cardRepository) FindUserCardByCardNumber(ctx context.Context, card_number string) (*record.CardEmailRecord, error) {
	res, err := r.db.GetUserEmailByCardNumber(ctx, card_number)

	if err != nil {
		return nil, card_errors.ErrFindCardByCardNumberFailed
	}

	return r.mapper.ToCardEmailRecord(res), nil
}

// UpdateCard updates an existing card's data.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The updated card information.
//
// Returns:
//   - *record.CardRecord: The updated card record.
//   - error: Error if update fails.
func (r *cardRepository) UpdateCard(ctx context.Context, request *requests.UpdateCardRequest) (*record.CardRecord, error) {
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
