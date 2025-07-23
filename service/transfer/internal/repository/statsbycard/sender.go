package transferstatsbycardrepository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/transfer/statsbycard"
)

type transferStatsAmountSenderRepository struct {
	db     *db.Queries
	mapper recordmapper.TransferStatisticSenderAmountByCardRecordMapper
}

func NewTransferStatsAmountSenderRepository(db *db.Queries, mapper recordmapper.TransferStatisticSenderAmountByCardRecordMapper) TransferStatsByCardAmountSenderRepository {
	return &transferStatsAmountSenderRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyTransferAmountsBySenderCardNumber retrieves monthly transfer amounts where the card is the sender.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing the card number and year.
//
// Returns:
//   - []*record.TransferMonthAmount: List of monthly transfer amount records.
//   - error: Any error encountered during the operation.
func (r *transferStatsAmountSenderRepository) GetMonthlyTransferAmountsBySenderCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*record.TransferMonthAmount, error) {
	res, err := r.db.GetMonthlyTransferAmountsBySenderCardNumber(ctx, db.GetMonthlyTransferAmountsBySenderCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	if err != nil {
		return nil, transfer_errors.ErrGetMonthlyTransferAmountsBySenderCardFailed
	}

	so := r.mapper.ToTransferMonthAmountsSender(res)

	return so, nil
}

// GetYearlyTransferAmountsBySenderCardNumber retrieves yearly transfer amounts where the card is the sender.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing the card number and year.
//
// Returns:
//   - []*record.TransferYearAmount: List of yearly transfer amount records.
//   - error: Any error encountered during the operation.
func (r *transferStatsAmountSenderRepository) GetYearlyTransferAmountsBySenderCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*record.TransferYearAmount, error) {
	res, err := r.db.GetYearlyTransferAmountsBySenderCardNumber(ctx, db.GetYearlyTransferAmountsBySenderCardNumberParams{
		TransferFrom: req.CardNumber,
		Column2:      req.Year,
	})

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferAmountsBySenderCardFailed
	}

	so := r.mapper.ToTransferYearAmountsSender(res)

	return so, nil
}
