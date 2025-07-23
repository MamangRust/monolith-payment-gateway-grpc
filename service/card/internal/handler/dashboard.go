package handler

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-card/internal/service"
	pbhelpers "github.com/MamangRust/monolith-payment-gateway-pb"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/card"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/grpc"

	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/card"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type cardDashboardHandleGrpc struct {
	pb.UnimplementedCardDashboardServiceServer

	cardDashboard service.CardDashboardService

	logger logger.LoggerInterface

	mapper protomapper.CardDashboardProtoMapper
}

func NewCardDashboardHandleGrpc(cardDashboard service.CardDashboardService, logger logger.LoggerInterface, mapper protomapper.CardDashboardProtoMapper) CardDashboardService {
	return &cardDashboardHandleGrpc{
		cardDashboard: cardDashboard,
		logger:        logger,
		mapper:        mapper,
	}
}

// DashboardCard retrieves dashboard card statistics.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: An empty request object.
//
// Returns:
//   - An ApiResponseDashboardCard containing the dashboard card statistics retrieved from the
//     database.
//   - An error if the operation fails.
func (s *cardDashboardHandleGrpc) DashboardCard(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseDashboardCard, error) {
	s.logger.Info("Fetching dashboard card")

	dashboardCard, err := s.cardDashboard.DashboardCard(ctx)
	if err != nil {
		s.logger.Error("DashboardCard failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseDashboardCard("success", "Dashboard card retrieved successfully", dashboardCard)

	s.logger.Info("Successfully fetched dashboard card", zap.Bool("success", true))

	return so, nil
}

// DashboardCardNumber retrieves dashboard card statistics for a specific card number.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindByCardNumberRequest object containing the card number to fetch the dashboard card statistics for.
//
// Returns:
//   - An ApiResponseDashboardCardNumber containing the dashboard card statistics retrieved from the database.
//   - An error if the operation fails, or if the provided card number is invalid.
func (s *cardDashboardHandleGrpc) DashboardCardNumber(ctx context.Context, req *pbhelpers.FindByCardNumberRequest) (*pb.ApiResponseDashboardCardNumber, error) {
	card_number := req.GetCardNumber()

	s.logger.Info("Fetching dashboard card for card number", zap.String("card.card_number", card_number))

	if card_number == "" {
		s.logger.Error("DashboardCardNumber failed", zap.Any("error", card_errors.ErrGrpcInvalidCardNumber))
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	dashboardCard, err := s.cardDashboard.DashboardCardCardNumber(ctx, card_number)

	if err != nil {
		s.logger.Error("DashboardCardNumber failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseDashboardCardCardNumber("success", "Dashboard card for card number retrieved successfully", dashboardCard)

	s.logger.Info("Successfully fetched dashboard card for card number", zap.Bool("success", true))

	return so, nil
}
