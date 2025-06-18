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
	"go.opentelemetry.io/otel/codes"
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
	year := req.Year
	month := req.Month

	const method = "FindMonthWithdrawStatusSuccess"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("month", month))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthWithdrawStatusSuccessCache(req); found {
		logSuccess("Successfully fetched monthly withdraw status success from cache", zap.Int("year", year), zap.Int("month", month))
		return data, nil
	}

	records, err := s.withdrawStatisticRepository.GetMonthWithdrawStatusSuccess(req)

	if err != nil {
		return s.errorhandler.HandleMonthWithdrawStatusSuccessError(err, method, "FAILED_GET_MONTH_WITHDRAW_STATUS_SUCCESS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToWithdrawResponsesMonthStatusSuccess(records)

	s.mencache.SetCachedMonthWithdrawStatusSuccessCache(req, so)

	logSuccess("Successfully fetched monthly withdraw status success", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

func (s *withdrawStatisticService) FindYearlyWithdrawStatusSuccess(year int) ([]*response.WithdrawResponseYearStatusSuccess, *response.ErrorResponse) {
	const method = "FindYearlyWithdrawStatusSuccess"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearlyWithdrawStatusSuccessCache(year); found {
		s.logger.Debug("Successfully fetched yearly withdraw status success from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.withdrawStatisticRepository.GetYearlyWithdrawStatusSuccess(year)

	if err != nil {
		return s.errorhandler.HandleYearWithdrawStatusSuccessError(err, method, "FAILED_GET_YEARLY_WITHDRAW_STATUS_SUCCESS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToWithdrawResponsesYearStatusSuccess(records)

	s.mencache.SetCachedYearlyWithdrawStatusSuccessCache(year, so)

	logSuccess("Successfully fetched yearly withdraw status success", zap.Int("year", year))

	return so, nil
}

func (s *withdrawStatisticService) FindMonthWithdrawStatusFailed(req *requests.MonthStatusWithdraw) ([]*response.WithdrawResponseMonthStatusFailed, *response.ErrorResponse) {
	year := req.Year
	month := req.Month

	const method = "FindMonthWithdrawStatusFailed"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("month", month))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthWithdrawStatusFailedCache(req); found {
		logSuccess("Successfully fetched monthly Withdraw status Failed from cache", zap.Int("year", year), zap.Int("month", month))
		return data, nil
	}

	records, err := s.withdrawStatisticRepository.GetMonthWithdrawStatusFailed(req)

	if err != nil {
		return s.errorhandler.HandleMonthWithdrawStatusFailedError(err, method, "FAILED_GET_MONTH_WITHDRAW_STATUS_FAILED", span, &status, zap.Error(err))
	}

	so := s.mapping.ToWithdrawResponsesMonthStatusFailed(records)

	s.mencache.SetCachedMonthWithdrawStatusFailedCache(req, so)

	return so, nil
}

func (s *withdrawStatisticService) FindYearlyWithdrawStatusFailed(year int) ([]*response.WithdrawResponseYearStatusFailed, *response.ErrorResponse) {
	const method = "FindYearlyWithdrawStatusFailed"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearlyWithdrawStatusFailedCache(year); found {
		logSuccess("Successfully fetched yearly Withdraw status Failed from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.withdrawStatisticRepository.GetYearlyWithdrawStatusFailed(year)

	if err != nil {
		return s.errorhandler.HandleYearWithdrawStatusFailedError(err, method, "FAILED_GET_YEARLY_WITHDRAW_STATUS_FAILED", span, &status, zap.Error(err))
	}
	so := s.mapping.ToWithdrawResponsesYearStatusFailed(records)

	s.mencache.SetCachedYearlyWithdrawStatusFailedCache(year, so)

	logSuccess("Successfully fetched yearly Withdraw status Failed", zap.Int("year", year))

	return so, nil
}

func (s *withdrawStatisticService) FindMonthlyWithdraws(year int) ([]*response.WithdrawMonthlyAmountResponse, *response.ErrorResponse) {
	const method = "FindMonthlyWithdraws"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthlyWithdraws(year); found {
		logSuccess("Successfully fetched monthly withdraws from cache", zap.Int("year", year))
		return data, nil
	}

	withdraws, err := s.withdrawStatisticRepository.GetMonthlyWithdraws(year)

	if err != nil {
		return s.errorhandler.HandleMonthlyWithdrawAmountsError(err, method, "FAILED_GET_MONTHLY_WITHDRAW", span, &status, zap.Error(err))
	}

	responseWithdraws := s.mapping.ToWithdrawsAmountMonthlyResponses(withdraws)

	s.mencache.SetCachedMonthlyWithdraws(year, responseWithdraws)

	logSuccess("Successfully fetched monthly withdraws", zap.Int("year", year))

	return responseWithdraws, nil
}

func (s *withdrawStatisticService) FindYearlyWithdraws(year int) ([]*response.WithdrawYearlyAmountResponse, *response.ErrorResponse) {
	const method = "FindYearlyWithdraws"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearlyWithdraws(year); found {
		logSuccess("Successfully fetched yearly withdraws from cache", zap.Int("year", year))
		return data, nil
	}

	withdraws, err := s.withdrawStatisticRepository.GetYearlyWithdraws(year)
	if err != nil {
		return s.errorhandler.HandleYearlyWithdrawAmountsError(err, method, "FAILED_GET_YEARLY_WITHDRAW", span, &status, zap.Error(err))
	}

	responseWithdraws := s.mapping.ToWithdrawsAmountYearlyResponses(withdraws)

	s.mencache.SetCachedYearlyWithdraws(year, responseWithdraws)

	logSuccess("Successfully fetched yearly withdraws", zap.Int("year", year))

	return responseWithdraws, nil
}

func (s *withdrawStatisticService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
	trace.Span,
	func(string),
	string,
	func(string, ...zap.Field),
) {
	start := time.Now()
	status := "success"

	_, span := s.trace.Start(s.ctx, method)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.AddEvent("Start: " + method)

	s.logger.Info("Start: " + method)

	end := func(status string) {
		s.recordMetrics(method, status, start)
		code := codes.Ok
		if status != "success" {
			code = codes.Error
		}
		span.SetStatus(code, status)
		span.End()
	}

	logSuccess := func(msg string, fields ...zap.Field) {
		span.AddEvent(msg)
		s.logger.Info(msg, fields...)
	}

	return span, end, status, logSuccess
}

func (s *withdrawStatisticService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
