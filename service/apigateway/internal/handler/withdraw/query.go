package withdrawhandler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-apigateway/internal/shared"
	pbhelpers "github.com/MamangRust/monolith-payment-gateway-pb"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/withdraw"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	withdraw_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors/api"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/withdraw"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type withdrawQueryHandleApi struct {
	client pb.WithdrawQueryServiceClient

	logger logger.LoggerInterface

	mapper apimapper.WithdrawQueryResponseMapper

	trace trace.Tracer

	requestCounter *prometheus.CounterVec

	requestDuration *prometheus.HistogramVec
}

type withdrawQueryHandleDeps struct {
	client pb.WithdrawQueryServiceClient

	router *echo.Echo

	logger logger.LoggerInterface

	mapper apimapper.WithdrawQueryResponseMapper
}

func NewWithdrawQueryHandleApi(params *withdrawQueryHandleDeps) *withdrawQueryHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "withdraw_query_handler_requests_total",
			Help: "Total number of withdraw query requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "withdraw_query_handler_request_duration_seconds",
			Help:    "Duration of withdraw query requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	withdrawQueryHandleApi := &withdrawQueryHandleApi{
		client:          params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("withdraw-query-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerWithdraw := params.router.Group("/api/withdraw-query-query")

	routerWithdraw.GET("", withdrawQueryHandleApi.FindAll)
	routerWithdraw.GET("/card-number/:card_number", withdrawQueryHandleApi.FindAllByCardNumber)
	routerWithdraw.GET("/card/:card_number", withdrawQueryHandleApi.FindByCardNumber)

	routerWithdraw.GET("/:id", withdrawQueryHandleApi.FindById)

	routerWithdraw.GET("/active", withdrawQueryHandleApi.FindByActive)
	routerWithdraw.GET("/trashed", withdrawQueryHandleApi.FindByTrashed)

	return withdrawQueryHandleApi
}

// @Summary Find all withdraw records
// @Tags Withdraw Query
// @Security Bearer
// @Description Retrieve a list of all withdraw records with pagination and search
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationWithdraw "List of withdraw records"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve withdraw data"
// @Router /api/withdraw-query [get]
func (h *withdrawQueryHandleApi) FindAll(c echo.Context) error {
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

	req := &pb.FindAllWithdrawRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAllWithdraw(ctx, req)

	if err != nil {
		logError("failed to find all withdraw", err, zap.Error(err))

		return withdraw_errors.ErrApiFailedFindAllWithdraw(c)
	}

	so := h.mapper.ToApiResponsePaginationWithdraw(res)

	logSuccess("success find all withdraw", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Find all withdraw records by card number
// @Tags Withdraw Query
// @Security Bearer
// @Description Retrieve a list of withdraw records for a specific card number with pagination and search
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationWithdraw "List of withdraw records"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve withdraw data"
// @Router /api/withdraw-query/card-number/{card_number} [get]
func (h *withdrawQueryHandleApi) FindAllByCardNumber(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAllByCardNumber"
	)

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	cardNumber, err := shared.ParseQueryCard(c, h.logger)

	if err != nil {
		return err
	}

	page := shared.ParseQueryInt(c, "page", defaultPage)
	pageSize := shared.ParseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	req := &pb.FindAllWithdrawByCardNumberRequest{
		CardNumber: cardNumber,
		Page:       int32(page),
		PageSize:   int32(pageSize),
		Search:     search,
	}

	res, err := h.client.FindAllWithdrawByCardNumber(ctx, req)

	if err != nil {
		logError("failed to find all withdraw by card number", err, zap.Error(err))

		return withdraw_errors.ErrApiFailedFindAllWithdrawByCardNumber(c)
	}

	so := h.mapper.ToApiResponsePaginationWithdraw(res)

	logSuccess("success find all withdraw by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Retrieve all active withdraw data
// @Tags Withdraw Query
// @Security Bearer
// @Description Retrieve a list of all active withdraw data
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponsesWithdraw "List of withdraw data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve withdraw data"
// @Router /api/withdraw-query-query/active [get]
func (h *withdrawQueryHandleApi) FindByActive(c echo.Context) error {
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

	req := &pb.FindAllWithdrawRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		logError("Failed to retrieve withdraw data", err, zap.Error(err))

		return withdraw_errors.ErrApiFailedFindByActiveWithdraw(c)
	}

	so := h.mapper.ToApiResponsePaginationWithdrawDeleteAt(res)

	logSuccess("Success retrieve withdraw data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Retrieve trashed withdraw data
// @Tags Withdraw Query
// @Security Bearer
// @Description Retrieve a list of trashed withdraw data
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponsesWithdraw "List of trashed withdraw data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve withdraw data"
// @Router /api/withdraw-query-query/trashed [get]
func (h *withdrawQueryHandleApi) FindByTrashed(c echo.Context) error {
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

	req := &pb.FindAllWithdrawRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		logError("Failed to retrieve withdraw data", err, zap.Error(err))

		return withdraw_errors.ErrApiFailedFindByTrashedWithdraw(c)
	}

	so := h.mapper.ToApiResponsePaginationWithdrawDeleteAt(res)

	logSuccess("Success retrieve withdraw data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Find a withdraw by ID
// @Tags Withdraw Query
// @Security Bearer
// @Description Retrieve a withdraw record using its ID
// @Accept json
// @Produce json
// @Param id path int true "Withdraw ID"
// @Success 200 {object} response.ApiResponseWithdraw "Withdraw data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve withdraw data"
// @Router /api/withdraw-query/{id} [get]
func (h *withdrawQueryHandleApi) FindById(c echo.Context) error {
	const method = "FindById"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logError("failed to retrieve withdraw data", err, zap.Error(err))

		return withdraw_errors.ErrApiWithdrawInvalidID(c)
	}

	req := &pb.FindByIdWithdrawRequest{
		WithdrawId: int32(id),
	}

	withdraw, err := h.client.FindByIdWithdraw(ctx, req)

	if err != nil {
		logError("failed to retrieve withdraw data", err, zap.Error(err))

		return withdraw_errors.ErrApiFailedFindByIdWithdraw(c)
	}

	so := h.mapper.ToApiResponseWithdraw(withdraw)

	logSuccess("success retrieve withdraw data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Find a withdraw by card number
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve a withdraw record using its card number
// @Accept json
// @Produce json
// @Param card_number query string true "Card number"
// @Success 200 {object} response.ApiResponsesWithdraw "Withdraw data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid card number"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve withdraw data"
// @Router /api/withdraw-query/card/{card_number} [get]
func (h *withdrawQueryHandleApi) FindByCardNumber(c echo.Context) error {
	const method = "FindByCardNumber"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	cardNumber, err := shared.ParseQueryCard(c, h.logger)

	if err != nil {
		return err
	}

	req := &pbhelpers.FindByCardNumberRequest{
		CardNumber: cardNumber,
	}

	withdraw, err := h.client.FindByCardNumber(ctx, req)

	if err != nil {
		logError("Failed to retrieve withdraw data", err, zap.Error(err))

		return withdraw_errors.ErrApiFailedFindByCardNumber(c)
	}

	so := h.mapper.ToApiResponsesWithdraw(withdraw)

	logSuccess("Success retrieve withdraw data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *withdrawQueryHandleApi) startTracingAndLogging(
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

func (s *withdrawQueryHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
