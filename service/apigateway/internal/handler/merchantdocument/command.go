package merchantdocumenthandler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	pb "github.com/MamangRust/monolith-payment-gateway-pb/merchantdocument"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	merchantdocument_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_document_errors/api"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/merchantdocument"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantCommandDocumentHandleApi struct {
	// merchantDocument is the gRPC client used to communicate with the MerchantDocumentService.
	merchantDocument pb.MerchantDocumentCommandServiceClient

	// logger is used for structured and contextual logging.
	logger logger.LoggerInterface

	// mapper is responsible for converting gRPC responses to API-compliant responses.
	mapper apimapper.MerchantDocumentCommandResponseMapper

	// trace provides distributed tracing capabilities using OpenTelemetry.
	trace trace.Tracer

	// requestCounter counts the number of requests received by the handler.
	requestCounter *prometheus.CounterVec

	// requestDuration records how long each handler request takes in seconds.
	requestDuration *prometheus.HistogramVec
}

// merchantDocumentHandleDeps defines the parameters required to initialize
// the merchant document handler and register its HTTP routes.
type merchantCommandDocumentHandleDeps struct {
	// client is the gRPC client for the MerchantDocumentService.
	client pb.MerchantDocumentCommandServiceClient

	// router is the Echo HTTP router for endpoint registration.
	router *echo.Echo

	// logger is the logging interface used throughout the handler.
	logger logger.LoggerInterface

	// mapper maps internal service responses to HTTP API response formats.
	mapper apimapper.MerchantDocumentCommandResponseMapper
}

func NewMerchantCommandDocumentHandler(params *merchantCommandDocumentHandleDeps) {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_document_command_handler_requests_total",
			Help: "Total number of merchant document requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_document_command_handler_request_duration_seconds",
			Help:    "Duration of merchant document requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	merchantDocumentHandler := &merchantCommandDocumentHandleApi{
		merchantDocument: params.client,
		logger:           params.logger,
		mapper:           params.mapper,
		trace:            otel.Tracer("merchant-document-command-handler"),
		requestCounter:   requestCounter,
		requestDuration:  requestDuration,
	}

	routerMerchantDocument := params.router.Group("/api/merchant-document-command")

	routerMerchantDocument.POST("/create", merchantDocumentHandler.Create)
	routerMerchantDocument.POST("/updates/:id", merchantDocumentHandler.Update)
	routerMerchantDocument.POST("/update-status/:id", merchantDocumentHandler.UpdateStatus)

	routerMerchantDocument.POST("/trashed/:id", merchantDocumentHandler.TrashedDocument)
	routerMerchantDocument.POST("/restore/:id", merchantDocumentHandler.RestoreDocument)
	routerMerchantDocument.DELETE("/permanent/:id", merchantDocumentHandler.Delete)

	routerMerchantDocument.POST("/restore/all", merchantDocumentHandler.RestoreAllDocuments)
	routerMerchantDocument.POST("/permanent/all", merchantDocumentHandler.DeleteAllDocumentsPermanent)
}

