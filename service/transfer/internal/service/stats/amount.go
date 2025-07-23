package transferstatsservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
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

type transferStatsAmountDeps struct {
	ErrorHandler errorhandler.TransferStatisticErrorHandler

	Cache mencache.TransferStatsAmountCache

	Repository repository.TransferStatsAmountRepository

	Logger logger.LoggerInterface

	Mapper responseservice.TransferAmountResponseMapper
}

type transferStatsAmountService struct {
	errorHandler errorhandler.TransferStatisticErrorHandler

	cache mencache.TransferStatsAmountCache

	repository repository.TransferStatsAmountRepository

	logger logger.LoggerInterface

	mapper responseservice.TransferAmountResponseMapper

	observability observability.TraceLoggerObservability
}

func NewTransferStatsAmountService(params *transferStatsAmountDeps) TransferStatsAmountService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transfer_stats_amount_service_request_total",
			Help: "Total number of requests to the TransferStatsAmountService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transfer_stats_amount_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TransferStatsAmountService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(
		otel.Tracer("transfer-stats-amount-service"), params.Logger, requestCounter, requestDuration)

	return &transferStatsAmountService{
		errorHandler:  params.ErrorHandler,
		cache:         params.Cache,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		observability: observability,
	}
}

// FindMonthlyTransferAmounts retrieves monthly total transfer amounts.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which monthly data is requested.
//
// Returns:
//   - []*response.TransferMonthAmountResponse: List of monthly transfer amount statistics.
//   - *response.ErrorResponse: Error response if an error occurs.
func (s *transferStatsAmountService) FindMonthlyTransferAmounts(ctx context.Context, year int) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse) {
	const method = "FindMonthlyTransferAmounts"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMonthTransferAmounts(ctx, year); found {
		logSuccess("Successfully fetched monthly transfer amounts from cache", zap.Int("year", year))
		return data, nil
	}

	amounts, err := s.repository.GetMonthlyTransferAmounts(ctx, year)
	if err != nil {
		return s.errorHandler.HandleMonthlyTransferAmountsError(err, method, "FAILED_FIND_MONTHLY_TRANSFER_AMOUNTS", span, &status, zap.Error(err))
	}

	responseAmounts := s.mapper.ToTransferResponsesMonthAmount(amounts)

	s.cache.SetCachedMonthTransferAmounts(ctx, year, responseAmounts)

	logSuccess("Successfully fetched monthly transfer amounts", zap.Int("year", year))

	return responseAmounts, nil
}

// FindYearlyTransferAmounts retrieves yearly total transfer amounts.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The year for which data is requested.
//
// Returns:
//   - []*response.TransferYearAmountResponse: List of yearly transfer amount statistics.
//   - *response.ErrorResponse: Error response if an error occurs.
func (s *transferStatsAmountService) FindYearlyTransferAmounts(ctx context.Context, year int) ([]*response.TransferYearAmountResponse, *response.ErrorResponse) {
	const method = "FindYearlyTransferAmounts"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedYearlyTransferAmounts(ctx, year); found {
		logSuccess("Successfully fetched yearly transfer amounts from cache", zap.Int("year", year))
		return data, nil
	}

	amounts, err := s.repository.GetYearlyTransferAmounts(ctx, year)
	if err != nil {
		return s.errorHandler.HandleYearlyTransferAmountsError(err, method, "FAILED_FIND_YEARLY_TRANSFER_AMOUNTS", span, &status, zap.Error(err))
	}

	responseAmounts := s.mapper.ToTransferResponsesYearAmount(amounts)

	s.cache.SetCachedYearlyTransferAmounts(ctx, year, responseAmounts)

	logSuccess("Successfully fetched yearly transfer amounts", zap.Int("year", year))

	return responseAmounts, nil
}
