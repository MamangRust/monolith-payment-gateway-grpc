package topuphandler

import (
	pb "github.com/MamangRust/monolith-payment-gateway-pb/topup"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/topup"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsTopup struct {
	Client *grpc.ClientConn

	E *echo.Echo

	Logger logger.LoggerInterface
}

func RegisterTopupHandler(deps *DepsTopup) {
	mapper := apimapper.NewTopupResponseMapper()

	handlers := []func(){
		setupTopupQueryHandler(deps, mapper.QueryMapper()),
		setupTopupCommandHandler(deps, mapper.CommandMapper()),
		setupTopupStatsMethodHandler(deps, mapper.MethodStatsMapper()),
		setupTopupStatsStatusHandler(deps, mapper.StatusStatsMapper()),
		setupTopupStatsAmountHandler(deps, mapper.AmountStatsMapper()),
	}

	for _, h := range handlers {
		h()
	}
}

func setupTopupQueryHandler(deps *DepsTopup, mapper apimapper.TopupQueryResponseMapper) func() {
	return func() {
		NewTopupQueryHandleApi(
			&topupQueryHandleDeps{
				client: pb.NewTopupQueryServiceClient(deps.Client),
				router: deps.E,
				logger: deps.Logger,
				mapper: mapper,
			},
		)
	}
}

func setupTopupCommandHandler(deps *DepsTopup, mapper apimapper.TopupCommandResponseMapper) func() {
	return func() {
		NewTopupCommandHandleApi(
			&topupCommandHandleDeps{
				client: pb.NewTopupCommandServiceClient(deps.Client),
				router: deps.E,
				logger: deps.Logger,
				mapper: mapper,
			},
		)
	}
}

func setupTopupStatsMethodHandler(deps *DepsTopup, mapper apimapper.TopupStatsMethodResponseMapper) func() {
	return func() {
		NewTopupStatsMethodHandleApi(
			&topupStatsMethodHandleDeps{
				client: pb.NewTopupStatsMethodServiceClient(deps.Client),
				router: deps.E,
				logger: deps.Logger,
				mapper: mapper,
			},
		)
	}
}

func setupTopupStatsAmountHandler(deps *DepsTopup, mapper apimapper.TopupStatsAmountResponseMapper) func() {
	return func() {
		NewTopupStatsAmountHandleApi(
			&topupStatsAmountHandleDeps{
				client: pb.NewTopupStatsAmountServiceClient(deps.Client),
				router: deps.E,
				logger: deps.Logger,
				mapper: mapper,
			},
		)
	}
}

func setupTopupStatsStatusHandler(deps *DepsTopup, mapper apimapper.TopupStatsStatusResponseMapper) func() {
	return func() {
		NewTopupStatsStatusHandleApi(
			&topupStatsStatusHandleDeps{
				client: pb.NewTopupStatsStatusServiceClient(deps.Client),
				router: deps.E,
				logger: deps.Logger,
				mapper: mapper,
			},
		)
	}
}
