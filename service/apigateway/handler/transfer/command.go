package transferhandler

import (
	"fmt"
	"net/http"
	"strconv"

	transfer_cache "github.com/MamangRust/monolith-payment-gateway-apigateway/redis/api/transfer"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/transfer"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/transfer"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type transferCommandHandleApi struct {
	client pb.TransferCommandServiceClient

	logger logger.LoggerInterface

	mapper apimapper.TransferCommandResponseMapper

	cache transfer_cache.TransferMencache

	apiHandler errors.ApiHandler
}

type transferCommandHandleDeps struct {
	client pb.TransferCommandServiceClient

	router *echo.Echo

	logger logger.LoggerInterface

	mapper apimapper.TransferCommandResponseMapper

	cache transfer_cache.TransferMencache

	apiHandler errors.ApiHandler
}

func NewTransferCommandHandleApi(params *transferCommandHandleDeps) *transferCommandHandleApi {

	transferCommandHandleApi := &transferCommandHandleApi{
		client:     params.client,
		logger:     params.logger,
		mapper:     params.mapper,
		cache:      params.cache,
		apiHandler: params.apiHandler,
	}

	routerTransfer := params.router.Group("/api/transfer-command")

	routerTransfer.POST("/create", params.apiHandler.Handle("create-transfer", transferCommandHandleApi.CreateTransfer))
	routerTransfer.POST("/update/:id", params.apiHandler.Handle("update-transfer", transferCommandHandleApi.UpdateTransfer))
	routerTransfer.POST("/trashed/:id", params.apiHandler.Handle("trash-transfer", transferCommandHandleApi.TrashTransfer))
	routerTransfer.POST("/restore/:id", params.apiHandler.Handle("restore-transfer", transferCommandHandleApi.RestoreTransfer))
	routerTransfer.DELETE("/permanent/:id", params.apiHandler.Handle("delete-transfer-permanent", transferCommandHandleApi.DeleteTransferPermanent))

	routerTransfer.POST("/restore/all", params.apiHandler.Handle("restore-all-transfers", transferCommandHandleApi.RestoreAllTransfer))
	routerTransfer.POST("/permanent/all", params.apiHandler.Handle("delete-all-transfers-permanent", transferCommandHandleApi.DeleteAllTransferPermanent))

	return transferCommandHandleApi
}

// @Summary Create a transfer
// @Tags Transfer Command
// @Security Bearer
// @Description Create a new transfer record
// @Accept json
// @Produce json
// @Param body body requests.CreateTransferRequest true "Transfer request"
// @Success 200 {object} response.ApiResponseTransfer "Transfer data"
// @Failure 400 {object} response.ErrorResponse "Validation Error"
// @Failure 500 {object} response.ErrorResponse "Failed to create transfer"
// @Router /api/transfer-command/create [post]
func (h *transferCommandHandleApi) CreateTransfer(c echo.Context) error {
	var body requests.CreateTransferRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Invalid request body: ", zap.Error(err))

		return err
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation Error: ", zap.Error(err))

		return err
	}

	ctx := c.Request().Context()

	res, err := h.client.CreateTransfer(ctx, &pb.CreateTransferRequest{
		TransferFrom:   body.TransferFrom,
		TransferTo:     body.TransferTo,
		TransferAmount: int32(body.TransferAmount),
	})

	if err != nil {
		h.logger.Debug("Failed to create transfer: ", zap.Error(err))

		return err
	}

	so := h.mapper.ToApiResponseTransfer(res)

	h.cache.SetCachedTransferCache(ctx, so)

	return c.JSON(http.StatusOK, so)
}

// @Summary Update a transfer
// @Tags Transfer Command
// @Security Bearer
// @Description Update an existing transfer record
// @Accept json
// @Produce json
// @Param id path int true "Transfer ID"
// @Param body body requests.UpdateTransferRequest true "Transfer request"
// @Success 200 {object} response.ApiResponseTransfer "Transfer data"
// @Failure 400 {object} response.ErrorResponse "Validation Error"
// @Failure 500 {object} response.ErrorResponse "Failed to update transfer"
// @Router /api/transfer-command/update/{id} [post]
func (h *transferCommandHandleApi) UpdateTransfer(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		h.logger.Debug("Bad Request: Invalid ID", zap.Error(err))
		return err
	}

	var body requests.UpdateTransferRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Invalid request body: ", zap.Error(err))

		return err
	}

	body.TransferID = &idInt

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation Error: ", zap.Error(err))

		return err
	}

	ctx := c.Request().Context()

	res, err := h.client.UpdateTransfer(ctx, &pb.UpdateTransferRequest{
		TransferId:     int32(idInt),
		TransferFrom:   body.TransferFrom,
		TransferTo:     body.TransferTo,
		TransferAmount: int32(body.TransferAmount),
	})

	if err != nil {
		h.logger.Debug("Failed to update transfer: ", zap.Error(err))

		return err
	}

	so := h.mapper.ToApiResponseTransfer(res)

	h.cache.DeleteTransferCache(ctx, idInt)
	h.cache.SetCachedTransferCache(ctx, so)

	return c.JSON(http.StatusOK, so)
}

