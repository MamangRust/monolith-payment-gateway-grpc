package merchantstatsbymerchantservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type MerchantStatsByMerchantAmountService interface {
	// FindMonthlyAmountByMerchants retrieves monthly transaction amount statistics for a merchant.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing merchant identifier and year.
	//
	// Returns:
	//   - []*response.MerchantResponseMonthlyAmount: A slice of monthly transaction amount statistics.
	//   - *response.ErrorResponse: An error returned if the retrieval fails.
	FindMonthlyAmountByMerchants(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse)

	// FindYearlyAmountByMerchants retrieves yearly transaction amount statistics for a merchant.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing merchant identifier and year.
	//
	// Returns:
	//   - []*response.MerchantResponseYearlyAmount: A slice of yearly transaction amount statistics.
	//   - *response.ErrorResponse: An error returned if the retrieval fails.
	FindYearlyAmountByMerchants(ctx context.Context, req *requests.MonthYearAmountMerchant) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse)
}

type MerchantStatsByMerchantMethodService interface {
	// FindMonthlyPaymentMethodByMerchants retrieves monthly payment method statistics for a merchant.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing merchant identifier and year.
	//
	// Returns:
	//   - []*response.MerchantResponseMonthlyPaymentMethod: A slice of monthly payment method statistics.
	//   - *response.ErrorResponse: An error returned if the retrieval fails.
	FindMonthlyPaymentMethodByMerchants(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse)

	// FindYearlyPaymentMethodByMerchants retrieves yearly payment method statistics for a merchant.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing merchant identifier and year.
	//
	// Returns:
	//   - []*response.MerchantResponseYearlyPaymentMethod: A slice of yearly payment method statistics.
	//   - *response.ErrorResponse: An error returned if the retrieval fails.
	FindYearlyPaymentMethodByMerchants(ctx context.Context, req *requests.MonthYearPaymentMethodMerchant) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse)
}

type MerchantStatsByMerchantTotalAmountService interface {
	// FindMonthlyTotalAmountByMerchants retrieves monthly total transaction amounts for a merchant.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing merchant identifier and year.
	//
	// Returns:
	//   - []*response.MerchantResponseMonthlyTotalAmount: A slice of monthly total amount statistics.
	//   - *response.ErrorResponse: An error returned if the retrieval fails.
	FindMonthlyTotalAmountByMerchants(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse)

	// FindYearlyTotalAmountByMerchants retrieves yearly total transaction amounts for a merchant.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing merchant identifier and year.
	//
	// Returns:
	//   - []*response.MerchantResponseYearlyTotalAmount: A slice of yearly total amount statistics.
	//   - *response.ErrorResponse: An error returned if the retrieval fails.
	FindYearlyTotalAmountByMerchants(ctx context.Context, req *requests.MonthYearTotalAmountMerchant) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse)
}
