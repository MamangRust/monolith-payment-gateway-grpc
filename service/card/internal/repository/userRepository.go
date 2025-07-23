package repository

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	user_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/user"
)

// userRepository is a struct that represents a user repository
type userRepository struct {
	db     *db.Queries
	mapper recordmapper.UserQueryRecordMapper
}

// NewUserRepository creates a new instance of userRepository.
//
// Parameters:
//   - db: A pointer to the db.Queries object for executing database queries.
//   - mapper: A UserRecordMapping that provides methods to map database rows to User domain models.
//
// Returns:
//   - A pointer to the newly created userRepository instance.
func NewUserRepository(db *db.Queries, mapper recordmapper.UserQueryRecordMapper) UserRepository {
	return &userRepository{
		db:     db,
		mapper: mapper,
	}
}

// FindById retrieves a user by their unique identifier
//
// Parameters:
//   - ctx: the context for the database operation
//   - user_id: the integer unique identifier for the user to retrieve
//
// Returns:
//   - A pointer to the UserRecord if the user is found, or an error if operation fails.
//   - ErrUserNotFound if the user is not found in the database
func (r *userRepository) FindById(ctx context.Context, user_id int) (*record.UserRecord, error) {
	res, err := r.db.GetUserByID(ctx, int32(user_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound
		}

		return nil, user_errors.ErrUserNotFound
	}

	return r.mapper.ToUserRecord(res), nil
}
