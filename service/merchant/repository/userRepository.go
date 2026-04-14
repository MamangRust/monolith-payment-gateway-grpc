package repository

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	sharedErrors "github.com/MamangRust/monolith-payment-gateway-shared/errors"
	user_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors/repository"
)

// userRepository is a struct that represents a repository for user operations.
type userRepository struct {
	db *db.Queries
}

// NewUserRepository creates a new instance of userRepository.
//
// Parameters:
//   - db: A pointer to the db.Queries object for executing database queries.
//   - mapper: A UserRecordMapping that provides methods to map database rows to User domain models.
//
// Returns:
//   - A pointer to the newly created userRepository instance.
func NewUserRepository(db *db.Queries) UserRepository {
	return &userRepository{
		db: db,
	}
}

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
