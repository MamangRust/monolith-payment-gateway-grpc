package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	refreshtoken_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/refresh_token_errors/repository"
)

// refreshTokenRepository is a struct that implements the RefreshTokenRepository interface
type refreshTokenRepository struct {
	db *db.Queries
}

// NewRefreshTokenRepository creates a new RefreshTokenRepository instance
//
// Args:
// db: a pointer to the database queries
// mapper: a RefreshTokenRecordMapping object
//
// Returns:
// a pointer to the refreshTokenRepository struct
func NewRefreshTokenRepository(db *db.Queries) *refreshTokenRepository {
	return &refreshTokenRepository{
		db: db,
	}
}

// FindByToken retrieves a refresh token record by the token string.
//
// Parameters:
//   - ctx: the context for the database operation
//   - token: the refresh token to search for
//
// Returns:
//   - A RefreshTokenRecord if found, or an error if not found or operation fails.
func (r *refreshTokenRepository) FindByToken(ctx context.Context, token string) (*db.RefreshToken, error) {
	res, err := r.db.FindRefreshTokenByToken(ctx, token)

	if err != nil {
		return nil, refreshtoken_errors.ErrTokenNotFound
	}

	return res, nil
}

// FindByUserId retrieves a refresh token record by the associated user ID.
//
// Parameters:
//   - ctx: the context for the database operation
//   - user_id: the ID of the user
//
// Returns:
//   - A RefreshTokenRecord if found, or an error if not found or operation fails.
func (r *refreshTokenRepository) FindByUserId(ctx context.Context, user_id int) (*db.RefreshToken, error) {
	res, err := r.db.FindRefreshTokenByUserId(ctx, int32(user_id))

	if err != nil {
		return nil, refreshtoken_errors.ErrFindByUserID
	}

	return res, nil
}

// CreateRefreshToken inserts a new refresh token record into the database.
//
// Parameters:
//   - ctx: the context for the database operation
//   - req: the request payload containing token and user information
//
// Returns:
//   - The created RefreshTokenRecord, or an error if the operation fails.
func (r *refreshTokenRepository) CreateRefreshToken(ctx context.Context, req *requests.CreateRefreshToken) (*db.RefreshToken, error) {
	layout := "2006-01-02 15:04:05"
	expirationTime, err := time.Parse(layout, req.ExpiresAt)
	if err != nil {
		return nil, refreshtoken_errors.ErrParseDate
	}

	res, err := r.db.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		UserID:     int32(req.UserId),
		Token:      req.Token,
		Expiration: expirationTime,
	})

	if err != nil {
		return nil, refreshtoken_errors.ErrCreateRefreshToken
	}

	return res, nil
}

// UpdateRefreshToken updates an existing refresh token record.
//
// Parameters:
//   - ctx: the context for the database operation
//   - req: the request payload with updated token data
//
// Returns:
//   - The updated RefreshTokenRecord, or an error if the operation fails.
func (r *refreshTokenRepository) UpdateRefreshToken(ctx context.Context, req *requests.UpdateRefreshToken) (*db.RefreshToken, error) {
	layout := "2006-01-02 15:04:05"
	expirationTime, err := time.Parse(layout, req.ExpiresAt)
	if err != nil {
		return nil, refreshtoken_errors.ErrParseDate
	}

	res, err := r.db.UpdateRefreshTokenByUserId(ctx, db.UpdateRefreshTokenByUserIdParams{
		UserID:     int32(req.UserId),
		Token:      req.Token,
		Expiration: expirationTime,
	})
	if err != nil {
		return nil, refreshtoken_errors.ErrUpdateRefreshToken
	}

	return res, nil
}

// DeleteRefreshToken removes a refresh token by its token string.
//
// Parameters:
//   - ctx: the context for the database operation
//   - token: the refresh token string to delete
//
// Returns:
//   - An error if the deletion fails.
func (r *refreshTokenRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	err := r.db.DeleteRefreshToken(ctx, token)

	if err != nil {
		return refreshtoken_errors.ErrDeleteRefreshToken
	}

	return nil
}

// DeleteRefreshTokenByUserId removes a refresh token by the associated user ID.
//
// Parameters:
//   - ctx: the context for the database operation
//   - user_id: the ID of the user whose token will be deleted
//
// Returns:
//   - An error if the deletion fails.
func (r *refreshTokenRepository) DeleteRefreshTokenByUserId(ctx context.Context, user_id int) error {
	err := r.db.DeleteRefreshTokenByUserId(ctx, int32(user_id))

	if err != nil {
		return refreshtoken_errors.ErrDeleteByUserID
	}

	return nil
}
