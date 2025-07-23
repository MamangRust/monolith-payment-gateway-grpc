package cardstatsmencache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// CardStatsBalanceCache defines the caching behavior for card balance statistics.
type CardStatsBalanceCache interface {
	// GetMonthlyBalanceCache retrieves the global monthly balance statistics
	// (across all cards) for a given year from the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which monthly balance data is requested.
	//
	// Returns:
	//   - []*response.CardResponseMonthBalance: Slice of monthly balance statistics.
	//   - bool: Whether the data was found in the cache.
	GetMonthlyBalanceCache(ctx context.Context, year int) ([]*response.CardResponseMonthBalance, bool)

	// SetMonthlyBalanceCache stores the global monthly balance statistics
	// (across all cards) for a given year in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is being cached.
	//   - data: The data to be cached.
	SetMonthlyBalanceCache(ctx context.Context, year int, data []*response.CardResponseMonthBalance)

	// GetYearlyBalanceCache retrieves the global yearly balance statistics
	// (across all cards) for a given year from the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which yearly balance data is requested.
	//
	// Returns:
	//   - []*response.CardResponseYearlyBalance: Slice of yearly balance statistics.
	//   - bool: Whether the data was found in the cache.
	GetYearlyBalanceCache(ctx context.Context, year int) ([]*response.CardResponseYearlyBalance, bool)

	// SetYearlyBalanceCache stores the global yearly balance statistics
	// (across all cards) for a given year in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is being cached.
	//   - data: The data to be cached.
	SetYearlyBalanceCache(ctx context.Context, year int, data []*response.CardResponseYearlyBalance)
}

// CardStatsTopupCache defines the caching behavior for top-up statistics.
type CardStatsTopupCache interface {

	// GetMonthlyTopupCache retrieves the global monthly top-up statistics
	// (across all cards) for a given year from the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which monthly top-up data is requested.
	//
	// Returns:
	//   - []*response.CardResponseMonthAmount: Slice of monthly top-up statistics.
	//   - bool: Whether the data was found in the cache.
	GetMonthlyTopupCache(ctx context.Context, year int) ([]*response.CardResponseMonthAmount, bool)

	// SetMonthlyTopupCache stores the global monthly top-up statistics
	// (across all cards) for a given year in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is being cached.
	//   - data: The data to be cached.
	SetMonthlyTopupCache(ctx context.Context, year int, data []*response.CardResponseMonthAmount)

	// GetYearlyTopupCache retrieves the global yearly top-up statistics
	// (across all cards) for a given year from the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which yearly top-up data is requested.
	//
	// Returns:
	//   - []*response.CardResponseYearAmount: Slice of yearly top-up statistics.
	//   - bool: Whether the data was found in the cache.
	GetYearlyTopupCache(ctx context.Context, year int) ([]*response.CardResponseYearAmount, bool)

	// SetYearlyTopupCache stores the global yearly top-up statistics
	// (across all cards) for a given year in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is being cached.
	//   - data: The data to be cached.
	SetYearlyTopupCache(ctx context.Context, year int, data []*response.CardResponseYearAmount)
}

// CardStatsWithdrawCache defines the caching behavior for withdraw statistics.
type CardStatsWithdrawCache interface {
	// GetMonthlyWithdrawCache retrieves the global monthly withdraw statistics
	// (across all cards) for a given year from the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which monthly withdraw data is requested.
	//
	// Returns:
	//   - []*response.CardResponseMonthAmount: Slice of monthly withdraw statistics.
	//   - bool: Whether the data was found in the cache.
	GetMonthlyWithdrawCache(ctx context.Context, year int) ([]*response.CardResponseMonthAmount, bool)

	// SetMonthlyWithdrawCache stores the global monthly withdraw statistics
	// (across all cards) for a given year in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is being cached.
	//   - data: The data to be cached.
	SetMonthlyWithdrawCache(ctx context.Context, year int, data []*response.CardResponseMonthAmount)

	// GetYearlyWithdrawCache retrieves the global yearly withdraw statistics
	// (across all cards) for a given year from the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which yearly withdraw data is requested.
	//
	// Returns:
	//   - []*response.CardResponseYearAmount: Slice of yearly withdraw statistics.
	//   - bool: Whether the data was found in the cache.
	GetYearlyWithdrawCache(ctx context.Context, year int) ([]*response.CardResponseYearAmount, bool)

	// SetYearlyWithdrawCache stores the global yearly withdraw statistics
	// (across all cards) for a given year in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is being cached.
	//   - data: The data to be cached.
	SetYearlyWithdrawCache(ctx context.Context, year int, data []*response.CardResponseYearAmount)
}

