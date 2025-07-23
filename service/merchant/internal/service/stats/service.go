package merchantstatsservice

import (
	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-merchant/internal/redis/stats"
	repository "github.com/MamangRust/monolith-payment-gateway-merchant/internal/repository/stats"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/merchant"
)

type MerchantStatsService interface {
	MerchantStatsAmountService
	MerchantStatsTotalAmountService
	MerchantStatsMethodService
}

type merchantStatsService struct {
	MerchantStatsAmountService
	MerchantStatsTotalAmountService
	MerchantStatsMethodService
}

type DepsStats struct {
	Mencache          mencache.MerchantStatsCache
	ErrorHandler      errorhandler.MerchantStatisticErrorHandler
	Logger            logger.LoggerInterface
	Repository        repository.MerchantStatsRepository
	MapperAmount      responseservice.MerchantAmountResponseMapper
	MapperMethod      responseservice.MerchantPaymentMethodResponseMapper
	MapperTotalAmount responseservice.MerchantTotalAmountResponseMapper
}

func NewMerchantStatsService(deps *DepsStats) MerchantStatsService {
	return &merchantStatsService{
		NewMerchantStatsAmountService(&merchantStatsAmountDeps{
			Cache:        deps.Mencache,
			ErrorHandler: deps.ErrorHandler,
			Repository:   deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MapperAmount,
		}),
		NewMerchantStatsTotalAmountService(&merchantStatsTotalAmountDeps{
			Cache:        deps.Mencache,
			ErrorHandler: deps.ErrorHandler,
			Repository:   deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MapperTotalAmount,
		}),
		NewMerchantStatsMethodService(&merchantStatsMethodDeps{
			Cache:        deps.Mencache,
			ErrorHandler: deps.ErrorHandler,
			Repository:   deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MapperMethod,
		}),
	}
}
