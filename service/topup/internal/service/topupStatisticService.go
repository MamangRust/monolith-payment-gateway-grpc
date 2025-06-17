package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-topup/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type topupStasticService struct {
	ctx                      context.Context
	mencache                 mencache.TopupStatisticCache
	errorhandler             errorhandler.TopupStatisticErrorHandler
	trace                    trace.Tracer
	topupStatisticRepository repository.TopupStatisticRepository
	logger                   logger.LoggerInterface
	mapping                  responseservice.TopupResponseMapper
	requestCounter           *prometheus.CounterVec
	requestDuration          *prometheus.HistogramVec
}

func NewTopupStasticService(ctx context.Context, mencache mencache.TopupStatisticCache,
	errorhandler errorhandler.TopupStatisticErrorHandler, topupStatistic repository.TopupStatisticRepository, logger logger.LoggerInterface, mapping responseservice.TopupResponseMapper) *topupStasticService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "topup_statistic_service_request_total",
			Help: "Total number of requests to the TopupStatisticService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "topup_statistic_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TopupStatisticService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &topupStasticService{
		ctx:                      ctx,
		mencache:                 mencache,
		errorhandler:             errorhandler,
		trace:                    otel.Tracer("topup-statistic-service"),
		topupStatisticRepository: topupStatistic,
		logger:                   logger,
		mapping:                  mapping,
		requestCounter:           requestCounter,
		requestDuration:          requestDuration,
	}
}

