package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/service"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto"
	"github.com/MamangRust/monolith-payment-gateway-shared/pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type saldoHandleGrpc struct {
	pb.UnimplementedSaldoServiceServer
	saldoQueryService     service.SaldoQueryService
	saldoStatisticService service.SaldoStatisticService
	saldoCommandService   service.SaldoCommandService
	mapping               protomapper.SaldoProtoMapper
	logger                logger.LoggerInterface
}

func NewSaldoHandleGrpc(service service.Service, logger logger.LoggerInterface) *saldoHandleGrpc {
	return &saldoHandleGrpc{
		saldoQueryService:     service.SaldoQuery,
		saldoStatisticService: service.SaldoStats,
		saldoCommandService:   service.SaldoCommand,
		mapping:               protomapper.NewSaldoProtoMapper(),
		logger:                logger,
	}
}

func (s *saldoHandleGrpc) FindAllSaldo(ctx context.Context, req *pb.FindAllSaldoRequest) (*pb.ApiResponsePaginationSaldo, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	s.logger.Debug("Fetching saldo records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

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

	res, totalRecords, err := s.saldoQueryService.FindAll(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationSaldo(paginationMeta, "success", "Successfully fetched saldo record", res)

	return so, nil
}

func (s *saldoHandleGrpc) FindByIdSaldo(ctx context.Context, req *pb.FindByIdSaldoRequest) (*pb.ApiResponseSaldo, error) {
	id := int(req.GetSaldoId())

	s.logger.Debug("Fetching saldo record", zap.Int("id", id))

	if id == 0 {
		s.logger.Debug("FindById failed", zap.Any("error", saldo_errors.ErrGrpcSaldoInvalidID))
		return nil, saldo_errors.ErrGrpcSaldoInvalidID
	}

	saldo, err := s.saldoQueryService.FindById(id)

	if err != nil {
		s.logger.Debug("FindById failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseSaldo("success", "Successfully fetched saldo record", saldo)

	return so, nil
}

func (s *saldoHandleGrpc) FindMonthlyTotalSaldoBalance(ctx context.Context, req *pb.FindMonthlySaldoTotalBalance) (*pb.ApiResponseMonthTotalSaldo, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	s.logger.Debug("Fetching monthly total saldo balance", zap.Int("year", year), zap.Int("month", month))

	if year <= 0 {
		s.logger.Debug("FindMonthlyTotalSaldoBalance failed", zap.Any("error", saldo_errors.ErrGrpcSaldoInvalidYear))
		return nil, saldo_errors.ErrGrpcSaldoInvalidYear
	}

	if month <= 0 {
		s.logger.Debug("FindMonthlyTotalSaldoBalance failed", zap.Any("error", saldo_errors.ErrGrpcSaldoInvalidMonth))
		return nil, saldo_errors.ErrGrpcSaldoInvalidMonth
	}

	reqService := requests.MonthTotalSaldoBalance{
		Year:  year,
		Month: month,
	}

	res, err := s.saldoStatisticService.FindMonthlyTotalSaldoBalance(&reqService)

	if err != nil {
		s.logger.Debug("FindMonthlyTotalSaldoBalance failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMonthTotalSaldo("success", "Successfully fetched monthly total saldo balance", res)

	return so, nil
}

func (s *saldoHandleGrpc) FindYearTotalSaldoBalance(ctx context.Context, req *pb.FindYearlySaldo) (*pb.ApiResponseYearTotalSaldo, error) {
	year := int(req.GetYear())

	s.logger.Debug("Fetching yearly total saldo balance", zap.Int("year", year))

	if year <= 0 {
		s.logger.Debug("FindYearTotalSaldoBalance failed", zap.Any("error", saldo_errors.ErrGrpcSaldoInvalidYear))
		return nil, saldo_errors.ErrGrpcSaldoInvalidYear
	}

	res, err := s.saldoStatisticService.FindYearTotalSaldoBalance(year)

	if err != nil {
		s.logger.Debug("FindYearTotalSaldoBalance failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseYearTotalSaldo("success", "Successfully fetched yearly total saldo balance", res)

	return so, nil
}

func (s *saldoHandleGrpc) FindMonthlySaldoBalances(ctx context.Context, req *pb.FindYearlySaldo) (*pb.ApiResponseMonthSaldoBalances, error) {
	year := int(req.GetYear())

	if year <= 0 {
		s.logger.Debug("FindMonthlySaldoBalances failed", zap.Any("error", saldo_errors.ErrGrpcSaldoInvalidYear))
		return nil, saldo_errors.ErrGrpcSaldoInvalidYear
	}

	res, err := s.saldoStatisticService.FindMonthlySaldoBalances(year)

	if err != nil {
		s.logger.Debug("FindMonthlySaldoBalances failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMonthSaldoBalances("success", "Successfully fetched monthly saldo balances", res)

	return so, nil
}

func (s *saldoHandleGrpc) FindYearlySaldoBalances(ctx context.Context, req *pb.FindYearlySaldo) (*pb.ApiResponseYearSaldoBalances, error) {
	year := int(req.GetYear())

	if year <= 0 {
		s.logger.Debug("FindYearlySaldoBalances failed", zap.Any("error", saldo_errors.ErrGrpcSaldoInvalidYear))
		return nil, saldo_errors.ErrGrpcSaldoInvalidYear
	}

	res, err := s.saldoStatisticService.FindYearlySaldoBalances(year)

	if err != nil {
		s.logger.Debug("FindYearlySaldoBalances failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseYearSaldoBalances("success", "Successfully fetched yearly saldo balances", res)

	return so, nil
}

func (s *saldoHandleGrpc) FindByCardNumber(ctx context.Context, req *pb.FindByCardNumberRequest) (*pb.ApiResponseSaldo, error) {
	cardNumber := req.GetCardNumber()

	s.logger.Debug("Fetching saldo records", zap.String("card_number", cardNumber))

	if cardNumber == "" {
		s.logger.Debug("FindByCardNumber failed", zap.Any("error", saldo_errors.ErrGrpcSaldoInvalidCardNumber))
		return nil, saldo_errors.ErrGrpcSaldoInvalidCardNumber
	}

	saldo, err := s.saldoQueryService.FindByCardNumber(cardNumber)

	if err != nil {
		s.logger.Debug("FindByCardNumber failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseSaldo("success", "Successfully fetched saldo record", saldo)

	return so, nil
}

func (s *saldoHandleGrpc) FindByActive(ctx context.Context, req *pb.FindAllSaldoRequest) (*pb.ApiResponsePaginationSaldoDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	s.logger.Debug("Fetching active saldo records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

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

	res, totalRecords, err := s.saldoQueryService.FindByActive(&reqService)

	if err != nil {
		s.logger.Debug("FindByActive failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapping.ToProtoResponsePaginationSaldoDeleteAt(paginationMeta, "success", "Successfully fetched saldo record", res)

	return so, nil
}

func (s *saldoHandleGrpc) FindByTrashed(ctx context.Context, req *pb.FindAllSaldoRequest) (*pb.ApiResponsePaginationSaldoDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	s.logger.Debug("Fetching trashed saldo records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

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

	res, totalRecords, err := s.saldoQueryService.FindByTrashed(&reqService)

	if err != nil {
		s.logger.Debug("FindByTrashed failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapping.ToProtoResponsePaginationSaldoDeleteAt(paginationMeta, "success", "Successfully fetched saldo record", res)

	return so, nil
}

func (s *saldoHandleGrpc) CreateSaldo(ctx context.Context, req *pb.CreateSaldoRequest) (*pb.ApiResponseSaldo, error) {
	request := requests.CreateSaldoRequest{
		CardNumber:   req.GetCardNumber(),
		TotalBalance: int(req.GetTotalBalance()),
	}

	s.logger.Debug("Creating saldo record", zap.Any("request", request))

	if err := request.Validate(); err != nil {
		s.logger.Debug("Create failed", zap.Any("error", err))
		return nil, saldo_errors.ErrGrpcValidateCreateSaldo
	}

	saldo, err := s.saldoCommandService.CreateSaldo(&request)

	if err != nil {
		s.logger.Debug("Create failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseSaldo("success", "Successfully created saldo record", saldo)

	return so, nil

}

func (s *saldoHandleGrpc) UpdateSaldo(ctx context.Context, req *pb.UpdateSaldoRequest) (*pb.ApiResponseSaldo, error) {
	id := int(req.GetSaldoId())

	s.logger.Debug("Updating saldo record", zap.Int("id", id))

	if id == 0 {
		return nil, saldo_errors.ErrGrpcSaldoInvalidID
	}

	request := requests.UpdateSaldoRequest{
		SaldoID:      &id,
		CardNumber:   req.GetCardNumber(),
		TotalBalance: int(req.GetTotalBalance()),
	}

	if err := request.Validate(); err != nil {
		s.logger.Debug("Update failed", zap.Any("error", err))
		return nil, saldo_errors.ErrGrpcValidateUpdateSaldo
	}

	saldo, err := s.saldoCommandService.UpdateSaldo(&request)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseSaldo("success", "Successfully updated saldo record", saldo)

	return so, nil
}

func (s *saldoHandleGrpc) TrashedSaldo(ctx context.Context, req *pb.FindByIdSaldoRequest) (*pb.ApiResponseSaldo, error) {
	id := int(req.GetSaldoId())

	s.logger.Debug("Trashing saldo record", zap.Int("id", id))

	if id == 0 {
		s.logger.Debug("TrashedSaldo failed", zap.Any("error", saldo_errors.ErrGrpcSaldoInvalidID))
		return nil, saldo_errors.ErrGrpcSaldoInvalidID
	}

	saldo, err := s.saldoCommandService.TrashSaldo(id)

	if err != nil {
		s.logger.Debug("TrashedSaldo failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseSaldo("success", "Successfully trashed saldo record", saldo)

	return so, nil
}

func (s *saldoHandleGrpc) RestoreSaldo(ctx context.Context, req *pb.FindByIdSaldoRequest) (*pb.ApiResponseSaldo, error) {
	id := int(req.GetSaldoId())

	s.logger.Debug("Restoring saldo record", zap.Int("id", id))

	if id == 0 {
		s.logger.Debug("RestoreSaldo failed", zap.Any("error", saldo_errors.ErrGrpcSaldoInvalidID))
		return nil, saldo_errors.ErrGrpcSaldoInvalidID
	}

	saldo, err := s.saldoCommandService.RestoreSaldo(id)

	if err != nil {
		s.logger.Debug("RestoreSaldo failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseSaldo("success", "Successfully restored saldo record", saldo)

	return so, nil
}

func (s *saldoHandleGrpc) DeleteSaldo(ctx context.Context, req *pb.FindByIdSaldoRequest) (*pb.ApiResponseSaldoDelete, error) {
	id := int(req.GetSaldoId())

	s.logger.Debug("Deleting saldo record", zap.Int("id", id))

	if id == 0 {
		s.logger.Debug("DeleteSaldo failed", zap.Any("error", saldo_errors.ErrGrpcSaldoInvalidID))
		return nil, saldo_errors.ErrGrpcSaldoInvalidID
	}

	_, err := s.saldoCommandService.DeleteSaldoPermanent(id)

	if err != nil {
		s.logger.Debug("DeleteSaldo failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseSaldoDelete("success", "Successfully deleted saldo record")

	return so, nil
}

func (s *saldoHandleGrpc) RestoreAllSaldo(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseSaldoAll, error) {
	s.logger.Debug("Restoring all saldo record")

	_, err := s.saldoCommandService.RestoreAllSaldo()

	if err != nil {
		s.logger.Debug("RestoreAllSaldo failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseSaldoAll("success", "Successfully restore all saldo")

	return so, nil
}

func (s *saldoHandleGrpc) DeleteAllSaldoPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseSaldoAll, error) {
	s.logger.Debug("Deleting all saldo record")

	_, err := s.saldoCommandService.DeleteAllSaldoPermanent()

	if err != nil {
		s.logger.Debug("DeleteAllSaldoPermanent failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseSaldoAll("success", "delete saldo permanent")

	return so, nil
}
