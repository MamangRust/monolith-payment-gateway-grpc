package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto"
	"github.com/MamangRust/monolith-payment-gateway-shared/pb"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/service"
	"google.golang.org/protobuf/types/known/emptypb"
)

type transactionHandleGrpc struct {
	pb.UnimplementedTransactionServiceServer
	transactionQueryService           service.TransactionQueryService
	transactionCommandService         service.TransactionCommandService
	transactionStatisticService       service.TransactionStatisticService
	transactionStatisticByCardService service.TransactionsStatistcByCardService
	mapping                           protomapper.TransactionProtoMapper
}

func NewTransactionHandleGrpc(service service.Service) *transactionHandleGrpc {

	return &transactionHandleGrpc{
		transactionQueryService:           service.TransactionQuery,
		transactionCommandService:         service.TransactionCommand,
		transactionStatisticService:       service.TransactionStatistic,
		transactionStatisticByCardService: service.TransactionStatisticByCard,
		mapping:                           protomapper.NewTransactionProtoMapper(),
	}
}

func (t *transactionHandleGrpc) FindAllTransaction(ctx context.Context, request *pb.FindAllTransactionRequest) (*pb.ApiResponsePaginationTransaction, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllTransactions{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	transactions, totalRecords, err := t.transactionQueryService.FindAll(&reqService)

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
	so := t.mapping.ToProtoResponsePaginationTransaction(paginationMeta, "", "", transactions)

	return so, nil
}

func (t *transactionHandleGrpc) FindAllTransactionByCardNumber(ctx context.Context, request *pb.FindAllTransactionCardNumberRequest) (*pb.ApiResponsePaginationTransaction, error) {
	card_number := request.GetCardNumber()
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	reqService := requests.FindAllTransactionCardNumber{
		CardNumber: card_number,
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
	}

	transactions, totalRecords, err := t.transactionQueryService.FindAllByCardNumber(&reqService)

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
	so := t.mapping.ToProtoResponsePaginationTransaction(paginationMeta, "", "", transactions)

	return so, nil
}

func (t *transactionHandleGrpc) FindByIdTransaction(ctx context.Context, req *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransaction, error) {
	id := int(req.GetTransactionId())

	if id == 0 {
		return nil, transaction_errors.ErrGrpcTransactionInvalidID
	}

	transaction, err := t.transactionQueryService.FindById(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapping.ToProtoResponseTransaction("success", "Transaction fetched successfully", transaction)

	return so, nil
}

func (s *transactionHandleGrpc) FindMonthlyTransactionStatusSuccess(ctx context.Context, req *pb.FindMonthlyTransactionStatus) (*pb.ApiResponseTransactionMonthStatusSuccess, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthStatusTransaction{
		Year:  year,
		Month: month,
	}

	records, err := s.transactionStatisticService.FindMonthTransactionStatusSuccess(&reqService)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransactionMonthStatusSuccess("success", "Successfully fetched monthly Transaction status success", records)

	return so, nil
}

func (s *transactionHandleGrpc) FindYearlyTransactionStatusSuccess(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionYearStatusSuccess, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	records, err := s.transactionStatisticService.FindYearlyTransactionStatusSuccess(year)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransactionYearStatusSuccess("success", "Successfully fetched yearly Transaction status success", records)

	return so, nil
}

func (s *transactionHandleGrpc) FindMonthlyTransactionStatusFailed(ctx context.Context, req *pb.FindMonthlyTransactionStatus) (*pb.ApiResponseTransactionMonthStatusFailed, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthStatusTransaction{
		Year:  year,
		Month: month,
	}

	records, err := s.transactionStatisticService.FindMonthTransactionStatusFailed(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransactionMonthStatusFailed("success", "success fetched monthly Transaction status Failed", records)

	return so, nil
}

func (s *transactionHandleGrpc) FindYearlyTransactionStatusFailed(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionYearStatusFailed, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	records, err := s.transactionStatisticService.FindYearlyTransactionStatusFailed(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransactionYearStatusFailed("success", "success fetched yearly Transaction status Failed", records)

	return so, nil
}

func (s *transactionHandleGrpc) FindMonthlyTransactionStatusSuccessByCardNumber(ctx context.Context, req *pb.FindMonthlyTransactionStatusCardNumber) (*pb.ApiResponseTransactionMonthStatusSuccess, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	if cardNumber == "" {
		return nil, transaction_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthStatusTransactionCardNumber{
		CardNumber: cardNumber,
		Year:       year,
		Month:      month,
	}

	records, err := s.transactionStatisticByCardService.FindMonthTransactionStatusSuccessByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransactionMonthStatusSuccess("success", "Successfully fetched monthly Transaction status success", records)

	return so, nil
}

func (s *transactionHandleGrpc) FindYearlyTransactionStatusSuccessByCardNumber(ctx context.Context, req *pb.FindYearTransactionStatusCardNumber) (*pb.ApiResponseTransactionYearStatusSuccess, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		return nil, transaction_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.YearStatusTransactionCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	records, err := s.transactionStatisticByCardService.FindYearlyTransactionStatusSuccessByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransactionYearStatusSuccess("success", "Successfully fetched yearly Transaction status success", records)

	return so, nil
}

func (s *transactionHandleGrpc) FindMonthlyTransactionStatusFailedByCardNumber(ctx context.Context, req *pb.FindMonthlyTransactionStatusCardNumber) (*pb.ApiResponseTransactionMonthStatusFailed, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	if cardNumber == "" {
		return nil, transaction_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthStatusTransactionCardNumber{
		CardNumber: cardNumber,
		Year:       year,
		Month:      month,
	}

	records, err := s.transactionStatisticByCardService.FindMonthTransactionStatusFailedByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransactionMonthStatusFailed("success", "success fetched monthly Transaction status Failed", records)

	return so, nil
}

func (s *transactionHandleGrpc) FindYearlyTransactionStatusFailedByCardNumber(ctx context.Context, req *pb.FindYearTransactionStatusCardNumber) (*pb.ApiResponseTransactionYearStatusFailed, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		return nil, transaction_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.YearStatusTransactionCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	records, err := s.transactionStatisticByCardService.FindYearlyTransactionStatusFailedByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransactionYearStatusFailed("success", "success fetched yearly Transaction status Failed", records)

	return so, nil
}

func (t *transactionHandleGrpc) FindMonthlyPaymentMethods(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionMonthMethod, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	methods, err := t.transactionStatisticService.FindMonthlyPaymentMethods(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapping.ToProtoResponseTransactionMonthMethod("success", "Successfully fetched monthly payment methods", methods)

	return so, nil
}

func (t *transactionHandleGrpc) FindYearlyPaymentMethods(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionYearMethod, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	methods, err := t.transactionStatisticService.FindYearlyPaymentMethods(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapping.ToProtoResponseTransactionYearMethod("success", "Successfully fetched yearly payment methods", methods)

	return so, nil
}

func (t *transactionHandleGrpc) FindMonthlyAmounts(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionMonthAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	amounts, err := t.transactionStatisticService.FindMonthlyAmounts(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapping.ToProtoResponseTransactionMonthAmount("success", "Successfully fetched monthly amounts", amounts)

	return so, nil
}

func (t *transactionHandleGrpc) FindYearlyAmounts(ctx context.Context, req *pb.FindYearTransactionStatus) (*pb.ApiResponseTransactionYearAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	amounts, err := t.transactionStatisticService.FindYearlyAmounts(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapping.ToProtoResponseTransactionYearAmount("success", "Successfully fetched yearly amounts", amounts)

	return so, nil
}

func (t *transactionHandleGrpc) FindMonthlyPaymentMethodsByCardNumber(ctx context.Context, req *pb.FindByYearCardNumberTransactionRequest) (*pb.ApiResponseTransactionMonthMethod, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		return nil, transaction_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearPaymentMethod{
		Year:       year,
		CardNumber: cardNumber,
	}

	methods, err := t.transactionStatisticByCardService.FindMonthlyPaymentMethodsByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapping.ToProtoResponseTransactionMonthMethod("success", "Successfully fetched monthly payment methods by card number", methods)

	return so, nil
}

func (t *transactionHandleGrpc) FindYearlyPaymentMethodsByCardNumber(ctx context.Context, req *pb.FindByYearCardNumberTransactionRequest) (*pb.ApiResponseTransactionYearMethod, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}
	if cardNumber == "" {
		return nil, transaction_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearPaymentMethod{
		Year:       year,
		CardNumber: cardNumber,
	}

	methods, err := t.transactionStatisticByCardService.FindYearlyPaymentMethodsByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapping.ToProtoResponseTransactionYearMethod("success", "Successfully fetched yearly payment methods by card number", methods)

	return so, nil
}

func (t *transactionHandleGrpc) FindMonthlyAmountsByCardNumber(ctx context.Context, req *pb.FindByYearCardNumberTransactionRequest) (*pb.ApiResponseTransactionMonthAmount, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		return nil, transaction_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearPaymentMethod{
		Year:       year,
		CardNumber: cardNumber,
	}

	amounts, err := t.transactionStatisticByCardService.FindMonthlyAmountsByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapping.ToProtoResponseTransactionMonthAmount("success", "Successfully fetched monthly amounts by card number", amounts)

	return so, nil
}

func (t *transactionHandleGrpc) FindYearlyAmountsByCardNumber(ctx context.Context, req *pb.FindByYearCardNumberTransactionRequest) (*pb.ApiResponseTransactionYearAmount, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		return nil, transaction_errors.ErrGrpcInvalidCardNumber
	}

	reqService := requests.MonthYearPaymentMethod{
		Year:       year,
		CardNumber: cardNumber,
	}

	amounts, err := t.transactionStatisticByCardService.FindYearlyAmountsByCardNumber(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapping.ToProtoResponseTransactionYearAmount("success", "Successfully fetched yearly amounts by card number", amounts)

	return so, nil
}

func (t *transactionHandleGrpc) FindTransactionByMerchantIdRequest(ctx context.Context, req *pb.FindTransactionByMerchantIdRequest) (*pb.ApiResponseTransactions, error) {
	id := int(req.GetMerchantId())

	if id == 0 {
		return nil, transaction_errors.ErrGrpcTransactionInvalidMerchantID
	}

	transactions, err := t.transactionQueryService.FindTransactionByMerchantId(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapping.ToProtoResponseTransactions("success", "Successfully fetch transactions", transactions)

	return so, nil
}

func (t *transactionHandleGrpc) FindByActiveTransaction(ctx context.Context, req *pb.FindAllTransactionRequest) (*pb.ApiResponsePaginationTransactionDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllTransactions{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	transactions, totalRecords, err := t.transactionQueryService.FindByActive(&reqService)

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
	so := t.mapping.ToProtoResponsePaginationTransactionDeleteAt(paginationMeta, "success", "Successfully fetch transactions", transactions)

	return so, nil
}

func (t *transactionHandleGrpc) FindByTrashedTransaction(ctx context.Context, req *pb.FindAllTransactionRequest) (*pb.ApiResponsePaginationTransactionDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllTransactions{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	transactions, totalRecords, err := t.transactionQueryService.FindByTrashed(&reqService)

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
	so := t.mapping.ToProtoResponsePaginationTransactionDeleteAt(paginationMeta, "success", "Successfully fetch transactions", transactions)

	return so, nil
}

func (t *transactionHandleGrpc) CreateTransaction(ctx context.Context, request *pb.CreateTransactionRequest) (*pb.ApiResponseTransaction, error) {
	transactionTime := request.GetTransactionTime().AsTime()
	merchantID := int(request.GetMerchantId())

	req := requests.CreateTransactionRequest{
		CardNumber:      request.GetCardNumber(),
		Amount:          int(request.GetAmount()),
		PaymentMethod:   request.GetPaymentMethod(),
		MerchantID:      &merchantID,
		TransactionTime: transactionTime,
	}

	if err := req.Validate(); err != nil {
		return nil, transaction_errors.ErrGrpcFailedCreateTransaction
	}

	res, err := t.transactionCommandService.Create(request.ApiKey, &req)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapping.ToProtoResponseTransaction("success", "Successfully created transaction", res)

	return so, nil
}

func (t *transactionHandleGrpc) UpdateTransaction(ctx context.Context, request *pb.UpdateTransactionRequest) (*pb.ApiResponseTransaction, error) {
	id := int(request.GetTransactionId())

	if id == 0 {
		return nil, transaction_errors.ErrGrpcTransactionInvalidID
	}

	transactionTime := request.GetTransactionTime().AsTime()
	merchantID := int(request.GetMerchantId())

	req := requests.UpdateTransactionRequest{
		TransactionID:   &id,
		CardNumber:      request.GetCardNumber(),
		Amount:          int(request.GetAmount()),
		PaymentMethod:   request.GetPaymentMethod(),
		MerchantID:      &merchantID,
		TransactionTime: transactionTime,
	}

	res, err := t.transactionCommandService.Update(request.ApiKey, &req)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapping.ToProtoResponseTransaction("success", "Successfully updated transaction", res)

	return so, nil
}

func (t *transactionHandleGrpc) TrashedTransaction(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransaction, error) {
	id := int(request.GetTransactionId())

	if id == 0 {
		return nil, transaction_errors.ErrGrpcTransactionInvalidID
	}

	res, err := t.transactionCommandService.TrashedTransaction(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapping.ToProtoResponseTransaction("success", "Successfully trashed transaction", res)

	return so, nil
}

func (t *transactionHandleGrpc) RestoreTransaction(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransaction, error) {
	id := int(request.GetTransactionId())

	if id == 0 {
		return nil, transaction_errors.ErrGrpcTransactionInvalidID
	}

	res, err := t.transactionCommandService.RestoreTransaction(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapping.ToProtoResponseTransaction("success", "Successfully restored transaction", res)

	return so, nil
}

func (t *transactionHandleGrpc) DeleteTransaction(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransactionDelete, error) {
	id := int(request.GetTransactionId())

	if id == 0 {
		return nil, transaction_errors.ErrGrpcTransactionInvalidID
	}

	_, err := t.transactionCommandService.DeleteTransactionPermanent(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapping.ToProtoResponseTransactionDelete("success", "Successfully deleted transaction")

	return so, nil

}

func (t *transactionHandleGrpc) RestoreAllTransaction(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTransactionAll, error) {
	_, err := t.transactionCommandService.RestoreAllTransaction()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapping.ToProtoResponseTransactionAll("success", "Successfully restore all transaction")

	return so, nil
}

func (t *transactionHandleGrpc) DeleteAllTransactionPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTransactionAll, error) {
	_, err := t.transactionCommandService.DeleteAllTransactionPermanent()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapping.ToProtoResponseTransactionAll("success", "Successfully delete transaction permanent")

	return so, nil
}
