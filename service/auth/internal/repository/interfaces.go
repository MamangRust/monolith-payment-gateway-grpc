package repository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

// UserRepository defines the data access layer for user-related operations.
//
//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go
type UserRepository interface {
	// FindByEmail retrieves a user by their email address.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - email: the email address to search for
	//
	// Returns:
	//   - A UserRecord if found, or an error if the operation fails or user does not exist.
	FindByEmail(ctx context.Context, email string) (*record.UserRecord, error)

	// FindByEmailAndVerify retrieves a verified user by their email address.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - email: the email address to search for
	//
	// Returns:
	//   - A UserRecord if found and verified, or an error otherwise.
	FindByEmailAndVerify(ctx context.Context, email string) (*record.UserRecord, error)

	// FindById retrieves a user by their unique ID.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - id: the user's unique identifier
	//
	// Returns:
	//   - A UserRecord if found, or an error if the operation fails or user is not found.
	FindById(ctx context.Context, id int) (*record.UserRecord, error)

	// CreateUser inserts a new user into the database.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - request: the user registration data
	//
	// Returns:
	//   - The created UserRecord, or an error if the operation fails.
	CreateUser(ctx context.Context, request *requests.RegisterRequest) (*record.UserRecord, error)

	// UpdateUserIsVerified updates the verification status of a user.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - userID: the user's ID
	//   - isVerified: the updated verification status
	//
	// Returns:
	//   - The updated UserRecord, or an error if the operation fails.
	UpdateUserIsVerified(ctx context.Context, userID int, isVerified bool) (*record.UserRecord, error)

	// UpdateUserPassword updates a user's password.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - userID: the user's ID
	//   - password: the new password (hashed)
	//
	// Returns:
	//   - The updated UserRecord, or an error if the operation fails.
	UpdateUserPassword(ctx context.Context, userID int, password string) (*record.UserRecord, error)

	// FindByVerificationCode retrieves a user by their verification code.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - verificationCode: the verification code string
	//
	// Returns:
	//   - A UserRecord if found, or an error if invalid or not found.
	FindByVerificationCode(ctx context.Context, verificationCode string) (*record.UserRecord, error)
}

// ResetTokenRepository defines the data access layer for password reset token operations.
type ResetTokenRepository interface {
	// FindByToken retrieves a reset token record by token string.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - token: the reset token to search for
	//
	// Returns:
	//   - A ResetTokenRecord if found, or an error if not found or operation fails.
	FindByToken(ctx context.Context, token string) (*record.ResetTokenRecord, error)

	// CreateResetToken inserts a new reset token into the database.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - req: the request payload containing user ID and token info
	//
	// Returns:
	//   - The created ResetTokenRecord, or an error if the operation fails.
	CreateResetToken(ctx context.Context, req *requests.CreateResetTokenRequest) (*record.ResetTokenRecord, error)

	// DeleteResetToken removes the reset token associated with the given user ID.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - userID: the user ID whose token should be deleted
	//
	// Returns:
	//   - An error if the deletion fails.
	DeleteResetToken(ctx context.Context, userID int) error
}

// RefreshTokenRepository defines the data access layer for managing refresh tokens.
type RefreshTokenRepository interface {
	// FindByToken retrieves a refresh token record by the token string.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - token: the refresh token to search for
	//
	// Returns:
	//   - A RefreshTokenRecord if found, or an error if not found or operation fails.
	FindByToken(ctx context.Context, token string) (*record.RefreshTokenRecord, error)

	// FindByUserId retrieves a refresh token record by the associated user ID.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - user_id: the ID of the user
	//
	// Returns:
	//   - A RefreshTokenRecord if found, or an error if not found or operation fails.
	FindByUserId(ctx context.Context, user_id int) (*record.RefreshTokenRecord, error)

	// CreateRefreshToken inserts a new refresh token record into the database.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - req: the request payload containing token and user information
	//
	// Returns:
	//   - The created RefreshTokenRecord, or an error if the operation fails.
	CreateRefreshToken(ctx context.Context, req *requests.CreateRefreshToken) (*record.RefreshTokenRecord, error)

	// UpdateRefreshToken updates an existing refresh token record.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - req: the request payload with updated token data
	//
	// Returns:
	//   - The updated RefreshTokenRecord, or an error if the operation fails.
	UpdateRefreshToken(ctx context.Context, req *requests.UpdateRefreshToken) (*record.RefreshTokenRecord, error)

	// DeleteRefreshToken removes a refresh token by its token string.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - token: the refresh token string to delete
	//
	// Returns:
	//   - An error if the deletion fails.
	DeleteRefreshToken(ctx context.Context, token string) error

	// DeleteRefreshTokenByUserId removes a refresh token by the associated user ID.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - user_id: the ID of the user whose token will be deleted
	//
	// Returns:
	//   - An error if the deletion fails.
	DeleteRefreshTokenByUserId(ctx context.Context, user_id int) error
}

// UserRoleRepository defines the data access layer for assigning and removing user roles.
type UserRoleRepository interface {
	// AssignRoleToUser assigns a role to a user.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - req: the request payload containing user ID and role ID
	//
	// Returns:
	//   - The created UserRoleRecord if successful, or an error if the operation fails.
	AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*record.UserRoleRecord, error)

	// RemoveRoleFromUser removes a role assigned to a user.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - req: the request payload containing user ID and role ID
	//
	// Returns:
	//   - An error if the operation fails.
	RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error
}

// RoleRepository defines the data access layer for role-related operations.
type RoleRepository interface {
	// FindById retrieves a role by its unique ID.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - role_id: the unique identifier of the role
	//
	// Returns:
	//   - A RoleRecord if found, or an error if the role does not exist or operation fails.
	FindById(ctx context.Context, role_id int) (*record.RoleRecord, error)

	// FindByName retrieves a role by its name.
	//
	// Parameters:
	//   - ctx: the context for the database operation
	//   - name: the name of the role to search for
	//
	// Returns:
	//   - A RoleRecord if found, or an error if the role does not exist or operation fails.
	FindByName(ctx context.Context, name string) (*record.RoleRecord, error)
}
