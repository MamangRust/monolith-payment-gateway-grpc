package mencache

import (
	"context"
	"fmt"
	"time"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// Constants for cache keys
const (
	userAllCacheKey     = "user:all:page:%d:pageSize:%d:search:%s"
	userByIdCacheKey    = "user:id:%d"
	userActiveCacheKey  = "user:active:page:%d:pageSize:%d:search:%s"
	userTrashedCacheKey = "user:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

// UserResponse represents the structure of the cached user data.
type userCacheResponse struct {
	Data         []*response.UserResponse `json:"data"`
	TotalRecords *int                     `json:"total_records"`
}

// UserCachedResponseDeleteAt represents the structure of the cached user data.
type userCacheResponseDeleteAt struct {
	Data         []*response.UserResponseDeleteAt `json:"data"`
	TotalRecords *int                             `json:"total_records"`
}

// userQueryCache represents the cache for user queries.
type userQueryCache struct {
	store *sharedcachehelpers.CacheStore
}

// NewUserQueryCache creates a new user query cache using the provided cache store.
func NewUserQueryCache(store *sharedcachehelpers.CacheStore) UserQueryCache {
	return &userQueryCache{store: store}
}

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
func (s *userQueryCache) GetCachedUsersCache(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponse, *int, bool) {
	key := fmt.Sprintf(userAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[userCacheResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

// SetCachedUsersCache sets the cached list of users.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing filter parameters.
//   - data: The list of users to be cached.
//   - total: The total count of users.
func (s *userQueryCache) SetCachedUsersCache(ctx context.Context, req *requests.FindAllUsers, data []*response.UserResponse, total *int) {
	if total == nil {
		zero := 0

		total = &zero
	}

	if data == nil {
		data = []*response.UserResponse{}
	}

	key := fmt.Sprintf(userAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &userCacheResponse{Data: data, TotalRecords: total}
	sharedcachehelpers.SetToCache(ctx, s.store, key, payload, ttlDefault)

}

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
func (s *userQueryCache) GetCachedUserActiveCache(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(userActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[userCacheResponseDeleteAt](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

// SetCachedUserActiveCache sets the cached active users.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request with filter parameters.
//   - data: The list of active users.
//   - total: The total count of active users.
func (s *userQueryCache) SetCachedUserActiveCache(ctx context.Context, req *requests.FindAllUsers, data []*response.UserResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*response.UserResponseDeleteAt{}
	}

	key := fmt.Sprintf(userActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &userCacheResponseDeleteAt{Data: data, TotalRecords: total}
	sharedcachehelpers.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

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
func (s *userQueryCache) GetCachedUserTrashedCache(ctx context.Context, req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, bool) {
	key := fmt.Sprintf(userTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := sharedcachehelpers.GetFromCache[userCacheResponseDeleteAt](ctx, s.store, key)

	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

// SetCachedUserTrashedCache sets the cached trashed (soft deleted) users.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request with filter parameters.
//   - data: The list of trashed users.
//   - total: The total count of trashed users.
func (s *userQueryCache) SetCachedUserTrashedCache(ctx context.Context, req *requests.FindAllUsers, data []*response.UserResponseDeleteAt, total *int) {
	if total == nil {
		zero := 0

		total = &zero
	}

	if data == nil {
		return
	}

	key := fmt.Sprintf(userTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &userCacheResponseDeleteAt{Data: data, TotalRecords: total}
	sharedcachehelpers.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

// GetCachedUserCache retrieves cached user by ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - id: The user ID to retrieve from cache.
//
// Returns:
//   - *response.UserResponse: The cached user response.
//   - bool: Whether the cache was found.
func (s *userQueryCache) GetCachedUserCache(ctx context.Context, id int) (*response.UserResponse, bool) {
	key := fmt.Sprintf(userByIdCacheKey, id)

	result, found := sharedcachehelpers.GetFromCache[*response.UserResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetCachedUserCache sets the cached data for a user by ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - data: The user data to be cached.
func (s *userQueryCache) SetCachedUserCache(ctx context.Context, data *response.UserResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(userByIdCacheKey, data.ID)
	sharedcachehelpers.SetToCache(ctx, s.store, key, data, ttlDefault)
}
