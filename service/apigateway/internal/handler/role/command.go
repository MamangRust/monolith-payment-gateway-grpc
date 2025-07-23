package rolehandler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-apigateway/internal/middlewares"
	mencache "github.com/MamangRust/monolith-payment-gateway-apigateway/internal/redis"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/role"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	role_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/role_errors/api"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/role"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

// roleCommandHandleApi provides HTTP handlers for role-related operations.
//
// This struct integrates the RoleService gRPC client, logging, Kafka event publishing,
// response mapper, tracing, and Prometheus metrics for complete observability and functionality.
type roleCommandHandleApi struct {
	// kafka is the Kafka client used for publishing role-related events.
	kafka *kafka.Kafka

	// role is the gRPC client to communicate with the RoleService.
	role pb.RoleCommandServiceClient

	// logger provides structured logging capabilities for debugging and tracing.
	logger logger.LoggerInterface

	// mapper maps gRPC responses into API-compliant response formats.
	mapper apimapper.RoleCommandResponseMapper

	// trace enables distributed tracing for handler operations via OpenTelemetry.
	trace trace.Tracer

	// requestCounter counts the number of incoming requests handled.
	requestCounter *prometheus.CounterVec

	// requestDuration records the duration of each request in seconds.
	requestDuration *prometheus.HistogramVec
}

// roleHandleDeps contains the dependencies required to initialize and register
// the role handler routes into the Echo router.
type roleCommandHandleDeps struct {
	// client is the gRPC client for communicating with the RoleService.
	client pb.RoleCommandServiceClient

	// router is the Echo instance used to register HTTP routes.
	router *echo.Echo

	// logger provides structured and contextual logging throughout the handler.
	logger logger.LoggerInterface

	// mapper transforms internal gRPC responses into HTTP-friendly API responses.
	mapper apimapper.RoleCommandResponseMapper

	// kafka is used to publish domain events to Kafka topics (e.g., for auditing or async workflows).
	kafka *kafka.Kafka

	cache mencache.RoleCache
}

func NewRoleCommandHandleApi(params *roleCommandHandleDeps) *roleCommandHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "role_command_handler_requests_total",
			Help: "Total number of role requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "role_command_handler_request_duration_seconds",
			Help:    "Duration of role requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	roleHandler := &roleCommandHandleApi{
		role:            params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("role-command-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
		kafka:           params.kafka,
	}

	roleMiddleware := middlewares.NewRoleValidator(params.kafka, "request-role", "response-role", 5*time.Second, params.logger, params.cache)

	routerRole := params.router.Group("/api/role")

	roleMiddlewareChain := roleMiddleware.Middleware()
	requireAdmin := middlewares.RequireRoles("Admin_Admin_14")

	routerRole.POST("/",
		roleMiddlewareChain(requireAdmin(roleHandler.Create)),
	)

	routerRole.POST("/:id", roleMiddlewareChain(requireAdmin(roleHandler.Update)))

	routerRole.DELETE("/:id", roleMiddlewareChain(requireAdmin(roleHandler.DeletePermanent)))
	routerRole.PUT("/restore/:id", roleMiddlewareChain(requireAdmin(roleHandler.Restore)))
	routerRole.DELETE("/permanent/:id", roleMiddlewareChain(requireAdmin(roleHandler.DeletePermanent)))

	routerRole.POST("/restore/all", roleMiddlewareChain(requireAdmin(roleHandler.RestoreAll)))
	routerRole.POST("/permanent/all", roleMiddlewareChain(requireAdmin(roleHandler.DeleteAllPermanent)))

	return roleHandler
}

