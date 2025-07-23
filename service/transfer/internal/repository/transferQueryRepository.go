package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/transfer"
)

// transferQueryRepository is a struct that implements the TransferQueryRepository interface
type transferQueryRepository struct {
	db     *db.Queries
	mapper recordmapper.TransferQueryRecordMapper
}

// NewTransferQueryRepository initializes a new instance of transferQueryRepository with the provided
// database queries, context, and transfer record mapper. This repository is responsible for executing
// query operations related to transfer records in the database.
//
// Parameters:
//   - db: A pointer to the db.Queries object for executing database queries.
//   - mapper: A TransferRecordMapping that provides methods to map database rows to TransferRecord domain models.
//
// Returns:
//   - A pointer to the newly created transferQueryRepository instance.
func NewTransferQueryRepository(db *db.Queries, mapper recordmapper.TransferQueryRecordMapper) TransferQueryRepository {
	return &transferQueryRepository{
		db:     db,
		mapper: mapper,
	}
}

// FindAll retrieves all transfer records with optional filtering and pagination.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request parameters for filtering and pagination.
//
// Returns:
//   - []*record.TransferRecord: List of transfer records.
//   - *int: Total number of records (for pagination).
//   - error: Any error encountered during the operation.
func (r *transferQueryRepository) FindAll(ctx context.Context, req *requests.FindAllTransfers) ([]*record.TransferRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransfersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTransfers(ctx, reqDb)

	if err != nil {
		return nil, nil, transfer_errors.ErrFindAllTransfersFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	so := r.mapper.ToTransfersRecordAll(res)

	return so, &totalCount, nil
}

// FindByActive retrieves all active (non-trashed) transfer records.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request parameters for filtering and pagination.
//
// Returns:
//   - []*record.TransferRecord: List of active transfer records.
//   - *int: Total number of records (for pagination).
//   - error: Any error encountered during the operation.
func (r *transferQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllTransfers) ([]*record.TransferRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveTransfersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveTransfers(ctx, reqDb)

	if err != nil {
		return nil, nil, transfer_errors.ErrFindActiveTransfersFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	so := r.mapper.ToTransfersRecordActive(res)

	return so, &totalCount, nil
}

// FindByTrashed retrieves all soft-deleted (trashed) transfer records.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request parameters for filtering and pagination.
//
// Returns:
//   - []*record.TransferRecord: List of trashed transfer records.
//   - *int: Total number of records (for pagination).
//   - error: Any error encountered during the operation.
func (r *transferQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllTransfers) ([]*record.TransferRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedTransfersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedTransfers(ctx, reqDb)

	if err != nil {
		return nil, nil, transfer_errors.ErrFindTrashedTransfersFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	so := r.mapper.ToTransfersRecordTrashed(res)

	return so, &totalCount, nil
}

// FindById retrieves a single transfer record by its ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - id: The ID of the transfer record.
//
// Returns:
//   - *record.TransferRecord: The transfer record, if found.
//   - error: Any error encountered during the operation.
func (r *transferQueryRepository) FindById(ctx context.Context, id int) (*record.TransferRecord, error) {
	transfer, err := r.db.GetTransferByID(ctx, int32(id))

	if err != nil {
		return nil, transfer_errors.ErrFindTransferByIdFailed
	}

	so := r.mapper.ToTransferRecord(transfer)

	return so, nil
}

// FindTransferByTransferFrom retrieves all transfer records where the given card is the sender.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - transfer_from: The sender card number.
//
// Returns:
//   - []*record.TransferRecord: List of transfer records from the specified sender.
//   - error: Any error encountered during the operation.
func (r *transferQueryRepository) FindTransferByTransferFrom(ctx context.Context, transfer_from string) ([]*record.TransferRecord, error) {
	res, err := r.db.GetTransfersBySourceCard(ctx, transfer_from)

	if err != nil {
		return nil, transfer_errors.ErrFindTransferByTransferFromFailed
	}

	so := r.mapper.ToTransfersRecord(res)

	return so, nil
}

// FindTransferByTransferTo retrieves all transfer records where the given card is the receiver.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - transfer_to: The receiver card number.
//
// Returns:
//   - []*record.TransferRecord: List of transfer records to the specified receiver.
//   - error: Any error encountered during the operation.
func (r *transferQueryRepository) FindTransferByTransferTo(ctx context.Context, transfer_to string) ([]*record.TransferRecord, error) {
	res, err := r.db.GetTransfersByDestinationCard(ctx, transfer_to)

	if err != nil {
		return nil, transfer_errors.ErrFindTransferByTransferToFailed
	}

	so := r.mapper.ToTransfersRecord(res)

	return so, nil
}
