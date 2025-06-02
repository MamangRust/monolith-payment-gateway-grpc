package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-user/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-user/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-user/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type userQueryService struct {
	ctx                 context.Context
	errorhandler        errorhandler.UserQueryError
	mencache            mencache.UserQueryCache
	trace               trace.Tracer
	userQueryRepository repository.UserQueryRepository
	logger              logger.LoggerInterface
	mapping             responseservice.UserResponseMapper
	requestCounter      *prometheus.CounterVec
	requestDuration     *prometheus.HistogramVec
}

func NewUserQueryService(ctx context.Context, errorhandler errorhandler.UserQueryError,
	mencache mencache.UserQueryCache, userQueryRepository repository.UserQueryRepository, logger logger.LoggerInterface, mapping responseservice.UserResponseMapper) *userQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_query_service_request_total",
			Help: "Total number of requests to the UserQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "user_query_service_request_duration_seconds",
			Help:    "Histogram of request durations for the UserQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &userQueryService{
		ctx:                 ctx,
		errorhandler:        errorhandler,
		mencache:            mencache,
		trace:               otel.Tracer("user-query-service"),
		userQueryRepository: userQueryRepository,
		logger:              logger,
		mapping:             mapping,
		requestCounter:      requestCounter,
		requestDuration:     requestDuration,
	}
}

func (s *userQueryService) FindAll(req *requests.FindAllUsers) ([]*response.UserResponse, *int, *response.ErrorResponse) {
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

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	if data, total, found := s.mencache.GetCachedUsersCache(req); found {
		return data, total, nil
	}

	users, totalRecords, err := s.userQueryRepository.FindAllUsers(req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, "FindAll", "FAILED_FIND_ALL_USERS", span, &status, zap.Error(err))
	}

	userResponses := s.mapping.ToUsersResponse(users)

	s.mencache.SetCachedUsersCache(req, userResponses, totalRecords)

	s.logger.Debug("Successfully fetched users", zap.Int("total_records", *totalRecords))
	return userResponses, totalRecords, nil
}

func (s *userQueryService) FindByID(id int) (*response.UserResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByID", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindByID")
	defer span.End()
	span.SetAttributes(attribute.Int("user_id", id))

	if data := s.mencache.GetCachedUserCache(id); data != nil {
		return data, nil
	}

	user, err := s.userQueryRepository.FindById(id)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "FindByID", "FAILED_FIND_USER", span, &status, user_errors.ErrUserNotFoundRes, zap.Int("user_id", id))
	}

	userRes := s.mapping.ToUserResponse(user)

	s.mencache.SetCachedUserCache(userRes)

	s.logger.Debug("Successfully fetched user", zap.Int("user_id", id))
	return userRes, nil
}

func (s *userQueryService) FindByActive(req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse) {
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

	s.logger.Debug("Fetching active user",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	if data, total, found := s.mencache.GetCachedUserActiveCache(req); found {
		return data, total, nil
	}

	users, totalRecords, err := s.userQueryRepository.FindByActive(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, "FindByActive", "FAILED_FIND_ACTIVE_USERS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToUsersResponseDeleteAt(users)

	s.mencache.SetCachedUserActiveCache(req, so, totalRecords)

	s.logger.Debug("Successfully fetched active user",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *userQueryService) FindByTrashed(req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse) {
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

	s.logger.Debug("Fetching trashed user",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	if data, total, found := s.mencache.GetCachedUserTrashedCache(req); found {
		return data, total, nil
	}

	users, totalRecords, err := s.userQueryRepository.FindByTrashed(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, "FindByTrashed", "FAILED_FIND_TRASHED_USERS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToUsersResponseDeleteAt(users)

	s.logger.Debug("Successfully fetched trashed user",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *userQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
