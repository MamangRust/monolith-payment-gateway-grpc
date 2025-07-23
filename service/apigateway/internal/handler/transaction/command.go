package transactionhandler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-apigateway/internal/middlewares"
	mencache "github.com/MamangRust/monolith-payment-gateway-apigateway/internal/redis"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/transaction"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/api"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/transaction"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type transactionCommandHandleApi struct {
	kafka *kafka.Kafka

	client pb.TransactionCommandServiceClient

	logger logger.LoggerInterface

	mapper apimapper.TransactionCommandResponseMapper

	trace trace.Tracer

	requestCounter *prometheus.CounterVec

	requestDuration *prometheus.HistogramVec
}

type transactionCommandHandleDeps struct {
	kafka *kafka.Kafka

	client pb.TransactionCommandServiceClient

	router *echo.Echo

	logger logger.LoggerInterface

	mapper apimapper.TransactionCommandResponseMapper

	cache mencache.MerchantCache
}

func NewTransactionCommandHandleApi(params *transactionCommandHandleDeps) *transactionCommandHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transaction_command_handler_requests_total",
			Help: "Total number of transaction command requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transaction_command_handler_request_duration_seconds",
			Help:    "Duration of transaction command requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	transactionCommandHandleApi := &transactionCommandHandleApi{
		kafka:           params.kafka,
		client:          params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("transaction-command-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	transactionMiddleware := middlewares.NewApiKeyValidator(params.kafka, "request-transaction", "response-transaction", 5*time.Second, params.logger, params.cache)

	routerTransaction := params.router.Group("/api/transaction-command")

	routerTransaction.POST("/create", transactionMiddleware.Middleware()(transactionCommandHandleApi.Create))
	routerTransaction.POST("/update/:id", transactionMiddleware.Middleware()(transactionCommandHandleApi.Update))

	routerTransaction.POST("/restore/:id", transactionCommandHandleApi.RestoreTransaction)
	routerTransaction.POST("/trashed/:id", transactionCommandHandleApi.TrashedTransaction)
	routerTransaction.DELETE("/permanent/:id", transactionCommandHandleApi.DeletePermanent)

	routerTransaction.POST("/restore/all", transactionCommandHandleApi.RestoreAllTransaction)
	routerTransaction.POST("/permanent/all", transactionCommandHandleApi.DeleteAllTransactionPermanent)

	return transactionCommandHandleApi
}

