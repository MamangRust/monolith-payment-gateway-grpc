package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type withdrawQueryRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.WithdrawRecordMapping
}

func NewWithdrawQueryRepository(db *db.Queries, ctx context.Context, mapping recordmapper.WithdrawRecordMapping) *withdrawQueryRepository {
	return &withdrawQueryRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *withdrawQueryRepository) FindAll(req *requests.FindAllWithdraws) ([]*record.WithdrawRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetWithdrawsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	withdraw, err := r.db.GetWithdraws(r.ctx, reqDb)

	if err != nil {
		return nil, nil, withdraw_errors.ErrFindAllWithdrawsFailed
	}

	var totalCount int
	if len(withdraw) > 0 {
		totalCount = int(withdraw[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToWithdrawsRecordALl(withdraw), &totalCount, nil

}

func (r *withdrawQueryRepository) FindByActive(req *requests.FindAllWithdraws) ([]*record.WithdrawRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveWithdrawsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveWithdraws(r.ctx, reqDb)

	if err != nil {
		return nil, nil, withdraw_errors.ErrFindActiveWithdrawsFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToWithdrawsRecordActive(res), &totalCount, nil
}

func (r *withdrawQueryRepository) FindByTrashed(req *requests.FindAllWithdraws) ([]*record.WithdrawRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedWithdrawsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedWithdraws(r.ctx, reqDb)

	if err != nil {
		return nil, nil, withdraw_errors.ErrFindTrashedWithdrawsFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToWithdrawsRecordTrashed(res), &totalCount, nil
}

func (r *withdrawQueryRepository) FindAllByCardNumber(req *requests.FindAllWithdrawCardNumber) ([]*record.WithdrawRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetWithdrawsByCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    req.Search,
		Limit:      int32(req.PageSize),
		Offset:     int32(offset),
	}

	withdraw, err := r.db.GetWithdrawsByCardNumber(r.ctx, reqDb)

	if err != nil {
		return nil, nil, withdraw_errors.ErrFindWithdrawsByCardNumberFailed
	}
	var totalCount int
	if len(withdraw) > 0 {
		totalCount = int(withdraw[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToWithdrawsByCardNumberRecord(withdraw), &totalCount, nil

}

func (r *withdrawQueryRepository) FindById(id int) (*record.WithdrawRecord, error) {
	withdraw, err := r.db.GetWithdrawByID(r.ctx, int32(id))

	if err != nil {
		return nil, withdraw_errors.ErrFindWithdrawByIdFailed
	}

	return r.mapping.ToWithdrawRecord(withdraw), nil
}
