package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type transferQueryRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.TransferRecordMapping
}

func NewTransferQueryRepository(db *db.Queries, ctx context.Context, mapping recordmapper.TransferRecordMapping) *transferQueryRepository {
	return &transferQueryRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *transferQueryRepository) FindAll(req *requests.FindAllTranfers) ([]*record.TransferRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransfersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTransfers(r.ctx, reqDb)

	if err != nil {
		return nil, nil, transfer_errors.ErrFindAllTransfersFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransfersRecordAll(res), &totalCount, nil
}

func (r *transferQueryRepository) FindByActive(req *requests.FindAllTranfers) ([]*record.TransferRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveTransfersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveTransfers(r.ctx, reqDb)

	if err != nil {
		return nil, nil, transfer_errors.ErrFindActiveTransfersFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransfersRecordActive(res), &totalCount, nil
}

func (r *transferQueryRepository) FindByTrashed(req *requests.FindAllTranfers) ([]*record.TransferRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedTransfersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedTransfers(r.ctx, reqDb)

	if err != nil {
		return nil, nil, transfer_errors.ErrFindTrashedTransfersFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransfersRecordTrashed(res), &totalCount, nil
}

func (r *transferQueryRepository) FindById(id int) (*record.TransferRecord, error) {
	transfer, err := r.db.GetTransferByID(r.ctx, int32(id))

	if err != nil {
		return nil, transfer_errors.ErrFindTransferByIdFailed
	}

	return r.mapping.ToTransferRecord(transfer), nil
}

func (r *transferQueryRepository) FindTransferByTransferFrom(transfer_from string) ([]*record.TransferRecord, error) {
	res, err := r.db.GetTransfersBySourceCard(r.ctx, transfer_from)

	if err != nil {
		return nil, transfer_errors.ErrFindTransferByTransferFromFailed
	}

	return r.mapping.ToTransfersRecord(res), nil
}

func (r *transferQueryRepository) FindTransferByTransferTo(transfer_to string) ([]*record.TransferRecord, error) {
	res, err := r.db.GetTransfersByDestinationCard(r.ctx, transfer_to)

	if err != nil {
		return nil, transfer_errors.ErrFindTransferByTransferToFailed
	}
	return r.mapping.ToTransfersRecord(res), nil
}
