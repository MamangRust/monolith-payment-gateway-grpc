package handler

import (
	"net/http"
	"strconv"

	"github.com/MamangRust/payment-gateway-monolith-grpc/pkg/logger"
	"github.com/MamangRust/payment-gateway-monolith-grpc/shared/domain/requests"
	"github.com/MamangRust/payment-gateway-monolith-grpc/shared/errors/withdraw_errors"
	apimapper "github.com/MamangRust/payment-gateway-monolith-grpc/shared/mapper/response/api"
	"github.com/MamangRust/payment-gateway-monolith-grpc/shared/pb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type withdrawHandleApi struct {
	client  pb.WithdrawServiceClient
	logger  logger.LoggerInterface
	mapping apimapper.WithdrawResponseMapper
}

func NewHandlerWithdraw(client pb.WithdrawServiceClient, router *echo.Echo, logger logger.LoggerInterface, mapping apimapper.WithdrawResponseMapper) *withdrawHandleApi {
	withdrawHandler := &withdrawHandleApi{
		client:  client,
		logger:  logger,
		mapping: mapping,
	}
	routerWithdraw := router.Group("/api/withdraws")

	routerWithdraw.GET("", withdrawHandler.FindAll)
	routerWithdraw.GET("/card-number/:card_number", withdrawHandler.FindAllByCardNumber)

	routerWithdraw.GET("/:id", withdrawHandler.FindById)

	routerWithdraw.GET("/monthly-success", withdrawHandler.FindMonthlyWithdrawStatusSuccess)
	routerWithdraw.GET("/yearly-success", withdrawHandler.FindYearlyWithdrawStatusSuccess)
	routerWithdraw.GET("/monthly-failed", withdrawHandler.FindMonthlyWithdrawStatusFailed)
	routerWithdraw.GET("/yearly-failed", withdrawHandler.FindYearlyWithdrawStatusFailed)

	routerWithdraw.GET("/monthly-success-by-card", withdrawHandler.FindMonthlyWithdrawStatusSuccessByCardNumber)
	routerWithdraw.GET("/yearly-success-by-card", withdrawHandler.FindYearlyWithdrawStatusSuccessByCardNumber)
	routerWithdraw.GET("/monthly-failed-by-card", withdrawHandler.FindMonthlyWithdrawStatusFailedByCardNumber)
	routerWithdraw.GET("/yearly-failed-by-card", withdrawHandler.FindYearlyWithdrawStatusFailedByCardNumber)

	routerWithdraw.GET("/monthly-amount", withdrawHandler.FindMonthlyWithdraws)
	routerWithdraw.GET("/yearly-amount", withdrawHandler.FindYearlyWithdraws)

	routerWithdraw.GET("/monthly-amount-card", withdrawHandler.FindMonthlyWithdrawsByCardNumber)
	routerWithdraw.GET("/yearly-amount-card", withdrawHandler.FindYearlyWithdrawsByCardNumber)

	routerWithdraw.GET("/active", withdrawHandler.FindByActive)
	routerWithdraw.GET("/trashed", withdrawHandler.FindByTrashed)
	routerWithdraw.POST("/create", withdrawHandler.Create)
	routerWithdraw.POST("/update/:id", withdrawHandler.Update)

	routerWithdraw.POST("/trashed/:id", withdrawHandler.TrashWithdraw)
	routerWithdraw.POST("/restore/:id", withdrawHandler.RestoreWithdraw)
	routerWithdraw.DELETE("/permanent/:id", withdrawHandler.DeleteWithdrawPermanent)

	routerWithdraw.POST("/restore/all", withdrawHandler.RestoreAllWithdraw)
	routerWithdraw.POST("/permanent/all", withdrawHandler.DeleteAllWithdrawPermanent)

	return withdrawHandler
}

