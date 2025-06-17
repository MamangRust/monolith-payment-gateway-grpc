package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-transfer/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transferStatisticService struct {
	ctx                         context.Context
	errorhandler                errorhandler.TransferStatisticErrorHandler
	mencache                    mencache.TransferStatisticCache
	trace                       trace.Tracer
	transferStatisticRepository repository.TransferStatisticRepository
	logger                      logger.LoggerInterface
	mapping                     responseservice.TransferResponseMapper
	requestCounter              *prometheus.CounterVec
	requestDuration             *prometheus.HistogramVec
}

func NewTransferStatisticService(ctx context.Context, errorhandler errorhandler.TransferStatisticErrorHandler,
	mencache mencache.TransferStatisticCache, transferStatisticRepository repository.TransferStatisticRepository, logger logger.LoggerInterface, mapping responseservice.TransferResponseMapper) *transferStatisticService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transfer_statistic_service_request_total",
			Help: "Total number of requests to the TransferStatisticService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transfer_statistic_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TransferStatisticService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &transferStatisticService{
		ctx:                         ctx,
		errorhandler:                errorhandler,
		mencache:                    mencache,
		trace:                       otel.Tracer("transfer-statistic-service"),
		transferStatisticRepository: transferStatisticRepository,
		logger:                      logger,
		mapping:                     mapping,
		requestCounter:              requestCounter,
		requestDuration:             requestDuration,
	}
}

func (s *transferStatisticService) FindMonthTransferStatusSuccess(req *requests.MonthStatusTransfer) ([]*response.TransferResponseMonthStatusSuccess, *response.ErrorResponse) {

	year := req.Year
	month := req.Month

	const method = "FindMonthTransferStatusSuccess"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("month", month))

	defer func() {
		end(status)
	}()

	if data := s.mencache.GetCachedMonthTransferStatusSuccess(req); data != nil {
		logSuccess("Successfully fetched monthly Transfer status success from cache", zap.Int("year", year), zap.Int("month", month))
		return data, nil
	}

	records, err := s.transferStatisticRepository.GetMonthTransferStatusSuccess(req)

	if err != nil {
		return s.errorhandler.HandleMonthTransferStatusSuccessError(err, method, "FAILED_FIND_MONTH_TRANSFER_STATUS_SUCCESS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransferResponsesMonthStatusSuccess(records)

	s.mencache.SetCachedMonthTransferStatusSuccess(req, so)

	logSuccess("Successfully fetched monthly Transfer status success", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

func (s *transferStatisticService) FindYearlyTransferStatusSuccess(year int) ([]*response.TransferResponseYearStatusSuccess, *response.ErrorResponse) {
	const method = "FindYearlyTransferStatusSuccess"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data := s.mencache.GetCachedYearlyTransferStatusSuccess(year); data != nil {
		logSuccess("Successfully fetched yearly Transfer status success from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.transferStatisticRepository.GetYearlyTransferStatusSuccess(year)

	if err != nil {
		return s.errorhandler.HandleYearTransferStatusSuccessError(err, method, "FAILED_YEARLY_TRANSFER_STATUS_SUCCESS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransferResponsesYearStatusSuccess(records)

	s.mencache.SetCachedYearlyTransferStatusSuccess(year, so)

	logSuccess("Successfully fetched yearly Transfer status success", zap.Int("year", year))

	return so, nil

}

func (s *transferStatisticService) FindMonthTransferStatusFailed(req *requests.MonthStatusTransfer) ([]*response.TransferResponseMonthStatusFailed, *response.ErrorResponse) {
	year := req.Year
	month := req.Month

	const method = "FindMonthTransferStatusFailed"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("month", month))

	defer func() {
		end(status)
	}()

	if data := s.mencache.GetCachedMonthTransferStatusFailed(req); data != nil {
		logSuccess("Successfully fetched monthly Transfer status Failed from cache", zap.Int("year", year), zap.Int("month", month))
		return data, nil
	}

	records, err := s.transferStatisticRepository.GetMonthTransferStatusFailed(req)

	if err != nil {
		return s.errorhandler.HandleMonthTransferStatusFailedError(err, method, "FAILED_MONTHLY_TRANSFER_STATUS_FAILED", span, &status, zap.Error(err))
	}
	so := s.mapping.ToTransferResponsesMonthStatusFailed(records)

	s.mencache.SetCachedMonthTransferStatusFailed(req, so)

	logSuccess("Successfully fetched monthly Transfer status Failed", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

func (s *transferStatisticService) FindYearlyTransferStatusFailed(year int) ([]*response.TransferResponseYearStatusFailed, *response.ErrorResponse) {
	const method = "FindYearlyTransferStatusFailed"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data := s.mencache.GetCachedYearlyTransferStatusFailed(year); data != nil {
		logSuccess("Successfully fetched yearly Transfer status Failed from cache", zap.Int("year", year))
		return data, nil
	}

	records, err := s.transferStatisticRepository.GetYearlyTransferStatusFailed(year)
	if err != nil {
		return s.errorhandler.HandleYearTransferStatusFailedError(err, method, "FAILED_FIND_YEARLY_TRANSFER_STATUS_FAILED", span, &status, zap.Error(err))
	}
	so := s.mapping.ToTransferResponsesYearStatusFailed(records)

	s.mencache.SetCachedYearlyTransferStatusFailed(year, so)

	logSuccess("Successfully fetched yearly Transfer status Failed", zap.Int("year", year))

	return so, nil
}

func (s *transferStatisticService) FindMonthlyTransferAmounts(year int) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse) {
	const method = "FindMonthlyTransferAmounts"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data := s.mencache.GetCachedMonthTransferAmounts(year); data != nil {
		logSuccess("Successfully fetched monthly transfer amounts from cache", zap.Int("year", year))
		return data, nil
	}

	amounts, err := s.transferStatisticRepository.GetMonthlyTransferAmounts(year)
	if err != nil {
		return s.errorhandler.HandleMonthlyTransferAmountsError(err, method, "FAILED_FIND_MONTHLY_TRANSFER_AMOUNTS", span, &status, zap.Error(err))
	}

	responseAmounts := s.mapping.ToTransferResponsesMonthAmount(amounts)

	s.mencache.SetCachedMonthTransferAmounts(year, responseAmounts)

	logSuccess("Successfully fetched monthly transfer amounts", zap.Int("year", year))

	return responseAmounts, nil
}

func (s *transferStatisticService) FindYearlyTransferAmounts(year int) ([]*response.TransferYearAmountResponse, *response.ErrorResponse) {
	const method = "FindYearlyTransferAmounts"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data := s.mencache.GetCachedYearlyTransferAmounts(year); data != nil {
		logSuccess("Successfully fetched yearly transfer amounts from cache", zap.Int("year", year))
		return data, nil
	}

	amounts, err := s.transferStatisticRepository.GetYearlyTransferAmounts(year)
	if err != nil {
		return s.errorhandler.HandleYearlyTransferAmountsError(err, method, "FAILED_FIND_YEARLY_TRANSFER_AMOUNTS", span, &status, zap.Error(err))
	}

	responseAmounts := s.mapping.ToTransferResponsesYearAmount(amounts)

	logSuccess("Successfully fetched yearly transfer amounts", zap.Int("year", year))

	return responseAmounts, nil
}

func (s *transferStatisticService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *transferStatisticService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
