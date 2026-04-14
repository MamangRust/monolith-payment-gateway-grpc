package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	role_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/role_errors/repository"
)

// roleCommandRepository is a struct that implements the RoleCommandRepository interface
type roleCommandRepository struct {
	db *db.Queries
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
func NewRoleCommandRepository(db *db.Queries) RoleCommandRepository {
	return &roleCommandRepository{
		db: db,
	}
}

func (r *roleCommandRepository) CreateRole(ctx context.Context, req *requests.CreateRoleRequest) (*db.Role, error) {
	res, err := r.db.CreateRole(ctx, req.Name)

	if err != nil {
		return nil, role_errors.ErrCreateRole.WithInternal(err)
	}

	return res, nil
}

func (r *roleCommandRepository) UpdateRole(ctx context.Context, req *requests.UpdateRoleRequest) (*db.Role, error) {
	res, err := r.db.UpdateRole(ctx, db.UpdateRoleParams{
		RoleID:   int32(*req.ID),
		RoleName: req.Name,
	})

	if err != nil {
		return nil, role_errors.ErrUpdateRole.WithInternal(err)
	}

	return res, nil
}

func (r *roleCommandRepository) TrashedRole(ctx context.Context, id int) (*db.Role, error) {
	res, err := r.db.TrashRole(ctx, int32(id))
	if err != nil {
		return nil, role_errors.ErrTrashedRole.WithInternal(err)
	}
	return res, nil
}

func (r *roleCommandRepository) RestoreRole(ctx context.Context, id int) (*db.Role, error) {
	res, err := r.db.RestoreRole(ctx, int32(id))
	if err != nil {
		return nil, role_errors.ErrRestoreRole.WithInternal(err)
	}
	return res, nil
}

func (r *roleCommandRepository) DeleteRolePermanent(ctx context.Context, role_id int) (bool, error) {
	err := r.db.DeletePermanentRole(ctx, int32(role_id))
	if err != nil {
		return false, role_errors.ErrDeleteRolePermanent.WithInternal(err)
	}
	return true, nil
}

func (r *roleCommandRepository) RestoreAllRole(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllRoles(ctx)

	if err != nil {
		return false, role_errors.ErrRestoreAllRoles.WithInternal(err)
	}

	return true, nil
}

func (r *roleCommandRepository) DeleteAllRolePermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentRoles(ctx)

	if err != nil {
		return false, role_errors.ErrDeleteAllRoles.WithInternal(err)
	}

	return true, nil
}
