package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto"
	"github.com/MamangRust/monolith-payment-gateway-shared/pb"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/service"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type transferHandleGrpc struct {
	pb.UnimplementedTransferServiceServer
	transferQueryService            service.TransferQueryService
	transferStatisticsService       service.TransferStatisticsService
	transferStatisticsByCardService service.TransferStatisticByCardService
	transferCommandService          service.TransferCommandService
	logger                          logger.LoggerInterface
	mapping                         protomapper.TransferProtoMapper
}

func NewTransferHandleGrpc(service service.Service, logger logger.LoggerInterface) *transferHandleGrpc {
	return &transferHandleGrpc{
		transferQueryService:            service.TransferQuery,
		transferStatisticsService:       service.TransferStatistic,
		transferStatisticsByCardService: service.TransferStatisticByCard,
		transferCommandService:          service.TransferCommand,
		logger:                          logger,
		mapping:                         protomapper.NewTransferProtoMapper(),
	}
}

func (s *transferHandleGrpc) FindAllTransfer(ctx context.Context, request *pb.FindAllTransferRequest) (*pb.ApiResponsePaginationTransfer, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	s.logger.Debug("Fetching transfer", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllTranfers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchants, totalRecords, err := s.transferQueryService.FindAll(&reqService)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapping.ToProtoResponsePaginationTransfer(paginationMeta, "success", "Successfully fetch transfer records", merchants)

	return so, nil
}

func (s *transferHandleGrpc) FindByIdTransfer(ctx context.Context, request *pb.FindByIdTransferRequest) (*pb.ApiResponseTransfer, error) {
	id := int(request.GetTransferId())

	s.logger.Debug("Fetching transfer", zap.Int("id", id))

	if id == 0 {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", transfer_errors.ErrGrpcTransferInvalidID))
		return nil, transfer_errors.ErrGrpcTransferInvalidID
	}

	transfer, err := s.transferQueryService.FindById(id)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransfer("success", "Successfully fetch transfer record", transfer)

	return so, nil
}

func (s *transferHandleGrpc) FindMonthlyTransferStatusSuccess(ctx context.Context, req *pb.FindMonthlyTransferStatus) (*pb.ApiResponseTransferMonthStatusSuccess, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	s.logger.Debug("Fetching transfer", zap.Int("year", year), zap.Int("month", month))

	if year <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("year", year))
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("month", month))
		return nil, transfer_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthStatusTransfer{
		Year:  year,
		Month: month,
	}

	records, err := s.transferStatisticsService.FindMonthTransferStatusSuccess(&reqService)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferMonthStatusSuccess("success", "Successfully fetched monthly Transfer status success", records)

	return so, nil
}

func (s *transferHandleGrpc) FindYearlyTransferStatusSuccess(ctx context.Context, req *pb.FindYearTransferStatus) (*pb.ApiResponseTransferYearStatusSuccess, error) {
	year := int(req.GetYear())

	s.logger.Debug("Fetching transfer", zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("year", year))
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	records, err := s.transferStatisticsService.FindYearlyTransferStatusSuccess(year)
	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferYearStatusSuccess("success", "Successfully fetched yearly Transfer status success", records)

	return so, nil
}

func (s *transferHandleGrpc) FindMonthlyTransferStatusFailed(ctx context.Context, req *pb.FindMonthlyTransferStatus) (*pb.ApiResponseTransferMonthStatusFailed, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	s.logger.Debug("Fetching transfer", zap.Int("year", year), zap.Int("month", month))

	if year <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("year", year))
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("month", month))
		return nil, transfer_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthStatusTransfer{
		Year:  year,
		Month: month,
	}

	records, err := s.transferStatisticsService.FindMonthTransferStatusFailed(&reqService)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferMonthStatusFailed("success", "success fetched monthly Transfer status Failed", records)

	return so, nil
}

func (s *transferHandleGrpc) FindYearlyTransferStatusFailed(ctx context.Context, req *pb.FindYearTransferStatus) (*pb.ApiResponseTransferYearStatusFailed, error) {
	year := int(req.GetYear())

	s.logger.Debug("Fetching transfer", zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("year", year))
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	records, err := s.transferStatisticsService.FindYearlyTransferStatusFailed(year)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferYearStatusFailed("success", "success fetched yearly Transfer status Failed", records)

	return so, nil
}

