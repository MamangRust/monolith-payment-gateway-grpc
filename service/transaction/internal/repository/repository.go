package repository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	mappercard "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/card"
	mappermerchant "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/merchant"
	mappersaldo "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/saldo"
	mapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/transaction"
	transactionstatsrepository "github.com/MamangRust/monolith-payment-gateway-transaction/internal/repository/stats"
	transactionbycardrepository "github.com/MamangRust/monolith-payment-gateway-transaction/internal/repository/statsbycard"
)

type Repositories interface {
	SaldoRepository
	MerchantRepository
	CardRepository
	TransactionQueryRepository
	TransactionCommandRepository
	transactionstatsrepository.TransactionStatsRepository
	transactionbycardrepository.TransactionStatsByCardRepository
}

// Repositories is a struct that contains all the repositories used in the transaction service.
type repositories struct {
	SaldoRepository
	MerchantRepository
	CardRepository
	TransactionQueryRepository
	TransactionCommandRepository
	transactionstatsrepository.TransactionStatsRepository
	transactionbycardrepository.TransactionStatsByCardRepository
}

// NewRepositories creates a new instance of Repositories with the provided database
// queries, context, and record mappers. This repository is responsible for
// executing command and query operations related to topup records in the database.
//
// Parameters:
//   - deps: A pointer to Deps containing the required dependencies.
//
// Returns:
//   - A pointer to the newly created Repositories instance.
func NewRepositories(db *db.Queries) Repositories {
	mapper := mapper.NewTransactionRecordMapper()
	mappersaldo := mappersaldo.NewSaldoQueryRecordMapper()
	mappercard := mappercard.NewCardQueryRecordMapper()
	mappermerchant := mappermerchant.NewMerchantQueryRecordMapper()

	return &repositories{
		SaldoRepository:                  NewSaldoRepository(db, mappersaldo),
		MerchantRepository:               NewMerchantRepository(db, mappermerchant),
		CardRepository:                   NewCardRepository(db, mappercard),
		TransactionQueryRepository:       NewTransactionQueryRepository(db, mapper.QueryMapper()),
		TransactionCommandRepository:     NewTransactionCommandRepository(db, mapper.CommandMapper()),
		TransactionStatsRepository:       transactionstatsrepository.NewTransactionStatsRepository(db, mapper.StatsMapper()),
		TransactionStatsByCardRepository: transactionbycardrepository.NewTransactionStatsRepository(db, mapper.StatsByCardMapper()),
	}
}
