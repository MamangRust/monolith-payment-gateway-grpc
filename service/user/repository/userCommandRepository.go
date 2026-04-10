package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	user_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors/repository"
	"github.com/google/uuid"
)

// userCommandRepository is a struct that implements the UserCommandRepository interface
type userCommandRepository struct {
	db *db.Queries
}

// NewUserCommandRepository creates a new instance of userCommandRepository with the provided
// database queries, context, and user record mapper. This repository is responsible for
// executing command operations related to user records in the database.
//
// Parameters:
//   - db: A pointer to the db.Queries object for executing database queries.
//   - mapper: A UserRecordMapping that provides methods to map database rows to UserRecord domain models.
//
// Returns:
//   - A pointer to the newly created userCommandRepository instance.
func NewUserCommandRepository(db *db.Queries) UserCommandRepository {
	return &userCommandRepository{
		db: db,
	}
}

// CreateUser inserts a new user record into the database.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The request payload containing user information.
//
// Returns:
//   - *record.UserRecord: The created user record.
//   - error: Error if creation fails.
func (r *userCommandRepository) CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*db.CreateUserRow, error) {
	verified := false
	verifyCode := uuid.New().String()

	req := db.CreateUserParams{
		Firstname:        request.FirstName,
		Lastname:         request.LastName,
		Email:            request.Email,
		Password:         request.Password,
		VerificationCode: verifyCode,
		IsVerified:       &verified,
	}

	user, err := r.db.CreateUser(ctx, req)

	if err != nil {
		return nil, user_errors.ErrCreateUser
	}

	return user, nil
}

// UpdateUser updates an existing user record in the database.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The request payload containing updated user information.
//
// Returns:
//   - *record.UserRecord: The updated user record.
//   - error: Error if update fails.
func (r *userCommandRepository) UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*db.UpdateUserRow, error) {
	req := db.UpdateUserParams{
		UserID:    int32(*request.UserID),
		Firstname: request.FirstName,
		Lastname:  request.LastName,
		Email:     request.Email,
		Password:  request.Password,
	}

	res, err := r.db.UpdateUser(ctx, req)

	if err != nil {
		return nil, user_errors.ErrUpdateUser
	}

	return res, nil
}

// TrashedUser soft-deletes a user by marking it as trashed.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - user_id: The ID of the user to be trashed.
//
// Returns:
//   - *record.UserRecord: The trashed user record.
//   - error: Error if deletion fails.
func (r *userCommandRepository) TrashedUser(ctx context.Context, user_id int) (*db.TrashUserRow, error) {
	res, err := r.db.TrashUser(ctx, int32(user_id))

	if err != nil {
		return nil, user_errors.ErrTrashedUser
	}

	return res, nil
}

// RestoreUser restores a soft-deleted (trashed) user.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - user_id: The ID of the user to be restored.
//
// Returns:
//   - *record.UserRecord: The restored user record.
//   - error: Error if restoration fails.
func (r *userCommandRepository) RestoreUser(ctx context.Context, user_id int) (*db.RestoreUserRow, error) {
	res, err := r.db.RestoreUser(ctx, int32(user_id))

	if err != nil {
		return nil, user_errors.ErrRestoreUser
	}

	return res, nil
}

// DeleteUserPermanent permanently deletes a user from the database.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - user_id: The ID of the user to be permanently deleted.
//
// Returns:
//   - bool: Whether the operation was successful.
//   - error: Error if deletion fails.
func (r *userCommandRepository) DeleteUserPermanent(ctx context.Context, user_id int) (bool, error) {
	err := r.db.DeleteUserPermanently(ctx, int32(user_id))

	if err != nil {
		return false, user_errors.ErrDeleteUserPermanent
	}

	return true, nil
}

// RestoreAllUser restores all soft-deleted users in the database.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//
// Returns:
//   - bool: Whether the operation was successful.
//   - error: Error if restoration fails.
func (r *userCommandRepository) RestoreAllUser(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllUsers(ctx)

	if err != nil {
		return false, user_errors.ErrRestoreAllUsers
	}

	return true, nil
}

// DeleteAllUserPermanent permanently deletes all trashed users from the database.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//
// Returns:
//   - bool: Whether the operation was successful.
//   - error: Error if deletion fails.
func (r *userCommandRepository) DeleteAllUserPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentUsers(ctx)

	if err != nil {
		return false, user_errors.ErrDeleteAllUsers
	}
	return true, nil
}
