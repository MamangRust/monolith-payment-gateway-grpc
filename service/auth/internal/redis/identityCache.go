package mencache

import (
	"fmt"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

var (
	keyIdentityRefreshToken = "identity:refresh_token:%s"
	keyIdentityUserInfo     = "identity:user_info:%s"
)

type identityCache struct {
	store *CacheStore
}

func NewidentityCache(store *CacheStore) *identityCache {
	return &identityCache{store: store}
}

func (c *identityCache) SetRefreshToken(token string, expiration time.Duration) {
	key := keyIdentityRefreshToken
	key = fmt.Sprintf(key, token)

	SetToCache(c.store, key, &token, expiration)
}

func (c *identityCache) GetRefreshToken(token string) (string, bool) {
	key := keyIdentityRefreshToken
	key = fmt.Sprintf(key, token)

	result, found := GetFromCache[string](c.store, key)
	if !found {
		return "", false
	}
	return *result, true
}

func (c *identityCache) DeleteRefreshToken(token string) {
	key := fmt.Sprintf(keyIdentityRefreshToken, token)
	DeleteFromCache(c.store, key)
}

func (c *identityCache) SetCachedUserInfo(user *response.UserResponse, expiration time.Duration) {
	key := fmt.Sprintf(keyIdentityUserInfo, user.ID)

	SetToCache(c.store, key, user, expiration)
}

func (c *identityCache) GetCachedUserInfo(userId string) (*response.UserResponse, bool) {
	key := fmt.Sprintf(keyIdentityUserInfo, userId)

	return GetFromCache[response.UserResponse](c.store, key)
}

func (c *identityCache) DeleteCachedUserInfo(userId string) {
	key := fmt.Sprintf(keyIdentityUserInfo, userId)

	DeleteFromCache(c.store, key)
}
