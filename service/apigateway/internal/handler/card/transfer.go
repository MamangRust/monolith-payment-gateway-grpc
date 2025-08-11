package cardhandler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-apigateway/internal/shared"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/card"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	cardapierrors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/api"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/card"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type cardStatsTransferHandleApi struct {
	// card is the gRPC client used to interact with the CardService.
	card pb.CardStatsTransferServiceClient

	// logger provides structured and leveled logging capabilities.
	logger logger.LoggerInterface

	// mapper transforms gRPC responses into standardized HTTP API responses.
	mapper apimapper.CardStatsAmountResponseMapper

	// trace is the OpenTelemetry tracer for distributed tracing.
	trace trace.Tracer

	// requestCounter records the number of HTTP requests handled by this service.
	requestCounter *prometheus.CounterVec

	// requestDuration records the duration of HTTP request handling in seconds.
	requestDuration *prometheus.HistogramVec
}

// cardStatsTransferHandleApiDeps holds the dependencies required to create a new cardStatsTransferHandleApi.
type cardStatsTransferHandleApiDeps struct {
	// client is the gRPC client used to interact with the CardService.
	client pb.CardStatsTransferServiceClient

	// router is the Echo router used to register HTTP routes.
	router *echo.Echo

	// logger provides structured and leveled logging capabilities.
	logger logger.LoggerInterface

	// mapper transforms gRPC responses into standardized HTTP API responses.
	mapper apimapper.CardStatsAmountResponseMapper
}

// NewCardStatsTransferHandleApi initializes a new cardStatsTransferHandleApi and sets up the routes for
// card stats transfer-related operations.
//
// This function registers various HTTP endpoints related to card stats transfer management, including
// retrieval of monthly and yearly transfer amounts. It also tracks metrics like request count and
// duration using Prometheus metrics. The routes are grouped under "/api/card-stats-transfer".
//
// Parameters:
// - params: A pointer to cardStatsTransferHandleApiDeps containing the necessary dependencies.
//
// Returns:
// - A pointer to a newly created cardStatsTransferHandleApi.
func NewCardStatsTransferHandleApi(
	params *cardStatsTransferHandleApiDeps,
) *cardStatsTransferHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "card_stats_transfer_handler_requests_total",
			Help: "Total number of Card Stats Transfer requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "card_stats_transfer_handler_request_duration_seconds",
			Help:    "Duration of Card Stats Transfer requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	cardStatsTransferHandler := &cardStatsTransferHandleApi{
		card:            params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("card-stats-transfer-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerCard := params.router.Group("/api/card-stats-transfer")

	routerCard.GET("/monthly-transfer-sender-amount", cardStatsTransferHandler.FindMonthlyTransferSenderAmount)
	routerCard.GET("/yearly-transfer-sender-amount", cardStatsTransferHandler.FindYearlyTransferSenderAmount)
	routerCard.GET("/monthly-transfer-receiver-amount", cardStatsTransferHandler.FindMonthlyTransferReceiverAmount)
	routerCard.GET("/yearly-transfer-receiver-amount", cardStatsTransferHandler.FindYearlyTransferReceiverAmount)

	routerCard.GET("/monthly-transfer-sender-amount-by-card", cardStatsTransferHandler.FindMonthlyTransferSenderAmountByCardNumber)
	routerCard.GET("/yearly-transfer-sender-amount-by-card", cardStatsTransferHandler.FindYearlyTransferSenderAmountByCardNumber)
	routerCard.GET("/monthly-transfer-receiver-amount-by-card", cardStatsTransferHandler.FindMonthlyTransferReceiverAmountByCardNumber)
	routerCard.GET("/yearly-transfer-receiver-amount-by-card", cardStatsTransferHandler.FindYearlyTransferReceiverAmountByCardNumber)

	return cardStatsTransferHandler
}

