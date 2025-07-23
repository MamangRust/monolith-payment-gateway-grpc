package withdrawhandler

import (
	pb "github.com/MamangRust/monolith-payment-gateway-pb/withdraw"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/withdraw"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsWithdraw struct {
	Client *grpc.ClientConn

	E *echo.Echo

	Logger logger.LoggerInterface
}

func RegisterWithdrawHandler(deps *DepsWithdraw) {
	mapper := apimapper.NewWithdrawResponseMapper()

	handlers := []func(){
		setupWithdrawQueryHandler(deps, mapper.QueryMapper()),
		setupWithdrawCommandHandler(deps, mapper.CommandMapper()),
		setupWithdrawStatsAmountHandler(deps, mapper.AmountStatsMapper()),
		setupWithdrawStatsStatusHandler(deps, mapper.StatusStatsMapper()),
	}

	for _, h := range handlers {
		h()
	}
}

func setupWithdrawQueryHandler(deps *DepsWithdraw, mapper apimapper.WithdrawQueryResponseMapper) func() {
	return func() {
		NewWithdrawQueryHandleApi(&withdrawQueryHandleDeps{
			client: pb.NewWithdrawQueryServiceClient(deps.Client),
			router: deps.E,
			logger: deps.Logger,
			mapper: mapper,
		})
	}
}

func setupWithdrawCommandHandler(deps *DepsWithdraw, mapper apimapper.WithdrawCommandResponseMapper) func() {
	return func() {
		NewWithdrawCommandHandleApi(&withdrawCommandHandleDeps{
			client: pb.NewWithdrawCommandServiceClient(deps.Client),
			router: deps.E,
			logger: deps.Logger,
			mapper: mapper,
		})
	}
}

func setupWithdrawStatsAmountHandler(deps *DepsWithdraw, mapper apimapper.WithdrawStatsAmountResponseMapper) func() {
	return func() {
		NewWithdrawStatsAmountHandleApi(&withdrawStatsAmountHandleDeps{
			client: pb.NewWithdrawStatsAmountServiceClient(deps.Client),
			router: deps.E,
			logger: deps.Logger,
			mapper: mapper,
		})
	}
}

func setupWithdrawStatsStatusHandler(deps *DepsWithdraw, mapper apimapper.WithdrawStatsStatusResponseMapper) func() {
	return func() {
		NewWithdrawStatsStatusHandleApi(&withdrawStatsStatusHandleDeps{
			client: pb.NewWithdrawStatsStatusClient(deps.Client),
			router: deps.E,
			logger: deps.Logger,
			mapper: mapper,
		})
	}
}
