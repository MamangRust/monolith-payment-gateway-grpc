package repository

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	role_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/role_errors/repository"
	rolerecordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/role"
)

// roleQueryRepository is a struct that implements the RoleQueryRepository interface
type roleQueryRepository struct {
	db     *db.Queries
	mapper rolerecordmapper.RoleQueryRecordMapping
}

// NewRoleQueryRepository creates a new instance of roleQueryRepository with the provided
// database queries, context, and role record mapper. This repository is responsible for executing
// query operations related to role records in the database.
//
// Parameters:
//   - db: A pointer to the db.Queries object for executing database queries.
//   - mapper: A RoleRecordMapping that provides methods to map database rows to Role domain models.
//
// Returns:
//   - A pointer to the newly created roleQueryRepository instance.
func NewRoleQueryRepository(db *db.Queries, mapper rolerecordmapper.RoleQueryRecordMapping) RoleQueryRepository {
	return &roleQueryRepository{
		db:     db,
		mapper: mapper,
	}
}

// FindAllRoles retrieves all roles based on the given filter parameters.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request object containing pagination and search filters.
//
// Returns:
//   - []*record.RoleRecord: The list of role records.
//   - *int: The total number of roles.
//   - error: An error if any occurred during the query.
func (r *roleQueryRepository) FindAllRoles(ctx context.Context, req *requests.FindAllRoles) ([]*record.RoleRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetRolesParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetRoles(ctx, reqDb)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, nil, errRoleDeadlineExceeded
		}
		if ctx.Err() == context.Canceled {
			return nil, nil, errRoleCancelled
		}

		return nil, nil, role_errors.ErrFindAllRoles
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToRolesRecordAll(res), &totalCount, nil
}

// FindById retrieves a role by its ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - role_id: The ID of the role to retrieve.
//
// Returns:
//   - *record.RoleRecord: The role record.
//   - error: An error if any occurred during the query.
func (r *roleQueryRepository) FindById(ctx context.Context, id int) (*record.RoleRecord, error) {
	res, err := r.db.GetRole(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, role_errors.ErrRoleNotFound
		}

		if ctx.Err() == context.DeadlineExceeded {
			return nil, errRoleDeadlineExceeded
		}
		if ctx.Err() == context.Canceled {
			return nil, errRoleCancelled
		}

		return nil, role_errors.ErrRoleNotFound
	}

	return r.mapper.ToRoleRecord(res), nil
}

// FindByName retrieves a role by its name.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - name: The name of the role.
//
// Returns:
//   - *record.RoleRecord: The role record matching the name.
//   - error: An error if any occurred during the query.
func (r *roleQueryRepository) FindByName(ctx context.Context, name string) (*record.RoleRecord, error) {

	res, err := r.db.GetRoleByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, role_errors.ErrRoleNotFound
		}

		if ctx.Err() == context.DeadlineExceeded {
			return nil, errRoleDeadlineExceeded
		}
		if ctx.Err() == context.Canceled {
			return nil, errRoleCancelled
		}

		return nil, role_errors.ErrRoleNotFound
	}
	return r.mapper.ToRoleRecord(res), nil
}

// FindByUserId retrieves all roles assigned to a specific user ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - user_id: The user ID to filter roles by.
//
// Returns:
//   - []*record.RoleRecord: The list of role records assigned to the user.
//   - error: An error if any occurred during the query.
func (r *roleQueryRepository) FindByUserId(ctx context.Context, user_id int) ([]*record.RoleRecord, error) {
	res, err := r.db.GetUserRoles(ctx, int32(user_id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, role_errors.ErrRoleNotFound
		}

		if ctx.Err() == context.DeadlineExceeded {
			return nil, errRoleDeadlineExceeded
		}
		if ctx.Err() == context.Canceled {
			return nil, errRoleCancelled
		}

		return nil, role_errors.ErrRoleNotFound
	}
	return r.mapper.ToRolesRecord(res), nil
}

// FindByActiveRole retrieves only active roles (non-deleted) based on the given filters.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request object containing pagination and filters.
//
// Returns:
//   - []*record.RoleRecord: The list of active role records.
//   - *int: The total number of active roles.
//   - error: An error if any occurred during the query.
func (r *roleQueryRepository) FindByActiveRole(ctx context.Context, req *requests.FindAllRoles) ([]*record.RoleRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveRolesParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveRoles(ctx, reqDb)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, nil, errRoleDeadlineExceeded
		}
		if ctx.Err() == context.Canceled {
			return nil, nil, errRoleCancelled
		}

		return nil, nil, role_errors.ErrFindActiveRoles
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToRolesRecordActive(res), &totalCount, nil
}

// FindByTrashedRole retrieves only trashed (soft-deleted) roles based on the given filters.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request object containing pagination and filters.
//
// Returns:
//   - []*record.RoleRecord: The list of trashed role records.
//   - *int: The total number of trashed roles.
//   - error: An error if any occurred during the query.
func (r *roleQueryRepository) FindByTrashedRole(ctx context.Context, req *requests.FindAllRoles) ([]*record.RoleRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedRolesParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedRoles(ctx, reqDb)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, nil, errRoleDeadlineExceeded
		}
		if ctx.Err() == context.Canceled {
			return nil, nil, errRoleCancelled
		}

		return nil, nil, role_errors.ErrFindTrashedRoles
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToRolesRecordTrashed(res), &totalCount, nil
}
