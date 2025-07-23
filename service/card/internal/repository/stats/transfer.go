package repositorystats

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/card/stats"
)

type cardStaticsTransferRepository struct {
	db     *db.Queries
	mapper recordmapper.CardStatisticTransferRecordMapper
}

func NewCardStatsTransferRepository(db *db.Queries, mapper recordmapper.CardStatisticTransferRecordMapper) CardStatsTransferRepository {
	return &cardStaticsTransferRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyTransferAmountSender retrieves the monthly transfer amount data
// for all cards where the card is the sender for a given year.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - year: The year for which the transfer amount is requested.
//
// Returns:
//   - A slice of pointers to CardMonthAmount containing the transfer amount data for each month of the given year.
//   - An error if the retrieval fails, of type ErrGetMonthlyTransferAmountSenderFailed.
func (r *cardStaticsTransferRepository) GetMonthlyTransferAmountSender(ctx context.Context, year int) ([]*record.CardMonthAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransferAmountSender(ctx, yearStart)

	if err != nil {
		return nil, card_errors.ErrGetMonthlyTransferAmountSenderFailed
	}

	return r.mapper.ToMonthlyTransferSenderAmounts(res), nil
}

// GetYearlyTransferAmountSender retrieves the yearly transfer amount data
// for all cards where the card is the sender for a given year.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - year: The year for which the transfer amount is requested.
//
// Returns:
//   - A slice of pointers to CardYearAmount containing the transfer amount data for the specified year.
//   - An error if the retrieval fails, of type ErrGetYearlyTransferAmountSenderFailed.
func (r *cardStaticsTransferRepository) GetYearlyTransferAmountSender(ctx context.Context, year int) ([]*record.CardYearAmount, error) {
	res, err := r.db.GetYearlyTransferAmountSender(ctx, int32(year))

	if err != nil {
		return nil, card_errors.ErrGetYearlyTransferAmountSenderFailed
	}

	return r.mapper.ToYearlyTransferSenderAmounts(res), nil
}

// GetMonthlyTransferAmountReceiver retrieves the monthly transfer amount data
// for all cards where the card is the receiver for a given year.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - year: The year for which the transfer amount is requested.
//
// Returns:
//   - A slice of pointers to CardMonthAmount containing the transfer amount data for each month of the given year.
//   - An error if the retrieval fails, of type ErrGetMonthlyTransferAmountReceiverFailed.
func (r *cardStaticsTransferRepository) GetMonthlyTransferAmountReceiver(ctx context.Context, year int) ([]*record.CardMonthAmount, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyTransferAmountReceiver(ctx, yearStart)

	if err != nil {
		return nil, card_errors.ErrGetMonthlyTransferAmountReceiverFailed
	}

	return r.mapper.ToMonthlyTransferReceiverAmounts(res), nil
}

// GetYearlyTransferAmountReceiver retrieves the yearly transfer amount data
// for all cards where the card is the receiver for a specific year.
//
// Parameters:
//   - ctx: The context for the database operation.
//   - year: The year for which the transfer amount is requested.
//
// Returns:
//   - A slice of pointers to CardYearAmount containing the transfer amount data for the specified year.
//   - An error if the retrieval fails, of type ErrGetYearlyTransferAmountReceiverFailed.
func (r *cardStaticsTransferRepository) GetYearlyTransferAmountReceiver(ctx context.Context, year int) ([]*record.CardYearAmount, error) {
	res, err := r.db.GetYearlyTransferAmountReceiver(ctx, int32(year))

	if err != nil {
		return nil, card_errors.ErrGetYearlyTransferAmountReceiverFailed
	}

	return r.mapper.ToYearlyTransferReceiverAmounts(res), nil
}
