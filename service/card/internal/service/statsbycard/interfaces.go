package cardstatsbycard

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// CardStatsBalanceByCardService provides methods for retrieving card balance statistics by card number.
type CardStatsBalanceByCardService interface {
	// FindMonthlyBalanceByCardNumber retrieves monthly balance statistics for a specific card number and year.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - req: a request object containing the month, year, and card number
	//
	// Returns:
	//   - A slice of CardResponseMonthBalance or an error response if the operation fails.
	FindMonthlyBalanceByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthBalance, *response.ErrorResponse)

	// FindYearlyBalanceByCardNumber retrieves yearly balance statistics for a specific card number and year.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - req: a request object containing the year and card number
	//
	// Returns:
	//   - A slice of CardResponseYearlyBalance or an error response if the operation fails.
	FindYearlyBalanceByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearlyBalance, *response.ErrorResponse)
}

type CardStatsTopupByCardService interface {
	// FindMonthlyTopupAmountByCardNumber retrieves monthly top-up statistics for a specific card number and year.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - req: a request object containing the month, year, and card number
	//
	// Returns:
	//   - A slice of CardResponseMonthAmount or an error response if the operation fails.
	FindMonthlyTopupAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)

	// FindYearlyTopupAmountByCardNumber retrieves yearly top-up statistics for a specific card number and year.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - req: a request object containing the year and card number
	//
	// Returns:
	//   - A slice of CardResponseYearAmount or an error response if the operation fails.
	FindYearlyTopupAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
}

type CardStatsWithdrawByCardService interface {
	// FindMonthlyWithdrawAmountByCardNumber retrieves monthly withdraw statistics for a specific card number and year.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - req: a request object containing the month, year, and card number
	//
	// Returns:
	//   - A slice of CardResponseMonthAmount or an error response if the operation fails.
	FindMonthlyWithdrawAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)

	// FindYearlyWithdrawAmountByCardNumber retrieves yearly withdraw statistics for a specific card number and year.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - req: a request object containing the year and card number
	//
	// Returns:
	//   - A slice of CardResponseYearAmount or an error response if the operation fails.
	FindYearlyWithdrawAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
}

type CardStatsTransactionByCardService interface {
	// FindMonthlyTransactionAmountByCardNumber retrieves monthly transaction statistics for a specific card number.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - req: a request object containing the month, year, and card number
	//
	// Returns:
	//   - A slice of CardResponseMonthAmount or an error response if the operation fails.
	FindMonthlyTransactionAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)

	// FindYearlyTransactionAmountByCardNumber retrieves yearly transaction statistics for a specific card number.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - req: a request object containing the year and card number
	//
	// Returns:
	//   - A slice of CardResponseYearAmount or an error response if the operation fails.
	FindYearlyTransactionAmountByCardNumber(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
}

type CardStatsTransferByCardService interface {
	// FindMonthlyTransferAmountBySender retrieves monthly transfer statistics for a specific sender card number.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - req: request containing year, month, and sender card number
	//
	// Returns:
	//   - A slice of CardResponseMonthAmount or an error response if the operation fails.
	FindMonthlyTransferAmountBySender(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)

	// FindYearlyTransferAmountBySender retrieves yearly transfer statistics for a specific sender card number.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - req: request containing year and sender card number
	//
	// Returns:
	//   - A slice of CardResponseYearAmount or an error response if the operation fails.
	FindYearlyTransferAmountBySender(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse)

	// FindMonthlyTransferAmountByReceiver retrieves monthly transfer statistics for a specific receiver card number.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - req: request containing year, month, and receiver card number
	//
	// Returns:
	//   - A slice of CardResponseMonthAmount or an error response if the operation fails.
	FindMonthlyTransferAmountByReceiver(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse)

	// FindYearlyTransferAmountByReceiver retrieves yearly transfer statistics for a specific receiver card number.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - req: request containing year and receiver card number
	//
	// Returns:
	//   - A slice of CardResponseYearAmount or an error response if the operation fails.
	FindYearlyTransferAmountByReceiver(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse)
}
