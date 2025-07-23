package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/card"
)

// cardRepository represents a repository for card operations.
type cardRepository struct {
	db     *db.Queries
	mapper recordmapper.CardQueryRecordMapper
}

// NewCardRepository initializes a new instance of cardRepository responsible
// for executing card-related database operations. It takes a database query
// executor, a context for managing request deadlines and cancellations, and a
// mapper interface to convert database rows to Card domain models.
//
// Parameters:
//   - db: A pointer to the db.Queries object used for executing database queries.
//   - mapper: A CardRecordMapping instance for mapper database rows to Card domain models.
//
// Returns:
//   - A pointer to a newly created cardRepository instance.
func NewCardRepository(db *db.Queries, mapper recordmapper.CardQueryRecordMapper) CardRepository {
	return &cardRepository{
		db:     db,
		mapper: mapper,
	}
}

// FindCardByUserId retrieves a card associated with a specific user ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - user_id: The user ID to lookup.
//
// Returns:
//   - *record.CardRecord: The card record if found.
//   - error: Error if something went wrong during the query.
func (r *cardRepository) FindCardByUserId(ctx context.Context, user_id int) (*record.CardRecord, error) {
	res, err := r.db.GetCardByUserID(ctx, int32(user_id))

	if err != nil {
		return nil, card_errors.ErrFindCardByUserIdFailed
	}

	return r.mapper.ToCardRecord(res), nil
}

// FindUserCardByCardNumber retrieves a user's card including email by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - card_number: The card number to lookup.
//
// Returns:
//   - *record.CardEmailRecord: The card and user email record.
//   - error: Error if something went wrong during the query.
func (r *cardRepository) FindUserCardByCardNumber(ctx context.Context, card_number string) (*record.CardEmailRecord, error) {
	res, err := r.db.GetUserEmailByCardNumber(ctx, card_number)

	if err != nil {
		return nil, card_errors.ErrFindCardByCardNumberFailed
	}

	return r.mapper.ToCardEmailRecord(res), nil
}

// FindCardByCardNumber retrieves a card by its card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - card_number: The card number to lookup.
//
// Returns:
//   - *record.CardRecord: The card record if found.
//   - error: Error if something went wrong during the query.
func (r *cardRepository) FindCardByCardNumber(ctx context.Context, card_number string) (*record.CardRecord, error) {
	res, err := r.db.GetCardByCardNumber(ctx, card_number)

	if err != nil {
		return nil, card_errors.ErrFindCardByCardNumberFailed
	}

	return r.mapper.ToCardRecord(res), nil
}
