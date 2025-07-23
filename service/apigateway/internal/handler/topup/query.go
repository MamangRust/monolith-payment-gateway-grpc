package topuphandler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-apigateway/internal/shared"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/topup"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/api"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/topup"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type topupQueryHandleApi struct {
	client pb.TopupQueryServiceClient

	logger logger.LoggerInterface

	mapper apimapper.TopupQueryResponseMapper

	trace trace.Tracer

	requestCounter *prometheus.CounterVec

	requestDuration *prometheus.HistogramVec
}

type topupQueryHandleDeps struct {
	client pb.TopupQueryServiceClient

	router *echo.Echo

	logger logger.LoggerInterface

	mapper apimapper.TopupQueryResponseMapper
}

func NewTopupQueryHandleApi(params *topupQueryHandleDeps) *topupQueryHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "topup_query_handler_requests_total",
			Help: "Total number of topup query requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "topup_query_handler_request_duration_seconds",
			Help:    "Duration of topup query requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	topupHandler := &topupQueryHandleApi{
		client:          params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("topup-query-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerTopup := params.router.Group("/api/topup-query")

	routerTopup.GET("", topupHandler.FindAll)
	routerTopup.GET("/card-number/:card_number", topupHandler.FindAllByCardNumber)
	routerTopup.GET("/:id", topupHandler.FindById)
	routerTopup.GET("/active", topupHandler.FindByActive)
	routerTopup.GET("/trashed", topupHandler.FindByTrashed)

	return topupHandler
}

// @Summary Retrieve a list of all topup data
// @Tags Topup Query
// @Security Bearer
// @Description Retrieve a list of all topup data with pagination and search
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTopup "List of topup data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve topup data"
// @Router /api/topup-query [get]
func (h topupQueryHandleApi) FindAll(c echo.Context) error {
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

	req := &pb.FindAllTopupRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAllTopup(ctx, req)

	if err != nil {
		logError("failed to find all topups", err, zap.Error(err))

		return topup_errors.ErrApiFailedFindAllTopups(c)
	}

	so := h.mapper.ToApiResponsePaginationTopup(res)

	logSuccess("success find all topups", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Find all topup by card number
// @Tags Transaction
// @Security Bearer
// @Description Retrieve a list of transactions for a specific card number
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTopup "List of topups"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve topups data"
// @Router /api/topup-query/card-number/{card_number} [get]
func (h *topupQueryHandleApi) FindAllByCardNumber(c echo.Context) error {
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

	req := &pb.FindAllTopupByCardNumberRequest{
		CardNumber: cardNumber,
		Page:       int32(page),
		PageSize:   int32(pageSize),
		Search:     search,
	}

	res, err := h.client.FindAllTopupByCardNumber(ctx, req)

	if err != nil {
		logError("failed to find all topups by card number", err, zap.Error(err))

		return topup_errors.ErrApiFailedFindAllByCardNumberTopup(c)
	}

	so := h.mapper.ToApiResponsePaginationTopup(res)

	logSuccess("success find all topups by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Find a topup by ID
// @Tags Topup Query
// @Security Bearer
// @Description Retrieve a topup record using its ID
// @Accept json
// @Produce json
// @Param id path string true "Topup ID"
// @Success 200 {object} response.ApiResponseTopup "Topup data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve topup data"
// @Router /api/topup-query/{id} [get]
func (h topupQueryHandleApi) FindById(c echo.Context) error {
	const method = "FindById"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		err := errors.New("invalid topup id")

		logError("Invalid topup id", err, zap.Error(err))

		return topup_errors.ErrApiInvalidTopupID(c)
	}

	res, err := h.client.FindByIdTopup(ctx, &pb.FindByIdTopupRequest{
		TopupId: int32(idInt),
	})

	if err != nil {
		logError("failed to find topup by id", err, zap.Error(err))

		return topup_errors.ErrApiFailedFindByIdTopup(c)
	}

	so := h.mapper.ToApiResponseTopup(res)

	logSuccess("success find topup by id", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Find active topups
// @Tags Topup Query
// @Security Bearer
// @Description Retrieve a list of active topup records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsesTopup "Active topup data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve topup data"
// @Router /api/topup-query/active [get]
func (h *topupQueryHandleApi) FindByActive(c echo.Context) error {
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

	req := &pb.FindAllTopupRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		logError("Failed to retrieve active topups", err, zap.Error(err))

		return topup_errors.ErrApiFailedFindAllTopupsActive(c)
	}

	so := h.mapper.ToApiResponsePaginationTopupDeleteAt(res)

	logSuccess("success find active topups", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Retrieve trashed topups
// @Tags Topup Query
// @Security Bearer
// @Description Retrieve a list of trashed topup records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsesTopup "List of trashed topup data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve topup data"
// @Router /api/topup-query/trashed [get]
func (h *topupQueryHandleApi) FindByTrashed(c echo.Context) error {
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

	req := &pb.FindAllTopupRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		logError("Failed to retrieve trashed topups", err, zap.Error(err))

		return topup_errors.ErrApiFailedFindAllTopupsTrashed(c)
	}

	so := h.mapper.ToApiResponsePaginationTopupDeleteAt(res)

	logSuccess("success find trashed topups", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// startTracingAndLogging starts a tracing span and returns functions to log the outcome of the call.
// The returned functions are logSuccess and logError, which log the outcome of the call to the trace span.
// The returned end function records the metrics and ends the trace span.
func (s *topupQueryHandleApi) startTracingAndLogging(
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

// recordMetrics records a Prometheus metric for the given method and status.
// It increments a counter and records the duration since the provided start time.
func (s *topupQueryHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
