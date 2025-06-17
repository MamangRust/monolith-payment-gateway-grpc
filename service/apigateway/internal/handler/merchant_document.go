package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	merchantdocument_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_document_errors"
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

type merchantDocumentHandleApi struct {
	merchantDocument pb.MerchantDocumentServiceClient
	logger           logger.LoggerInterface
	mapping          apimapper.MerchantDocumentResponseMapper
	trace            trace.Tracer
	requestCounter   *prometheus.CounterVec
	requestDuration  *prometheus.HistogramVec
}

func NewHandlerMerchantDocument(merchantDocument pb.MerchantDocumentServiceClient, router *echo.Echo, logger logger.LoggerInterface, ma apimapper.MerchantDocumentResponseMapper) *merchantDocumentHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_document_handler_requests_total",
			Help: "Total number of card requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_document_handler_request_duration_seconds",
			Help:    "Duration of card requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	merchantDocumentHandler := &merchantDocumentHandleApi{
		merchantDocument: merchantDocument,
		logger:           logger,
		mapping:          ma,
		trace:            otel.Tracer("merchant-document-handler"),
		requestCounter:   requestCounter,
		requestDuration:  requestDuration,
	}

	routerMerchantDocument := router.Group("/api/merchant-documents")

	routerMerchantDocument.GET("", merchantDocumentHandler.FindAll)
	routerMerchantDocument.GET("/:id", merchantDocumentHandler.FindById)
	routerMerchantDocument.GET("/active", merchantDocumentHandler.FindAllActive)
	routerMerchantDocument.GET("/trashed", merchantDocumentHandler.FindAllTrashed)

	routerMerchantDocument.POST("/create", merchantDocumentHandler.Create)
	routerMerchantDocument.POST("/updates/:id", merchantDocumentHandler.Update)
	routerMerchantDocument.POST("/update-status/:id", merchantDocumentHandler.UpdateStatus)

	routerMerchantDocument.POST("/trashed/:id", merchantDocumentHandler.TrashedDocument)
	routerMerchantDocument.POST("/restore/:id", merchantDocumentHandler.RestoreDocument)
	routerMerchantDocument.DELETE("/permanent/:id", merchantDocumentHandler.Delete)

	routerMerchantDocument.POST("/restore/all", merchantDocumentHandler.RestoreAllDocuments)
	routerMerchantDocument.POST("/permanent/all", merchantDocumentHandler.DeleteAllDocumentsPermanent)

	return merchantDocumentHandler
}

