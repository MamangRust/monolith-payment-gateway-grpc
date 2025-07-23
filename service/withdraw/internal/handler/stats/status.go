package withdrawstatshandler

import (
	"context"

	pb "github.com/MamangRust/monolith-payment-gateway-pb/withdraw"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	withdraw_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors/grpc"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/withdraw"

	service "github.com/MamangRust/monolith-payment-gateway-withdraw/internal/service"

	servicestats "github.com/MamangRust/monolith-payment-gateway-withdraw/internal/service/stats"
	servicestatsbycard "github.com/MamangRust/monolith-payment-gateway-withdraw/internal/service/statsbycard"

	"go.uber.org/zap"
)

type withdrawStatusHandleGrpc struct {
	pb.UnimplementedWithdrawStatsStatusServer

	withdrawStatus       servicestats.WithdrawStatsService
	withdrawStatusByCard servicestatsbycard.WithdrawStatsByCardService

	logger logger.LoggerInterface
	mapper protomapper.WithdrawStatsStatusProtoMapper
}

func NewWithdrawStatsStatusHandleGrpc(
	service service.Service,
	logger logger.LoggerInterface,
	mapper protomapper.WithdrawStatsStatusProtoMapper,
) WithdrawStatsStatusHandleGrpc {
	return &withdrawStatusHandleGrpc{
		withdrawStatus:       service,
		withdrawStatusByCard: service,
		logger:               logger,
		mapper:               mapper,
	}
}

