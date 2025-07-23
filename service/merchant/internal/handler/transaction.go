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
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/merchant"
	"go.uber.org/zap"
)

type merchantTransactionHandleGrpc struct {
	pbmerchant.UnimplementedMerchantTransactionServiceServer

	merchantTransaction service.MerchantTransactionService
	logger              logger.LoggerInterface
	mapper              protomapper.MerchantTransactionProtoMapper
}

func NewMerchantTransactionHandleGrpc(merchantTransaction service.MerchantTransactionService, logger logger.LoggerInterface, mapper protomapper.MerchantTransactionProtoMapper) MerchantTransactionHandleGrpc {
	return &merchantTransactionHandleGrpc{
		merchantTransaction: merchantTransaction,
		logger:              logger,
		mapper:              mapper,
	}
}

// FindAllTransactionMerchant retrieves a paginated list of transactions for all merchants based on the provided request parameters.
// It handles pagination, search criteria, and returns a gRPC response with pagination metadata.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancelation signals, and deadlines.
//   - req: A pointer to a FindAllMerchantRequest containing pagination and search details.
//
// Returns:
//   - A pointer to ApiResponsePaginationMerchantTransaction containing the list of transactions
//     and pagination metadata on success.
//   - An error if the retrieval operation fails.
func (s *merchantTransactionHandleGrpc) FindAllTransactionMerchant(ctx context.Context, req *pbmerchant.FindAllMerchantRequest) (*pbmerchant.ApiResponsePaginationMerchantTransaction, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchantTransactions{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchants, totalRecords, err := s.merchantTransaction.FindAllTransactions(ctx, &reqService)

	if err != nil {
		s.logger.Error("FindAllTransactionMerchant failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapper.ToProtoResponsePaginationMerchantTransaction(paginationMeta, "success", "Successfully fetched merchant record", merchants)

	s.logger.Info("Successfully fetched merchant records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, nil
}

// FindAllTransactionByMerchant retrieves all transactions for a merchant by ID.
// It validates the page, page size, and search string, and returns a gRPC response containing the transactions
// or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to a FindAllMerchantTransaction containing the merchant ID, page, page size, and search string.
//
// Returns:
//   - A pointer to ApiResponsePaginationMerchantTransaction containing the transactions on success.
//   - An error if the retrieval operation fails.
func (s *merchantTransactionHandleGrpc) FindAllTransactionByMerchant(ctx context.Context, req *pbmerchant.FindAllMerchantTransaction) (*pbmerchant.ApiResponsePaginationMerchantTransaction, error) {
	merchant_id := int(req.MerchantId)
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchantTransactionsById{
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
		MerchantID: merchant_id,
	}

	merchants, totalRecords, err := s.merchantTransaction.FindAllTransactionsByMerchant(ctx, &reqService)

	if err != nil {
		s.logger.Error("FindAllTransactionByMerchant failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapper.ToProtoResponsePaginationMerchantTransaction(paginationMeta, "success", "Successfully fetched merchant record", merchants)

	s.logger.Info("Successfully fetched merchant record", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, nil
}

// FindAllTransactionByApikey retrieves paginated transaction records for a specific merchant by API key.
// It validates the API key, page, page size, and search query, and returns a gRPC response containing the
// paginated transaction records on success or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to a FindAllMerchantApikey containing the API key, page, page size, and search query.
//
// Returns:
//   - A pointer to ApiResponsePaginationMerchantTransaction containing the paginated transaction records
//     on success.
//   - An error if the retrieval operation fails.
func (s *merchantTransactionHandleGrpc) FindAllTransactionByApikey(ctx context.Context, req *pbmerchant.FindAllMerchantApikey) (*pbmerchant.ApiResponsePaginationMerchantTransaction, error) {
	api_key := req.GetApiKey()
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchantTransactionsByApiKey{
		ApiKey:   api_key,
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchants, totalRecords, err := s.merchantTransaction.FindAllTransactionsByApikey(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAllTransactionsByApikey failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapper.ToProtoResponsePaginationMerchantTransaction(paginationMeta, "success", "Successfully fetched merchant record", merchants)

	s.logger.Info("Successfully fetched merchant record", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, nil
}
