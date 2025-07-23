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

type transferStatsAmountReceiverRepository struct {
	db     *db.Queries
	mapper recordmapper.TransferStatisticReceiverAmountByCardRecordMapper
}

func NewTransferStatsAmountReceiverRepository(db *db.Queries, mapper recordmapper.TransferStatisticReceiverAmountByCardRecordMapper) TransferStatsByCardAmountReceiverRepository {
	return &transferStatsAmountReceiverRepository{
		db:     db,
		mapper: mapper,
	}
}

// GetMonthlyTransferAmountsByReceiverCardNumber retrieves monthly transfer amounts where the card is the receiver.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing the card number and year.
//
// Returns:
//   - []*record.TransferMonthAmount: List of monthly transfer amount records.
//   - error: Any error encountered during the operation.
func (r *transferStatsAmountReceiverRepository) GetMonthlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*record.TransferMonthAmount, error) {
	res, err := r.db.GetMonthlyTransferAmountsByReceiverCardNumber(ctx, db.GetMonthlyTransferAmountsByReceiverCardNumberParams{
		TransferTo: req.CardNumber,
		Column2:    time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC),
	})

	if err != nil {
		return nil, transfer_errors.ErrGetMonthlyTransferAmountsByReceiverCardFailed
	}

	so := r.mapper.ToTransferMonthAmountsReceiver(res)

	return so, nil
}

// GetYearlyTransferAmountsByReceiverCardNumber retrieves yearly transfer amounts where the card is the receiver.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing the card number and year.
//
// Returns:
//   - []*record.TransferYearAmount: List of yearly transfer amount records.
//   - error: Any error encountered during the operation.
func (r *transferStatsAmountReceiverRepository) GetYearlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*record.TransferYearAmount, error) {
	res, err := r.db.GetYearlyTransferAmountsByReceiverCardNumber(ctx, db.GetYearlyTransferAmountsByReceiverCardNumberParams{
		TransferTo: req.CardNumber,
		Column2:    req.Year,
	})

	if err != nil {
		return nil, transfer_errors.ErrGetYearlyTransferAmountsByReceiverCardFailed
	}

	so := r.mapper.ToTransferYearAmountsReceiver(res)

	return so, nil
}
