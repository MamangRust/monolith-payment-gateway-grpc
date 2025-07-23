package withdrawstatsbycardrepository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type WithdrawStatsByCardStatusRepository interface {
	// GetMonthWithdrawStatusSuccessByCardNumber retrieves monthly withdraw success statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request containing card number, month, and year.
	//
	// Returns:
	//   - []*record.WithdrawRecordMonthStatusSuccess: List of successful monthly withdraw records.
	//   - error: An error if the operation fails.
	GetMonthWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*record.WithdrawRecordMonthStatusSuccess, error)

	// GetYearlyWithdrawStatusSuccessByCardNumber retrieves yearly withdraw success statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request containing card number and year.
	//
	// Returns:
	//   - []*record.WithdrawRecordYearStatusSuccess: List of successful yearly withdraw records.
	//   - error: An error if the operation fails.
	GetYearlyWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*record.WithdrawRecordYearStatusSuccess, error)

	// GetMonthWithdrawStatusFailedByCardNumber retrieves monthly withdraw failed statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request containing card number, month, and year.
	//
	// Returns:
	//   - []*record.WithdrawRecordMonthStatusFailed: List of failed monthly withdraw records.
	//   - error: An error if the operation fails.
	GetMonthWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*record.WithdrawRecordMonthStatusFailed, error)

	// GetYearlyWithdrawStatusFailedByCardNumber retrieves yearly withdraw failed statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request containing card number and year.
	//
	// Returns:
	//   - []*record.WithdrawRecordYearStatusFailed: List of failed yearly withdraw records.
	//   - error: An error if the operation fails.
	GetYearlyWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*record.WithdrawRecordYearStatusFailed, error)
}

type WithdrawStatsByCardAmountRepository interface {
	// GetMonthlyWithdrawsByCardNumber retrieves total monthly withdraw amounts by card number for a specific year.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request containing card number and year.
	//
	// Returns:
	//   - []*record.WithdrawMonthlyAmount: List of monthly withdraw amount records.
	//   - error: An error if the operation fails.
	GetMonthlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*record.WithdrawMonthlyAmount, error)

	// GetYearlyWithdrawsByCardNumber retrieves total yearly withdraw amounts by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Request containing card number and year.
	//
	// Returns:
	//   - []*record.WithdrawYearlyAmount: List of yearly withdraw amount records.
	//   - error: An error if the operation fails.
	GetYearlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*record.WithdrawYearlyAmount, error)
}
