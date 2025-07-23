package merchantstatsservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type MerchantStatsAmountService interface {
	// FindMonthlyAmountMerchant retrieves the monthly transaction amount statistics for a merchant.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the monthly amount statistics should be retrieved.
	//
	// Returns:
	//   - []*response.MerchantResponseMonthlyAmount: A slice of monthly amount statistics.
	//   - *response.ErrorResponse: An error returned if the retrieval fails.
	FindMonthlyAmountMerchant(ctx context.Context, year int) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse)

	// FindYearlyAmountMerchant retrieves the yearly transaction amount statistics for a merchant.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the yearly amount statistics should be retrieved.
	//
	// Returns:
	//   - []*response.MerchantResponseYearlyAmount: A slice of yearly amount statistics.
	//   - *response.ErrorResponse: An error returned if the retrieval fails.
	FindYearlyAmountMerchant(ctx context.Context, year int) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse)
}

type MerchantStatsMethodService interface {
	// FindMonthlyPaymentMethodsMerchant retrieves monthly payment method statistics for a merchant.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the monthly payment method statistics should be retrieved.
	//
	// Returns:
	//   - []*response.MerchantResponseMonthlyPaymentMethod: A slice of monthly payment method statistics.
	//   - *response.ErrorResponse: An error returned if the retrieval fails.
	FindMonthlyPaymentMethodsMerchant(ctx context.Context, year int) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse)

	// FindYearlyPaymentMethodMerchant retrieves the yearly payment methods for a merchant.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the yearly payment methods should be retrieved.
	//
	// Returns:
	//   - []*response.MerchantResponseYearlyPaymentMethod: A slice of yearly payment method statistics.
	//   - *response.ErrorResponse: An error returned if the retrieval fails.
	FindYearlyPaymentMethodMerchant(ctx context.Context, year int) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse)
}

type MerchantStatsTotalAmountService interface {
	// FindMonthlyTotalAmountMerchant retrieves the monthly total transaction amounts for a merchant.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the monthly total amount statistics should be retrieved.
	//
	// Returns:
	//   - []*response.MerchantResponseMonthlyTotalAmount: A slice of monthly total amount statistics.
	//   - *response.ErrorResponse: An error returned if the retrieval fails.
	FindMonthlyTotalAmountMerchant(ctx context.Context, year int) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse)

	// FindYearlyTotalAmountMerchant retrieves the yearly total transaction amounts for a merchant.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the yearly total amount statistics should be retrieved.
	//
	// Returns:
	//   - []*response.MerchantResponseYearlyTotalAmount: A slice of yearly total amount statistics.
	//   - *response.ErrorResponse: An error returned if the retrieval fails.
	FindYearlyTotalAmountMerchant(ctx context.Context, year int) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse)
}
