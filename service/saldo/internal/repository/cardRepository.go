package repository

import (
	"context"
	"fmt"

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
func NewCardRepository(db *db.Queries, mapper recordmapper.CardQueryRecordMapper) CardRepository {
	return &cardRepository{
		db:     db,
		mapper: mapper,
	}
}

// FindCardByCardNumber retrieves a card record based on the card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - card_number: The card number to look up.
//
// Returns:
//   - *record.CardRecord: The found card record if exists.
//   - error: An error if the card is not found or the query fails.
func (r *cardRepository) FindCardByCardNumber(ctx context.Context, card_number string) (*record.CardRecord, error) {
	res, err := r.db.GetCardByCardNumber(ctx, card_number)

	fmt.Println("hello res", res)
	fmt.Println("hello err", err)

	if err != nil {
		return nil, card_errors.ErrFindCardByCardNumberFailed
	}

	return r.mapper.ToCardRecord(res), nil
}
