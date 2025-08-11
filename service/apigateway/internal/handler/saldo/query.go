package saldohandler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-apigateway/internal/shared"
	pbhelpers "github.com/MamangRust/monolith-payment-gateway-pb"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/saldo"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/api"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/saldo"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type saldoQueryHandleApi struct {
	saldo pb.SaldoQueryServiceClient

	logger logger.LoggerInterface

	mapper apimapper.SaldoQueryResponseMapper

	trace trace.Tracer

	requestCounter *prometheus.CounterVec

	requestDuration *prometheus.HistogramVec
}

type saldoQueryHandleDeps struct {
	client pb.SaldoQueryServiceClient

	router *echo.Echo

	logger logger.LoggerInterface

	mapper apimapper.SaldoQueryResponseMapper
}

func NewSaldoQueryHandleApi(params *saldoQueryHandleDeps) *saldoQueryHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "saldo_query_handler_requests_total",
			Help: "Total number of saldo query requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "saldo_query_handler_request_duration_seconds",
			Help:    "Duration of saldo query requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	saldoHandler := &saldoQueryHandleApi{
		saldo:           params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("saldo-query-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerSaldo := params.router.Group("/api/saldo-query")

	routerSaldo.GET("", saldoHandler.FindAll)
	routerSaldo.GET("/:id", saldoHandler.FindById)

	routerSaldo.GET("/active", saldoHandler.FindByActive)
	routerSaldo.GET("/trashed", saldoHandler.FindByTrashed)
	routerSaldo.GET("/card_number/:card_number", saldoHandler.FindByCardNumber)

	return saldoHandler
}

// @Summary Find all saldo data
// @Tags Saldo-Query
// @Security Bearer
// @Description Retrieve a list of all saldo data with pagination and search
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationSaldo "List of saldo data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve saldo data"
// @Router /api/saldo-query [get]
func (h *saldoQueryHandleApi) FindAll(c echo.Context) error {
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

	req := &pb.FindAllSaldoRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.saldo.FindAllSaldo(ctx, req)

	if err != nil {
		logError("Failed to retrieve saldo data", err, zap.Error(err))

		return saldo_errors.ErrApiFailedFindAllSaldo(c)
	}

	so := h.mapper.ToApiResponsePaginationSaldo(res)

	logSuccess("Successfully retrieve saldo data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Find a saldo by ID
// @Tags Saldo-Query
// @Security Bearer
// @Description Retrieve a saldo by its ID
// @Accept json
// @Produce json
// @Param id path int true "Saldo ID"
// @Success 200 {object} response.ApiResponseSaldo "Saldo data"
// @Failure 400 {object} response.ErrorResponse "Invalid saldo ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve saldo data"
// @Router /api/saldo-query/{id} [get]
func (h *saldoQueryHandleApi) FindById(c echo.Context) error {
	const method = "FindById"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logError("Invalid saldo ID", err, zap.Error(err))

		return saldo_errors.ErrApiInvalidSaldoID(c)
	}

	req := &pb.FindByIdSaldoRequest{
		SaldoId: int32(id),
	}

	res, err := h.saldo.FindByIdSaldo(ctx, req)

	if err != nil {
		logError("Failed to retrieve saldo data", err, zap.Error(err))

		return saldo_errors.ErrApiFailedFindByIdSaldo(c)
	}

	so := h.mapper.ToApiResponseSaldo(res)

	logSuccess("Successfully retrieve saldo data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Find a saldo by card number
// @Tags Saldo-Query
// @Security Bearer
// @Description Retrieve a saldo by its card number
// @Accept json
// @Produce json
// @Param card_number path string true "Card number"
// @Success 200 {object} response.ApiResponseSaldo "Saldo data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve saldo data"
// @Router /api/saldo-query/card_number/{card_number} [get]
func (h *saldoQueryHandleApi) FindByCardNumber(c echo.Context) error {
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

	res, err := h.saldo.FindByCardNumber(ctx, req)

	if err != nil {
		logError("Failed to retrieve saldo data", err, zap.Error(err))

		return saldo_errors.ErrApiFailedFindByCardNumberSaldo(c)
	}

	so := h.mapper.ToApiResponseSaldo(res)

	logSuccess("Successfully retrieve saldo data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Retrieve all active saldo data
// @Tags Saldo-Query
// @Security Bearer
// @Description Retrieve a list of all active saldo data
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsesSaldo "List of saldo data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve saldo data"
// @Router /api/saldo-query/active [get]
func (h *saldoQueryHandleApi) FindByActive(c echo.Context) error {
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

	req := &pb.FindAllSaldoRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.saldo.FindByActive(ctx, req)

	if err != nil {
		logError("Failed to retrieve saldo data", err, zap.Error(err))
		return saldo_errors.ErrApiFailedFindAllSaldoActive(c)
	}

	so := h.mapper.ToApiResponsePaginationSaldoDeleteAt(res)

	logSuccess("Successfully retrieve saldo data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Retrieve trashed saldo data
// @Tags Saldo-Query
// @Security Bearer
// @Description Retrieve a list of all trashed saldo data
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsesSaldo "List of trashed saldo data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve saldo data"
// @Router /api/saldo-query/trashed [get]
func (h *saldoQueryHandleApi) FindByTrashed(c echo.Context) error {
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

	req := &pb.FindAllSaldoRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.saldo.FindByTrashed(ctx, req)

	if err != nil {
		logError("Failed to retrieve trashed saldo data", err, zap.Error(err))

		return saldo_errors.ErrApiFailedFindAllSaldoTrashed(c)
	}

	so := h.mapper.ToApiResponsePaginationSaldoDeleteAt(res)

	logSuccess("Successfully retrieve trashed saldo data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *saldoQueryHandleApi) startTracingAndLogging(
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

func (s *saldoQueryHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
