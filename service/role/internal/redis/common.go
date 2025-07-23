package mencache

import "time"

// roleQueryCache is a struct that represents the cache store
const (
	roleAllCacheKey     = "role:all:page:%d:pageSize:%d:search:%s"
	roleByIdCacheKey    = "role:id:%d"
	roleActiveCacheKey  = "role:active:page:%d:pageSize:%d:search:%s"
	roleTrashedCacheKey = "role:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)
