package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record"
)

type transferStatisticByCardRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.TransferRecordMapping
}

func NewTransferStatisticByCardRepository(db *db.Queries, ctx context.Context, mapping recordmapper.TransferRecordMapping) *transferStatisticByCardRepository {
	return &transferStatisticByCardRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *transferStatisticByCardRepository) GetMonthTransferStatusSuccessByCardNumber(req *requests.MonthStatusTransferCardNumber) ([]*record.TransferRecordMonthStatusSuccess, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransferStatusSuccessCardNumber(r.ctx, db.GetMonthTransferStatusSuccessCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      currentDate,
		Column3:      lastDayCurrentMonth,
		Column4:      prevDate,
		Column5:      lastDayPrevMonth,
	})

	if err != nil {
		return nil, transfer_errors.ErrGetMonthTransferStatusSuccessByCardFailed
	}

	so := r.mapping.ToTransferRecordsMonthStatusSuccessCardNumber(res)

	return so, nil
}

func (r *transferStatisticByCardRepository) GetYearlyTransferStatusSuccessByCardNumber(req *requests.YearStatusTransferCardNumber) ([]*record.TransferRecordYearStatusSuccess, error) {
	res, err := r.db.GetYearlyTransferStatusSuccessCardNumber(r.ctx, db.GetYearlyTransferStatusSuccessCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      int32(req.Year),
	})

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferStatusSuccessByCardFailed
	}

	so := r.mapping.ToTransferRecordsYearStatusSuccessCardNumber(res)

	return so, nil
}

func (r *transferStatisticByCardRepository) GetMonthTransferStatusFailedByCardNumber(req *requests.MonthStatusTransferCardNumber) ([]*record.TransferRecordMonthStatusFailed, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthTransferStatusFailedCardNumber(r.ctx, db.GetMonthTransferStatusFailedCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      currentDate,
		Column3:      lastDayCurrentMonth,
		Column4:      prevDate,
		Column5:      lastDayPrevMonth,
	})

	if err != nil {
		return nil, transfer_errors.ErrGetMonthTransferStatusFailedByCardFailed
	}

	so := r.mapping.ToTransferRecordsMonthStatusFailedCardNumber(res)

	return so, nil
}

func (r *transferStatisticByCardRepository) GetYearlyTransferStatusFailedByCardNumber(req *requests.YearStatusTransferCardNumber) ([]*record.TransferRecordYearStatusFailed, error) {
	res, err := r.db.GetYearlyTransferStatusFailedCardNumber(r.ctx, db.GetYearlyTransferStatusFailedCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      int32(req.Year),
	})

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferStatusFailedByCardFailed
	}

	so := r.mapping.ToTransferRecordsYearStatusFailedCardNumber(res)

	return so, nil
}

func (r *transferStatisticByCardRepository) GetMonthlyTransferAmountsBySenderCardNumber(req *requests.MonthYearCardNumber) ([]*record.TransferMonthAmount, error) {
	res, err := r.db.GetMonthlyTransferAmountsBySenderCardNumber(r.ctx, db.GetMonthlyTransferAmountsBySenderCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	if err != nil {
		return nil, transfer_errors.ErrGetMonthlyTransferAmountsBySenderCardFailed
	}

	return r.mapping.ToTransferMonthAmountsSender(res), nil
}

func (r *transferStatisticByCardRepository) GetMonthlyTransferAmountsByReceiverCardNumber(req *requests.MonthYearCardNumber) ([]*record.TransferMonthAmount, error) {
	res, err := r.db.GetMonthlyTransferAmountsByReceiverCardNumber(r.ctx, db.GetMonthlyTransferAmountsByReceiverCardNumberParams{
		TransferTo: req.CardNumber,
		Column2:    time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	if err != nil {
		return nil, transfer_errors.ErrGetMonthlyTransferAmountsByReceiverCardFailed
	}
	return r.mapping.ToTransferMonthAmountsReceiver(res), nil
}

func (r *transferStatisticByCardRepository) GetYearlyTransferAmountsBySenderCardNumber(req *requests.MonthYearCardNumber) ([]*record.TransferYearAmount, error) {
	res, err := r.db.GetYearlyTransferAmountsBySenderCardNumber(r.ctx, db.GetYearlyTransferAmountsBySenderCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      req.Year,
	})

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferAmountsBySenderCardFailed
	}

	return r.mapping.ToTransferYearAmountsSender(res), nil
}

func (r *transferStatisticByCardRepository) GetYearlyTransferAmountsByReceiverCardNumber(req *requests.MonthYearCardNumber) ([]*record.TransferYearAmount, error) {
	res, err := r.db.GetYearlyTransferAmountsByReceiverCardNumber(r.ctx, db.GetYearlyTransferAmountsByReceiverCardNumberParams{
		TransferTo: req.CardNumber,
		Column2:    req.Year,
	})

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferAmountsByReceiverCardFailed
	}

	return r.mapping.ToTransferYearAmountsReceiver(res), nil
}

func (r *transferStatisticByCardRepository) FindTransferByTransferFrom(transfer_from string) ([]*record.TransferRecord, error) {
	res, err := r.db.GetTransfersBySourceCard(r.ctx, transfer_from)

	if err != nil {
		return nil, transfer_errors.ErrFindTransferByTransferFromFailed
	}

	return r.mapping.ToTransfersRecord(res), nil
}

func (r *transferStatisticByCardRepository) FindTransferByTransferTo(transfer_to string) ([]*record.TransferRecord, error) {
	res, err := r.db.GetTransfersByDestinationCard(r.ctx, transfer_to)

	if err != nil {
		return nil, transfer_errors.ErrFindTransferByTransferToFailed
	}
	return r.mapping.ToTransfersRecord(res), nil
}
