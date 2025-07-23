package handler

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/service"
	pbdocument "github.com/MamangRust/monolith-payment-gateway-pb/merchantdocument"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	merchantdocument_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_document_errors/grpc"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/merchantdocument"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantDocumentCommandHandleGrpc struct {
	pbdocument.UnimplementedMerchantDocumentCommandServiceServer

	merchantDocumentCommand service.MerchantDocumentCommandService
	logger                  logger.LoggerInterface
	mapper                  protomapper.MerchantDocumentCommandProtoMapper
}

func NewMerchantDocumentCommandHandleGrpc(merchantCommand service.MerchantDocumentCommandService, logger logger.LoggerInterface, mapper protomapper.MerchantDocumentCommandProtoMapper) MerchantDocumentCommandHandleGrpc {
	return &merchantDocumentCommandHandleGrpc{
		merchantDocumentCommand: merchantCommand,
		logger:                  logger,
		mapper:                  mapper,
	}
}

// Create creates a new merchant document with the given request parameters and returns a gRPC response containing the newly created document on success.
// It logs the operation's success or failure. It returns an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancelation signals, and deadlines.
//   - req: A pointer to a CreateMerchantDocumentRequest containing the merchant ID, document type, and document URL.
//
// Returns:
//   - A pointer to ApiResponseMerchantDocument containing the newly created merchant document on success.
//   - An error if the creation operation fails.
func (s *merchantDocumentCommandHandleGrpc) Create(ctx context.Context, req *pbdocument.CreateMerchantDocumentRequest) (*pbdocument.ApiResponseMerchantDocument, error) {
	request := requests.CreateMerchantDocumentRequest{
		MerchantID:   int(req.GetMerchantId()),
		DocumentType: req.GetDocumentType(),
		DocumentUrl:  req.GetDocumentUrl(),
	}

	s.logger.Info("Creating merchant document", zap.Any("request", request))

	if err := request.Validate(); err != nil {
		s.logger.Error("Create failed", zap.Any("error", err))
		return nil, merchantdocument_errors.ErrGrpcValidateCreateMerchantDocument
	}

	document, err := s.merchantDocumentCommand.CreateMerchantDocument(ctx, &request)
	if err != nil {
		s.logger.Error("Create failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMerchantDocument("success", "Successfully created merchant document", document)

	s.logger.Info("Create success", zap.Bool("success", true))

	return so, nil
}

// Update updates an existing merchant document with the given request parameters and returns a gRPC response containing the newly updated document on success.
// It logs the operation's success or failure. It returns an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancelation signals, and deadlines.
//   - req: A pointer to an UpdateMerchantDocumentRequest containing the document ID, merchant ID, document type, document URL, status, and note.
//
// Returns:
//   - A pointer to ApiResponseMerchantDocument containing the newly updated merchant document on success.
//   - An error if the update operation fails.
func (s *merchantDocumentCommandHandleGrpc) Update(ctx context.Context, req *pbdocument.UpdateMerchantDocumentRequest) (*pbdocument.ApiResponseMerchantDocument, error) {
	id := int(req.GetDocumentId())

	s.logger.Info("Updating merchant document", zap.Int("id", id))

	if id == 0 {
		s.logger.Error("Update failed", zap.Any("error", merchantdocument_errors.ErrGrpcMerchantInvalidID))
		return nil, merchantdocument_errors.ErrGrpcMerchantInvalidID
	}

	request := requests.UpdateMerchantDocumentRequest{
		DocumentID:   &id,
		MerchantID:   int(req.GetMerchantId()),
		DocumentType: req.GetDocumentType(),
		DocumentUrl:  req.GetDocumentUrl(),
		Status:       req.GetStatus(),
		Note:         req.GetNote(),
	}

	if err := request.Validate(); err != nil {
		s.logger.Error("Update failed", zap.Any("error", err))
		return nil, merchantdocument_errors.ErrGrpcFailedUpdateMerchantDocument
	}

	document, err := s.merchantDocumentCommand.UpdateMerchantDocument(ctx, &request)
	if err != nil {
		s.logger.Error("Update failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMerchantDocument("success", "Successfully updated merchant document", document)

	s.logger.Info("Update success", zap.Bool("success", true))

	return so, nil
}

// UpdateStatus updates the status and note of an existing merchant document with the given request parameters.
// It logs the operation's success or failure. It returns a gRPC response containing the updated document on success.
// It returns an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancelation signals, and deadlines.
//   - req: A pointer to an UpdateMerchantDocumentStatusRequest containing the document ID, merchant ID, status, and note.
//
// Returns:
//   - A pointer to ApiResponseMerchantDocument containing the updated merchant document on success.
//   - An error if the update operation fails.
func (s *merchantDocumentCommandHandleGrpc) UpdateStatus(ctx context.Context, req *pbdocument.UpdateMerchantDocumentStatusRequest) (*pbdocument.ApiResponseMerchantDocument, error) {
	id := int(req.GetDocumentId())

	s.logger.Info("Updating merchant document status", zap.Int("id", id))

	if id == 0 {
		s.logger.Error("UpdateStatus failed", zap.Any("error", merchantdocument_errors.ErrGrpcMerchantInvalidID))
		return nil, merchantdocument_errors.ErrGrpcMerchantInvalidID
	}

	request := requests.UpdateMerchantDocumentStatusRequest{
		DocumentID: &id,
		MerchantID: int(req.GetMerchantId()),
		Status:     req.GetStatus(),
		Note:       req.GetNote(),
	}

	if err := request.Validate(); err != nil {
		s.logger.Error("UpdateStatus failed", zap.Any("error", err))
		return nil, merchantdocument_errors.ErrGrpcFailedUpdateMerchantDocument
	}

	document, err := s.merchantDocumentCommand.UpdateMerchantDocumentStatus(ctx, &request)
	if err != nil {
		s.logger.Error("UpdateStatus failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMerchantDocument("success", "Successfully updated merchant document status", document)

	s.logger.Info("Successfully updated merchant document status", zap.Bool("success", true))

	return so, nil
}

// Trashed marks a merchant document as trashed by its ID.
// It logs the operation's success or failure and returns a gRPC response containing
// the trashed merchant document or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancelation signals, and deadlines.
//   - req: A pointer to a TrashedMerchantDocumentRequest containing the document ID.
//
// Returns:
//   - A pointer to ApiResponseMerchantDocument containing the trashed merchant document on success.
//   - An error if the trashing operation fails.

func (s *merchantDocumentCommandHandleGrpc) Trashed(ctx context.Context, req *pbdocument.FindMerchantDocumentByIdRequest) (*pbdocument.ApiResponseMerchantDocumentDeleteAt, error) {
	id := int(req.GetDocumentId())

	s.logger.Info("Trashing merchant document", zap.Int("id", id))

	if id == 0 {
		s.logger.Error("Trashed failed", zap.Any("error", merchantdocument_errors.ErrGrpcMerchantInvalidID))
		return nil, merchantdocument_errors.ErrGrpcMerchantInvalidID
	}

	document, err := s.merchantDocumentCommand.TrashedMerchantDocument(ctx, id)
	if err != nil {
		s.logger.Error("Trashed failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMerchantDocumentDeleteAt("success", "Successfully trashed merchant document", document)

	s.logger.Info("Successfully trashed merchant document", zap.Bool("success", true))

	return so, nil
}

// Restore restores a merchant document by its ID.
// It logs the operation's success or failure and returns a gRPC response containing
// the restored merchant document or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancelation signals, and deadlines.
//   - req: A pointer to a RestoreMerchantDocumentRequest containing the document ID.
//
// Returns:
//   - A pointer to ApiResponseMerchantDocument containing the restored merchant document on success.
//   - An error if the restoration operation fails.
func (s *merchantDocumentCommandHandleGrpc) Restore(ctx context.Context, req *pbdocument.FindMerchantDocumentByIdRequest) (*pbdocument.ApiResponseMerchantDocument, error) {
	id := int(req.GetDocumentId())

	s.logger.Info("Restoring merchant document", zap.Int("id", id))

	if id == 0 {
		s.logger.Error("Restore failed", zap.Any("error", merchantdocument_errors.ErrGrpcMerchantInvalidID))
		return nil, merchantdocument_errors.ErrGrpcMerchantInvalidID
	}

	document, err := s.merchantDocumentCommand.RestoreMerchantDocument(ctx, id)
	if err != nil {
		s.logger.Error("Restore failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMerchantDocument("success", "Successfully restored merchant document", document)

	s.logger.Info("Successfully restored merchant document", zap.Bool("success", true))

	return so, nil
}

// DeletePermanent permanently deletes a merchant document by its ID.
// It logs the operation's success or failure and returns a gRPC response containing
// the result of the deletion or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancelation signals, and deadlines.
//   - req: A pointer to a DeleteMerchantDocumentPermanentRequest containing the document ID.
//
// Returns:
//   - A pointer to ApiResponseMerchantDocumentDelete indicating the success of the deletion.
//   - An error if the deletion operation fails.
func (s *merchantDocumentCommandHandleGrpc) DeletePermanent(ctx context.Context, req *pbdocument.FindMerchantDocumentByIdRequest) (*pbdocument.ApiResponseMerchantDocumentDelete, error) {
	id := int(req.GetDocumentId())

	s.logger.Info("Permanently deleting merchant document", zap.Int("id", id))

	if id == 0 {
		s.logger.Error("DeletePermanent failed", zap.Any("error", merchantdocument_errors.ErrGrpcMerchantInvalidID))
		return nil, merchantdocument_errors.ErrGrpcMerchantInvalidID
	}

	_, err := s.merchantDocumentCommand.DeleteMerchantDocumentPermanent(ctx, id)
	if err != nil {
		s.logger.Error("DeletePermanent failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMerchantDocumentDelete("success", "Successfully permanently deleted merchant document")

	s.logger.Info("Successfully permanently deleted merchant document", zap.Bool("success", true))

	return so, nil
}

// RestoreAll restores all merchant documents that were previously deleted.
// It logs the operation's success or failure and returns a gRPC response containing
// the result of the restoration or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancelation signals, and deadlines.
//   - _ *emptypb.Empty: Unused parameter.
//
// Returns:
//   - A pointer to ApiResponseMerchantDocumentAll containing the restored merchant documents on success.
//   - An error if the restoration operation fails.
func (s *merchantDocumentCommandHandleGrpc) RestoreAll(ctx context.Context, _ *emptypb.Empty) (*pbdocument.ApiResponseMerchantDocumentAll, error) {
	s.logger.Info("Restoring all merchant documents")

	_, err := s.merchantDocumentCommand.RestoreAllMerchantDocument(ctx)
	if err != nil {
		s.logger.Error("RestoreAll failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMerchantDocumentAll("success", "Successfully restored all merchant documents")

	s.logger.Info("Successfully restored all merchant documents", zap.Bool("success", true))

	return so, nil
}

// DeleteAllPermanent permanently deletes all merchant documents that were previously deleted.
// It logs the operation's success or failure and returns a gRPC response containing
// the result of the deletion or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancelation signals, and deadlines.
//   - _ *emptypb.Empty: Unused parameter.
//
// Returns:
//   - A pointer to ApiResponseMerchantDocumentAll containing the deleted merchant documents on success.
//   - An error if the deletion operation fails.
func (s *merchantDocumentCommandHandleGrpc) DeleteAllPermanent(ctx context.Context, _ *emptypb.Empty) (*pbdocument.ApiResponseMerchantDocumentAll, error) {
	s.logger.Info("Permanently deleting all merchant documents")

	_, err := s.merchantDocumentCommand.DeleteAllMerchantDocumentPermanent(ctx)
	if err != nil {
		s.logger.Error("DeleteAllPermanent failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMerchantDocumentAll("success", "Successfully permanently deleted all merchant documents")

	s.logger.Info("Successfully permanently deleted all merchant documents", zap.Bool("success", true))

	return so, nil
}
