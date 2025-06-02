package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-card/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-card/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-card/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type cardQueryService struct {
	ctx                 context.Context
	errorhandler        errorhandler.CardQueryErrorHandler
	mencache            mencache.CardQueryCache
	trace               trace.Tracer
	cardQueryRepository repository.CardQueryRepository
	logger              logger.LoggerInterface
	mapping             responseservice.CardResponseMapper
	requestCounter      *prometheus.CounterVec
	requestDuration     *prometheus.HistogramVec
}

func NewCardQueryService(
	ctx context.Context,
	errorhandler errorhandler.CardQueryErrorHandler,
	mencache mencache.CardQueryCache,
	cardQueryRepository repository.CardQueryRepository,
	logger logger.LoggerInterface,
	mapper responseservice.CardResponseMapper,
) *cardQueryService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "card_query_request_count",
		Help: "Number of card query requests CardQueryService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "card_query_request_duration_seconds",
		Help:    "Duration of card query requests CardQueryService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	return &cardQueryService{
		ctx:                 ctx,
		errorhandler:        errorhandler,
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

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching card records",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if data, total, found := s.mencache.GetFindAllCache(req); found {
		s.logger.Debug("Successfully fetched card records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	cards, totalRecords, err := s.cardQueryRepository.FindAllCards(req)

	if err != nil {
		return s.errorhandler.HandleFindAllError(err, "FindAll", "FAILED_FIND_ALL_CARD", span, &status, zap.Error(err))
	}

	responseData := s.mapping.ToCardsResponse(cards)

	s.mencache.SetFindAllCache(req, responseData, totalRecords)

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

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching active card records",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if data, total, found := s.mencache.GetByActiveCache(req); found {
		s.logger.Debug("Successfully fetched card records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	res, totalRecords, err := s.cardQueryRepository.FindByActive(req)

	if err != nil {
		return s.errorhandler.HandleFindByActiveError(err, "FindByActive", "FAILED_FIND_ACTIVE_CARD", span, &status, zap.Error(err))
	}

	responseData := s.mapping.ToCardsResponseDeleteAt(res)

	s.mencache.SetByActiveCache(req, responseData, totalRecords)

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

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching trashed card records",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if data, total, found := s.mencache.GetByTrashedCache(req); found {
		s.logger.Debug("Successfully fetched card records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	res, totalRecords, err := s.cardQueryRepository.FindByTrashed(req)
	if err != nil {
		return s.errorhandler.HandleFindByTrashedError(err, "FindByTrashed", "FAILED_FIND_TRASHED_CARD", span, &status, zap.Error(err))
	}

	responseData := s.mapping.ToCardsResponseDeleteAt(res)

	s.mencache.SetByTrashedCache(req, responseData, totalRecords)

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

	if data, found := s.mencache.GetByIdCache(card_id); found {
		s.logger.Debug("Successfully fetched card from cache", zap.Int("card_id", card_id))
		return data, nil
	}

	res, err := s.cardQueryRepository.FindById(card_id)

	if err != nil {
		return s.errorhandler.HandleFindByIdError(err, "FindById", "FAILED_TO_FIND_CARD", span, &status, zap.Error(err))
	}

	so := s.mapping.ToCardResponse(res)

	s.mencache.SetByIdCache(card_id, so)

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

	if data, found := s.mencache.GetByUserIDCache(user_id); found {
		s.logger.Debug("Successfully fetched card from cache by user ID", zap.Int("user_id", user_id))
		return data, nil
	}

	res, err := s.cardQueryRepository.FindCardByUserId(user_id)

	if err != nil {
		return s.errorhandler.HandleFindByUserIdError(err, "FindByUserID", "FAILED_FIND_CARD_BY_USER_ID", span, &status, zap.Error(err))
	}

	so := s.mapping.ToCardResponse(res)

	s.mencache.SetByUserIDCache(user_id, so)

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

	if data, found := s.mencache.GetByCardNumberCache(card_number); found {
		s.logger.Debug("Successfully fetched card record from cache by card number", zap.String("card_number", card_number))
		return data, nil
	}

	res, err := s.cardQueryRepository.FindCardByCardNumber(card_number)

	if err != nil {
		return s.errorhandler.HandleFindByCardNumberError(err, "FindByCardNumber", "FAILED_FIND_CARD_BY_CARD_NUMBER", span, &status, zap.Error(err))
	}

	so := s.mapping.ToCardResponse(res)

	s.mencache.SetByCardNumberCache(card_number, so)

	s.logger.Debug("Successfully fetched card record by card number", zap.String("card_number", card_number))

	return so, nil
}

func (s *cardQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func (s *cardQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
