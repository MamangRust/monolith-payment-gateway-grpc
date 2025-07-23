package carddashboardmencache

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// CardDashboardTotalCache handles caching for aggregated (global) card dashboard statistics.
type CardDashboardTotalCache interface {
	// GetDashboardCardCache retrieves the cached aggregated dashboard data for all cards.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//
	// Returns:
	//   - A pointer to DashboardCard if found in cache, or false if not present.
	GetDashboardCardCache(ctx context.Context) (*response.DashboardCard, bool)

	// SetDashboardCardCache caches the aggregated dashboard data for all cards.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - data: the aggregated dashboard data to cache
	SetDashboardCardCache(ctx context.Context, data *response.DashboardCard)

	// DeleteDashboardCardCache removes the aggregated dashboard cache entry.
	//
	// Parameters:
	//   - ctx: the context for the operation
	DeleteDashboardCardCache(ctx context.Context)
}

// CardDashboardByCardNumberCache handles caching of dashboard data for individual cards by card number.
type CardDashboardByCardNumberCache interface {
	// GetDashboardCardCardNumberCache retrieves cached dashboard data for a specific card number.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - cardNumber: the specific card number to look up
	//
	// Returns:
	//   - A pointer to DashboardCardCardNumber if found, or false if not present.
	GetDashboardCardCardNumberCache(ctx context.Context, cardNumber string) (*response.DashboardCardCardNumber, bool)

	// SetDashboardCardCardNumberCache caches dashboard data for a specific card number.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - cardNumber: the card number to associate with the cached data
	//   - data: the dashboard data to cache
	SetDashboardCardCardNumberCache(ctx context.Context, cardNumber string, data *response.DashboardCardCardNumber)

	// DeleteDashboardCardCardNumberCache removes cached dashboard data for a specific card number.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - cardNumber: the card number whose cache should be cleared
	DeleteDashboardCardCardNumberCache(ctx context.Context, cardNumber string)
}
