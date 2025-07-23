package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/merchant"
)

// merchantQueryRepository is a struct that implements the MerchantQueryRepository interface
type merchantQueryRepository struct {
	db     *db.Queries
	mapper recordmapper.MerchantQueryRecordMapper
}

// NewMerchantQueryRepository creates a new instance of merchantQueryRepository with the provided
// database queries, context, and merchant record mapper. This repository is responsible for executing
// query operations related to merchant records in the database.
//
// Parameters:
//   - db: A pointer to the db.Queries object for executing database queries.
//   - mapper: A MerchantRecordMapping that provides methods to map database rows to Merchant domain models.
//
// Returns:
//   - A pointer to the newly created merchantQueryRepository instance.
func NewMerchantQueryRepository(db *db.Queries, mapper recordmapper.MerchantQueryRecordMapper) MerchantQueryRepository {
	return &merchantQueryRepository{
		db:     db,
		mapper: mapper,
	}
}

// FindAllMerchants retrieves a list of merchants based on the provided request.
//
// This function implements the MerchantQueryRepository interface and is responsible for executing the
// database query to retrieve a list of merchant records.
//
// Parameters:
//   - req: A pointer to a FindAllMerchants request object containing the page, page size, and search string.
//
// Returns:
//   - A slice of pointers to MerchantRecord objects, containing the merchant data retrieved from the database.
//   - A pointer to an integer, containing the total records count that matches the search query.
//   - An error, if the query execution fails.
func (r *merchantQueryRepository) FindAllMerchants(ctx context.Context, req *requests.FindAllMerchants) ([]*record.MerchantRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	merchant, err := r.db.GetMerchants(ctx, reqDb)

	if err != nil {
		return nil, nil, merchant_errors.ErrFindAllMerchantsFailed
	}

	var totalCount int
	if len(merchant) > 0 {
		totalCount = int(merchant[0].TotalCount)
	} else {
		totalCount = 0
	}
	return r.mapper.ToMerchantsGetAllRecord(merchant), &totalCount, nil
}

// FindByActive retrieves a list of active merchants based on the provided request.
//
// This function implements the MerchantQueryRepository interface and is responsible for executing the
// database query to retrieve a list of active merchant records.
//
// Parameters:
//   - req: A pointer to a FindAllMerchants request object containing the page, page size, and search string.
//
// Returns:
//   - A slice of pointers to MerchantRecord objects, containing the active merchant data retrieved from the database.
//   - A pointer to an integer, containing the total records count that matches the search query.
//   - An error, if the query execution fails.
func (r *merchantQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllMerchants) ([]*record.MerchantRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveMerchantsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveMerchants(ctx, reqDb)

	if err != nil {
		return nil, nil, merchant_errors.ErrFindActiveMerchantsFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToMerchantsActiveRecord(res), &totalCount, nil
}

// FindByTrashed retrieves a list of trashed merchants based on the provided request.
//
// This function implements the MerchantQueryRepository interface and is responsible for executing the
// database query to retrieve a list of trashed merchant records.
//
// Parameters:
//   - req: A pointer to a FindAllMerchants request object containing the page, page size, and search string.
//
// Returns:
//   - A slice of pointers to MerchantRecord objects, containing the trashed merchant data retrieved from the database.
//   - A pointer to an integer, containing the total records count that matches the search query.
//   - An error, if the query execution fails.
func (r *merchantQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllMerchants) ([]*record.MerchantRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedMerchantsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedMerchants(ctx, reqDb)

	if err != nil {
		return nil, nil, merchant_errors.ErrFindTrashedMerchantsFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToMerchantsTrashedRecord(res), &totalCount, nil
}

// FindById retrieves a single merchant record based on the provided merchant ID.
//
// This function implements the MerchantQueryRepository interface and is responsible for executing the
// database query to retrieve a single merchant record.
//
// Parameters:
//   - merchant_id: The merchant ID to retrieve
//
// Returns:
//   - A pointer to a MerchantRecord object, containing the merchant data retrieved from the database.
//   - An error, if the query execution fails.
func (r *merchantQueryRepository) FindById(ctx context.Context, merchant_id int) (*record.MerchantRecord, error) {
	res, err := r.db.GetMerchantByID(ctx, int32(merchant_id))

	if err != nil {
		return nil, merchant_errors.ErrFindMerchantByIdFailed
	}

	return r.mapper.ToMerchantRecord(res), nil
}

// FindByApiKey retrieves a merchant record based on the provided API key.
//
// This function implements the MerchantQueryRepository interface and is responsible for executing the
// database query to retrieve a single merchant record associated with the given API key.
//
// Parameters:
//   - api_key: A string containing the API key used to identify the merchant.
//
// Returns:
//   - A pointer to a MerchantRecord object containing the merchant data retrieved from the database.
//   - An error, if the query execution fails.
func (r *merchantQueryRepository) FindByApiKey(ctx context.Context, api_key string) (*record.MerchantRecord, error) {
	res, err := r.db.GetMerchantByApiKey(ctx, api_key)

	if err != nil {
		return nil, merchant_errors.ErrFindMerchantByApiKeyFailed
	}

	return r.mapper.ToMerchantRecord(res), nil
}

// FindByName retrieves a merchant record based on the provided name.
//
// This function implements the MerchantQueryRepository interface and is responsible for executing the
// database query to retrieve a single merchant record associated with the given name.
//
// Parameters:
//   - name: A string containing the name used to identify the merchant.
//
// Returns:
//   - A pointer to a MerchantRecord object containing the merchant data retrieved from the database.
//   - An error, if the query execution fails.
func (r *merchantQueryRepository) FindByName(ctx context.Context, name string) (*record.MerchantRecord, error) {
	res, err := r.db.GetMerchantByName(ctx, name)

	if err != nil {
		return nil, merchant_errors.ErrFindMerchantByNameFailed
	}

	return r.mapper.ToMerchantRecord(res), nil
}

// FindByMerchantUserId retrieves a list of merchant records associated with the provided user ID.
//
// This function implements the MerchantQueryRepository interface and is responsible for executing the
// database query to retrieve a list of merchant records associated with the given user ID.
//
// Parameters:
//   - user_id: An integer containing the user ID used to identify the merchant records.
//
// Returns:
//   - A slice of pointers to MerchantRecord objects containing the merchant data retrieved from the
//     database.
//   - An error, if the query execution fails.
func (r *merchantQueryRepository) FindByMerchantUserId(ctx context.Context, user_id int) ([]*record.MerchantRecord, error) {
	res, err := r.db.GetMerchantsByUserID(ctx, int32(user_id))

	if err != nil {
		return nil, merchant_errors.ErrFindMerchantByUserIdFailed
	}

	return r.mapper.ToMerchantsRecord(res), nil
}
