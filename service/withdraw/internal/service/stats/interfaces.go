package withdrawstatsservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type WithdrawStatsStatusService interface {
	// FindMonthWithdrawStatusSuccess retrieves monthly successful withdraw statistics.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing the month and year for filtering.
	//
	// Returns:
	//   - []*response.WithdrawResponseMonthStatusSuccess: List of successful monthly withdraw statistics.
	//   - *response.ErrorResponse: Error information if any occurred.
	FindMonthWithdrawStatusSuccess(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*response.WithdrawResponseMonthStatusSuccess, *response.ErrorResponse)

	// FindYearlyWithdrawStatusSuccess retrieves yearly successful withdraw statistics.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year to filter the data.
	//
	// Returns:
	//   - []*response.WithdrawResponseYearStatusSuccess: List of successful yearly withdraw statistics.
	//   - *response.ErrorResponse: Error information if any occurred.
	FindYearlyWithdrawStatusSuccess(ctx context.Context, year int) ([]*response.WithdrawResponseYearStatusSuccess, *response.ErrorResponse)

	// FindMonthWithdrawStatusFailed retrieves monthly failed withdraw statistics.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing the month and year for filtering.
	//
	// Returns:
	//   - []*response.WithdrawResponseMonthStatusFailed: List of failed monthly withdraw statistics.
	//   - *response.ErrorResponse: Error information if any occurred.
	FindMonthWithdrawStatusFailed(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*response.WithdrawResponseMonthStatusFailed, *response.ErrorResponse)

	// FindYearlyWithdrawStatusFailed retrieves yearly failed withdraw statistics.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year to filter the data.
	//
	// Returns:
	//   - []*response.WithdrawResponseYearStatusFailed: List of failed yearly withdraw statistics.
	//   - *response.ErrorResponse: Error information if any occurred.
	FindYearlyWithdrawStatusFailed(ctx context.Context, year int) ([]*response.WithdrawResponseYearStatusFailed, *response.ErrorResponse)
}

type WithdrawStatsAmountService interface {
	// FindMonthlyWithdraws retrieves total amount statistics of monthly withdraws.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year to filter the monthly data.
	//
	// Returns:
	//   - []*response.WithdrawMonthlyAmountResponse: List of total monthly withdraw amounts.
	//   - *response.ErrorResponse: Error information if any occurred.
	FindMonthlyWithdraws(ctx context.Context, year int) ([]*response.WithdrawMonthlyAmountResponse, *response.ErrorResponse)

	// FindYearlyWithdraws retrieves total amount statistics of yearly withdraws.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*response.WithdrawYearlyAmountResponse: List of total yearly withdraw amounts.
	//   - *response.ErrorResponse: Error information if any occurred.
	FindYearlyWithdraws(ctx context.Context, year int) ([]*response.WithdrawYearlyAmountResponse, *response.ErrorResponse)
}
