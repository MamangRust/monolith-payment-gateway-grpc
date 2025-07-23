package transferhandler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-apigateway/internal/shared"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/transfer"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	transfer_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors/api"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/transfer"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transferQueryHandleApi struct {
	client pb.TransferQueryServiceClient

	logger logger.LoggerInterface

	mapper apimapper.TransferQueryResponseMapper

	trace trace.Tracer

	requestCounter *prometheus.CounterVec

	requestDuration *prometheus.HistogramVec
}

type transferQueryHandleDeps struct {
	client pb.TransferQueryServiceClient

	router *echo.Echo

	logger logger.LoggerInterface

	mapper apimapper.TransferQueryResponseMapper
}

func NewTransferQueryHandleApi(params *transferQueryHandleDeps) *transferQueryHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transfer_query_handler_requests_total",
			Help: "Total number of transfer query requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transfer_query_handler_request_duration_seconds",
			Help:    "Duration of transfer query requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	transferQueryHandleApi := &transferQueryHandleApi{
		client:          params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("transfer-query-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerTransfer := params.router.Group("/api/transfer-query")

	routerTransfer.GET("", transferQueryHandleApi.FindAll)
	routerTransfer.GET("/:id", transferQueryHandleApi.FindById)

	routerTransfer.GET("/transfer_from/:transfer_from", transferQueryHandleApi.FindByTransferByTransferFrom)
	routerTransfer.GET("/transfer_to/:transfer_to", transferQueryHandleApi.FindByTransferByTransferTo)

	routerTransfer.GET("/active", transferQueryHandleApi.FindByActiveTransfer)
	routerTransfer.GET("/trashed", transferQueryHandleApi.FindByTrashedTransfer)

	return transferQueryHandleApi
}

