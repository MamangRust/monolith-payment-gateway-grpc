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
	startTime := time.Now()

	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyPaymentMethodsMerchant", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyPaymentMethodsMerchant")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Finding monthly payment methods for merchant", zap.Int("year", year))

	if cachedMerchant := s.mencache.GetMonthlyPaymentMethodsMerchantCache(year); cachedMerchant != nil {
		s.logger.Debug("Successfully fetched merchant from cache",
			zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.merchantStatisRepository.GetMonthlyPaymentMethodsMerchant(year)

	if err != nil {
		return s.errorHandler.HandleMonthlyPaymentMethodsMerchantError(
			err, "FindMonthlyPaymentMethodsMerchant", "FAILED_FIND_MONTHLY_PAYMENT_METHODS_MERCHANT", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapping.ToMerchantMonthlyPaymentMethods(res)

	s.mencache.SetMonthlyPaymentMethodsMerchantCache(year, so)

	s.logger.Debug("Successfully found monthly payment methods for merchant", zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisService) FindYearlyPaymentMethodMerchant(year int) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse) {
	startTime := time.Now()

	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyPaymentMethodMerchant", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyPaymentMethodMerchant")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Finding yearly payment methods for merchant", zap.Int("year", year))

	if cachedMerchant := s.mencache.GetYearlyPaymentMethodMerchantCache(year); cachedMerchant != nil {
		s.logger.Debug("Successfully fetched merchant from cache",
			zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.merchantStatisRepository.GetYearlyPaymentMethodMerchant(year)

	if err != nil {
		return s.errorHandler.HandleYearlyPaymentMethodMerchantError(
			err, "FindYearlyPaymentMethodMerchant", "FAILED_FIND_YEARLY_PAYMENT_METHOD_MERCHANT", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapping.ToMerchantYearlyPaymentMethods(res)

	s.mencache.SetYearlyPaymentMethodMerchantCache(year, so)

	s.logger.Debug("Successfully found yearly payment methods for merchant", zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisService) FindMonthlyAmountMerchant(year int) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse) {
	startTime := time.Now()

	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyAmountMerchant", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyAmountMerchant")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Finding monthly amount for merchant", zap.Int("year", year))

	if cachedMerchant := s.mencache.GetMonthlyAmountMerchantCache(year); cachedMerchant != nil {
		s.logger.Debug("Successfully fetched merchant from cache",
			zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.merchantStatisRepository.GetMonthlyAmountMerchant(year)

	if err != nil {
		return s.errorHandler.HandleMonthlyAmountMerchantError(
			err, "FindMonthlyAmountMerchant", "FAILED_FIND_MONTHLY_AMOUNT_MERCHANT", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapping.ToMerchantMonthlyAmounts(res)

	s.mencache.SetMonthlyAmountMerchantCache(year, so)

	s.logger.Debug("Successfully found monthly amount for merchant", zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisService) FindYearlyAmountMerchant(year int) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse) {
	startTime := time.Now()

	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyAmountMerchant", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyAmountMerchant")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Finding yearly amount for merchant", zap.Int("year", year))

	if cachedMerchant := s.mencache.GetYearlyAmountMerchantCache(year); cachedMerchant != nil {
		s.logger.Debug("Successfully fetched merchant from cache",
			zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.merchantStatisRepository.GetYearlyAmountMerchant(year)

	if err != nil {
		return s.errorHandler.HandleYearlyAmountMerchantError(
			err, "FindYearlyAmountMerchant", "FAILED_FIND_YEARLY_AMOUNT_MERCHANT", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapping.ToMerchantYearlyAmounts(res)

	s.mencache.SetYearlyAmountMerchantCache(year, so)

	s.logger.Debug("Successfully found yearly amount for merchant", zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisService) FindMonthlyTotalAmountMerchant(year int) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyTotalAmountMerchant", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyTotalAmountMerchant")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Finding monthly amount for merchant", zap.Int("year", year))

	if cachedMerchant := s.mencache.GetMonthlyTotalAmountMerchantCache(year); cachedMerchant != nil {
		s.logger.Debug("Successfully fetched merchant from cache",
			zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.merchantStatisRepository.GetMonthlyTotalAmountMerchant(year)

	if err != nil {
		return s.errorHandler.HandleMonthlyTotalAmountMerchantError(
			err, "FindMonthlyTotalAmountMerchant", "FAILED_FIND_MONTHLY_TOTAL_AMOUNT_MERCHANT", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapping.ToMerchantMonthlyTotalAmounts(res)

	s.mencache.SetMonthlyTotalAmountMerchantCache(year, so)

	s.logger.Debug("Successfully found monthly amount for merchant", zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisService) FindYearlyTotalAmountMerchant(year int) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTotalAmountMerchant", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTotalAmountMerchant")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("Finding yearly amount for merchant", zap.Int("year", year))

	if cachedMerchant := s.mencache.GetYearlyTotalAmountMerchantCache(year); cachedMerchant != nil {
		s.logger.Debug("Successfully fetched merchant from cache",
			zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.merchantStatisRepository.GetYearlyTotalAmountMerchant(year)

	if err != nil {
		return s.errorHandler.HandleYearlyTotalAmountMerchantError(
			err, "FindYearlyTotalAmountMerchant", "FAILED_FIND_YEARLY_TOTAL_AMOUNT_MERCHANT", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapping.ToMerchantYearlyTotalAmounts(res)

	s.mencache.SetYearlyTotalAmountMerchantCache(year, so)

	s.logger.Debug("Successfully found yearly amount for merchant", zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
