package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	traceunic "github.com/MamangRust/monolith-payment-gateway-pkg/trace_unic"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transactionQueryService struct {
	ctx                        context.Context
	trace                      trace.Tracer
	transactionQueryRepository repository.TransactionQueryRepository
	logger                     logger.LoggerInterface
	mapping                    responseservice.TransactionResponseMapper
	requestCounter             *prometheus.CounterVec
	requestDuration            *prometheus.HistogramVec
}

func NewTransactionQueryService(ctx context.Context, transactionQueryRepository repository.TransactionQueryRepository, logger logger.LoggerInterface, mapping responseservice.TransactionResponseMapper) *transactionQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transaction_query_service_request_total",
			Help: "Total number of requests to the TransactionQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transaction_query_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TransactionQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &transactionQueryService{
		ctx:                        ctx,
		trace:                      otel.Tracer("transaction-query-service"),
		transactionQueryRepository: transactionQueryRepository,
		logger:                     logger,
		mapping:                    mapping,
		requestCounter:             requestCounter,
		requestDuration:            requestDuration,
	}
}

func (s *transactionQueryService) FindAll(req *requests.FindAllTransactions) ([]*response.TransactionResponse, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindAll", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindAll")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching transaction",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	transactions, totalRecords, err := s.transactionQueryRepository.FindAllTransactions(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_TRANSACTION")

		s.logger.Error("Failed to fetch transaction",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch transaction")

		status = "failed_find_all_transaction"

		return nil, nil, transaction_errors.ErrFailedFindAllTransactions
	}

	so := s.mapping.ToTransactionsResponse(transactions)

	s.logger.Debug("Successfully fetched transaction",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *transactionQueryService) FindAllByCardNumber(req *requests.FindAllTransactionCardNumber) ([]*response.TransactionResponse, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindAllByCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindAllByCardNumber")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching transaction",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	transactions, totalRecords, err := s.transactionQueryRepository.FindAllTransactionByCardNumber(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_TRANSACTION_BY_CARD_NUMBER")

		s.logger.Error("Failed to fetch transaction",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch transaction")

		status = "failed_find_all_transaction_by_card_number"

		return nil, nil, transaction_errors.ErrFailedFindAllByCardNumber
	}

	so := s.mapping.ToTransactionsResponse(transactions)

	s.logger.Debug("Successfully fetched transaction",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *transactionQueryService) FindById(transactionID int) (*response.TransactionResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindById", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindById")
	defer span.End()

	span.SetAttributes(
		attribute.Int("transaction_id", transactionID),
	)

	s.logger.Debug("Fetching transaction by ID", zap.Int("transaction_id", transactionID))

	transaction, err := s.transactionQueryRepository.FindById(transactionID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TRANSACTION_BY_ID")

		s.logger.Error("failed to fetch transaction by ID", zap.String("trace.id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch transaction by ID")
		status = "failed_find_transaction_by_id"

		return nil, transaction_errors.ErrTransactionNotFound
	}

	so := s.mapping.ToTransactionResponse(transaction)

	s.logger.Debug("Successfully fetched transaction", zap.Int("transaction_id", transactionID))

	return so, nil
}

func (s *transactionQueryService) FindByActive(req *requests.FindAllTransactions) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByActive", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindByActive")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching active transaction",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	transactions, totalRecords, err := s.transactionQueryRepository.FindByActive(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_BY_ACTIVE_TRANSACTION")

		s.logger.Error("Failed to fetch active transaction",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch active transaction")
		status = "failed_find_by_active_transaction"

		return nil, nil, transaction_errors.ErrFailedFindByActiveTransactions
	}

	so := s.mapping.ToTransactionsResponseDeleteAt(transactions)

	s.logger.Debug("Successfully fetched active transaction",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *transactionQueryService) FindByTrashed(req *requests.FindAllTransactions) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByTrashed", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindByTrashed")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching trashed transaction",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	transactions, totalRecords, err := s.transactionQueryRepository.FindByTrashed(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_BY_TRASHED_TRANSACTION")

		s.logger.Error("Failed to fetch trashed transaction",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch trashed transaction")
		status = "failed_find_by_trashed_transaction"

		return nil, nil, transaction_errors.ErrFailedFindByTrashedTransactions
	}

	so := s.mapping.ToTransactionsResponseDeleteAt(transactions)

	s.logger.Debug("Successfully fetched trashed transaction",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *transactionQueryService) FindTransactionByMerchantId(merchant_id int) ([]*response.TransactionResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindTransactionByMerchantId", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindTransactionByMerchantId")
	defer span.End()

	span.SetAttributes(
		attribute.Int("merchant_id", merchant_id),
	)

	s.logger.Debug("Starting FindTransactionByMerchantId process",
		zap.Int("merchantID", merchant_id),
	)

	res, err := s.transactionQueryRepository.FindTransactionByMerchantId(merchant_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TRANSACTION_BY_MERCHANT_ID")

		s.logger.Error("Failed to find transaction by merchant ID", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find transaction by merchant ID")
		status = "failed_find_transaction_by_merchant_id"

		return nil, transaction_errors.ErrFailedFindByMerchantID
	}

	so := s.mapping.ToTransactionsResponse(res)

	s.logger.Debug("Successfully fetched transaction by merchant ID", zap.Int("merchant_id", merchant_id))

	return so, nil
}

func (s *transactionQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
