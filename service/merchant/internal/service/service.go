package service

import (
	"context"

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
}

func NewService(deps Deps) *Service {
	merchantMapper := responseservice.NewMerchantResponseMapper()
	merchantDocument := responseservice.NewMerchantDocumentResponseMapper()

	return &Service{
		MerchantQuery:               NewMerchantQueryService(deps.Ctx, deps.Repositories.MerchantQuery, deps.Logger, merchantMapper),
		MerchantTransaction:         NewMerchantTransactionService(deps.Ctx, deps.Repositories.MerchantTrans, deps.Logger, merchantMapper),
		MerchantCommand:             NewMerchantCommandService(deps.Kafka, deps.Ctx, deps.Repositories.User, deps.Repositories.MerchantQuery, deps.Repositories.MerchantCommand, deps.Logger, merchantMapper),
		MerchantStatistic:           NewMerchantStatisService(deps.Ctx, deps.Repositories.MerchantStat, deps.Logger, merchantMapper),
		MerchantStatisticByMerchant: NewMerchantStatisByMerchantService(deps.Ctx, deps.Repositories.MerchantStatByMerchant, deps.Logger, merchantMapper),
		MerchantStatisByApiKey:      NewMerchantStatisByApiKeyService(deps.Ctx, deps.Repositories.MerchantStatByApiKey, deps.Logger, merchantMapper),
		MerchantDocumentCommand:     NewMerchantDocumentCommandService(deps.Kafka, deps.Ctx, deps.Repositories.MerchantDocumentCommand, deps.Repositories.MerchantQuery, deps.Repositories.User, deps.Logger, merchantDocument),
		MerchantDocumentQuery:       NewMerchantDocumentQueryService(deps.Ctx, deps.Repositories.MerchantDocumentQuery, deps.Logger, merchantDocument),
	}
}
