package repository

import (
	"context"
	"errors"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	role_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/role_errors/repository"
	rolerecordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/role"
)

var errRoleDeadlineExceeded = errors.New("request deadline exceeded while processing role operation")
var errRoleCancelled = errors.New("request was cancelled during role operation")

// roleCommandRepository is a struct that implements the RoleCommandRepository interface
type roleCommandRepository struct {
	db     *db.Queries
	mapper rolerecordmapper.RoleCommandRecordMapping
}

// NewRoleCommandRepository creates a new RoleCommandRepository instance with the provided
// database queries, context, and role record mapper. This repository is responsible for
// executing command operations related to role records in the database.
//
// Parameters:
//   - db: A pointer to the db.Queries object for executing database queries.
//   - mapper: A RoleRecordMapping that provides methods to map database rows to Role domain models.
//
// Returns:
//   - A pointer to the newly created roleCommandRepository instance.
func NewRoleCommandRepository(db *db.Queries, mapper rolerecordmapper.RoleCommandRecordMapping) RoleCommandRepository {
	return &roleCommandRepository{
		db:     db,
		mapper: mapper,
	}
}

// CreateRole creates a new role in the database.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The request payload containing role data to be created.
//
// Returns:
//   - *record.RoleRecord: The newly created role record.
//   - error: An error if the creation failed.
func (r *roleCommandRepository) CreateRole(ctx context.Context, req *requests.CreateRoleRequest) (*record.RoleRecord, error) {
	res, err := r.db.CreateRole(ctx, req.Name)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, errRoleDeadlineExceeded
		}
		if ctx.Err() == context.Canceled {
			return nil, errRoleCancelled
		}
		return nil, role_errors.ErrCreateRole
	}

	return r.mapper.ToRoleRecord(res), nil
}

// UpdateRole updates an existing role in the database.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The request payload containing updated role data.
//
// Returns:
//   - *record.RoleRecord: The updated role record.
//   - error: An error if the update failed.
func (r *roleCommandRepository) UpdateRole(ctx context.Context, req *requests.UpdateRoleRequest) (*record.RoleRecord, error) {
	res, err := r.db.UpdateRole(ctx, db.UpdateRoleParams{
		RoleID:   int32(*req.ID),
		RoleName: req.Name,
	})

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, errRoleDeadlineExceeded
		}
		if ctx.Err() == context.Canceled {
			return nil, errRoleCancelled
		}
		return nil, role_errors.ErrUpdateRole
	}

	return r.mapper.ToRoleRecord(res), nil
}

// TrashedRole performs a soft-delete (move to trash) for a role by its ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - role_id: The ID of the role to be trashed.
//
// Returns:
//   - *record.RoleRecord: The trashed role record.
//   - error: An error if the operation failed.
func (r *roleCommandRepository) TrashedRole(ctx context.Context, id int) (*record.RoleRecord, error) {
	res, err := r.db.TrashRole(ctx, int32(id))

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, errRoleDeadlineExceeded
		}
		if ctx.Err() == context.Canceled {
			return nil, errRoleCancelled
		}
		return nil, role_errors.ErrTrashedRole
	}

	return r.mapper.ToRoleRecord(res), nil
}

// RestoreRole restores a soft-deleted (trashed) role by its ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - role_id: The ID of the role to restore.
//
// Returns:
//   - *record.RoleRecord: The restored role record.
//   - error: An error if the restoration failed.
func (r *roleCommandRepository) RestoreRole(ctx context.Context, id int) (*record.RoleRecord, error) {
	res, err := r.db.RestoreRole(ctx, int32(id))

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, errRoleDeadlineExceeded
		}
		if ctx.Err() == context.Canceled {
			return nil, errRoleCancelled
		}
		return nil, role_errors.ErrRestoreRole
	}

	return r.mapper.ToRoleRecord(res), nil
}

// DeleteRolePermanent permanently deletes a role by its ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - role_id: The ID of the role to delete permanently.
//
// Returns:
//   - bool: True if deletion was successful.
//   - error: An error if the deletion failed.
func (r *roleCommandRepository) DeleteRolePermanent(ctx context.Context, role_id int) (bool, error) {
	err := r.db.DeletePermanentRole(ctx, int32(role_id))

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return false, errRoleDeadlineExceeded
		}
		if ctx.Err() == context.Canceled {
			return false, errRoleCancelled
		}
		return false, role_errors.ErrDeleteRolePermanent
	}

	return true, nil
}

// RestoreAllRole restores all trashed (soft-deleted) roles.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//
// Returns:
//   - bool: True if the restore was successful.
//   - error: An error if the operation failed.
func (r *roleCommandRepository) RestoreAllRole(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllRoles(ctx)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return false, errRoleDeadlineExceeded
		}
		if ctx.Err() == context.Canceled {
			return false, errRoleCancelled
		}
		return false, role_errors.ErrRestoreAllRoles
	}

	return true, nil
}

// DeleteAllRolePermanent permanently deletes all roles that are trashed.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//
// Returns:
//   - bool: True if the deletion was successful.
//   - error: An error if the operation failed.
func (r *roleCommandRepository) DeleteAllRolePermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentRoles(ctx)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return false, errRoleDeadlineExceeded
		}
		if ctx.Err() == context.Canceled {
			return false, errRoleCancelled
		}

		return false, role_errors.ErrDeleteAllRoles
	}

	return true, nil
}
