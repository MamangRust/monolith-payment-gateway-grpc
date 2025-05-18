package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type cardDashboardRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.CardRecordMapping
}

func NewCardDashboardRepository(db *db.Queries, ctx context.Context, mapping recordmapper.CardRecordMapping) *cardDashboardRepository {
	return &cardDashboardRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *cardDashboardRepository) GetTotalBalances() (*int64, error) {
	res, err := r.db.GetTotalBalance(r.ctx)

	if err != nil {
		return nil, card_errors.ErrGetTotalBalancesFailed
	}

	return &res, nil
}

func (r *cardDashboardRepository) GetTotalTopAmount() (*int64, error) {
	res, err := r.db.GetTotalTopupAmount(r.ctx)

	if err != nil {
		return nil, card_errors.ErrGetTotalTopAmountFailed
	}

	return &res, nil
}

func (r *cardDashboardRepository) GetTotalWithdrawAmount() (*int64, error) {
	res, err := r.db.GetTotalWithdrawAmount(r.ctx)

	if err != nil {
		return nil, card_errors.ErrGetTotalWithdrawAmountFailed
	}

	return &res, nil
}

func (r *cardDashboardRepository) GetTotalTransactionAmount() (*int64, error) {
	res, err := r.db.GetTotalTransactionAmount(r.ctx)

	if err != nil {
		return nil, card_errors.ErrGetTotalTransactionAmountFailed
	}

	return &res, nil
}

func (r *cardDashboardRepository) GetTotalTransferAmount() (*int64, error) {
	res, err := r.db.GetTotalTransferAmount(r.ctx)

	if err != nil {
		return nil, card_errors.ErrGetTotalTransferAmountFailed
	}

	return &res, nil
}

func (r *cardDashboardRepository) GetTotalBalanceByCardNumber(cardNumber string) (*int64, error) {
	res, err := r.db.GetTotalBalanceByCardNumber(r.ctx, cardNumber)

	if err != nil {
		return nil, card_errors.ErrGetTotalBalanceByCardFailed
	}

	return &res, nil
}

func (r *cardDashboardRepository) GetTotalTopupAmountByCardNumber(cardNumber string) (*int64, error) {
	res, err := r.db.GetTotalTopupAmountByCardNumber(r.ctx, cardNumber)

	if err != nil {
		return nil, card_errors.ErrGetTotalTopupAmountByCardFailed
	}

	return &res, nil
}

func (r *cardDashboardRepository) GetTotalWithdrawAmountByCardNumber(cardNumber string) (*int64, error) {
	res, err := r.db.GetTotalWithdrawAmountByCardNumber(r.ctx, cardNumber)

	if err != nil {
		return nil, card_errors.ErrGetTotalWithdrawAmountByCardFailed
	}

	return &res, nil
}

func (r *cardDashboardRepository) GetTotalTransactionAmountByCardNumber(cardNumber string) (*int64, error) {
	res, err := r.db.GetTotalTransactionAmountByCardNumber(r.ctx, cardNumber)

	if err != nil {
		return nil, card_errors.ErrGetTotalTransactionAmountByCardFailed
	}

	return &res, nil
}

func (r *cardDashboardRepository) GetTotalTransferAmountBySender(senderCardNumber string) (*int64, error) {
	res, err := r.db.GetTotalTransferAmountBySender(r.ctx, senderCardNumber)

	if err != nil {
		return nil, card_errors.ErrGetTotalTransferAmountBySenderFailed
	}

	return &res, nil
}

func (r *cardDashboardRepository) GetTotalTransferAmountByReceiver(receiverCardNumber string) (*int64, error) {
	res, err := r.db.GetTotalTransferAmountByReceiver(r.ctx, receiverCardNumber)

	if err != nil {
		return nil, card_errors.ErrGetTotalTransferAmountByReceiverFailed
	}

	return &res, nil
}
