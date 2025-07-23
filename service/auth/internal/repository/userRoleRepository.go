package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	userrole_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_role_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/userrole"
)

// userRoleRepository is a struct that implements the UserRoleRepository interface
type userRoleRepository struct {
	db     *db.Queries
	mapper recordmapper.UserRoleRecordMapping
}

// NewUserRoleRepository creates a new UserRoleRepository instance
//
// Args:
// db: a pointer to the database queries
// mapper: a UserRoleRecordMapping object
//
// Returns:
// a pointer to the userRoleRepository struct
func NewUserRoleRepository(db *db.Queries, mapper recordmapper.UserRoleRecordMapping) *userRoleRepository {
	return &userRoleRepository{
		db:     db,
		mapper: mapper,
	}
}

// AssignRoleToUser assigns a role to a user.
//
// Parameters:
//   - ctx: the context for the database operation
//   - req: the request payload containing user ID and role ID
//
// Returns:
//   - The created UserRoleRecord if successful, or an error if the operation fails.
func (r *userRoleRepository) AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*record.UserRoleRecord, error) {
	res, err := r.db.AssignRoleToUser(ctx, db.AssignRoleToUserParams{
		UserID: int32(req.UserId),
		RoleID: int32(req.RoleId),
	})

	if err != nil {
		return nil, userrole_errors.ErrAssignRoleToUser
	}

	return r.mapper.ToUserRoleRecord(res), nil
}

// RemoveRoleFromUser removes a role assigned to a user.
//
// Parameters:
//   - ctx: the context for the database operation
//   - req: the request payload containing user ID and role ID
//
// Returns:
//   - An error if the operation fails.
func (r *userRoleRepository) RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error {
	err := r.db.RemoveRoleFromUser(ctx, db.RemoveRoleFromUserParams{
		UserID: int32(req.UserId),
		RoleID: int32(req.RoleId),
	})

	if err != nil {
		return userrole_errors.ErrRemoveRole
	}

	return nil
}
