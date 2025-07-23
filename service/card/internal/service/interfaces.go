package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// CardQueryService is an interface for querying cards
//
// CardQueryService defines the business logic for querying card records.
//
//go:generate mockgen -source=interfaces.go -destination=mocks/service.go
type CardQueryService interface {
	// FindAll retrieves a paginated list of card records based on the search criteria
	// specified in the request. It queries the database and returns the results.
	//
	// Parameters:
	//   - ctx: The context for the database operation.
	//   - req: A FindAllCards request object containing the search parameters
	//     such as search keyword, page number, and page size.
	//
	// Returns:
	//   - []*response.CardResponse: A slice of card records fetched from the database.
	//   - *int: The total number of records matching the search criteria.
	//   - *response.ErrorResponse: An error if the operation fails, nil otherwise.
	FindAll(ctx context.Context, req *requests.FindAllCards) ([]*response.CardResponse, *int, *response.ErrorResponse)

	// FindByActive retrieves a paginated list of active (non-deleted) card records
	// based on the search criteria specified in the request.
	//
	// Parameters:
	//   - ctx: The context for the database operation.
	//   - req: A FindAllCards request object containing search filters such as keyword, page, and limit.
	//
	// Returns:
	//   - []*response.CardResponseDeleteAt: A slice of active cards.
	//   - *int: The total number of matching records.
	//   - *response.ErrorResponse: An error if the operation fails, nil otherwise.
	FindByActive(ctx context.Context, req *requests.FindAllCards) ([]*response.CardResponseDeleteAt, *int, *response.ErrorResponse)

	// FindByTrashed retrieves a paginated list of soft-deleted (trashed) card records
	// based on the search criteria specified in the request.
	//
	// Parameters:
	//   - ctx: The context for the database operation.
	//   - req: A FindAllCards request object containing search filters such as keyword, page, and limit.
	//
	// Returns:
	//   - []*response.CardResponseDeleteAt: A slice of trashed cards.
	//   - *int: The total number of matching records.
	//   - *response.ErrorResponse: An error if the operation fails, nil otherwise.
	FindByTrashed(ctx context.Context, req *requests.FindAllCards) ([]*response.CardResponseDeleteAt, *int, *response.ErrorResponse)

	// FindById retrieves a card by its unique ID.
	//
	// Parameters:
	//   - ctx: The context for the database operation.
	//   - cardID: The ID of the card to be retrieved.
	//
	// Returns:
	//   - *response.CardResponse: The card record if found.
	//   - *response.ErrorResponse: An error if the operation fails, nil otherwise.
	FindById(ctx context.Context, cardID int) (*response.CardResponse, *response.ErrorResponse)

	// FindByUserID retrieves a card by the associated user ID.
	//
	// Parameters:
	//   - ctx: The context for the database operation.
	//   - userID: The ID of the user whose card is being retrieved.
	//
	// Returns:
	//   - *response.CardResponse: The card record if found.
	//   - *response.ErrorResponse: An error if the operation fails, nil otherwise.
	FindByUserID(ctx context.Context, userID int) (*response.CardResponse, *response.ErrorResponse)

	// FindByCardNumber retrieves a card by its card number.
	//
	// Parameters:
	//   - ctx: The context for the database operation.
	//   - cardNumber: The card number to search.
	//
	// Returns:
	//   - *response.CardResponse: The card record if found.
	//   - *response.ErrorResponse: An error if the operation fails, nil otherwise.
	FindByCardNumber(ctx context.Context, cardNumber string) (*response.CardResponse, *response.ErrorResponse)
}

