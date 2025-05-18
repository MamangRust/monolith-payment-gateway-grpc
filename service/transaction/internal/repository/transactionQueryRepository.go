package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type transactionQueryRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.TransactionRecordMapping
}

func NewTransactionQueryRepository(db *db.Queries, ctx context.Context, mapping recordmapper.TransactionRecordMapping) *transactionQueryRepository {
	return &transactionQueryRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *transactionQueryRepository) FindAllTransactions(req *requests.FindAllTransactions) ([]*record.TransactionRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransactionsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	transactions, err := r.db.GetTransactions(r.ctx, reqDb)

	if err != nil {
		return nil, nil, transaction_errors.ErrFindAllTransactionsFailed
	}

	var totalCount int
	if len(transactions) > 0 {
		totalCount = int(transactions[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransactionsRecordAll(transactions), &totalCount, nil
}

func (r *transactionQueryRepository) FindAllTransactionByCardNumber(req *requests.FindAllTransactionCardNumber) ([]*record.TransactionRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransactionsByCardNumberParams{
		CardNumber: req.CardNumber,
		Column2:    req.Search,
		Limit:      int32(req.PageSize),
		Offset:     int32(offset),
	}

	transactions, err := r.db.GetTransactionsByCardNumber(r.ctx, reqDb)

	if err != nil {
		return nil, nil, transaction_errors.ErrFindTransactionsByCardNumberFailed
	}

	var totalCount int
	if len(transactions) > 0 {
		totalCount = int(transactions[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransactionsByCardNumberRecord(transactions), &totalCount, nil
}

func (r *transactionQueryRepository) FindByActive(req *requests.FindAllTransactions) ([]*record.TransactionRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveTransactionsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveTransactions(r.ctx, reqDb)

	if err != nil {
		return nil, nil, transaction_errors.ErrFindActiveTransactionsFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransactionsRecordActive(res), &totalCount, nil
}

func (r *transactionQueryRepository) FindByTrashed(req *requests.FindAllTransactions) ([]*record.TransactionRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedTransactionsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedTransactions(r.ctx, reqDb)

	if err != nil {
		return nil, nil, transaction_errors.ErrFindTrashedTransactionsFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransactionsRecordTrashed(res), &totalCount, nil
}

func (r *transactionQueryRepository) FindById(transaction_id int) (*record.TransactionRecord, error) {
	res, err := r.db.GetTransactionByID(r.ctx, int32(transaction_id))

	if err != nil {
		return nil, transaction_errors.ErrFindTransactionByIdFailed
	}

	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionQueryRepository) FindTransactionByMerchantId(merchant_id int) ([]*record.TransactionRecord, error) {
	res, err := r.db.GetTransactionsByMerchantID(r.ctx, int32(merchant_id))

	if err != nil {
		return nil, transaction_errors.ErrFindTransactionByMerchantIdFailed
	}

	return r.mapping.ToTransactionsRecord(res), nil
}
