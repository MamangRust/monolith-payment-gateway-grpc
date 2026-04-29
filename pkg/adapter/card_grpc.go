package adapter

import (
	"context"
	"time"

	pbcard "github.com/MamangRust/monolith-payment-gateway-pb/card"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CardAdapter struct {
	QueryClient   pbcard.CardQueryServiceClient
	CommandClient pbcard.CardCommandServiceClient
}

func NewCardAdapter(queryClient pbcard.CardQueryServiceClient, commandClient pbcard.CardCommandServiceClient) *CardAdapter {
	return &CardAdapter{
		QueryClient:   queryClient,
		CommandClient: commandClient,
	}
}

func (a *CardAdapter) FindCardByUserId(ctx context.Context, user_id int) (*db.GetCardByUserIDRow, error) {
	resp, err := a.QueryClient.FindByUserIdCard(ctx, &pbcard.FindByUserIdCardRequest{
		UserId: int32(user_id),
	})
	if err != nil {
		return nil, err
	}

	return &db.GetCardByUserIDRow{
		CardID:       resp.Data.Id,
		UserID:       resp.Data.UserId,
		CardNumber:   resp.Data.CardNumber,
		CardType:     resp.Data.CardType,
		ExpireDate:   parseDate(resp.Data.ExpireDate),
		Cvv:          resp.Data.Cvv,
		CardProvider: resp.Data.CardProvider,
	}, nil
}

func (a *CardAdapter) FindUserCardByCardNumber(ctx context.Context, card_number string) (*db.GetUserEmailByCardNumberRow, error) {
	resp, err := a.QueryClient.FindByCardNumber(ctx, &pbcard.FindByCardNumberRequest{
		CardNumber: card_number,
	})
	if err != nil {
		return nil, err
	}

	return &db.GetUserEmailByCardNumberRow{
		CardNumber: resp.Data.CardNumber,
		Email:      "mapped@example.com",
	}, nil
}

func (a *CardAdapter) FindCardByCardNumber(ctx context.Context, card_number string) (*db.GetCardByCardNumberRow, error) {
	resp, err := a.QueryClient.FindByCardNumber(ctx, &pbcard.FindByCardNumberRequest{
		CardNumber: card_number,
	})
	if err != nil {
		return nil, err
	}

	return &db.GetCardByCardNumberRow{
		CardID:       resp.Data.Id,
		UserID:       resp.Data.UserId,
		CardNumber:   resp.Data.CardNumber,
		CardType:     resp.Data.CardType,
		ExpireDate:   parseDate(resp.Data.ExpireDate),
		Cvv:          resp.Data.Cvv,
		CardProvider: resp.Data.CardProvider,
	}, nil
}

func (a *CardAdapter) UpdateCard(ctx context.Context, request *requests.UpdateCardRequest) (*db.UpdateCardRow, error) {
	resp, err := a.CommandClient.UpdateCard(ctx, &pbcard.UpdateCardRequest{
		CardId:       int32(request.CardID),
		UserId:       int32(request.UserID),
		CardType:     request.CardType,
		ExpireDate:   timestamppb.New(request.ExpireDate),
		Cvv:          request.CVV,
		CardProvider: request.CardProvider,
	})

	if err != nil {
		return nil, err
	}

	return &db.UpdateCardRow{
		CardID:       resp.Data.Id,
		UserID:       resp.Data.UserId,
		CardNumber:   resp.Data.CardNumber,
		CardType:     resp.Data.CardType,
		ExpireDate:   parseDate(resp.Data.ExpireDate),
		Cvv:          resp.Data.Cvv,
		CardProvider: resp.Data.CardProvider,
	}, nil
}

func parseDate(ts string) pgtype.Date {
	if ts == "" {
		return pgtype.Date{Valid: false}
	}
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return pgtype.Date{Valid: false}
	}
	return pgtype.Date{Time: t, Valid: true}
}
