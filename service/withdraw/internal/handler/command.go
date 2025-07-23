package handler

import (
	"context"

	pb "github.com/MamangRust/monolith-payment-gateway-pb/withdraw"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	withdraw_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors/grpc"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/withdraw"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/service"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type withdrawCommandHandleGrpc struct {
	pb.UnimplementedWithdrawCommandServiceServer

	withdrawCommand service.WithdrawCommandService

	logger logger.LoggerInterface

	mapper protomapper.WithdrawCommandProtoMapper
}

func NewWithdrawCommandHandleGrpc(
	withdrawCommand service.WithdrawCommandService,
	logger logger.LoggerInterface,
	mapper protomapper.WithdrawCommandProtoMapper,
) WithdrawCommandHandlerGrpc {
	return &withdrawCommandHandleGrpc{
		withdrawCommand: withdrawCommand,
		logger:          logger,
		mapper:          mapper,
	}
}

// CreateWithdraw creates a new withdraw record in the database.
//
// Parameters:
//   - ctx: the context.Context object for request-scoped values, cancellation, and deadlines.
//   - req: a pointer to a CreateWithdrawRequest message containing the withdraw details.
//
// Returns:
//   - A pointer to an ApiResponseWithdraw message containing the newly created withdraw record.
//   - An error, which is non-nil if the operation fails.
func (w *withdrawCommandHandleGrpc) CreateWithdraw(ctx context.Context, req *pb.CreateWithdrawRequest) (*pb.ApiResponseWithdraw, error) {
	request := &requests.CreateWithdrawRequest{
		CardNumber:     req.CardNumber,
		WithdrawAmount: int(req.WithdrawAmount),
		WithdrawTime:   req.WithdrawTime.AsTime(),
	}

	w.logger.Info("CreateWithdraw", zap.Any("request", request))

	if err := request.Validate(); err != nil {
		w.logger.Error("CreateWithdraw", zap.Any("error", err))
		return nil, withdraw_errors.ErrGrpcValidateCreateWithdrawRequest
	}

	withdraw, err := w.withdrawCommand.Create(ctx, request)

	if err != nil {
		w.logger.Error("CreateWithdraw", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := w.mapper.ToProtoResponseWithdraw("success", "Successfully created withdraw", withdraw)

	return so, nil

}

// UpdateWithdraw updates a withdraw record in the database.
//
// Parameters:
//   - ctx: the context.Context object for request-scoped values, cancellation, and deadlines.
//   - req: a pointer to a UpdateWithdrawRequest message containing the withdraw details.
//
// Returns:
//   - A pointer to an ApiResponseWithdraw message containing the updated withdraw record.
//   - An error, which is non-nil if the operation fails.
func (w *withdrawCommandHandleGrpc) UpdateWithdraw(ctx context.Context, req *pb.UpdateWithdrawRequest) (*pb.ApiResponseWithdraw, error) {
	id := int(req.GetWithdrawId())

	w.logger.Info("UpdateWithdraw", zap.Int("id", id))

	if id == 0 {
		w.logger.Error("UpdateWithdraw", zap.Int("id", id))
		return nil, withdraw_errors.ErrGrpcWithdrawInvalidID
	}

	request := &requests.UpdateWithdrawRequest{
		WithdrawID:     &id,
		CardNumber:     req.CardNumber,
		WithdrawAmount: int(req.WithdrawAmount),
		WithdrawTime:   req.WithdrawTime.AsTime(),
	}

	w.logger.Info("UpdateWithdraw", zap.Any("request", request))

	if err := request.Validate(); err != nil {
		w.logger.Error("UpdateWithdraw", zap.Any("error", err))
		return nil, withdraw_errors.ErrGrpcValidateUpdateWithdrawRequest
	}

	withdraw, err := w.withdrawCommand.Update(ctx, request)

	if err != nil {
		w.logger.Error("UpdateWithdraw", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := w.mapper.ToProtoResponseWithdraw("success", "Successfully updated withdraw", withdraw)

	return so, nil
}

// TrashedWithdraw marks a withdrawal record as trashed based on the provided request,
// which includes the withdraw ID.
//
// Parameters:
//   - ctx: the context.Context object for tracing and cancellation.
//   - req: a pointer to a FindByIdWithdrawRequest message containing the withdraw ID.
//
// Returns:
//   - A pointer to an ApiResponseWithdraw message indicating the result of the operation.
//   - An error, which is non-nil if the operation fails.
func (w *withdrawCommandHandleGrpc) TrashedWithdraw(ctx context.Context, req *pb.FindByIdWithdrawRequest) (*pb.ApiResponseWithdrawDeleteAt, error) {
	id := int(req.GetWithdrawId())

	w.logger.Info("TrashedWithdraw", zap.Int("id", id))

	if id == 0 {
		w.logger.Error("TrashedWithdraw", zap.Int("id", id))
		return nil, withdraw_errors.ErrGrpcWithdrawInvalidID
	}

	withdraw, err := w.withdrawCommand.TrashedWithdraw(ctx, id)

	if err != nil {
		w.logger.Error("TrashedWithdraw", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := w.mapper.ToProtoResponseWithdrawDeleteAt("success", "Successfully trashed withdraw", withdraw)

	return so, nil
}

// RestoreWithdraw restores a previously trashed withdrawal record based on the provided request,
// which includes the withdraw ID.
//
// Parameters:
//   - ctx: the context.Context object for tracing and cancellation.
//   - req: a pointer to a FindByIdWithdrawRequest message containing the withdraw ID.
//
// Returns:
//   - A pointer to an ApiResponseWithdraw message indicating the result of the operation.
//   - An error, which is non-nil if the operation fails.
func (w *withdrawCommandHandleGrpc) RestoreWithdraw(ctx context.Context, req *pb.FindByIdWithdrawRequest) (*pb.ApiResponseWithdraw, error) {
	id := int(req.GetWithdrawId())

	w.logger.Info("RestoreWithdraw", zap.Int("id", id))

	if id == 0 {
		w.logger.Error("RestoreWithdraw", zap.Int("id", id))
		return nil, withdraw_errors.ErrGrpcWithdrawInvalidID
	}

	withdraw, err := w.withdrawCommand.RestoreWithdraw(ctx, id)

	if err != nil {
		w.logger.Error("RestoreWithdraw", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := w.mapper.ToProtoResponseWithdraw("success", "Successfully restored withdraw", withdraw)

	return so, nil
}

// DeleteWithdrawPermanent permanently deletes a withdrawal record from the database by its ID.
//
// Parameters:
//   - ctx: the context.Context object for tracing and cancellation.
//   - req: a pointer to a FindByIdWithdrawRequest message containing the withdraw ID.
//
// Returns:
//   - A pointer to an ApiResponseWithdrawDelete message indicating the result of the operation.
//   - An error, which is non-nil if the operation fails.
func (w *withdrawCommandHandleGrpc) DeleteWithdrawPermanent(ctx context.Context, req *pb.FindByIdWithdrawRequest) (*pb.ApiResponseWithdrawDelete, error) {
	id := int(req.GetWithdrawId())

	w.logger.Info("DeleteWithdrawPermanent", zap.Int("id", id))

	if id == 0 {
		w.logger.Error("DeleteWithdrawPermanent", zap.Int("id", id))
		return nil, withdraw_errors.ErrGrpcWithdrawInvalidID
	}

	_, err := w.withdrawCommand.DeleteWithdrawPermanent(ctx, id)

	if err != nil {
		w.logger.Error("DeleteWithdrawPermanent", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := w.mapper.ToProtoResponseWithdrawDelete("success", "Successfully deleted withdraw permanently")

	return so, nil
}

// RestoreAllWithdraw is a gRPC handler that restores all trashed withdrawal records.
// It logs the operation's success or failure and returns a gRPC response containing
// the result of the restoration or an error if the operation fails.
//
// Parameters:
//   - ctx: the context.Context object for tracing and cancellation.
//   - _: an empty protobuf message, as no request data is needed.
//
// Returns:
//   - A pointer to an ApiResponseWithdrawAll message, which indicates the success of the restoration operation.
//   - An error, which is non-nil if the operation fails.
func (s *withdrawCommandHandleGrpc) RestoreAllWithdraw(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseWithdrawAll, error) {
	s.logger.Info("RestoreAllWithdraw")

	_, err := s.withdrawCommand.RestoreAllWithdraw(ctx)

	if err != nil {
		s.logger.Error("RestoreAllWithdraw", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseWithdrawAll("success", "Successfully restore all withdraw")

	return so, nil
}

// DeleteAllWithdrawPermanent permanently deletes all trashed withdrawal records.
// It logs the operation's success or failure and returns a gRPC response containing
// the result of the deletion or an error if the operation fails.
//
// Parameters:
//   - ctx: the context.Context object for tracing and cancellation.
//   - _: an empty protobuf message, as no request data is needed.
//
// Returns:
//   - A pointer to an ApiResponseWithdrawAll message, which indicates the success of the deletion operation.
//   - An error, which is non-nil if the operation fails.
func (s *withdrawCommandHandleGrpc) DeleteAllWithdrawPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseWithdrawAll, error) {
	s.logger.Info("DeleteAllWithdrawPermanent")

	_, err := s.withdrawCommand.DeleteAllWithdrawPermanent(ctx)

	if err != nil {
		s.logger.Error("DeleteAllWithdrawPermanent", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseWithdrawAll("success", "Successfully delete withdraw permanent")

	return so, nil
}
