package repository

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	role_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/role_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/role"
)

// roleRepository is a struct that implements the RoleRepository interface
type roleRepository struct {
	db     *db.Queries
	mapper recordmapper.RoleQueryRecordMapping
}

// NewRoleRepository creates a new RoleRepository instance
//
// Args:
// db: a pointer to the database queries
// ctx: a context.Context object
// mapper: a RoleRecordMapping object
//
// Returns:
// a pointer to the roleRepository struct
func NewRoleRepository(db *db.Queries, mapper recordmapper.RoleQueryRecordMapping) *roleRepository {
	return &roleRepository{
		db:     db,
		mapper: mapper,
	}
}

// FindById retrieves a role by its unique ID.
//
// Parameters:
//   - ctx: the context for the database operation
//   - role_id: the unique identifier of the role
//
// Returns:
//   - A RoleRecord if found, or an error if the role does not exist or operation fails.
func (r *roleRepository) FindById(ctx context.Context, id int) (*record.RoleRecord, error) {
	res, err := r.db.GetRole(ctx, int32(id))
	if err != nil {
		return nil, role_errors.ErrRoleNotFound
	}
	return r.mapper.ToRoleRecord(res), nil
}

// FindByName retrieves a role by its name from the database.
//
// Args:
// name: The name of the role to retrieve.
// // FindByName retrieves a role by its name.
//
// Parameters:
//   - ctx: the context for the database operation
//   - name: the name of the role to search for
//
// Returns:
//   - A RoleRecord if found, or an error if the rol
func (r *roleRepository) FindByName(ctx context.Context, name string) (*record.RoleRecord, error) {
	res, err := r.db.GetRoleByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, role_errors.ErrRoleNotFound
	}
	return r.mapper.ToRoleRecord(res), nil
}
