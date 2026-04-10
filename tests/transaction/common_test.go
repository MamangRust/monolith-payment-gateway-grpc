package transaction_test

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	card_repo "github.com/MamangRust/monolith-payment-gateway-card/repository"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

// Wrapper to satisfy transaction repository requirements
type transactionCardRepo struct {
	query   card_repo.CardQueryRepository
	command card_repo.CardCommandRepository
}

func (r *transactionCardRepo) FindCardByUserId(ctx context.Context, user_id int) (*db.GetCardByUserIDRow, error) {
	return r.query.FindCardByUserId(ctx, user_id)
}
func (r *transactionCardRepo) FindUserCardByCardNumber(ctx context.Context, card_number string) (*db.GetUserEmailByCardNumberRow, error) {
	return r.query.FindUserCardByCardNumber(ctx, card_number)
}
func (r *transactionCardRepo) FindCardByCardNumber(ctx context.Context, card_number string) (*db.GetCardByCardNumberRow, error) {
	return r.query.FindCardByCardNumber(ctx, card_number)
}
func (r *transactionCardRepo) UpdateCard(ctx context.Context, request *requests.UpdateCardRequest) (*db.UpdateCardRow, error) {
	return r.command.UpdateCard(ctx, request)
}
