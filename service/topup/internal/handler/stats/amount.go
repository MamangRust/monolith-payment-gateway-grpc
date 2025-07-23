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

type topupStatsAmountHandleGrpc struct {
	pb.UnimplementedTopupStatsAmountServiceServer

	servicestats servicestats.TopupStatsService

	servicestatsbycard servicestatsbycard.TopupStatsByCardService

	logger logger.LoggerInterface

	mapper protomapper.TopupStatsAmountProtoMapper
}

func NewTopupStatsAmountHandleGrpc(
	service service.Service,
	logger logger.LoggerInterface,
	mapper protomapper.TopupStatsAmountProtoMapper,
) TopupStatsAmountHandleGrpc {
	return &topupStatsAmountHandleGrpc{
		servicestats:       service,
		servicestatsbycard: service,
		logger:             logger,
		mapper:             mapper,
	}
}

// FindMonthlyTopupAmounts fetches monthly topup amounts for a given year.
//
// Parameters:
//   - ctx: the context.Context object passed through the gRPC request.
//   - req: a pointer to a FindYearTopupStatus message containing the year.
//
// Returns:
//   - A pointer to an ApiResponseTopupMonthAmount message containing the amounts.
//   - An error, if the topup query service returns an error or if the year is invalid.
func (s *topupStatsAmountHandleGrpc) FindMonthlyTopupAmounts(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupMonthAmount, error) {
	year := int(req.GetYear())

	s.logger.Info("Fetching monthly topup amounts",
		zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("Failed to fetch monthly topup amounts", zap.Int("year", year))
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	amounts, err := s.servicestats.FindMonthlyTopupAmounts(ctx, year)

	if err != nil {
		s.logger.Error("Failed to fetch monthly topup amounts", zap.Any("error", err), zap.Int("year", year))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTopupMonthAmount("success", "Successfully fetched monthly topup amounts", amounts)

	s.logger.Info("Successfully fetched monthly topup amounts",
		zap.Int("year", year))

	return so, nil
}

// FindYearlyTopupAmounts fetches yearly topup amounts for a given year.
//
// Parameters:
//   - ctx: the context.Context object passed through the gRPC request.
//   - req: a pointer to a FindYearTopupStatus message containing the year.
//
// Returns:
//   - A pointer to an ApiResponseTopupYearAmount message containing the amounts.
//   - An error, if the topup query service returns an error or if the year is invalid.
func (s *topupStatsAmountHandleGrpc) FindYearlyTopupAmounts(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupYearAmount, error) {
	year := int(req.GetYear())

	s.logger.Info("Fetching yearly topup amounts",
		zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("Failed to fetch yearly topup amounts", zap.Int("year", year))
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	amounts, err := s.servicestats.FindYearlyTopupAmounts(ctx, year)

	if err != nil {
		s.logger.Error("Failed to fetch yearly topup amounts", zap.Any("error", err), zap.Int("year", year))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTopupYearAmount("success", "Successfully fetched yearly topup amounts", amounts)

	s.logger.Info("Successfully fetched yearly topup amounts",
		zap.Int("year", year))

	return so, nil
}

// FindMonthlyTopupAmountsByCardNumber fetches monthly topup amounts for a specific card number and year.
//
// Parameters:
//   - ctx: the context.Context object passed through the gRPC request.
//   - req: a pointer to a FindYearTopupCardNumber message containing the year and card number.
//
// Returns:
//   - A pointer to an ApiResponseTopupMonthAmount message containing the amounts.
//   - An error, if the topup query service returns an error or if the year or card number is invalid.
func (s *topupStatsAmountHandleGrpc) FindMonthlyTopupAmountsByCardNumber(ctx context.Context, req *pb.FindYearTopupCardNumber) (*pb.ApiResponseTopupMonthAmount, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	s.logger.Info("Fetching monthly topup amounts by card number",
		zap.Int("year", year),
		zap.String("card_number", cardNumber))

	if year <= 0 {
		s.logger.Error("Failed to fetch monthly topup amounts by card number", zap.Int("year", year))
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	if cardNumber == "" {
		s.logger.Error("Failed to fetch monthly topup amounts by card number", zap.String("card_number", cardNumber))
		return nil, topup_errors.ErrGrpcInvalidCardNumber
	}

	reqService := &requests.YearMonthMethod{
		Year:       year,
		CardNumber: cardNumber,
	}

	amounts, err := s.servicestatsbycard.FindMonthlyTopupAmountsByCardNumber(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to fetch monthly topup amounts by card number", zap.Any("error", err), zap.Int("year", year), zap.String("card_number", cardNumber))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTopupMonthAmount("success", "Successfully fetched monthly topup amounts by card number", amounts)

	s.logger.Info("Successfully fetched monthly topup amounts by card number", zap.Int("year", year), zap.String("card_number", cardNumber))

	return so, nil
}

// FindYearlyTopupAmountsByCardNumber fetches yearly topup amounts for a specific card number and year.
//
// Parameters:
//   - ctx: the context.Context object passed through the gRPC request.
//   - req: a pointer to a FindYearTopupCardNumber message containing the year and card number.
//
// Returns:
//   - A pointer to an ApiResponseTopupYearAmount message containing the amounts.
//   - An error, if the topup query service returns an error or if the year or card number is invalid.
func (s *topupStatsAmountHandleGrpc) FindYearlyTopupAmountsByCardNumber(ctx context.Context, req *pb.FindYearTopupCardNumber) (*pb.ApiResponseTopupYearAmount, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	s.logger.Info("Fetching yearly topup amounts by card number",
		zap.Int("year", year),
		zap.String("card_number", cardNumber))

	if year <= 0 {
		s.logger.Error("Failed to fetch yearly topup amounts by card number", zap.Int("year", year))
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	if cardNumber == "" {
		s.logger.Error("Failed to fetch yearly topup amounts by card number", zap.String("card_number", cardNumber))
		return nil, topup_errors.ErrGrpcInvalidCardNumber
	}

	reqService := &requests.YearMonthMethod{
		Year:       year,
		CardNumber: cardNumber,
	}

	amounts, err := s.servicestatsbycard.FindYearlyTopupAmountsByCardNumber(ctx, reqService)
	if err != nil {
		s.logger.Error("Failed to fetch yearly topup amounts by card number", zap.Any("error", err), zap.Int("year", year), zap.String("card_number", cardNumber))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTopupYearAmount("success", "Successfully fetched yearly topup amounts by card number", amounts)

	s.logger.Info("Successfully fetched yearly topup amounts by card number",
		zap.Int("year", year),
		zap.String("card_number", cardNumber))

	return so, nil
}
