package transactionstatshandler

import (
	"context"

	pb "github.com/MamangRust/monolith-payment-gateway-pb/transaction"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/grpc"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/transaction"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/service"

	servicestats "github.com/MamangRust/monolith-payment-gateway-transaction/internal/service/stats"
	servicestatsbycard "github.com/MamangRust/monolith-payment-gateway-transaction/internal/service/statsbycard"

	"go.uber.org/zap"
)

type transactionStatsAmountHandleGrpc struct {
	pb.UnimplementedTransactionsStatsAmountServiceServer

	servicestats       servicestats.TransactionStatsService
	servicestatsbycard servicestatsbycard.TransactionStatsByCardService
	logger             logger.LoggerInterface
	mapper             protomapper.TransactionStatsAmountProtoMapper
}

func NewTransactionStatsAmountHandleGrpc(
	service service.Service,
	logger logger.LoggerInterface,
	mapper protomapper.TransactionStatsAmountProtoMapper,
) TransactionStatsAmountHandlerGrpc {
	return &transactionStatsAmountHandleGrpc{
		servicestats:       service,
		servicestatsbycard: service,
		logger:             logger,
		mapper:             mapper,
	}
}

// FindMonthlyAmounts retrieves the monthly transaction amounts for a specified year.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindYearTransactionStatus object containing the year for which to fetch the transaction amounts.
//
// Returns:
//   - An ApiResponseTransactionMonthAmount containing the monthly transaction amounts for the specified year.
//   - An error if the operation fails, or if the provided year is invalid.
func (t *transactionStatsAmountHandleGrpc) FindMonthlyAmounts(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionMonthAmount, error) {
	year := int(req.GetYear())

	t.logger.Info("Fetching transaction",
		zap.Int("year", year))

	if year <= 0 {
		t.logger.Error("Failed to fetch transaction", zap.Int("year", year))
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	amounts, err := t.servicestats.FindMonthlyAmounts(ctx, year)

	if err != nil {
		t.logger.Error("failed to fetch monthly amounts", zap.Any("error", err), zap.Int("year", year))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapper.ToProtoResponseTransactionMonthAmount("success", "Successfully fetched monthly amounts", amounts)

	t.logger.Info("Successfully fetch monthly amounts", zap.Int("year", year))

	return so, nil
}

// FindYearlyAmounts retrieves the yearly transaction amounts for a specified year.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindYearTransactionStatus object containing the year for which to fetch the transaction amounts.
//
// Returns:
//   - An ApiResponseTransactionYearAmount containing the yearly transaction amounts for the specified year.
//   - An error if the operation fails, or if the provided year is invalid.
func (t *transactionStatsAmountHandleGrpc) FindYearlyAmounts(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionYearAmount, error) {
	year := int(req.GetYear())

	t.logger.Info("Fetching transaction",
		zap.Int("year", year))

	if year <= 0 {
		t.logger.Error("Failed to fetch transaction", zap.Int("year", year))
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	amounts, err := t.servicestats.FindYearlyAmounts(ctx, year)

	if err != nil {
		t.logger.Error("failed to fetch yearly amounts", zap.Any("error", err), zap.Int("year", year))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapper.ToProtoResponseTransactionYearAmount("success", "Successfully fetched yearly amounts", amounts)

	return so, nil
}

// FindMonthlyAmountsByCardNumber retrieves the monthly transaction amounts for a specific card number and year.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindByYearCardNumberTransactionRequest object containing the card number and year for which to fetch the transaction amounts.
//
// Returns:
//   - An ApiResponseTransactionMonthAmount containing the transaction amounts for each month of the specified year.
//   - An error if the operation fails, or if the provided year or card number is invalid.
func (t *transactionStatsAmountHandleGrpc) FindMonthlyAmountsByCardNumber(ctx context.Context, req *pb.FindByYearCardNumberTransactionRequest) (*pb.ApiResponseTransactionMonthAmount, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	t.logger.Info("Fetching transaction",
		zap.Int("year", year),
		zap.String("cardNumber", cardNumber))

	if year <= 0 {
		t.logger.Error("Failed to fetch transaction", zap.Int("year", year))
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		t.logger.Error("Failed to fetch transaction", zap.String("cardNumber", cardNumber))
		return nil, transaction_errors.ErrGrpcInvalidCardNumber
	}

	reqService := &requests.MonthYearPaymentMethod{
		Year:       year,
		CardNumber: cardNumber,
	}

	amounts, err := t.servicestatsbycard.FindMonthlyAmountsByCardNumber(ctx, reqService)

	if err != nil {
		t.logger.Error("failed to fetch monthly amounts by card number", zap.Any("error", err), zap.Int("year", year), zap.String("card_number", cardNumber))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapper.ToProtoResponseTransactionMonthAmount("success", "Successfully fetched monthly amounts by card number", amounts)

	t.logger.Info("Successfully fetch monthly amounts by card number", zap.Int("year", year), zap.String("card_number", cardNumber))

	return so, nil
}

// FindYearlyAmountsByCardNumber fetches the yearly transaction amount statistics for a given card number from the database.
//
// Parameters:
//   - ctx: The context.Context object for the gRPC request.
//   - req: A FindByYearCardNumber object containing the year and card number to fetch the yearly transaction amount statistics for.
//
// Returns:
//   - An ApiResponseTransactionYearAmount containing the yearly transaction amount statistics retrieved from the database.
//   - An error if the operation fails, or if the provided year or card number is invalid.
func (t *transactionStatsAmountHandleGrpc) FindYearlyAmountsByCardNumber(ctx context.Context, req *pb.FindByYearCardNumberTransactionRequest) (*pb.ApiResponseTransactionYearAmount, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	t.logger.Info("Fetching transaction",
		zap.Int("year", year),
		zap.String("cardNumber", cardNumber))

	if year <= 0 {
		t.logger.Error("Failed to fetch transaction", zap.Int("year", year))
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		t.logger.Error("Failed to fetch transaction", zap.String("cardNumber", cardNumber))
		return nil, transaction_errors.ErrGrpcInvalidCardNumber
	}

	reqService := &requests.MonthYearPaymentMethod{
		Year:       year,
		CardNumber: cardNumber,
	}

	amounts, err := t.servicestatsbycard.FindYearlyAmountsByCardNumber(ctx, reqService)

	if err != nil {
		t.logger.Error("failed to fetch yearly amounts by card number", zap.Any("error", err), zap.Int("year", year), zap.String("card_number", cardNumber))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapper.ToProtoResponseTransactionYearAmount("success", "Successfully fetched yearly amounts by card number", amounts)

	t.logger.Info("Successfully fetch yearly amounts by card number", zap.Int("year", year), zap.String("card_number", cardNumber))

	return so, nil
}
