package repositorystatsbycard

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/card/statsbycard"
)

type cardStatsTransferByCardRepository struct {
	db     *db.Queries
	mapper recordmapper.CardStatisticTransferByCardRecordMapper
}

func NewCardStatsTransferByCardRepository(db *db.Queries, mapper recordmapper.CardStatisticTransferByCardRecordMapper) CardStatsTransferByCardRepository {
	return &cardStatsTransferByCardRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyTransferAmountBySender retrieves the monthly transfer amount data
// for a given card number and year where the card is the sender.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - req: A pointer to a MonthYearCardNumberCard request object, containing the year and card number.
//
// Returns:
//   - A slice of pointers to CardMonthAmount containing the transfer amount data for each month of the given year.
//   - An error if the retrieval fails, of type ErrGetMonthlyTransferAmountBySenderFailed.
func (r *cardStatsTransferByCardRepository) GetMonthlyTransferAmountBySender(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*record.CardMonthAmount, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransferAmountBySender(ctx, db.GetMonthlyTransferAmountBySenderParams{
		Column2:      yearStart,
		TransferFrom: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetMonthlyTransferAmountBySenderFailed
	}

	return r.mapper.ToMonthlyTransferSenderAmountsByCardNumber(res), nil
}

// GetYearlyTransferAmountBySender retrieves the yearly transfer amount data
// for a given card number and year where the card is the sender.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - req: A pointer to a MonthYearCardNumberCard request object, containing the year and card number.
//
// Returns:
//   - A slice of pointers to CardYearAmount containing the transfer amount data for the given year.
//   - An error if the retrieval fails, of type ErrGetYearlyTransferAmountBySenderFailed.
func (r *cardStatsTransferByCardRepository) GetYearlyTransferAmountBySender(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*record.CardYearAmount, error) {
	res, err := r.db.GetYearlyTransferAmountBySender(ctx, db.GetYearlyTransferAmountBySenderParams{
		Column2:      int32(req.Year),
		TransferFrom: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetYearlyTransferAmountBySenderFailed
	}

	return r.mapper.ToYearlyTransferSenderAmountsByCardNumber(res), nil
}

// GetMonthlyTransferAmountByReceiver retrieves the monthly transfer amount data
// for a given card number and year where the card is the receiver.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - req: A pointer to a MonthYearCardNumberCard request object, containing the year and card number.
//
// Returns:
//   - A slice of pointers to CardMonthAmount containing the transfer amount data for each month of the given year.
//   - An error if the retrieval fails, of type ErrGetMonthlyTransferAmountByReceiverFailed.
func (r *cardStatsTransferByCardRepository) GetMonthlyTransferAmountByReceiver(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*record.CardMonthAmount, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransferAmountByReceiver(ctx, db.GetMonthlyTransferAmountByReceiverParams{
		Column2:    yearStart,
		TransferTo: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetMonthlyTransferAmountByReceiverFailed
	}

	return r.mapper.ToMonthlyTransferReceiverAmountsByCardNumber(res), nil
}

// GetYearlyTransferAmountByReceiver retrieves the yearly transfer amount data
// for a given card number and year where the card is the receiver.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - req: A pointer to a MonthYearCardNumberCard request object, containing the year and card number.
//
// Returns:
//   - A slice of pointers to CardYearAmount containing the transfer amount data for the given year.
//   - An error if the retrieval fails, of type ErrGetYearlyTransferAmountByReceiverFailed.
func (r *cardStatsTransferByCardRepository) GetYearlyTransferAmountByReceiver(ctx context.Context, req *requests.MonthYearCardNumberCard) ([]*record.CardYearAmount, error) {
	res, err := r.db.GetYearlyTransferAmountByReceiver(ctx, db.GetYearlyTransferAmountByReceiverParams{
		Column2:    int32(req.Year),
		TransferTo: req.CardNumber,
	})

	if err != nil {
		return nil, card_errors.ErrGetYearlyTransferAmountByReceiverFailed
	}

	return r.mapper.ToYearlyTransferReceiverAmountsByCardNumber(res), nil
}
