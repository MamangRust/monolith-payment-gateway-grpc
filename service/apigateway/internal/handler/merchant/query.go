package merchanthandler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-apigateway/internal/shared"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/merchant"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/api"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/merchant"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// merchantQueryHandleApi is a struct that handles merchant query requests.
type merchantQueryHandleApi struct {
	// The client for making gRPC requests to the merchant query service.
	client pb.MerchantQueryServiceClient

	// The logger for logging messages.
	logger logger.LoggerInterface

	// The mapper for mapping gRPC responses to API responses.
	mapper apimapper.MerchantQueryResponseMapper

	// The tracer for tracing requests.
	trace trace.Tracer

	// The router for handling HTTP requests.
	requestCounter *prometheus.CounterVec

	// The router for handling HTTP requests.
	requestDuration *prometheus.HistogramVec
}

// merchantQueryHandleDeps is a struct that holds the necessary dependencies for the merchantQueryHandleApi.
type merchantQueryHandleDeps struct {
	// The client for making gRPC requests to the merchant query service.
	client pb.MerchantQueryServiceClient

	// The router for handling HTTP requests.
	router *echo.Echo

	// The logger for logging messages.
	logger logger.LoggerInterface

	// The mapper for mapping gRPC responses to API responses.
	mapper apimapper.MerchantQueryResponseMapper
}

// NewMerchantQueryHandleApi initializes a new merchantQueryHandleApi and sets up the routes for merchant query-related operations.
//
// This function registers various HTTP endpoints related to merchant query management, including retrieval of all merchants, a single merchant by ID, and merchants by API key.
// It also collects Prometheus metrics for tracking the number of requests and their durations, helping to monitor the performance and reliability of the handler.
// The routes are grouped under "/api/merchant-query".
//
// Parameters:
// - params: A pointer to merchantQueryHandleDeps containing the necessary dependencies such as router, logger, client, and mapper.
//
// Returns:
// - A pointer to a newly created merchantQueryHandleApi.
func NewMerchantQueryHandleApi(params *merchantQueryHandleDeps) *merchantQueryHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_query_handler_requests_total",
			Help: "Total number of merchant stats amount requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_query_handler_request_duration_seconds",
			Help:    "Duration of merchant query requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	merchantHandler := &merchantQueryHandleApi{
		client:          params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("merchant-query-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerMerchant := params.router.Group("/api/merchant-query")

	routerMerchant.GET("", merchantHandler.FindAll)
	routerMerchant.GET("/:id", merchantHandler.FindById)
	routerMerchant.GET("/api-key", merchantHandler.FindByApiKey)
	routerMerchant.GET("/merchant-user", merchantHandler.FindByMerchantUserId)

	routerMerchant.GET("/active", merchantHandler.FindByActive)
	routerMerchant.GET("/trashed", merchantHandler.FindByTrashed)

	return merchantHandler
}

// FindAll godoc
// @Summary Find all merchants
// @Tags Merchant Query
// @Security Bearer
// @Description Retrieve a list of all merchants
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchant "List of merchants"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-query [get]
func (h *merchantQueryHandleApi) FindAll(c echo.Context) error {
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

	req := &pb.FindAllMerchantRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAllMerchant(ctx, req)

	if err != nil {
		logError("failed to retrieve merchant data", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindAllMerchants(c)
	}

	so := h.mapper.ToApiResponsesMerchant(res)

	logSuccess("merchant data retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)

}

// FindById godoc
// @Summary Find a merchant by ID
// @Tags Merchant Query
// @Security Bearer
// @Description Retrieve a merchant by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchant "Merchant data"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-query/{id} [get]
func (h *merchantQueryHandleApi) FindById(c echo.Context) error {
	const method = "FindById"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logError("failed to retrieve merchant data", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidMerchantID(c)
	}

	req := &pb.FindByIdMerchantRequest{
		MerchantId: int32(id),
	}

	res, err := h.client.FindByIdMerchant(ctx, req)

	if err != nil {
		logError("failed to retrieve merchant data", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindByIdMerchant(c)
	}

	so := h.mapper.ToApiResponseMerchant(res)

	logSuccess("merchant data retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindByApiKey godoc
// @Summary Find a merchant by API key
// @Tags Merchant Query
// @Security Bearer
// @Description Retrieve a merchant by its API key
// @Accept json
// @Produce json
// @Param api_key query string true "API key"
// @Success 200 {object} response.ApiResponseMerchant "Merchant data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-query/api-key [get]
func (h *merchantQueryHandleApi) FindByApiKey(c echo.Context) error {
	const method = "FindByApiKey"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	apiKey := c.QueryParam("api_key")

	if apiKey == "" {
		err := errors.New("api key is empty")
		logError("failed to find merchant by api key", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidApiKey(c)
	}

	req := &pb.FindByApiKeyRequest{
		ApiKey: apiKey,
	}

	res, err := h.client.FindByApiKey(ctx, req)

	if err != nil {
		logError("failed to find merchant by api key", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindByApiKeyMerchant(c)
	}

	so := h.mapper.ToApiResponseMerchant(res)

	logSuccess("merchant retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindByMerchantUserId godoc.
// @Summary Find a merchant by user ID
// @Tags Merchant Query
// @Security Bearer
// @Description Retrieve a merchant by its user ID
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.ApiResponsesMerchant "Merchant data"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-query/merchant-user [get]
func (h *merchantQueryHandleApi) FindByMerchantUserId(c echo.Context) error {
	const method = "FindByMerchantUserId"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id, ok := c.Get("user_id").(int32)

	if !ok {
		err := errors.New("user id not found in context")
		logError("failed to find merchant by user id", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidUserID(c)
	}

	req := &pb.FindByMerchantUserIdRequest{
		UserId: id,
	}

	res, err := h.client.FindByMerchantUserId(ctx, req)

	if err != nil {
		logError("failed to find merchant by user id", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindByUserId(c)
	}

	so := h.mapper.ToApiResponseMerchants(res)

	logSuccess("merchant retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindByActive godoc
// @Summary Find active merchants
// @Tags Merchant Query
// @Security Bearer
// @Description Retrieve a list of active merchants
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsesMerchant "List of active merchants"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-query/active [get]
func (h *merchantQueryHandleApi) FindByActive(c echo.Context) error {
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

	req := &pb.FindAllMerchantRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		logError("failed to find merchant by active", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindAllMerchantsActive(c)
	}

	so := h.mapper.ToApiResponsesMerchantDeleteAt(res)

	logSuccess("merchant retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindByTrashed godoc
// @Summary Find trashed merchants
// @Tags Merchant Query
// @Security Bearer
// @Description Retrieve a list of trashed merchants
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsesMerchant "List of trashed merchants"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-query/trashed [get]
func (h *merchantQueryHandleApi) FindByTrashed(c echo.Context) error {
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

	req := &pb.FindAllMerchantRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		logError("failed to find merchant by trashed", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindAllMerchantsTrashed(c)
	}

	so := h.mapper.ToApiResponsesMerchantDeleteAt(res)

	logSuccess("merchant retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *merchantQueryHandleApi) startTracingAndLogging(
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

func (s *merchantQueryHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
