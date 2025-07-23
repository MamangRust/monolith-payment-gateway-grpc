package merchantstatsrepository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
)

type MerchantStatsMethodRepository interface {
	// GetMonthlyPaymentMethodsMerchant retrieves monthly merchant payment method statistics for a given year.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the data is requested.
	//
	// Returns:
	//   - []*record.MerchantMonthlyPaymentMethod: The list of monthly payment method records.
	//   - error: An error if any occurred during the query.
	GetMonthlyPaymentMethodsMerchant(ctx context.Context, year int) ([]*record.MerchantMonthlyPaymentMethod, error)

	// GetYearlyPaymentMethodMerchant retrieves yearly merchant payment method statistics.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the data is requested.
	//
	// Returns:
	//   - []*record.MerchantYearlyPaymentMethod: The list of yearly payment method records.
	//   - error: An error if any occurred during the query.
	GetYearlyPaymentMethodMerchant(ctx context.Context, year int) ([]*record.MerchantYearlyPaymentMethod, error)
}

type MerchantStatsAmountRepository interface {
	// GetMonthlyAmountMerchant retrieves monthly transaction amount statistics for merchants.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the data is requested.
	//
	// Returns:
	//   - []*record.MerchantMonthlyAmount: The list of monthly amount records.
	//   - error: An error if any occurred during the query.
	GetMonthlyAmountMerchant(ctx context.Context, year int) ([]*record.MerchantMonthlyAmount, error)

	// GetYearlyAmountMerchant retrieves yearly transaction amount statistics for merchants.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the data is requested.
	//
	// Returns:
	//   - []*record.MerchantYearlyAmount: The list of yearly amount records.
	//   - error: An error if any occurred during the query.
	GetYearlyAmountMerchant(ctx context.Context, year int) ([]*record.MerchantYearlyAmount, error)
}

type MerchantStatsTotalAmountRepository interface {
	// GetMonthlyTotalAmountMerchant retrieves monthly total transaction amount statistics across all merchants.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the data is requested.
	//
	// Returns:
	//   - []*record.MerchantMonthlyTotalAmount: The list of monthly total amount records.
	//   - error: An error if any occurred during the query.
	GetMonthlyTotalAmountMerchant(ctx context.Context, year int) ([]*record.MerchantMonthlyTotalAmount, error)

	// GetYearlyTotalAmountMerchant retrieves yearly total transaction amount statistics across all merchants.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the data is requested.
	//
	// Returns:
	//   - []*record.MerchantYearlyTotalAmount: The list of yearly total amount records.
	//   - error: An error if any occurred during the query.
	GetYearlyTotalAmountMerchant(ctx context.Context, year int) ([]*record.MerchantYearlyTotalAmount, error)
}
