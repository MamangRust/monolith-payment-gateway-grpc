package mencache

import (
	"context"
	"fmt"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// roleCachedResponse is a struct that represents the cached response
type roleCachedResponse struct {
	Data         []*response.RoleResponse `json:"data"`
	TotalRecords *int                     `json:"total_records"`
}

// roleCachedResponseDeleteAt is a struct that represents the cached response
type roleCachedResponseDeleteAt struct {
	Data         []*response.RoleResponseDeleteAt `json:"data"`
	TotalRecords *int                             `json:"total_records"`
}

// roleQueryCache is a struct that represents the cache store
type roleQueryCache struct {
	store *sharedcachehelpers.CacheStore
}

// NewRoleQueryCache creates a new roleQueryCache instance.
//
// Parameters:
//   - store: The cache store to use for caching.
//
// Returns:
//   - *roleQueryCache: The newly created roleQueryCache instance.
func NewRoleQueryCache(store *sharedcachehelpers.CacheStore) *roleQueryCache {
	return &roleQueryCache{store: store}
}

// SetCachedRoles stores a list of roles and their total count in cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request object used to fetch the data.
//   - data: The list of role responses to cache.
//   - total: The total number of roles.
func (m *roleQueryCache) SetCachedRoles(ctx context.Context, req *requests.FindAllRoles, data []*response.RoleResponse, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.RoleResponse{}
	}

	key := fmt.Sprintf(roleAllCacheKey, req.Page, req.PageSize, req.Search)

	payload := &roleCachedResponse{Data: data, TotalRecords: total}
	sharedcachehelpers.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

// SetCachedRoleById stores a single role by ID in cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - id: The role ID.
//   - data: The role response to store in cache.
func (m *roleQueryCache) SetCachedRoleById(ctx context.Context, id int, data *response.RoleResponse) {
	if data == nil {
		data = &response.RoleResponse{}
	}

	key := fmt.Sprintf(roleByIdCacheKey, id)
	sharedcachehelpers.SetToCache(ctx, m.store, key, data, ttlDefault)
}

// SetCachedRoleByUserId stores roles associated with a user ID in cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - userId: The user ID.
//   - data: The list of roles associated with the user.
func (m *roleQueryCache) SetCachedRoleByUserId(ctx context.Context, userId int, data []*response.RoleResponse) {
	if data == nil {
		data = []*response.RoleResponse{}
	}

	key := fmt.Sprintf(roleByIdCacheKey, userId)

	sharedcachehelpers.SetToCache(ctx, m.store, key, &data, ttlDefault)
}

// SetCachedRoleActive stores a list of active roles in cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request used for filtering.
//   - data: The list of active role responses.
//   - total: The total number of active roles.
func (m *roleQueryCache) SetCachedRoleActive(ctx context.Context, req *requests.FindAllRoles, data []*response.RoleResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.RoleResponseDeleteAt{}
	}

	key := fmt.Sprintf(roleActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload := &roleCachedResponseDeleteAt{Data: data, TotalRecords: total}

	sharedcachehelpers.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

// SetCachedRoleTrashed stores a list of trashed roles in cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request used for filtering.
//   - data: The list of trashed role responses.
//   - total: The total number of trashed roles.
func (m *roleQueryCache) SetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRoles, data []*response.RoleResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.RoleResponseDeleteAt{}
	}

	key := fmt.Sprintf(roleTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &roleCachedResponseDeleteAt{Data: data, TotalRecords: total}
	sharedcachehelpers.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

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
func (m *roleQueryCache) GetCachedRoles(ctx context.Context, req *requests.FindAllRoles) ([]*response.RoleResponse, *int, bool) {
	key := fmt.Sprintf(roleAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[roleCachedResponse](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

// GetCachedRoleById retrieves a role by its ID from cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - id: The ID of the role.
//
// Returns:
//   - *response.RoleResponse: The cached role.
//   - bool: Whether the cache was found.
func (m *roleQueryCache) GetCachedRoleById(ctx context.Context, id int) (*response.RoleResponse, bool) {
	key := fmt.Sprintf(roleByIdCacheKey, id)

	result, found := sharedcachehelpers.GetFromCache[*response.RoleResponse](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// GetCachedRoleByUserId retrieves roles associated with a user ID from cache.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - userId: The user ID to search by.
//
// Returns:
//   - []*response.RoleResponse: The cached roles.
//   - bool: Whether the cache was found.
func (m *roleQueryCache) GetCachedRoleByUserId(ctx context.Context, userId int) ([]*response.RoleResponse, bool) {
	key := fmt.Sprintf(roleByIdCacheKey, userId)

	result, found := sharedcachehelpers.GetFromCache[[]*response.RoleResponse](ctx, m.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

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
func (m *roleQueryCache) GetCachedRoleActive(ctx context.Context, req *requests.FindAllRoles) ([]*response.RoleResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(roleActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[roleCachedResponseDeleteAt](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

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
func (m *roleQueryCache) GetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRoles) ([]*response.RoleResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(roleTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[roleCachedResponseDeleteAt](ctx, m.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}
