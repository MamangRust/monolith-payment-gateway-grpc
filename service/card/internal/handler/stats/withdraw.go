package handlerstats

import (
	"context"

	pb "github.com/MamangRust/monolith-payment-gateway-pb/card"

	cardstatsservice "github.com/MamangRust/monolith-payment-gateway-card/internal/service/stats"
	cardstatsbycard "github.com/MamangRust/monolith-payment-gateway-card/internal/service/statsbycard"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/grpc"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/card"
	"go.uber.org/zap"
)

type cardStatsWithdrawGrpc struct {
	pb.UnimplementedCardStatsWithdrawServiceServer

	cardStatsWithdraw cardstatsservice.CardStatsWithdrawService

	cardStatsWithdrawByCard cardstatsbycard.CardStatsWithdrawByCardService

	logger logger.LoggerInterface

	mapper protomapper.CardStatsAmountProtoMapper
}

func NewCardStatsWithdrawGrpc(cardStatsWithdraw cardstatsservice.CardStatsService, cardStatsWithdrawByCard cardstatsbycard.CardStatsByCardService, logger logger.LoggerInterface, mapper protomapper.CardStatsAmountProtoMapper) CardStatsWithdrawService {
	return &cardStatsWithdrawGrpc{
		cardStatsWithdraw:       cardStatsWithdraw,
		cardStatsWithdrawByCard: cardStatsWithdrawByCard,
		logger:                  logger,
		mapper:                  mapper,
	}
}

// FindMonthlyWithdrawAmount retrieves the monthly withdraw amount statistics for a given year.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindYearAmount object containing the year to fetch the monthly withdraw amount statistics for.
//
// Returns:
//   - An ApiResponseMonthlyAmount containing the monthly withdraw amount statistics retrieved from the database.
//   - An error if the operation fails, or if the provided year is invalid.
func (s *cardStatsWithdrawGrpc) FindMonthlyWithdrawAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseMonthlyAmount, error) {
	year := int(req.GetYear())

	s.logger.Info("Fetching monthly withdraw amount", zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("FindMonthlyWithdrawAmount failed", zap.Any("error", card_errors.ErrGrpcInvalidYear))
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardStatsWithdraw.FindMonthlyWithdrawAmount(ctx, year)

	if err != nil {
		s.logger.Error("FindMonthlyWithdrawAmount failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMonthlyAmounts("success", "Monthly withdraw amount retrieved successfully", res)

	s.logger.Info("Successfully fetched monthly withdraw amount", zap.Bool("success", true))

	return so, nil
}

// FindYearlyWithdrawAmount retrieves the yearly withdraw amount statistics for a given year.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindYearAmount object containing the year to fetch the yearly withdraw amount statistics for.
//
// Returns:
//   - An ApiResponseYearlyAmount containing the yearly withdraw amount statistics retrieved from the database.
//   - An error if the operation fails, or if the provided year is invalid.
func (s *cardStatsWithdrawGrpc) FindYearlyWithdrawAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseYearlyAmount, error) {
	year := int(req.GetYear())

	s.logger.Info("Fetching yearly withdraw amount", zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("FindYearlyWithdrawAmount failed", zap.Any("error", card_errors.ErrGrpcInvalidYear))
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardStatsWithdraw.FindYearlyWithdrawAmount(ctx, year)

	if err != nil {
		s.logger.Error("FindYearlyWithdrawAmount failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseYearlyAmounts("success", "Yearly withdraw amount retrieved successfully", res)

	s.logger.Info("Successfully fetched yearly withdraw amount", zap.Bool("success", true))

	return so, nil
}

// FindMonthlyWithdrawAmountByCardNumber retrieves the monthly withdraw amount statistics for a given card number and year.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindYearAmountCardNumber object containing the card number and year to fetch the monthly withdraw amount statistics for.
//
// Returns:
//   - An ApiResponseMonthlyAmount containing the monthly withdraw amount statistics retrieved from the database.
//   - An error if the operation fails, or if the provided year or card number is invalid.
func (s *cardStatsWithdrawGrpc) FindMonthlyWithdrawAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseMonthlyAmount, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	s.logger.Info("Fetching monthly withdraw amount by card number", zap.String("card_number", card_number), zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("FindMonthlyWithdrawAmountByCardNumber failed", zap.Any("error", card_errors.ErrGrpcInvalidYear))
		return nil, card_errors.ErrGrpcInvalidYear
	}

	if card_number == "" {
		s.logger.Error("FindMonthlyWithdrawAmountByCardNumber failed", zap.Any("error", card_errors.ErrGrpcInvalidCardNumber))
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardStatsWithdrawByCard.FindMonthlyWithdrawAmountByCardNumber(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthlyWithdrawAmountByCardNumber failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMonthlyAmounts("success", "Monthly withdraw amount by card number retrieved successfully", res)

	s.logger.Info("Successfully fetched monthly withdraw amount by card number", zap.Bool("success", true))

	return so, nil
}

// FindYearlyWithdrawAmountByCardNumber retrieves the yearly withdraw amount statistics for a given card number and year.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindYearAmountCardNumber object containing the card number and year to fetch the yearly withdraw amount statistics for.
//
// Returns:
//   - An ApiResponseYearlyAmount containing the yearly withdraw amount statistics retrieved from the database.
//   - An error if the operation fails, or if the provided year or card number is invalid.
func (s *cardStatsWithdrawGrpc) FindYearlyWithdrawAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseYearlyAmount, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	s.logger.Info("Fetching yearly withdraw amount by card number", zap.String("card_number", card_number), zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("FindYearlyWithdrawAmountByCardNumber failed", zap.Any("error", card_errors.ErrGrpcInvalidYear))
		return nil, card_errors.ErrGrpcInvalidYear
	}

	if card_number == "" {
		s.logger.Error("FindYearlyWithdrawAmountByCardNumber failed", zap.Any("error", card_errors.ErrGrpcInvalidCardNumber))
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardStatsWithdrawByCard.FindYearlyWithdrawAmountByCardNumber(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearlyWithdrawAmountByCardNumber failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseYearlyAmounts("success", "Yearly withdraw amount by card number retrieved successfully", res)

	s.logger.Info("Successfully fetched yearly withdraw amount by card number", zap.Bool("success", true))

	return so, nil
}
