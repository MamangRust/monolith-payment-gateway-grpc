package transactionstatsservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type TransactionStatsAmountService interface {
	// FindMonthlyAmounts retrieves the total monthly transaction amounts.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*response.TransactionMonthAmountResponse: List of monthly transaction amounts.
	//   - *response.ErrorResponse: Error detail if the operation fails.
	FindMonthlyAmounts(ctx context.Context, year int) ([]*response.TransactionMonthAmountResponse, *response.ErrorResponse)

	// FindYearlyAmounts retrieves the total yearly transaction amounts.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*response.TransactionYearlyAmountResponse: List of yearly transaction amounts.
	//   - *response.ErrorResponse: Error detail if the operation fails.
	FindYearlyAmounts(ctx context.Context, year int) ([]*response.TransactionYearlyAmountResponse, *response.ErrorResponse)
}

type TransactionStatsMethodService interface {
	// FindMonthlyPaymentMethods retrieves monthly usage statistics for each payment method.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*response.TransactionMonthMethodResponse: List of monthly method statistics.
	//   - *response.ErrorResponse: Error detail if the operation fails.
	FindMonthlyPaymentMethods(ctx context.Context, year int) ([]*response.TransactionMonthMethodResponse, *response.ErrorResponse)

	// FindYearlyPaymentMethods retrieves yearly usage statistics for each payment method.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*response.TransactionYearMethodResponse: List of yearly method statistics.
	//   - *response.ErrorResponse: Error detail if the operation fails.
	FindYearlyPaymentMethods(ctx context.Context, year int) ([]*response.TransactionYearMethodResponse, *response.ErrorResponse)
}

type TransactionStatsStatusService interface {
	// FindMonthTransactionStatusSuccess retrieves monthly success statistics for all transactions.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Contains the year and month of the transaction.
	//
	// Returns:
	//   - []*response.TransactionResponseMonthStatusSuccess: List of success statistics.
	//   - *response.ErrorResponse: Error detail if the operation fails.
	FindMonthTransactionStatusSuccess(ctx context.Context, req *requests.MonthStatusTransaction) ([]*response.TransactionResponseMonthStatusSuccess, *response.ErrorResponse)

	// FindYearlyTransactionStatusSuccess retrieves yearly success statistics for all transactions.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*response.TransactionResponseYearStatusSuccess: List of success statistics.
	//   - *response.ErrorResponse: Error detail if the operation fails.
	FindYearlyTransactionStatusSuccess(ctx context.Context, year int) ([]*response.TransactionResponseYearStatusSuccess, *response.ErrorResponse)

	// FindMonthTransactionStatusFailed retrieves monthly failed statistics for all transactions.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Contains the year and month of the transaction.
	//
	// Returns:
	//   - []*response.TransactionResponseMonthStatusFailed: List of failed statistics.
	//   - *response.ErrorResponse: Error detail if the operation fails.
	FindMonthTransactionStatusFailed(ctx context.Context, req *requests.MonthStatusTransaction) ([]*response.TransactionResponseMonthStatusFailed, *response.ErrorResponse)

	// FindYearlyTransactionStatusFailed retrieves yearly failed statistics for all transactions.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - year: The year for which data is requested.
	//
	// Returns:
	//   - []*response.TransactionResponseYearStatusFailed: List of failed statistics.
	//   - *response.ErrorResponse: Error detail if the operation fails.
	FindYearlyTransactionStatusFailed(ctx context.Context, year int) ([]*response.TransactionResponseYearStatusFailed, *response.ErrorResponse)
}
