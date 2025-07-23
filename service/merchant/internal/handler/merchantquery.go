package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/service"
	pb "github.com/MamangRust/monolith-payment-gateway-pb"
	pbmerchant "github.com/MamangRust/monolith-payment-gateway-pb/merchant"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/grpc"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/merchant"
	"go.uber.org/zap"
)

type merchantQueryHandleGrpc struct {
	pbmerchant.UnimplementedMerchantQueryServiceServer

	merchantQuery service.MerchantQueryService
	logger        logger.LoggerInterface
	mapper        protomapper.MerchantQueryProtoMapper
}

func NewMerchantQueryHandleGrpc(merchantQuery service.MerchantQueryService, logger logger.LoggerInterface, mapper protomapper.MerchantQueryProtoMapper) MerchantQueryHandleGrpc {
	return &merchantQueryHandleGrpc{
		merchantQuery: merchantQuery,
		logger:        logger,
		mapper:        mapper,
	}
}

// FindAllMerchant retrieves a paginated list of merchants based on the provided request parameters.
// It handles pagination, search criteria, and returns a gRPC response with pagination metadata.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancelation signals, and deadlines.
//   - req: A pointer to a FindAllMerchantRequest containing pagination and search details.
//
// Returns:
//   - A pointer to ApiResponsePaginationMerchant containing the list of merchants
//     and pagination metadata on success.
//   - An error if the retrieval operation fails.
func (s *merchantQueryHandleGrpc) FindAllMerchant(ctx context.Context, req *pbmerchant.FindAllMerchantRequest) (*pbmerchant.ApiResponsePaginationMerchant, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchants{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchants, totalRecords, err := s.merchantQuery.FindAll(ctx, &reqService)

	if err != nil {
		s.logger.Error("FindAllMerchant failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapper.ToProtoResponsePaginationMerchant(paginationMeta, "success", "Successfully fetched merchant record", merchants)

	s.logger.Info("Successfully fetched merchant records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, nil
}

// FindByIdMerchant retrieves a merchant by its ID.
// It handles invalid merchant IDs and returns a gRPC response containing the merchant data
// or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancelation signals, and deadlines.
//   - req: A pointer to a FindByIdMerchantRequest containing the document ID.
//
// Returns:
//   - A pointer to ApiResponseMerchant containing the merchant data on success.
//   - An error if the retrieval operation fails.
func (s *merchantQueryHandleGrpc) FindByIdMerchant(ctx context.Context, req *pbmerchant.FindByIdMerchantRequest) (*pbmerchant.ApiResponseMerchant, error) {
	id := int(req.GetMerchantId())

	if id == 0 {
		s.logger.Error("invalid id failed", zap.Any("error", merchant_errors.ErrGrpcMerchantInvalidID))
		return nil, merchant_errors.ErrGrpcMerchantInvalidID
	}

	merchant, err := s.merchantQuery.FindById(ctx, id)

	if err != nil {
		s.logger.Error("FindByIdMerchant failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMerchant("success", "Successfully fetched merchant record", merchant)

	s.logger.Info("Successfully fetched merchant record", zap.Bool("success", true))

	return so, nil
}

// FindByApiKey retrieves a merchant's information using the provided API key.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to a FindByApiKeyRequest containing the API key.
//
// Returns:
//   - A pointer to ApiResponseMerchant containing the merchant data on success.
//   - An error if the retrieval operation fails or if the provided API key is invalid.
func (s *merchantQueryHandleGrpc) FindByApiKey(ctx context.Context, req *pbmerchant.FindByApiKeyRequest) (*pbmerchant.ApiResponseMerchant, error) {
	api_key := req.GetApiKey()

	if api_key == "" {
		s.logger.Error("invalid api key failed", zap.Any("error", merchant_errors.ErrGrpcMerchantInvalidApiKey))
		return nil, merchant_errors.ErrGrpcMerchantInvalidApiKey
	}

	merchant, err := s.merchantQuery.FindByApiKey(ctx, api_key)

	if err != nil {
		s.logger.Error("FindByApiKey failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMerchant("success", "Successfully fetched merchant record", merchant)

	s.logger.Info("Successfully fetched merchant record", zap.Bool("success", true))

	return so, nil
}

// FindByMerchantUserId retrieves merchant information based on the provided user ID.
// It validates the user ID and returns a gRPC response containing the merchant data
// on success or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to a FindByMerchantUserIdRequest containing the user ID.
//
// Returns:
//   - A pointer to ApiResponsesMerchant containing the merchant data on success.
//   - An error if the retrieval operation fails or if the provided user ID is invalid.
func (s *merchantQueryHandleGrpc) FindByMerchantUserId(ctx context.Context, req *pbmerchant.FindByMerchantUserIdRequest) (*pbmerchant.ApiResponsesMerchant, error) {
	user_id := req.GetUserId()

	if user_id <= 0 {
		s.logger.Error("invalid user id failed", zap.Any("error", merchant_errors.ErrGrpcMerchantInvalidUserID))
		return nil, merchant_errors.ErrGrpcMerchantInvalidUserID
	}

	res, err := s.merchantQuery.FindByMerchantUserId(ctx, int(user_id))

	if err != nil {
		s.logger.Error("FindByMerchantUserId failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMerchants("success", "Successfully fetched merchant record", res)

	return so, nil
}

// FindByActive retrieves a paginated list of active merchants based on the provided request parameters.
// It validates the page and page size, then constructs a request to the merchant query service.
// On success, it returns a gRPC response containing the paginated merchant data or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to a FindAllMerchantRequest containing the page, page size, and search query.
//
// Returns:
//   - A pointer to ApiResponsePaginationMerchantDeleteAt containing the paginated merchant data on success.
//   - An error if the retrieval operation fails.
func (s *merchantQueryHandleGrpc) FindByActive(ctx context.Context, req *pbmerchant.FindAllMerchantRequest) (*pbmerchant.ApiResponsePaginationMerchantDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchants{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.merchantQuery.FindByActive(ctx, &reqService)

	if err != nil {
		s.logger.Error("FindByActive failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapper.ToProtoResponsePaginationMerchantDeleteAt(paginationMeta, "success", "Successfully fetched merchant record", res)

	s.logger.Info("Successfully fetched merchant record", zap.Bool("success", true))

	return so, nil
}

// FindByTrashed retrieves a paginated list of trashed merchants based on the provided request parameters.
// It validates the page and page size, then constructs a request to the merchant query service.
// On success, it returns a gRPC response containing the paginated merchant data or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to a FindAllMerchantRequest containing the page, page size, and search query.
//
// Returns:
//   - A pointer to ApiResponsePaginationMerchantDeleteAt containing the paginated merchant data on success.
//   - An error if the retrieval operation fails.
func (s *merchantQueryHandleGrpc) FindByTrashed(ctx context.Context, req *pbmerchant.FindAllMerchantRequest) (*pbmerchant.ApiResponsePaginationMerchantDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchants{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.merchantQuery.FindByTrashed(ctx, &reqService)

	if err != nil {
		s.logger.Error("FindByTrashed failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapper.ToProtoResponsePaginationMerchantDeleteAt(paginationMeta, "success", "Successfully fetched merchant record", res)

	s.logger.Info("Successfully fetched merchant record", zap.Bool("success", true))

	return so, nil
}
