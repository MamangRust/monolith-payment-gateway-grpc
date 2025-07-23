package transferstatshandler

import (
	"context"

	pb "github.com/MamangRust/monolith-payment-gateway-pb/transfer"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/grpc"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/transfer"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/service"
	"go.uber.org/zap"

	servicestats "github.com/MamangRust/monolith-payment-gateway-transfer/internal/service/stats"
	servicestatsbycard "github.com/MamangRust/monolith-payment-gateway-transfer/internal/service/statsbycard"
)

type transferStatsAmountHandleGrpc struct {
	pb.UnimplementedTransferStatsAmountServiceServer

	servicestats       servicestats.TransferStatsService
	servicestatsbycard servicestatsbycard.TransferStatsByCardService
	logger             logger.LoggerInterface
	mapper             protomapper.TransferStatsAmountProtoMapper
}

func NewTransferStatsAmountHandler(service service.Service, logger logger.LoggerInterface, mapper protomapper.TransferStatsAmountProtoMapper) TransferStatsAmountHandleGrpc {
	return &transferStatsAmountHandleGrpc{
		servicestats:       service,
		servicestatsbycard: service,
		logger:             logger,
		mapper:             mapper,
	}
}

// FindMonthlyTransferAmounts retrieves the monthly transfer amounts for the specified year.
//
// It logs the process of fetching the transfer amounts and validates the input,
// ensuring that the year is a positive integer. If the validation fails, it logs
// the error and returns an appropriate error response. The function queries the
// transfer statistics service to fetch the transfer amounts and maps the result
// to a protocol buffer response object.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindYearTransferStatus object containing the year for which to fetch
//     the transfer amounts.
//
// Returns:
//   - An ApiResponseTransferMonthAmount containing the transfer amounts for
//     the specified year.
//   - An error if the operation fails, or if the provided year is invalid.
func (s *transferStatsAmountHandleGrpc) FindMonthlyTransferAmounts(ctx context.Context, req *pb.FindYearTransferStatus) (*pb.ApiResponseTransferMonthAmount, error) {
	year := int(req.GetYear())

	s.logger.Info("Fetching transfer", zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("year", year))
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	amounts, err := s.servicestats.FindMonthlyTransferAmounts(ctx, year)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTransferMonthAmount("success", "Successfully fetched monthly transfer amounts", amounts)

	s.logger.Info("Successfully fetched monthly transfer amounts", zap.Int("year", year))

	return so, nil
}

// FindYearlyTransferAmounts retrieves the yearly transfer amounts for a specified year.
//
// It logs the process of fetching the transfer amounts and validates the input,
// ensuring that the year is a positive integer. If the validation fails, it logs
// the error and returns an appropriate error response. The function queries the
// transfer statistics service to fetch the transfer amounts and maps the result
// to a protocol buffer response object.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindYearTransferStatus object containing the year for which to fetch
//     the transfer amounts.
//
// Returns:
//   - An ApiResponseTransferYearAmount containing the transfer amounts for
//     the specified year.
//   - An error if the operation fails, or if the provided year is invalid.
func (s *transferStatsAmountHandleGrpc) FindYearlyTransferAmounts(ctx context.Context, req *pb.FindYearTransferStatus) (*pb.ApiResponseTransferYearAmount, error) {
	year := int(req.GetYear())

	s.logger.Info("Fetching transfer", zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("year", year))
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	amounts, err := s.servicestats.FindYearlyTransferAmounts(ctx, year)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTransferYearAmount("success", "Successfully fetched yearly transfer amounts", amounts)

	s.logger.Info("Successfully fetched yearly transfer amounts", zap.Int("year", year))

	return so, nil
}

// FindMonthlyTransferAmountsBySenderCardNumber retrieves the monthly transfer amounts
// for a specified year and sender card number.
//
// The function logs the process of fetching the transfer data and validates the inputs,
// ensuring that the year is a positive integer and the card number is not empty. If
// any of these validations fail, it logs the error and returns an appropriate error
// response. The function queries the transfer statistics service by card number to
// fetch the transfer amounts and maps the result to a protocol buffer response object.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindByCardNumberTransferRequest object containing the year and card number
//     for which to fetch the transfer amounts.
//
// Returns:
//   - An ApiResponseTransferMonthAmount containing the transfer amounts for the specified
//     year and card number.
//   - An error if the operation fails, or if the provided year or card number are invalid.
func (s *transferStatsAmountHandleGrpc) FindMonthlyTransferAmountsBySenderCardNumber(ctx context.Context, req *pb.FindByCardNumberTransferRequest) (*pb.ApiResponseTransferMonthAmount, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	s.logger.Info("Fetching transfer", zap.Int("year", year), zap.String("card_number", cardNumber))

	if year <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("year", year))
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		s.logger.Error("Failed to fetch transfer", zap.String("card_number", cardNumber))
		return nil, transfer_errors.ErrGrpcInvalidCardNumber
	}

	reqService := &requests.MonthYearCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	amounts, err := s.servicestatsbycard.FindMonthlyTransferAmountsBySenderCardNumber(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTransferMonthAmount("success", "Successfully fetched monthly transfer amounts by sender card number", amounts)

	s.logger.Info("Successfully fetched monthly transfer amounts by sender card number", zap.Int("year", year), zap.String("card_number", cardNumber))

	return so, nil
}

