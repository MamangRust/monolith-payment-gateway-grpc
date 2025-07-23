package transferstatsservice

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/transfer"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-transfer/internal/redis/stats"
	repository "github.com/MamangRust/monolith-payment-gateway-transfer/internal/repository/stats"
)

type TransferStatsService interface {
	TransferStatsAmountService
	TransferStatsStatusService
}

type transferStatsService struct {
	TransferStatsAmountService
	TransferStatsStatusService
}

type DepsStats struct {
	Cache        mencache.TransferStatsCache
	ErrorHandler errorhandler.TransferStatisticErrorHandler
	Logger       logger.LoggerInterface
	Repository   repository.TransferStatsRepository
	MappeAmount  responseservice.TransferAmountResponseMapper
	MapperStatus responseservice.TransferStatsStatusResponseMapper
}

func NewTransferStatsService(deps *DepsStats) TransferStatsService {
	return &transferStatsService{
		TransferStatsAmountService: NewTransferStatsAmountService(&transferStatsAmountDeps{
			ErrorHandler: deps.ErrorHandler,
			Cache:        deps.Cache,
			Repository:   deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MappeAmount,
		}),
		TransferStatsStatusService: NewTransferStatsStatusService(&transferStatsStatusDeps{
			ErrorHandler: deps.ErrorHandler,
			Cache:        deps.Cache,
			Repository:   deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MapperStatus,
		}),
	}
}
