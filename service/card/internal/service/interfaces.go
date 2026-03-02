package service

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/service.go
type CardQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllCards) ([]*db.GetCardsRow, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllCards) ([]*db.GetActiveCardsWithCountRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllCards) ([]*db.GetTrashedCardsWithCountRow, *int, error)
	FindById(ctx context.Context, card_id int) (*db.GetCardByIDRow, error)
	FindByUserID(ctx context.Context, userID int) (*db.GetCardByUserIDRow, error)
	FindByCardNumber(ctx context.Context, card_number string) (*db.GetCardByCardNumberRow, error)
}

type CardDashboardService interface {
	DashboardCard(ctx context.Context) (*response.DashboardCard, error)
	DashboardCardCardNumber(ctx context.Context, cardNumber string) (*response.DashboardCardCardNumber, error)
}

type CardCommandService interface {
	CreateCard(ctx context.Context, request *requests.CreateCardRequest) (*db.CreateCardRow, error)
	UpdateCard(ctx context.Context, request *requests.UpdateCardRequest) (*db.UpdateCardRow, error)
	TrashedCard(ctx context.Context, cardId int) (*db.Card, error)
	RestoreCard(ctx context.Context, cardId int) (*db.Card, error)
	DeleteCardPermanent(ctx context.Context, cardId int) (bool, error)

	RestoreAllCard(ctx context.Context) (bool, error)
	DeleteAllCardPermanent(ctx context.Context) (bool, error)
}
