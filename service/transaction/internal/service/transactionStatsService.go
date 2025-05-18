package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	traceunic "github.com/MamangRust/monolith-payment-gateway-pkg/trace_unic"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transactionStatisticService struct {
	ctx                            context.Context
	trace                          trace.Tracer
	transactionStatisticRepository repository.TransactionStatisticsRepository
	logger                         logger.LoggerInterface
	mapping                        responseservice.TransactionResponseMapper
	requestCounter                 *prometheus.CounterVec
	requestDuration                *prometheus.HistogramVec
}

func NewTransactionStatisticService(ctx context.Context, transactionStatisticRepository repository.TransactionStatisticsRepository, logger logger.LoggerInterface, mapping responseservice.TransactionResponseMapper) *transactionStatisticService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transaction_statistic_service_request_total",
			Help: "Total number of requests to the TransactionStatisticService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transaction_statistic_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TransactionStatisticService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &transactionStatisticService{
		ctx:                            ctx,
		trace:                          otel.Tracer("transaction-statistic-service"),
		transactionStatisticRepository: transactionStatisticRepository,
		logger:                         logger,
		mapping:                        mapping,
		requestCounter:                 requestCounter,
		requestDuration:                requestDuration,
	}
}

func (s *transactionStatisticService) FindMonthTransactionStatusSuccess(req *requests.MonthStatusTransaction) ([]*response.TransactionResponseMonthStatusSuccess, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthTransactionStatusSuccess", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthTransactionStatusSuccess")
	defer span.End()

	month := req.Month
	year := req.Year

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
	)

	s.logger.Debug("Fetching monthly Transaction status success", zap.Int("year", year), zap.Int("month", month))

	records, err := s.transactionStatisticRepository.GetMonthTransactionStatusSuccess(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTH_TRANSACTION_SUCCESS")

		s.logger.Error("failed to fetch monthly Transaction status success", zap.Error(err), zap.Int("year", year), zap.Int("month", month))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch monthly Transaction status success")

		status = "failed_find_month_transaction_success"

		return nil, transaction_errors.ErrFailedFindMonthTransactionSuccess
	}

	s.logger.Debug("Successfully fetched monthly Transaction status success", zap.Int("year", year), zap.Int("month", month))

	so := s.mapping.ToTransactionResponsesMonthStatusSuccess(records)

	return so, nil
}

func (s *transactionStatisticService) FindYearlyTransactionStatusSuccess(year int) ([]*response.TransactionResponseYearStatusSuccess, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTransactionStatusSuccess", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTransactionStatusSuccess")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Fetching yearly Transaction status success", zap.Int("year", year))

	records, err := s.transactionStatisticRepository.GetYearlyTransactionStatusSuccess(year)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEAR_TRANSACTION_SUCCESS")

		s.logger.Error("failed to fetch yearly Transaction status success", zap.Error(err), zap.Int("year", year))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch yearly Transaction status success")

		status = "failed_find_year_transaction_success"

		return nil, transaction_errors.ErrFailedFindYearTransactionSuccess
	}

	s.logger.Debug("Successfully fetched yearly Transaction status success", zap.Int("year", year))

	so := s.mapping.ToTransactionResponsesYearStatusSuccess(records)

	return so, nil
}

func (s *transactionStatisticService) FindMonthTransactionStatusFailed(req *requests.MonthStatusTransaction) ([]*response.TransactionResponseMonthStatusFailed, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthTransactionStatusFailed", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthTransactionStatusFailed")
	defer span.End()

	year := req.Year
	month := req.Month

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
	)

	s.logger.Debug("Fetching monthly Transaction status Failed", zap.Int("year", year), zap.Int("month", month))

	records, err := s.transactionStatisticRepository.GetMonthTransactionStatusFailed(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTH_TRANSACTION_FAILED")

		s.logger.Error("failed to fetch monthly Transaction status Failed", zap.Error(err), zap.Int("year", year), zap.Int("month", month))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch monthly Transaction status Failed")

		status = "failed_find_month_transaction_failed"

		return nil, transaction_errors.ErrFailedFindMonthTransactionFailed
	}

	s.logger.Debug("Failedfully fetched monthly Transaction status Failed", zap.Int("year", year), zap.Int("month", month))

	so := s.mapping.ToTransactionResponsesMonthStatusFailed(records)

	return so, nil
}

