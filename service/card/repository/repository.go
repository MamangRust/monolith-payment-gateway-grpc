package repository

import (
	repositorydashboard "github.com/MamangRust/monolith-payment-gateway-card/repository/dashboard"
	repositorystats "github.com/MamangRust/monolith-payment-gateway-card/repository/stats"
	repositorystatsbycard "github.com/MamangRust/monolith-payment-gateway-card/repository/statsbycard"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
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

func NewRepositories(db *db.Queries) *Repositories {

	return &Repositories{
		CardQuery:           NewCardQueryRepository(db),
		CardCommand:         NewCardCommandRepository(db),
		CardDashboard:       repositorydashboard.NewCardDashboardRepository(db),
		CardStatistic:       repositorystats.NewCardStatsRepository(db),
		CardStatisticByCard: repositorystatsbycard.NewCardStatsByCardRepository(db),
		User:                NewUserRepository(db),
	}
}
