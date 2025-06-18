package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-topup/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type topupQueryService struct {
	ctx                  context.Context
	errorhandler         errorhandler.TopupQueryErrorHandler
	mencache             mencache.TopupQueryCache
	trace                trace.Tracer
	topupQueryRepository repository.TopupQueryRepository
	logger               logger.LoggerInterface
	mapping              responseservice.TopupResponseMapper
	requestCounter       *prometheus.CounterVec
	requestDuration      *prometheus.HistogramVec
}

func NewTopupQueryService(
	ctx context.Context, errorhandler errorhandler.TopupQueryErrorHandler,
	mencache mencache.TopupQueryCache, topupQueryRepository repository.TopupQueryRepository, logger logger.LoggerInterface, mapping responseservice.TopupResponseMapper) *topupQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "topup_query_service_request_total",
			Help: "Total number of requests to the TopupQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "topup_query_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TopupQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &topupQueryService{
		ctx:                  ctx,
		errorhandler:         errorhandler,
		mencache:             mencache,
		trace:                otel.Tracer("topup-query-service"),
		topupQueryRepository: topupQueryRepository,
		logger:               logger,
		mapping:              mapping,
		requestCounter:       requestCounter,
		requestDuration:      requestDuration,
	}
}

func (s *topupQueryService) FindAll(req *requests.FindAllTopups) ([]*response.TopupResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindAll"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedTopupsCache(req); found {
		logSuccess("Successfully retrieved all topup records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	topups, totalRecords, err := s.topupQueryRepository.FindAllTopups(req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_TOPUPS", span, &status, topup_errors.ErrFailedFindAllTopups, zap.Error(err))
	}

	so := s.mapping.ToTopupResponses(topups)

	s.mencache.SetCachedTopupsCache(req, so, totalRecords)

	logSuccess("Successfully retrieved all topup records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *topupQueryService) FindAllByCardNumber(req *requests.FindAllTopupsByCardNumber) ([]*response.TopupResponse, *int, *response.ErrorResponse) {

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search
	cardNumber := req.CardNumber

	const method = "FindAll"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search), attribute.String("cardNumber", cardNumber))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCacheTopupByCardCache(req); found {
		logSuccess("Successfully retrieved all topup records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	topups, totalRecords, err := s.topupQueryRepository.FindAllTopupByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_TOPUPS_BY_CARD", span, &status, topup_errors.ErrFailedFindAllTopupsByCardNumber)
	}

	so := s.mapping.ToTopupResponses(topups)

	s.mencache.SetCacheTopupByCardCache(req, so, totalRecords)

	logSuccess("Successfully retrieved all topup records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *topupQueryService) FindById(topupID int) (*response.TopupResponse, *response.ErrorResponse) {
	const method = "FindById"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("topup.id", topupID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedTopupCache(topupID); found {
		logSuccess("Successfully retrieved topup from cache", zap.Int("topup.id", topupID))
		return data, nil
	}

	span.SetAttributes(attribute.String("cache.status", "miss"))

	topup, err := s.topupQueryRepository.FindById(topupID)
	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.TopupResponse](s.logger, err, method, "FAILED_FIND_TOPUP", span, &status, topup_errors.ErrFailedFindTopupById, zap.Error(err))
	}

	so := s.mapping.ToTopupResponse(topup)

	s.mencache.SetCachedTopupCache(so)

	logSuccess("Successfully retrieved topup", zap.Int("topup.id", topupID))

	return so, nil
}

func (s *topupQueryService) FindByActive(req *requests.FindAllTopups) ([]*response.TopupResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindAll"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedTopupActiveCache(req); found {
		logSuccess("Successfully retrieved all topup records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	topups, totalRecords, err := s.topupQueryRepository.FindByActive(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_ACTIVE_TOPUPS", span, &status, topup_errors.ErrFailedFindActiveTopups)
	}

	so := s.mapping.ToTopupResponsesDeleteAt(topups)

	s.mencache.SetCachedTopupActiveCache(req, so, totalRecords)

	logSuccess("Successfully retrieved all topup records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *topupQueryService) FindByTrashed(req *requests.FindAllTopups) ([]*response.TopupResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindByTrashed"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedTopupTrashedCache(req); found {
		logSuccess("Successfully retrieved all topup records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	topups, totalRecords, err := s.topupQueryRepository.FindByTrashed(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_TRASHED_TOPUPS", span, &status, topup_errors.ErrFailedFindTrashedTopups)
	}

	so := s.mapping.ToTopupResponsesDeleteAt(topups)

	logSuccess("Successfully retrieved all topup records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *topupQueryService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *topupQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func (s *topupQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
