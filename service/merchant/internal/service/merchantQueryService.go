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
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindAll", status, startTime)
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

	s.logger.Debug("Fetching all merchant records",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if data, total, found := s.mencache.GetCachedMerchants(req); found {
		s.logger.Debug("Successfully fetched merchants from cache",
			zap.Int("totalRecords", *total))
		return data, total, nil
	}

	merchants, totalRecords, err := s.merchantQueryRepository.FindAllMerchants(req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, "FindAll", "FAILED_FIND_ALL_MERCHANTS", span, &status,
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search))
	}

	merchantResponses := s.mapping.ToMerchantsResponse(merchants)

	s.mencache.SetCachedMerchants(req, merchantResponses, totalRecords)

	s.logger.Debug("Successfully retrieved all merchant records",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchantResponses, totalRecords, nil
}

func (s *merchantQueryService) FindById(merchantID int) (*response.MerchantResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindById", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindById")
	defer span.End()

	span.SetAttributes(attribute.Int("merchant_id", merchantID))

	s.logger.Debug("Finding merchant by ID", zap.Int("merchant_id", merchantID))

	if cachedMerchant := s.mencache.GetCachedMerchant(merchantID); cachedMerchant != nil {
		s.logger.Debug("Successfully fetched merchant from cache",
			zap.Int("merchant_id", merchantID))
		return cachedMerchant, nil
	}

	merchant, err := s.merchantQueryRepository.FindById(merchantID)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "FindById", "FAILED_FIND_MERCHANT_BY_ID", span, &status, zap.Error(err))
	}

	merchantResponse := s.mapping.ToMerchantResponse(merchant)

	s.mencache.SetCachedMerchant(merchantResponse)

	s.logger.Debug("Successfully found merchant by ID", zap.Int("merchant_id", merchantID))

	return merchantResponse, nil
}

func (s *merchantQueryService) FindByActive(req *requests.FindAllMerchants) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByActive", status, startTime)
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

	s.logger.Debug("Fetching active merchants",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if data, total, found := s.mencache.GetCachedMerchantActive(req); found {
		s.logger.Debug("Successfully fetched active merchants from cache",
			zap.Int("totalRecords", *total))
		return data, total, nil
	}

	merchants, totalRecords, err := s.merchantQueryRepository.FindByActive(req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, "FindByActive", "FAILED_FIND_ACTIVE_MERCHANTS", span, &status, merchant_errors.ErrFailedFindActiveMerchants,
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search))
	}

	merchantResponses := s.mapping.ToMerchantsResponseDeleteAt(merchants)

	s.mencache.SetCachedMerchantActive(req, merchantResponses, totalRecords)

	s.logger.Debug("Successfully fetched active merchants",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchantResponses, totalRecords, nil
}

func (s *merchantQueryService) FindByTrashed(req *requests.FindAllMerchants) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByTrashed", status, startTime)
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

	s.logger.Debug("Fetching trashed merchants",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if data, total, found := s.mencache.GetCachedMerchantTrashed(req); found {
		s.logger.Debug("Successfully fetched trashed merchants from cache",
			zap.Int("totalRecords", *total))
		return data, total, nil
	}

	merchants, totalRecords, err := s.merchantQueryRepository.FindByTrashed(req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, "FindByTrashed", "FAILED_FIND_TRASHED_MERCHANTS", span, &status, merchant_errors.ErrFailedFindTrashedMerchants,
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search))
	}

	merchantResponses := s.mapping.ToMerchantsResponseDeleteAt(merchants)

	s.mencache.SetCachedMerchantTrashed(req, merchantResponses, totalRecords)

	s.logger.Debug("Successfully fetched trashed merchants",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchantResponses, totalRecords, nil
}

func (s *merchantQueryService) FindByApiKey(apiKey string) (*response.MerchantResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByApiKey", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindByApiKey")
	defer span.End()

	span.SetAttributes(attribute.String("api_key", apiKey))

	s.logger.Debug("Finding merchant by API key", zap.String("api_key", apiKey))

	if cachedMerchant := s.mencache.GetCachedMerchantByApiKey(apiKey); cachedMerchant != nil {
		s.logger.Debug("Successfully fetched merchant from cache by API key",
			zap.String("api_key", apiKey))
		return cachedMerchant, nil
	}

	merchant, err := s.merchantQueryRepository.FindByApiKey(apiKey)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "FindByApiKey", "FAILED_FIND_MERCHANT_BY_API_KEY", span, &status, zap.Error(err))
	}

	merchantResponse := s.mapping.ToMerchantResponse(merchant)

	s.mencache.SetCachedMerchantByApiKey(apiKey, merchantResponse)

	s.logger.Debug("Successfully found merchant by API key", zap.String("api_key", apiKey))

	return merchantResponse, nil
}

func (s *merchantQueryService) FindByMerchantUserId(userID int) ([]*response.MerchantResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByMerchantUserId", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindByMerchantUserId")
	defer span.End()

	span.SetAttributes(attribute.Int("user_id", userID))

	s.logger.Debug("Finding merchants by user ID", zap.Int("user_id", userID))

	if cachedMerchants := s.mencache.GetCachedMerchantsByUserId(userID); cachedMerchants != nil {
		s.logger.Debug("Successfully fetched merchants from cache by user ID",
			zap.Int("user_id", userID),
			zap.Int("count", len(cachedMerchants)))
		return cachedMerchants, nil
	}

	merchants, err := s.merchantQueryRepository.FindByMerchantUserId(userID)
	if err != nil {
		return s.errorhandler.HandleRepositoryListError(err, "FindByMerchantUserId", "FAILED_FIND_MERCHANT_BY_USER_ID", span, &status, zap.Error(err))

	}

	merchantResponses := s.mapping.ToMerchantsResponse(merchants)

	s.mencache.SetCachedMerchantsByUserId(userID, merchantResponses)

	s.logger.Debug("Successfully found merchants by user ID",
		zap.Int("user_id", userID),
		zap.Int("count", len(merchantResponses)))

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

func (s *merchantQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
