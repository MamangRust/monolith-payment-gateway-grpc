package mencache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// RoleCommandCache is an interface that represents the cache store
type RoleCommandCache interface {
	// DeleteCachedRole removes a cached role by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - id: The ID of the role to delete from cache.
	DeleteCachedRole(ctx context.Context, id int)
}

// RoleQueryCache is an interface that represents the cache store
type RoleQueryCache interface {
	// SetCachedRoles stores a list of roles and their total count in cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object used to fetch the data.
	//   - data: The list of role responses to cache.
	//   - total: The total number of roles.
	SetCachedRoles(ctx context.Context, req *requests.FindAllRoles, data []*response.RoleResponse, total *int)

	// SetCachedRoleById stores a single role by ID in cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - id: The role ID.
	//   - data: The role response to store in cache.
	SetCachedRoleById(ctx context.Context, id int, data *response.RoleResponse)

	// SetCachedRoleByUserId stores roles associated with a user ID in cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - userId: The user ID.
	//   - data: The list of roles associated with the user.
	SetCachedRoleByUserId(ctx context.Context, userId int, data []*response.RoleResponse)

	// SetCachedRoleActive stores a list of active roles in cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request used for filtering.
	//   - data: The list of active role responses.
	//   - total: The total number of active roles.
	SetCachedRoleActive(ctx context.Context, req *requests.FindAllRoles, data []*response.RoleResponseDeleteAt, total *int)

	// SetCachedRoleTrashed stores a list of trashed roles in cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request used for filtering.
	//   - data: The list of trashed role responses.
	//   - total: The total number of trashed roles.
	SetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRoles, data []*response.RoleResponseDeleteAt, total *int)

	// GetCachedRoles retrieves a cached list of roles if available.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object containing filters for roles.
	//
	// Returns:
	//   - []*response.RoleResponse: The cached role responses.
	//   - *int: The total count of cached roles.
	//   - bool: Whether the cache was found.
	GetCachedRoles(ctx context.Context, req *requests.FindAllRoles) ([]*response.RoleResponse, *int, bool)

	// GetCachedRoleByUserId retrieves roles associated with a user ID from cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - userId: The user ID to search by.
	//
	// Returns:
	//   - []*response.RoleResponse: The cached roles.
	//   - bool: Whether the cache was found.
	GetCachedRoleByUserId(ctx context.Context, userId int) ([]*response.RoleResponse, bool)

	// GetCachedRoleById retrieves a role by its ID from cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - id: The ID of the role.
	//
	// Returns:
	//   - *response.RoleResponse: The cached role.
	//   - bool: Whether the cache was found.
	GetCachedRoleById(ctx context.Context, id int) (*response.RoleResponse, bool)

	// GetCachedRoleActive retrieves active roles from cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object with filter criteria.
	//
	// Returns:
	//   - []*response.RoleResponseDeleteAt: The list of active roles.
	//   - *int: The total number of records.
	//   - bool: Whether the cache was found.
	GetCachedRoleActive(ctx context.Context, req *requests.FindAllRoles) ([]*response.RoleResponseDeleteAt, *int, bool)

	// GetCachedRoleTrashed retrieves trashed roles from cache.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request object with filter criteria.
	//
	// Returns:
	//   - []*response.RoleResponseDeleteAt: The list of trashed roles.
	//   - *int: The total number of records.
	//   - bool: Whether the cache was found.
	GetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRoles) ([]*response.RoleResponseDeleteAt, *int, bool)
}
