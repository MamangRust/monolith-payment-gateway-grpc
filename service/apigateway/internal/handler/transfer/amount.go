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

type transferStatsAmountHandleApi struct {
	client pb.TransferStatsAmountServiceClient

	logger logger.LoggerInterface

	mapper apimapper.TransferStatsAmountResponseMapper

	trace trace.Tracer

	requestCounter *prometheus.CounterVec

	requestDuration *prometheus.HistogramVec
}

type transferStatsAmountHandleDeps struct {
	client pb.TransferStatsAmountServiceClient

	router *echo.Echo

	logger logger.LoggerInterface

	mapper apimapper.TransferStatsAmountResponseMapper
}

func NewTransferStatsAmountHandleApi(params *transferStatsAmountHandleDeps) *transferStatsAmountHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transfer_stats_amount_handler_requests_total",
			Help: "Total number of transfer stats amount requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transfer_stats_amount_handler_request_duration_seconds",
			Help:    "Duration of transfer stats amount requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	transferStatsAmountHandleApi := &transferStatsAmountHandleApi{
		client:          params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("transfer-stats-amount-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerTransfer := params.router.Group("/api/transfer-stats-amount")

	routerTransfer.GET("/monthly-amount", transferStatsAmountHandleApi.FindMonthlyTransferAmounts)
	routerTransfer.GET("/yearly-amount", transferStatsAmountHandleApi.FindYearlyTransferAmounts)

	routerTransfer.GET("/monthly-by-sender", transferStatsAmountHandleApi.FindMonthlyTransferAmountsBySenderCardNumber)
	routerTransfer.GET("/monthly-by-receiver", transferStatsAmountHandleApi.FindMonthlyTransferAmountsByReceiverCardNumber)
	routerTransfer.GET("/yearly-by-sender", transferStatsAmountHandleApi.FindYearlyTransferAmountsBySenderCardNumber)
	routerTransfer.GET("/yearly-by-receiver", transferStatsAmountHandleApi.FindYearlyTransferAmountsByReceiverCardNumber)

	return transferStatsAmountHandleApi
}

