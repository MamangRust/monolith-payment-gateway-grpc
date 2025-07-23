package transferstatsbycardrepository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type TransferStatsByCardAmountSenderRepository interface {
	// GetMonthlyTransferAmountsBySenderCardNumber retrieves monthly transfer amounts where the card is the sender.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing the card number and year.
	//
	// Returns:
	//   - []*record.TransferMonthAmount: List of monthly transfer amount records.
	//   - error: Any error encountered during the operation.
	GetMonthlyTransferAmountsBySenderCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*record.TransferMonthAmount, error)

	// GetYearlyTransferAmountsBySenderCardNumber retrieves yearly transfer amounts where the card is the sender.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing the card number and year.
	//
	// Returns:
	//   - []*record.TransferYearAmount: List of yearly transfer amount records.
	//   - error: Any error encountered during the operation.
	GetYearlyTransferAmountsBySenderCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*record.TransferYearAmount, error)
}

type TransferStatsByCardAmountReceiverRepository interface {

	// GetMonthlyTransferAmountsByReceiverCardNumber retrieves monthly transfer amounts where the card is the receiver.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing the card number and year.
	//
	// Returns:
	//   - []*record.TransferMonthAmount: List of monthly transfer amount records.
	//   - error: Any error encountered during the operation.
	GetMonthlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*record.TransferMonthAmount, error)

	// GetYearlyTransferAmountsByReceiverCardNumber retrieves yearly transfer amounts where the card is the receiver.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing the card number and year.
	//
	// Returns:
	//   - []*record.TransferYearAmount: List of yearly transfer amount records.
	//   - error: Any error encountered during the operation.
	GetYearlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*record.TransferYearAmount, error)
}

type TransferStatsByCardStatusRepository interface {
	// GetMonthTransferStatusSuccessByCardNumber retrieves monthly successful transfer statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing the card number and date filters.
	//
	// Returns:
	//   - []*record.TransferRecordMonthStatusSuccess: List of monthly successful transfer statistics.
	//   - error: Any error encountered during the operation.
	GetMonthTransferStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*record.TransferRecordMonthStatusSuccess, error)

	// GetYearlyTransferStatusSuccessByCardNumber retrieves yearly successful transfer statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing the card number and year.
	//
	// Returns:
	//   - []*record.TransferRecordYearStatusSuccess: List of yearly successful transfer statistics.
	//   - error: Any error encountered during the operation.
	GetYearlyTransferStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*record.TransferRecordYearStatusSuccess, error)

	// GetMonthTransferStatusFailedByCardNumber retrieves monthly failed transfer statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing the card number and date filters.
	//
	// Returns:
	//   - []*record.TransferRecordMonthStatusFailed: List of monthly failed transfer statistics.
	//   - error: Any error encountered during the operation.
	GetMonthTransferStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*record.TransferRecordMonthStatusFailed, error)

	// GetYearlyTransferStatusFailedByCardNumber retrieves yearly failed transfer statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing the card number and year.
	//
	// Returns:
	//   - []*record.TransferRecordYearStatusFailed: List of yearly failed transfer statistics.
	//   - error: Any error encountered during the operation.
	GetYearlyTransferStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*record.TransferRecordYearStatusFailed, error)
}
