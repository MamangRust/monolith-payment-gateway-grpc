package mencache

import (
	"context"
	"fmt"
	"time"

	sharedcachehelpers "github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

var keylogin = "auth:login:%s"

type loginCache struct {
	store *sharedcachehelpers.CacheStore
}

// NewLoginCache creates and returns a new instance of loginCache.
// It initializes the cache with the provided CacheStore, which is used
// for storing and retrieving cached login-related data.
func NewLoginCache(store *sharedcachehelpers.CacheStore) *loginCache {
	return &loginCache{store: store}
}

// GetCachedLogin retrieves a cached token response by email.
//
// Parameters:
//   - ctx: the context for the caching operation
//   - email: the email address for which to retrieve the cached login data
//
// Returns:
//   - The cached TokenResponse and a boolean indicating whether it was found.
func (s *loginCache) GetCachedLogin(ctx context.Context, email string) (*response.TokenResponse, bool) {
	key := fmt.Sprintf(keylogin, email)

	result, found := sharedcachehelpers.GetFromCache[*response.TokenResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

// SetCachedLogin stores login token response in cache with expiration.
//
// Parameters:
//   - ctx: the context for the caching operation
//   - email: the email address used for login
//   - data: the token response to cache
//   - expiration: the duration until the cache entry expires
func (s *loginCache) SetCachedLogin(ctx context.Context, email string, data *response.TokenResponse, expiration time.Duration) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(keylogin, email)

	sharedcachehelpers.SetToCache(ctx, s.store, key, data, expiration)
}
