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

type cardStatsTransferGrpc struct {
	pb.UnimplementedCardStatsTransferServiceServer

	cardStatsTransfer cardstatsservice.CardStatsTransferService

	cardStatsTransferByCard cardstatsbycard.CardStatsTransferByCardService

	logger logger.LoggerInterface

	mapper protomapper.CardStatsAmountProtoMapper
}

func NewCardStatsTransferGrpc(cardStatsTransfer cardstatsservice.CardStatsService, cardStatsTransferByCard cardstatsbycard.CardStatsByCardService, logger logger.LoggerInterface, mapper protomapper.CardStatsAmountProtoMapper) CardStatsTransferService {
	return &cardStatsTransferGrpc{
		cardStatsTransfer:       cardStatsTransfer,
		cardStatsTransferByCard: cardStatsTransferByCard,
		logger:                  logger,
		mapper:                  mapper,
	}
}

// FindMonthlyTransferSenderAmount retrieves the monthly transfer sender amount statistics for a given year.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindYearAmount object containing the year to fetch the monthly transfer sender amount statistics for.
//
// Returns:
//   - An ApiResponseMonthlyAmount containing the monthly transfer sender amount statistics retrieved from the database.
//   - An error if the operation fails, or if the provided year is invalid.
func (s *cardStatsTransferGrpc) FindMonthlyTransferSenderAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseMonthlyAmount, error) {
	year := int(req.GetYear())

	s.logger.Info("Fetching monthly transfer sender amount", zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("FindMonthlyTransferSenderAmount failed", zap.Any("error", card_errors.ErrGrpcInvalidYear))
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardStatsTransfer.FindMonthlyTransferAmountSender(ctx, year)

	if err != nil {
		s.logger.Error("FindMonthlyTransferSenderAmount failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMonthlyAmounts("success", "Monthly transfer sender amount retrieved successfully", res)

	s.logger.Info("Successfully fetched monthly transfer sender amount", zap.Bool("success", true))

	return so, nil
}

// FindYearlyTransferSenderAmount retrieves the yearly transfer sender amount statistics for a given year.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindYearAmount object containing the year to fetch the yearly transfer sender amount statistics for.
//
// Returns:
//   - An ApiResponseYearlyAmount containing the yearly transfer sender amount statistics retrieved from the database.
//   - An error if the operation fails, or if the provided year is invalid.
func (s *cardStatsTransferGrpc) FindYearlyTransferSenderAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseYearlyAmount, error) {
	year := int(req.GetYear())

	s.logger.Info("Fetching yearly transfer sender amount", zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("FindYearlyTransferSenderAmount failed", zap.Any("error", card_errors.ErrGrpcInvalidYear))
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardStatsTransfer.FindYearlyTransferAmountSender(ctx, year)

	if err != nil {
		s.logger.Error("FindYearlyTransferSenderAmount failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseYearlyAmounts("success", "transfer sender amount retrieved successfully", res)

	s.logger.Info("Successfully fetched yearly transfer sender amount", zap.Bool("success", true))

	return so, nil
}

// FindMonthlyTransferReceiverAmount retrieves the monthly transfer receiver amount statistics for a given year.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindYearAmount object containing the year to fetch the monthly transfer receiver amount statistics for.
//
// Returns:
//   - An ApiResponseMonthlyAmount containing the monthly transfer receiver amount statistics retrieved from the database.
//   - An error if the operation fails, or if the provided year is invalid.
func (s *cardStatsTransferGrpc) FindMonthlyTransferReceiverAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseMonthlyAmount, error) {
	year := int(req.GetYear())

	s.logger.Info("Fetching monthly transfer receiver amount", zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("FindMonthlyTransferReceiverAmount failed", zap.Any("error", card_errors.ErrGrpcInvalidYear))
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardStatsTransfer.FindMonthlyTransferAmountReceiver(ctx, year)

	if err != nil {
		s.logger.Error("FindMonthlyTransferReceiverAmount failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMonthlyAmounts("success", "Monthly transfer receiver amount retrieved successfully", res)

	s.logger.Info("Successfully fetched monthly transfer receiver amount", zap.Bool("success", true))

	return so, nil
}

// FindYearlyTransferReceiverAmount retrieves the yearly transfer receiver amount statistics for a given year.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindYearAmount object containing the year to fetch the yearly transfer receiver amount statistics for.
//
// Returns:
//   - An ApiResponseYearlyAmount containing the yearly transfer receiver amount statistics retrieved from the database.
//   - An error if the operation fails, or if the provided year is invalid.
func (s *cardStatsTransferGrpc) FindYearlyTransferReceiverAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseYearlyAmount, error) {
	year := int(req.GetYear())

	s.logger.Info("Fetching yearly transfer receiver amount", zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("FindYearlyTransferReceiverAmount failed", zap.Any("error", card_errors.ErrGrpcInvalidYear))
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardStatsTransfer.FindYearlyTransferAmountReceiver(ctx, year)

	if err != nil {
		s.logger.Error("FindYearlyTransferReceiverAmount failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseYearlyAmounts("success", "Yearly transfer receiver amount retrieved successfully", res)

	s.logger.Info("Successfully fetched yearly transfer receiver amount", zap.Bool("success", true))

	return so, nil
}

// FindMonthlyTransferSenderAmountByCardNumber fetches the monthly transfer sender amount statistics for a given card number from the database.
//
// Parameters:
//   - ctx: The context.Context object for the gRPC request.
//   - req: A FindYearAmountCardNumber object containing the card number and year to fetch the monthly transfer sender amount statistics for.
//
// Returns:
//   - An ApiResponseMonthlyAmount containing the monthly transfer sender amount statistics retrieved from the database.
//   - An error if the operation fails, or if the provided year or card number is invalid.
func (s *cardStatsTransferGrpc) FindMonthlyTransferSenderAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseMonthlyAmount, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	s.logger.Info("Fetching monthly transfer sender amount by card number", zap.String("card_number", card_number), zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("FindMonthlyTransferSenderAmountByCardNumber failed", zap.Any("error", card_errors.ErrGrpcInvalidYear))
		return nil, card_errors.ErrGrpcInvalidYear
	}

	if card_number == "" {
		s.logger.Error("FindMonthlyTransferSenderAmountByCardNumber failed", zap.Any("error", card_errors.ErrGrpcInvalidCardNumber))
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardStatsTransferByCard.FindMonthlyTransferAmountBySender(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthlyTransferSenderAmountByCardNumber failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMonthlyAmounts("success", "Monthly transfer sender amount by card number retrieved successfully", res)

	s.logger.Info("Successfully fetched monthly transfer sender amount by card number", zap.Bool("success", true))

	return so, nil
}

// FindYearlyTransferSenderAmountByCardNumber retrieves the yearly transfer sender amount statistics for a given card number
// and year.
//
// Parameters:
//   - ctx: The context.Context object for the gRPC request.
//   - req: A FindYearAmountCardNumber object containing the card number and year to fetch the yearly transfer sender amount
//     statistics for.
//
// Returns:
//   - An ApiResponseYearlyAmount containing the yearly transfer sender amount statistics retrieved from the database.
//   - An error if the operation fails, or if the provided year or card number is invalid.
func (s *cardStatsTransferGrpc) FindYearlyTransferSenderAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseYearlyAmount, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	if year <= 0 {
		s.logger.Error("FindYearlyTransferSenderAmountByCardNumber failed", zap.Any("error", card_errors.ErrGrpcInvalidYear))
		return nil, card_errors.ErrGrpcInvalidYear
	}

	if card_number == "" {
		s.logger.Error("FindYearlyTransferSenderAmountByCardNumber failed", zap.Any("error", card_errors.ErrGrpcInvalidCardNumber))
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardStatsTransferByCard.FindYearlyTransferAmountBySender(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearlyTransferSenderAmountByCardNumber failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseYearlyAmounts("success", "Yearly transfer sender amount by card number retrieved successfully", res)

	s.logger.Info("Successfully fetched yearly transfer sender amount by card number", zap.Bool("success", true))

	return so, nil
}

// FindMonthlyTransferReceiverAmountByCardNumber retrieves the monthly transfer receiver amount statistics for a given
// card number and year.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindYearAmountCardNumber object containing the card number and year to fetch the monthly transfer receiver
//     amount statistics for.
//
// Returns:
//   - An ApiResponseMonthlyAmount containing the monthly transfer receiver amount statistics retrieved from the database.
//   - An error if the operation fails, or if the provided year or card number is invalid.
func (s *cardStatsTransferGrpc) FindMonthlyTransferReceiverAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseMonthlyAmount, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	if year <= 0 {
		s.logger.Error("FindMonthlyTransferReceiverAmountByCardNumber failed", zap.Any("error", card_errors.ErrGrpcInvalidYear))
		return nil, card_errors.ErrGrpcInvalidYear
	}

	if card_number == "" {
		s.logger.Error("FindMonthlyTransferReceiverAmountByCardNumber failed", zap.Any("error", card_errors.ErrGrpcInvalidCardNumber))
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardStatsTransferByCard.FindMonthlyTransferAmountByReceiver(ctx, &reqService)

	if err != nil {
		s.logger.Error("FindMonthlyTransferReceiverAmountByCardNumber failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMonthlyAmounts("success", "Monthly transfer receiver amount by card number retrieved successfully", res)

	s.logger.Info("Successfully fetched monthly transfer receiver amount by card number", zap.Bool("success", true))

	return so, nil
}

// FindYearlyTransferReceiverAmountByCardNumber retrieves the yearly transfer receiver amount statistics for a given card number
// and year.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindYearAmountCardNumber object containing the card number and year to fetch the yearly transfer receiver amount
//     statistics for.
//
// Returns:
//   - An ApiResponseYearlyAmount containing the yearly transfer receiver amount statistics retrieved from the database.
//   - An error if the operation fails, or if the provided year or card number is invalid.
func (s *cardStatsTransferGrpc) FindYearlyTransferReceiverAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseYearlyAmount, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	s.logger.Info("Fetching yearly transfer receiver amount by card number", zap.String("card_number", card_number), zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("FindYearlyTransferReceiverAmountByCardNumber failed", zap.Any("error", card_errors.ErrGrpcInvalidYear))
		return nil, card_errors.ErrGrpcInvalidYear
	}

	if card_number == "" {
		s.logger.Error("FindYearlyTransferReceiverAmountByCardNumber failed", zap.Any("error", card_errors.ErrGrpcInvalidCardNumber))
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardStatsTransferByCard.FindYearlyTransferAmountByReceiver(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearlyTransferReceiverAmountByCardNumber failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseYearlyAmounts("success", "Yearly transfer receiver amount by card number retrieved successfully", res)

	s.logger.Info("Successfully fetched yearly transfer receiver amount by card number", zap.Bool("success", true))

	return so, nil
}