// FindMonthlyTransferSenderAmount godoc
// @Summary Get monthly transfer sender amount data
// @Description Retrieve monthly transfer sender amount data for a specific year
// @Tags Card-Stats-Transfer
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMonthlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card-stats-transfer/monthly-transfer-sender-amount [get]
func (h *cardStatsTransferHandleApi) FindMonthlyTransferSenderAmount(c echo.Context) error {
	const method = "FindMonthlyTransferSenderAmount"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	req := &pb.FindYearAmount{
		Year: int32(year),
	}

	res, err := h.card.FindMonthlyTransferSenderAmount(ctx, req)

	if err != nil {
		logError("Failed to retrieve monthly transfer sender amount data", err, zap.Error(err))

		return cardapierrors.ErrApiFailedFindMonthlyTransferSenderAmount(c)
	}

	so := h.mapper.ToApiResponseMonthlyAmounts(res)

	logSuccess("Success retrieve monthly transfer sender amount data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTransferSenderAmount godoc
// @Summary Get yearly transfer sender amount data
// @Description Retrieve yearly transfer sender amount data for a specific year
// @Tags Card-Stats-Transfer
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card-stats-transfer/yearly-transfer-sender-amount [get]
func (h *cardStatsTransferHandleApi) FindYearlyTransferSenderAmount(c echo.Context) error {
	const method = "FindYearlyTransferSenderAmount"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	req := &pb.FindYearAmount{
		Year: int32(year),
	}

	res, err := h.card.FindYearlyTransferSenderAmount(ctx, req)

	if err != nil {
		logError("Failed to retrieve yearly transfer sender amount data", err, zap.Error(err))

		return cardapierrors.ErrApiFailedFindYearlyTransferSenderAmount(c)
	}

	so := h.mapper.ToApiResponseYearlyAmounts(res)

	logSuccess("Success retrieve yearly transfer sender amount data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTransferReceiverAmount godoc
// @Summary Get monthly transfer receiver amount data
// @Description Retrieve monthly transfer receiver amount data for a specific year
// @Tags Card-Stats-Transfer
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMonthlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card-stats-transfer/monthly-transfer-receiver-amount [get]
func (h *cardStatsTransferHandleApi) FindMonthlyTransferReceiverAmount(c echo.Context) error {
	const method = "FindMonthlyTransferReceiverAmount"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	req := &pb.FindYearAmount{
		Year: int32(year),
	}

	res, err := h.card.FindMonthlyTransferReceiverAmount(ctx, req)

	if err != nil {
		logError("Failed to retrieve monthly transfer receiver amount data", err, zap.Error(err))

		return cardapierrors.ErrApiFailedFindMonthlyTransferReceiverAmount(c)
	}

	so := h.mapper.ToApiResponseMonthlyAmounts(res)

	logSuccess("Success retrieve monthly transfer receiver amount data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTransferReceiverAmount godoc
// @Summary Get yearly transfer receiver amount data
// @Description Retrieve yearly transfer receiver amount data for a specific year
// @Tags Card-Stats-Transfer
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card-stats-transfer/yearly-transfer-receiver-amount [get]
func (h *cardStatsTransferHandleApi) FindYearlyTransferReceiverAmount(c echo.Context) error {
	const method = "FindYearlyTransferReceiverAmount"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	year, err := shared.ParseQueryYear(c, h.logger)

	if err != nil {
		return err
	}

	req := &pb.FindYearAmount{
		Year: int32(year),
	}

	res, err := h.card.FindYearlyTransferReceiverAmount(ctx, req)

	if err != nil {
		logError("Failed to retrieve yearly transfer receiver amount data", err, zap.Error(err))

		return cardapierrors.ErrApiFailedFindYearlyTransferReceiverAmount(c)
	}

	so := h.mapper.ToApiResponseYearlyAmounts(res)

	logSuccess("Success retrieve yearly transfer receiver amount data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTransferSenderAmountByCardNumber godoc
// @Summary Get monthly transfer sender amount data by card number
// @Description Retrieve monthly transfer sender amount data for a specific year and card number
// @Tags Card-Stats-Transfer
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseMonthlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card-stats-transfer/monthly-transfer-sender-amount-by-card [get]
func (h *cardStatsTransferHandleApi) FindMonthlyTransferSenderAmountByCardNumber(c echo.Context) error {
	const method = "FindMonthlyTransferSenderAmountByCardNumber"

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

	req := &pb.FindYearAmountCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindMonthlyTransferSenderAmountByCardNumber(ctx, req)

	if err != nil {
		logError("Failed to retrieve monthly transfer sender amount data by card number", err, zap.Error(err))

		return cardapierrors.ErrApiFailedFindMonthlyTransferSenderAmountByCard(c)
	}

	so := h.mapper.ToApiResponseMonthlyAmounts(res)

	logSuccess("Success retrieve monthly transfer sender amount data by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTransferSenderAmountByCardNumber godoc
// @Summary Get yearly transfer sender amount by card number
// @Description Retrieve the total amount sent by a specific card number in a given year
// @Tags Card-Stats-Transfer
// @Security Bearer
// @Accept json
// @Produce json
// @Param year query int true "Year for which the data is requested"
// @Param card_number query string true "Card number for which the data is requested"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card-stats-transfer/yearly-transfer-sender-amount-by-card [get]
func (h *cardStatsTransferHandleApi) FindYearlyTransferSenderAmountByCardNumber(c echo.Context) error {
	const method = "FindYearlyTransferSenderAmountByCardNumber"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		logError("Invalid year", err, zap.Error(err))

		return cardapierrors.ErrApiInvalidYear(c)
	}

	cardNumber, err := shared.ParseQueryCard(c, h.logger)

	if err != nil {
		return err
	}

	req := &pb.FindYearAmountCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindYearlyTransferSenderAmountByCardNumber(ctx, req)

	if err != nil {
		logError("Failed to retrieve yearly transfer sender amount data by card number", err, zap.Error(err))

		return cardapierrors.ErrApiFailedFindYearlyTransferSenderAmountByCard(c)
	}

	so := h.mapper.ToApiResponseYearlyAmounts(res)

	logSuccess("Success retrieve yearly transfer sender amount data by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTransferReceiverAmountByCardNumber godoc
// @Summary Get monthly transfer receiver amount by card number
// @Description Retrieve the total amount received by a specific card number in a given year, broken down by month
// @Tags Card-Stats-Transfer
// @Security Bearer
// @Accept json
// @Produce json
// @Param year query int true "Year for which the data is requested"
// @Param card_number query string true "Card number for which the data is requested"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card-stats-transfer/monthly-transfer-receiver-amount-by-card [get]
func (h *cardStatsTransferHandleApi) FindMonthlyTransferReceiverAmountByCardNumber(c echo.Context) error {
	const method = "FindMonthlyTransferReceiverAmountByCardNumber"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		logError("Invalid year", err, zap.Error(err))

		return cardapierrors.ErrApiInvalidYear(c)
	}

	cardNumber, err := shared.ParseQueryCard(c, h.logger)

	if err != nil {
		return err
	}

	req := &pb.FindYearAmountCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindMonthlyTransferReceiverAmountByCardNumber(ctx, req)

	if err != nil {
		logError("Failed to retrieve monthly transfer receiver amount data by card number", err, zap.Error(err))

		return cardapierrors.ErrApiFailedFindMonthlyTransferReceiverAmountByCard(c)
	}

	so := h.mapper.ToApiResponseMonthlyAmounts(res)

	logSuccess("Success retrieve monthly transfer receiver amount data by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTransferReceiverAmountByCardNumber godoc
// @Summary Get yearly transfer receiver amount by card number
// @Description Retrieve the total amount received by a specific card number in a given year
// @Tags Card-Stats-Transfer
// @Security Bearer
// @Accept json
// @Produce json
// @Param year query int true "Year for which the data is requested"
// @Param card_number query string true "Card number for which the data is requested"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card-stats-transfer/yearly-transfer-receiver-amount-by-card [get]
func (h *cardStatsTransferHandleApi) FindYearlyTransferReceiverAmountByCardNumber(c echo.Context) error {
	const method = "FindYearlyTransferReceiverAmountByCardNumber"
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

	req := &pb.FindYearAmountCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindYearlyTransferReceiverAmountByCardNumber(ctx, req)

	if err != nil {
		logError("Failed to retrieve yearly transfer receiver amount data by card number", err, zap.Error(err))

		return cardapierrors.ErrApiFailedFindYearlyTransferReceiverAmountByCard(c)
	}

	logSuccess("Success retrieve yearly transfer receiver amount data by card number", zap.Bool("success", true))

	so := h.mapper.ToApiResponseYearlyAmounts(res)

	return c.JSON(http.StatusOK, so)
}

// startTracingAndLogging returns three functions: end, logSuccess and logError.
// The end function must be called to end the tracing and logging.
// The logSuccess and logError functions can be used to log success and error messages.
// If an error is passed to logError, it will be recorded as an error event on the span.
func (s *cardStatsTransferHandleApi) startTracingAndLogging(
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

// recordMetrics records Prometheus metrics for the specified method and status.
// It increments the request counter for the provided method and status,
// and observes the request duration by calculating the time elapsed since the provided start time.
func (s *cardStatsTransferHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
