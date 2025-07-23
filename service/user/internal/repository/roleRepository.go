package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

// NewRoleRepository creates a new RoleRepository instance with the provided
// database queries, context, and role record mapper. This repository is
// responsible for executing query operations related to role records in the
// database.
//
// Parameters:
//   - db: A pointer to the db.Queries object for executing database queries.
//   - mapper: A RoleRecordMapping that provides methods to map database rows to Role domain models.
//
// Returns:
//   - A pointer to the newly created roleRepository instance.
func NewRoleRepository(db *db.Queries, mapper recordmapper.RoleQueryRecordMapping) RoleRepository {
	return &roleRepository{
		db:     db,
		mapper: mapper,
	}
}

// FindById retrieves a role by its ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - role_id: The ID of the role.
//
// Returns:
//   - *record.RoleRecord: The role record if found.
//   - error: Error if retrieval fails.
func (r *roleRepository) FindById(ctx context.Context, id int) (*record.RoleRecord, error) {
	res, err := r.db.GetRole(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("role not found with ID: %d", id)
		}
		return nil, fmt.Errorf("failed to find role by ID %d: %w", id, err)
	}

	so := r.mapper.ToRoleRecord(res)

	return so, nil
}

// FindByName retrieves a role by its name.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - name: The name of the role.
//
// Returns:
//   - *record.RoleRecord: The role record if found.
//   - error: Error if retrieval fails.
func (r *roleRepository) FindByName(ctx context.Context, name string) (*record.RoleRecord, error) {
	res, err := r.db.GetRoleByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, role_errors.ErrRoleNotFound
		}

		return nil, role_errors.ErrRoleNotFound
	}

	so := r.mapper.ToRoleRecord(res)

	return so, nil
}
