package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type merchantTransactionRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.MerchantRecordMapping
}

func NewMerchantTransactionRepository(db *db.Queries, ctx context.Context, mapping recordmapper.MerchantRecordMapping) *merchantTransactionRepository {
	return &merchantTransactionRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *merchantTransactionRepository) FindAllTransactions(req *requests.FindAllMerchantTransactions) ([]*record.MerchantTransactionsRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.FindAllTransactionsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	merchant, err := r.db.FindAllTransactions(r.ctx, reqDb)

	if err != nil {
		return nil, nil, merchant_errors.ErrFindAllTransactionsFailed
	}

	var totalCount int
	if len(merchant) > 0 {
		totalCount = int(merchant[0].TotalCount)
	} else {
		totalCount = 0
	}
	return r.mapping.ToMerchantsTransactionRecord(merchant), &totalCount, nil
}

func (r *merchantTransactionRepository) FindAllTransactionsByMerchant(req *requests.FindAllMerchantTransactionsById) ([]*record.MerchantTransactionsRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.FindAllTransactionsByMerchantParams{
		MerchantID: int32(req.MerchantID),
		Column2:    req.Search,
		Limit:      int32(req.PageSize),
		Offset:     int32(offset),
	}

	merchant, err := r.db.FindAllTransactionsByMerchant(r.ctx, reqDb)

	if err != nil {
		return nil, nil, merchant_errors.ErrFindAllTransactionsByMerchantFailed
	}

	var totalCount int
	if len(merchant) > 0 {
		totalCount = int(merchant[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantsTransactionByMerchantRecord(merchant), &totalCount, nil
}

func (r *merchantTransactionRepository) FindAllTransactionsByApikey(req *requests.FindAllMerchantTransactionsByApiKey) ([]*record.MerchantTransactionsRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.FindAllTransactionsByApikeyParams{
		ApiKey:  req.ApiKey,
		Column2: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	merchant, err := r.db.FindAllTransactionsByApikey(r.ctx, reqDb)

	if err != nil {
		return nil, nil, merchant_errors.ErrFindAllTransactionsByApiKeyFailed
	}

	var totalCount int
	if len(merchant) > 0 {
		totalCount = int(merchant[0].TotalCount)
	} else {
		totalCount = 0
	}
	return r.mapping.ToMerchantsTransactionByApikeyRecord(merchant), &totalCount, nil
}
