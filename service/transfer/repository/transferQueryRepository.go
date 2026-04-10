package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/repository"
)

type transferQueryRepository struct {
	db *db.Queries
}

func NewTransferQueryRepository(db *db.Queries) TransferQueryRepository {
	return &transferQueryRepository{
		db: db,
	}
}

func (r *transferQueryRepository) FindAll(ctx context.Context, req *requests.FindAllTransfers) ([]*db.GetTransfersRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransfersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTransfers(ctx, reqDb)

	if err != nil {
		return nil, transfer_errors.ErrFindAllTransfersFailed
	}

	return res, nil
}

func (r *transferQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllTransfers) ([]*db.GetActiveTransfersRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveTransfersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveTransfers(ctx, reqDb)

	if err != nil {
		return nil, transfer_errors.ErrFindActiveTransfersFailed
	}

	return res, nil
}

func (r *transferQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllTransfers) ([]*db.GetTrashedTransfersRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedTransfersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedTransfers(ctx, reqDb)

	if err != nil {
		return nil, transfer_errors.ErrFindTrashedTransfersFailed
	}

	return res, nil
}

func (r *transferQueryRepository) FindById(ctx context.Context, id int) (*db.GetTransferByIDRow, error) {
	transfer, err := r.db.GetTransferByID(ctx, int32(id))

	if err != nil {
		return nil, transfer_errors.ErrFindTransferByIdFailed
	}

	return transfer, nil
}

func (r *transferQueryRepository) FindTransferByTransferFrom(ctx context.Context, transfer_from string) ([]*db.GetTransfersBySourceCardRow, error) {
	res, err := r.db.GetTransfersBySourceCard(ctx, transfer_from)

	if err != nil {
		return nil, transfer_errors.ErrFindTransferByTransferFromFailed
	}

	return res, nil
}

func (r *transferQueryRepository) FindTransferByTransferTo(ctx context.Context, transfer_to string) ([]*db.GetTransfersByDestinationCardRow, error) {
	res, err := r.db.GetTransfersByDestinationCard(ctx, transfer_to)

	if err != nil {
		return nil, transfer_errors.ErrFindTransferByTransferToFailed
	}
	return res, nil
}
