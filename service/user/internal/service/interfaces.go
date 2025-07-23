package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// UserQueryService handles query operations related to user data.
type UserQueryService interface {
	// FindAll retrieves all users based on the given request filter.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing filter parameters for users.
	//
	// Returns:
	//   - []*response.UserResponse: List of user data.
	//   - *int: Total count of users.
	//   - *response.ErrorResponse: Error response if query fails.
	FindAll(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponse, *int, *response.ErrorResponse)

	// FindByID retrieves a specific user by ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - id: The ID of the user.
	//
	// Returns:
	//   - *response.UserResponse: The user data.
	//   - *response.ErrorResponse: Error response if retrieval fails.
	FindByID(ctx context.Context, id int) (*response.UserResponse, *response.ErrorResponse)

	// FindByActive retrieves all active users (not soft-deleted).
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing filter parameters.
	//
	// Returns:
	//   - []*response.UserResponseDeleteAt: List of active user data.
	//   - *int: Total count of active users.
	//   - *response.ErrorResponse: Error response if query fails.
	FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse)

	// FindByTrashed retrieves all soft-deleted (trashed) users.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing filter parameters.
	//
	// Returns:
	//   - []*response.UserResponseDeleteAt: List of trashed user data.
	//   - *int: Total count of trashed users.
	//   - *response.ErrorResponse: Error response if query fails.
	FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse)
}

// UserCommandService handles command operations related to user management.
type UserCommandService interface {
	// CreateUser creates a new user with the provided request data.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The user creation request payload.
	//
	// Returns:
	//   - *response.UserResponse: The created user response.
	//   - *response.ErrorResponse: Error response if creation fails.
	CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*response.UserResponse, *response.ErrorResponse)

	// UpdateUser updates an existing user with the provided request data.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The user update request payload.
	//
	// Returns:
	//   - *response.UserResponse: The updated user response.
	//   - *response.ErrorResponse: Error response if update fails.
	UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*response.UserResponse, *response.ErrorResponse)

	// TrashedUser soft-deletes a user by marking the user as trashed.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - user_id: The ID of the user to be trashed.
	//
	// Returns:
	//   - *response.UserResponseDeleteAt: Response including soft-delete timestamp.
	//   - *response.ErrorResponse: Error response if trash operation fails.
	TrashedUser(ctx context.Context, user_id int) (*response.UserResponseDeleteAt, *response.ErrorResponse)

	// RestoreUser restores a previously soft-deleted user.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - user_id: The ID of the user to be restored.
	//
	// Returns:
	//   - *response.UserResponseDeleteAt: Response with restoration details.
	//   - *response.ErrorResponse: Error response if restore fails.
	RestoreUser(ctx context.Context, user_id int) (*response.UserResponse, *response.ErrorResponse)

	// DeleteUserPermanent permanently deletes a user by ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - user_id: The ID of the user to delete permanently.
	//
	// Returns:
	//   - bool: Whether the deletion was successful.
	//   - *response.ErrorResponse: Error response if deletion fails.
	DeleteUserPermanent(ctx context.Context, user_id int) (bool, *response.ErrorResponse)

	// RestoreAllUser restores all soft-deleted users.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: Whether all users were successfully restored.
	//   - *response.ErrorResponse: Error response if restoration fails.
	RestoreAllUser(ctx context.Context) (bool, *response.ErrorResponse)

	// DeleteAllUserPermanent permanently deletes all trashed users.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: Whether all users were successfully deleted permanently.
	//   - *response.ErrorResponse: Error response if deletion fails.
	DeleteAllUserPermanent(ctx context.Context) (bool, *response.ErrorResponse)
}
