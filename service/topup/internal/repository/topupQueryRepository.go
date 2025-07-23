package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/topup"
)

// topupQueryRepository is a struct that implements the TopupQueryRepository interface
type topupQueryRepository struct {
	db     *db.Queries
	mapper recordmapper.TopupQueryRecordMapping
}

// NewTopupQueryRepository initializes a new instance of topupQueryRepository.
//
// Parameters:
//   - db: A pointer to the db.Queries object for executing database queries.
//   - mapper: A TopupRecordMapping that provides methods to map database rows to Topup domain models.
//
// Returns:
//   - A pointer to the newly created topupQueryRepository instance.
func NewTopupQueryRepository(db *db.Queries, mapper recordmapper.TopupQueryRecordMapping) TopupQueryRepository {
	return &topupQueryRepository{
		db:     db,
		mapper: mapper,
	}
}

// FindAllTopups retrieves a paginated list of all topups.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request parameters for pagination and filters.
//
// Returns:
//   - []*record.TopupRecord: List of topup records.
//   - *int: Total number of records.
//   - error: Error if the query fails.
func (r *topupQueryRepository) FindAllTopups(ctx context.Context, req *requests.FindAllTopups) ([]*record.TopupRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTopupsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTopups(ctx, reqDb)

	if err != nil {
		return nil, nil, topup_errors.ErrFindAllTopupsFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToTopupRecordsAll(res), &totalCount, nil
}

// FindByActive retrieves a paginated list of active topup records.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request parameters for pagination and filters.
//
// Returns:
//   - []*record.TopupRecord: List of active topup records.
//   - *int: Total number of records.
//   - error: Error if the query fails.
func (r *topupQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllTopups) ([]*record.TopupRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveTopupsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveTopups(ctx, reqDb)

	if err != nil {
		return nil, nil, topup_errors.ErrFindTopupsByActiveFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToTopupRecordsActive(res), &totalCount, nil
}

// FindByTrashed retrieves a paginated list of trashed topup records.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request parameters for pagination and filters.
//
// Returns:
//   - []*record.TopupRecord: List of trashed topup records.
//   - *int: Total number of records.
//   - error: Error if the query fails.
func (r *topupQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllTopups) ([]*record.TopupRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedTopupsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedTopups(ctx, reqDb)

	if err != nil {
		return nil, nil, topup_errors.ErrFindTopupsByTrashedFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToTopupRecordsTrashed(res), &totalCount, nil
}

// FindAllTopupByCardNumber retrieves all topups associated with a specific card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing the card number and pagination info.
//
// Returns:
//   - []*record.TopupRecord: List of topups associated with the card.
//   - *int: Total number of records.
//   - error: Error if the query fails.
func (r *topupQueryRepository) FindAllTopupByCardNumber(ctx context.Context, req *requests.FindAllTopupsByCardNumber) ([]*record.TopupRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTopupsByCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    req.Search,
		Limit:      int32(req.PageSize),
		Offset:     int32(offset),
	}

	res, err := r.db.GetTopupsByCardNumber(ctx, reqDb)

	if err != nil {
		return nil, nil, topup_errors.ErrFindTopupsByCardNumberFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToTopupByCardNumberRecords(res), &totalCount, nil
}

// FindById retrieves a topup record by its ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - topup_id: The unique ID of the topup.
//
// Returns:
//   - *record.TopupRecord: The found topup record.
//   - error: Error if the query fails or topup is not found.
func (r *topupQueryRepository) FindById(ctx context.Context, topup_id int) (*record.TopupRecord, error) {
	res, err := r.db.GetTopupByID(ctx, int32(topup_id))
	if err != nil {
		return nil, topup_errors.ErrFindTopupByIdFailed
	}
	return r.mapper.ToTopupRecord(res), nil
}
