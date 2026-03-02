package repository

import (
	merchantstatsrepository "github.com/MamangRust/monolith-payment-gateway-merchant/internal/repository/stats"
	merchantstatsapikeyrepository "github.com/MamangRust/monolith-payment-gateway-merchant/internal/repository/statsbyapikey"
	merchantstatsmerchantrepository "github.com/MamangRust/monolith-payment-gateway-merchant/internal/repository/statsbymerchant"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
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

func NewRepositories(db *db.Queries) Repositories {
	return &repositories{
		NewMerchantQueryRepository(db),
		NewMerchantCommandRepository(db),
		NewMerchantDocumentQueryRepository(db),
		NewMerchantDocumentCommandRepository(db),
		NewMerchantTransactionRepository(db),
		merchantstatsrepository.NewMerchantStatsRepository(db),
		merchantstatsapikeyrepository.NewMerchantStatsByApiKeyRepository(db),
		merchantstatsmerchantrepository.NewMerchantStatsByMerchantRepository(db),
		NewUserRepository(db),
	}
}
