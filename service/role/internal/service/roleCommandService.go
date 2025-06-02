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
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("CreateRole", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "CreateRole")
	defer span.End()

	span.SetAttributes(
		attribute.String("name", request.Name),
	)

	s.logger.Debug("Starting CreateRole process",
		zap.String("roleName", request.Name),
	)

	role, err := s.roleCommand.CreateRole(request)

	if err != nil {
		return s.errorhandler.HandleCreateRoleError(err, "CreateRole", "FAILED_CREATE_ROLE", span, &status, zap.Error(err))
	}

	so := s.mapping.ToRoleResponse(role)

	s.logger.Debug("CreateRole process completed",
		zap.String("roleName", request.Name),
		zap.Int("roleID", role.ID),
	)

	return so, nil
}

func (s *roleCommandService) UpdateRole(request *requests.UpdateRoleRequest) (*response.RoleResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("UpdateRole", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateRole")
	defer span.End()

	span.SetAttributes(
		attribute.Int("id", *request.ID),
		attribute.String("name", request.Name),
	)

	s.logger.Debug("Starting UpdateRole process",
		zap.Int("roleID", *request.ID),
		zap.String("newRoleName", request.Name),
	)

	role, err := s.roleCommand.UpdateRole(request)
	if err != nil {
		return s.errorhandler.HandleUpdateRoleError(err, "UpdateRole", "FAILED_UPDATE_ROLE", span, &status, zap.Error(err))
	}

	so := s.mapping.ToRoleResponse(role)

	s.mencache.DeleteCachedRole(*request.ID)

	s.logger.Debug("UpdateRole process completed",
		zap.Int("roleID", *request.ID),
		zap.String("newRoleName", request.Name),
	)

	return so, nil
}

func (s *roleCommandService) TrashedRole(id int) (*response.RoleResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("TrashedRole", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "TrashedRole")
	defer span.End()

	span.SetAttributes(
		attribute.Int("id", id),
	)

	s.logger.Debug("Starting TrashedRole process",
		zap.Int("roleID", id),
	)

	role, err := s.roleCommand.TrashedRole(id)

	if err != nil {
		return s.errorhandler.HandleTrashedRoleError(err, "TrashedRole", "FAILED_TRASH_ROLE", span, &status, zap.Error(err))
	}

	so := s.mapping.ToRoleResponse(role)

	s.mencache.DeleteCachedRole(id)

	s.logger.Debug("TrashedRole process completed",
		zap.Int("roleID", id),
	)

	return so, nil
}

func (s *roleCommandService) RestoreRole(id int) (*response.RoleResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreRole", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreRole")
	defer span.End()

	span.SetAttributes(
		attribute.Int("id", id),
	)

	s.logger.Debug("Starting RestoreRole process",
		zap.Int("roleID", id),
	)

	role, err := s.roleCommand.RestoreRole(id)

	if err != nil {
		return s.errorhandler.HandleRestoreRoleError(err, "RestoreRole", "FAILED_RESTORE_ROLE", span, &status, zap.Error(err))
	}

	so := s.mapping.ToRoleResponse(role)

	s.logger.Debug("RestoreRole process completed",
		zap.Int("roleID", id),
	)

	return so, nil
}

func (s *roleCommandService) DeleteRolePermanent(id int) (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteRolePermanent", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteRolePermanent")
	defer span.End()

	span.SetAttributes(
		attribute.Int("id", id),
	)

	s.logger.Debug("Starting DeleteRolePermanent process",
		zap.Int("roleID", id),
	)

	_, err := s.roleCommand.DeleteRolePermanent(id)
	if err != nil {
		return s.errorhandler.HandleDeleteRolePermanentError(err, "DeleteRolePermanent", "FAILED_DELETE_ROLE_PERMANENT", span, &status, zap.Error(err))
	}

	s.logger.Debug("DeleteRolePermanent process completed",
		zap.Int("roleID", id),
	)

	return true, nil
}

func (s *roleCommandService) RestoreAllRole() (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreAllRole", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreAllRole")
	defer span.End()

	_, err := s.roleCommand.RestoreAllRole()
	if err != nil {
		return s.errorhandler.HandleRestoreAllRoleError(err, "RestoreAllRole", "FAILED_RESTORE_ALL_ROLE", span, &status, zap.Error(err))
	}

	s.logger.Debug("Successfully restored all roles")
	return true, nil
}

func (s *roleCommandService) DeleteAllRolePermanent() (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteAllRolePermanent", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteAllRolePermanent")
	defer span.End()

	s.logger.Debug("Permanently deleting all roles")

	_, err := s.roleCommand.DeleteAllRolePermanent()

	if err != nil {
		return s.errorhandler.HandleDeleteAllRolePermanentError(err, "DeleteAllRolePermanent", "FAILED_DELETE_ALL_ROLE_PERMANENT", span, &status, zap.Error(err))
	}

	s.logger.Debug("Successfully deleted all roles permanently")
	return true, nil
}

func (s *roleCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
