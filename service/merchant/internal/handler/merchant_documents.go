package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/service"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	merchantdocument_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_document_errors"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto"
	"github.com/MamangRust/monolith-payment-gateway-shared/pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantDocumentHandleGrpc struct {
	pb.UnimplementedMerchantDocumentServiceServer
	merchantDocumentQuery   service.MerchantDocumentQueryService
	merchantDocumentCommand service.MerchantDocumentCommandService
	logger                  logger.LoggerInterface
	mapping                 protomapper.MerchantDocumentProtoMapper
}

func NewMerchantDocumentHandleGrpc(
	service service.Service,
	mapping protomapper.MerchantDocumentProtoMapper,
	logger logger.LoggerInterface,
) pb.MerchantDocumentServiceServer {
	return &merchantDocumentHandleGrpc{
		merchantDocumentQuery:   service.MerchantDocumentQuery,
		merchantDocumentCommand: service.MerchantDocumentCommand,
		mapping:                 mapping,
		logger:                  logger,
	}
}

func (s *merchantDocumentHandleGrpc) FindAll(ctx context.Context, req *pb.FindAllMerchantDocumentsRequest) (*pb.ApiResponsePaginationMerchantDocument, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	s.logger.Debug("Fetching merchant document records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchantDocuments{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	documents, totalRecords, err := s.merchantDocumentQuery.FindAll(&reqService)
	if err != nil {
		s.logger.Debug("FindAll failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return s.mapping.ToProtoResponsePaginationMerchantDocument(paginationMeta, "success", "Successfully fetched merchant documents", documents), nil
}

func (s *merchantDocumentHandleGrpc) FindById(ctx context.Context, req *pb.FindMerchantDocumentByIdRequest) (*pb.ApiResponseMerchantDocument, error) {
	id := int(req.GetDocumentId())

	s.logger.Debug("Fetching merchant document by id", zap.Int("id", id))

	if id == 0 {
		s.logger.Debug("FindById failed", zap.Any("error", merchantdocument_errors.ErrGrpcMerchantInvalidID))
		return nil, merchantdocument_errors.ErrGrpcMerchantInvalidID
	}

	document, err := s.merchantDocumentQuery.FindById(id)
	if err != nil {
		s.logger.Debug("FindById failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseMerchantDocument("success", "Successfully fetched merchant document", document), nil

}

func (s *merchantDocumentHandleGrpc) FindAllActive(ctx context.Context, req *pb.FindAllMerchantDocumentsRequest) (*pb.ApiResponsePaginationMerchantDocument, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	s.logger.Debug("Fetching active merchant document records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchantDocuments{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	documents, totalRecords, err := s.merchantDocumentQuery.FindByActive(&reqService)
	if err != nil {
		s.logger.Debug("FindAllActive failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return s.mapping.ToProtoResponsePaginationMerchantDocument(paginationMeta, "success", "Successfully fetched active merchant documents", documents), nil
}

func (s *merchantDocumentHandleGrpc) FindAllTrashed(ctx context.Context, req *pb.FindAllMerchantDocumentsRequest) (*pb.ApiResponsePaginationMerchantDocumentAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	s.logger.Debug("Fetching trashed merchant document records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchantDocuments{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	documents, totalRecords, err := s.merchantDocumentQuery.FindByTrashed(&reqService)
	if err != nil {
		s.logger.Debug("FindAllTrashed failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return s.mapping.ToProtoResponsePaginationMerchantDocumentDeleteAt(paginationMeta, "success", "Successfully fetched trashed merchant documents", documents), nil
}

func (s *merchantDocumentHandleGrpc) Create(ctx context.Context, req *pb.CreateMerchantDocumentRequest) (*pb.ApiResponseMerchantDocument, error) {
	request := requests.CreateMerchantDocumentRequest{
		MerchantID:   int(req.GetMerchantId()),
		DocumentType: req.GetDocumentType(),
		DocumentUrl:  req.GetDocumentUrl(),
	}

	s.logger.Debug("Creating merchant document", zap.Any("request", request))

	if err := request.Validate(); err != nil {
		s.logger.Debug("Create failed", zap.Any("error", err))
		return nil, merchantdocument_errors.ErrGrpcValidateCreateMerchantDocument
	}

	document, err := s.merchantDocumentCommand.CreateMerchantDocument(&request)
	if err != nil {
		s.logger.Debug("Create failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseMerchantDocument("success", "Successfully created merchant document", document), nil
}

func (s *merchantDocumentHandleGrpc) Update(ctx context.Context, req *pb.UpdateMerchantDocumentRequest) (*pb.ApiResponseMerchantDocument, error) {
	id := int(req.GetDocumentId())

	s.logger.Debug("Updating merchant document", zap.Int("id", id))

	if id == 0 {
		s.logger.Debug("Update failed", zap.Any("error", merchantdocument_errors.ErrGrpcMerchantInvalidID))
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
		s.logger.Debug("Update failed", zap.Any("error", err))
		return nil, merchantdocument_errors.ErrGrpcFailedUpdateMerchantDocument
	}

	document, err := s.merchantDocumentCommand.UpdateMerchantDocument(&request)
	if err != nil {
		s.logger.Debug("Update failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseMerchantDocument("success", "Successfully updated merchant document", document), nil
}

func (s *merchantDocumentHandleGrpc) UpdateStatus(ctx context.Context, req *pb.UpdateMerchantDocumentStatusRequest) (*pb.ApiResponseMerchantDocument, error) {
	id := int(req.GetDocumentId())

	s.logger.Debug("Updating merchant document status", zap.Int("id", id))

	if id == 0 {
		s.logger.Debug("UpdateStatus failed", zap.Any("error", merchantdocument_errors.ErrGrpcMerchantInvalidID))
		return nil, merchantdocument_errors.ErrGrpcMerchantInvalidID
	}

	request := requests.UpdateMerchantDocumentStatusRequest{
		DocumentID: &id,
		MerchantID: int(req.GetMerchantId()),
		Status:     req.GetStatus(),
		Note:       req.GetNote(),
	}

	if err := request.Validate(); err != nil {
		s.logger.Debug("UpdateStatus failed", zap.Any("error", err))
		return nil, merchantdocument_errors.ErrGrpcFailedUpdateMerchantDocument
	}

	document, err := s.merchantDocumentCommand.UpdateMerchantDocumentStatus(&request)
	if err != nil {
		s.logger.Debug("UpdateStatus failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseMerchantDocument("success", "Successfully updated merchant document status", document), nil
}

func (s *merchantDocumentHandleGrpc) Trashed(ctx context.Context, req *pb.TrashedMerchantDocumentRequest) (*pb.ApiResponseMerchantDocument, error) {
	id := int(req.GetDocumentId())

	s.logger.Debug("Trashing merchant document", zap.Int("id", id))

	if id == 0 {
		s.logger.Debug("Trashed failed", zap.Any("error", merchantdocument_errors.ErrGrpcMerchantInvalidID))
		return nil, merchantdocument_errors.ErrGrpcMerchantInvalidID
	}

	document, err := s.merchantDocumentCommand.TrashedMerchantDocument(id)
	if err != nil {
		s.logger.Debug("Trashed failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseMerchantDocument("success", "Successfully trashed merchant document", document), nil
}

func (s *merchantDocumentHandleGrpc) Restore(ctx context.Context, req *pb.RestoreMerchantDocumentRequest) (*pb.ApiResponseMerchantDocument, error) {
	id := int(req.GetDocumentId())

	s.logger.Debug("Restoring merchant document", zap.Int("id", id))

	if id == 0 {
		s.logger.Debug("Restore failed", zap.Any("error", merchantdocument_errors.ErrGrpcMerchantInvalidID))
		return nil, merchantdocument_errors.ErrGrpcMerchantInvalidID
	}

	document, err := s.merchantDocumentCommand.RestoreMerchantDocument(id)
	if err != nil {
		s.logger.Debug("Restore failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseMerchantDocument("success", "Successfully restored merchant document", document), nil
}

func (s *merchantDocumentHandleGrpc) DeletePermanent(ctx context.Context, req *pb.DeleteMerchantDocumentPermanentRequest) (*pb.ApiResponseMerchantDocumentDelete, error) {
	id := int(req.GetDocumentId())

	s.logger.Debug("Permanently deleting merchant document", zap.Int("id", id))

	if id == 0 {
		s.logger.Debug("DeletePermanent failed", zap.Any("error", merchantdocument_errors.ErrGrpcMerchantInvalidID))
		return nil, merchantdocument_errors.ErrGrpcMerchantInvalidID
	}

	_, err := s.merchantDocumentCommand.DeleteMerchantDocumentPermanent(id)
	if err != nil {
		s.logger.Debug("DeletePermanent failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseMerchantDocumentDelete("success", "Successfully permanently deleted merchant document"), nil
}

func (s *merchantDocumentHandleGrpc) RestoreAll(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantDocumentAll, error) {
	s.logger.Debug("Restoring all merchant documents")

	_, err := s.merchantDocumentCommand.RestoreAllMerchantDocument()
	if err != nil {
		s.logger.Debug("RestoreAll failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseMerchantDocumentAll("success", "Successfully restored all merchant documents"), nil
}

func (s *merchantDocumentHandleGrpc) DeleteAllPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantDocumentAll, error) {
	s.logger.Debug("Permanently deleting all merchant documents")

	_, err := s.merchantDocumentCommand.DeleteAllMerchantDocumentPermanent()
	if err != nil {
		s.logger.Debug("DeleteAllPermanent failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseMerchantDocumentAll("success", "Successfully permanently deleted all merchant documents"), nil
}
