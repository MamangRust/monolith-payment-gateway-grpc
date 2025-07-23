package transferhandler

import (
	"context"
	"net/http"
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

type transferStatsStatusHandleApi struct {
	client pb.TransferStatsStatusServiceClient

	logger logger.LoggerInterface

	mapper apimapper.TransferStatsStatusResponseMapper

	trace trace.Tracer

	requestCounter *prometheus.CounterVec

	requestDuration *prometheus.HistogramVec
}

type transferStatsStatusHandleDeps struct {
	client pb.TransferStatsStatusServiceClient

	router *echo.Echo

	logger logger.LoggerInterface

	mapper apimapper.TransferStatsStatusResponseMapper
}

func NewTransferStatsStatusHandleApi(params *transferStatsStatusHandleDeps) *transferStatsStatusHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transfer_stats_status_handler_requests_total",
			Help: "Total number of transfer stats status requests",
		},
		[]string{"status", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transfer_stats_status_handler_request_duration_seconds",
			Help:    "Duration of transfer stats status requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"status", "status"},
	)

	transferStatsStatusHandleApi := &transferStatsStatusHandleApi{
		client:          params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("transfer-stats-status-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerTransfer := params.router.Group("/api/transfer-stats-status")

	routerTransfer.GET("/monthly-success", transferStatsStatusHandleApi.FindMonthlyTransferStatusSuccess)
	routerTransfer.GET("/yearly-success", transferStatsStatusHandleApi.FindYearlyTransferStatusSuccess)
	routerTransfer.GET("/monthly-failed", transferStatsStatusHandleApi.FindMonthlyTransferStatusFailed)
	routerTransfer.GET("/yearly-failed", transferStatsStatusHandleApi.FindYearlyTransferStatusFailed)

	routerTransfer.GET("/monthly-success-by-card", transferStatsStatusHandleApi.FindMonthlyTransferStatusSuccessByCardNumber)
	routerTransfer.GET("/yearly-success-by-card", transferStatsStatusHandleApi.FindYearlyTransferStatusSuccessByCardNumber)
	routerTransfer.GET("/monthly-failed-by-card", transferStatsStatusHandleApi.FindMonthlyTransferStatusFailedByCardNumber)
	routerTransfer.GET("/yearly-failed-by-card", transferStatsStatusHandleApi.FindYearlyTransferStatusFailedByCardNumber)

	return transferStatsStatusHandleApi
}

// FindMonthlyTransferStatusSuccess retrieves the monthly transfer status for successful transactions.
// @Summary Get monthly transfer status for successful transactions
// @Tags Transfer Stats Status
// @Security Bearer
// @Description Retrieve the monthly transfer status for successful transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseTransferMonthStatusSuccess "Monthly transfer status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transfer status for successful transactions"
// @Router /api/transfer-stats-status/monthly-success [get]
func (h *transferStatsStatusHandleApi) FindMonthlyTransferStatusSuccess(c echo.Context) error {
	const method = "FindMonthlyTransferStatusSuccess"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	month, err := shared.ParseQueryMonth(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindMonthlyTransferStatusSuccess(ctx, &pb.FindMonthlyTransferStatus{
		Year:  int32(year),
		Month: int32(month),
	})

	if err != nil {
		logError("Failed to retrieve monthly Transfer status success", err, zap.Error(err))

		return transfer_errors.ErrApiFailedFindMonthlyTransferStatusSuccess(c)
	}

	so := h.mapper.ToApiResponseTransferMonthStatusSuccess(res)

	logSuccess("Successfully retrieved monthly Transfer status success", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTransferStatusSuccess retrieves the yearly transfer status for successful transactions.
// @Summary Get yearly transfer status for successful transactions
// @Tags Transfer Stats Status
// @Security Bearer
// @Description Retrieve the yearly transfer status for successful transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransferYearStatusSuccess "Yearly transfer status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transfer status for successful transactions"
// @Router /api/transfer-stats-status/yearly-success [get]
func (h *transferStatsStatusHandleApi) FindYearlyTransferStatusSuccess(c echo.Context) error {
	const method = "FindYearlyTransferStatusSuccess"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindYearlyTransferStatusSuccess(ctx, &pb.FindYearTransferStatus{
		Year: int32(year),
	})

	if err != nil {
		logError("Failed to retrieve yearly Transfer status success", err, zap.Error(err))

		return transfer_errors.ErrApiFailedFindYearlyTransferStatusSuccess(c)
	}

	so := h.mapper.ToApiResponseTransferYearStatusSuccess(res)

	logSuccess("Successfully retrieved yearly Transfer status success", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTransferStatusFailed retrieves the monthly transfer status for failed transactions.
// @Summary Get monthly transfer status for failed transactions
// @Tags Transfer Stats Status
// @Security Bearer
// @Description Retrieve the monthly transfer status for failed transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseTransferMonthStatusFailed "Monthly transfer status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transfer status for failed transactions"
// @Router /api/transfer-stats-status/monthly-failed [get]
func (h *transferStatsStatusHandleApi) FindMonthlyTransferStatusFailed(c echo.Context) error {
	const method = "FindMonthlyTransferStatusFailed"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	month, err := shared.ParseQueryMonth(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindMonthlyTransferStatusFailed(ctx, &pb.FindMonthlyTransferStatus{
		Year:  int32(year),
		Month: int32(month),
	})

	if err != nil {
		logError("Failed to retrieve monthly Transfer status Failed", err, zap.Error(err))

		return transfer_errors.ErrApiFailedFindMonthlyTransferStatusFailed(c)
	}

	so := h.mapper.ToApiResponseTransferMonthStatusFailed(res)

	logSuccess("Successfully retrieved monthly Transfer status Failed", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTransferStatusFailed retrieves the yearly transfer status for failed transactions.
// @Summary Get yearly transfer status for failed transactions
// @Tags Transfer Stats Status
// @Security Bearer
// @Description Retrieve the yearly transfer status for failed transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransferYearStatusFailed "Yearly transfer status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transfer status for failed transactions"
// @Router /api/transfer-stats-status/yearly-failed [get]
func (h *transferStatsStatusHandleApi) FindYearlyTransferStatusFailed(c echo.Context) error {
	const method = "FindYearlyTransferStatusFailed"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindYearlyTransferStatusFailed(ctx, &pb.FindYearTransferStatus{
		Year: int32(year),
	})

	if err != nil {
		logError("Failed to retrieve yearly Transfer status Failed", err, zap.Error(err))

		return transfer_errors.ErrApiFailedFindYearlyTransferStatusFailed(c)
	}

	so := h.mapper.ToApiResponseTransferYearStatusFailed(res)

	logSuccess("Successfully retrieved yearly Transfer status Failed", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTransferStatusSuccessByCardNumber retrieves the monthly transfer status for successful transactions.
// @Summary Get monthly transfer status for successful transactions
// @Tags Transfer Stats Status
// @Security Bearer
// @Description Retrieve the monthly transfer status for successful transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseTransferMonthStatusSuccess "Monthly transfer status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transfer status for successful transactions"
// @Router /api/transfer-stats-status/monthly-success-by-card [get]
func (h *transferStatsStatusHandleApi) FindMonthlyTransferStatusSuccessByCardNumber(c echo.Context) error {
	const method = "FindMonthlyTransferStatusSuccessByCardNumber"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	cardNumber, err := shared.ParseQueryCard(c, h.logger)

	if err != nil {
		return err
	}

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	month, err := shared.ParseQueryMonth(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindMonthlyTransferStatusSuccessByCardNumber(ctx, &pb.FindMonthlyTransferStatusCardNumber{
		Year:       int32(year),
		Month:      int32(month),
		CardNumber: cardNumber,
	})

	if err != nil {
		logError("Failed to retrieve monthly Transfer status success", err, zap.Error(err))

		return transfer_errors.ErrApiFailedFindMonthlyTransferStatusSuccessByCardNumber(c)
	}

	so := h.mapper.ToApiResponseTransferMonthStatusSuccess(res)

	logSuccess("Successfully retrieved monthly Transfer status success", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTransferStatusSuccessByCardNumber retrieves the yearly transfer status for successful transactions.
// @Summary Get yearly transfer status for successful transactions
// @Tags Transfer Stats Status
// @Security Bearer
// @Description Retrieve the yearly transfer status for successful transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseTransferYearStatusSuccess "Yearly transfer status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transfer status for successful transactions"
// @Router /api/transfer-stats-status/yearly-success-by-card [get]
func (h *transferStatsStatusHandleApi) FindYearlyTransferStatusSuccessByCardNumber(c echo.Context) error {
	const method = "FindYearlyTransferStatusSuccessByCardNumber"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	cardNumber, err := shared.ParseQueryCard(c, h.logger)

	if err != nil {
		return err
	}

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindYearlyTransferStatusSuccessByCardNumber(ctx, &pb.FindYearTransferStatusCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	})

	if err != nil {
		logError("Failed to retrieve yearly Transfer status success", err, zap.Error(err))

		return transfer_errors.ErrApiFailedFindYearlyTransferStatusSuccessByCardNumber(c)
	}

	so := h.mapper.ToApiResponseTransferYearStatusSuccess(res)

	logSuccess("Successfully retrieved yearly Transfer status success", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTransferStatusFailedByCardNumber retrieves the monthly transfer status for failed transactions.
// @Summary Get monthly transfer status for failed transactions
// @Tags Transfer Stats Status
// @Security Bearer
// @Description Retrieve the monthly transfer status for failed transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseTransferMonthStatusFailed "Monthly transfer status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transfer status for failed transactions"
// @Router /api/transfer-stats-status/monthly-failed-by-card [get]
func (h *transferStatsStatusHandleApi) FindMonthlyTransferStatusFailedByCardNumber(c echo.Context) error {
	const method = "FindMonthlyTransferStatusFailedByCardNumber"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	cardNumber, err := shared.ParseQueryCard(c, h.logger)

	if err != nil {
		return err
	}

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	month, err := shared.ParseQueryMonth(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindMonthlyTransferStatusFailedByCardNumber(ctx, &pb.FindMonthlyTransferStatusCardNumber{
		Year:       int32(year),
		Month:      int32(month),
		CardNumber: cardNumber,
	})

	if err != nil {
		logError("Failed to retrieve monthly Transfer status failed", err, zap.Error(err))

		return transfer_errors.ErrApiFailedFindMonthlyTransferStatusFailedByCardNumber(c)
	}

	so := h.mapper.ToApiResponseTransferMonthStatusFailed(res)

	logSuccess("Successfully retrieved monthly Transfer status failed", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTransferStatusFailedByCardNumber retrieves the yearly transfer status for failed transactions.
// @Summary Get yearly transfer status for failed transactions
// @Tags Transfer Stats Status
// @Security Bearer
// @Description Retrieve the yearly transfer status for failed transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseTransferYearStatusFailed "Yearly transfer status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transfer status for failed transactions"
// @Router /api/transfer-stats-status/yearly-failed-by-card [get]
func (h *transferStatsStatusHandleApi) FindYearlyTransferStatusFailedByCardNumber(c echo.Context) error {
	const method = "FindYearlyTransferStatusFailedByCardNumber"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	cardNumber, err := shared.ParseQueryCard(c, h.logger)

	if err != nil {
		return err
	}

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindYearlyTransferStatusFailedByCardNumber(ctx, &pb.FindYearTransferStatusCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	})

	if err != nil {
		logError("Failed to retrieve yearly Transfer status failed", err, zap.Error(err))

		return transfer_errors.ErrApiFailedFindYearlyTransferStatusFailedByCardNumber(c)
	}

	so := h.mapper.ToApiResponseTransferYearStatusFailed(res)

	logSuccess("Successfully retrieved yearly Transfer status failed", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *transferStatsStatusHandleApi) startTracingAndLogging(
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

func (s *transferStatsStatusHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
