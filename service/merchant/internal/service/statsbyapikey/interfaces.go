package merchantstatsbyapikeyservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type MerchantStatsByApiKeyAmountService interface {
	// FindMonthlyAmountByApikeys retrieves monthly transaction amount statistics for a merchant using API key.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing the API key and the target year.
	//
	// Returns:
	//   - []*response.MerchantResponseMonthlyAmount: A slice of monthly transaction amount statistics.
	//   - *response.ErrorResponse: An error returned if the retrieval fails.
	FindMonthlyAmountByApikeys(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse)

	// FindYearlyAmountByApikeys retrieves yearly transaction amount statistics for a merchant using API key.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing the API key and the target year.
	//
	// Returns:
	//   - []*response.MerchantResponseYearlyAmount: A slice of yearly transaction amount statistics.
	//   - *response.ErrorResponse: An error returned if the retrieval fails.
	FindYearlyAmountByApikeys(ctx context.Context, req *requests.MonthYearAmountApiKey) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse)
}

type MerchantStatsByApiKeyMethodService interface {
	// FindMonthlyPaymentMethodByApikeys retrieves monthly payment method statistics for a merchant using API key.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing the API key and the target year.
	//
	// Returns:
	//   - []*response.MerchantResponseMonthlyPaymentMethod: A slice of monthly payment method statistics.
	//   - *response.ErrorResponse: An error returned if the retrieval fails.
	FindMonthlyPaymentMethodByApikeys(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse)

	// FindYearlyPaymentMethodByApikeys retrieves yearly payment method statistics for a merchant using API key.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing the API key and the target year.
	//
	// Returns:
	//   - []*response.MerchantResponseYearlyPaymentMethod: A slice of yearly payment method statistics.
	//   - *response.ErrorResponse: An error returned if the retrieval fails.
	FindYearlyPaymentMethodByApikeys(ctx context.Context, req *requests.MonthYearPaymentMethodApiKey) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse)
}

type MerchantStatsByApiKeyTotalAmountService interface {
	// FindMonthlyTotalAmountByApikeys retrieves monthly total transaction amounts for a merchant using API key.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing the API key and the target year.
	//
	// Returns:
	//   - []*response.MerchantResponseMonthlyTotalAmount: A slice of monthly total amount statistics.
	//   - *response.ErrorResponse: An error returned if the retrieval fails.
	FindMonthlyTotalAmountByApikeys(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse)

	// FindYearlyTotalAmountByApikeys retrieves yearly total transaction amounts for a merchant using API key.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing the API key and the target year.
	//
	// Returns:
	//   - []*response.MerchantResponseYearlyTotalAmount: A slice of yearly total amount statistics.
	//   - *response.ErrorResponse: An error returned if the retrieval fails.
	FindYearlyTotalAmountByApikeys(ctx context.Context, req *requests.MonthYearTotalAmountApiKey) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse)
}
