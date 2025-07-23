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

// NewCardRepository creates a new instance of cardRepository.
//
// Parameters:
//   - db: A pointer to the db.Queries object for executing database queries.
//   - mapper: A CardRecordMapping that provides methods to map database rows to Card domain models.
//
// Returns:
//   - A pointer to the newly created cardRepository instance.
func NewCardRepository(db *db.Queries, mapper recordmapper.CardQueryRecordMapper) CardRepository {
	return &cardRepository{
		db:     db,
		mapper: mapper,
	}
}

// FindUserCardByCardNumber retrieves a card record along with associated user email by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - card_number: The card number to look up.
//
// Returns:
//   - *record.CardEmailRecord: The card and user email record if found.
//   - error: An error if the operation fails or no record is found.
func (r *cardRepository) FindUserCardByCardNumber(ctx context.Context, card_number string) (*record.CardEmailRecord, error) {
	res, err := r.db.GetUserEmailByCardNumber(ctx, card_number)

	if err != nil {
		return nil, card_errors.ErrFindCardByCardNumberFailed
	}

	so := r.mapper.ToCardEmailRecord(res)

	return so, nil
}
