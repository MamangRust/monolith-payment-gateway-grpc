package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/card"
)

// cardQueryRepository is a struct that implements the CardQueryRepository interface
type cardQueryRepository struct {
	db     *db.Queries
	mapper recordmapper.CardQueryRecordMapper
}

// NewCardQueryRepository initializes a new instance of cardQueryRepository with the provided
// database queries, context, and card record mapper. This repository is responsible for executing
// query operations related to card records in the database.
//
// Parameters:
//   - db: A pointer to the db.Queries object for executing database queries.
//   - mapper: A CardRecordMapping that provides methods to map database rows to Card domain models.
//
// Returns:
//   - A pointer to the newly created cardQueryRepository instance.
func NewCardQueryRepository(db *db.Queries, mapper recordmapper.CardQueryRecordMapper) CardQueryRepository {
	return &cardQueryRepository{
		db:     db,
		mapper: mapper,
	}
}

// FindAllCards retrieves a paginated list of card records based on the search criteria
// specified in the request. It queries the database and returns a slice of CardRecord,
// the total count of records, and an error if any occurred.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - req: A FindAllCards request object containing the search parameters
//     such as search keyword, page number, and page size.
//
// Returns:
//   - A slice of CardRecord representing the card records fetched from the database.
//   - A pointer to an int representing the total number of records matching the search criteria.
//   - An error if the operation fails, nil otherwise.
func (r *cardQueryRepository) FindAllCards(ctx context.Context, req *requests.FindAllCards) ([]*record.CardRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCardsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	cards, err := r.db.GetCards(ctx, reqDb)

	if err != nil {
		return nil, nil, card_errors.ErrFindAllCardsFailed
	}

	var totalCount int

	if len(cards) > 0 {
		totalCount = int(cards[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToCardsRecord(cards), &totalCount, nil
}

// FindByActive retrieves a paginated list of active card records based on the search criteria
// specified in the request. It queries the database and returns a slice of CardRecord,
// the total count of records, and an error if any occurred.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - req: A FindAllCards request object containing the search parameters
//     such as search keyword, page number, and page size.
//
// Returns:
//   - A slice of CardRecord representing the card records fetched from the database.
//   - A pointer to an int representing the total number of records matching the search criteria.
//   - An error if the operation fails, nil otherwise.
func (r *cardQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllCards) ([]*record.CardRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveCardsWithCountParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveCardsWithCount(ctx, reqDb)

	if err != nil {
		return nil, nil, card_errors.ErrFindActiveCardsFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToCardRecordsActive(res), &totalCount, nil
}

// FindByTrashed retrieves a paginated list of trashed card records based on the search
// criteria specified in the request. It queries the database and returns a slice of
// CardRecord, the total count of records, and an error if any occurred.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - req: A FindAllCards request object containing the search parameters
//     such as search keyword, page number, and page size.
//
// Returns:
//   - A slice of CardRecord representing the card records fetched from the database.
//   - A pointer to an int representing the total number of records matching the search
//     criteria.
//   - An error if the operation fails, nil otherwise.
func (r *cardQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllCards) ([]*record.CardRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedCardsWithCountParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedCardsWithCount(ctx, reqDb)

	if err != nil {
		return nil, nil, card_errors.ErrFindTrashedCardsFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToCardRecordsTrashed(res), &totalCount, nil
}

// FindById retrieves a card record by its ID from the database.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - card_id: The ID of the card to be retrieved.
//
// Returns:
//   - A pointer to a CardRecord representing the card record fetched from the database.
//   - An error if the operation fails, nil otherwise.
func (r *cardQueryRepository) FindById(ctx context.Context, card_id int) (*record.CardRecord, error) {
	res, err := r.db.GetCardByID(ctx, int32(card_id))

	if err != nil {
		return nil, card_errors.ErrFindCardByIdFailed
	}

	return r.mapper.ToCardRecord(res), nil
}

// FindCardByUserId retrieves a card record by its user ID from the database.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - user_id: The ID of the user to be retrieved.
//
// Returns:
//   - A pointer to a CardRecord representing the card record fetched from the database.
//   - An error if the operation fails, nil otherwise.
func (r *cardQueryRepository) FindCardByUserId(ctx context.Context, user_id int) (*record.CardRecord, error) {
	res, err := r.db.GetCardByUserID(ctx, int32(user_id))

	if err != nil {
		return nil, card_errors.ErrFindCardByUserIdFailed
	}

	return r.mapper.ToCardRecord(res), nil
}

// FindCardByCardNumber retrieves a card record by its card number from the database.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - card_number: The card number of the card to be retrieved.
//
// Returns:
//   - A pointer to a CardRecord representing the card record fetched from the database.
//   - An error if the operation fails, nil otherwise.
func (r *cardQueryRepository) FindCardByCardNumber(ctx context.Context, card_number string) (*record.CardRecord, error) {
	res, err := r.db.GetCardByCardNumber(ctx, card_number)

	if err != nil {
		return nil, card_errors.ErrFindCardByCardNumberFailed
	}

	return r.mapper.ToCardRecord(res), nil
}