// CardDashboardService defines the business logic for retrieving card dashboard statistics.
type CardDashboardService interface {
	// DashboardCard retrieves aggregated dashboard statistics across all cards.
	//
	// Parameters:
	//   - ctx: The context for the database operation.
	//
	// Returns:
	//   - *response.DashboardCard: The aggregated dashboard data across all cards.
	//   - *response.ErrorResponse: An error if the operation fails, nil otherwise.
	DashboardCard(ctx context.Context) (*response.DashboardCard, *response.ErrorResponse)

	// DashboardCardCardNumber retrieves dashboard statistics for a specific card number.
	//
	// Parameters:
	//   - ctx: The context for the database operation.
	//   - cardNumber: The card number to retrieve dashboard data for.
	//
	// Returns:
	//   - *response.DashboardCardCardNumber: The dashboard data for the specified card number.
	//   - *response.ErrorResponse: An error if the operation fails, nil otherwise.
	DashboardCardCardNumber(ctx context.Context, cardNumber string) (*response.DashboardCardCardNumber, *response.ErrorResponse)
}

// CardCommandService defines the business logic for creating, updating,
// deleting, and restoring card records.
type CardCommandService interface {
	// CreateCard creates a new card based on the provided request data.
	//
	// Parameters:
	//   - ctx: The context for the database operation.
	//   - request: The CreateCardRequest object containing card details to be created.
	//
	// Returns:
	//   - *response.CardResponse: The newly created card.
	//   - *response.ErrorResponse: An error if the operation fails, nil otherwise.
	CreateCard(ctx context.Context, request *requests.CreateCardRequest) (*response.CardResponse, *response.ErrorResponse)

	// UpdateCard updates an existing card based on the provided request data.
	//
	// Parameters:
	//   - ctx: The context for the database operation.
	//   - request: The UpdateCardRequest object containing updated card details.
	//
	// Returns:
	//   - *response.CardResponse: The updated card.
	//   - *response.ErrorResponse: An error if the operation fails, nil otherwise.
	UpdateCard(ctx context.Context, request *requests.UpdateCardRequest) (*response.CardResponse, *response.ErrorResponse)

	// TrashedCard soft-deletes (trashes) a card by its ID.
	//
	// Parameters:
	//   - ctx: The context for the database operation.
	//   - cardId: The ID of the card to be trashed.
	//
	// Returns:
	//   - *response.CardResponse: The trashed card.
	//   - *response.ErrorResponse: An error if the operation fails, nil otherwise.
	TrashedCard(ctx context.Context, cardId int) (*response.CardResponseDeleteAt, *response.ErrorResponse)

	// RestoreCard restores a soft-deleted card by its ID.
	//
	// Parameters:
	//   - ctx: The context for the database operation.
	//   - cardId: The ID of the card to be restored.
	//
	// Returns:
	//   - *response.CardResponse: The restored card.
	//   - *response.ErrorResponse: An error if the operation fails, nil otherwise.
	RestoreCard(ctx context.Context, cardId int) (*response.CardResponse, *response.ErrorResponse)

	// DeleteCardPermanent permanently deletes a card by its ID.
	//
	// Parameters:
	//   - ctx: The context for the database operation.
	//   - cardId: The ID of the card to be permanently deleted.
	//
	// Returns:
	//   - bool: True if deletion was successful.
	//   - *response.ErrorResponse: An error if the operation fails, nil otherwise.
	DeleteCardPermanent(ctx context.Context, cardId int) (bool, *response.ErrorResponse)

	// RestoreAllCard restores all soft-deleted cards.
	//
	// Parameters:
	//   - ctx: The context for the database operation.
	//
	// Returns:
	//   - bool: True if restoration was successful.
	//   - *response.ErrorResponse: An error if the operation fails, nil otherwise.
	RestoreAllCard(ctx context.Context) (bool, *response.ErrorResponse)

	// DeleteAllCardPermanent permanently deletes all trashed cards.
	//
	// Parameters:
	//   - ctx: The context for the database operation.
	//
	// Returns:
	//   - bool: True if deletion was successful.
	//   - *response.ErrorResponse: An error if the operation fails, nil otherwise.
	DeleteAllCardPermanent(ctx context.Context) (bool, *response.ErrorResponse)
}
