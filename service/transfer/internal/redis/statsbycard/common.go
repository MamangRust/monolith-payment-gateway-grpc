package transferstatsbycardcache

import "time"

// Constants for cache keys
const (
	transferMonthTransferStatusSuccessByCardKey = "transfer:month:transfer_status:success:card_number:%s:month:%d:year:%d"
	transferMonthTransferStatusFailedByCardKey  = "transfer:month:transfer_status:failed:card_number:%s:month:%d:year:%d"

	transferYearTransferStatusSuccessByCardKey = "transfer:year:transfer_status:success:card_number:%s:year:%d"
	transferYearTransferStatusFailedByCardKey  = "transfer:year:transfer_status:failed:card_number:%s:year:%d"

	transferMonthTransferAmountByCardKey = "transfer:month:transfer_amount:card_number:%s:year:%d"

	transferYearTransferAmountByCardKey = "transfer:year:transfer_amount:card_number:%s:year:%d"

	ttlDefault = 5 * time.Minute
)
