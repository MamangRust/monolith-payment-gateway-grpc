package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-withdraw/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type withdrawStatisticService struct {
	ctx                         context.Context
	errorhandler                errorhandler.WithdrawStatisticErrorHandler
	mencache                    mencache.WithdrawStatisticCache
	trace                       trace.Tracer
	withdrawStatisticRepository repository.WithdrawStatisticRepository
	logger                      logger.LoggerInterface
	mapping                     responseservice.WithdrawResponseMapper
	requestCounter              *prometheus.CounterVec
	requestDuration             *prometheus.HistogramVec
}

func NewWithdrawStatisticService(ctx context.Context, errorhandler errorhandler.WithdrawStatisticErrorHandler,
	mencache mencache.WithdrawStatisticCache, withdrawStatisticRepository repository.WithdrawStatisticRepository, logger logger.LoggerInterface, mapping responseservice.WithdrawResponseMapper) *withdrawStatisticService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "withdraw_statistic_service_request_total",
			Help: "Total number of requests to the WithdrawStatisticService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "withdraw_statistic_service_request_duration_seconds",
			Help:    "Histogram of request durations for the WithdrawStatisticService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &withdrawStatisticService{
		ctx:                         ctx,
		trace:                       otel.Tracer("withdraw-statistic-service"),
		withdrawStatisticRepository: withdrawStatisticRepository,
		logger:                      logger,
		mapping:                     mapping,
		requestCounter:              requestCounter,
		requestDuration:             requestDuration,
		mencache:                    mencache,
		errorhandler:                errorhandler,
	}
}

