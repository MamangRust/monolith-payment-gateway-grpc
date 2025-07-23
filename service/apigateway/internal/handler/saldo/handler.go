package saldohandler

import (
	pb "github.com/MamangRust/monolith-payment-gateway-pb/saldo"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/saldo"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsSaldo struct {
	Client *grpc.ClientConn

	E *echo.Echo

	Logger logger.LoggerInterface
}

func RegisterSaldoHandler(deps *DepsSaldo) {
	mapper := apimapper.NewSaldoResponseMapper()

	handlers := []func(){
		setupSaldoQueryHandler(deps, mapper.QueryMapper()),
		setupSaldoCommandHandler(deps, mapper.CommandMapper()),
		setupSaldoStatsBalanceHandler(deps, mapper.BalanceStatsMapper()),
		setupStatsSaldoTotalBalanceHandler(deps, mapper.TotalStatsMapper()),
	}

	for _, h := range handlers {
		h()
	}
}

func setupSaldoQueryHandler(deps *DepsSaldo, mapper apimapper.SaldoQueryResponseMapper) func() {
	return func() {
		NewSaldoQueryHandleApi(
			&saldoQueryHandleDeps{
				client: pb.NewSaldoQueryServiceClient(deps.Client),
				router: deps.E,
				logger: deps.Logger,
				mapper: mapper,
			},
		)
	}
}

func setupSaldoCommandHandler(deps *DepsSaldo, mapper apimapper.SaldoCommandResponseMapper) func() {
	return func() {
		NewSaldoCommandHandleApi(
			&saldoCommandHandleDeps{
				client: pb.NewSaldoCommandServiceClient(deps.Client),
				router: deps.E,
				logger: deps.Logger,
				mapper: mapper,
			},
		)
	}
}

func setupSaldoStatsBalanceHandler(deps *DepsSaldo, mapper apimapper.SaldoStatsBalanceResponseMapper) func() {
	return func() {
		NewSaldoStatsBalanceHandleApi(
			&saldoStatsBalanceHandleDeps{
				client: pb.NewSaldoStatsBalanceServiceClient(deps.Client),
				router: deps.E,
				logger: deps.Logger,
				mapper: mapper,
			},
		)
	}
}

func setupStatsSaldoTotalBalanceHandler(deps *DepsSaldo, mapper apimapper.SaldoStatsTotalResponseMapper) func() {
	return func() {
		NewSaldoTotalBalanceHandleApi(
			&saldoTotalBalanceHandleDeps{
				client: pb.NewSaldoStatsTotalBalanceClient(deps.Client),
				router: deps.E,
				logger: deps.Logger,
				mapper: mapper,
			},
		)
	}
}
