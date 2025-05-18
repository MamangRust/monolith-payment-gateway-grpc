package handler

import (
	"net/http"
	"strconv"

	"github.com/MamangRust/payment-gateway-monolith-grpc/pkg/logger"
	"github.com/MamangRust/payment-gateway-monolith-grpc/shared/domain/requests"
	"github.com/MamangRust/payment-gateway-monolith-grpc/shared/domain/response"
	"github.com/MamangRust/payment-gateway-monolith-grpc/shared/errors/topup_errors"
	apimapper "github.com/MamangRust/payment-gateway-monolith-grpc/shared/mapper/response/api"
	"github.com/MamangRust/payment-gateway-monolith-grpc/shared/pb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type topupHandleApi struct {
	client  pb.TopupServiceClient
	logger  logger.LoggerInterface
	mapping apimapper.TopupResponseMapper
}

func NewHandlerTopup(client pb.TopupServiceClient, router *echo.Echo, logger logger.LoggerInterface, mapping apimapper.TopupResponseMapper) *topupHandleApi {
	topupHandler := &topupHandleApi{
		client:  client,
		logger:  logger,
		mapping: mapping,
	}
	routerTopup := router.Group("/api/topups")

	routerTopup.GET("", topupHandler.FindAll)
	routerTopup.GET("/card-number/:card_number", topupHandler.FindAllByCardNumber)
	routerTopup.GET("/:id", topupHandler.FindById)

	routerTopup.GET("/monthly-success", topupHandler.FindMonthlyTopupStatusSuccess)
	routerTopup.GET("/yearly-success", topupHandler.FindYearlyTopupStatusSuccess)
	routerTopup.GET("/monthly-failed", topupHandler.FindMonthlyTopupStatusFailed)
	routerTopup.GET("/yearly-failed", topupHandler.FindYearlyTopupStatusFailed)

	routerTopup.GET("/monthly-success-by-card", topupHandler.FindMonthlyTopupStatusSuccessByCardNumber)
	routerTopup.GET("/yearly-success-by-card", topupHandler.FindYearlyTopupStatusSuccessByCardNumber)
	routerTopup.GET("/monthly-failed-by-card", topupHandler.FindMonthlyTopupStatusFailedByCardNumber)
	routerTopup.GET("/yearly-failed-by-card", topupHandler.FindYearlyTopupStatusFailedByCardNumber)

	routerTopup.GET("/monthly-methods", topupHandler.FindMonthlyTopupMethods)
	routerTopup.GET("/yearly-methods", topupHandler.FindYearlyTopupMethods)
	routerTopup.GET("/monthly-amounts", topupHandler.FindMonthlyTopupAmounts)
	routerTopup.GET("/yearly-amounts", topupHandler.FindYearlyTopupAmounts)

	routerTopup.GET("/monthly-methods-by-card", topupHandler.FindMonthlyTopupMethodsByCardNumber)
	routerTopup.GET("/yearly-methods-by-card", topupHandler.FindYearlyTopupMethodsByCardNumber)
	routerTopup.GET("/monthly-amounts-by-card", topupHandler.FindMonthlyTopupAmountsByCardNumber)
	routerTopup.GET("/yearly-amounts-by-card", topupHandler.FindYearlyTopupAmountsByCardNumber)

	routerTopup.GET("/active", topupHandler.FindByActive)
	routerTopup.GET("/trashed", topupHandler.FindByTrashed)

	routerTopup.POST("/create", topupHandler.Create)
	routerTopup.POST("/update/:id", topupHandler.Update)
	routerTopup.POST("/trashed/:id", topupHandler.TrashTopup)
	routerTopup.POST("/restore/:id", topupHandler.RestoreTopup)
	routerTopup.DELETE("/permanent/:id", topupHandler.DeleteTopupPermanent)

	routerTopup.POST("/trashed/all", topupHandler.DeleteAllTopupPermanent)
	routerTopup.POST("/restore/all", topupHandler.RestoreAllTopup)

	return topupHandler

}

