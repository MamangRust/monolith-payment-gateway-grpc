package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/errorhandler"
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

type merchantTransactionService struct {
	ctx                           context.Context
	errorHandler                  errorhandler.MerchantTransactionErrorHandler
	trace                         trace.Tracer
	merchantTransactionRepository repository.MerchantTransactionRepository
	logger                        logger.LoggerInterface
	mapping                       responseservice.MerchantResponseMapper
	requestCounter                *prometheus.CounterVec
	requestDuration               *prometheus.HistogramVec
}

func NewMerchantTransactionService(ctx context.Context,
	errorHandler errorhandler.MerchantTransactionErrorHandler,
	merchantTransactionRepository repository.MerchantTransactionRepository, logger logger.LoggerInterface, mapping responseservice.MerchantResponseMapper) *merchantTransactionService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_transaction_service_requests_total",
			Help: "Total number of requests to the MerchantTransactionService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_transaction_service_request_duration_seconds",
			Help:    "Histogram of request durations for the MerchantTransactionService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &merchantTransactionService{
		ctx:                           ctx,
		errorHandler:                  errorHandler,
		trace:                         otel.Tracer("merchant-transaction-service"),
		merchantTransactionRepository: merchantTransactionRepository,
		logger:                        logger,
		mapping:                       mapping,
		requestCounter:                requestCounter,
		requestDuration:               requestDuration,
	}
}

func (s *merchantTransactionService) FindAllTransactions(req *requests.FindAllMerchantTransactions) ([]*response.MerchantTransactionResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindAllTransactions"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	merchants, totalRecords, err := s.merchantTransactionRepository.FindAllTransactions(req)

	if err != nil {
		return s.errorHandler.HandleRepositoryAllError(
			err, method, "FAILED_FIND_ALL_TRANSACTIONS", span, &status, zap.Error(err),
		)
	}

	merchantResponses := s.mapping.ToMerchantsTransactionResponse(merchants)

	logSuccess("Successfully retrieved all merchant records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return merchantResponses, totalRecords, nil
}

func (s *merchantTransactionService) FindAllTransactionsByMerchant(req *requests.FindAllMerchantTransactionsById) ([]*response.MerchantTransactionResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindAllTransactionsByMerchant"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	merchants, totalRecords, err := s.merchantTransactionRepository.FindAllTransactionsByMerchant(req)

	if err != nil {
		return s.errorHandler.HandleRepositoryAllError(
			err, method, "FAILED_FIND_ALL_TRANSACTIONS_BY_MERCHANT", span, &status, zap.Error(err),
		)
	}

	merchantResponses := s.mapping.ToMerchantsTransactionResponse(merchants)

	logSuccess("Successfully retrieved all merchant records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return merchantResponses, totalRecords, nil
}

func (s *merchantTransactionService) FindAllTransactionsByApikey(req *requests.FindAllMerchantTransactionsByApiKey) ([]*response.MerchantTransactionResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindAllTransactionsByApikey"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	merchants, totalRecords, err := s.merchantTransactionRepository.FindAllTransactionsByApikey(req)

	if err != nil {
		return s.errorHandler.HandleRepositoryAllError(
			err, method, "FAILED_FIND_ALL_TRANSACTIONS_BY_APIKEY", span, &status, zap.Error(err),
		)
	}

	merchantResponses := s.mapping.ToMerchantsTransactionResponse(merchants)

	logSuccess("Successfully retrieved all merchant records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return merchantResponses, totalRecords, nil
}

func (s *merchantTransactionService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func (s *merchantTransactionService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *merchantTransactionService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
