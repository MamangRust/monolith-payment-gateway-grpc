package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-merchant/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantStatisByMerchantService struct {
	ctx                                context.Context
	trace                              trace.Tracer
	mencache                           mencache.MerchantStatisticByMerchantCache
	errorHandler                       errorhandler.MerchantStatisticByMerchantErrorHandler
	merchantStatisByMerchantRepository repository.MerchantStatisticByMerchantRepository
	mapping                            responseservice.MerchantResponseMapper
	logger                             logger.LoggerInterface
	requestCounter                     *prometheus.CounterVec
	requestDuration                    *prometheus.HistogramVec
}

func NewMerchantStatisByMerchantService(
	ctx context.Context,
	mencache mencache.MerchantStatisticByMerchantCache,
	errorHandler errorhandler.MerchantStatisticByMerchantErrorHandler,
	merchantStatisBymerchantStatisByMerchantRepository repository.MerchantStatisticByMerchantRepository,
	logger logger.LoggerInterface,
	mapping responseservice.MerchantResponseMapper,
) *merchantStatisByMerchantService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_statis_by_merchant_service_requests_total",
			Help: "Total number of requests to the MerchantStatisByMerchantService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_statis_by_merchant_service_request_duration_seconds",
			Help:    "Histogram of request durations for the MerchantStatisByMerchantService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &merchantStatisByMerchantService{
		ctx:                                ctx,
		mencache:                           mencache,
		errorHandler:                       errorHandler,
		trace:                              otel.Tracer("merchant-statis-by-merchant-service"),
		merchantStatisByMerchantRepository: merchantStatisBymerchantStatisByMerchantRepository,
		mapping:                            mapping,
		logger:                             logger,
		requestCounter:                     requestCounter,
		requestDuration:                    requestDuration,
	}
}

func (s *merchantStatisByMerchantService) FindMonthlyPaymentMethodByMerchants(req *requests.MonthYearPaymentMethodMerchant) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse) {
	year := req.Year
	merchantID := req.MerchantID

	const method = "FindMonthlyPaymentMethodByMerchants"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("merchant.id", merchantID), attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetMonthlyPaymentMethodByMerchantsCache(req); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("merchant_id", merchantID))
		return cachedMerchant, nil
	}

	res, err := s.merchantStatisByMerchantRepository.GetMonthlyPaymentMethodByMerchants(req)

	if err != nil {
		return s.errorHandler.HandleMonthlyPaymentMethodByMerchantsError(
			err, method, "FAILED_FIND_MONTHLY_PAYMENT_METHOD_BY_MERCHANTS", span, &status,
			zap.Any("error", err),
		)
	}

	so := s.mapping.ToMerchantMonthlyPaymentMethods(res)

	s.mencache.SetMonthlyPaymentMethodByMerchantsCache(req, so)

	logSuccess("Successfully fetched merchant", zap.Int("merchant.id", merchantID))

	return so, nil
}

