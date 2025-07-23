package topupstatsbycardservice

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/topup"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-topup/internal/redis/statsbycard"
	repository "github.com/MamangRust/monolith-payment-gateway-topup/internal/repository/statsbycard"
)

type TopupStatsByCardService interface {
	TopupStatsByCardAmountService
	TopupStatsByCardMethodService
	TopupStatsByCardStatusService
}

type topupStatsByCardService struct {
	TopupStatsByCardAmountService
	TopupStatsByCardMethodService
	TopupStatsByCardStatusService
}

type DepsStatsByCard struct {
	Cache        mencache.TopupStatsByCardCache
	ErrorHandler errorhandler.TopupStatisticByCardErrorHandler
	Logger       logger.LoggerInterface
	Repository   repository.TopupStatsByCardRepository
	MappeAmount  responseservice.TopupStatsAmountResponseMapper
	MapperMethod responseservice.TopupStatsMethodResponseMapper
	MapperStatus responseservice.TopupStatsStatusResponseMapper
}

func NewTopupStatsByCardService(deps *DepsStatsByCard) TopupStatsByCardService {
	return &topupStatsByCardService{
		TopupStatsByCardAmountService: NewTopupStatsByCardAmountService(&topupStatsByCardAmountDeps{
			Cache:        deps.Cache,
			ErrorHandler: deps.ErrorHandler,
			Repository:   deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MappeAmount,
		}),
		TopupStatsByCardMethodService: NewTopupStatsByCardMethodService(&topupStatsByCardMethodDeps{
			Cache:        deps.Cache,
			ErrorHandler: deps.ErrorHandler,
			Repository:   deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MapperMethod,
		}),
		TopupStatsByCardStatusService: NewTopupStatsByCardStatusService(&topupStatsByCardStatusDeps{
			Cache:        deps.Cache,
			ErrorHandler: deps.ErrorHandler,
			Repository:   deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MapperStatus,
		}),
	}
}
