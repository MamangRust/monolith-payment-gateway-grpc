package withdrawstatsservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
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

type withdrawStatsStatusDeps struct {
	ErrorHandler errorhandler.WithdrawStatisticErrorHandler

	Cache mencache.WithdrawStatsStatusCache

	Repository repository.WithdrawStatsStatusRepository

	Logger logger.LoggerInterface

	Mapper responseservice.WithdrawStatsStatusResponseMapper
}

type withdrawStatsStatusService struct {
	errorhandler errorhandler.WithdrawStatisticErrorHandler

	cache mencache.WithdrawStatsStatusCache

	repository repository.WithdrawStatsStatusRepository

	logger logger.LoggerInterface

	mapper responseservice.WithdrawStatsStatusResponseMapper

	observability observability.TraceLoggerObservability
}

func NewWithdrawStatsStatusService(deps *withdrawStatsStatusDeps) WithdrawStatsStatusService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "withdraw_stats_status_service_request_total",
			Help: "Total number of requests to the WithdrawStatsStatusService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "withdraw_stats_status_service_request_duration_seconds",
			Help:    "Histogram of request durations for the WithdrawStatsStatusService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(
		otel.Tracer("withdraw-stats-status-service"), deps.Logger, requestCounter, requestDuration)

	return &withdrawStatsStatusService{
		errorhandler:  deps.ErrorHandler,
		cache:         deps.Cache,
		repository:    deps.Repository,
		logger:        deps.Logger,
		mapper:        deps.Mapper,
		observability: observability,
	}
}

// FindMonthWithdrawStatusSuccess retrieves monthly successful withdraw statistics.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing the month and year for filtering.
//
// Returns:
//   - []*response.WithdrawResponseMonthStatusSuccess: List of successful monthly withdraw statistics.
//   - *response.ErrorResponse: Error information if any occurred.
func (s *withdrawStatsStatusService) FindMonthWithdrawStatusSuccess(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*response.WithdrawResponseMonthStatusSuccess, *response.ErrorResponse) {
	year := req.Year
	month := req.Month

	const method = "FindMonthWithdrawStatusSuccess"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("month", month))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMonthWithdrawStatusSuccessCache(ctx, req); found {
		logSuccess("Successfully fetched monthly withdraw status success from cache", zap.Int("year", year), zap.Int("month", month))
		return data, nil
	}

	records, err := s.repository.GetMonthWithdrawStatusSuccess(ctx, req)

	if err != nil {
		return s.errorhandler.HandleMonthWithdrawStatusSuccessError(err, method, "FAILED_GET_MONTH_WITHDRAW_STATUS_SUCCESS", span, &status, zap.Error(err))
	}

	so := s.mapper.ToWithdrawResponsesMonthStatusSuccess(records)

	s.cache.SetCachedMonthWithdrawStatusSuccessCache(ctx, req, so)

	logSuccess("Successfully fetched monthly withdraw status success", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

// FindYearlyWithdrawStatusSuccess retrieves yearly successful withdraw statistics.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year to filter the data.
//
// Returns:
//   - []*response.WithdrawResponseYearStatusSuccess: List of successful yearly withdraw statistics.
//   - *response.ErrorResponse: Error information if any occurred.
func (s *withdrawStatsStatusService) FindYearlyWithdrawStatusSuccess(ctx context.Context, year int) ([]*response.WithdrawResponseYearStatusSuccess, *response.ErrorResponse) {
	const method = "FindYearlyWithdrawStatusSuccess"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedYearlyWithdrawStatusSuccessCache(ctx, year); found {
		s.logger.Debug("Successfully fetched yearly withdraw status success from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.repository.GetYearlyWithdrawStatusSuccess(ctx, year)

	if err != nil {
		return s.errorhandler.HandleYearWithdrawStatusSuccessError(err, method, "FAILED_GET_YEARLY_WITHDRAW_STATUS_SUCCESS", span, &status, zap.Error(err))
	}

	so := s.mapper.ToWithdrawResponsesYearStatusSuccess(records)

	s.cache.SetCachedYearlyWithdrawStatusSuccessCache(ctx, year, so)

	logSuccess("Successfully fetched yearly withdraw status success", zap.Int("year", year))

	return so, nil
}

// FindMonthWithdrawStatusFailed retrieves monthly failed withdraw statistics.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing the month and year for filtering.
//
// Returns:
//   - []*response.WithdrawResponseMonthStatusFailed: List of failed monthly withdraw statistics.
//   - *response.ErrorResponse: Error information if any occurred.
func (s *withdrawStatsStatusService) FindMonthWithdrawStatusFailed(ctx context.Context, req *requests.MonthStatusWithdraw) ([]*response.WithdrawResponseMonthStatusFailed, *response.ErrorResponse) {
	year := req.Year
	month := req.Month

	const method = "FindMonthWithdrawStatusFailed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("month", month))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMonthWithdrawStatusFailedCache(ctx, req); found {
		logSuccess("Successfully fetched monthly Withdraw status Failed from cache", zap.Int("year", year), zap.Int("month", month))
		return data, nil
	}

	records, err := s.repository.GetMonthWithdrawStatusFailed(ctx, req)

	if err != nil {
		return s.errorhandler.HandleMonthWithdrawStatusFailedError(err, method, "FAILED_GET_MONTH_WITHDRAW_STATUS_FAILED", span, &status, zap.Error(err))
	}

	so := s.mapper.ToWithdrawResponsesMonthStatusFailed(records)

	s.cache.SetCachedMonthWithdrawStatusFailedCache(ctx, req, so)

	return so, nil
}

// FindYearlyWithdrawStatusFailed retrieves yearly failed withdraw statistics.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year to filter the data.
//
// Returns:
//   - []*response.WithdrawResponseYearStatusFailed: List of failed yearly withdraw statistics.
//   - *response.ErrorResponse: Error information if any occurred.
func (s *withdrawStatsStatusService) FindYearlyWithdrawStatusFailed(ctx context.Context, year int) ([]*response.WithdrawResponseYearStatusFailed, *response.ErrorResponse) {
	const method = "FindYearlyWithdrawStatusFailed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedYearlyWithdrawStatusFailedCache(ctx, year); found {
		logSuccess("Successfully fetched yearly Withdraw status Failed from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.repository.GetYearlyWithdrawStatusFailed(ctx, year)

	if err != nil {
		return s.errorhandler.HandleYearWithdrawStatusFailedError(err, method, "FAILED_GET_YEARLY_WITHDRAW_STATUS_FAILED", span, &status, zap.Error(err))
	}
	so := s.mapper.ToWithdrawResponsesYearStatusFailed(records)

	s.cache.SetCachedYearlyWithdrawStatusFailedCache(ctx, year, so)

	logSuccess("Successfully fetched yearly Withdraw status Failed", zap.Int("year", year))

	return so, nil
}
