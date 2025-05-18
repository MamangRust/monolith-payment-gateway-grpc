package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-card/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	traceunic "github.com/MamangRust/monolith-payment-gateway-pkg/trace_unic"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type cardQueryService struct {
	ctx                 context.Context
	trace               trace.Tracer
	cardQueryRepository repository.CardQueryRepository
	logger              logger.LoggerInterface
	mapping             responseservice.CardResponseMapper
	requestCounter      *prometheus.CounterVec
	requestDuration     *prometheus.HistogramVec
}

func NewCardQueryService(
	ctx context.Context,
	cardQueryRepository repository.CardQueryRepository,
	logger logger.LoggerInterface,
	mapper responseservice.CardResponseMapper,
) *cardQueryService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "card_query_request_count",
		Help: "Number of card query requests CardQueryService",
	}, []string{"status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "card_query_request_duration_seconds",
		Help:    "Duration of card query requests CardQueryService",
		Buckets: prometheus.DefBuckets,
	}, []string{"status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	return &cardQueryService{
		ctx:                 ctx,
		trace:               otel.Tracer("card-query-service"),
		cardQueryRepository: cardQueryRepository,
		logger:              logger,
		mapping:             mapper,
		requestCounter:      requestCounter,
		requestDuration:     requestDuration,
	}
}

func (s *cardQueryService) FindAll(req *requests.FindAllCards) ([]*response.CardResponse, *int, *response.ErrorResponse) {
	startTime := time.Now()

	status := "success"

	defer func() {
		s.recordMetrics("FindAll", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindAll")
	defer span.End()

	span.SetAttributes(
		attribute.Int("page", req.Page),
		attribute.Int("pageSize", req.PageSize),
		attribute.String("search", req.Search),
	)

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching card records",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	cards, totalRecords, err := s.cardQueryRepository.FindAllCards(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_CARDS")

		s.logger.Error("Failed to fetch card records",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.String("traceID", traceID))

		span.SetAttributes(attribute.String("trace.id", traceID))

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch card records")
		status = "failed_find_all_cards"

		return nil, nil, card_errors.ErrFailedFindAllCards
	}

	responseData := s.mapping.ToCardsResponse(cards)

	s.logger.Debug("Successfully fetched card records",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return responseData, totalRecords, nil
}

func (s *cardQueryService) FindByActive(req *requests.FindAllCards) ([]*response.CardResponseDeleteAt, *int, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByActive", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindByActive")
	defer span.End()

	span.SetAttributes(
		attribute.Int("page", req.Page),
		attribute.Int("pageSize", req.PageSize),
		attribute.String("search", req.Search),
	)

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching active card records",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	res, totalRecords, err := s.cardQueryRepository.FindByActive(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ACTIVE_CARDS")

		s.logger.Error("Failed to fetch active card records",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.String("traceID", traceID))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch active card records")
		status = "failed_find_active_cards"

		return nil, nil, card_errors.ErrFailedFindActiveCards
	}

	responseData := s.mapping.ToCardsResponseDeleteAt(res)

	s.logger.Debug("Successfully fetched active card records",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return responseData, totalRecords, nil
}

func (s *cardQueryService) FindByTrashed(req *requests.FindAllCards) ([]*response.CardResponseDeleteAt, *int, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByTrashed", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindByTrashed")
	defer span.End()

	span.SetAttributes(
		attribute.Int("page", req.Page),
		attribute.Int("pageSize", req.PageSize),
		attribute.String("search", req.Search),
	)

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching trashed card records",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	res, totalRecords, err := s.cardQueryRepository.FindByTrashed(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TRASHED_CARDS")

		s.logger.Error("Failed to fetch trashed card records",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.String("traceID", traceID))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch trashed card records")
		status = "failed_find_trashed_cards"

		return nil, nil, card_errors.ErrFailedFindTrashedCards
	}

	responseData := s.mapping.ToCardsResponseDeleteAt(res)

	s.logger.Debug("Successfully fetched trashed card records",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return responseData, totalRecords, nil
}

func (s *cardQueryService) FindById(card_id int) (*response.CardResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("FindById", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindById")
	defer span.End()

	s.logger.Debug("Fetching card by ID", zap.Int("card_id", card_id))

	res, err := s.cardQueryRepository.FindById(card_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_BY_ID")

		s.logger.Error("Failed to retrieve Card details by ID",
			zap.Error(err),
			zap.Int("card_id", card_id),
			zap.String("traceID", traceID))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve Card details by ID")
		status = "failed_find_by_id"

		return nil, card_errors.ErrFailedFindById
	}

	so := s.mapping.ToCardResponse(res)

	s.logger.Debug("Successfully fetched card", zap.Int("card_id", card_id))

	return so, nil
}

func (s *cardQueryService) FindByUserID(user_id int) (*response.CardResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("FindByUserID", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindByUserID")
	defer span.End()

	s.logger.Debug("Fetching card by user ID", zap.Int("user_id", user_id))

	res, err := s.cardQueryRepository.FindCardByUserId(user_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_BY_USER_ID")

		s.logger.Error("Failed to retrieve Card details by user ID",
			zap.Error(err),
			zap.Int("user_id", user_id),
			zap.String("traceID", traceID))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve Card details by user ID")
		status = "failed_find_by_user_id"

		return nil, card_errors.ErrFailedFindByUserID
	}

	so := s.mapping.ToCardResponse(res)

	s.logger.Debug("Successfully fetched card records by user ID", zap.Int("user_id", user_id))

	return so, nil
}

func (s *cardQueryService) FindByCardNumber(card_number string) (*response.CardResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("FindByCardNumber", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindByCardNumber")
	defer span.End()

	s.logger.Debug("Fetching card record by card number", zap.String("card_number", card_number))

	res, err := s.cardQueryRepository.FindCardByCardNumber(card_number)

	if err != nil {
		traceID := traceunic.GenerateTraceID("CARD_NOT_FOUND")

		s.logger.Error("Failed to retrieve Card details by card number",
			zap.Error(err),
			zap.String("card_number", card_number),
			zap.String("traceID", traceID))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve Card details by card number")
		status = "card_not_found"

		return nil, card_errors.ErrCardNotFoundRes
	}

	so := s.mapping.ToCardResponse(res)

	s.logger.Debug("Successfully fetched card record by card number", zap.String("card_number", card_number))

	return so, nil
}

func (s *cardQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
