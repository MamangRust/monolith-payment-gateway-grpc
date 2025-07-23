package service

import (
	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/errorhandler"

	mencache "github.com/MamangRust/monolith-payment-gateway-merchant/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/repository"
	merchantstatsservice "github.com/MamangRust/monolith-payment-gateway-merchant/internal/service/stats"
	merchantstatsbyapikeyservice "github.com/MamangRust/monolith-payment-gateway-merchant/internal/service/statsbyapikey"
	merchantstatsbymerchantservice "github.com/MamangRust/monolith-payment-gateway-merchant/internal/service/statsbymerchant"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	mapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/merchant"
	mapperdocument "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/merchantdocument"
)

type Service interface {
	MerchantQueryService() MerchantQueryService
	MerchantTransactionService() MerchantTransactionService
	MerchantCommandService() MerchantCommandService
	MerchantDocumentCommandService() MerchantDocumentCommandService
	MerchantDocumentQueryService() MerchantDocumentQueryService
	MerchantStatsService() merchantstatsservice.MerchantStatsService
	MerchantStatsByMerchantService() merchantstatsbymerchantservice.MerchantStatsByMerchantService
	MerchantStatsByApiKeyService() merchantstatsbyapikeyservice.MerchantStatsByApiKeyService
}

type service struct {
	merchantQuery           MerchantQueryService
	merchantTransaction     MerchantTransactionService
	merchantCommand         MerchantCommandService
	merchantDocumentCommand MerchantDocumentCommandService
	merchantDocumentQuery   MerchantDocumentQueryService
	merchantStats           merchantstatsservice.MerchantStatsService
	merchantStatsByMerchant merchantstatsbymerchantservice.MerchantStatsByMerchantService
	merchantStatsByApiKey   merchantstatsbyapikeyservice.MerchantStatsByApiKeyService
}

// Deps contains the shared dependencies required to initialize the merchant services.
type Deps struct {
	// Kafka is the Kafka client used for producing and consuming events.
	Kafka *kafka.Kafka

	// Repositories contains all data layer interfaces required by the services.
	Repositories repository.Repositories

	// Logger provides structured logging capability.
	Logger logger.LoggerInterface

	// ErrorHander holds centralized error handler implementations for each service.
	ErrorHander *errorhandler.ErrorHandler

	// Mencache provides access to in-memory caches used across merchant services.
	Mencache mencache.Mencache
}

// NewService initializes and returns a new instance of the Service struct,
// which provides a comprehensive suite of merchant-related business logic services.
// It sets up all necessary sub-services, including query, transaction, command,
// and document services, using the provided dependencies and response mappers.
func NewService(deps *Deps) Service {
	merchantMapper := mapper.NewMerchantResponseMapper()
	merchantDocument := mapperdocument.NewMerchantDocumentResponseMapper()

	return &service{
		merchantQuery:       newMerchantQueryService(deps, merchantMapper.QueryMapper()),
		merchantTransaction: newMerchantTransactionService(deps, merchantMapper.TransactionMapper()),
		merchantCommand:     newMerchantCommandService(deps, merchantMapper.CommandMapper()),
		merchantDocumentCommand: newMerchantDocumentCommandService(
			deps, merchantDocument.CommandMapper()),
		merchantDocumentQuery: newMerchantDocumentQueryService(
			deps, merchantDocument.QueryMapper()),
		merchantStats: merchantstatsservice.NewMerchantStatsService(&merchantstatsservice.DepsStats{
			Mencache:          deps.Mencache,
			ErrorHandler:      deps.ErrorHander.MerchantStatisticError,
			Repository:        deps.Repositories,
			Logger:            deps.Logger,
			MapperAmount:      merchantMapper.AmountMapper(),
			MapperTotalAmount: merchantMapper.TotalAmountMapper(),
			MapperMethod:      merchantMapper.MethodMapper(),
		}),
		merchantStatsByMerchant: merchantstatsbymerchantservice.NewMerchantStatsByMerchantService(&merchantstatsbymerchantservice.DepsStatsByMerchant{
			Mencache:          deps.Mencache,
			ErrorHandler:      deps.ErrorHander.MerchantStatisticByMerchantError,
			Repository:        deps.Repositories,
			Logger:            deps.Logger,
			MapperAmount:      merchantMapper.AmountMapper(),
			MapperTotalAmount: merchantMapper.TotalAmountMapper(),
			MapperMethod:      merchantMapper.MethodMapper(),
		}),
		merchantStatsByApiKey: merchantstatsbyapikeyservice.NewMerchantStatsByApiKeyService(&merchantstatsbyapikeyservice.DepsStatsByApiKey{
			Mencache:          deps.Mencache,
			ErrorHandler:      deps.ErrorHander.MerchantStatisticByApiKeyError,
			Repository:        deps.Repositories,
			Logger:            deps.Logger,
			MapperAmount:      merchantMapper.AmountMapper(),
			MapperTotalAmount: merchantMapper.TotalAmountMapper(),
			MapperMethod:      merchantMapper.MethodMapper(),
		}),
	}
}

