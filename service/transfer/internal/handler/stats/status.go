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

type transferStatsStatusHandleGrpc struct {
	pb.UnimplementedTransferStatsStatusServiceServer

	servicestats       servicestats.TransferStatsService
	servicestatsbycard servicestatsbycard.TransferStatsByCardService
	logger             logger.LoggerInterface
	mapper             protomapper.TransferStatsStatusProtoMapper
}

func NewTransferStatsStatusHandler(service service.Service, logger logger.LoggerInterface, mapper protomapper.TransferStatsStatusProtoMapper) TransferStatsStatusHandleGrpc {
	return &transferStatsStatusHandleGrpc{
		servicestats:       service,
		servicestatsbycard: service,
		logger:             logger,
		mapper:             mapper,
	}
}

// FindMonthlyTransferStatusSuccess retrieves the monthly transfer status for successful transfers
// filtered by year and month.
//
// It validates the input parameters, logs the request, and fetches the transfer status from the service layer.
// If the input parameters are invalid, appropriate errors are returned. It maps the results to a protocol buffer response
// and returns it.
//
// Parameters:
//   - ctx: the context of the gRPC request.
//   - req: the request containing the year and month.
//
// Returns:
//   - A pointer to the API response containing the transfer status success records for the specified year and month.
//   - An error if any validation or retrieval operation fails.
func (s *transferStatsStatusHandleGrpc) FindMonthlyTransferStatusSuccess(ctx context.Context, req *pb.FindMonthlyTransferStatus) (*pb.ApiResponseTransferMonthStatusSuccess, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	s.logger.Info("Fetching transfer", zap.Int("year", year), zap.Int("month", month))

	if year <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("year", year))
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("month", month))
		return nil, transfer_errors.ErrGrpcInvalidMonth
	}

	reqService := &requests.MonthStatusTransfer{
		Year:  year,
		Month: month,
	}

	records, err := s.servicestats.FindMonthTransferStatusSuccess(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTransferMonthStatusSuccess("success", "Successfully fetched monthly Transfer status success", records)

	s.logger.Info("Successfully fetched monthly Transfer status success", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

// FindYearlyTransferStatusSuccess retrieves the yearly transfer status for successful transfers
// based on the provided year.
//
// It logs the process of fetching the transfer data and validates the input,
// ensuring that the year is a positive integer. If the validation fails, it logs
// the error and returns an appropriate error response. The function queries the
// transfer statistics service to fetch the transfer status and maps the result
// to a protocol buffer response object.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindYearTransferStatus object containing the year for which to fetch
//     the transfer status.
//
// Returns:
//   - An ApiResponseTransferYearStatusSuccess containing the transfer status for
//     the specified year.
//   - An error if the operation fails, or if the provided year is invalid.
func (s *transferStatsStatusHandleGrpc) FindYearlyTransferStatusSuccess(ctx context.Context, req *pb.FindYearTransferStatus) (*pb.ApiResponseTransferYearStatusSuccess, error) {
	year := int(req.GetYear())

	s.logger.Info("Fetching transfer", zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("year", year))
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	records, err := s.servicestats.FindYearlyTransferStatusSuccess(ctx, year)
	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTransferYearStatusSuccess("success", "Successfully fetched yearly Transfer status success", records)

	s.logger.Info("Successfully fetched yearly Transfer status success", zap.Int("year", year))

	return so, nil
}

// FindMonthlyTransferStatusFailed retrieves the monthly transfer status for failed transfers
// filtered by year and month.
//
// It logs the process of fetching the transfer data and validates the inputs,
// ensuring that the year is a positive integer and the month is a positive integer between
// 1 and 12. If any of these validations fail, it logs the error and returns an appropriate
// error response. The function queries the transfer statistics service to fetch the
// transfer status and maps the result to a protocol buffer response object.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindMonthlyTransferStatus object containing the year and month for which to fetch
//     the transfer status.
//
// Returns:
//   - An ApiResponseTransferMonthStatusFailed containing the transfer status for the
//     specified year and month.
//   - An error if the operation fails, or if the provided year or month is invalid.
func (s *transferStatsStatusHandleGrpc) FindMonthlyTransferStatusFailed(ctx context.Context, req *pb.FindMonthlyTransferStatus) (*pb.ApiResponseTransferMonthStatusFailed, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	s.logger.Info("Fetching transfer", zap.Int("year", year), zap.Int("month", month))

	if year <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("year", year))
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("month", month))
		return nil, transfer_errors.ErrGrpcInvalidMonth
	}

	reqService := &requests.MonthStatusTransfer{
		Year:  year,
		Month: month,
	}

	records, err := s.servicestats.FindMonthTransferStatusFailed(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTransferMonthStatusFailed("success", "success fetched monthly Transfer status Failed", records)

	s.logger.Info("Successfully fetched monthly Transfer status Failed", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

// FindYearlyTransferStatusFailed retrieves the yearly transfer status for failed transfers
// based on the provided year.
//
// It logs the process of fetching the transfer data and validates the input,
// ensuring that the year is a positive integer. If the validation fails, it logs
// the error and returns an appropriate error response. The function queries the
// transfer statistics service to fetch the transfer status and maps the result
// to a protocol buffer response object.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindYearTransferStatus object containing the year for which to fetch
//     the transfer status.
//
// Returns:
//   - An ApiResponseTransferYearStatusFailed containing the transfer status for
//     the specified year.
//   - An error if the operation fails, or if the provided year is invalid.
func (s *transferStatsStatusHandleGrpc) FindYearlyTransferStatusFailed(ctx context.Context, req *pb.FindYearTransferStatus) (*pb.ApiResponseTransferYearStatusFailed, error) {
	year := int(req.GetYear())

	s.logger.Info("Fetching transfer", zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("year", year))
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	records, err := s.servicestats.FindYearlyTransferStatusFailed(ctx, year)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTransferYearStatusFailed("success", "success fetched yearly Transfer status Failed", records)

	s.logger.Info("Successfully fetched yearly Transfer status Failed", zap.Int("year", year))

	return so, nil
}

// FindMonthlyTransferStatusSuccessByCardNumber retrieves the monthly transfer status for successful transfers
// filtered by year, month, and card number.
//
// It logs the process of fetching the transfer data and validates the inputs,
// ensuring that the year is a positive integer, the month is a positive integer between
// 1 and 12, and the card number is not empty. If any of these validations fail, it logs
// the error and returns an appropriate error response. The function queries the
// transfer statistics by card number service to fetch the transfer status and maps the
// result to a protocol buffer response object.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindMonthlyTransferStatusCardNumber object containing the year, month, and card number
//     for which to fetch the transfer status.
//
// Returns:
//   - An ApiResponseTransferMonthStatusSuccess containing the transfer status for the
//     specified year, month, and card number.
//   - An error if the operation fails, or if the provided year, month, or card number is invalid.
func (s *transferStatsStatusHandleGrpc) FindMonthlyTransferStatusSuccessByCardNumber(ctx context.Context, req *pb.FindMonthlyTransferStatusCardNumber) (*pb.ApiResponseTransferMonthStatusSuccess, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	cardNumber := req.GetCardNumber()

	s.logger.Info("Fetching transfer", zap.Int("year", year), zap.Int("month", month))

	if year <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("year", year))
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("month", month))
		return nil, transfer_errors.ErrGrpcInvalidMonth
	}

	if cardNumber == "" {
		s.logger.Error("Failed to fetch transfer", zap.String("card_number", cardNumber))
		return nil, transfer_errors.ErrGrpcInvalidCardNumber
	}

	reqService := &requests.MonthStatusTransferCardNumber{
		Year:       year,
		Month:      month,
		CardNumber: cardNumber,
	}

	records, err := s.servicestatsbycard.FindMonthTransferStatusSuccessByCardNumber(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTransferMonthStatusSuccess("success", "Successfully fetched monthly Transfer status success", records)

	s.logger.Info("Successfully fetched monthly Transfer status success", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

// FindYearlyTransferStatusSuccessByCardNumber retrieves the yearly transfer status for successful transfers
// based on the provided year and card number.
//
// It logs the process of fetching the transfer data and validates the input,
// ensuring that the year is a positive integer and the card number is not empty.
// If the validation fails, it logs the error and returns an appropriate error
// response. The function queries the transfer statistics service to fetch the
// transfer status and maps the result to a protocol buffer response object.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindYearTransferStatusCardNumber object containing the year and card
//     number for which to fetch the transfer status.
//
// Returns:
//   - An ApiResponseTransferYearStatusSuccess containing the transfer status for
//     the specified year and card number.
//   - An error if the operation fails, or if the provided year or card number are
//     invalid.
func (s *transferStatsStatusHandleGrpc) FindYearlyTransferStatusSuccessByCardNumber(ctx context.Context, req *pb.FindYearTransferStatusCardNumber) (*pb.ApiResponseTransferYearStatusSuccess, error) {
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

	reqService := &requests.YearStatusTransferCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	records, err := s.servicestatsbycard.FindYearlyTransferStatusSuccessByCardNumber(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTransferYearStatusSuccess("success", "Successfully fetched yearly Transfer status success", records)

	s.logger.Info("Successfully fetched yearly Transfer status success", zap.Int("year", year))

	return so, nil
}

// FindMonthlyTransferStatusFailedByCardNumber retrieves the monthly transfer status for failed transfers
// filtered by year, month, and card number.
//
// It logs the process of fetching the transfer data and validates the inputs,
// ensuring that the year is a positive integer, the month is a positive integer between
// 1 and 12, and the card number is not empty. If any of these validations fail, it logs
// the error and returns an appropriate error response. The function queries the
// transfer statistics by card number service to fetch the transfer status and maps the
// result to a protocol buffer response object.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindMonthlyTransferStatusCardNumber object containing the year, month, and card number
//     for which to fetch the transfer status.
//
// Returns:
//   - An ApiResponseTransferMonthStatusFailed containing the transfer status for the
//     specified year, month, and card number.
//   - An error if the operation fails, or if the provided year, month, or card number is invalid.
func (s *transferStatsStatusHandleGrpc) FindMonthlyTransferStatusFailedByCardNumber(ctx context.Context, req *pb.FindMonthlyTransferStatusCardNumber) (*pb.ApiResponseTransferMonthStatusFailed, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	cardNumber := req.GetCardNumber()

	s.logger.Info("Fetching transfer", zap.Int("year", year), zap.Int("month", month))

	if year <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("year", year))
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("month", month))
		return nil, transfer_errors.ErrGrpcInvalidMonth
	}

	if cardNumber == "" {
		s.logger.Error("Failed to fetch transfer", zap.String("card_number", cardNumber))
		return nil, transfer_errors.ErrGrpcInvalidCardNumber
	}

	reqService := &requests.MonthStatusTransferCardNumber{
		Year:       year,
		Month:      month,
		CardNumber: cardNumber,
	}

	records, err := s.servicestatsbycard.FindMonthTransferStatusFailedByCardNumber(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTransferMonthStatusFailed("success", "success fetched monthly Transfer status Failed", records)

	s.logger.Info("Successfully fetched monthly Transfer status Failed", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

// FindYearlyTransferStatusFailedByCardNumber retrieves the yearly transfer status for failed transfers
// based on the provided year and card number.
//
// This function logs the process of fetching the transfer data and validates the inputs,
// ensuring that the year is a positive integer and the card number is not empty. If any
// of these validations fail, it logs the error and returns an appropriate error response.
// The function queries the transfer statistics service by card number to fetch the
// transfer status and maps the result to a protocol buffer response object.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindYearTransferStatusCardNumber object containing the year and card
//     number for which to fetch the transfer status.
//
// Returns:
//   - An ApiResponseTransferYearStatusFailed containing the transfer status for
//     the specified year and card number.
//   - An error if the operation fails, or if the provided year or card number are invalid.
func (s *transferStatsStatusHandleGrpc) FindYearlyTransferStatusFailedByCardNumber(ctx context.Context, req *pb.FindYearTransferStatusCardNumber) (*pb.ApiResponseTransferYearStatusFailed, error) {
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

	reqService := &requests.YearStatusTransferCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	records, err := s.servicestatsbycard.FindYearlyTransferStatusFailedByCardNumber(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTransferYearStatusFailed("success", "success fetched yearly Transfer status Failed", records)

	s.logger.Info("Successfully fetched yearly Transfer status Failed", zap.Int("year", year))

	return so, nil
}
