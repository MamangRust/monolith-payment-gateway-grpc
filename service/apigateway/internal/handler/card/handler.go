package cardhandler

import (
	pb "github.com/MamangRust/monolith-payment-gateway-pb/card"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/card"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsCard struct {
	Client *grpc.ClientConn

	E *echo.Echo

	Logger logger.LoggerInterface
}

// NewCardHandler initializes handlers for various card-related operations.
//
// This function sets up multiple handlers for card operations including query,
// command, dashboard, and statistics (balance, transaction, top-up, withdrawal, transfer).
// It takes a DepsCard struct which contains the necessary dependencies such as
// gRPC client connection, Echo router, and logger. Each handler is initialized
// with the corresponding response mapper and added to a slice of handler functions,
// which are executed sequentially to set up the routes.
func RegisterCardHandler(deps *DepsCard) {
	mapper := apimapper.NewCardResponseMapper()

	handlers := []func(){
		setupCardQueryHandler(deps, mapper.QueryMapper()),
		setupCardCommandHandler(deps, mapper.CommandMapper()),
		setupCardDashboardHandler(deps, mapper.DashboardMapper()),
		setupCardStatsBalanceHandler(deps, mapper.BalanceStatsMapper()),
		setupCardStatsTransactionHandler(deps, mapper.AmountStatsMapper()),
		setupCardStatsTopupHandler(deps, mapper.AmountStatsMapper()),
		setupCardStatsWithdrawHandler(deps, mapper.AmountStatsMapper()),
		setupCardStatsTransferHandler(deps, mapper.AmountStatsMapper()),
	}

	for _, h := range handlers {
		h()
	}
}

// setupCardQueryHandler sets up the handler for the card query service.
//
// It creates a new instance of the CardQueryHandleApi and registers the
// handler with the Echo router. It takes a pointer to DepsCard and a
// mapper for CardResponse. It returns a function that can be executed to
// set up the handler.
func setupCardQueryHandler(deps *DepsCard, mapper apimapper.CardQueryResponseMapper) func() {
	return func() {
		NewCardQueryHandleApi(&cardQueryHandleApiDeps{
			client: pb.NewCardQueryServiceClient(deps.Client),
			router: deps.E,
			logger: deps.Logger,
			mapper: mapper,
		})
	}
}

// setupCardCommandHandler sets up the handler for the card command service.
//
// It creates a new instance of the CardCommandHandleApi and registers the
// handler with the Echo router. It takes a pointer to DepsCard and a
// mapper for CardResponse. It returns a function that can be executed to
// set up the handler.
func setupCardCommandHandler(deps *DepsCard, mapper apimapper.CardCommandResponseMapper) func() {
	return func() {
		NewCardCommandHandleApi(&cardCommandHandleApiDeps{
			client: pb.NewCardCommandServiceClient(deps.Client),
			router: deps.E,
			logger: deps.Logger,
			mapper: mapper,
		})
	}
}

// setupCardDashboardHandler sets up the handler for the card dashboard service.
//
// It creates a new instance of the CardDashboardHandleApi and registers the
// handler with the Echo router. It takes a pointer to DepsCard and a
// mapper for CardResponse. It returns a function that can be executed to
// set up the handler.
func setupCardDashboardHandler(deps *DepsCard, mapper apimapper.CardDashboardResponseMapper) func() {
	return func() {
		NewCardDashboardHandleApi(&cardDashboardHandleApiDeps{
			client: pb.NewCardDashboardServiceClient(deps.Client),
			router: deps.E,
			logger: deps.Logger,
			mapper: mapper,
		})
	}
}

// setupCardStatsBalanceHandler sets up the handler for the card statistics balance service.
//
// It creates a new instance of the CardStatsBalanceHandleApi and registers the
// handler with the Echo router. It takes a pointer to DepsCard and a
// mapper for CardStatsBalanceResponse. It returns a function that can be executed to
// set up the handler.
func setupCardStatsBalanceHandler(deps *DepsCard, mapper apimapper.CardStatsBalanceResponseMapper) func() {
	return func() {
		NewCardStatsBalanceHandleApi(&cardStatsBalanceHandleApiDeps{
			client: pb.NewCardStatsBalanceServiceClient(deps.Client),
			router: deps.E,
			logger: deps.Logger,
			mapper: mapper,
		})
	}
}

// setupCardStatsTopupHandler sets up the handler for the card statistics top-up service.
//
// It creates a new instance of the CardStatsTopupHandleApi and registers the
// handler with the Echo router. It takes a pointer to DepsCard and a
// mapper for CardStatsAmountResponse. It returns a function that can be executed to
// set up the handler.
func setupCardStatsTopupHandler(deps *DepsCard, mapper apimapper.CardStatsAmountResponseMapper) func() {
	return func() {
		NewCardStatsTopupHandleApi(&cardStatsTopupHandleApiDeps{
			client: pb.NewCardStatsTopupServiceClient(deps.Client),
			router: deps.E,
			logger: deps.Logger,
			mapper: mapper,
		})
	}
}

// setupCardStatsTransactionHandler sets up the handler for the card statistics transaction service.
//
// It creates a new instance of the CardStatsTransactionHandleApi and registers the
// handler with the Echo router. It takes a pointer to DepsCard and a
// mapper for CardStatsAmountResponse. It returns a function that can be executed to
// set up the handler.
func setupCardStatsTransactionHandler(deps *DepsCard, mapper apimapper.CardStatsAmountResponseMapper) func() {
	return func() {
		NewCardStatsTransactionHandleApi(&cardStatsTransactionHandleApiDeps{
			client: pb.NewCardStatsTransactonServiceClient(deps.Client),
			router: deps.E,
			logger: deps.Logger,
			mapper: mapper,
		})
	}
}

// setupCardStatsTransferHandler sets up the handler for the card statistics transfer service.
//
// It creates a new instance of the CardStatsTransferHandleApi and registers the
// handler with the Echo router. It takes a pointer to DepsCard and a
// mapper for CardStatsAmountResponse. It returns a function that can be executed to
// set up the handler.
func setupCardStatsTransferHandler(deps *DepsCard, mapper apimapper.CardStatsAmountResponseMapper) func() {
	return func() {
		NewCardStatsTransferHandleApi(&cardStatsTransferHandleApiDeps{
			client: pb.NewCardStatsTransferServiceClient(deps.Client),
			router: deps.E,
			logger: deps.Logger,
			mapper: mapper,
		})
	}
}

// setupCardStatsWithdrawHandler sets up the handler for the card statistics withdraw service.
//
// It creates a new instance of the CardStatsWithdrawHandleApi and registers the
// handler with the Echo router. It takes a pointer to DepsCard and a
// mapper for CardStatsAmountResponse. It returns a function that can be executed to
// set up the handler.
func setupCardStatsWithdrawHandler(deps *DepsCard, mapper apimapper.CardStatsAmountResponseMapper) func() {
	return func() {
		NewCardStatsWithdrawHandleApi(&cardStatsWithdrawHandleApiDeps{
			client: pb.NewCardStatsWithdrawServiceClient(deps.Client),
			router: deps.E,
			logger: deps.Logger,
			mapper: mapper,
		})
	}
}
