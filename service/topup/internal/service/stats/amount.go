package topupstatsservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/topup"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-topup/internal/redis/stats"
	repository "github.com/MamangRust/monolith-payment-gateway-topup/internal/repository/stats"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type topupStatsAmountDeps struct {
	Cache mencache.TopupStatsAmountCache

	ErrorHandler errorhandler.TopupStatisticErrorHandler

	Repository repository.TopupStatsAmountRepository

	Logger logger.LoggerInterface

	Mapper responseservice.TopupStatsAmountResponseMapper
}

type topupStatsAmountService struct {
	cache mencache.TopupStatsAmountCache

	errorHandler errorhandler.TopupStatisticErrorHandler

	repository repository.TopupStatsAmountRepository

	logger logger.LoggerInterface

	mapper responseservice.TopupStatsAmountResponseMapper

	observability observability.TraceLoggerObservability
}

func NewTopupStatsAmountService(params *topupStatsAmountDeps) TopupStatsAmountService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "topup_stats_amount_service_request_total",
			Help: "Total number of requests to the TopupStatsAmountService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "topup_stats_amount_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TopupStatsAmountService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(
		otel.Tracer("topup-stats-amount-service"), params.Logger, requestCounter, requestDuration)

	return &topupStatsAmountService{
		cache:         params.Cache,
		errorHandler:  params.ErrorHandler,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		observability: observability,
	}
}

// FindMonthlyTopupAmounts retrieves monthly statistics of topup amounts.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which the statistics are requested.
//
// Returns:
//   - []*response.TopupMonthAmountResponse: List of monthly topup amounts.
//   - *response.ErrorResponse: Error details if retrieval fails.
func (s *topupStatsAmountService) FindMonthlyTopupAmounts(ctx context.Context, year int) ([]*response.TopupMonthAmountResponse, *response.ErrorResponse) {
	const method = "FindMonthlyTopupAmounts"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyTopupAmountsCache(ctx, year); found {
		logSuccess("Successfully fetched monthly topup amounts from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.repository.GetMonthlyTopupAmounts(ctx, year)
	if err != nil {
		return s.errorHandler.HandleMonthlyTopupAmounts(err, method, "FAILED_FIND_MONTHLY_TOPUP_AMOUNT", span, &status, zap.Error(err))
	}

	responses := s.mapper.ToTopupMonthlyAmountResponses(records)

	s.cache.SetMonthlyTopupAmountsCache(ctx, year, responses)

	logSuccess("Successfully fetched monthly topup amounts", zap.Int("year", year))

	return responses, nil
}

// FindYearlyTopupAmounts retrieves yearly statistics of topup amounts.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which the statistics are requested.
//
// Returns:
//   - []*response.TopupYearlyAmountResponse: List of yearly topup amounts.
//   - *response.ErrorResponse: Error details if retrieval fails.
func (s *topupStatsAmountService) FindYearlyTopupAmounts(ctx context.Context, year int) ([]*response.TopupYearlyAmountResponse, *response.ErrorResponse) {
	const method = "FindYearlyTopupAmounts"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTopupAmountsCache(ctx, year); found {
		logSuccess("Successfully fetched yearly topup amounts from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.repository.GetYearlyTopupAmounts(ctx, year)
	if err != nil {
		return s.errorHandler.HandleYearlyTopupAmounts(err, method, "FAILED_FIND_YEARLY_TOPUP_AMOUNTS", span, &status, zap.Error(err))
	}

	responses := s.mapper.ToTopupYearlyAmountResponses(records)

	s.cache.SetYearlyTopupAmountsCache(ctx, year, responses)

	logSuccess("Successfully fetched yearly topup amounts", zap.Int("year", year))

	return responses, nil
}
