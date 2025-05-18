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

type merchantQueryService struct {
	ctx                     context.Context
	trace                   trace.Tracer
	merchantQueryRepository repository.MerchantQueryRepository
	logger                  logger.LoggerInterface
	mapping                 responseservice.MerchantResponseMapper
	requestCounter          *prometheus.CounterVec
	requestDuration         *prometheus.HistogramVec
}

func NewMerchantQueryService(ctx context.Context, merchantQueryRepository repository.MerchantQueryRepository, logger logger.LoggerInterface, mapping responseservice.MerchantResponseMapper) *merchantQueryService {
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
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &merchantQueryService{
		ctx:                     ctx,
		trace:                   otel.Tracer("merchant-query-service"),
		merchantQueryRepository: merchantQueryRepository,
		logger:                  logger,
		mapping:                 mapping,
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

	merchants, totalRecords, err := s.merchantQueryRepository.FindAllMerchants(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_MERCHANTS")

		s.logger.Error("Failed to retrieve all merchants",
			zap.Error(err),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve all merchants")
		status = "failed_to_find_all_merchants"

		return nil, nil, merchant_errors.ErrFailedFindAllMerchants
	}

	merchantResponses := s.mapping.ToMerchantsResponse(merchants)

	s.logger.Debug("Successfully all merchant records",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchantResponses, totalRecords, nil
}

func (s *merchantQueryService) FindById(merchant_id int) (*response.MerchantResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindById", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindById")
	defer span.End()

	span.SetAttributes(
		attribute.Int("merchant_id", merchant_id),
	)

	s.logger.Debug("Finding merchant by ID", zap.Int("merchant_id", merchant_id))

	res, err := s.merchantQueryRepository.FindById(merchant_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MERCHANT_BY_ID")

		s.logger.Error("Failed to find merchant by ID",
			zap.Error(err),
			zap.Int("merchant_id", merchant_id),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find merchant by ID")
		status = "failed_to_find_merchant_by_id"

		return nil, merchant_errors.ErrMerchantNotFoundRes
	}

	so := s.mapping.ToMerchantResponse(res)

	return so, nil
}

func (s *merchantQueryService) FindByActive(req *requests.FindAllMerchants) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByActive", status, startTime)
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

	s.logger.Debug("Fetching all merchant active",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	merchants, totalRecords, err := s.merchantQueryRepository.FindByActive(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ACTIVE_MERCHANTS")

		s.logger.Error("Failed to retrieve active merchant",
			zap.Error(err),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve active merchant")
		status = "failed_to_find_active_merchants"

		return nil, nil, merchant_errors.ErrFailedFindActiveMerchants
	}

	so := s.mapping.ToMerchantsResponseDeleteAt(merchants)

	s.logger.Debug("Successfully fetched active merchants",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *merchantQueryService) FindByTrashed(req *requests.FindAllMerchants) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse) {
	startTime := time.Now()

	status := "success"

	defer func() {
		s.recordMetrics("FindByTrashed", status, startTime)
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

	s.logger.Debug("Fetching fetched trashed merchants",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	merchants, totalRecords, err := s.merchantQueryRepository.FindByTrashed(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TRASHED_MERCHANTS")

		s.logger.Error("Failed to retrieve trashed merchant",
			zap.Error(err),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve trashed merchant")
		status = "failed_to_find_trashed_merchants"

		return nil, nil, merchant_errors.ErrFailedFindTrashedMerchants
	}

	so := s.mapping.ToMerchantsResponseDeleteAt(merchants)

	s.logger.Debug("Successfully fetched trashed merchants",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *merchantQueryService) FindByApiKey(api_key string) (*response.MerchantResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByApiKey", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindByApiKey")
	defer span.End()

	span.SetAttributes(
		attribute.String("api_key", api_key),
	)

	s.logger.Debug("Finding merchant by API key", zap.String("api_key", api_key))

	res, err := s.merchantQueryRepository.FindByApiKey(api_key)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MERCHANT_BY_API_KEY")

		s.logger.Error("Failed to retrieve merchant by api_key",
			zap.Error(err),
			zap.String("api_key", api_key),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve merchant by api_key")
		status = "failed_to_find_merchant_by_api_key"

		return nil, merchant_errors.ErrMerchantNotFoundRes
	}

	so := s.mapping.ToMerchantResponse(res)

	s.logger.Debug("Successfully found merchant by API key", zap.String("api_key", api_key))

	return so, nil
}

func (s *merchantQueryService) FindByMerchantUserId(user_id int) ([]*response.MerchantResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByMerchantUserId", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindByMerchantUserId")
	defer span.End()

	span.SetAttributes(
		attribute.Int("user_id", user_id),
	)

	s.logger.Debug("Finding merchant by user ID", zap.Int("user_id", user_id))

	res, err := s.merchantQueryRepository.FindByMerchantUserId(user_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MERCHANT_BY_USER_ID")

		s.logger.Error("Failed to retrieve merchant by user_id",
			zap.Error(err),
			zap.Int("user_id", user_id),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve merchant by user_id")
		status = "failed_to_find_merchant_by_user_id"

		return nil, merchant_errors.ErrMerchantNotFoundRes
	}

	so := s.mapping.ToMerchantsResponse(res)

	s.logger.Debug("Successfully found merchant by user ID", zap.Int("user_id", user_id))

	return so, nil
}

func (s *merchantQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
