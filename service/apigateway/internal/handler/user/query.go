package userhandler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-apigateway/internal/shared"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/user"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	user_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors/api"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/user"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type userQueryHandleApi struct {
	client pb.UserQueryServiceClient

	logger logger.LoggerInterface

	mapper apimapper.UserQueryResponseMapper

	trace trace.Tracer

	requestCounter *prometheus.CounterVec

	requestDuration *prometheus.HistogramVec
}

type userQueryHandleDeps struct {
	client pb.UserQueryServiceClient

	router *echo.Echo

	logger logger.LoggerInterface

	mapper apimapper.UserQueryResponseMapper
}

func NewUserQueryHandleApi(params *userQueryHandleDeps) *userQueryHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_query_handler_requests_total",
			Help: "Total number of user query requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "user_query_handler_request_duration_seconds",
			Help:    "Duration of user query requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	userQueryHandleApi := &userQueryHandleApi{
		client:          params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("user-query-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerUser := params.router.Group("/api/user-query")

	routerUser.GET("", userQueryHandleApi.FindAllUser)
	routerUser.GET("/:id", userQueryHandleApi.FindById)
	routerUser.GET("/active", userQueryHandleApi.FindByActive)
	routerUser.GET("/trashed", userQueryHandleApi.FindByTrashed)

	return userQueryHandleApi
}

// @Security Bearer
// @Summary Find all users
// @Tags User Command
// @Description Retrieve a list of all users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationUser "List of users"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve user data"
// @Router /api/user-query [get]
func (h *userQueryHandleApi) FindAllUser(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAllUser"
	)

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	page := shared.ParseQueryInt(c, "page", defaultPage)
	pageSize := shared.ParseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	req := &pb.FindAllUserRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)

	if err != nil {
		logError("Failed to retrieve user data", err, zap.Error(err))

		return user_errors.ErrApiFailedFindAll(c)
	}

	so := h.mapper.ToApiResponsePaginationUser(res)

	logSuccess("Successfully retrieve user data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find user by ID
// @Tags User Command
// @Description Retrieve a user by ID
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.ApiResponseUser "User data"
// @Failure 400 {object} response.ErrorResponse "Invalid user ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve user data"
// @Router /api/user-query/{id} [get]
func (h *userQueryHandleApi) FindById(c echo.Context) error {
	const method = "FindById"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logError("Invalid user ID", err, zap.Error(err))

		return user_errors.ErrApiUserInvalidId(c)
	}

	req := &pb.FindByIdUserRequest{
		Id: int32(id),
	}

	user, err := h.client.FindById(ctx, req)

	if err != nil {
		logError("Failed to retrieve user data", err, zap.Error(err))

		return user_errors.ErrApiUserNotFound(c)
	}

	so := h.mapper.ToApiResponseUser(user)

	logSuccess("Successfully retrieve user data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active users
// @Tags User Command
// @Description Retrieve a list of active users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsesUser "List of active users"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve user data"
// @Router /api/user-query/active [get]
func (h *userQueryHandleApi) FindByActive(c echo.Context) error {
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

	req := &pb.FindAllUserRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		logError("Failed to retrieve user data", err, zap.Error(err))

		return user_errors.ErrApiFailedFindActive(c)
	}

	so := h.mapper.ToApiResponsePaginationUserDeleteAt(res)

	logSuccess("Successfully retrieve user data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed user records.
// @Summary Retrieve trashed users
// @Tags User Command
// @Description Retrieve a list of trashed user records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsesUser "List of trashed user data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve user data"
// @Router /api/user-query/trashed [get]
func (h *userQueryHandleApi) FindByTrashed(c echo.Context) error {
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

	req := &pb.FindAllUserRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		logError("Failed to retrieve user data", err, zap.Error(err))

		return user_errors.ErrApiFailedFindTrashed(c)
	}

	so := h.mapper.ToApiResponsePaginationUserDeleteAt(res)

	logSuccess("Successfully retrieve user data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *userQueryHandleApi) startTracingAndLogging(
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

func (s *userQueryHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
