package handler

import (
	"net/http"
	"strconv"

	"github.com/MamangRust/payment-gateway-monolith-grpc/pkg/logger"
	"github.com/MamangRust/payment-gateway-monolith-grpc/shared/domain/requests"
	merchantdocument_errors "github.com/MamangRust/payment-gateway-monolith-grpc/shared/errors/merchant_document_errors"
	apimapper "github.com/MamangRust/payment-gateway-monolith-grpc/shared/mapper/response/api"
	"github.com/MamangRust/payment-gateway-monolith-grpc/shared/pb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantDocumentHandleApi struct {
	merchantDocument pb.MerchantDocumentServiceClient
	logger           logger.LoggerInterface
	mapping          apimapper.MerchantDocumentResponseMapper
}

func NewHandlerMerchantDocument(merchantDocument pb.MerchantDocumentServiceClient, router *echo.Echo, logger logger.LoggerInterface, ma apimapper.MerchantDocumentResponseMapper) *merchantDocumentHandleApi {
	merchantDocumentHandler := &merchantDocumentHandleApi{
		merchantDocument: merchantDocument,
		logger:           logger,
		mapping:          ma,
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
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")

	ctx := c.Request().Context()

	req := &pb.FindAllMerchantDocumentsRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.merchantDocument.FindAll(ctx, req)
	if err != nil {
		h.logger.Debug("Failed to retrieve merchant document data", zap.Error(err))
		return merchantdocument_errors.ErrApiFailedFindAllMerchantDocuments(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantDocument(res)

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
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return merchantdocument_errors.ErrApiInvalidMerchantDocumentID(c)
	}

	ctx := c.Request().Context()

	res, err := h.merchantDocument.FindById(ctx, &pb.FindMerchantDocumentByIdRequest{
		DocumentId: int32(id),
	})

	if err != nil {
		h.logger.Debug("Failed to get merchant document", zap.Error(err))
		return merchantdocument_errors.ErrApiFailedFindByIdMerchantDocument(c)
	}

	so := h.mapping.ToApiResponseMerchantDocument(res)

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
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")

	ctx := c.Request().Context()

	req := &pb.FindAllMerchantDocumentsRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.merchantDocument.FindAllActive(ctx, req)
	if err != nil {
		h.logger.Debug("Failed to retrieve active merchant document data", zap.Error(err))
		return merchantdocument_errors.ErrApiFailedFindAllActiveMerchantDocuments(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantDocument(res)

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
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")

	ctx := c.Request().Context()

	req := &pb.FindAllMerchantDocumentsRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.merchantDocument.FindAllTrashed(ctx, req)
	if err != nil {
		h.logger.Debug("Failed to retrieve trashed merchant document data", zap.Error(err))
		return merchantdocument_errors.ErrApiFailedFindAllTrashedMerchantDocuments(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantDocumentDeleteAt(res)

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
	var body requests.CreateMerchantDocumentRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Bad Request", zap.Error(err))
		return merchantdocument_errors.ErrApiBindCreateMerchantDocument(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation Error", zap.Error(err))
		return merchantdocument_errors.ErrApiBindCreateMerchantDocument(c)
	}

	ctx := c.Request().Context()

	req := &pb.CreateMerchantDocumentRequest{
		MerchantId:   int32(body.MerchantID),
		DocumentType: body.DocumentType,
		DocumentUrl:  body.DocumentUrl,
	}

	res, err := h.merchantDocument.Create(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to create merchant document", zap.Error(err))
		return merchantdocument_errors.ErrApiFailedCreateMerchantDocument(c)
	}

	so := h.mapping.ToApiResponseMerchantDocument(res)

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
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return merchantdocument_errors.ErrApiFailedUpdateMerchantDocument(c)
	}

	var body requests.UpdateMerchantDocumentRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Bad Request", zap.Error(err))
		return merchantdocument_errors.ErrApiBindUpdateMerchantDocument(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation Error", zap.Error(err))
		return merchantdocument_errors.ErrApiValidateUpdateMerchantDocument(c)
	}

	ctx := c.Request().Context()
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
		h.logger.Debug("Failed to update merchant document", zap.Error(err))
		return merchantdocument_errors.ErrApiFailedUpdateMerchantDocument(c)
	}

	so := h.mapping.ToApiResponseMerchantDocument(res)

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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return merchantdocument_errors.ErrApiInvalidMerchantDocumentID(c)
	}

	var body requests.UpdateMerchantDocumentStatusRequest
	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Bad Request", zap.Error(err))
		return merchantdocument_errors.ErrApiBindUpdateMerchantDocumentStatus(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation Error", zap.Error(err))
		return merchantdocument_errors.ErrApiBindUpdateMerchantDocumentStatus(c)
	}

	ctx := c.Request().Context()
	req := &pb.UpdateMerchantDocumentStatusRequest{
		DocumentId: int32(id),
		MerchantId: int32(body.MerchantID),
		Status:     body.Status,
		Note:       body.Note,
	}

	res, err := h.merchantDocument.UpdateStatus(ctx, req)
	if err != nil {
		h.logger.Debug("Failed to update merchant document status", zap.Error(err))
		return merchantdocument_errors.ErrApiBindUpdateMerchantDocumentStatus(c)
	}

	so := h.mapping.ToApiResponseMerchantDocument(res)
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
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		h.logger.Debug("Bad Request", zap.Error(err))
		return merchantdocument_errors.ErrApiInvalidMerchantDocumentID(c)
	}

	ctx := c.Request().Context()

	res, err := h.merchantDocument.Trashed(ctx, &pb.TrashedMerchantDocumentRequest{
		DocumentId: int32(idInt),
	})

	if err != nil {
		h.logger.Debug("Failed to trashed merchant document", zap.Error(err))
		return merchantdocument_errors.ErrApiFailedTrashMerchantDocument(c)
	}

	so := h.mapping.ToApiResponseMerchantDocument(res)
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
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		h.logger.Debug("Bad Request", zap.Error(err))
		return merchantdocument_errors.ErrApiInvalidMerchantDocumentID(c)
	}

	ctx := c.Request().Context()

	res, err := h.merchantDocument.Restore(ctx, &pb.RestoreMerchantDocumentRequest{
		DocumentId: int32(idInt),
	})

	if err != nil {
		h.logger.Debug("Failed to restore merchant document", zap.Error(err))
		return merchantdocument_errors.ErrApiFailedRestoreMerchantDocument(c)
	}

	so := h.mapping.ToApiResponseMerchantDocument(res)
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
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return merchantdocument_errors.ErrApiInvalidMerchantDocumentID(c)
	}

	ctx := c.Request().Context()

	res, err := h.merchantDocument.DeletePermanent(ctx, &pb.DeleteMerchantDocumentPermanentRequest{
		DocumentId: int32(id),
	})

	if err != nil {
		h.logger.Debug("Failed to delete merchant document", zap.Error(err))
		return merchantdocument_errors.ErrApiFailedDeleteMerchantDocumentPermanent(c)
	}

	so := h.mapping.ToApiResponseMerchantDocumentDeleteAt(res)

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
	ctx := c.Request().Context()

	res, err := h.merchantDocument.RestoreAll(ctx, &emptypb.Empty{})
	if err != nil {
		h.logger.Error("Failed to restore all merchant documents",
			zap.Error(err),
			zap.String("method", "POST"),
		)
		return merchantdocument_errors.ErrApiFailedRestoreAllMerchantDocuments(c)
	}

	response := h.mapping.ToApiResponseMerchantDocumentAll(res)
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
	ctx := c.Request().Context()

	res, err := h.merchantDocument.DeleteAllPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Failed to permanently delete all merchant documents",
			zap.Error(err),
			zap.String("method", "POST"),
		)
		return merchantdocument_errors.ErrApiFailedDeleteAllMerchantDocumentsPermanent(c)
	}

	response := h.mapping.ToApiResponseMerchantDocumentAll(res)
	return c.JSON(http.StatusOK, response)
}
