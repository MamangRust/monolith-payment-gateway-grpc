package repository

import (
	repositorydashboard "github.com/MamangRust/monolith-payment-gateway-card/internal/repository/dashboard"
	repositorystats "github.com/MamangRust/monolith-payment-gateway-card/internal/repository/stats"
	repositorystatsbycard "github.com/MamangRust/monolith-payment-gateway-card/internal/repository/statsbycard"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/card"
	recordmapperuser "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/user"
)

// Repositories contains all the repositories used in the application
type Repositories struct {
	CardCommand         CardCommandRepository
	CardQuery           CardQueryRepository
	CardDashboard       repositorydashboard.CardDashboardRepository
	CardStatistic       repositorystats.CardStatsRepository
	CardStatisticByCard repositorystatsbycard.CardStatsByCardRepository
	User                UserRepository
}

// NewRepositories creates a new instance of the Repositories struct
func NewRepositories(db *db.Queries) *Repositories {
	mapper := recordmapper.NewCardRecordMapper()
	mapperuser := recordmapperuser.NewUserQueryRecordMapper()

	return &Repositories{
		CardQuery:           NewCardQueryRepository(db, mapper.QueryMapper()),
		CardCommand:         NewCardCommandRepository(db, mapper.CommandMapper()),
		CardDashboard:       repositorydashboard.NewCardDashboardRepository(db),
		CardStatistic:       repositorystats.NewCardStatsRepository(db, mapper.StatsMapper()),
		CardStatisticByCard: repositorystatsbycard.NewCardStatsByCardRepository(db, mapper.StatsByCardMapper()),
		User:                NewUserRepository(db, mapperuser),
	}
}
