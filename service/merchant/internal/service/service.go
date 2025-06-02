package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-merchant/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
)

type Service struct {
	MerchantQuery               MerchantQueryService
	MerchantTransaction         MerchantTransactionService
	MerchantCommand             MerchantCommandService
	MerchantStatistic           MerchantStatisticService
	MerchantStatisticByMerchant MerchantStatisticByMerchantService
	MerchantStatisByApiKey      MerchantStatisticByApikeyService
	MerchantDocumentCommand     MerchantDocumentCommandService
	MerchantDocumentQuery       MerchantDocumentQueryService
}

type Deps struct {
	Kafka        kafka.Kafka
	Ctx          context.Context
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
	ErrorHander  errorhandler.ErrorHandler
	Mencache     mencache.Mencache
}

func NewService(deps Deps) *Service {
	merchantMapper := responseservice.NewMerchantResponseMapper()
	merchantDocument := responseservice.NewMerchantDocumentResponseMapper()

	return &Service{
		MerchantQuery:               NewMerchantQueryService(deps.Ctx, deps.Repositories.MerchantQuery, deps.ErrorHander.MerchantQueryError, deps.Mencache.MerchantQueryCache, deps.Logger, merchantMapper),
		MerchantTransaction:         NewMerchantTransactionService(deps.Ctx, deps.ErrorHander.MerchantTransactionError, deps.Repositories.MerchantTrans, deps.Logger, merchantMapper),
		MerchantCommand:             NewMerchantCommandService(deps.Kafka, deps.Ctx, deps.ErrorHander.MerchantCommandError, deps.Mencache.MerchantCommandCache, deps.Repositories.User, deps.Repositories.MerchantQuery, deps.Repositories.MerchantCommand, deps.Logger, merchantMapper),
		MerchantStatistic:           NewMerchantStatisService(deps.Ctx, deps.Mencache.MerchantStatisticCache, deps.ErrorHander.MerchantStatisticError, deps.Repositories.MerchantStat, deps.Logger, merchantMapper),
		MerchantStatisticByMerchant: NewMerchantStatisByMerchantService(deps.Ctx, deps.Mencache.MerchantStatisticByMerchantCache, deps.ErrorHander.MerchantStatisticByMerchantError, deps.Repositories.MerchantStatByMerchant, deps.Logger, merchantMapper),
		MerchantStatisByApiKey:      NewMerchantStatisByApiKeyService(deps.Ctx, deps.Mencache.MerchantStatisticByApiCache, deps.ErrorHander.MerchantStatisticByApiKeyError, deps.Repositories.MerchantStatByApiKey, deps.Logger, merchantMapper),
		MerchantDocumentCommand:     NewMerchantDocumentCommandService(deps.Kafka, deps.Ctx, deps.Mencache.MerchantDocumentCommandCache, deps.ErrorHander.MerchantDocumentCommandError, deps.Repositories.MerchantDocumentCommand, deps.Repositories.MerchantQuery, deps.Repositories.User, deps.Logger, merchantDocument),
		MerchantDocumentQuery:       NewMerchantDocumentQueryService(deps.Ctx, deps.Mencache.MerchantDocumentQueryCache, deps.ErrorHander.MerchantDocumentQueryError, deps.Repositories.MerchantDocumentQuery, deps.Logger, merchantDocument),
	}
}
