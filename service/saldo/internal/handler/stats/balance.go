package saldostatshandler

import (
	"context"

	pb "github.com/MamangRust/monolith-payment-gateway-pb/saldo"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	saldostatsservice "github.com/MamangRust/monolith-payment-gateway-saldo/internal/service/stats"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/grpc"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/saldo"
	"go.uber.org/zap"
)

type saldoStatsBalanceHandleGrpc struct {
	pb.UnimplementedSaldoStatsBalanceServiceServer

	service saldostatsservice.SaldoStatsBalanceService

	mapper protomapper.SaldoStatsBalanceProtoMapper

	logger logger.LoggerInterface
}

func NewSaldoStatsBalanceHandleGrpc(service saldostatsservice.SaldoStatsBalanceService, logger logger.LoggerInterface, mapper protomapper.SaldoStatsBalanceProtoMapper) SaldoStatsBalanceHandleGrpc {
	return &saldoStatsBalanceHandleGrpc{
		service: service,
		mapper:  mapper,
		logger:  logger,
	}
}

// FindMonthlySaldoBalances is a gRPC handler that fetches the monthly saldo balances for a specific year.
//
// Parameters:
//   - ctx: the context.Context object for tracing and cancellation.
//   - req: a pointer to a FindYearlySaldo message, which contains the year of the saldo balances to be fetched.
//
// Returns:
//   - A pointer to a ApiResponseMonthSaldoBalances message, which contains the monthly saldo balances for the given year.
//   - An error, which is non-nil if the operation fails.
func (s *saldoStatsBalanceHandleGrpc) FindMonthlySaldoBalances(ctx context.Context, req *pb.FindYearlySaldo) (*pb.ApiResponseMonthSaldoBalances, error) {
	year := int(req.GetYear())

	if year <= 0 {
		s.logger.Error("FindMonthlySaldoBalances failed", zap.Any("error", saldo_errors.ErrGrpcSaldoInvalidYear))
		return nil, saldo_errors.ErrGrpcSaldoInvalidYear
	}

	res, err := s.service.FindMonthlySaldoBalances(ctx, year)

	if err != nil {
		s.logger.Error("FindMonthlySaldoBalances failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMonthSaldoBalances("success", "Successfully fetched monthly saldo balances", res)

	s.logger.Info("Successfully fetched monthly saldo balances", zap.Bool("success", true))

	return so, nil
}

// FindYearlySaldoBalances is a gRPC handler that fetches the yearly saldo balances for a specific year.
//
// Parameters:
//   - ctx: the context.Context object for tracing and cancellation.
//   - req: a pointer to a FindYearlySaldo message, which contains the year of the saldo balances to be fetched.
//
// Returns:
//   - A pointer to a ApiResponseYearSaldoBalances message, which contains the yearly saldo balances for the given year.
//   - An error, which is non-nil if the operation fails.
func (s *saldoStatsBalanceHandleGrpc) FindYearlySaldoBalances(ctx context.Context, req *pb.FindYearlySaldo) (*pb.ApiResponseYearSaldoBalances, error) {
	year := int(req.GetYear())

	if year <= 0 {
		s.logger.Error("FindYearlySaldoBalances failed", zap.Any("error", saldo_errors.ErrGrpcSaldoInvalidYear))
		return nil, saldo_errors.ErrGrpcSaldoInvalidYear
	}

	res, err := s.service.FindYearlySaldoBalances(ctx, year)

	if err != nil {
		s.logger.Error("FindYearlySaldoBalances failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseYearSaldoBalances("success", "Successfully fetched yearly saldo balances", res)

	s.logger.Info("Successfully fetched yearly saldo balances", zap.Bool("success", true))

	return so, nil
}
