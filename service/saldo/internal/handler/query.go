package handler

import (
	"context"
	"math"

	pbhelpers "github.com/MamangRust/monolith-payment-gateway-pb"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/saldo"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/service"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/grpc"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/saldo"
	"go.uber.org/zap"
)

type saldoQueryHandleGrpc struct {
	pb.UnimplementedSaldoQueryServiceServer

	service service.SaldoQueryService
	logger  logger.LoggerInterface
	mapper  protomapper.SaldoQueryProtoMapper
}

func NewSaldoQueryHandleGrpc(query service.SaldoQueryService, logger logger.LoggerInterface, mapper protomapper.SaldoQueryProtoMapper) SaldoQueryHandleGrpc {
	return &saldoQueryHandleGrpc{
		service: query,
		logger:  logger,
		mapper:  mapper,
	}
}

// FindAllSaldo is a gRPC handler that fetches all saldo records according to the given request.
//
// Parameters:
//   - ctx: the context.Context object for tracing and cancellation.
//   - req: a pointer to a FindAllSaldoRequest message, which contains the pagination and search parameters.
//
// Returns:
//   - A pointer to a ApiResponsePaginationSaldo message, which contains the pagination metadata and the fetched saldo records.
//   - An error, which is non-nil if the operation fails.
func (s *saldoQueryHandleGrpc) FindAllSaldo(ctx context.Context, req *pb.FindAllSaldoRequest) (*pb.ApiResponsePaginationSaldo, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	s.logger.Info("Fetching saldo records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllSaldos{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.service.FindAll(ctx, &reqService)

	if err != nil {
		s.logger.Error("FindAll failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbhelpers.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapper.ToProtoResponsePaginationSaldo(paginationMeta, "success", "Successfully fetched saldo record", res)

	s.logger.Info("Successfully fetched saldo records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, nil
}

// FindByIdSaldo is a gRPC handler that fetches a saldo record by its ID.
//
// Parameters:
//   - ctx: the context.Context object for tracing and cancellation.
//   - req: a pointer to a FindByIdSaldoRequest message, which contains the ID of the saldo record to be fetched.
//
// Returns:
//   - A pointer to a ApiResponseSaldo message, which contains the fetched saldo record.
//   - An error, which is non-nil if the operation fails.
func (s *saldoQueryHandleGrpc) FindByIdSaldo(ctx context.Context, req *pb.FindByIdSaldoRequest) (*pb.ApiResponseSaldo, error) {
	id := int(req.GetSaldoId())

	s.logger.Info("Fetching saldo record", zap.Int("id", id))

	if id == 0 {
		s.logger.Error("FindById failed", zap.Any("error", saldo_errors.ErrGrpcSaldoInvalidID))
		return nil, saldo_errors.ErrGrpcSaldoInvalidID
	}

	saldo, err := s.service.FindById(ctx, id)

	if err != nil {
		s.logger.Error("FindById failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseSaldo("success", "Successfully fetched saldo record", saldo)

	s.logger.Info("Successfully fetched saldo record", zap.Bool("success", true))

	return so, nil
}

// FindByCardNumber is a gRPC handler that fetches a saldo record by its card number.
//
// Parameters:
//   - ctx: the context.Context object for tracing and cancellation.
//   - req: a pointer to a FindByCardNumberRequest message, which contains the card number of the saldo record to be fetched.
//
// Returns:
//   - A pointer to a ApiResponseSaldo message, which contains the fetched saldo record.
//   - An error, which is non-nil if the operation fails.
func (s *saldoQueryHandleGrpc) FindByCardNumber(ctx context.Context, req *pbhelpers.FindByCardNumberRequest) (*pb.ApiResponseSaldo, error) {
	cardNumber := req.GetCardNumber()

	s.logger.Info("Fetching saldo records", zap.String("card_number", cardNumber))

	if cardNumber == "" {
		s.logger.Error("FindByCardNumber failed", zap.Any("error", saldo_errors.ErrGrpcSaldoInvalidCardNumber))
		return nil, saldo_errors.ErrGrpcSaldoInvalidCardNumber
	}

	saldo, err := s.service.FindByCardNumber(ctx, cardNumber)

	if err != nil {
		s.logger.Error("FindByCardNumber failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseSaldo("success", "Successfully fetched saldo record", saldo)

	s.logger.Info("Successfully fetched saldo record", zap.Bool("success", true))

	return so, nil
}

// FindByActive is a gRPC handler that fetches active saldo records according to the given request.
//
// Parameters:
//   - ctx: the context.Context object for tracing and cancellation.
//   - req: a pointer to a FindAllSaldoRequest message, which contains the pagination and search parameters.
//
// Returns:
//   - A pointer to a ApiResponsePaginationSaldoDeleteAt message, which contains the pagination metadata and the fetched saldo records.
//   - An error, which is non-nil if the operation fails.
func (s *saldoQueryHandleGrpc) FindByActive(ctx context.Context, req *pb.FindAllSaldoRequest) (*pb.ApiResponsePaginationSaldoDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	s.logger.Info("Fetching active saldo records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := &requests.FindAllSaldos{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.service.FindByActive(ctx, reqService)

	if err != nil {
		s.logger.Error("FindByActive failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbhelpers.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapper.ToProtoResponsePaginationSaldoDeleteAt(paginationMeta, "success", "Successfully fetched saldo record", res)

	s.logger.Info("Successfully fetched active saldo records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, nil
}

// FindByTrashed is a gRPC handler that fetches all trashed saldo records according to the given request.
//
// Parameters:
//   - ctx: the context.Context object for tracing and cancellation.
//   - req: a pointer to a FindAllSaldoRequest message, which contains the pagination and search parameters.
//
// Returns:
//   - A pointer to a ApiResponsePaginationSaldoDeleteAt message, which contains the pagination metadata and the fetched saldo records.
//   - An error, which is non-nil if the operation fails.
func (s *saldoQueryHandleGrpc) FindByTrashed(ctx context.Context, req *pb.FindAllSaldoRequest) (*pb.ApiResponsePaginationSaldoDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	s.logger.Info("Fetching trashed saldo records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := &requests.FindAllSaldos{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.service.FindByTrashed(ctx, reqService)

	if err != nil {
		s.logger.Error("FindByTrashed failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbhelpers.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapper.ToProtoResponsePaginationSaldoDeleteAt(paginationMeta, "success", "Successfully fetched saldo record", res)

	s.logger.Info("Successfully fetched trashed saldo record", zap.Bool("success", true))

	return so, nil
}