func (s *service) MerchantQueryService() MerchantQueryService {
	return s.merchantQuery
}
func (s *service) MerchantTransactionService() MerchantTransactionService {
	return s.merchantTransaction
}
func (s *service) MerchantCommandService() MerchantCommandService {
	return s.merchantCommand
}
func (s *service) MerchantDocumentCommandService() MerchantDocumentCommandService {
	return s.merchantDocumentCommand
}
func (s *service) MerchantDocumentQueryService() MerchantDocumentQueryService {
	return s.merchantDocumentQuery
}
func (s *service) MerchantStatsService() merchantstatsservice.MerchantStatsService {
	return s.merchantStats
}
func (s *service) MerchantStatsByMerchantService() merchantstatsbymerchantservice.MerchantStatsByMerchantService {
	return s.merchantStatsByMerchant
}
func (s *service) MerchantStatsByApiKeyService() merchantstatsbyapikeyservice.MerchantStatsByApiKeyService {
	return s.merchantStatsByApiKey
}

// newMerchantQueryService initializes and returns a new instance of
// MerchantQueryService. It sets up necessary components, like context,
// repository, error handler, cache, logger, and response mapper, using
// the provided dependencies to perform merchant query operations.
//
// Parameters:
// - deps: A pointer to Deps containing the shared dependencies used by the merchant services.
// - mapper: A MerchantResponseMapper to map domain models to API-compatible response formats.
//
// Returns:
// - A MerchantQueryService, which is responsible for handling merchant-related query operations.
func newMerchantQueryService(deps *Deps, mapper mapper.MerchantQueryResponseMapper) MerchantQueryService {
	return NewMerchantQueryService(&merchantQueryDeps{
		Repository:   deps.Repositories,
		ErrorHandler: deps.ErrorHander.MerchantQueryError,
		Cache:        deps.Mencache,
		Logger:       deps.Logger,
		Mapper:       mapper,
	})
}

// newMerchantDocumentQueryService initializes and returns a new instance of
// MerchantDocumentQueryService. It sets up the required components, including
// context, cache, error handler, repository, logger, and response mapper,
// using the provided dependencies to handle merchant document query operations.
//
// Parameters:
// - deps: A pointer to Deps containing shared dependencies required by the services.
// - mapper: A MerchantDocumentResponseMapper to map domain models to API-compatible response formats.
//
// Returns:
// - A MerchantDocumentQueryService, responsible for handling merchant document read/query operations.
func newMerchantDocumentQueryService(deps *Deps, mapper mapperdocument.MerchantDocumentQueryResponseMapper) MerchantDocumentQueryService {
	return NewMerchantDocumentQueryService(&merchantDocumentQueryDeps{
		Cache:        deps.Mencache,
		ErrorHandler: deps.ErrorHander.MerchantDocumentQueryError,
		Repository:   deps.Repositories,
		Logger:       deps.Logger,
		Mapper:       mapper,
	})
}

// newMerchantTransactionService initializes and returns a new instance of
// MerchantTransactionService. It sets up the required components, including
// context, error handler, repository, logger, and response mapper, using the
// provided dependencies to perform merchant transaction operations.
//
// Parameters:
// - deps: A pointer to Deps containing shared dependencies required by the services.
// - mapper: A MerchantResponseMapper to map domain models to API-compatible response formats.
//
// Returns:
// - A MerchantTransactionService, responsible for handling merchant transaction operations.
func newMerchantTransactionService(deps *Deps, mapper mapper.MerchantTransactionResponseMapper) MerchantTransactionService {
	return NewMerchantTransactionService(&merchantTransactionDeps{
		ErrorHandler: deps.ErrorHander.MerchantTransactionError,
		Repository:   deps.Repositories,
		Logger:       deps.Logger,
		Mapper:       mapper,
		Cache:        deps.Mencache,
	})
}

// newMerchantCommandService initializes and returns a new instance of
// MerchantCommandService. It sets up the required components, including Kafka
// connection, context, error handler, cache, logger, and response mapper,
// using the provided dependencies to perform merchant command operations.
//
// Parameters:
// - deps: A pointer to Deps containing shared dependencies required by the services.
// - mapper: A MerchantResponseMapper to map domain models to API-compatible response formats.
//
// Returns:
// - A MerchantCommandService, responsible for handling merchant command operations.
func newMerchantCommandService(deps *Deps, mapper mapper.MerchantCommandResponseMapper) MerchantCommandService {
	return NewMerchantCommandService(&merchantCommandServiceDeps{
		Kafka:                     deps.Kafka,
		ErrorHandler:              deps.ErrorHander.MerchantCommandError,
		Cache:                     deps.Mencache,
		UserRepository:            deps.Repositories,
		MerchantQueryRepository:   deps.Repositories,
		MerchantCommandRepository: deps.Repositories,
		Logger:                    deps.Logger,
		Mapper:                    mapper,
	})
}

// newMerchantDocumentCommandService initializes and returns a new instance of
// MerchantDocumentCommandService. It sets up the required components, including
// Kafka connection, context, error handler, cache, logger, and response mapper,
// using the provided dependencies to perform merchant document command operations.
//
// Parameters:
// - deps: A pointer to Deps containing shared dependencies required by the services.
// - mapper: A MerchantDocumentResponseMapper to map domain models to API-compatible response formats.
//
// Returns:
// - A MerchantDocumentCommandService, responsible for handling merchant document command operations.
func newMerchantDocumentCommandService(deps *Deps, mapper mapperdocument.MerchantDocumentCommandResponseMapper) MerchantDocumentCommandService {
	return NewMerchantDocumentCommandService(&merchantDocumentCommandDeps{
		Kafka:                   deps.Kafka,
		Cache:                   deps.Mencache,
		ErrorHandler:            deps.ErrorHander.MerchantDocumentCommandError,
		CommandRepository:       deps.Repositories,
		MerchantQueryRepository: deps.Repositories,
		UserRepository:          deps.Repositories,
		Logger:                  deps.Logger,
		Mapper:                  mapper,
	})
}