func (s *topupStasticService) FindMonthTopupStatusSuccess(req *requests.MonthTopupStatus) ([]*response.TopupResponseMonthStatusSuccess, *response.ErrorResponse) {

	year := req.Year
	month := req.Month

	const method = "FindMonthTopupStatusSuccess"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("month", month))

	defer func() {
		end(status)
	}()

	if data := s.mencache.GetMonthTopupStatusSuccessCache(req); data != nil {
		logSuccess("Successfully fetched monthly topup status success from cache", zap.Int("year", year), zap.Int("month", month))
		return data, nil
	}

	records, err := s.topupStatisticRepository.GetMonthTopupStatusSuccess(req)

	if err != nil {
		return s.errorhandler.HandleMonthTopupStatusSuccess(err, method, "FAILED_MONTHLY_TOPUP_STATUS_SUCCESS", span, &status, zap.Int("year", year), zap.Int("month", month))
	}
	so := s.mapping.ToTopupResponsesMonthStatusSuccess(records)

	s.mencache.SetMonthTopupStatusSuccessCache(req, so)

	logSuccess("Successfully fetched monthly topup status success", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

func (s *topupStasticService) FindYearlyTopupStatusSuccess(year int) ([]*response.TopupResponseYearStatusSuccess, *response.ErrorResponse) {
	const method = "FindYearlyTopupStatusSuccess"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data := s.mencache.GetYearlyTopupStatusSuccessCache(year); data != nil {
		logSuccess("Successfully fetched yearly topup status success from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.topupStatisticRepository.GetYearlyTopupStatusSuccess(year)
	if err != nil {
		return s.errorhandler.HandleYearlyTopupStatusSuccess(err, method, "FAILED_FIND_YEARLY_TOPUP_STATUS_SUCCESS", span, &status, zap.Error(err))
	}
	so := s.mapping.ToTopupResponsesYearStatusSuccess(records)

	s.mencache.SetYearlyTopupStatusSuccessCache(year, so)

	logSuccess("Successfully fetched yearly topup status success", zap.Int("year", year))

	return so, nil
}

func (s *topupStasticService) FindMonthTopupStatusFailed(req *requests.MonthTopupStatus) ([]*response.TopupResponseMonthStatusFailed, *response.ErrorResponse) {
	year := req.Year
	month := req.Month

	const method = "FindMonthTopupStatusFailed"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("month", month))

	defer func() {
		end(status)
	}()

	if data := s.mencache.GetMonthTopupStatusFailedCache(req); data != nil {
		logSuccess("Successfully fetched monthly topup status Failed from cache", zap.Int("year", year), zap.Int("month", month))
		return data, nil
	}

	records, err := s.topupStatisticRepository.GetMonthTopupStatusFailed(req)
	if err != nil {
		return s.errorhandler.HandleMonthTopupStatusFailed(err, method, "FAILED_MONTHLY_TOPUP_STATUS_FAILED", span, &status, zap.Error(err))
	}
	so := s.mapping.ToTopupResponsesMonthStatusFailed(records)

	s.mencache.SetMonthTopupStatusFailedCache(req, so)

	logSuccess("Successfully fetched monthly topup status Failed", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

func (s *topupStasticService) FindYearlyTopupStatusFailed(year int) ([]*response.TopupResponseYearStatusFailed, *response.ErrorResponse) {
	const method = "FindYearlyTopupStatusFailed"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data := s.mencache.GetYearlyTopupStatusFailedCache(year); data != nil {
		logSuccess("Successfully fetched yearly topup status Failed from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.topupStatisticRepository.GetYearlyTopupStatusFailed(year)
	if err != nil {
		return s.errorhandler.HandleYearlyTopupStatusFailed(err, method, "FAILED_FIND_YEARLY_TOPUP_STATUS_FAILED", span, &status, zap.Int("year", year))
	}
	so := s.mapping.ToTopupResponsesYearStatusFailed(records)

	s.mencache.SetYearlyTopupStatusFailedCache(year, so)

	logSuccess("Successfully fetched yearly topup status Failed", zap.Int("year", year))

	return so, nil
}

func (s *topupStasticService) FindMonthlyTopupMethods(year int) ([]*response.TopupMonthMethodResponse, *response.ErrorResponse) {
	const method = "FindMonthlyTopupMethods"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data := s.mencache.GetMonthlyTopupMethodsCache(year); data != nil {
		logSuccess("Successfully fetched monthly topup methods from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.topupStatisticRepository.GetMonthlyTopupMethods(year)
	if err != nil {
		return s.errorhandler.HandleMonthlyTopupMethods(err, method, "FAILED_FIND_MONTHLY_TOPUP_METHODS", span, &status, zap.Error(err))
	}

	responses := s.mapping.ToTopupMonthlyMethodResponses(records)

	s.mencache.SetMonthlyTopupMethodsCache(year, responses)

	logSuccess("Successfully fetched monthly topup methods", zap.Int("year", year))

	return responses, nil
}

func (s *topupStasticService) FindYearlyTopupMethods(year int) ([]*response.TopupYearlyMethodResponse, *response.ErrorResponse) {
	const method = "FindYearlyTopupMethods"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data := s.mencache.GetYearlyTopupMethodsCache(year); data != nil {
		logSuccess("Successfully fetched yearly topup methods from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.topupStatisticRepository.GetYearlyTopupMethods(year)
	if err != nil {
		return s.errorhandler.HandleYearlyTopupMethods(err, method, "FAILED_FIND_YEARLY_TOPUP_METHODS", span, &status, zap.Error(err))
	}

	responses := s.mapping.ToTopupYearlyMethodResponses(records)

	s.mencache.SetYearlyTopupMethodsCache(year, responses)

	logSuccess("Successfully fetched yearly topup methods", zap.Int("year", year))

	return responses, nil
}

func (s *topupStasticService) FindMonthlyTopupAmounts(year int) ([]*response.TopupMonthAmountResponse, *response.ErrorResponse) {
	const method = "FindMonthlyTopupAmounts"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data := s.mencache.GetMonthlyTopupAmountsCache(year); data != nil {
		logSuccess("Successfully fetched monthly topup amounts from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.topupStatisticRepository.GetMonthlyTopupAmounts(year)
	if err != nil {
		return s.errorhandler.HandleMonthlyTopupAmounts(err, method, "FAILED_FIND_MONTHLY_TOPUP_AMOUNT", span, &status, zap.Error(err))
	}

	responses := s.mapping.ToTopupMonthlyAmountResponses(records)

	s.mencache.SetMonthlyTopupAmountsCache(year, responses)

	logSuccess("Successfully fetched monthly topup amounts", zap.Int("year", year))

	return responses, nil
}

func (s *topupStasticService) FindYearlyTopupAmounts(year int) ([]*response.TopupYearlyAmountResponse, *response.ErrorResponse) {
	const method = "FindYearlyTopupAmounts"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data := s.mencache.GetYearlyTopupAmountsCache(year); data != nil {
		logSuccess("Successfully fetched yearly topup amounts from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.topupStatisticRepository.GetYearlyTopupAmounts(year)
	if err != nil {
		return s.errorhandler.HandleYearlyTopupAmounts(err, method, "FAILED_FIND_YEARLY_TOPUP_AMOUNTS", span, &status, zap.Error(err))
	}

	responses := s.mapping.ToTopupYearlyAmountResponses(records)

	s.mencache.SetYearlyTopupAmountsCache(year, responses)

	logSuccess("Successfully fetched yearly topup amounts", zap.Int("year", year))

	return responses, nil
}

func (s *topupStasticService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *topupStasticService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
