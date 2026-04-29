package transaction_test

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	card_repo "github.com/MamangRust/monolith-payment-gateway-card/repository"
	saldo_repo "github.com/MamangRust/monolith-payment-gateway-saldo/repository"
	merchant_repo "github.com/MamangRust/monolith-payment-gateway-merchant/repository"
)

// We use the real repository implementations from other services to ensure full integration
type realMerchantRepo struct {
	repo merchant_repo.MerchantQueryRepository
}
func (r *realMerchantRepo) FindByApiKey(ctx context.Context, api_key string) (*db.GetMerchantByApiKeyRow, error) {
	return r.repo.FindByApiKey(ctx, api_key)
}

// Rename to match old tests requirement
type transactionCardRepo struct {
	query   card_repo.CardQueryRepository
	command card_repo.CardCommandRepository
}
func (r *transactionCardRepo) FindCardByUserId(ctx context.Context, user_id int) (*db.GetCardByUserIDRow, error) {
	return r.query.FindCardByUserId(ctx, user_id)
}
func (r *transactionCardRepo) FindUserCardByCardNumber(ctx context.Context, card_number string) (*db.GetUserEmailByCardNumberRow, error) {
	return r.query.FindUserCardByCardNumber(ctx, card_number)
}
func (r *transactionCardRepo) FindCardByCardNumber(ctx context.Context, card_number string) (*db.GetCardByCardNumberRow, error) {
	return r.query.FindCardByCardNumber(ctx, card_number)
}
func (r *transactionCardRepo) UpdateCard(ctx context.Context, request *requests.UpdateCardRequest) (*db.UpdateCardRow, error) {
	return r.command.UpdateCard(ctx, request)
}

// Keep realCardRepo as alias or just use the same struct
type realCardRepo = transactionCardRepo

type realSaldoRepo struct {
	repo saldo_repo.Repositories
}
func (r *realSaldoRepo) FindByCardNumber(ctx context.Context, card_number string) (*db.Saldo, error) {
	return r.repo.FindByCardNumber(ctx, card_number)
}
func (r *realSaldoRepo) UpdateSaldoBalance(ctx context.Context, request *requests.UpdateSaldoBalance) (*db.UpdateSaldoBalanceRow, error) {
	return r.repo.UpdateSaldoBalance(ctx, request)
}

type dummyCacheMetrics struct{}
func (d *dummyCacheMetrics) RecordCacheHit(ctx context.Context, key string) {}
func (d *dummyCacheMetrics) RecordCacheMiss(ctx context.Context, key string) {}
func (d *dummyCacheMetrics) RecordCacheSet(ctx context.Context, key string, success bool) {}
func (d *dummyCacheMetrics) RecordCacheDelete(ctx context.Context, key string, success bool) {}
func (d *dummyCacheMetrics) RecordCacheOperationLatency(ctx context.Context, operation string, duration time.Duration) {}
func (d *dummyCacheMetrics) RecordCacheError(ctx context.Context, operation, key string, err error) {}
