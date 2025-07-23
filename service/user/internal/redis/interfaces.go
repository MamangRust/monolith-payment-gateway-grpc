package mencache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// UserQueryCache is an interface for caching user queries.

// UserQueryCache handles caching operations related to user queries.
type UserQueryCache interface {
	// GetCachedUsersCache retrieves cached list of users based on filter.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing filter parameters.
	//
	// Returns:
	//   - []*response.UserResponse: List of user responses.
	//   - *int: Total number of users.
	//   - bool: Whether the cache was found.
	GetCachedUsersCache(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponse, *int, bool)

	// SetCachedUsersCache sets the cached list of users.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request containing filter parameters.
	//   - data: The list of users to be cached.
	//   - total: The total count of users.
	SetCachedUsersCache(ctx context.Context, req *requests.FindAllUsers, data []*response.UserResponse, total *int)

	// GetCachedUserActiveCache retrieves cached active users.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request with filter parameters.
	//
	// Returns:
	//   - []*response.UserResponseDeleteAt: List of active users.
	//   - *int: Total count of active users.
	//   - bool: Whether the cache was found.
	GetCachedUserActiveCache(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, bool)

	// SetCachedUserActiveCache sets the cached active users.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request with filter parameters.``
	//   - data: The list of active users.
	//   - total: The total count of active users.
	SetCachedUserActiveCache(ctx context.Context, req *requests.FindAllUsers, data []*response.UserResponseDeleteAt, total *int)

	// GetCachedUserTrashedCache retrieves cached trashed (soft deleted) users.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request with filter parameters.
	//
	// Returns:
	//   - []*response.UserResponseDeleteAt: List of trashed users.
	//   - *int: Total count of trashed users.
	//   - bool: Whether the cache was found.
	GetCachedUserTrashedCache(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, bool)

	// SetCachedUserTrashedCache sets the cached trashed (soft deleted) users.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - req: The request with filter parameters.
	//   - data: The list of trashed users.
	//   - total: The total count of trashed users.
	SetCachedUserTrashedCache(ctx context.Context, req *requests.FindAllUsers, data []*response.UserResponseDeleteAt, total *int)

	// GetCachedUserCache retrieves cached user by ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - id: The user ID to retrieve from cache.
	//
	// Returns:
	//   - *response.UserResponse: The cached user response.
	//   - bool: Whether the cache was found.
	GetCachedUserCache(ctx context.Context, id int) (*response.UserResponse, bool)

	// SetCachedUserCache sets the cached data for a user by ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - data: The user data to be cached.
	SetCachedUserCache(ctx context.Context, data *response.UserResponse)
}

// UserCommandCache is an interface for user command cache operations.
type UserCommandCache interface {
	// DeleteUserCache removes a cached user by its ID.
	//
	// Parameters:
	//   - ctx: The context for timeout and cancellation.
	//   - id: The user ID.
	DeleteUserCache(ctx context.Context, id int)
}
