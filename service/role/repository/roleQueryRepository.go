package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	role_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/role_errors/repository"
)

// roleQueryRepository is a struct that implements the RoleQueryRepository interface
type roleQueryRepository struct {
	db *db.Queries
}

func NewRoleQueryRepository(db *db.Queries) RoleQueryRepository {
	return &roleQueryRepository{
		db: db,
	}
}

func (r *roleQueryRepository) FindAllRoles(ctx context.Context, req *requests.FindAllRoles) ([]*db.GetRolesRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetRolesParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetRoles(ctx, reqDb)

	if err != nil {
		return nil, role_errors.ErrFindAllRoles
	}

	return res, nil
}

func (r *roleQueryRepository) FindById(ctx context.Context, id int) (*db.Role, error) {
	res, err := r.db.GetRole(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("role not found with ID: %d", id)
		}
		return nil, fmt.Errorf("failed to find role by ID %d: %w", id, err)
	}
	return res, nil
}

func (r *roleQueryRepository) FindByName(ctx context.Context, name string) (*db.Role, error) {
	res, err := r.db.GetRoleByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, role_errors.ErrRoleNotFound
		}

		return nil, role_errors.ErrRoleNotFound
	}
	return res, nil
}

func (r *roleQueryRepository) FindByUserId(ctx context.Context, user_id int) ([]*db.Role, error) {
	res, err := r.db.GetUserRoles(ctx, int32(user_id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, role_errors.ErrRoleNotFound
		}

		return nil, role_errors.ErrRoleNotFound
	}
	return res, nil
}

func (r *roleQueryRepository) FindByActiveRole(ctx context.Context, req *requests.FindAllRoles) ([]*db.GetActiveRolesRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveRolesParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveRoles(ctx, reqDb)

	if err != nil {
		return nil, role_errors.ErrFindActiveRoles
	}

	return res, nil
}

func (r *roleQueryRepository) FindByTrashedRole(ctx context.Context, req *requests.FindAllRoles) ([]*db.GetTrashedRolesRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedRolesParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedRoles(ctx, reqDb)

	if err != nil {
		return nil, role_errors.ErrFindTrashedRoles
	}

	return res, nil
}
