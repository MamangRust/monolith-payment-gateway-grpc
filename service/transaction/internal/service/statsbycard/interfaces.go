package transactionstatsbycardservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type TransactionStatsByCardAmountService interface {
	// FindMonthlyAmountsByCardNumber retrieves monthly transaction amounts by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Contains the card number, year, and month.
	//
	// Returns:
	//   - []*response.TransactionMonthAmountResponse: List of monthly transaction amounts.
	//   - *response.ErrorResponse: Error detail if the operation fails.
	FindMonthlyAmountsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*response.TransactionMonthAmountResponse, *response.ErrorResponse)

	// FindYearlyAmountsByCardNumber retrieves yearly transaction amounts by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Contains the card number and year.
	//
	// Returns:
	//   - []*response.TransactionYearlyAmountResponse: List of yearly transaction amounts.
	//   - *response.ErrorResponse: Error detail if the operation fails.
	FindYearlyAmountsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*response.TransactionYearlyAmountResponse, *response.ErrorResponse)
}

type TransactionStatsByCardMethodService interface {
	// FindMonthlyPaymentMethodsByCardNumber retrieves monthly payment method usage by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Contains the card number, year, and month.
	//
	// Returns:
	//   - []*response.TransactionMonthMethodResponse: List of monthly payment method usage.
	//   - *response.ErrorResponse: Error detail if the operation fails.
	FindMonthlyPaymentMethodsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*response.TransactionMonthMethodResponse, *response.ErrorResponse)

	// FindYearlyPaymentMethodsByCardNumber retrieves yearly payment method usage by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Contains the card number and year.
	//
	// Returns:
	//   - []*response.TransactionYearMethodResponse: List of yearly payment method usage.
	//   - *response.ErrorResponse: Error detail if the operation fails.
	FindYearlyPaymentMethodsByCardNumber(ctx context.Context, req *requests.MonthYearPaymentMethod) ([]*response.TransactionYearMethodResponse, *response.ErrorResponse)
}

type TransactionStatsByCardStatusService interface {
	// FindMonthTransactionStatusSuccessByCardNumber retrieves monthly successful transaction statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Contains the card number, year, and month.
	//
	// Returns:
	//   - []*response.TransactionResponseMonthStatusSuccess: List of successful transactions by month.
	//   - *response.ErrorResponse: Error detail if the operation fails.
	FindMonthTransactionStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*response.TransactionResponseMonthStatusSuccess, *response.ErrorResponse)

	// FindYearlyTransactionStatusSuccessByCardNumber retrieves yearly successful transaction statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Contains the card number and year.
	//
	// Returns:
	//   - []*response.TransactionResponseYearStatusSuccess: List of successful transactions by year.
	//   - *response.ErrorResponse: Error detail if the operation fails.
	FindYearlyTransactionStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*response.TransactionResponseYearStatusSuccess, *response.ErrorResponse)

	// FindMonthTransactionStatusFailedByCardNumber retrieves monthly failed transaction statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Contains the card number, year, and month.
	//
	// Returns:
	//   - []*response.TransactionResponseMonthStatusFailed: List of failed transactions by month.
	//   - *response.ErrorResponse: Error detail if the operation fails.
	FindMonthTransactionStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusTransactionCardNumber) ([]*response.TransactionResponseMonthStatusFailed, *response.ErrorResponse)

	// FindYearlyTransactionStatusFailedByCardNumber retrieves yearly failed transaction statistics by card number.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: Contains the card number and year.
	//
	// Returns:
	//   - []*response.TransactionResponseYearStatusFailed: List of failed transactions by year.
	//   - *response.ErrorResponse: Error detail if the operation fails.
	FindYearlyTransactionStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusTransactionCardNumber) ([]*response.TransactionResponseYearStatusFailed, *response.ErrorResponse)
}
