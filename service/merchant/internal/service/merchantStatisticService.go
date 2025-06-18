package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-merchant/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantStatisService struct {
	ctx                      context.Context
	mencache                 mencache.MerchantStatisticCache
	errorHandler             errorhandler.MerchantStatisticErrorHandler
	trace                    trace.Tracer
	merchantStatisRepository repository.MerchantStatisticRepository
	logger                   logger.LoggerInterface
	mapping                  responseservice.MerchantResponseMapper
	requestCounter           *prometheus.CounterVec
	requestDuration          *prometheus.HistogramVec
}

func NewMerchantStatisService(
	ctx context.Context,
	mencache mencache.MerchantStatisticCache,
	errorHandler errorhandler.MerchantStatisticErrorHandler,
	merchantStatisRepository repository.MerchantStatisticRepository, logger logger.LoggerInterface, mapping responseservice.MerchantResponseMapper) *merchantStatisService {

	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "merchant_statistic_service_request_total",
		Help: "The total number of requests MerchantStatisticService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "merchant_statistic_service_request_duration_seconds",
		Help:    "The duration of requests MerchantStatisticService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	return &merchantStatisService{
		ctx:                      ctx,
		mencache:                 mencache,
		errorHandler:             errorHandler,
		trace:                    otel.Tracer("merchant-statistic-service"),
		merchantStatisRepository: merchantStatisRepository,
		logger:                   logger,
		mapping:                  mapping,
		requestCounter:           requestCounter,
		requestDuration:          requestDuration,
	}
}

func (s *merchantStatisService) FindMonthlyPaymentMethodsMerchant(year int) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse) {
	const method = "FindMonthlyPaymentMethodsMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetMonthlyPaymentMethodsMerchantCache(year); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.merchantStatisRepository.GetMonthlyPaymentMethodsMerchant(year)

	if err != nil {
		return s.errorHandler.HandleMonthlyPaymentMethodsMerchantError(
			err, method, "FAILED_FIND_MONTHLY_PAYMENT_METHODS_MERCHANT", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapping.ToMerchantMonthlyPaymentMethods(res)

	s.mencache.SetMonthlyPaymentMethodsMerchantCache(year, so)

	logSuccess("Successfully found monthly payment methods for merchant", zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisService) FindYearlyPaymentMethodMerchant(year int) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse) {
	const method = "FindYearlyPaymentMethodMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetYearlyPaymentMethodMerchantCache(year); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.merchantStatisRepository.GetYearlyPaymentMethodMerchant(year)

	if err != nil {
		return s.errorHandler.HandleYearlyPaymentMethodMerchantError(
			err, method, "FAILED_FIND_YEARLY_PAYMENT_METHOD_MERCHANT", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapping.ToMerchantYearlyPaymentMethods(res)

	s.mencache.SetYearlyPaymentMethodMerchantCache(year, so)

	logSuccess("Successfully found yearly payment methods for merchant", zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisService) FindMonthlyAmountMerchant(year int) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse) {
	const method = "FindMonthlyAmountMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetMonthlyAmountMerchantCache(year); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.merchantStatisRepository.GetMonthlyAmountMerchant(year)

	if err != nil {
		return s.errorHandler.HandleMonthlyAmountMerchantError(
			err, method, "FAILED_FIND_MONTHLY_AMOUNT_MERCHANT", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapping.ToMerchantMonthlyAmounts(res)

	s.mencache.SetMonthlyAmountMerchantCache(year, so)

	logSuccess("Successfully found monthly amount for merchant", zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisService) FindYearlyAmountMerchant(year int) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse) {
	const method = "FindYearlyAmountMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetYearlyAmountMerchantCache(year); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.merchantStatisRepository.GetYearlyAmountMerchant(year)

	if err != nil {
		return s.errorHandler.HandleYearlyAmountMerchantError(
			err, method, "FAILED_FIND_YEARLY_AMOUNT_MERCHANT", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapping.ToMerchantYearlyAmounts(res)

	s.mencache.SetYearlyAmountMerchantCache(year, so)

	logSuccess("Successfully found yearly amount for merchant", zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisService) FindMonthlyTotalAmountMerchant(year int) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse) {
	const method = "FindMonthlyTotalAmountMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetMonthlyTotalAmountMerchantCache(year); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.merchantStatisRepository.GetMonthlyTotalAmountMerchant(year)

	if err != nil {
		return s.errorHandler.HandleMonthlyTotalAmountMerchantError(
			err, method, "FAILED_FIND_MONTHLY_TOTAL_AMOUNT_MERCHANT", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapping.ToMerchantMonthlyTotalAmounts(res)

	s.mencache.SetMonthlyTotalAmountMerchantCache(year, so)

	logSuccess("Successfully found monthly total amount for merchant", zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisService) FindYearlyTotalAmountMerchant(year int) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse) {
	const method = "FindYearlyTotalAmountMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetYearlyTotalAmountMerchantCache(year); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.merchantStatisRepository.GetYearlyTotalAmountMerchant(year)

	if err != nil {
		return s.errorHandler.HandleYearlyTotalAmountMerchantError(
			err, method, "FAILED_FIND_YEARLY_TOTAL_AMOUNT_MERCHANT", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapping.ToMerchantYearlyTotalAmounts(res)

	s.mencache.SetYearlyTotalAmountMerchantCache(year, so)

	logSuccess("Successfully found yearly total amount for merchant", zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *merchantStatisService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
