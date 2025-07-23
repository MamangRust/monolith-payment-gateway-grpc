package transactionstatscache

import "time"

const (
	monthTopupStatusSuccessCacheKey = "transaction:month:status:success:month:%d:year:%d"
	yearTopupStatusSuccessCacheKey  = "transaction:year:status:success:year:%d"
	monthTopupStatusFailedCacheKey  = "transaction:month:status:failed:month:%d:year:%d"
	yearTopupStatusFailedCacheKey   = "transaction:year:status:failed:year:%d"

	monthTopupAmountCacheKey = "transaction:month:amount:year:%d"
	yearTopupAmountCacheKey  = "transaction:year:amount:year:%d"

	monthTopupMethodCacheKey = "transaction:month:method:year:%d"
	yearTopupMethodCacheKey  = "transaction:year:method:year:%d"

	ttlDefault = 5 * time.Minute
)
