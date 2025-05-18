package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	traceunic "github.com/MamangRust/monolith-payment-gateway-pkg/trace_unic"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors"
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
	trace                    trace.Tracer
	merchantStatisRepository repository.MerchantStatisticRepository
	logger                   logger.LoggerInterface
	mapping                  responseservice.MerchantResponseMapper
	requestCounter           *prometheus.CounterVec
	requestDuration          *prometheus.HistogramVec
}

func NewMerchantStatisService(
	ctx context.Context,
	merchantStatisRepository repository.MerchantStatisticRepository, logger logger.LoggerInterface, mapping responseservice.MerchantResponseMapper) *merchantStatisService {

	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "merchant_statistic_service_request_total",
		Help: "The total number of requests MerchantStatisticService",
	}, []string{"method", "path"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "merchant_statistic_service_request_duration_seconds",
		Help:    "The duration of requests MerchantStatisticService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})

	prometheus.MustRegister(requestCounter, requestDuration)

	return &merchantStatisService{
		ctx:                      ctx,
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

	res, err := s.merchantStatisRepository.GetMonthlyPaymentMethodsMerchant(year)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_PAYMENT_METHODS_MERCHANT")

		span.SetAttributes(attribute.String("trace.id", traceID))

		s.logger.Error("Failed to find monthly payment methods for merchant", zap.Error(err), zap.Int("year", year))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find monthly payment methods for merchant")
		status = "failed_find_monthly_payment_methods_merchant"

		return nil, merchant_errors.ErrFailedFindMonthlyPaymentMethodsMerchant
	}

	so := s.mapping.ToMerchantMonthlyPaymentMethods(res)

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

	res, err := s.merchantStatisRepository.GetYearlyPaymentMethodMerchant(year)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_PAYMENT_METHOD_MERCHANT")

		span.SetAttributes(attribute.String("trace.id", traceID))

		s.logger.Error("Failed to find yearly payment methods for merchant", zap.Error(err), zap.Int("year", year))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find yearly payment methods for merchant")
		status = "failed_find_yearly_payment_method_merchant"
		return nil, merchant_errors.ErrFailedFindYearlyPaymentMethodMerchant
	}

	so := s.mapping.ToMerchantYearlyPaymentMethods(res)

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

	res, err := s.merchantStatisRepository.GetMonthlyAmountMerchant(year)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_AMOUNT_MERCHANT")

		span.SetAttributes(attribute.String("trace.id", traceID))

		s.logger.Error("Failed to find monthly amount for merchant", zap.Error(err), zap.Int("year", year))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find monthly amount for merchant")
		status = "failed_find_monthly_amount_merchant"
		return nil, merchant_errors.ErrFailedFindMonthlyAmountMerchant
	}

	so := s.mapping.ToMerchantMonthlyAmounts(res)

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

	res, err := s.merchantStatisRepository.GetYearlyAmountMerchant(year)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_AMOUNT_MERCHANT")

		span.SetAttributes(attribute.String("trace.id", traceID))

		s.logger.Error("Failed to find yearly amount for merchant", zap.Error(err), zap.Int("year", year))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find yearly amount for merchant")
		status = "failed_find_yearly_amount_merchant"
		return nil, merchant_errors.ErrFailedFindYearlyAmountMerchant
	}

	so := s.mapping.ToMerchantYearlyAmounts(res)

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

	res, err := s.merchantStatisRepository.GetMonthlyTotalAmountMerchant(year)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_TOTAL_AMOUNT_MERCHANT")

		span.SetAttributes(attribute.String("trace.id", traceID))

		s.logger.Error("Failed to find monthly amount for merchant", zap.Error(err), zap.Int("year", year))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find monthly amount for merchant")
		status = "failed_find_monthly_total_amount_merchant"
		return nil, merchant_errors.ErrFailedFindMonthlyTotalAmountMerchant
	}

	so := s.mapping.ToMerchantMonthlyTotalAmounts(res)

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

	res, err := s.merchantStatisRepository.GetYearlyTotalAmountMerchant(year)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_TOTAL_AMOUNT_MERCHANT")

		span.SetAttributes(attribute.String("trace.id", traceID))

		s.logger.Error("Failed to find yearly amount for merchant", zap.Error(err), zap.Int("year", year))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find yearly amount for merchant")
		status = "failed_find_yearly_total_amount_merchant"
		return nil, merchant_errors.ErrFailedFindYearlyTotalAmountMerchant
	}

	so := s.mapping.ToMerchantYearlyTotalAmounts(res)

	s.logger.Debug("Successfully found yearly amount for merchant", zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
