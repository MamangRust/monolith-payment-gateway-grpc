package transferhandler

import (
	pb "github.com/MamangRust/monolith-payment-gateway-pb/transfer"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/transfer"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsTransfer struct {
	client *grpc.ClientConn

	E *echo.Echo

	logger logger.LoggerInterface
}

func RegisterTransferHandler(deps *DepsTransfer) {
	mapper := apimapper.NewTransferResponseMapper()

	handlers := []func(){
		setupTransferQueryHandler(deps, mapper.QueryMapper()),
		setupTransferCommandHandler(deps, mapper.CommandMapper()),
		setupTransferStatsAmountHandler(deps, mapper.AmountStatsMapper()),
		setupTransferStatsStatusHandler(deps, mapper.StatusStatsMapper()),
	}

	for _, h := range handlers {
		h()
	}
}

func setupTransferQueryHandler(deps *DepsTransfer, mapper apimapper.TransferQueryResponseMapper) func() {
	return func() {
		NewTransferQueryHandleApi(&transferQueryHandleDeps{
			client: pb.NewTransferQueryServiceClient(deps.client),
			router: deps.E,
			logger: deps.logger,
			mapper: mapper,
		})
	}
}

func setupTransferCommandHandler(deps *DepsTransfer, mapper apimapper.TransferCommandResponseMapper) func() {
	return func() {
		NewTransferCommandHandleApi(&transferCommandHandleDeps{
			client: pb.NewTransferCommandServiceClient(deps.client),
			router: deps.E,
			logger: deps.logger,
			mapper: mapper,
		})
	}
}

func setupTransferStatsAmountHandler(deps *DepsTransfer, mapper apimapper.TransferStatsAmountResponseMapper) func() {
	return func() {
		NewTransferStatsAmountHandleApi(&transferStatsAmountHandleDeps{
			client: pb.NewTransferStatsAmountServiceClient(deps.client),
			router: deps.E,
			logger: deps.logger,
			mapper: mapper,
		})
	}
}

func setupTransferStatsStatusHandler(deps *DepsTransfer, mapper apimapper.TransferStatsStatusResponseMapper) func() {
	return func() {
		NewTransferStatsStatusHandleApi(&transferStatsStatusHandleDeps{
			client: pb.NewTransferStatsStatusServiceClient(deps.client),
			router: deps.E,
			logger: deps.logger,
			mapper: mapper,
		})
	}
}
