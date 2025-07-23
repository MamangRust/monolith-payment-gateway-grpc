package merchanthandler

import (
	"context"
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

type merchantTransactionHandleApi struct {
	client pb.MerchantTransactionServiceClient

	logger logger.LoggerInterface

	mapper apimapper.MerchantTransactionResponseMapper

	trace trace.Tracer

	requestCounter *prometheus.CounterVec

	requestDuration *prometheus.HistogramVec
}

type merchantTransactionHandleDeps struct {
	client pb.MerchantTransactionServiceClient

	router *echo.Echo

	logger logger.LoggerInterface

	mapper apimapper.MerchantTransactionResponseMapper
}

func NewMerchantTransactionHandleApi(params *merchantTransactionHandleDeps) *merchantTransactionHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_transactions_handler_requests_total",
			Help: "Total number of merchant transactions requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_transactions_handler_request_duration_seconds",
			Help:    "Duration of merchant transactions requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	merchantHandler := &merchantTransactionHandleApi{
		client:          params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("merchant-transactions-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerMerchant := params.router.Group("/api/merchant-transactions")

	routerMerchant.GET("/transactions", merchantHandler.FindAllTransactions)
	routerMerchant.GET("/transactions/:merchant_id", merchantHandler.FindAllTransactionByMerchant)
	routerMerchant.GET("/transactions/api-key/:api_key", merchantHandler.FindAllTransactionByApikey)
	return merchantHandler
}

// FindAllTransactions godoc
// @Summary Find all transactions
// @Tags Merchant
// @Security Bearer
// @Description Retrieve a list of all transactions
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTransaction "List of transactions"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/merchants/transaction [get]
func (h *merchantTransactionHandleApi) FindAllTransactions(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAllTransactions"
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

	res, err := h.client.FindAllTransactionMerchant(ctx, req)

	if err != nil {
		logError("failed to retrieve transaction data", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindAllTransactions(c)
	}

	so := h.mapper.ToApiResponseMerchantsTransactionResponse(res)

	logSuccess("transaction data retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindAllTransactionByMerchant godoc
// @Summary Find all transactions by merchant ID
// @Tags Merchant
// @Security Bearer
// @Description Retrieve a list of transactions for a specific merchant
// @Accept json
// @Produce json
// @Param merchant_id path int true "Merchant ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTransaction "List of transactions"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/merchant-transactions/:merchant_id [get]
func (h *merchantTransactionHandleApi) FindAllTransactionByMerchant(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAllTransactionByMerchant"
	)

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	merchantID, err := strconv.Atoi(c.Param("merchant_id"))

	if err != nil || merchantID <= 0 {
		logError("failed to retrieve transaction data", err, zap.Error(err))
		return merchant_errors.ErrApiInvalidMerchantID(c)
	}

	page := shared.ParseQueryInt(c, "page", defaultPage)
	pageSize := shared.ParseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	req := &pb.FindAllMerchantTransaction{
		MerchantId: int32(merchantID),
		Page:       int32(page),
		PageSize:   int32(pageSize),
		Search:     search,
	}

	res, err := h.client.FindAllTransactionByMerchant(ctx, req)

	if err != nil {
		logError("failed to retrieve transaction data", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindAllTransactionByMerchant(c)
	}

	so := h.mapper.ToApiResponseMerchantsTransactionResponse(res)

	logSuccess("transaction data retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindAllTransactionByApikey godoc
// @Summary Find all transactions by api_key
// @Tags Merchant
// @Security Bearer
// @Description Retrieve a list of transactions for a specific merchant
// @Accept json
// @Produce json
// @Param api_key path string true "Api key"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTransaction "List of transactions"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/merchant-transactions/api-key/:api_key [get]
func (h *merchantTransactionHandleApi) FindAllTransactionByApikey(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAllTransactionByApikey"
	)

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	api_key := c.Param("api_key")

	page := shared.ParseQueryInt(c, "page", defaultPage)
	pageSize := shared.ParseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	req := &pb.FindAllMerchantApikey{
		ApiKey:   api_key,
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAllTransactionByApikey(ctx, req)

	if err != nil {
		logError("failed to find all transaction by api key", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindAllTransactionByApikey(c)
	}

	so := h.mapper.ToApiResponseMerchantsTransactionResponse(res)

	logSuccess("transaction retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *merchantTransactionHandleApi) startTracingAndLogging(
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

func (s *merchantTransactionHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