// @Summary Find all withdraw records
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve a list of all withdraw records with pagination and search
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationWithdraw "List of withdraw records"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve withdraw data"
// @Router /api/withdraw [get]
func (h *withdrawHandleApi) FindAll(c echo.Context) error {
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

	req := &pb.FindAllWithdrawRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAllWithdraw(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to retrieve withdraw data", zap.Error(err))

		return withdraw_errors.ErrApiFailedFindAllWithdraw(c)
	}

	so := h.mapping.ToApiResponsePaginationWithdraw(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Find all withdraw records by card number
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve a list of withdraw records for a specific card number with pagination and search
// @Accept json
// @Produce json
// @Param card_number path string true "Card Number"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationWithdraw "List of withdraw records"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve withdraw data"
// @Router /api/withdraw/card-number/{card_number} [get]
func (h *withdrawHandleApi) FindAllByCardNumber(c echo.Context) error {
	cardNumber := c.Param("card_number")
	if cardNumber == "" {
		return withdraw_errors.ErrApiInvalidCardNumber(c)
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

	req := &pb.FindAllWithdrawByCardNumberRequest{
		CardNumber: cardNumber,
		Page:       int32(page),
		PageSize:   int32(pageSize),
		Search:     search,
	}

	res, err := h.client.FindAllWithdrawByCardNumber(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to retrieve withdraw data", zap.Error(err))

		return withdraw_errors.ErrApiFailedFindAllWithdrawByCardNumber(c)
	}

	so := h.mapping.ToApiResponsePaginationWithdraw(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Find a withdraw by ID
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve a withdraw record using its ID
// @Accept json
// @Produce json
// @Param id path int true "Withdraw ID"
// @Success 200 {object} response.ApiResponseWithdraw "Withdraw data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve withdraw data"
// @Router /api/withdraw/{id} [get]
func (h *withdrawHandleApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid withdraw ID", zap.Error(err))

		return withdraw_errors.ErrApiWithdrawInvalidID(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdWithdrawRequest{
		WithdrawId: int32(id),
	}

	withdraw, err := h.client.FindByIdWithdraw(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to retrieve withdraw data", zap.Error(err))

		return withdraw_errors.ErrApiFailedFindByIdWithdraw(c)
	}

	so := h.mapping.ToApiResponseWithdraw(withdraw)

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyWithdrawStatusSuccess retrieves the monthly withdraw status for successful transactions.
// @Summary Get monthly withdraw status for successful transactions
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve the monthly withdraw status for successful transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseWithdrawMonthStatusSuccess "Monthly withdraw status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly withdraw status for successful transactions"
// @Router /api/withdraws/monthly-success [get]
func (h *withdrawHandleApi) FindMonthlyWithdrawStatusSuccess(c echo.Context) error {
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return withdraw_errors.ErrApiInvalidYear(c)
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return withdraw_errors.ErrApiInvalidMonth(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthlyWithdrawStatusSuccess(ctx, &pb.FindMonthlyWithdrawStatus{
		Year:  int32(year),
		Month: int32(month),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve monthly Withdraw status success", zap.Error(err))

		return withdraw_errors.ErrApiFailedFindMonthlyWithdrawStatusSuccess(c)
	}

	so := h.mapping.ToApiResponseWithdrawMonthStatusSuccess(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearlyWithdrawStatusSuccess retrieves the yearly withdraw status for successful transactions.
// @Summary Get yearly withdraw status for successful transactions
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve the yearly withdraw status for successful transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseWithdrawYearStatusSuccess "Yearly withdraw status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly withdraw status for successful transactions"
// @Router /api/withdraws/yearly-success [get]
func (h *withdrawHandleApi) FindYearlyWithdrawStatusSuccess(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return withdraw_errors.ErrApiInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearlyWithdrawStatusSuccess(ctx, &pb.FindYearWithdrawStatus{
		Year: int32(year),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve yearly Withdraw status success", zap.Error(err))

		return withdraw_errors.ErrApiFailedFindYearlyWithdrawStatusSuccess(c)
	}

	so := h.mapping.ToApiResponseWithdrawYearStatusSuccess(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyWithdrawStatusFailed retrieves the monthly withdraw status for failed transactions.
// @Summary Get monthly withdraw status for failed transactions
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve the monthly withdraw status for failed transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseWithdrawMonthStatusFailed "Monthly withdraw status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly withdraw status for failed transactions"
// @Router /api/withdraws/monthly-failed [get]
func (h *withdrawHandleApi) FindMonthlyWithdrawStatusFailed(c echo.Context) error {
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return withdraw_errors.ErrApiInvalidYear(c)
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return withdraw_errors.ErrApiInvalidMonth(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthlyWithdrawStatusFailed(ctx, &pb.FindMonthlyWithdrawStatus{
		Year:  int32(year),
		Month: int32(month),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve monthly Withdraw status Failed", zap.Error(err))

		return withdraw_errors.ErrApiFailedFindMonthlyWithdrawStatusFailed(c)
	}

	so := h.mapping.ToApiResponseWithdrawMonthStatusFailed(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearlyWithdrawStatusFailed retrieves the yearly withdraw status for failed transactions.
// @Summary Get yearly withdraw status for failed transactions
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve the yearly withdraw status for failed transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseWithdrawYearStatusSuccess "Yearly withdraw status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly withdraw status for failed transactions"
// @Router /api/withdraws/yearly-failed [get]
func (h *withdrawHandleApi) FindYearlyWithdrawStatusFailed(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return withdraw_errors.ErrApiInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearlyWithdrawStatusFailed(ctx, &pb.FindYearWithdrawStatus{
		Year: int32(year),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve yearly Withdraw status Failed", zap.Error(err))

		return withdraw_errors.ErrApiFailedFindYearlyWithdrawStatusFailed(c)
	}

	so := h.mapping.ToApiResponseWithdrawYearStatusFailed(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyWithdrawStatusSuccessByCardNumber retrieves the monthly withdraw status for successful transactions.
// @Summary Get monthly withdraw status for successful transactions
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve the monthly withdraw status for successful transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseWithdrawMonthStatusSuccess "Monthly withdraw status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly withdraw status for successful transactions"
// @Router /api/withdraws/monthly-success-by-card [get]
func (h *withdrawHandleApi) FindMonthlyWithdrawStatusSuccessByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")
	cardNumber := c.QueryParam("card_number")

	if cardNumber == "" {
		return withdraw_errors.ErrApiInvalidCardNumber(c)
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return withdraw_errors.ErrApiInvalidYear(c)
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return withdraw_errors.ErrApiInvalidMonth(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthlyWithdrawStatusSuccessCardNumber(ctx, &pb.FindMonthlyWithdrawStatusCardNumber{
		Year:       int32(year),
		Month:      int32(month),
		CardNumber: cardNumber,
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve monthly Withdraw status success", zap.Error(err))

		return withdraw_errors.ErrApiFailedFindMonthlyWithdrawStatusSuccessCardNumber(c)
	}

	so := h.mapping.ToApiResponseWithdrawMonthStatusSuccess(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearlyWithdrawStatusSuccessByCardNumber retrieves the yearly withdraw status for successful transactions.
// @Summary Get yearly withdraw status for successful transactions
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve the yearly withdraw status for successful transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseWithdrawYearStatusSuccess "Yearly withdraw status for successful transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly withdraw status for successful transactions"
// @Router /api/withdraws/yearly-success-by-card-number [get]
func (h *withdrawHandleApi) FindYearlyWithdrawStatusSuccessByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	card_number := c.QueryParam("card_number")

	if card_number == "" {
		return withdraw_errors.ErrApiInvalidCardNumber(c)
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return withdraw_errors.ErrApiInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearlyWithdrawStatusSuccessCardNumber(ctx, &pb.FindYearWithdrawStatusCardNumber{
		CardNumber: card_number,
		Year:       int32(year),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve yearly Withdraw status success", zap.Error(err))

		return withdraw_errors.ErrApiFailedFindYearlyWithdrawStatusSuccessCardNumber(c)
	}

	so := h.mapping.ToApiResponseWithdrawYearStatusSuccess(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyWithdrawStatusFailedByCardNumber retrieves the monthly withdraw status for failed transactions.
// @Summary Get monthly withdraw status for failed transactions
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve the monthly withdraw status for failed transactions by year and month.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseWithdrawMonthStatusFailed "Monthly withdraw status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly withdraw status for failed transactions"
// @Router /api/withdraws/monthly-failed-by-card [get]
func (h *withdrawHandleApi) FindMonthlyWithdrawStatusFailedByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")
	card_number := c.QueryParam("card_number")

	if card_number == "" {
		return withdraw_errors.ErrApiInvalidCardNumber(c)
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return withdraw_errors.ErrApiInvalidYear(c)
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return withdraw_errors.ErrApiInvalidMonth(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthlyWithdrawStatusFailedCardNumber(ctx, &pb.FindMonthlyWithdrawStatusCardNumber{
		Year:       int32(year),
		Month:      int32(month),
		CardNumber: card_number,
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve monthly Withdraw status Failed", zap.Error(err))

		return withdraw_errors.ErrApiFailedFindMonthlyWithdrawStatusFailedCardNumber(c)
	}

	so := h.mapping.ToApiResponseWithdrawMonthStatusFailed(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearlyWithdrawStatusFailedByCardNumber retrieves the yearly withdraw status for failed transactions.
// @Summary Get yearly withdraw status for failed transactions
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve the yearly withdraw status for failed transactions by year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseWithdrawYearStatusSuccess "Yearly withdraw status for failed transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid year"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly withdraw status for failed transactions"
// @Router /api/withdraws/yearly-failed-by-card [get]
func (h *withdrawHandleApi) FindYearlyWithdrawStatusFailedByCardNumber(c echo.Context) error {
	yearStr := c.QueryParam("year")
	cardNumber := c.QueryParam("card_number")

	if cardNumber == "" {
		return withdraw_errors.ErrApiInvalidCardNumber(c)
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return withdraw_errors.ErrApiInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearlyWithdrawStatusFailedCardNumber(ctx, &pb.FindYearWithdrawStatusCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve yearly Withdraw status Failed", zap.Error(err))

		return withdraw_errors.ErrApiFailedFindYearlyWithdrawStatusFailedCardNumber(c)
	}

	so := h.mapping.ToApiResponseWithdrawYearStatusFailed(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyWithdraws retrieves the monthly withdraws for a specific year.
// @Summary Get monthly withdraws
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve the monthly withdraws for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseWithdrawMonthAmount "Monthly withdraws"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly withdraws"
// @Router /api/withdraws/monthly [get]
func (h *withdrawHandleApi) FindMonthlyWithdraws(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return withdraw_errors.ErrApiInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthlyWithdraws(ctx, &pb.FindYearWithdrawStatus{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly withdraws", zap.Error(err))
		return withdraw_errors.ErrApiFailedFindMonthlyWithdraws(c)
	}

	so := h.mapping.ToApiResponseWithdrawMonthAmount(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearlyWithdraws retrieves the yearly withdraws for a specific year.
// @Summary Get yearly withdraws
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve the yearly withdraws for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseWithdrawYearAmount "Yearly withdraws"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly withdraws"
// @Router /api/withdraws/yearly [get]
func (h *withdrawHandleApi) FindYearlyWithdraws(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return withdraw_errors.ErrApiInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearlyWithdraws(ctx, &pb.FindYearWithdrawStatus{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly withdraws", zap.Error(err))
		return withdraw_errors.ErrApiFailedFindYearlyWithdraws(c)
	}

	so := h.mapping.ToApiResponseWithdrawYearAmount(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyWithdrawsByCardNumber retrieves the monthly withdraws for a specific card number and year.
// @Summary Get monthly withdraws by card number
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve the monthly withdraws for a specific card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseWithdrawMonthAmount "Monthly withdraws by card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly withdraws by card number"
// @Router /api/withdraws/monthly-by-card [get]
func (h *withdrawHandleApi) FindMonthlyWithdrawsByCardNumber(c echo.Context) error {
	cardNumber := c.QueryParam("card_number")
	yearStr := c.QueryParam("year")

	if cardNumber == "" {
		h.logger.Debug("Invalid card number parameter")
		return withdraw_errors.ErrApiInvalidCardNumber(c)
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return withdraw_errors.ErrApiInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthlyWithdrawsByCardNumber(ctx, &pb.FindYearWithdrawCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly withdraws by card number", zap.Error(err))
		return withdraw_errors.ErrApiFailedFindMonthlyWithdrawsByCardNumber(c)
	}

	so := h.mapping.ToApiResponseWithdrawMonthAmount(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearlyWithdrawsByCardNumber retrieves the yearly withdraws for a specific card number and year.
// @Summary Get yearly withdraws by card number
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve the yearly withdraws for a specific card number and year.
// @Accept json
// @Produce json
// @Param card_number query string true "Card Number"
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseWithdrawYearAmount "Yearly withdraws by card number"
// @Failure 400 {object} response.ErrorResponse "Invalid card number or year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly withdraws by card number"
// @Router /api/withdraws/yearly-by-card [get]
func (h *withdrawHandleApi) FindYearlyWithdrawsByCardNumber(c echo.Context) error {
	cardNumber := c.QueryParam("card_number")
	yearStr := c.QueryParam("year")

	if cardNumber == "" {
		h.logger.Debug("Invalid card number parameter")
		return withdraw_errors.ErrApiInvalidCardNumber(c)
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return withdraw_errors.ErrApiInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearlyWithdrawsByCardNumber(ctx, &pb.FindYearWithdrawCardNumber{
		CardNumber: cardNumber,
		Year:       int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly withdraws by card number", zap.Error(err))
		return withdraw_errors.ErrApiFailedFindYearlyWithdrawsByCardNumber(c)
	}

	so := h.mapping.ToApiResponseWithdrawYearAmount(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Find a withdraw by card number
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve a withdraw record using its card number
// @Accept json
// @Produce json
// @Param card_number query string true "Card number"
// @Success 200 {object} response.ApiResponsesWithdraw "Withdraw data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid card number"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve withdraw data"
// @Router /api/withdraws/card/{card_number} [get]
func (h *withdrawHandleApi) FindByCardNumber(c echo.Context) error {
	cardNumber := c.QueryParam("card_number")

	ctx := c.Request().Context()

	req := &pb.FindByCardNumberRequest{
		CardNumber: cardNumber,
	}

	withdraw, err := h.client.FindByCardNumber(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to retrieve withdraw data", zap.Error(err))

		return withdraw_errors.ErrApiFailedFindByCardNumber(c)
	}

	so := h.mapping.ToApiResponsesWithdraw(withdraw)

	return c.JSON(http.StatusOK, so)
}

// @Summary Retrieve all active withdraw data
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve a list of all active withdraw data
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponsesWithdraw "List of withdraw data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve withdraw data"
// @Router /api/withdraws/active [get]
func (h *withdrawHandleApi) FindByActive(c echo.Context) error {
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

	req := &pb.FindAllWithdrawRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to retrieve withdraw data", zap.Error(err))

		return withdraw_errors.ErrApiFailedFindByActiveWithdraw(c)
	}

	so := h.mapping.ToApiResponsePaginationWithdrawDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Retrieve trashed withdraw data
// @Tags Withdraw
// @Security Bearer
// @Description Retrieve a list of trashed withdraw data
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponsesWithdraw "List of trashed withdraw data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve withdraw data"
// @Router /api/withdraws/trashed [get]
func (h *withdrawHandleApi) FindByTrashed(c echo.Context) error {
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

	req := &pb.FindAllWithdrawRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to retrieve withdraw data", zap.Error(err))

		return withdraw_errors.ErrApiFailedFindByTrashedWithdraw(c)
	}

	so := h.mapping.ToApiResponsePaginationWithdrawDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Create a new withdraw
// @Tags Withdraw
// @Security Bearer
// @Description Create a new withdraw record with the provided details.
// @Accept json
// @Produce json
// @Param CreateWithdrawRequest body requests.CreateWithdrawRequest true "Create Withdraw Request"
// @Success 200 {object} response.ApiResponseWithdraw "Successfully created withdraw record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create withdraw"
// @Router /api/withdraws/create [post]
func (h *withdrawHandleApi) Create(c echo.Context) error {
	var body requests.CreateWithdrawRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Invalid request body", zap.Error(err))

		return withdraw_errors.ErrApiBindCreateWithdraw(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation Error: " + err.Error())

		return withdraw_errors.ErrApiValidateCreateWithdraw(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.CreateWithdraw(ctx, &pb.CreateWithdrawRequest{
		CardNumber:     body.CardNumber,
		WithdrawAmount: int32(body.WithdrawAmount),
		WithdrawTime:   timestamppb.New(body.WithdrawTime),
	})

	if err != nil {
		h.logger.Debug("Failed to create withdraw", zap.Error(err))

		return withdraw_errors.ErrApiFailedCreateWithdraw(c)
	}

	so := h.mapping.ToApiResponseWithdraw(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Update an existing withdraw
// @Tags Withdraw
// @Security Bearer
// @Description Update an existing withdraw record with the provided details.
// @Accept json
// @Produce json
// @Param id path int true "Withdraw ID"
// @Param UpdateWithdrawRequest body requests.UpdateWithdrawRequest true "Update Withdraw Request"
// @Success 200 {object} response.ApiResponseWithdraw "Successfully updated withdraw record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update withdraw"
// @Router /api/withdraws/update/{id} [post]
func (h *withdrawHandleApi) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid withdraw ID", zap.Error(err))

		return withdraw_errors.ErrApiWithdrawInvalidID(c)
	}

	var body requests.UpdateWithdrawRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Invalid request body", zap.Error(err))

		return withdraw_errors.ErrApiBindUpdateWithdraw(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation Error: " + err.Error())

		return withdraw_errors.ErrApiValidateUpdateWithdraw(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.UpdateWithdraw(ctx, &pb.UpdateWithdrawRequest{
		WithdrawId:     int32(id),
		CardNumber:     body.CardNumber,
		WithdrawAmount: int32(body.WithdrawAmount),
		WithdrawTime:   timestamppb.New(body.WithdrawTime),
	})

	if err != nil {
		h.logger.Debug("Failed to update withdraw", zap.Error(err))

		return withdraw_errors.ErrApiFailedUpdateWithdraw(c)
	}

	so := h.mapping.ToApiResponseWithdraw(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Trash a withdraw by ID
// @Tags Withdraw
// @Security Bearer
// @Description Trash a withdraw using its ID
// @Accept json
// @Produce json
// @Param id path int true "Withdraw ID"
// @Success 200 {object} response.ApiResponseWithdraw "Withdaw data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to trash withdraw"
// @Router /api/withdraws/trashed/{id} [post]
func (h *withdrawHandleApi) TrashWithdraw(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid withdraw ID", zap.Error(err))

		return withdraw_errors.ErrApiWithdrawInvalidID(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.TrashedWithdraw(ctx, &pb.FindByIdWithdrawRequest{
		WithdrawId: int32(id),
	})

	if err != nil {
		h.logger.Debug("Failed to trash withdraw", zap.Error(err))

		return withdraw_errors.ErrApiFailedTrashedWithdraw(c)
	}

	so := h.mapping.ToApiResponseWithdraw(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Restore a withdraw by ID
// @Tags Withdraw
// @Security Bearer
// @Description Restore a withdraw by its ID
// @Accept json
// @Produce json
// @Param id path int true "Withdraw ID"
// @Success 200 {object} response.ApiResponseWithdraw "Withdraw data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore withdraw"
// @Router /api/withdraws/restore/{id} [post]
func (h *withdrawHandleApi) RestoreWithdraw(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid withdraw ID", zap.Error(err))
		return withdraw_errors.ErrApiWithdrawInvalidID(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.RestoreWithdraw(ctx, &pb.FindByIdWithdrawRequest{
		WithdrawId: int32(id),
	})

	if err != nil {
		h.logger.Debug("Failed to restore withdraw", zap.Error(err))

		return withdraw_errors.ErrApiFailedRestoreWithdraw(c)
	}

	so := h.mapping.ToApiResponseWithdraw(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Permanently delete a withdraw by ID
// @Tags Withdraw
// @Security Bearer
// @Description Permanently delete a withdraw by its ID
// @Accept json
// @Produce json
// @Param id path int true "Withdraw ID"
// @Success 200 {object} response.ApiResponseWithdrawDelete "Successfully deleted withdraw permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete withdraw permanently:"
// @Router /api/withdraws/permanent/{id} [delete]
func (h *withdrawHandleApi) DeleteWithdrawPermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid withdraw ID", zap.Error(err))

		return withdraw_errors.ErrApiWithdrawInvalidID(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.DeleteWithdrawPermanent(ctx, &pb.FindByIdWithdrawRequest{
		WithdrawId: int32(id),
	})

	if err != nil {
		h.logger.Debug("Failed to delete withdraw permanently", zap.Error(err))

		return withdraw_errors.ErrApiFailedDeleteWithdrawPermanent(c)
	}

	so := h.mapping.ToApiResponseWithdrawDelete(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Restore a withdraw all
// @Tags Withdraw
// @Security Bearer
// @Description Restore a withdraw all
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseWithdrawAll "Withdraw data"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore withdraw"
// @Router /api/withdraws/restore/all [post]
func (h *withdrawHandleApi) RestoreAllWithdraw(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllWithdraw(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Failed to restore all withdraw", zap.Error(err))
		return withdraw_errors.ErrApiFailedRestoreAllWithdraw(c)
	}

	h.logger.Debug("Successfully restored all withdraw")

	so := h.mapping.ToApiResponseWithdrawAll(res)

	return c.JSON(http.StatusOK, so)
}

// @Summary Permanently delete a withdraw by ID
// @Tags Withdraw
// @Security Bearer
// @Description Permanently delete a withdraw by its ID
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseWithdrawAll "Successfully deleted withdraw permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete withdraw permanently:"
// @Router /api/withdraws/permanent/all [post]
func (h *withdrawHandleApi) DeleteAllWithdrawPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.DeleteAllWithdrawPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Failed to permanently delete all withdraw", zap.Error(err))

		return withdraw_errors.ErrApiFailedDeleteAllWithdrawPermanent(c)
	}

	h.logger.Debug("Successfully deleted all withdraw permanently")

	so := h.mapping.ToApiResponseWithdrawAll(res)

	return c.JSON(http.StatusOK, so)
}
