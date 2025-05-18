package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type merchantStatisticByApiKeyRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.MerchantRecordMapping
}

func NewMerchantStatisticByApiKeyRepository(db *db.Queries, ctx context.Context, mapping recordmapper.MerchantRecordMapping) *merchantStatisticByApiKeyRepository {
	return &merchantStatisticByApiKeyRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *merchantStatisticByApiKeyRepository) GetMonthlyPaymentMethodByApikey(req *requests.MonthYearPaymentMethodApiKey) ([]*record.MerchantMonthlyPaymentMethod, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyPaymentMethodByApikey(r.ctx, db.GetMonthlyPaymentMethodByApikeyParams{
		ApiKey:  req.Apikey,
		Column1: yearStart,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyPaymentMethodByApikeyFailed
	}

	return r.mapping.ToMerchantMonthlyPaymentMethodsByApikey(res), nil
}

func (r *merchantStatisticByApiKeyRepository) GetYearlyPaymentMethodByApikey(req *requests.MonthYearPaymentMethodApiKey) ([]*record.MerchantYearlyPaymentMethod, error) {
	res, err := r.db.GetYearlyPaymentMethodByApikey(r.ctx, db.GetYearlyPaymentMethodByApikeyParams{
		ApiKey:  req.Apikey,
		Column2: req.Year,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyPaymentMethodByApikeyFailed
	}

	return r.mapping.ToMerchantYearlyPaymentMethodsByApikey(res), nil
}

func (r *merchantStatisticByApiKeyRepository) GetMonthlyAmountByApikey(req *requests.MonthYearAmountApiKey) ([]*record.MerchantMonthlyAmount, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyAmountByApikey(r.ctx, db.GetMonthlyAmountByApikeyParams{
		ApiKey:  req.Apikey,
		Column1: yearStart,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyAmountByApikeyFailed
	}

	return r.mapping.ToMerchantMonthlyAmountsByApikey(res), nil
}

func (r *merchantStatisticByApiKeyRepository) GetYearlyAmountByApikey(req *requests.MonthYearAmountApiKey) ([]*record.MerchantYearlyAmount, error) {
	res, err := r.db.GetYearlyAmountByApikey(r.ctx, db.GetYearlyAmountByApikeyParams{
		ApiKey:  req.Apikey,
		Column2: req.Year,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyAmountByApikeyFailed
	}

	return r.mapping.ToMerchantYearlyAmountsByApikey(res), nil
}

func (r *merchantStatisticByApiKeyRepository) GetMonthlyTotalAmountByApikey(req *requests.MonthYearTotalAmountApiKey) ([]*record.MerchantMonthlyTotalAmount, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyTotalAmountByApikey(r.ctx, db.GetMonthlyTotalAmountByApikeyParams{
		ApiKey:  req.Apikey,
		Column1: yearStart,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyTotalAmountByApikeyFailed
	}

	return r.mapping.ToMerchantMonthlyTotalAmountsByApikey(res), nil
}

func (r *merchantStatisticByApiKeyRepository) GetYearlyTotalAmountByApikey(req *requests.MonthYearTotalAmountApiKey) ([]*record.MerchantYearlyTotalAmount, error) {
	res, err := r.db.GetYearlyTotalAmountByApikey(r.ctx, db.GetYearlyTotalAmountByApikeyParams{
		ApiKey:  req.Apikey,
		Column1: int32(req.Year),
	})

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyTotalAmountByApikeyFailed
	}

	return r.mapping.ToMerchantYearlyTotalAmountsByApikey(res), nil
}
