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

	s.logger.Debug("Fetching topups",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search),
	)

	if data, total, found := s.mencache.GetCachedTopupsCache(req); found {
		s.logger.Debug("Successfully fetched topups from cache",
			zap.Int("totalRecords", *total))
		return data, total, nil
	}

	topups, totalRecords, err := s.topupQueryRepository.FindAllTopups(req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, "FindAll", "FAILED_FIND_ALL_TOPUPS", span, &status, topup_errors.ErrFailedFindAllTopups)
	}

	so := s.mapping.ToTopupResponses(topups)

	s.mencache.SetCachedTopupsCache(req, so, totalRecords)

	s.logger.Debug("Fetched topups from DB",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
	)

	return so, totalRecords, nil
}

func (s *topupQueryService) FindAllByCardNumber(req *requests.FindAllTopupsByCardNumber) ([]*response.TopupResponse, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindAllByCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindAllByCardNumber")
	defer span.End()

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search
	cardNumber := req.CardNumber

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
		attribute.String("card_number", cardNumber),
	)

	s.logger.Debug("Fetching topup by card number",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search),
		zap.String("card_number", cardNumber),
	)

	if data, total, found := s.mencache.GetCacheTopupByCardCache(req); found {
		span.SetAttributes(attribute.String("cache.status", "hit"))
		s.logger.Debug("Successfully fetched topup by card number from cache",
			zap.Int("totalRecords", *total))
		return data, total, nil
	}

	topups, totalRecords, err := s.topupQueryRepository.FindAllTopupByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, "FindAllByCardNumber", "FAILED_FIND_ALL_TOPUPS_BY_CARD", span, &status, topup_errors.ErrFailedFindAllTopupsByCardNumber)
	}

	so := s.mapping.ToTopupResponses(topups)

	s.mencache.SetCacheTopupByCardCache(req, so, totalRecords)

	s.logger.Debug("Successfully fetched topup",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
	)

	return so, totalRecords, nil
}

func (s *topupQueryService) FindById(topupID int) (*response.TopupResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindById", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindById")
	defer span.End()

	span.SetAttributes(attribute.Int("topup_id", topupID))
	s.logger.Debug("Fetching topup by ID", zap.Int("topup_id", topupID))

	if data := s.mencache.GetCachedTopupCache(topupID); data != nil {
		s.logger.Debug("Successfully fetched topup from cache", zap.Int("topup_id", topupID))
		return data, nil
	}

	span.SetAttributes(attribute.String("cache.status", "miss"))

	topup, err := s.topupQueryRepository.FindById(topupID)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "FindById", "FAILED_FIND_TOPUP", span, &status, topup_errors.ErrFailedFindTopupById)
	}

	so := s.mapping.ToTopupResponse(topup)

	s.mencache.SetCachedTopupCache(so)

	s.logger.Debug("Topup success", zap.Int("topup_id", topupID))
	return so, nil
}

func (s *topupQueryService) FindByActive(req *requests.FindAllTopups) ([]*response.TopupResponseDeleteAt, *int, *response.ErrorResponse) {
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

	s.logger.Debug("Fetching active topup",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if data, total, found := s.mencache.GetCachedTopupActiveCache(req); found {
		s.logger.Debug("Successfully fetched active topup from cache",
			zap.Int("totalRecords", *total))
		return data, total, nil
	}

	topups, totalRecords, err := s.topupQueryRepository.FindByActive(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, "FindByActive", "FAILED_FIND_ACTIVE_TOPUPS", span, &status, topup_errors.ErrFailedFindActiveTopups)
	}

	so := s.mapping.ToTopupResponsesDeleteAt(topups)

	s.mencache.SetCachedTopupActiveCache(req, so, totalRecords)

	s.logger.Debug("Successfully fetched active topup",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *topupQueryService) FindByTrashed(req *requests.FindAllTopups) ([]*response.TopupResponseDeleteAt, *int, *response.ErrorResponse) {
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

	s.logger.Debug("Fetching trashed topup",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if data, total, found := s.mencache.GetCachedTopupTrashedCache(req); found {
		s.logger.Debug("Successfully fetched trashed topup from cache",
			zap.Int("totalRecords", *total))
		return data, total, nil
	}

	topups, totalRecords, err := s.topupQueryRepository.FindByTrashed(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, "FindByTrashed", "FAILED_FIND_TRASHED_TOPUPS", span, &status, topup_errors.ErrFailedFindTrashedTopups)
	}

	so := s.mapping.ToTopupResponsesDeleteAt(topups)

	s.logger.Debug("Successfully fetched trashed topup",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *topupQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
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
