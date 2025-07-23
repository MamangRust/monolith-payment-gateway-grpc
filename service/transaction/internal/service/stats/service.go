package transactionstatsservice

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/transaction"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-transaction/internal/redis/stats"
	repository "github.com/MamangRust/monolith-payment-gateway-transaction/internal/repository/stats"
)

type TransactionStatsService interface {
	TransactionStatsAmountService
	TransactionStatsStatusService
	TransactionStatsMethodService
}

type transactionStatsService struct {
	TransactionStatsAmountService
	TransactionStatsStatusService
	TransactionStatsMethodService
}

type DepsStats struct {
	Cache mencache.TransactionStatsCache

	ErrorHandler errorhandler.TransactionStatisticErrorHandler

	Repository repository.TransactionStatsRepository

	Logger logger.LoggerInterface

	MapperAmount responseservice.TransactionStatsAmountResponseMapper

	MapperMethod responseservice.TransactionStatsMethodResponseMapper

	MapperStatus responseservice.TransactionStatsStatusResponseMapper
}

func NewTransactionStatsService(deps *DepsStats) TransactionStatsService {
	return &transactionStatsService{
		TransactionStatsAmountService: NewTransactionStatsAmountService(&transactionStatsAmountServiceDeps{
			Cache:        deps.Cache,
			ErrorHandler: deps.ErrorHandler,
			Repository:   deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MapperAmount,
		}),
		TransactionStatsStatusService: NewTransactionStatsStatusService(&transactionStatsStatusServiceDeps{
			Cache:        deps.Cache,
			ErrorHandler: deps.ErrorHandler,
			Repository:   deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MapperStatus,
		}),
		TransactionStatsMethodService: NewTransactionStatsMethodService(&transactionStatsMethodServiceDeps{
			Cache:        deps.Cache,
			ErrorHandler: deps.ErrorHandler,
			Repository:   deps.Repository,
			Logger:       deps.Logger,
			Mapper:       deps.MapperMethod,
		}),
	}
}
