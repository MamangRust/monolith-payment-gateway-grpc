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
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantStatisByApiKeyService struct {
	ctx                              context.Context
	trace                            trace.Tracer
	mencache                         mencache.MerchantStatisticByApikeyCache
	errorHandler                     errorhandler.MerchantStatisticByApikeyErrorHandler
	merchantStatisByApiKeyRepository repository.MerchantStatisticByApiKeyRepository
	logger                           logger.LoggerInterface
	mapping                          responseservice.MerchantResponseMapper
	requestCounter                   *prometheus.CounterVec
	requestDuration                  *prometheus.HistogramVec
}

func NewMerchantStatisByApiKeyService(
	ctx context.Context,
	mencache mencache.MerchantStatisticByApikeyCache,
	errorHandler errorhandler.MerchantStatisticByApikeyErrorHandler,
	merchantStatisByApiKeyRepository repository.MerchantStatisticByApiKeyRepository,
	logger logger.LoggerInterface,
	mapping responseservice.MerchantResponseMapper,
) *merchantStatisByApiKeyService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_statis_by_apikey_service_requests_total",
			Help: "Total number of requests to the MerchantStatisByApiKeyService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_statis_by_apikey_service_request_duration_seconds",
			Help:    "Histogram of request durations for the MerchantStatisByApiKeyService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &merchantStatisByApiKeyService{
		ctx:                              ctx,
		mencache:                         mencache,
		errorHandler:                     errorHandler,
		trace:                            otel.Tracer("merchant-statis-by-apikey-service"),
		merchantStatisByApiKeyRepository: merchantStatisByApiKeyRepository,
		logger:                           logger,
		mapping:                          mapping,
		requestCounter:                   requestCounter,
		requestDuration:                  requestDuration,
	}
}

