package handler

import (
	"context"

	pb "github.com/MamangRust/monolith-payment-gateway-pb/saldo"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/service"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/grpc"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/saldo"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type saldoCommandHandleGrpc struct {
	pb.UnimplementedSaldoCommandServiceServer

	service service.SaldoCommandService
	logger  logger.LoggerInterface
	mapper  protomapper.SaldoCommandProtoMapper
}

func NewSaldoCommandHandleGrpc(query service.SaldoCommandService, logger logger.LoggerInterface, mapper protomapper.SaldoCommandProtoMapper) SaldoCommandHandleGrpc {
	return &saldoCommandHandleGrpc{
		service: query,
		logger:  logger,
		mapper:  mapper,
	}
}

// CreateSaldo is a gRPC handler that creates a new saldo record.
//
// Parameters:
//   - ctx: the context.Context object for tracing and cancellation.
//   - req: a pointer to a CreateSaldoRequest message, which contains the card number and total balance for the new saldo.
//
// Returns:
//   - A pointer to a ApiResponseSaldo message, which contains the created saldo record.
//   - An error, which is non-nil if the operation fails.
func (s *saldoCommandHandleGrpc) CreateSaldo(ctx context.Context, req *pb.CreateSaldoRequest) (*pb.ApiResponseSaldo, error) {
	request := requests.CreateSaldoRequest{
		CardNumber:   req.GetCardNumber(),
		TotalBalance: int(req.GetTotalBalance()),
	}

	s.logger.Info("Creating saldo record", zap.String("card_number", request.CardNumber), zap.Int("total_balance", request.TotalBalance))

	if err := request.Validate(); err != nil {
		s.logger.Error("Create failed", zap.Any("error", err))
		return nil, saldo_errors.ErrGrpcValidateCreateSaldo
	}

	saldo, err := s.service.CreateSaldo(ctx, &request)

	if err != nil {
		s.logger.Error("Create failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseSaldo("success", "Successfully created saldo record", saldo)

	s.logger.Info("Successfully created saldo record", zap.Bool("success", true))

	return so, nil

}

// UpdateSaldo is a gRPC handler that updates a saldo record.
//
// Parameters:
//   - ctx: the context.Context object for tracing and cancellation.
//   - req: a pointer to a UpdateSaldoRequest message, which contains the ID of the saldo record to be updated,
//     as well as the new card number and total balance.
//
// Returns:
//   - A pointer to a ApiResponseSaldo message, which contains the updated saldo record.
//   - An error, which is non-nil if the operation fails.
func (s *saldoCommandHandleGrpc) UpdateSaldo(ctx context.Context, req *pb.UpdateSaldoRequest) (*pb.ApiResponseSaldo, error) {
	id := int(req.GetSaldoId())

	s.logger.Info("Updating saldo record", zap.Int("id", id))

	if id == 0 {
		s.logger.Error("Update failed", zap.Any("error", saldo_errors.ErrGrpcSaldoInvalidID))
		return nil, saldo_errors.ErrGrpcSaldoInvalidID
	}

	request := requests.UpdateSaldoRequest{
		SaldoID:      &id,
		CardNumber:   req.GetCardNumber(),
		TotalBalance: int(req.GetTotalBalance()),
	}

	if err := request.Validate(); err != nil {
		s.logger.Error("Update failed", zap.Any("error", err))
		return nil, saldo_errors.ErrGrpcValidateUpdateSaldo
	}

	saldo, err := s.service.UpdateSaldo(ctx, &request)

	if err != nil {
		s.logger.Error("Update failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseSaldo("success", "Successfully updated saldo record", saldo)

	s.logger.Info("Successfully updated saldo record", zap.Bool("success", true))

	return so, nil
}

// TrashedSaldo is a gRPC handler that trashes a saldo record.
//
// Parameters:
//   - ctx: the context.Context object for tracing and cancellation.
//   - req: a pointer to a FindByIdSaldoRequest message, which contains the ID of the saldo record to be trashed.
//
// Returns:
//   - A pointer to a ApiResponseSaldo message, which contains the trashed saldo record.
//   - An error, which is non-nil if the operation fails.
func (s *saldoCommandHandleGrpc) TrashedSaldo(ctx context.Context, req *pb.FindByIdSaldoRequest) (*pb.ApiResponseSaldoDeleteAt, error) {
	id := int(req.GetSaldoId())

	s.logger.Info("Trashing saldo record", zap.Int("id", id))

	if id == 0 {
		s.logger.Error("TrashedSaldo failed", zap.Any("error", saldo_errors.ErrGrpcSaldoInvalidID))
		return nil, saldo_errors.ErrGrpcSaldoInvalidID
	}

	saldo, err := s.service.TrashSaldo(ctx, id)

	if err != nil {
		s.logger.Error("TrashedSaldo failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseSaldoDeleteAt("success", "Successfully trashed saldo record", saldo)

	s.logger.Info("Successfully trashed saldo record", zap.Bool("success", true))

	return so, nil
}

// RestoreSaldo is a gRPC handler that restores a trashed saldo record.
//
// Parameters:
//   - ctx: the context.Context object for tracing and cancellation.
//   - req: a pointer to a FindByIdSaldoRequest message, which contains the ID of the saldo record to be restored.
//
// Returns:
//   - A pointer to a ApiResponseSaldo message, which contains the restored saldo record.
//   - An error, which is non-nil if the operation fails.
func (s *saldoCommandHandleGrpc) RestoreSaldo(ctx context.Context, req *pb.FindByIdSaldoRequest) (*pb.ApiResponseSaldo, error) {
	id := int(req.GetSaldoId())

	s.logger.Info("Restoring saldo record", zap.Int("id", id))

	if id == 0 {
		s.logger.Error("RestoreSaldo failed", zap.Any("error", saldo_errors.ErrGrpcSaldoInvalidID))
		return nil, saldo_errors.ErrGrpcSaldoInvalidID
	}

	saldo, err := s.service.RestoreSaldo(ctx, id)

	if err != nil {
		s.logger.Error("RestoreSaldo failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseSaldo("success", "Successfully restored saldo record", saldo)

	s.logger.Info("Successfully restored saldo", zap.Bool("success", true))

	return so, nil
}

// DeleteSaldo is a gRPC handler that deletes a saldo record permanently.
//
// Parameters:
//   - ctx: the context.Context object for tracing and cancellation.
//   - req: a pointer to a FindByIdSaldoRequest message, which contains the ID of the saldo record to be deleted.
//
// Returns:
//   - A pointer to a ApiResponseSaldoDelete message, which contains the deleted saldo record.
//   - An error, which is non-nil if the operation fails.
func (s *saldoCommandHandleGrpc) DeleteSaldo(ctx context.Context, req *pb.FindByIdSaldoRequest) (*pb.ApiResponseSaldoDelete, error) {
	id := int(req.GetSaldoId())

	s.logger.Info("Deleting saldo record", zap.Int("id", id))

	if id == 0 {
		s.logger.Error("DeleteSaldo failed", zap.Any("error", saldo_errors.ErrGrpcSaldoInvalidID))
		return nil, saldo_errors.ErrGrpcSaldoInvalidID
	}

	_, err := s.service.DeleteSaldoPermanent(ctx, id)

	if err != nil {
		s.logger.Error("DeleteSaldo failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseSaldoDelete("success", "Successfully deleted saldo record")

	s.logger.Info("Successfully deleted saldo record", zap.Bool("success", true))

	return so, nil
}

// RestoreAllSaldo is a gRPC handler that restores all saldo records.
//
// Parameters:
//   - ctx: the context.Context object for tracing and cancellation.
//   - _: an empty protobuf message, as no request data is needed.
//
// Returns:
//   - A pointer to a ApiResponseSaldoAll message, which indicates the success of the restoration operation.
//   - An error, which is non-nil if the operation fails.
func (s *saldoCommandHandleGrpc) RestoreAllSaldo(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseSaldoAll, error) {
	s.logger.Info("Restoring all saldo record")

	_, err := s.service.RestoreAllSaldo(ctx)

	if err != nil {
		s.logger.Error("RestoreAllSaldo failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseSaldoAll("success", "Successfully restore all saldo")

	s.logger.Info("Successfully restored all saldo", zap.Bool("success", true))

	return so, nil
}

// DeleteAllSaldoPermanent is a gRPC handler that deletes all saldo records permanently.
//
// Parameters:
//   - ctx: the context.Context object for tracing and cancellation.
//   - _: an empty protobuf message, as no request data is needed.
//
// Returns:
//   - A pointer to a ApiResponseSaldoAll message, which indicates the success of the deletion operation.
//   - An error, which is non-nil if the operation fails.
func (s *saldoCommandHandleGrpc) DeleteAllSaldoPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseSaldoAll, error) {
	s.logger.Info("Deleting all saldo record")

	_, err := s.service.DeleteAllSaldoPermanent(ctx)

	if err != nil {
		s.logger.Error("DeleteAllSaldoPermanent failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseSaldoAll("success", "delete saldo permanent")

	s.logger.Info("Successfully deleted all permanent saldo", zap.Bool("success", true))

	return so, nil
}
