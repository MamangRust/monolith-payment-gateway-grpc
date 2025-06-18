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
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantQueryService struct {
	ctx                     context.Context
	trace                   trace.Tracer
	merchantQueryRepository repository.MerchantQueryRepository
	logger                  logger.LoggerInterface
	mapping                 responseservice.MerchantResponseMapper
	errorhandler            errorhandler.MerchantQueryErrorHandler
	mencache                mencache.MerchantQueryCache
	requestCounter          *prometheus.CounterVec
	requestDuration         *prometheus.HistogramVec
}

func NewMerchantQueryService(ctx context.Context, merchantQueryRepository repository.MerchantQueryRepository, errorhandler errorhandler.MerchantQueryErrorHandler, mencache mencache.MerchantQueryCache, logger logger.LoggerInterface, mapping responseservice.MerchantResponseMapper) *merchantQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_query_service_requests_total",
			Help: "Total number of requests to the MerchantQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_query_service_request_duration_seconds",
			Help:    "Histogram of request durations for the MerchantQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &merchantQueryService{
		ctx:                     ctx,
		trace:                   otel.Tracer("merchant-query-service"),
		merchantQueryRepository: merchantQueryRepository,
		logger:                  logger,
		mapping:                 mapping,
		errorhandler:            errorhandler,
		mencache:                mencache,
		requestCounter:          requestCounter,
		requestDuration:         requestDuration,
	}
}

func (s *merchantQueryService) FindAll(req *requests.FindAllMerchants) ([]*response.MerchantResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindAll"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedMerchants(req); found {
		logSuccess("Successfully retrieved all merchant records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))

		return data, total, nil
	}

	merchants, totalRecords, err := s.merchantQueryRepository.FindAllMerchants(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_MERCHANTS", span, &status, zap.Error(err))
	}

	merchantResponses := s.mapping.ToMerchantsResponse(merchants)

	s.mencache.SetCachedMerchants(req, merchantResponses, totalRecords)

	logSuccess("Successfully retrieved all merchant records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return merchantResponses, totalRecords, nil
}

func (s *merchantQueryService) FindById(merchantID int) (*response.MerchantResponse, *response.ErrorResponse) {
	const method = "FindById"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("merchant.id", merchantID))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetCachedMerchant(merchantID); found {
		logSuccess("Successfully retrieved merchant from cache", zap.Int("merchant.id", merchantID))
		return cachedMerchant, nil
	}

	merchant, err := s.merchantQueryRepository.FindById(merchantID)
	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantResponse](s.logger, err, method, "FAILED_FIND_MERCHANT_BY_ID", span, &status, merchant_errors.ErrFailedFindMerchantById, zap.Error(err))
	}

	merchantResponse := s.mapping.ToMerchantResponse(merchant)

	s.mencache.SetCachedMerchant(merchantResponse)

	logSuccess("Successfully retrieved merchant", zap.Int("merchant.id", merchantID))

	return merchantResponse, nil
}

func (s *merchantQueryService) FindByActive(req *requests.FindAllMerchants) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByActive"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedMerchantActive(req); found {
		logSuccess("Successfully fetched active merchants from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	merchants, totalRecords, err := s.merchantQueryRepository.FindByActive(req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_ACTIVE_MERCHANTS", span, &status, merchant_errors.ErrFailedFindActiveMerchants,
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search))
	}

	merchantResponses := s.mapping.ToMerchantsResponseDeleteAt(merchants)

	s.mencache.SetCachedMerchantActive(req, merchantResponses, totalRecords)

	logSuccess("Successfully fetched active merchants", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return merchantResponses, totalRecords, nil
}

func (s *merchantQueryService) FindByTrashed(req *requests.FindAllMerchants) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByTrashed"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedMerchantTrashed(req); found {
		logSuccess("Successfully fetched trashed merchants from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	merchants, totalRecords, err := s.merchantQueryRepository.FindByTrashed(req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_TRASHED_MERCHANTS", span, &status, merchant_errors.ErrFailedFindTrashedMerchants,
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search))
	}

	merchantResponses := s.mapping.ToMerchantsResponseDeleteAt(merchants)

	s.mencache.SetCachedMerchantTrashed(req, merchantResponses, totalRecords)

	logSuccess("Successfully fetched trashed merchants", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return merchantResponses, totalRecords, nil
}

func (s *merchantQueryService) FindByApiKey(apiKey string) (*response.MerchantResponse, *response.ErrorResponse) {
	const method = "FindByApiKey"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.String("api_key", apiKey))

	defer func() {
		end(status)
	}()

	if cachedMerchant, found := s.mencache.GetCachedMerchantByApiKey(apiKey); found {
		logSuccess("Successfully found merchant by API key from cache", zap.String("api_key", apiKey))
		return cachedMerchant, nil
	}

	merchant, err := s.merchantQueryRepository.FindByApiKey(apiKey)
	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantResponse](s.logger, err, method, "FAILED_FIND_MERCHANT_BY_API_KEY", span, &status, merchant_errors.ErrMerchantNotFoundRes, zap.Error(err))
	}

	merchantResponse := s.mapping.ToMerchantResponse(merchant)

	s.mencache.SetCachedMerchantByApiKey(apiKey, merchantResponse)

	logSuccess("Successfully found merchant by API key", zap.String("api_key", apiKey))

	return merchantResponse, nil
}

func (s *merchantQueryService) FindByMerchantUserId(userID int) ([]*response.MerchantResponse, *response.ErrorResponse) {
	const method = "FindByMerchantUserId"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("user.id", userID))

	defer func() {
		end(status)
	}()

	if cachedMerchants, found := s.mencache.GetCachedMerchantsByUserId(userID); found {
		logSuccess("Successfully found merchants by user ID from cache", zap.Int("user.id", userID), zap.Int("count", len(cachedMerchants)))

		return cachedMerchants, nil
	}

	merchants, err := s.merchantQueryRepository.FindByMerchantUserId(userID)
	if err != nil {
		return s.errorhandler.HandleRepositoryListError(err, method, "FAILED_FIND_MERCHANT_BY_USER_ID", span, &status, zap.Error(err))

	}

	merchantResponses := s.mapping.ToMerchantsResponse(merchants)

	s.mencache.SetCachedMerchantsByUserId(userID, merchantResponses)

	logSuccess("Successfully found merchants by user ID", zap.Int("user.id", userID), zap.Int("count", len(merchantResponses)))

	return merchantResponses, nil
}

func (s *merchantQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func (s *merchantQueryService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *merchantQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
