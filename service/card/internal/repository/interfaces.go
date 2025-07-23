package repository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

// CardCommandRepository provides methods for creating, updating, and deleting card records in the database.
//
//go:generate mockgen -source=interfaces.go -destination=mocks/card.go
type CardCommandRepository interface {
	// CreateCard generates a new card number, constructs a CreateCardParams object
	// from the provided CreateCardRequest, and inserts a new card record into the database.
	// It returns the created CardRecord or an error if the operation fails.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - request: A CreateCardRequest object containing the details of the card to be created.
	//
	// Returns:
	//   - A pointer to the created CardRecord, or an error if the operation fails.
	CreateCard(ctx context.Context, request *requests.CreateCardRequest) (*record.CardRecord, error)

	// UpdateCard updates a card record in the database.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - request: An UpdateCardRequest object containing the details of the card to be updated.
	//
	// Returns:
	//   - A pointer to the updated CardRecord, or an error if the operation fails.
	UpdateCard(ctx context.Context, request *requests.UpdateCardRequest) (*record.CardRecord, error)

	// TrashedCard permanently deletes a card record from the database.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - card_id: The ID of the card to be trashed.
	//
	// Returns:
	//   - A pointer to the trashed CardRecord, or an error if the operation fails.
	TrashedCard(ctx context.Context, cardId int) (*record.CardRecord, error)

	// RestoreCard restores a previously trashed card by setting its deleted_at field to NULL.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - card_id: The ID of the card to be restored.
	//
	// Returns:
	//   - A pointer to the restored CardRecord, or an error if the operation fails.
	RestoreCard(ctx context.Context, cardId int) (*record.CardRecord, error)

	// DeleteCardPermanent permanently deletes a card record from the database.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - card_id: The ID of the card to be deleted permanently.
	//
	// Returns:
	//   - A boolean indicating if the operation was successful, and an error if the operation fails.
	DeleteCardPermanent(ctx context.Context, card_id int) (bool, error)

	// RestoreAllCard restores all previously trashed card records by setting their deleted_at fields to NULL.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//
	// Returns:
	//   - A boolean indicating if the operation was successful.
	//   - An error if the operation fails.
	RestoreAllCard(ctx context.Context) (bool, error)

	// DeleteAllCardPermanent permanently deletes all card records from the database.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//
	// Returns:
	//   - A boolean indicating if the operation was successful, and an error if the operation fails.
	DeleteAllCardPermanent(ctx context.Context) (bool, error)
}

// CardQueryRepository provides methods for retrieving card records from the database.
type CardQueryRepository interface {
	// FindAllCards retrieves a paginated list of card records based on the search criteria
	// specified in the request. It queries the database and returns a slice of CardRecord,
	// the total count of records, and an error if any occurred.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - req: A FindAllCards request object containing the search parameters
	//     such as search keyword, page number, and page size.
	//
	// Returns:
	//   - A slice of CardRecord representing the card records fetched from the database.
	//   - A pointer to an int representing the total number of records matching the search criteria.
	//   - An error if the operation fails, nil otherwise.
	FindAllCards(ctx context.Context, req *requests.FindAllCards) ([]*record.CardRecord, *int, error)

	// FindByActive retrieves a paginated list of active card records based on the search criteria
	// specified in the request. It queries the database and returns a slice of CardRecord,
	// the total count of records, and an error if any occurred.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - req: A FindAllCards request object containing the search parameters
	//     such as search keyword, page number, and page size.
	//
	// Returns:
	//   - A slice of CardRecord representing the active card records fetched from the database.
	//   - A pointer to an int representing the total number of active card records matching the search criteria.
	//   - An error if the operation fails, nil otherwise.
	FindByActive(ctx context.Context, req *requests.FindAllCards) ([]*record.CardRecord, *int, error)

	// FindByTrashed retrieves a paginated list of trashed card records based on the search criteria
	// specified in the request. It queries the database and returns a slice of CardRecord,
	// the total count of records, and an error if any occurred.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - req: A FindAllCards request object containing the search parameters
	//     such as search keyword, page number, and page size.
	//
	// Returns:
	//   - A slice of CardRecord representing the trashed card records fetched from the database.
	//   - A pointer to an int representing the total number of trashed card records matching the search criteria.
	//   - An error if the operation fails, nil otherwise.
	FindByTrashed(ctx context.Context, req *requests.FindAllCards) ([]*record.CardRecord, *int, error)

	// FindById retrieves a card record by its ID from the database.
	//
	// Parameters:
	//   - card_id: The ID of the card to be retrieved.
	//
	// Returns:
	//   - A pointer to a CardRecord representing the card record fetched from the database.
	//   - An error if the operation fails, nil otherwise.
	FindById(ctx context.Context, card_id int) (*record.CardRecord, error)

	// FindCardByUserId retrieves a card record by its user ID from the database.
	//
	// Parameters:
	//   - user_id: The ID of the user who owns the card.
	//
	// Returns:
	//   - A pointer to a CardRecord representing the card record fetched from the database.
	//   - An error if the operation fails, nil otherwise.
	FindCardByUserId(ctx context.Context, user_id int) (*record.CardRecord, error)

	// FindCardByCardNumber retrieves a card record by its card number from the database.
	//
	// Parameters:
	//   - card_number: The card number of the card to be retrieved.
	//
	// Returns:
	//   - A pointer to a CardRecord representing the card record fetched from the database.
	//   - An error if the operation fails, nil otherwise.
	FindCardByCardNumber(ctx context.Context, card_number string) (*record.CardRecord, error)
}


// UserRepository provides methods for retrieving user data.
type UserRepository interface {
	// FindById retrieves a user by their unique identifier
	//
	// Parameters:
	//   - user_id: the integer unique identifier for the user to retrieve
	//
	// Returns:
	//   - A pointer to the UserRecord if the user is found, or an error if operation fails.
	//   - ErrUserNotFound if the user is not found in the database
	FindById(ctx context.Context, user_id int) (*record.UserRecord, error)
}
