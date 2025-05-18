package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	traceunic "github.com/MamangRust/monolith-payment-gateway-pkg/trace_unic"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
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

type merchantStatisByMerchantService struct {
	ctx                                context.Context
	trace                              trace.Tracer
	merchantStatisByMerchantRepository repository.MerchantStatisticByMerchantRepository
	mapping                            responseservice.MerchantResponseMapper
	logger                             logger.LoggerInterface
	requestCounter                     *prometheus.CounterVec
	requestDuration                    *prometheus.HistogramVec
}

func NewMerchantStatisByMerchantService(
	ctx context.Context,
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
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &merchantStatisByMerchantService{
		ctx:                                ctx,
		trace:                              otel.Tracer("merchant-statis-by-merchant-service"),
		merchantStatisByMerchantRepository: merchantStatisBymerchantStatisByMerchantRepository,
		mapping:                            mapping,
		logger:                             logger,
		requestCounter:                     requestCounter,
		requestDuration:                    requestDuration,
	}
}

func (s *merchantStatisByMerchantService) FindMonthlyPaymentMethodByMerchants(req *requests.MonthYearPaymentMethodMerchant) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyPaymentMethodByMerchants", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyPaymentMethodByMerchants")
	defer span.End()

	year := req.Year
	merchantID := req.MerchantID

	span.SetAttributes(
		attribute.Int("merchant_id", merchantID),
		attribute.Int("year", year),
	)

	s.logger.Debug("Finding monthly payment methods by merchant", zap.Int("merchant_id", merchantID), zap.Int("year", year))

	res, err := s.merchantStatisByMerchantRepository.GetMonthlyPaymentMethodByMerchants(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_PAYMENT_METHOD_BY_MERCHANTS")

		span.SetAttributes(attribute.String("trace.id", traceID))

		s.logger.Error("Failed to find monthly payment methods by merchant", zap.Error(err), zap.Int("merchant_id", merchantID), zap.Int("year", year))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find monthly payment methods by merchant")
		status = "failed_find_monthly_payment_method_by_merchants"

		return nil, merchant_errors.ErrFailedFindMonthlyPaymentMethodByMerchants
	}

	so := s.mapping.ToMerchantMonthlyPaymentMethods(res)

	s.logger.Debug("Successfully found monthly payment methods by merchant", zap.Int("merchantID", merchantID), zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisByMerchantService) FindYearlyPaymentMethodByMerchants(req *requests.MonthYearPaymentMethodMerchant) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyPaymentMethodByMerchants", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyPaymentMethodByMerchants")
	defer span.End()

	year := req.Year
	merchantID := req.MerchantID

	span.SetAttributes(
		attribute.Int("merchant_id", merchantID),
		attribute.Int("year", year),
	)

	s.logger.Debug("Finding yearly payment methods by merchant", zap.Int("merchant_id", merchantID), zap.Int("year", year))

	res, err := s.merchantStatisByMerchantRepository.GetYearlyPaymentMethodByMerchants(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_PAYMENT_METHOD_BY_MERCHANTS")

		span.SetAttributes(attribute.String("trace.id", traceID))

		s.logger.Error("Failed to find yearly payment methods by merchant", zap.Error(err), zap.Int("merchant_id", merchantID), zap.Int("year", year))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find yearly payment methods by merchant")
		status = "failed_find_yearly_payment_method_by_merchants"

		return nil, merchant_errors.ErrFailedFindYearlyPaymentMethodByMerchants
	}

	so := s.mapping.ToMerchantYearlyPaymentMethods(res)

	s.logger.Debug("Successfully found yearly payment methods by merchant", zap.Int("merchantID", merchantID), zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisByMerchantService) FindMonthlyAmountByMerchants(req *requests.MonthYearAmountMerchant) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyAmountByMerchants", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyAmountByMerchants")
	defer span.End()

	span.SetAttributes(
		attribute.Int("merchant_id", req.MerchantID),
		attribute.Int("year", req.Year),
	)

	year := req.Year
	merchantID := req.MerchantID

	span.SetAttributes(
		attribute.Int("merchant_id", merchantID),
		attribute.Int("year", year),
	)

	s.logger.Debug("Finding monthly amount by merchant", zap.Int("merchant_id", merchantID), zap.Int("year", year))

	res, err := s.merchantStatisByMerchantRepository.GetMonthlyAmountByMerchants(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_AMOUNT_BY_MERCHANTS")

		span.SetAttributes(attribute.String("trace.id", traceID))

		s.logger.Error("Failed to find monthly amount by merchant", zap.Error(err), zap.Int("merchant_id", merchantID), zap.Int("year", year))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find monthly amount by merchant")
		status = "failed_find_monthly_amount_by_merchants"

		return nil, merchant_errors.ErrFailedFindMonthlyAmountByMerchants
	}

	so := s.mapping.ToMerchantMonthlyAmounts(res)

	s.logger.Debug("Successfully found monthly amount by merchant", zap.Int("merchantID", merchantID), zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisByMerchantService) FindYearlyAmountByMerchants(req *requests.MonthYearAmountMerchant) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyAmountByMerchants", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyAmountByMerchants")
	defer span.End()

	year := req.Year
	merchantID := req.MerchantID

	span.SetAttributes(
		attribute.Int("merchant_id", merchantID),
		attribute.Int("year", year),
	)

	s.logger.Debug("Finding yearly amount by merchant", zap.Int("merchantID", merchantID), zap.Int("year", year))

	res, err := s.merchantStatisByMerchantRepository.GetYearlyAmountByMerchants(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_AMOUNT_BY_MERCHANTS")

		span.SetAttributes(attribute.String("trace.id", traceID))

		s.logger.Error("Failed to find yearly amount by merchant", zap.Error(err), zap.Int("merchantID", merchantID), zap.Int("year", year))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find yearly amount by merchant")
		status = "failed_find_yearly_amount_by_merchants"

		return nil, merchant_errors.ErrFailedFindYearlyAmountByMerchants
	}

	so := s.mapping.ToMerchantYearlyAmounts(res)

	s.logger.Debug("Successfully found yearly amount by merchant", zap.Int("merchantID", merchantID), zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisByMerchantService) FindMonthlyTotalAmountByMerchants(req *requests.MonthYearTotalAmountMerchant) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyTotalAmountByMerchants", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyTotalAmountByMerchants")
	defer span.End()

	year := req.Year
	merchantID := req.MerchantID

	span.SetAttributes(
		attribute.Int("merchant_id", merchantID),
		attribute.Int("year", year),
	)

	s.logger.Debug("Finding monthly total amount by merchant", zap.Int("merchant_id", merchantID), zap.Int("year", year))

	res, err := s.merchantStatisByMerchantRepository.GetMonthlyTotalAmountByMerchants(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_TOTAL_AMOUNT_BY_MERCHANTS")

		span.SetAttributes(attribute.String("trace.id", traceID))

		s.logger.Error("Failed to find monthly total amount by merchant", zap.Error(err), zap.Int("merchant_id", merchantID), zap.Int("year", year))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find monthly total amount by merchant")
		status = "failed_find_monthly_total_amount_by_merchants"
		return nil, merchant_errors.ErrFailedFindMonthlyTotalAmountByMerchants
	}

	s.logger.Debug("Example", zap.Any("response month", res))

	so := s.mapping.ToMerchantMonthlyTotalAmounts(res)

	s.logger.Debug("Successfully found monthly total amount by merchant", zap.Int("merchantID", merchantID), zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisByMerchantService) FindYearlyTotalAmountByMerchants(req *requests.MonthYearTotalAmountMerchant) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTotalAmountByMerchants", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTotalAmountByMerchants")
	defer span.End()

	year := req.Year
	merchantID := req.MerchantID

	span.SetAttributes(
		attribute.Int("merchant_id", merchantID),
		attribute.Int("year", year),
	)

	s.logger.Debug("Finding yearly amount by merchant", zap.Int("merchantID", merchantID), zap.Int("year", year))

	res, err := s.merchantStatisByMerchantRepository.GetYearlyTotalAmountByMerchants(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_TOTAL_AMOUNT_BY_MERCHANTS")

		span.SetAttributes(attribute.String("trace.id", traceID))

		s.logger.Error("Failed to find yearly amount by merchant", zap.Error(err), zap.Int("merchantID", merchantID), zap.Int("year", year))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find yearly amount by merchant")
		status = "failed_find_yearly_total_amount_by_merchants"
		return nil, merchant_errors.ErrFailedFindYearlyTotalAmountByMerchants
	}

	so := s.mapping.ToMerchantYearlyTotalAmounts(res)

	s.logger.Debug("Successfully found yearly amount by merchant", zap.Int("merchantID", merchantID), zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisByMerchantService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
