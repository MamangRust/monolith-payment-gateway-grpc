package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	traceunic "github.com/MamangRust/monolith-payment-gateway-pkg/trace_unic"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type saldoStatisticsService struct {
	ctx                  context.Context
	trace                trace.Tracer
	logger               logger.LoggerInterface
	mapping              responseservice.SaldoResponseMapper
	saldoStatsRepository repository.SaldoStatisticsRepository
	requestCounter       *prometheus.CounterVec
	requestDuration      *prometheus.HistogramVec
}

func NewSaldoStatisticsService(ctx context.Context, saldoStatsRepository repository.SaldoStatisticsRepository, logger logger.LoggerInterface, mapping responseservice.SaldoResponseMapper) *saldoStatisticsService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "saldo_statistics_service_request_total",
			Help: "Total number of requests to the SaldoStatisticsService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "saldo_statistics_service_request_duration_seconds",
			Help:    "Histogram of request durations for the SaldoStatisticsService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &saldoStatisticsService{
		ctx:                  ctx,
		trace:                otel.Tracer("saldo-statistics-service"),
		saldoStatsRepository: saldoStatsRepository,
		logger:               logger,
		mapping:              mapping,
		requestCounter:       requestCounter,
		requestDuration:      requestDuration,
	}
}

func (s *saldoStatisticsService) FindMonthlyTotalSaldoBalance(req *requests.MonthTotalSaldoBalance) ([]*response.SaldoMonthTotalBalanceResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyTotalSaldoBalance", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyTotalSaldoBalance")
	defer span.End()

	year := req.Year
	month := req.Month

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
	)

	s.logger.Debug("Fetching monthly total saldo balance", zap.Int("year", year), zap.Int("month", month))

	res, err := s.saldoStatsRepository.GetMonthlyTotalSaldoBalance(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_TOTAL_SALDO_BALANCE")

		s.logger.Error("Failed to fetch monthly total saldo balance", zap.Error(err), zap.Int("year", year), zap.Int("month", month))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch monthly total saldo balance")
		status = "failed_to_fetch_monthly_total_saldo_balance"

		return nil, saldo_errors.ErrFailedFindMonthlyTotalSaldoBalance
	}

	responses := s.mapping.ToSaldoMonthTotalBalanceResponses(res)

	s.logger.Debug("Successfully fetched monthly total saldo balance", zap.Int("year", year), zap.Int("month", month))

	return responses, nil
}

func (s *saldoStatisticsService) FindYearTotalSaldoBalance(year int) ([]*response.SaldoYearTotalBalanceResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearTotalSaldoBalance", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearTotalSaldoBalance")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Fetching yearly total saldo balance", zap.Int("year", year))

	res, err := s.saldoStatsRepository.GetYearTotalSaldoBalance(year)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEAR_TOTAL_SALDO_BALANCE")

		s.logger.Error("Failed to fetch yearly total saldo balance", zap.Error(err), zap.Int("year", year))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch yearly total saldo balance")
		status = "failed_to_fetch_yearly_total_saldo_balance"

		return nil, saldo_errors.ErrFailedFindYearTotalSaldoBalance
	}

	s.logger.Debug("Successfully fetched yearly total saldo balance", zap.Int("year", year))

	so := s.mapping.ToSaldoYearTotalBalanceResponses(res)

	return so, nil
}

func (s *saldoStatisticsService) FindMonthlySaldoBalances(year int) ([]*response.SaldoMonthBalanceResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlySaldoBalances", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlySaldoBalances")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Fetching monthly saldo balances", zap.Int("year", year))

	res, err := s.saldoStatsRepository.GetMonthlySaldoBalances(year)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_SALDO_BALANCES")

		s.logger.Error("Failed to fetch monthly saldo balances", zap.Error(err), zap.Int("year", year))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch monthly saldo balances")
		status = "failed_to_fetch_monthly_saldo_balances"

		return nil, saldo_errors.ErrFailedFindMonthlySaldoBalances
	}

	responses := s.mapping.ToSaldoMonthBalanceResponses(res)

	s.logger.Debug("Successfully fetched monthly saldo balances", zap.Int("year", year))

	return responses, nil
}

func (s *saldoStatisticsService) FindYearlySaldoBalances(year int) ([]*response.SaldoYearBalanceResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlySaldoBalances", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlySaldoBalances")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Fetching yearly saldo balances", zap.Int("year", year))

	res, err := s.saldoStatsRepository.GetYearlySaldoBalances(year)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_SALDO_BALANCES")

		s.logger.Error("Failed to fetch yearly saldo balances", zap.Error(err), zap.Int("year", year))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch yearly saldo balances")
		status = "failed_to_fetch_yearly_saldo_balances"

		return nil, saldo_errors.ErrFailedFindYearlySaldoBalances
	}

	responses := s.mapping.ToSaldoYearBalanceResponses(res)

	s.logger.Debug("Successfully fetched yearly saldo balances", zap.Int("year", year))

	return responses, nil
}

func (s *saldoStatisticsService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