func (s *transactionStatisticService) FindYearlyTransactionStatusFailed(year int) ([]*response.TransactionResponseYearStatusFailed, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTransactionStatusFailed", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTransactionStatusFailed")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Fetching yearly Transaction status Failed", zap.Int("year", year))

	records, err := s.transactionStatisticRepository.GetYearlyTransactionStatusFailed(year)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEAR_TRANSACTION_FAILED")

		s.logger.Error("failed to fetch yearly Transaction status Failed", zap.Error(err), zap.Int("year", year))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch yearly Transaction status Failed")

		status = "failed_find_year_transaction_failed"

		return nil, transaction_errors.ErrFailedFindYearTransactionFailed
	}

	s.logger.Debug("Failedfully fetched yearly Transaction status Failed", zap.Int("year", year))

	so := s.mapping.ToTransactionResponsesYearStatusFailed(records)

	return so, nil
}

func (s *transactionStatisticService) FindMonthlyPaymentMethods(year int) ([]*response.TransactionMonthMethodResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyPaymentMethods", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyPaymentMethods")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Fetching monthly payment methods", zap.Int("year", year))

	records, err := s.transactionStatisticRepository.GetMonthlyPaymentMethods(year)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_PAYMENT_METHODS")

		s.logger.Error("Failed to fetch monthly payment methods", zap.Error(err), zap.Int("year", year))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch monthly payment methods")

		status = "failed_to_fetch_monthly_payment_methods"

		return nil, transaction_errors.ErrFailedFindMonthlyPaymentMethods
	}

	responses := s.mapping.ToTransactionMonthlyMethodResponses(records)

	s.logger.Debug("Successfully fetched monthly payment methods", zap.Int("year", year))

	return responses, nil
}

func (s *transactionStatisticService) FindYearlyPaymentMethods(year int) ([]*response.TransactionYearMethodResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyPaymentMethods", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyPaymentMethods")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Fetching yearly payment methods", zap.Int("year", year))

	records, err := s.transactionStatisticRepository.GetYearlyPaymentMethods(year)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_PAYMENT_METHODS")

		s.logger.Error("Failed to fetch yearly payment methods", zap.Error(err), zap.Int("year", year))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch yearly payment methods")

		status = "failed_to_fetch_yearly_payment_methods"

		return nil, transaction_errors.ErrFailedFindYearlyPaymentMethods
	}

	responses := s.mapping.ToTransactionYearlyMethodResponses(records)

	s.logger.Debug("Successfully fetched yearly payment methods", zap.Int("year", year))

	return responses, nil
}

func (s *transactionStatisticService) FindMonthlyAmounts(year int) ([]*response.TransactionMonthAmountResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyAmounts", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyAmounts")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Fetching monthly amounts", zap.Int("year", year))

	records, err := s.transactionStatisticRepository.GetMonthlyAmounts(year)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_AMOUNTS")

		s.logger.Error("Failed to fetch monthly amounts", zap.Error(err), zap.Int("year", year))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch monthly amounts")

		status = "failed_to_fetch_monthly_amounts"

		return nil, transaction_errors.ErrFailedFindMonthlyAmounts
	}

	responses := s.mapping.ToTransactionMonthlyAmountResponses(records)

	s.logger.Debug("Successfully fetched monthly amounts", zap.Int("year", year))

	return responses, nil
}

func (s *transactionStatisticService) FindYearlyAmounts(year int) ([]*response.TransactionYearlyAmountResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyAmounts", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyAmounts")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Fetching yearly amounts", zap.Int("year", year))

	records, err := s.transactionStatisticRepository.GetYearlyAmounts(year)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_AMOUNTS")

		s.logger.Error("Failed to fetch yearly amounts", zap.Error(err), zap.Int("year", year))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch yearly amounts")

		status = "failed_to_fetch_yearly_amounts"

		return nil, transaction_errors.ErrFailedFindYearlyAmounts
	}

	responses := s.mapping.ToTransactionYearlyAmountResponses(records)

	s.logger.Debug("Successfully fetched yearly amounts", zap.Int("year", year))

	return responses, nil
}

func (s *transactionStatisticService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
