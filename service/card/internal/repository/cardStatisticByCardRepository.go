package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type cardStatisticByCardRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.CardRecordMapping
}

func NewCardStatisticByCardRepository(db *db.Queries, ctx context.Context, mapping recordmapper.CardRecordMapping) *cardStatisticByCardRepository {
	return &cardStatisticByCardRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *cardStatisticByCardRepository) GetMonthlyBalancesByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardMonthBalance, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyBalancesByCardNumber(r.ctx, db.GetMonthlyBalancesByCardNumberParams{
		Column1:    yearStart,
		CardNumber: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetMonthlyBalanceByCardFailed
	}

	return r.mapping.ToMonthlyBalancesCardNumber(res), nil
}

func (r *cardStatisticByCardRepository) GetYearlyBalanceByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardYearlyBalance, error) {
	res, err := r.db.GetYearlyBalancesByCardNumber(r.ctx, db.GetYearlyBalancesByCardNumberParams{
		Column1:    req.Year,
		CardNumber: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetYearlyBalanceByCardFailed
	}

	return r.mapping.ToYearlyBalancesCardNumber(res), nil
}

func (r *cardStatisticByCardRepository) GetMonthlyTopupAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardMonthAmount, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTopupAmountByCardNumber(r.ctx, db.GetMonthlyTopupAmountByCardNumberParams{
		Column2:    yearStart,
		CardNumber: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetMonthlyTopupAmountByCardFailed
	}

	return r.mapping.ToMonthlyTopupAmountsByCardNumber(res), nil
}

func (r *cardStatisticByCardRepository) GetYearlyTopupAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardYearAmount, error) {
	res, err := r.db.GetYearlyTopupAmountByCardNumber(r.ctx, db.GetYearlyTopupAmountByCardNumberParams{
		Column2:    int32(req.Year),
		CardNumber: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetYearlyTopupAmountByCardFailed
	}

	return r.mapping.ToYearlyTopupAmountsByCardNumber(res), nil
}

func (r *cardStatisticByCardRepository) GetMonthlyWithdrawAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardMonthAmount, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyWithdrawAmountByCardNumber(r.ctx, db.GetMonthlyWithdrawAmountByCardNumberParams{
		Column2:    yearStart,
		CardNumber: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetMonthlyWithdrawAmountByCardFailed
	}

	return r.mapping.ToMonthlyWithdrawAmountsByCardNumber(res), nil
}

func (r *cardStatisticByCardRepository) GetYearlyWithdrawAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardYearAmount, error) {
	res, err := r.db.GetYearlyWithdrawAmountByCardNumber(r.ctx, db.GetYearlyWithdrawAmountByCardNumberParams{
		Column2:    int32(req.Year),
		CardNumber: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetYearlyWithdrawAmountByCardFailed
	}

	return r.mapping.ToYearlyWithdrawAmountsByCardNumber(res), nil
}

func (r *cardStatisticByCardRepository) GetMonthlyTransactionAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardMonthAmount, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransactionAmountByCardNumber(r.ctx, db.GetMonthlyTransactionAmountByCardNumberParams{
		Column2:    yearStart,
		CardNumber: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetMonthlyTransactionAmountByCardFailed
	}

	return r.mapping.ToMonthlyTransactionAmountsByCardNumber(res), nil
}

func (r *cardStatisticByCardRepository) GetYearlyTransactionAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*record.CardYearAmount, error) {
	res, err := r.db.GetYearlyTransactionAmountByCardNumber(r.ctx, db.GetYearlyTransactionAmountByCardNumberParams{
		Column2:    int32(req.Year),
		CardNumber: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetYearlyTransactionAmountByCardFailed
	}

	return r.mapping.ToYearlyTransactionAmountsByCardNumber(res), nil
}

func (r *cardStatisticByCardRepository) GetMonthlyTransferAmountBySender(req *requests.MonthYearCardNumberCard) ([]*record.CardMonthAmount, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransferAmountBySender(r.ctx, db.GetMonthlyTransferAmountBySenderParams{
		Column2:      yearStart,
		TransferFrom: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetMonthlyTransferAmountBySenderFailed
	}

	return r.mapping.ToMonthlyTransferSenderAmountsByCardNumber(res), nil
}

func (r *cardStatisticByCardRepository) GetYearlyTransferAmountBySender(req *requests.MonthYearCardNumberCard) ([]*record.CardYearAmount, error) {
	res, err := r.db.GetYearlyTransferAmountBySender(r.ctx, db.GetYearlyTransferAmountBySenderParams{
		Column2:      int32(req.Year),
		TransferFrom: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetYearlyTransferAmountBySenderFailed
	}

	return r.mapping.ToYearlyTransferSenderAmountsByCardNumber(res), nil
}

func (r *cardStatisticByCardRepository) GetMonthlyTransferAmountByReceiver(req *requests.MonthYearCardNumberCard) ([]*record.CardMonthAmount, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransferAmountByReceiver(r.ctx, db.GetMonthlyTransferAmountByReceiverParams{
		Column2:    yearStart,
		TransferTo: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetMonthlyTransferAmountByReceiverFailed
	}

	return r.mapping.ToMonthlyTransferReceiverAmountsByCardNumber(res), nil
}

func (r *cardStatisticByCardRepository) GetYearlyTransferAmountByReceiver(req *requests.MonthYearCardNumberCard) ([]*record.CardYearAmount, error) {
	res, err := r.db.GetYearlyTransferAmountByReceiver(r.ctx, db.GetYearlyTransferAmountByReceiverParams{
		Column2:    int32(req.Year),
		TransferTo: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetYearlyTransferAmountByReceiverFailed
	}

	return r.mapping.ToYearlyTransferReceiverAmountsByCardNumber(res), nil
}
