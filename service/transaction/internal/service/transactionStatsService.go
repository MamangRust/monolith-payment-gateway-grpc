package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-transaction/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transactionStatisticService struct {
	ctx                            context.Context
	mencache                       mencache.TransactonStatistcCache
	errorhandler                   errorhandler.TransactionStatisticErrorHandler
	trace                          trace.Tracer
	transactionStatisticRepository repository.TransactionStatisticsRepository
	logger                         logger.LoggerInterface
	mapping                        responseservice.TransactionResponseMapper
	requestCounter                 *prometheus.CounterVec
	requestDuration                *prometheus.HistogramVec
}

func NewTransactionStatisticService(ctx context.Context, mencache mencache.TransactonStatistcCache,
	errorhandler errorhandler.TransactionStatisticErrorHandler, transactionStatisticRepository repository.TransactionStatisticsRepository, logger logger.LoggerInterface, mapping responseservice.TransactionResponseMapper) *transactionStatisticService {
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
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &transactionStatisticService{
		ctx:                            ctx,
		mencache:                       mencache,
		errorhandler:                   errorhandler,
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

	if data := s.mencache.GetMonthTransactonStatusSuccessCache(req); data != nil {
		s.logger.Debug("Successfully fetched monthly Transaction status success from cache", zap.Int("year", year), zap.Int("month", month))
		return data, nil
	}

	records, err := s.transactionStatisticRepository.GetMonthTransactionStatusSuccess(req)
	if err != nil {
		return s.errorhandler.HandleMonthTransactionStatusSuccessError(err, "FindMonthTransactionStatusSuccess", "FAILED_FIND_MONTH_TRANSACTION_STATUS_SUCCESS", span, &status, zap.Int("year", year), zap.Int("month", month))
	}
	so := s.mapping.ToTransactionResponsesMonthStatusSuccess(records)

	s.mencache.SetMonthTransactonStatusSuccessCache(req, so)

	s.logger.Debug("Successfully fetched monthly Transaction status success", zap.Int("year", year), zap.Int("month", month))

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

	if data := s.mencache.GetYearTransactonStatusSuccessCache(year); data != nil {
		s.logger.Debug("Successfully fetched yearly Transaction status success from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.transactionStatisticRepository.GetYearlyTransactionStatusSuccess(year)
	if err != nil {
		return s.errorhandler.HandleYearlyTransactionStatusSuccessError(err, "FindYearlyTransactionStatusSuccess", "FAILED_FIND_YEARLY_TRANSACTION_STATUS_SUCCESS", span, &status, zap.Int("year", year))
	}
	so := s.mapping.ToTransactionResponsesYearStatusSuccess(records)

	s.mencache.SetYearTransactonStatusSuccessCache(year, so)

	s.logger.Debug("Successfully fetched yearly Transaction status success", zap.Int("year", year))

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

	if data := s.mencache.GetMonthTransactonStatusFailedCache(req); data != nil {
		s.logger.Debug("Successfully fetched monthly Transaction status Failed from cache", zap.Int("year", year), zap.Int("month", month))
		return data, nil
	}

	records, err := s.transactionStatisticRepository.GetMonthTransactionStatusFailed(req)
	if err != nil {
		return s.errorhandler.HandleMonthTransactionStatusFailedError(err, "FindMonthTransactionStatusFailed", "FAILED_MONTHLY_TRANSACTION_STATUS_FAILED", span, &status, zap.Int("year", year), zap.Int("month", month))
	}
	so := s.mapping.ToTransactionResponsesMonthStatusFailed(records)

	s.mencache.SetMonthTransactonStatusFailedCache(req, so)

	s.logger.Debug("Failedfully fetched monthly Transaction status Failed", zap.Int("year", year), zap.Int("month", month))

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

	if data := s.mencache.GetYearTransactonStatusFailedCache(year); data != nil {
		s.logger.Debug("Successfully fetched yearly Transaction status Failed from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.transactionStatisticRepository.GetYearlyTransactionStatusFailed(year)
	if err != nil {
		return s.errorhandler.HandleYearlyTransactionStatusFailedError(err, "FindYearlyTransactionStatusFailed", "FAILED_FIND_YEARLY_TRANSACTION_STATUS_FAILED", span, &status, zap.Int("year", year))
	}
	so := s.mapping.ToTransactionResponsesYearStatusFailed(records)

	s.mencache.SetYearTransactonStatusFailedCache(year, so)

	s.logger.Debug("Failedfully fetched yearly Transaction status Failed", zap.Int("year", year))

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

	if data := s.mencache.GetMonthlyPaymentMethodsCache(year); data != nil {
		s.logger.Debug("Successfully fetched monthly payment methods from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.transactionStatisticRepository.GetMonthlyPaymentMethods(year)
	if err != nil {
		return s.errorhandler.HandleMonthlyPaymentMethodsError(err, "FindMonthlyPaymentMethods", "FAILED_FIND_MONTHLY_PAYMENT_METHODS", span, &status, zap.Int("year", year))
	}

	responses := s.mapping.ToTransactionMonthlyMethodResponses(records)

	s.mencache.SetMonthlyPaymentMethodsCache(year, responses)

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

	if data := s.mencache.GetYearlyPaymentMethodsCache(year); data != nil {
		s.logger.Debug("Successfully fetched yearly payment methods from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.transactionStatisticRepository.GetYearlyPaymentMethods(year)
	if err != nil {
		return s.errorhandler.HandleYearlyPaymentMethodsError(err, "FindYearlyPaymentMethods", "FAILED_FIND_YEARLY_PAYMENT_METHODS", span, &status, zap.Int("year", year))
	}

	responses := s.mapping.ToTransactionYearlyMethodResponses(records)

	s.mencache.SetYearlyPaymentMethodsCache(year, responses)

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

	if data := s.mencache.GetMonthlyAmountsCache(year); data != nil {
		s.logger.Debug("Successfully fetched monthly amounts from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.transactionStatisticRepository.GetMonthlyAmounts(year)
	if err != nil {
		return s.errorhandler.HandleMonthlyAmountsError(err, "FindMonthlyAmounts", "FAILED_FIND_MONTHLY_AMOUNTS", span, &status, zap.Int("year", year))
	}

	responses := s.mapping.ToTransactionMonthlyAmountResponses(records)

	s.mencache.SetMonthlyAmountsCache(year, responses)

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

	if data := s.mencache.GetYearlyAmountsCache(year); data != nil {
		s.logger.Debug("Successfully fetched yearly amounts from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.transactionStatisticRepository.GetYearlyAmounts(year)
	if err != nil {
		return s.errorhandler.HandleYearlyAmountsError(err, "FindYearlyAmounts", "FAILED_FIND_YEARLY_AMOUNTS", span, &status, zap.Int("year", year))
	}

	responses := s.mapping.ToTransactionYearlyAmountResponses(records)

	s.mencache.SetYearlyAmountsCache(year, responses)

	s.logger.Debug("Successfully fetched yearly amounts", zap.Int("year", year))

	return responses, nil
}

func (s *transactionStatisticService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
