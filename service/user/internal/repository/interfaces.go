package repository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

// UserQueryRepository defines query operations for retrieving user data from the database.
type UserQueryRepository interface {
	// FindAllUsers retrieves all users with optional filters and pagination.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing filter and pagination information.
	//
	// Returns:
	//   - []*record.UserRecord: List of user records.
	//   - *int: Total count of records.
	//   - error: Error if retrieval fails.
	FindAllUsers(ctx context.Context, req *requests.FindAllUsers) ([]*record.UserRecord, *int, error)

	// FindByActive retrieves all active (non-deleted) users.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing filter and pagination information.
	//
	// Returns:
	//   - []*record.UserRecord: List of active user records.
	//   - *int: Total count of records.
	//   - error: Error if retrieval fails.
	FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*record.UserRecord, *int, error)

	// FindByTrashed retrieves all trashed (soft-deleted) users.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing filter and pagination information.
	//
	// Returns:
	//   - []*record.UserRecord: List of trashed user records.
	//   - *int: Total count of records.
	//   - error: Error if retrieval fails.
	FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*record.UserRecord, *int, error)

	// FindById retrieves a user by their ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - user_id: The unique identifier of the user.
	//
	// Returns:
	//   - *record.UserRecord: The user record if found.
	//   - error: Error if retrieval fails.
	FindById(ctx context.Context, user_id int) (*record.UserRecord, error)

	// FindByEmail retrieves a user by their email address.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - email: The email address of the user.
	//
	// Returns:
	//   - *record.UserRecord: The user record if found.
	//   - error: Error if retrieval fails.
	FindByEmail(ctx context.Context, email string) (*record.UserRecord, error)
}

// UserCommandRepository defines commands for modifying user records in the database.
type UserCommandRepository interface {
	// CreateUser inserts a new user record into the database.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The request payload containing user information.
	//
	// Returns:
	//   - *record.UserRecord: The created user record.
	//   - error: Error if creation fails.
	CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*record.UserRecord, error)

	// UpdateUser updates an existing user record in the database.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The request payload containing updated user information.
	//
	// Returns:
	//   - *record.UserRecord: The updated user record.
	//   - error: Error if update fails.
	UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*record.UserRecord, error)

	// TrashedUser soft-deletes a user by marking it as trashed.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - user_id: The ID of the user to be trashed.
	//
	// Returns:
	//   - *record.UserRecord: The trashed user record.
	//   - error: Error if deletion fails.
	TrashedUser(ctx context.Context, user_id int) (*record.UserRecord, error)

	// RestoreUser restores a soft-deleted (trashed) user.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - user_id: The ID of the user to be restored.
	//
	// Returns:
	//   - *record.UserRecord: The restored user record.
	//   - error: Error if restoration fails.
	RestoreUser(ctx context.Context, user_id int) (*record.UserRecord, error)

	// DeleteUserPermanent permanently deletes a user from the database.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - user_id: The ID of the user to be permanently deleted.
	//
	// Returns:
	//   - bool: Whether the operation was successful.
	//   - error: Error if deletion fails.
	DeleteUserPermanent(ctx context.Context, user_id int) (bool, error)

	// RestoreAllUser restores all soft-deleted users in the database.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: Whether the operation was successful.
	//   - error: Error if restoration fails.
	RestoreAllUser(ctx context.Context) (bool, error)

	// DeleteAllUserPermanent permanently deletes all trashed users from the database.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: Whether the operation was successful.
	//   - error: Error if deletion fails.
	DeleteAllUserPermanent(ctx context.Context) (bool, error)
}

// RoleRepository defines operations for retrieving role information.
type RoleRepository interface {
	// FindById retrieves a role by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - role_id: The ID of the role.
	//
	// Returns:
	//   - *record.RoleRecord: The role record if found.
	//   - error: Error if retrieval fails.
	FindById(ctx context.Context, role_id int) (*record.RoleRecord, error)

	// FindByName retrieves a role by its name.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - name: The name of the role.
	//
	// Returns:
	//   - *record.RoleRecord: The role record if found.
	//   - error: Error if retrieval fails.
	FindByName(ctx context.Context, name string) (*record.RoleRecord, error)
}
