package transactionhandler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-apigateway/internal/shared"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/transaction"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/api"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/transaction"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transactionQueryHandleApi struct {
	client pb.TransactionQueryServiceClient

	logger logger.LoggerInterface

	mapper apimapper.TransactionQueryResponseMapper

	trace trace.Tracer

	requestCounter *prometheus.CounterVec

	requestDuration *prometheus.HistogramVec
}

type transactionQueryHandleDeps struct {
	client pb.TransactionQueryServiceClient

	router *echo.Echo

	logger logger.LoggerInterface

	mapper apimapper.TransactionQueryResponseMapper
}

func NewTransactionQueryHandleApi(params *transactionQueryHandleDeps) *transactionQueryHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transaction_query_handler_requests_total",
			Help: "Total number of transaction query requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transaction_query_handler_request_duration_seconds",
			Help:    "Duration of transaction query requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	transactionQueryHandleApi := &transactionQueryHandleApi{
		client:          params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("transaction-query-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerTransaction := params.router.Group("/api/transaction-query")

	routerTransaction.GET("", transactionQueryHandleApi.FindAll)
	routerTransaction.GET("/card-number/:card_number", transactionQueryHandleApi.FindAllTransactionByCardNumber)
	routerTransaction.GET("/:id", transactionQueryHandleApi.FindById)
	routerTransaction.GET("/merchant/:merchant_id", transactionQueryHandleApi.FindByTransactionMerchantId)
	routerTransaction.GET("/active", transactionQueryHandleApi.FindByActiveTransaction)
	routerTransaction.GET("/trashed", transactionQueryHandleApi.FindByTrashedTransaction)

	return transactionQueryHandleApi
}

// @Summary Find all
// @Tags Transaction Query
// @Security Bearer
// @Description Retrieve a list of all transactions
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTransaction "List of transactions"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transaction-query [get]
func (h *transactionQueryHandleApi) FindAll(c echo.Context) error {
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

	req := &pb.FindAllTransactionRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAllTransaction(ctx, req)

	if err != nil {
		logError("failed to retrieve transaction data", err, zap.Error(err))

		return transaction_errors.ErrApiFailedFindAllTransactions(c)
	}

	so := h.mapper.ToApiResponsePaginationTransaction(res)

	logSuccess("success retrieve transaction data", zap.Any("data", so))

	return c.JSON(http.StatusOK, so)
}

// @Summary Find all transactions by card number
// @Tags Transaction Query
// @Security Bearer
// @Description Retrieve a list of transactions for a specific card number
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTransaction "List of transactions"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transaction-query/card-number/{card_number} [get]
func (h *transactionQueryHandleApi) FindAllTransactionByCardNumber(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAllTransactionByCardNumber"
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

	req := &pb.FindAllTransactionCardNumberRequest{
		CardNumber: cardNumber,
		Page:       int32(page),
		PageSize:   int32(pageSize),
		Search:     search,
	}

	res, err := h.client.FindAllTransactionByCardNumber(ctx, req)

	if err != nil {
		logError("failed to retrieve transaction data", err, zap.Error(err))

		return transaction_errors.ErrApiFailedFindByCardNumber(c)
	}

	so := h.mapper.ToApiResponsePaginationTransaction(res)

	logSuccess("success retrieve transaction data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Find a transaction by ID
// @Tags Transaction Query
// @Security Bearer
// @Description Retrieve a transaction record using its ID
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransaction "Transaction data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transaction-query/{id} [get]
func (h *transactionQueryHandleApi) FindById(c echo.Context) error {
	const method = "FindById"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		logError("failed to retrieve transaction data", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidTransactionID(c)
	}

	res, err := h.client.FindByIdTransaction(ctx, &pb.FindByIdTransactionRequest{
		TransactionId: int32(idInt),
	})

	if err != nil {
		logError("failed to retrieve transaction data", err, zap.Error(err))

		return transaction_errors.ErrApiFailedFindById(c)
	}

	so := h.mapper.ToApiResponseTransaction(res)

	logSuccess("success retrieve transaction data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Find transactions by merchant ID
// @Tags Transaction Query
// @Security Bearer
// @Description Retrieve a list of transactions using the merchant ID
// @Accept json
// @Produce json
// @Param merchant_id query string true "Merchant ID"
// @Success 200 {object} response.ApiResponseTransactions "Transaction data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transaction-query/merchant/{merchant_id} [get]
func (h *transactionQueryHandleApi) FindByTransactionMerchantId(c echo.Context) error {
	const method = "FindByTransactionMerchantId"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	merchantId := c.QueryParam("merchant_id")

	merchantIdInt, err := strconv.Atoi(merchantId)

	if err != nil {
		logError("Failed to retrieve transaction data", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidTransactionMerchantID(c)
	}

	req := &pb.FindTransactionByMerchantIdRequest{
		MerchantId: int32(merchantIdInt),
	}

	res, err := h.client.FindTransactionByMerchantId(ctx, req)

	if err != nil {
		logError("Failed to retrieve transaction data", err, zap.Error(err))

		return transaction_errors.ErrApiFailedFindByMerchantID(c)
	}

	so := h.mapper.ToApiResponseTransactions(res)

	logSuccess("success retrieve transaction data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Find active transactions
// @Tags Transaction Query
// @Security Bearer
// @Description Retrieve a list of active transactions
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Number of items per page (default: 10)"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponseTransactions "List of active transactions"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transaction-query/active [get]
func (h *transactionQueryHandleApi) FindByActiveTransaction(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindByActiveTransaction"
	)

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	page := shared.ParseQueryInt(c, "page", defaultPage)
	pageSize := shared.ParseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	req := &pb.FindAllTransactionRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActiveTransaction(ctx, req)

	if err != nil {
		logError("Failed to retrieve transaction data", err, zap.Error(err))

		return transaction_errors.ErrApiFailedFindActive(c)
	}

	so := h.mapper.ToApiResponsePaginationTransactionDeleteAt(res)

	logSuccess("success retrieve transaction data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Retrieve trashed transactions
// @Tags Transaction Query
// @Security Bearer
// @Description Retrieve a list of trashed transactions
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Number of items per page (default: 10)"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponseTransactions "List of trashed transactions"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transaction-query/trashed [get]
func (h *transactionQueryHandleApi) FindByTrashedTransaction(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindByTrashedTransaction"
	)

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	page := shared.ParseQueryInt(c, "page", defaultPage)
	pageSize := shared.ParseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	req := &pb.FindAllTransactionRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashedTransaction(ctx, req)

	if err != nil {
		logError("Failed to retrieve transaction data", err, zap.Error(err))

		return transaction_errors.ErrApiFailedFindTrashed(c)
	}

	so := h.mapper.ToApiResponsePaginationTransactionDeleteAt(res)

	logSuccess("success retrieve transaction data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *transactionQueryHandleApi) startTracingAndLogging(
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

func (s *transactionQueryHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
