package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	traceunic "github.com/MamangRust/monolith-payment-gateway-pkg/trace_unic"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
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
	trace                trace.Tracer
	topupQueryRepository repository.TopupQueryRepository
	logger               logger.LoggerInterface
	mapping              responseservice.TopupResponseMapper
	requestCounter       *prometheus.CounterVec
	requestDuration      *prometheus.HistogramVec
}

func NewTopupQueryService(
	ctx context.Context, topupQueryRepository repository.TopupQueryRepository, logger logger.LoggerInterface, mapping responseservice.TopupResponseMapper) *topupQueryService {
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
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &topupQueryService{
		ctx:                  ctx,
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

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching topup",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	topups, totalRecords, err := s.topupQueryRepository.FindAllTopups(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_TOPUPS")

		s.logger.Error("Failed to fetch topup",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch topup")
		status = "failed_to_retrieve_topups"

		return nil, nil, topup_errors.ErrFailedFindAllTopups
	}

	so := s.mapping.ToTopupResponses(topups)

	s.logger.Debug("Successfully fetched topup",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

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

	page := req.Page
	pageSize := req.PageSize
	search := req.Search
	card_number := req.CardNumber

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
		attribute.String("card_number", card_number),
	)

	s.logger.Debug("Fetching topup by card number",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search),
		zap.String("card_number", card_number),
	)

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	topups, totalRecords, err := s.topupQueryRepository.FindAllTopupByCardNumber(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_TOPUPS_BY_CARD_NUMBER")

		s.logger.Error("Failed to fetch topup by card number",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.String("card_number", card_number))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch topup by card number")
		status = "failed_to_retrieve_topups_by_card_number"
		return nil, nil, topup_errors.ErrFailedFindAllTopupsByCardNumber
	}

	so := s.mapping.ToTopupResponses(topups)

	s.logger.Debug("Successfully fetched topup",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

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

	span.SetAttributes(
		attribute.Int("topup_id", topupID),
	)

	s.logger.Debug("Fetching topup by ID", zap.Int("topup_id", topupID))

	topup, err := s.topupQueryRepository.FindById(topupID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TOPUP_BY_ID")

		s.logger.Error("Failed to fetch topup by ID", zap.Int("topup_id", topupID), zap.String("trace_id", traceID), zap.Error(err))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch topup by ID")
		status = "failed_to_retrieve_topup_by_id"

		return nil, topup_errors.ErrTopupNotFoundRes
	}

	so := s.mapping.ToTopupResponse(topup)

	s.logger.Debug("Successfully fetched topup", zap.Int("topup_id", topupID))

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

	page := req.Page
	pageSize := req.PageSize
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

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	topups, totalRecords, err := s.topupQueryRepository.FindByActive(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_TOPUPS")

		s.logger.Error("Failed to fetch active topup",
			zap.Error(err),
			zap.String("trace_id", traceID),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch active topup")
		status = "failed_to_retrieve_topups"

		return nil, nil, topup_errors.ErrFailedFindAllTopups
	}

	so := s.mapping.ToTopupResponsesDeleteAt(topups)

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

	page := req.Page
	pageSize := req.PageSize
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

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	topups, totalRecords, err := s.topupQueryRepository.FindByTrashed(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_TOPUPS")

		s.logger.Error("Failed to fetch trashed topup",
			zap.Error(err),
			zap.String("trace_id", traceID),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch trashed topup")
		status = "failed_to_retrieve_topups"

		return nil, nil, topup_errors.ErrFailedFindAllTopups
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
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
