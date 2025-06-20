package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-role/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-role/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-role/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/role_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type roleQueryService struct {
	ctx             context.Context
	errorhandler    errorhandler.RoleQueryErrorHandler
	mencache        mencache.RoleQueryCache
	trace           trace.Tracer
	roleQuery       repository.RoleQueryRepository
	logger          logger.LoggerInterface
	mapping         responseservice.RoleResponseMapper
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewRoleQueryService(ctx context.Context, errorhandler errorhandler.RoleQueryErrorHandler,
	mencache mencache.RoleQueryCache, roleQuery repository.RoleQueryRepository, logger logger.LoggerInterface, mapping responseservice.RoleResponseMapper) *roleQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "role_query_service_request_total",
			Help: "Total number of requests to the RoleQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "role_query_service_request_duration_seconds",
			Help:    "Histogram of request durations for the RoleQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &roleQueryService{
		ctx:             ctx,
		errorhandler:    errorhandler,
		mencache:        mencache,
		trace:           otel.Tracer("role-query-service"),
		roleQuery:       roleQuery,
		logger:          logger,
		mapping:         mapping,
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}
}

func (s *roleQueryService) FindAll(req *requests.FindAllRoles) ([]*response.RoleResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindAll"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedRoles(req); found {
		logSuccess("Successfully retrieved all role records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	res, totalRecords, err := s.roleQuery.FindAllRoles(req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(
			err, method, "FAILED_FIND_ALL_ROLE", span, &status, zap.Error(err))
	}
	so := s.mapping.ToRolesResponse(res)

	s.mencache.SetCachedRoles(req, so, totalRecords)

	logSuccess("Successfully retrieved all role records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *roleQueryService) FindById(id int) (*response.RoleResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindById", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindById")
	defer span.End()

	span.SetAttributes(
		attribute.Int("id", id),
	)

	s.logger.Debug("Fetching role by ID", zap.Int("id", id))

	if data, found := s.mencache.GetCachedRoleById(id); found {
		s.logger.Debug("Successfully fetched role from cache")
		return data, nil
	}

	res, err := s.roleQuery.FindById(id)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(
			err, "FindById", "FAILED_FIND_ROLE_BY_ID", span, &status, role_errors.ErrRoleNotFoundRes, zap.Error(err))
	}

	so := s.mapping.ToRoleResponse(res)

	s.mencache.SetCachedRoleById(id, so)

	s.logger.Debug("Successfully fetched role", zap.Int("id", id))

	return so, nil
}

func (s *roleQueryService) FindByUserId(id int) ([]*response.RoleResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByUserId", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindByUserId")
	defer span.End()

	span.SetAttributes(
		attribute.Int("id", id),
	)

	s.logger.Debug("Fetching role by user ID", zap.Int("id", id))

	if data, found := s.mencache.GetCachedRoleByUserId(id); found {
		s.logger.Debug("Successfully fetched role from cache")
		return data, nil
	}

	res, err := s.roleQuery.FindByUserId(id)

	if err != nil {
		return s.errorhandler.HandleRepositoryListError(err, "FindByUserId", "FAILED_FIND_ROLE_BY_USER_ID", span, &status, zap.Error(err))
	}

	so := s.mapping.ToRolesResponse(res)

	s.mencache.SetCachedRoleByUserId(id, so)

	s.logger.Debug("Successfully fetched role by user ID", zap.Int("id", id))

	return so, nil
}

func (s *roleQueryService) FindByActiveRole(req *requests.FindAllRoles) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindByActiveRole"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedRoleActive(req); found {
		logSuccess("Successfully retrieved active role records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	res, totalRecords, err := s.roleQuery.FindByActiveRole(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeletedError(err, method, "FAILED_FIND_BY_ACTIVE_ROLE", span, &status, role_errors.ErrRoleNotFoundRes, zap.Error(err))
	}

	so := s.mapping.ToRolesResponseDeleteAt(res)

	s.mencache.SetCachedRoleActive(req, so, totalRecords)

	logSuccess("Successfully retrieved active role records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *roleQueryService) FindByTrashedRole(req *requests.FindAllRoles) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindByTrashedRole"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedRoleTrashed(req); found {
		logSuccess("Successfully fetched trashed role from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	res, totalRecords, err := s.roleQuery.FindByTrashedRole(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeletedError(err, method, "FAILED_FIND_BY_TRASHED_ROLE", span, &status, role_errors.ErrRoleNotFoundRes)
	}
	so := s.mapping.ToRolesResponseDeleteAt(res)

	s.mencache.SetCachedRoleTrashed(req, so, totalRecords)

	logSuccess("Successfully fetched trashed role", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *roleQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func (s *roleQueryService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *roleQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