// @Summary Find all transfer records
// @Tags Transfer Query
// @Security Bearer
// @Description Retrieve a list of all transfer records with pagination
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTransfer "List of transfer records"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transfer data"
// @Router /api/transfer-query [get]
func (h *transferQueryHandleApi) FindAll(c echo.Context) error {
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

	req := &pb.FindAllTransferRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAllTransfer(ctx, req)

	if err != nil {
		logError("Failed to retrieve transfer data", err, zap.Error(err))

		return transfer_errors.ErrApiFailedFindAllTransfers(c)
	}

	so := h.mapper.ToApiResponsePaginationTransfer(res)

	logSuccess("Successfully retrieved transfer data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Find a transfer by ID
// @Tags Transfer Query
// @Security Bearer
// @Description Retrieve a transfer record using its ID
// @Accept json
// @Produce json
// @Param id path string true "Transfer ID"
// @Success 200 {object} response.ApiResponseTransfer "Transfer data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transfer data"
// @Router /api/transfer-query/{id} [get]
func (h *transferQueryHandleApi) FindById(c echo.Context) error {
	const method = "FindById"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		logError("Failed to retrieve transfer data", err, zap.Error(err))

		return transfer_errors.ErrApiTransferInvalidID(c)

	}

	res, err := h.client.FindByIdTransfer(ctx, &pb.FindByIdTransferRequest{
		TransferId: int32(idInt),
	})

	if err != nil {
		logError("Failed to retrieve transfer data", err, zap.Error(err))

		return transfer_errors.ErrApiFailedFindByIdTransfer(c)
	}

	so := h.mapper.ToApiResponseTransfer(res)

	logSuccess("Successfully retrieved transfer data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Find transfers by transfer_from
// @Tags Transfer Query
// @Security Bearer
// @Description Retrieve a list of transfer records using the transfer_from parameter
// @Accept json
// @Produce json
// @Param transfer_from path string true "Transfer From"
// @Success 200 {object} response.ApiResponseTransfers "Transfer data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transfer data"
// @Router /api/transfer-query/transfer_from/{transfer_from} [get]
func (h *transferQueryHandleApi) FindByTransferByTransferFrom(c echo.Context) error {
	const method = "FindByTransferByTransferFrom"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	transfer_from := c.Param("transfer_from")

	if transfer_from == "" {
		err := errors.New("transfer_from is required")

		logError("Failed to retrieve transfer data", err, zap.Error(err))

		return transfer_errors.ErrApiInvalidCardNumber(c)
	}

	res, err := h.client.FindTransferByTransferFrom(ctx, &pb.FindTransferByTransferFromRequest{
		TransferFrom: transfer_from,
	})

	if err != nil {
		logError("Failed to retrieve transfer data", err, zap.Error(err))

		return transfer_errors.ErrApiFailedFindByTransferFrom(c)
	}

	so := h.mapper.ToApiResponseTransfers(res)

	logSuccess("Successfully retrieved transfer data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Find transfers by transfer_to
// @Tags Transfer Query
// @Security Bearer
// @Description Retrieve a list of transfer records using the transfer_to parameter
// @Accept json
// @Produce json
// @Param transfer_to path string true "Transfer To"
// @Success 200 {object} response.ApiResponseTransfers "Transfer data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transfer data"
// @Router /api/transfer-query/transfer_to/{transfer_to} [get]
func (h *transferQueryHandleApi) FindByTransferByTransferTo(c echo.Context) error {
	const method = "FindByTransferByTransferTo"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	transfer_to := c.Param("transfer_to")

	if transfer_to == "" {
		err := errors.New("transfer_to is required")

		logError("Failed to retrieve transfer data", err, zap.Error(err))

		return transfer_errors.ErrApiInvalidCardNumber(c)
	}

	res, err := h.client.FindTransferByTransferTo(ctx, &pb.FindTransferByTransferToRequest{
		TransferTo: transfer_to,
	})

	if err != nil {
		logError("Failed to retrieve transfer data", err, zap.Error(err))

		return transfer_errors.ErrApiFailedFindByTransferTo(c)
	}

	so := h.mapper.ToApiResponseTransfers(res)

	logSuccess("Successfully retrieved transfer data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Find active transfers
// @Tags Transfer Query
// @Security Bearer
// @Description Retrieve a list of active transfer records
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Number of items per page (default: 10)"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponseTransfers "Active transfer data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transfer data"
// @Router /api/transfer-query/active [get]
func (h *transferQueryHandleApi) FindByActiveTransfer(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindByActiveTransfer"
	)

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	page := shared.ParseQueryInt(c, "page", defaultPage)
	pageSize := shared.ParseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	req := &pb.FindAllTransferRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActiveTransfer(ctx, req)

	if err != nil {
		logError("Failed to retrieve transfer data", err, zap.Error(err))

		return transfer_errors.ErrApiFailedFindByActiveTransfer(c)
	}

	so := h.mapper.ToApiResponsePaginationTransferDeleteAt(res)

	logSuccess("Successfully retrieved transfer data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Retrieve trashed transfers
// @Tags Transfer Query
// @Security Bearer
// @Description Retrieve a list of trashed transfer records
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Number of items per page (default: 10)"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponseTransfers "List of trashed transfer records"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transfer data"
// @Router /api/transfer-query/trashed [get]
func (h *transferQueryHandleApi) FindByTrashedTransfer(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindByTrashedTransfer"
	)

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	page := shared.ParseQueryInt(c, "page", defaultPage)
	pageSize := shared.ParseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	req := &pb.FindAllTransferRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashedTransfer(ctx, req)

	if err != nil {
		logError("Failed to retrieve transfer data", err, zap.Error(err))

		return transfer_errors.ErrApiFailedFindByTrashedTransfer(c)
	}

	so := h.mapper.ToApiResponsePaginationTransferDeleteAt(res)

	logSuccess("Successfully retrieved transfer data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *transferQueryHandleApi) startTracingAndLogging(
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

func (s *transferQueryHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
