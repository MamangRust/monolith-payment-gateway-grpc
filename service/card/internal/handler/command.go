package handler

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-card/internal/service"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/card"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/grpc"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/card"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type cardCommandService struct {
	pb.UnimplementedCardCommandServiceServer

	cardCommand service.CardCommandService

	logger logger.LoggerInterface

	mapper protomapper.CardCommandProtoMapper
}

func NewCardCommandHandleGrpc(cardCommand service.CardCommandService, logger logger.LoggerInterface, mapper protomapper.CardCommandProtoMapper) CardCommandService {
	return &cardCommandService{
		cardCommand: cardCommand,
		logger:      logger,
		mapper:      mapper,
	}
}

// CreateCard creates a new card for a user.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A CreateCardRequest object containing the details of the card to be created.
//
// Returns:
//   - A pointer to an ApiResponseCard object containing the created card's info.
//   - An error if the operation fails.
func (s *cardCommandService) CreateCard(ctx context.Context, req *pb.CreateCardRequest) (*pb.ApiResponseCard, error) {
	request := requests.CreateCardRequest{
		UserID:       int(req.UserId),
		CardType:     req.CardType,
		ExpireDate:   req.ExpireDate.AsTime(),
		CVV:          req.Cvv,
		CardProvider: req.CardProvider,
	}

	if err := request.Validate(); err != nil {
		s.logger.Error("CreateCard failed", zap.Any("error", err))
		return nil, card_errors.ErrGrpcValidateCreateCardRequest
	}

	res, err := s.cardCommand.CreateCard(ctx, &request)

	if err != nil {
		s.logger.Error("CreateCard failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseCard("success", "Successfully created card", res)

	s.logger.Info("Successfully created card", zap.Bool("success", true))

	return so, nil
}

// UpdateCard updates an existing card with the specified details.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: An UpdateCardRequest object containing the details to update the card with.
//
// Returns:
//   - A pointer to an ApiResponseCard representing the updated card.
//   - An error if the operation fails or the request validation fails.
func (s *cardCommandService) UpdateCard(ctx context.Context, req *pb.UpdateCardRequest) (*pb.ApiResponseCard, error) {
	request := requests.UpdateCardRequest{
		CardID:       int(req.CardId),
		UserID:       int(req.UserId),
		CardType:     req.CardType,
		ExpireDate:   req.ExpireDate.AsTime(),
		CVV:          req.Cvv,
		CardProvider: req.CardProvider,
	}

	if err := request.Validate(); err != nil {
		s.logger.Error("UpdateCard failed", zap.Any("error", err))
		return nil, card_errors.ErrGrpcValidateUpdateCardRequest
	}

	res, err := s.cardCommand.UpdateCard(ctx, &request)

	if err != nil {
		s.logger.Error("UpdateCard failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseCard("success", "Successfully updated card", res)

	s.logger.Info("Successfully updated card", zap.Bool("success", true))

	return so, nil
}

// TrashedCard marks a card as trashed.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindByIdCardRequest object containing the ID of the card to be trashed.
//
// Returns:
//   - An ApiResponseCard containing the trashed card record.
//   - An error if the operation fails.
func (s *cardCommandService) TrashedCard(ctx context.Context, req *pb.FindByIdCardRequest) (*pb.ApiResponseCardDeleteAt, error) {
	id := int(req.GetCardId())

	s.logger.Info("Trashing card", zap.Int("cardId", id))

	if id == 0 {
		s.logger.Error("TrashedCard failed", zap.Any("error", card_errors.ErrGrpcInvalidCardID))
		return nil, card_errors.ErrGrpcInvalidCardID
	}

	res, err := s.cardCommand.TrashedCard(ctx, id)

	if err != nil {
		s.logger.Error("TrashedCard failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseCardDeleteAt("success", "Successfully trashed card", res)

	s.logger.Info("Successfully trashed card", zap.Bool("success", true))

	return so, nil
}

// RestoreCard restores a trashed card.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindByIdCardRequest object containing the ID of the card to be restored.
//
// Returns:
//   - An ApiResponseCard containing the restored card record.
//   - An error if the operation fails.
func (s *cardCommandService) RestoreCard(ctx context.Context, req *pb.FindByIdCardRequest) (*pb.ApiResponseCard, error) {
	id := int(req.GetCardId())

	s.logger.Info("Restoring card", zap.Int("cardId", id))

	if id == 0 {
		s.logger.Error("RestoreCard failed", zap.Any("error", card_errors.ErrGrpcInvalidCardID))
		return nil, card_errors.ErrGrpcInvalidCardID
	}

	res, err := s.cardCommand.RestoreCard(ctx, id)

	if err != nil {
		s.logger.Error("RestoreCard failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseCard("success", "Successfully restored card", res)

	s.logger.Info("Successfully restored card", zap.Bool("success", true))

	return so, nil
}

// DeleteCardPermanent permanently deletes a card record from the database using its ID.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindByIdCardRequest object containing the ID of the card to be permanently deleted.
//
// Returns:
//   - An ApiResponseCardDelete indicating the result of the deletion operation.
//   - An error if the operation fails or if the provided card ID is invalid.
func (s *cardCommandService) DeleteCardPermanent(ctx context.Context, req *pb.FindByIdCardRequest) (*pb.ApiResponseCardDelete, error) {
	id := int(req.GetCardId())

	s.logger.Info("Deleting card", zap.Int("cardId", id))

	if id == 0 {
		s.logger.Error("DeleteCardPermanent failed", zap.Any("error", card_errors.ErrGrpcInvalidCardID))
		return nil, card_errors.ErrGrpcInvalidCardID
	}

	_, err := s.cardCommand.DeleteCardPermanent(ctx, id)

	if err != nil {
		s.logger.Error("DeleteCardPermanent failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseCardDelete("success", "Successfully deleted card")

	s.logger.Info("Successfully deleted card", zap.Bool("success", true))

	return so, nil
}

// RestoreAllCard restores all trashed cards.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: An empty request object.
//
// Returns:
//   - An ApiResponseCardAll containing the restored cards.
//   - An error if the operation fails.
func (s *cardCommandService) RestoreAllCard(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseCardAll, error) {
	s.logger.Info("Restoring all card")

	_, err := s.cardCommand.RestoreAllCard(ctx)

	if err != nil {
		s.logger.Error("RestoreAllCard failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseCardAll("success", "Successfully restore card")

	s.logger.Info("Successfully restore card", zap.Bool("success", true))

	return so, nil
}

// DeleteAllCardPermanent permanently deletes all trashed cards.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: An empty request object.
//
// Returns:
//   - An ApiResponseCardAll containing the deleted cards.
//   - An error if the operation fails.
func (s *cardCommandService) DeleteAllCardPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseCardAll, error) {
	s.logger.Info("Deleting all card")

	_, err := s.cardCommand.DeleteAllCardPermanent(ctx)

	if err != nil {
		s.logger.Info("DeleteAllCardPermanent failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseCardAll("success", "Successfully delete card permanent")

	s.logger.Info("Successfully delete card permanent", zap.Bool("success", true))

	return so, nil
}
