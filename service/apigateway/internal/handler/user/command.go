package userhandler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	pb "github.com/MamangRust/monolith-payment-gateway-pb/user"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	user_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors/api"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/user"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userCommandHandleApi struct {
	client pb.UserCommandServiceClient

	logger logger.LoggerInterface

	mapper apimapper.UserCommandResponseMapper

	trace trace.Tracer

	requestCounter *prometheus.CounterVec

	requestDuration *prometheus.HistogramVec
}

type userCommandHandleDeps struct {
	client pb.UserCommandServiceClient

	router *echo.Echo

	logger logger.LoggerInterface

	mapper apimapper.UserCommandResponseMapper
}

func NewUserCommandHandleApi(params *userCommandHandleDeps) *userCommandHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_command_handler_requests_total",
			Help: "Total number of user command requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "user_command_handler_request_duration_seconds",
			Help:    "Duration of user command requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	userCommandHandleApi := &userCommandHandleApi{
		client:          params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("user-command-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerUser := params.router.Group("/api/user-command")

	routerUser.POST("/create", userCommandHandleApi.Create)
	routerUser.POST("/update/:id", userCommandHandleApi.Update)

	routerUser.POST("/trashed/:id", userCommandHandleApi.TrashedUser)
	routerUser.POST("/restore/:id", userCommandHandleApi.RestoreUser)
	routerUser.DELETE("/permanent/:id", userCommandHandleApi.DeleteUserPermanent)

	routerUser.POST("/restore/all", userCommandHandleApi.RestoreAllUser)
	routerUser.POST("/permanent/all", userCommandHandleApi.DeleteAllUserPermanent)

	return userCommandHandleApi
}

