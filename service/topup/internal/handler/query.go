package handler

import (
	"context"
	"math"

	pbhelpers "github.com/MamangRust/monolith-payment-gateway-pb"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/topup"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/grpc"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/topup"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/service"
	"go.uber.org/zap"
)

type topupQueryHandleGrpc struct {
	pb.UnimplementedTopupQueryServiceServer

	service service.TopupQueryService

	mapper protomapper.TopupQueryProtoMapper

	logger logger.LoggerInterface
}

func NewTopupQueryHandleGrpc(service service.TopupQueryService, logger logger.LoggerInterface, mapper protomapper.TopupQueryProtoMapper) TopupQueryHandleGrpc {
	return &topupQueryHandleGrpc{
		service: service,
		mapper:  mapper,
		logger:  logger,
	}
}

// FindAllTopup is a gRPC handler function that fetches all topups using the topup query service.
//
// Parameters:
//   - ctx: the context.Context object passed through the gRPC request.
//   - req: a pointer to a FindAllTopupRequest message, containing the page, pageSize, and search query.
//
// Returns:
//   - A pointer to an ApiResponsePaginationTopup message, containing the list of topups and pagination metadata.
//   - An error, if the topup query service returns an error.
func (s *topupQueryHandleGrpc) FindAllTopup(ctx context.Context, req *pb.FindAllTopupRequest) (*pb.ApiResponsePaginationTopup, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	s.logger.Info("Fetching topup",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := &requests.FindAllTopups{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	topups, totalRecords, err := s.service.FindAll(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to fetch topup", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbhelpers.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapper.ToProtoResponsePaginationTopup(paginationMeta, "success", "Successfully fetch topups", topups)

	s.logger.Info("Successfully fetched topups",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, nil
}

// FindAllTopupByCardNumber fetches top-up records associated with a specific card number.
// It uses pagination and search criteria to filter the results.
//
// Parameters:
//   - ctx: the context.Context object passed through the gRPC request.
//   - req: a pointer to a FindAllTopupByCardNumberRequest message, containing the card number,
//     page, pageSize, and search query.
//
// Returns:
//   - A pointer to an ApiResponsePaginationTopup message, containing the list of top-ups and pagination metadata.
//   - An error, if the topup query service returns an error.
func (s *topupQueryHandleGrpc) FindAllTopupByCardNumber(ctx context.Context, req *pb.FindAllTopupByCardNumberRequest) (*pb.ApiResponsePaginationTopup, error) {
	card_number := req.GetCardNumber()
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	s.logger.Info("Fetching topup by card number",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search),
		zap.String("card_number", card_number))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := &requests.FindAllTopupsByCardNumber{
		CardNumber: card_number,
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
	}

	topups, totalRecords, err := s.service.FindAllByCardNumber(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to fetch topup", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbhelpers.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapper.ToProtoResponsePaginationTopup(paginationMeta, "success", "Successfully fetch topups", topups)

	s.logger.Info("Successfully fetched topups",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, nil
}

// FindByIdTopup retrieves a top-up record associated with a specific ID.
//
// Parameters:
//   - ctx: the context.Context object passed through the gRPC request.
//   - req: a pointer to a FindByIdTopupRequest message containing the top-up ID.
//
// Returns:
//   - A pointer to an ApiResponseTopup message containing the top-up details.
//   - An error, if the top-up query service returns an error or if the ID is invalid.
func (s *topupQueryHandleGrpc) FindByIdTopup(ctx context.Context, req *pb.FindByIdTopupRequest) (*pb.ApiResponseTopup, error) {
	id := int(req.GetTopupId())

	s.logger.Info("Fetching topup by id",
		zap.Int("id", id))

	if id == 0 {
		s.logger.Error("Failed to fetch topup by id", zap.Int("id", id))
		return nil, topup_errors.ErrGrpcTopupInvalidID
	}

	topup, err := s.service.FindById(ctx, id)

	if err != nil {
		s.logger.Error("Failed to fetch topup by id", zap.Int("id", id), zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseTopup("success", "Successfully fetch topup", topup)

	s.logger.Info("Successfully fetched topup by id",
		zap.Int("id", id))

	return so, nil
}

// FindByActive fetches active topup records using pagination and search criteria.
//
// Parameters:
//   - ctx: the context.Context object passed through the gRPC request.
//   - req: a pointer to a FindAllTopupRequest message, containing the page, pageSize, and search query.
//
// Returns:
//   - A pointer to an ApiResponsePaginationTopupDeleteAt message, containing the list of topups and pagination metadata.
//   - An error, if the topup query service returns an error.
func (s *topupQueryHandleGrpc) FindByActive(ctx context.Context, req *pb.FindAllTopupRequest) (*pb.ApiResponsePaginationTopupDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	s.logger.Debug("Fetching active topup",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := &requests.FindAllTopups{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.service.FindByActive(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to fetch topup", zap.Any("error", err))

		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbhelpers.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapper.ToProtoResponsePaginationTopupDeleteAt(paginationMeta, "success", "Successfully fetch topups", res)

	s.logger.Info("Successfully fetched active topup",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	return so, nil
}

// FindByTrashed fetches trashed topup records using pagination and search criteria.
//
// Parameters:
//   - ctx: the context.Context object passed through the gRPC request.
//   - req: a pointer to a FindAllTopupRequest message, containing the page, pageSize, and search query.
//
// Returns:
//   - A pointer to an ApiResponsePaginationTopupDeleteAt message, containing the list of topups and pagination metadata.
//   - An error, if the topup query service returns an error.
func (s *topupQueryHandleGrpc) FindByTrashed(ctx context.Context, req *pb.FindAllTopupRequest) (*pb.ApiResponsePaginationTopupDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	s.logger.Info("Fetching trashed topup",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := &requests.FindAllTopups{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.service.FindByTrashed(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to fetch topup", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbhelpers.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapper.ToProtoResponsePaginationTopupDeleteAt(paginationMeta, "success", "Successfully fetch topups", res)

	s.logger.Info("Successfully fetch topups",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	return so, nil
}
