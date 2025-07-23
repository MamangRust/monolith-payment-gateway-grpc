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

type topupMethodHandleGrpc struct {
	pb.UnimplementedTopupStatsMethodServiceServer

	servicestats servicestats.TopupStatsService

	servicestatsbycard servicestatsbycard.TopupStatsByCardService

	logger logger.LoggerInterface

	mapper protomapper.TopupStatsMethodProtoMapper
}

func NewTopupStatsMethodHandleGrpc(
	service service.Service,
	logger logger.LoggerInterface,
	mapper protomapper.TopupStatsMethodProtoMapper,
) TopupStatsMethodHandleGrpc {
	return &topupMethodHandleGrpc{
		servicestats:       service,
		servicestatsbycard: service,
		logger:             logger,
		mapper:             mapper,
	}
}

// FindMonthlyTopupMethods fetches monthly topup methods for a given year.
//
// Parameters:
//   - ctx: the context.Context object passed through the gRPC request.
//   - req: a pointer to a FindYearTopupStatus message containing the year.
//
// Returns:
//   - A pointer to an ApiResponseTopupMonthMethod message containing the methods.
//   - An error, if the topup query service returns an error or if the year is invalid.
func (s *topupMethodHandleGrpc) FindMonthlyTopupMethods(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupMonthMethod, error) {
	year := int(req.GetYear())

	s.logger.Info("Fetching monthly topup methods",
		zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("Failed to fetch monthly topup methods", zap.Int("year", year))
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	methods, err := s.servicestats.FindMonthlyTopupMethods(ctx, year)

	if err != nil {
		s.logger.Error("Failed to fetch monthly topup methods", zap.Any("error", err), zap.Int("year", year))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTopupMonthMethod("success", "Successfully fetched monthly topup methods", methods)

	s.logger.Info("Successfully fetched monthly topup methods",
		zap.Int("year", year))

	return so, nil
}

// FindYearlyTopupMethods fetches yearly topup methods for a given year.
//
// Parameters:
//   - ctx: the context.Context object passed through the gRPC request.
//   - req: a pointer to a FindYearTopupStatus message containing the year.
//
// Returns:
//   - A pointer to an ApiResponseTopupYearMethod message containing the methods.
//   - An error, if the topup query service returns an error or if the year is invalid.
func (s *topupMethodHandleGrpc) FindYearlyTopupMethods(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupYearMethod, error) {
	year := int(req.GetYear())

	s.logger.Info("Fetching yearly topup methods",
		zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("Failed to fetch yearly topup methods", zap.Int("year", year))
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	methods, err := s.servicestats.FindYearlyTopupMethods(ctx, year)

	if err != nil {
		s.logger.Error("Failed to fetch yearly topup methods", zap.Any("error", err), zap.Int("year", year))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTopupYearMethod("success", "Successfully fetched yearly topup methods", methods)

	s.logger.Info("Successfully fetched yearly topup methods",
		zap.Int("year", year))

	return so, nil
}

// FindMonthlyTopupMethodsByCardNumber fetches monthly topup methods for a specific card number and year.
//
// Parameters:
//   - ctx: the context.Context object passed through the gRPC request.
//   - req: a pointer to a FindYearTopupCardNumber message containing the year and card number.
//
// Returns:
//   - A pointer to an ApiResponseTopupMonthMethod message containing the methods.
//   - An error, if the topup query service returns an error or if the year or card number is invalid.
func (s *topupMethodHandleGrpc) FindMonthlyTopupMethodsByCardNumber(ctx context.Context, req *pb.FindYearTopupCardNumber) (*pb.ApiResponseTopupMonthMethod, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	s.logger.Info("Fetching monthly topup methods by card number",
		zap.Int("year", year),
		zap.String("card_number", cardNumber))

	if year <= 0 {
		s.logger.Error("Failed to fetch monthly topup methods by card number", zap.Int("year", year))
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	if cardNumber == "" {
		s.logger.Error("Failed to fetch monthly topup methods by card number", zap.String("card_number", cardNumber))
		return nil, topup_errors.ErrGrpcInvalidCardNumber
	}

	reqService := &requests.YearMonthMethod{
		Year:       year,
		CardNumber: cardNumber,
	}

	methods, err := s.servicestatsbycard.FindMonthlyTopupMethodsByCardNumber(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to fetch monthly topup methods by card number", zap.Any("error", err), zap.Int("year", year), zap.String("card_number", cardNumber))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTopupMonthMethod("success", "Successfully fetched monthly topup methods by card number", methods)

	s.logger.Info("Successfully fetched monthly topup methods by card number",
		zap.Int("year", year),
		zap.String("card_number", cardNumber))

	return so, nil
}

// FindYearlyTopupMethodsByCardNumber fetches yearly topup methods for a specific card number and year.
//
// Parameters:
//   - ctx: the context.Context object passed through the gRPC request.
//   - req: a pointer to a FindYearTopupCardNumber message containing the year and card number.
//
// Returns:
//   - A pointer to an ApiResponseTopupYearMethod message containing the methods.
//   - An error, if the topup query service returns an error or if the year or card number is invalid.
func (s *topupMethodHandleGrpc) FindYearlyTopupMethodsByCardNumber(ctx context.Context, req *pb.FindYearTopupCardNumber) (*pb.ApiResponseTopupYearMethod, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	s.logger.Info("Fetching yearly topup methods by card number",
		zap.Int("year", year),
		zap.String("card_number", cardNumber))

	if year <= 0 {
		s.logger.Error("Failed to fetch yearly topup methods by card number", zap.Int("year", year))
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	if cardNumber == "" {
		s.logger.Error("Failed to fetch yearly topup methods by card number", zap.String("card_number", cardNumber))
		return nil, topup_errors.ErrGrpcInvalidCardNumber
	}

	reqService := &requests.YearMonthMethod{
		Year:       year,
		CardNumber: cardNumber,
	}

	methods, err := s.servicestatsbycard.FindYearlyTopupMethodsByCardNumber(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to fetch yearly topup methods by card number", zap.Any("error", err), zap.Int("year", year), zap.String("card_number", cardNumber))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTopupYearMethod("success", "Successfully fetched yearly topup methods by card number", methods)

	s.logger.Info("Successfully fetched yearly topup methods by card number",
		zap.Int("year", year),
		zap.String("card_number", cardNumber))

	return so, nil
}
