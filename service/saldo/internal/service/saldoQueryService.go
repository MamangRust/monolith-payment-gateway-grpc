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
	"go.opentelemetry.io/otel/codes"
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
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindAll"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedSaldos(req); found {
		logSuccess("Successfully retrieved all saldo records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	res, totalRecords, err := s.saldoQueryRepository.FindAllSaldos(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_SALDOS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToSaldoResponses(res)

	logSuccess("Successfully retrieved all saldo records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *saldoQueryService) FindByActive(req *requests.FindAllSaldos) ([]*response.SaldoResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindByActive"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedSaldoByActive(req); found {
		logSuccess("Successfully fetched active saldo from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	res, totalRecords, err := s.saldoQueryRepository.FindByActive(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_ACTIVE_SALDOS", span, &status, saldo_errors.ErrFailedFindActiveSaldos, zap.Error(err))
	}

	so := s.mapping.ToSaldoResponsesDeleteAt(res)

	s.mencache.SetCachedSaldoByActive(req, so, totalRecords)

	logSuccess("Successfully fetched active saldo", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *saldoQueryService) FindByTrashed(req *requests.FindAllSaldos) ([]*response.SaldoResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindByTrashed"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedSaldoByTrashed(req); found {
		logSuccess("Successfully fetched trashed saldo from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	res, totalRecords, err := s.saldoQueryRepository.FindByTrashed(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_TRASHED_SALDOS", span, &status, saldo_errors.ErrFailedFindTrashedSaldos, zap.Error(err))
	}
	so := s.mapping.ToSaldoResponsesDeleteAt(res)

	s.mencache.SetCachedSaldoByTrashed(req, so, totalRecords)

	logSuccess("Successfully fetched trashed saldo", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *saldoQueryService) FindById(saldo_id int) (*response.SaldoResponse, *response.ErrorResponse) {
	const method = "FindById"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("saldo.id", saldo_id))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedSaldoById(saldo_id); found {
		logSuccess("Successfully fetched saldo from cache", zap.Int("saldo.id", saldo_id))
		return data, nil
	}

	res, err := s.saldoQueryRepository.FindById(saldo_id)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_SALDO", span, &status, saldo_errors.ErrFailedSaldoNotFound, zap.Error(err))
	}

	so := s.mapping.ToSaldoResponse(res)

	s.mencache.SetCachedSaldoById(saldo_id, so)

	logSuccess("Successfully fetched saldo", zap.Int("saldo.id", saldo_id))

	return so, nil
}

func (s *saldoQueryService) FindByCardNumber(card_number string) (*response.SaldoResponse, *response.ErrorResponse) {
	const method = "FindByCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.String("card_number", card_number))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedSaldoByCardNumber(card_number); found {
		logSuccess("Successfully fetched saldo by card number from cache", zap.String("card_number", card_number))
		return data, nil
	}

	res, err := s.saldoQueryRepository.FindByCardNumber(card_number)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_SALDO_BY_CARD_NUMBER", span, &status, saldo_errors.ErrFailedSaldoNotFound, zap.Error(err))
	}

	so := s.mapping.ToSaldoResponse(res)

	s.mencache.SetCachedSaldoByCardNumber(card_number, so)

	logSuccess("Successfully fetched saldo by card number", zap.String("card_number", card_number))

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

func (s *saldoQueryService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *saldoQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
