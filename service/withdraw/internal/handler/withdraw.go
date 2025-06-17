package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto"
	"github.com/MamangRust/monolith-payment-gateway-shared/pb"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/service"
	"go.uber.org/zap"

	"google.golang.org/protobuf/types/known/emptypb"
)

type withdrawHandleGrpc struct {
	pb.UnimplementedWithdrawServiceServer
	withdrawQuery       service.WithdrawQueryService
	withdrawCommand     service.WithdrawCommandService
	withdrawStats       service.WithdrawStatisticService
	withdrawStatsByCard service.WithdrawStatisticByCardService
	logger              logger.LoggerInterface
	mapping             protomapper.WithdrawalProtoMapper
}

func NewWithdrawHandleGrpc(service *service.Service, logger logger.LoggerInterface) *withdrawHandleGrpc {
	return &withdrawHandleGrpc{
		withdrawQuery:       service.WithdrawQuery,
		withdrawCommand:     service.WithdrawCommand,
		withdrawStats:       service.WithdrawStats,
		withdrawStatsByCard: service.WithdrawStatsByCard,
		logger:              logger,
		mapping:             protomapper.NewWithdrawProtoMapper(),
	}
}

func (w *withdrawHandleGrpc) FindAllWithdraw(ctx context.Context, req *pb.FindAllWithdrawRequest) (*pb.ApiResponsePaginationWithdraw, error) {
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

	reqService := requests.FindAllWithdraws{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	withdraws, totalRecords, err := w.withdrawQuery.FindAll(&reqService)

	if err != nil {
		w.logger.Error("FindAllWithdraw", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := w.mapping.ToProtoResponsePaginationWithdraw(paginationMeta, "success", "withdraw", withdraws)

	return so, nil
}

func (w *withdrawHandleGrpc) FindAllWithdrawByCardNumber(ctx context.Context, req *pb.FindAllWithdrawByCardNumberRequest) (*pb.ApiResponsePaginationWithdraw, error) {
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

	reqService := requests.FindAllWithdrawCardNumber{
		CardNumber: card_number,
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
	}

	withdraws, totalRecords, err := w.withdrawQuery.FindAllByCardNumber(&reqService)

	if err != nil {
		w.logger.Error("FindAllWithdrawByCardNumber", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := w.mapping.ToProtoResponsePaginationWithdraw(paginationMeta, "success", "Withdraws fetched successfully", withdraws)

	return so, nil
}

func (w *withdrawHandleGrpc) FindByIdWithdraw(ctx context.Context, req *pb.FindByIdWithdrawRequest) (*pb.ApiResponseWithdraw, error) {
	id := int(req.GetWithdrawId())

	w.logger.Debug("FindByIdWithdraw", zap.Int("id", id))

	if id == 0 {
		w.logger.Error("FindByIdWithdraw", zap.Any("error", withdraw_errors.ErrGrpcWithdrawInvalidID))
		return nil, withdraw_errors.ErrGrpcWithdrawInvalidID
	}

	withdraw, err := w.withdrawQuery.FindById(id)

	if err != nil {
		w.logger.Error("FindByIdWithdraw", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := w.mapping.ToProtoResponseWithdraw("success", "Successfully fetched withdraw", withdraw)

	return so, nil
}

func (s *withdrawHandleGrpc) FindMonthlyWithdrawStatusSuccess(ctx context.Context, req *pb.FindMonthlyWithdrawStatus) (*pb.ApiResponseWithdrawMonthStatusSuccess, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	s.logger.Debug("FindMonthlyWithdrawStatusSuccess", zap.Int("year", year), zap.Int("month", month))

	if year <= 0 {
		s.logger.Error("FindMonthlyWithdrawStatusSuccess", zap.Any("error", withdraw_errors.ErrGrpcInvalidYear))
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		s.logger.Error("FindMonthlyWithdrawStatusSuccess", zap.Any("error", withdraw_errors.ErrGrpcInvalidMonth))
		return nil, withdraw_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthStatusWithdraw{
		Year:  year,
		Month: month,
	}

	records, err := s.withdrawStats.FindMonthWithdrawStatusSuccess(&reqService)

	if err != nil {
		s.logger.Error("FindMonthlyWithdrawStatusSuccess", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseWithdrawMonthStatusSuccess("success", "Successfully fetched withdraw", records)

	return so, nil
}

func (s *withdrawHandleGrpc) FindYearlyWithdrawStatusSuccess(ctx context.Context, req *pb.FindYearWithdrawStatus) (*pb.ApiResponseWithdrawYearStatusSuccess, error) {
	year := int(req.GetYear())

	s.logger.Debug("FindYearlyWithdrawStatusSuccess", zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("FindYearlyWithdrawStatusSuccess", zap.Any("error", withdraw_errors.ErrGrpcInvalidYear))
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	records, err := s.withdrawStats.FindYearlyWithdrawStatusSuccess(year)

	if err != nil {
		s.logger.Error("FindYearlyWithdrawStatusSuccess", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseWithdrawYearStatusSuccess("success", "Successfully fetched yearly Withdraw status success", records)

	return so, nil
}

func (s *withdrawHandleGrpc) FindMonthlyWithdrawStatusFailed(ctx context.Context, req *pb.FindMonthlyWithdrawStatus) (*pb.ApiResponseWithdrawMonthStatusFailed, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	s.logger.Debug("FindMonthlyWithdrawStatusFailed", zap.Int("year", year), zap.Int("month", month))

	if year <= 0 {
		s.logger.Error("FindMonthlyWithdrawStatusFailed", zap.Any("error", withdraw_errors.ErrGrpcInvalidYear))
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		s.logger.Error("FindMonthlyWithdrawStatusFailed", zap.Any("error", withdraw_errors.ErrGrpcInvalidMonth))
		return nil, withdraw_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthStatusWithdraw{
		Year:  year,
		Month: month,
	}

	records, err := s.withdrawStats.FindMonthWithdrawStatusFailed(&reqService)

	if err != nil {
		s.logger.Error("FindMonthlyWithdrawStatusFailed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseWithdrawMonthStatusFailed("success", "success fetched monthly Withdraw status Failed", records)

	return so, nil
}

func (s *withdrawHandleGrpc) FindYearlyWithdrawStatusFailed(ctx context.Context, req *pb.FindYearWithdrawStatus) (*pb.ApiResponseWithdrawYearStatusFailed, error) {
	year := int(req.GetYear())

	s.logger.Debug("FindYearlyWithdrawStatusFailed", zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("FindYearlyWithdrawStatusFailed", zap.Any("error", withdraw_errors.ErrGrpcInvalidYear))
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	records, err := s.withdrawStats.FindYearlyWithdrawStatusFailed(year)

	if err != nil {
		s.logger.Error("FindYearlyWithdrawStatusFailed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseWithdrawYearStatusFailed("success", "success fetched yearly Withdraw status Failed", records)

	return so, nil
}

func (s *withdrawHandleGrpc) FindMonthlyWithdrawStatusSuccessCardNumber(ctx context.Context, req *pb.FindMonthlyWithdrawStatusCardNumber) (*pb.ApiResponseWithdrawMonthStatusSuccess, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	cardNumber := req.GetCardNumber()

	s.logger.Debug("FindMonthlyWithdrawStatusSuccessCardNumber", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", cardNumber))

	if year <= 0 {
		s.logger.Error("FindMonthlyWithdrawStatusSuccessCardNumber", zap.Any("error", withdraw_errors.ErrGrpcInvalidYear))
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		s.logger.Error("FindMonthlyWithdrawStatusSuccessCardNumber", zap.Any("error", withdraw_errors.ErrGrpcInvalidMonth))
		return nil, withdraw_errors.ErrGrpcInvalidMonth
	}

	if cardNumber == "" {
		s.logger.Error("FindMonthlyWithdrawStatusSuccessCardNumber", zap.Any("error", withdraw_errors.ErrGrpcInvalidCardNumber))
		return nil, withdraw_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthStatusWithdrawCardNumber{
		Year:       year,
		Month:      month,
		CardNumber: cardNumber,
	}

	records, err := s.withdrawStatsByCard.FindMonthWithdrawStatusSuccessByCardNumber(&reqService)

	if err != nil {
		s.logger.Error("FindMonthlyWithdrawStatusSuccessCardNumber", zap.Any("error", err))

		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseWithdrawMonthStatusSuccess("success", "Successfully fetched withdraw", records)

	return so, nil
}

func (s *withdrawHandleGrpc) FindYearlyWithdrawStatusSuccessCardNumber(ctx context.Context, req *pb.FindYearWithdrawStatusCardNumber) (*pb.ApiResponseWithdrawYearStatusSuccess, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	s.logger.Debug("FindYearlyWithdrawStatusSuccessCardNumber", zap.Int("year", year), zap.String("card_number", cardNumber))

	if year <= 0 {
		s.logger.Error("FindYearlyWithdrawStatusSuccessCardNumber", zap.Any("error", withdraw_errors.ErrGrpcInvalidYear))
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		s.logger.Error("FindYearlyWithdrawStatusSuccessCardNumber", zap.Any("error", withdraw_errors.ErrGrpcInvalidCardNumber))
		return nil, withdraw_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.YearStatusWithdrawCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	records, err := s.withdrawStatsByCard.FindYearlyWithdrawStatusSuccessByCardNumber(&reqService)

	if err != nil {
		s.logger.Error("FindYearlyWithdrawStatusSuccessCardNumber", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseWithdrawYearStatusSuccess("success", "Successfully fetched yearly Withdraw status success", records)

	return so, nil
}

func (s *withdrawHandleGrpc) FindMonthlyWithdrawStatusFailedCardNumber(ctx context.Context, req *pb.FindMonthlyWithdrawStatusCardNumber) (*pb.ApiResponseWithdrawMonthStatusFailed, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	cardNumber := req.GetCardNumber()

	s.logger.Debug("FindMonthlyWithdrawStatusFailedCardNumber", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", cardNumber))

	if year <= 0 {
		s.logger.Error("FindMonthlyWithdrawStatusFailedCardNumber", zap.Any("error", withdraw_errors.ErrGrpcInvalidYear))
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		s.logger.Error("FindMonthlyWithdrawStatusFailedCardNumber", zap.Any("error", withdraw_errors.ErrGrpcInvalidMonth))
		return nil, withdraw_errors.ErrGrpcInvalidMonth
	}

	if cardNumber == "" {
		s.logger.Error("FindMonthlyWithdrawStatusFailedCardNumber", zap.Any("error", withdraw_errors.ErrGrpcInvalidCardNumber))
		return nil, withdraw_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthStatusWithdrawCardNumber{
		Year:       year,
		Month:      month,
		CardNumber: cardNumber,
	}

	records, err := s.withdrawStatsByCard.FindMonthWithdrawStatusFailedByCardNumber(&reqService)

	if err != nil {
		s.logger.Error("FindMonthlyWithdrawStatusFailedCardNumber", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseWithdrawMonthStatusFailed("success", "Successfully fetched monthly Withdraw status failed", records)

	return so, nil
}

func (s *withdrawHandleGrpc) FindYearlyWithdrawStatusFailedCardNumber(ctx context.Context, req *pb.FindYearWithdrawStatusCardNumber) (*pb.ApiResponseWithdrawYearStatusFailed, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	s.logger.Debug("FindYearlyWithdrawStatusFailedCardNumber", zap.Int("year", year), zap.String("card_number", cardNumber))

	if year <= 0 {
		s.logger.Error("FindYearlyWithdrawStatusFailedCardNumber", zap.Any("error", withdraw_errors.ErrGrpcInvalidYear))
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	reqService := requests.YearStatusWithdrawCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	records, err := s.withdrawStatsByCard.FindYearlyWithdrawStatusFailedByCardNumber(&reqService)
	if err != nil {
		s.logger.Error("FindYearlyWithdrawStatusFailedCardNumber", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseWithdrawYearStatusFailed("success", "Successfully fetched yearly Withdraw status failed", records)

	return so, nil
}

func (w *withdrawHandleGrpc) FindMonthlyWithdraws(ctx context.Context, req *pb.FindYearWithdrawStatus) (*pb.ApiResponseWithdrawMonthAmount, error) {
	year := int(req.GetYear())

	w.logger.Debug("FindMonthlyWithdraws", zap.Int("year", year))

	if year <= 0 {
		w.logger.Error("FindMonthlyWithdraws", zap.Any("error", withdraw_errors.ErrGrpcInvalidYear))
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	withdraws, err := w.withdrawStats.FindMonthlyWithdraws(year)

	if err != nil {
		w.logger.Error("FindMonthlyWithdraws", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := w.mapping.ToProtoResponseWithdrawMonthAmount("success", "Successfully fetched monthly withdraws", withdraws)

	return so, nil
}

func (w *withdrawHandleGrpc) FindYearlyWithdraws(ctx context.Context, req *pb.FindYearWithdrawStatus) (*pb.ApiResponseWithdrawYearAmount, error) {
	year := int(req.GetYear())

	w.logger.Debug("FindYearlyWithdraws", zap.Int("year", year))

	if year <= 0 {
		w.logger.Error("FindYearlyWithdraws", zap.Any("error", withdraw_errors.ErrGrpcInvalidYear))
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	withdraws, err := w.withdrawStats.FindYearlyWithdraws(year)

	if err != nil {
		w.logger.Error("FindYearlyWithdraws", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := w.mapping.ToProtoResponseWithdrawYearAmount("success", "Successfully fetched yearly withdraws", withdraws)

	return so, nil
}

func (w *withdrawHandleGrpc) FindMonthlyWithdrawsByCardNumber(ctx context.Context, req *pb.FindYearWithdrawCardNumber) (*pb.ApiResponseWithdrawMonthAmount, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	w.logger.Debug("FindMonthlyWithdrawsByCardNumber", zap.Int("year", year), zap.String("card_number", cardNumber))

	if year <= 0 {
		w.logger.Error("FindMonthlyWithdrawsByCardNumber", zap.Any("error", withdraw_errors.ErrGrpcInvalidYear))
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		w.logger.Error("FindMonthlyWithdrawsByCardNumber", zap.Any("error", withdraw_errors.ErrGrpcInvalidCardNumber))
		return nil, withdraw_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.YearMonthCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	withdraws, err := w.withdrawStatsByCard.FindMonthlyWithdrawsByCardNumber(&reqService)

	if err != nil {
		w.logger.Error("FindMonthlyWithdrawsByCardNumber", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := w.mapping.ToProtoResponseWithdrawMonthAmount("success", "Successfully fetched monthly withdraws by card number", withdraws)

	return so, nil
}

func (w *withdrawHandleGrpc) FindYearlyWithdrawsByCardNumber(ctx context.Context, req *pb.FindYearWithdrawCardNumber) (*pb.ApiResponseWithdrawYearAmount, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	w.logger.Debug("FindYearlyWithdrawsByCardNumber", zap.Int("year", year), zap.String("card_number", cardNumber))

	if year <= 0 {
		w.logger.Error("FindYearlyWithdrawsByCardNumber", zap.Any("error", withdraw_errors.ErrGrpcInvalidYear))
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		w.logger.Error("FindYearlyWithdrawsByCardNumber", zap.Any("error", withdraw_errors.ErrGrpcInvalidCardNumber))
		return nil, withdraw_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.YearMonthCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	withdraws, err := w.withdrawStatsByCard.FindYearlyWithdrawsByCardNumber(&reqService)

	if err != nil {
		w.logger.Error("FindYearlyWithdrawsByCardNumber", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := w.mapping.ToProtoResponseWithdrawYearAmount("success", "Successfully fetched yearly withdraws by card number", withdraws)

	return so, nil
}

func (w *withdrawHandleGrpc) FindByActive(ctx context.Context, req *pb.FindAllWithdrawRequest) (*pb.ApiResponsePaginationWithdrawDeleteAt, error) {
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

	reqService := requests.FindAllWithdraws{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := w.withdrawQuery.FindByActive(&reqService)

	if err != nil {
		w.logger.Error("FindByActive", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := w.mapping.ToProtoResponsePaginationWithdrawDeleteAt(paginationMeta, "success", "Successfully fetched withdraws", res)

	return so, nil
}

func (w *withdrawHandleGrpc) FindByTrashed(ctx context.Context, req *pb.FindAllWithdrawRequest) (*pb.ApiResponsePaginationWithdrawDeleteAt, error) {
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

	reqService := requests.FindAllWithdraws{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := w.withdrawQuery.FindByTrashed(&reqService)

	if err != nil {
		w.logger.Error("FindByTrashed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := w.mapping.ToProtoResponsePaginationWithdrawDeleteAt(paginationMeta, "success", "Successfully fetched withdraws", res)

	return so, nil
}

func (w *withdrawHandleGrpc) CreateWithdraw(ctx context.Context, req *pb.CreateWithdrawRequest) (*pb.ApiResponseWithdraw, error) {
	request := &requests.CreateWithdrawRequest{
		CardNumber:     req.CardNumber,
		WithdrawAmount: int(req.WithdrawAmount),
		WithdrawTime:   req.WithdrawTime.AsTime(),
	}

	w.logger.Debug("CreateWithdraw", zap.Any("request", request))

	if err := request.Validate(); err != nil {
		w.logger.Error("CreateWithdraw", zap.Any("error", err))
		return nil, withdraw_errors.ErrGrpcValidateCreateWithdrawRequest
	}

	withdraw, err := w.withdrawCommand.Create(request)

	if err != nil {
		w.logger.Error("CreateWithdraw", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := w.mapping.ToProtoResponseWithdraw("success", "Successfully created withdraw", withdraw)

	return so, nil

}

func (w *withdrawHandleGrpc) UpdateWithdraw(ctx context.Context, req *pb.UpdateWithdrawRequest) (*pb.ApiResponseWithdraw, error) {
	id := int(req.GetWithdrawId())

	w.logger.Debug("UpdateWithdraw", zap.Int("id", id))

	if id == 0 {
		w.logger.Error("UpdateWithdraw", zap.Int("id", id))
		return nil, withdraw_errors.ErrGrpcWithdrawInvalidID
	}

	request := &requests.UpdateWithdrawRequest{
		WithdrawID:     &id,
		CardNumber:     req.CardNumber,
		WithdrawAmount: int(req.WithdrawAmount),
		WithdrawTime:   req.WithdrawTime.AsTime(),
	}

	w.logger.Debug("UpdateWithdraw", zap.Any("request", request))

	if err := request.Validate(); err != nil {
		w.logger.Error("UpdateWithdraw", zap.Any("error", err))
		return nil, withdraw_errors.ErrGrpcValidateUpdateWithdrawRequest
	}

	withdraw, err := w.withdrawCommand.Update(request)

	if err != nil {
		w.logger.Error("UpdateWithdraw", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := w.mapping.ToProtoResponseWithdraw("success", "Successfully updated withdraw", withdraw)

	return so, nil
}

func (w *withdrawHandleGrpc) TrashedWithdraw(ctx context.Context, req *pb.FindByIdWithdrawRequest) (*pb.ApiResponseWithdraw, error) {
	id := int(req.GetWithdrawId())

	w.logger.Debug("TrashedWithdraw", zap.Int("id", id))

	if id == 0 {
		w.logger.Error("TrashedWithdraw", zap.Int("id", id))
		return nil, withdraw_errors.ErrGrpcWithdrawInvalidID
	}

	withdraw, err := w.withdrawCommand.TrashedWithdraw(id)

	if err != nil {
		w.logger.Error("TrashedWithdraw", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := w.mapping.ToProtoResponseWithdraw("success", "Successfully trashed withdraw", withdraw)

	return so, nil
}

func (w *withdrawHandleGrpc) RestoreWithdraw(ctx context.Context, req *pb.FindByIdWithdrawRequest) (*pb.ApiResponseWithdraw, error) {
	id := int(req.GetWithdrawId())

	w.logger.Debug("RestoreWithdraw", zap.Int("id", id))

	if id == 0 {
		w.logger.Error("RestoreWithdraw", zap.Int("id", id))
		return nil, withdraw_errors.ErrGrpcWithdrawInvalidID
	}

	withdraw, err := w.withdrawCommand.RestoreWithdraw(id)

	if err != nil {
		w.logger.Error("RestoreWithdraw", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := w.mapping.ToProtoResponseWithdraw("success", "Successfully restored withdraw", withdraw)

	return so, nil
}

func (w *withdrawHandleGrpc) DeleteWithdrawPermanent(ctx context.Context, req *pb.FindByIdWithdrawRequest) (*pb.ApiResponseWithdrawDelete, error) {
	id := int(req.GetWithdrawId())

	w.logger.Debug("DeleteWithdrawPermanent", zap.Int("id", id))

	if id == 0 {
		w.logger.Error("DeleteWithdrawPermanent", zap.Int("id", id))
		return nil, withdraw_errors.ErrGrpcWithdrawInvalidID
	}

	_, err := w.withdrawCommand.DeleteWithdrawPermanent(id)

	if err != nil {
		w.logger.Error("DeleteWithdrawPermanent", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := w.mapping.ToProtoResponseWithdrawDelete("success", "Successfully deleted withdraw permanently")

	return so, nil
}

func (s *withdrawHandleGrpc) RestoreAllWithdraw(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseWithdrawAll, error) {
	s.logger.Debug("RestoreAllWithdraw")

	_, err := s.withdrawCommand.RestoreAllWithdraw()

	if err != nil {
		s.logger.Error("RestoreAllWithdraw", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseWithdrawAll("success", "Successfully restore all withdraw")

	return so, nil
}

func (s *withdrawHandleGrpc) DeleteAllWithdrawPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseWithdrawAll, error) {
	s.logger.Debug("DeleteAllWithdrawPermanent")

	_, err := s.withdrawCommand.DeleteAllWithdrawPermanent()

	if err != nil {
		s.logger.Error("DeleteAllWithdrawPermanent", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseWithdrawAll("success", "Successfully delete withdraw permanent")

	return so, nil
}
