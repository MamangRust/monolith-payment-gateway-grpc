package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/transaction"
)

// transactionQueryRepository is a repository for handling transaction query operations.
type transactionQueryRepository struct {
	db     *db.Queries
	mapper recordmapper.TransactionQueryRecordMapper
}

// NewTransactionQueryRepository initializes a new instance of transactionQueryRepository.
//
// Parameters:
//   - db: A pointer to the db.Queries object for executing database queries.
//   - ctx: The context to be used for database operations, allowing for cancellation and timeout.
//   - mapper: A TransactionRecordMapping that provides methods to map database rows to Transaction domain models.
//
// Returns:
//   - A pointer to the newly created transactionQueryRepository instance.
func NewTransactionQueryRepository(db *db.Queries, mapper recordmapper.TransactionQueryRecordMapper) TransactionQueryRepository {
	return &transactionQueryRepository{
		db:     db,
		mapper: mapper,
	}
}

// FindAllTransactions retrieves all transactions with optional filtering and pagination.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The filter and pagination parameters.
//
// Returns:
//   - []*record.TransactionRecord: A list of transaction records.
//   - *int: The total count of transactions.
//   - error: Error if something went wrong during the query.
func (r *transactionQueryRepository) FindAllTransactions(ctx context.Context, req *requests.FindAllTransactions) ([]*record.TransactionRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransactionsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	transactions, err := r.db.GetTransactions(ctx, reqDb)

	if err != nil {
		return nil, nil, transaction_errors.ErrFindAllTransactionsFailed
	}

	var totalCount int
	if len(transactions) > 0 {
		totalCount = int(transactions[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToTransactionsRecordAll(transactions), &totalCount, nil
}

// FindAllTransactionByCardNumber retrieves all transactions by a specific card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing the card number and optional filters.
//
// Returns:
//   - []*record.TransactionRecord: A list of transactions associated with the card number.
//   - *int: The total count of transactions.
//   - error: Error if something went wrong during the query.
func (r *transactionQueryRepository) FindAllTransactionByCardNumber(ctx context.Context, req *requests.FindAllTransactionCardNumber) ([]*record.TransactionRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransactionsByCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    req.Search,
		Limit:      int32(req.PageSize),
		Offset:     int32(offset),
	}

	transactions, err := r.db.GetTransactionsByCardNumber(ctx, reqDb)

	if err != nil {
		return nil, nil, transaction_errors.ErrFindTransactionsByCardNumberFailed
	}

	var totalCount int
	if len(transactions) > 0 {
		totalCount = int(transactions[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToTransactionsByCardNumberRecord(transactions), &totalCount, nil
}

// FindByActive retrieves all active (non-deleted) transactions with optional filtering and pagination.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The filter and pagination parameters.
//
// Returns:
//   - []*record.TransactionRecord: A list of active transaction records.
//   - *int: The total count of active transactions.
//   - error: Error if something went wrong during the query.
func (r *transactionQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllTransactions) ([]*record.TransactionRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveTransactionsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveTransactions(ctx, reqDb)

	if err != nil {
		return nil, nil, transaction_errors.ErrFindActiveTransactionsFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToTransactionsRecordActive(res), &totalCount, nil
}

// FindByTrashed retrieves all soft-deleted (trashed) transactions with optional filtering and pagination.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The filter and pagination parameters.
//
// Returns:
//   - []*record.TransactionRecord: A list of trashed transaction records.
//   - *int: The total count of trashed transactions.
//   - error: Error if something went wrong during the query.
func (r *transactionQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllTransactions) ([]*record.TransactionRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedTransactionsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedTransactions(ctx, reqDb)

	if err != nil {
		return nil, nil, transaction_errors.ErrFindTrashedTransactionsFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToTransactionsRecordTrashed(res), &totalCount, nil
}

// FindById retrieves a transaction by its ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - transaction_id: The ID of the transaction.
//
// Returns:
//   - *record.TransactionRecord: The transaction record if found.
//   - error: Error if something went wrong during the query.
func (r *transactionQueryRepository) FindById(ctx context.Context, transaction_id int) (*record.TransactionRecord, error) {
	res, err := r.db.GetTransactionByID(ctx, int32(transaction_id))

	if err != nil {
		return nil, transaction_errors.ErrFindTransactionByIdFailed
	}

	return r.mapper.ToTransactionRecord(res), nil
}

// FindTransactionByMerchantId retrieves all transactions associated with a specific merchant.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - merchant_id: The ID of the merchant.
//
// Returns:
//   - []*record.TransactionRecord: A list of transactions for the given merchant.
//   - error: Error if something went wrong during the query.
func (r *transactionQueryRepository) FindTransactionByMerchantId(ctx context.Context, merchant_id int) ([]*record.TransactionRecord, error) {
	res, err := r.db.GetTransactionsByMerchantID(ctx, int32(merchant_id))

	if err != nil {
		return nil, transaction_errors.ErrFindTransactionByMerchantIdFailed
	}

	return r.mapper.ToTransactionsRecord(res), nil
}
