package transactionhandler

import (
	mencache "github.com/MamangRust/monolith-payment-gateway-apigateway/internal/redis"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/transaction"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/transaction"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsTransaction struct {
	Client *grpc.ClientConn

	E *echo.Echo

	Kafka *kafka.Kafka

	Logger logger.LoggerInterface

	Cache mencache.MerchantCache
}

func RegisterTransactionHandler(deps *DepsTransaction) {
	mapper := apimapper.NewTransactionResponseMapper()

	handlers := []func(){
		setupTransactionQueryHandler(deps, mapper.QueryMapper()),
		setupTransactionCommandHandler(deps, deps.Cache, mapper.CommandMapper()),
		setupTransactionStatsAmountHandler(deps, mapper.AmountStatsMapper()),
		setupTransactionStatsStatusHandler(deps, mapper.StatusStatsMapper()),
		setupTransactionStatsMethodHandler(deps, mapper.MethodStatsMapper()),
	}

	for _, h := range handlers {
		h()
	}
}

func setupTransactionQueryHandler(deps *DepsTransaction, mapper apimapper.TransactionQueryResponseMapper) func() {
	return func() {
		NewTransactionQueryHandleApi(&transactionQueryHandleDeps{
			client: pb.NewTransactionQueryServiceClient(deps.Client),
			router: deps.E,
			logger: deps.Logger,
			mapper: mapper,
		})
	}
}

func setupTransactionCommandHandler(deps *DepsTransaction, cache mencache.MerchantCache, mapper apimapper.TransactionCommandResponseMapper) func() {
	return func() {
		NewTransactionCommandHandleApi(&transactionCommandHandleDeps{
			kafka:  deps.Kafka,
			client: pb.NewTransactionCommandServiceClient(deps.Client),
			router: deps.E,
			logger: deps.Logger,
			mapper: mapper,
			cache:  cache,
		})
	}
}

func setupTransactionStatsAmountHandler(deps *DepsTransaction, mapper apimapper.TransactionStatsAmountResponseMapper) func() {
	return func() {
		NewTransactionStatsAmountHandleApi(&transactionStatsAmountHandleDeps{
			client: pb.NewTransactionsStatsAmountServiceClient(deps.Client),
			router: deps.E,
			logger: deps.Logger,
			mapper: mapper,
		})
	}
}

func setupTransactionStatsMethodHandler(deps *DepsTransaction, mapper apimapper.TransactionStatsMethodResponseMapper) func() {
	return func() {
		NewTransactionStatsMethodHandleApi(&transactionStatsMethodHandleDeps{
			client: pb.NewTransactionStatsMethodServiceClient(deps.Client),
			router: deps.E,
			logger: deps.Logger,
			mapper: mapper,
		})
	}
}

func setupTransactionStatsStatusHandler(deps *DepsTransaction, mapper apimapper.TransactionStatsStatusResponseMapper) func() {
	return func() {
		NewTransactionStatsStatusHandleApi(&transactionStatsStatusHandleDeps{
			client: pb.NewTransactionStatsStatusServiceClient(deps.Client),
			router: deps.E,
			logger: deps.Logger,
			mapper: mapper,
		})
	}
}