// Create godoc
// @Summary Create a new merchant document
// @Tags Merchant Document Command
// @Security Bearer
// @Description Create a new document for a merchant
// @Accept json
// @Produce json
// @Param body body requests.CreateMerchantDocumentRequest true "Create merchant document request"
// @Success 200 {object} response.ApiResponseMerchantDocument "Created document"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create document"
// @Router /api/merchant-document-command/create [post]
func (h *merchantCommandDocumentHandleApi) Create(c echo.Context) error {
	const method = "Create"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	var body requests.CreateMerchantDocumentRequest

	if err := c.Bind(&body); err != nil {
		logError("failed bind create", err, zap.Error(err))

		return merchantdocument_errors.ErrApiBindCreateMerchantDocument(c)
	}

	if err := body.Validate(); err != nil {
		logError("failed validate create", err, zap.Error(err))

		return merchantdocument_errors.ErrApiBindCreateMerchantDocument(c)
	}

	req := &pb.CreateMerchantDocumentRequest{
		MerchantId:   int32(body.MerchantID),
		DocumentType: body.DocumentType,
		DocumentUrl:  body.DocumentUrl,
	}

	res, err := h.merchantDocument.Create(ctx, req)

	if err != nil {
		logError("failed create", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedCreateMerchantDocument(c)
	}

	so := h.mapper.ToApiResponseMerchantDocument(res)

	logSuccess("success create", zap.Error(err))

	return c.JSON(http.StatusOK, so)
}

// Update godoc
// @Summary Update a merchant document
// @Tags Merchant Document Command
// @Security Bearer
// @Description Update a merchant document with the given ID
// @Accept json
// @Produce json
// @Param id path int true "Document ID"
// @Param body body requests.UpdateMerchantDocumentRequest true "Update merchant document request"
// @Success 200 {object} response.ApiResponseMerchantDocument "Updated document"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update document"
// @Router /api/merchant-document-command/update/{id} [post]
func (h *merchantCommandDocumentHandleApi) Update(c echo.Context) error {
	const method = "Update"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return merchantdocument_errors.ErrApiFailedUpdateMerchantDocument(c)
	}

	var body requests.UpdateMerchantDocumentRequest

	if err := c.Bind(&body); err != nil {
		logError("failed bind update", err, zap.Error(err))

		return merchantdocument_errors.ErrApiBindUpdateMerchantDocument(c)
	}

	if err := body.Validate(); err != nil {
		logError("failed validate update", err, zap.Error(err))

		return merchantdocument_errors.ErrApiValidateUpdateMerchantDocument(c)
	}

	req := &pb.UpdateMerchantDocumentRequest{
		DocumentId:   int32(id),
		MerchantId:   int32(body.MerchantID),
		DocumentType: body.DocumentType,
		DocumentUrl:  body.DocumentUrl,
		Status:       body.Status,
		Note:         body.Note,
	}

	res, err := h.merchantDocument.Update(ctx, req)

	if err != nil {
		logError("failed update", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedUpdateMerchantDocument(c)
	}

	so := h.mapper.ToApiResponseMerchantDocument(res)

	logSuccess("success update", zap.Error(err))

	return c.JSON(http.StatusOK, so)
}

// UpdateStatus godoc
// @Summary Update merchant document status
// @Tags Merchant Document Command
// @Security Bearer
// @Description Update the status of a merchant document
// @Accept json
// @Produce json
// @Param id path int true "Document ID"
// @Param body body requests.UpdateMerchantDocumentStatusRequest true "Update status request"
// @Success 200 {object} response.ApiResponseMerchantDocument "Updated document"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update document status"
// @Router /api/merchants-documents/update-status/{id} [post]
func (h *merchantCommandDocumentHandleApi) UpdateStatus(c echo.Context) error {
	const method = "UpdateStatus"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logError("failed parse id", err, zap.Error(err))

		return merchantdocument_errors.ErrApiInvalidMerchantDocumentID(c)
	}

	var body requests.UpdateMerchantDocumentStatusRequest

	if err := c.Bind(&body); err != nil {
		logError("failed bind update status", err, zap.Error(err))

		return merchantdocument_errors.ErrApiBindUpdateMerchantDocumentStatus(c)
	}

	if err := body.Validate(); err != nil {
		logError("failed validate update status", err, zap.Error(err))

		return merchantdocument_errors.ErrApiBindUpdateMerchantDocumentStatus(c)
	}

	req := &pb.UpdateMerchantDocumentStatusRequest{
		DocumentId: int32(id),
		MerchantId: int32(body.MerchantID),
		Status:     body.Status,
		Note:       body.Note,
	}

	res, err := h.merchantDocument.UpdateStatus(ctx, req)

	if err != nil {
		logError("failed update status", err, zap.Error(err))

		return merchantdocument_errors.ErrApiBindUpdateMerchantDocumentStatus(c)
	}

	so := h.mapper.ToApiResponseMerchantDocument(res)

	logSuccess("success update status", zap.Error(err))

	return c.JSON(http.StatusOK, so)
}