func (s *withdrawStatisticService) FindMonthWithdrawStatusSuccess(req *requests.MonthStatusWithdraw) ([]*response.WithdrawResponseMonthStatusSuccess, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("find_month_withdraw_status_success", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "find_month_withdraw_status_success")
	defer span.End()

	year := req.Year
	month := req.Month

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
	)

	s.logger.Debug("Fetching monthly Withdraw status success", zap.Int("year", year), zap.Int("month", month))

	if data := s.mencache.GetCachedMonthWithdrawStatusSuccessCache(req); data != nil {
		s.logger.Debug("Successfully fetched monthly withdraw status success from cache", zap.Int("year", year), zap.Int("month", month))
		return data, nil
	}

	records, err := s.withdrawStatisticRepository.GetMonthWithdrawStatusSuccess(req)

	if err != nil {
		return s.errorhandler.HandleMonthWithdrawStatusSuccessError(err, "find_month_withdraw_status_success", "FAILED_GET_MONTH_WITHDRAW_STATUS_SUCCESS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToWithdrawResponsesMonthStatusSuccess(records)

	s.mencache.SetCachedMonthWithdrawStatusSuccessCache(req, so)

	s.logger.Debug("Successfully fetched monthly Withdraw status success", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

func (s *withdrawStatisticService) FindYearlyWithdrawStatusSuccess(year int) ([]*response.WithdrawResponseYearStatusSuccess, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyWithdrawStatusSuccess", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyWithdrawStatusSuccess")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Fetching yearly Withdraw status success", zap.Int("year", year))

	if data := s.mencache.GetCachedYearlyWithdrawStatusSuccessCache(year); data != nil {
		s.logger.Debug("Successfully fetched yearly withdraw status success from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.withdrawStatisticRepository.GetYearlyWithdrawStatusSuccess(year)

	if err != nil {
		return s.errorhandler.HandleYearWithdrawStatusSuccessError(err, "FindYearlyWithdrawStatusSuccess", "FAILED_GET_YEARLY_WITHDRAW_STATUS_SUCCESS", span, &status, zap.Error(err))
	}
	so := s.mapping.ToWithdrawResponsesYearStatusSuccess(records)

	s.mencache.SetCachedYearlyWithdrawStatusSuccessCache(year, so)

	s.logger.Debug("Successfully fetched yearly Withdraw status success", zap.Int("year", year))

	return so, nil
}

func (s *withdrawStatisticService) FindMonthWithdrawStatusFailed(req *requests.MonthStatusWithdraw) ([]*response.WithdrawResponseMonthStatusFailed, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthWithdrawStatusFailed", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthWithdrawStatusFailed")
	defer span.End()

	year := req.Year
	month := req.Month

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
	)
	s.logger.Debug("Fetching monthly Withdraw status Failed", zap.Int("year", year), zap.Int("month", month))

	if data := s.mencache.GetCachedMonthWithdrawStatusFailedCache(req); data != nil {
		s.logger.Debug("Successfully fetched monthly withdraw status Failed from cache", zap.Int("year", year), zap.Int("month", month))
		return data, nil
	}

	records, err := s.withdrawStatisticRepository.GetMonthWithdrawStatusFailed(req)

	if err != nil {
		return s.errorhandler.HandleMonthWithdrawStatusFailedError(err, "FindMonthWithdrawStatusFailed", "FAILED_GET_MONTH_WITHDRAW_STATUS_FAILED", span, &status, zap.Error(err))
	}

	so := s.mapping.ToWithdrawResponsesMonthStatusFailed(records)

	s.mencache.SetCachedMonthWithdrawStatusFailedCache(req, so)

	s.logger.Debug("Failedfully fetched monthly Withdraw status Failed", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

func (s *withdrawStatisticService) FindYearlyWithdrawStatusFailed(year int) ([]*response.WithdrawResponseYearStatusFailed, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyWithdrawStatusFailed", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyWithdrawStatusFailed")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Fetching yearly Withdraw status Failed", zap.Int("year", year))

	if data := s.mencache.GetCachedYearlyWithdrawStatusFailedCache(year); data != nil {
		s.logger.Debug("Successfully fetched yearly withdraw status Failed from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.withdrawStatisticRepository.GetYearlyWithdrawStatusFailed(year)

	if err != nil {
		return s.errorhandler.HandleYearWithdrawStatusFailedError(err, "FindYearlyWithdrawStatusFailed", "FAILED_GET_YEARLY_WITHDRAW_STATUS_FAILED", span, &status, zap.Error(err))
	}
	so := s.mapping.ToWithdrawResponsesYearStatusFailed(records)

	s.mencache.SetCachedYearlyWithdrawStatusFailedCache(year, so)

	s.logger.Debug("Failedfully fetched yearly Withdraw status Failed", zap.Int("year", year))

	return so, nil
}

func (s *withdrawStatisticService) FindMonthlyWithdraws(year int) ([]*response.WithdrawMonthlyAmountResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyWithdraws", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyWithdraws")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Fetching monthly withdraws", zap.Int("year", year))

	if data := s.mencache.GetCachedMonthlyWithdraws(year); data != nil {
		s.logger.Debug("Successfully fetched monthly withdraws from cache", zap.Int("year", year))
		return data, nil
	}

	withdraws, err := s.withdrawStatisticRepository.GetMonthlyWithdraws(year)

	if err != nil {
		return s.errorhandler.HandleMonthlyWithdrawAmountsError(err, "FindMonthlyWithdraws", "FAILED_GET_MONTHLY_WITHDRAW", span, &status, zap.Error(err))
	}

	responseWithdraws := s.mapping.ToWithdrawsAmountMonthlyResponses(withdraws)

	s.mencache.SetCachedMonthlyWithdraws(year, responseWithdraws)

	s.logger.Debug("Successfully fetched monthly withdraws", zap.Int("year", year))

	return responseWithdraws, nil
}

func (s *withdrawStatisticService) FindYearlyWithdraws(year int) ([]*response.WithdrawYearlyAmountResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyWithdraws", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyWithdraws")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Fetching yearly withdraws", zap.Int("year", year))

	if data := s.mencache.GetCachedYearlyWithdraws(year); data != nil {
		s.logger.Debug("Successfully fetched yearly withdraws from cache", zap.Int("year", year))
		return data, nil
	}

	withdraws, err := s.withdrawStatisticRepository.GetYearlyWithdraws(year)
	if err != nil {
		return s.errorhandler.HandleYearlyWithdrawAmountsError(err, "FindYearlyWithdraws", "FAILED_GET_YEARLY_WITHDRAW", span, &status, zap.Error(err))
	}

	responseWithdraws := s.mapping.ToWithdrawsAmountYearlyResponses(withdraws)

	s.mencache.SetCachedYearlyWithdraws(year, responseWithdraws)

	s.logger.Debug("Successfully fetched yearly withdraws", zap.Int("year", year))

	return responseWithdraws, nil
}

func (s *withdrawStatisticService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