// FindMonthlyWithdrawStatusSuccess retrieves a withdraw record based on the
// provided request, which includes the year and month.
//
// Parameters:
//   - ctx: the context.Context object for tracing and cancellation.
//   - req: a pointer to a FindMonthlyWithdrawStatus message containing the
//     year and month.
//
// Returns:
//   - A pointer to an ApiResponseWithdrawMonthStatusSuccess message containing
//     the withdraw record.
//   - An error, which is non-nil if the operation fails.
func (s *withdrawStatusHandleGrpc) FindMonthlyWithdrawStatusSuccess(ctx context.Context, req *pb.FindMonthlyWithdrawStatus) (*pb.ApiResponseWithdrawMonthStatusSuccess, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	s.logger.Info("FindMonthlyWithdrawStatusSuccess", zap.Int("year", year), zap.Int("month", month))

	if year <= 0 {
		s.logger.Error("FindMonthlyWithdrawStatusSuccess", zap.Any("error", withdraw_errors.ErrGrpcInvalidYear))
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		s.logger.Error("FindMonthlyWithdrawStatusSuccess", zap.Any("error", withdraw_errors.ErrGrpcInvalidMonth))
		return nil, withdraw_errors.ErrGrpcInvalidMonth
	}

	reqService := &requests.MonthStatusWithdraw{
		Year:  year,
		Month: month,
	}

	records, err := s.withdrawStatus.FindMonthWithdrawStatusSuccess(ctx, reqService)

	if err != nil {
		s.logger.Error("FindMonthlyWithdrawStatusSuccess", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseWithdrawMonthStatusSuccess("success", "Successfully fetched withdraw", records)

	return so, nil
}

// FindYearlyWithdrawStatusSuccess retrieves the yearly withdraw status for successful transactions.
//
// Parameters:
//   - ctx: the context.Context object for request-scoped values, cancellation, and deadlines.
//   - req: a pointer to a FindYearWithdrawStatus message containing the year.
//
// Returns:
//   - A pointer to an ApiResponseWithdrawYearStatusSuccess message containing the yearly withdraw status for successful transactions.
//   - An error, which is non-nil if the operation fails or if the provided year is invalid.
func (s *withdrawStatusHandleGrpc) FindYearlyWithdrawStatusSuccess(ctx context.Context, req *pb.FindYearWithdrawStatus) (*pb.ApiResponseWithdrawYearStatusSuccess, error) {
	year := int(req.GetYear())

	s.logger.Info("FindYearlyWithdrawStatusSuccess", zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("FindYearlyWithdrawStatusSuccess", zap.Any("error", withdraw_errors.ErrGrpcInvalidYear))
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	records, err := s.withdrawStatus.FindYearlyWithdrawStatusSuccess(ctx, year)

	if err != nil {
		s.logger.Error("FindYearlyWithdrawStatusSuccess", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseWithdrawYearStatusSuccess("success", "Successfully fetched yearly Withdraw status success", records)

	return so, nil
}

// FindMonthlyWithdrawStatusFailed retrieves the monthly withdraw status for failed transactions.
//
// Parameters:
//   - ctx: the context.Context object for request-scoped values, cancellation, and deadlines.
//   - req: a pointer to a FindMonthlyWithdrawStatus message containing the year and month.
//
// Returns:
//   - A pointer to an ApiResponseWithdrawMonthStatusFailed message containing the monthly withdraw status for failed transactions.
//   - An error, which is non-nil if the operation fails or if the provided year or month is invalid.
func (s *withdrawStatusHandleGrpc) FindMonthlyWithdrawStatusFailed(ctx context.Context, req *pb.FindMonthlyWithdrawStatus) (*pb.ApiResponseWithdrawMonthStatusFailed, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	s.logger.Info("FindMonthlyWithdrawStatusFailed", zap.Int("year", year), zap.Int("month", month))

	if year <= 0 {
		s.logger.Error("FindMonthlyWithdrawStatusFailed", zap.Any("error", withdraw_errors.ErrGrpcInvalidYear))
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	if month <= 0 {
		s.logger.Error("FindMonthlyWithdrawStatusFailed", zap.Any("error", withdraw_errors.ErrGrpcInvalidMonth))
		return nil, withdraw_errors.ErrGrpcInvalidMonth
	}

	reqService := &requests.MonthStatusWithdraw{
		Year:  year,
		Month: month,
	}

	records, err := s.withdrawStatus.FindMonthWithdrawStatusFailed(ctx, reqService)

	if err != nil {
		s.logger.Error("FindMonthlyWithdrawStatusFailed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseWithdrawMonthStatusFailed("success", "success fetched monthly Withdraw status Failed", records)

	return so, nil
}

// FindYearlyWithdrawStatusFailed retrieves the yearly withdraw status for failed transactions.
//
// Parameters:
//   - ctx: the context.Context object for request-scoped values, cancellation, and deadlines.
//   - req: a pointer to a FindYearWithdrawStatus message containing the year.
//
// Returns:
//   - A pointer to an ApiResponseWithdrawYearStatusFailed message containing the yearly withdraw status for failed transactions.
//   - An error, which is non-nil if the operation fails or if the provided year is invalid.
func (s *withdrawStatusHandleGrpc) FindYearlyWithdrawStatusFailed(ctx context.Context, req *pb.FindYearWithdrawStatus) (*pb.ApiResponseWithdrawYearStatusFailed, error) {
	year := int(req.GetYear())

	s.logger.Info("FindYearlyWithdrawStatusFailed", zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("FindYearlyWithdrawStatusFailed", zap.Any("error", withdraw_errors.ErrGrpcInvalidYear))
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	records, err := s.withdrawStatus.FindYearlyWithdrawStatusFailed(ctx, year)

	if err != nil {
		s.logger.Error("FindYearlyWithdrawStatusFailed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseWithdrawYearStatusFailed("success", "success fetched yearly Withdraw status Failed", records)

	return so, nil
}

// FindMonthlyWithdrawStatusSuccessCardNumber retrieves the monthly withdraw status for successful transactions by card number.
//
// Parameters:
//   - ctx: the context.Context object for request-scoped values, cancellation, and deadlines.
//   - req: a pointer to a FindMonthlyWithdrawStatusCardNumber message containing the year, month, and card number.
//
// Returns:
//   - A pointer to an ApiResponseWithdrawMonthStatusSuccess message containing the monthly withdraw status for successful transactions.
//   - An error, which is non-nil if the operation fails or if the provided year, month, or card number is invalid.
func (s *withdrawStatusHandleGrpc) FindMonthlyWithdrawStatusSuccessCardNumber(ctx context.Context, req *pb.FindMonthlyWithdrawStatusCardNumber) (*pb.ApiResponseWithdrawMonthStatusSuccess, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	cardNumber := req.GetCardNumber()

	s.logger.Info("FindMonthlyWithdrawStatusSuccessCardNumber", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", cardNumber))

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

	reqService := &requests.MonthStatusWithdrawCardNumber{
		Year:       year,
		Month:      month,
		CardNumber: cardNumber,
	}

	records, err := s.withdrawStatusByCard.FindMonthWithdrawStatusSuccessByCardNumber(ctx, reqService)

	if err != nil {
		s.logger.Error("FindMonthlyWithdrawStatusSuccessCardNumber", zap.Any("error", err))

		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseWithdrawMonthStatusSuccess("success", "Successfully fetched withdraw", records)

	return so, nil
}

// FindYearlyWithdrawStatusSuccessCardNumber retrieves the yearly withdraw status for successful transactions by card number.
//
// Parameters:
//   - ctx: the context.Context object for request-scoped values, cancellation, and deadlines.
//   - req: a pointer to a FindYearWithdrawStatusCardNumber message containing the year and card number.
//
// Returns:
//   - A pointer to an ApiResponseWithdrawYearStatusSuccess message containing the yearly withdraw status for successful transactions.
//   - An error, which is non-nil if the operation fails or if the provided year or card number is invalid.
func (s *withdrawStatusHandleGrpc) FindYearlyWithdrawStatusSuccessByCardNumber(ctx context.Context, req *pb.FindYearWithdrawStatusCardNumber) (*pb.ApiResponseWithdrawYearStatusSuccess, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	s.logger.Info("FindYearlyWithdrawStatusSuccessCardNumber", zap.Int("year", year), zap.String("card_number", cardNumber))

	if year <= 0 {
		s.logger.Error("FindYearlyWithdrawStatusSuccessCardNumber", zap.Any("error", withdraw_errors.ErrGrpcInvalidYear))
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	if cardNumber == "" {
		s.logger.Error("FindYearlyWithdrawStatusSuccessCardNumber", zap.Any("error", withdraw_errors.ErrGrpcInvalidCardNumber))
		return nil, withdraw_errors.ErrGrpcInvalidCardNumber
	}

	reqService := &requests.YearStatusWithdrawCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	records, err := s.withdrawStatusByCard.FindYearlyWithdrawStatusSuccessByCardNumber(ctx, reqService)

	if err != nil {
		s.logger.Error("FindYearlyWithdrawStatusSuccessCardNumber", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseWithdrawYearStatusSuccess("success", "Successfully fetched yearly Withdraw status success", records)

	return so, nil
}

// FindMonthlyWithdrawStatusFailedCardNumber retrieves the monthly withdraw status for failed transactions
// by card number, year, and month.
//
// Parameters:
//   - ctx: the context.Context object for request-scoped values, cancellation, and deadlines.
//   - req: a pointer to a FindMonthlyWithdrawStatusCardNumber message containing the year, month, and card number.
//
// Returns:
//   - A pointer to an ApiResponseWithdrawMonthStatusFailed message containing the monthly withdraw status for failed transactions.
//   - An error, which is non-nil if the operation fails or if the provided year, month, or card number is invalid.
func (s *withdrawStatusHandleGrpc) FindMonthlyWithdrawStatusFailedCardNumber(ctx context.Context, req *pb.FindMonthlyWithdrawStatusCardNumber) (*pb.ApiResponseWithdrawMonthStatusFailed, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	cardNumber := req.GetCardNumber()

	s.logger.Info("FindMonthlyWithdrawStatusFailedCardNumber", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", cardNumber))

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

	reqService := &requests.MonthStatusWithdrawCardNumber{
		Year:       year,
		Month:      month,
		CardNumber: cardNumber,
	}

	records, err := s.withdrawStatusByCard.FindMonthWithdrawStatusFailedByCardNumber(ctx, reqService)

	if err != nil {
		s.logger.Error("FindMonthlyWithdrawStatusFailedCardNumber", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseWithdrawMonthStatusFailed("success", "Successfully fetched monthly Withdraw status failed", records)

	return so, nil
}

// FindYearlyWithdrawStatusFailedCardNumber retrieves the yearly withdraw status for failed transactions by card number.
//
// Parameters:
//   - ctx: the context.Context object for request-scoped values, cancellation, and deadlines.
//   - req: a pointer to a FindYearWithdrawStatusCardNumber message containing the year and card number.
//
// Returns:
//   - A pointer to an ApiResponseWithdrawYearStatusFailed message containing the yearly withdraw status for failed transactions.
//   - An error, which is non-nil if the operation fails or if the provided year or card number is invalid.
func (s *withdrawStatusHandleGrpc) FindYearlyWithdrawStatusFailedByCardNumber(ctx context.Context, req *pb.FindYearWithdrawStatusCardNumber) (*pb.ApiResponseWithdrawYearStatusFailed, error) {
	year := int(req.GetYear())
	cardNumber := req.GetCardNumber()

	s.logger.Info("FindYearlyWithdrawStatusFailedCardNumber", zap.Int("year", year), zap.String("card_number", cardNumber))

	if year <= 0 {
		s.logger.Error("FindYearlyWithdrawStatusFailedCardNumber", zap.Any("error", withdraw_errors.ErrGrpcInvalidYear))
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	reqService := &requests.YearStatusWithdrawCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	records, err := s.withdrawStatusByCard.FindYearlyWithdrawStatusFailedByCardNumber(ctx, reqService)
	if err != nil {
		s.logger.Error("FindYearlyWithdrawStatusFailedCardNumber", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseWithdrawYearStatusFailed("success", "Successfully fetched yearly Withdraw status failed", records)

	return so, nil
}
