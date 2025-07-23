package repository

import (
	"context"

	merchantstatsrepository "github.com/MamangRust/monolith-payment-gateway-merchant/internal/repository/stats"
	merchantstatsapikeyrepository "github.com/MamangRust/monolith-payment-gateway-merchant/internal/repository/statsbyapikey"
	merchantstatsmerchantrepository "github.com/MamangRust/monolith-payment-gateway-merchant/internal/repository/statsbymerchant"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	mapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/merchant"
	mapperdocument "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/merchantdocument"
	mapperuser "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/user"
)

type Repositories interface {
	MerchantQueryRepository
	MerchantCommandRepository
	MerchantDocumentQueryRepository
	MerchantDocumentCommandRepository
	MerchantTransactionRepository
	merchantstatsrepository.MerchantStatsRepository
	merchantstatsapikeyrepository.MerchantStatsByApiKeyRepository
	merchantstatsmerchantrepository.MerchantStatsByMerchantRepository
	UserRepository
}

type repositories struct {
	MerchantQueryRepository
	MerchantCommandRepository
	MerchantDocumentQueryRepository
	MerchantDocumentCommandRepository
	MerchantTransactionRepository
	merchantstatsrepository.MerchantStatsRepository
	merchantstatsapikeyrepository.MerchantStatsByApiKeyRepository
	merchantstatsmerchantrepository.MerchantStatsByMerchantRepository
	UserRepository
}

// Deps is a struct that contains all the dependencies needed to create the repositories
type Deps struct {
	DB  *db.Queries
	Ctx context.Context
}

// NewRepositories creates a new instance of the Repositories interface, which contains all the repositories used in the merchant service.
func NewRepositories(db *db.Queries) Repositories {
	mapper := mapper.NewMerchantRecordMapper()
	mapperdocument := mapperdocument.NewMerchantDocumentRecordMapper()
	mapperuser := mapperuser.NewUserQueryRecordMapper()

	return &repositories{
		NewMerchantQueryRepository(db, mapper.QueryMapper()),
		NewMerchantCommandRepository(db, mapper.CommandMapper()),
		NewMerchantDocumentQueryRepository(db, mapperdocument.QueryMapper()),
		NewMerchantDocumentCommandRepository(db, mapperdocument.CommandMapper()),
		NewMerchantTransactionRepository(db, mapper.TransactionMapper()),
		merchantstatsrepository.NewMerchantStatsRepository(db, mapper.StatisticMapper()),
		merchantstatsapikeyrepository.NewMerchantStatsByApiKeyRepository(db, mapper.ByApiKeyMapper()),
		merchantstatsmerchantrepository.NewMerchantStatsByMerchantRepository(db, mapper.ByMerchantMapper()),
		NewUserRepository(db, mapperuser),
	}
}
