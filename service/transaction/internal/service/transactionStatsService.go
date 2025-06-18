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
	"go.opentelemetry.io/otel/codes"
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
	month := req.Month
	year := req.Year

	const method = "FindMonthTransactionStatusSuccess"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("month", month))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthTransactonStatusSuccessCache(req); found {
		logSuccess("Successfully fetched monthly Transaction status success from cache", zap.Int("year", year), zap.Int("month", month))
		return data, nil
	}

	records, err := s.transactionStatisticRepository.GetMonthTransactionStatusSuccess(req)
	if err != nil {
		return s.errorhandler.HandleMonthTransactionStatusSuccessError(err, method, "FAILED_FIND_MONTH_TRANSACTION_STATUS_SUCCESS", span, &status, zap.Error(err))
	}
	so := s.mapping.ToTransactionResponsesMonthStatusSuccess(records)

	s.mencache.SetMonthTransactonStatusSuccessCache(req, so)

	logSuccess("Successfully fetched monthly Transaction status success", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

func (s *transactionStatisticService) FindYearlyTransactionStatusSuccess(year int) ([]*response.TransactionResponseYearStatusSuccess, *response.ErrorResponse) {
	const method = "FindYearlyTransactionStatusSuccess"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearTransactonStatusSuccessCache(year); found {
		logSuccess("Successfully fetched yearly Transaction status success from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.transactionStatisticRepository.GetYearlyTransactionStatusSuccess(year)
	if err != nil {
		return s.errorhandler.HandleYearlyTransactionStatusSuccessError(err, method, "FAILED_FIND_YEARLY_TRANSACTION_STATUS_SUCCESS", span, &status, zap.Error(err))
	}
	so := s.mapping.ToTransactionResponsesYearStatusSuccess(records)

	s.mencache.SetYearTransactonStatusSuccessCache(year, so)

	logSuccess("Successfully fetched yearly Transaction status success", zap.Int("year", year))

	return so, nil
}

func (s *transactionStatisticService) FindMonthTransactionStatusFailed(req *requests.MonthStatusTransaction) ([]*response.TransactionResponseMonthStatusFailed, *response.ErrorResponse) {

	year := req.Year
	month := req.Month

	const method = "FindMonthTransactionStatusFailed"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("month", month))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthTransactonStatusFailedCache(req); found {
		logSuccess("Successfully fetched monthly Transaction status Failed from cache", zap.Int("year", year), zap.Int("month", month))
		return data, nil
	}

	records, err := s.transactionStatisticRepository.GetMonthTransactionStatusFailed(req)
	if err != nil {
		return s.errorhandler.HandleMonthTransactionStatusFailedError(err, method, "FAILED_MONTHLY_TRANSACTION_STATUS_FAILED", span, &status, zap.Error(err))
	}
	so := s.mapping.ToTransactionResponsesMonthStatusFailed(records)

	s.mencache.SetMonthTransactonStatusFailedCache(req, so)

	logSuccess("Successfully fetched monthly Transaction status Failed", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

func (s *transactionStatisticService) FindYearlyTransactionStatusFailed(year int) ([]*response.TransactionResponseYearStatusFailed, *response.ErrorResponse) {
	const method = "FindYearlyTransactionStatusFailed"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearTransactonStatusFailedCache(year); found {
		logSuccess("Successfully fetched yearly Transaction status Failed from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.transactionStatisticRepository.GetYearlyTransactionStatusFailed(year)
	if err != nil {
		return s.errorhandler.HandleYearlyTransactionStatusFailedError(err, method, "FAILED_FIND_YEARLY_TRANSACTION_STATUS_FAILED", span, &status, zap.Error(err))
	}
	so := s.mapping.ToTransactionResponsesYearStatusFailed(records)

	s.mencache.SetYearTransactonStatusFailedCache(year, so)

	logSuccess("Successfully fetched yearly Transaction status Failed", zap.Int("year", year))

	return so, nil
}

func (s *transactionStatisticService) FindMonthlyPaymentMethods(year int) ([]*response.TransactionMonthMethodResponse, *response.ErrorResponse) {
	const method = "FindMonthlyPaymentMethods"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyPaymentMethodsCache(year); found {
		logSuccess("Successfully fetched monthly payment methods from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.transactionStatisticRepository.GetMonthlyPaymentMethods(year)
	if err != nil {
		return s.errorhandler.HandleMonthlyPaymentMethodsError(err, method, "FAILED_FIND_MONTHLY_PAYMENT_METHODS", span, &status, zap.Error(err))
	}

	responses := s.mapping.ToTransactionMonthlyMethodResponses(records)

	s.mencache.SetMonthlyPaymentMethodsCache(year, responses)

	logSuccess("Successfully fetched monthly payment methods", zap.Int("year", year))

	return responses, nil
}

func (s *transactionStatisticService) FindYearlyPaymentMethods(year int) ([]*response.TransactionYearMethodResponse, *response.ErrorResponse) {
	const method = "FindYearlyPaymentMethods"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyPaymentMethodsCache(year); found {
		logSuccess("Successfully fetched yearly payment methods from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.transactionStatisticRepository.GetYearlyPaymentMethods(year)
	if err != nil {
		return s.errorhandler.HandleYearlyPaymentMethodsError(err, method, "FAILED_FIND_YEARLY_PAYMENT_METHODS", span, &status, zap.Error(err))
	}

	responses := s.mapping.ToTransactionYearlyMethodResponses(records)

	s.mencache.SetYearlyPaymentMethodsCache(year, responses)

	logSuccess("Successfully fetched yearly payment methods", zap.Int("year", year))

	return responses, nil
}

func (s *transactionStatisticService) FindMonthlyAmounts(year int) ([]*response.TransactionMonthAmountResponse, *response.ErrorResponse) {
	const method = "FindMonthlyAmounts"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()
	if data, found := s.mencache.GetMonthlyAmountsCache(year); found {
		s.logger.Debug("Successfully fetched monthly amounts from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.transactionStatisticRepository.GetMonthlyAmounts(year)
	if err != nil {
		return s.errorhandler.HandleMonthlyAmountsError(err, method, "FAILED_FIND_MONTHLY_AMOUNTS", span, &status, zap.Error(err))
	}

	responses := s.mapping.ToTransactionMonthlyAmountResponses(records)

	s.mencache.SetMonthlyAmountsCache(year, responses)

	logSuccess("Successfully fetched monthly amounts", zap.Int("year", year))

	return responses, nil
}

func (s *transactionStatisticService) FindYearlyAmounts(year int) ([]*response.TransactionYearlyAmountResponse, *response.ErrorResponse) {
	const method = "FindYearlyAmounts"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyAmountsCache(year); found {
		logSuccess("Successfully fetched yearly amounts from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.transactionStatisticRepository.GetYearlyAmounts(year)
	if err != nil {
		return s.errorhandler.HandleYearlyAmountsError(err, method, "FAILED_FIND_YEARLY_AMOUNTS", span, &status, zap.Error(err))
	}

	responses := s.mapping.ToTransactionYearlyAmountResponses(records)

	s.mencache.SetYearlyAmountsCache(year, responses)

	logSuccess("Successfully fetched yearly amounts", zap.Int("year", year))

	return responses, nil
}

func (s *transactionStatisticService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *transactionStatisticService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
