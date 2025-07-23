package topupstatsservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type TopupStatsAmountService interface {
	// FindMonthlyTopupAmounts retrieves monthly statistics of topup amounts.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the statistics are requested.
	//
	// Returns:
	//   - []*response.TopupMonthAmountResponse: List of monthly topup amounts.
	//   - *response.ErrorResponse: Error details if retrieval fails.
	FindMonthlyTopupAmounts(ctx context.Context, year int) ([]*response.TopupMonthAmountResponse, *response.ErrorResponse)

	// FindYearlyTopupAmounts retrieves yearly statistics of topup amounts.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the statistics are requested.
	//
	// Returns:
	//   - []*response.TopupYearlyAmountResponse: List of yearly topup amounts.
	//   - *response.ErrorResponse: Error details if retrieval fails.
	FindYearlyTopupAmounts(ctx context.Context, year int) ([]*response.TopupYearlyAmountResponse, *response.ErrorResponse)
}

type TopupStatsMethodService interface {
	// FindMonthlyTopupMethods retrieves monthly statistics grouped by topup methods.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the statistics are requested.
	//
	// Returns:
	//   - []*response.TopupMonthMethodResponse: List of monthly method usage.
	//   - *response.ErrorResponse: Error details if retrieval fails.
	FindMonthlyTopupMethods(ctx context.Context, year int) ([]*response.TopupMonthMethodResponse, *response.ErrorResponse)

	// FindYearlyTopupMethods retrieves yearly statistics grouped by topup methods.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the statistics are requested.
	//
	// Returns:
	//   - []*response.TopupYearlyMethodResponse: List of yearly method usage.
	//   - *response.ErrorResponse: Error details if retrieval fails.
	FindYearlyTopupMethods(ctx context.Context, year int) ([]*response.TopupYearlyMethodResponse, *response.ErrorResponse)
}

type TopupStatsStatusService interface {
	// FindMonthTopupStatusSuccess retrieves monthly statistics of successful topups.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing month and year filters.
	//
	// Returns:
	//   - []*response.TopupResponseMonthStatusSuccess: List of monthly success statistics.
	//   - *response.ErrorResponse: Error details if retrieval fails.
	FindMonthTopupStatusSuccess(ctx context.Context, req *requests.MonthTopupStatus) ([]*response.TopupResponseMonthStatusSuccess, *response.ErrorResponse)

	// FindYearlyTopupStatusSuccess retrieves yearly statistics of successful topups.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the statistics are requested.
	//
	// Returns:
	//   - []*response.TopupResponseYearStatusSuccess: List of yearly success statistics.
	//   - *response.ErrorResponse: Error details if retrieval fails.
	FindYearlyTopupStatusSuccess(ctx context.Context, year int) ([]*response.TopupResponseYearStatusSuccess, *response.ErrorResponse)

	// FindMonthTopupStatusFailed retrieves monthly statistics of failed topups.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing month and year filters.
	//
	// Returns:
	//   - []*response.TopupResponseMonthStatusFailed: List of monthly failed statistics.
	//   - *response.ErrorResponse: Error details if retrieval fails.
	FindMonthTopupStatusFailed(ctx context.Context, req *requests.MonthTopupStatus) ([]*response.TopupResponseMonthStatusFailed, *response.ErrorResponse)

	// FindYearlyTopupStatusFailed retrieves yearly statistics of failed topups.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which the statistics are requested.
	//
	// Returns:
	//   - []*response.TopupResponseYearStatusFailed: List of yearly failed statistics.
	//   - *response.ErrorResponse: Error details if retrieval fails.
	FindYearlyTopupStatusFailed(ctx context.Context, year int) ([]*response.TopupResponseYearStatusFailed, *response.ErrorResponse)
}
