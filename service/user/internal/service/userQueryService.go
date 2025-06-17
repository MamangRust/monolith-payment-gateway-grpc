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
	"go.opentelemetry.io/otel/codes"
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
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindAll"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()
	if data, total, found := s.mencache.GetCachedUsersCache(req); found {
		logSuccess("Successfully retrieved all user records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))

		return data, total, nil
	}

	users, totalRecords, err := s.userQueryRepository.FindAllUsers(req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_USERS", span, &status, zap.Error(err))
	}

	userResponses := s.mapping.ToUsersResponse(users)

	s.mencache.SetCachedUsersCache(req, userResponses, totalRecords)

	logSuccess("Successfully retrieved all user records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return userResponses, totalRecords, nil
}

func (s *userQueryService) FindByID(id int) (*response.UserResponse, *response.ErrorResponse) {
	const method = "FindByID"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("user.id", id))

	defer func() {
		end(status)
	}()

	if data := s.mencache.GetCachedUserCache(id); data != nil {
		return data, nil
	}

	user, err := s.userQueryRepository.FindById(id)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_USER", span, &status, user_errors.ErrUserNotFoundRes, zap.Error(err))
	}

	userRes := s.mapping.ToUserResponse(user)

	s.mencache.SetCachedUserCache(userRes)

	logSuccess("Successfully retrieved user record", zap.Int("user.id", id))

	return userRes, nil
}

func (s *userQueryService) FindByActive(req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindByActive"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedUserActiveCache(req); found {
		logSuccess("Successfully retrieved active user records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	users, totalRecords, err := s.userQueryRepository.FindByActive(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_ACTIVE_USERS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToUsersResponseDeleteAt(users)

	s.mencache.SetCachedUserActiveCache(req, so, totalRecords)

	logSuccess("Successfully retrieved active user records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *userQueryService) FindByTrashed(req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindByTrashed"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedUserTrashedCache(req); found {
		return data, total, nil
	}

	users, totalRecords, err := s.userQueryRepository.FindByTrashed(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_TRASHED_USERS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToUsersResponseDeleteAt(users)

	logSuccess("Successfully retrieved trashed user records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *userQueryService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *userQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func (s *userQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
