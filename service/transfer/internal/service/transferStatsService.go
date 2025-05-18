package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	traceunic "github.com/MamangRust/monolith-payment-gateway-pkg/trace_unic"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transferStatisticService struct {
	ctx                         context.Context
	trace                       trace.Tracer
	transferStatisticRepository repository.TransferStatisticRepository
	logger                      logger.LoggerInterface
	mapping                     responseservice.TransferResponseMapper
	requestCounter              *prometheus.CounterVec
	requestDuration             *prometheus.HistogramVec
}

func NewTransferStatisticService(ctx context.Context, transferStatisticRepository repository.TransferStatisticRepository, logger logger.LoggerInterface, mapping responseservice.TransferResponseMapper) *transferStatisticService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transfer_statistic_service_request_total",
			Help: "Total number of requests to the TransferStatisticService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transfer_statistic_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TransferStatisticService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &transferStatisticService{
		ctx:                         ctx,
		trace:                       otel.Tracer("transfer-statistic-service"),
		transferStatisticRepository: transferStatisticRepository,
		logger:                      logger,
		mapping:                     mapping,
		requestCounter:              requestCounter,
		requestDuration:             requestDuration,
	}
}

func (s *transferStatisticService) FindMonthTransferStatusSuccess(req *requests.MonthStatusTransfer) ([]*response.TransferResponseMonthStatusSuccess, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthTransferStatusSuccess", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthTransferStatusSuccess")
	defer span.End()

	year := req.Year
	month := req.Month

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
	)

	s.logger.Debug("Fetching monthly Transfer status success", zap.Int("year", year), zap.Int("month", month))

	records, err := s.transferStatisticRepository.GetMonthTransferStatusSuccess(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTH_TRANSFER_SUCCESS")

		s.logger.Error("Failed to fetch monthly Transfer status success", zap.Error(err), zap.Int("year", year), zap.Int("month", month))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)

		status = "failed_find_month_transfer_status_success"

		return nil, transfer_errors.ErrFailedFindMonthTransferStatusSuccess
	}

	s.logger.Debug("Successfully fetched monthly Transfer status success", zap.Int("year", year), zap.Int("month", month))

	so := s.mapping.ToTransferResponsesMonthStatusSuccess(records)

	return so, nil
}

func (s *transferStatisticService) FindYearlyTransferStatusSuccess(year int) ([]*response.TransferResponseYearStatusSuccess, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTransferStatusSuccess", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTransferStatusSuccess")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Fetching yearly Transfer status success", zap.Int("year", year))

	records, err := s.transferStatisticRepository.GetYearlyTransferStatusSuccess(year)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEAR_TRANSFER_SUCCESS")

		s.logger.Error("Failed to fetch yearly Transfer status success", zap.Error(err), zap.Int("year", year))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)

		status = "failed_find_year_transfer_status_success"

		return nil, transfer_errors.ErrFailedFindYearTransferStatusSuccess
	}

	s.logger.Debug("Successfully fetched yearly Transfer status success", zap.Int("year", year))

	so := s.mapping.ToTransferResponsesYearStatusSuccess(records)

	return so, nil
}

func (s *transferStatisticService) FindMonthTransferStatusFailed(req *requests.MonthStatusTransfer) ([]*response.TransferResponseMonthStatusFailed, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthTransferStatusFailed", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthTransferStatusFailed")
	defer span.End()

	year := req.Year
	month := req.Month

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
	)

	s.logger.Debug("Fetching monthly Transfer status Failed", zap.Int("year", year), zap.Int("month", month))

	records, err := s.transferStatisticRepository.GetMonthTransferStatusFailed(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTH_TRANSFER_STATUS_FAILED")

		s.logger.Error("Failed to fetch monthly Transfer status Failed", zap.Error(err), zap.Int("year", year), zap.Int("month", month))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)

		status = "failed_find_month_transfer_status_failed"

		return nil, transfer_errors.ErrFailedFindMonthTransferStatusFailed
	}

	s.logger.Debug("Failedfully fetched monthly Transfer status Failed", zap.Int("year", year), zap.Int("month", month))

	so := s.mapping.ToTransferResponsesMonthStatusFailed(records)

	return so, nil
}

func (s *transferStatisticService) FindYearlyTransferStatusFailed(year int) ([]*response.TransferResponseYearStatusFailed, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTransferStatusFailed", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTransferStatusFailed")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Fetching yearly Transfer status Failed", zap.Int("year", year))

	records, err := s.transferStatisticRepository.GetYearlyTransferStatusFailed(year)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEAR_TRANSFER_STATUS_FAILED")

		s.logger.Error("Failed to fetch yearly Transfer status Failed", zap.Error(err), zap.Int("year", year))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)

		status = "failed_find_year_transfer_status_failed"

		return nil, transfer_errors.ErrFailedFindYearTransferStatusFailed
	}

	s.logger.Debug("Failedfully fetched yearly Transfer status Failed", zap.Int("year", year))

	so := s.mapping.ToTransferResponsesYearStatusFailed(records)

	return so, nil
}

func (s *transferStatisticService) FindMonthlyTransferAmounts(year int) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyTransferAmounts", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyTransferAmounts")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Fetching monthly transfer amounts", zap.Int("year", year))

	amounts, err := s.transferStatisticRepository.GetMonthlyTransferAmounts(year)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_TRANSFER_AMOUNTS")

		s.logger.Error("failed to find monthly transfer amounts", zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)

		status = "failed_find_monthly_transfer_amounts"

		return nil, transfer_errors.ErrFailedFindMonthlyTransferAmounts
	}

	responseAmounts := s.mapping.ToTransferResponsesMonthAmount(amounts)

	s.logger.Debug("Successfully fetched monthly transfer amounts", zap.Int("year", year))

	return responseAmounts, nil
}

func (s *transferStatisticService) FindYearlyTransferAmounts(year int) ([]*response.TransferYearAmountResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTransferAmounts", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTransferAmounts")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Fetching yearly transfer amounts", zap.Int("year", year))

	amounts, err := s.transferStatisticRepository.GetYearlyTransferAmounts(year)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_TRANSFER_AMOUNTS")

		s.logger.Error("failed to find yearly transfer amounts", zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)

		status = "failed_find_yearly_transfer_amounts"

		return nil, transfer_errors.ErrFailedFindYearlyTransferAmounts
	}

	responseAmounts := s.mapping.ToTransferResponsesYearAmount(amounts)

	s.logger.Debug("Successfully fetched yearly transfer amounts", zap.Int("year", year))

	return responseAmounts, nil
}

func (s *transferStatisticService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
