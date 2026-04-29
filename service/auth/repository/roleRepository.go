package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"errors"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	sharedErrors "github.com/MamangRust/monolith-payment-gateway-shared/errors"
	role_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/role_errors/repository"
)

// roleRepository is a struct that implements the RoleRepository interface
type roleRepository struct {
	db *db.Queries
}

// NewRoleRepository creates a new RoleRepository instance
func NewRoleRepository(db *db.Queries) *roleRepository {
	return &roleRepository{
		db: db,
	}
}

// FindById retrieves a role by its unique ID.
func (r *roleRepository) FindById(ctx context.Context, id int) (*db.Role, error) {
	res, err := r.db.GetRole(ctx, int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, role_errors.ErrRoleNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return res, nil
}

// FindByName retrieves a role by its name from the database.
func (r *roleRepository) FindByName(ctx context.Context, name string) (*db.Role, error) {
	res, err := r.db.GetRoleByName(ctx, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, role_errors.ErrRoleNotFound.WithInternal(err)
		}

		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return res, nil
}
