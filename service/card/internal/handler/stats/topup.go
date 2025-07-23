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

type cardStatsTopupGrpc struct {
	pb.UnimplementedCardStatsTopupServiceServer

	cardStatsTopup cardstatsservice.CardStatsTopupService

	cardStatsTopupByCard cardstatsbycard.CardStatsTopupByCardService

	logger logger.LoggerInterface

	mapper protomapper.CardStatsAmountProtoMapper
}

func NewCardStatsTopupGrpc(cardStatsTopup cardstatsservice.CardStatsService, cardStatsTopupByCard cardstatsbycard.CardStatsByCardService, logger logger.LoggerInterface, mapper protomapper.CardStatsAmountProtoMapper) CardStatsTopupService {
	return &cardStatsTopupGrpc{
		cardStatsTopup:       cardStatsTopup,
		cardStatsTopupByCard: cardStatsTopupByCard,
		logger:               logger,
		mapper:               mapper,
	}
}

// FindMonthlyTopupAmount retrieves the monthly topup amount statistics for a given year.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindYearAmount object containing the year to fetch the monthly topup amount statistics for.
//
// Returns:
//   - An ApiResponseMonthlyAmount containing the monthly topup amount statistics retrieved from the database.
//   - An error if the operation fails, or if the provided year is invalid.
func (s *cardStatsTopupGrpc) FindMonthlyTopupAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseMonthlyAmount, error) {
	year := int(req.GetYear())

	s.logger.Info("Fetching monthly topup amount", zap.Int("year", year))

	if year <= 0 {
		s.logger.Info("FindMonthlyTopupAmount failed", zap.Any("error", card_errors.ErrGrpcInvalidYear))
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardStatsTopup.FindMonthlyTopupAmount(ctx, year)

	if err != nil {
		s.logger.Error("FindMonthlyTopupAmount failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMonthlyAmounts("success", "Monthly topup amount retrieved successfully", res)

	s.logger.Info("Successfully fetched monthly topup amount", zap.Bool("success", true))

	return so, nil
}

// FindYearlyTopupAmount retrieves the yearly topup amount statistics for a given year.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindYearAmount object containing the year to fetch the yearly topup amount statistics for.
//
// Returns:
//   - An ApiResponseYearlyAmount containing the yearly topup amount statistics retrieved from the database.
//   - An error if the operation fails, or if the provided year is invalid.
func (s *cardStatsTopupGrpc) FindYearlyTopupAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseYearlyAmount, error) {
	year := int(req.GetYear())

	s.logger.Info("Fetching yearly topup amount", zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("FindYearlyTopupAmount failed", zap.Any("error", card_errors.ErrGrpcInvalidYear))
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardStatsTopup.FindYearlyTopupAmount(ctx, year)

	if err != nil {
		s.logger.Error("FindYearlyTopupAmount failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseYearlyAmounts("success", "Yearly topup amount retrieved successfully", res)

	s.logger.Info("Successfully fetched yearly topup amount", zap.Bool("success", true))

	return so, nil
}

// FindMonthlyTopupAmountByCardNumber retrieves the monthly topup amount statistics for a given card number and year.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindYearAmountCardNumber object containing the card number and year to fetch the monthly topup amount statistics for.
//
// Returns:
//   - An ApiResponseMonthlyAmount containing the monthly topup amount statistics retrieved from the database.
//   - An error if the operation fails, or if the provided year or card number is invalid.
func (s *cardStatsTopupGrpc) FindMonthlyTopupAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseMonthlyAmount, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	s.logger.Info("Fetching monthly topup amount by card number", zap.String("card_number", card_number), zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("FindMonthlyTopupAmountByCardNumber failed", zap.Any("error", card_errors.ErrGrpcInvalidYear))
		return nil, card_errors.ErrGrpcInvalidYear
	}

	if card_number == "" {
		s.logger.Error("FindMonthlyTopupAmountByCardNumber failed", zap.Any("error", card_errors.ErrGrpcInvalidCardNumber))
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardStatsTopupByCard.FindMonthlyTopupAmountByCardNumber(ctx, &reqService)

	if err != nil {
		s.logger.Error("FindMonthlyTopupAmountByCardNumber failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMonthlyAmounts("success", "Monthly topup amount by card number retrieved successfully", res)

	s.logger.Info("Successfully fetched monthly topup amount by card number", zap.Bool("success", true))

	return so, nil
}

// FindYearlyTopupAmountByCardNumber retrieves the yearly topup amount statistics for a given card number and year.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindYearAmountCardNumber object containing the card number and year to fetch the yearly topup amount statistics for.
//
// Returns:
//   - An ApiResponseYearlyAmount containing the yearly topup amount statistics retrieved from the database.
//   - An error if the operation fails, or if the provided year or card number is invalid.
func (s *cardStatsTopupGrpc) FindYearlyTopupAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseYearlyAmount, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	s.logger.Info("Fetching yearly topup amount by card number", zap.String("card_number", card_number), zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("FindYearlyTopupAmountByCardNumber failed", zap.Any("error", card_errors.ErrGrpcInvalidYear))
		return nil, card_errors.ErrGrpcInvalidYear
	}

	if card_number == "" {
		s.logger.Error("FindYearlyTopupAmountByCardNumber failed", zap.Any("error", card_errors.ErrGrpcInvalidCardNumber))
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       int(year),
	}

	res, err := s.cardStatsTopupByCard.FindYearlyTopupAmountByCardNumber(ctx, &reqService)

	if err != nil {
		s.logger.Error("FindYearlyTopupAmountByCardNumber failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseYearlyAmounts("success", "Yearly topup amount by card number retrieved successfully", res)

	s.logger.Info("Successfully fetched yearly topup amount by card number", zap.Bool("success", true))

	return so, nil
}
