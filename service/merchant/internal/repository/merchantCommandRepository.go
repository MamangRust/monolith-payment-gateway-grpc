package repository

import (
	"context"

	apikey "github.com/MamangRust/monolith-payment-gateway-pkg/api-key"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type merchantCommandRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.MerchantRecordMapping
}

func NewMerchantCommandRepository(db *db.Queries, ctx context.Context, mapping recordmapper.MerchantRecordMapping) *merchantCommandRepository {
	return &merchantCommandRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *merchantCommandRepository) CreateMerchant(request *requests.CreateMerchantRequest) (*record.MerchantRecord, error) {
	req := db.CreateMerchantParams{
		Name:   request.Name,
		ApiKey: apikey.GenerateApiKey(),
		UserID: int32(request.UserID),
		Status: "inactive",
	}

	res, err := r.db.CreateMerchant(r.ctx, req)

	if err != nil {
		return nil, merchant_errors.ErrCreateMerchantFailed
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantCommandRepository) UpdateMerchant(request *requests.UpdateMerchantRequest) (*record.MerchantRecord, error) {
	req := db.UpdateMerchantParams{
		MerchantID: int32(*request.MerchantID),
		Name:       request.Name,
		UserID:     int32(request.UserID),
		Status:     request.Status,
	}

	res, err := r.db.UpdateMerchant(r.ctx, req)

	if err != nil {
		return nil, merchant_errors.ErrUpdateMerchantFailed
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantCommandRepository) UpdateMerchantStatus(request *requests.UpdateMerchantStatusRequest) (*record.MerchantRecord, error) {
	req := db.UpdateMerchantStatusParams{
		MerchantID: int32(*request.MerchantID),
		Status:     request.Status,
	}

	res, err := r.db.UpdateMerchantStatus(r.ctx, req)

	if err != nil {
		return nil, merchant_errors.ErrUpdateMerchantStatusFailed
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantCommandRepository) TrashedMerchant(merchant_id int) (*record.MerchantRecord, error) {
	res, err := r.db.TrashMerchant(r.ctx, int32(merchant_id))

	if err != nil {
		return nil, merchant_errors.ErrTrashedMerchantFailed
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantCommandRepository) RestoreMerchant(merchant_id int) (*record.MerchantRecord, error) {
	res, err := r.db.RestoreMerchant(r.ctx, int32(merchant_id))

	if err != nil {
		return nil, merchant_errors.ErrRestoreMerchantFailed
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantCommandRepository) DeleteMerchantPermanent(merchant_id int) (bool, error) {
	err := r.db.DeleteMerchantPermanently(r.ctx, int32(merchant_id))

	if err != nil {
		return false, merchant_errors.ErrDeleteMerchantPermanentFailed
	}

	return true, nil
}

func (r *merchantCommandRepository) RestoreAllMerchant() (bool, error) {
	err := r.db.RestoreAllMerchants(r.ctx)

	if err != nil {
		return false, merchant_errors.ErrRestoreAllMerchantFailed
	}

	return true, nil
}

func (r *merchantCommandRepository) DeleteAllMerchantPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentMerchants(r.ctx)

	if err != nil {
		return false, merchant_errors.ErrDeleteAllMerchantPermanentFailed
	}

	return true, nil
}
