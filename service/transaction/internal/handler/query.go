package handler

import (
	"context"
	"math"

	pbhelpers "github.com/MamangRust/monolith-payment-gateway-pb"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/transaction"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/grpc"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/transaction"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/service"
	"go.uber.org/zap"
)

type transactionQueryHandleGrpc struct {
	pb.UnimplementedTransactionQueryServiceServer

	service service.TransactionQueryService
	logger  logger.LoggerInterface
	mapper  protomapper.TransactionQueryProtoMapper
}

func NewTransactionQueryHandleGrpc(service service.TransactionQueryService, logger logger.LoggerInterface, mapper protomapper.TransactionQueryProtoMapper) TransactionQueryHandleGrpc {
	return &transactionQueryHandleGrpc{
		service: service,
		logger:  logger,
		mapper:  mapper,
	}
}

// FindAllTransaction implements the gRPC service for fetching all transactions.
//
// This function takes the request containing page, page size, and search query.
// It fetches the transactions from the query service and returns a pagination
// response containing the transactions list and pagination metadata.
//
// Parameters:
//   - ctx: the context.Context of the gRPC request.
//   - request: the request containing page, page size, and search query.
//
// Returns:
//   - A pointer to the pagination response containing transactions list and pagination metadata.
//   - An error if any.
func (t *transactionQueryHandleGrpc) FindAllTransaction(ctx context.Context, request *pb.FindAllTransactionRequest) (*pb.ApiResponsePaginationTransaction, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	t.logger.Info("Fetching transaction",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := &requests.FindAllTransactions{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	transactions, totalRecords, err := t.service.FindAll(ctx, reqService)

	if err != nil {
		t.logger.Error("Failed to fetch transaction", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbhelpers.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := t.mapper.ToProtoResponsePaginationTransaction(paginationMeta, "success", "successfully fetch transaction", transactions)

	t.logger.Info("Successfully fetch transaction", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, nil
}

// FindAllTransactionByCardNumber implements the gRPC service for fetching all transactions by card number.
//
// This function takes the request containing card number, page, page size, and search query.
// It fetches the transactions from the query service and returns a pagination
// response containing the transactions list and pagination metadata.
//
// Parameters:
//   - ctx: the context.Context of the gRPC request.
//   - request: the request containing card number, page, page size, and search query.
//
// Returns:
//   - A pointer to the pagination response containing transactions list and pagination metadata.
//   - An error if any.
func (t *transactionQueryHandleGrpc) FindAllTransactionByCardNumber(ctx context.Context, request *pb.FindAllTransactionCardNumberRequest) (*pb.ApiResponsePaginationTransaction, error) {
	card_number := request.GetCardNumber()
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	t.logger.Info("Fetching transaction",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	reqService := &requests.FindAllTransactionCardNumber{
		CardNumber: card_number,
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
	}

	transactions, totalRecords, err := t.service.FindAllByCardNumber(ctx, reqService)

	if err != nil {
		t.logger.Error("Failed to fetch transaction", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbhelpers.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := t.mapper.ToProtoResponsePaginationTransaction(paginationMeta, "", "", transactions)

	return so, nil
}

// FindByIdTransaction implements the gRPC service for fetching a single transaction by id.
//
// Parameters:
//   - ctx: the context.Context of the gRPC request.
//   - req: the request containing the transaction id.
//
// Returns:
//   - A pointer to the transaction response containing the transaction details.
//   - An error if any.
func (t *transactionQueryHandleGrpc) FindByIdTransaction(ctx context.Context, req *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransaction, error) {
	id := int(req.GetTransactionId())

	t.logger.Info("Fetching transaction",
		zap.Int("transaction.id", id))

	if id == 0 {
		t.logger.Error("Failed to fetch transaction", zap.Int("transaction.id", id))
		return nil, transaction_errors.ErrGrpcTransactionInvalidID
	}

	transaction, err := t.service.FindById(ctx, id)

	if err != nil {
		t.logger.Error("Failed to fetch transaction", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapper.ToProtoResponseTransaction("success", "Transaction fetched successfully", transaction)

	t.logger.Info("Successfully fetch transaction", zap.Int("transaction.id", id))

	return so, nil
}

// FindTransactionByMerchantIdRequest retrieves a list of transactions by merchant id.
//
// Parameters:
//   - ctx: The context.Context object for the gRPC request.
//   - req: A FindTransactionByMerchantIdRequest object containing the merchant id to fetch the transactions for.
//
// Returns:
//   - An ApiResponseTransactions containing the list of transactions retrieved from the database.
//   - An error if the operation fails, or if the provided merchant id is invalid.
func (t *transactionQueryHandleGrpc) FindTransactionByMerchantIdRequest(ctx context.Context, req *pb.FindTransactionByMerchantIdRequest) (*pb.ApiResponseTransactions, error) {
	id := int(req.GetMerchantId())

	t.logger.Info("Fetching transaction",
		zap.Int("merchant_id", id))

	if id == 0 {
		t.logger.Error("Failed to fetch transaction", zap.Int("merchant_id", id))
		return nil, transaction_errors.ErrGrpcTransactionInvalidMerchantID
	}

	transactions, err := t.service.FindTransactionByMerchantId(ctx, id)

	if err != nil {
		t.logger.Error("failed to fetch transaction by merchant id", zap.Any("error", err), zap.Int("merchant_id", id))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapper.ToProtoResponseTransactions("success", "Successfully fetch transactions", transactions)

	t.logger.Info("Successfully fetch transactions", zap.Int("merchant_id", id))

	return so, nil
}

// FindByActiveTransaction retrieves a paginated list of active transactions based on the provided request parameters.
// It validates the page and page size, then constructs a request to the transaction query service.
// On success, it returns a gRPC response containing the paginated transaction data or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to a FindAllTransactionRequest containing the page, page size, and search query.
//
// Returns:
//   - A pointer to ApiResponsePaginationTransactionDeleteAt containing the paginated transaction data on success.
//   - An error if the retrieval operation fails.
func (t *transactionQueryHandleGrpc) FindByActiveTransaction(ctx context.Context, req *pb.FindAllTransactionRequest) (*pb.ApiResponsePaginationTransactionDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	t.logger.Info("Fetching transaction",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := &requests.FindAllTransactions{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	transactions, totalRecords, err := t.service.FindByActive(ctx, reqService)

	if err != nil {
		t.logger.Error("failed to fetch active transaction", zap.Any("error", err), zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbhelpers.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := t.mapper.ToProtoResponsePaginationTransactionDeleteAt(paginationMeta, "success", "Successfully fetch transactions", transactions)

	t.logger.Info("Successfully fetch transactions", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, nil
}

// FindByTrashedTransaction retrieves a paginated list of trashed transactions based on the provided request parameters.
// It validates the page and page size, then constructs a request to the transaction query service.
// On success, it returns a gRPC response containing the paginated transaction data or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to a FindAllTransactionRequest containing the page, page size, and search query.
//
// Returns:
//   - A pointer to ApiResponsePaginationTransactionDeleteAt containing the paginated transaction data on success.
//   - An error if the retrieval operation fails.
func (t *transactionQueryHandleGrpc) FindByTrashedTransaction(ctx context.Context, req *pb.FindAllTransactionRequest) (*pb.ApiResponsePaginationTransactionDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	t.logger.Info("Fetching transaction",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := &requests.FindAllTransactions{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	transactions, totalRecords, err := t.service.FindByTrashed(ctx, reqService)

	if err != nil {
		t.logger.Error("failed to fetch trashed transaction", zap.Any("error", err), zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbhelpers.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := t.mapper.ToProtoResponsePaginationTransactionDeleteAt(paginationMeta, "success", "Successfully fetch transactions", transactions)

	t.logger.Info("Successfully fetch transactions", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, nil
}
