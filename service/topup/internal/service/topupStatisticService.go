package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	traceunic "github.com/MamangRust/monolith-payment-gateway-pkg/trace_unic"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
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
	trace                    trace.Tracer
	topupStatisticRepository repository.TopupStatisticRepository
	logger                   logger.LoggerInterface
	mapping                  responseservice.TopupResponseMapper
	requestCounter           *prometheus.CounterVec
	requestDuration          *prometheus.HistogramVec
}

func NewTopupStasticService(ctx context.Context, topupStatistic repository.TopupStatisticRepository, logger logger.LoggerInterface, mapping responseservice.TopupResponseMapper) *topupStasticService {
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
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &topupStasticService{
		ctx:                      ctx,
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

	records, err := s.topupStatisticRepository.GetMonthTopupStatusSuccess(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTH_TOPUP_STATUS_SUCCESS")

		s.logger.Error("Failed to fetch monthly topup status success", zap.Error(err), zap.Int("year", year), zap.Int("month", month))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)

		status = "failed_find_month_topup_status_success"
		return nil, topup_errors.ErrFailedFindMonthTopupStatusSuccess
	}

	s.logger.Debug("Successfully fetched monthly topup status success", zap.Int("year", year), zap.Int("month", month))

	so := s.mapping.ToTopupResponsesMonthStatusSuccess(records)

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

	records, err := s.topupStatisticRepository.GetYearlyTopupStatusSuccess(year)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_TOPUP_STATUS_SUCCESS")

		s.logger.Error("Failed to fetch yearly topup status success", zap.Error(err), zap.Int("year", year))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch yearly topup status success")
		status = "failed_find_yearly_topup_status_success"
		return nil, topup_errors.ErrFailedFindYearlyTopupStatusSuccess
	}

	s.logger.Debug("Successfully fetched yearly topup status success", zap.Int("year", year))

	so := s.mapping.ToTopupResponsesYearStatusSuccess(records)

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

	records, err := s.topupStatisticRepository.GetMonthTopupStatusFailed(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTH_TOPUP_STATUS_FAILED")

		s.logger.Error("Failed to fetch monthly topup status Failed", zap.Error(err), zap.Int("year", year), zap.Int("month", month))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)

		status = "failed_find_month_topup_status_failed"

		return nil, topup_errors.ErrFailedFindMonthTopupStatusFailed
	}

	s.logger.Debug("Failedfully fetched monthly topup status Failed", zap.Int("year", year), zap.Int("month", month))

	so := s.mapping.ToTopupResponsesMonthStatusFailed(records)

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

	records, err := s.topupStatisticRepository.GetYearlyTopupStatusFailed(year)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_TOPUP_STATUS_FAILED")

		s.logger.Error("Failed to fetch yearly topup status Failed", zap.Error(err), zap.Int("year", year))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch yearly topup status Failed")
		status = "failed_find_yearly_topup_status_failed"

		return nil, topup_errors.ErrFailedFindYearlyTopupStatusFailed
	}

	s.logger.Debug("Failedfully fetched yearly topup status Failed", zap.Int("year", year))

	so := s.mapping.ToTopupResponsesYearStatusFailed(records)

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

	records, err := s.topupStatisticRepository.GetMonthlyTopupMethods(year)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_TOPUP_METHODS")

		s.logger.Error("Failed to fetch monthly topup methods", zap.Error(err), zap.Int("year", year))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch monthly topup methods")
		status = "failed_find_monthly_topup_methods"

		return nil, topup_errors.ErrFailedFindMonthlyTopupMethods
	}

	responses := s.mapping.ToTopupMonthlyMethodResponses(records)

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

	records, err := s.topupStatisticRepository.GetYearlyTopupMethods(year)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_TOPUP_METHODS")

		s.logger.Error("Failed to fetch yearly topup methods", zap.Error(err), zap.Int("year", year))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch yearly topup methods")
		status = "failed_find_yearly_topup_methods"

		return nil, topup_errors.ErrFailedFindYearlyTopupMethods
	}

	responses := s.mapping.ToTopupYearlyMethodResponses(records)

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

	records, err := s.topupStatisticRepository.GetMonthlyTopupAmounts(year)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTH_TOPUP_AMOUNTS")

		s.logger.Error("Failed to fetch monthly topup amounts", zap.Error(err), zap.Int("year", year))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch monthly topup amounts")
		status = "failed_find_month_topup_amounts"

		return nil, topup_errors.ErrFailedFindMonthlyTopupAmounts
	}

	responses := s.mapping.ToTopupMonthlyAmountResponses(records)

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

	records, err := s.topupStatisticRepository.GetYearlyTopupAmounts(year)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_TOPUP_AMOUNTS")

		s.logger.Error("Failed to fetch yearly topup amounts", zap.Error(err), zap.Int("year", year))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch yearly topup amounts")
		status = "failed_find_yearly_topup_amounts"

		return nil, topup_errors.ErrFailedFindYearlyTopupAmounts
	}

	responses := s.mapping.ToTopupYearlyAmountResponses(records)

	s.logger.Debug("Successfully fetched yearly topup amounts", zap.Int("year", year))

	return responses, nil
}

func (s *topupStasticService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
