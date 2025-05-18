package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type cardStatisticRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.CardRecordMapping
}

func NewCardStatisticRepository(db *db.Queries, ctx context.Context, mapping recordmapper.CardRecordMapping) *cardStatisticRepository {
	return &cardStatisticRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *cardStatisticRepository) GetMonthlyBalance(year int) ([]*record.CardMonthBalance, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyBalances(r.ctx, yearStart)

	if err != nil {
		return nil, card_errors.ErrGetMonthlyBalanceFailed
	}

	return r.mapping.ToMonthlyBalances(res), nil
}

func (r *cardStatisticRepository) GetYearlyBalance(year int) ([]*record.CardYearlyBalance, error) {
	res, err := r.db.GetYearlyBalances(r.ctx, int32(year))

	if err != nil {
		return nil, card_errors.ErrGetYearlyBalanceFailed
	}

	return r.mapping.ToYearlyBalances(res), nil
}

func (r *cardStatisticRepository) GetMonthlyTopupAmount(year int) ([]*record.CardMonthAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupAmount(r.ctx, yearStart)

	if err != nil {
		return nil, card_errors.ErrGetMonthlyTopupAmountFailed
	}

	return r.mapping.ToMonthlyTopupAmounts(res), nil
}

func (r *cardStatisticRepository) GetYearlyTopupAmount(year int) ([]*record.CardYearAmount, error) {
	res, err := r.db.GetYearlyTopupAmount(r.ctx, int32(year))

	if err != nil {
		return nil, card_errors.ErrGetYearlyTopupAmountFailed
	}

	return r.mapping.ToYearlyTopupAmounts(res), nil
}

func (r *cardStatisticRepository) GetMonthlyWithdrawAmount(year int) ([]*record.CardMonthAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyWithdrawAmount(r.ctx, yearStart)

	if err != nil {
		return nil, card_errors.ErrGetMonthlyWithdrawAmountFailed
	}

	return r.mapping.ToMonthlyWithdrawAmounts(res), nil
}

func (r *cardStatisticRepository) GetYearlyWithdrawAmount(year int) ([]*record.CardYearAmount, error) {
	res, err := r.db.GetYearlyWithdrawAmount(r.ctx, int32(year))

	if err != nil {
		return nil, card_errors.ErrGetYearlyWithdrawAmountFailed
	}

	return r.mapping.ToYearlyWithdrawAmounts(res), nil
}

func (r *cardStatisticRepository) GetMonthlyTransactionAmount(year int) ([]*record.CardMonthAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransactionAmount(r.ctx, yearStart)

	if err != nil {
		return nil, card_errors.ErrGetMonthlyTransactionAmountFailed
	}

	return r.mapping.ToMonthlyTransactionAmounts(res), nil
}

func (r *cardStatisticRepository) GetYearlyTransactionAmount(year int) ([]*record.CardYearAmount, error) {
	res, err := r.db.GetYearlyTransactionAmount(r.ctx, int32(year))

	if err != nil {
		return nil, card_errors.ErrGetYearlyTransactionAmountFailed
	}

	return r.mapping.ToYearlyTransactionAmounts(res), nil
}

func (r *cardStatisticRepository) GetMonthlyTransferAmountSender(year int) ([]*record.CardMonthAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransferAmountSender(r.ctx, yearStart)

	if err != nil {
		return nil, card_errors.ErrGetMonthlyTransferAmountSenderFailed
	}

	return r.mapping.ToMonthlyTransferSenderAmounts(res), nil
}

func (r *cardStatisticRepository) GetYearlyTransferAmountSender(year int) ([]*record.CardYearAmount, error) {
	res, err := r.db.GetYearlyTransferAmountSender(r.ctx, int32(year))

	if err != nil {
		return nil, card_errors.ErrGetYearlyTransferAmountSenderFailed
	}

	return r.mapping.ToYearlyTransferSenderAmounts(res), nil
}

func (r *cardStatisticRepository) GetMonthlyTransferAmountReceiver(year int) ([]*record.CardMonthAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransferAmountReceiver(r.ctx, yearStart)

	if err != nil {
		return nil, card_errors.ErrGetMonthlyTransferAmountReceiverFailed
	}

	return r.mapping.ToMonthlyTransferReceiverAmounts(res), nil
}

func (r *cardStatisticRepository) GetYearlyTransferAmountReceiver(year int) ([]*record.CardYearAmount, error) {
	res, err := r.db.GetYearlyTransferAmountReceiver(r.ctx, int32(year))

	if err != nil {
		return nil, card_errors.ErrGetYearlyTransferAmountReceiverFailed
	}

	return r.mapping.ToYearlyTransferReceiverAmounts(res), nil
}
