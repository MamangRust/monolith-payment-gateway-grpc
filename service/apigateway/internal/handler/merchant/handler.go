package merchanthandler

import (
	pb "github.com/MamangRust/monolith-payment-gateway-pb/merchant"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/merchant"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsMerchant struct {
	Client *grpc.ClientConn
	E      *echo.Echo
	Logger logger.LoggerInterface
}

// RegisterMerchantHandler registers the merchant handler.
//
// This function is responsible for setting up all merchant handlers and their
// corresponding routes.
func RegisterMerchantHandler(deps *DepsMerchant) {
	mapper := apimapper.NewMerchantResponseMapper()

	handlers := []func(){
		setupMerchantQueryHandler(deps, mapper.QueryMapper()),
		setupMerchantCommandHandler(deps, mapper.CommandMapper()),
		setupMerchantStatsAmountHandler(deps, mapper.AmountStatsMapper()),
		setupMerchantStatsMethodHandler(deps, mapper.MethodStatsMapper()),
		setupMerchantStatsTotalAmountHandler(deps, mapper.TotalAmountStatsMapper()),
		setupMerchantTransactionHandler(deps, mapper.TransactionMapper()),
	}

	for _, h := range handlers {
		h()
	}
}

// setupMerchantQueryHandler sets up the merchant query handler and its route.
//
// The handler will be registered with the given echo router and will use the
// given client, logger, and mapper to handle incoming requests.
//
// The returned function is a setup function that can be used to register the
// handler with the given router.
func setupMerchantQueryHandler(deps *DepsMerchant, mapper apimapper.MerchantQueryResponseMapper) func() {
	return func() {
		NewMerchantQueryHandleApi(&merchantQueryHandleDeps{
			client: pb.NewMerchantQueryServiceClient(deps.Client),
			router: deps.E,
			logger: deps.Logger,
			mapper: mapper,
		})
	}
}

// setupMerchantCommandHandler sets up the merchant command handler and its route.
//
// This handler is responsible for processing merchant command operations such
// as creation, updating, and deletion of merchant entities. It utilizes the
// provided dependencies to initialize the handler and register the routes
// with the given Echo router.
//
// Parameters:
//   - deps: A pointer to DepsMerchant, which contains shared dependencies such as
//     a gRPC client connection, an Echo router, and a logger interface.
//   - mapper: A MerchantCommandResponseMapper that translates domain models into
//     API-compatible response formats.
//
// Returns:
//   - A function that, when executed, initializes the merchant command handler
//     and registers its routes with the Echo router.
func setupMerchantCommandHandler(deps *DepsMerchant, mapper apimapper.MerchantCommandResponseMapper) func() {
	return func() {
		NewMerchantCommandHandleApi(&merchantCommandHandleDeps{
			client: pb.NewMerchantCommandServiceClient(deps.Client),
			router: deps.E,
			logger: deps.Logger,
			mapper: mapper,
		})
	}
}

// setupMerchantStatsAmountHandler sets up the merchant stats amount handler and its route.
//
// This handler is responsible for processing merchant statistics related to transaction amounts,
// such as monthly and yearly summaries, grouped by various criteria (e.g., merchant or API key).
// It utilizes the provided dependencies to initialize the handler and register the routes
// with the given Echo router.
//
// Parameters:
//   - deps: A pointer to DepsMerchant, which contains shared dependencies such as
//     a gRPC client connection, an Echo router, and a logger interface.
//   - mapper: A MerchantStatsAmountResponseMapper that translates domain models into
//     API-compatible response formats.
//
// Returns:
//   - A function that, when executed, initializes the merchant stats amount handler
//     and registers its routes with the Echo router.
func setupMerchantStatsAmountHandler(deps *DepsMerchant, mapper apimapper.MerchantStatsAmountResponseMapper) func() {
	return func() {
		NewMerchantStatsAmountHandleApi(&merchantStatsAmountHandleDeps{
			client: pb.NewMerchantStatsAmountServiceClient(deps.Client),
			router: deps.E,
			logger: deps.Logger,
			mapper: mapper,
		})
	}
}

// setupMerchantStatsMethodHandler sets up the merchant stats method handler and its route.
//
// This handler manages statistics related to the usage of different payment methods by merchants.
// It uses the provided dependencies to create the handler and register the relevant endpoints
// under the Echo router.
//
// Parameters:
//   - deps: A pointer to DepsMerchant, which contains shared dependencies such as
//     a gRPC client connection, an Echo router, and a logger interface.
//   - mapper: A MerchantStatsMethodResponseMapper that translates domain models into
//     API-compatible response formats.
//
// Returns:
//   - A function that, when executed, initializes the merchant stats method handler
//     and registers its routes with the Echo router.
func setupMerchantStatsMethodHandler(deps *DepsMerchant, mapper apimapper.MerchantStatsMethodResponseMapper) func() {
	return func() {
		NewMerchantStatsMethodHandleApi(&merchantStatsMethodHandleDeps{
			client: pb.NewMerchantStatsMethodServiceClient(deps.Client),
			router: deps.E,
			logger: deps.Logger,
			mapper: mapper,
		})
	}
}

// setupMerchantStatsTotalAmountHandler sets up the merchant stats total amount handler and its route.
//
// This handler is responsible for providing aggregated total amount statistics for merchant transactions.
// It uses the given dependencies to initialize the handler and bind the HTTP routes to the Echo router.
//
// Parameters:
//   - deps: A pointer to DepsMerchant, which contains shared dependencies such as
//     a gRPC client connection, an Echo router, and a logger interface.
//   - mapper: A MerchantStatsTotalAmountResponseMapper that maps internal data structures to API responses.
//
// Returns:
//   - A function that, when executed, initializes the merchant stats total amount handler
//     and registers its routes with the Echo router.
func setupMerchantStatsTotalAmountHandler(deps *DepsMerchant, mapper apimapper.MerchantStatsTotalAmountResponseMapper) func() {
	return func() {
		NewMerchantStatsTotalAmountHandleApi(&merchantStatsTotalAmountHandleDeps{
			client: pb.NewMerchantStatsTotalAmountServiceClient(deps.Client),
			router: deps.E,
			logger: deps.Logger,
			mapper: mapper,
		})
	}
}

// setupMerchantTransactionHandler sets up the merchant transaction handler and its route.
//
// This handler processes operations related to merchant transaction retrieval,
// providing access to transaction records based on various filters.
// The handler is initialized using the provided dependencies and registered with the Echo router.
//
// Parameters:
//   - deps: A pointer to DepsMerchant, which contains shared dependencies such as
//     a gRPC client connection, an Echo router, and a logger interface.
//   - mapper: A MerchantTransactionResponseMapper that converts internal models into
//     response structures compatible with the API.
//
// Returns:
//   - A function that, when executed, initializes the merchant transaction handler
//     and registers its routes with the Echo router.
func setupMerchantTransactionHandler(deps *DepsMerchant, mapper apimapper.MerchantTransactionResponseMapper) func() {
	return func() {
		NewMerchantTransactionHandleApi(&merchantTransactionHandleDeps{
			client: pb.NewMerchantTransactionServiceClient(deps.Client),
			router: deps.E,
			logger: deps.Logger,
			mapper: mapper,
		})
	}
}