func (s *merchantStatisByApiKeyService) FindMonthlyPaymentMethodByApikeys(req *requests.MonthYearPaymentMethodApiKey) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse) {
	startTime := time.Now()

	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyPaymentMethodByApikeys", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyPaymentMethodByApikeys")
	defer span.End()

	api_key := req.Apikey
	year := req.Year

	s.logger.Debug("Finding monthly payment methods by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	if cachedMerchant := s.mencache.GetMonthlyPaymentMethodByApikeysCache(req); cachedMerchant != nil {
		s.logger.Debug("Successfully fetched merchant from cache",
			zap.String("api_key", api_key), zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.merchantStatisByApiKeyRepository.GetMonthlyPaymentMethodByApikey(req)

	if err != nil {
		return s.errorHandler.HandleMonthlyPaymentMethodByApikeysError(
			err, "FindMonthlyPaymentMethodByApikeys", "FAILED_FIND_MONTHLY_PAYMENT_METHOD_BY_APIKEYS", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapping.ToMerchantMonthlyPaymentMethods(res)

	s.mencache.SetMonthlyPaymentMethodByApikeysCache(req, so)

	s.logger.Debug("Successfully found monthly payment methods by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisByApiKeyService) FindYearlyPaymentMethodByApikeys(req *requests.MonthYearPaymentMethodApiKey) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyPaymentMethodByApikeys", status, startTime)
	}()
	_, span := s.trace.Start(s.ctx, "FindYearlyPaymentMethodByApikeys")
	defer span.End()

	api_key := req.Apikey
	year := req.Year

	span.SetAttributes(
		attribute.String("api_key", api_key),
		attribute.Int("year", year),
	)

	s.logger.Debug("Finding yearly payment methods by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	if cachedMerchant := s.mencache.GetYearlyPaymentMethodByApikeysCache(req); cachedMerchant != nil {
		s.logger.Debug("Successfully fetched merchant from cache",
			zap.String("api_key", api_key), zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.merchantStatisByApiKeyRepository.GetYearlyPaymentMethodByApikey(req)

	if err != nil {
		return s.errorHandler.HandleYearlyPaymentMethodByApikeysError(
			err, "FindYearlyPaymentMethodByApikeys", "FAILED_FIND_YEARLY_PAYMENT_METHOD_BY_APIKEYS", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapping.ToMerchantYearlyPaymentMethods(res)

	s.mencache.SetYearlyPaymentMethodByApikeysCache(req, so)

	s.logger.Debug("Successfully found yearly payment methods by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisByApiKeyService) FindMonthlyAmountByApikeys(req *requests.MonthYearAmountApiKey) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyAmountByApikeys", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyAmountByApikeys")
	defer span.End()

	api_key := req.Apikey
	year := req.Year

	span.SetAttributes(
		attribute.String("api_key", api_key),
		attribute.Int("year", year),
	)

	s.logger.Debug("Finding monthly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	if cachedMerchant := s.mencache.GetMonthlyAmountByApikeysCache(req); cachedMerchant != nil {
		s.logger.Debug("Successfully fetched merchant from cache",
			zap.String("api_key", api_key), zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.merchantStatisByApiKeyRepository.GetMonthlyAmountByApikey(req)

	if err != nil {
		return s.errorHandler.HandleMonthlyAmountByApikeysError(
			err, "FindMonthlyAmountByApikeys", "FAILED_FIND_MONTHLY_AMOUNT_BY_APIKEYS", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapping.ToMerchantMonthlyAmounts(res)

	s.mencache.SetMonthlyAmountByApikeysCache(req, so)

	s.logger.Debug("Successfully found monthly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisByApiKeyService) FindYearlyAmountByApikeys(req *requests.MonthYearAmountApiKey) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyAmountByApikeys", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyAmountByApikeys")
	defer span.End()

	api_key := req.Apikey
	year := req.Year

	span.SetAttributes(
		attribute.String("api_key", api_key),
		attribute.Int("year", year),
	)

	s.logger.Debug("Finding yearly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	if cachedMerchant := s.mencache.GetYearlyAmountByApikeysCache(req); cachedMerchant != nil {
		s.logger.Debug("Successfully fetched merchant from cache",
			zap.String("api_key", api_key), zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.merchantStatisByApiKeyRepository.GetYearlyAmountByApikey(req)

	if err != nil {
		return s.errorHandler.HandleYearlyAmountByApikeysError(
			err, "FindYearlyAmountByApikeys", "FAILED_FIND_YEARLY_AMOUNT_BY_APIKEYS", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapping.ToMerchantYearlyAmounts(res)

	s.mencache.SetYearlyAmountByApikeysCache(req, so)

	s.logger.Debug("Successfully found yearly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisByApiKeyService) FindMonthlyTotalAmountByApikeys(req *requests.MonthYearTotalAmountApiKey) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyTotalAmountByApikeys", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyTotalAmountByApikeys")
	defer span.End()

	api_key := req.Apikey
	year := req.Year

	span.SetAttributes(
		attribute.String("api_key", api_key),
		attribute.Int("year", year),
	)

	s.logger.Debug("Finding monthly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	if cachedMerchant := s.mencache.GetMonthlyTotalAmountByApikeysCache(req); cachedMerchant != nil {
		s.logger.Debug("Successfully fetched merchant from cache",
			zap.String("api_key", api_key), zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.merchantStatisByApiKeyRepository.GetMonthlyTotalAmountByApikey(req)

	if err != nil {
		return s.errorHandler.HandleMonthlyTotalAmountByApikeysError(
			err, "FindMonthlyTotalAmountByApikeys", "FAILED_FIND_MONTHLY_TOTAL_AMOUNT_BY_APIKEYS", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapping.ToMerchantMonthlyTotalAmounts(res)

	s.mencache.SetMonthlyTotalAmountByApikeysCache(req, so)

	s.logger.Debug("Successfully found monthly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisByApiKeyService) FindYearlyTotalAmountByApikeys(req *requests.MonthYearTotalAmountApiKey) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTotalAmountByApikeys", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTotalAmountByApikeys")
	defer span.End()

	api_key := req.Apikey
	year := req.Year

	span.SetAttributes(
		attribute.String("api_key", api_key),
		attribute.Int("year", year),
	)

	s.logger.Debug("Finding yearly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	if cachedMerchant := s.mencache.GetYearlyTotalAmountByApikeysCache(req); cachedMerchant != nil {
		s.logger.Debug("Successfully fetched merchant from cache",
			zap.String("api_key", api_key), zap.Int("year", year))
		return cachedMerchant, nil
	}

	res, err := s.merchantStatisByApiKeyRepository.GetYearlyTotalAmountByApikey(req)

	if err != nil {
		return s.errorHandler.HandleYearlyTotalAmountByApikeysError(
			err, "FindYearlyTotalAmountByApikeys", "FAILED_FIND_YEARLY_TOTAL_AMOUNT_BY_APIKEYS", span, &status,
			zap.Error(err),
		)
	}

	so := s.mapping.ToMerchantYearlyTotalAmounts(res)

	s.mencache.SetYearlyTotalAmountByApikeysCache(req, so)

	s.logger.Debug("Successfully found yearly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisByApiKeyService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
