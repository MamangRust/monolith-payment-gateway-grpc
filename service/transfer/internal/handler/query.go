package handler

import (
	"context"
	"math"

	pbhelper "github.com/MamangRust/monolith-payment-gateway-pb"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/transfer"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/grpc"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/transfer"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/service"
	"go.uber.org/zap"
)

// transferQueryHandleGrpc represents the gRPC handler for transfer operations.
type transferQueryHandleGrpc struct {
	pb.UnimplementedTransferQueryServiceServer

	transferQueryService service.TransferQueryService
	logger               logger.LoggerInterface
	mapper               protomapper.TransferQueryProtoMapper
}

func NewTransferQueryHandler(service service.TransferQueryService, logger logger.LoggerInterface, mapper protomapper.TransferQueryProtoMapper) TransferQueryHandleGrpc {
	return &transferQueryHandleGrpc{
		transferQueryService: service,
		logger:               logger,
		mapper:               mapper,
	}
}

// FindAllTransfer implements the gRPC service for fetching all transfer records.
//
// This function handles the "FindAllTransfer" gRPC request by calling the
// underlying transfer query service with the provided pagination parameters
// and search query. It returns a protobuf response containing the transfer
// records, pagination metadata, and a success message.
//
// Parameters:
//   - ctx: the context object for the gRPC request.
//   - request: a protobuf request containing the pagination parameters and search query.
//
// Returns:
//   - A protobuf response containing the transfer records, pagination metadata, and a success message.
//   - An error object if the request fails for any reason.
func (s *transferQueryHandleGrpc) FindAllTransfer(ctx context.Context, request *pb.FindAllTransferRequest) (*pb.ApiResponsePaginationTransfer, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	s.logger.Info("Fetching transfer", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := &requests.FindAllTransfers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchants, totalRecords, err := s.transferQueryService.FindAll(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbhelper.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapper.ToProtoResponsePaginationTransfer(paginationMeta, "success", "Successfully fetch transfer records", merchants)

	s.logger.Info("Successfully fetch transfer records", zap.Int("totalRecords", *totalRecords))

	return so, nil
}

// FindByIdTransfer retrieves a transfer by its ID.
// It validates the ID, then constructs a request to the transfer query service.
// On success, it returns a gRPC response containing the transfer data or an error if the retrieval operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - request: A pointer to a FindByIdTransferRequest containing the transfer ID.
//
// Returns:
//   - A pointer to ApiResponseTransfer containing the transfer data on success.
//   - An error if the retrieval operation fails.
func (s *transferQueryHandleGrpc) FindByIdTransfer(ctx context.Context, request *pb.FindByIdTransferRequest) (*pb.ApiResponseTransfer, error) {
	id := int(request.GetTransferId())

	s.logger.Info("Fetching transfer", zap.Int("id", id))

	if id == 0 {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", transfer_errors.ErrGrpcTransferInvalidID))
		return nil, transfer_errors.ErrGrpcTransferInvalidID
	}

	transfer, err := s.transferQueryService.FindById(ctx, id)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTransfer("success", "Successfully fetch transfer record", transfer)

	s.logger.Info("Successfully fetch transfer record", zap.Int("id", id))

	return so, nil
}

// FindByTransferByTransferFrom retrieves the transfer records with the specified transfer_from value.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - request: A FindTransferByTransferFromRequest object containing the card number
//     for which to fetch the transfer records.
//
// Returns:
//   - An ApiResponseTransfers containing the transfer records.
//   - An error if the operation fails, or if the provided card number is invalid.
func (s *transferQueryHandleGrpc) FindByTransferByTransferFrom(ctx context.Context, request *pb.FindTransferByTransferFromRequest) (*pb.ApiResponseTransfers, error) {
	transfer_from := request.GetTransferFrom()

	s.logger.Info("Fetching transfer", zap.String("transfer_from", transfer_from))

	if transfer_from == "" {
		s.logger.Error("Failed to fetch transfer", zap.String("transfer_from", transfer_from))
		return nil, transfer_errors.ErrGrpcInvalidCardNumber
	}

	merchants, err := s.transferQueryService.FindTransferByTransferFrom(ctx, transfer_from)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTransfers("success", "Successfully fetch transfer records", merchants)

	s.logger.Info("Successfully fetched transfer records", zap.String("transfer_from", transfer_from))

	return so, nil
}

// FindByTransferByTransferTo retrieves the transfer records with the specified transfer_to value.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - request: A FindTransferByTransferToRequest object containing the card number
//     for which to fetch the transfer records.
//
// Returns:
//   - An ApiResponseTransfers containing the transfer records.
//   - An error if the operation fails, or if the provided card number is invalid.
func (s *transferQueryHandleGrpc) FindByTransferByTransferTo(ctx context.Context, request *pb.FindTransferByTransferToRequest) (*pb.ApiResponseTransfers, error) {
	transfer_to := request.GetTransferTo()

	s.logger.Info("Fetching transfer", zap.String("transfer_to", transfer_to))

	if transfer_to == "" {
		s.logger.Error("Failed to fetch transfer", zap.String("transfer_to", transfer_to))
		return nil, transfer_errors.ErrGrpcInvalidCardNumber
	}

	merchants, err := s.transferQueryService.FindTransferByTransferTo(ctx, transfer_to)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTransfers("success", "Successfully fetch transfer records", merchants)

	s.logger.Info("Successfully fetched transfer records", zap.String("transfer_to", transfer_to))

	return so, nil
}

// FindByActiveTransfer retrieves a paginated list of active transfer records based on the provided request parameters.
// It validates the page and page size, then constructs a request to the transfer query service.
// On success, it returns a gRPC response containing the paginated transfer data or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to a FindAllTransferRequest containing the page, page size, and search query.
//
// Returns:
//   - A pointer to ApiResponsePaginationTransferDeleteAt containing the paginated transfer data on success.
//   - An error if the retrieval operation fails.
func (s *transferQueryHandleGrpc) FindByActiveTransfer(ctx context.Context, req *pb.FindAllTransferRequest) (*pb.ApiResponsePaginationTransferDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	s.logger.Info("Fetching transfer",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := &requests.FindAllTransfers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.transferQueryService.FindByActive(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbhelper.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapper.ToProtoResponsePaginationTransferDeleteAt(paginationMeta, "success", "Successfully fetch transfer records", res)

	s.logger.Info("Successfully fetched transfer records", zap.String("search", search))

	return so, nil
}

// FindByTrashedTransfer retrieves a paginated list of trashed transfer records based on the request parameters.
// It validates pagination inputs, constructs a request to the transfer query service, and processes the response.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation, and deadlines.
//   - req: A pointer to a FindAllTransferRequest containing pagination and search query.
//
// Returns:
//   - A pointer to ApiResponsePaginationTransferDeleteAt with the paginated trashed transfer records.
//   - An error if the operation fails.
func (s *transferQueryHandleGrpc) FindByTrashedTransfer(ctx context.Context, req *pb.FindAllTransferRequest) (*pb.ApiResponsePaginationTransferDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	s.logger.Info("Fetching transfer",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := &requests.FindAllTransfers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.transferQueryService.FindByTrashed(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}
	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbhelper.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapper.ToProtoResponsePaginationTransferDeleteAt(paginationMeta, "success", "Successfully fetch transfer records", res)

	s.logger.Info("Successfully fetched transfer records", zap.String("search", search))

	return so, nil
}
