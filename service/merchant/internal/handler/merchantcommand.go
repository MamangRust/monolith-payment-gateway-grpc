package handler

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/service"
	pbmerchant "github.com/MamangRust/monolith-payment-gateway-pb/merchant"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/grpc"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/merchant"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantCommandHandleGrpc struct {
	pbmerchant.UnimplementedMerchantCommandServiceServer

	merchantCommand service.MerchantCommandService
	logger          logger.LoggerInterface
	mapper          protomapper.MerchantCommandProtoMapper
}

func NewMerchantCommandHandleGrpc(merchantCommand service.MerchantCommandService, logger logger.LoggerInterface, mapper protomapper.MerchantCommandProtoMapper) MerchantCommandHandleGrpc {
	return &merchantCommandHandleGrpc{
		merchantCommand: merchantCommand,
		logger:          logger,
		mapper:          mapper,
	}
}

// CreateMerchant creates a new merchant record with the provided name and user ID and returns the newly
// created record. The status of the merchant is set to "inactive" by default.
//
// Parameters:
//   - req: A CreateMerchantRequest containing the name, user ID and status of the merchant.
//
// Returns:
//   - A pointer to an ApiResponseMerchant containing the newly created record.
//   - An error if the record could not be created.
func (s *merchantCommandHandleGrpc) CreateMerchant(ctx context.Context, req *pbmerchant.CreateMerchantRequest) (*pbmerchant.ApiResponseMerchant, error) {
	request := requests.CreateMerchantRequest{
		Name:   req.GetName(),
		UserID: int(req.GetUserId()),
	}

	if err := request.Validate(); err != nil {
		s.logger.Error("CreateMerchant failed", zap.Any("error", err))
		return nil, merchant_errors.ErrGrpcValidateCreateMerchant
	}

	merchant, err := s.merchantCommand.CreateMerchant(ctx, &request)

	if err != nil {
		s.logger.Error("CreateMerchant failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMerchant("success", "Successfully created merchant", merchant)

	s.logger.Info("Successfully created merchant", zap.Bool("success", true))

	return so, nil

}

// UpdateMerchant updates an existing merchant record with the provided name and user ID and returns
// the newly updated record.
//
// Parameters:
//   - req: A UpdateMerchantRequest containing the name, user ID and status of the merchant.
//
// Returns:
//   - A pointer to an ApiResponseMerchant containing the newly updated record.
//   - An error if the record could not be updated.
func (s *merchantCommandHandleGrpc) UpdateMerchant(ctx context.Context, req *pbmerchant.UpdateMerchantRequest) (*pbmerchant.ApiResponseMerchant, error) {
	id := int(req.GetMerchantId())

	if id == 0 {
		s.logger.Error("UpdateMerchant failed", zap.Any("error", merchant_errors.ErrGrpcMerchantInvalidID))
		return nil, merchant_errors.ErrGrpcMerchantInvalidID
	}

	request := requests.UpdateMerchantRequest{
		MerchantID: &id,
		Name:       req.GetName(),
		UserID:     int(req.GetUserId()),
		Status:     req.GetStatus(),
	}

	if err := request.Validate(); err != nil {
		s.logger.Error("UpdateMerchant failed", zap.Any("error", err))
		return nil, merchant_errors.ErrGrpcValidateUpdateMerchant
	}

	merchant, err := s.merchantCommand.UpdateMerchant(ctx, &request)

	if err != nil {
		s.logger.Error("UpdateMerchant failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMerchant("success", "Successfully updated merchant", merchant)

	s.logger.Info("Successfully updated merchant", zap.Bool("success", true))

	return so, nil
}

// UpdateMerchantStatus updates the status of a merchant with the given ID.
// It validates the input request, ensuring the merchant ID and status are
// properly provided and updates the merchant's status accordingly.
//
// Parameters:
//   - ctx: Context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to UpdateMerchantStatusRequest containing the merchant ID and status.
//
// Returns:
//   - A pointer to ApiResponseMerchant containing the updated merchant information on success.
//   - An error if the update operation fails or if the input request is invalid.
func (s *merchantCommandHandleGrpc) UpdateMerchantStatus(ctx context.Context, req *pbmerchant.UpdateMerchantStatusRequest) (*pbmerchant.ApiResponseMerchant, error) {
	id := int(req.GetMerchantId())

	if id == 0 {
		s.logger.Error("UpdateMerchantStatus failed", zap.Any("error", merchant_errors.ErrGrpcMerchantInvalidID))
		return nil, merchant_errors.ErrGrpcMerchantInvalidID
	}

	request := requests.UpdateMerchantStatusRequest{
		MerchantID: &id,
		Status:     req.GetStatus(),
	}

	if err := request.Validate(); err != nil {
		s.logger.Error("UpdateMerchantStatus failed", zap.Any("error", err))
		return nil, merchant_errors.ErrGrpcValidateUpdateMerchantStatus
	}

	merchant, err := s.merchantCommand.UpdateMerchantStatus(ctx, &request)

	if err != nil {
		s.logger.Error("UpdateMerchantStatus failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMerchant("success", "Successfully updated merchant status", merchant)

	s.logger.Info("Successfully updated merchant status", zap.Bool("success", true))

	return so, nil
}

// TrashedMerchant soft-deletes a merchant by its ID.
//
// Parameters:
//   - ctx: Context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to FindByIdMerchantRequest containing the merchant ID.
//
// Returns:
//   - A pointer to ApiResponseMerchant containing the updated merchant information on success.
//   - An error if the delete operation fails or if the input request is invalid.
func (s *merchantCommandHandleGrpc) TrashedMerchant(ctx context.Context, req *pbmerchant.FindByIdMerchantRequest) (*pbmerchant.ApiResponseMerchantDeleteAt, error) {
	id := int(req.GetMerchantId())

	if id == 0 {
		s.logger.Error("TrashedMerchant failed", zap.Any("error", merchant_errors.ErrGrpcMerchantInvalidID))
		return nil, merchant_errors.ErrGrpcMerchantInvalidID
	}

	merchant, err := s.merchantCommand.TrashedMerchant(ctx, id)

	if err != nil {
		s.logger.Error("TrashedMerchant failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMerchantDeleteAt("success", "Successfully trashed merchant", merchant)

	s.logger.Info("Successfully trashed merchant", zap.Bool("success", true))

	return so, nil
}

// RestoreMerchant restores a merchant by its ID.
// It validates the input request to ensure the merchant ID is provided,
// then restores the merchant using the merchant command service.
// If successful, it returns a gRPC response containing the restored
// merchant information. If the operation fails or the input is invalid,
// an appropriate error is returned.
//
// Parameters:
//   - ctx: Context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to FindByIdMerchantRequest containing the merchant ID.
//
// Returns:
//   - A pointer to ApiResponseMerchant containing the restored merchant information on success.
//   - An error if the restore operation fails or if the input request is invalid.
func (s *merchantCommandHandleGrpc) RestoreMerchant(ctx context.Context, req *pbmerchant.FindByIdMerchantRequest) (*pbmerchant.ApiResponseMerchant, error) {
	id := int(req.GetMerchantId())

	if id == 0 {
		s.logger.Error("RestoreMerchant failed", zap.Any("error", merchant_errors.ErrGrpcMerchantInvalidID))
		return nil, merchant_errors.ErrGrpcMerchantInvalidID
	}

	merchant, err := s.merchantCommand.RestoreMerchant(ctx, id)

	if err != nil {
		s.logger.Error("RestoreMerchant failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMerchant("success", "Successfully restored merchant", merchant)

	s.logger.Info("Successfully restored merchant", zap.Bool("success", true))

	return so, nil
}

// DeleteMerchant permanently deletes a merchant by its ID.
// It validates the input request to ensure the merchant ID is provided,
// then deletes the merchant using the merchant command service.
// If successful, it returns a gRPC response containing the deleted
// merchant information. If the operation fails or the input is invalid,
// an appropriate error is returned.
//
// Parameters:
//   - ctx: Context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to FindByIdMerchantRequest containing the merchant ID.
//
// Returns:
//   - A pointer to ApiResponseMerchantDelete containing the deleted merchant information on success.
//   - An error if the delete operation fails or if the input request is invalid.
func (s *merchantCommandHandleGrpc) DeleteMerchant(ctx context.Context, req *pbmerchant.FindByIdMerchantRequest) (*pbmerchant.ApiResponseMerchantDelete, error) {
	id := int(req.GetMerchantId())

	if id == 0 {
		s.logger.Error("DeleteMerchant failed", zap.Any("error", merchant_errors.ErrGrpcMerchantInvalidID))
		return nil, merchant_errors.ErrGrpcMerchantInvalidID
	}

	_, err := s.merchantCommand.DeleteMerchantPermanent(ctx, id)

	if err != nil {
		s.logger.Error("DeleteMerchant failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMerchantDelete("success", "Successfully deleted merchant")

	s.logger.Info("Successfully deleted merchant", zap.Bool("success", true))

	return so, nil
}

// RestoreAllMerchant restores all trashed merchants.
//
// Args:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - _ *emptypbmerchant.Empty: Unused parameter.
//
// Returns:
//   - An ApiResponseMerchantAll containing the restored merchant information on success.
//   - An error if the restore operation fails.
func (s *merchantCommandHandleGrpc) RestoreAllMerchant(ctx context.Context, _ *emptypb.Empty) (*pbmerchant.ApiResponseMerchantAll, error) {
	_, err := s.merchantCommand.RestoreAllMerchant(ctx)

	if err != nil {
		s.logger.Error("RestoreAllMerchant failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMerchantAll("success", "Successfully restore all merchant")

	s.logger.Info("Successfully restore all merchant", zap.Bool("success", true))

	return so, nil
}

// DeleteAllMerchantPermanent permanently deletes all merchants.
// It logs the operation's success or failure and returns a gRPC response containing
// the result of the deletion or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - _ *emptypbmerchant.Empty: Unused parameter.
//
// Returns:
//   - A pointer to ApiResponseMerchantAll containing the result of the deletion on success.
//   - An error if the deletion operation fails.
func (s *merchantCommandHandleGrpc) DeleteAllMerchantPermanent(ctx context.Context, _ *emptypb.Empty) (*pbmerchant.ApiResponseMerchantAll, error) {
	_, err := s.merchantCommand.DeleteAllMerchantPermanent(ctx)

	if err != nil {
		s.logger.Error("DeleteAllMerchantPermanent failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMerchantAll("success", "Successfully delete all merchant")

	s.logger.Info("Successfully delete all merchant", zap.Bool("success", true))

	return so, nil
}