// @Tags Topup
// @Security Bearer
// @Description Retrieve a list of all topup data with pagination and search
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTopup "List of topup data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve topup data"
// @Router /api/topups [get]
func (h topupHandleApi) FindAll(c echo.Context) error {
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

	req := &pb.FindAllTopupRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAllTopup(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to retrieve topup data", zap.Error(err))

		return topup_errors.ErrApiFailedFindAllTopups(c)
	}

	so := h.mapping.ToApiResponsePaginationTopup(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Find all topup by card number
// @Tags Transaction
// @Security Bearer
// @Description Retrieve a list of transactions for a specific card number
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTopup "List of topups"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve topups data"
// @Router /api/topups/card-number/{card_number} [get]
func (h *topupHandleApi) FindAllByCardNumber(c echo.Context) error {
	cardNumber := c.Param("card_number")
	if cardNumber == "" {
		return topup_errors.ErrApiInvalidCardNumber(c)
	}

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

	req := &pb.FindAllTopupByCardNumberRequest{
		CardNumber: cardNumber,
		Page:       int32(page),
		PageSize:   int32(pageSize),
		Search:     search,
	}

	res, err := h.client.FindAllTopupByCardNumber(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to retrieve transaction data", zap.Error(err))

		return topup_errors.ErrApiFailedFindAllByCardNumberTopup(c)
	}

	so := h.mapping.ToApiResponsePaginationTopup(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Find a topup by ID
// @Tags Topup
// @Security Bearer
// @Description Retrieve a topup record using its ID
// @Accept json
// @Produce json
// @Param id path string true "Topup ID"
// @Success 200 {object} response.ApiResponseTopup "Topup data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve topup data"
// @Router /api/topups/{id} [get]
func (h topupHandleApi) FindById(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return topup_errors.ErrApiInvalidTopupID(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindByIdTopup(ctx, &pb.FindByIdTopupRequest{
		TopupId: int32(idInt),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve topup data", zap.Error(err))

		return topup_errors.ErrApiFailedFindByIdTopup(c)
	}

	so := h.mapping.ToApiResponseTopup(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTopupStatusSuccess retrieves the monthly top-up status for successful transactions.
// @Summary Get monthly top-up status for successful transactions
// @Tags Topup
// @Security Bearer
// @Description Retrieve the monthly top-up status for successful transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseTopupMonthStatusSuccess "Monthly top-up status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly top-up status for successful transactions"
// @Router /api/topups/monthly-success [get]
func (h *topupHandleApi) FindMonthlyTopupStatusSuccess(c echo.Context) error {
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return topup_errors.ErrApiInvalidYear(c)
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return topup_errors.ErrApiInvalidMonth(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthlyTopupStatusSuccess(ctx, &pb.FindMonthlyTopupStatus{
		Year:  int32(year),
		Month: int32(month),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve monthly topup status success", zap.Error(err))

		return topup_errors.ErrApiFailedFindMonthlyTopupStatusSuccess(c)
	}

	so := h.mapping.ToApiResponseTopupMonthStatusSuccess(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTopupStatusSuccess retrieves the yearly top-up status for successful transactions.
// @Summary Get yearly top-up status for successful transactions
// @Tags Topup
// @Security Bearer
// @Description Retrieve the yearly top-up status for successful transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupYearStatusSuccess "Yearly top-up status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly top-up status for successful transactions"
// @Router /api/topups/yearly-success [get]
func (h *topupHandleApi) FindYearlyTopupStatusSuccess(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return topup_errors.ErrApiInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearlyTopupStatusSuccess(ctx, &pb.FindYearTopupStatus{
		Year: int32(year),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve yearly topup status success", zap.Error(err))

		return topup_errors.ErrApiFailedFindYearlyTopupStatusSuccess(c)
	}

	so := h.mapping.ToApiResponseTopupYearStatusSuccess(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTopupStatusFailed retrieves the monthly top-up status for failed transactions.
// @Summary Get monthly top-up status for failed transactions
// @Tags Topup
// @Security Bearer
// @Description Retrieve the monthly top-up status for failed transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseTopupMonthStatusFailed "Monthly top-up status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly top-up status for failed transactions"
// @Router /api/topups/monthly-failed [get]
func (h *topupHandleApi) FindMonthlyTopupStatusFailed(c echo.Context) error {
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return topup_errors.ErrApiInvalidYear(c)
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return topup_errors.ErrApiInvalidMonth(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthlyTopupStatusFailed(ctx, &pb.FindMonthlyTopupStatus{
		Year:  int32(year),
		Month: int32(month),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve monthly topup status failed", zap.Error(err))

		return topup_errors.ErrApiFailedFindMonthlyTopupStatusFailed(c)
	}

	so := h.mapping.ToApiResponseTopupMonthStatusFailed(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTopupStatusFailed retrieves the yearly top-up status for failed transactions.
// @Summary Get yearly top-up status for failed transactions
// @Tags Topup
// @Security Bearer
// @Description Retrieve the yearly top-up status for failed transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupYearStatusFailed "Yearly top-up status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly top-up status for failed transactions"
// @Router /api/topups/yearly-failed [get]
func (h *topupHandleApi) FindYearlyTopupStatusFailed(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return topup_errors.ErrApiInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearlyTopupStatusFailed(ctx, &pb.FindYearTopupStatus{
		Year: int32(year),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve yearly topup status failed", zap.Error(err))

		return topup_errors.ErrApiFailedFindYearlyTopupStatusFailed(c)
	}

	so := h.mapping.ToApiResponseTopupYearStatusFailed(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTopupStatusSuccess retrieves the monthly top-up status for successful transactions.
// @Summary Get monthly top-up status for successful transactions
// @Tags Topup
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
// @Router /api/topups/monthly-success [get]
func (h *topupHandleApi) FindMonthlyTopupStatusSuccessByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")
	cardNumber := c.QueryParam("card_number")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return topup_errors.ErrApiInvalidYear(c)
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return topup_errors.ErrApiInvalidMonth(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthlyTopupStatusSuccessByCardNumber(ctx, &pb.FindMonthlyTopupStatusCardNumber{
		Year:       int32(year),
		Month:      int32(month),
		CardNumber: cardNumber,
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve monthly topup status success", zap.Error(err))

		return topup_errors.ErrApiFailedFindMonthlyTopupStatusSuccess(c)
	}

	so := h.mapping.ToApiResponseTopupMonthStatusSuccess(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTopupStatusSuccess retrieves the yearly top-up status for successful transactions.
// @Summary Get yearly top-up status for successful transactions
// @Tags Topup
// @Security Bearer
// @Description Retrieve the yearly top-up status for successful transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseTopupYearStatusSuccess "Yearly top-up status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly top-up status for successful transactions"
// @Router /api/topups/yearly-success [get]
func (h *topupHandleApi) FindYearlyTopupStatusSuccessByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	cardNumber := c.QueryParam("card_number")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return topup_errors.ErrApiInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearlyTopupStatusSuccessByCardNumber(ctx, &pb.FindYearTopupStatusCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve yearly topup status success", zap.Error(err))

		return topup_errors.ErrApiFailedFindYearlyTopupStatusSuccess(c)
	}

	so := h.mapping.ToApiResponseTopupYearStatusSuccess(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTopupStatusFailed retrieves the monthly top-up status for failed transactions.
// @Summary Get monthly top-up status for failed transactions
// @Tags Topup
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
// @Router /api/topups/monthly-failed [get]
func (h *topupHandleApi) FindMonthlyTopupStatusFailedByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")
	cardNumber := c.QueryParam("card_number")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return topup_errors.ErrApiInvalidYear(c)
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return topup_errors.ErrApiInvalidMonth(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthlyTopupStatusFailedByCardNumber(ctx, &pb.FindMonthlyTopupStatusCardNumber{
		Year:       int32(year),
		Month:      int32(month),
		CardNumber: cardNumber,
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve monthly topup status failed", zap.Error(err))

		return topup_errors.ErrApiFailedFindMonthlyTopupStatusFailed(c)
	}

	so := h.mapping.ToApiResponseTopupMonthStatusFailed(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTopupStatusFailedByCardNumber retrieves the yearly top-up status for failed transactions.
// @Summary Get yearly top-up status for failed transactions
// @Tags Topup
// @Security Bearer
// @Description Retrieve the yearly top-up status for failed transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseTopupYearStatusFailed "Yearly top-up status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly top-up status for failed transactions"
// @Router /api/topups/yearly-failed [get]
func (h *topupHandleApi) FindYearlyTopupStatusFailedByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	cardNumber := c.QueryParam("card_number")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return topup_errors.ErrApiInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearlyTopupStatusFailedByCardNumber(ctx, &pb.FindYearTopupStatusCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve yearly topup status failed", zap.Error(err))

		return topup_errors.ErrApiFailedFindYearlyTopupStatusFailed(c)
	}

	so := h.mapping.ToApiResponseTopupYearStatusFailed(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTopupMethods retrieves the monthly top-up methods for a specific year.
// @Summary Get monthly top-up methods
// @Tags Topup
// @Security Bearer
// @Description Retrieve the monthly top-up methods for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupMonthMethod "Monthly top-up methods"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly top-up methods"
// @Router /api/topups/monthly-methods [get]
func (h *topupHandleApi) FindMonthlyTopupMethods(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return topup_errors.ErrApiInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthlyTopupMethods(ctx, &pb.FindYearTopupStatus{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly topup methods", zap.Error(err))
		return topup_errors.ErrApiFailedFindMonthlyTopupMethods(c)
	}

	so := h.mapping.ToApiResponseTopupMonthMethod(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTopupMethods retrieves the yearly top-up methods for a specific year.
// @Summary Get yearly top-up methods
// @Tags Topup
// @Security Bearer
// @Description Retrieve the yearly top-up methods for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupYearMethod "Yearly top-up methods"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly top-up methods"
// @Router /api/topups/yearly-methods [get]
func (h *topupHandleApi) FindYearlyTopupMethods(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return topup_errors.ErrApiInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearlyTopupMethods(ctx, &pb.FindYearTopupStatus{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly topup methods", zap.Error(err))
		return topup_errors.ErrApiFailedFindYearlyTopupMethods(c)
	}

	so := h.mapping.ToApiResponseTopupYearMethod(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTopupAmounts retrieves the monthly top-up amounts for a specific year.
// @Summary Get monthly top-up amounts
// @Tags Topup
// @Security Bearer
// @Description Retrieve the monthly top-up amounts for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupMonthAmount "Monthly top-up amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly top-up amounts"
// @Router /api/topup/monthly-amounts [get]
func (h *topupHandleApi) FindMonthlyTopupAmounts(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return topup_errors.ErrApiInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthlyTopupAmounts(ctx, &pb.FindYearTopupStatus{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly topup amounts", zap.Error(err))
		return topup_errors.ErrApiFailedFindMonthlyTopupAmounts(c)
	}

	so := h.mapping.ToApiResponseTopupMonthAmount(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTopupAmounts retrieves the yearly top-up amounts for a specific year.
// @Summary Get yearly top-up amounts
// @Tags Topup
// @Security Bearer
// @Description Retrieve the yearly top-up amounts for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupYearAmount "Yearly top-up amounts"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly top-up amounts"
// @Router /api/topups/yearly-amounts [get]
func (h *topupHandleApi) FindYearlyTopupAmounts(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "error",
			Message: "Invalid year parameter",
			Code:    http.StatusBadRequest,
		})
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearlyTopupAmounts(ctx, &pb.FindYearTopupStatus{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly topup amounts", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve yearly topup amounts",
			Code:    http.StatusInternalServerError,
		})
	}

	so := h.mapping.ToApiResponseTopupYearAmount(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTopupMethodsByCardNumber retrieves the monthly top-up methods for a specific card number and year.
// @Summary Get monthly top-up methods by card number
// @Tags Topup
// @Security Bearer
// @Description Retrieve the monthly top-up methods for a specific card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupMonthMethod "Monthly top-up methods by card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly top-up methods by card number"
// @Router /api/topups/monthly-methods-by-card [get]
func (h *topupHandleApi) FindMonthlyTopupMethodsByCardNumber(c echo.Context) error {
	cardNumber := c.QueryParam("card_number")
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return topup_errors.ErrApiInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthlyTopupMethodsByCardNumber(ctx, &pb.FindYearTopupCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly topup methods by card number", zap.Error(err))
		return topup_errors.ErrApiFailedFindMonthlyTopupMethodsByCardNumber(c)
	}

	so := h.mapping.ToApiResponseTopupMonthMethod(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTopupMethodsByCardNumber retrieves the yearly top-up methods for a specific card number and year.
// @Summary Get yearly top-up methods by card number
// @Tags Topup
// @Security Bearer
// @Description Retrieve the yearly top-up methods for a specific card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupYearMethod "Yearly top-up methods by card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly top-up methods by card number"
// @Router /api/topups/yearly-methods-by-card [get]
func (h *topupHandleApi) FindYearlyTopupMethodsByCardNumber(c echo.Context) error {
	cardNumber := c.QueryParam("card_number")
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return topup_errors.ErrApiInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearlyTopupMethodsByCardNumber(ctx, &pb.FindYearTopupCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly topup methods by card number", zap.Error(err))
		return topup_errors.ErrApiFailedFindYearlyTopupMethodsByCardNumber(c)
	}

	so := h.mapping.ToApiResponseTopupYearMethod(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTopupAmountsByCardNumber retrieves the monthly top-up amounts for a specific card number and year.
// @Summary Get monthly top-up amounts by card number
// @Tags Topup
// @Security Bearer
// @Description Retrieve the monthly top-up amounts for a specific card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupMonthAmount "Monthly top-up amounts by card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly top-up amounts by card number"
// @Router /api/topups/monthly-amounts-by-card [get]
func (h *topupHandleApi) FindMonthlyTopupAmountsByCardNumber(c echo.Context) error {
	cardNumber := c.QueryParam("card_number")
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)

	if cardNumber == "" {
		h.logger.Debug("Invalid card number parameter")
		return topup_errors.ErrApiInvalidCardNumber(c)
	}

	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return topup_errors.ErrApiInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthlyTopupAmountsByCardNumber(ctx, &pb.FindYearTopupCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly topup amounts by card number", zap.Error(err))
		return topup_errors.ErrApiFailedFindMonthlyTopupAmountsByCardNumber(c)
	}

	so := h.mapping.ToApiResponseTopupMonthAmount(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTopupAmountsByCardNumber retrieves the yearly top-up amounts for a specific card number and year.
// @Summary Get yearly top-up amounts by card number
// @Tags Topup
// @Security Bearer
// @Description Retrieve the yearly top-up amounts for a specific card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseTopupYearAmount "Yearly top-up amounts by card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly top-up amounts by card number"
// @Router /api/topups/yearly-amounts-by-card [get]
func (h *topupHandleApi) FindYearlyTopupAmountsByCardNumber(c echo.Context) error {
	cardNumber := c.QueryParam("card_number")
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return topup_errors.ErrApiInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearlyTopupAmountsByCardNumber(ctx, &pb.FindYearTopupCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly topup amounts by card number", zap.Error(err))
		return topup_errors.ErrApiFailedFindYearlyTopupAmountsByCardNumber(c)
	}

	so := h.mapping.ToApiResponseTopupYearAmount(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Find active topups
// @Tags Topup
// @Security Bearer
// @Description Retrieve a list of active topup records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsesTopup "Active topup data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve topup data"
// @Router /api/topups/active [get]
func (h *topupHandleApi) FindByActive(c echo.Context) error {
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

	req := &pb.FindAllTopupRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to retrieve topup data", zap.Error(err))

		return topup_errors.ErrApiFailedFindAllTopupsActive(c)
	}

	so := h.mapping.ToApiResponsePaginationTopupDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Retrieve trashed topups
// @Tags Topup
// @Security Bearer
// @Description Retrieve a list of trashed topup records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsesTopup "List of trashed topup data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve topup data"
// @Router /api/topups/trashed [get]
func (h *topupHandleApi) FindByTrashed(c echo.Context) error {
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

	req := &pb.FindAllTopupRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to retrieve topup data", zap.Error(err))

		return topup_errors.ErrApiFailedFindAllTopupsTrashed(c)
	}

	so := h.mapping.ToApiResponsePaginationTopupDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Create topup
// @Tags Topup
// @Security Bearer
// @Description Create a new topup record
// @Accept json
// @Produce json
// @Param CreateTopupRequest body requests.CreateTopupRequest true "Create topup request"
// @Success 200 {object} response.ApiResponseTopup "Created topup data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: "
// @Failure 500 {object} response.ErrorResponse "Failed to create topup: "
// @Router /api/topups/create [post]
func (h *topupHandleApi) Create(c echo.Context) error {
	var body requests.CreateTopupRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Bad Request", zap.Error(err))

		return topup_errors.ErrApiBindCreateTopup(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation Error", zap.Error(err))

		return topup_errors.ErrApiValidateCreateTopup(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.CreateTopup(ctx, &pb.CreateTopupRequest{
		CardNumber:  body.CardNumber,
		TopupAmount: int32(body.TopupAmount),
		TopupMethod: body.TopupMethod,
	})

	if err != nil {
		h.logger.Debug("Failed to create topup", zap.Error(err))

		return topup_errors.ErrApiFailedCreateTopup(c)
	}

	so := h.mapping.ToApiResponseTopup(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Update topup
// @Tags Topup
// @Security Bearer
// @Description Update an existing topup record with the provided details
// @Accept json
// @Produce json
// @Param id path int true "Topup ID"
// @Param UpdateTopupRequest body requests.UpdateTopupRequest true "Update topup request"
// @Success 200 {object} response.ApiResponseTopup "Updated topup data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid input data"
// @Failure 500 {object} response.ErrorResponse "Failed to update topup: "
// @Router /api/topups/update/{id} [post]
func (h *topupHandleApi) Update(c echo.Context) error {
	idint, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Bad Request", zap.Error(err))

		return topup_errors.ErrApiInvalidTopupID(c)
	}

	var body requests.UpdateTopupRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Bad Request", zap.Error(err))

		return topup_errors.ErrApiBindUpdateTopup(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation Error", zap.Error(err))

		return topup_errors.ErrApiValidateUpdateTopup(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.UpdateTopup(ctx, &pb.UpdateTopupRequest{
		TopupId:     int32(idint),
		CardNumber:  body.CardNumber,
		TopupAmount: int32(body.TopupAmount),
		TopupMethod: body.TopupMethod,
	})

	if err != nil {
		h.logger.Debug("Failed to update topup", zap.Error(err))

		return topup_errors.ErrApiFailedUpdateTopup(c)
	}

	so := h.mapping.ToApiResponseTopup(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Trash a topup
// @Tags Topup
// @Security Bearer
// @Description Trash a topup record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Topup ID"
// @Success 200 {object} response.ApiResponseTopup "Successfully trashed topup record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to trash topup:"
// @Router /api/topups/trash/{id} [post]
func (h *topupHandleApi) TrashTopup(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return topup_errors.ErrApiInvalidTopupID(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.TrashedTopup(ctx, &pb.FindByIdTopupRequest{
		TopupId: int32(idInt),
	})

	if err != nil {
		h.logger.Debug("Failed to trashed topup", zap.Error(err))

		return topup_errors.ErrApiFailedTrashTopup(c)
	}

	so := h.mapping.ToApiResponseTopupDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Restore a trashed topup
// @Tags Topup
// @Security Bearer
// @Description Restore a trashed topup record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Topup ID"
// @Success 200 {object} response.ApiResponseTopup "Successfully restored topup record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore topup:"
// @Router /api/topups/restore/{id} [post]
func (h *topupHandleApi) RestoreTopup(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		h.logger.Debug("Bad Request: Invalid ID", zap.Error(err))

		return topup_errors.ErrApiInvalidTopupID(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.RestoreTopup(ctx, &pb.FindByIdTopupRequest{
		TopupId: int32(idInt),
	})

	if err != nil {
		h.logger.Debug("Failed to restore topup", zap.Error(err))

		return topup_errors.ErrApiFailedRestoreTopup(c)
	}

	so := h.mapping.ToApiResponseTopupDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Permanently delete a topup
// @Tags Topup
// @Security Bearer
// @Description Permanently delete a topup record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Topup ID"
// @Success 200 {object} response.ApiResponseTopupDelete "Successfully deleted topup record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete topup:"
// @Router /api/topups/permanent/{id} [delete]
func (h *topupHandleApi) DeleteTopupPermanent(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		h.logger.Debug("Bad Request: Invalid ID", zap.Error(err))

		return topup_errors.ErrApiInvalidTopupID(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.DeleteTopupPermanent(ctx, &pb.FindByIdTopupRequest{
		TopupId: int32(idInt),
	})

	if err != nil {
		h.logger.Debug("Failed to delete topup", zap.Error(err))

		return topup_errors.ErrApiFailedDeletePermanentTopup(c)
	}

	so := h.mapping.ToApiResponseTopupDelete(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Restore all topup records
// @Tags Topup
// @Security Bearer
// @Description Restore all topup records that were previously deleted.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseTopupAll "Successfully restored all topup records"
// @Failure 500 {object} response.ErrorResponse "Failed to restore all topup records"
// @Router /api/topups/restore/all [post]
func (h *topupHandleApi) RestoreAllTopup(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllTopup(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Failed to restore all topup", zap.Error(err))
		return topup_errors.ErrApiFailedRestoreAllTopup(c)
	}

	h.logger.Debug("Successfully restored all topup")

	so := h.mapping.ToApiResponseTopupAll(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Permanently delete all topup records
// @Tags Topup
// @Security Bearer
// @Description Permanently delete all topup records from the database.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseTopupAll "Successfully deleted all topup records permanently"
// @Failure 500 {object} response.ErrorResponse "Failed to permanently delete all topup records"
// @Router /api/topups/permanent/all [post]
func (h *topupHandleApi) DeleteAllTopupPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.DeleteAllTopupPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Failed to permanently delete all topup", zap.Error(err))

		return topup_errors.ErrApiFailedDeleteAllTopupPermanent(c)
	}

	h.logger.Debug("Successfully deleted all topup permanently")

	so := h.mapping.ToApiResponseTopupAll(res)

	return c.JSON(http.StatusOK, so)
}
