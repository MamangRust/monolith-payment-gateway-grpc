package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-payment-gateway-card/internal/service"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto"
	"github.com/MamangRust/monolith-payment-gateway-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type cardHandleGrpc struct {
	pb.UnimplementedCardServiceServer
	cardQuery           service.CardQueryService
	cardDashboard       service.CardDashboardService
	cardStatistic       service.CardStatisticService
	cardStatisticByCard service.CardStatisticByNumberService
	cardCommand         service.CardCommandService

	mapping protomapper.CardProtoMapper
}

func NewCardHandleGrpc(
	cardService service.Service,
) *cardHandleGrpc {
	return &cardHandleGrpc{
		cardQuery:           cardService.CardQuery,
		cardDashboard:       cardService.CardDashboard,
		cardStatistic:       cardService.CardStatistic,
		cardStatisticByCard: cardService.CardStatisticByCard,
		cardCommand:         cardService.CardCommand,
		mapping:             protomapper.NewCardProtoMapper(),
	}
}

func (s *cardHandleGrpc) FindAllCard(ctx context.Context, req *pb.FindAllCardRequest) (*pb.ApiResponsePaginationCard, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCards{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cards, totalRecords, err := s.cardQuery.FindAll(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationCard(paginationMeta, "success", "Successfully fetched card records", cards)

	return so, nil
}

func (s *cardHandleGrpc) FindByIdCard(ctx context.Context, req *pb.FindByIdCardRequest) (*pb.ApiResponseCard, error) {
	id := int(req.GetCardId())

	if id == 0 {
		return nil, card_errors.ErrGrpcInvalidCardID
	}

	card, err := s.cardQuery.FindById(id)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseCard("success", "Successfully fetched card record", card)

	return so, nil
}

func (s *cardHandleGrpc) FindByUserIdCard(ctx context.Context, req *pb.FindByUserIdCardRequest) (*pb.ApiResponseCard, error) {
	id := int(req.GetUserId())

	if id == 0 {
		return nil, card_errors.ErrGrpcInvalidUserID
	}
	res, err := s.cardQuery.FindByUserID(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseCard("success", "Successfully fetched card record", res)

	return so, nil
}

func (s *cardHandleGrpc) DashboardCard(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseDashboardCard, error) {
	dashboardCard, err := s.cardDashboard.DashboardCard()
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseDashboardCard("success", "Dashboard card retrieved successfully", dashboardCard)

	return so, nil
}

func (s *cardHandleGrpc) DashboardCardNumber(ctx context.Context, req *pb.FindByCardNumberRequest) (*pb.ApiResponseDashboardCardNumber, error) {
	card_number := req.GetCardNumber()

	if card_number == "" {
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	dashboardCard, err := s.cardDashboard.DashboardCardCardNumber(card_number)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseDashboardCardCardNumber("success", "Dashboard card for card number retrieved successfully", dashboardCard)

	return so, nil
}

func (s *cardHandleGrpc) FindMonthlyBalance(ctx context.Context, req *pb.FindYearBalance) (*pb.ApiResponseMonthlyBalance, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}
	res, err := s.cardStatistic.FindMonthlyBalance(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMonthlyBalances("success", "Monthly balance retrieved successfully", res)

	return so, nil
}

func (s *cardHandleGrpc) FindYearlyBalance(ctx context.Context, req *pb.FindYearBalance) (*pb.ApiResponseYearlyBalance, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardStatistic.FindYearlyBalance(year)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseYearlyBalances("success", "Yearly balance retrieved successfully", res)

	return so, nil
}

func (s *cardHandleGrpc) FindMonthlyTopupAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseMonthlyAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardStatistic.FindMonthlyTopupAmount(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMonthlyAmounts("success", "Monthly topup amount retrieved successfully", res)

	return so, nil
}

func (s *cardHandleGrpc) FindYearlyTopupAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseYearlyAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardStatistic.FindYearlyTopupAmount(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseYearlyAmounts("success", "Yearly topup amount retrieved successfully", res)

	return so, nil
}

func (s *cardHandleGrpc) FindMonthlyWithdrawAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseMonthlyAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardStatistic.FindMonthlyWithdrawAmount(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMonthlyAmounts("success", "Monthly withdraw amount retrieved successfully", res)

	return so, nil
}

func (s *cardHandleGrpc) FindYearlyWithdrawAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseYearlyAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardStatistic.FindYearlyWithdrawAmount(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseYearlyAmounts("success", "Yearly withdraw amount retrieved successfully", res)

	return so, nil
}

func (s *cardHandleGrpc) FindMonthlyTransactionAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseMonthlyAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardStatistic.FindMonthlyTransactionAmount(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMonthlyAmounts("success", "Monthly transaction amount retrieved successfully", res)

	return so, nil
}

func (s *cardHandleGrpc) FindYearlyTransactionAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseYearlyAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardStatistic.FindYearlyTransactionAmount(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseYearlyAmounts("success", "Yearly transaction amount retrieved successfully", res)

	return so, nil
}

func (s *cardHandleGrpc) FindMonthlyTransferSenderAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseMonthlyAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardStatistic.FindMonthlyTransferAmountSender(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMonthlyAmounts("success", "Monthly transfer sender amount retrieved successfully", res)

	return so, nil
}

func (s *cardHandleGrpc) FindYearlyTransferSenderAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseYearlyAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardStatistic.FindYearlyTransferAmountSender(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseYearlyAmounts("success", "transfer sender amount retrieved successfully", res)

	return so, nil
}

func (s *cardHandleGrpc) FindMonthlyTransferReceiverAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseMonthlyAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardStatistic.FindMonthlyTransferAmountReceiver(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMonthlyAmounts("success", "Monthly transfer receiver amount retrieved successfully", res)

	return so, nil
}

func (s *cardHandleGrpc) FindYearlyTransferReceiverAmount(ctx context.Context, req *pb.FindYearAmount) (*pb.ApiResponseYearlyAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	res, err := s.cardStatistic.FindYearlyTransferAmountReceiver(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseYearlyAmounts("success", "Yearly transfer receiver amount retrieved successfully", res)

	return so, nil
}

func (s *cardHandleGrpc) FindMonthlyBalanceByCardNumber(ctx context.Context, req *pb.FindYearBalanceCardNumber) (*pb.ApiResponseMonthlyBalance, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	if card_number == "" {
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardStatisticByCard.FindMonthlyBalanceByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMonthlyBalances("success", "Monthly balance retrieved successfully", res)

	return so, nil
}

func (s *cardHandleGrpc) FindYearlyBalanceByCardNumber(ctx context.Context, req *pb.FindYearBalanceCardNumber) (*pb.ApiResponseYearlyBalance, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	if card_number == "" {
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardStatisticByCard.FindYearlyBalanceByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseYearlyBalances("success", "Yearly balance retrieved successfully", res)

	return so, nil
}

func (s *cardHandleGrpc) FindMonthlyTopupAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseMonthlyAmount, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	if card_number == "" {
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardStatisticByCard.FindMonthlyTopupAmountByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMonthlyAmounts("success", "Monthly topup amount by card number retrieved successfully", res)

	return so, nil
}

func (s *cardHandleGrpc) FindYearlyTopupAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseYearlyAmount, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	if card_number == "" {
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       int(year),
	}

	res, err := s.cardStatisticByCard.FindYearlyTopupAmountByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseYearlyAmounts("success", "Yearly topup amount by card number retrieved successfully", res)

	return so, nil
}

func (s *cardHandleGrpc) FindMonthlyWithdrawAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseMonthlyAmount, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	if card_number == "" {
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardStatisticByCard.FindMonthlyWithdrawAmountByCardNumber(&reqService)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMonthlyAmounts("success", "Monthly withdraw amount by card number retrieved successfully", res)

	return so, nil
}

func (s *cardHandleGrpc) FindYearlyWithdrawAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseYearlyAmount, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	if card_number == "" {
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardStatisticByCard.FindYearlyWithdrawAmountByCardNumber(&reqService)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseYearlyAmounts("success", "Yearly withdraw amount by card number retrieved successfully", res)

	return so, nil
}

func (s *cardHandleGrpc) FindMonthlyTransactionAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseMonthlyAmount, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	if card_number == "" {
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardStatisticByCard.FindMonthlyTransactionAmountByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMonthlyAmounts("success", "Monthly transaction amount by card number retrieved successfully", res)

	return so, nil
}

func (s *cardHandleGrpc) FindYearlyTransactionAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseYearlyAmount, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	if card_number == "" {
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardStatisticByCard.FindYearlyTransactionAmountByCardNumber(&reqService)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseYearlyAmounts("success", "Yearly transaction amount by card number retrieved successfully", res)

	return so, nil
}

func (s *cardHandleGrpc) FindMonthlyTransferSenderAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseMonthlyAmount, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	if card_number == "" {
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardStatisticByCard.FindMonthlyTransferAmountBySender(&reqService)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMonthlyAmounts("success", "Monthly transfer sender amount by card number retrieved successfully", res)

	return so, nil
}

func (s *cardHandleGrpc) FindYearlyTransferSenderAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseYearlyAmount, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	if card_number == "" {
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardStatisticByCard.FindYearlyTransferAmountBySender(&reqService)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseYearlyAmounts("success", "Yearly transfer sender amount by card number retrieved successfully", res)

	return so, nil
}

func (s *cardHandleGrpc) FindMonthlyTransferReceiverAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseMonthlyAmount, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	if card_number == "" {
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardStatisticByCard.FindMonthlyTransferAmountByReceiver(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMonthlyAmounts("success", "Monthly transfer receiver amount by card number retrieved successfully", res)

	return so, nil
}

func (s *cardHandleGrpc) FindYearlyTransferReceiverAmountByCardNumber(ctx context.Context, req *pb.FindYearAmountCardNumber) (*pb.ApiResponseYearlyAmount, error) {
	card_number := req.GetCardNumber()
	year := int(req.GetYear())

	if year <= 0 {
		return nil, card_errors.ErrGrpcInvalidYear
	}

	if card_number == "" {
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearCardNumberCard{
		CardNumber: card_number,
		Year:       year,
	}

	res, err := s.cardStatisticByCard.FindYearlyTransferAmountByReceiver(&reqService)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseYearlyAmounts("success", "Yearly transfer receiver amount by card number retrieved successfully", res)

	return so, nil
}

func (s *cardHandleGrpc) FindByActiveCard(ctx context.Context, req *pb.FindAllCardRequest) (*pb.ApiResponsePaginationCardDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCards{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.cardQuery.FindByActive(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationCardDeletedAt(paginationMeta, "success", "Successfully fetched card record", res)

	return so, nil
}

func (s *cardHandleGrpc) FindByTrashedCard(ctx context.Context, req *pb.FindAllCardRequest) (*pb.ApiResponsePaginationCardDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCards{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.cardQuery.FindByTrashed(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationCardDeletedAt(paginationMeta, "success", "Successfully fetched card record", res)

	return so, nil

}

func (s *cardHandleGrpc) FindByCardNumber(ctx context.Context, req *pb.FindByCardNumberRequest) (*pb.ApiResponseCard, error) {
	card_number := req.GetCardNumber()

	if card_number == "" {
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	res, err := s.cardQuery.FindByCardNumber(card_number)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseCard("success", "Successfully fetched card record", res)

	return so, nil

}

func (s *cardHandleGrpc) CreateCard(ctx context.Context, req *pb.CreateCardRequest) (*pb.ApiResponseCard, error) {
	request := requests.CreateCardRequest{
		UserID:       int(req.UserId),
		CardType:     req.CardType,
		ExpireDate:   req.ExpireDate.AsTime(),
		CVV:          req.Cvv,
		CardProvider: req.CardProvider,
	}

	if err := request.Validate(); err != nil {
		return nil, card_errors.ErrGrpcValidateCreateCardRequest
	}

	res, err := s.cardCommand.CreateCard(&request)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseCard("success", "Successfully created card", res)

	return so, nil
}

func (s *cardHandleGrpc) UpdateCard(ctx context.Context, req *pb.UpdateCardRequest) (*pb.ApiResponseCard, error) {
	request := requests.UpdateCardRequest{
		CardID:       int(req.CardId),
		UserID:       int(req.UserId),
		CardType:     req.CardType,
		ExpireDate:   req.ExpireDate.AsTime(),
		CVV:          req.Cvv,
		CardProvider: req.CardProvider,
	}

	if err := request.Validate(); err != nil {
		return nil, card_errors.ErrGrpcValidateUpdateCardRequest
	}

	res, err := s.cardCommand.UpdateCard(&request)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseCard("success", "Successfully updated card", res)

	return so, nil
}

func (s *cardHandleGrpc) TrashedCard(ctx context.Context, req *pb.FindByIdCardRequest) (*pb.ApiResponseCard, error) {
	id := int(req.GetCardId())

	if id == 0 {
		return nil, card_errors.ErrGrpcInvalidCardID
	}

	res, err := s.cardCommand.TrashedCard(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseCard("success", "Successfully trashed card", res)

	return so, nil
}

func (s *cardHandleGrpc) RestoreCard(ctx context.Context, req *pb.FindByIdCardRequest) (*pb.ApiResponseCard, error) {
	id := int(req.GetCardId())

	if id == 0 {
		return nil, card_errors.ErrGrpcInvalidCardID
	}

	res, err := s.cardCommand.RestoreCard(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseCard("success", "Successfully restored card", res)

	return so, nil
}

func (s *cardHandleGrpc) DeleteCardPermanent(ctx context.Context, req *pb.FindByIdCardRequest) (*pb.ApiResponseCardDelete, error) {
	id := int(req.GetCardId())

	if id == 0 {
		return nil, card_errors.ErrGrpcInvalidCardID
	}

	_, err := s.cardCommand.DeleteCardPermanent(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseCardDeleteAt("success", "Successfully deleted card")

	return so, nil
}

func (s *cardHandleGrpc) RestoreAllCard(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseCardAll, error) {
	_, err := s.cardCommand.RestoreAllCard()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseCardAll("success", "Successfully restore card")

	return so, nil
}

func (s *cardHandleGrpc) DeleteAllCardPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseCardAll, error) {
	_, err := s.cardCommand.DeleteAllCardPermanent()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseCardAll("success", "Successfully delete card permanent")

	return so, nil
}
