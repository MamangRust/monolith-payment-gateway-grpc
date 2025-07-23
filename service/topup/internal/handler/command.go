package handler

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
	"google.golang.org/protobuf/types/known/emptypb"
)

type topupCommandHandleGrpc struct {
	pb.UnimplementedTopupCommandServiceServer

	service service.TopupCommandService

	mapper protomapper.TopupCommandProtoMapper

	logger logger.LoggerInterface
}

func NewTopupCommandHandleGrpc(service service.TopupCommandService, mapper protomapper.TopupCommandProtoMapper, logger logger.LoggerInterface) TopupCommandHandleGrpc {
	return &topupCommandHandleGrpc{
		service: service,
		mapper:  mapper,
		logger:  logger,
	}
}

// CreateTopup is a gRPC handler function that creates a new topup record.
//
// Parameters:
//   - ctx: the context.Context object passed through the gRPC request.
//   - req: a pointer to a CreateTopupRequest message, containing the card number, topup amount, and topup method.
//
// Returns:
//   - A pointer to an ApiResponseTopup message, containing the created topup record.
//   - An error, if the topup command service returns an error.
func (s *topupCommandHandleGrpc) CreateTopup(ctx context.Context, req *pb.CreateTopupRequest) (*pb.ApiResponseTopup, error) {
	request := &requests.CreateTopupRequest{
		CardNumber:  req.GetCardNumber(),
		TopupAmount: int(req.GetTopupAmount()),
		TopupMethod: req.GetTopupMethod(),
	}

	s.logger.Debug("Creating topup",
		zap.String("card_number", request.CardNumber),
		zap.Int("topup_amount", request.TopupAmount),
		zap.String("topup_method", request.TopupMethod))

	if err := request.Validate(); err != nil {
		s.logger.Error("Failed to create topup", zap.Any("error", err))
		return nil, topup_errors.ErrGrpcValidateCreateTopup
	}

	res, err := s.service.CreateTopup(ctx, request)

	if err != nil {
		s.logger.Error("Failed to create topup", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTopup("success", "Successfully created topup", res)

	s.logger.Info("Successfully created topup", zap.Bool("success", true))

	return so, nil
}

// UpdateTopup updates an existing topup record with the provided details.
//
// Parameters:
//   - ctx: the context.Context object passed through the gRPC request.
//   - req: a pointer to a UpdateTopupRequest message, containing the topup ID,
//     card number, topup amount, and topup method.
//
// Returns:
//   - A pointer to an ApiResponseTopup message, containing the updated topup record.
//   - An error, if the topup command service returns an error.
func (s *topupCommandHandleGrpc) UpdateTopup(ctx context.Context, req *pb.UpdateTopupRequest) (*pb.ApiResponseTopup, error) {
	id := int(req.GetTopupId())

	s.logger.Info("Updating topup",
		zap.Int("topup_id", id),
		zap.String("card_number", req.GetCardNumber()),
		zap.Int("topup_amount", int(req.GetTopupAmount())),
		zap.String("topup_method", req.GetTopupMethod()))

	if id == 0 {
		s.logger.Error("Failed to update topup", zap.Int("topup_id", id))
		return nil, topup_errors.ErrGrpcTopupInvalidID
	}

	request := &requests.UpdateTopupRequest{
		TopupID:     &id,
		CardNumber:  req.GetCardNumber(),
		TopupAmount: int(req.GetTopupAmount()),
		TopupMethod: req.GetTopupMethod(),
	}

	if err := request.Validate(); err != nil {
		s.logger.Error("Failed to update topup", zap.Any("error", err))
		return nil, topup_errors.ErrGrpcValidateUpdateTopup
	}

	res, err := s.service.UpdateTopup(ctx, request)

	if err != nil {
		s.logger.Error("Failed to update topup", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTopup("success", "Successfully updated topup", res)

	s.logger.Info("Successfully updated topup", zap.Bool("success", true))

	return so, nil
}

// TrashedTopup marks a top-up record as trashed by its ID.
//
// Parameters:
//   - ctx: the context.Context object passed through the gRPC request.
//   - req: a pointer to a FindByIdTopupRequest message containing the top-up ID.
//
// Returns:
//   - A pointer to an ApiResponseTopupDeleteAt message, indicating the top-up has been trashed.
//   - An error, if the top-up command service returns an error or if the ID is invalid.
func (s *topupCommandHandleGrpc) TrashedTopup(ctx context.Context, req *pb.FindByIdTopupRequest) (*pb.ApiResponseTopupDeleteAt, error) {
	id := int(req.GetTopupId())

	s.logger.Info("Trashing topup",
		zap.Int("topup.id", id))

	if id == 0 {
		s.logger.Error("Failed to trash topup", zap.Int("topup.id", id))
		return nil, topup_errors.ErrGrpcTopupInvalidID
	}

	res, err := s.service.TrashedTopup(ctx, id)

	if err != nil {
		s.logger.Error("Failed to trash topup", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTopupDeleteAt("success", "Successfully trashed topup", res)

	s.logger.Info("Successfully trashed topup", zap.Bool("success", true))

	return so, nil
}

// RestoreTopup is a gRPC handler function that restores a trashed top-up record.
//
// Parameters:
//   - ctx: the context.Context object passed through the gRPC request.
//   - req: a pointer to a FindByIdTopupRequest message containing the top-up ID.
//
// Returns:
//   - A pointer to an ApiResponseTopupDeleteAt message, indicating the top-up has been restored.
//   - An error, if the top-up command service returns an error or if the ID is invalid.
func (s *topupCommandHandleGrpc) RestoreTopup(ctx context.Context, req *pb.FindByIdTopupRequest) (*pb.ApiResponseTopup, error) {
	id := int(req.GetTopupId())

	s.logger.Info("Restoring topup",
		zap.Int("topup.id", id))

	if id == 0 {
		s.logger.Error("Failed to restore topup", zap.Int("topup.id", id))
		return nil, topup_errors.ErrGrpcTopupInvalidID
	}

	res, err := s.service.RestoreTopup(ctx, id)

	if err != nil {
		s.logger.Error("Failed to restore topup", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTopup("success", "Successfully restored topup", res)

	s.logger.Info("Successfully restored topup", zap.Bool("success", true))

	return so, nil
}

// DeleteTopupPermanent is a gRPC handler function that permanently deletes a top-up record by its ID.
//
// Parameters:
//   - ctx: the context.Context object passed through the gRPC request.
//   - req: a pointer to a FindByIdTopupRequest message containing the top-up ID.
//
// Returns:
//   - A pointer to an ApiResponseTopupDelete message, indicating the top-up has been deleted.
//   - An error, if the top-up command service returns an error or if the ID is invalid.
func (s *topupCommandHandleGrpc) DeleteTopupPermanent(ctx context.Context, req *pb.FindByIdTopupRequest) (*pb.ApiResponseTopupDelete, error) {
	id := int(req.GetTopupId())

	s.logger.Info("Deleting topup permanently",
		zap.Int("topup.id", id))

	if id == 0 {
		s.logger.Error("Failed to delete topup permanently", zap.Int("topup.id", id))
		return nil, topup_errors.ErrGrpcTopupInvalidID
	}

	_, err := s.service.DeleteTopupPermanent(ctx, id)

	if err != nil {
		s.logger.Error("Failed to delete topup permanently", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTopupDelete("success", "Successfully deleted topup permanently")

	s.logger.Info("Successfully deleted topup permanently", zap.Bool("success", true))

	return so, nil
}

// RestoreAllTopup is a gRPC handler function that restores all trashed top-up records in the database.
//
// Parameters:
//   - ctx: the context.Context object passed through the gRPC request.
//   - req: an empty protobuf message.
//
// Returns:
//   - A pointer to an ApiResponseTopupAll message, indicating all top-up records have been restored.
//   - An error, if the top-up command service returns an error.
func (s *topupCommandHandleGrpc) RestoreAllTopup(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTopupAll, error) {
	s.logger.Info("Restoring all topup")

	_, err := s.service.RestoreAllTopup(ctx)

	if err != nil {
		s.logger.Error("Failed to restore all topup", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTopupAll("success", "Successfully restore all topup")

	s.logger.Info("Successfully restored all topup", zap.Bool("success", true))

	return so, nil
}

// DeleteAllTopupPermanent permanently deletes all topup records from the database.
//
// Returns:
//   - A pointer to an ApiResponseTopupAll message, indicating all topup records have been deleted.
//   - An error, if the topup command service returns an error.
func (s *topupCommandHandleGrpc) DeleteAllTopupPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTopupAll, error) {
	s.logger.Info("Deleting all topup permanently")

	_, err := s.service.DeleteAllTopupPermanent(ctx)

	if err != nil {
		s.logger.Error("Failed to delete all topup permanently", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTopupAll("success", "Successfully delete topup permanent")

	s.logger.Info("Successfully deleted all topup permanently", zap.Bool("success", true))

	return so, nil
}
