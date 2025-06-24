package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-apigateway/internal/middlewares"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/role_errors"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api"
	"github.com/MamangRust/monolith-payment-gateway-shared/pb"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type roleHandleApi struct {
	kafka           *kafka.Kafka
	role            pb.RoleServiceClient
	logger          logger.LoggerInterface
	mapping         apimapper.RoleResponseMapper
	trace           trace.Tracer
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewHandlerRole(role pb.RoleServiceClient, router *echo.Echo, logger logger.LoggerInterface, mapping apimapper.RoleResponseMapper, kafka *kafka.Kafka) *roleHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "role_handler_requests_total",
			Help: "Total number of role requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "role_handler_request_duration_seconds",
			Help:    "Duration of role requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	roleHandler := &roleHandleApi{
		role:            role,
		logger:          logger,
		mapping:         mapping,
		trace:           otel.Tracer("role-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
		kafka:           kafka,
	}

	roleMiddleware := middlewares.NewRoleValidator(kafka, "request-role", "response-role", 5*time.Second, logger)

	routerRole := router.Group("/api/role")

	roleMiddlewareChain := roleMiddleware.Middleware()
	requireAdmin := middlewares.RequireRoles("Admin_Admin_14")

	routerRole.GET("", roleMiddlewareChain(requireAdmin(roleHandler.FindAll)))

	routerRole.GET("/:id", roleHandler.FindById)

	routerRole.GET("/active", roleMiddlewareChain(requireAdmin(roleHandler.FindAll)))

	routerRole.GET("/trashed", roleHandler.FindByTrashed)

	routerRole.GET("/user/:user_id", roleHandler.FindByUserId)

	routerRole.POST("/",
		roleHandler.Create,
		roleMiddleware.Middleware(),
		middlewares.RequireRoles("Admin_Admin_14"),
	)

	routerRole.POST("/:id", roleHandler.Update, roleMiddleware.Middleware(), middlewares.RequireRoles("Admin_Admin_14"))

	routerRole.DELETE("/:id", roleHandler.Trashed)
	routerRole.PUT("/restore/:id", roleHandler.Restore)
	routerRole.DELETE("/permanent/:id", roleHandler.DeletePermanent)

	routerRole.POST("/restore/all", roleHandler.RestoreAll)
	routerRole.POST("/permanent/all", roleHandler.DeleteAllPermanent)

	return roleHandler
}

