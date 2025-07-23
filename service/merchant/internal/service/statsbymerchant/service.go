package merchantstatsbymerchantservice

import (
	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-merchant/internal/redis/statsbymerchant"
	repository "github.com/MamangRust/monolith-payment-gateway-merchant/internal/repository/statsbymerchant"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/merchant"
)

type MerchantStatsByMerchantService interface {
	MerchantStatsByMerchantAmountService
	MerchantStatsByMerchantTotalAmountService
	MerchantStatsByMerchantMethodService
}

type merchantStatsByMerchantService struct {
	MerchantStatsByMerchantAmountService
	MerchantStatsByMerchantTotalAmountService
	MerchantStatsByMerchantMethodService
}

type DepsStatsByMerchant struct {
	Mencache          mencache.MerchantStatsByMerchantCache
	ErrorHandler      errorhandler.MerchantStatisticByMerchantErrorHandler
	Logger            logger.LoggerInterface
	Repository        repository.MerchantStatsByMerchantRepository
	MapperAmount      responseservice.MerchantAmountResponseMapper
	MapperMethod      responseservice.MerchantPaymentMethodResponseMapper
	MapperTotalAmount responseservice.MerchantTotalAmountResponseMapper
}

func NewMerchantStatsByMerchantService(deps *DepsStatsByMerchant) MerchantStatsByMerchantService {
	return &merchantStatsByMerchantService{
		NewMerchantStatsByMerchantAmountService(&merchantStatsByMerchantAmountServiceDeps{
			Cache:        deps.Mencache,
			ErrorHandler: deps.ErrorHandler,
			Repository:   deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MapperAmount,
		}),
		NewMerchantStatsByMerchantTotalAmountService(&merchantStatsByMerchantTotalAmountDeps{
			Cache:        deps.Mencache,
			ErrorHandler: deps.ErrorHandler,
			Repository:   deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MapperTotalAmount,
		}),
		NewMerchantStatsByMerchantMethodService(&merchantStatsByMerchantMethodDeps{
			Cache:        deps.Mencache,
			ErrorHandler: deps.ErrorHandler,
			Repository:   deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MapperMethod,
		}),
	}
}
