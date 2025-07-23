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

type cardStatsTransactionGrpc struct {
	pb.UnimplementedCardStatsTransactonServiceServer

	cardStatsTransaction cardstatsservice.CardStatsTransactionService

	cardStatsTransactionByCard cardstatsbycard.CardStatsTransactionByCardService

	logger logger.LoggerInterface

	mapper protomapper.CardStatsAmountProtoMapper
}

func NewCardStatsTransactionGrpc(cardStatsTransaction cardstatsservice.CardStatsService, cardStatsTransactionByCard cardstatsbycard.CardStatsByCardService, logger logger.LoggerInterface, mapper protomapper.CardStatsAmountProtoMapper) CardStatsTransactionService {
	return &cardStatsTransactionGrpc{
		cardStatsTransaction:       cardStatsTransaction,
		cardStatsTransactionByCard: cardStatsTransactionByCard,
		logger:                     logger,
		mapper:                     mapper,
	}
}

// FindMonthlyTransactionAmount retrieves the monthly transaction amount statistics for a given year.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindYearAmount object containing the year to fetch the monthly transaction amount statistics for.
//
// Returns:
//   - An ApiResponseMonthlyAmount containing the monthly transaction amount statistics retrieved from the database.
//   - An error if the operation fails, or if the provided year is invalid.
func (s *cardStatsTransactionGrpc) FindMonthlyTransactionAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseMonthlyAmount, error) {
	year := int(req.GetYear())

	s.logger.Info("Fetching monthly transaction amount", zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("FindMonthlyTransactionAmount failed", zap.Any("error", card_errors.ErrGrpcInvalidYear))
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardStatsTransaction.FindMonthlyTransactionAmount(ctx, year)

	if err != nil {
		s.logger.Error("FindMonthlyTransactionAmount failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMonthlyAmounts("success", "Monthly transaction amount retrieved successfully", res)

	s.logger.Info("Successfully fetched monthly transaction amount", zap.Bool("success", true))

	return so, nil
}

// FindYearlyTransactionAmount retrieves the yearly transaction amount statistics for a given year.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindYearAmount object containing the year to fetch the yearly transaction amount statistics for.
//
// Returns:
//   - An ApiResponseYearlyAmount containing the yearly transaction amount statistics retrieved from the database.
//   - An error if the operation fails, or if the provided year is invalid.
func (s *cardStatsTransactionGrpc) FindYearlyTransactionAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseYearlyAmount, error) {
	year := int(req.GetYear())

	s.logger.Info("Fetching yearly transaction amount", zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("FindYearlyTransactionAmount failed", zap.Any("error", card_errors.ErrGrpcInvalidYear))
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardStatsTransaction.FindYearlyTransactionAmount(ctx, year)

	if err != nil {
		s.logger.Error("FindYearlyTransactionAmount failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseYearlyAmounts("success", "Yearly transaction amount retrieved successfully", res)

	s.logger.Info("Successfully fetched yearly transaction amount", zap.Bool("success", true))

	return so, nil
}

// FindMonthlyTransactionAmountByCardNumber finds the monthly transaction amount statistics for a given card number and year.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindYearAmountCardNumber object containing the card number and year to fetch the monthly transaction amount statistics for.
//
// Returns:
//   - An ApiResponseMonthlyAmount containing the monthly transaction amount statistics retrieved from the database.
//   - An error if the operation fails, or if the provided year or card number is invalid.
func (s *cardStatsTransactionGrpc) FindMonthlyTransactionAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseMonthlyAmount, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	s.logger.Info("Fetching monthly transaction amount by card number", zap.String("card_number", card_number), zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("FindMonthlyTransactionAmountByCardNumber failed", zap.Any("error", card_errors.ErrGrpcInvalidYear))
		return nil, card_errors.ErrGrpcInvalidYear
	}

	if card_number == "" {
		s.logger.Error("FindMonthlyTransactionAmountByCardNumber failed", zap.Any("error", card_errors.ErrGrpcInvalidCardNumber))
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardStatsTransactionByCard.FindMonthlyTransactionAmountByCardNumber(ctx, &reqService)

	if err != nil {
		s.logger.Error("FindMonthlyTransactionAmountByCardNumber failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMonthlyAmounts("success", "Monthly transaction amount by card number retrieved successfully", res)

	s.logger.Info("Successfully fetched monthly transaction amount by card number", zap.Bool("success", true))

	return so, nil
}

// FindYearlyTransactionAmountByCardNumber fetches the yearly transaction amount statistics for a given card number from the database.
//
// Parameters:
//   - ctx: The context.Context object for the gRPC request.
//   - req: A FindYearAmountCardNumber object containing the card number and year to fetch the yearly transaction amount statistics for.
//
// Returns:
//   - An ApiResponseYearlyAmount containing the yearly transaction amount statistics retrieved from the database.
//   - An error if the operation fails, or if the provided year or card number is invalid.
func (s *cardStatsTransactionGrpc) FindYearlyTransactionAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseYearlyAmount, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	s.logger.Info("Fetching yearly transaction amount by card number", zap.String("card_number", card_number), zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("FindYearlyTransactionAmountByCardNumber failed", zap.Any("error", card_errors.ErrGrpcInvalidYear))
		return nil, card_errors.ErrGrpcInvalidYear
	}

	if card_number == "" {
		s.logger.Error("FindYearlyTransactionAmountByCardNumber failed", zap.Any("error", card_errors.ErrGrpcInvalidCardNumber))
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardStatsTransactionByCard.FindYearlyTransactionAmountByCardNumber(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearlyTransactionAmountByCardNumber failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseYearlyAmounts("success", "Yearly transaction amount by card number retrieved successfully", res)

	s.logger.Info("Successfully fetched yearly transaction amount by card number", zap.Bool("success", true))

	return so, nil
}