// @Security Bearer
// Create handles the creation of a new user.
// @Summary Create a new user
// @Tags User Command
// @Description Create a new user with the provided details
// @Accept json
// @Produce json
// @Param request body requests.CreateUserRequest true "Create user request"
// @Success 200 {object} response.ApiResponseUser "Successfully created user"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create user"
// @Router /api/user-command/create [post]
func (h *userCommandHandleApi) Create(c echo.Context) error {
	const method = "Create"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	var body requests.CreateUserRequest

	if err := c.Bind(&body); err != nil {
		logError("Failed to bind CreateUser request", err, zap.Error(err))

		return user_errors.ErrApiBindCreateUser(c)
	}

	if err := body.Validate(); err != nil {
		logError("Failed to validate CreateUser request", err, zap.Error(err))

		return user_errors.ErrApiValidateCreateUser(c)
	}

	req := &pb.CreateUserRequest{
		Firstname:       body.FirstName,
		Lastname:        body.LastName,
		Email:           body.Email,
		Password:        body.Password,
		ConfirmPassword: body.ConfirmPassword,
	}

	res, err := h.client.Create(ctx, req)

	if err != nil {
		logError("Failed to create user", err, zap.Error(err))

		return user_errors.ErrApiFailedCreateUser(c)
	}

	so := h.mapper.ToApiResponseUser(res)

	logSuccess("Successfully create user", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// Update handles the update of an existing user record.
// @Summary Update an existing user
// @Tags User Command
// @Description Update an existing user record with the provided details
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param UpdateUserRequest body requests.UpdateUserRequest true "Update user request"
// @Success 200 {object} response.ApiResponseUser "Successfully updated user"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update user"
// @Router /api/user-command/update/{id} [post]
func (h *userCommandHandleApi) Update(c echo.Context) error {
	const method = "Update"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		logError("Invalid user id", err, zap.Error(err))

		return user_errors.ErrApiUserInvalidId(c)
	}

	var body requests.UpdateUserRequest

	if err := c.Bind(&body); err != nil {
		logError("Failed to bind UpdateUser request", err, zap.Error(err))

		return user_errors.ErrApiBindUpdateUser(c)
	}

	if err := body.Validate(); err != nil {
		logError("Failed to validate UpdateUser request", err, zap.Error(err))

		return user_errors.ErrApiValidateUpdateUser(c)
	}

	req := &pb.UpdateUserRequest{
		Id:              int32(idInt),
		Firstname:       body.FirstName,
		Lastname:        body.LastName,
		Email:           body.Email,
		Password:        body.Password,
		ConfirmPassword: body.ConfirmPassword,
	}

	res, err := h.client.Update(ctx, req)

	if err != nil {
		logError("Failed to update user", err, zap.Error(err))

		return user_errors.ErrApiFailedUpdateUser(c)
	}

	so := h.mapper.ToApiResponseUser(res)

	logSuccess("Successfully update user", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// TrashedUser retrieves a trashed user record by its ID.
// @Summary Retrieve a trashed user
// @Tags User Command
// @Description Retrieve a trashed user record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.ApiResponseUser "Successfully retrieved trashed user"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed user"
// @Router /api/user-command/trashed/{id} [post]
func (h *userCommandHandleApi) TrashedUser(c echo.Context) error {
	const method = "TrashedUser"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logError("Invalid user id", err, zap.Error(err))

		return user_errors.ErrApiUserInvalidId(c)
	}

	req := &pb.FindByIdUserRequest{
		Id: int32(id),
	}

	user, err := h.client.TrashedUser(ctx, req)

	if err != nil {
		logError("Failed to trashed user", err, zap.Error(err))

		return user_errors.ErrApiFailedTrashedUser(c)
	}

	so := h.mapper.ToApiResponseUserDeleteAt(user)

	logSuccess("Successfully trashed user", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreUser restores a user record from the trash by its ID.
// @Summary Restore a trashed user
// @Tags User Command
// @Description Restore a trashed user record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.ApiResponseUser "Successfully restored user"
// @Failure 400 {object} response.ErrorResponse "Invalid user ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore user"
// @Router /api/user-command/restore/{id} [post]
func (h *userCommandHandleApi) RestoreUser(c echo.Context) error {
	const method = "RestoreUser"
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

	user, err := h.client.RestoreUser(ctx, req)

	if err != nil {
		logError("Failed to restore user", err, zap.Error(err))

		return user_errors.ErrApiFailedRestoreUser(c)
	}

	so := h.mapper.ToApiResponseUser(user)

	logSuccess("Successfully restore user", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteUserPermanent permanently deletes a user record by its ID.
// @Summary Permanently delete a user
// @Tags User Command
// @Description Permanently delete a user record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.ApiResponseUserDelete "Successfully deleted user record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete user:"
// @Router /api/user-command/delete/{id} [delete]
func (h *userCommandHandleApi) DeleteUserPermanent(c echo.Context) error {
	const method = "DeletUserPermanent"
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

	user, err := h.client.DeleteUserPermanent(ctx, req)

	if err != nil {
		logError("Failed to delete user", err, zap.Error(err))

		return user_errors.ErrApiFailedDeletePermanent(c)
	}

	so := h.mapper.ToApiResponseUserDelete(user)

	logSuccess("Successfully deleted user record permanently", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreUser restores a user record from the trash by its ID.
// @Summary Restore a trashed user
// @Tags User Command
// @Description Restore a trashed user record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.ApiResponseUserAll "Successfully restored user all"
// @Failure 400 {object} response.ErrorResponse "Invalid user ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore user"
// @Router /api/user-command/restore/all [post]
func (h *userCommandHandleApi) RestoreAllUser(c echo.Context) error {
	const method = "RestoreAllUser"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	res, err := h.client.RestoreAllUser(ctx, &emptypb.Empty{})

	if err != nil {
		logError("Failed to restore all user", err, zap.Error(err))

		return user_errors.ErrApiFailedRestoreAll(c)
	}

	so := h.mapper.ToApiResponseUserAll(res)

	logSuccess("Successfully restore all user", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteUserPermanent permanently deletes a user record by its ID.
// @Summary Permanently delete a user
// @Tags User Command
// @Description Permanently delete a user record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.ApiResponseUserDelete "Successfully deleted user record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete user:"
// @Router /api/user-command/delete/all [post]
func (h *userCommandHandleApi) DeleteAllUserPermanent(c echo.Context) error {
	const method = "FindAll"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	res, err := h.client.DeleteAllUserPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		logError("Failed to delete all user permanently", err, zap.Error(err))

		return user_errors.ErrApiFailedDeleteAll(c)
	}

	so := h.mapper.ToApiResponseUserAll(res)

	logSuccess("Successfully deleted all user permanently", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *userCommandHandleApi) startTracingAndLogging(
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

func (s *userCommandHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
