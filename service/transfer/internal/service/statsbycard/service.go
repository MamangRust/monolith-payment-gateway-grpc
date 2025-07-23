package transferstatsbycardservice

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/transfer"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-transfer/internal/redis/statsbycard"
	repository "github.com/MamangRust/monolith-payment-gateway-transfer/internal/repository/statsbycard"
)

type TransferStatsByCardService interface {
	TransferStatsByCardAmountService
	TransferStatsByCardStatusService
}

type transferStatsByCardService struct {
	TransferStatsByCardAmountService
	TransferStatsByCardStatusService
}

type DepsStats struct {
	Cache        mencache.TransferStatsByCardCache
	ErrorHandler errorhandler.TransferStatisticByCardErrorHandler
	Logger       logger.LoggerInterface
	Repository   repository.TransferStatsByCardRepository
	MappeAmount  responseservice.TransferAmountResponseMapper
	MapperStatus responseservice.TransferStatsStatusResponseMapper
}

func NewTransferStatsByCardService(deps *DepsStats) TransferStatsByCardService {
	return &transferStatsByCardService{
		TransferStatsByCardAmountService: NewTransferStatsByCardAmountService(&transferStatsByCardAmountDeps{
			ErrorHandler: deps.ErrorHandler,
			Cache:        deps.Cache,
			Sender:       deps.Repository,
			Receiver:     deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MappeAmount,
		}),
		TransferStatsByCardStatusService: NewTransferStatsByCardStatusService(&transferStatsByCardStatusDeps{
			ErrorHandler: deps.ErrorHandler,
			Cache:        deps.Cache,
			Repository:   deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MapperStatus,
		}),
	}
}
