package transferstatsbycardservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type TransferStatsByCardAmountService interface {

	// FindMonthlyTransferAmountsBySenderCardNumber retrieves monthly transfer amounts by sender card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request with card number, month, and year.
	//
	// Returns:
	//   - []*response.TransferMonthAmountResponse: Monthly transfer amount stats (sent).
	//   - *response.ErrorResponse: Error response if any.
	FindMonthlyTransferAmountsBySenderCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse)

	// FindMonthlyTransferAmountsByReceiverCardNumber retrieves monthly transfer amounts by receiver card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request with card number, month, and year.
	//
	// Returns:
	//   - []*response.TransferMonthAmountResponse: Monthly transfer amount stats (received).
	//   - *response.ErrorResponse: Error response if any.
	FindMonthlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse)

	// FindYearlyTransferAmountsBySenderCardNumber retrieves yearly transfer amounts by sender card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request with card number and year.
	//
	// Returns:
	//   - []*response.TransferYearAmountResponse: Yearly transfer amount stats (sent).
	//   - *response.ErrorResponse: Error response if any.
	FindYearlyTransferAmountsBySenderCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*response.TransferYearAmountResponse, *response.ErrorResponse)

	// FindYearlyTransferAmountsByReceiverCardNumber retrieves yearly transfer amounts by receiver card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request with card number and year.
	//
	// Returns:
	//   - []*response.TransferYearAmountResponse: Yearly transfer amount stats (received).
	//   - *response.ErrorResponse: Error response if any.
	FindYearlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*response.TransferYearAmountResponse, *response.ErrorResponse)
}

type TransferStatsByCardStatusService interface {
	// FindMonthTransferStatusSuccessByCardNumber retrieves monthly successful transfer stats by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing card number, month, and year.
	//
	// Returns:
	//   - []*response.TransferResponseMonthStatusSuccess: Monthly success transfer stats.
	//   - *response.ErrorResponse: Error response if any.
	FindMonthTransferStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*response.TransferResponseMonthStatusSuccess, *response.ErrorResponse)

	// FindYearlyTransferStatusSuccessByCardNumber retrieves yearly successful transfer stats by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing card number and year.
	//
	// Returns:
	//   - []*response.TransferResponseYearStatusSuccess: Yearly success transfer stats.
	//   - *response.ErrorResponse: Error response if any.
	FindYearlyTransferStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*response.TransferResponseYearStatusSuccess, *response.ErrorResponse)

	// FindMonthTransferStatusFailedByCardNumber retrieves monthly failed transfer stats by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing card number, month, and year.
	//
	// Returns:
	//   - []*response.TransferResponseMonthStatusFailed: Monthly failed transfer stats.
	//   - *response.ErrorResponse: Error response if any.
	FindMonthTransferStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*response.TransferResponseMonthStatusFailed, *response.ErrorResponse)

	// FindYearlyTransferStatusFailedByCardNumber retrieves yearly failed transfer stats by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing card number and year.
	//
	// Returns:
	//   - []*response.TransferResponseYearStatusFailed: Yearly failed transfer stats.
	//   - *response.ErrorResponse: Error response if any.
	FindYearlyTransferStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*response.TransferResponseYearStatusFailed, *response.ErrorResponse)
}