// @Summary Soft delete a transfer
// @Tags Transfer Command
// @Security Bearer
// @Description Soft delete a transfer record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Transfer ID"
// @Success 200 {object} response.ApiResponseTransfer "Successfully trashed transfer record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to trashed transfer"
// @Router /api/transfer-command/trash/{id} [post]
func (h *transferCommandHandleApi) TrashTransfer(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		h.logger.Debug("Bad Request: Invalid ID", zap.Error(err))

		return err
	}

	ctx := c.Request().Context()

	res, err := h.client.TrashedTransfer(ctx, &pb.FindByIdTransferRequest{
		TransferId: int32(idInt),
	})

	if err != nil {
		h.logger.Debug("Failed to trash transfer: ", zap.Error(err))

		return err
	}

	so := h.mapper.ToApiResponseTransferDeleteAt(res)

	h.cache.DeleteTransferCache(ctx, idInt)

	return c.JSON(http.StatusOK, so)
}

// @Summary Restore a trashed transfer
// @Tags Transfer Command
// @Security Bearer
// @Description Restore a trashed transfer record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Transfer ID"
// @Success 200 {object} response.ApiResponseTransfer "Successfully restored transfer record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore transfer:"
// @Router /api/transfer-command/restore/{id} [post]
func (h *transferCommandHandleApi) RestoreTransfer(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		h.logger.Debug("Bad Request: Invalid ID", zap.Error(err))

		return err
	}

	ctx := c.Request().Context()

	res, err := h.client.RestoreTransfer(ctx, &pb.FindByIdTransferRequest{
		TransferId: int32(idInt),
	})

	if err != nil {
		h.logger.Debug("Failed to restore transfer: ", zap.Error(err))

		return err
	}

	so := h.mapper.ToApiResponseTransferDeleteAt(res)

	h.cache.DeleteTransferCache(ctx, idInt)

	return c.JSON(http.StatusOK, so)
}

// @Summary Permanently delete a transfer
// @Tags Transfer Command
// @Security Bearer
// @Description Permanently delete a transfer record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Transfer ID"
// @Success 200 {object} response.ApiResponseTransferDelete "Successfully deleted transfer record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete transfer:"
// @Router /api/transfer-command/permanent/{id} [delete]
func (h *transferCommandHandleApi) DeleteTransferPermanent(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		h.logger.Debug("Bad Request: Invalid ID", zap.Error(err))

		return err
	}

	ctx := c.Request().Context()

	res, err := h.client.DeleteTransferPermanent(ctx, &pb.FindByIdTransferRequest{
		TransferId: int32(idInt),
	})

	if err != nil {
		return err
	}

	so := h.mapper.ToApiResponseTransferDelete(res)

	h.cache.DeleteTransferCache(ctx, idInt)

	return c.JSON(http.StatusOK, so)
}

// @Summary Restore a trashed transfer
// @Tags Transfer Command
// @Security Bearer
// @Description Restore a trashed transfer all
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseTransferAll "Successfully restored transfer record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore transfer:"
// @Router /api/transfer-command/restore/all [post]
func (h *transferCommandHandleApi) RestoreAllTransfer(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllTransfer(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Failed to restore all transfer", zap.Error(err))
		return err
	}

	h.logger.Debug("Successfully restored all transfer")

	so := h.mapper.ToApiResponseTransferAll(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Permanently delete a transfer
// @Tags Transfer Command
// @Security Bearer
// @Description Permanently delete a transfer record all.
// @Accept json
// @Produce json
// @Param id path int true "Transfer ID"
// @Success 200 {object} response.ApiResponseTransferAll "Successfully deleted transfer all"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete transfer:"
// @Router /api/transfer-command/permanent/all [post]
func (h *transferCommandHandleApi) DeleteAllTransferPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.DeleteAllTransferPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Failed to permanently delete all transfer", zap.Error(err))

		return err
	}

	h.logger.Debug("Successfully deleted all transfer permanently")

	so := h.mapper.ToApiResponseTransferAll(res)

	return c.JSON(http.StatusOK, so)
}

func (h *transferCommandHandleApi) handleGrpcError(err error, operation string) *errors.AppError {
	st, ok := status.FromError(err)
	if !ok {
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}

	switch st.Code() {
	case codes.NotFound:
		return errors.NewNotFoundError("Transfer").WithInternal(err)

	case codes.AlreadyExists:
		return errors.NewConflictError("Transfer already exists").WithInternal(err)

	case codes.InvalidArgument:
		return errors.NewBadRequestError(st.Message()).WithInternal(err)

	case codes.PermissionDenied:
		return errors.ErrForbidden.WithInternal(err)

	case codes.Unauthenticated:
		return errors.ErrUnauthorized.WithInternal(err)

	case codes.ResourceExhausted:
		return errors.ErrTooManyRequests.WithInternal(err)

	case codes.Unavailable:
		return errors.NewServiceUnavailableError("Transfer service").WithInternal(err)

	case codes.DeadlineExceeded:
		return errors.ErrTimeout.WithInternal(err)

	default:
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}
}

func (h *transferCommandHandleApi) parseValidationErrors(err error) []errors.ValidationError {
	var validationErrs []errors.ValidationError

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			validationErrs = append(validationErrs, errors.ValidationError{
				Field:   fe.Field(),
				Message: h.getValidationMessage(fe),
			})
		}
		return validationErrs
	}

	return []errors.ValidationError{
		{
			Field:   "general",
			Message: err.Error(),
		},
	}
}

func (h *transferCommandHandleApi) getValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Must be at least %s", fe.Param())
	case "max":
		return fmt.Sprintf("Must be at most %s", fe.Param())
	case "gte":
		return fmt.Sprintf("Must be greater than or equal to %s", fe.Param())
	case "lte":
		return fmt.Sprintf("Must be less than or equal to %s", fe.Param())
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", fe.Param())
	default:
		return fmt.Sprintf("Validation failed on '%s' tag", fe.Tag())
	}
}
