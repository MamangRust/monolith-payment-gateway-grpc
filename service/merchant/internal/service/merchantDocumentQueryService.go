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
	merchantdocument_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_document_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantDocumentQueryService struct {
	ctx                             context.Context
	trace                           trace.Tracer
	mencache                        mencache.MerchantDocumentQueryCache
	errorhandler                    errorhandler.MerchantDocumentQueryErrorHandler
	merchantDocumentQueryRepository repository.MerchantDocumentQueryRepository
	logger                          logger.LoggerInterface
	mapping                         responseservice.MerchantDocumentResponseMapper
	requestCounter                  *prometheus.CounterVec
	requestDuration                 *prometheus.HistogramVec
}

func NewMerchantDocumentQueryService(
	ctx context.Context,
	mencache mencache.MerchantDocumentQueryCache,
	errorhandler errorhandler.MerchantDocumentQueryErrorHandler,
	merchantDocumentQueryRepository repository.MerchantDocumentQueryRepository,
	logger logger.LoggerInterface,
	mapping responseservice.MerchantDocumentResponseMapper,
) *merchantDocumentQueryService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "merchant_document_query_request_count",
		Help: "Number of merchant document query requests MerchantDocumentQueryService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "merchant_document_query_request_duration_seconds",
		Help:    "The duration of requests MerchantDocumentQueryService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	return &merchantDocumentQueryService{
		ctx:                             ctx,
		mencache:                        mencache,
		errorhandler:                    errorhandler,
		trace:                           otel.Tracer("merchant-document-query-service"),
		merchantDocumentQueryRepository: merchantDocumentQueryRepository,
		logger:                          logger,
		mapping:                         mapping,
		requestCounter:                  requestCounter,
		requestDuration:                 requestDuration,
	}
}

func (s *merchantDocumentQueryService) FindAll(req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindAll"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedMerchantDocuments(req); found {
		logSuccess("Successfully retrieved all merchant document records from cache", zap.Int("total", *total), zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	merchantDocuments, total, err := s.merchantDocumentQueryRepository.FindAllDocuments(req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_MERCHANT_DOCUMENTS", span, &status, zap.Error(err))
	}

	merchantResponse := s.mapping.ToMerchantDocumentsResponse(merchantDocuments)

	s.mencache.SetCachedMerchantDocuments(req, merchantResponse, total)

	logSuccess("Successfully retrieved all merchant document records", zap.Int("total", *total), zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return merchantResponse, total, nil
}

func (s *merchantDocumentQueryService) FindById(merchant_id int) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	const method = "FindById"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("merchantDocument.id", merchant_id))

	defer func() {
		end(status)
	}()

	if cachedMerchant := s.mencache.GetCachedMerchantDocument(merchant_id); cachedMerchant != nil {
		logSuccess("Successfully found merchant document by ID from cache", zap.Int("merchantDocument.id", merchant_id))
		return cachedMerchant, nil
	}

	merchantDocument, err := s.merchantDocumentQueryRepository.FindById(merchant_id)
	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantDocumentResponse](s.logger, err, method, "FAILED_FIND_MERCHANT_DOCUMENT_BY_ID", span, &status, merchantdocument_errors.ErrFailedFindMerchantDocumentById, zap.Error(err))

	}

	merchantResponse := s.mapping.ToMerchantDocumentResponse(merchantDocument)

	s.mencache.SetCachedMerchantDocument(merchant_id, merchantResponse)

	logSuccess("Successfully found merchant document by ID", zap.Int("merchantDocument.id", merchant_id))

	return merchantResponse, nil
}

func (s *merchantDocumentQueryService) FindByActive(req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponse, *int, *response.ErrorResponse) {
	const method = "FindByActive"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedMerchantDocuments(req); found {
		logSuccess("Successfully retrieved active merchant document records from cache", zap.Int("total", *total), zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	merchantDocuments, total, err := s.merchantDocumentQueryRepository.FindByActive(req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_MERCHANT_DOCUMENTS_ACTIVE", span, &status)
	}

	merchantResponse := s.mapping.ToMerchantDocumentsResponse(merchantDocuments)

	s.mencache.SetCachedMerchantDocuments(req, merchantResponse, total)

	logSuccess("Successfully retrieved active merchant document records", zap.Int("total", *total), zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return merchantResponse, total, nil
}

func (s *merchantDocumentQueryService) FindByTrashed(req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByTrashed"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedMerchantDocumentsTrashed(req); found {
		logSuccess("Successfully retrieved trashed merchant document records from cache", zap.Int("total", *total), zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	merchantDocuments, total, err := s.merchantDocumentQueryRepository.FindByTrashed(req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_ALL_MERCHANT_DOCUMENTS_TRASHED", span, &status, zap.Error(err))
	}

	merchantResponse := s.mapping.ToMerchantDocumentsResponseDeleteAt(merchantDocuments)

	s.mencache.SetCachedMerchantDocumentsTrashed(req, merchantResponse, total)

	logSuccess("Successfully retrieved trashed merchant document records", zap.Int("total", *total), zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return merchantResponse, total, nil
}

func (s *merchantDocumentQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func (s *merchantDocumentQueryService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *merchantDocumentQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
