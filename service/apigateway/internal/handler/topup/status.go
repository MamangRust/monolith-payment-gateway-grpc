package topuphandler

import (
	"context"
	"net/http"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-apigateway/internal/shared"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/topup"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/api"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/topup"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type topupStatsStatusHandleApi struct {
	client pb.TopupStatsStatusServiceClient

	logger logger.LoggerInterface

	mapper apimapper.TopupStatsStatusResponseMapper

	trace trace.Tracer

	requestCounter *prometheus.CounterVec

	requestDuration *prometheus.HistogramVec
}

type topupStatsStatusHandleDeps struct {
	client pb.TopupStatsStatusServiceClient

	router *echo.Echo

	logger logger.LoggerInterface

	mapper apimapper.TopupStatsStatusResponseMapper
}

func NewTopupStatsStatusHandleApi(params *topupStatsStatusHandleDeps) *topupStatsStatusHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "topup_stats_status_handler_requests_total",
			Help: "Total number of topup stats status requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "topup_stats_status_handler_request_duration_seconds",
			Help:    "Duration of topup stats status requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	topupHandler := &topupStatsStatusHandleApi{
		client:          params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("topup-stats-status-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerTopup := params.router.Group("/api/topup-stats-status")

	routerTopup.GET("/monthly-success", topupHandler.FindMonthlyTopupStatusSuccess)
	routerTopup.GET("/yearly-success", topupHandler.FindYearlyTopupStatusSuccess)
	routerTopup.GET("/monthly-failed", topupHandler.FindMonthlyTopupStatusFailed)
	routerTopup.GET("/yearly-failed", topupHandler.FindYearlyTopupStatusFailed)

	routerTopup.GET("/monthly-success-by-card", topupHandler.FindMonthlyTopupStatusSuccessByCardNumber)
	routerTopup.GET("/yearly-success-by-card", topupHandler.FindYearlyTopupStatusSuccessByCardNumber)
	routerTopup.GET("/monthly-failed-by-card", topupHandler.FindMonthlyTopupStatusFailedByCardNumber)
	routerTopup.GET("/yearly-failed-by-card", topupHandler.FindYearlyTopupStatusFailedByCardNumber)

	return topupHandler
}

// FindMonthlyTopupStatusSuccess retrieves the monthly top-up status for successful transactions.
// @Summary Get monthly top-up status for successful transactions
// @Tags Topup Stats Status
// @Security Bearer
// @Description Retrieve the monthly top-up status for successful transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseTopupMonthStatusSuccess "Monthly top-up status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly top-up status for successful transactions"
// @Router /api/topup-stats-status/monthly-success [get]
func (h *topupStatsStatusHandleApi) FindMonthlyTopupStatusSuccess(c echo.Context) error {
	const method = "FindMonthlyTopupStatusSuccess"
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

	res, err := h.client.FindMonthlyTopupStatusSuccess(ctx, &pb.FindMonthlyTopupStatus{
		Year:  int32(year),
		Month: int32(month),
	})

	if err != nil {
		logError("Failed to retrieve monthly topup status success", err, zap.Error(err))

		return topup_errors.ErrApiFailedFindMonthlyTopupStatusSuccess(c)
	}

	so := h.mapper.ToApiResponseTopupMonthStatusSuccess(res)

	logSuccess("success find monthly topup status success", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTopupStatusSuccess retrieves the yearly top-up status for successful transactions.
// @Summary Get yearly top-up status for successful transactions
// @Tags Topup Stats Status
// @Security Bearer
// @Description Retrieve the yearly top-up status for successful transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupYearStatusSuccess "Yearly top-up status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly top-up status for successful transactions"
// @Router /api/topup-stats-status/yearly-success [get]
func (h *topupStatsStatusHandleApi) FindYearlyTopupStatusSuccess(c echo.Context) error {
	const method = "FindYearlyTopupStatusSuccess"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindYearlyTopupStatusSuccess(ctx, &pb.FindYearTopupStatus{
		Year: int32(year),
	})

	if err != nil {
		logError("Failed to retrieve yearly topup status success", err, zap.Error(err))

		return topup_errors.ErrApiFailedFindYearlyTopupStatusSuccess(c)
	}

	so := h.mapper.ToApiResponseTopupYearStatusSuccess(res)

	logSuccess("success find yearly topup status success", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTopupStatusFailed retrieves the monthly top-up status for failed transactions.
// @Summary Get monthly top-up status for failed transactions
// @Tags Topup Stats Status
// @Security Bearer
// @Description Retrieve the monthly top-up status for failed transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseTopupMonthStatusFailed "Monthly top-up status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly top-up status for failed transactions"
// @Router /api/topup-stats-status/monthly-failed [get]
func (h *topupStatsStatusHandleApi) FindMonthlyTopupStatusFailed(c echo.Context) error {
	const method = "FindMonthlyTopupStatusFailed"
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

	res, err := h.client.FindMonthlyTopupStatusFailed(ctx, &pb.FindMonthlyTopupStatus{
		Year:  int32(year),
		Month: int32(month),
	})

	if err != nil {
		logError("Failed to retrieve monthly topup status failed", err, zap.Error(err))

		return topup_errors.ErrApiFailedFindMonthlyTopupStatusFailed(c)
	}

	so := h.mapper.ToApiResponseTopupMonthStatusFailed(res)

	logSuccess("success find monthly topup status failed", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTopupStatusFailed retrieves the yearly top-up status for failed transactions.
// @Summary Get yearly top-up status for failed transactions
// @Tags Topup Stats Status
// @Security Bearer
// @Description Retrieve the yearly top-up status for failed transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupYearStatusFailed "Yearly top-up status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly top-up status for failed transactions"
// @Router /api/topup-stats-status/yearly-failed [get]
func (h *topupStatsStatusHandleApi) FindYearlyTopupStatusFailed(c echo.Context) error {
	const method = "FindYearlyTopupStatusFailed"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindYearlyTopupStatusFailed(ctx, &pb.FindYearTopupStatus{
		Year: int32(year),
	})

	if err != nil {
		logError("Failed to retrieve yearly topup status failed", err, zap.Error(err))

		return topup_errors.ErrApiFailedFindYearlyTopupStatusFailed(c)
	}

	so := h.mapper.ToApiResponseTopupYearStatusFailed(res)

	logSuccess("success find yearly topup status failed", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTopupStatusSuccess retrieves the monthly top-up status for successful transactions.
// @Summary Get monthly top-up status for successful transactions
// @Tags Topup Stats Status
// @Security Bearer
// @Description Retrieve the monthly top-up status for successful transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseTopupMonthStatusSuccess "Monthly top-up status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly top-up status for successful transactions"
// @Router /api/topup-stats-status/monthly-success-by-card [get]
func (h *topupStatsStatusHandleApi) FindMonthlyTopupStatusSuccessByCardNumber(c echo.Context) error {
	const method = "FindMonthlyTopupStatusSuccessByCardNumber"
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

	cardNumber, err := shared.ParseQueryCard(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindMonthlyTopupStatusSuccessByCardNumber(ctx, &pb.FindMonthlyTopupStatusCardNumber{
		Year:       int32(year),
		Month:      int32(month),
		CardNumber: cardNumber,
	})

	if err != nil {
		logError("Failed to retrieve monthly topup status success", err, zap.Error(err))

		return topup_errors.ErrApiFailedFindMonthlyTopupStatusSuccess(c)
	}

	so := h.mapper.ToApiResponseTopupMonthStatusSuccess(res)

	logSuccess("success find monthly topup status success", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTopupStatusSuccess retrieves the yearly top-up status for successful transactions.
// @Summary Get yearly top-up status for successful transactions
// @Tags Topup Stats Status
// @Security Bearer
// @Description Retrieve the yearly top-up status for successful transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseTopupYearStatusSuccess "Yearly top-up status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly top-up status for successful transactions"
// @Router /api/topup-stats-status/yearly-success-by-card [get]
func (h *topupStatsStatusHandleApi) FindYearlyTopupStatusSuccessByCardNumber(c echo.Context) error {
	const method = "FindYearlyTopupStatusSuccessByCardNumber"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	cardNumber, err := shared.ParseQueryCard(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindYearlyTopupStatusSuccessByCardNumber(ctx, &pb.FindYearTopupStatusCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	})

	if err != nil {
		logError("Failed to retrieve yearly topup status success", err, zap.Error(err))

		return topup_errors.ErrApiFailedFindYearlyTopupStatusSuccess(c)
	}

	so := h.mapper.ToApiResponseTopupYearStatusSuccess(res)

	logSuccess("success find yearly topup status success", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTopupStatusFailed retrieves the monthly top-up status for failed transactions.
// @Summary Get monthly top-up status for failed transactions
// @Tags Topup Stats Status
// @Security Bearer
// @Description Retrieve the monthly top-up status for failed transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseTopupMonthStatusFailed "Monthly top-up status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly top-up status for failed transactions"
// @Router /api/topup-stats-status/monthly-failed-by-card [get]
func (h *topupStatsStatusHandleApi) FindMonthlyTopupStatusFailedByCardNumber(c echo.Context) error {
	const method = "FindMonthlyTopupStatusFailedByCardNumber"
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

	cardNumber, err := shared.ParseQueryCard(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindMonthlyTopupStatusFailedByCardNumber(ctx, &pb.FindMonthlyTopupStatusCardNumber{
		Year:       int32(year),
		Month:      int32(month),
		CardNumber: cardNumber,
	})

	if err != nil {
		logError("Failed to retrieve monthly topup status failed", err, zap.Error(err))

		return topup_errors.ErrApiFailedFindMonthlyTopupStatusFailed(c)
	}

	so := h.mapper.ToApiResponseTopupMonthStatusFailed(res)

	logSuccess("success find monthly topup status failed", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTopupStatusFailedByCardNumber retrieves the yearly top-up status for failed transactions.
// @Summary Get yearly top-up status for failed transactions
// @Tags Topup Stats Status
// @Security Bearer
// @Description Retrieve the yearly top-up status for failed transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseTopupYearStatusFailed "Yearly top-up status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly top-up status for failed transactions"
// @Router /api/topup-stats-status/yearly-failed-by-card [get]
func (h *topupStatsStatusHandleApi) FindYearlyTopupStatusFailedByCardNumber(c echo.Context) error {
	const method = "FindYearlyTopupStatusFailedByCardNumber"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	cardNumber, err := shared.ParseQueryCard(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindYearlyTopupStatusFailedByCardNumber(ctx, &pb.FindYearTopupStatusCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	})

	if err != nil {
		logError("Failed to retrieve yearly topup status failed", err, zap.Error(err))

		return topup_errors.ErrApiFailedFindYearlyTopupStatusFailed(c)
	}

	so := h.mapper.ToApiResponseTopupYearStatusFailed(res)

	logSuccess("success find yearly topup status failed", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *topupStatsStatusHandleApi) startTracingAndLogging(
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

func (s *topupStatsStatusHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
