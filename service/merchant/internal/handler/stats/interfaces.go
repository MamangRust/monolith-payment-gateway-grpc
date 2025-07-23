package merchantstatshandler

import (
	pbmerchant "github.com/MamangRust/monolith-payment-gateway-pb/merchant"
)

type MerchantStatsAmountHandleGrpc interface {
	pbmerchant.MerchantStatsAmountServiceServer
}

type MerchantStatsMethodHandleGrpc interface {
	pbmerchant.MerchantStatsMethodServiceServer
}

type MerchantStatsTotalAmountHandleGrpc interface {
	pbmerchant.MerchantStatsTotalAmountServiceServer
}
