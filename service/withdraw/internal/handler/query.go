package handler

import (
	"context"
	"math"

	pbhelper "github.com/MamangRust/monolith-payment-gateway-pb"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/withdraw"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	withdraw_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors/grpc"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/withdraw"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/service"
	"go.uber.org/zap"
)

type withdrawQueryHandleGrpc struct {
	pb.UnimplementedWithdrawQueryServiceServer

	withdrawQuery service.WithdrawQueryService

	logger logger.LoggerInterface

	mapper protomapper.WithdrawQueryProtoMapper
}

func NewWithdrawQueryHandleGrpc(
	withdrawQuery service.WithdrawQueryService,
	logger logger.LoggerInterface,
	mapper protomapper.WithdrawQueryProtoMapper,
) WithdrawQueryHandlerGrpc {
	return &withdrawQueryHandleGrpc{
		withdrawQuery: withdrawQuery,
		logger:        logger,
		mapper:        mapper,
	}
}

// FindAllWithdraw retrieves all withdraw records based on the provided request,
// which includes pagination and search parameters.
//
// Parameters:
//   - ctx: the context.Context object for tracing and cancellation.
//   - req: a pointer to a FindAllWithdrawRequest message containing pagination
//     and search parameters.
//
// Returns:
//   - A pointer to an ApiResponsePaginationWithdraw message containing the pagination
//     metadata and the list of withdraw records.
//   - An error, which is non-nil if the operation fails.
func (w *withdrawQueryHandleGrpc) FindAllWithdraw(ctx context.Context, req *pb.FindAllWithdrawRequest) (*pb.ApiResponsePaginationWithdraw, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	w.logger.Debug("FindAllWithdraw", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := &requests.FindAllWithdraws{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	withdraws, totalRecords, err := w.withdrawQuery.FindAll(ctx, reqService)

	if err != nil {
		w.logger.Error("FindAllWithdraw", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbhelper.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := w.mapper.ToProtoResponsePaginationWithdraw(paginationMeta, "success", "withdraw", withdraws)

	return so, nil
}

// FindAllWithdrawByCardNumber retrieves all withdraw records associated with a
// specific card number based on the provided request, which includes pagination
// and search parameters.
//
// Parameters:
//   - ctx: the context.Context object for tracing and cancellation.
//   - req: a pointer to a FindAllWithdrawByCardNumberRequest message containing
//     the card number, pagination, and search parameters.
//
// Returns:
//   - A pointer to an ApiResponsePaginationWithdraw message containing the
//     pagination metadata and the list of withdraw records.
//   - An error, which is non-nil if the operation fails.
func (w *withdrawQueryHandleGrpc) FindAllWithdrawByCardNumber(ctx context.Context, req *pb.FindAllWithdrawByCardNumberRequest) (*pb.ApiResponsePaginationWithdraw, error) {
	card_number := req.GetCardNumber()
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	w.logger.Debug("FindAllWithdrawByCardNumber", zap.String("card_number", card_number), zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := &requests.FindAllWithdrawCardNumber{
		CardNumber: card_number,
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
	}

	withdraws, totalRecords, err := w.withdrawQuery.FindAllByCardNumber(ctx, reqService)

	if err != nil {
		w.logger.Error("FindAllWithdrawByCardNumber", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbhelper.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := w.mapper.ToProtoResponsePaginationWithdraw(paginationMeta, "success", "Withdraws fetched successfully", withdraws)

	return so, nil
}

// FindByIdWithdraw retrieves a withdraw record based on the provided request,
// which includes the withdraw ID.
//
// Parameters:
//   - ctx: the context.Context object for tracing and cancellation.
//   - req: a pointer to a FindByIdWithdrawRequest message containing the
//     withdraw ID.
//
// Returns:
//   - A pointer to an ApiResponseWithdraw message containing the withdraw
//     record.
//   - An error, which is non-nil if the operation fails.
func (w *withdrawQueryHandleGrpc) FindByIdWithdraw(ctx context.Context, req *pb.FindByIdWithdrawRequest) (*pb.ApiResponseWithdraw, error) {
	id := int(req.GetWithdrawId())

	w.logger.Debug("FindByIdWithdraw", zap.Int("id", id))

	if id == 0 {
		w.logger.Error("FindByIdWithdraw", zap.Any("error", withdraw_errors.ErrGrpcWithdrawInvalidID))
		return nil, withdraw_errors.ErrGrpcWithdrawInvalidID
	}

	withdraw, err := w.withdrawQuery.FindById(ctx, id)

	if err != nil {
		w.logger.Error("FindByIdWithdraw", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := w.mapper.ToProtoResponseWithdraw("success", "Successfully fetched withdraw", withdraw)

	return so, nil
}

// FindByActive retrieves active withdraw records with pagination and search criteria.
//
// Parameters:
//   - ctx: the context.Context object for request-scoped values, cancellation, and deadlines.
//   - req: a pointer to a FindAllWithdrawRequest message containing pagination and search parameters.
//
// Returns:
//   - A pointer to an ApiResponsePaginationWithdrawDeleteAt message containing the list of active withdraw records and pagination metadata.
//   - An error, which is non-nil if the operation fails.
func (w *withdrawQueryHandleGrpc) FindByActive(ctx context.Context, req *pb.FindAllWithdrawRequest) (*pb.ApiResponsePaginationWithdrawDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	w.logger.Debug("FindByActive", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := &requests.FindAllWithdraws{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := w.withdrawQuery.FindByActive(ctx, reqService)

	if err != nil {
		w.logger.Error("FindByActive", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbhelper.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := w.mapper.ToProtoResponsePaginationWithdrawDeleteAt(paginationMeta, "success", "Successfully fetched withdraws", res)

	return so, nil
}

// FindByTrashed retrieves trashed withdraw records with pagination and search criteria.
//
// Parameters:
//   - ctx: the context.Context object for request-scoped values, cancellation, and deadlines.
//   - req: a pointer to a FindAllWithdrawRequest message containing pagination and search parameters.
//
// Returns:
//   - A pointer to an ApiResponsePaginationWithdrawDeleteAt message containing the list of trashed withdraw records and pagination metadata.
//   - An error, which is non-nil if the operation fails.
func (w *withdrawQueryHandleGrpc) FindByTrashed(ctx context.Context, req *pb.FindAllWithdrawRequest) (*pb.ApiResponsePaginationWithdrawDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	w.logger.Debug("FindByTrashed", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := &requests.FindAllWithdraws{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := w.withdrawQuery.FindByTrashed(ctx, reqService)

	if err != nil {
		w.logger.Error("FindByTrashed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbhelper.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := w.mapper.ToProtoResponsePaginationWithdrawDeleteAt(paginationMeta, "success", "Successfully fetched withdraws", res)

	return so, nil
}
