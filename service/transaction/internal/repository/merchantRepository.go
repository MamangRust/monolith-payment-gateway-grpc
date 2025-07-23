package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/merchant"
)

// merchantRepository is a struct that implements the MerchantRepository interface.
type merchantRepository struct {
	db     *db.Queries
	mapper recordmapper.MerchantQueryRecordMapper
}

// NewMerchantRepository creates a new instance of merchantRepository with the provided database queries, context, and merchant record mapper.
// It is responsible for providing methods to query merchant records from the database.
//
// Parameters:
//   - db: A pointer to the db.Queries object for executing database queries.
//   - mapper: A MerchantRecordMapping that provides methods to map database rows to Merchant domain models.
//
// Returns:
//   - A pointer to the newly created merchantRepository instance.
func NewMerchantRepository(db *db.Queries, mapper recordmapper.MerchantQueryRecordMapper) MerchantRepository {
	return &merchantRepository{
		db:     db,
		mapper: mapper,
	}
}

// FindByApiKey retrieves a merchant by its API key.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - api_key: The API key associated with the merchant.
//
// Returns:
//   - *record.MerchantRecord: The merchant record if found.
//   - error: Error if something went wrong during the query.
func (r *merchantRepository) FindByApiKey(ctx context.Context, api_key string) (*record.MerchantRecord, error) {
	res, err := r.db.GetMerchantByApiKey(ctx, api_key)

	if err != nil {
		return nil, merchant_errors.ErrFindMerchantByApiKeyFailed
	}

	return r.mapper.ToMerchantRecord(res), nil
}
