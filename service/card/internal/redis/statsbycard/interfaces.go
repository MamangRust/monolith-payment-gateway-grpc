package cardstatsbycardmencache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// CardStatsBalanceByCardCache defines the caching behavior for card balance statistics
// specific to an individual card number.
type CardStatsBalanceByCardCache interface {
	// GetMonthlyBalanceByNumberCache retrieves the cached monthly balance statistics
	// for a specific card number based on the given month and year.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: A request object containing the card number, month, and year.
	//
	// Returns:
	//   - []*response.CardResponseMonthBalance: Slice of monthly balance statistics for the specified card.
	//   - bool: Whether the data was found in the cache.
	GetMonthlyBalanceByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthBalance, bool)

	// GetYearlyBalanceByNumberCache retrieves the cached yearly balance statistics
	// for a specific card number based on the given year.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: A request object containing the card number and year.
	//
	// Returns:
	//   - []*response.CardResponseYearlyBalance: Slice of yearly balance statistics for the specified card.
	//   - bool: Whether the data was found in the cache.
	GetYearlyBalanceByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearlyBalance, bool)

	// SetMonthlyBalanceByNumberCache stores the monthly balance statistics
	// for a specific card number in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: A request object containing the card number, month, and year.
	//   - data: The data to be cached.
	SetMonthlyBalanceByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*response.CardResponseMonthBalance)

	// SetYearlyBalanceByNumberCache stores the yearly balance statistics
	// for a specific card number in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: A request object containing the card number and year.
	//   - data: The data to be cached.
	SetYearlyBalanceByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*response.CardResponseYearlyBalance)
}

// CardStatsTopupByCardCache defines the caching behavior for top-up statistics
// specific to an individual card number.
type CardStatsTopupByCardCache interface {
	// GetMonthlyTopupByNumberCache retrieves the cached monthly top-up statistics
	// for a specific card number based on the given month and year.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: A request object containing the card number, month, and year.
	//
	// Returns:
	//   - []*response.CardResponseMonthAmount: Slice of monthly top-up statistics for the specified card.
	//   - bool: Whether the data was found in the cache.
	GetMonthlyTopupByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, bool)

	// GetYearlyTopupByNumberCache retrieves the cached yearly top-up statistics
	// for a specific card number based on the given year.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: A request object containing the card number and year.
	//
	// Returns:
	//   - []*response.CardResponseYearAmount: Slice of yearly top-up statistics for the specified card.
	//   - bool: Whether the data was found in the cache.
	GetYearlyTopupByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, bool)

	// SetMonthlyTopupByNumberCache stores the monthly top-up statistics
	// for a specific card number in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: A request object containing the card number, month, and year.
	//   - data: The data to be cached.
	SetMonthlyTopupByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*response.CardResponseMonthAmount)

	// SetYearlyTopupByNumberCache stores the yearly top-up statistics
	// for a specific card number in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: A request object containing the card number and year.
	//   - data: The data to be cached.
	SetYearlyTopupByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*response.CardResponseYearAmount)
}

// CardStatsWithdrawByCardCache defines the caching behavior for withdraw statistics
// specific to an individual card number.
type CardStatsWithdrawByCardCache interface {
	// GetMonthlyWithdrawByNumberCache retrieves the cached monthly withdraw statistics
	// for a specific card number based on the given month and year.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: A request object containing the card number, month, and year.
	//
	// Returns:
	//   - []*response.CardResponseMonthAmount: Slice of monthly withdraw statistics for the specified card.
	//   - bool: Whether the data was found in the cache.
	GetMonthlyWithdrawByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, bool)

	// GetYearlyWithdrawByNumberCache retrieves the cached yearly withdraw statistics
	// for a specific card number based on the given year.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: A request object containing the card number and year.
	//
	// Returns:
	//   - []*response.CardResponseYearAmount: Slice of yearly withdraw statistics for the specified card.
	//   - bool: Whether the data was found in the cache.
	GetYearlyWithdrawByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, bool)

	// SetMonthlyWithdrawByNumberCache stores the monthly withdraw statistics
	// for a specific card number in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: A request object containing the card number, month, and year.
	//   - data: The data to be cached.
	SetMonthlyWithdrawByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*response.CardResponseMonthAmount)

	// SetYearlyWithdrawByNumberCache stores the yearly withdraw statistics
	// for a specific card number in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: A request object containing the card number and year.
	//   - data: The data to be cached.
	SetYearlyWithdrawByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*response.CardResponseYearAmount)
}