func (s *transferHandleGrpc) FindMonthlyTransferStatusSuccessByCardNumber(ctx context.Context, req *pb.FindMonthlyTransferStatusCardNumber) (*pb.ApiResponseTransferMonthStatusSuccess, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	cardNumber := req.GetCardNumber()

	s.logger.Debug("Fetching transfer", zap.Int("year", year), zap.Int("month", month))

	if year <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("year", year))
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("month", month))
		return nil, transfer_errors.ErrGrpcInvalidMonth
	}

	if cardNumber == "" {
		s.logger.Error("Failed to fetch transfer", zap.String("card_number", cardNumber))
		return nil, transfer_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthStatusTransferCardNumber{
		Year:       year,
		Month:      month,
		CardNumber: cardNumber,
	}

	records, err := s.transferStatisticsByCardService.FindMonthTransferStatusSuccessByCardNumber(&reqService)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferMonthStatusSuccess("success", "Successfully fetched monthly Transfer status success", records)

	return so, nil
}

func (s *transferHandleGrpc) FindYearlyTransferStatusSuccessByCardNumber(ctx context.Context, req *pb.FindYearTransferStatusCardNumber) (*pb.ApiResponseTransferYearStatusSuccess, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	s.logger.Debug("Fetching transfer", zap.Int("year", year), zap.String("card_number", cardNumber))

	if year <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("year", year))
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		s.logger.Error("Failed to fetch transfer", zap.String("card_number", cardNumber))
		return nil, transfer_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.YearStatusTransferCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	records, err := s.transferStatisticsByCardService.FindYearlyTransferStatusSuccessByCardNumber(&reqService)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferYearStatusSuccess("success", "Successfully fetched yearly Transfer status success", records)

	return so, nil
}

func (s *transferHandleGrpc) FindMonthlyTransferStatusFailedByCardNumber(ctx context.Context, req *pb.FindMonthlyTransferStatusCardNumber) (*pb.ApiResponseTransferMonthStatusFailed, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	cardNumber := req.GetCardNumber()

	s.logger.Debug("Fetching transfer", zap.Int("year", year), zap.Int("month", month))

	if year <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("year", year))
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("month", month))
		return nil, transfer_errors.ErrGrpcInvalidMonth
	}

	if cardNumber == "" {
		s.logger.Error("Failed to fetch transfer", zap.String("card_number", cardNumber))
		return nil, transfer_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthStatusTransferCardNumber{
		Year:       year,
		Month:      month,
		CardNumber: cardNumber,
	}

	records, err := s.transferStatisticsByCardService.FindMonthTransferStatusFailedByCardNumber(&reqService)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferMonthStatusFailed("success", "success fetched monthly Transfer status Failed", records)

	return so, nil
}

func (s *transferHandleGrpc) FindYearlyTransferStatusFailedByCardNumber(ctx context.Context, req *pb.FindYearTransferStatusCardNumber) (*pb.ApiResponseTransferYearStatusFailed, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	s.logger.Debug("Fetching transfer", zap.Int("year", year), zap.String("card_number", cardNumber))

	if year <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("year", year))
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		s.logger.Error("Failed to fetch transfer", zap.String("card_number", cardNumber))
		return nil, transfer_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.YearStatusTransferCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	records, err := s.transferStatisticsByCardService.FindYearlyTransferStatusFailedByCardNumber(&reqService)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferYearStatusFailed("success", "success fetched yearly Transfer status Failed", records)

	return so, nil
}

func (s *transferHandleGrpc) FindMonthlyTransferAmounts(ctx context.Context, req *pb.FindYearTransferStatus) (*pb.ApiResponseTransferMonthAmount, error) {
	year := int(req.GetYear())

	s.logger.Debug("Fetching transfer", zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("year", year))
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	amounts, err := s.transferStatisticsService.FindMonthlyTransferAmounts(year)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferMonthAmount("success", "Successfully fetched monthly transfer amounts", amounts)

	return so, nil
}

func (s *transferHandleGrpc) FindYearlyTransferAmounts(ctx context.Context, req *pb.FindYearTransferStatus) (*pb.ApiResponseTransferYearAmount, error) {
	year := int(req.GetYear())

	s.logger.Debug("Fetching transfer", zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("year", year))
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	amounts, err := s.transferStatisticsService.FindYearlyTransferAmounts(year)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferYearAmount("success", "Successfully fetched yearly transfer amounts", amounts)

	return so, nil
}

func (s *transferHandleGrpc) FindMonthlyTransferAmountsBySenderCardNumber(ctx context.Context, req *pb.FindByCardNumberTransferRequest) (*pb.ApiResponseTransferMonthAmount, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	s.logger.Debug("Fetching transfer", zap.Int("year", year), zap.String("card_number", cardNumber))

	if year <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("year", year))
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		s.logger.Error("Failed to fetch transfer", zap.String("card_number", cardNumber))
		return nil, transfer_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	amounts, err := s.transferStatisticsByCardService.FindMonthlyTransferAmountsBySenderCardNumber(&reqService)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, transfer_errors.ErrGrpcFailedFindMonthlyTransferAmountsBySenderCardNumber
	}

	so := s.mapping.ToProtoResponseTransferMonthAmount("success", "Successfully fetched monthly transfer amounts by sender card number", amounts)

	return so, nil
}

