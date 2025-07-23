package handlerstats

import (
	"context"

	cardstatsservice "github.com/MamangRust/monolith-payment-gateway-card/internal/service/stats"
	cardstatsbycard "github.com/MamangRust/monolith-payment-gateway-card/internal/service/statsbycard"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/card"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/grpc"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/card"
	"go.uber.org/zap"
)

type cardStatsBalanceGrpc struct {
	pb.UnimplementedCardStatsBalanceServiceServer

	cardStatsBalance cardstatsservice.CardStatsBalanceService

	cardStatsBalanceByCard cardstatsbycard.CardStatsBalanceByCardService

	logger logger.LoggerInterface

	mapper protomapper.CardStatsBalanceProtoMapper
}

func NewCardStatsBalanceGrpc(cardStatsBalance cardstatsservice.CardStatsService, cardStatsBalanceByCard cardstatsbycard.CardStatsByCardService, logger logger.LoggerInterface, mapper protomapper.CardStatsBalanceProtoMapper) CardStatsBalanceService {
	return &cardStatsBalanceGrpc{
		cardStatsBalance:       cardStatsBalance,
		cardStatsBalanceByCard: cardStatsBalanceByCard,
		logger:                 logger,
		mapper:                 mapper,
	}
}

// FindMonthlyBalance retrieves the monthly balance statistics for a given year.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindYearBalance object containing the year to fetch the monthly balance statistics for.
//
// Returns:
//   - An ApiResponseMonthlyBalance containing the monthly balance statistics retrieved from the database.
//   - An error if the operation fails, or if the provided year is invalid.
func (s *cardStatsBalanceGrpc) FindMonthlyBalance(ctx context.Context, req *pb.FindYearBalance) (*pb.ApiResponseMonthlyBalance, error) {
	year := int(req.GetYear())

	s.logger.Info("Fetching monthly balance", zap.Int("year", year))

	if year <= 0 {
		s.logger.Info("FindMonthlyBalance failed", zap.Any("error", card_errors.ErrGrpcInvalidYear))
		return nil, card_errors.ErrGrpcInvalidYear
	}
	res, err := s.cardStatsBalance.FindMonthlyBalance(ctx, year)

	if err != nil {
		s.logger.Error("FindMonthlyBalance failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMonthlyBalances("success", "Monthly balance retrieved successfully", res)

	s.logger.Info("Successfully fetched monthly balance", zap.Bool("success", true))

	return so, nil
}

// FindYearlyBalance retrieves the yearly balance statistics for a given year.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindYearBalance object containing the year to fetch the yearly balance statistics for.
//
// Returns:
//   - An ApiResponseYearlyBalance containing the yearly balance statistics retrieved from the database.
//   - An error if the operation fails, or if the provided year is invalid.
func (s *cardStatsBalanceGrpc) FindYearlyBalance(ctx context.Context, req *pb.FindYearBalance) (*pb.ApiResponseYearlyBalance, error) {
	year := int(req.GetYear())

	s.logger.Info("Fetching yearly balance", zap.Int("year", year))

	if year <= 0 {
		s.logger.Info("FindYearlyBalance failed", zap.Any("error", card_errors.ErrGrpcInvalidYear))
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardStatsBalance.FindYearlyBalance(ctx, year)
	if err != nil {
		s.logger.Error("FindYearlyBalance failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseYearlyBalances("success", "Yearly balance retrieved successfully", res)

	s.logger.Info("Successfully fetched yearly balance", zap.Bool("success", true))

	return so, nil
}

// FindMonthlyBalanceByCardNumber retrieves the monthly balance statistics for a given card number and year.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindYearBalanceCardNumber object containing the card number and year to fetch the monthly balance statistics for.
//
// Returns:
//   - An ApiResponseMonthlyBalance containing the monthly balance statistics retrieved from the database.
//   - An error if the operation fails, or if the provided year or card number is invalid.
func (s *cardStatsBalanceGrpc) FindMonthlyBalanceByCardNumber(ctx context.Context, req *pb.FindYearBalanceCardNumber) (*pb.ApiResponseMonthlyBalance, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	s.logger.Info("Fetching monthly balance by card number", zap.String("card_number", card_number), zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("FindMonthlyBalanceByCardNumber failed", zap.Any("error", card_errors.ErrGrpcInvalidYear))
		return nil, card_errors.ErrGrpcInvalidYear
	}

	if card_number == "" {
		s.logger.Error("FindMonthlyBalanceByCardNumber failed", zap.Any("error", card_errors.ErrGrpcInvalidCardNumber))
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardStatsBalanceByCard.FindMonthlyBalanceByCardNumber(ctx, &reqService)

	if err != nil {
		s.logger.Error("FindMonthlyBalanceByCardNumber failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMonthlyBalances("success", "Monthly balance retrieved successfully", res)

	s.logger.Info("Successfully fetched monthly balance by card number", zap.Bool("success", true))

	return so, nil
}

// FindYearlyBalanceByCardNumber retrieves the yearly balance statistics for a given card number and year.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindYearBalanceCardNumber object containing the card number and year to fetch the yearly balance statistics for.
//
// Returns:
//   - An ApiResponseYearlyBalance containing the yearly balance statistics retrieved from the database.
//   - An error if the operation fails, or if the provided year or card number is invalid.
func (s *cardStatsBalanceGrpc) FindYearlyBalanceByCardNumber(ctx context.Context, req *pb.FindYearBalanceCardNumber) (*pb.ApiResponseYearlyBalance, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	s.logger.Info("Fetching yearly balance by card number", zap.String("card_number", card_number), zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("FindYearlyBalanceByCardNumber failed", zap.Any("error", card_errors.ErrGrpcInvalidYear))
		return nil, card_errors.ErrGrpcInvalidYear
	}

	if card_number == "" {
		s.logger.Error("FindYearlyBalanceByCardNumber failed", zap.Any("error", card_errors.ErrGrpcInvalidCardNumber))
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardStatsBalanceByCard.FindYearlyBalanceByCardNumber(ctx, &reqService)

	if err != nil {
		s.logger.Error("FindYearlyBalanceByCardNumber failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseYearlyBalances("success", "Yearly balance retrieved successfully", res)

	s.logger.Info("Successfully fetched yearly balance by card number", zap.Bool("success", true))

	return so, nil
}
