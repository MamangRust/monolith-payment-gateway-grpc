package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type merchantStatisticRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.MerchantRecordMapping
}

func NewMerchantStatisticRepository(db *db.Queries, ctx context.Context, mapping recordmapper.MerchantRecordMapping) *merchantStatisticRepository {
	return &merchantStatisticRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *merchantStatisticRepository) GetMonthlyPaymentMethodsMerchant(year int) ([]*record.MerchantMonthlyPaymentMethod, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyPaymentMethodsMerchant(r.ctx, yearStart)

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyPaymentMethodsMerchantFailed
	}

	return r.mapping.ToMerchantMonthlyPaymentMethods(res), nil
}

func (r *merchantStatisticRepository) GetYearlyPaymentMethodMerchant(year int) ([]*record.MerchantYearlyPaymentMethod, error) {
	res, err := r.db.GetYearlyPaymentMethodMerchant(r.ctx, year)

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyPaymentMethodMerchantFailed
	}

	return r.mapping.ToMerchantYearlyPaymentMethods(res), nil

}

func (r *merchantStatisticRepository) GetMonthlyAmountMerchant(year int) ([]*record.MerchantMonthlyAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyAmountMerchant(r.ctx, yearStart)

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyAmountMerchantFailed
	}

	return r.mapping.ToMerchantMonthlyAmounts(res), nil
}

func (r *merchantStatisticRepository) GetYearlyAmountMerchant(year int) ([]*record.MerchantYearlyAmount, error) {
	res, err := r.db.GetYearlyAmountMerchant(r.ctx, year)

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyAmountMerchantFailed
	}

	return r.mapping.ToMerchantYearlyAmounts(res), nil
}

func (r *merchantStatisticRepository) GetMonthlyTotalAmountMerchant(year int) ([]*record.MerchantMonthlyTotalAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyTotalAmountMerchant(r.ctx, yearStart)

	if err != nil {
		return nil, merchant_errors.ErrGetMonthlyTotalAmountMerchantFailed
	}

	return r.mapping.ToMerchantMonthlyTotalAmounts(res), nil
}

func (r *merchantStatisticRepository) GetYearlyTotalAmountMerchant(year int) ([]*record.MerchantYearlyTotalAmount, error) {
	res, err := r.db.GetYearlyTotalAmountMerchant(r.ctx, int32(year))

	if err != nil {
		return nil, merchant_errors.ErrGetYearlyTotalAmountMerchantFailed
	}

	return r.mapping.ToMerchantYearlyTotalAmounts(res), nil
}
