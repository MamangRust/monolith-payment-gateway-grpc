package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type topupQueryRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.TopupRecordMapping
}

func NewTopupQueryRepository(db *db.Queries, ctx context.Context, mapping recordmapper.TopupRecordMapping) *topupQueryRepository {
	return &topupQueryRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *topupQueryRepository) FindAllTopups(req *requests.FindAllTopups) ([]*record.TopupRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTopupsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTopups(r.ctx, reqDb)

	if err != nil {
		return nil, nil, topup_errors.ErrFindAllTopupsFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTopupRecordsAll(res), &totalCount, nil
}

func (r *topupQueryRepository) FindByActive(req *requests.FindAllTopups) ([]*record.TopupRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveTopupsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveTopups(r.ctx, reqDb)

	if err != nil {
		return nil, nil, topup_errors.ErrFindTopupsByActiveFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTopupRecordsActive(res), &totalCount, nil
}

func (r *topupQueryRepository) FindByTrashed(req *requests.FindAllTopups) ([]*record.TopupRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedTopupsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedTopups(r.ctx, reqDb)

	if err != nil {
		return nil, nil, topup_errors.ErrFindTopupsByTrashedFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTopupRecordsTrashed(res), &totalCount, nil
}

func (r *topupQueryRepository) FindAllTopupByCardNumber(req *requests.FindAllTopupsByCardNumber) ([]*record.TopupRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTopupsByCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    req.Search,
		Limit:      int32(req.PageSize),
		Offset:     int32(offset),
	}

	res, err := r.db.GetTopupsByCardNumber(r.ctx, reqDb)

	if err != nil {
		return nil, nil, topup_errors.ErrFindTopupsByCardNumberFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTopupByCardNumberRecords(res), &totalCount, nil
}

func (r *topupQueryRepository) FindById(topup_id int) (*record.TopupRecord, error) {
	res, err := r.db.GetTopupByID(r.ctx, int32(topup_id))
	if err != nil {
		return nil, topup_errors.ErrFindTopupByIdFailed
	}
	return r.mapping.ToTopupRecord(res), nil
}
