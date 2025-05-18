package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	traceunic "github.com/MamangRust/monolith-payment-gateway-pkg/trace_unic"
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
	trace                trace.Tracer
	saldoQueryRepository repository.SaldoQueryRepository
	logger               logger.LoggerInterface
	mapping              responseservice.SaldoResponseMapper
	requestCounter       *prometheus.CounterVec
	requestDuration      *prometheus.HistogramVec
}

func NewSaldoQueryService(ctx context.Context, saldo repository.SaldoQueryRepository, logger logger.LoggerInterface, mapping responseservice.SaldoResponseMapper) *saldoQueryService {
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
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &saldoQueryService{
		ctx:                  ctx,
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

	page := req.Page
	pageSize := req.PageSize
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

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Debug("Fetching all saldo records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	res, totalRecords, err := s.saldoQueryRepository.FindAllSaldos(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_SALDOS")

		s.logger.Error("Failed to retrieve saldo",
			zap.Error(err),
			zap.String("traceID", traceID),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		span.SetAttributes(attribute.String("traceID", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve saldo")
		status = "failed_to_retrieve_saldo"

		return nil, nil, saldo_errors.ErrFailedFindAllSaldos
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

	page := req.Page
	pageSize := req.PageSize
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

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	res, totalRecords, err := s.saldoQueryRepository.FindByActive(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ACTIVE_SALDOS")

		s.logger.Error("Failed to retrieve active saldo",
			zap.Error(err),
			zap.String("traceID", traceID),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		span.SetAttributes(attribute.String("traceID", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve active saldo")
		status = "failed_to_retrieve_active_saldo"

		return nil, nil, saldo_errors.ErrFailedFindActiveSaldos
	}

	so := s.mapping.ToSaldoResponsesDeleteAt(res)

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

	page := req.Page
	pageSize := req.PageSize
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

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	res, totalRecords, err := s.saldoQueryRepository.FindByTrashed(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TRASHED_SALDOS")

		s.logger.Error("Failed to retrieve trashed saldo",
			zap.Error(err),
			zap.String("traceID", traceID),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		span.SetAttributes(attribute.String("traceID", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve trashed saldo")
		status = "failed_to_retrieve_trashed_saldo"

		return nil, nil, saldo_errors.ErrFailedFindTrashedSaldos
	}

	s.logger.Debug("Successfully fetched trashed saldo",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize))

	so := s.mapping.ToSaldoResponsesDeleteAt(res)

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

	res, err := s.saldoQueryRepository.FindById(saldo_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_SALDO_BY_ID")

		s.logger.Error("Failed to retrieve saldo details",
			zap.Int("saldo_id", saldo_id),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve saldo details")
		status = "failed_to_retrieve_saldo_details"

		return nil, saldo_errors.ErrFailedSaldoNotFound
	}

	so := s.mapping.ToSaldoResponse(res)

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

	res, err := s.saldoQueryRepository.FindByCardNumber(card_number)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_SALDO_BY_CARD_NUMBER")

		s.logger.Error("Failed to retrieve saldo details",
			zap.String("card_number", card_number),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve saldo details")
		status = "failed_to_retrieve_saldo_details"

		return nil, saldo_errors.ErrFailedSaldoNotFound
	}

	so := s.mapping.ToSaldoResponse(res)

	s.logger.Debug("Successfully fetched saldo by card number", zap.String("card_number", card_number))

	return so, nil
}

func (s *saldoQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
