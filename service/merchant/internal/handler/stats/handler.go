package merchantstatshandler

import (
	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/service"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/merchant"
)

type DepsStats struct {
	Service service.Service

	Logger            logger.LoggerInterface
	MapperAmount      protomapper.MerchantStatsAmountProtoMapper
	MapperMethod      protomapper.MerchantStatsMethodProtoMapper
	MapperTotalAmount protomapper.MerchantStatsTotalAmountProtoMapper
}

type HandleStats interface {
	MerchantStatsAmountHandleGrpc
	MerchantStatsMethodHandleGrpc
	MerchantStatsTotalAmountHandleGrpc
}

type handlerStats struct {
	MerchantStatsAmountHandleGrpc
	MerchantStatsMethodHandleGrpc
	MerchantStatsTotalAmountHandleGrpc
}

func NewMerchantStatsHandler(deps *DepsStats) HandleStats {
	return &handlerStats{
		MerchantStatsAmountHandleGrpc:      NewMerchantStatsAmountHandler(deps.Service.MerchantStatsService(), deps.Service.MerchantStatsByMerchantService(), deps.Service.MerchantStatsByApiKeyService(), deps.Logger, deps.MapperAmount),
		MerchantStatsMethodHandleGrpc:      NewMerchantStatsMethodHandler(deps.Service.MerchantStatsService(), deps.Service.MerchantStatsByMerchantService(), deps.Service.MerchantStatsByApiKeyService(), deps.Logger, deps.MapperMethod),
		MerchantStatsTotalAmountHandleGrpc: NewMerchantStatsTotalAmountHandler(deps.Service.MerchantStatsService(), deps.Service.MerchantStatsByMerchantService(), deps.Service.MerchantStatsByApiKeyService(), deps.Logger, deps.MapperTotalAmount),
	}
}
