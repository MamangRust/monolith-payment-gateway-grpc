package adapter

import (
	"context"
	"time"

	pbcard "github.com/MamangRust/monolith-payment-gateway-pb/card"
	pbsaldo "github.com/MamangRust/monolith-payment-gateway-pb/saldo"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/jackc/pgx/v5/pgtype"
)

type SaldoAdapter struct {
	QueryClient   pbsaldo.SaldoQueryServiceClient
	CommandClient pbsaldo.SaldoCommandServiceClient
}

func NewSaldoAdapter(queryClient pbsaldo.SaldoQueryServiceClient, commandClient pbsaldo.SaldoCommandServiceClient) *SaldoAdapter {
	return &SaldoAdapter{
		QueryClient:   queryClient,
		CommandClient: commandClient,
	}
}

func (a *SaldoAdapter) FindByCardNumber(ctx context.Context, card_number string) (*db.Saldo, error) {
	resp, err := a.QueryClient.FindByCardNumber(ctx, &pbcard.FindByCardNumberRequest{
		CardNumber: card_number,
	})
	if err != nil {
		return nil, err
	}

	return MapSaldoResponseToDB(resp.Data), nil
}

func (a *SaldoAdapter) UpdateSaldoBalance(ctx context.Context, request *requests.UpdateSaldoBalance) (*db.UpdateSaldoBalanceRow, error) {
	resp, err := a.CommandClient.UpdateSaldoBalance(ctx, &pbsaldo.UpdateSaldoBalanceRequest{
		CardNumber:   request.CardNumber,
		TotalBalance: int32(request.TotalBalance),
	})
	if err != nil {
		return nil, err
	}

	return &db.UpdateSaldoBalanceRow{
		SaldoID:      resp.Data.SaldoId,
		CardNumber:   resp.Data.CardNumber,
		TotalBalance: resp.Data.TotalBalance,
	}, nil
}

func (a *SaldoAdapter) UpdateSaldoWithdraw(ctx context.Context, request *requests.UpdateSaldoWithdraw) (*db.UpdateSaldoWithdrawRow, error) {
	resp, err := a.CommandClient.UpdateSaldoWithdraw(ctx, &pbsaldo.UpdateSaldoWithdrawRequest{
		CardNumber:     request.CardNumber,
		TotalBalance:   int32(request.TotalBalance),
		WithdrawAmount: int32(*request.WithdrawAmount),
		WithdrawTime:   request.WithdrawTime.Format(time.RFC3339),
	})
	if err != nil {
		return nil, err
	}

	return &db.UpdateSaldoWithdrawRow{
		SaldoID:      resp.Data.SaldoId,
		CardNumber:   resp.Data.CardNumber,
		TotalBalance: resp.Data.TotalBalance,
	}, nil
}

func MapSaldoResponseToDB(s *pbsaldo.SaldoResponse) *db.Saldo {
	if s == nil {
		return nil
	}

	saldo := &db.Saldo{
		SaldoID:      s.SaldoId,
		CardNumber:   s.CardNumber,
		TotalBalance: s.TotalBalance,
	}

	if s.WithdrawAmount != 0 {
		wa := s.WithdrawAmount
		saldo.WithdrawAmount = &wa
	}

	parseTime := func(ts string) pgtype.Timestamp {
		if ts == "" {
			return pgtype.Timestamp{Valid: false}
		}
		t, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			return pgtype.Timestamp{Valid: false}
		}
		return pgtype.Timestamp{Time: t, Valid: true}
	}

	saldo.WithdrawTime = parseTime(s.WithdrawTime)
	saldo.CreatedAt = parseTime(s.CreatedAt)
	saldo.UpdatedAt = parseTime(s.UpdatedAt)

	return saldo
}
