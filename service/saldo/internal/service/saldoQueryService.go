package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-saldo/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type saldoQueryService struct {
	ctx                  context.Context
	errorhandler         errorhandler.SaldoQueryErrorHandler
	mencache             mencache.SaldoQueryCache
	trace                trace.Tracer
	saldoQueryRepository repository.SaldoQueryRepository
	logger               logger.LoggerInterface
	mapping              responseservice.SaldoResponseMapper
	requestCounter       *prometheus.CounterVec
	requestDuration      *prometheus.HistogramVec
}

func NewSaldoQueryService(ctx context.Context, errorhandler errorhandler.SaldoQueryErrorHandler,
	mencache mencache.SaldoQueryCache, saldo repository.SaldoQueryRepository, logger logger.LoggerInterface, mapping responseservice.SaldoResponseMapper) *saldoQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "saldo_query_service_request_total",
			Help: "Total number of requests to the SaldoQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "saldo_query_service_request_duration_seconds",
			Help:    "Histogram of request durations for the SaldoQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &saldoQueryService{
		ctx:                  ctx,
		errorhandler:         errorhandler,
		mencache:             mencache,
		trace:                otel.Tracer("saldo-query-service"),
		saldoQueryRepository: saldo,
		logger:               logger,
		mapping:              mapping,
		requestCounter:       requestCounter,
		requestDuration:      requestDuration,
	}
}

func (s *saldoQueryService) FindAll(req *requests.FindAllSaldos) ([]*response.SaldoResponse, *int, *response.ErrorResponse) {
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

	s.logger.Debug("Fetching saldo",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	s.logger.Debug("Fetching all saldo records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	if data, total, found := s.mencache.GetCachedSaldos(req); found {
		s.logger.Debug("Successfully fetched saldo from cache",
			zap.Int("totalRecords", *total))
		return data, total, nil
	}

	res, totalRecords, err := s.saldoQueryRepository.FindAllSaldos(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, "FindAll", "FAILED_FIND_SALDOS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToSaldoResponses(res)

	s.logger.Debug("Successfully fetched saldo",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize))

	return so, totalRecords, nil
}

func (s *saldoQueryService) FindByActive(req *requests.FindAllSaldos) ([]*response.SaldoResponseDeleteAt, *int, *response.ErrorResponse) {
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

	s.logger.Debug("Fetching active saldo",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if data, total, found := s.mencache.GetCachedSaldoByActive(req); found {
		s.logger.Debug("Successfully fetched active saldo from cache",
			zap.Int("totalRecords", *total))
		return data, total, nil
	}

	res, totalRecords, err := s.saldoQueryRepository.FindByActive(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, "FindByActive", "FAILED_FIND_ACTIVE_SALDOS", span, &status, saldo_errors.ErrFailedFindActiveSaldos, zap.Error(err))
	}

	so := s.mapping.ToSaldoResponsesDeleteAt(res)

	s.mencache.SetCachedSaldoByActive(req, so, totalRecords)

	s.logger.Debug("Successfully fetched active saldo",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *saldoQueryService) FindByTrashed(req *requests.FindAllSaldos) ([]*response.SaldoResponseDeleteAt, *int, *response.ErrorResponse) {
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

	s.logger.Debug("Fetching saldo record",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if data, total, found := s.mencache.GetCachedSaldoByTrashed(req); found {
		s.logger.Debug("Successfully fetched trashed saldo from cache",
			zap.Int("totalRecords", *total))
		return data, total, nil
	}

	res, totalRecords, err := s.saldoQueryRepository.FindByTrashed(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, "FindByTrashed", "FAILED_FIND_TRASHED_SALDOS", span, &status, saldo_errors.ErrFailedFindTrashedSaldos, zap.Error(err))
	}
	so := s.mapping.ToSaldoResponsesDeleteAt(res)

	s.mencache.SetCachedSaldoByTrashed(req, so, totalRecords)

	s.logger.Debug("Successfully fetched trashed saldo",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize))

	return so, totalRecords, nil
}

func (s *saldoQueryService) FindById(saldo_id int) (*response.SaldoResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindById", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindById")
	defer span.End()

	span.SetAttributes(
		attribute.Int("saldo_id", saldo_id),
	)

	s.logger.Debug("Fetching saldo record by ID", zap.Int("saldo_id", saldo_id))

	if data, found := s.mencache.GetCachedSaldoById(saldo_id); found {
		s.logger.Debug("Successfully fetched saldo from cache", zap.Int("saldo_id", saldo_id))
		return data, nil
	}

	res, err := s.saldoQueryRepository.FindById(saldo_id)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "FindById", "FAILED_FIND_SALDO", span, &status, saldo_errors.ErrFailedSaldoNotFound, zap.Error(err))
	}

	so := s.mapping.ToSaldoResponse(res)

	s.mencache.SetCachedSaldoById(saldo_id, so)

	s.logger.Debug("Successfully fetched saldo", zap.Int("saldo_id", saldo_id))

	return so, nil
}

func (s *saldoQueryService) FindByCardNumber(card_number string) (*response.SaldoResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByCardNumber", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindByCardNumber")
	defer span.End()

	span.SetAttributes(
		attribute.String("card_number", card_number),
	)

	s.logger.Debug("Fetching saldo record by card number", zap.String("card_number", card_number))

	if data, found := s.mencache.GetCachedSaldoByCardNumber(card_number); found {
		s.logger.Debug("Successfully fetched saldo from cache by card number", zap.String("card_number", card_number))
		return data, nil
	}

	res, err := s.saldoQueryRepository.FindByCardNumber(card_number)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "FindByCardNumber", "FAILED_FIND_SALDO_BY_CARD_NUMBER", span, &status, saldo_errors.ErrFailedSaldoNotFound, zap.Error(err))
	}

	so := s.mapping.ToSaldoResponse(res)

	s.mencache.SetCachedSaldoByCardNumber(card_number, so)

	s.logger.Debug("Successfully fetched saldo by card number", zap.String("card_number", card_number))

	return so, nil
}

func (s *saldoQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func (s *saldoQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
