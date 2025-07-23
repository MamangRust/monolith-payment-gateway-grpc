package topupstatsservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/topup"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/errorhandler"
	cache "github.com/MamangRust/monolith-payment-gateway-topup/internal/redis/stats"
	repository "github.com/MamangRust/monolith-payment-gateway-topup/internal/repository/stats"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type topupStatsStatusDeps struct {
	Cache cache.TopupStatsStatusCache

	ErrorHandler errorhandler.TopupStatisticErrorHandler

	Repository repository.TopupStatsStatusRepository

	Logger logger.LoggerInterface

	Mapper responseservice.TopupStatsStatusResponseMapper
}

type topupStatsStatusService struct {
	cache cache.TopupStatsStatusCache

	errorHandler errorhandler.TopupStatisticErrorHandler

	repository repository.TopupStatsStatusRepository

	logger logger.LoggerInterface

	mapper responseservice.TopupStatsStatusResponseMapper

	observability observability.TraceLoggerObservability
}

func NewTopupStatsStatusService(params *topupStatsStatusDeps) TopupStatsStatusService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "topup_stats_status_service_request_total",
			Help: "Total number of requests to the TopupStatsStatusService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "topup_stats_status_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TopupStatsStatusService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(
		otel.Tracer("topup-stats-status-service"), params.Logger, requestCounter, requestDuration)

	return &topupStatsStatusService{
		cache:         params.Cache,
		errorHandler:  params.ErrorHandler,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		observability: observability,
	}
}

// FindMonthTopupStatusSuccess retrieves monthly statistics of successful topups.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing month and year filters.
//
// Returns:
//   - []*response.TopupResponseMonthStatusSuccess: List of monthly succe
func (s *topupStatsStatusService) FindMonthTopupStatusSuccess(ctx context.Context, req *requests.MonthTopupStatus) ([]*response.TopupResponseMonthStatusSuccess, *response.ErrorResponse) {
	year := req.Year
	month := req.Month

	const method = "FindMonthTopupStatusSuccess"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("month", month))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthTopupStatusSuccessCache(ctx, req); found {
		logSuccess("Successfully fetched monthly topup status success from cache", zap.Int("year", year), zap.Int("month", month))
		return data, nil
	}

	records, err := s.repository.GetMonthTopupStatusSuccess(ctx, req)

	if err != nil {
		return s.errorHandler.HandleMonthTopupStatusSuccess(err, method, "FAILED_MONTHLY_TOPUP_STATUS_SUCCESS", span, &status, zap.Int("year", year), zap.Int("month", month))
	}
	so := s.mapper.ToTopupResponsesMonthStatusSuccess(records)

	s.cache.SetMonthTopupStatusSuccessCache(ctx, req, so)

	logSuccess("Successfully fetched monthly topup status success", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

// FindYearlyTopupStatusSuccess retrieves yearly statistics of successful topups.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which the statistics are requested.
//
// Returns:
//   - []*response.TopupResponseYearStatusSuccess: List of yearly success statistics.
//   - *response.ErrorResponse: Error details if retrieval fails.
func (s *topupStatsStatusService) FindYearlyTopupStatusSuccess(ctx context.Context, year int) ([]*response.TopupResponseYearStatusSuccess, *response.ErrorResponse) {
	const method = "FindYearlyTopupStatusSuccess"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTopupStatusSuccessCache(ctx, year); found {
		logSuccess("Successfully fetched yearly topup status success from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.repository.GetYearlyTopupStatusSuccess(ctx, year)
	if err != nil {
		return s.errorHandler.HandleYearlyTopupStatusSuccess(err, method, "FAILED_FIND_YEARLY_TOPUP_STATUS_SUCCESS", span, &status, zap.Error(err))
	}
	so := s.mapper.ToTopupResponsesYearStatusSuccess(records)

	s.cache.SetYearlyTopupStatusSuccessCache(ctx, year, so)

	logSuccess("Successfully fetched yearly topup status success", zap.Int("year", year))

	return so, nil
}

// FindMonthTopupStatusFailed retrieves monthly statistics of failed topups.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing month and year filters.
//
// Returns:
//   - []*response.TopupResponseMonthStatusFailed: List of monthly failed statistics.
//   - *response.ErrorResponse: Error details if retrieval fails.
func (s *topupStatsStatusService) FindMonthTopupStatusFailed(ctx context.Context, req *requests.MonthTopupStatus) ([]*response.TopupResponseMonthStatusFailed, *response.ErrorResponse) {
	year := req.Year
	month := req.Month

	const method = "FindMonthTopupStatusFailed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("month", month))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthTopupStatusFailedCache(ctx, req); found {
		logSuccess("Successfully fetched monthly topup status Failed from cache", zap.Int("year", year), zap.Int("month", month))
		return data, nil
	}

	records, err := s.repository.GetMonthTopupStatusFailed(ctx, req)
	if err != nil {
		return s.errorHandler.HandleMonthTopupStatusFailed(err, method, "FAILED_MONTHLY_TOPUP_STATUS_FAILED", span, &status, zap.Error(err))
	}
	so := s.mapper.ToTopupResponsesMonthStatusFailed(records)

	s.cache.SetMonthTopupStatusFailedCache(ctx, req, so)

	logSuccess("Successfully fetched monthly topup status Failed", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

// FindYearlyTopupStatusFailed retrieves yearly statistics of failed topups.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which the statistics are requested.
//
// Returns:
//   - []*response.TopupResponseYearStatusFailed: List of yearly failed statistics.
//   - *response.ErrorResponse: Error details if retrieval fails.
func (s *topupStatsStatusService) FindYearlyTopupStatusFailed(ctx context.Context, year int) ([]*response.TopupResponseYearStatusFailed, *response.ErrorResponse) {
	const method = "FindYearlyTopupStatusFailed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTopupStatusFailedCache(ctx, year); found {
		logSuccess("Successfully fetched yearly topup status Failed from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.repository.GetYearlyTopupStatusFailed(ctx, year)
	if err != nil {
		return s.errorHandler.HandleYearlyTopupStatusFailed(err, method, "FAILED_FIND_YEARLY_TOPUP_STATUS_FAILED", span, &status, zap.Int("year", year))
	}
	so := s.mapper.ToTopupResponsesYearStatusFailed(records)

	s.cache.SetYearlyTopupStatusFailedCache(ctx, year, so)

	logSuccess("Successfully fetched yearly topup status Failed", zap.Int("year", year))

	return so, nil
}
