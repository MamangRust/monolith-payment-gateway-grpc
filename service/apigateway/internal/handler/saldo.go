package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors"
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

type saldoHandleApi struct {
	saldo           pb.SaldoServiceClient
	logger          logger.LoggerInterface
	mapping         apimapper.SaldoResponseMapper
	trace           trace.Tracer
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewHandlerSaldo(client pb.SaldoServiceClient, router *echo.Echo, logger logger.LoggerInterface, mapping apimapper.SaldoResponseMapper) *saldoHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "saldo_handler_requests_total",
			Help: "Total number of card requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "saldo_handler_request_duration_seconds",
			Help:    "Duration of card requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	saldoHandler := &saldoHandleApi{
		saldo:           client,
		logger:          logger,
		mapping:         mapping,
		trace:           otel.Tracer("saldo-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerSaldo := router.Group("/api/saldos")

	routerSaldo.GET("", saldoHandler.FindAll)
	routerSaldo.GET("/:id", saldoHandler.FindById)

	routerSaldo.GET("/monthly-total-balance", saldoHandler.FindMonthlyTotalSaldoBalance)
	routerSaldo.GET("/yearly-total-balance", saldoHandler.FindYearTotalSaldoBalance)
	routerSaldo.GET("/monthly-balances", saldoHandler.FindMonthlySaldoBalances)
	routerSaldo.GET("/yearly-balances", saldoHandler.FindYearlySaldoBalances)

	routerSaldo.GET("/active", saldoHandler.FindByActive)
	routerSaldo.GET("/trashed", saldoHandler.FindByTrashed)
	routerSaldo.GET("/card_number/:card_number", saldoHandler.FindByCardNumber)

	routerSaldo.POST("/create", saldoHandler.Create)
	routerSaldo.POST("/update/:id", saldoHandler.Update)
	routerSaldo.POST("/trashed/:id", saldoHandler.TrashSaldo)
	routerSaldo.POST("/restore/:id", saldoHandler.RestoreSaldo)
	routerSaldo.DELETE("/permanent/:id", saldoHandler.Delete)

	routerSaldo.POST("/restore/all", saldoHandler.RestoreAllSaldo)
	routerSaldo.POST("/permanent/all", saldoHandler.DeleteAllSaldoPermanent)

	return saldoHandler

}

// @Summary Find all saldo data
// @Tags Saldo
// @Security Bearer
// @Description Retrieve a list of all saldo data with pagination and search
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationSaldo "List of saldo data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve saldo data"
// @Router /api/saldos [get]
func (h *saldoHandleApi) FindAll(c echo.Context) error {
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

	req := &pb.FindAllSaldoRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.saldo.FindAllSaldo(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve saldo data", err, zap.Error(err))

		return saldo_errors.ErrApiFailedFindAllSaldo(c)
	}

	so := h.mapping.ToApiResponsePaginationSaldo(res)

	logSuccess("Successfully retrieve saldo data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Find a saldo by ID
// @Tags Saldo
// @Security Bearer
// @Description Retrieve a saldo by its ID
// @Accept json
// @Produce json
// @Param id path int true "Saldo ID"
// @Success 200 {object} response.ApiResponseSaldo "Saldo data"
// @Failure 400 {object} response.ErrorResponse "Invalid saldo ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve saldo data"
// @Router /api/saldos/{id} [get]
func (h *saldoHandleApi) FindById(c echo.Context) error {
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

		logError("Invalid saldo ID", err, zap.Error(err))

		return saldo_errors.ErrApiInvalidSaldoID(c)
	}

	req := &pb.FindByIdSaldoRequest{
		SaldoId: int32(id),
	}

	res, err := h.saldo.FindByIdSaldo(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve saldo data", err, zap.Error(err))

		return saldo_errors.ErrApiFailedFindByIdSaldo(c)
	}

	so := h.mapping.ToApiResponseSaldo(res)

	logSuccess("Successfully retrieve saldo data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTotalSaldoBalance retrieves the total saldo balance for a specific month and year.
// @Summary Get monthly total saldo balance
// @Tags Saldo
// @Security Bearer
// @Description Retrieve the total saldo balance for a specific month and year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseMonthTotalSaldo "Monthly total saldo balance"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly total saldo balance"
// @Router /api/saldos/monthly-total-balance [get]
func (h *saldoHandleApi) FindMonthlyTotalSaldoBalance(c echo.Context) error {
	const method = "FindMonthlyTotalSaldoBalance"
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

		logError("Invalid year parameter", err, zap.Error(err))

		return saldo_errors.ErrApiInvalidYear(c)
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		status = "error"

		logError("Invalid month parameter", err, zap.Error(err))

		return saldo_errors.ErrApiInvalidMonth(c)
	}

	res, err := h.saldo.FindMonthlyTotalSaldoBalance(ctx, &pb.FindMonthlySaldoTotalBalance{
		Year:  int32(year),
		Month: int32(month),
	})
	if err != nil {
		status = "error"

		logError("Failed to retrieve monthly total saldo balance", err, zap.Error(err))

		return saldo_errors.ErrApiFailedFindMonthlyTotalSaldoBalance(c)
	}

	so := h.mapping.ToApiResponseMonthTotalSaldo(res)

	logSuccess("Successfully retrieve monthly total saldo balance", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearTotalSaldoBalance retrieves the total saldo balance for a specific year.
// @Summary Get yearly total saldo balance
// @Tags Saldo
// @Security Bearer
// @Description Retrieve the total saldo balance for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseYearTotalSaldo "Yearly total saldo balance"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly total saldo balance"
// @Router /api/saldos/yearly-total-balance [get]
func (h *saldoHandleApi) FindYearTotalSaldoBalance(c echo.Context) error {
	const method = "FindYearTotalSaldoBalance"
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

		logError("Invalid year parameter", err, zap.Error(err))

		return saldo_errors.ErrApiInvalidYear(c)
	}

	res, err := h.saldo.FindYearTotalSaldoBalance(ctx, &pb.FindYearlySaldo{
		Year: int32(year),
	})
	if err != nil {
		status = "error"

		logError("Failed to retrieve yearly total saldo balance", err, zap.Error(err))

		return saldo_errors.ErrApiFailedFindYearTotalSaldoBalance(c)
	}

	so := h.mapping.ToApiResponseYearTotalSaldo(res)

	logSuccess("Successfully retrieve yearly total saldo balance", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlySaldoBalances retrieves monthly saldo balances for a specific year.
// @Summary Get monthly saldo balances
// @Tags Saldo
// @Security Bearer
// @Description Retrieve monthly saldo balances for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMonthSaldoBalances "Monthly saldo balances"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly saldo balances"
// @Router /api/saldos/monthly-balances [get]
func (h *saldoHandleApi) FindMonthlySaldoBalances(c echo.Context) error {
	const method = "FindMonthlySaldoBalances"
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

		logError("Invalid year parameter", err, zap.Error(err))

		return saldo_errors.ErrApiInvalidYear(c)
	}

	res, err := h.saldo.FindMonthlySaldoBalances(ctx, &pb.FindYearlySaldo{
		Year: int32(year),
	})
	if err != nil {
		status = "error"

		logError("Failed to retrieve monthly saldo balances", err, zap.Error(err))

		return saldo_errors.ErrApiFailedFindMonthlySaldoBalances(c)
	}

	so := h.mapping.ToApiResponseMonthSaldoBalances(res)

	logSuccess("Successfully retrieve monthly saldo balances", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlySaldoBalances retrieves yearly saldo balances for a specific year.
// @Summary Get yearly saldo balances
// @Tags Saldo
// @Security Bearer
// @Description Retrieve yearly saldo balances for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseYearSaldoBalances "Yearly saldo balances"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly saldo balances"
// @Router /api/saldo/yearly-balances [get]
func (h *saldoHandleApi) FindYearlySaldoBalances(c echo.Context) error {
	const method = "FindYearlySaldoBalances"
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

		logError("Invalid year parameter", err, zap.Error(err))

		return saldo_errors.ErrApiInvalidYear(c)
	}

	res, err := h.saldo.FindYearlySaldoBalances(ctx, &pb.FindYearlySaldo{
		Year: int32(year),
	})
	if err != nil {
		status = "error"

		logError("Failed to retrieve yearly saldo balances", err, zap.Error(err))

		return saldo_errors.ErrApiFailedFindYearlySaldoBalances(c)
	}

	so := h.mapping.ToApiResponseYearSaldoBalances(res)

	logSuccess("Successfully retrieve yearly saldo balances", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Find a saldo by card number
// @Tags Saldo
// @Security Bearer
// @Description Retrieve a saldo by its card number
// @Accept json
// @Produce json
// @Param card_number path string true "Card number"
// @Success 200 {object} response.ApiResponseSaldo "Saldo data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve saldo data"
// @Router /api/saldos/card_number/{card_number} [get]
func (h *saldoHandleApi) FindByCardNumber(c echo.Context) error {
	const method = "FindByCardNumber"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	cardNumber := c.Param("card_number")

	if cardNumber == "" {
		status = "error"
		err := errors.New("invalid card number")

		logError("Invalid card number", err, zap.Error(err))

		return saldo_errors.ErrApiInvalidCardNumber(c)
	}

	req := &pb.FindByCardNumberRequest{
		CardNumber: cardNumber,
	}

	res, err := h.saldo.FindByCardNumber(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve saldo data", err, zap.Error(err))

		return saldo_errors.ErrApiFailedFindByCardNumberSaldo(c)
	}

	so := h.mapping.ToApiResponseSaldo(res)

	logSuccess("Successfully retrieve saldo data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Retrieve all active saldo data
// @Tags Saldo
// @Security Bearer
// @Description Retrieve a list of all active saldo data
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsesSaldo "List of saldo data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve saldo data"
// @Router /api/saldos/active [get]
func (h *saldoHandleApi) FindByActive(c echo.Context) error {
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

	req := &pb.FindAllSaldoRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.saldo.FindByActive(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve saldo data", err, zap.Error(err))
		return saldo_errors.ErrApiFailedFindAllSaldoActive(c)
	}

	so := h.mapping.ToApiResponsePaginationSaldoDeleteAt(res)

	logSuccess("Successfully retrieve saldo data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Retrieve trashed saldo data
// @Tags Saldo
// @Security Bearer
// @Description Retrieve a list of all trashed saldo data
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsesSaldo "List of trashed saldo data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve saldo data"
// @Router /api/saldos/trashed [get]
func (h *saldoHandleApi) FindByTrashed(c echo.Context) error {
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

	req := &pb.FindAllSaldoRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.saldo.FindByTrashed(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve trashed saldo data", err, zap.Error(err))

		return saldo_errors.ErrApiFailedFindAllSaldoTrashed(c)
	}

	so := h.mapping.ToApiResponsePaginationSaldoDeleteAt(res)

	logSuccess("Successfully retrieve trashed saldo data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Create a new saldo
// @Tags Saldo
// @Security Bearer
// @Description Create a new saldo record with the provided card number and total balance.
// @Accept json
// @Produce json
// @Param CreateSaldoRequest body requests.CreateSaldoRequest true "Create Saldo Request"
// @Success 200 {object} response.ApiResponseSaldo "Successfully created saldo record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create saldo"
// @Router /api/saldos/create [post]
func (h *saldoHandleApi) Create(c echo.Context) error {
	const method = "Create"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	var body requests.CreateSaldoRequest

	if err := c.Bind(&body); err != nil {
		status = "error"

		logError("Failed to bind CreateSaldo request", err, zap.Error(err))

		return saldo_errors.ErrApiBindCreateSaldo(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"

		logError("Failed to validate CreateSaldo request", err, zap.Error(err))

		return saldo_errors.ErrApiValidateCreateSaldo(c)
	}

	res, err := h.saldo.CreateSaldo(ctx, &pb.CreateSaldoRequest{
		CardNumber:   body.CardNumber,
		TotalBalance: int32(body.TotalBalance),
	})

	if err != nil {
		status = "error"

		logError("Failed to create saldo", err, zap.Error(err))

		return saldo_errors.ErrApiFailedCreateSaldo(c)
	}

	so := h.mapping.ToApiResponseSaldo(res)

	logSuccess("Successfully create saldo", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Update an existing saldo
// @Tags Saldo
// @Security Bearer
// @Description Update an existing saldo record with the provided card number and total balance.
// @Accept json
// @Produce json
// @Param id path int true "Saldo ID"
// @Param UpdateSaldoRequest body requests.UpdateSaldoRequest true "Update Saldo Request"
// @Success 200 {object} response.ApiResponseSaldo "Successfully updated saldo record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update saldo"
// @Router /api/saldos/update/{id} [post]
func (h *saldoHandleApi) Update(c echo.Context) error {
	const method = "Update"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	idint, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Bad Request: Invalid ID", err, zap.Error(err))

		return saldo_errors.ErrApiInvalidSaldoID(c)
	}

	var body requests.UpdateSaldoRequest

	if err := c.Bind(&body); err != nil {
		status = "error"

		logError("Failed to bind UpdateSaldo request", err, zap.Error(err))

		return saldo_errors.ErrApiBindUpdateSaldo(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"

		logError("Failed to validate UpdateSaldo request", err, zap.Error(err))

		return saldo_errors.ErrApiValidateUpdateSaldo(c)
	}

	res, err := h.saldo.UpdateSaldo(ctx, &pb.UpdateSaldoRequest{
		SaldoId:      int32(idint),
		CardNumber:   body.CardNumber,
		TotalBalance: int32(body.TotalBalance),
	})

	if err != nil {
		status = "error"

		logError("Failed to update saldo", err, zap.Error(err))

		return saldo_errors.ErrApiFailedUpdateSaldo(c)
	}

	so := h.mapping.ToApiResponseSaldo(res)

	logSuccess("Successfully update saldo", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Soft delete a saldo
// @Tags Saldo
// @Security Bearer
// @Description Soft delete an existing saldo record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Saldo ID"
// @Success 200 {object} response.ApiResponseSaldo "Successfully trashed saldo record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to trashed saldo"
// @Router /api/saldos/trashed/{id} [post]
func (h *saldoHandleApi) TrashSaldo(c echo.Context) error {
	const method = "TrashSaldo"
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

		return saldo_errors.ErrApiInvalidSaldoID(c)
	}

	res, err := h.saldo.TrashedSaldo(ctx, &pb.FindByIdSaldoRequest{
		SaldoId: int32(idInt),
	})

	if err != nil {
		status = "error"

		logError("Failed to trashed saldo", err, zap.Error(err))

		return saldo_errors.ErrApiFailedTrashSaldo(c)
	}

	so := h.mapping.ToApiResponseSaldo(res)

	logSuccess("Successfully trashed saldo", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Restore a trashed saldo
// @Tags Saldo
// @Security Bearer
// @Description Restore an existing saldo record from the trash by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Saldo ID"
// @Success 200 {object} response.ApiResponseSaldo "Successfully restored saldo record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore saldo"
// @Router /api/saldos/restore/{id} [post]
func (h *saldoHandleApi) RestoreSaldo(c echo.Context) error {
	const method = "RestoreSaldo"
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

		return saldo_errors.ErrApiInvalidSaldoID(c)
	}

	res, err := h.saldo.RestoreSaldo(ctx, &pb.FindByIdSaldoRequest{
		SaldoId: int32(idInt),
	})

	if err != nil {
		status = "error"

		logError("Failed to restore saldo", err, zap.Error(err))

		return saldo_errors.ErrApiFailedRestoreSaldo(c)
	}

	so := h.mapping.ToApiResponseSaldo(res)

	logSuccess("Successfully restored saldo", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Permanently delete a saldo
// @Tags Saldo
// @Security Bearer
// @Description Permanently delete an existing saldo record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Saldo ID"
// @Success 200 {object} response.ApiResponseSaldoDelete "Successfully deleted saldo record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete saldo"
// @Router /api/saldos/permanent/{id} [delete]
func (h *saldoHandleApi) Delete(c echo.Context) error {
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

		logError("Bad Request: Invalid ID", err, zap.Error(err))

		return saldo_errors.ErrApiInvalidSaldoID(c)
	}

	res, err := h.saldo.DeleteSaldoPermanent(ctx, &pb.FindByIdSaldoRequest{
		SaldoId: int32(idInt),
	})

	if err != nil {
		status = "error"

		logError("Failed to delete saldo", err, zap.Error(err))

		return saldo_errors.ErrApiFailedDeleteSaldoPermanent(c)
	}

	so := h.mapping.ToApiResponseSaldoDelete(res)

	logSuccess("Successfully deleted saldo", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// RestoreAllSaldo restores all saldo records.
// @Summary Restore all saldo records
// @Tags Saldo
// @Security Bearer
// @Description Restore all saldo records that were previously deleted.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseSaldoAll "Successfully restored all saldo records"
// @Failure 500 {object} response.ErrorResponse "Failed to restore all saldo records"
// @Router /api/saldos/restore/all [post]
func (h *saldoHandleApi) RestoreAllSaldo(c echo.Context) error {
	const method = "RestoreAllSaldo"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	res, err := h.saldo.RestoreAllSaldo(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"

		logError("Failed to restore all saldo", err, zap.Error(err))

		return saldo_errors.ErrApiFailedRestoreAllSaldo(c)
	}

	so := h.mapping.ToApiResponseSaldoAll(res)

	logSuccess("Successfully restored all saldo", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Permanently delete all saldo records
// @Tags Saldo
// @Security Bearer
// @Description Permanently delete all saldo records from the database.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseSaldoAll "Successfully deleted all saldo records permanently"
// @Failure 500 {object} response.ErrorResponse "Failed to permanently delete all saldo records"
// @Router /api/saldos/permanent/all [post]
func (h *saldoHandleApi) DeleteAllSaldoPermanent(c echo.Context) error {
	const method = "DeleteAllSaldoPermanent"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	res, err := h.saldo.DeleteAllSaldoPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"

		logError("Failed to delete all saldo permanently", err, zap.Error(err))

		return saldo_errors.ErrApiFailedDeleteAllSaldoPermanent(c)
	}

	so := h.mapping.ToApiResponseSaldoAll(res)

	logSuccess("Successfully deleted all saldo", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *saldoHandleApi) startTracingAndLogging(
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

func (s *saldoHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
