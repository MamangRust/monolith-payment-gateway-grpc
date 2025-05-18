package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type topupCommandRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.TopupRecordMapping
}

func NewTopupCommandRepository(db *db.Queries, ctx context.Context, mapping recordmapper.TopupRecordMapping) *topupCommandRepository {
	return &topupCommandRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *topupCommandRepository) CreateTopup(request *requests.CreateTopupRequest) (*record.TopupRecord, error) {
	req := db.CreateTopupParams{
		CardNumber:  request.CardNumber,
		TopupAmount: int32(request.TopupAmount),
		TopupMethod: request.TopupMethod,
	}

	res, err := r.db.CreateTopup(r.ctx, req)

	if err != nil {
		return nil, topup_errors.ErrCreateTopupFailed
	}

	return r.mapping.ToTopupRecord(res), nil
}

func (r *topupCommandRepository) UpdateTopup(request *requests.UpdateTopupRequest) (*record.TopupRecord, error) {
	req := db.UpdateTopupParams{
		TopupID:     int32(*request.TopupID),
		CardNumber:  request.CardNumber,
		TopupAmount: int32(request.TopupAmount),
		TopupMethod: request.TopupMethod,
	}

	res, err := r.db.UpdateTopup(r.ctx, req)

	if err != nil {
		return nil, topup_errors.ErrUpdateTopupFailed
	}

	return r.mapping.ToTopupRecord(res), nil
}

func (r *topupCommandRepository) UpdateTopupAmount(request *requests.UpdateTopupAmount) (*record.TopupRecord, error) {
	req := db.UpdateTopupAmountParams{
		TopupID:     int32(request.TopupID),
		TopupAmount: int32(request.TopupAmount),
	}

	res, err := r.db.UpdateTopupAmount(r.ctx, req)

	if err != nil {
		return nil, topup_errors.ErrUpdateTopupAmountFailed
	}

	return r.mapping.ToTopupRecord(res), nil
}

func (r *topupCommandRepository) UpdateTopupStatus(request *requests.UpdateTopupStatus) (*record.TopupRecord, error) {
	req := db.UpdateTopupStatusParams{
		TopupID: int32(request.TopupID),
		Status:  request.Status,
	}

	res, err := r.db.UpdateTopupStatus(r.ctx, req)

	if err != nil {
		return nil, topup_errors.ErrUpdateTopupStatusFailed
	}

	return r.mapping.ToTopupRecord(res), nil
}

func (r *topupCommandRepository) TrashedTopup(topup_id int) (*record.TopupRecord, error) {
	res, err := r.db.TrashTopup(r.ctx, int32(topup_id))
	if err != nil {
		return nil, topup_errors.ErrTrashedTopupFailed
	}
	return r.mapping.ToTopupRecord(res), nil
}

func (r *topupCommandRepository) RestoreTopup(topup_id int) (*record.TopupRecord, error) {
	res, err := r.db.RestoreTopup(r.ctx, int32(topup_id))
	if err != nil {
		return nil, topup_errors.ErrRestoreTopupFailed
	}
	return r.mapping.ToTopupRecord(res), nil
}

func (r *topupCommandRepository) DeleteTopupPermanent(topup_id int) (bool, error) {
	err := r.db.DeleteTopupPermanently(r.ctx, int32(topup_id))
	if err != nil {
		return false, topup_errors.ErrDeleteTopupPermanentFailed
	}
	return true, nil
}

func (r *topupCommandRepository) RestoreAllTopup() (bool, error) {
	err := r.db.RestoreAllTopups(r.ctx)

	if err != nil {
		return false, topup_errors.ErrRestoreAllTopupFailed
	}

	return true, nil
}

func (r *topupCommandRepository) DeleteAllTopupPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentTopups(r.ctx)

	if err != nil {
		return false, topup_errors.ErrDeleteAllTopupPermanentFailed
	}

	return true, nil
}
