package topupstatsservice

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/topup"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-topup/internal/redis/stats"
	repository "github.com/MamangRust/monolith-payment-gateway-topup/internal/repository/stats"
)

type TopupStatsService interface {
	TopupStatsAmountService
	TopupStatsMethodService
	TopupStatsStatusService
}

type topupStatsService struct {
	TopupStatsAmountService
	TopupStatsMethodService
	TopupStatsStatusService
}

type DepsStats struct {
	Cache        mencache.TopupStatsCache
	ErrorHandler errorhandler.TopupStatisticErrorHandler
	Logger       logger.LoggerInterface
	Repository   repository.TopupStatsRepository
	MappeAmount  responseservice.TopupStatsAmountResponseMapper
	MapperMethod responseservice.TopupStatsMethodResponseMapper
	MapperStatus responseservice.TopupStatsStatusResponseMapper
}

func NewTopupStatsService(deps *DepsStats) TopupStatsService {
	return &topupStatsService{
		TopupStatsAmountService: NewTopupStatsAmountService(&topupStatsAmountDeps{
			Cache:        deps.Cache,
			ErrorHandler: deps.ErrorHandler,
			Repository:   deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MappeAmount,
		}),
		TopupStatsMethodService: NewTopupStatsMethodService(&topupStatsMethodDeps{
			Cache:        deps.Cache,
			ErrorHandler: deps.ErrorHandler,
			Repository:   deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MapperMethod,
		}),
		TopupStatsStatusService: NewTopupStatsStatusService(&topupStatsStatusDeps{
			Cache:        deps.Cache,
			ErrorHandler: deps.ErrorHandler,
			Repository:   deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MapperStatus,
		}),
	}
}
