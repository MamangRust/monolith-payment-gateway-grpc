package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// RoleQueryService is an interface for querying role records
type RoleQueryService interface {
	// FindAll retrieves all roles based on the given filters.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing pagination and search filters.
	//
	// Returns:
	//   - []*response.RoleResponse: The list of role responses.
	//   - *int: The total number of roles.
	//   - *response.ErrorResponse: An error response if any occurred.
	FindAll(ctx context.Context, req *requests.FindAllRoles) ([]*response.RoleResponse, *int, *response.ErrorResponse)

	// FindByActiveRole retrieves only active (non-deleted) roles based on the given filters.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing filters.
	//
	// Returns:
	//   - []*response.RoleResponseDeleteAt: The list of active role responses.
	//   - *int: The total number of active roles.
	//   - *response.ErrorResponse: An error response if any occurred.
	FindByActiveRole(ctx context.Context, req *requests.FindAllRoles) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse)

	// FindByTrashedRole retrieves only trashed (soft-deleted) roles based on the given filters.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing filters.
	//
	// Returns:
	//   - []*response.RoleResponseDeleteAt: The list of trashed role responses.
	//   - *int: The total number of trashed roles.
	//   - *response.ErrorResponse: An error response if any occurred.
	FindByTrashedRole(ctx context.Context, req *requests.FindAllRoles) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse)

	// FindById retrieves a role by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - role_id: The ID of the role.
	//
	// Returns:
	//   - *response.RoleResponse: The role response.
	//   - *response.ErrorResponse: An error response if any occurred.
	FindById(ctx context.Context, role_id int) (*response.RoleResponse, *response.ErrorResponse)

	// FindByUserId retrieves all roles assigned to a specific user.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - id: The ID of the user.
	//
	// Returns:
	//   - []*response.RoleResponse: The list of roles assigned to the user.
	//   - *response.ErrorResponse: An error response if any occurred.
	FindByUserId(ctx context.Context, id int) ([]*response.RoleResponse, *response.ErrorResponse)
}

// RoleCommandService is an interface for creating, updating, and deleting role records
type RoleCommandService interface {
	// CreateRole creates a new role based on the given request.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The request payload containing the role details.
	//
	// Returns:
	//   - *response.RoleResponse: The created role.
	//   - *response.ErrorResponse: An error response if creation failed.
	CreateRole(ctx context.Context, request *requests.CreateRoleRequest) (*response.RoleResponse, *response.ErrorResponse)

	// UpdateRole updates an existing role based on the given request.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - request: The request payload containing updated role details.
	//
	// Returns:
	//   - *response.RoleResponse: The updated role.
	//   - *response.ErrorResponse: An error response if update failed.
	UpdateRole(ctx context.Context, request *requests.UpdateRoleRequest) (*response.RoleResponse, *response.ErrorResponse)

	// TrashedRole soft-deletes (moves to trash) a role by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - role_id: The ID of the role to be trashed.
	//
	// Returns:
	//   - *response.RoleResponse: The trashed role.
	//   - *response.ErrorResponse: An error response if the operation failed.
	TrashedRole(ctx context.Context, role_id int) (*response.RoleResponseDeleteAt, *response.ErrorResponse)

	// RestoreRole restores a soft-deleted role by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - role_id: The ID of the role to be restored.
	//
	// Returns:
	//   - *response.RoleResponse: The restored role.
	//   - *response.ErrorResponse: An error response if the restoration failed.
	RestoreRole(ctx context.Context, role_id int) (*response.RoleResponse, *response.ErrorResponse)

	// DeleteRolePermanent permanently deletes a role by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - role_id: The ID of the role to be permanently deleted.
	//
	// Returns:
	//   - bool: True if deletion was successful.
	//   - *response.ErrorResponse: An error response if the deletion failed.
	DeleteRolePermanent(ctx context.Context, role_id int) (bool, *response.ErrorResponse)

	// RestoreAllRole restores all trashed roles.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: True if restoration was successful.
	//   - *response.ErrorResponse: An error response if the operation failed.
	RestoreAllRole(ctx context.Context) (bool, *response.ErrorResponse)

	// DeleteAllRolePermanent permanently deletes all trashed roles.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//
	// Returns:
	//   - bool: True if deletion was successful.
	//   - *response.ErrorResponse: An error response if the operation failed.
	DeleteAllRolePermanent(ctx context.Context) (bool, *response.ErrorResponse)
}
