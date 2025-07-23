package transferstatsservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/transfer"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-transfer/internal/redis/stats"
	repository "github.com/MamangRust/monolith-payment-gateway-transfer/internal/repository/stats"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type transferStatsStatusDeps struct {
	ErrorHandler errorhandler.TransferStatisticErrorHandler

	Cache mencache.TransferStatsStatusCache

	Repository repository.TransferStatsStatusRepository

	Logger logger.LoggerInterface

	Mapper responseservice.TransferStatsStatusResponseMapper
}

type transferStatsStatusService struct {
	errorHandler errorhandler.TransferStatisticErrorHandler

	cache mencache.TransferStatsStatusCache

	repository repository.TransferStatsStatusRepository

	logger logger.LoggerInterface

	mapper responseservice.TransferStatsStatusResponseMapper

	observability observability.TraceLoggerObservability
}

func NewTransferStatsStatusService(params *transferStatsStatusDeps) TransferStatsStatusService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transfer_stats_status_service_request_total",
			Help: "Total number of requests to the TransferStatsStatusService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transfer_stats_status_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TransferStatsStatusService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(
		otel.Tracer("transfer-stats-status-service"), params.Logger, requestCounter, requestDuration)

	return &transferStatsStatusService{
		errorHandler:  params.ErrorHandler,
		cache:         params.Cache,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		observability: observability,
	}
}

// FindMonthTransferStatusSuccess retrieves monthly successful transfer statistics.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing month and year filters.
//
// Returns:
//   - []*response.TransferResponseMonthStatusSuccess: List of monthly success transfer statistics.
//   - *response.ErrorResponse: Error response if an error occurs.
func (s *transferStatsStatusService) FindMonthTransferStatusSuccess(ctx context.Context, req *requests.MonthStatusTransfer) ([]*response.TransferResponseMonthStatusSuccess, *response.ErrorResponse) {

	year := req.Year
	month := req.Month

	const method = "FindMonthTransferStatusSuccess"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("month", month))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMonthTransferStatusSuccess(ctx, req); found {
		logSuccess("Successfully fetched monthly Transfer status success from cache", zap.Int("year", year), zap.Int("month", month))
		return data, nil
	}

	records, err := s.repository.GetMonthTransferStatusSuccess(ctx, req)

	if err != nil {
		return s.errorHandler.HandleMonthTransferStatusSuccessError(err, method, "FAILED_FIND_MONTH_TRANSFER_STATUS_SUCCESS", span, &status, zap.Error(err))
	}

	so := s.mapper.ToTransferResponsesMonthStatusSuccess(records)

	s.cache.SetCachedMonthTransferStatusSuccess(ctx, req, so)

	logSuccess("Successfully fetched monthly Transfer status success", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

// FindYearlyTransferStatusSuccess retrieves yearly successful transfer statistics.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*response.TransferResponseYearStatusSuccess: List of yearly success transfer statistics.
//   - *response.ErrorResponse: Error response if an error occurs.
func (s *transferStatsStatusService) FindYearlyTransferStatusSuccess(ctx context.Context, year int) ([]*response.TransferResponseYearStatusSuccess, *response.ErrorResponse) {
	const method = "FindYearlyTransferStatusSuccess"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedYearlyTransferStatusSuccess(ctx, year); found {
		logSuccess("Successfully fetched yearly Transfer status success from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.repository.GetYearlyTransferStatusSuccess(ctx, year)

	if err != nil {
		return s.errorHandler.HandleYearTransferStatusSuccessError(err, method, "FAILED_YEARLY_TRANSFER_STATUS_SUCCESS", span, &status, zap.Error(err))
	}

	so := s.mapper.ToTransferResponsesYearStatusSuccess(records)

	s.cache.SetCachedYearlyTransferStatusSuccess(ctx, year, so)

	logSuccess("Successfully fetched yearly Transfer status success", zap.Int("year", year))

	return so, nil

}

// FindMonthTransferStatusFailed retrieves monthly failed transfer statistics.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing month and year filters.
//
// Returns:
//   - []*response.TransferResponseMonthStatusFailed: List of monthly failed transfer statistics.
//   - *response.ErrorResponse: Error response if an error occurs.
func (s *transferStatsStatusService) FindMonthTransferStatusFailed(ctx context.Context, req *requests.MonthStatusTransfer) ([]*response.TransferResponseMonthStatusFailed, *response.ErrorResponse) {
	year := req.Year
	month := req.Month

	const method = "FindMonthTransferStatusFailed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("month", month))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMonthTransferStatusFailed(ctx, req); found {
		logSuccess("Successfully fetched monthly Transfer status Failed from cache", zap.Int("year", year), zap.Int("month", month))
		return data, nil
	}

	records, err := s.repository.GetMonthTransferStatusFailed(ctx, req)

	if err != nil {
		return s.errorHandler.HandleMonthTransferStatusFailedError(err, method, "FAILED_MONTHLY_TRANSFER_STATUS_FAILED", span, &status, zap.Error(err))
	}
	so := s.mapper.ToTransferResponsesMonthStatusFailed(records)

	s.cache.SetCachedMonthTransferStatusFailed(ctx, req, so)

	logSuccess("Successfully fetched monthly Transfer status Failed", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

// FindYearlyTransferStatusFailed retrieves yearly failed transfer statistics.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*response.TransferResponseYearStatusFailed: List of yearly failed transfer statistics.
//   - *response.ErrorResponse: Error response if an error occurs.
func (s *transferStatsStatusService) FindYearlyTransferStatusFailed(ctx context.Context, year int) ([]*response.TransferResponseYearStatusFailed, *response.ErrorResponse) {
	const method = "FindYearlyTransferStatusFailed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedYearlyTransferStatusFailed(ctx, year); found {
		logSuccess("Successfully fetched yearly Transfer status Failed from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.repository.GetYearlyTransferStatusFailed(ctx, year)
	if err != nil {
		return s.errorHandler.HandleYearTransferStatusFailedError(err, method, "FAILED_FIND_YEARLY_TRANSFER_STATUS_FAILED", span, &status, zap.Error(err))
	}
	so := s.mapper.ToTransferResponsesYearStatusFailed(records)

	s.cache.SetCachedYearlyTransferStatusFailed(ctx, year, so)

	logSuccess("Successfully fetched yearly Transfer status Failed", zap.Int("year", year))

	return so, nil
}
