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

type transactionStatsMethodHandleGrpc struct {
	pb.UnimplementedTransactionStatsMethodServiceServer

	servicestats       servicestats.TransactionStatsService
	servicestatsbycard servicestatsbycard.TransactionStatsByCardService
	logger             logger.LoggerInterface
	mapper             protomapper.TransactionStatsMethodProtoMapper
}

func NewTransactionStatsMethodHandleGrpc(
	service service.Service,
	logger logger.LoggerInterface,
	mapper protomapper.TransactionStatsMethodProtoMapper,
) TransactionStatsMethodHandleGrpc {
	return &transactionStatsMethodHandleGrpc{
		servicestats:       service,
		servicestatsbycard: service,
		logger:             logger,
		mapper:             mapper,
	}
}

// FindMonthlyPaymentMethods retrieves the monthly payment methods used in the specified year.
//
// Parameters:
//   - ctx: The context of the request.
//   - req: A pointer to the FindYearTransactionStatus containing the year to be fetched.
//
// Returns:
//   - An ApiResponseTransactionMonthMethod containing the monthly payment methods used in the specified year.
//   - An error if the operation fails, or if the provided year is invalid.
func (t *transactionStatsMethodHandleGrpc) FindMonthlyPaymentMethods(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionMonthMethod, error) {
	year := int(req.GetYear())

	t.logger.Info("Fetching transaction",
		zap.Int("year", year))

	if year <= 0 {
		t.logger.Error("Failed to fetch transaction", zap.Int("year", year))
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	methods, err := t.servicestats.FindMonthlyPaymentMethods(ctx,year)

	if err != nil {
		t.logger.Error("failed to fetch monthly payment methods", zap.Any("error", err), zap.Int("year", year))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapper.ToProtoResponseTransactionMonthMethod("success", "Successfully fetched monthly payment methods", methods)

	return so, nil
}

// FindYearlyPaymentMethods retrieves the yearly payment methods used in the specified year.
//
// Parameters:
//   - ctx: The context of the request.
//   - req: A pointer to the FindYearTransactionStatus containing the year to be fetched.
//
// Returns:
//   - An ApiResponseTransactionYearMethod containing the yearly payment methods used in the specified year.
//   - An error if the operation fails, or if the provided year is invalid.
func (t *transactionStatsMethodHandleGrpc) FindYearlyPaymentMethods(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionYearMethod, error) {
	year := int(req.GetYear())

	t.logger.Info("Fetching transaction",
		zap.Int("year", year))

	if year <= 0 {
		t.logger.Error("Failed to fetch transaction", zap.Int("year", year))
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	methods, err := t.servicestats.FindYearlyPaymentMethods(ctx,year)

	if err != nil {
		t.logger.Error("failed to fetch yearly payment methods", zap.Any("error", err), zap.Int("year", year))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapper.ToProtoResponseTransactionYearMethod("success", "Successfully fetched yearly payment methods", methods)

	t.logger.Info("Successfully fetch yearly payment methods", zap.Int("year", year))

	return so, nil
}



// FindMonthlyPaymentMethodsByCardNumber retrieves the monthly payment methods for transactions by card number and year.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindByYearCardNumberTransactionRequest object containing the card number and year for which to fetch the payment methods.
//
// Returns:
//   - An ApiResponseTransactionMonthMethod containing the payment methods for each month of the specified year.
//   - An error if the operation fails, or if the provided year or card number is invalid.
func (t *transactionStatsMethodHandleGrpc) FindMonthlyPaymentMethodsByCardNumber(ctx context.Context, req *pb.FindByYearCardNumberTransactionRequest) (*pb.ApiResponseTransactionMonthMethod, error) {
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

	methods, err := t.servicestatsbycard.FindMonthlyPaymentMethodsByCardNumber(ctx,reqService)

	if err != nil {
		t.logger.Error("failed to fetch monthly payment methods by card number", zap.Any("error", err), zap.Int("year", year), zap.String("card_number", cardNumber))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapper.ToProtoResponseTransactionMonthMethod("success", "Successfully fetched monthly payment methods by card number", methods)

	t.logger.Info("Successfully fetch monthly payment methods by card number", zap.Int("year", year), zap.String("card_number", cardNumber))

	return so, nil
}

// FindYearlyPaymentMethodsByCardNumber retrieves the yearly payment methods for transactions by card number and year.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindByYearCardNumberTransactionRequest object containing the card number and year for which to fetch the payment methods.
//
// Returns:
//   - An ApiResponseTransactionYearMethod containing the yearly payment methods for the specified year.
//   - An error if the operation fails, or if the provided year or card number is invalid.
func (t *transactionStatsMethodHandleGrpc) FindYearlyPaymentMethodsByCardNumber(ctx context.Context, req *pb.FindByYearCardNumberTransactionRequest) (*pb.ApiResponseTransactionYearMethod, error) {
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

	methods, err := t.servicestatsbycard.FindYearlyPaymentMethodsByCardNumber(ctx,reqService)

	if err != nil {
		t.logger.Error("failed to fetch yearly payment methods by card number", zap.Any("error", err), zap.Int("year", year), zap.String("card_number", cardNumber))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapper.ToProtoResponseTransactionYearMethod("success", "Successfully fetched yearly payment methods by card number", methods)

	t.logger.Info("Successfully fetch yearly payment methods by card number", zap.Int("year", year), zap.String("card_number", cardNumber))

	return so, nil
}