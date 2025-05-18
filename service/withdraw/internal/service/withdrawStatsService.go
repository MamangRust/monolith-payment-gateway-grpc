package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	traceunic "github.com/MamangRust/monolith-payment-gateway-pkg/trace_unic"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
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
	trace                       trace.Tracer
	withdrawStatisticRepository repository.WithdrawStatisticRepository
	logger                      logger.LoggerInterface
	mapping                     responseservice.WithdrawResponseMapper
	requestCounter              *prometheus.CounterVec
	requestDuration             *prometheus.HistogramVec
}

func NewWithdrawStatisticService(ctx context.Context, withdrawStatisticRepository repository.WithdrawStatisticRepository, logger logger.LoggerInterface, mapping responseservice.WithdrawResponseMapper) *withdrawStatisticService {
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
		[]string{"method"},
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

	records, err := s.withdrawStatisticRepository.GetMonthWithdrawStatusSuccess(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTH_WITHDRAW_SUCCESS")

		s.logger.Error("failed to fetch monthly Withdraw status success", zap.Error(err), zap.String("traceID", traceID))

		span.SetAttributes(attribute.String("traceID", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to fetch monthly Withdraw status success")

		status = "failed_to_fetch_monthly_withdraw_status_success"

		return nil, withdraw_errors.ErrFailedFindMonthWithdrawStatusSuccess
	}

	s.logger.Debug("Successfully fetched monthly Withdraw status success", zap.Int("year", year), zap.Int("month", month))

	so := s.mapping.ToWithdrawResponsesMonthStatusSuccess(records)

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

	records, err := s.withdrawStatisticRepository.GetYearlyWithdrawStatusSuccess(year)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEAR_WITHDRAW_SUCCESS")

		s.logger.Error("failed to fetch yearly Withdraw status success", zap.Error(err), zap.String("traceID", traceID))

		span.SetAttributes(attribute.String("traceID", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to fetch yearly Withdraw status success")

		status = "failed_to_fetch_yearly_withdraw_status_success"

		return nil, withdraw_errors.ErrFailedFindYearWithdrawStatusSuccess
	}

	s.logger.Debug("Successfully fetched yearly Withdraw status success", zap.Int("year", year))

	so := s.mapping.ToWithdrawResponsesYearStatusSuccess(records)

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

	span.SetAttributes(
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month),
	)

	year := req.Year
	month := req.Month

	s.logger.Debug("Fetching monthly Withdraw status Failed", zap.Int("year", year), zap.Int("month", month))

	records, err := s.withdrawStatisticRepository.GetMonthWithdrawStatusFailed(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTH_WITHDRAW_FAILED")

		s.logger.Error("failed to fetch monthly Withdraw status Failed", zap.Error(err), zap.String("traceID", traceID))

		span.SetAttributes(attribute.String("traceID", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to fetch monthly Withdraw status Failed")

		status = "failed_to_fetch_monthly_withdraw_status_failed"

		return nil, withdraw_errors.ErrFailedFindMonthWithdrawStatusFailed
	}

	s.logger.Debug("Failedfully fetched monthly Withdraw status Failed", zap.Int("year", year), zap.Int("month", month))

	so := s.mapping.ToWithdrawResponsesMonthStatusFailed(records)

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

	records, err := s.withdrawStatisticRepository.GetYearlyWithdrawStatusFailed(year)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEAR_WITHDRAW_FAILED")

		s.logger.Error("failed to fetch yearly Withdraw status Failed", zap.Error(err), zap.String("traceID", traceID))

		span.SetAttributes(attribute.String("traceID", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to fetch yearly Withdraw status Failed")

		status = "failed_to_fetch_yearly_withdraw_status_failed"

		return nil, withdraw_errors.ErrFailedFindYearWithdrawStatusFailed
	}

	s.logger.Debug("Failedfully fetched yearly Withdraw status Failed", zap.Int("year", year))

	so := s.mapping.ToWithdrawResponsesYearStatusFailed(records)

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

	s.logger.Debug("Fetching monthly withdraws", zap.Int("year", year))

	withdraws, err := s.withdrawStatisticRepository.GetMonthlyWithdraws(year)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTH_WITHDRAW")

		s.logger.Error("failed to find monthly withdraws", zap.Error(err), zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to find monthly withdraws")
		status = "failed_to_find_monthly_withdraws"

		return nil, withdraw_errors.ErrFailedFindMonthlyWithdraws
	}

	responseWithdraws := s.mapping.ToWithdrawsAmountMonthlyResponses(withdraws)

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

	withdraws, err := s.withdrawStatisticRepository.GetYearlyWithdraws(year)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEAR_WITHDRAW")

		s.logger.Error("failed to find yearly withdraws", zap.Error(err), zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to find yearly withdraws")
		status = "failed_to_find_yearly_withdraws"

		return nil, withdraw_errors.ErrFailedFindYearlyWithdraws
	}

	responseWithdraws := s.mapping.ToWithdrawsAmountYearlyResponses(withdraws)

	s.logger.Debug("Successfully fetched yearly withdraws", zap.Int("year", year))

	return responseWithdraws, nil
}

func (s *withdrawStatisticService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
