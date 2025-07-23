package saldostatshandler

import (
	"context"

	pb "github.com/MamangRust/monolith-payment-gateway-pb/saldo"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	saldostatsservice "github.com/MamangRust/monolith-payment-gateway-saldo/internal/service/stats"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/grpc"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/saldo"
	"go.uber.org/zap"
)

type saldoStatsTotalBalanceHandleGrpc struct {
	pb.UnimplementedSaldoStatsTotalBalanceServer

	service saldostatsservice.SaldoStatsTotalBalanceService

	mapper protomapper.SaldoStatsTotalSaldoProtoMapper

	logger logger.LoggerInterface
}

func NewSaldoStatsTotalBalanceHandleGrpc(service saldostatsservice.SaldoStatsTotalBalanceService, logger logger.LoggerInterface, mapper protomapper.SaldoStatsTotalSaldoProtoMapper) SaldoStatsTotalBalanceHandleGrpc {
	return &saldoStatsTotalBalanceHandleGrpc{
		service: service,
		mapper:  mapper,
		logger:  logger,
	}
}

// FindMonthlyTotalSaldoBalance is a gRPC handler that fetches the total saldo balance for a specific month and year.
//
// Parameters:
//   - ctx: the context.Context object for tracing and cancellation.
//   - req: a pointer to a FindMonthlySaldoTotalBalance message, which contains the year and month of the saldo balance to be fetched.
//
// Returns:
//   - A pointer to a ApiResponseMonthTotalSaldo message, which contains the total saldo balance for the given month and year.
//   - An error, which is non-nil if the operation fails.
func (s *saldoStatsTotalBalanceHandleGrpc) FindMonthlyTotalSaldoBalance(ctx context.Context, req *pb.FindMonthlySaldoTotalBalance) (*pb.ApiResponseMonthTotalSaldo, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	s.logger.Info("Fetching monthly total saldo balance", zap.Int("year", year), zap.Int("month", month))

	if year <= 0 {
		s.logger.Error("FindMonthlyTotalSaldoBalance failed", zap.Any("error", saldo_errors.ErrGrpcSaldoInvalidYear))
		return nil, saldo_errors.ErrGrpcSaldoInvalidYear
	}

	if month <= 0 {
		s.logger.Error("FindMonthlyTotalSaldoBalance failed", zap.Any("error", saldo_errors.ErrGrpcSaldoInvalidMonth))
		return nil, saldo_errors.ErrGrpcSaldoInvalidMonth
	}

	reqService := &requests.MonthTotalSaldoBalance{
		Year:  year,
		Month: month,
	}

	res, err := s.service.FindMonthlyTotalSaldoBalance(ctx, reqService)

	if err != nil {
		s.logger.Error("FindMonthlyTotalSaldoBalance failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMonthTotalSaldo("success", "Successfully fetched monthly total saldo balance", res)

	s.logger.Info("Successfully fetched monthly total saldo balance", zap.Bool("success", true))

	return so, nil
}

// FindYearTotalSaldoBalance is a gRPC handler that fetches the total saldo balance for a specific year.
//
// Parameters:
//   - ctx: the context.Context object for tracing and cancellation.
//   - req: a pointer to a FindYearlySaldo message, which contains the year of the saldo balance to be fetched.
//
// Returns:
//   - A pointer to a ApiResponseYearTotalSaldo message, which contains the total saldo balance for the given year.
//   - An error, which is non-nil if the operation fails.
func (s *saldoStatsTotalBalanceHandleGrpc) FindYearTotalSaldoBalance(ctx context.Context, req *pb.FindYearlySaldo) (*pb.ApiResponseYearTotalSaldo, error) {
	year := int(req.GetYear())

	s.logger.Info("Fetching yearly total saldo balance", zap.Int("year", year))

	if year <= 0 {
		s.logger.Error("FindYearTotalSaldoBalance failed", zap.Any("error", saldo_errors.ErrGrpcSaldoInvalidYear))
		return nil, saldo_errors.ErrGrpcSaldoInvalidYear
	}

	res, err := s.service.FindYearTotalSaldoBalance(ctx, year)

	if err != nil {
		s.logger.Error("FindYearTotalSaldoBalance failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseYearTotalSaldo("success", "Successfully fetched yearly total saldo balance", res)

	s.logger.Info("Successfully fetched yearly total saldo balance", zap.Bool("success", true))

	return so, nil
}