func (s *transferHandleGrpc) FindMonthlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *pb.FindByCardNumberTransferRequest) (*pb.ApiResponseTransferMonthAmount, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	s.logger.Debug("Fetching transfer", zap.Int("year", year), zap.String("card_number", cardNumber))

	if year <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("year", year))
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		s.logger.Error("Failed to fetch transfer", zap.String("card_number", cardNumber))
		return nil, transfer_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	amounts, err := s.transferStatisticsByCardService.FindMonthlyTransferAmountsByReceiverCardNumber(&reqService)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferMonthAmount("success", "Successfully fetched monthly transfer amounts by receiver card number", amounts)

	return so, nil
}

func (s *transferHandleGrpc) FindYearlyTransferAmountsBySenderCardNumber(ctx context.Context, req *pb.FindByCardNumberTransferRequest) (*pb.ApiResponseTransferYearAmount, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	s.logger.Debug("Fetching transfer", zap.Int("year", year), zap.String("card_number", cardNumber))

	if year <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("year", year))
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		s.logger.Error("Failed to fetch transfer", zap.String("card_number", cardNumber))
		return nil, transfer_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	amounts, err := s.transferStatisticsByCardService.FindYearlyTransferAmountsBySenderCardNumber(&reqService)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferYearAmount("success", "Successfully fetched yearly transfer amounts by sender card number", amounts)

	return so, nil
}

