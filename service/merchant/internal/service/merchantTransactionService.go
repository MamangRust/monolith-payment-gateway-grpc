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

type merchantTransactionService struct {
	ctx                           context.Context
	trace                         trace.Tracer
	merchantTransactionRepository repository.MerchantTransactionRepository
	logger                        logger.LoggerInterface
	mapping                       responseservice.MerchantResponseMapper
	requestCounter                *prometheus.CounterVec
	requestDuration               *prometheus.HistogramVec
}

func NewMerchantTransactionService(ctx context.Context, merchantTransactionRepository repository.MerchantTransactionRepository, logger logger.LoggerInterface, mapping responseservice.MerchantResponseMapper) *merchantTransactionService {

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
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &merchantTransactionService{
		ctx:                           ctx,
		trace:                         otel.Tracer("merchant-transaction-service"),
		merchantTransactionRepository: merchantTransactionRepository,
		logger:                        logger,
		mapping:                       mapping,
		requestCounter:                requestCounter,
		requestDuration:               requestDuration,
	}
}

func (s *merchantTransactionService) FindAllTransactions(req *requests.FindAllMerchantTransactions) ([]*response.MerchantTransactionResponse, *int, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindAllTransactions", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindAllTransactions")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching all merchant records",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	merchants, totalRecords, err := s.merchantTransactionRepository.FindAllTransactions(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_TRANSACTIONS")

		s.logger.Error("Failed to retrieve active merchant",
			zap.Error(err),
			zap.String("traceID", traceID),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve active merchant")
		status = "failed_find_all_transactions"

		return nil, nil, merchant_errors.ErrFailedFindAllTransactions
	}

	merchantResponses := s.mapping.ToMerchantsTransactionResponse(merchants)

	s.logger.Debug("Successfully all merchant records",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchantResponses, totalRecords, nil
}

func (s *merchantTransactionService) FindAllTransactionsByMerchant(req *requests.FindAllMerchantTransactionsById) ([]*response.MerchantTransactionResponse, *int, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindAllTransactionsByMerchant", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindAllTransactionsByMerchant")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching all merchant records",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	merchants, totalRecords, err := s.merchantTransactionRepository.FindAllTransactionsByMerchant(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_TRANSACTIONS_BY_MERCHANT")

		s.logger.Error("Failed to retrieve active merchant",
			zap.Error(err),
			zap.String("traceID", traceID),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve active merchant")
		status = "failed_find_all_transactions_by_merchant"

		return nil, nil, merchant_errors.ErrFailedFindAllTransactionsByMerchant
	}

	merchantResponses := s.mapping.ToMerchantsTransactionResponse(merchants)

	s.logger.Debug("Successfully fetched active merchant",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchantResponses, totalRecords, nil
}

func (s *merchantTransactionService) FindAllTransactionsByApikey(req *requests.FindAllMerchantTransactionsByApiKey) ([]*response.MerchantTransactionResponse, *int, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindAllTransactionsByApikey", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindAllTransactionsByApikey")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching all transaction merchant records",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	merchants, totalRecords, err := s.merchantTransactionRepository.FindAllTransactionsByApikey(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_TRANSACTIONS_BY_APIKEY")

		s.logger.Error("Failed to retrieve active merchant",
			zap.Error(err),
			zap.String("traceID", traceID),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve active merchant")
		status = "failed_find_all_transactions_by_apikey"
		return nil, nil, merchant_errors.ErrFailedFindAllTransactionsByApikey
	}

	merchantResponses := s.mapping.ToMerchantsTransactionResponse(merchants)

	s.logger.Debug("Successfully all transaction merchant records",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchantResponses, totalRecords, nil
}

func (s *merchantTransactionService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
