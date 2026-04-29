package adapter

import (
	"context"

	pbmerchant "github.com/MamangRust/monolith-payment-gateway-pb/merchant"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
)

type MerchantAdapter struct {
	QueryClient pbmerchant.MerchantQueryServiceClient
}

func NewMerchantAdapter(queryClient pbmerchant.MerchantQueryServiceClient) *MerchantAdapter {
	return &MerchantAdapter{
		QueryClient: queryClient,
	}
}

func (a *MerchantAdapter) FindByApiKey(ctx context.Context, api_key string) (*db.GetMerchantByApiKeyRow, error) {
	resp, err := a.QueryClient.FindByApiKey(ctx, &pbmerchant.FindByApiKeyRequest{
		ApiKey: api_key,
	})
	if err != nil {
		return nil, err
	}

	return &db.GetMerchantByApiKeyRow{
		MerchantID: resp.Data.Id,
		Name:       resp.Data.Name,
		ApiKey:     resp.Data.ApiKey,
		UserID:     resp.Data.UserId,
		Status:     resp.Data.Status,
	}, nil
}
