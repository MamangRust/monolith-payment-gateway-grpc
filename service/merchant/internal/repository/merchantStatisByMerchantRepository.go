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

type merchantStatisticByMerchantRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.MerchantRecordMapping
}

func NewMerchantStatisticByMerchantRepository(db *db.Queries, ctx context.Context, mapping recordmapper.MerchantRecordMapping) *merchantStatisticByMerchantRepository {
	return &merchantStatisticByMerchantRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *merchantStatisticByMerchantRepository) GetMonthlyPaymentMethodByMerchants(req *requests.MonthYearPaymentMethodMerchant) ([]*record.MerchantMonthlyPaymentMethod, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyPaymentMethodByMerchants(r.ctx, db.GetMonthlyPaymentMethodByMerchantsParams{
		MerchantID: int32(req.MerchantID),
		Column1:    yearStart,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyPaymentMethodByMerchantsFailed
	}

	return r.mapping.ToMerchantMonthlyPaymentMethodsByMerchant(res), nil
}

func (r *merchantStatisticByMerchantRepository) GetYearlyPaymentMethodByMerchants(req *requests.MonthYearPaymentMethodMerchant) ([]*record.MerchantYearlyPaymentMethod, error) {
	res, err := r.db.GetYearlyPaymentMethodByMerchants(r.ctx, db.GetYearlyPaymentMethodByMerchantsParams{
		MerchantID: int32(req.MerchantID),
		Column2:    req.Year,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyPaymentMethodByMerchantsFailed
	}

	return r.mapping.ToMerchantYearlyPaymentMethodsByMerchant(res), nil
}

func (r *merchantStatisticByMerchantRepository) GetMonthlyAmountByMerchants(req *requests.MonthYearAmountMerchant) ([]*record.MerchantMonthlyAmount, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyAmountByMerchants(r.ctx, db.GetMonthlyAmountByMerchantsParams{
		MerchantID: int32(req.MerchantID),
		Column1:    yearStart,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyAmountByMerchantsFailed
	}

	return r.mapping.ToMerchantMonthlyAmountsByMerchant(res), nil
}

func (r *merchantStatisticByMerchantRepository) GetYearlyAmountByMerchants(req *requests.MonthYearAmountMerchant) ([]*record.MerchantYearlyAmount, error) {
	res, err := r.db.GetYearlyAmountByMerchants(r.ctx, db.GetYearlyAmountByMerchantsParams{
		MerchantID: int32(req.MerchantID),
		Column2:    req.Year,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyAmountByMerchantsFailed
	}

	return r.mapping.ToMerchantYearlyAmountsByMerchant(res), nil
}

func (r *merchantStatisticByMerchantRepository) GetMonthlyTotalAmountByMerchants(req *requests.MonthYearTotalAmountMerchant) ([]*record.MerchantMonthlyTotalAmount, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyTotalAmountByMerchant(r.ctx, db.GetMonthlyTotalAmountByMerchantParams{
		Column2: int32(req.MerchantID),
		Column1: yearStart,
	})

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyTotalAmountByMerchantsFailed
	}

	return r.mapping.ToMerchantMonthlyTotalAmountsByMerchant(res), nil
}

func (r *merchantStatisticByMerchantRepository) GetYearlyTotalAmountByMerchants(req *requests.MonthYearTotalAmountMerchant) ([]*record.MerchantYearlyTotalAmount, error) {
	res, err := r.db.GetYearlyTotalAmountByMerchant(r.ctx, db.GetYearlyTotalAmountByMerchantParams{
		Column2: int32(req.MerchantID),
		Column1: int32(req.Year),
	})

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyTotalAmountByMerchantsFailed
	}

	return r.mapping.ToMerchantYearlyTotalAmountsByMerchant(res), nil
}
