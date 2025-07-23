package repositorystatsbycard

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

type CardStatsBalanceByCardRepository interface {
	// GetMonthlyBalancesByCardNumber retrieves total balances grouped by month
	// for a specific card number in the specified year.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - req: the request containing month, year, and card number
	//
	// Returns:
	//   - A slice of CardMonthBalance records or an error if the operation fails.
	GetMonthlyBalancesByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*record.CardMonthBalance, error)

	// GetYearlyBalanceByCardNumber retrieves total balances grouped by year
	// for a specific card number.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - req: the request containing year and card number
	//
	// Returns:
	//   - A slice of CardYearlyBalance records or an error if the operation fails.
	GetYearlyBalanceByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*record.CardYearlyBalance, error)
}

type CardStatsTopupByCardRepository interface {
	// GetMonthlyTopupAmountByCardNumber retrieves total top-up amounts grouped by month
	// for a specific card number in the specified year.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - req: the request containing month, year, and card number
	//
	// Returns:
	//   - A slice of CardMonthAmount records or an error if the operation fails.
	GetMonthlyTopupAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*record.CardMonthAmount, error)

	// GetYearlyTopupAmountByCardNumber retrieves total top-up amounts grouped by year
	// for a specific card number.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - req: the request containing year and card number
	//
	// Returns:
	//   - A slice of CardYearAmount records or an error if the operation fails.
	GetYearlyTopupAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*record.CardYearAmount, error)
}

type CardStatsWithdrawByCardRepository interface {
	// GetMonthlyWithdrawAmountByCardNumber retrieves total withdraw amounts grouped by month
	// for a specific card number in the specified year.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - req: the request containing month, year, and card number
	//
	// Returns:
	//   - A slice of CardMonthAmount records or an error if the operation fails.
	GetMonthlyWithdrawAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*record.CardMonthAmount, error)

	// GetYearlyWithdrawAmountByCardNumber retrieves total withdraw amounts grouped by year
	// for a specific card number.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - req: the request containing year and card number
	//
	// Returns:
	//   - A slice of CardYearAmount records or an error if the operation fails.
	GetYearlyWithdrawAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*record.CardYearAmount, error)
}

type CardStatsTransactionByCardRepository interface {
	// GetMonthlyTransactionAmountByCardNumber retrieves total transaction amounts grouped by month
	// for a specific card number in the given year.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - req: the request containing card number, month, and year
	//
	// Returns:
	//   - A slice of CardMonthAmount records or an error if the operation fails.
	GetMonthlyTransactionAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*record.CardMonthAmount, error)

	// GetYearlyTransactionAmountByCardNumber retrieves total transaction amounts grouped by year
	// for a specific card number.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - req: the request containing card number and year
	//
	// Returns:
	//   - A slice of CardYearAmount records or an error if the operation fails.
	GetYearlyTransactionAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*record.CardYearAmount, error)
}

type CardStatsTransferByCardRepository interface {
	// GetMonthlyTransferAmountBySender retrieves total transfer amounts grouped by month
	// for a specific sender card number in a given year.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - req: request containing month, year, and sender card number
	//
	// Returns:
	//   - A slice of CardMonthAmount records or an error if the operation fails.
	GetMonthlyTransferAmountBySender(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*record.CardMonthAmount, error)

	// GetYearlyTransferAmountBySender retrieves total transfer amounts grouped by year
	// for a specific sender card number.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - req: request containing year and sender card number
	//
	// Returns:
	//   - A slice of CardYearAmount records or an error if the operation fails.
	GetYearlyTransferAmountBySender(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*record.CardYearAmount, error)

	// GetMonthlyTransferAmountByReceiver retrieves total transfer amounts grouped by month
	// for a specific receiver card number in a given year.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - req: request containing month, year, and receiver card number
	//
	// Returns:
	//   - A slice of CardMonthAmount records or an error if the operation fails.
	GetMonthlyTransferAmountByReceiver(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*record.CardMonthAmount, error)

	// GetYearlyTransferAmountByReceiver retrieves total transfer amounts grouped by year
	// for a specific receiver card number.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - req: request containing year and receiver card number
	//
	// Returns:
	//   - A slice of CardYearAmount records or an error if the operation fails.
	GetYearlyTransferAmountByReceiver(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*record.CardYearAmount, error)
}
