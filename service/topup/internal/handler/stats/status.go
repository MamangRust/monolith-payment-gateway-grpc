package topupstatshandler

import (
	"context"

	pb "github.com/MamangRust/monolith-payment-gateway-pb/topup"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/grpc"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/topup"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/service"
	"go.uber.org/zap"

	servicestats "github.com/MamangRust/monolith-payment-gateway-topup/internal/service/stats"
	servicestatsbycard "github.com/MamangRust/monolith-payment-gateway-topup/internal/service/statsbycard"
)

type topupStatusHandleGrpc struct {
	pb.UnimplementedTopupStatsStatusServiceServer

	servicestats servicestats.TopupStatsService

	servicestatsbycard servicestatsbycard.TopupStatsByCardService

	logger logger.LoggerInterface

	mapper protomapper.TopupStatsStatusProtoMapper
}

func NewTopupStatsStatusHandleGrpc(
	service service.Service,
	logger logger.LoggerInterface,
	mapper protomapper.TopupStatsStatusProtoMapper,
) TopupStatsStatusHandleGrpc {
	return &topupStatusHandleGrpc{
		servicestats:       service,
		servicestatsbycard: service,
		logger:             logger,
		mapper:             mapper,
	}
}

// FindMonthlyTopupStatusSuccess fetches monthly topup status success records.
//
// Parameters:
//   - ctx: the context.Context object passed through the gRPC request.
//   - req: a pointer to a FindMonthlyTopupStatus message containing the year and month.
//
// Returns:
//   - A pointer to an ApiResponseTopupMonthStatusSuccess message containing the records.
//   - An error, if the topup query service returns an error or if the year or month is invalid.
func (s *topupStatusHandleGrpc) FindMonthlyTopupStatusSuccess(ctx context.Context, req *pb.FindMonthlyTopupStatus) (*pb.ApiResponseTopupMonthStatusSuccess, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	s.logger.Info("Fetching monthly topup status success",
		zap.Int("year", year),
		zap.Int("month", month))

	if year <= 0 {
		s.logger.Error("Failed to fetch monthly topup status success", zap.Int("year", year))
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	if month <= 0 {
		s.logger.Error("Failed to fetch monthly topup status success", zap.Int("month", month))
		return nil, topup_errors.ErrGrpcTopupInvalidMonth
	}

	reqService := &requests.MonthTopupStatus{
		Year:  year,
		Month: month,
	}

	records, err := s.servicestats.FindMonthTopupStatusSuccess(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to fetch monthly topup status success", zap.Any("error", err), zap.Int("year", year), zap.Int("month", month))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTopupMonthStatusSuccess("success", "Successfully fetched monthly topup status success", records)

	s.logger.Info("Successfully fetched monthly topup status success",
		zap.Int("year", year),
		zap.Int("month", month))

	return so, nil
}

// FindYearlyTopupStatusSuccess fetches yearly topup status success records.
//
// Parameters:
//   - ctx: the context.Context object passed through the gRPC request.
//   - req: a pointer to a FindYearTopupStatus message containing the year.
//
// Returns:
//   - A pointer to an ApiResponseTopupYearStatusSuccess message containing the records.
//   - An error, if the topup query service returns an error or if the year is invalid.
func (s *topupStatusHandleGrpc) FindYearlyTopupStatusSuccess(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupYearStatusSuccess, error) {
	year := int(req.GetYear())

	s.logger.Info("Fetching yearly topup status success",
		zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("Failed to fetch yearly topup status success", zap.Int("year", year))
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	records, err := s.servicestats.FindYearlyTopupStatusSuccess(ctx, year)

	if err != nil {
		s.logger.Error("Failed to fetch yearly topup status success", zap.Any("error", err), zap.Int("year", year))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTopupYearStatusSuccess("success", "Successfully fetched yearly topup status success", records)

	s.logger.Info("Successfully fetched yearly topup status success", zap.Int("year", year))

	return so, nil
}

// FindMonthlyTopupStatusFailed fetches monthly topup status failed records.
//
// Parameters:
//   - ctx: the context.Context object passed through the gRPC request.
//   - req: a pointer to a FindMonthlyTopupStatus message containing the year and month.
//
// Returns:
//   - A pointer to an ApiResponseTopupMonthStatusFailed message containing the records.
//   - An error, if the topup query service returns an error or if the year or month is invalid.
func (s *topupStatusHandleGrpc) FindMonthlyTopupStatusFailed(ctx context.Context, req *pb.FindMonthlyTopupStatus) (*pb.ApiResponseTopupMonthStatusFailed, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	s.logger.Info("Fetching monthly topup status Failed",
		zap.Int("year", year),
		zap.Int("month", month))

	if year <= 0 {
		s.logger.Error("Failed to fetch monthly topup status Failed", zap.Int("year", year))
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	if month <= 0 {
		s.logger.Error("Failed to fetch monthly topup status Failed", zap.Int("month", month))
		return nil, topup_errors.ErrGrpcTopupInvalidMonth
	}

	reqService := &requests.MonthTopupStatus{
		Year:  year,
		Month: month,
	}

	records, err := s.servicestats.FindMonthTopupStatusFailed(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to fetch monthly topup status Failed", zap.Any("error", err), zap.Int("year", year), zap.Int("month", month))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTopupMonthStatusFailed("Successfully", "Successfully fetched monthly topup status Failed", records)

	s.logger.Info("Successfully fetched monthly topup status Failed",
		zap.Int("year", year),
		zap.Int("month", month))

	return so, nil
}

// FindYearlyTopupStatusFailed fetches yearly topup status failed records.
//
// Parameters:
//   - ctx: the context.Context object passed through the gRPC request.
//   - req: a pointer to a FindYearTopupStatus message containing the year.
//
// Returns:
//   - A pointer to an ApiResponseTopupYearStatusFailed message containing the records.
//   - An error, if the topup query service returns an error or if the year is invalid.
func (s *topupStatusHandleGrpc) FindYearlyTopupStatusFailed(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupYearStatusFailed, error) {
	year := int(req.GetYear())

	s.logger.Info("Fetching yearly topup status Failed",
		zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("Failed to fetch yearly topup status Failed", zap.Int("year", year))
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	records, err := s.servicestats.FindYearlyTopupStatusFailed(ctx, year)

	if err != nil {
		s.logger.Error("Failed to fetch yearly topup status Failed", zap.Any("error", err), zap.Int("year", year))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTopupYearStatusFailed("Successfully", "Successfully fetched yearly topup status Failed", records)

	s.logger.Info("Successfully fetched yearly topup status Failed", zap.Int("year", year))

	return so, nil
}

// FindMonthlyTopupStatusSuccessByCardNumber fetches monthly topup status success records by card number.
//
// Parameters:
//   - ctx: the context.Context object passed through the gRPC request.
//   - req: a pointer to a FindMonthlyTopupStatusCardNumber message containing the year, month, and card number.
//
// Returns:
//   - A pointer to an ApiResponseTopupMonthStatusSuccess message containing the records.
//   - An error, if the topup query service returns an error or if the year, month, or card number is invalid.
func (s *topupStatusHandleGrpc) FindMonthlyTopupStatusSuccessByCardNumber(ctx context.Context, req *pb.FindMonthlyTopupStatusCardNumber) (*pb.ApiResponseTopupMonthStatusSuccess, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	cardNumber := req.GetCardNumber()

	s.logger.Info("Fetching monthly topup status success",
		zap.Int("year", year),
		zap.Int("month", month),
		zap.String("card_number", cardNumber))

	if year <= 0 {
		s.logger.Error("Failed to fetch monthly topup status success", zap.Int("year", year))
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	if month <= 0 {
		s.logger.Error("Failed to fetch monthly topup status success", zap.Int("month", month))
		return nil, topup_errors.ErrGrpcTopupInvalidMonth
	}

	if cardNumber == "" {
		s.logger.Error("Failed to fetch monthly topup status success", zap.String("card_number", cardNumber))
		return nil, topup_errors.ErrGrpcInvalidCardNumber
	}

	reqService := &requests.MonthTopupStatusCardNumber{
		Year:       year,
		Month:      month,
		CardNumber: cardNumber,
	}

	records, err := s.servicestatsbycard.FindMonthTopupStatusSuccessByCardNumber(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to fetch monthly topup status success", zap.Any("error", err), zap.Int("year", year), zap.Int("month", month), zap.String("card_number", cardNumber))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTopupMonthStatusSuccess("success", "Successfully fetched monthly topup status success", records)

	s.logger.Info("Successfully fetched monthly topup status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", cardNumber))

	return so, nil
}

// FindYearlyTopupStatusSuccessByCardNumber fetches yearly topup status success records by card number.
//
// Parameters:
//   - ctx: the context.Context object passed through the gRPC request.
//   - req: a pointer to a FindYearTopupStatusCardNumber message containing the year and card number.
//
// Returns:
//   - A pointer to an ApiResponseTopupYearStatusSuccess message containing the records.
//   - An error, if the topup query service returns an error or if the year or card number is invalid.
func (s *topupStatusHandleGrpc) FindYearlyTopupStatusSuccessByCardNumber(ctx context.Context, req *pb.FindYearTopupStatusCardNumber) (*pb.ApiResponseTopupYearStatusSuccess, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	s.logger.Info("Fetching yearly topup status success",
		zap.Int("year", year),
		zap.String("card_number", cardNumber))

	if year <= 0 {
		s.logger.Error("Failed to fetch yearly topup status success", zap.Int("year", year))
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	if cardNumber == "" {
		s.logger.Error("Failed to fetch yearly topup status success", zap.String("card_number", cardNumber))
		return nil, topup_errors.ErrGrpcInvalidCardNumber
	}

	reqService := &requests.YearTopupStatusCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	records, err := s.servicestatsbycard.FindYearlyTopupStatusSuccessByCardNumber(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to fetch yearly topup status success", zap.Any("error", err), zap.Int("year", year), zap.String("card_number", cardNumber))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTopupYearStatusSuccess("success", "Successfully fetched yearly topup status success", records)

	s.logger.Info("Successfully fetched yearly topup status success", zap.Int("year", year), zap.String("card_number", cardNumber))

	return so, nil
}

// FindMonthlyTopupStatusFailedByCardNumber fetches monthly topup status failed records by card number.
//
// Parameters:
//   - ctx: the context.Context object passed through the gRPC request.
//   - req: a pointer to a FindMonthlyTopupStatusCardNumber message containing the year, month, and card number.
//
// Returns:
//   - A pointer to an ApiResponseTopupMonthStatusFailed message containing the records.
//   - An error, if the topup query service returns an error or if the year, month, or card number is invalid.
func (s *topupStatusHandleGrpc) FindMonthlyTopupStatusFailedByCardNumber(ctx context.Context, req *pb.FindMonthlyTopupStatusCardNumber) (*pb.ApiResponseTopupMonthStatusFailed, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	cardNumber := req.GetCardNumber()

	s.logger.Info("Fetching monthly topup status failed",
		zap.Int("year", year),
		zap.Int("month", month),
		zap.String("card_number", cardNumber))

	if year <= 0 {
		s.logger.Error("Failed to fetch monthly topup status failed", zap.Int("year", year))
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	if month <= 0 {
		s.logger.Error("Failed to fetch monthly topup status failed", zap.Int("month", month))
		return nil, topup_errors.ErrGrpcTopupInvalidMonth
	}

	if cardNumber == "" {
		s.logger.Error("Failed to fetch monthly topup status failed", zap.String("card_number", cardNumber))
		return nil, topup_errors.ErrGrpcInvalidCardNumber
	}

	reqService := &requests.MonthTopupStatusCardNumber{
		Year:       year,
		Month:      month,
		CardNumber: cardNumber,
	}

	records, err := s.servicestatsbycard.FindMonthTopupStatusFailedByCardNumber(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to fetch monthly topup status failed", zap.Any("error", err), zap.Int("year", year), zap.Int("month", month), zap.String("card_number", cardNumber))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTopupMonthStatusFailed("success", "Successfully fetched monthly topup status failed", records)

	s.logger.Info("Successfully fetched monthly topup status failed",
		zap.Int("year", year),
		zap.Int("month", month),
		zap.String("card_number", cardNumber))

	return so, nil
}

// FindYearlyTopupStatusFailedByCardNumber fetches yearly topup status failed records by card number.
//
// Parameters:
//   - ctx: the context.Context object passed through the gRPC request.
//   - req: a pointer to a FindYearTopupStatusCardNumber message containing the year and card number.
//
// Returns:
//   - A pointer to an ApiResponseTopupYearStatusFailed message containing the records.
//   - An error, if the topup query service returns an error or if the year or card number is invalid.
func (s *topupStatusHandleGrpc) FindYearlyTopupStatusFailedByCardNumber(ctx context.Context, req *pb.FindYearTopupStatusCardNumber) (*pb.ApiResponseTopupYearStatusFailed, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	s.logger.Info("Fetching yearly topup status failed",
		zap.Int("year", year),
		zap.String("card_number", cardNumber))

	if year <= 0 {
		s.logger.Error("Failed to fetch yearly topup status failed", zap.Int("year", year))
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	if cardNumber == "" {
		s.logger.Error("Failed to fetch yearly topup status failed", zap.String("card_number", cardNumber))
		return nil, topup_errors.ErrGrpcInvalidCardNumber
	}

	reqService := &requests.YearTopupStatusCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	records, err := s.servicestatsbycard.FindYearlyTopupStatusFailedByCardNumber(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to fetch yearly topup status failed", zap.Any("error", err), zap.Int("year", year), zap.String("card_number", cardNumber))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTopupYearStatusFailed("success", "Successfully fetched yearly topup status failed", records)

	s.logger.Info("Successfully fetched yearly topup status failed", zap.Int("year", year), zap.String("card_number", cardNumber))

	return so, nil
}
