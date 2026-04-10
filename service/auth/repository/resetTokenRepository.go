package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	refresh_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/refresh_token_errors/repository"
)

// resetTokenRepository is a struct that implements the ResetTokenRepository interface
type resetTokenRepository struct {
	db *db.Queries
}

// NewResetTokenRepository creates a new instance of resetTokenRepository.
//
// Args:
// db: a pointer to the db.Queries object, providing database query capabilities.
// mapper: a ResetTokenRecordMapping object to map database records to domain records.
//
// Returns:
// A pointer to a newly initialized resetTokenRepository struct.
func NewResetTokenRepository(db *db.Queries) *resetTokenRepository {
	return &resetTokenRepository{
		db: db,
	}
}

// FindByToken retrieves a reset token record by token string.
//
// Parameters:
//   - ctx: the context for the database operation
//   - token: the reset token to search for
//
// Returns:
//   - A ResetTokenRecord if found, or an error if not found or operation fails.
func (r *resetTokenRepository) FindByToken(ctx context.Context, code string) (*db.GetResetTokenRow, error) {
	res, err := r.db.GetResetToken(ctx, code)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// CreateResetToken inserts a new reset token into the database.
//
// Parameters:
//   - ctx: the context for the database operation
//   - req: the request payload containing user ID and token info
//
// Returns:
//   - The created ResetTokenRecord, or an error if the operation fails.
func (r *resetTokenRepository) CreateResetToken(ctx context.Context, req *requests.CreateResetTokenRequest) (*db.CreateResetTokenRow, error) {
	expiryDate, err := time.Parse("2006-01-02 15:04:05", req.ExpiredAt)
	if err != nil {
		return nil, err
	}
	res, err := r.db.CreateResetToken(ctx, db.CreateResetTokenParams{
		UserID:     int64(req.UserID),
		Token:      req.ResetToken,
		ExpiryDate: expiryDate,
	})
	if err != nil {
		return nil, refresh_errors.ErrCreateRefreshToken
	}
	return res, nil
}

// DeleteResetToken removes the reset token associated with the given user ID.
//
// Parameters:
//   - ctx: the context for the database operation
//   - userID: the user ID whose token should be deleted
//
// Returns:
//   - An error if the deletion fails.
func (r *resetTokenRepository) DeleteResetToken(ctx context.Context, user_id int) error {
	err := r.db.DeleteResetToken(ctx, int64(user_id))
	if err != nil {
		return refresh_errors.ErrDeleteByUserID
	}
	return nil
}
