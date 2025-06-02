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
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthTopupStatusSuccess", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthTopupStatusSuccess")
	defer span.End()

	year := req.Year
	month := req.Month

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
	)

	s.logger.Debug("Fetching monthly topup status success", zap.Int("year", year), zap.Int("month", month))

	if data := s.mencache.GetMonthTopupStatusSuccessCache(req); data != nil {
		s.logger.Debug("Successfully fetched monthly topup status success from cache", zap.Int("year", year), zap.Int("month", month))
		return data, nil
	}

	records, err := s.topupStatisticRepository.GetMonthTopupStatusSuccess(req)

	if err != nil {
		return s.errorhandler.HandleMonthTopupStatusSuccess(err, "FindMonthTopupStatusSuccess", "FAILED_MONTHLY_TOPUP_STATUS_SUCCESS", span, &status, zap.Int("year", year), zap.Int("month", month))
	}
	so := s.mapping.ToTopupResponsesMonthStatusSuccess(records)

	s.mencache.SetMonthTopupStatusSuccessCache(req, so)

	s.logger.Debug("Successfully fetched monthly topup status success", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

func (s *topupStasticService) FindYearlyTopupStatusSuccess(year int) ([]*response.TopupResponseYearStatusSuccess, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTopupStatusSuccess", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTopupStatusSuccess")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Fetching yearly topup status success", zap.Int("year", year))

	if data := s.mencache.GetYearlyTopupStatusSuccessCache(year); data != nil {
		s.logger.Debug("Successfully fetched yearly topup status success", zap.Int("year", year))
		return data, nil
	}

	records, err := s.topupStatisticRepository.GetYearlyTopupStatusSuccess(year)
	if err != nil {
		return s.errorhandler.HandleYearlyTopupStatusSuccess(err, "FindYearlyTopupStatusSuccess", "FAILED_FIND_YEARLY_TOPUP_STATUS_SUCCESS", span, &status, zap.Int("year", year))
	}
	so := s.mapping.ToTopupResponsesYearStatusSuccess(records)

	s.mencache.SetYearlyTopupStatusSuccessCache(year, so)

	s.logger.Debug("Successfully fetched yearly topup status success", zap.Int("year", year))

	return so, nil
}

func (s *topupStasticService) FindMonthTopupStatusFailed(req *requests.MonthTopupStatus) ([]*response.TopupResponseMonthStatusFailed, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthTopupStatusFailed", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthTopupStatusFailed")
	defer span.End()

	year := req.Year
	month := req.Month

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
	)

	s.logger.Debug("Fetching monthly topup status Failed", zap.Int("year", year), zap.Int("month", month))

	if data := s.mencache.GetMonthTopupStatusFailedCache(req); data != nil {
		s.logger.Debug("Successfully fetched monthly topup status Failed from cache", zap.Int("year", year), zap.Int("month", month))
		return data, nil
	}

	records, err := s.topupStatisticRepository.GetMonthTopupStatusFailed(req)
	if err != nil {
		return s.errorhandler.HandleMonthTopupStatusFailed(err, "FindMonthTopupStatusFailed", "FAILED_MONTHLY_TOPUP_STATUS_FAILED", span, &status, zap.Int("year", year), zap.Int("month", month))
	}
	so := s.mapping.ToTopupResponsesMonthStatusFailed(records)

	s.mencache.SetMonthTopupStatusFailedCache(req, so)

	s.logger.Debug("Failedfully fetched monthly topup status Failed", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

func (s *topupStasticService) FindYearlyTopupStatusFailed(year int) ([]*response.TopupResponseYearStatusFailed, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTopupStatusFailed", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTopupStatusFailed")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Fetching yearly topup status Failed", zap.Int("year", year))

	if data := s.mencache.GetYearlyTopupStatusFailedCache(year); data != nil {
		s.logger.Debug("Successfully fetched yearly topup status Failed", zap.Int("year", year))
		return data, nil
	}

	records, err := s.topupStatisticRepository.GetYearlyTopupStatusFailed(year)
	if err != nil {
		return s.errorhandler.HandleYearlyTopupStatusFailed(err, "FindYearlyTopupStatusFailed", "FAILED_FIND_YEARLY_TOPUP_STATUS_FAILED", span, &status, zap.Int("year", year))
	}
	so := s.mapping.ToTopupResponsesYearStatusFailed(records)

	s.mencache.SetYearlyTopupStatusFailedCache(year, so)

	s.logger.Debug("Failedfully fetched yearly topup status Failed", zap.Int("year", year))

	return so, nil
}

