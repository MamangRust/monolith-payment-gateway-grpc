package repositorystats

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
)

// CardStatsBalanceRepository handles balance statistics across all cards,
// supporting both global aggregation and per-card statistics.
type CardStatsBalanceRepository interface {
	// GetMonthlyBalance retrieves total balances grouped by month
	// for all cards in the specified year.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - year: the integer year for which to retrieve the monthly balance statistics
	//
	// Returns:
	//   - A slice of CardMonthBalance records or an error if the operation fails.
	GetMonthlyBalance(ctx context.Context, year int) ([]*record.CardMonthBalance, error)

	// GetYearlyBalance retrieves total balances grouped by year
	// for all cards in the specified year.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - year: the integer year for which to retrieve the yearly balance statistics
	//
	// Returns:
	//   - A slice of CardYearlyBalance records or an error if the operation fails.
	GetYearlyBalance(ctx context.Context, year int) ([]*record.CardYearlyBalance, error)
}

// CardStatsTopupRepository handles top-up statistics across all cards,
// supporting both global aggregation and per-card statistics.
type CardStatsTopupRepository interface {
	// GetMonthlyTopupAmount retrieves total top-up amounts grouped by month
	// for all cards in the specified year.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - year: the integer year for which to retrieve the monthly statistics
	//
	// Returns:
	//   - A slice of CardMonthAmount records or an error if the operation fails.
	GetMonthlyTopupAmount(ctx context.Context, year int) ([]*record.CardMonthAmount, error)

	// GetYearlyTopupAmount retrieves total top-up amounts grouped by year
	// for all cards in the specified year.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - year: the integer year for which to retrieve the yearly statistics
	//
	// Returns:
	//   - A slice of CardYearAmount records or an error if the operation fails.
	GetYearlyTopupAmount(ctx context.Context, year int) ([]*record.CardYearAmount, error)
}

// CardStatsWithdrawRepository handles withdraw statistics across all cards,
// supporting both global aggregation and per-card statistics.
type CardStatsWithdrawRepository interface {
	// GetMonthlyWithdrawAmount retrieves total withdraw amounts grouped by month
	// for all cards in the specified year.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - year: the integer year for which to retrieve the monthly statistics
	//
	// Returns:
	//   - A slice of CardMonthAmount records or an error if the operation fails.
	GetMonthlyWithdrawAmount(ctx context.Context, year int) ([]*record.CardMonthAmount, error)

	// GetYearlyWithdrawAmount retrieves total withdraw amounts grouped by year
	// for all cards in the specified year.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - year: the integer year for which to retrieve the yearly statistics
	//
	// Returns:
	//   - A slice of CardYearAmount records or an error if the operation fails.
	GetYearlyWithdrawAmount(ctx context.Context, year int) ([]*record.CardYearAmount, error)
}

// CardStatsTransactionRepository handles transaction statistics across all cards (global and per card number).
type CardStatsTransactionRepository interface {
	// GetMonthlyTransactionAmount retrieves total transaction amounts grouped by month
	// for all cards in the given year.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - year: the integer year for which to retrieve the monthly statistics
	//
	// Returns:
	//   - A slice of CardMonthAmount records or an error if the operation fails.
	GetMonthlyTransactionAmount(ctx context.Context, year int) ([]*record.CardMonthAmount, error)

	// GetYearlyTransactionAmount retrieves total transaction amounts grouped by year
	// for all cards.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - year: the target year
	//
	// Returns:
	//   - A slice of CardYearAmount records or an error if the operation fails.
	GetYearlyTransactionAmount(ctx context.Context, year int) ([]*record.CardYearAmount, error)
}

// CardStatsTransferRepository handles transfer statistics across all cards,
// including sender and receiver perspectives, both globally and per card number.
type CardStatsTransferRepository interface {
	// GetMonthlyTransferAmountSender retrieves total transfer amounts grouped by month
	// for all cards acting as sender in the given year.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - year: the integer year for which to retrieve the monthly statistics
	//
	// Returns:
	//   - A slice of CardMonthAmount records or an error if the operation fails.
	GetMonthlyTransferAmountSender(ctx context.Context, year int) ([]*record.CardMonthAmount, error)

	// GetYearlyTransferAmountSender retrieves total transfer amounts grouped by year
	// for all cards acting as sender.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - year: the integer year for which to retrieve the statistics
	//
	// Returns:
	//   - A slice of CardYearAmount records or an error if the operation fails.
	GetYearlyTransferAmountSender(ctx context.Context, year int) ([]*record.CardYearAmount, error)

	// GetMonthlyTransferAmountReceiver retrieves total transfer amounts grouped by month
	// for all cards acting as receiver in the given year.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - year: the integer year for which to retrieve the monthly statistics
	//
	// Returns:
	//   - A slice of CardMonthAmount records or an error if the operation fails.
	GetMonthlyTransferAmountReceiver(ctx context.Context, year int) ([]*record.CardMonthAmount, error)

	// GetYearlyTransferAmountReceiver retrieves total transfer amounts grouped by year
	// for all cards acting as receiver.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - year: the integer year for which to retrieve the statistics
	//
	// Returns:
	//   - A slice of CardYearAmount records or an error if the operation fails.
	GetYearlyTransferAmountReceiver(ctx context.Context, year int) ([]*record.CardYearAmount, error)
}