// @Summary Create a new transaction
// @Tags Transaction Command
// @Security Bearer
// @Description Create a new transaction record with the provided details.
// @Accept json
// @Produce json
// @Param CreateTransactionRequest body requests.CreateTransactionRequest true "Create Transaction Request"
// @Success 200 {object} response.ApiResponseTransaction "Successfully created transaction record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create transaction"
// @Router /api/transaction-command/create [post]
func (h *transactionCommandHandleApi) Create(c echo.Context) error {
	const method = "Create"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	var body requests.CreateTransactionRequest

	apiKeyRaw := c.Get("apiKey")

	if apiKeyRaw == nil {
		err := errors.New("api key not found")

		logError("Failed to create transaction", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidTransactionApiKey(c)
	}

	apiKey, ok := apiKeyRaw.(string)
	if !ok || apiKey == "" {
		err := errors.New("invalid api key")

		logError("Failed to create transaction", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidTransactionApiKey(c)
	}

	if err := c.Bind(&body); err != nil {
		logError("Failed to create transaction", err, zap.Error(err))

		return transaction_errors.ErrApiBindCreateTransaction(c)
	}

	if err := body.Validate(); err != nil {
		logError("Failed to create transaction", err, zap.Error(err))

		return transaction_errors.ErrApiValidateCreateTransaction(c)
	}

	merchantIDRaw := c.Get("merchant_id")
	h.logger.Debug("Merchant ID raw", zap.Any("merchantID", merchantIDRaw), zap.String("type", fmt.Sprintf("%T", merchantIDRaw)))

	var merchantID int
	switch id := merchantIDRaw.(type) {
	case int:
		merchantID = id
	case int32:
		merchantID = int(id)
	case int64:
		merchantID = int(id)
	case float64:
		merchantID = int(id)
	case string:
		parsed, err := strconv.Atoi(id)
		if err != nil {
			h.logger.Error("Failed to parse merchant ID string", zap.String("value", id), zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Invalid merchant ID format")
		}
		merchantID = parsed
	default:
		h.logger.Error("Unknown merchant ID type", zap.Any("merchantID", merchantIDRaw))
		return echo.NewHTTPError(http.StatusInternalServerError, "Unknown merchant ID type")
	}

	h.logger.Debug("Merchant ID parsed", zap.Int("merchant.id", merchantID))

	res, err := h.client.CreateTransaction(ctx, &pb.CreateTransactionRequest{
		ApiKey:          apiKey,
		CardNumber:      body.CardNumber,
		Amount:          int32(body.Amount),
		PaymentMethod:   body.PaymentMethod,
		MerchantId:      int32(merchantID),
		TransactionTime: timestamppb.New(body.TransactionTime),
	})

	if err != nil {
		logError("Failed to create transaction", err, zap.Error(err))

		return transaction_errors.ErrApiFailedCreateTransaction(c)
	}

	so := h.mapper.ToApiResponseTransaction(res)

	logSuccess("success create transaction", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Update a transaction
// @Tags Transaction Command
// @Security Bearer
// @Description Update an existing transaction record using its ID
// @Accept json
// @Produce json
// @Param transaction body requests.UpdateTransactionRequest true "Transaction data"
// @Success 200 {object} response.ApiResponseTransaction "Updated transaction data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update transaction"
// @Router /api/transaction-command/update [post]
func (h *transactionCommandHandleApi) Update(c echo.Context) error {
	const method = "Update"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logError("Failed to update transaction", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidTransactionID(c)
	}

	var body requests.UpdateTransactionRequest

	body.MerchantID = &id

	apiKey, ok := c.Get("apiKey").(string)
	if !ok {
		err := errors.New("invalid api key")

		logError("Failed to update transaction", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidTransactionApiKey(c)
	}

	if err := c.Bind(&body); err != nil {
		logError("Failed to update transaction", err, zap.Error(err))

		return transaction_errors.ErrApiBindUpdateTransaction(c)
	}

	if err := body.Validate(); err != nil {
		logError("Failed to update transaction", err, zap.Error(err))

		return transaction_errors.ErrApiValidateUpdateTransaction(c)
	}

	res, err := h.client.UpdateTransaction(ctx, &pb.UpdateTransactionRequest{
		TransactionId:   int32(id),
		CardNumber:      body.CardNumber,
		ApiKey:          apiKey,
		Amount:          int32(body.Amount),
		PaymentMethod:   body.PaymentMethod,
		MerchantId:      int32(*body.MerchantID),
		TransactionTime: timestamppb.New(body.TransactionTime),
	})

	if err != nil {
		logError("Failed to update transaction", err, zap.Error(err))

		return transaction_errors.ErrApiFailedUpdateTransaction(c)
	}

	so := h.mapper.ToApiResponseTransaction(res)

	logSuccess("success update transaction", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Trash a transaction
// @Tags Transaction Command
// @Security Bearer
// @Description Trash a transaction record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransaction "Successfully trashed transaction record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to trashed transaction"
// @Router /api/transaction-command/trashed/{id} [post]
func (h *transactionCommandHandleApi) TrashedTransaction(c echo.Context) error {
	const method = "TrashedTransaction"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		logError("Bad Request: Invalid ID", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidTransactionID(c)
	}

	res, err := h.client.TrashedTransaction(ctx, &pb.FindByIdTransactionRequest{
		TransactionId: int32(idInt),
	})

	if err != nil {
		logError("Failed to trashed transaction", err, zap.Error(err))

		return transaction_errors.ErrApiFailedTrashTransaction(c)
	}

	so := h.mapper.ToApiResponseTransactionDeleteAt(res)

	logSuccess("success trashed transaction", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Restore a trashed transaction
// @Tags Transaction Command
// @Security Bearer
// @Description Restore a trashed transaction record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransaction "Successfully restored transaction record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore transaction:"
// @Router /api/transaction-command/restore/{id} [post]
func (h *transactionCommandHandleApi) RestoreTransaction(c echo.Context) error {
	const method = "RestoreTransaction"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		logError("Bad Request: Invalid ID", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidTransactionID(c)
	}

	res, err := h.client.RestoreTransaction(ctx, &pb.FindByIdTransactionRequest{
		TransactionId: int32(idInt),
	})

	if err != nil {
		logError("Failed to restore transaction", err, zap.Error(err))

		return transaction_errors.ErrApiFailedRestoreTransaction(c)
	}

	so := h.mapper.ToApiResponseTransaction(res)

	logSuccess("success restore transaction", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Permanently delete a transaction
// @Tags Transaction Command
// @Security Bearer
// @Description Permanently delete a transaction record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransactionDelete "Successfully deleted transaction record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete transaction:"
// @Router /api/transaction-command/permanent/{id} [delete]
func (h *transactionCommandHandleApi) DeletePermanent(c echo.Context) error {
	const method = "DeletePermanent"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		logError("Bad Request: Invalid ID", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidTransactionID(c)
	}

	res, err := h.client.DeleteTransactionPermanent(ctx, &pb.FindByIdTransactionRequest{
		TransactionId: int32(idInt),
	})

	if err != nil {
		logError("Failed to delete transaction", err, zap.Error(err))

		return transaction_errors.ErrApiFailedDeletePermanent(c)
	}

	so := h.mapper.ToApiResponseTransactionDelete(res)

	logSuccess("success delete transaction", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Restore a trashed transaction
// @Tags Transaction Command
// @Security Bearer
// @Description Restore a trashed transaction all.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseTransactionAll "Successfully restored transaction record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore transaction:"
// @Router /api/transaction-command/restore/all [post]
func (h *transactionCommandHandleApi) RestoreAllTransaction(c echo.Context) error {
	const method = "RestoreAllTransaction"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	res, err := h.client.RestoreAllTransaction(ctx, &emptypb.Empty{})

	if err != nil {
		logError("Failed to restore all transaction", err, zap.Error(err))

		return transaction_errors.ErrApiFailedRestoreAllTransactions(c)
	}

	so := h.mapper.ToApiResponseTransactionAll(res)

	logSuccess("success restore all transaction", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Permanently delete a transaction
// @Tags Transaction Command
// @Security Bearer
// @Description Permanently delete a transaction all.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseTransactionAll "Successfully deleted transaction record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete transaction:"
// @Router /api/transaction-command/delete/all [post]
func (h *transactionCommandHandleApi) DeleteAllTransactionPermanent(c echo.Context) error {
	const method = "DeleteAllTransactionPermanent"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	res, err := h.client.DeleteAllTransactionPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		logError("Failed to delete all transaction permanently", err, zap.Error(err))

		return transaction_errors.ErrApiFailedDeleteAllPermanent(c)
	}

	so := h.mapper.ToApiResponseTransactionAll(res)

	logSuccess("success delete all transaction", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *transactionCommandHandleApi) startTracingAndLogging(
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

func (s *transactionCommandHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
