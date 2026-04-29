package topup_test

import (
	"context"
	"time"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	card_repo "github.com/MamangRust/monolith-payment-gateway-card/repository"
	saldo_repo "github.com/MamangRust/monolith-payment-gateway-saldo/repository"
)

type realCardRepo struct {
	query card_repo.CardQueryRepository
	cmd   card_repo.CardCommandRepository
}
func (r *realCardRepo) FindUserCardByCardNumber(ctx context.Context, card_number string) (*db.GetUserEmailByCardNumberRow, error) {
	return r.query.FindUserCardByCardNumber(ctx, card_number)
}
func (r *realCardRepo) FindCardByCardNumber(ctx context.Context, card_number string) (*db.GetCardByCardNumberRow, error) {
	return r.query.FindCardByCardNumber(ctx, card_number)
}
func (r *realCardRepo) UpdateCard(ctx context.Context, request *requests.UpdateCardRequest) (*db.UpdateCardRow, error) {
	return r.cmd.UpdateCard(ctx, request)
}

type realSaldoRepo struct {
	repo saldo_repo.Repositories
}
func (r *realSaldoRepo) FindByCardNumber(ctx context.Context, card_number string) (*db.Saldo, error) {
	return r.repo.FindByCardNumber(ctx, card_number)
}
func (r *realSaldoRepo) UpdateSaldoBalance(ctx context.Context, request *requests.UpdateSaldoBalance) (*db.UpdateSaldoBalanceRow, error) {
	return r.repo.UpdateSaldoBalance(ctx, request)
}

// Dummy cache metrics
type dummyCacheMetrics struct{}
func (d *dummyCacheMetrics) RecordCacheHit(ctx context.Context, key string) {}
func (d *dummyCacheMetrics) RecordCacheMiss(ctx context.Context, key string) {}
func (d *dummyCacheMetrics) RecordCacheSet(ctx context.Context, key string, success bool) {}
func (d *dummyCacheMetrics) RecordCacheDelete(ctx context.Context, key string, success bool) {}
func (d *dummyCacheMetrics) RecordCacheOperationLatency(ctx context.Context, operation string, duration time.Duration) {}
func (d *dummyCacheMetrics) RecordCacheError(ctx context.Context, operation, key string, err error) {}