func (s *merchantStatisByMerchantService) FindYearlyPaymentMethodByMerchants(req *requests.MonthYearPaymentMethodMerchant) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse) {
	year := req.Year
	merchantID := req.MerchantID

	const method = "FindYearlyPaymentMethodByMerchants"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("merchant.id", merchantID), attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetYearlyPaymentMethodByMerchantsCache(req); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("merchant_id", merchantID))
		return cachedMerchant, nil
	}

	res, err := s.merchantStatisByMerchantRepository.GetYearlyPaymentMethodByMerchants(req)

	if err != nil {
		return s.errorHandler.HandleYearlyPaymentMethodByMerchantsError(
			err, method, "FAILED_FIND_YEARLY_PAYMENT_METHOD_BY_MERCHANTS", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapping.ToMerchantYearlyPaymentMethods(res)

	s.mencache.SetYearlyPaymentMethodByMerchantsCache(req, so)

	logSuccess("Successfully fetched merchant", zap.Int("merchant.id", merchantID))

	return so, nil
}

func (s *merchantStatisByMerchantService) FindMonthlyAmountByMerchants(req *requests.MonthYearAmountMerchant) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse) {
	year := req.Year
	merchantID := req.MerchantID

	const method = "FindMonthlyAmountByMerchants"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("merchant.id", merchantID), attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetMonthlyAmountByMerchantsCache(req); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("merchant.id", merchantID))

		return cachedMerchant, nil
	}

	res, err := s.merchantStatisByMerchantRepository.GetMonthlyAmountByMerchants(req)
	if err != nil {
		return s.errorHandler.HandleMonthlyAmountByMerchantsError(
			err, method, "FAILED_FIND_MONTHLY_AMOUNT_BY_MERCHANTS", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapping.ToMerchantMonthlyAmounts(res)

	s.mencache.SetMonthlyAmountByMerchantsCache(req, so)

	logSuccess("Successfully fetched merchant", zap.Int("merchant.id", merchantID))

	return so, nil
}

func (s *merchantStatisByMerchantService) FindYearlyAmountByMerchants(req *requests.MonthYearAmountMerchant) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse) {
	year := req.Year
	merchantID := req.MerchantID

	const method = "FindYearlyAmountByMerchants"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("merchant.id", merchantID), attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetYearlyAmountByMerchantsCache(req); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("merchant.id", merchantID))

		return cachedMerchant, nil
	}

	res, err := s.merchantStatisByMerchantRepository.GetYearlyAmountByMerchants(req)

	if err != nil {
		return s.errorHandler.HandleYearlyAmountByMerchantsError(
			err, method, "FAILED_FIND_YEARLY_AMOUNT_BY_MERCHANTS", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapping.ToMerchantYearlyAmounts(res)

	s.mencache.SetYearlyAmountByMerchantsCache(req, so)

	logSuccess("Successfully fetched merchant", zap.Int("merchant.id", merchantID))

	return so, nil
}

func (s *merchantStatisByMerchantService) FindMonthlyTotalAmountByMerchants(req *requests.MonthYearTotalAmountMerchant) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse) {
	year := req.Year
	merchantID := req.MerchantID

	const method = "FindMonthlyTotalAmountByMerchants"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("merchant.id", merchantID), attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetMonthlyTotalAmountByMerchantsCache(req); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("merchant.id", merchantID), zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.merchantStatisByMerchantRepository.GetMonthlyTotalAmountByMerchants(req)

	if err != nil {
		return s.errorHandler.HandleMonthlyTotalAmountByMerchantsError(
			err, method, "FAILED_FIND_MONTHLY_TOTAL_AMOUNT_BY_MERCHANTS", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapping.ToMerchantMonthlyTotalAmounts(res)

	s.mencache.SetMonthlyTotalAmountByMerchantsCache(req, so)

	logSuccess("Successfully fetched merchant", zap.Int("merchant.id", merchantID), zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisByMerchantService) FindYearlyTotalAmountByMerchants(req *requests.MonthYearTotalAmountMerchant) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse) {
	year := req.Year
	merchantID := req.MerchantID

	const method = "FindYearlyTotalAmountByMerchants"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("merchant.id", merchantID), attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetYearlyTotalAmountByMerchantsCache(req); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("merchant.id", merchantID), zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.merchantStatisByMerchantRepository.GetYearlyTotalAmountByMerchants(req)
	if err != nil {
		return s.errorHandler.HandleYearlyTotalAmountByMerchantsError(
			err, method, "FAILED_FIND_YEARLY_TOTAL_AMOUNT_BY_MERCHANTS", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapping.ToMerchantYearlyTotalAmounts(res)

	s.mencache.SetYearlyTotalAmountByMerchantsCache(req, so)

	logSuccess("Successfully fetched merchant", zap.Int("merchant.id", merchantID), zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisByMerchantService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *merchantStatisByMerchantService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
