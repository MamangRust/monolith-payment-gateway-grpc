package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"errors"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	user_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors/repository"
)

type userRepository struct {
	db *db.Queries
}

func NewUserRepository(db *db.Queries) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) FindById(ctx context.Context, user_id int) (*db.GetUserByIDRow, error) {
	res, err := r.db.GetUserByID(ctx, int32(user_id))

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound
		}

		return nil, user_errors.ErrUserNotFound
	}

	return res, nil
}