// Create godoc.
// @Summary Create a new role
// @Tags Role
// @Security Bearer
// @Description Create a new role with the provided details.
// @Accept json
// @Produce json
// @Param request body requests.CreateRoleRequest true "Role data"
// @Success 200 {object} response.ApiResponseRole "Created role data"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 500 {object} response.ErrorResponse "Failed to create role"
// @Router /api/role [post]
func (h *roleCommandHandleApi) Create(c echo.Context) error {
	const method = "Create"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	var req requests.CreateRoleRequest

	if err := c.Bind(&req); err != nil {
		logError("Failed to bind CreateRole request", err, zap.Error(err))

		return role_errors.ErrApiBindCreateRole(c)
	}

	if err := req.Validate(); err != nil {
		logError("Failed to validate CreateRole request", err, zap.Error(err))

		return role_errors.ErrApiValidateCreateRole(c)
	}

	reqPb := &pb.CreateRoleRequest{
		Name: req.Name,
	}

	res, err := h.role.CreateRole(ctx, reqPb)

	if err != nil {
		logError("Failed to create role", err, zap.Error(err))

		return role_errors.ErrApiFailedCreateRole(c)
	}

	so := h.mapper.ToApiResponseRole(res)

	logSuccess("Create role successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// Update godoc.
// @Summary Update a role
// @Tags Role
// @Security Bearer
// @Description Update an existing role with the provided details.
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param request body requests.UpdateRoleRequest true "Role data"
// @Success 200 {object} response.ApiResponseRole "Updated role data"
// @Failure 400 {object} response.ErrorResponse "Invalid role ID or request body"
// @Failure 500 {object} response.ErrorResponse "Failed to update role"
// @Router /api/role/{id} [post]
func (h *roleCommandHandleApi) Update(c echo.Context) error {
	const method = "Update"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	roleID, err := strconv.Atoi(c.Param("id"))

	if err != nil || roleID <= 0 {
		logError("Failed to update role", err, zap.Error(err))

		return role_errors.ErrInvalidRoleId(c)
	}

	var req requests.UpdateRoleRequest

	if err := c.Bind(&req); err != nil {
		logError("Failed to bind UpdateRole request", err, zap.Error(err))

		return role_errors.ErrApiBindUpdateRole(c)
	}

	if err := req.Validate(); err != nil {
		logError("Failed to validate UpdateRole request", err, zap.Error(err))

		return role_errors.ErrApiValidateUpdateRole(c)
	}

	reqPb := &pb.UpdateRoleRequest{
		Id:   int32(roleID),
		Name: req.Name,
	}

	res, err := h.role.UpdateRole(ctx, reqPb)

	if err != nil {
		logError("Failed to update role", err, zap.Error(err))

		return role_errors.ErrApiFailedUpdateRole(c)
	}

	so := h.mapper.ToApiResponseRole(res)

	logSuccess("Update role successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// Trashed godoc.
// @Summary Soft-delete a role
// @Tags Role
// @Security Bearer
// @Description Soft-delete a role by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} response.ApiResponseRole "Soft-deleted role data"
// @Failure 400 {object} response.ErrorResponse "Invalid role ID"
// @Failure 500 {object} response.ErrorResponse "Failed to soft-delete role"
// @Router /api/role/{id} [delete]
func (h *roleCommandHandleApi) Trashed(c echo.Context) error {
	const method = "Trashed"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	roleID, err := strconv.Atoi(c.Param("id"))

	if err != nil || roleID <= 0 {
		logError("Failed to trash role", err, zap.Error(err))

		return role_errors.ErrInvalidRoleId(c)
	}

	req := &pb.FindByIdRoleRequest{
		RoleId: int32(roleID),
	}

	res, err := h.role.TrashedRole(ctx, req)

	if err != nil {
		logError("Failed to trash role", err, zap.Error(err))

		return role_errors.ErrApiFailedTrashedRole(c)
	}

	so := h.mapper.ToApiResponseRoleDeleteAt(res)

	logSuccess("Trash role successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// Restore godoc.
// @Summary Restore a soft-deleted role
// @Tags Role
// @Security Bearer
// @Description Restore a soft-deleted role by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} response.ApiResponseRole "Restored role data"
// @Failure 400 {object} response.ErrorResponse "Invalid role ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore role"
// @Router /api/role/restore/{id} [put]
func (h *roleCommandHandleApi) Restore(c echo.Context) error {
	const method = "Restore"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	roleID, err := strconv.Atoi(c.Param("id"))

	if err != nil || roleID <= 0 {
		logError("Failed to restore role", err, zap.Error(err))

		return role_errors.ErrInvalidRoleId(c)
	}

	req := &pb.FindByIdRoleRequest{
		RoleId: int32(roleID),
	}

	res, err := h.role.RestoreRole(ctx, req)

	if err != nil {
		logError("Failed to restore role", err, zap.Error(err))

		return role_errors.ErrApiFailedRestoreRole(c)
	}

	so := h.mapper.ToApiResponseRole(res)

	logSuccess("Restore role successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// DeletePermanent godoc.
// @Summary Permanently delete a role
// @Tags Role
// @Security Bearer
// @Description Permanently delete a role by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} response.ApiResponseRole "Permanently deleted role data"
// @Failure 400 {object} response.ErrorResponse "Invalid role ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete role permanently"
// @Router /api/role/permanent/{id} [delete]
func (h *roleCommandHandleApi) DeletePermanent(c echo.Context) error {
	const method = "DeletePermanent"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	roleID, err := strconv.Atoi(c.Param("id"))

	if err != nil || roleID <= 0 {
		logError("Failed to delete role permanently", err, zap.Error(err))

		return role_errors.ErrInvalidRoleId(c)
	}

	req := &pb.FindByIdRoleRequest{
		RoleId: int32(roleID),
	}

	res, err := h.role.DeleteRolePermanent(ctx, req)

	if err != nil {
		logError("Failed to delete role permanently", err, zap.Error(err))

		return role_errors.ErrApiFailedDeletePermanent(c)
	}

	so := h.mapper.ToApiResponseRoleDelete(res)

	logSuccess("Delete role permanently successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// RestoreAll godoc.
// @Summary Restore all soft-deleted roles
// @Tags Role
// @Security Bearer
// @Description Restore all soft-deleted roles.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseRoleAll "Restored roles data"
// @Failure 500 {object} response.ErrorResponse "Failed to restore all roles"
// @Router /api/role/restore/all [post]
func (h *roleCommandHandleApi) RestoreAll(c echo.Context) error {
	const method = "RestoreAll"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	res, err := h.role.RestoreAllRole(ctx, &emptypb.Empty{})

	if err != nil {
		logError("Failed to restore all roles", err, zap.Error(err))

		return role_errors.ErrApiFailedRestoreAll(c)
	}

	so := h.mapper.ToApiResponseRoleAll(res)

	logSuccess("Restore all roles successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// DeleteAllPermanent godoc.
// @Summary Permanently delete all roles
// @Tags Role
// @Security Bearer
// @Description Permanently delete all roles.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseRoleAll "Permanently deleted roles data"
// @Failure 500 {object} response.ErrorResponse "Failed to delete all roles permanently"
// @Router /api/role/permanent/all [post]
func (h *roleCommandHandleApi) DeleteAllPermanent(c echo.Context) error {
	const method = "DeleteAllPermanent"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	res, err := h.role.DeleteAllRolePermanent(ctx, &emptypb.Empty{})
	if err != nil {
		logError("Failed to delete all roles permanently", err, zap.Error(err))

		return role_errors.ErrApiFailedDeleteAll(c)
	}

	so := h.mapper.ToApiResponseRoleAll(res)

	logSuccess("Delete all roles permanently successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *roleCommandHandleApi) startTracingAndLogging(
	ctx context.Context,
	method string,
	attrs ...attribute.KeyValue,
) (
	end func(),
	logSuccess func(string, ...zap.Field),
	logError func(string, error, ...zap.Field),
) {
	start := time.Now()
	_, span := s.trace.Start(ctx, method)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.AddEvent("Start: " + method)
	s.logger.Debug("Start: " + method)

	status := "success"

	end = func() {
		s.recordMetrics(method, status, start)
		code := otelcode.Ok
		if status != "success" {
			code = otelcode.Error
		}
		span.SetStatus(code, status)
		span.End()
	}

	logSuccess = func(msg string, fields ...zap.Field) {
		status = "success"
		span.AddEvent(msg)
		s.logger.Debug(msg, fields...)
	}

	logError = func(msg string, err error, fields ...zap.Field) {
		status = "error"
		span.RecordError(err)
		span.SetStatus(otelcode.Error, msg)
		span.AddEvent(msg)
		allFields := append([]zap.Field{zap.Error(err)}, fields...)
		s.logger.Error(msg, allFields...)
	}

	return end, logSuccess, logError
}

func (s *roleCommandHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
