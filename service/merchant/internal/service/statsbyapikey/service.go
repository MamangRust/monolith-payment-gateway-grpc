package merchantstatsbyapikeyservice

import (
	mencache "github.com/MamangRust/monolith-payment-gateway-merchant/internal/redis/statsbyapikey"
	repository "github.com/MamangRust/monolith-payment-gateway-merchant/internal/repository/statsbyapikey"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
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
	Mencache      mencache.MerchantStatsByApiKeyCache
	Logger        logger.LoggerInterface
	Repository    repository.MerchantStatsByApiKeyRepository
	Observability observability.TraceLoggerObservability
}

func NewMerchantStatsByApiKeyService(deps *DepsStatsByApiKey) MerchantStatsByApiKeyService {
	return &merchantStatsByApiKeyService{
		NewMerchantStatsAmountByApiKeyService(&merchantStatsAmountByApiKeyDeps{
			Cache:         deps.Mencache,
			Repository:    deps.Repository,
			Logger:        deps.Logger,
			Observability: deps.Observability,
		}),
		NewMerchantStatsTotalAmountByApiKeyService(&merchantStatsTotalAmountByApiKeyDeps{
			Cache:         deps.Mencache,
			Repository:    deps.Repository,
			Logger:        deps.Logger,
			Observability: deps.Observability,
		}),
		NewMerchantStatsMethodByApiKeyService(&merchantStatsMethodByApiKeyDeps{
			Cache:         deps.Mencache,
			Repository:    deps.Repository,
			Logger:        deps.Logger,
			Observability: deps.Observability,
		}),
	}
}
