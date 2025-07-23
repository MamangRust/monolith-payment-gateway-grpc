package repository

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
)

// RoleQueryRepository is an interface for querying role records
type RoleQueryRepository interface {
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
	FindAllRoles(ctx context.Context, req *requests.FindAllRoles) ([]*record.RoleRecord, *int, error)

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
	FindByActiveRole(ctx context.Context, req *requests.FindAllRoles) ([]*record.RoleRecord, *int, error)

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
	FindByTrashedRole(ctx context.Context, req *requests.FindAllRoles) ([]*record.RoleRecord, *int, error)

	// FindById retrieves a role by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - role_id: The ID of the role to retrieve.
	//
	// Returns:
	//   - *record.RoleRecord: The role record.
	//   - error: An error if any occurred during the query.
	FindById(ctx context.Context, role_id int) (*record.RoleRecord, error)

	// FindByName retrieves a role by its name.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - name: The name of the role.
	//
	// Returns:
	//   - *record.RoleRecord: The role record matching the name.
	//   - error: An error if any occurred during the query.
	FindByName(ctx context.Context, name string) (*record.RoleRecord, error)

	// FindByUserId retrieves all roles assigned to a specific user ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - user_id: The user ID to filter roles by.
	//
	// Returns:
	//   - []*record.RoleRecord: The list of role records assigned to the user.
	//   - error: An error if any occurred during the query.
	FindByUserId(ctx context.Context, user_id int) ([]*record.RoleRecord, error)
}

// RoleCommandRepository is an interface for creating, updating, and deleting role records
type RoleCommandRepository interface {
	// CreateRole creates a new role in the database.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The request payload containing role data to be created.
	//
	// Returns:
	//   - *record.RoleRecord: The newly created role record.
	//   - error: An error if the creation failed.
	CreateRole(ctx context.Context, request *requests.CreateRoleRequest) (*record.RoleRecord, error)

	// UpdateRole updates an existing role in the database.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The request payload containing updated role data.
	//
	// Returns:
	//   - *record.RoleRecord: The updated role record.
	//   - error: An error if the update failed.
	UpdateRole(ctx context.Context, request *requests.UpdateRoleRequest) (*record.RoleRecord, error)

	// TrashedRole performs a soft-delete (move to trash) for a role by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - role_id: The ID of the role to be trashed.
	//
	// Returns:
	//   - *record.RoleRecord: The trashed role record.
	//   - error: An error if the operation failed.
	TrashedRole(ctx context.Context, role_id int) (*record.RoleRecord, error)

	// RestoreRole restores a soft-deleted (trashed) role by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - role_id: The ID of the role to restore.
	//
	// Returns:
	//   - *record.RoleRecord: The restored role record.
	//   - error: An error if the restoration failed.
	RestoreRole(ctx context.Context, role_id int) (*record.RoleRecord, error)

	// DeleteRolePermanent permanently deletes a role by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - role_id: The ID of the role to delete permanently.
	//
	// Returns:
	//   - bool: True if deletion was successful.
	//   - error: An error if the deletion failed.
	DeleteRolePermanent(ctx context.Context, role_id int) (bool, error)

	// RestoreAllRole restores all trashed (soft-deleted) roles.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: True if the restore was successful.
	//   - error: An error if the operation failed.
	RestoreAllRole(ctx context.Context) (bool, error)

	// DeleteAllRolePermanent permanently deletes all roles that are trashed.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: True if the deletion was successful.
	//   - error: An error if the operation failed.
	DeleteAllRolePermanent(ctx context.Context) (bool, error)
}
