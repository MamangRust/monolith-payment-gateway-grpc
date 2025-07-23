package mencache

import (
	"context"
	"fmt"
	"time"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

var (
	keyIdentityRefreshToken = "identity:refresh_token:%s"
	keyIdentityUserInfo     = "identity:user_info:%s"
)

type identityCache struct {
	store *sharedcachehelpers.CacheStore
}

// NewidentityCache creates and returns a new instance of identityCache.
// It initializes the cache with the provided CacheStore, which is used
// for storing and retrieving cached identity-related data.
func NewidentityCache(store *sharedcachehelpers.CacheStore) *identityCache {
	return &identityCache{store: store}
}

// SetRefreshToken stores a refresh token in the cache with expiration.
//
// Parameters:
//   - ctx: the context for the caching operation
//   - token: the refresh token string
//   - expiration: the duration until the token expires
func (c *identityCache) SetRefreshToken(ctx context.Context, token string, expiration time.Duration) {
	key := keyIdentityRefreshToken
	key = fmt.Sprintf(key, token)

	sharedcachehelpers.SetToCache(ctx, c.store, key, &token, expiration)
}

// GetRefreshToken retrieves a refresh token from the cache.
//
// Parameters:
//   - ctx: the context for the caching operation
//   - token: the refresh token string to retrieve
//
// Returns:
//   - The stored token string and a boolean indicating whether it was found.
func (c *identityCache) GetRefreshToken(ctx context.Context, token string) (string, bool) {
	key := keyIdentityRefreshToken
	key = fmt.Sprintf(key, token)

	result, found := sharedcachehelpers.GetFromCache[string](ctx, c.store, key)
	if !found || result == nil {
		return "", false
	}
	return *result, true
}

// DeleteRefreshToken removes a refresh token from the cache.
//
// Parameters:
//   - ctx: the context for the caching operation
//   - token: the refresh token to delete
func (c *identityCache) DeleteRefreshToken(ctx context.Context, token string) {
	key := fmt.Sprintf(keyIdentityRefreshToken, token)
	sharedcachehelpers.DeleteFromCache(ctx, c.store, key)
}

// SetCachedUserInfo stores user information in the cache with expiration.
//
// Parameters:
//   - ctx: the context for the caching operation
//   - user: the user data to cache
//   - expiration: the duration until the cache entry expires
func (c *identityCache) SetCachedUserInfo(ctx context.Context, user *response.UserResponse, expiration time.Duration) {
	if user == nil {
		return
	}

	key := fmt.Sprintf(keyIdentityUserInfo, user.ID)

	sharedcachehelpers.SetToCache(ctx, c.store, key, user, expiration)
}

// GetCachedUserInfo retrieves cached user information by user ID.
//
// Parameters:
//   - ctx: the context for the caching operation
//   - userId: the unique identifier of the user
//
// Returns:
//   - The cached UserResponse and a boolean indicating whether it was found.
func (c *identityCache) GetCachedUserInfo(ctx context.Context, userId string) (*response.UserResponse, bool) {
	key := fmt.Sprintf(keyIdentityUserInfo, userId)

	return sharedcachehelpers.GetFromCache[response.UserResponse](ctx, c.store, key)
}

// DeleteCachedUserInfo removes cached user information from the cache.
//
// Parameters:
//   - ctx: the context for the caching operation
//   - userId: the unique identifier of the user to remove from cache
func (c *identityCache) DeleteCachedUserInfo(ctx context.Context, userId string) {
	key := fmt.Sprintf(keyIdentityUserInfo, userId)

	sharedcachehelpers.DeleteFromCache(ctx, c.store, key)
}
