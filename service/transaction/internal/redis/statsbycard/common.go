package transactionstatsbycarcache

import "time"

const (
	monthTopupStatusSuccessByCardCacheKey = "transaction:bycard:month:status:success:card:%s:month:%d:year:%d"
	yearTopupStatusSuccessByCardCacheKey  = "transaction:bycard:year:status:success:card:%s:year:%d"
	monthTopupStatusFailedByCardCacheKey  = "transaction:bycard:month:status:failed:card:%s:month:%d:year:%d"
	yearTopupStatusFailedByCardCacheKey   = "transaction:bycard:year:status:failed:card:%s:year:%d"

	monthTopupAmountByCardCacheKey = "transaction:bycard:month:amount:card:%s:year:%d"
	yearTopupAmountByCardCacheKey  = "transaction:bycard:year:amount:card:%s:year:%d"

	monthTopupMethodByCardCacheKey = "transaction:bycard:month:method:card:%s:year:%d"
	yearTopupMethodByCardCacheKey  = "transaction:bycard:year:method:card:%s:year:%d"

	ttlDefault = 5 * time.Minute
)
