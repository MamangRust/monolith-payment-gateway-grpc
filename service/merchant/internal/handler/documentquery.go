package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/service"
	pb "github.com/MamangRust/monolith-payment-gateway-pb"
	pbdocument "github.com/MamangRust/monolith-payment-gateway-pb/merchantdocument"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	merchantdocument_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_document_errors/grpc"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/merchantdocument"
	"go.uber.org/zap"
)

type merchantDocumentQueryHandleGrpc struct {
	pbdocument.UnsafeMerchantDocumentServiceServer

	merchantDocumentQuery service.MerchantDocumentQueryService
	logger                logger.LoggerInterface
	mapper                protomapper.MerchantDocumentQueryProtoMapper
}

func NewMerchantDocumentQueryHandleGrpc(merchantQuery service.MerchantDocumentQueryService, logger logger.LoggerInterface, mapper protomapper.MerchantDocumentQueryProtoMapper) MerchantDocumentQueryHandleGrpc {
	return &merchantDocumentQueryHandleGrpc{
		merchantDocumentQuery: merchantQuery,
		logger:                logger,
		mapper:                mapper,
	}
}

// FindAll retrieves a paginated list of merchant documents based on the provided request parameters.
// It supports pagination and search functionalities.
//
// The function uses the merchantDocumentQuery service to fetch the records and logs the operation's
// success or failure. It returns a gRPC response containing the paginated list of merchant documents
// or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancelation signals, and deadlines.
//   - req: A pointer to a FindAllMerchantDocumentsRequest containing pagination and search details.
//
// Returns:
//   - A pointer to ApiResponsePaginationMerchantDocument containing the list of merchant documents
//     and pagination metadata on success.
//   - An error if the retrieval operation fails.
func (s *merchantDocumentQueryHandleGrpc) FindAll(ctx context.Context, req *pbdocument.FindAllMerchantDocumentsRequest) (*pbdocument.ApiResponsePaginationMerchantDocument, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	s.logger.Info("Fetching merchant document records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

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

	documents, totalRecords, err := s.merchantDocumentQuery.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAll failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapper.ToProtoResponsePaginationMerchantDocument(paginationMeta, "success", "Successfully fetched merchant documents", documents)

	s.logger.Info("Successfully fetched merchant documents", zap.Bool("success", true))

	return so, nil
}

// FindById retrieves a merchant document by its ID.
// It logs the operation's success or failure. It returns a gRPC response containing the merchant document
// or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancelation signals, and deadlines.
//   - req: A pointer to a FindMerchantDocumentByIdRequest containing the document ID.
//
// Returns:
//   - A pointer to ApiResponseMerchantDocument containing the merchant document on success.
//   - An error if the retrieval operation fails.
func (s *merchantDocumentQueryHandleGrpc) FindById(ctx context.Context, req *pbdocument.FindMerchantDocumentByIdRequest) (*pbdocument.ApiResponseMerchantDocument, error) {
	id := int(req.GetDocumentId())

	s.logger.Info("Fetching merchant document by id", zap.Int("id", id))

	if id == 0 {
		s.logger.Error("FindById failed", zap.Any("error", merchantdocument_errors.ErrGrpcMerchantInvalidID))
		return nil, merchantdocument_errors.ErrGrpcMerchantInvalidID
	}

	document, err := s.merchantDocumentQuery.FindById(ctx, id)
	if err != nil {
		s.logger.Error("FindById failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMerchantDocument("success", "Successfully fetched merchant document", document)

	s.logger.Info("Successfully fetched merchant document", zap.Bool("success", true))

	return so, nil
}

// FindAllActive retrieves a list of all active merchant documents, paginated by the given page and page size,
// and filtered by the search query. It logs the operation's success or failure. It returns a gRPC response containing
// the merchant documents or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancelation signals, and deadlines.
//   - req: A pointer to a FindAllMerchantDocumentsRequest containing the page, page size, and search query.
//
// Returns:
//   - A pointer to ApiResponsePaginationMerchantDocument containing the active merchant documents on success.
//   - An error if the retrieval operation fails.
func (s *merchantDocumentQueryHandleGrpc) FindAllActive(ctx context.Context, req *pbdocument.FindAllMerchantDocumentsRequest) (*pbdocument.ApiResponsePaginationMerchantDocument, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	s.logger.Info("Fetching active merchant document records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

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

	documents, totalRecords, err := s.merchantDocumentQuery.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAllActive failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapper.ToProtoResponsePaginationMerchantDocument(paginationMeta, "success", "Successfully fetched active merchant documents", documents)

	s.logger.Info("Successfully fetched active merchant documents", zap.Bool("success", true))

	return so, nil
}

// FindAllTrashed retrieves a list of all trashed merchant documents, paginated by the given page and page size,
// and filtered by the search query. It logs the operation's success or failure. It returns a gRPC response containing
// the trashed merchant documents or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancelation signals, and deadlines.
//   - req: A pointer to a FindAllMerchantDocumentsRequest containing the page, page size, and search query.
//
// Returns:
//   - A pointer to ApiResponsePaginationMerchantDocumentAt containing the trashed merchant documents on success.
//   - An error if the retrieval operation fails.
func (s *merchantDocumentQueryHandleGrpc) FindAllTrashed(ctx context.Context, req *pbdocument.FindAllMerchantDocumentsRequest) (*pbdocument.ApiResponsePaginationMerchantDocumentAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	s.logger.Info("Fetching trashed merchant document records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

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

	documents, totalRecords, err := s.merchantDocumentQuery.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAllTrashed failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapper.ToProtoResponsePaginationMerchantDocumentDeleteAt(paginationMeta, "success", "Successfully fetched trashed merchant documents", documents)

	s.logger.Info("Successfully fetched trashed merchant documents", zap.Bool("success", true))

	return so, nil
}