// FindAll godoc
// @Summary Find all merchant documents
// @Tags Merchant Document
// @Security Bearer
// @Description Retrieve a list of all merchant documents
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantDocument "List of merchant documents"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant document data"
// @Router /api/merchant-documents [get]
func (h *merchantDocumentHandleApi) FindAll(c echo.Context) error {
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

	req := &pb.FindAllMerchantDocumentsRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.merchantDocument.FindAll(ctx, req)
	if err != nil {
		status = "error"

		logError("failed find all", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedFindAllMerchantDocuments(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantDocument(res)

	logSuccess("success find all", zap.Any("response", so))

	return c.JSON(http.StatusOK, so)
}

// FindById godoc
// @Summary Get merchant document by ID
// @Tags Merchant Document
// @Security Bearer
// @Description Get a merchant document by its ID
// @Accept json
// @Produce json
// @Param id path int true "Document ID"
// @Success 200 {object} response.ApiResponseMerchantDocument "Document details"
// @Failure 400 {object} response.ErrorResponse "Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to get document"
// @Router /api/merchant-documents/{id} [get]
func (h *merchantDocumentHandleApi) FindById(c echo.Context) error {
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

		logError("failed find by id", err, zap.Error(err))

		return merchantdocument_errors.ErrApiInvalidMerchantDocumentID(c)
	}

	res, err := h.merchantDocument.FindById(ctx, &pb.FindMerchantDocumentByIdRequest{
		DocumentId: int32(id),
	})

	if err != nil {
		status = "error"

		logError("failed find by id", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedFindByIdMerchantDocument(c)
	}

	so := h.mapping.ToApiResponseMerchantDocument(res)

	logSuccess("success find by id", zap.Error(err))

	return c.JSON(http.StatusOK, so)
}

// FindAllActive godoc
// @Summary Find all active merchant documents
// @Tags Merchant Document
// @Security Bearer
// @Description Retrieve a list of all active merchant documents
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantDocument "List of active merchant documents"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve active merchant documents"
// @Router /api/merchant-documents/active [get]
func (h *merchantDocumentHandleApi) FindAllActive(c echo.Context) error {
	const method = "FindAllActive"

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

	req := &pb.FindAllMerchantDocumentsRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.merchantDocument.FindAllActive(ctx, req)
	if err != nil {
		status = "error"

		logError("failed find all active", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedFindAllActiveMerchantDocuments(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantDocument(res)

	logSuccess("success find all active", zap.Error(err))

	return c.JSON(http.StatusOK, so)
}

// FindAllTrashed godoc
// @Summary Find all trashed merchant documents
// @Tags Merchant Document
// @Security Bearer
// @Description Retrieve a list of all trashed merchant documents
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantDocumentDeleteAt "List of trashed merchant documents"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed merchant documents"
// @Router /api/merchant-documents/trashed [get]
func (h *merchantDocumentHandleApi) FindAllTrashed(c echo.Context) error {
	const method = "FindAllTrashed"

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

	req := &pb.FindAllMerchantDocumentsRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.merchantDocument.FindAllTrashed(ctx, req)
	if err != nil {
		status = "error"

		logError("failed find all trashed", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedFindAllTrashedMerchantDocuments(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantDocumentDeleteAt(res)

	logSuccess("success find all trashed", zap.Error(err))

	return c.JSON(http.StatusOK, so)
}

// Create godoc
// @Summary Create a new merchant document
// @Tags Merchant Document
// @Security Bearer
// @Description Create a new document for a merchant
// @Accept json
// @Produce json
// @Param body body requests.CreateMerchantDocumentRequest true "Create merchant document request"
// @Success 200 {object} response.ApiResponseMerchantDocument "Created document"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create document"
// @Router /api/merchant-documents/create [post]
func (h *merchantDocumentHandleApi) Create(c echo.Context) error {
	const method = "Create"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	var body requests.CreateMerchantDocumentRequest

	if err := c.Bind(&body); err != nil {
		status = "error"

		logError("failed bind create", err, zap.Error(err))

		return merchantdocument_errors.ErrApiBindCreateMerchantDocument(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"

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
		status = "error"

		logError("failed create", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedCreateMerchantDocument(c)
	}

	so := h.mapping.ToApiResponseMerchantDocument(res)

	logSuccess("success create", zap.Error(err))

	return c.JSON(http.StatusOK, so)
}

// Update godoc
// @Summary Update a merchant document
// @Tags Merchant Document
// @Security Bearer
// @Description Update a merchant document with the given ID
// @Accept json
// @Produce json
// @Param id path int true "Document ID"
// @Param body body requests.UpdateMerchantDocumentRequest true "Update merchant document request"
// @Success 200 {object} response.ApiResponseMerchantDocument "Updated document"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update document"
// @Router /api/merchant-documents/update/{id} [post]
func (h *merchantDocumentHandleApi) Update(c echo.Context) error {
	const method = "Update"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return merchantdocument_errors.ErrApiFailedUpdateMerchantDocument(c)
	}

	var body requests.UpdateMerchantDocumentRequest

	if err := c.Bind(&body); err != nil {
		status = "error"

		logError("failed bind update", err, zap.Error(err))

		return merchantdocument_errors.ErrApiBindUpdateMerchantDocument(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"

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
		status = "error"

		logError("failed update", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedUpdateMerchantDocument(c)
	}

	so := h.mapping.ToApiResponseMerchantDocument(res)

	logSuccess("success update", zap.Error(err))

	return c.JSON(http.StatusOK, so)
}

// UpdateStatus godoc
// @Summary Update merchant document status
// @Tags Merchant Document
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
func (h *merchantDocumentHandleApi) UpdateStatus(c echo.Context) error {
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

		logError("failed parse id", err, zap.Error(err))

		return merchantdocument_errors.ErrApiInvalidMerchantDocumentID(c)
	}

	var body requests.UpdateMerchantDocumentStatusRequest

	if err := c.Bind(&body); err != nil {
		status = "error"

		logError("failed bind update status", err, zap.Error(err))

		return merchantdocument_errors.ErrApiBindUpdateMerchantDocumentStatus(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"

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
		status = "error"

		logError("failed update status", err, zap.Error(err))

		return merchantdocument_errors.ErrApiBindUpdateMerchantDocumentStatus(c)
	}

	so := h.mapping.ToApiResponseMerchantDocument(res)

	logSuccess("success update status", zap.Error(err))

	return c.JSON(http.StatusOK, so)
}

// TrashedDocument godoc
// @Summary Trashed a merchant document
// @Tags Merchant Document
// @Security Bearer
// @Description Trashed a merchant document by its ID
// @Accept json
// @Produce json
// @Param id path int true "Document ID"
// @Success 200 {object} response.ApiResponseMerchantDocument "Trashed document"
// @Failure 400 {object} response.ErrorResponse "Bad request or invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to trashed document"
// @Router /api/merchant-documents/trashed/{id} [post]
func (h *merchantDocumentHandleApi) TrashedDocument(c echo.Context) error {
	const method = "TrashedDocument"

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

		logError("failed parse id", err, zap.Error(err))

		return merchantdocument_errors.ErrApiInvalidMerchantDocumentID(c)
	}

	res, err := h.merchantDocument.Trashed(ctx, &pb.TrashedMerchantDocumentRequest{
		DocumentId: int32(idInt),
	})

	if err != nil {
		status = "error"

		logError("failed trashed", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedTrashMerchantDocument(c)
	}

	so := h.mapping.ToApiResponseMerchantDocument(res)

	logSuccess("success trashed", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// RestoreDocument godoc
// @Summary Restore a merchant document
// @Tags Merchant Document
// @Security Bearer
// @Description Restore a merchant document by its ID
// @Accept json
// @Produce json
// @Param id path int true "Document ID"
// @Success 200 {object} response.ApiResponseMerchantDocument "Restored document"
// @Failure 400 {object} response.ErrorResponse "Bad request or invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore document"
// @Router /api/merchant-documents/restore/{id} [post]
func (h *merchantDocumentHandleApi) RestoreDocument(c echo.Context) error {
	const method = "RestoreDocument"

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

		logError("failed parse id", err, zap.Error(err))

		return merchantdocument_errors.ErrApiInvalidMerchantDocumentID(c)
	}

	res, err := h.merchantDocument.Restore(ctx, &pb.RestoreMerchantDocumentRequest{
		DocumentId: int32(idInt),
	})

	if err != nil {
		status = "error"

		logError("failed restore", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedRestoreMerchantDocument(c)
	}

	so := h.mapping.ToApiResponseMerchantDocument(res)

	logSuccess("Success restore merchant document", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// Delete godoc
// @Summary Delete a merchant document
// @Tags Merchant Document
// @Security Bearer
// @Description Delete a merchant document by its ID
// @Accept json
// @Produce json
// @Param id path int true "Document ID"
// @Success 200 {object} response.ApiResponseMerchantDocumentDelete "Deleted document"
// @Failure 400 {object} response.ErrorResponse "Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete document"
// @Router /api/merchant-documents/permanent/{id} [delete]
func (h *merchantDocumentHandleApi) Delete(c echo.Context) error {
	const method = "Delete"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("failed parse id", err, zap.Error(err))

		return merchantdocument_errors.ErrApiInvalidMerchantDocumentID(c)
	}

	res, err := h.merchantDocument.DeletePermanent(ctx, &pb.DeleteMerchantDocumentPermanentRequest{
		DocumentId: int32(id),
	})

	if err != nil {
		status = "error"

		logError("failed delete", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedDeleteMerchantDocumentPermanent(c)
	}

	so := h.mapping.ToApiResponseMerchantDocumentDeleteAt(res)

	logSuccess("Successfully deleted merchant document", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// RestoreAllDocuments godoc
// @Summary Restore all merchant documents
// @Tags Merchant Document
// @Security Bearer
// @Description Restore all merchant documents that were previously deleted
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseMerchantDocumentAll "Successfully restored all documents"
// @Failure 500 {object} response.ErrorResponse "Failed to restore all documents"
// @Router /api/merchant-documents/restore/all [post]
func (h *merchantDocumentHandleApi) RestoreAllDocuments(c echo.Context) error {
	const method = "RestoreAllDocuments"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	res, err := h.merchantDocument.RestoreAll(ctx, &emptypb.Empty{})
	if err != nil {
		status = "error"

		logError("failed restore all", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedRestoreAllMerchantDocuments(c)
	}

	response := h.mapping.ToApiResponseMerchantDocumentAll(res)

	logSuccess("Successfully restored all merchant documents")

	return c.JSON(http.StatusOK, response)
}

// DeleteAllDocumentsPermanent godoc
// @Summary Permanently delete all merchant documents
// @Tags Merchant Document
// @Security Bearer
// @Description Permanently delete all merchant documents from the database
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseMerchantDocumentAll "Successfully deleted all documents permanently"
// @Failure 500 {object} response.ErrorResponse "Failed to permanently delete all documents"
// @Router /api/merchant-documents/permanent/all [post]
func (h *merchantDocumentHandleApi) DeleteAllDocumentsPermanent(c echo.Context) error {
	const method = "DeleteAllDocumentsPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	res, err := h.merchantDocument.DeleteAllPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"

		logError("failed delete all", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedDeleteAllMerchantDocumentsPermanent(c)
	}

	response := h.mapping.ToApiResponseMerchantDocumentAll(res)

	logSuccess("Successfully deleted all merchant documents permanently")

	return c.JSON(http.StatusOK, response)
}

func (s *merchantDocumentHandleApi) startTracingAndLogging(
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

func (s *merchantDocumentHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
