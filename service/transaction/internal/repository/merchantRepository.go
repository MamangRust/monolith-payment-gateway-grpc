package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type merchantRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.MerchantRecordMapping
}

func NewMerchantRepository(db *db.Queries, ctx context.Context, mapping recordmapper.MerchantRecordMapping) *merchantRepository {
	return &merchantRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *merchantRepository) FindByApiKey(api_key string) (*record.MerchantRecord, error) {
	res, err := r.db.GetMerchantByApiKey(r.ctx, api_key)

	if err != nil {
		return nil, merchant_errors.ErrFindMerchantByApiKeyFailed
	}

	return r.mapping.ToMerchantRecord(res), nil
}
