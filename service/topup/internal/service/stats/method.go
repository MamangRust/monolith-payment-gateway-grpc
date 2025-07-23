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

type topupStatsMethodDeps struct {
	Cache mencache.TopupStatsMethodCache

	ErrorHandler errorhandler.TopupStatisticErrorHandler

	Repository repository.TOpupStatsMethodRepository

	Logger logger.LoggerInterface

	Mapper responseservice.TopupStatsMethodResponseMapper
}

type topupStatsMethodService struct {
	cache mencache.TopupStatsMethodCache

	errorHandler errorhandler.TopupStatisticErrorHandler

	repository repository.TOpupStatsMethodRepository

	logger logger.LoggerInterface

	mapper responseservice.TopupStatsMethodResponseMapper

	observability observability.TraceLoggerObservability
}

func NewTopupStatsMethodService(params *topupStatsMethodDeps) TopupStatsMethodService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "topup_stats_method_service_request_total",
			Help: "Total number of requests to the TopupStatsMethodService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "topup_stats_method_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TopupStatsMethodService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(
		otel.Tracer("topup-stats-method-service"), params.Logger, requestCounter, requestDuration)

	return &topupStatsMethodService{
		cache:         params.Cache,
		errorHandler:  params.ErrorHandler,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		observability: observability,
	}
}

// FindMonthlyTopupMethods retrieves monthly statistics grouped by topup methods.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which the statistics are requested.
//
// Returns:
//   - []*response.TopupMonthMethodResponse: List of monthly method usage.
//   - *response.ErrorResponse: Error details if retrieval fails.
func (s *topupStatsMethodService) FindMonthlyTopupMethods(ctx context.Context, year int) ([]*response.TopupMonthMethodResponse, *response.ErrorResponse) {
	const method = "FindMonthlyTopupMethods"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyTopupMethodsCache(ctx, year); found {
		logSuccess("Successfully fetched monthly topup methods from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.repository.GetMonthlyTopupMethods(ctx, year)
	if err != nil {
		return s.errorHandler.HandleMonthlyTopupMethods(err, method, "FAILED_FIND_MONTHLY_TOPUP_METHODS", span, &status, zap.Error(err))
	}

	responses := s.mapper.ToTopupMonthlyMethodResponses(records)

	s.cache.SetMonthlyTopupMethodsCache(ctx, year, responses)

	logSuccess("Successfully fetched monthly topup methods", zap.Int("year", year))

	return responses, nil
}

// FindYearlyTopupMethods retrieves yearly statistics grouped by topup methods.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which the statistics are requested.
//
// Returns:
//   - []*response.TopupYearlyMethodResponse: List of yearly method usage.
//   - *response.ErrorResponse: Error details if retrieval fails.
func (s *topupStatsMethodService) FindYearlyTopupMethods(ctx context.Context, year int) ([]*response.TopupYearlyMethodResponse, *response.ErrorResponse) {
	const method = "FindYearlyTopupMethods"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTopupMethodsCache(ctx, year); found {
		logSuccess("Successfully fetched yearly topup methods from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.repository.GetYearlyTopupMethods(ctx, year)
	if err != nil {
		return s.errorHandler.HandleYearlyTopupMethods(err, method, "FAILED_FIND_YEARLY_TOPUP_METHODS", span, &status, zap.Error(err))
	}

	responses := s.mapper.ToTopupYearlyMethodResponses(records)

	s.cache.SetYearlyTopupMethodsCache(ctx, year, responses)

	logSuccess("Successfully fetched yearly topup methods", zap.Int("year", year))

	return responses, nil
}
