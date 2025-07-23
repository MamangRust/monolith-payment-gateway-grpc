package withdrawstatsbycardservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type WithdrawStatsByCardStatusService interface {
	// FindMonthWithdrawStatusSuccessByCardNumber retrieves monthly successful withdraw statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing the card number, month, and year.
	//
	// Returns:
	//   - []*response.WithdrawResponseMonthStatusSuccess: List of successful monthly withdraw statistics for the given card number.
	//   - *response.ErrorResponse: Error information if any occurred.
	FindMonthWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*response.WithdrawResponseMonthStatusSuccess, *response.ErrorResponse)

	// FindYearlyWithdrawStatusSuccessByCardNumber retrieves yearly successful withdraw statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing the card number and year.
	//
	// Returns:
	//   - []*response.WithdrawResponseYearStatusSuccess: List of successful yearly withdraw statistics for the given card number.
	//   - *response.ErrorResponse: Error information if any occurred.
	FindYearlyWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*response.WithdrawResponseYearStatusSuccess, *response.ErrorResponse)

	// FindMonthWithdrawStatusFailedByCardNumber retrieves monthly failed withdraw statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing the card number, month, and year.
	//
	// Returns:
	//   - []*response.WithdrawResponseMonthStatusFailed: List of failed monthly withdraw statistics for the given card number.
	//   - *response.ErrorResponse: Error information if any occurred.
	FindMonthWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*response.WithdrawResponseMonthStatusFailed, *response.ErrorResponse)

	// FindYearlyWithdrawStatusFailedByCardNumber retrieves yearly failed withdraw statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing the card number and year.
	//
	// Returns:
	//   - []*response.WithdrawResponseYearStatusFailed: List of failed yearly withdraw statistics for the given card number.
	//   - *response.ErrorResponse: Error information if any occurred.
	FindYearlyWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*response.WithdrawResponseYearStatusFailed, *response.ErrorResponse)
}

type WithdrawStatsByCardAmountService interface {
	// FindMonthlyWithdrawsByCardNumber retrieves total monthly withdraw amounts by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing the card number, month, and year.
	//
	// Returns:
	//   - []*response.WithdrawMonthlyAmountResponse: List of monthly withdraw amounts for the given card number.
	//   - *response.ErrorResponse: Error information if any occurred.
	FindMonthlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*response.WithdrawMonthlyAmountResponse, *response.ErrorResponse)

	// FindYearlyWithdrawsByCardNumber retrieves total yearly withdraw amounts by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing the card number and year.
	//
	// Returns:
	//   - []*response.WithdrawYearlyAmountResponse: List of yearly withdraw amounts for the given card number.
	//   - *response.ErrorResponse: Error information if any occurred.
	FindYearlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*response.WithdrawYearlyAmountResponse, *response.ErrorResponse)
}
