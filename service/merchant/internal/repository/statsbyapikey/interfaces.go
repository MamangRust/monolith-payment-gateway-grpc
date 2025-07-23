package merchantstatsapikeyrepository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type MerchantStatsMethodByApiKeyRepository interface {
	// GetMonthlyPaymentMethodByApikey retrieves monthly payment method statistics for a specific merchant using API key.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing month, year, and API key.
	//
	// Returns:
	//   - []*record.MerchantMonthlyPaymentMethod: The list of monthly payment method records.
	//   - error: An error if any occurred during the query.
	GetMonthlyPaymentMethodByApikey(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*record.MerchantMonthlyPaymentMethod, error)

	// GetYearlyPaymentMethodByApikey retrieves yearly payment method statistics for a specific merchant using API key.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing year and API key.
	//
	// Returns:
	//   - []*record.MerchantYearlyPaymentMethod: The list of yearly payment method records.
	//   - error: An error if any occurred during the query.
	GetYearlyPaymentMethodByApikey(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*record.MerchantYearlyPaymentMethod, error)
}

type MerchantStatsAmountByApiKeyRepository interface {
	// GetMonthlyAmountByApikey retrieves monthly transaction amount statistics for a specific merchant using API key.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing month, year, and API key.
	//
	// Returns:
	//   - []*record.MerchantMonthlyAmount: The list of monthly amount records.
	//   - error: An error if any occurred during the query.
	GetMonthlyAmountByApikey(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*record.MerchantMonthlyAmount, error)

	// GetYearlyAmountByApikey retrieves yearly transaction amount statistics for a specific merchant using API key.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing year and API key.
	//
	// Returns:
	//   - []*record.MerchantYearlyAmount: The list of yearly amount records.
	//   - error: An error if any occurred during the query.
	GetYearlyAmountByApikey(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*record.MerchantYearlyAmount, error)
}

type MerchantStatsTotalAmountByApiKeyRepository interface {
	// GetMonthlyTotalAmountByApikey retrieves monthly total transaction amount statistics for a specific merchant using API key.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing month, year, and API key.
	//
	// Returns:
	//   - []*record.MerchantMonthlyTotalAmount: The list of monthly total amount records.
	//   - error: An error if any occurred during the query.
	GetMonthlyTotalAmountByApikey(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*record.MerchantMonthlyTotalAmount, error)

	// GetYearlyTotalAmountByApikey retrieves yearly total transaction amount statistics for a specific merchant using API key.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing year and API key.
	//
	// Returns:
	//   - []*record.MerchantYearlyTotalAmount: The list of yearly total amount records.
	//   - error: An error if any occurred during the query.
	GetYearlyTotalAmountByApikey(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*record.MerchantYearlyTotalAmount, error)
}
