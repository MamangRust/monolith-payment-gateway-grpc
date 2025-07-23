package withdrawstatsservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/withdraw"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-withdraw/internal/redis/stats"
	repository "github.com/MamangRust/monolith-payment-gateway-withdraw/internal/repository/stats"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type withdrawStatsAmountDeps struct {
	ErrorHandler errorhandler.WithdrawStatisticErrorHandler

	Cache mencache.WithdrawStatsAmountCache

	Repository repository.WithdrawStatsAmountRepository

	Logger logger.LoggerInterface

	Mapper responseservice.WithdrawStatsAmountResponseMapper
}

type withdrawStatsAmountService struct {
	errorhandler errorhandler.WithdrawStatisticErrorHandler

	cache mencache.WithdrawStatsAmountCache

	repository repository.WithdrawStatsAmountRepository

	logger logger.LoggerInterface

	mapper responseservice.WithdrawStatsAmountResponseMapper

	observability observability.TraceLoggerObservability
}

func NewWithdrawStatsAmountService(deps *withdrawStatsAmountDeps) WithdrawStatsAmountService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "withdraw_stats_amount_service_request_total",
			Help: "Total number of requests to the WithdrawStatsAmountService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "withdraw_stats_amount_service_request_duration_seconds",
			Help:    "Histogram of request durations for the WithdrawStatsAmountService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(
		otel.Tracer("withdraw-stats-amount-service"), deps.Logger, requestCounter, requestDuration)

	return &withdrawStatsAmountService{
		errorhandler:  deps.ErrorHandler,
		cache:         deps.Cache,
		repository:    deps.Repository,
		logger:        deps.Logger,
		mapper:        deps.Mapper,
		observability: observability,
	}
}

// FindMonthlyWithdraws retrieves total amount statistics of monthly withdraws.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year to filter the monthly data.
//
// Returns:
//   - []*response.WithdrawMonthlyAmountResponse: List of total monthly withdraw amounts.
//   - *response.ErrorResponse: Error information if any occurred.
func (s *withdrawStatsAmountService) FindMonthlyWithdraws(ctx context.Context, year int) ([]*response.WithdrawMonthlyAmountResponse, *response.ErrorResponse) {
	const method = "FindMonthlyWithdraws"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMonthlyWithdraws(ctx, year); found {
		logSuccess("Successfully fetched monthly withdraws from cache", zap.Int("year", year))
		return data, nil
	}

	withdraws, err := s.repository.GetMonthlyWithdraws(ctx, year)

	if err != nil {
		return s.errorhandler.HandleMonthlyWithdrawAmountsError(err, method, "FAILED_GET_MONTHLY_WITHDRAW", span, &status, zap.Error(err))
	}

	responseWithdraws := s.mapper.ToWithdrawsAmountMonthlyResponses(withdraws)

	s.cache.SetCachedMonthlyWithdraws(ctx, year, responseWithdraws)

	logSuccess("Successfully fetched monthly withdraws", zap.Int("year", year))

	return responseWithdraws, nil
}

// FindYearlyWithdraws retrieves total amount statistics of yearly withdraws.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*response.WithdrawYearlyAmountResponse: List of total yearly withdraw amounts.
//   - *response.ErrorResponse: Error information if any occurred.
func (s *withdrawStatsAmountService) FindYearlyWithdraws(ctx context.Context, year int) ([]*response.WithdrawYearlyAmountResponse, *response.ErrorResponse) {
	const method = "FindYearlyWithdraws"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedYearlyWithdraws(ctx, year); found {
		logSuccess("Successfully fetched yearly withdraws from cache", zap.Int("year", year))
		return data, nil
	}

	withdraws, err := s.repository.GetYearlyWithdraws(ctx, year)
	if err != nil {
		return s.errorhandler.HandleYearlyWithdrawAmountsError(err, method, "FAILED_GET_YEARLY_WITHDRAW", span, &status, zap.Error(err))
	}

	responseWithdraws := s.mapper.ToWithdrawsAmountYearlyResponses(withdraws)

	s.cache.SetCachedYearlyWithdraws(ctx, year, responseWithdraws)

	logSuccess("Successfully fetched yearly withdraws", zap.Int("year", year))

	return responseWithdraws, nil
}
