package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto"
	"github.com/MamangRust/monolith-payment-gateway-shared/pb"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/service"
	"google.golang.org/protobuf/types/known/emptypb"
)

type topupHandleGrpc struct {
	pb.UnimplementedTopupServiceServer
	topupQuery        service.TopupQueryService
	topupStatistic    service.TopupStatisticService
	topupStatisByCard service.TopupStatisticByCardService
	topupCommand      service.TopupCommandService
	mapping           protomapper.TopupProtoMapper
}

func NewTopupHandleGrpc(service service.Service) *topupHandleGrpc {
	return &topupHandleGrpc{
		topupQuery:        service.TopupQuery,
		topupStatistic:    service.TopupStatistic,
		topupStatisByCard: service.TopupStatisticByCard,
		topupCommand:      service.TopupCommand,
		mapping:           protomapper.NewTopupProtoMapper(),
	}
}

func (s *topupHandleGrpc) FindAllTopup(ctx context.Context, req *pb.FindAllTopupRequest) (*pb.ApiResponsePaginationTopup, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllTopups{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	topups, totalRecords, err := s.topupQuery.FindAll(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapping.ToProtoResponsePaginationTopup(paginationMeta, "success", "Successfully fetch topups", topups)

	return so, nil
}

func (s *topupHandleGrpc) FindAllTopupByCardNumber(ctx context.Context, req *pb.FindAllTopupByCardNumberRequest) (*pb.ApiResponsePaginationTopup, error) {
	card_number := req.GetCardNumber()
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllTopupsByCardNumber{
		CardNumber: card_number,
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
	}

	topups, totalRecords, err := s.topupQuery.FindAllByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapping.ToProtoResponsePaginationTopup(paginationMeta, "success", "Successfully fetch topups", topups)

	return so, nil
}

func (s *topupHandleGrpc) FindByIdTopup(ctx context.Context, req *pb.FindByIdTopupRequest) (*pb.ApiResponseTopup, error) {
	id := int(req.GetTopupId())

	if id == 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidID
	}

	topup, err := s.topupQuery.FindById(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTopup("success", "Successfully fetch topup", topup)

	return so, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupStatusSuccess(ctx context.Context, req *pb.FindMonthlyTopupStatus) (*pb.ApiResponseTopupMonthStatusSuccess, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	if month <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidMonth
	}

	reqService := requests.MonthTopupStatus{
		Year:  year,
		Month: month,
	}

	records, err := s.topupStatistic.FindMonthTopupStatusSuccess(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTopupMonthStatusSuccess("success", "Successfully fetched monthly topup status success", records)

	return so, nil
}

func (s *topupHandleGrpc) FindYearlyTopupStatusSuccess(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupYearStatusSuccess, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	records, err := s.topupStatistic.FindYearlyTopupStatusSuccess(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTopupYearStatusSuccess("success", "Successfully fetched yearly topup status success", records)

	return so, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupStatusFailed(ctx context.Context, req *pb.FindMonthlyTopupStatus) (*pb.ApiResponseTopupMonthStatusFailed, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	if month <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidMonth
	}

	reqService := requests.MonthTopupStatus{
		Year:  year,
		Month: month,
	}

	records, err := s.topupStatistic.FindMonthTopupStatusFailed(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTopupMonthStatusFailed("Successfully", "Successfully fetched monthly topup status Failed", records)

	return so, nil
}

func (s *topupHandleGrpc) FindYearlyTopupStatusFailed(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupYearStatusFailed, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	records, err := s.topupStatistic.FindYearlyTopupStatusFailed(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTopupYearStatusFailed("Successfully", "Successfully fetched yearly topup status Failed", records)

	return so, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupStatusSuccessByCardNumber(ctx context.Context, req *pb.FindMonthlyTopupStatusCardNumber) (*pb.ApiResponseTopupMonthStatusSuccess, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	if month <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidMonth
	}

	if cardNumber == "" {
		return nil, topup_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthTopupStatusCardNumber{
		Year:       year,
		Month:      month,
		CardNumber: cardNumber,
	}

	records, err := s.topupStatisByCard.FindMonthTopupStatusSuccessByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTopupMonthStatusSuccess("success", "Successfully fetched monthly topup status success", records)
	return so, nil
}

func (s *topupHandleGrpc) FindYearlyTopupStatusSuccessByCardNumber(ctx context.Context, req *pb.FindYearTopupStatusCardNumber) (*pb.ApiResponseTopupYearStatusSuccess, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	if cardNumber == "" {
		return nil, topup_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.YearTopupStatusCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	records, err := s.topupStatisByCard.FindYearlyTopupStatusSuccessByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTopupYearStatusSuccess("success", "Successfully fetched yearly topup status success", records)
	return so, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupStatusFailedByCardNumber(ctx context.Context, req *pb.FindMonthlyTopupStatusCardNumber) (*pb.ApiResponseTopupMonthStatusFailed, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	if month <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidMonth
	}

	if cardNumber == "" {
		return nil, topup_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthTopupStatusCardNumber{
		Year:       year,
		Month:      month,
		CardNumber: cardNumber,
	}

	records, err := s.topupStatisByCard.FindMonthTopupStatusFailedByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTopupMonthStatusFailed("success", "Successfully fetched monthly topup status failed", records)
	return so, nil
}

func (s *topupHandleGrpc) FindYearlyTopupStatusFailedByCardNumber(ctx context.Context, req *pb.FindYearTopupStatusCardNumber) (*pb.ApiResponseTopupYearStatusFailed, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	if cardNumber == "" {
		return nil, topup_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.YearTopupStatusCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	records, err := s.topupStatisByCard.FindYearlyTopupStatusFailedByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTopupYearStatusFailed("success", "Successfully fetched yearly topup status failed", records)
	return so, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupMethods(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupMonthMethod, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	methods, err := s.topupStatistic.FindMonthlyTopupMethods(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTopupMonthMethod("success", "Successfully fetched monthly topup methods", methods)

	return so, nil
}

func (s *topupHandleGrpc) FindYearlyTopupMethods(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupYearMethod, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	methods, err := s.topupStatistic.FindYearlyTopupMethods(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTopupYearMethod("success", "Successfully fetched yearly topup methods", methods)

	return so, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupAmounts(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupMonthAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	amounts, err := s.topupStatistic.FindMonthlyTopupAmounts(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTopupMonthAmount("success", "Successfully fetched monthly topup amounts", amounts)

	return so, nil
}

func (s *topupHandleGrpc) FindYearlyTopupAmounts(ctx context.Context, req *pb.FindYearTopupStatus) (*pb.ApiResponseTopupYearAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	amounts, err := s.topupStatistic.FindYearlyTopupAmounts(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTopupYearAmount("success", "Successfully fetched yearly topup amounts", amounts)

	return so, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupMethodsByCardNumber(ctx context.Context, req *pb.FindYearTopupCardNumber) (*pb.ApiResponseTopupMonthMethod, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	if cardNumber == "" {
		return nil, topup_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.YearMonthMethod{
		Year:       year,
		CardNumber: cardNumber,
	}

	methods, err := s.topupStatisByCard.FindMonthlyTopupMethodsByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTopupMonthMethod("success", "Successfully fetched monthly topup methods by card number", methods)

	return so, nil
}

func (s *topupHandleGrpc) FindYearlyTopupMethodsByCardNumber(ctx context.Context, req *pb.FindYearTopupCardNumber) (*pb.ApiResponseTopupYearMethod, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	if cardNumber == "" {
		return nil, topup_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.YearMonthMethod{
		Year:       year,
		CardNumber: cardNumber,
	}

	methods, err := s.topupStatisByCard.FindYearlyTopupMethodsByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTopupYearMethod("success", "Successfully fetched yearly topup methods by card number", methods)

	return so, nil
}

func (s *topupHandleGrpc) FindMonthlyTopupAmountsByCardNumber(ctx context.Context, req *pb.FindYearTopupCardNumber) (*pb.ApiResponseTopupMonthAmount, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	if cardNumber == "" {
		return nil, topup_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.YearMonthMethod{
		Year:       year,
		CardNumber: cardNumber,
	}

	amounts, err := s.topupStatisByCard.FindMonthlyTopupAmountsByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTopupMonthAmount("success", "Successfully fetched monthly topup amounts by card number", amounts)

	return so, nil
}

func (s *topupHandleGrpc) FindYearlyTopupAmountsByCardNumber(ctx context.Context, req *pb.FindYearTopupCardNumber) (*pb.ApiResponseTopupYearAmount, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidYear
	}

	if cardNumber == "" {
		return nil, topup_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.YearMonthMethod{
		Year:       year,
		CardNumber: cardNumber,
	}

	amounts, err := s.topupStatisByCard.FindYearlyTopupAmountsByCardNumber(&reqService)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTopupYearAmount("success", "Successfully fetched yearly topup amounts by card number", amounts)

	return so, nil
}

func (s *topupHandleGrpc) FindByActive(ctx context.Context, req *pb.FindAllTopupRequest) (*pb.ApiResponsePaginationTopupDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllTopups{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.topupQuery.FindByActive(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapping.ToProtoResponsePaginationTopupDeleteAt(paginationMeta, "success", "Successfully fetch topups", res)

	return so, nil
}

func (s *topupHandleGrpc) FindByTrashed(ctx context.Context, req *pb.FindAllTopupRequest) (*pb.ApiResponsePaginationTopupDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllTopups{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.topupQuery.FindByTrashed(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationTopupDeleteAt(paginationMeta, "success", "Successfully fetch topups", res)

	return so, nil
}

func (s *topupHandleGrpc) CreateTopup(ctx context.Context, req *pb.CreateTopupRequest) (*pb.ApiResponseTopup, error) {
	request := requests.CreateTopupRequest{
		CardNumber:  req.GetCardNumber(),
		TopupAmount: int(req.GetTopupAmount()),
		TopupMethod: req.GetTopupMethod(),
	}

	if err := request.Validate(); err != nil {
		return nil, topup_errors.ErrGrpcValidateCreateTopup
	}

	res, err := s.topupCommand.CreateTopup(&request)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTopup("success", "Successfully created topup", res)

	return so, nil
}

func (s *topupHandleGrpc) UpdateTopup(ctx context.Context, req *pb.UpdateTopupRequest) (*pb.ApiResponseTopup, error) {
	id := int(req.GetTopupId())

	if id == 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidID
	}

	request := requests.UpdateTopupRequest{
		TopupID:     &id,
		CardNumber:  req.GetCardNumber(),
		TopupAmount: int(req.GetTopupAmount()),
		TopupMethod: req.GetTopupMethod(),
	}

	if err := request.Validate(); err != nil {
		return nil, topup_errors.ErrGrpcValidateUpdateTopup
	}

	res, err := s.topupCommand.UpdateTopup(&request)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTopup("success", "Successfully updated topup", res)

	return so, nil
}

func (s *topupHandleGrpc) TrashedTopup(ctx context.Context, req *pb.FindByIdTopupRequest) (*pb.ApiResponseTopupDeleteAt, error) {
	id := int(req.GetTopupId())

	if id == 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidID
	}

	res, err := s.topupCommand.TrashedTopup(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTopupDeletAt("success", "Successfully trashed topup", res)

	return so, nil
}

func (s *topupHandleGrpc) RestoreTopup(ctx context.Context, req *pb.FindByIdTopupRequest) (*pb.ApiResponseTopupDeleteAt, error) {
	id := int(req.GetTopupId())

	if id == 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidID
	}

	res, err := s.topupCommand.RestoreTopup(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTopupDeletAt("success", "Successfully restored topup", res)

	return so, nil
}

func (s *topupHandleGrpc) DeleteTopupPermanent(ctx context.Context, req *pb.FindByIdTopupRequest) (*pb.ApiResponseTopupDelete, error) {
	id := int(req.GetTopupId())

	if id == 0 {
		return nil, topup_errors.ErrGrpcTopupInvalidID
	}

	_, err := s.topupCommand.DeleteTopupPermanent(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTopupDelete("success", "Successfully deleted topup permanently")

	return so, nil
}

func (s *topupHandleGrpc) RestoreAllTopup(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTopupAll, error) {
	_, err := s.topupCommand.RestoreAllTopup()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTopupAll("success", "Successfully restore all topup")

	return so, nil
}

func (s *topupHandleGrpc) DeleteAllTopupPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTopupAll, error) {
	_, err := s.topupCommand.DeleteAllTopupPermanent()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTopupAll("success", "Successfully delete topup permanent")

	return so, nil
}
