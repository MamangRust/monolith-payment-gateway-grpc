package cardhandler

import (
	"context"
	"net/http"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-apigateway/internal/shared"
	pbhelper "github.com/MamangRust/monolith-payment-gateway-pb"
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
	"google.golang.org/protobuf/types/known/emptypb"
)

type cardDashboardHandleApi struct {
	// card is the gRPC client used to interact with the CardService.
	card pb.CardDashboardServiceClient

	// logger provides structured and leveled logging capabilities.
	logger logger.LoggerInterface

	// mapper transforms gRPC responses into standardized HTTP API responses.
	mapper apimapper.CardDashboardResponseMapper

	// trace is the OpenTelemetry tracer for distributed tracing.
	trace trace.Tracer

	// requestCounter records the number of HTTP requests handled by this service.
	requestCounter *prometheus.CounterVec

	// requestDuration records the duration of HTTP request handling in seconds.
	requestDuration *prometheus.HistogramVec
}

// cardDashboardHandleApiDeps contains the necessary dependencies for the cardDashboardHandleApi.
type cardDashboardHandleApiDeps struct {
	// client is the gRPC client used to interact with the CardService.
	client pb.CardDashboardServiceClient

	// router is the Echo HTTP router used to register routes.
	router *echo.Echo

	// logger provides structured and leveled logging capabilities.
	logger logger.LoggerInterface

	// mapper transforms gRPC responses into standardized HTTP API responses.
	mapper apimapper.CardDashboardResponseMapper
}

// NewCardDashboardHandleApi initializes a new cardDashboardHandleApi and sets up the routes for
// card dashboard-related operations.
//
// This function registers various HTTP endpoints related to card dashboard management, including
// retrieval of dashboard card data. It also tracks metrics like request count and duration using
// Prometheus metrics. The routes are grouped under "/api/card-dashboard".
//
// Parameters:
// - params: A pointer to cardDashboardHandleApiDeps containing the necessary dependencies.
//
// Returns:
// - A pointer to a newly created cardDashboardHandleApi.
func NewCardDashboardHandleApi(
	params *cardDashboardHandleApiDeps,
) *cardDashboardHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "card_dashboard_handler_requests_total",
			Help: "Total number of Card Dashboard requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "card_dashboard_handler_request_duration_seconds",
			Help:    "Duration of Card Dashboard requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	cardDashboardHandler := &cardDashboardHandleApi{
		card:            params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("card-dashboard-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerCard := params.router.Group("/api/card-dashboard")

	routerCard.GET("", cardDashboardHandler.DashboardCard)
	routerCard.GET("/:cardNumber", cardDashboardHandler.DashboardCardCardNumber)

	return cardDashboardHandler
}

// DashboardCard godoc
// @Summary Get dashboard card data
// @Description Retrieve dashboard card data
// @Tags Card-Dashboard
// @Security Bearer
// @Produce json
// @Success 200 {object} response.ApiResponseDashboardCard
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/dashboard [get]
func (h *cardDashboardHandleApi) DashboardCard(c echo.Context) error {
	const method = "DashboardCard"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	res, err := h.card.DashboardCard(ctx, &emptypb.Empty{})

	if err != nil {
		logError("Failed to retrieve dashboard card data", err, zap.Error(err))

		return cardapierrors.ErrApiFailedDashboardCard(c)
	}

	so := h.mapper.ToApiResponseDashboardCard(res)

	logSuccess("Success retrieve dashboard card data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// DashboardCardCardNumber godoc
// @Summary Get dashboard card data by card number
// @Description Retrieve dashboard card data for a specific card number
// @Tags Card-Dashboard
// @Security Bearer
// @Produce json
// @Param cardNumber path string true "Card Number"
// @Success 200 {object} response.ApiResponseDashboardCardNumber
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/dashboard/{cardNumber} [get]
func (h *cardDashboardHandleApi) DashboardCardCardNumber(c echo.Context) error {
	const method = "DashboardCardCardNumber"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	cardNumber, err := shared.ParseQueryCard(c, h.logger)

	if err != nil {
		return err
	}

	req := &pbhelper.FindByCardNumberRequest{
		CardNumber: cardNumber,
	}

	res, err := h.card.DashboardCardNumber(ctx, req)

	if err != nil {
		logError("Failed to retrieve dashboard card data by card number", err, zap.Error(err))

		return cardapierrors.ErrApiFailedDashboardCardByCardNumber(c)
	}

	so := h.mapper.ToApiResponseDashboardCardCardNumber(res)

	logSuccess("Success retrieve dashboard card data by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *cardDashboardHandleApi) startTracingAndLogging(
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

func (s *cardDashboardHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
