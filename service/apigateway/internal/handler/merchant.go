package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors"
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
)

type merchantHandleApi struct {
	merchant        pb.MerchantServiceClient
	logger          logger.LoggerInterface
	mapping         apimapper.MerchantResponseMapper
	trace           trace.Tracer
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewHandlerMerchant(merchant pb.MerchantServiceClient, router *echo.Echo, logger logger.LoggerInterface, mapping apimapper.MerchantResponseMapper) *merchantHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_handler_requests_total",
			Help: "Total number of card requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_handler_request_duration_seconds",
			Help:    "Duration of card requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	merchantHandler := &merchantHandleApi{
		merchant:        merchant,
		mapping:         mapping,
		logger:          logger,
		trace:           otel.Tracer("merchant-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerMerchant := router.Group("/api/merchants")

	routerMerchant.GET("", merchantHandler.FindAll)
	routerMerchant.GET("/:id", merchantHandler.FindById)

	// TODO: ini belum dibuatkan metrics grafananya
	routerMerchant.GET("/transactions", merchantHandler.FindAllTransactions)
	routerMerchant.GET("/transactions/:merchant_id", merchantHandler.FindAllTransactionByMerchant)
	routerMerchant.GET("/transactions/api-key/:api_key", merchantHandler.FindAllTransactionByApikey)

	routerMerchant.GET("/monthly-payment-methods", merchantHandler.FindMonthlyPaymentMethodsMerchant)
	routerMerchant.GET("/yearly-payment-methods", merchantHandler.FindYearlyPaymentMethodMerchant)
	routerMerchant.GET("/monthly-amount", merchantHandler.FindMonthlyAmountMerchant)
	routerMerchant.GET("/yearly-amount", merchantHandler.FindYearlyAmountMerchant)
	routerMerchant.GET("/monthly-total-amount", merchantHandler.FindMonthlyTotalAmountMerchant)
	routerMerchant.GET("/yearly-total-amount", merchantHandler.FindYearlyTotalAmountMerchant)

	routerMerchant.GET("/monthly-payment-methods-by-merchant", merchantHandler.FindMonthlyPaymentMethodByMerchants)
	routerMerchant.GET("/yearly-payment-methods-by-merchant", merchantHandler.FindYearlyPaymentMethodByMerchants)
	routerMerchant.GET("/monthly-amount-by-merchant", merchantHandler.FindMonthlyAmountByMerchants)
	routerMerchant.GET("/yearly-amount-by-merchant", merchantHandler.FindYearlyAmountByMerchants)
	routerMerchant.GET("/monthly-totalamount-by-merchant", merchantHandler.FindMonthlyTotalAmountByMerchants)
	routerMerchant.GET("/yearly-totalamount-by-merchant", merchantHandler.FindYearlyTotalAmountByMerchants)

	routerMerchant.GET("/monthly-payment-methods-by-apikey", merchantHandler.FindMonthlyPaymentMethodByApikeys)
	routerMerchant.GET("/yearly-payment-methods-by-apikey", merchantHandler.FindYearlyPaymentMethodByApikeys)
	routerMerchant.GET("/monthly-amount-by-apikey", merchantHandler.FindMonthlyAmountByApikeys)
	routerMerchant.GET("/yearly-amount-by-apikey", merchantHandler.FindYearlyAmountByApikeys)
	routerMerchant.GET("/monthly-totalamount-by-apikey", merchantHandler.FindMonthlyTotalAmountByApikeys)
	routerMerchant.GET("/yearly-totalamount-by-apikey", merchantHandler.FindYearlyAmountByApikeys)

	// TODO: ini belum juga di buatkan metrics grafananya
	routerMerchant.GET("/api-key", merchantHandler.FindByApiKey)
	routerMerchant.GET("/merchant-user", merchantHandler.FindByMerchantUserId)

	routerMerchant.GET("/active", merchantHandler.FindByActive)
	routerMerchant.GET("/trashed", merchantHandler.FindByTrashed)

	routerMerchant.POST("/create", merchantHandler.Create)
	routerMerchant.POST("/updates/:id", merchantHandler.Update)
	routerMerchant.POST("/update-status/:id", merchantHandler.UpdateStatus)

	routerMerchant.POST("/trashed/:id", merchantHandler.TrashedMerchant)
	routerMerchant.POST("/restore/:id", merchantHandler.RestoreMerchant)
	routerMerchant.DELETE("/permanent/:id", merchantHandler.Delete)

	routerMerchant.POST("/restore/all", merchantHandler.RestoreAllMerchant)
	routerMerchant.POST("/permanent/all", merchantHandler.DeleteAllMerchantPermanent)

	return merchantHandler
}

// FindAll godoc
// @Summary Find all merchants
// @Tags Merchant
// @Security Bearer
// @Description Retrieve a list of all merchants
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchant "List of merchants"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchants [get]
func (h *merchantHandleApi) FindAll(c echo.Context) error {
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

	req := &pb.FindAllMerchantRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.merchant.FindAllMerchant(ctx, req)

	if err != nil {
		status = "error"

		logError("failed to retrieve merchant data", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindAllMerchants(c)
	}

	so := h.mapping.ToApiResponsesMerchant(res)

	logSuccess("merchant data retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)

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
func (h *merchantHandleApi) FindAllTransactions(c echo.Context) error {
	const method = "FindAllTransactions"

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

	req := &pb.FindAllMerchantRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.merchant.FindAllTransactionMerchant(ctx, req)

	if err != nil {
		status = "error"

		logError("failed to retrieve transaction data", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindAllTransactions(c)
	}

	so := h.mapping.ToApiResponseMerchantsTransactionResponse(res)

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
// @Router /api/merchants/transactions/:merchant_id [get]
func (h *merchantHandleApi) FindAllTransactionByMerchant(c echo.Context) error {
	const method = "FindAllTransactionByMerchant"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	merchantID, err := strconv.Atoi(c.Param("merchant_id"))
	if err != nil || merchantID <= 0 {
		status = "error"

		logError("failed to retrieve transaction data", err, zap.Error(err))
		return merchant_errors.ErrApiInvalidMerchantID(c)
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

	req := &pb.FindAllMerchantTransaction{
		MerchantId: int32(merchantID),
		Page:       int32(page),
		PageSize:   int32(pageSize),
		Search:     search,
	}

	res, err := h.merchant.FindAllTransactionByMerchant(ctx, req)

	if err != nil {
		status = "error"

		logError("failed to retrieve transaction data", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindAllTransactionByMerchant(c)
	}

	so := h.mapping.ToApiResponseMerchantsTransactionResponse(res)

	logSuccess("transaction data retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindById godoc
// @Summary Find a merchant by ID
// @Tags Merchant
// @Security Bearer
// @Description Retrieve a merchant by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchant "Merchant data"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchants/{id} [get]
func (h *merchantHandleApi) FindById(c echo.Context) error {
	const method = "FindById"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("failed to retrieve merchant data", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidMerchantID(c)
	}

	req := &pb.FindByIdMerchantRequest{
		MerchantId: int32(id),
	}

	res, err := h.merchant.FindByIdMerchant(ctx, req)

	if err != nil {
		status = "error"

		logError("failed to retrieve merchant data", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindByIdMerchant(c)
	}

	so := h.mapping.ToApiResponseMerchant(res)

	logSuccess("merchant data retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyPaymentMethodsMerchant godoc
// @Summary Find monthly payment methods for a merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve monthly payment methods for a merchant by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantMonthlyPaymentMethod "Monthly payment methods"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly payment methods"
// @Router /api/merchants/monthly-payment-methods [get]
func (h *merchantHandleApi) FindMonthlyPaymentMethodsMerchant(c echo.Context) error {
	const method = "FindMonthlyPaymentMethodsMerchant"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)

	if err != nil || year <= 0 {
		status = "error"

		logError("failed to retrieve monthly payment methods", err, zap.Error(err))
		return merchant_errors.ErrApiInvalidYear(c)
	}

	req := &pb.FindYearMerchant{
		Year: int32(year),
	}

	res, err := h.merchant.FindMonthlyPaymentMethodsMerchant(ctx, req)
	if err != nil {
		status = "error"

		logError("failed to retrieve monthly payment methods", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindMonthlyPaymentMethodsMerchant(c)
	}

	so := h.mapping.ToApiResponseMonthlyPaymentMethods(res)

	logSuccess("monthly payment methods retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyPaymentMethodMerchant godoc.
// @Summary Find yearly payment methods for a merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve yearly payment methods for a merchant by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantYearlyPaymentMethod "Yearly payment methods"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly payment methods"
// @Router /api/merchants/monthly-amount [get]
func (h *merchantHandleApi) FindYearlyPaymentMethodMerchant(c echo.Context) error {
	const method = "FindYearlyPaymentMethodMerchant"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)

	if err != nil || year <= 0 {
		status = "error"

		logError("failed to retrieve yearly payment methods", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidYear(c)
	}

	req := &pb.FindYearMerchant{
		Year: int32(year),
	}

	res, err := h.merchant.FindYearlyPaymentMethodMerchant(ctx, req)

	if err != nil {
		status = "error"

		logError("failed to retrieve yearly payment methods", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindYearlyPaymentMethodMerchant(c)
	}

	so := h.mapping.ToApiResponseYearlyPaymentMethods(res)

	logSuccess("yearly payment methods retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyAmountMerchant godoc
// @Summary Find monthly transaction amounts for a merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve monthly transaction amounts for a merchant by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantMonthlyAmount "Monthly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction amounts"
// @Router /api/merchants/monthly-amount [get]
func (h *merchantHandleApi) FindMonthlyAmountMerchant(c echo.Context) error {
	const method = "FindMonthlyAmountMerchant"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)

	if err != nil || year <= 0 {
		status = "error"

		logError("failed to retrieve monthly transaction amounts", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidYear(c)
	}

	req := &pb.FindYearMerchant{
		Year: int32(year),
	}

	res, err := h.merchant.FindMonthlyAmountMerchant(ctx, req)
	if err != nil {
		status = "error"

		logError("failed to retrieve monthly transaction amounts", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindMonthlyAmountMerchant(c)
	}

	so := h.mapping.ToApiResponseMonthlyAmounts(res)

	logSuccess("monthly transaction amounts retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyAmountMerchant godoc.
// @Summary Find yearly transaction amounts for a merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve yearly transaction amounts for a merchant by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseYearlyAmount "Yearly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction amounts"
// @Router /api/merchants/yearly-amount [get]
func (h *merchantHandleApi) FindYearlyAmountMerchant(c echo.Context) error {
	const method = "FindYearlyAmountMerchant"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)

	if err != nil || year <= 0 {
		status = "error"

		logError("failed to retrieve yearly transaction amounts", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidYear(c)
	}

	req := &pb.FindYearMerchant{
		Year: int32(year),
	}

	res, err := h.merchant.FindYearlyAmountMerchant(ctx, req)
	if err != nil {
		status = "error"

		logError("failed to retrieve yearly transaction amounts", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindYearlyAmountMerchant(c)
	}

	so := h.mapping.ToApiResponseYearlyAmounts(res)

	logSuccess("yearly transaction amounts retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyAmountMerchant godoc
// @Summary Find monthly transaction amounts for a merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve monthly transaction amounts for a merchant by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantMonthlyAmount "Monthly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction amounts"
// @Router /api/merchants/monthly-total-amount [get]
func (h *merchantHandleApi) FindMonthlyTotalAmountMerchant(c echo.Context) error {
	const method = "FindMonthlyTotalAmountMerchant"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil || year <= 0 {
		status = "error"

		logError("failed to retrieve monthly transaction amounts", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidYear(c)
	}

	req := &pb.FindYearMerchant{
		Year: int32(year),
	}

	res, err := h.merchant.FindMonthlyTotalAmountMerchant(ctx, req)
	if err != nil {
		status = "error"

		logError("failed to retrieve monthly transaction amounts", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindMonthlyTotalAmountMerchant(c)
	}

	so := h.mapping.ToApiResponseMonthlyTotalAmounts(res)

	logSuccess("monthly transaction amounts retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyAmountMerchant godoc.
// @Summary Find yearly transaction amounts for a merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve yearly transaction amounts for a merchant by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseYearlyAmount "Yearly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction amounts"
// @Router /api/merchants/yearly-total-amount [get]
func (h *merchantHandleApi) FindYearlyTotalAmountMerchant(c echo.Context) error {
	const method = "FindYearlyTotalAmountMerchant"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		status = "error"

		logError("failed to retrieve yearly transaction amounts", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidYear(c)
	}

	req := &pb.FindYearMerchant{
		Year: int32(year),
	}

	res, err := h.merchant.FindYearlyTotalAmountMerchant(ctx, req)
	if err != nil {
		status = "error"

		logError("failed to retrieve yearly transaction amounts", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindYearlyTotalAmountMerchant(c)
	}

	so := h.mapping.ToApiResponseYearlyTotalAmounts(res)

	logSuccess("yearly transaction amounts retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyPaymentMethodByMerchants godoc.
// @Summary Find monthly payment methods for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve monthly payment methods for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantMonthlyPaymentMethod "Monthly payment methods"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly payment methods"
// @Router /api/merchants/monthly-payment-methods-by-merchant [get]
func (h *merchantHandleApi) FindMonthlyPaymentMethodByMerchants(c echo.Context) error {
	const method = "FindMonthlyPaymentMethodByMerchants"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	merchantIDStr := c.QueryParam("merchant_id")
	yearStr := c.QueryParam("year")

	merchantID, err := strconv.Atoi(merchantIDStr)
	if err != nil || merchantID <= 0 {
		status = "error"

		logError("failed to find monthly payment methods by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidMerchantID(c)
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		status = "error"

		logError("failed to find monthly payment methods by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidYear(c)
	}

	req := &pb.FindYearMerchantById{
		MerchantId: int32(merchantID),
		Year:       int32(year),
	}

	res, err := h.merchant.FindMonthlyPaymentMethodByMerchants(ctx, req)
	if err != nil {
		status = "error"

		logError("failed to find monthly payment methods by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindMonthlyPaymentMethodByMerchants(c)
	}

	so := h.mapping.ToApiResponseMonthlyPaymentMethods(res)

	logSuccess("monthly payment methods retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyPaymentMethodByMerchants godoc.
// @Summary Find yearly payment methods for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve yearly payment methods for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantYearlyPaymentMethod "Yearly payment methods"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly payment methods"
// @Router /api/merchants/yearly-payment-methods-by-merchant [get]
func (h *merchantHandleApi) FindYearlyPaymentMethodByMerchants(c echo.Context) error {
	const method = "FindYearlyPaymentMethodByMerchants"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	merchantIDStr := c.QueryParam("merchant_id")
	yearStr := c.QueryParam("year")

	merchantID, err := strconv.Atoi(merchantIDStr)
	if err != nil || merchantID <= 0 {
		status = "error"

		logError("failed to find yearly payment methods by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidMerchantID(c)
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		status = "error"

		logError("failed to find yearly payment methods by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidYear(c)
	}

	req := &pb.FindYearMerchantById{
		MerchantId: int32(merchantID),
		Year:       int32(year),
	}

	res, err := h.merchant.FindYearlyPaymentMethodByMerchants(ctx, req)
	if err != nil {
		status = "error"

		logError("failed to find yearly payment methods by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindYearlyPaymentMethodByMerchants(c)
	}

	so := h.mapping.ToApiResponseYearlyPaymentMethods(res)

	logSuccess("yearly payment methods retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyAmountByMerchants godoc.
// @Summary Find monthly transaction amounts for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve monthly transaction amounts for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantMonthlyAmount "Monthly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction amounts"
// @Router /api/merchants/monthly-amount-by-merchant [get]
func (h *merchantHandleApi) FindMonthlyAmountByMerchants(c echo.Context) error {
	const method = "FindMonthlyAmountByMerchants"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	merchantIDStr := c.QueryParam("merchant_id")
	yearStr := c.QueryParam("year")

	merchantID, err := strconv.Atoi(merchantIDStr)
	if err != nil || merchantID <= 0 {
		status = "error"

		logError("failed to find monthly amount by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidMerchantID(c)
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		status = "error"

		logError("failed to find monthly amount by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidYear(c)
	}

	req := &pb.FindYearMerchantById{
		MerchantId: int32(merchantID),
		Year:       int32(year),
	}

	res, err := h.merchant.FindMonthlyAmountByMerchants(ctx, req)
	if err != nil {
		status = "error"

		logError("failed to find monthly amount by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindMonthlyAmountByMerchants(c)
	}

	so := h.mapping.ToApiResponseMonthlyAmounts(res)

	logSuccess("monthly amount retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyAmountByMerchants godoc.
// @Summary Find yearly transaction amounts for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve yearly transaction amounts for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantYearlyAmount "Yearly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction amounts"
// @Router /api/merchants/yearly-amount-by-merchant [get]
func (h *merchantHandleApi) FindYearlyAmountByMerchants(c echo.Context) error {
	const method = "FindYearlyAmountByMerchants"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	merchantIDStr := c.QueryParam("merchant_id")
	yearStr := c.QueryParam("year")

	merchantID, err := strconv.Atoi(merchantIDStr)

	if err != nil || merchantID <= 0 {
		status = "error"

		logError("failed to find yearly amount by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidMerchantID(c)
	}

	year, err := strconv.Atoi(yearStr)

	if err != nil || year <= 0 {
		status = "error"

		logError("failed to find yearly amount by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidYear(c)
	}

	req := &pb.FindYearMerchantById{
		MerchantId: int32(merchantID),
		Year:       int32(year),
	}

	res, err := h.merchant.FindYearlyAmountByMerchants(ctx, req)

	if err != nil {
		status = "error"

		logError("failed to find yearly amount by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindYearlyAmountByMerchants(c)
	}

	so := h.mapping.ToApiResponseYearlyAmounts(res)

	logSuccess("yearly amount retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyAmountByMerchants godoc.
// @Summary Find monthly transaction amounts for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve monthly transaction amounts for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantMonthlyAmount "Monthly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction amounts"
// @Router /api/merchants/monthly-totalamount-by-merchant [get]
func (h *merchantHandleApi) FindMonthlyTotalAmountByMerchants(c echo.Context) error {
	const method = "FindMonthlyTotalAmountByMerchants"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	merchantIDStr := c.QueryParam("merchant_id")
	yearStr := c.QueryParam("year")

	merchantID, err := strconv.Atoi(merchantIDStr)
	if err != nil || merchantID <= 0 {
		status = "error"

		logError("failed to find monthly total amount by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidMerchantID(c)
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		status = "error"

		logError("failed to find monthly total amount by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidYear(c)
	}

	req := &pb.FindYearMerchantById{
		MerchantId: int32(merchantID),
		Year:       int32(year),
	}

	res, err := h.merchant.FindMonthlyTotalAmountByMerchants(ctx, req)
	if err != nil {
		status = "error"

		logError("failed to find monthly total amount by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindMonthlyTotalAmountMerchant(c)
	}

	so := h.mapping.ToApiResponseMonthlyTotalAmounts(res)

	logSuccess("monthly total amount retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyAmountByMerchants godoc.
// @Summary Find yearly transaction amounts for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve yearly transaction amounts for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantYearlyAmount "Yearly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction amounts"
// @Router /api/merchants/yearly-totalamount-by-merchant [get]
func (h *merchantHandleApi) FindYearlyTotalAmountByMerchants(c echo.Context) error {
	const method = "FindYearlyTotalAmountByMerchants"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	merchantIDStr := c.QueryParam("merchant_id")
	yearStr := c.QueryParam("year")

	merchantID, err := strconv.Atoi(merchantIDStr)
	if err != nil || merchantID <= 0 {
		status = "error"

		logError("failed to find yearly total amount by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidMerchantID(c)
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		status = "error"

		logError("failed to find yearly total amount by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidYear(c)
	}

	req := &pb.FindYearMerchantById{
		MerchantId: int32(merchantID),
		Year:       int32(year),
	}

	res, err := h.merchant.FindYearlyTotalAmountByMerchants(ctx, req)
	if err != nil {
		status = "error"

		logError("failed to find yearly total amount by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindYearlyTotalAmountMerchant(c)
	}

	so := h.mapping.ToApiResponseYearlyTotalAmounts(res)

	logSuccess("yearly total amount retrieved successfully", zap.Bool("success", true))

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
// @Router /api/merchants/transactions/api-key/:api_key [get]
func (h *merchantHandleApi) FindAllTransactionByApikey(c echo.Context) error {
	const method = "FindAllTransactionByApikey"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	api_key := c.Param("api_key")

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")

	req := &pb.FindAllMerchantApikey{
		ApiKey:   api_key,
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.merchant.FindAllTransactionByApikey(ctx, req)

	if err != nil {
		status = "error"

		logError("failed to find all transaction by api key", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindAllTransactionByApikey(c)
	}

	so := h.mapping.ToApiResponseMerchantsTransactionResponse(res)

	logSuccess("transaction retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyPaymentMethodByApikeys godoc.
// @Summary Find monthly payment methods for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve monthly payment methods for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantMonthlyPaymentMethod "Monthly payment methods"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly payment methods"
// @Router /api/merchants/monthly-payment-methods-by-apikey [get]
func (h *merchantHandleApi) FindMonthlyPaymentMethodByApikeys(c echo.Context) error {
	const method = "FindMonthlyPaymentMethodByApikeys"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	api_key := c.QueryParam("api_key")
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		status = "error"

		logError("failed to find monthly payment methods by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidYear(c)
	}

	req := &pb.FindYearMerchantByApikey{
		ApiKey: api_key,
		Year:   int32(year),
	}

	res, err := h.merchant.FindMonthlyPaymentMethodByApikey(ctx, req)
	if err != nil {
		status = "error"

		logError("failed to find monthly payment methods by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindMonthlyPaymentMethodsMerchant(c)
	}

	so := h.mapping.ToApiResponseMonthlyPaymentMethods(res)

	logSuccess("monthly payment methods retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyPaymentMethodByApikeys godoc.
// @Summary Find yearly payment methods for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve yearly payment methods for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantYearlyPaymentMethod "Yearly payment methods"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly payment methods"
// @Router /api/merchants/yearly-payment-methods-by-apikey [get]
func (h *merchantHandleApi) FindYearlyPaymentMethodByApikeys(c echo.Context) error {
	const method = "FindYearlyPaymentMethodByApikeys"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	api_key := c.QueryParam("api_key")
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		status = "error"

		logError("failed to find yearly payment methods by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidYear(c)
	}

	req := &pb.FindYearMerchantByApikey{
		ApiKey: api_key,
		Year:   int32(year),
	}

	res, err := h.merchant.FindYearlyPaymentMethodByApikey(ctx, req)
	if err != nil {
		status = "error"

		logError("failed to find yearly payment methods by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindYearlyPaymentMethodMerchant(c)
	}

	so := h.mapping.ToApiResponseYearlyPaymentMethods(res)

	logSuccess("yearly payment methods retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyAmountByApikeys godoc.
// @Summary Find monthly transaction amounts for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve monthly transaction amounts for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantMonthlyAmount "Monthly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction amounts"
// @Router /api/merchants/monthly-amount-by-apikey [get]
func (h *merchantHandleApi) FindMonthlyAmountByApikeys(c echo.Context) error {
	const method = "FindMonthlyAmountByApikeys"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	api_key := c.QueryParam("api_key")
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		status = "error"

		logError("failed to find monthly amount by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidYear(c)
	}

	req := &pb.FindYearMerchantByApikey{
		ApiKey: api_key,
		Year:   int32(year),
	}

	res, err := h.merchant.FindMonthlyAmountByApikey(ctx, req)
	if err != nil {
		status = "error"

		logError("failed to find monthly amount by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindMonthlyAmountMerchant(c)
	}

	so := h.mapping.ToApiResponseMonthlyAmounts(res)

	logSuccess("monthly amounts retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyAmountByApikeys godoc.
// @Summary Find yearly transaction amounts for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve yearly transaction amounts for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantYearlyAmount "Yearly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction amounts"
// @Router /api/merchants/yearly-amount-by-apikey [get]
func (h *merchantHandleApi) FindYearlyAmountByApikeys(c echo.Context) error {
	const method = "FindYearlyAmountByApikeys"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()
	api_key := c.QueryParam("api_key")
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		status = "error"

		logError("failed to find yearly amount by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidYear(c)
	}

	req := &pb.FindYearMerchantByApikey{
		ApiKey: api_key,
		Year:   int32(year),
	}

	res, err := h.merchant.FindYearlyAmountByApikey(ctx, req)
	if err != nil {
		status = "error"

		logError("failed to find yearly amount by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindYearlyAmountMerchant(c)
	}

	so := h.mapping.ToApiResponseYearlyAmounts(res)

	logSuccess("yearly amounts retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyAmountByApikeys godoc.
// @Summary Find monthly transaction amounts for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve monthly transaction amounts for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantMonthlyAmount "Monthly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transaction amounts"
// @Router /api/merchants/monthly-totalamount-by-apikey [get]
func (h *merchantHandleApi) FindMonthlyTotalAmountByApikeys(c echo.Context) error {
	const method = "FindMonthlyTotalAmountByApikeys"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	api_key := c.QueryParam("api_key")
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		status = "error"

		logError("failed to find monthly amount by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidYear(c)
	}

	req := &pb.FindYearMerchantByApikey{
		ApiKey: api_key,
		Year:   int32(year),
	}

	res, err := h.merchant.FindMonthlyTotalAmountByApikey(ctx, req)
	if err != nil {
		status = "error"

		logError("failed to find monthly amount by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindMonthlyTotalAmountMerchant(c)
	}

	so := h.mapping.ToApiResponseMonthlyTotalAmounts(res)

	logSuccess("monthly amounts retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyAmountByApikeys godoc.
// @Summary Find yearly transaction amounts for a specific merchant
// @Tags Merchant
// @Security Bearer
// @Description Retrieve yearly transaction amounts for a specific merchant by year.
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMerchantYearlyAmount "Yearly transaction amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transaction amounts"
// @Router /api/merchants/yearly-totalamount-by-apikey [get]
func (h *merchantHandleApi) FindYearlyTotalAmountByApikeys(c echo.Context) error {
	const method = "FindYearlyAmountByApikeys"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	api_key := c.QueryParam("api_key")
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		status = "error"

		logError("failed to find yearly amount by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidYear(c)
	}

	req := &pb.FindYearMerchantByApikey{
		ApiKey: api_key,
		Year:   int32(year),
	}

	res, err := h.merchant.FindYearlyTotalAmountByApikey(ctx, req)
	if err != nil {
		status = "error"

		logError("failed to find yearly amount by merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindYearlyTotalAmountMerchant(c)
	}

	so := h.mapping.ToApiResponseYearlyTotalAmounts(res)

	logSuccess("yearly amounts retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

//

// FindByApiKey godoc
// @Summary Find a merchant by API key
// @Tags Merchant
// @Security Bearer
// @Description Retrieve a merchant by its API key
// @Accept json
// @Produce json
// @Param api_key query string true "API key"
// @Success 200 {object} response.ApiResponseMerchant "Merchant data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchants/api-key [get]
func (h *merchantHandleApi) FindByApiKey(c echo.Context) error {
	const method = "FindByApiKey"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	apiKey := c.QueryParam("api_key")

	if apiKey == "" {
		status = "error"

		err := errors.New("api key is empty")
		logError("failed to find merchant by api key", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidApiKey(c)
	}

	req := &pb.FindByApiKeyRequest{
		ApiKey: apiKey,
	}

	res, err := h.merchant.FindByApiKey(ctx, req)

	if err != nil {
		status = "error"

		logError("failed to find merchant by api key", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindByApiKeyMerchant(c)
	}

	so := h.mapping.ToApiResponseMerchant(res)

	logSuccess("merchant retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindByMerchantUserId godoc.
// @Summary Find a merchant by user ID
// @Tags Merchant
// @Security Bearer
// @Description Retrieve a merchant by its user ID
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.ApiResponsesMerchant "Merchant data"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchants/merchant-user [get]
func (h *merchantHandleApi) FindByMerchantUserId(c echo.Context) error {
	const method = "FindByMerchantUserId"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	id, ok := c.Get("user_id").(int32)

	if !ok {
		status = "error"

		err := errors.New("user id not found in context")
		logError("failed to find merchant by user id", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidUserID(c)
	}

	req := &pb.FindByMerchantUserIdRequest{
		UserId: id,
	}

	res, err := h.merchant.FindByMerchantUserId(ctx, req)

	if err != nil {
		status = "error"

		logError("failed to find merchant by user id", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindByUserId(c)
	}

	so := h.mapping.ToApiResponseMerchants(res)

	logSuccess("merchant retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindByActive godoc
// @Summary Find active merchants
// @Tags Merchant
// @Security Bearer
// @Description Retrieve a list of active merchants
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsesMerchant "List of active merchants"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchants/active [get]
func (h *merchantHandleApi) FindByActive(c echo.Context) error {
	const method = "FindByActive"
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

	req := &pb.FindAllMerchantRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.merchant.FindByActive(ctx, req)

	if err != nil {
		status = "error"

		logError("failed to find merchant by active", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindAllMerchantsActive(c)
	}

	so := h.mapping.ToApiResponsesMerchantDeleteAt(res)

	logSuccess("merchant retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindByTrashed godoc
// @Summary Find trashed merchants
// @Tags Merchant
// @Security Bearer
// @Description Retrieve a list of trashed merchants
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsesMerchant "List of trashed merchants"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchants/trashed [get]
func (h *merchantHandleApi) FindByTrashed(c echo.Context) error {
	const method = "FindByTrashed"
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

	req := &pb.FindAllMerchantRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.merchant.FindByTrashed(ctx, req)

	if err != nil {
		status = "error"

		logError("failed to find merchant by trashed", err, zap.Error(err))

		return merchant_errors.ErrApiFailedFindAllMerchantsTrashed(c)
	}

	so := h.mapping.ToApiResponsesMerchantDeleteAt(res)

	logSuccess("merchant retrieved successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// Create godoc
// @Summary Create a new merchant
// @Tags Merchant
// @Security Bearer
// @Description Create a new merchant with the given name and user ID
// @Accept json
// @Produce json
// @Param body body requests.CreateMerchantRequest true "Create merchant request"
// @Success 200 {object} response.ApiResponseMerchant "Created merchant"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create merchant"
// @Router /api/merchants/create [post]
func (h *merchantHandleApi) Create(c echo.Context) error {
	const method = "Create"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	var body requests.CreateMerchantRequest

	if err := c.Bind(&body); err != nil {
		status = "error"

		logError("Failed to bind CreateMerchant request", err, zap.Error(err))

		return merchant_errors.ErrApiBindCreateMerchant(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"

		logError("Validation Error", err, zap.Error(err))

		return merchant_errors.ErrApiValidateCreateMerchant(c)
	}

	req := &pb.CreateMerchantRequest{
		Name:   body.Name,
		UserId: int32(body.UserID),
	}

	res, err := h.merchant.CreateMerchant(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to create merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedCreateMerchant(c)
	}

	so := h.mapping.ToApiResponseMerchant(res)

	logSuccess("Merchant created successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// Update godoc
// @Summary Update a merchant
// @Tags Merchant
// @Security Bearer
// @Description Update a merchant with the given ID
// @Accept json
// @Produce json
// @Param body body requests.UpdateMerchantRequest true "Update merchant request"
// @Success 200 {object} response.ApiResponseMerchant "Updated merchant"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update merchant"
// @Router /api/merchants/update/{id} [post]
func (h *merchantHandleApi) Update(c echo.Context) error {
	const method = "FindMonthlyAmountByApikeys"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Invalid merchant ID", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidMerchantID(c)
	}

	var body requests.UpdateMerchantRequest

	if err := c.Bind(&body); err != nil {
		status = "error"

		logError("Failed to bind UpdateMerchant request", err, zap.Error(err))

		return merchant_errors.ErrApiBindUpdateMerchant(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"

		logError("Validation Error", err, zap.Error(err))

		return merchant_errors.ErrApiValidateUpdateMerchant(c)
	}

	req := &pb.UpdateMerchantRequest{
		MerchantId: int32(id),
		Name:       body.Name,
		UserId:     int32(body.UserID),
		Status:     body.Status,
	}

	res, err := h.merchant.UpdateMerchant(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to update merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedUpdateMerchant(c)
	}

	so := h.mapping.ToApiResponseMerchant(res)

	logSuccess("Merchant updated successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// UpdateStatus godoc
// @Summary Update merchant status
// @Tags Merchant
// @Security Bearer
// @Description Update the status of a merchant with the given ID
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Param body body requests.UpdateMerchantStatusRequest true "Update merchant status request"
// @Success 200 {object} response.ApiResponseMerchant "Updated merchant status"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update merchant status"
// @Router /api/merchants/update-status/{id} [post]
func (h *merchantHandleApi) UpdateStatus(c echo.Context) error {
	const method = "UpdateStatus"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		status = "error"

		logError("Invalid merchant ID", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidMerchantID(c)
	}

	var body requests.UpdateMerchantStatusRequest
	if err := c.Bind(&body); err != nil {
		status = "error"

		logError("Failed to bind UpdateMerchantStatus request", err, zap.Error(err))

		return merchant_errors.ErrApiBindUpdateMerchantStatus(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"

		logError("Validation Error", err, zap.Error(err))

		return merchant_errors.ErrApiValidateUpdateMerchantStatus(c)
	}

	req := &pb.UpdateMerchantStatusRequest{
		MerchantId: int32(id),
		Status:     body.Status,
	}

	res, err := h.merchant.UpdateMerchantStatus(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to update merchant status", err, zap.Error(err))

		return merchant_errors.ErrApiFailedUpdateMerchantStatus(c)
	}

	so := h.mapping.ToApiResponseMerchant(res)

	logSuccess("Merchant status updated successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// TrashedMerchant godoc
// @Summary Trashed a merchant
// @Tags Merchant
// @Security Bearer
// @Description Trashed a merchant by its ID
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchant "Trashed merchant"
// @Failure 400 {object} response.ErrorResponse "Bad request or invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to trashed merchant"
// @Router /api/merchants/trashed/{id} [post]
func (h *merchantHandleApi) TrashedMerchant(c echo.Context) error {
	const method = "TrashedMerchant"
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

		logError("Invalid merchant ID", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidMerchantID(c)
	}

	res, err := h.merchant.TrashedMerchant(ctx, &pb.FindByIdMerchantRequest{
		MerchantId: int32(idInt),
	})

	if err != nil {
		status = "error"

		logError("Failed to trashed merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedTrashMerchant(c)
	}

	so := h.mapping.ToApiResponseMerchant(res)

	logSuccess("Merchant trashed successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// RestoreMerchant godoc
// @Summary Restore a merchant
// @Tags Merchant
// @Security Bearer
// @Description Restore a merchant by its ID
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchant "Restored merchant"
// @Failure 400 {object} response.ErrorResponse "Bad request or invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore merchant"
// @Router /api/merchants/restore/{id} [post]
func (h *merchantHandleApi) RestoreMerchant(c echo.Context) error {
	const method = "FindMonthlyAmountByApikeys"
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

		logError("Invalid merchant ID", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidMerchantID(c)
	}

	res, err := h.merchant.RestoreMerchant(ctx, &pb.FindByIdMerchantRequest{
		MerchantId: int32(idInt),
	})

	if err != nil {
		status = "error"

		logError("Failed to restore merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedRestoreMerchant(c)
	}

	so := h.mapping.ToApiResponseMerchant(res)

	logSuccess("Merchant restored successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// Delete godoc
// @Summary Delete a merchant permanently
// @Tags Merchant
// @Security Bearer
// @Description Delete a merchant by its ID permanently
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchantDelete "Deleted merchant"
// @Failure 400 {object} response.ErrorResponse "Bad request or invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete merchant"
// @Router /api/merchants/{id} [delete]
func (h *merchantHandleApi) Delete(c echo.Context) error {
	const method = "Delete"
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

		logError("Invalid merchant ID", err, zap.Error(err))

		return merchant_errors.ErrApiInvalidMerchantID(c)
	}

	res, err := h.merchant.DeleteMerchantPermanent(ctx, &pb.FindByIdMerchantRequest{
		MerchantId: int32(idInt),
	})

	if err != nil {
		status = "error"

		logError("Failed to delete merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedDeleteMerchantPermanent(c)
	}

	so := h.mapping.ToApiResponseMerchantDeleteAt(res)

	logSuccess("Merchant deleted successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// RestoreAllMerchant godoc.
// @Summary Restore all merchant records
// @Tags Merchant
// @Security Bearer
// @Description Restore all merchant records that were previously deleted.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully restored all merchant records"
// @Failure 500 {object} response.ErrorResponse "Failed to restore all merchant records"
// @Router /api/merchants/restore/all [post]
func (h *merchantHandleApi) RestoreAllMerchant(c echo.Context) error {
	const method = "FindMonthlyAmountByApikeys"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	res, err := h.merchant.RestoreAllMerchant(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"

		logError("Failed to restore all merchant", err, zap.Error(err))

		return merchant_errors.ErrApiFailedRestoreAllMerchant(c)
	}

	so := h.mapping.ToApiResponseMerchantAll(res)

	logSuccess("Merchant restored successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// DeleteAllMerchantPermanent godoc.
// @Summary Permanently delete all merchant records
// @Tags Merchant
// @Security Bearer
// @Description Permanently delete all merchant records from the database.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully deleted all merchant records permanently"
// @Failure 500 {object} response.ErrorResponse "Failed to permanently delete all merchant records"
// @Router /api/merchants/permanent/all [post]
func (h *merchantHandleApi) DeleteAllMerchantPermanent(c echo.Context) error {
	const method = "DeleteAllMerchantPermanent"

	end, logSuccess, logError := h.startTracingAndLogging(c.Request().Context(), method)
	status := "success"

	defer func() {

		end(status)
	}()

	ctx := c.Request().Context()

	res, err := h.merchant.DeleteAllMerchantPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"

		logError("Failed to delete all merchant permanently", err, zap.Error(err))

		return merchant_errors.ErrApiFailedDeleteAllMerchantPermanent(c)
	}

	so := h.mapping.ToApiResponseMerchantAll(res)

	logSuccess("Merchant deleted successfully", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *merchantHandleApi) startTracingAndLogging(
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

func (s *merchantHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
