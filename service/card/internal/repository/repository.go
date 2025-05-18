package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type Repositories struct {
	CardCommand         CardCommandRepository
	CardQuery           CardQueryRepository
	CardDashboard       CardDashboardRepository
	CardStatistic       CardStatisticRepository
	CardStatisticByCard CardStatisticByCardRepository

	User UserRepository
}

type Deps struct {
	DB           *db.Queries
	Ctx          context.Context
	MapperRecord *recordmapper.RecordMapper
}

func NewRepositories(deps Deps) *Repositories {
	return &Repositories{
		CardCommand:         NewCardCommandRepository(deps.DB, deps.Ctx, deps.MapperRecord.CardRecordMapper),
		CardQuery:           NewCardQueryRepository(deps.DB, deps.Ctx, deps.MapperRecord.CardRecordMapper),
		CardDashboard:       NewCardDashboardRepository(deps.DB, deps.Ctx, deps.MapperRecord.CardRecordMapper),
		CardStatistic:       NewCardStatisticRepository(deps.DB, deps.Ctx, deps.MapperRecord.CardRecordMapper),
		CardStatisticByCard: NewCardStatisticByCardRepository(deps.DB, deps.Ctx, deps.MapperRecord.CardRecordMapper),
		User:                NewUserRepository(deps.DB, deps.Ctx, deps.MapperRecord.UserRecordMapper),
	}
}