// FindMonthlyTransferAmountsByReceiverCardNumber retrieves the monthly transfer amounts
// for a specified year and receiver card number.
//
// The function logs the process of fetching the transfer data and validates the inputs,
// ensuring that the year is a positive integer and the card number is not empty. If
// any of these validations fail, it logs the error and returns an appropriate error
// response. The function queries the transfer statistics service by card number to
// fetch the transfer amounts and maps the result to a protocol buffer response object.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindByCardNumberTransferRequest object containing the year and card number
//     for which to fetch the transfer amounts.
//
// Returns:
//   - An ApiResponseTransferMonthAmount containing the transfer amounts for the specified
//     year and card number.
//   - An error if the operation fails, or if the provided year or card number are invalid.
func (s *transferStatsAmountHandleGrpc) FindMonthlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *pb.FindByCardNumberTransferRequest) (*pb.ApiResponseTransferMonthAmount, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	s.logger.Info("Fetching transfer", zap.Int("year", year), zap.String("card_number", cardNumber))

	if year <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("year", year))
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		s.logger.Error("Failed to fetch transfer", zap.String("card_number", cardNumber))
		return nil, transfer_errors.ErrGrpcInvalidCardNumber
	}

	reqService := &requests.MonthYearCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	amounts, err := s.servicestatsbycard.FindMonthlyTransferAmountsByReceiverCardNumber(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTransferMonthAmount("success", "Successfully fetched monthly transfer amounts by receiver card number", amounts)

	s.logger.Info("Successfully fetched monthly transfer amounts by receiver card number", zap.Int("year", year), zap.String("card_number", cardNumber))

	return so, nil
}

// FindYearlyTransferAmountsBySenderCardNumber retrieves the yearly transfer amounts
// for a specified year and sender card number.
//
// It logs the process of fetching the transfer data and validates the inputs,
// ensuring that the year is a positive integer and the card number is not empty.
// If the validation fails, it logs the error and returns an appropriate error
// response. The function queries the transfer statistics by card number service
// to fetch the transfer amounts and maps the result to a protocol buffer response
// object.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindByCardNumberTransferRequest object containing the year and card
//     number for which to fetch the transfer amounts.
//
// Returns:
//   - An ApiResponseTransferYearAmount containing the transfer amounts for the
//     specified year and card number.
//   - An error if the operation fails, or if the provided year or card number are
//     invalid.
func (s *transferStatsAmountHandleGrpc) FindYearlyTransferAmountsBySenderCardNumber(ctx context.Context, req *pb.FindByCardNumberTransferRequest) (*pb.ApiResponseTransferYearAmount, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	s.logger.Info("Fetching transfer", zap.Int("year", year), zap.String("card_number", cardNumber))

	if year <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("year", year))
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		s.logger.Error("Failed to fetch transfer", zap.String("card_number", cardNumber))
		return nil, transfer_errors.ErrGrpcInvalidCardNumber
	}

	reqService := &requests.MonthYearCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	amounts, err := s.servicestatsbycard.FindYearlyTransferAmountsBySenderCardNumber(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTransferYearAmount("success", "Successfully fetched yearly transfer amounts by sender card number", amounts)

	s.logger.Info("Successfully fetched yearly transfer amounts by sender card number", zap.Int("year", year), zap.String("card_number", cardNumber))

	return so, nil
}

// FindYearlyTransferAmountsByReceiverCardNumber retrieves the yearly transfer amounts
// for a specified year and receiver card number.
//
// It logs the process of fetching the transfer data and validates the inputs,
// ensuring that the year is a positive integer and the card number is not empty.
// If the validation fails, it logs the error and returns an appropriate error
// response. The function queries the transfer statistics by card number service
// to fetch the transfer amounts and maps the result to a protocol buffer response
// object.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindByCardNumberTransferRequest object containing the year and card
//     number for which to fetch the transfer amounts.
//
// Returns:
//   - An ApiResponseTransferYearAmount containing the transfer amounts for the
//     specified year and card number.
//   - An error if the operation fails, or if the provided year or card number are
//     invalid.
func (s *transferStatsAmountHandleGrpc) FindYearlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *pb.FindByCardNumberTransferRequest) (*pb.ApiResponseTransferYearAmount, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	s.logger.Info("Fetching transfer", zap.Int("year", year), zap.String("card_number", cardNumber))

	if year <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("year", year))
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		s.logger.Error("Failed to fetch transfer", zap.String("card_number", cardNumber))
		return nil, transfer_errors.ErrGrpcInvalidCardNumber
	}

	reqService := &requests.MonthYearCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	amounts, err := s.servicestatsbycard.FindYearlyTransferAmountsByReceiverCardNumber(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTransferYearAmount("success", "Successfully fetched yearly transfer amounts by receiver card number", amounts)

	s.logger.Info("Successfully fetched yearly transfer amounts by receiver card number", zap.Int("year", year), zap.String("card_number", cardNumber))

	return so, nil
}
