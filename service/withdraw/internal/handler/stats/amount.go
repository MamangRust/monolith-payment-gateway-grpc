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

type withdrawAmountHandleGrpc struct {
	pb.UnimplementedWithdrawStatsAmountServiceServer

	withdrawAmount       servicestats.WithdrawStatsService
	withdrawAmounyByCard servicestatsbycard.WithdrawStatsByCardService

	logger logger.LoggerInterface
	mapper protomapper.WithdrawaStatsAmountProtoMapper
}

func NewWithdrawStatsAmountHandleGrpc(
	service service.Service,
	logger logger.LoggerInterface,
	mapper protomapper.WithdrawaStatsAmountProtoMapper,
) WithdrawStatsAmountHandlerGrpc {
	return &withdrawAmountHandleGrpc{
		withdrawAmount:       service,
		withdrawAmounyByCard: service,
		logger:               logger,
		mapper:               mapper,
	}
}

// FindMonthlyWithdraws retrieves the monthly withdraws for the given year.
//
// Parameters:
//   - ctx: the context.Context object for request-scoped values, cancellation, and deadlines.
//   - req: a pointer to a FindYearWithdrawStatus message containing the year.
//
// Returns:
//   - A pointer to an ApiResponseWithdrawMonthAmount message containing the monthly withdraws.
//   - An error, which is non-nil if the operation fails or if the provided year is invalid.
func (w *withdrawAmountHandleGrpc) FindMonthlyWithdraws(ctx context.Context, req *pb.FindYearWithdrawStatus) (*pb.ApiResponseWithdrawMonthAmount, error) {
	year := int(req.GetYear())

	w.logger.Debug("FindMonthlyWithdraws", zap.Int("year", year))

	if year <= 0 {
		w.logger.Error("FindMonthlyWithdraws", zap.Any("error", withdraw_errors.ErrGrpcInvalidYear))
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	withdraws, err := w.withdrawAmount.FindMonthlyWithdraws(ctx, year)

	if err != nil {
		w.logger.Error("FindMonthlyWithdraws", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := w.mapper.ToProtoResponseWithdrawMonthAmount("success", "Successfully fetched monthly withdraws", withdraws)

	return so, nil
}

// FindYearlyWithdraws retrieves the yearly withdraws for the given year.
//
// Parameters:
//   - ctx: the context.Context object for request-scoped values, cancellation, and deadlines.
//   - req: a pointer to a FindYearWithdrawStatus message containing the year.
//
// Returns:
//   - A pointer to an ApiResponseWithdrawYearAmount message containing the yearly withdraws.
//   - An error, which is non-nil if the operation fails or if the provided year is invalid.
func (w *withdrawAmountHandleGrpc) FindYearlyWithdraws(ctx context.Context, req *pb.FindYearWithdrawStatus) (*pb.ApiResponseWithdrawYearAmount, error) {
	year := int(req.GetYear())

	w.logger.Info("FindYearlyWithdraws", zap.Int("year", year))

	if year <= 0 {
		w.logger.Error("FindYearlyWithdraws", zap.Any("error", withdraw_errors.ErrGrpcInvalidYear))
		return nil, withdraw_errors.ErrGrpcInvalidYear
	}

	withdraws, err := w.withdrawAmount.FindYearlyWithdraws(ctx, year)

	if err != nil {
		w.logger.Error("FindYearlyWithdraws", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := w.mapper.ToProtoResponseWithdrawYearAmount("success", "Successfully fetched yearly withdraws", withdraws)

	return so, nil
}

// FindMonthlyWithdrawsByCardNumber retrieves the monthly withdraw amounts for a specific card number and year.
//
// Parameters:
//   - ctx: the context.Context object for request-scoped values, cancellation, and deadlines.
//   - req: a pointer to a FindYearWithdrawCardNumber message containing the year and card number.
//
// Returns:
//   - A pointer to an ApiResponseWithdrawMonthAmount message containing the monthly withdraw amounts.
//   - An error, which is non-nil if the operation fails or if the provided year or card number is invalid.
func (w *withdrawAmountHandleGrpc) FindMonthlyWithdrawsByCardNumber(ctx context.Context, req *pb.FindYearWithdrawCardNumber) (*pb.ApiResponseWithdrawMonthAmount, error) {
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

	reqService := &requests.YearMonthCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	withdraws, err := w.withdrawAmounyByCard.FindMonthlyWithdrawsByCardNumber(ctx, reqService)

	if err != nil {
		w.logger.Error("FindMonthlyWithdrawsByCardNumber", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := w.mapper.ToProtoResponseWithdrawMonthAmount("success", "Successfully fetched monthly withdraws by card number", withdraws)

	return so, nil
}

// FindYearlyWithdrawsByCardNumber retrieves the yearly withdraw amount data
// for a given card number and year.
//
// Parameters:
//   - ctx: the context.Context object for request-scoped values, cancellation, and deadlines.
//   - req: a pointer to a FindYearWithdrawCardNumber message containing the year and card number.
//
// Returns:
//   - A pointer to an ApiResponseWithdrawYearAmount message containing the yearly withdraw amounts.
//   - An error, which is non-nil if the operation fails or if the provided year or card number is invalid.
func (w *withdrawAmountHandleGrpc) FindYearlyWithdrawsByCardNumber(ctx context.Context, req *pb.FindYearWithdrawCardNumber) (*pb.ApiResponseWithdrawYearAmount, error) {
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

	reqService := &requests.YearMonthCardNumber{
		Year:       year,
		CardNumber: cardNumber,
	}

	withdraws, err := w.withdrawAmounyByCard.FindYearlyWithdrawsByCardNumber(ctx, reqService)

	if err != nil {
		w.logger.Error("FindYearlyWithdrawsByCardNumber", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := w.mapper.ToProtoResponseWithdrawYearAmount("success", "Successfully fetched yearly withdraws by card number", withdraws)

	return so, nil
}
