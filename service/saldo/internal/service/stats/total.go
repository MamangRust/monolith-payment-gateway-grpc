package saldostatsservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-saldo/internal/redis/stats"
	repository "github.com/MamangRust/monolith-payment-gateway-saldo/internal/repository/stats"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/saldo"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type saldoStatsTotalBalanceDeps struct {
	Cache mencache.SaldoStatsTotalCache

	ErrorHandler errorhandler.SaldoStatisticErrorHandler

	Repository repository.SaldoStatsTotalSaldoRepository

	Logger logger.LoggerInterface

	Mapper responseservice.SaldoStatisticTotalBalanceResponseMapper
}

type saldoStatsTotalBalanceService struct {
	mencache mencache.SaldoStatsTotalCache

	repository repository.SaldoStatsTotalSaldoRepository

	errorhandler errorhandler.SaldoStatisticErrorHandler

	logger logger.LoggerInterface

	mapper responseservice.SaldoStatisticTotalBalanceResponseMapper

	observability observability.TraceLoggerObservability
}

func NewSaldoStatsTotalBalanceService(params *saldoStatsTotalBalanceDeps) SaldoStatsTotalBalanceService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "saldo_stats_total_balance_service_request_total",
		Help: "The total number of requests SaldoStatsTotalBalanceService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "saldo_stats_total_balance_service_request_duration_seconds",
		Help:    "The duration of requests SaldoStatsTotalBalanceService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("saldo-stats-balance-service"), params.Logger, requestCounter, requestDuration)

	return &saldoStatsTotalBalanceService{
		mencache:      params.Cache,
		repository:    params.Repository,
		errorhandler:  params.ErrorHandler,
		logger:        params.Logger,
		mapper:        params.Mapper,
		observability: observability,
	}
}

// FindMonthlyTotalSaldoBalance retrieves the total saldo balance grouped by month
// based on the provided request, which may include year and optional filters.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request object containing filter criteria (e.g., year, card number).
//
// Returns:
//   - []*response.SaldoMonthTotalBalanceResponse: List of monthly total saldo balances.
//   - *response.ErrorResponse: An error response if the operation fails.
func (s *saldoStatsTotalBalanceService) FindMonthlyTotalSaldoBalance(ctx context.Context, req *requests.MonthTotalSaldoBalance) ([]*response.SaldoMonthTotalBalanceResponse, *response.ErrorResponse) {

	year := req.Year
	month := req.Month

	const method = "FindMonthlyTotalSaldoBalance"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("month", month))

	defer func() {
		end(status)
	}()

	if cache, found := s.mencache.GetMonthlyTotalSaldoBalanceCache(ctx, req); found {
		logSuccess("Successfully fetched monthly total saldo balance from cache", zap.Int("year", year), zap.Int("month", month))
		return cache, nil
	}

	res, err := s.repository.GetMonthlyTotalSaldoBalance(ctx, req)
	if err != nil {
		return s.errorhandler.HandleMonthlyTotalSaldoBalanceError(err, method, "FAILED_FIND_MONTHLY_TOTAL_SALDO_BALANCE", span, &status, zap.Error(err))
	}

	responses := s.mapper.ToSaldoMonthTotalBalanceResponses(res)

	s.mencache.SetMonthlyTotalSaldoCache(ctx, req, responses)

	logSuccess("Successfully fetched monthly total saldo balance", zap.Int("year", year), zap.Int("month", month))

	return responses, nil
}

// FindYearTotalSaldoBalance retrieves the total saldo balance aggregated for a specific year.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - year: The target year to retrieve the total saldo balance.
//
// Returns:
//   - []*response.SaldoYearTotalBalanceResponse: List of yearly total saldo balances.
//   - *response.ErrorResponse: An error response if the operation fails.
func (s *saldoStatsTotalBalanceService) FindYearTotalSaldoBalance(ctx context.Context, year int) ([]*response.SaldoYearTotalBalanceResponse, *response.ErrorResponse) {
	const method = "FindYearTotalSaldoBalance"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cache, found := s.mencache.GetYearTotalSaldoBalanceCache(ctx, year); found {
		logSuccess("Successfully fetched yearly total saldo balance from cache", zap.Int("year", year))
		return cache, nil
	}

	res, err := s.repository.GetYearTotalSaldoBalance(ctx, year)

	if err != nil {
		return s.errorhandler.HandleYearlyTotalSaldoBalanceError(err, method, "FAILED_FIND_YEARLY_TOTAL_SALDO_BALANCE", span, &status, zap.Error(err))
	}
	so := s.mapper.ToSaldoYearTotalBalanceResponses(res)

	s.mencache.SetYearTotalSaldoBalanceCache(ctx, year, so)

	logSuccess("Successfully fetched yearly total saldo balance", zap.Int("year", year))

	return so, nil
}
