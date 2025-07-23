package transactionstatsbycardservice

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/transaction"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-transaction/internal/redis/statsbycard"
	repository "github.com/MamangRust/monolith-payment-gateway-transaction/internal/repository/statsbycard"
)

type TransactionStatsByCardService interface {
	TransactionStatsByCardAmountService
	TransactionStatsByCardStatusService
	TransactionStatsByCardMethodService
}

type transactionStatsByCardService struct {
	TransactionStatsByCardAmountService
	TransactionStatsByCardStatusService
	TransactionStatsByCardMethodService
}

type DepsStats struct {
	Cache mencache.TransactionStatsByCardCache

	ErrorHandler errorhandler.TransactionStatisticByCardErrorHandler

	Repository repository.TransactionStatsByCardRepository

	Logger logger.LoggerInterface

	MapperAmount responseservice.TransactionStatsAmountResponseMapper

	MapperMethod responseservice.TransactionStatsMethodResponseMapper

	MapperStatus responseservice.TransactionStatsStatusResponseMapper
}

func NewTransactionStatsByCardService(deps *DepsStats) TransactionStatsByCardService {
	return &transactionStatsByCardService{
		TransactionStatsByCardAmountService: NewTransactionStatsByCardAmountService(&transactionStatsByCardAmountServiceDeps{
			Cache:        deps.Cache,
			ErrorHandler: deps.ErrorHandler,
			Repository:   deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MapperAmount,
		}),
		TransactionStatsByCardStatusService: NewTransactionStatsByCardStatusService(&transactionStatsByCardStatusServiceDeps{
			Cache:        deps.Cache,
			ErrorHandler: deps.ErrorHandler,
			Repository:   deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MapperStatus,
		}),
		TransactionStatsByCardMethodService: NewTransactionStatsByCardMethodService(&transactionStatsByCardMethodServiceDeps{
			Cache:        deps.Cache,
			ErrorHandler: deps.ErrorHandler,
			Repository:   deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MapperMethod,
		}),
	}
}
