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

type transactionStatsStatusHandleGrpc struct {
	pb.UnimplementedTransactionStatsStatusServiceServer

	servicestats       servicestats.TransactionStatsService
	servicestatsbycard servicestatsbycard.TransactionStatsByCardService
	logger             logger.LoggerInterface
	mapper             protomapper.TransactionStatsStatusProtoMapper
}

func NewTransactionStatsStatusHandleGrpc(
	service service.Service,
	logger logger.LoggerInterface,
	mapper protomapper.TransactionStatsStatusProtoMapper,
) TransactionStatsStatusHandleGrpc {
	return &transactionStatsStatusHandleGrpc{
		servicestats:       service,
		servicestatsbycard: service,
		logger:             logger,
		mapper:             mapper,
	}
}

// FindMonthlyTransactionStatusSuccess implements the gRPC service for fetching the monthly transaction status success.
//
// Parameters:
//   - ctx: the context.Context of the gRPC request.
//   - req: the request containing the year and month.
//
// Returns:
//   - A pointer to the API response containing the transaction status success records.
//   - An error if any.
func (t *transactionStatsStatusHandleGrpc) FindMonthlyTransactionStatusSuccess(ctx context.Context, req *pb.FindMonthlyTransactionStatus) (*pb.ApiResponseTransactionMonthStatusSuccess, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	t.logger.Info("Fetching transaction",
		zap.Int("year", year),
		zap.Int("month", month))

	if year <= 0 {
		t.logger.Error("Failed to fetch transaction", zap.Int("year", year))
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		t.logger.Error("Failed to fetch transaction", zap.Int("month", month))
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	reqService := &requests.MonthStatusTransaction{
		Year:  year,
		Month: month,
	}

	records, err := t.servicestats.FindMonthTransactionStatusSuccess(ctx, reqService)
	if err != nil {
		t.logger.Error("failed to fetch monthly Transaction status success", zap.Any("error", err), zap.Int("year", year), zap.Int("month", month))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapper.ToProtoResponseTransactionMonthStatusSuccess("success", "Successfully fetched monthly Transaction status success", records)

	t.logger.Info("Successfully fetch monthly Transaction status success", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

// FindYearlyTransactionStatusSuccess retrieves the yearly transaction status for successful transactions.
//
// Parameters:
//   - req: A pointer to a FindYearTransactionStatus request object, containing the year.
//
// Returns:
//   - A pointer to the API response containing the transaction status success records.
//   - An error if any.
func (t *transactionStatsStatusHandleGrpc) FindYearlyTransactionStatusSuccess(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionYearStatusSuccess, error) {
	year := int(req.GetYear())

	t.logger.Info("Fetching transaction",
		zap.Int("year", year))

	if year <= 0 {
		t.logger.Error("Failed to fetch transaction", zap.Int("year", year))
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	records, err := t.servicestats.FindYearlyTransactionStatusSuccess(ctx, year)
	if err != nil {
		t.logger.Error("failed to fetch yearly Transaction status success", zap.Any("error", err), zap.Int("year", year))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapper.ToProtoResponseTransactionYearStatusSuccess("success", "Successfully fetched yearly Transaction status success", records)

	t.logger.Info("Successfully fetch yearly Transaction status success", zap.Int("year", year))

	return so, nil
}

// FindMonthlyTransactionStatusFailed implements the gRPC service for fetching the monthly transaction status failed.
//
// Parameters:
//   - ctx: the context.Context of the gRPC request.
//   - req: the request containing the year and month.
//
// Returns:
//   - A pointer to the API response containing the transaction status failed records.
//   - An error if any.
func (t *transactionStatsStatusHandleGrpc) FindMonthlyTransactionStatusFailed(ctx context.Context, req *pb.FindMonthlyTransactionStatus) (*pb.ApiResponseTransactionMonthStatusFailed, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	t.logger.Info("Fetching transaction",
		zap.Int("year", year),
		zap.Int("month", month))

	if year <= 0 {
		t.logger.Error("Failed to fetch transaction", zap.Int("year", year))
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		t.logger.Error("Failed to fetch transaction", zap.Int("month", month))
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	reqService := &requests.MonthStatusTransaction{
		Year:  year,
		Month: month,
	}

	records, err := t.servicestats.FindMonthTransactionStatusFailed(ctx, reqService)

	if err != nil {
		t.logger.Error("failed to fetch monthly Transaction status Failed", zap.Any("error", err), zap.Int("year", year), zap.Int("month", month))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapper.ToProtoResponseTransactionMonthStatusFailed("success", "success fetched monthly Transaction status Failed", records)

	t.logger.Info("Successfully fetch monthly Transaction status Failed", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

// FindYearlyTransactionStatusFailed retrieves the yearly transaction status for failed transactions.
//
// Args:
//
//	ctx: The context.Context object for the current request.
//	req: The gRPC request containing the year to filter by.
//
// Returns:
//   - A pointer to the API response containing the transaction status failed records.
//   - An error if any.
func (t *transactionStatsStatusHandleGrpc) FindYearlyTransactionStatusFailed(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionYearStatusFailed, error) {
	year := int(req.GetYear())

	t.logger.Info("Fetching transaction",
		zap.Int("year", year))

	if year <= 0 {
		t.logger.Error("Failed to fetch transaction", zap.Int("year", year))
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	records, err := t.servicestats.FindYearlyTransactionStatusFailed(ctx, year)

	if err != nil {
		t.logger.Error("failed to fetch yearly Transaction status Failed", zap.Any("error", err), zap.Int("year", year))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapper.ToProtoResponseTransactionYearStatusFailed("success", "success fetched yearly Transaction status Failed", records)

	t.logger.Info("Successfully fetch yearly Transaction status Failed", zap.Int("year", year))

	return so, nil
}

// FindMonthlyTransactionStatusSuccessByCardNumber retrieves the monthly transaction status for successful transactions
// filtered by card number, year, and month.
//
// It validates the input parameters, logs the request, and fetches the transaction status from the service layer.
// If the input parameters are invalid, appropriate errors are returned. It maps the results to a protocol buffer response
// and returns it.
//
// Parameters:
//   - ctx: the context of the gRPC request.
//   - req: the request containing the year, month, and card number.
//
// Returns:
//   - A pointer to the API response containing the transaction status success records for the specified card number, year, and month.
//   - An error if any validation or retrieval operation fails.
func (t *transactionStatsStatusHandleGrpc) FindMonthlyTransactionStatusSuccessByCardNumber(ctx context.Context, req *pb.FindMonthlyTransactionStatusCardNumber) (*pb.ApiResponseTransactionMonthStatusSuccess, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	cardNumber := req.GetCardNumber()

	t.logger.Info("Fetching transaction",
		zap.Int("year", year),
		zap.Int("month", month),
		zap.String("cardNumber", cardNumber))

	if year <= 0 {
		t.logger.Error("Failed to fetch transaction", zap.Int("year", year))
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		t.logger.Error("Failed to fetch transaction", zap.Int("month", month))
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	if cardNumber == "" {
		t.logger.Error("Failed to fetch transaction", zap.String("cardNumber", cardNumber))
		return nil, transaction_errors.ErrGrpcInvalidCardNumber
	}

	reqService := &requests.MonthStatusTransactionCardNumber{
		CardNumber: cardNumber,
		Year:       year,
		Month:      month,
	}

	records, err := t.servicestatsbycard.FindMonthTransactionStatusSuccessByCardNumber(ctx, reqService)

	if err != nil {
		t.logger.Error("failed to fetch monthly Transaction status success", zap.Any("error", err), zap.Int("year", year), zap.Int("month", month), zap.String("card_number", cardNumber))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapper.ToProtoResponseTransactionMonthStatusSuccess("success", "Successfully fetched monthly Transaction status success", records)

	t.logger.Info("Successfully fetch monthly Transaction status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", cardNumber))

	return so, nil
}

// FindYearlyTransactionStatusSuccessByCardNumber retrieves the yearly transaction status for successful transactions
// filtered by card number.
//
// Parameters:
//
//	ctx: The context.Context object for the current request.
//	req: The gRPC request containing the year and card number.
//
// Returns:
//   - A pointer to the API response containing the transaction status success records for the specified card number and year.
//   - An error if any validation or retrieval operation fails.
func (s *transactionStatsStatusHandleGrpc) FindYearlyTransactionStatusSuccessByCardNumber(ctx context.Context, req *pb.FindYearTransactionStatusCardNumber) (*pb.ApiResponseTransactionYearStatusSuccess, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	s.logger.Info("Fetching transaction",
		zap.Int("year", year),
		zap.String("cardNumber", cardNumber))

	if year <= 0 {
		s.logger.Error("Failed to fetch transaction", zap.Int("year", year))
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		s.logger.Error("Failed to fetch transaction", zap.String("cardNumber", cardNumber))
		return nil, transaction_errors.ErrGrpcInvalidCardNumber
	}

	reqService := &requests.YearStatusTransactionCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	records, err := s.servicestatsbycard.FindYearlyTransactionStatusSuccessByCardNumber(ctx,reqService)

	if err != nil {
		s.logger.Error("failed to fetch yearly Transaction status success", zap.Any("error", err), zap.Int("year", year), zap.String("card_number", cardNumber))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTransactionYearStatusSuccess("success", "Successfully fetched yearly Transaction status success", records)

	s.logger.Info("Successfully fetch yearly Transaction status success", zap.Int("year", year), zap.String("card_number", cardNumber))

	return so, nil
}

// FindMonthlyTransactionStatusFailedByCardNumber retrieves the monthly transaction status for failed transactions
// filtered by card number.
//
// Parameters:
//
//	ctx: The context.Context object for the current request.
//	req: The gRPC request containing the year, month and card number.
//
// Returns:
//   - A pointer to the API response containing the transaction status failed records for the specified card number and year.
//   - An error if any validation or retrieval operation fails.
func (t *transactionStatsStatusHandleGrpc) FindMonthlyTransactionStatusFailedByCardNumber(ctx context.Context, req *pb.FindMonthlyTransactionStatusCardNumber) (*pb.ApiResponseTransactionMonthStatusFailed, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	cardNumber := req.GetCardNumber()

	t.logger.Info("Fetching transaction",
		zap.Int("year", year),
		zap.Int("month", month),
		zap.String("cardNumber", cardNumber))

	if year <= 0 {
		t.logger.Error("Failed to fetch transaction", zap.Int("year", year))
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		t.logger.Error("Failed to fetch transaction", zap.Int("month", month))
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	if cardNumber == "" {
		t.logger.Error("Failed to fetch transaction", zap.String("cardNumber", cardNumber))
		return nil, transaction_errors.ErrGrpcInvalidCardNumber
	}

	reqService := &requests.MonthStatusTransactionCardNumber{
		CardNumber: cardNumber,
		Year:       year,
		Month:      month,
	}

	records, err := t.servicestatsbycard.FindMonthTransactionStatusFailedByCardNumber(ctx,reqService)

	if err != nil {
		t.logger.Error("failed to fetch monthly Transaction status Failed", zap.Any("error", err), zap.Int("year", year), zap.Int("month", month), zap.String("card_number", cardNumber))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapper.ToProtoResponseTransactionMonthStatusFailed("success", "success fetched monthly Transaction status Failed", records)

	t.logger.Info("Successfully fetch monthly Transaction status Failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", cardNumber))

	return so, nil
}

// FindYearlyTransactionStatusFailedByCardNumber retrieves the yearly transaction status for failed transactions
// based on the provided card number and year.
//
// It logs the process of fetching the transaction data and validates the inputs, ensuring that the year is a positive
// integer and the card number is not empty. If any of these validations fail, it logs the error and returns an
// appropriate error response. The function queries the transaction statistics service to fetch the transaction status
// and maps the result to a protocol buffer response object.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindYearTransactionStatusCardNumber object containing the card number and year for which to fetch
//     the transaction status.
//
// Returns:
//   - An ApiResponseTransactionYearStatusFailed containing the transaction status for the specified year and card number.
//   - An error if the operation fails, or if the provided year or card number is invalid.
func (s *transactionStatsStatusHandleGrpc) FindYearlyTransactionStatusFailedByCardNumber(ctx context.Context, req *pb.FindYearTransactionStatusCardNumber) (*pb.ApiResponseTransactionYearStatusFailed, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	s.logger.Info("Fetching transaction",
		zap.Int("year", year),
		zap.String("cardNumber", cardNumber))

	if year <= 0 {
		s.logger.Error("Failed to fetch transaction", zap.Int("year", year))
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		s.logger.Error("Failed to fetch transaction", zap.String("cardNumber", cardNumber))
		return nil, transaction_errors.ErrGrpcInvalidCardNumber
	}

	reqService := &requests.YearStatusTransactionCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	records, err := s.servicestatsbycard.FindYearlyTransactionStatusFailedByCardNumber(ctx,reqService)

	if err != nil {
		s.logger.Error("failed to fetch yearly Transaction status Failed", zap.Any("error", err), zap.Int("year", year), zap.String("card_number", cardNumber))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTransactionYearStatusFailed("success", "success fetched yearly Transaction status Failed", records)

	s.logger.Info("Successfully fetch yearly Transaction status Failed", zap.Int("year", year), zap.String("card_number", cardNumber))

	return so, nil
}
