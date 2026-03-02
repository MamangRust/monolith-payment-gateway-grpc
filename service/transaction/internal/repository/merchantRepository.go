package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/repository"
)

type merchantRepository struct {
	db *db.Queries
}

func NewMerchantRepository(db *db.Queries) MerchantRepository {
	return &merchantRepository{
		db: db,
	}
}

func (r *merchantRepository) FindByApiKey(ctx context.Context, api_key string) (*db.GetMerchantByApiKeyRow, error) {
	res, err := r.db.GetMerchantByApiKey(ctx, api_key)

	if err != nil {
		return nil, merchant_errors.ErrFindMerchantByApiKeyFailed
	}

	return res, nil
}
