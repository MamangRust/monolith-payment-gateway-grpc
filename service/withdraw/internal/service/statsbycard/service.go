package withdrawstatsbycardservice

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/withdraw"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-withdraw/internal/redis/statsbycard"
	repository "github.com/MamangRust/monolith-payment-gateway-withdraw/internal/repository/statsbycard"
)

type WithdrawStatsByCardService interface {
	WithdrawStatsByCardAmountService
	WithdrawStatsByCardStatusService
}

type withdrawStatsByCardStatsByCardService struct {
	WithdrawStatsByCardAmountService
	WithdrawStatsByCardStatusService
}

type DepsStatsByCard struct {
	Cache mencache.WithdrawStatsByCardCache

	ErrorHandler errorhandler.WithdrawStatisticByCardErrorHandler

	Logger logger.LoggerInterface

	Repository repository.WithdrawStatsByCardRepository

	MapperAmount responseservice.WithdrawStatsAmountResponseMapper
	MapperStatus responseservice.WithdrawStatsStatusResponseMapper
}

func NewWithdrawStatsByCardService(deps *DepsStatsByCard) WithdrawStatsByCardService {
	return &withdrawStatsByCardStatsByCardService{
		WithdrawStatsByCardAmountService: NewWithdrawStatsByCardAmountService(&withdrawStatsByCardAmountDeps{
			ErrorHandler: deps.ErrorHandler,
			Cache:        deps.Cache,
			Repository:   deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MapperAmount,
		}),
		WithdrawStatsByCardStatusService: NewWithdrawStatsByCardStatusService(&withdrawStatsByCardStatusDeps{
			ErrorHandler: deps.ErrorHandler,
			Cache:        deps.Cache,
			Repository:   deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MapperStatus,
		}),
	}
}