// CardStatsTransactionByCardCache defines the caching behavior for transaction statistics
// specific to an individual card number.
type CardStatsTransactionByCardCache interface {
	// GetMonthlyTransactionByNumberCache retrieves the cached monthly transaction statistics
	// for a specific card number based on the given month and year.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: A request object containing the card number, month, and year.
	//
	// Returns:
	//   - []*response.CardResponseMonthAmount: Slice of monthly transaction statistics for the specified card.
	//   - bool: Whether the data was found in the cache.
	GetMonthlyTransactionByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, bool)

	// GetYearlyTransactionByNumberCache retrieves the cached yearly transaction statistics
	// for a specific card number based on the given year.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: A request object containing the card number and year.
	//
	// Returns:
	//   - []*response.CardResponseYearAmount: Slice of yearly transaction statistics for the specified card.
	//   - bool: Whether the data was found in the cache.
	GetYearlyTransactionByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, bool)

	// SetMonthlyTransactionByNumberCache stores the monthly transaction statistics
	// for a specific card number in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: A request object containing the card number, month, and year.
	//   - data: The data to be cached.
	SetMonthlyTransactionByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*response.CardResponseMonthAmount)

	// SetYearlyTransactionByNumberCache stores the yearly transaction statistics
	// for a specific card number in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: A request object containing the card number and year.
	//   - data: The data to be cached.
	SetYearlyTransactionByNumberCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*response.CardResponseYearAmount)
}

// CardStatsTransferByCardCache defines the caching behavior for transfer statistics
// specific to individual sender and receiver card numbers.
type CardStatsTransferByCardCache interface {
	// GetMonthlyTransferBySenderCache retrieves the cached monthly transfer-out statistics
	// for a specific sender card number based on the given month and year.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: A request object containing the sender card number, month, and year.
	//
	// Returns:
	//   - []*response.CardResponseMonthAmount: Slice of monthly transfer-out statistics for the sender card.
	//   - bool: Whether the data was found in the cache.
	GetMonthlyTransferBySenderCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, bool)

	// GetYearlyTransferBySenderCache retrieves the cached yearly transfer-out statistics
	// for a specific sender card number based on the given year.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: A request object containing the sender card number and year.
	//
	// Returns:
	//   - []*response.CardResponseYearAmount: Slice of yearly transfer-out statistics for the sender card.
	//   - bool: Whether the data was found in the cache.
	GetYearlyTransferBySenderCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, bool)

	// SetMonthlyTransferBySenderCache stores the monthly transfer-out statistics
	// for a specific sender card number in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: A request object containing the sender card number, month, and year.
	//   - data: The data to be cached.
	SetMonthlyTransferBySenderCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*response.CardResponseMonthAmount)

	// SetYearlyTransferBySenderCache stores the yearly transfer-out statistics
	// for a specific sender card number in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: A request object containing the sender card number and year.
	//   - data: The data to be cached.
	SetYearlyTransferBySenderCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*response.CardResponseYearAmount)

	// GetMonthlyTransferByReceiverCache retrieves the cached monthly transfer-in statistics
	// for a specific receiver card number based on the given month and year.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: A request object containing the receiver card number, month, and year.
	//
	// Returns:
	//   - []*response.CardResponseMonthAmount: Slice of monthly transfer-in statistics for the receiver card.
	//   - bool: Whether the data was found in the cache.
	GetMonthlyTransferByReceiverCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, bool)

	// GetYearlyTransferByReceiverCache retrieves the cached yearly transfer-in statistics
	// for a specific receiver card number based on the given year.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: A request object containing the receiver card number and year.
	//
	// Returns:
	//   - []*response.CardResponseYearAmount: Slice of yearly transfer-in statistics for the receiver card.
	//   - bool: Whether the data was found in the cache.
	GetYearlyTransferByReceiverCache(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, bool)

	// SetMonthlyTransferByReceiverCache stores the monthly transfer-in statistics
	// for a specific receiver card number in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: A request object containing the receiver card number, month, and year.
	//   - data: The data to be cached.
	SetMonthlyTransferByReceiverCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*response.CardResponseMonthAmount)

	// SetYearlyTransferByReceiverCache stores the yearly transfer-in statistics
	// for a specific receiver card number in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: A request object containing the receiver card number and year.
	//   - data: The data to be cached.
	SetYearlyTransferByReceiverCache(ctx context.Context, req *requests.MonthYearCardNumberCard, data []*response.CardResponseYearAmount)
}