func (s *transferHandleGrpc) FindYearlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *pb.FindByCardNumberTransferRequest) (*pb.ApiResponseTransferYearAmount, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	s.logger.Debug("Fetching transfer", zap.Int("year", year), zap.String("card_number", cardNumber))

	if year <= 0 {
		s.logger.Error("Failed to fetch transfer", zap.Int("year", year))
		return nil, transfer_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		s.logger.Error("Failed to fetch transfer", zap.String("card_number", cardNumber))
		return nil, transfer_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	amounts, err := s.transferStatisticsByCardService.FindYearlyTransferAmountsByReceiverCardNumber(&reqService)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferYearAmount("success", "Successfully fetched yearly transfer amounts by receiver card number", amounts)

	return so, nil
}

func (s *transferHandleGrpc) FindByTransferByTransferFrom(ctx context.Context, request *pb.FindTransferByTransferFromRequest) (*pb.ApiResponseTransfers, error) {
	transfer_from := request.GetTransferFrom()

	s.logger.Debug("Fetching transfer", zap.String("transfer_from", transfer_from))

	if transfer_from == "" {
		s.logger.Error("Failed to fetch transfer", zap.String("transfer_from", transfer_from))
		return nil, transfer_errors.ErrGrpcInvalidCardNumber
	}

	merchants, err := s.transferQueryService.FindTransferByTransferFrom(transfer_from)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransfers("success", "Successfully fetch transfer records", merchants)

	return so, nil
}

func (s *transferHandleGrpc) FindByTransferByTransferTo(ctx context.Context, request *pb.FindTransferByTransferToRequest) (*pb.ApiResponseTransfers, error) {
	transfer_to := request.GetTransferTo()

	s.logger.Debug("Fetching transfer", zap.String("transfer_to", transfer_to))

	if transfer_to == "" {
		s.logger.Error("Failed to fetch transfer", zap.String("transfer_to", transfer_to))
		return nil, transfer_errors.ErrGrpcInvalidCardNumber
	}

	merchants, err := s.transferQueryService.FindTransferByTransferTo(transfer_to)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransfers("success", "Successfully fetch transfer records", merchants)

	return so, nil
}

func (s *transferHandleGrpc) FindByActiveTransfer(ctx context.Context, req *pb.FindAllTransferRequest) (*pb.ApiResponsePaginationTransferDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	s.logger.Debug("Fetching transfer",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllTranfers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.transferQueryService.FindByActive(&reqService)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapping.ToProtoResponsePaginationTransferDeleteAt(paginationMeta, "success", "Successfully fetch transfer records", res)

	return so, nil
}

func (s *transferHandleGrpc) FindByTrashedTransfer(ctx context.Context, req *pb.FindAllTransferRequest) (*pb.ApiResponsePaginationTransferDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	s.logger.Debug("Fetching transfer",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllTranfers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.transferQueryService.FindByTrashed(&reqService)

	if err != nil {
		s.logger.Error("Failed to fetch transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}
	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapping.ToProtoResponsePaginationTransferDeleteAt(paginationMeta, "success", "Successfully fetch transfer records", res)

	return so, nil
}

func (s *transferHandleGrpc) CreateTransfer(ctx context.Context, request *pb.CreateTransferRequest) (*pb.ApiResponseTransfer, error) {
	req := requests.CreateTransferRequest{
		TransferFrom:   request.GetTransferFrom(),
		TransferTo:     request.GetTransferTo(),
		TransferAmount: int(request.GetTransferAmount()),
	}

	s.logger.Debug("Starting create transfer process",
		zap.Any("request", req),
	)

	if err := req.Validate(); err != nil {
		s.logger.Error("Failed to create transfer", zap.Any("error", err))
		return nil, transfer_errors.ErrGrpcValidateCreateTransferRequest
	}

	res, err := s.transferCommandService.CreateTransaction(&req)

	if err != nil {
		s.logger.Error("Failed to create transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransfer("success", "Successfully created transfer", res)

	return so, nil
}

func (s *transferHandleGrpc) UpdateTransfer(ctx context.Context, request *pb.UpdateTransferRequest) (*pb.ApiResponseTransfer, error) {
	id := int(request.GetTransferId())

	s.logger.Debug("Starting update transfer process",
		zap.Any("request", id),
	)

	if id == 0 {
		s.logger.Error("Failed to update transfer", zap.Any("error", transfer_errors.ErrGrpcTransferInvalidID))
		return nil, transfer_errors.ErrGrpcTransferInvalidID
	}

	req := requests.UpdateTransferRequest{
		TransferID:     &id,
		TransferFrom:   request.GetTransferFrom(),
		TransferTo:     request.GetTransferTo(),
		TransferAmount: int(request.GetTransferAmount()),
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Failed to update transfer", zap.Any("error", err))
		return nil, transfer_errors.ErrGrpcValidateUpdateTransferRequest
	}

	res, err := s.transferCommandService.UpdateTransaction(&req)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransfer("success", "Successfully updated transfer", res)

	return so, nil
}

func (s *transferHandleGrpc) TrashedTransfer(ctx context.Context, request *pb.FindByIdTransferRequest) (*pb.ApiResponseTransfer, error) {
	id := int(request.GetTransferId())

	s.logger.Debug("Starting trashed transfer process",
		zap.Any("request", id),
	)

	if id == 0 {
		s.logger.Error("Failed to trashed transfer", zap.Any("error", transfer_errors.ErrGrpcTransferInvalidID))
		return nil, transfer_errors.ErrGrpcTransferInvalidID
	}

	res, err := s.transferCommandService.TrashedTransfer(id)

	if err != nil {
		s.logger.Error("Failed to trashed transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransfer("success", "Successfully trashed transfer", res)

	return so, nil
}

func (s *transferHandleGrpc) RestoreTransfer(ctx context.Context, request *pb.FindByIdTransferRequest) (*pb.ApiResponseTransfer, error) {
	id := int(request.GetTransferId())

	s.logger.Debug("Starting restore transfer process",
		zap.Any("request", id),
	)

	if id == 0 {
		s.logger.Error("Failed to restore transfer", zap.Any("error", transfer_errors.ErrGrpcTransferInvalidID))
		return nil, transfer_errors.ErrGrpcTransferInvalidID
	}

	res, err := s.transferCommandService.RestoreTransfer(id)

	if err != nil {
		s.logger.Error("Failed to restore transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransfer("success", "Successfully restored transfer", res)

	return so, nil
}

func (s *transferHandleGrpc) DeleteTransferPermanent(ctx context.Context, request *pb.FindByIdTransferRequest) (*pb.ApiResponseTransferDelete, error) {
	id := int(request.GetTransferId())

	s.logger.Debug("Starting delete transfer process",
		zap.Any("request", id),
	)

	if id == 0 {
		s.logger.Error("Failed to delete transfer", zap.Any("error", transfer_errors.ErrGrpcTransferInvalidID))
		return nil, transfer_errors.ErrGrpcTransferInvalidID
	}

	_, err := s.transferCommandService.DeleteTransferPermanent(id)

	if err != nil {
		s.logger.Error("Failed to delete transfer", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferDelete("success", "Successfully restored transfer")

	return so, nil
}

func (s *transferHandleGrpc) RestoreAllTransfer(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTransferAll, error) {
	s.logger.Debug("Starting restore all transfer process")

	_, err := s.transferCommandService.RestoreAllTransfer()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferAll("success", "Successfully restored transfer")

	return so, nil
}

func (s *transferHandleGrpc) DeleteAllTransferPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTransferAll, error) {
	s.logger.Debug("Starting delete all transfer process")

	_, err := s.transferCommandService.DeleteAllTransferPermanent()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransferAll("success", "delete transfer permanent")

	return so, nil
}
