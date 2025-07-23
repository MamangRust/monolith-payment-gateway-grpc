package withdrawstatsservice

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/withdraw"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-withdraw/internal/redis/stats"
	repository "github.com/MamangRust/monolith-payment-gateway-withdraw/internal/repository/stats"
)

type WithdrawStatsService interface {
	WithdrawStatsAmountService
	WithdrawStatsStatusService
}

type withdrawStatsStatsService struct {
	WithdrawStatsAmountService
	WithdrawStatsStatusService
}

type DepsStats struct {
	Cache mencache.WithdrawStatsCache

	ErrorHandler errorhandler.WithdrawStatisticErrorHandler

	Logger logger.LoggerInterface

	Repository repository.WithdrawStatsRepository

	MapperAmount responseservice.WithdrawStatsAmountResponseMapper
	MapperStatus responseservice.WithdrawStatsStatusResponseMapper
}

func NewWithdrawStatsService(deps *DepsStats) WithdrawStatsService {
	return &withdrawStatsStatsService{
		WithdrawStatsAmountService: NewWithdrawStatsAmountService(&withdrawStatsAmountDeps{
			ErrorHandler: deps.ErrorHandler,
			Cache:        deps.Cache,
			Repository:   deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MapperAmount,
		}),
		WithdrawStatsStatusService: NewWithdrawStatsStatusService(&withdrawStatsStatusDeps{
			ErrorHandler: deps.ErrorHandler,
			Cache:        deps.Cache,
			Repository:   deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MapperStatus,
		}),
	}
}
