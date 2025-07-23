package cardstatsservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// CardStatsBalanceService handles balance statistics globally and per specific card number.
type CardStatsBalanceService interface {
	// FindMonthlyBalance retrieves monthly balance statistics across all card numbers for a given year.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - year: the year for which the monthly balances are requested
	//
	// Returns:
	//   - A slice of CardResponseMonthBalance or an error response if the operation fails.
	FindMonthlyBalance(ctx context.Context, year int) ([]*response.CardResponseMonthBalance, *response.ErrorResponse)

	// FindYearlyBalance retrieves yearly balance statistics across all card numbers for a given year.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - year: the year for which the yearly balances are requested
	//
	// Returns:
	//   - A slice of CardResponseYearlyBalance or an error response if the operation fails.
	FindYearlyBalance(ctx context.Context, year int) ([]*response.CardResponseYearlyBalance, *response.ErrorResponse)
}

// CardStatsTopupService handles top-up statistics globally and per specific card number.
type CardStatsTopupService interface {
	// FindMonthlyTopupAmount retrieves monthly top-up statistics across all card numbers for a given year.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - year: the year for which the monthly top-up data is requested
	//
	// Returns:
	//   - A slice of CardResponseMonthAmount or an error response if the operation fails.
	FindMonthlyTopupAmount(ctx context.Context, year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)

	// FindYearlyTopupAmount retrieves yearly top-up statistics across all card numbers for a given year.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - year: the year for which the yearly top-up data is requested
	//
	// Returns:
	//   - A slice of CardResponseYearAmount or an error response if the operation fails.
	FindYearlyTopupAmount(ctx context.Context, year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
}

// CardStatsWithdrawService handles withdraw statistics globally and per specific card number.
type CardStatsWithdrawService interface {
	// FindMonthlyWithdrawAmount retrieves monthly withdraw statistics across all card numbers for a given year.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - year: the year for which the monthly withdraw data is requested
	//
	// Returns:
	//   - A slice of CardResponseMonthAmount or an error response if the operation fails.
	FindMonthlyWithdrawAmount(ctx context.Context, year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)

	// FindYearlyWithdrawAmount retrieves yearly withdraw statistics across all card numbers for a given year.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - year: the year for which the yearly withdraw data is requested
	//
	// Returns:
	//   - A slice of CardResponseYearAmount or an error response if the operation fails.
	FindYearlyWithdrawAmount(ctx context.Context, year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
}

// CardStatsTransactionService handles transaction statistics globally and per specific card number.
type CardStatsTransactionService interface {
	// FindMonthlyTransactionAmount retrieves monthly transaction statistics across all card numbers for a given year.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - year: the year for which the monthly transaction data is requested
	//
	// Returns:
	//   - A slice of CardResponseMonthAmount or an error response if the operation fails.
	FindMonthlyTransactionAmount(ctx context.Context, year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)

	// FindYearlyTransactionAmount retrieves yearly transaction statistics across all card numbers for a given year.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - year: the year for which the yearly transaction data is requested
	//
	// Returns:
	//   - A slice of CardResponseYearAmount or an error response if the operation fails.
	FindYearlyTransactionAmount(ctx context.Context, year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
}

// CardStatsTransferService handles transfer statistics globally and per specific card number (as sender or receiver).
type CardStatsTransferService interface {
	// FindMonthlyTransferAmountSender retrieves total monthly transfer amounts from all cards acting as sender.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - year: the year for which the monthly data is requested
	//
	// Returns:
	//   - A slice of CardResponseMonthAmount or an error response if the operation fails.
	FindMonthlyTransferAmountSender(ctx context.Context, year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)

	// FindYearlyTransferAmountSender retrieves total yearly transfer amounts from all cards acting as sender.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - year: the year for which the yearly data is requested
	//
	// Returns:
	//   - A slice of CardResponseYearAmount or an error response if the operation fails.
	FindYearlyTransferAmountSender(ctx context.Context, year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse)

	// FindMonthlyTransferAmountReceiver retrieves total monthly transfer amounts for all cards acting as receiver.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - year: the year for which the monthly data is requested
	//
	// Returns:
	//   - A slice of CardResponseMonthAmount or an error response if the operation fails.
	FindMonthlyTransferAmountReceiver(ctx context.Context, year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)

	// FindYearlyTransferAmountReceiver retrieves total yearly transfer amounts for all cards acting as receiver.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - year: the year for which the yearly data is requested
	//
	// Returns:
	//   - A slice of CardResponseYearAmount or an error response if the operation fails.
	FindYearlyTransferAmountReceiver(ctx context.Context, year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
}
