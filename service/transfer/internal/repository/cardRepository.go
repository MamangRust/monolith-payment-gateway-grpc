package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/card"
)

// cardRepository is a struct that implements the CardRepository interface
type cardRepository struct {
	db     *db.Queries
	mapper recordmapper.CardQueryRecordMapper
}

// NewCardRepository initializes a new instance of cardRepository with the provided
// database queries, context, and card record mapper. This repository is responsible
// for executing query operations related to card records in the database.
//
// Parameters:
//   - db: A pointer to the db.Queries object for executing database queries.
//   - mapper: A CardRecordMapping that provides methods to map database rows to
//     Card domain models.
//
// Returns:
//   - A pointer to the newly created cardRepository instance.
func NewCardRepository(db *db.Queries, mapper recordmapper.CardQueryRecordMapper) *cardRepository {
	return &cardRepository{
		db:     db,
		mapper: mapper,
	}
}

// FindCardByCardNumber retrieves a card record by its card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - card_number: The card number to search for.
//
// Returns:
//   - *record.CardRecord: The retrieved card record.
//   - error: Any error encountered during the operation.
func (r *cardRepository) FindCardByCardNumber(ctx context.Context, card_number string) (*record.CardRecord, error) {
	res, err := r.db.GetCardByCardNumber(ctx, card_number)

	if err != nil {
		return nil, card_errors.ErrFindCardByCardNumberFailed
	}

	so := r.mapper.ToCardRecord(res)

	return so, nil
}

// FindUserCardByCardNumber retrieves user-related card info by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - card_number: The card number to search for.
//
// Returns:
//   - *record.CardEmailRecord: The card and user email info.
//   - error: Any error encountered during the operation.
func (r *cardRepository) FindUserCardByCardNumber(ctx context.Context, card_number string) (*record.CardEmailRecord, error) {
	res, err := r.db.GetUserEmailByCardNumber(ctx, card_number)

	if err != nil {
		return nil, card_errors.ErrFindCardByCardNumberFailed
	}

	so := r.mapper.ToCardEmailRecord(res)

	return so, nil
}
