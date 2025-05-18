package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type merchantQueryRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.MerchantRecordMapping
}

func NewMerchantQueryRepository(db *db.Queries, ctx context.Context, mapping recordmapper.MerchantRecordMapping) *merchantQueryRepository {
	return &merchantQueryRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *merchantQueryRepository) FindAllMerchants(req *requests.FindAllMerchants) ([]*record.MerchantRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	merchant, err := r.db.GetMerchants(r.ctx, reqDb)

	if err != nil {
		return nil, nil, merchant_errors.ErrFindAllMerchantsFailed
	}

	var totalCount int
	if len(merchant) > 0 {
		totalCount = int(merchant[0].TotalCount)
	} else {
		totalCount = 0
	}
	return r.mapping.ToMerchantsGetAllRecord(merchant), &totalCount, nil
}

func (r *merchantQueryRepository) FindByActive(req *requests.FindAllMerchants) ([]*record.MerchantRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveMerchantsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveMerchants(r.ctx, reqDb)

	if err != nil {
		return nil, nil, merchant_errors.ErrFindActiveMerchantsFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantsActiveRecord(res), &totalCount, nil
}

func (r *merchantQueryRepository) FindByTrashed(req *requests.FindAllMerchants) ([]*record.MerchantRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedMerchantsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedMerchants(r.ctx, reqDb)

	if err != nil {
		return nil, nil, merchant_errors.ErrFindTrashedMerchantsFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantsTrashedRecord(res), &totalCount, nil
}

func (r *merchantQueryRepository) FindById(merchant_id int) (*record.MerchantRecord, error) {
	res, err := r.db.GetMerchantByID(r.ctx, int32(merchant_id))

	if err != nil {
		return nil, merchant_errors.ErrFindMerchantByIdFailed
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantQueryRepository) FindByApiKey(api_key string) (*record.MerchantRecord, error) {
	res, err := r.db.GetMerchantByApiKey(r.ctx, api_key)

	if err != nil {
		return nil, merchant_errors.ErrFindMerchantByApiKeyFailed
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantQueryRepository) FindByName(name string) (*record.MerchantRecord, error) {
	res, err := r.db.GetMerchantByName(r.ctx, name)

	if err != nil {
		return nil, merchant_errors.ErrFindMerchantByNameFailed
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantQueryRepository) FindByMerchantUserId(user_id int) ([]*record.MerchantRecord, error) {
	res, err := r.db.GetMerchantsByUserID(r.ctx, int32(user_id))

	if err != nil {
		return nil, merchant_errors.ErrFindMerchantByUserIdFailed
	}

	return r.mapping.ToMerchantsRecord(res), nil
}
