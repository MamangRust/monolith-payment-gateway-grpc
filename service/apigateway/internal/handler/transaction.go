package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-apigateway/internal/middlewares"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api"
	"github.com/MamangRust/monolith-payment-gateway-shared/pb"
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

type transactionHandler struct {
	kafka           *kafka.Kafka
	transaction     pb.TransactionServiceClient
	logger          logger.LoggerInterface
	mapping         apimapper.TransactionResponseMapper
	trace           trace.Tracer
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewHandlerTransaction(transaction pb.TransactionServiceClient, merchant pb.MerchantServiceClient, router *echo.Echo, logger logger.LoggerInterface, mapping apimapper.TransactionResponseMapper, kafka *kafka.Kafka) *transactionHandler {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transaction_handler_requests_total",
			Help: "Total number of transaction requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transaction_handler_request_duration_seconds",
			Help:    "Duration of transaction requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	transactionHandler := transactionHandler{
		transaction:     transaction,
		logger:          logger,
		mapping:         mapping,
		kafka:           kafka,
		trace:           otel.Tracer("transaction-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	transactionMiddleware := middlewares.NewApiKeyValidator(kafka, "request-transaction", "response-transaction", 5*time.Second)

	routerTransaction := router.Group("/api/transactions")

	routerTransaction.GET("", transactionHandler.FindAll)
	routerTransaction.GET("/card-number/:card_number", transactionHandler.FindAllTransactionByCardNumber)

	routerTransaction.GET("/:id", transactionHandler.FindById)

	routerTransaction.GET("/monthly-success", transactionHandler.FindMonthlyTransactionStatusSuccess)
	routerTransaction.GET("/yearly-success", transactionHandler.FindYearlyTransactionStatusSuccess)
	routerTransaction.GET("/monthly-failed", transactionHandler.FindMonthlyTransactionStatusFailed)
	routerTransaction.GET("/yearly-failed", transactionHandler.FindYearlyTransactionStatusFailed)

	routerTransaction.GET("/monthly-success-by-card", transactionHandler.FindMonthlyTransactionStatusSuccessByCardNumber)
	routerTransaction.GET("/yearly-success-by-card", transactionHandler.FindYearlyTransactionStatusSuccessByCardNumber)
	routerTransaction.GET("/monthly-failed-by-card", transactionHandler.FindMonthlyTransactionStatusFailedByCardNumber)
	routerTransaction.GET("/yearly-failed-by-card", transactionHandler.FindYearlyTransactionStatusFailedByCardNumber)

	routerTransaction.GET("/monthly-methods", transactionHandler.FindMonthlyPaymentMethods)
	routerTransaction.GET("/yearly-methods", transactionHandler.FindYearlyPaymentMethods)
	routerTransaction.GET("/monthly-amounts", transactionHandler.FindMonthlyAmounts)
	routerTransaction.GET("/yearly-amounts", transactionHandler.FindYearlyAmounts)

	routerTransaction.GET("/monthly-methods-by-card", transactionHandler.FindMonthlyPaymentMethodsByCardNumber)
	routerTransaction.GET("/yearly-methods-by-card", transactionHandler.FindYearlyPaymentMethodsByCardNumber)
	routerTransaction.GET("/monthly-amounts-by-card", transactionHandler.FindMonthlyAmountsByCardNumber)
	routerTransaction.GET("/yearly-amounts-by-card", transactionHandler.FindYearlyAmountsByCardNumber)

	routerTransaction.GET("/merchant/:merchant_id", transactionHandler.FindByTransactionMerchantId)
	routerTransaction.GET("/active", transactionHandler.FindByActiveTransaction)
	routerTransaction.GET("/trashed", transactionHandler.FindByTrashedTransaction)
	routerTransaction.POST("/create", transactionMiddleware.Middleware()(transactionHandler.Create))
	routerTransaction.POST("/update/:id", transactionMiddleware.Middleware()(transactionHandler.Update))

	routerTransaction.POST("/restore/:id", transactionHandler.RestoreTransaction)
	routerTransaction.POST("/trashed/:id", transactionHandler.TrashedTransaction)
	routerTransaction.DELETE("/permanent/:id", transactionHandler.DeletePermanent)

	routerTransaction.POST("/restore/all", transactionHandler.RestoreAllTransaction)
	routerTransaction.POST("/permanent/all", transactionHandler.DeleteAllTransactionPermanent)

	return &transactionHandler
}

// @Summary Find all
// @Tags Transaction
// @Security Bearer
// @Description Retrieve a list of all transactions
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTransaction "List of transactions"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transactions [get]
func (h *transactionHandler) FindAll(c echo.Context) error {
	const method = "FindAll"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")

	req := &pb.FindAllTransactionRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.transaction.FindAllTransaction(ctx, req)

	if err != nil {
		status = "error"

		logError("failed to retrieve transaction data", err, zap.Error(err))

		return transaction_errors.ErrApiFailedFindAllTransactions(c)
	}

	so := h.mapping.ToApiResponsePaginationTransaction(res)

	logSuccess("success retrieve transaction data", zap.Any("data", so))

	return c.JSON(http.StatusOK, so)
}

// @Summary Find all transactions by card number
// @Tags Transaction
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
// @Router /api/transactions/card-number/{card_number} [get]
func (h *transactionHandler) FindAllTransactionByCardNumber(c echo.Context) error {
	const method = "FindAllTransactionByCardNumber"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	cardNumber := c.Param("card_number")
	if cardNumber == "" {
		status = "error"
		err := errors.New("card number is required")

		logError("failed to retrieve transaction data", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidTransactionCardNumber(c)
	}

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")

	req := &pb.FindAllTransactionCardNumberRequest{
		CardNumber: cardNumber,
		Page:       int32(page),
		PageSize:   int32(pageSize),
		Search:     search,
	}

	res, err := h.transaction.FindAllTransactionByCardNumber(ctx, req)

	if err != nil {
		status = "error"

		logError("failed to retrieve transaction data", err, zap.Error(err))

		return transaction_errors.ErrApiFailedFindByCardNumber(c)
	}

	so := h.mapping.ToApiResponsePaginationTransaction(res)

	logSuccess("success retrieve transaction data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Find a transaction by ID
// @Tags Transaction
// @Security Bearer
// @Description Retrieve a transaction record using its ID
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransaction "Transaction data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transactions/{id} [get]
func (h *transactionHandler) FindById(c echo.Context) error {
	const method = "FindById"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		status = "error"

		logError("failed to retrieve transaction data", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidTransactionID(c)
	}

	res, err := h.transaction.FindByIdTransaction(ctx, &pb.FindByIdTransactionRequest{
		TransactionId: int32(idInt),
	})

	if err != nil {
		status = "error"

		logError("failed to retrieve transaction data", err, zap.Error(err))

		return transaction_errors.ErrApiFailedFindById(c)
	}

	so := h.mapping.ToApiResponseTransaction(res)

	logSuccess("success retrieve transaction data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTransactionStatusSuccess retrieves the monthly transaction status for successful transactions.
// @Summary Get monthly transaction status for successful transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the monthly transaction status for successful transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseTransactionMonthStatusSuccess "Monthly transaction status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction status for successful transactions"
// @Router /api/transactions/monthly-success [get]
func (h *transactionHandler) FindMonthlyTransactionStatusSuccess(c echo.Context) error {
	const method = "FindMonthlyTransactionStatusSuccess"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		status = "error"

		logError("failed to retrieve monthly Transaction status success", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidYear(c)
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		status = "error"

		logError("failed to retrieve monthly Transaction status success", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidMonth(c)
	}

	res, err := h.transaction.FindMonthlyTransactionStatusSuccess(ctx, &pb.FindMonthlyTransactionStatus{
		Year:  int32(year),
		Month: int32(month),
	})

	if err != nil {
		status = "error"

		logError("failed to retrieve monthly Transaction status success", err, zap.Error(err))

		return transaction_errors.ErrApiFailedMonthlySuccess(c)
	}

	so := h.mapping.ToApiResponseTransactionMonthStatusSuccess(res)

	logSuccess("success retrieve monthly Transaction status success", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTransactionStatusSuccess retrieves the yearly transaction status for successful transactions.
// @Summary Get yearly transaction status for successful transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the yearly transaction status for successful transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionYearStatusSuccess "Yearly transaction status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction status for successful transactions"
// @Router /api/transactions/yearly-success [get]
func (h *transactionHandler) FindYearlyTransactionStatusSuccess(c echo.Context) error {
	const method = "FindYearlyTransactionStatusSuccess"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		status = "error"

		logError("failed to retrieve yearly Transaction status success", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidYear(c)
	}

	res, err := h.transaction.FindYearlyTransactionStatusSuccess(ctx, &pb.FindYearTransactionStatus{
		Year: int32(year),
	})

	if err != nil {
		status = "error"

		logError("failed to retrieve yearly Transaction status success", err, zap.Error(err))

		return transaction_errors.ErrApiFailedYearlySuccess(c)
	}

	so := h.mapping.ToApiResponseTransactionYearStatusSuccess(res)

	logSuccess("success retrieve yearly Transaction status success", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTransactionStatusFailed retrieves the monthly transaction status for failed transactions.
// @Summary Get monthly transaction status for failed transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the monthly transaction status for failed transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseTransactionMonthStatusFailed "Monthly transaction status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction status for failed transactions"
// @Router /api/transactions/monthly-failed [get]
func (h *transactionHandler) FindMonthlyTransactionStatusFailed(c echo.Context) error {
	const method = "FindMonthlyTransactionStatusFailed"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		status = "error"

		logError("failed to retrieve monthly Transaction status failed", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidYear(c)
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		status = "error"

		logError("failed to retrieve monthly Transaction status failed", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidMonth(c)
	}

	res, err := h.transaction.FindMonthlyTransactionStatusFailed(ctx, &pb.FindMonthlyTransactionStatus{
		Year:  int32(year),
		Month: int32(month),
	})

	if err != nil {
		status = "error"

		logError("failed to retrieve monthly Transaction status failed", err, zap.Error(err))

		return transaction_errors.ErrApiFailedMonthlyFailed(c)
	}

	so := h.mapping.ToApiResponseTransactionMonthStatusFailed(res)

	logSuccess("success retrieve monthly Transaction status failed", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTransactionStatusFailed retrieves the yearly transaction status for failed transactions.
// @Summary Get yearly transaction status for failed transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the yearly transaction status for failed transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionYearStatusFailed "Yearly transaction status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction status for failed transactions"
// @Router /api/transactions/yearly-failed [get]
func (h *transactionHandler) FindYearlyTransactionStatusFailed(c echo.Context) error {
	const method = "FindYearlyTransactionStatusFailed"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		status = "error"

		logError("failed to retrieve yearly Transaction status failed", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidYear(c)
	}

	res, err := h.transaction.FindYearlyTransactionStatusFailed(ctx, &pb.FindYearTransactionStatus{
		Year: int32(year),
	})

	if err != nil {
		status = "error"

		logError("failed to retrieve yearly Transaction status failed", err, zap.Error(err))

		return transaction_errors.ErrApiFailedYearlyFailed(c)
	}

	so := h.mapping.ToApiResponseTransactionYearStatusFailed(res)

	logSuccess("success retrieve yearly Transaction status failed", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTransactionStatusSuccess retrieves the monthly transaction status for successful transactions.
// @Summary Get monthly transaction status for successful transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the monthly transaction status for successful transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseTransactionMonthStatusSuccess "Monthly transaction status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction status for successful transactions"
// @Router /api/transactions/monthly-success-by-card [get]
func (h *transactionHandler) FindMonthlyTransactionStatusSuccessByCardNumber(c echo.Context) error {
	const method = "FindMonthlyTransactionStatusSuccessByCardNumber"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")
	cardNumber := c.QueryParam("card_number")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		status = "error"

		logError("failed to retrieve monthly Transaction status success", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidYear(c)
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		status = "error"

		logError("failed to retrieve monthly Transaction status success", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidMonth(c)
	}

	res, err := h.transaction.FindMonthlyTransactionStatusSuccessByCardNumber(ctx, &pb.FindMonthlyTransactionStatusCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
		Month:      int32(month),
	})

	if err != nil {
		status = "error"

		logError("failed to retrieve monthly Transaction status success", err, zap.Error(err))

		return transaction_errors.ErrApiFailedMonthlySuccess(c)
	}

	so := h.mapping.ToApiResponseTransactionMonthStatusSuccess(res)

	logSuccess("success retrieve monthly Transaction status success", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTransactionStatusSuccess retrieves the yearly transaction status for successful transactions.
// @Summary Get yearly transaction status for successful transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the yearly transaction status for successful transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param cardNumber query string true "Card Number"
// @Success 200 {object} response.ApiResponseTransactionYearStatusSuccess "Yearly transaction status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction status for successful transactions"
// @Router /api/transactions/yearly-success-by-card [get]
func (h *transactionHandler) FindYearlyTransactionStatusSuccessByCardNumber(c echo.Context) error {
	const method = "FindYearlyTransactionStatusSuccessByCardNumber"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")
	cardNumber := c.QueryParam("cardNumber")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		status = "error"

		logError("failed to retrieve yearly Transaction status success", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidYear(c)
	}

	res, err := h.transaction.FindYearlyTransactionStatusSuccessByCardNumber(ctx, &pb.FindYearTransactionStatusCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
	})

	if err != nil {
		status = "error"

		logError("failed to retrieve yearly Transaction status success", err, zap.Error(err))

		return transaction_errors.ErrApiFailedYearlySuccess(c)
	}

	so := h.mapping.ToApiResponseTransactionYearStatusSuccess(res)

	logSuccess("success retrieve yearly Transaction status success", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTransactionStatusFailed retrieves the monthly transaction status for failed transactions.
// @Summary Get monthly transaction status for failed transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the monthly transaction status for failed transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Param cardNumber query string true "Card Number"
// @Success 200 {object} response.ApiResponseTransactionMonthStatusFailed "Monthly transaction status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction status for failed transactions"
// @Router /api/transactions/monthly-failed-by-card [get]
func (h *transactionHandler) FindMonthlyTransactionStatusFailedByCardNumber(c echo.Context) error {
	const method = "FindMonthlyTransactionStatusFailedByCardNumber"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")
	cardNumber := c.QueryParam("cardNumber")

	if cardNumber == "" {
		status = "error"

		err := errors.New("card number is required")

		logError("failed to retrieve monthly Transaction status failed", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidTransactionCardNumber(c)
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		status = "error"

		logError("failed to retrieve monthly Transaction status failed", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidYear(c)
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		status = "error"

		logError("failed to retrieve monthly Transaction status failed", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidMonth(c)
	}

	res, err := h.transaction.FindMonthlyTransactionStatusFailedByCardNumber(ctx, &pb.FindMonthlyTransactionStatusCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
		Month:      int32(month),
	})

	if err != nil {
		status = "error"

		logError("failed to retrieve monthly Transaction status failed", err, zap.Error(err))

		return transaction_errors.ErrApiFailedMonthlyFailed(c)
	}

	so := h.mapping.ToApiResponseTransactionMonthStatusFailed(res)

	logSuccess("success retrieve monthly Transaction status failed", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTransactionStatusFailedByCardNumber retrieves the yearly transaction status for failed transactions.
// @Summary Get yearly transaction status for failed transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the yearly transaction status for failed transactions by year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionYearStatusFailed "Yearly transaction status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction status for failed transactions"
// @Router /api/transactions/yearly-failed-by-card [get]
func (h *transactionHandler) FindYearlyTransactionStatusFailedByCardNumber(c echo.Context) error {
	const method = "FindYearlyTransactionStatusFailedByCardNumber"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")
	card_number := c.QueryParam("card_number")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		status = "error"

		logError("failed to retrieve yearly Transaction status failed", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidYear(c)
	}

	res, err := h.transaction.FindYearlyTransactionStatusFailedByCardNumber(ctx, &pb.FindYearTransactionStatusCardNumber{
		CardNumber: card_number,
		Year:       int32(year),
	})

	if err != nil {
		status = "error"

		logError("failed to retrieve yearly Transaction status failed", err, zap.Error(err))

		return transaction_errors.ErrApiFailedYearlyFailed(c)
	}

	so := h.mapping.ToApiResponseTransactionYearStatusFailed(res)

	logSuccess("success retrieve yearly Transaction status failed", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyPaymentMethods retrieves the monthly payment methods for transactions.
// @Summary Get monthly payment methods
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the monthly payment methods for transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionMonthMethod "Monthly payment methods"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly payment methods"
// @Router /api/transactions/monthly-payment-methods [get]
func (h *transactionHandler) FindMonthlyPaymentMethods(c echo.Context) error {
	const method = "FindMonthlyPaymentMethods"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		status = "error"

		logError("failed to retrieve monthly payment methods", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidYear(c)
	}

	res, err := h.transaction.FindMonthlyPaymentMethods(ctx, &pb.FindYearTransactionStatus{
		Year: int32(year),
	})
	if err != nil {
		status = "error"

		logError("failed to retrieve monthly payment methods", err, zap.Error(err))

		return transaction_errors.ErrApiFailedMonthlyMethods(c)
	}

	so := h.mapping.ToApiResponseTransactionMonthMethod(res)

	logSuccess("success retrieve monthly payment methods", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyPaymentMethods retrieves the yearly payment methods for transactions.
// @Summary Get yearly payment methods
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the yearly payment methods for transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionYearMethod "Yearly payment methods"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly payment methods"
// @Router /api/transactions/yearly-payment-methods [get]
func (h *transactionHandler) FindYearlyPaymentMethods(c echo.Context) error {
	const method = "FindYearlyPaymentMethods"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		status = "error"

		logError("failed to retrieve yearly payment methods", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidYear(c)
	}

	res, err := h.transaction.FindYearlyPaymentMethods(ctx, &pb.FindYearTransactionStatus{
		Year: int32(year),
	})
	if err != nil {
		status = "error"

		logError("failed to retrieve yearly payment methods", err, zap.Error(err))

		return transaction_errors.ErrApiFailedYearlyMethods(c)
	}

	so := h.mapping.ToApiResponseTransactionYearMethod(res)

	logSuccess("success retrieve yearly payment methods", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyAmounts retrieves the monthly transaction amounts for a specific year.
// @Summary Get monthly transaction amounts
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the monthly transaction amounts for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionMonthAmount "Monthly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction amounts"
// @Router /api/transactions/monthly-amounts [get]
func (h *transactionHandler) FindMonthlyAmounts(c echo.Context) error {
	const method = "FindMonthlyAmounts"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		status = "error"

		logError("Failed to retrieve monthly amounts", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidYear(c)
	}

	res, err := h.transaction.FindMonthlyAmounts(ctx, &pb.FindYearTransactionStatus{
		Year: int32(year),
	})
	if err != nil {
		status = "error"

		logError("Failed to retrieve monthly amounts", err, zap.Error(err))

		return transaction_errors.ErrApiFailedMonthlyAmounts(c)
	}

	so := h.mapping.ToApiResponseTransactionMonthAmount(res)

	logSuccess("success retrieve monthly amounts", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyAmounts retrieves the yearly transaction amounts for a specific year.
// @Summary Get yearly transaction amounts
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the yearly transaction amounts for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionYearAmount "Yearly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction amounts"
// @Router /api/transactions/yearly-amounts [get]
func (h *transactionHandler) FindYearlyAmounts(c echo.Context) error {
	const method = "FindYearlyAmounts"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		status = "error"

		logError("Failed to retrieve yearly amounts", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidYear(c)
	}

	res, err := h.transaction.FindYearlyAmounts(ctx, &pb.FindYearTransactionStatus{
		Year: int32(year),
	})
	if err != nil {
		status = "error"

		logError("Failed to retrieve yearly amounts", err, zap.Error(err))

		return transaction_errors.ErrApiFailedYearlyAmounts(c)
	}

	so := h.mapping.ToApiResponseTransactionYearAmount(res)

	logSuccess("success retrieve yearly amounts", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyPaymentMethodsByCardNumber retrieves the monthly payment methods for transactions by card number and year.
// @Summary Get monthly payment methods by card number
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the monthly payment methods for transactions by card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionMonthMethod "Monthly payment methods by card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly payment methods by card number"
// @Router /api/transactions/monthly-payment-methods-by-card [get]
func (h *transactionHandler) FindMonthlyPaymentMethodsByCardNumber(c echo.Context) error {
	const method = "FindMonthlyPaymentMethodsByCardNumber"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	cardNumber := c.QueryParam("card_number")
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		status = "error"

		logError("failed to retrieve monthly payment methods by card number", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidYear(c)
	}

	res, err := h.transaction.FindMonthlyPaymentMethodsByCardNumber(ctx, &pb.FindByYearCardNumberTransactionRequest{
		CardNumber: cardNumber,
		Year:       int32(year),
	})
	if err != nil {
		status = "error"

		logError("failed to retrieve monthly payment methods by card number", err, zap.Error(err))

		return transaction_errors.ErrApiFailedMonthlyMethodsByCard(c)
	}

	so := h.mapping.ToApiResponseTransactionMonthMethod(res)

	logSuccess("success retrieve monthly payment methods by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyPaymentMethodsByCardNumber retrieves the yearly payment methods for transactions by card number and year.
// @Summary Get yearly payment methods by card number
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the yearly payment methods for transactions by card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionYearMethod "Yearly payment methods by card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly payment methods by card number"
// @Router /api/transactions/yearly-payment-methods-by-card [get]
func (h *transactionHandler) FindYearlyPaymentMethodsByCardNumber(c echo.Context) error {
	const method = "FindYearlyPaymentMethodsByCardNumber"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	cardNumber := c.QueryParam("card_number")
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		status = "error"

		logError("failed to retrieve yearly payment methods by card number", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidYear(c)
	}

	res, err := h.transaction.FindYearlyPaymentMethodsByCardNumber(ctx, &pb.FindByYearCardNumberTransactionRequest{
		CardNumber: cardNumber,
		Year:       int32(year),
	})
	if err != nil {
		status = "error"

		logError("failed to retrieve yearly payment methods by card number", err, zap.Error(err))

		return transaction_errors.ErrApiFailedYearlyMethodsByCard(c)
	}

	so := h.mapping.ToApiResponseTransactionYearMethod(res)

	logSuccess("success retrieve yearly payment methods by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyAmountsByCardNumber retrieves the monthly transaction amounts for a specific card number and year.
// @Summary Get monthly transaction amounts by card number
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the monthly transaction amounts for a specific card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionMonthAmount "Monthly transaction amounts by card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction amounts by card number"
// @Router /api/transactions/monthly-amounts-by-card [get]
func (h *transactionHandler) FindMonthlyAmountsByCardNumber(c echo.Context) error {
	const method = "FindMonthlyAmountsByCardNumber"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	cardNumber := c.QueryParam("card_number")

	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		status = "error"

		logError("Failed to retrieve monthly amounts by card number", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidYear(c)
	}

	res, err := h.transaction.FindMonthlyAmountsByCardNumber(ctx, &pb.FindByYearCardNumberTransactionRequest{
		CardNumber: cardNumber,
		Year:       int32(year),
	})
	if err != nil {
		status = "error"

		logError("Failed to retrieve monthly amounts by card number", err, zap.Error(err))

		return transaction_errors.ErrApiFailedMonthlyAmountsByCard(c)
	}

	so := h.mapping.ToApiResponseTransactionMonthAmount(res)

	logSuccess("success retrieve monthly amounts by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyAmountsByCardNumber retrieves the yearly transaction amounts for a specific card number and year.
// @Summary Get yearly transaction amounts by card number
// @Tags Transaction
// @Security Bearer
// @Description Retrieve the yearly transaction amounts for a specific card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransactionYearAmount "Yearly transaction amounts by card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction amounts by card number"
// @Router /api/transactions/yearly-amounts-by-card [get]
func (h *transactionHandler) FindYearlyAmountsByCardNumber(c echo.Context) error {
	const method = "FindYearlyAmountsByCardNumber"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	cardNumber := c.QueryParam("card_number")

	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		status = "error"

		logError("Failed to retrieve yearly amounts by card number", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidYear(c)
	}

	res, err := h.transaction.FindYearlyAmountsByCardNumber(ctx, &pb.FindByYearCardNumberTransactionRequest{
		CardNumber: cardNumber,
		Year:       int32(year),
	})
	if err != nil {
		status = "error"

		logError("Failed to retrieve yearly amounts by card number", err, zap.Error(err))

		return transaction_errors.ErrApiFailedYearlyAmountsByCard(c)
	}

	so := h.mapping.ToApiResponseTransactionYearAmount(res)

	logSuccess("success retrieve yearly amounts by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Find transactions by merchant ID
// @Tags Transaction
// @Security Bearer
// @Description Retrieve a list of transactions using the merchant ID
// @Accept json
// @Produce json
// @Param merchant_id query string true "Merchant ID"
// @Success 200 {object} response.ApiResponseTransactions "Transaction data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transactions/merchant/{merchant_id} [get]
func (h *transactionHandler) FindByTransactionMerchantId(c echo.Context) error {
	const method = "FindByTransactionMerchantId"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	merchantId := c.QueryParam("merchant_id")

	merchantIdInt, err := strconv.Atoi(merchantId)

	if err != nil {
		status = "error"

		logError("Failed to retrieve transaction data", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidTransactionMerchantID(c)
	}

	req := &pb.FindTransactionByMerchantIdRequest{
		MerchantId: int32(merchantIdInt),
	}

	res, err := h.transaction.FindTransactionByMerchantId(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve transaction data", err, zap.Error(err))

		return transaction_errors.ErrApiFailedFindByMerchantID(c)
	}

	so := h.mapping.ToApiResponseTransactions(res)

	logSuccess("success retrieve transaction data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Find active transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve a list of active transactions
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Number of items per page (default: 10)"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponseTransactions "List of active transactions"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transactions/active [get]
func (h *transactionHandler) FindByActiveTransaction(c echo.Context) error {
	const method = "FindByActiveTransaction"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")

	req := &pb.FindAllTransactionRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.transaction.FindByActiveTransaction(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve transaction data", err, zap.Error(err))

		return transaction_errors.ErrApiFailedFindActive(c)
	}

	so := h.mapping.ToApiResponsePaginationTransactionDeleteAt(res)

	logSuccess("success retrieve transaction data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Retrieve trashed transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve a list of trashed transactions
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Number of items per page (default: 10)"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponseTransactions "List of trashed transactions"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transactions/trashed [get]
func (h *transactionHandler) FindByTrashedTransaction(c echo.Context) error {
	const method = "FindByTrashedTransaction"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")

	req := &pb.FindAllTransactionRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.transaction.FindByTrashedTransaction(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve transaction data", err, zap.Error(err))

		return transaction_errors.ErrApiFailedFindTrashed(c)
	}

	so := h.mapping.ToApiResponsePaginationTransactionDeleteAt(res)

	logSuccess("success retrieve transaction data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Create a new transaction
// @Tags Transaction
// @Security Bearer
// @Description Create a new transaction record with the provided details.
// @Accept json
// @Produce json
// @Param CreateTransactionRequest body requests.CreateTransactionRequest true "Create Transaction Request"
// @Success 200 {object} response.ApiResponseTransaction "Successfully created transaction record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create transaction"
// @Router /api/transactions/create [post]
func (h *transactionHandler) Create(c echo.Context) error {
	const method = "Create"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	var body requests.CreateTransactionRequest

	apiKeyRaw := c.Get("apiKey")
	if apiKeyRaw == nil {
		status = "error"
		err := errors.New("api key not found")

		logError("Failed to create transaction", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidTransactionApiKey(c)
	}

	apiKey, ok := apiKeyRaw.(string)
	if !ok || apiKey == "" {
		status = "error"
		err := errors.New("invalid api key")

		logError("Failed to create transaction", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidTransactionApiKey(c)
	}

	if err := c.Bind(&body); err != nil {
		status = "error"
		logError("Failed to create transaction", err, zap.Error(err))

		return transaction_errors.ErrApiBindCreateTransaction(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"

		logError("Failed to create transaction", err, zap.Error(err))

		return transaction_errors.ErrApiValidateCreateTransaction(c)
	}

	merchantIDRaw := c.Get("merchant_id")

	if merchantIDRaw == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized: merchant ID not found")
	}

	var merchantID int
	switch v := merchantIDRaw.(type) {
	case float64:
		merchantID = int(v)
	case int:
		merchantID = v
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to convert merchant ID")
	}

	h.logger.Debug("Merchant ID", zap.Int("merchantID", merchantID))

	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to convert merchant ID to string")
	}

	h.logger.Debug("Merchant ID", zap.Int("merchantID", merchantID))

	res, err := h.transaction.CreateTransaction(ctx, &pb.CreateTransactionRequest{
		ApiKey:          apiKey,
		CardNumber:      body.CardNumber,
		Amount:          int32(body.Amount),
		PaymentMethod:   body.PaymentMethod,
		MerchantId:      int32(merchantID),
		TransactionTime: timestamppb.New(body.TransactionTime),
	})

	if err != nil {
		status = "error"
		logError("Failed to create transaction", err, zap.Error(err))

		return transaction_errors.ErrApiFailedCreateTransaction(c)
	}

	so := h.mapping.ToApiResponseTransaction(res)

	logSuccess("success create transaction", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Update a transaction
// @Tags Transaction
// @Security Bearer
// @Description Update an existing transaction record using its ID
// @Accept json
// @Produce json
// @Param transaction body requests.UpdateTransactionRequest true "Transaction data"
// @Success 200 {object} response.ApiResponseTransaction "Updated transaction data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update transaction"
// @Router /api/transactions/update [post]
func (h *transactionHandler) Update(c echo.Context) error {
	const method = "Update"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Failed to update transaction", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidTransactionID(c)
	}

	var body requests.UpdateTransactionRequest

	body.MerchantID = &id

	apiKey, ok := c.Get("apiKey").(string)
	if !ok {
		status = "error"
		err := errors.New("invalid api key")

		logError("Failed to update transaction", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidTransactionApiKey(c)
	}

	if err := c.Bind(&body); err != nil {
		status = "error"
		logError("Failed to update transaction", err, zap.Error(err))

		return transaction_errors.ErrApiBindUpdateTransaction(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"
		logError("Failed to update transaction", err, zap.Error(err))

		return transaction_errors.ErrApiValidateUpdateTransaction(c)
	}

	res, err := h.transaction.UpdateTransaction(ctx, &pb.UpdateTransactionRequest{
		TransactionId:   int32(id),
		CardNumber:      body.CardNumber,
		ApiKey:          apiKey,
		Amount:          int32(body.Amount),
		PaymentMethod:   body.PaymentMethod,
		MerchantId:      int32(*body.MerchantID),
		TransactionTime: timestamppb.New(body.TransactionTime),
	})

	if err != nil {
		status = "error"
		logError("Failed to update transaction", err, zap.Error(err))

		return transaction_errors.ErrApiFailedUpdateTransaction(c)
	}

	so := h.mapping.ToApiResponseTransaction(res)

	logSuccess("success update transaction", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Trash a transaction
// @Tags Transaction
// @Security Bearer
// @Description Trash a transaction record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransaction "Successfully trashed transaction record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to trashed transaction"
// @Router /api/transactions/trashed/{id} [post]
func (h *transactionHandler) TrashedTransaction(c echo.Context) error {
	const method = "TrashedTransaction"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		status = "error"
		logError("Bad Request: Invalid ID", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidTransactionID(c)
	}

	res, err := h.transaction.TrashedTransaction(ctx, &pb.FindByIdTransactionRequest{
		TransactionId: int32(idInt),
	})

	if err != nil {
		status = "error"
		logError("Failed to trashed transaction", err, zap.Error(err))

		return transaction_errors.ErrApiFailedTrashTransaction(c)
	}

	so := h.mapping.ToApiResponseTransaction(res)

	logSuccess("success trashed transaction", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Restore a trashed transaction
// @Tags Transaction
// @Security Bearer
// @Description Restore a trashed transaction record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransaction "Successfully restored transaction record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore transaction:"
// @Router /api/transactions/restore/{id} [post]
func (h *transactionHandler) RestoreTransaction(c echo.Context) error {
	const method = "RestoreTransaction"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		status = "error"
		logError("Bad Request: Invalid ID", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidTransactionID(c)
	}

	res, err := h.transaction.RestoreTransaction(ctx, &pb.FindByIdTransactionRequest{
		TransactionId: int32(idInt),
	})

	if err != nil {
		status = "error"
		logError("Failed to restore transaction", err, zap.Error(err))

		return transaction_errors.ErrApiFailedRestoreTransaction(c)
	}

	so := h.mapping.ToApiResponseTransaction(res)

	logSuccess("success restore transaction", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Permanently delete a transaction
// @Tags Transaction
// @Security Bearer
// @Description Permanently delete a transaction record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransactionDelete "Successfully deleted transaction record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete transaction:"
// @Router /api/transactions/permanent/{id} [delete]
func (h *transactionHandler) DeletePermanent(c echo.Context) error {
	const method = "DeletePermanent"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		status = "error"
		logError("Bad Request: Invalid ID", err, zap.Error(err))

		return transaction_errors.ErrApiInvalidTransactionID(c)
	}

	res, err := h.transaction.DeleteTransactionPermanent(ctx, &pb.FindByIdTransactionRequest{
		TransactionId: int32(idInt),
	})

	if err != nil {
		status = "error"
		logError("Failed to delete transaction", err, zap.Error(err))

		return transaction_errors.ErrApiFailedDeletePermanent(c)
	}

	so := h.mapping.ToApiResponseTransactionDelete(res)

	logSuccess("success delete transaction", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Restore a trashed transaction
// @Tags Transaction
// @Security Bearer
// @Description Restore a trashed transaction all.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseTransactionAll "Successfully restored transaction record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore transaction:"
// @Router /api/transactions/restore/all [post]
func (h *transactionHandler) RestoreAllTransaction(c echo.Context) error {
	const method = "RestoreAllTransaction"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	res, err := h.transaction.RestoreAllTransaction(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"
		logError("Failed to restore all transaction", err, zap.Error(err))

		return transaction_errors.ErrApiFailedRestoreAllTransactions(c)
	}

	so := h.mapping.ToApiResponseTransactionAll(res)

	logSuccess("success restore all transaction", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Permanently delete a transaction
// @Tags Transaction
// @Security Bearer
// @Description Permanently delete a transaction all.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseTransactionAll "Successfully deleted transaction record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete transaction:"
// @Router /api/transactions/delete/all [post]
func (h *transactionHandler) DeleteAllTransactionPermanent(c echo.Context) error {
	const method = "DeleteAllTransactionPermanent"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	res, err := h.transaction.DeleteAllTransactionPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"
		logError("Failed to delete all transaction permanently", err, zap.Error(err))

		return transaction_errors.ErrApiFailedDeleteAllPermanent(c)
	}

	so := h.mapping.ToApiResponseTransactionAll(res)

	logSuccess("success delete all transaction", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *transactionHandler) startTracingAndLogging(
	ctx context.Context,
	method string,
	attrs ...attribute.KeyValue,
) (func(string), func(string, ...zap.Field), func(string, error, ...zap.Field)) {
	start := time.Now()
	_, span := s.trace.Start(ctx, method)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.AddEvent("Start: " + method)
	s.logger.Debug("Start: " + method)

	end := func(status string) {
		s.recordMetrics(method, status, start)
		code := otelcode.Ok
		if status != "success" {
			code = otelcode.Error
		}
		span.SetStatus(code, status)
		span.End()
	}

	logSuccess := func(msg string, fields ...zap.Field) {
		span.AddEvent(msg)
		s.logger.Debug(msg, fields...)
	}

	logError := func(msg string, err error, fields ...zap.Field) {
		span.RecordError(err)
		span.SetStatus(otelcode.Error, msg)
		span.AddEvent(msg)
		allFields := append([]zap.Field{zap.Error(err)}, fields...)
		s.logger.Error(msg, allFields...)
	}

	return end, logSuccess, logError
}

func (s *transactionHandler) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