// CardStatsTransactionCache defines the caching behavior for transaction statistics,
// both globally and per specific card number.
type CardStatsTransactionCache interface {
	// GetMonthlyTransactionCache retrieves the global monthly transaction statistics
	// (across all cards) for a given year from the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which monthly transaction data is requested.
	//
	// Returns:
	//   - []*response.CardResponseMonthAmount: Slice of monthly transaction statistics.
	//   - bool: Whether the data was found in the cache.
	GetMonthlyTransactionCache(ctx context.Context, year int) ([]*response.CardResponseMonthAmount, bool)

	// SetMonthlyTransactionCache stores the global monthly transaction statistics
	// (across all cards) for a given year in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is being cached.
	//   - data: The data to be cached.
	SetMonthlyTransactionCache(ctx context.Context, year int, data []*response.CardResponseMonthAmount)

	// GetYearlyTransactionCache retrieves the global yearly transaction statistics
	// (across all cards) for a given year from the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which yearly transaction data is requested.
	//
	// Returns:
	//   - []*response.CardResponseYearAmount: Slice of yearly transaction statistics.
	//   - bool: Whether the data was found in the cache.
	GetYearlyTransactionCache(ctx context.Context, year int) ([]*response.CardResponseYearAmount, bool)

	// SetYearlyTransactionCache stores the global yearly transaction statistics
	// (across all cards) for a given year in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is being cached.
	//   - data: The data to be cached.
	SetYearlyTransactionCache(ctx context.Context, year int, data []*response.CardResponseYearAmount)
}

// CardStatsTransferCache defines the caching behavior for transfer statistics,
// including both transfer-out (sender) and transfer-in (receiver) data.
type CardStatsTransferCache interface {
	// GetMonthlyTransferSenderCache retrieves the global monthly transfer-out statistics
	// (across all sender card numbers) for a given year from the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which monthly transfer-out data is requested.
	//
	// Returns:
	//   - []*response.CardResponseMonthAmount: Slice of monthly transfer-out statistics.
	//   - bool: Whether the data was found in the cache.
	GetMonthlyTransferSenderCache(ctx context.Context, year int) ([]*response.CardResponseMonthAmount, bool)

	// SetMonthlyTransferSenderCache stores the global monthly transfer-out statistics
	// (across all sender card numbers) for a given year in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is being cached.
	//   - data: The data to be cached.
	SetMonthlyTransferSenderCache(ctx context.Context, year int, data []*response.CardResponseMonthAmount)

	// GetYearlyTransferSenderCache retrieves the global yearly transfer-out statistics
	// (across all sender card numbers) for a given year from the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which yearly transfer-out data is requested.
	//
	// Returns:
	//   - []*response.CardResponseYearAmount: Slice of yearly transfer-out statistics.
	//   - bool: Whether the data was found in the cache.
	GetYearlyTransferSenderCache(ctx context.Context, year int) ([]*response.CardResponseYearAmount, bool)

	// SetYearlyTransferSenderCache stores the global yearly transfer-out statistics
	// (across all sender card numbers) for a given year in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is being cached.
	//   - data: The data to be cached.
	SetYearlyTransferSenderCache(ctx context.Context, year int, data []*response.CardResponseYearAmount)

	// GetMonthlyTransferReceiverCache retrieves the global monthly transfer-in statistics
	// (across all receiver card numbers) for a given year from the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which monthly transfer-in data is requested.
	//
	// Returns:
	//   - []*response.CardResponseMonthAmount: Slice of monthly transfer-in statistics.
	//   - bool: Whether the data was found in the cache.
	GetMonthlyTransferReceiverCache(ctx context.Context, year int) ([]*response.CardResponseMonthAmount, bool)

	// SetMonthlyTransferReceiverCache stores the global monthly transfer-in statistics
	// (across all receiver card numbers) for a given year in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is being cached.
	//   - data: The data to be cached.
	SetMonthlyTransferReceiverCache(ctx context.Context, year int, data []*response.CardResponseMonthAmount)

	// GetYearlyTransferReceiverCache retrieves the global yearly transfer-in statistics
	// (across all receiver card numbers) for a given year from the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which yearly transfer-in data is requested.
	//
	// Returns:
	//   - []*response.CardResponseYearAmount: Slice of yearly transfer-in statistics.
	//   - bool: Whether the data was found in the cache.
	GetYearlyTransferReceiverCache(ctx context.Context, year int) ([]*response.CardResponseYearAmount, bool)

	// SetYearlyTransferReceiverCache stores the global yearly transfer-in statistics
	// (across all receiver card numbers) for a given year in the cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is being cached.
	//   - data: The data to be cached.
	SetYearlyTransferReceiverCache(ctx context.Context, year int, data []*response.CardResponseYearAmount)
}
