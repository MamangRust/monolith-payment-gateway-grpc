package topup_test

import (
	card_repo "github.com/MamangRust/monolith-payment-gateway-card/repository"
)

type topupCardRepoAdapter struct {
	card_repo.CardQueryRepository
	card_repo.CardCommandRepository
}
