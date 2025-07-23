package transferstatsservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type TransferStatsStatusService interface {
	// FindMonthTransferStatusSuccess retrieves monthly successful transfer statistics.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing month and year filters.
	//
	// Returns:
	//   - []*response.TransferResponseMonthStatusSuccess: List of monthly success transfer statistics.
	//   - *response.ErrorResponse: Error response if an error occurs.
	FindMonthTransferStatusSuccess(ctx context.Context, req *requests.MonthStatusTransfer) ([]*response.TransferResponseMonthStatusSuccess, *response.ErrorResponse)

	// FindYearlyTransferStatusSuccess retrieves yearly successful transfer statistics.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*response.TransferResponseYearStatusSuccess: List of yearly success transfer statistics.
	//   - *response.ErrorResponse: Error response if an error occurs.
	FindYearlyTransferStatusSuccess(ctx context.Context, year int) ([]*response.TransferResponseYearStatusSuccess, *response.ErrorResponse)

	// FindMonthTransferStatusFailed retrieves monthly failed transfer statistics.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing month and year filters.
	//
	// Returns:
	//   - []*response.TransferResponseMonthStatusFailed: List of monthly failed transfer statistics.
	//   - *response.ErrorResponse: Error response if an error occurs.
	FindMonthTransferStatusFailed(ctx context.Context, req *requests.MonthStatusTransfer) ([]*response.TransferResponseMonthStatusFailed, *response.ErrorResponse)

	// FindYearlyTransferStatusFailed retrieves yearly failed transfer statistics.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*response.TransferResponseYearStatusFailed: List of yearly failed transfer statistics.
	//   - *response.ErrorResponse: Error response if an error occurs.
	FindYearlyTransferStatusFailed(ctx context.Context, year int) ([]*response.TransferResponseYearStatusFailed, *response.ErrorResponse)
}

type TransferStatsAmountService interface {
	// FindMonthlyTransferAmounts retrieves monthly total transfer amounts.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which monthly data is requested.
	//
	// Returns:
	//   - []*response.TransferMonthAmountResponse: List of monthly transfer amount statistics.
	//   - *response.ErrorResponse: Error response if an error occurs.
	FindMonthlyTransferAmounts(ctx context.Context, year int) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse)

	// FindYearlyTransferAmounts retrieves yearly total transfer amounts.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*response.TransferYearAmountResponse: List of yearly transfer amount statistics.
	//   - *response.ErrorResponse: Error response if an error occurs.
	FindYearlyTransferAmounts(ctx context.Context, year int) ([]*response.TransferYearAmountResponse, *response.ErrorResponse)
}
