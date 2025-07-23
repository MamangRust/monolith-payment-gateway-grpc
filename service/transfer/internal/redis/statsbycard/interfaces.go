package transferstatsbycardcache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type TransferStatsByCardAmountCache interface {
	// GetMonthlyTransferAmountsBySenderCard retrieves cached monthly transfer amounts for a sender card.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Contains sender card number and year/month filter.
	//
	// Returns:
	//   - []*response.TransferMonthAmountResponse: Monthly transfer amount statistics.
	//   - bool: Whether the cache was found.
	GetMonthlyTransferAmountsBySenderCard(ctx context.Context, req *requests.MonthYearCardNumber) ([]*response.TransferMonthAmountResponse, bool)

	// SetMonthlyTransferAmountsBySenderCard stores monthly transfer amounts for a sender card into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Contains sender card number and year/month.
	//   - data: List of amounts to cache.
	SetMonthlyTransferAmountsBySenderCard(ctx context.Context, req *requests.MonthYearCardNumber, data []*response.TransferMonthAmountResponse)

	// GetMonthlyTransferAmountsByReceiverCard retrieves cached monthly transfer amounts for a receiver card.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Contains receiver card number and year/month filter.
	//
	// Returns:
	//   - []*response.TransferMonthAmountResponse: Monthly transfer amount statistics.
	//   - bool: Whether the cache was found.
	GetMonthlyTransferAmountsByReceiverCard(ctx context.Context, req *requests.MonthYearCardNumber) ([]*response.TransferMonthAmountResponse, bool)

	// SetMonthlyTransferAmountsByReceiverCard stores monthly transfer amounts for a receiver card into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Contains receiver card number and year/month.
	//   - data: List of amounts to cache.
	SetMonthlyTransferAmountsByReceiverCard(ctx context.Context, req *requests.MonthYearCardNumber, data []*response.TransferMonthAmountResponse)

	// GetYearlyTransferAmountsBySenderCard retrieves cached yearly transfer amounts for a sender card.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Contains sender card number and year.
	//
	// Returns:
	//   - []*response.TransferYearAmountResponse: Yearly transfer amount statistics.
	//   - bool: Whether the cache was found.
	GetYearlyTransferAmountsBySenderCard(ctx context.Context, req *requests.MonthYearCardNumber) ([]*response.TransferYearAmountResponse, bool)

	// SetYearlyTransferAmountsBySenderCard stores yearly transfer amounts for a sender card into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Contains sender card number and year.
	//   - data: List of yearly transfer amounts to cache.
	SetYearlyTransferAmountsBySenderCard(ctx context.Context, req *requests.MonthYearCardNumber, data []*response.TransferYearAmountResponse)

	// GetYearlyTransferAmountsByReceiverCard retrieves cached yearly transfer amounts for a receiver card.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Contains receiver card number and year.
	//
	// Returns:
	//   - []*response.TransferYearAmountResponse: Yearly transfer amount statistics.
	//   - bool: Whether the cache was found.
	GetYearlyTransferAmountsByReceiverCard(ctx context.Context, req *requests.MonthYearCardNumber) ([]*response.TransferYearAmountResponse, bool)

	// SetYearlyTransferAmountsByReceiverCard stores yearly transfer amounts for a receiver card into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Contains receiver card number and year.
	//   - data: List of yearly transfer amounts to cache.
	SetYearlyTransferAmountsByReceiverCard(ctx context.Context, req *requests.MonthYearCardNumber, data []*response.TransferYearAmountResponse)
}

type TransferStatsByCardStatusCache interface {
	// GetMonthTransferStatusSuccessByCard retrieves cached monthly successful transfers for a specific card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Contains card number and month.
	//
	// Returns:
	//   - []*response.TransferResponseMonthStatusSuccess: List of monthly successful transfers.
	//   - bool: Whether the cache was found.
	GetMonthTransferStatusSuccessByCard(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*response.TransferResponseMonthStatusSuccess, bool)

	// SetMonthTransferStatusSuccessByCard stores monthly successful transfers for a specific card number into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Contains card number and month.
	//   - data: List of monthly successful transfers to cache.
	SetMonthTransferStatusSuccessByCard(ctx context.Context, req *requests.MonthStatusTransferCardNumber, data []*response.TransferResponseMonthStatusSuccess)

	// GetYearlyTransferStatusSuccessByCard retrieves cached yearly successful transfers for a specific card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Contains card number and year.
	//
	// Returns:
	//   - []*response.TransferResponseYearStatusSuccess: List of yearly successful transfers.
	//   - bool: Whether the cache was found.
	GetYearlyTransferStatusSuccessByCard(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*response.TransferResponseYearStatusSuccess, bool)

	// SetYearlyTransferStatusSuccessByCard stores yearly successful transfers for a specific card number into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Contains card number and year.
	//   - data: List of yearly successful transfers to cache.
	SetYearlyTransferStatusSuccessByCard(ctx context.Context, req *requests.YearStatusTransferCardNumber, data []*response.TransferResponseYearStatusSuccess)

	// GetMonthTransferStatusFailedByCard retrieves cached monthly failed transfers for a specific card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Contains card number and month.
	//
	// Returns:
	//   - []*response.TransferResponseMonthStatusFailed: List of monthly failed transfers.
	//   - bool: Whether the cache was found.
	GetMonthTransferStatusFailedByCard(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*response.TransferResponseMonthStatusFailed, bool)

	// SetMonthTransferStatusFailedByCard stores monthly failed transfers for a specific card number into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Contains card number and month.
	//   - data: List of monthly failed transfers to cache.
	SetMonthTransferStatusFailedByCard(ctx context.Context, req *requests.MonthStatusTransferCardNumber, data []*response.TransferResponseMonthStatusFailed)

	// GetYearlyTransferStatusFailedByCard retrieves cached yearly failed transfers for a specific card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Contains card number and year.
	//
	// Returns:
	//   - []*response.TransferResponseYearStatusFailed: List of yearly failed transfers.
	//   - bool: Whether the cache was found.
	GetYearlyTransferStatusFailedByCard(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*response.TransferResponseYearStatusFailed, bool)

	// SetYearlyTransferStatusFailedByCard stores yearly failed transfers for a specific card number into cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Contains card number and year.
	//   - data: List of yearly failed transfers to cache.
	SetYearlyTransferStatusFailedByCard(ctx context.Context, req *requests.YearStatusTransferCardNumber, data []*response.TransferResponseYearStatusFailed)
}
