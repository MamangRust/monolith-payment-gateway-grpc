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
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type roleCommandService struct {
	ctx             context.Context
	errorhandler    errorhandler.RoleCommandErrorHandler
	mencache        mencache.RoleCommandCache
	trace           trace.Tracer
	roleCommand     repository.RoleCommandRepository
	logger          logger.LoggerInterface
	mapping         responseservice.RoleResponseMapper
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewRoleCommandService(ctx context.Context, errorhandler errorhandler.RoleCommandErrorHandler,
	mencache mencache.RoleCommandCache, roleCommand repository.RoleCommandRepository, logger logger.LoggerInterface, mapping responseservice.RoleResponseMapper) *roleCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "role_command_service_request_total",
			Help: "Total number of requests to the RoleCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "role_command_service_request_duration_seconds",
			Help:    "Histogram of request durations for the RoleCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &roleCommandService{
		ctx:             ctx,
		errorhandler:    errorhandler,
		mencache:        mencache,
		trace:           otel.Tracer("role-command-service"),
		roleCommand:     roleCommand,
		logger:          logger,
		mapping:         mapping,
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}
}

func (s *roleCommandService) CreateRole(request *requests.CreateRoleRequest) (*response.RoleResponse, *response.ErrorResponse) {
	const method = "CreateRole"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	role, err := s.roleCommand.CreateRole(request)

	if err != nil {
		return s.errorhandler.HandleCreateRoleError(err, method, "FAILED_CREATE_ROLE", span, &status, zap.Error(err))
	}

	so := s.mapping.ToRoleResponse(role)

	logSuccess("Successfully created role", zap.Int("role.id", role.ID))

	return so, nil
}

func (s *roleCommandService) UpdateRole(request *requests.UpdateRoleRequest) (*response.RoleResponse, *response.ErrorResponse) {
	const method = "UpdateRole"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	role, err := s.roleCommand.UpdateRole(request)
	if err != nil {
		return s.errorhandler.HandleUpdateRoleError(err, method, "FAILED_UPDATE_ROLE", span, &status, zap.Error(err))
	}

	so := s.mapping.ToRoleResponse(role)

	s.mencache.DeleteCachedRole(*request.ID)

	logSuccess("Successfully updated role", zap.Int("role.id", role.ID))

	return so, nil
}

func (s *roleCommandService) TrashedRole(id int) (*response.RoleResponse, *response.ErrorResponse) {
	const method = "TrashedRole"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	role, err := s.roleCommand.TrashedRole(id)

	if err != nil {
		return s.errorhandler.HandleTrashedRoleError(err, method, "FAILED_TRASH_ROLE", span, &status, zap.Error(err))
	}

	so := s.mapping.ToRoleResponse(role)

	s.mencache.DeleteCachedRole(id)

	logSuccess("Successfully trashed role", zap.Int("role.id", role.ID))

	return so, nil
}

func (s *roleCommandService) RestoreRole(id int) (*response.RoleResponse, *response.ErrorResponse) {
	const method = "RestoreRole"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	role, err := s.roleCommand.RestoreRole(id)

	if err != nil {
		return s.errorhandler.HandleRestoreRoleError(err, method, "FAILED_RESTORE_ROLE", span, &status, zap.Error(err))
	}

	so := s.mapping.ToRoleResponse(role)

	logSuccess("Successfully restored role", zap.Int("role.id", role.ID))

	return so, nil
}

func (s *roleCommandService) DeleteRolePermanent(id int) (bool, *response.ErrorResponse) {
	const method = "DeleteRolePermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	_, err := s.roleCommand.DeleteRolePermanent(id)
	if err != nil {
		return s.errorhandler.HandleDeleteRolePermanentError(err, method, "FAILED_DELETE_ROLE_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully permanently deleted role", zap.Int("role.id", id))

	return true, nil
}

func (s *roleCommandService) RestoreAllRole() (bool, *response.ErrorResponse) {
	const method = "RestoreAllRole"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	_, err := s.roleCommand.RestoreAllRole()
	if err != nil {
		return s.errorhandler.HandleRestoreAllRoleError(err, method, "FAILED_RESTORE_ALL_ROLE", span, &status, zap.Error(err))
	}

	logSuccess("Successfully restored all roles", zap.Bool("success", true))

	return true, nil
}

func (s *roleCommandService) DeleteAllRolePermanent() (bool, *response.ErrorResponse) {
	const method = "DeleteAllRolePermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	_, err := s.roleCommand.DeleteAllRolePermanent()

	if err != nil {
		return s.errorhandler.HandleDeleteAllRolePermanentError(err, method, "FAILED_DELETE_ALL_ROLE_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully permanently deleted all roles", zap.Bool("success", true))

	return true, nil
}

func (s *roleCommandService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *roleCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
