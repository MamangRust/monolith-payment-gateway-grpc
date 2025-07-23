package merchantdocumenthandler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-apigateway/internal/shared"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/merchantdocument"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	merchantdocument_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_document_errors/api"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/merchantdocument"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantQueryDocumentHandleApi struct {
	// merchantDocument is the gRPC client used to communicate with the MerchantDocumentService.
	merchantDocument pb.MerchantDocumentServiceClient

	// logger is used for structured and contextual logging.
	logger logger.LoggerInterface

	// mapper is responsible for converting gRPC responses to API-compliant responses.
	mapper apimapper.MerchantDocumentQueryResponseMapper

	// trace provides distributed tracing capabilities using OpenTelemetry.
	trace trace.Tracer

	// requestCounter counts the number of requests received by the shared.
	requestCounter *prometheus.CounterVec

	// requestDuration records how long each handler request takes in seconds.
	requestDuration *prometheus.HistogramVec
}

// merchantDocumentHandleDeps defines the parameters required to initialize
// the merchant document handler and register its HTTP routes.
type merchantDocumentQueryDocumentHandleDeps struct {
	// client is the gRPC client for the MerchantDocumentService.
	client pb.MerchantDocumentServiceClient

	// router is the Echo HTTP router for endpoint registration.
	router *echo.Echo

	// logger is the logging interface used throughout the shared.
	logger logger.LoggerInterface

	// mapper maps internal service responses to HTTP API response formats.
	mapper apimapper.MerchantDocumentQueryResponseMapper
}

func NewMerchantQueryDocumentHandler(params *merchantDocumentQueryDocumentHandleDeps) *merchantQueryDocumentHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_document_query_handler_requests_total",
			Help: "Total number of merchant document requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_document_query_handler_request_duration_seconds",
			Help:    "Duration of merchant document requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	merchantDocumentHandler := &merchantQueryDocumentHandleApi{
		merchantDocument: params.client,
		logger:           params.logger,
		mapper:           params.mapper,
		trace:            otel.Tracer("merchant-document-query-handler"),
		requestCounter:   requestCounter,
		requestDuration:  requestDuration,
	}

	routerMerchantDocument := params.router.Group("/api/merchant-document-query")

	routerMerchantDocument.GET("", merchantDocumentHandler.FindAll)
	routerMerchantDocument.GET("/:id", merchantDocumentHandler.FindById)
	routerMerchantDocument.GET("/active", merchantDocumentHandler.FindAllActive)
	routerMerchantDocument.GET("/trashed", merchantDocumentHandler.FindAllTrashed)

	return merchantDocumentHandler
}

// FindAll godoc
// @Summary Find all merchant documents
// @Tags Merchant Document Query
// @Security Bearer
// @Description Retrieve a list of all merchant documents
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantDocument "List of merchant documents"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant document data"
// @Router /api/merchant-document-query [get]
func (h *merchantQueryDocumentHandleApi) FindAll(c echo.Context) error {
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

	req := &pb.FindAllMerchantDocumentsRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.merchantDocument.FindAll(ctx, req)

	if err != nil {
		logError("failed find all", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedFindAllMerchantDocuments(c)
	}

	so := h.mapper.ToApiResponsePaginationMerchantDocument(res)

	logSuccess("success find all", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindById godoc
// @Summary Get merchant document by ID
// @Tags Merchant Document Query
// @Security Bearer
// @Description Get a merchant document by its ID
// @Accept json
// @Produce json
// @Param id path int true "Document ID"
// @Success 200 {object} response.ApiResponseMerchantDocument "Document details"
// @Failure 400 {object} response.ErrorResponse "Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to get document"
// @Router /api/merchant-document-query/{id} [get]
func (h *merchantQueryDocumentHandleApi) FindById(c echo.Context) error {
	const method = "FindById"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logError("failed find by id", err, zap.Error(err))

		return merchantdocument_errors.ErrApiInvalidMerchantDocumentID(c)
	}

	res, err := h.merchantDocument.FindById(ctx, &pb.FindMerchantDocumentByIdRequest{
		DocumentId: int32(id),
	})

	if err != nil {
		logError("failed find by id", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedFindByIdMerchantDocument(c)
	}

	so := h.mapper.ToApiResponseMerchantDocument(res)

	logSuccess("success find by id", zap.Error(err))

	return c.JSON(http.StatusOK, so)
}

// FindAllActive godoc
// @Summary Find all active merchant documents
// @Tags Merchant Document Query
// @Security Bearer
// @Description Retrieve a list of all active merchant documents
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantDocument "List of active merchant documents"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve active merchant documents"
// @Router /api/merchant-document-query/active [get]
func (h *merchantQueryDocumentHandleApi) FindAllActive(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAllActive"
	)

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	page := shared.ParseQueryInt(c, "page", defaultPage)
	pageSize := shared.ParseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	req := &pb.FindAllMerchantDocumentsRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.merchantDocument.FindAllActive(ctx, req)

	if err != nil {
		logError("failed find all active", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedFindAllActiveMerchantDocuments(c)
	}

	so := h.mapper.ToApiResponsePaginationMerchantDocument(res)

	logSuccess("success find all active", zap.Error(err))

	return c.JSON(http.StatusOK, so)
}

// FindAllTrashed godoc
// @Summary Find all trashed merchant documents
// @Tags Merchant Document Query
// @Security Bearer
// @Description Retrieve a list of all trashed merchant documents
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantDocumentDeleteAt "List of trashed merchant documents"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed merchant documents"
// @Router /api/merchant-document-query/trashed [get]
func (h *merchantQueryDocumentHandleApi) FindAllTrashed(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAllTrashed"
	)

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	page := shared.ParseQueryInt(c, "page", defaultPage)
	pageSize := shared.ParseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	req := &pb.FindAllMerchantDocumentsRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.merchantDocument.FindAllTrashed(ctx, req)

	if err != nil {
		logError("failed find all trashed", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedFindAllTrashedMerchantDocuments(c)
	}

	so := h.mapper.ToApiResponsePaginationMerchantDocumentDeleteAt(res)

	logSuccess("success find all trashed", zap.Error(err))

	return c.JSON(http.StatusOK, so)
}

func (s *merchantQueryDocumentHandleApi) startTracingAndLogging(
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

func (s *merchantQueryDocumentHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
