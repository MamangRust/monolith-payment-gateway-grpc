package merchantdocumenthandler

import (
	pb "github.com/MamangRust/monolith-payment-gateway-pb/merchantdocument"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/merchantdocument"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsMerchantDocument struct {
	Client *grpc.ClientConn

	E *echo.Echo

	Logger logger.LoggerInterface
}

func RegisterMerchantDocumentHandler(deps *DepsMerchantDocument) {
	mapper := apimapper.NewMerchantDocumentResponseMapper()

	handlers := []func(){
		setupMerchantDocumentQueryHandler(deps, mapper.QueryMapper()),
		setupMerchantDocumentCommandHandler(deps, mapper.CommandMapper()),
	}

	for _, h := range handlers {
		h()
	}
}

func setupMerchantDocumentQueryHandler(deps *DepsMerchantDocument, mapper apimapper.MerchantDocumentQueryResponseMapper) func() {
	return func() {
		NewMerchantQueryDocumentHandler(&merchantDocumentQueryDocumentHandleDeps{
			client: pb.NewMerchantDocumentServiceClient(deps.Client),
			router: deps.E,
			logger: deps.Logger,
			mapper: mapper,
		})
	}
}

func setupMerchantDocumentCommandHandler(deps *DepsMerchantDocument, mapper apimapper.MerchantDocumentCommandResponseMapper) func() {
	return func() {
		NewMerchantCommandDocumentHandler(&merchantCommandDocumentHandleDeps{
			client: pb.NewMerchantDocumentCommandServiceClient(deps.Client),
			router: deps.E,
			logger: deps.Logger,
			mapper: mapper,
		})
	}
}