// FindAll godoc.
// @Summary Get all roles
// @Tags Role
// @Security Bearer
// @Description Retrieve a paginated list of roles with optional search and pagination parameters.
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Number of items per page (default: 10)"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponsePaginationRole "List of roles"
// @Failure 400 {object} response.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} response.ErrorResponse "Failed to fetch roles"
// @Router /api/role [get]
func (h *roleHandleApi) FindAll(c echo.Context) error {
	const method = "FindAll"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")

	req := &pb.FindAllRoleRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.role.FindAllRole(ctx, req)
	if err != nil {
		status = "error"

		logError("failed to fetch roles", err, zap.Error(err))

		return role_errors.ErrApiFailedFindAll(c)
	}

	so := h.mapping.ToApiResponsePaginationRole(res)

	logSuccess("fetch roles successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindById godoc.
// @Summary Get a role by ID
// @Tags Role
// @Security Bearer
// @Description Retrieve a role by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} response.ApiResponseRole "Role data"
// @Failure 400 {object} response.ErrorResponse "Invalid role ID"
// @Failure 500 {object} response.ErrorResponse "Failed to fetch role"
// @Router /api/role/{id} [get]
func (h *roleHandleApi) FindById(c echo.Context) error {
	const method = "FindById"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil || roleID <= 0 {
		status = "error"

		logError("invalid role ID", err, zap.Error(err))

		return role_errors.ErrApiRoleInvalidId(c)
	}

	req := &pb.FindByIdRoleRequest{
		RoleId: int32(roleID),
	}

	res, err := h.role.FindByIdRole(ctx, req)
	if err != nil {
		status = "error"

		logError("failed to fetch role", err, zap.Error(err))

		return role_errors.ErrApiRoleNotFound(c)
	}

	so := h.mapping.ToApiResponseRole(res)

	logSuccess("fetch role successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindByActive godoc.
// @Summary Get active roles
// @Tags Role
// @Security Bearer
// @Description Retrieve a paginated list of active roles with optional search and pagination parameters.
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Number of items per page (default: 10)"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponsePaginationRoleDeleteAt "List of active roles"
// @Failure 400 {object} response.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} response.ErrorResponse "Failed to fetch active roles"
// @Router /api/role/active [get]
func (h *roleHandleApi) FindByActive(c echo.Context) error {
	const method = "FindByActive"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")

	req := &pb.FindAllRoleRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.role.FindByActive(ctx, req)
	if err != nil {
		status = "error"

		logError("failed to fetch active roles", err, zap.Error(err))

		return role_errors.ErrApiFailedFindActive(c)
	}

	so := h.mapping.ToApiResponsePaginationRoleDeleteAt(res)

	logSuccess("fetch active roles successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindByTrashed godoc.
// @Summary Get trashed roles
// @Tags Role
// @Security Bearer
// @Description Retrieve a paginated list of trashed roles with optional search and pagination parameters.
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Number of items per page (default: 10)"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponsePaginationRoleDeleteAt "List of trashed roles"
// @Failure 400 {object} response.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} response.ErrorResponse "Failed to fetch trashed roles"
// @Router /api/role/trashed [get]
func (h *roleHandleApi) FindByTrashed(c echo.Context) error {
	const method = "FindByTrashed"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")

	req := &pb.FindAllRoleRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.role.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"

		logError("failed to fetch trashed roles", err, zap.Error(err))

		return role_errors.ErrApiFailedFindTrashed(c)
	}

	so := h.mapping.ToApiResponsePaginationRoleDeleteAt(res)

	logSuccess("fetch trashed roles successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindByUserId godoc.
// @Summary Get role by user ID
// @Tags Role
// @Security Bearer
// @Description Retrieve a role by the associated user ID.
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} response.ApiResponseRole "Role data"
// @Failure 400 {object} response.ErrorResponse "Invalid user ID"
// @Failure 500 {object} response.ErrorResponse "Failed to fetch role by user ID"
// @Router /api/role/user/{user_id} [get]
func (h *roleHandleApi) FindByUserId(c echo.Context) error {
	const method = "FindAll"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil || userID <= 0 {
		status = "error"

		logError("invalid user id", err, zap.Error(err))

		return role_errors.ErrApiRoleInvalidId(c)
	}

	req := &pb.FindByIdUserRoleRequest{
		UserId: int32(userID),
	}

	res, err := h.role.FindByUserId(ctx, req)
	if err != nil {
		status = "error"

		logError("failed to fetch role by user id", err, zap.Error(err))

		return role_errors.ErrApiRoleNotFound(c)
	}

	so := h.mapping.ToApiResponsesRole(res)

	logSuccess("fetch role by user id successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
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
func (h *roleHandleApi) Create(c echo.Context) error {
	const method = "Create"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	var req requests.CreateRoleRequest

	if err := c.Bind(&req); err != nil {
		status = "error"

		logError("Failed to bind CreateRole request", err, zap.Error(err))

		return role_errors.ErrApiBindCreateRole(c)
	}

	if err := req.Validate(); err != nil {
		status = "error"

		logError("Failed to validate CreateRole request", err, zap.Error(err))

		return role_errors.ErrApiValidateCreateRole(c)
	}

	reqPb := &pb.CreateRoleRequest{
		Name: req.Name,
	}

	res, err := h.role.CreateRole(ctx, reqPb)

	if err != nil {
		status = "error"

		logError("Failed to create role", err, zap.Error(err))

		return role_errors.ErrApiFailedCreateRole(c)
	}

	so := h.mapping.ToApiResponseRole(res)

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
func (h *roleHandleApi) Update(c echo.Context) error {
	const method = "Update"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil || roleID <= 0 {
		status = "error"

		logError("Failed to update role", err, zap.Error(err))

		return role_errors.ErrInvalidRoleId(c)
	}

	var req requests.UpdateRoleRequest
	if err := c.Bind(&req); err != nil {
		status = "error"

		logError("Failed to bind UpdateRole request", err, zap.Error(err))

		return role_errors.ErrApiBindUpdateRole(c)
	}

	if err := req.Validate(); err != nil {
		status = "error"

		logError("Failed to validate UpdateRole request", err, zap.Error(err))

		return role_errors.ErrApiValidateUpdateRole(c)
	}

	reqPb := &pb.UpdateRoleRequest{
		Id:   int32(roleID),
		Name: req.Name,
	}

	res, err := h.role.UpdateRole(ctx, reqPb)
	if err != nil {
		status = "error"

		logError("Failed to update role", err, zap.Error(err))

		return role_errors.ErrApiFailedUpdateRole(c)
	}

	so := h.mapping.ToApiResponseRole(res)

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
func (h *roleHandleApi) Trashed(c echo.Context) error {
	const method = "Trashed"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	roleID, err := strconv.Atoi(c.Param("id"))

	if err != nil || roleID <= 0 {
		status = "error"

		logError("Failed to trash role", err, zap.Error(err))

		return role_errors.ErrInvalidRoleId(c)
	}

	req := &pb.FindByIdRoleRequest{
		RoleId: int32(roleID),
	}

	res, err := h.role.TrashedRole(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to trash role", err, zap.Error(err))

		return role_errors.ErrApiFailedTrashedRole(c)
	}

	so := h.mapping.ToApiResponseRole(res)

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
func (h *roleHandleApi) Restore(c echo.Context) error {
	const method = "Restore"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil || roleID <= 0 {
		status = "error"

		logError("Failed to restore role", err, zap.Error(err))

		return role_errors.ErrInvalidRoleId(c)
	}

	req := &pb.FindByIdRoleRequest{
		RoleId: int32(roleID),
	}

	res, err := h.role.RestoreRole(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to restore role", err, zap.Error(err))

		return role_errors.ErrApiFailedRestoreRole(c)
	}

	so := h.mapping.ToApiResponseRole(res)

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
func (h *roleHandleApi) DeletePermanent(c echo.Context) error {
	const method = "DeletePermanent"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil || roleID <= 0 {
		status = "error"

		logError("Failed to delete role permanently", err, zap.Error(err))

		return role_errors.ErrInvalidRoleId(c)
	}

	req := &pb.FindByIdRoleRequest{
		RoleId: int32(roleID),
	}

	res, err := h.role.DeleteRolePermanent(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to delete role permanently", err, zap.Error(err))

		return role_errors.ErrApiFailedDeletePermanent(c)
	}

	so := h.mapping.ToApiResponseRoleDelete(res)

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
func (h *roleHandleApi) RestoreAll(c echo.Context) error {
	const method = "RestoreAll"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	res, err := h.role.RestoreAllRole(ctx, &emptypb.Empty{})
	if err != nil {
		status = "error"
		logError("Failed to restore all roles", err, zap.Error(err))

		return role_errors.ErrApiFailedRestoreAll(c)
	}

	so := h.mapping.ToApiResponseRoleAll(res)

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
func (h *roleHandleApi) DeleteAllPermanent(c echo.Context) error {
	const method = "DeleteAllPermanent"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	res, err := h.role.DeleteAllRolePermanent(ctx, &emptypb.Empty{})
	if err != nil {
		status = "error"
		logError("Failed to delete all roles permanently", err, zap.Error(err))

		return role_errors.ErrApiFailedDeleteAll(c)
	}

	so := h.mapping.ToApiResponseRoleAll(res)

	logSuccess("Delete all roles permanently successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *roleHandleApi) startTracingAndLogging(
	ctx context.Context,
	method string,
	attrs ...attribute.KeyValue,
) (func(string), func(string, ...zap.Field), func(string, error, ...zap.Field)) {
	start := time.Now()
	_, span := s.trace.Start(ctx, method)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.AddEvent("Start: " + method)
	s.logger.Debug("Start: " + method)

	end := func(status string) {
		s.recordMetrics(method, status, start)
		code := otelcode.Ok
		if status != "success" {
			code = otelcode.Error
		}
		span.SetStatus(code, status)
		span.End()
	}

	logSuccess := func(msg string, fields ...zap.Field) {
		span.AddEvent(msg)
		s.logger.Debug(msg, fields...)
	}

	logError := func(msg string, err error, fields ...zap.Field) {
		span.RecordError(err)
		span.SetStatus(otelcode.Error, msg)
		span.AddEvent(msg)
		allFields := append([]zap.Field{zap.Error(err)}, fields...)
		s.logger.Error(msg, allFields...)
	}

	return end, logSuccess, logError
}

func (s *roleHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
