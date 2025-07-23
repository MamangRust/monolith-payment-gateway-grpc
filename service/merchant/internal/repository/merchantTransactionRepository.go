package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/merchant"
)

// NewMerchantTransactionRepository creates a new instance of MerchantTransactionRepository
type merchantTransactionRepository struct {
	db     *db.Queries
	mapper recordmapper.MerchantTransactionRecordMapper
}

// NewMerchantTransactionRepository initializes a new instance of merchantTransactionRepository
// with the provided database queries, context, and merchant record mapper. This repository
// is responsible for executing transaction-related operations on merchant records in the database.
//
// Parameters:
//   - db: A pointer to the db.Queries object for executing database queries.
//   - ctx: The context to be used for database operations, allowing for cancellation and timeout.
//   - mapper: A MerchantRecordMapping that provides methods to map database rows to Merchant domain models.
//
// Returns:
//   - A pointer to the newly created merchantTransactionRepository instance.
func NewMerchantTransactionRepository(db *db.Queries, mapper recordmapper.MerchantTransactionRecordMapper) MerchantTransactionRepository {
	return &merchantTransactionRepository{
		db:     db,
		mapper: mapper,
	}
}

// FindAllTransactions retrieves a list of merchant transactions with pagination support.
// It returns a list of MerchantTransactionsRecord objects, the total count of records, and an error.
//
// Parameters:
//   - req: A FindAllMerchantTransactions object containing the search query, page number, and page size.
//
// Returns:
//   - A list of MerchantTransactionsRecord objects, which contain the merchant transaction data.
//   - The total count of records, which is used to calculate the number of pages.
//   - An error, which is non-nil if the operation fails.
func (r *merchantTransactionRepository) FindAllTransactions(ctx context.Context, req *requests.FindAllMerchantTransactions) ([]*record.MerchantTransactionsRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.FindAllTransactionsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	merchant, err := r.db.FindAllTransactions(ctx, reqDb)

	if err != nil {
		return nil, nil, merchant_errors.ErrFindAllTransactionsFailed
	}

	var totalCount int
	if len(merchant) > 0 {
		totalCount = int(merchant[0].TotalCount)
	} else {
		totalCount = 0
	}
	return r.mapper.ToMerchantsTransactionRecord(merchant), &totalCount, nil
}

// FindAllTransactionsByMerchant retrieves a list of merchant transactions for a specific merchant ID with pagination support.
// It returns a list of MerchantTransactionsRecord objects, the total count of records, and an error.
//
// Parameters:
//   - req: A FindAllMerchantTransactionsById object containing the search query, page number, page size, and merchant ID.
//
// Returns:
//   - A list of MerchantTransactionsRecord objects, which contain the merchant transaction data.
//   - The total count of records, which is used to calculate the number of pages.
//   - An error, which is non-nil if the operation fails.
func (r *merchantTransactionRepository) FindAllTransactionsByMerchant(ctx context.Context, req *requests.FindAllMerchantTransactionsById) ([]*record.MerchantTransactionsRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.FindAllTransactionsByMerchantParams{
		MerchantID: int32(req.MerchantID),
		Column2:    req.Search,
		Limit:      int32(req.PageSize),
		Offset:     int32(offset),
	}

	merchant, err := r.db.FindAllTransactionsByMerchant(ctx, reqDb)

	if err != nil {
		return nil, nil, merchant_errors.ErrFindAllTransactionsByMerchantFailed
	}

	var totalCount int
	if len(merchant) > 0 {
		totalCount = int(merchant[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToMerchantsTransactionByMerchantRecord(merchant), &totalCount, nil
}

// FindAllTransactionsByApikey retrieves a list of merchant transactions filtered by API key with pagination support.
// It returns a list of MerchantTransactionsRecord objects, the total count of records, and an error.
//
// Parameters:
//   - req: A FindAllMerchantTransactionsByApiKey object containing the API key, search query, page number, and page size.
//
// Returns:
//   - A list of MerchantTransactionsRecord objects, which contain the merchant transaction data.
//   - The total count of records, which is used to calculate the number of pages.
//   - An error, which is non-nil if the operation fails.
func (r *merchantTransactionRepository) FindAllTransactionsByApikey(ctx context.Context, req *requests.FindAllMerchantTransactionsByApiKey) ([]*record.MerchantTransactionsRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.FindAllTransactionsByApikeyParams{
		ApiKey:  req.ApiKey,
		Column2: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	merchant, err := r.db.FindAllTransactionsByApikey(ctx, reqDb)

	if err != nil {
		return nil, nil, merchant_errors.ErrFindAllTransactionsByApiKeyFailed
	}

	var totalCount int
	if len(merchant) > 0 {
		totalCount = int(merchant[0].TotalCount)
	} else {
		totalCount = 0
	}
	return r.mapper.ToMerchantsTransactionByApikeyRecord(merchant), &totalCount, nil
}