// TrashedDocument godoc
// @Summary Trashed a merchant document
// @Tags Merchant Document Command
// @Security Bearer
// @Description Trashed a merchant document by its ID
// @Accept json
// @Produce json
// @Param id path int true "Document ID"
// @Success 200 {object} response.ApiResponseMerchantDocument "Trashed document"
// @Failure 400 {object} response.ErrorResponse "Bad request or invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to trashed document"
// @Router /api/merchant-document-command/trashed/{id} [post]
func (h *merchantCommandDocumentHandleApi) TrashedDocument(c echo.Context) error {
	const method = "TrashedDocument"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		logError("failed parse id", err, zap.Error(err))

		return merchantdocument_errors.ErrApiInvalidMerchantDocumentID(c)
	}

	res, err := h.merchantDocument.Trashed(ctx, &pb.FindMerchantDocumentByIdRequest{
		DocumentId: int32(idInt),
	})

	if err != nil {
		logError("failed trashed", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedTrashMerchantDocument(c)
	}

	so := h.mapper.ToApiResponseMerchantDocumentDeleteAt(res)

	logSuccess("success trashed", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// RestoreDocument godoc
// @Summary Restore a merchant document
// @Tags Merchant Document Command
// @Security Bearer
// @Description Restore a merchant document by its ID
// @Accept json
// @Produce json
// @Param id path int true "Document ID"
// @Success 200 {object} response.ApiResponseMerchantDocument "Restored document"
// @Failure 400 {object} response.ErrorResponse "Bad request or invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore document"
// @Router /api/merchant-document-command/restore/{id} [post]
func (h *merchantCommandDocumentHandleApi) RestoreDocument(c echo.Context) error {
	const method = "RestoreDocument"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		logError("failed parse id", err, zap.Error(err))

		return merchantdocument_errors.ErrApiInvalidMerchantDocumentID(c)
	}

	res, err := h.merchantDocument.Restore(ctx, &pb.FindMerchantDocumentByIdRequest{
		DocumentId: int32(idInt),
	})

	if err != nil {
		logError("failed restore", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedRestoreMerchantDocument(c)
	}

	so := h.mapper.ToApiResponseMerchantDocument(res)

	logSuccess("Success restore merchant document", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// Delete godoc
// @Summary Delete a merchant document
// @Tags Merchant Document Command
// @Security Bearer
// @Description Delete a merchant document by its ID
// @Accept json
// @Produce json
// @Param id path int true "Document ID"
// @Success 200 {object} response.ApiResponseMerchantDocumentDelete "Deleted document"
// @Failure 400 {object} response.ErrorResponse "Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete document"
// @Router /api/merchant-document-command/permanent/{id} [delete]
func (h *merchantCommandDocumentHandleApi) Delete(c echo.Context) error {
	const method = "Delete"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logError("failed parse id", err, zap.Error(err))

		return merchantdocument_errors.ErrApiInvalidMerchantDocumentID(c)
	}

	res, err := h.merchantDocument.DeletePermanent(ctx, &pb.FindMerchantDocumentByIdRequest{
		DocumentId: int32(id),
	})

	if err != nil {
		logError("failed delete", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedDeleteMerchantDocumentPermanent(c)
	}

	so := h.mapper.ToApiResponseMerchantDocumentDelete(res)

	logSuccess("Successfully deleted merchant document", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// RestoreAllDocuments godoc
// @Summary Restore all merchant documents
// @Tags Merchant Document Command
// @Security Bearer
// @Description Restore all merchant documents that were previously deleted
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseMerchantDocumentAll "Successfully restored all documents"
// @Failure 500 {object} response.ErrorResponse "Failed to restore all documents"
// @Router /api/merchant-document-command/restore/all [post]
func (h *merchantCommandDocumentHandleApi) RestoreAllDocuments(c echo.Context) error {
	const method = "RestoreAllDocuments"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	res, err := h.merchantDocument.RestoreAll(ctx, &emptypb.Empty{})

	if err != nil {
		logError("failed restore all", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedRestoreAllMerchantDocuments(c)
	}

	response := h.mapper.ToApiResponseMerchantDocumentAll(res)

	logSuccess("Successfully restored all merchant documents")

	return c.JSON(http.StatusOK, response)
}

// DeleteAllDocumentsPermanent godoc
// @Summary Permanently delete all merchant documents
// @Tags Merchant Document Command
// @Security Bearer
// @Description Permanently delete all merchant documents from the database
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseMerchantDocumentAll "Successfully deleted all documents permanently"
// @Failure 500 {object} response.ErrorResponse "Failed to permanently delete all documents"
// @Router /api/merchant-document-command/permanent/all [post]
func (h *merchantCommandDocumentHandleApi) DeleteAllDocumentsPermanent(c echo.Context) error {
	const method = "DeleteAllDocumentsPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	res, err := h.merchantDocument.DeleteAllPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		logError("failed delete all", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedDeleteAllMerchantDocumentsPermanent(c)
	}

	response := h.mapper.ToApiResponseMerchantDocumentAll(res)

	logSuccess("Successfully deleted all merchant documents permanently")

	return c.JSON(http.StatusOK, response)
}

// startTracingAndLogging starts tracing for a method, logs that the method has started,
// and returns a span, a function to end the span, the initial status of the span, and
// a function to log a success message.
func (s *merchantCommandDocumentHandleApi) startTracingAndLogging(
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

// recordMetrics records a Prometheus metric for the given method and status.
// It increments a counter and records the duration since the provided start time.
func (s *merchantCommandDocumentHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
