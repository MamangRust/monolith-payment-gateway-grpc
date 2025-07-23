package merchantstatsbyapikeyservice

import (
	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-merchant/internal/redis/statsbyapikey"
	repository "github.com/MamangRust/monolith-payment-gateway-merchant/internal/repository/statsbyapikey"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/merchant"
)

type MerchantStatsByApiKeyService interface {
	MerchantStatsByApiKeyAmountService
	MerchantStatsByApiKeyTotalAmountService
	MerchantStatsByApiKeyMethodService
}

type merchantStatsByApiKeyService struct {
	MerchantStatsByApiKeyAmountService
	MerchantStatsByApiKeyTotalAmountService
	MerchantStatsByApiKeyMethodService
}

type DepsStatsByApiKey struct {
	Mencache          mencache.MerchantStatsByApiKeyCache
	ErrorHandler      errorhandler.MerchantStatisticByApikeyErrorHandler
	Logger            logger.LoggerInterface
	Repository        repository.MerchantStatsByApiKeyRepository
	MapperAmount      responseservice.MerchantAmountResponseMapper
	MapperMethod      responseservice.MerchantPaymentMethodResponseMapper
	MapperTotalAmount responseservice.MerchantTotalAmountResponseMapper
}

func NewMerchantStatsByApiKeyService(deps *DepsStatsByApiKey) MerchantStatsByApiKeyService {
	return &merchantStatsByApiKeyService{
		NewMerchantStatsAmountByApiKeyService(&merchantStatsAmountByApiKeyDeps{
			Cache:        deps.Mencache,
			ErrorHandler: deps.ErrorHandler,
			Repository:   deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MapperAmount,
		}),
		NewMerchantStatsTotalAmountByApiKeyService(&merchantStatsTotalAmountByApiKeyDeps{
			Cache:        deps.Mencache,
			ErrorHandler: deps.ErrorHandler,
			Repository:   deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MapperTotalAmount,
		}),
		NewMerchantStatsMethodByApiKeyService(&merchantStatsMethodByApiKeyDeps{
			Cache:        deps.Mencache,
			ErrorHandler: deps.ErrorHandler,
			Repository:   deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MapperMethod,
		}),
	}
}
