package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	withdraw_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/withdraw"
)

// withdrawQueryRepository is a struct that implements the WithdrawQueryRepository interface
type withdrawQueryRepository struct {
	db     *db.Queries
	mapper recordmapper.WithdrawQueryRecordMapping
}

// NewWithdrawQueryRepository initializes a new instance of withdrawQueryRepository with the provided
// database queries, context, and withdraw record mapper. This repository is responsible for executing
// query operations related to withdraw records in the database.
//
// Parameters:
//   - db: A pointer to the db.Queries object for executing database queries.
//   - mapper: A WithdrawRecordMapping that provides methods to map database rows to WithdrawRecord domain models.
//
// Returns:
//   - A pointer to the newly created withdrawQueryRepository instance.
func NewWithdrawQueryRepository(db *db.Queries, mapper recordmapper.WithdrawQueryRecordMapping) WithdrawQueryRepository {
	return &withdrawQueryRepository{
		db:     db,
		mapper: mapper,
	}
}

// FindAll retrieves all withdraw records with pagination and filtering.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing filter and pagination information.
//
// Returns:
//   - []*record.WithdrawRecord: List of withdraw records.
//   - *int: Total count.
//   - error: An error if the operation fails.
func (r *withdrawQueryRepository) FindAll(ctx context.Context, req *requests.FindAllWithdraws) ([]*record.WithdrawRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetWithdrawsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	withdraw, err := r.db.GetWithdraws(ctx, reqDb)

	if err != nil {
		return nil, nil, withdraw_errors.ErrFindAllWithdrawsFailed
	}

	var totalCount int
	if len(withdraw) > 0 {
		totalCount = int(withdraw[0].TotalCount)
	} else {
		totalCount = 0
	}

	so := r.mapper.ToWithdrawsRecordAll(withdraw)

	return so, &totalCount, nil
}

// FindByActive retrieves active (non-deleted) withdraw records.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing filter and pagination information.
//
// Returns:
//   - []*record.WithdrawRecord: List of active withdraw records.
//   - *int: Total count.
//   - error: An error if the operation fails.
func (r *withdrawQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllWithdraws) ([]*record.WithdrawRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveWithdrawsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveWithdraws(ctx, reqDb)

	if err != nil {
		return nil, nil, withdraw_errors.ErrFindActiveWithdrawsFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	so := r.mapper.ToWithdrawsRecordActive(res)

	return so, &totalCount, nil
}

// FindByTrashed retrieves soft-deleted withdraw records.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing filter and pagination information.
//
// Returns:
//   - []*record.WithdrawRecord: List of trashed withdraw records.
//   - *int: Total count.
//   - error: An error if the operation fails.
func (r *withdrawQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllWithdraws) ([]*record.WithdrawRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedWithdrawsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedWithdraws(ctx, reqDb)

	if err != nil {
		return nil, nil, withdraw_errors.ErrFindTrashedWithdrawsFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	so := r.mapper.ToWithdrawsRecordTrashed(res)

	return so, &totalCount, nil
}

// FindAllByCardNumber retrieves all withdraw records associated with a specific card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing card number and filter info.
//
// Returns:
//   - []*record.WithdrawRecord: List of withdraw records for the card.
//   - *int: Total count.
//   - error: An error if the operation fails.
func (r *withdrawQueryRepository) FindAllByCardNumber(ctx context.Context, req *requests.FindAllWithdrawCardNumber) ([]*record.WithdrawRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetWithdrawsByCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    req.Search,
		Limit:      int32(req.PageSize),
		Offset:     int32(offset),
	}

	withdraw, err := r.db.GetWithdrawsByCardNumber(ctx, reqDb)

	if err != nil {
		return nil, nil, withdraw_errors.ErrFindWithdrawsByCardNumberFailed
	}
	var totalCount int
	if len(withdraw) > 0 {
		totalCount = int(withdraw[0].TotalCount)
	} else {
		totalCount = 0
	}

	so := r.mapper.ToWithdrawsByCardNumberRecord(withdraw)

	return so, &totalCount, nil
}

// FindById retrieves a withdraw record by its unique ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - id: The ID of the withdraw record.
//
// Returns:
//   - *record.WithdrawRecord: The withdraw record if found.
//   - error: An error if the operation fails or the record is not found.
func (r *withdrawQueryRepository) FindById(ctx context.Context, id int) (*record.WithdrawRecord, error) {
	withdraw, err := r.db.GetWithdrawByID(ctx, int32(id))

	if err != nil {
		return nil, withdraw_errors.ErrFindWithdrawByIdFailed
	}

	so := r.mapper.ToWithdrawRecord(withdraw)

	return so, nil
}
