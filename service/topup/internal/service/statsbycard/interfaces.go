package topupstatsbycardservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type TopupStatsByCardAmountService interface {
	// FindMonthlyTopupAmountsByCardNumber retrieves monthly topup amount statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//
	// Returns:
	//   - []*response.TopupMonthAmountResponse: List of monthly topup amount statistics.
	//   - *response.ErrorResponse: Error details if retrieval fails.
	FindMonthlyTopupAmountsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*response.TopupMonthAmountResponse, *response.ErrorResponse)

	// FindYearlyTopupAmountsByCardNumber retrieves yearly topup amount statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//
	// Returns:
	//   - []*response.TopupYearlyAmountResponse: List of yearly topup amount statistics.
	//   - *response.ErrorResponse: Error details if retrieval fails.
	FindYearlyTopupAmountsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*response.TopupYearlyAmountResponse, *response.ErrorResponse)
}

type TopupStatsByCardMethodService interface {
	// FindMonthlyTopupMethodsByCardNumber retrieves monthly topup method statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//
	// Returns:
	//   - []*response.TopupMonthMethodResponse: List of monthly topup method statistics.
	//   - *response.ErrorResponse: Error details if retrieval fails.
	FindMonthlyTopupMethodsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*response.TopupMonthMethodResponse, *response.ErrorResponse)

	// FindYearlyTopupMethodsByCardNumber retrieves yearly topup method statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//
	// Returns:
	//   - []*response.TopupYearlyMethodResponse: List of yearly topup method statistics.
	//   - *response.ErrorResponse: Error details if retrieval fails.
	FindYearlyTopupMethodsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*response.TopupYearlyMethodResponse, *response.ErrorResponse)
}

type TopupStatsByCardStatusService interface {
	// FindMonthTopupStatusSuccessByCardNumber retrieves monthly successful topup statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing month, year, and card number.
	//
	// Returns:
	//   - []*response.TopupResponseMonthStatusSuccess: List of monthly successful topup statistics.
	//   - *response.ErrorResponse: Error details if retrieval fails.
	FindMonthTopupStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*response.TopupResponseMonthStatusSuccess, *response.ErrorResponse)

	// FindYearlyTopupStatusSuccessByCardNumber retrieves yearly successful topup statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//
	// Returns:
	//   - []*response.TopupResponseYearStatusSuccess: List of yearly successful topup statistics.
	//   - *response.ErrorResponse: Error details if retrieval fails.
	FindYearlyTopupStatusSuccessByCardNumber(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*response.TopupResponseYearStatusSuccess, *response.ErrorResponse)

	// FindMonthTopupStatusFailedByCardNumber retrieves monthly failed topup statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing month, year, and card number.
	//
	// Returns:
	//   - []*response.TopupResponseMonthStatusFailed: List of monthly failed topup statistics.
	//   - *response.ErrorResponse: Error details if retrieval fails.
	FindMonthTopupStatusFailedByCardNumber(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*response.TopupResponseMonthStatusFailed, *response.ErrorResponse)

	// FindYearlyTopupStatusFailedByCardNumber retrieves yearly failed topup statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing year and card number.
	//
	// Returns:
	//   - []*response.TopupResponseYearStatusFailed: List of yearly failed topup statistics.
	//   - *response.ErrorResponse: Error details if retrieval fails.
	FindYearlyTopupStatusFailedByCardNumber(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*response.TopupResponseYearStatusFailed, *response.ErrorResponse)
}
