package handler

import (
	"context"

	pb "github.com/MamangRust/monolith-payment-gateway-pb/transfer"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/grpc"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/transfer"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/service"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

// transferCommandHandleGrpc represents the gRPC handler for transfer operations.
type transferCommandHandleGrpc struct {
	pb.UnimplementedTransferCommandServiceServer

	transferCommandService service.TransferCommandService
	logger                 logger.LoggerInterface
	mapper                 protomapper.TransferCommandProtoMapper
}

func NewTransferCommandHandler(service service.TransferCommandService, logger logger.LoggerInterface, mapper protomapper.TransferCommandProtoMapper) TransferCommandHandleGrpc {
	return &transferCommandHandleGrpc{
		transferCommandService: service,
		logger:                 logger,
		mapper:                 mapper,
	}
}

// CreateTransfer implements the gRPC service for creating a transfer record.
//
// This function handles the "CreateTransfer" gRPC request by creating a new transfer record
// with the provided parameters. It validates the input request, logs the process of creating
// a transfer, and maps the result to a protocol buffer response object.
//
// Parameters:
//   - ctx: the context object for the gRPC request.
//   - request: a protobuf request containing the transfer parameters.
//
// Returns:
//   - A protobuf response containing the created transfer data, or an error object if the
//     creation operation fails.
func (s *transferCommandHandleGrpc) CreateTransfer(ctx context.Context, request *pb.CreateTransferRequest) (*pb.ApiResponseTransfer, error) {
	req := &requests.CreateTransferRequest{
		TransferFrom:   request.GetTransferFrom(),
		TransferTo:     request.GetTransferTo(),
		TransferAmount: int(request.GetTransferAmount()),
	}

	s.logger.Info("Starting create transfer process",
		zap.Any("request", req),
	)

	if err := req.Validate(); err != nil {
		s.logger.Error("Failed to create transfer", zap.Any("error", err))
		return nil, transfer_errors.ErrGrpcValidateCreateTransferRequest
	}

	res, err := s.transferCommandService.CreateTransaction(ctx, req)

	if err != nil {
		s.logger.Error("Failed to create transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTransfer("success", "Successfully created transfer", res)

	s.logger.Info("Successfully created transfer", zap.Bool("success", true))

	return so, nil
}

// UpdateTransfer implements the gRPC service for updating a transfer record.
//
// This function handles the "UpdateTransfer" gRPC request by updating the transfer record
// with the provided parameters. It validates the input request, logs the process of updating
// a transfer, and maps the result to a protocol buffer response object.
//
// Parameters:
//   - ctx: the context object for the gRPC request.
//   - request: a protobuf request containing the transfer ID and updated transfer data.
//
// Returns:
//   - A protobuf response containing the updated transfer data, or an error object if the
//     update operation fails.
func (s *transferCommandHandleGrpc) UpdateTransfer(ctx context.Context, request *pb.UpdateTransferRequest) (*pb.ApiResponseTransfer, error) {
	id := int(request.GetTransferId())

	s.logger.Info("Starting update transfer process",
		zap.Any("request", id),
	)

	if id == 0 {
		s.logger.Error("Failed to update transfer", zap.Any("error", transfer_errors.ErrGrpcTransferInvalidID))
		return nil, transfer_errors.ErrGrpcTransferInvalidID
	}

	req := &requests.UpdateTransferRequest{
		TransferID:     &id,
		TransferFrom:   request.GetTransferFrom(),
		TransferTo:     request.GetTransferTo(),
		TransferAmount: int(request.GetTransferAmount()),
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Failed to update transfer", zap.Any("error", err))
		return nil, transfer_errors.ErrGrpcValidateUpdateTransferRequest
	}

	res, err := s.transferCommandService.UpdateTransaction(ctx, req)

	if err != nil {
		s.logger.Error("Failed to update transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTransfer("success", "Successfully updated transfer", res)

	s.logger.Info("Successfully updated transfer", zap.Bool("success", true))

	return so, nil
}

// TrashedTransfer implements the gRPC service for trashing a transfer record by its ID.
//
// It logs the process of trashing the transfer and validates the input ID,
// ensuring that it is a positive integer. If the ID is invalid, it logs the error
// and returns an appropriate error response. The function calls the transfer command
// service to trash the transfer and maps the result to a protocol buffer response object.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - request: A pointer to a FindByIdTransferRequest containing the transfer ID to trash.
//
// Returns:
//   - A pointer to ApiResponseTransfer containing the trashed transfer data on success.
//   - An error if the trashing operation fails, or if the provided ID is invalid.
func (s *transferCommandHandleGrpc) TrashedTransfer(ctx context.Context, request *pb.FindByIdTransferRequest) (*pb.ApiResponseTransferDeleteAt, error) {
	id := int(request.GetTransferId())

	s.logger.Info("Starting trashed transfer process",
		zap.Any("request", id),
	)

	if id == 0 {
		s.logger.Error("Failed to trashed transfer", zap.Any("error", transfer_errors.ErrGrpcTransferInvalidID))
		return nil, transfer_errors.ErrGrpcTransferInvalidID
	}

	res, err := s.transferCommandService.TrashedTransfer(ctx, id)

	if err != nil {
		s.logger.Error("Failed to trashed transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTransferDeleteAt("success", "Successfully trashed transfer", res)

	s.logger.Info("Successfully trashed transfer", zap.Bool("success", true))

	return so, nil
}

// RestoreTransfer implements the gRPC service for restoring a transfer record from trashed.
//
// It logs the process of restoring the transfer and validates the input ID,
// ensuring that it is a positive integer. If the ID is invalid, it logs the error
// and returns an appropriate error response. The function calls the transfer command
// service to restore the transfer and maps the result to a protocol buffer response object.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - request: A pointer to a FindByIdTransferRequest containing the transfer ID to restore.
//
// Returns:
//   - A pointer to ApiResponseTransfer containing the restored transfer data on success.
//   - An error if the restoration operation fails, or if the provided ID is invalid.
func (s *transferCommandHandleGrpc) RestoreTransfer(ctx context.Context, request *pb.FindByIdTransferRequest) (*pb.ApiResponseTransfer, error) {
	id := int(request.GetTransferId())

	s.logger.Info("Starting restore transfer process",
		zap.Any("request", id),
	)

	if id == 0 {
		s.logger.Error("Failed to restore transfer", zap.Any("error", transfer_errors.ErrGrpcTransferInvalidID))
		return nil, transfer_errors.ErrGrpcTransferInvalidID
	}

	res, err := s.transferCommandService.RestoreTransfer(ctx, id)

	if err != nil {
		s.logger.Error("Failed to restore transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTransfer("success", "Successfully restored transfer", res)

	s.logger.Info("Successfully restored transfer", zap.Bool("success", true))

	return so, nil
}

// DeleteTransferPermanent implements the gRPC service for deleting a transfer record permanently.
//
// It logs the process of deleting the transfer and validates the input ID,
// ensuring that it is a positive integer. If the ID is invalid, it logs the error
// and returns an appropriate error response. The function calls the transfer command
// service to delete the transfer permanently and maps the result to a protocol buffer response object.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - request: A pointer to a FindByIdTransferRequest containing the transfer ID to delete.
//
// Returns:
//   - A pointer to ApiResponseTransferDelete containing the deleted transfer data on success.
//   - An error if the deletion operation fails, or if the provided ID is invalid.
func (s *transferCommandHandleGrpc) DeleteTransferPermanent(ctx context.Context, request *pb.FindByIdTransferRequest) (*pb.ApiResponseTransferDelete, error) {
	id := int(request.GetTransferId())

	s.logger.Info("Starting delete transfer process",
		zap.Any("request", id),
	)

	if id == 0 {
		s.logger.Error("Failed to delete transfer", zap.Any("error", transfer_errors.ErrGrpcTransferInvalidID))
		return nil, transfer_errors.ErrGrpcTransferInvalidID
	}

	_, err := s.transferCommandService.DeleteTransferPermanent(ctx, id)

	if err != nil {
		s.logger.Error("Failed to delete transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTransferDelete("success", "Successfully restored transfer")

	s.logger.Info("Successfully deleted transfer", zap.Bool("success", true))

	return so, nil
}

// RestoreAllTransfer implements the gRPC service for restoring all trashed transfer records.
//
// It logs the process of restoring all transfer records and calls the transfer command service to restore all
// transfer records. The function maps the result to a protocol buffer response object.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - _ *emptypb.Empty: Unused parameter.
//
// Returns:
//   - A pointer to ApiResponseTransferAll containing the restored transfer records on success.
//   - An error if the restoration operation fails.
func (s *transferCommandHandleGrpc) RestoreAllTransfer(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTransferAll, error) {
	s.logger.Info("Starting restore all transfer process")

	_, err := s.transferCommandService.RestoreAllTransfer(ctx)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTransferAll("success", "Successfully restored transfer")

	s.logger.Info("Successfully restored all transfer", zap.Bool("success", true))

	return so, nil
}

// DeleteAllTransferPermanent implements the gRPC service for permanently deleting all trashed transfer records.
//
// It logs the process of deleting all transfer records and calls the transfer command service to delete all
// transfer records. The function maps the result to a protocol buffer response object.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - _ *emptypb.Empty: Unused parameter.
//
// Returns:
//   - A pointer to ApiResponseTransferAll containing the deleted transfer records on success.
//   - An error if the deletion operation fails.
func (s *transferCommandHandleGrpc) DeleteAllTransferPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTransferAll, error) {
	s.logger.Info("Starting delete all transfer process")

	_, err := s.transferCommandService.DeleteAllTransferPermanent(ctx)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTransferAll("success", "delete transfer permanent")

	s.logger.Info("Successfully deleted all transfer", zap.Bool("success", true))

	return so, nil
}
