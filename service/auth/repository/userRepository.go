package repository

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-payment-gateway-shared/errors"
	user_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors/repository"
)

// userRepository is a struct that represents a user repository
type userRepository struct {
	db *db.Queries
}

// NewUserRepository returns a new instance of userRepository.
func NewUserRepository(db *db.Queries) *userRepository {
	return &userRepository{
		db: db,
	}
}

// FindById retrieves a user by their unique ID.
func (r *userRepository) FindById(ctx context.Context, user_id int) (*db.GetUserByIDRow, error) {
	res, err := r.db.GetUserByID(ctx, int32(user_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}

		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}

// FindByEmail retrieves a user by their email address.
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*db.GetUserByEmailRow, error) {
	res, err := r.db.GetUserByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}

		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}

// FindByEmailAndVerify retrieves a verified user by their email address.
func (r *userRepository) FindByEmailAndVerify(ctx context.Context, email string) (*db.GetUserByEmailAndVerifiedRow, error) {
	res, err := r.db.GetUserByEmailAndVerified(ctx, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}

		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}

// FindByVerificationCode retrieves a user by their verification code.
func (r *userRepository) FindByVerificationCode(ctx context.Context, verification_code string) (*db.GetUserByVerificationCodeRow, error) {
	res, err := r.db.GetUserByVerificationCode(ctx, verification_code)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}

// CreateUser inserts a new user into the database.
func (r *userRepository) CreateUser(ctx context.Context, request *requests.RegisterRequest) (*db.CreateUserRow, error) {
	isVerified := true

	req := db.CreateUserParams{
		Firstname:        request.FirstName,
		Lastname:         request.LastName,
		Email:            request.Email,
		Password:         request.Password,
		VerificationCode: request.VerifiedCode,
		IsVerified:       &isVerified,
	}

	user, err := r.db.CreateUser(ctx, req)

	if err != nil {
		return nil, user_errors.ErrCreateUser.WithInternal(err)
	}

	return user, nil
}

// UpdateUserIsVerified updates the verification status of a user.
func (r *userRepository) UpdateUserIsVerified(ctx context.Context, user_id int, is_verified bool) (*db.UpdateUserIsVerifiedRow, error) {
	res, err := r.db.UpdateUserIsVerified(ctx, db.UpdateUserIsVerifiedParams{
		UserID:     int32(user_id),
		IsVerified: &is_verified,
	})

	if err != nil {
		return nil, user_errors.ErrUpdateUserVerificationCode.WithInternal(err)
	}

	return res, nil
}

// UpdateUserPassword updates a user's password.
func (r *userRepository) UpdateUserPassword(ctx context.Context, user_id int, password string) (*db.UpdateUserPasswordRow, error) {
	res, err := r.db.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		UserID:   int32(user_id),
		Password: password,
	})

	if err != nil {
		return nil, user_errors.ErrUpdateUserPassword.WithInternal(err)
	}

	return res, nil
}
