package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-transaction/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transactionQueryService struct {
	ctx                        context.Context
	mencache                   mencache.TransactinQueryCache
	errorhandler               errorhandler.TransactionQueryErrorHandler
	trace                      trace.Tracer
	transactionQueryRepository repository.TransactionQueryRepository
	logger                     logger.LoggerInterface
	mapping                    responseservice.TransactionResponseMapper
	requestCounter             *prometheus.CounterVec
	requestDuration            *prometheus.HistogramVec
}

func NewTransactionQueryService(ctx context.Context, mencache mencache.TransactinQueryCache,
	errorhandler errorhandler.TransactionQueryErrorHandler, transactionQueryRepository repository.TransactionQueryRepository, logger logger.LoggerInterface, mapping responseservice.TransactionResponseMapper) *transactionQueryService {
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
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &transactionQueryService{
		ctx:                        ctx,
		mencache:                   mencache,
		errorhandler:               errorhandler,
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

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
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

	if data, total, found := s.mencache.GetCachedTransactionsCache(req); found {
		s.logger.Debug("Successfully fetched transaction from cache",
			zap.Int("totalRecords", *total))
		return data, total, nil
	}

	transactions, totalRecords, err := s.transactionQueryRepository.FindAllTransactions(req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, "FindAll", "FAILED_FIND_ALL", span, &status, zap.Error(err))
	}

	responseData := s.mapping.ToTransactionsResponse(transactions)

	s.mencache.SetCachedTransactionsCache(req, responseData, totalRecords)

	s.logger.Debug("Successfully fetched transaction from database",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return responseData, totalRecords, nil
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

	if data, total, found := s.mencache.GetCachedTransactionByCardNumberCache(req); found {
		s.logger.Debug("Successfully fetched transaction from cache",
			zap.Int("totalRecords", *total))
		return data, total, nil
	}

	transactions, totalRecords, err := s.transactionQueryRepository.FindAllTransactionByCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, "FindAllByCardNumber", "FAILED_FIND_ALL_BYCARD", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionsResponse(transactions)

	s.mencache.SetCachedTransactionByCardNumberCache(req, so, totalRecords)

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

	if data := s.mencache.GetCachedTransactionCache(transactionID); data != nil {
		s.logger.Debug("Successfully fetched transaction from cache")
		return data, nil
	}

	transaction, err := s.transactionQueryRepository.FindById(transactionID)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "FindById", "FAILED_FIND_TRANSACTION", span, &status, transaction_errors.ErrTransactionNotFound, zap.Error(err))
	}

	so := s.mapping.ToTransactionResponse(transaction)

	s.mencache.SetCachedTransactionCache(so)

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

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
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

	if data, total, found := s.mencache.GetCachedTransactionActiveCache(req); found {
		s.logger.Debug("Successfully fetched active transaction from cache",
			zap.Int("totalRecords", *total))
		return data, total, nil
	}

	transactions, totalRecords, err := s.transactionQueryRepository.FindByActive(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, "FindByActive", "FAILED_FIND_BY_ACTIVE", span, &status, transaction_errors.ErrFailedFindByActiveTransactions, zap.Error(err))
	}

	so := s.mapping.ToTransactionsResponseDeleteAt(transactions)

	s.mencache.SetCachedTransactionActiveCache(req, so, totalRecords)

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

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
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

	if data, total, found := s.mencache.GetCachedTransactionTrashedCache(req); found {
		s.logger.Debug("Successfully fetched trashed transaction from cache",
			zap.Int("totalRecords", *total))
		return data, total, nil
	}

	transactions, totalRecords, err := s.transactionQueryRepository.FindByTrashed(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, "FindByTrashed", "FAILED_FIND_BY_TRASHED", span, &status, transaction_errors.ErrFailedFindByTrashedTransactions, zap.Error(err))
	}

	so := s.mapping.ToTransactionsResponseDeleteAt(transactions)

	s.mencache.SetCachedTransactionTrashedCache(req, so, totalRecords)

	s.logger.Debug("Successfully fetched trashed transaction",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *transactionQueryService) FindTransactionByMerchantId(merchantID int) ([]*response.TransactionResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("FindTransactionByMerchantId", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindTransactionByMerchantId")
	defer span.End()

	span.SetAttributes(attribute.Int("merchant_id", merchantID))
	s.logger.Debug("Starting FindTransactionByMerchantId process", zap.Int("merchant_id", merchantID))

	if data := s.mencache.GetCachedTransactionByMerchantIdCache(merchantID); data != nil {
		s.logger.Debug("Successfully fetched transaction from cache", zap.Int("merchant_id", merchantID))
		return data, nil
	}

	res, err := s.transactionQueryRepository.FindTransactionByMerchantId(merchantID)
	if err != nil {
		return s.errorhandler.HanldeRepositoryListError(err, "FindTransactionByMerchantId", "FAILED_FIND_TRANSACTION_BY_MERCHANT_ID", span, &status, transaction_errors.ErrFailedFindByMerchantID, zap.Error(err))
	}

	responseData := s.mapping.ToTransactionsResponse(res)

	s.mencache.SetCachedTransactionByMerchantIdCache(merchantID, responseData)

	s.logger.Debug("Successfully fetched transaction by merchant ID", zap.Int("merchant_id", merchantID))
	return responseData, nil
}

func (s *transactionQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}

func (s *transactionQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}
