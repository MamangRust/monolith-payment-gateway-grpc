package merchantstatsmerchantrepository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type MerchantStatsMethodByMerchantRepository interface {
	// GetMonthlyPaymentMethodByMerchants retrieves monthly payment method statistics for a specific merchant.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing month, year, and merchant ID.
	//
	// Returns:
	//   - []*record.MerchantMonthlyPaymentMethod: The list of monthly payment method records.
	//   - error: An error if any occurred during the query.
	GetMonthlyPaymentMethodByMerchants(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*record.MerchantMonthlyPaymentMethod, error)

	// GetYearlyPaymentMethodByMerchants retrieves yearly payment method statistics for a specific merchant.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing year and merchant ID.
	//
	// Returns:
	//   - []*record.MerchantYearlyPaymentMethod: The list of yearly payment method records.
	//   - error: An error if any occurred during the query.
	GetYearlyPaymentMethodByMerchants(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*record.MerchantYearlyPaymentMethod, error)
}

type MerchantStatsAmountByMerchantRepository interface {
	// GetMonthlyAmountByMerchants retrieves monthly transaction amount statistics for a specific merchant.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing month, year, and merchant ID.
	//
	// Returns:
	//   - []*record.MerchantMonthlyAmount: The list of monthly amount records.
	//   - error: An error if any occurred during the query.
	GetMonthlyAmountByMerchants(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*record.MerchantMonthlyAmount, error)

	// GetYearlyAmountByMerchants retrieves yearly transaction amount statistics for a specific merchant.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing year and merchant ID.
	//
	// Returns:
	//   - []*record.MerchantYearlyAmount: The list of yearly amount records.
	//   - error: An error if any occurred during the query.
	GetYearlyAmountByMerchants(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*record.MerchantYearlyAmount, error)
}

type MerchantStatsTotalAmountByMerchantRepository interface {
	// GetMonthlyTotalAmountByMerchants retrieves monthly total amount statistics for a specific merchant.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing month, year, and merchant ID.
	//
	// Returns:
	//   - []*record.MerchantMonthlyTotalAmount: The list of monthly total amount records.
	//   - error: An error if any occurred during the query.
	GetMonthlyTotalAmountByMerchants(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) ([]*record.MerchantMonthlyTotalAmount, error)

	// GetYearlyTotalAmountByMerchants retrieves yearly total amount statistics for a specific merchant.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing year and merchant ID.
	//
	// Returns:
	//   - []*record.MerchantYearlyTotalAmount: The list of yearly total amount records.
	//   - error: An error if any occurred during the query.
	GetYearlyTotalAmountByMerchants(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) ([]*record.MerchantYearlyTotalAmount, error)
}
