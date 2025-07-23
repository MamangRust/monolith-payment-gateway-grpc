package rolehandler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-apigateway/internal/middlewares"
	mencache "github.com/MamangRust/monolith-payment-gateway-apigateway/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-apigateway/internal/shared"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/role"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	role_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/role_errors/api"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/role"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type roleQueryHandlerApi struct {
	kafka *kafka.Kafka

	role pb.RoleQueryServiceClient

	// logger provides structured logging capabilities for debugging and tracing.
	logger logger.LoggerInterface

	// mapper maps gRPC responses into API-compliant response formats.
	mapper apimapper.RoleQueryResponseMapper

	// trace enables distributed tracing for handler operations via OpenTelemetry.
	trace trace.Tracer

	// requestCounter counts the number of incoming requests handled.
	requestCounter *prometheus.CounterVec

	// requestDuration records the duration of each request in seconds.
	requestDuration *prometheus.HistogramVec
}

type roleQueryHandleDeps struct {
	client pb.RoleQueryServiceClient
	router *echo.Echo
	logger logger.LoggerInterface
	mapper apimapper.RoleQueryResponseMapper
	kafka  *kafka.Kafka
	cache  mencache.RoleCache
}

func NewRoleQueryHandleApi(params *roleQueryHandleDeps) *roleQueryHandlerApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "role_query_handler_requests_total",
			Help: "Total number of role requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "role_query_handler_request_duration_seconds",
			Help:    "Duration of role requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	roleQueryHandler := &roleQueryHandlerApi{
		role:            params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		kafka:           params.kafka,
		trace:           otel.Tracer("role-query-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	roleMiddleware := middlewares.NewRoleValidator(params.kafka, "request-role", "response-role", 5*time.Second, params.logger, params.cache)

	routerRole := params.router.Group("/api/role-query")

	roleMiddlewareChain := roleMiddleware.Middleware()
	requireAdmin := middlewares.RequireRoles("Admin_Role_10")

	routerRole.GET("", roleMiddlewareChain(requireAdmin(roleQueryHandler.FindAll)))

	routerRole.GET("/:id", roleMiddlewareChain(requireAdmin(roleQueryHandler.FindById)))

	routerRole.GET("/active", roleMiddlewareChain(requireAdmin(roleQueryHandler.FindAll)))

	routerRole.GET("/trashed", roleMiddlewareChain(requireAdmin(roleQueryHandler.FindByTrashed)))

	routerRole.GET("/user/:user_id", roleMiddlewareChain(requireAdmin(roleQueryHandler.FindByUserId)))

	return roleQueryHandler
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
func (h *roleQueryHandlerApi) FindAll(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAll"
	)

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	page := shared.ParseQueryInt(c, "page", defaultPage)
	pageSize := shared.ParseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	req := &pb.FindAllRoleRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.role.FindAllRole(ctx, req)

	if err != nil {
		logError("failed to fetch roles", err, zap.Error(err))

		return role_errors.ErrApiFailedFindAll(c)
	}

	so := h.mapper.ToApiResponsePaginationRole(res)

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
func (h *roleQueryHandlerApi) FindById(c echo.Context) error {
	const method = "FindById"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	roleID, err := strconv.Atoi(c.Param("id"))

	if err != nil || roleID <= 0 {
		logError("invalid role ID", err, zap.Error(err))

		return role_errors.ErrApiRoleInvalidId(c)
	}

	req := &pb.FindByIdRoleRequest{
		RoleId: int32(roleID),
	}

	res, err := h.role.FindByIdRole(ctx, req)

	if err != nil {
		logError("failed to fetch role", err, zap.Error(err))

		return role_errors.ErrApiRoleNotFound(c)
	}

	so := h.mapper.ToApiResponseRole(res)

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
func (h *roleQueryHandlerApi) FindByActive(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindByActive"
	)

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	page := shared.ParseQueryInt(c, "page", defaultPage)
	pageSize := shared.ParseQueryInt(c, "page_size", defaultPageSize)

	search := c.QueryParam("search")

	req := &pb.FindAllRoleRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.role.FindByActive(ctx, req)

	if err != nil {
		logError("failed to fetch active roles", err, zap.Error(err))

		return role_errors.ErrApiFailedFindActive(c)
	}

	so := h.mapper.ToApiResponsePaginationRoleDeleteAt(res)

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
func (h *roleQueryHandlerApi) FindByTrashed(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindByTrashed"
	)

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	page := shared.ParseQueryInt(c, "page", defaultPage)
	pageSize := shared.ParseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	req := &pb.FindAllRoleRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.role.FindByTrashed(ctx, req)

	if err != nil {
		logError("failed to fetch trashed roles", err, zap.Error(err))

		return role_errors.ErrApiFailedFindTrashed(c)
	}

	so := h.mapper.ToApiResponsePaginationRoleDeleteAt(res)

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
func (h *roleQueryHandlerApi) FindByUserId(c echo.Context) error {
	const method = "FindAll"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	userID, err := strconv.Atoi(c.Param("user_id"))

	if err != nil || userID <= 0 {
		logError("invalid user id", err, zap.Error(err))

		return role_errors.ErrApiRoleInvalidId(c)
	}

	req := &pb.FindByIdUserRoleRequest{
		UserId: int32(userID),
	}

	res, err := h.role.FindByUserId(ctx, req)

	if err != nil {
		logError("failed to fetch role by user id", err, zap.Error(err))

		return role_errors.ErrApiRoleNotFound(c)
	}

	so := h.mapper.ToApiResponsesRole(res)

	logSuccess("fetch role by user id successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *roleQueryHandlerApi) startTracingAndLogging(
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

func (s *roleQueryHandlerApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
