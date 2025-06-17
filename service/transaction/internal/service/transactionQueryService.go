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
	"go.opentelemetry.io/otel/codes"
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
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindAll"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedTransactionsCache(req); found {
		logSuccess("Successfully retrieved all transaction records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	transactions, totalRecords, err := s.transactionQueryRepository.FindAllTransactions(req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL", span, &status, zap.Error(err))
	}

	responseData := s.mapping.ToTransactionsResponse(transactions)

	s.mencache.SetCachedTransactionsCache(req, responseData, totalRecords)

	logSuccess("Successfully retrieved all transaction records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return responseData, totalRecords, nil
}

func (s *transactionQueryService) FindAllByCardNumber(req *requests.FindAllTransactionCardNumber) ([]*response.TransactionResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search
	cardNumber := req.CardNumber

	const method = "FindAll"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search), attribute.String("cardNumber", cardNumber))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedTransactionByCardNumberCache(req); found {
		logSuccess("Successfully retrieved all transaction records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	transactions, totalRecords, err := s.transactionQueryRepository.FindAllTransactionByCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_BYCARD", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionsResponse(transactions)

	s.mencache.SetCachedTransactionByCardNumberCache(req, so, totalRecords)

	logSuccess("Successfully retrieved all transaction records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *transactionQueryService) FindById(transactionID int) (*response.TransactionResponse, *response.ErrorResponse) {
	const method = "FindById"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("transaction.id", transactionID))

	defer func() {
		end(status)
	}()

	if data := s.mencache.GetCachedTransactionCache(transactionID); data != nil {
		logSuccess("Successfully fetched transaction from cache", zap.Int("transaction.id", transactionID))
		return data, nil
	}

	transaction, err := s.transactionQueryRepository.FindById(transactionID)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_TRANSACTION", span, &status, transaction_errors.ErrTransactionNotFound, zap.Error(err))
	}

	so := s.mapping.ToTransactionResponse(transaction)

	s.mencache.SetCachedTransactionCache(so)

	logSuccess("Successfully fetched transaction", zap.Int("transaction_id", transactionID))

	return so, nil
}

func (s *transactionQueryService) FindByActive(req *requests.FindAllTransactions) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindByActive"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedTransactionActiveCache(req); found {
		logSuccess("Successfully fetched active transaction from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	transactions, totalRecords, err := s.transactionQueryRepository.FindByActive(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_BY_ACTIVE", span, &status, transaction_errors.ErrFailedFindByActiveTransactions, zap.Error(err))
	}

	so := s.mapping.ToTransactionsResponseDeleteAt(transactions)

	s.mencache.SetCachedTransactionActiveCache(req, so, totalRecords)

	logSuccess("Successfully fetched active transaction", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *transactionQueryService) FindByTrashed(req *requests.FindAllTransactions) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindByTrashed"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedTransactionTrashedCache(req); found {
		logSuccess("Successfully fetched trashed transaction from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	transactions, totalRecords, err := s.transactionQueryRepository.FindByTrashed(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_BY_TRASHED", span, &status, transaction_errors.ErrFailedFindByTrashedTransactions, zap.Error(err))
	}

	so := s.mapping.ToTransactionsResponseDeleteAt(transactions)

	s.mencache.SetCachedTransactionTrashedCache(req, so, totalRecords)

	logSuccess("Successfully fetched trashed transaction", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *transactionQueryService) FindTransactionByMerchantId(merchantID int) ([]*response.TransactionResponse, *response.ErrorResponse) {
	const method = "FindTransactionByMerchantId"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("merchant.id", merchantID))

	defer func() {
		end(status)
	}()

	if data := s.mencache.GetCachedTransactionByMerchantIdCache(merchantID); data != nil {
		logSuccess("Successfully fetched transaction by merchant ID from cache", zap.Int("merchant.id", merchantID))
		return data, nil
	}

	res, err := s.transactionQueryRepository.FindTransactionByMerchantId(merchantID)
	if err != nil {
		return s.errorhandler.HanldeRepositoryListError(err, method, "FAILED_FIND_TRANSACTION_BY_MERCHANT_ID", span, &status, transaction_errors.ErrFailedFindByMerchantID, zap.Error(err))
	}

	responseData := s.mapping.ToTransactionsResponse(res)

	s.mencache.SetCachedTransactionByMerchantIdCache(merchantID, responseData)

	logSuccess("Successfully fetched transaction by merchant ID", zap.Int("merchant.id", merchantID))

	return responseData, nil
}

func (s *transactionQueryService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *transactionQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func (s *transactionQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