func (s *topupStasticService) FindMonthlyTopupMethods(year int) ([]*response.TopupMonthMethodResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyTopupMethods", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyTopupMethods")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Fetching monthly topup methods", zap.Int("year", year))

	if data := s.mencache.GetMonthlyTopupMethodsCache(year); data != nil {
		s.logger.Debug("Successfully fetched monthly topup methods from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.topupStatisticRepository.GetMonthlyTopupMethods(year)
	if err != nil {
		return s.errorhandler.HandleMonthlyTopupMethods(err, "FindMonthlyTopupMethods", "FAILED_FIND_MONTHLY_TOPUP_METHODS", span, &status, zap.Int("year", year))
	}

	responses := s.mapping.ToTopupMonthlyMethodResponses(records)

	s.mencache.SetMonthlyTopupMethodsCache(year, responses)

	s.logger.Debug("Successfully fetched monthly topup methods", zap.Int("year", year))

	return responses, nil
}

func (s *topupStasticService) FindYearlyTopupMethods(year int) ([]*response.TopupYearlyMethodResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTopupMethods", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTopupMethods")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Fetching yearly topup methods", zap.Int("year", year))

	if data := s.mencache.GetYearlyTopupMethodsCache(year); data != nil {
		s.logger.Debug("Successfully fetched yearly topup methods", zap.Int("year", year))
		return data, nil
	}

	records, err := s.topupStatisticRepository.GetYearlyTopupMethods(year)
	if err != nil {
		return s.errorhandler.HandleYearlyTopupMethods(err, "FindYearlyTopupMethods", "FAILED_FIND_YEARLY_TOPUP_METHODS", span, &status, zap.Int("year", year))
	}

	responses := s.mapping.ToTopupYearlyMethodResponses(records)

	s.mencache.SetYearlyTopupMethodsCache(year, responses)

	s.logger.Debug("Successfully fetched yearly topup methods", zap.Int("year", year))

	return responses, nil
}

func (s *topupStasticService) FindMonthlyTopupAmounts(year int) ([]*response.TopupMonthAmountResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyTopupAmounts", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyTopupAmounts")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Fetching monthly topup amounts", zap.Int("year", year))

	if data := s.mencache.GetMonthlyTopupAmountsCache(year); data != nil {
		s.logger.Debug("Successfully fetched monthly topup amounts from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.topupStatisticRepository.GetMonthlyTopupAmounts(year)
	if err != nil {
		return s.errorhandler.HandleMonthlyTopupAmounts(err, "FindMonthlyTopupAmounts", "FAILED_FIND_MONTHLY_TOPUP_AMOUNT", span, &status, zap.Int("year", year))
	}

	responses := s.mapping.ToTopupMonthlyAmountResponses(records)

	s.mencache.SetMonthlyTopupAmountsCache(year, responses)

	s.logger.Debug("Successfully fetched monthly topup amounts", zap.Int("year", year))

	return responses, nil
}

func (s *topupStasticService) FindYearlyTopupAmounts(year int) ([]*response.TopupYearlyAmountResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTopupAmounts", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTopupAmounts")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Fetching yearly topup amounts", zap.Int("year", year))

	if data := s.mencache.GetYearlyTopupAmountsCache(year); data != nil {
		s.logger.Debug("Successfully fetched yearly topup amounts", zap.Int("year", year))
		return data, nil
	}

	records, err := s.topupStatisticRepository.GetYearlyTopupAmounts(year)
	if err != nil {
		return s.errorhandler.HandleYearlyTopupAmounts(err, "FindYearlyTopupAmounts", "FAILED_FIND_YEARLY_TOPUP_AMOUNTS", span, &status, zap.Int("year", year))
	}

	responses := s.mapping.ToTopupYearlyAmountResponses(records)

	s.mencache.SetYearlyTopupAmountsCache(year, responses)

	s.logger.Debug("Successfully fetched yearly topup amounts", zap.Int("year", year))

	return responses, nil
}

func (s *topupStasticService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