// FindMonthlyTransferAmounts retrieves the monthly transfer amounts for a specific year.
// @Summary Get monthly transfer amounts
// @Tags Transfer Stats Amount
// @Security Bearer
// @Description Retrieve the monthly transfer amounts for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransferMonthAmount "Monthly transfer amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transfer amounts"
// @Router /api/transfer-stats-amount/monthly-amounts [get]
func (h *transferStatsAmountHandleApi) FindMonthlyTransferAmounts(c echo.Context) error {
	const method = "FindMonthlyTransferAmounts"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindMonthlyTransferAmounts(ctx, &pb.FindYearTransferStatus{
		Year: int32(year),
	})
	if err != nil {
		logError("Failed to retrieve monthly transfer amounts", err, zap.Error(err))

		return transfer_errors.ErrApiFailedFindMonthlyTransferAmounts(c)
	}

	so := h.mapper.ToApiResponseTransferMonthAmount(res)

	logSuccess("Successfully retrieved monthly transfer amounts", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTransferAmounts retrieves the yearly transfer amounts for a specific year.
// @Summary Get yearly transfer amounts
// @Tags Transfer Stats Amount
// @Security Bearer
// @Description Retrieve the yearly transfer amounts for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransferYearAmount "Yearly transfer amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transfer amounts"
// @Router /api/transfer-stats-amount/yearly-amounts [get]
func (h *transferStatsAmountHandleApi) FindYearlyTransferAmounts(c echo.Context) error {
	const method = "FindYearlyTransferAmounts"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	res, err := h.client.FindYearlyTransferAmounts(ctx, &pb.FindYearTransferStatus{
		Year: int32(year),
	})

	if err != nil {
		logError("Failed to retrieve yearly transfer amounts", err, zap.Error(err))

		return transfer_errors.ErrApiFailedFindYearlyTransferAmounts(c)
	}

	so := h.mapper.ToApiResponseTransferYearAmount(res)

	logSuccess("Successfully retrieved yearly transfer amounts", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTransferAmountsBySenderCardNumber retrieves the monthly transfer amounts for a specific sender card number and year.
// @Summary Get monthly transfer amounts by sender card number
// @Tags Transfer Stats Amount
// @Security Bearer
// @Description Retrieve the monthly transfer amounts for a specific sender card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Sender Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransferMonthAmount "Monthly transfer amounts by sender card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transfer amounts by sender card number"
// @Router /api/transfer-stats-amount/monthly-amounts-by-sender-card [get]
func (h *transferStatsAmountHandleApi) FindMonthlyTransferAmountsBySenderCardNumber(c echo.Context) error {
	const method = "FindMonthlyTransferAmountsBySenderCardNumber"
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

	res, err := h.client.FindMonthlyTransferAmountsBySenderCardNumber(ctx, &pb.FindByCardNumberTransferRequest{
		CardNumber: cardNumber,
		Year:       int32(year),
	})

	if err != nil {
		logError("Failed to retrieve monthly transfer amounts by sender card number", err, zap.Error(err))

		return transfer_errors.ErrApiFailedFindMonthlyTransferAmountsBySenderCardNumber(c)
	}

	so := h.mapper.ToApiResponseTransferMonthAmount(res)

	logSuccess("Successfully retrieved monthly transfer amounts by sender card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTransferAmountsByReceiverCardNumber retrieves the monthly transfer amounts for a specific receiver card number and year.
// @Summary Get monthly transfer amounts by receiver card number
// @Tags Transfer Stats Amount
// @Security Bearer
// @Description Retrieve the monthly transfer amounts for a specific receiver card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Receiver Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransferMonthAmount "Monthly transfer amounts by receiver card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly transfer amounts by receiver card number"
// @Router /api/transfer-stats-amount/monthly-amounts-by-receiver-card [get]
func (h *transferStatsAmountHandleApi) FindMonthlyTransferAmountsByReceiverCardNumber(c echo.Context) error {
	const method = "FindMonthlyTransferAmountsByReceiverCardNumber"
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

	res, err := h.client.FindMonthlyTransferAmountsByReceiverCardNumber(ctx, &pb.FindByCardNumberTransferRequest{
		CardNumber: cardNumber,
		Year:       int32(year),
	})

	if err != nil {
		logError("Failed to retrieve monthly transfer amounts by receiver card number", err, zap.Error(err))

		return transfer_errors.ErrApiFailedFindMonthlyTransferAmountsByReceiverCardNumber(c)
	}

	so := h.mapper.ToApiResponseTransferMonthAmount(res)

	logSuccess("Successfully retrieved monthly transfer amounts by receiver card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTransferAmountsBySenderCardNumber retrieves the yearly transfer amounts for a specific sender card number and year.
// @Summary Get yearly transfer amounts by sender card number
// @Tags Transfer Stats Amount
// @Security Bearer
// @Description Retrieve the yearly transfer amounts for a specific sender card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Sender Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransferYearAmount "Yearly transfer amounts by sender card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transfer amounts by sender card number"
// @Router /api/transfer-stats-amount/yearly-amounts-by-sender-card [get]
func (h *transferStatsAmountHandleApi) FindYearlyTransferAmountsBySenderCardNumber(c echo.Context) error {
	const method = "FindYearlyTransferAmountsBySenderCardNumber"
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

	res, err := h.client.FindYearlyTransferAmountsBySenderCardNumber(ctx, &pb.FindByCardNumberTransferRequest{
		CardNumber: cardNumber,
		Year:       int32(year),
	})

	if err != nil {
		logError("Failed to retrieve yearly transfer amounts by sender card number", err, zap.Error(err))

		return transfer_errors.ErrApiFailedFindYearlyTransferAmountsBySenderCardNumber(c)
	}

	so := h.mapper.ToApiResponseTransferYearAmount(res)

	logSuccess("Successfully retrieved yearly transfer amounts by sender card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTransferAmountsByReceiverCardNumber retrieves the yearly transfer amounts for a specific receiver card number and year.
// @Summary Get yearly transfer amounts by receiver card number
// @Tags Transfer Stats Amount
// @Security Bearer
// @Description Retrieve the yearly transfer amounts for a specific receiver card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Receiver Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTransferYearAmount "Yearly transfer amounts by receiver card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly transfer amounts by receiver card number"
// @Router /api/transfer-stats-amount/yearly-amounts-by-receiver-card [get]
func (h *transferStatsAmountHandleApi) FindYearlyTransferAmountsByReceiverCardNumber(c echo.Context) error {
	const method = "FindYearlyTransferAmountsByReceiverCardNumber"
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

	res, err := h.client.FindYearlyTransferAmountsByReceiverCardNumber(ctx, &pb.FindByCardNumberTransferRequest{
		CardNumber: cardNumber,
		Year:       int32(year),
	})

	if err != nil {
		logError("Failed to retrieve yearly transfer amounts by receiver card number", err, zap.Error(err))

		return transfer_errors.ErrApiFailedFindYearlyTransferAmountsByReceiverCardNumber(c)
	}

	so := h.mapper.ToApiResponseTransferYearAmount(res)

	logSuccess("Successfully retrieved yearly transfer amounts by receiver card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *transferStatsAmountHandleApi) startTracingAndLogging(
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

func (s *transferStatsAmountHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
