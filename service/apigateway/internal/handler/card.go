package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api"
	"github.com/MamangRust/monolith-payment-gateway-shared/pb"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type cardHandleApi struct {
	card            pb.CardServiceClient
	logger          logger.LoggerInterface
	mapping         apimapper.CardResponseMapper
	trace           trace.Tracer
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewHandlerCard(card pb.CardServiceClient, router *echo.Echo, logger logger.LoggerInterface, mapper apimapper.CardResponseMapper) *cardHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "card_handler_requests_total",
			Help: "Total number of card requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "card_handler_request_duration_seconds",
			Help:    "Duration of card requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	cardHandler := &cardHandleApi{
		card:            card,
		logger:          logger,
		mapping:         mapper,
		trace:           otel.Tracer("card-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}
	routerCard := router.Group("/api/card")

	routerCard.GET("", cardHandler.FindAll)
	routerCard.GET("/:id", cardHandler.FindById)

	routerCard.GET("/dashboard", cardHandler.DashboardCard)
	routerCard.GET("/dashboard/:cardNumber", cardHandler.DashboardCardCardNumber)

	routerCard.GET("/monthly-balance", cardHandler.FindMonthlyBalance)
	routerCard.GET("/yearly-balance", cardHandler.FindYearlyBalance)

	routerCard.GET("/monthly-topup-amount", cardHandler.FindMonthlyTopupAmount)
	routerCard.GET("/yearly-topup-amount", cardHandler.FindYearlyTopupAmount)
	routerCard.GET("/monthly-withdraw-amount", cardHandler.FindMonthlyWithdrawAmount)
	routerCard.GET("/yearly-withdraw-amount", cardHandler.FindYearlyWithdrawAmount)

	routerCard.GET("/monthly-transaction-amount", cardHandler.FindMonthlyTransactionAmount)
	routerCard.GET("/yearly-transaction-amount", cardHandler.FindYearlyTransactionAmount)

	routerCard.GET("/monthly-transfer-sender-amount", cardHandler.FindMonthlyTransferSenderAmount)
	routerCard.GET("/yearly-transfer-sender-amount", cardHandler.FindYearlyTransferSenderAmount)
	routerCard.GET("/monthly-transfer-receiver-amount", cardHandler.FindMonthlyTransferReceiverAmount)
	routerCard.GET("/yearly-transfer-receiver-amount", cardHandler.FindYearlyTransferReceiverAmount)

	routerCard.GET("/monthly-balance-by-card", cardHandler.FindMonthlyBalanceByCardNumber)
	routerCard.GET("/yearly-balance-by-card", cardHandler.FindYearlyBalanceByCardNumber)
	routerCard.GET("/monthly-topup-amount-by-card", cardHandler.FindMonthlyTopupAmountByCardNumber)
	routerCard.GET("/yearly-topup-amount-by-card", cardHandler.FindYearlyTopupAmountByCardNumber)

	routerCard.GET("/monthly-withdraw-amount-by-card", cardHandler.FindMonthlyWithdrawAmountByCardNumber)
	routerCard.GET("/yearly-withdraw-amount-by-card", cardHandler.FindYearlyWithdrawAmountByCardNumber)
	routerCard.GET("/monthly-transaction-amount-by-card", cardHandler.FindMonthlyTransactionAmountByCardNumber)
	routerCard.GET("/yearly-transaction-amount-by-card", cardHandler.FindYearlyTransactionAmountByCardNumber)

	routerCard.GET("/monthly-transfer-sender-amount-by-card", cardHandler.FindMonthlyTransferSenderAmountByCardNumber)
	routerCard.GET("/yearly-transfer-sender-amount-by-card", cardHandler.FindYearlyTransferSenderAmountByCardNumber)
	routerCard.GET("/monthly-transfer-receiver-amount-by-card", cardHandler.FindMonthlyTransferReceiverAmountByCardNumber)
	routerCard.GET("/yearly-transfer-receiver-amount-by-card", cardHandler.FindYearlyTransferReceiverAmountByCardNumber)

	routerCard.GET("/user", cardHandler.FindByUserID)
	routerCard.GET("/active", cardHandler.FindByActive)
	routerCard.GET("/trashed", cardHandler.FindByTrashed)
	routerCard.GET("/card_number/:card_number", cardHandler.FindByCardNumber)

	routerCard.POST("/create", cardHandler.CreateCard)
	routerCard.POST("/update/:id", cardHandler.UpdateCard)
	routerCard.POST("/trashed/:id", cardHandler.TrashedCard)
	routerCard.POST("/restore/:id", cardHandler.RestoreCard)
	routerCard.DELETE("/permanent/:id", cardHandler.DeleteCardPermanent)

	routerCard.POST("/restore/all", cardHandler.RestoreAllCard)
	routerCard.POST("/permanent/all", cardHandler.DeleteAllCardPermanent)

	return cardHandler
}

// FindAll godoc
// @Summary Retrieve all cards
// @Tags Card
// @Security Bearer
// @Description Retrieve all cards with pagination
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Number of data per page"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponsePaginationCard "Card data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve card data"
// @Router /api/card [get]
func (h *cardHandleApi) FindAll(c echo.Context) error {
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

	req := &pb.FindAllCardRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	cards, err := h.card.FindAllCard(ctx, req)
	if err != nil {
		status = "error"
		logError("Failed to find all cards", err,
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search),
		)
		return card_errors.ErrApiFailedFindAllCards(c)
	}

	response := h.mapping.ToApiResponsesCard(cards)

	logSuccess("Successfully retrieved card list",
		zap.Int("count", len(response.Data)),
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.Bool("success", true),
	)

	return c.JSON(http.StatusOK, response)
}

// FindById godoc
// @Summary Retrieve card by ID
// @Tags Card
// @Security Bearer
// @Description Retrieve a card by its ID
// @Accept json
// @Produce json
// @Param id path int true "Card ID"
// @Success 200 {object} response.ApiResponseCard "Card data"
// @Failure 400 {object} response.ErrorResponse "Invalid card ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve card record"
// @Router /api/card/{id} [get]
func (h *cardHandleApi) FindById(c echo.Context) error {
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

		logError("Invalid card ID", err, zap.Error(err))
		return card_errors.ErrApiInvalidCardID(c)
	}

	req := &pb.FindByIdCardRequest{
		CardId: int32(id),
	}

	card, err := h.card.FindByIdCard(ctx, req)
	if err != nil {
		logError("FindByIdCard failed", err, zap.Int("card.id", id), zap.Error(err))
		return card_errors.ErrApiFailedFindByIdCard(c)
	}

	response := h.mapping.ToApiResponseCard(card)

	logSuccess("Successfully retrieved card record", zap.Int("card.id", id), zap.Bool("success", true))

	return c.JSON(http.StatusOK, response)
}

// FindByUserID godoc
// @Summary Retrieve cards by user ID
// @Tags Card
// @Security Bearer
// @Description Retrieve a list of cards associated with a user by their ID
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseCard "Card data"
// @Failure 400 {object} response.ErrorResponse "Invalid user ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve card record"
// @Router /api/card/user [get]
func (h *cardHandleApi) FindByUserID(c echo.Context) error {
	const method = "FindByUserID"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)
	status := "success"

	defer func() {
		end(status)
	}()

	userIDRaw := c.Get("user_id")
	userID, ok := userIDRaw.(int32)
	if !ok {
		err := errors.New("user id not found in context")

		logError("Invalid user id", err, zap.Error(err))

		return card_errors.ErrApiInvalidUserID(c)
	}

	req := &pb.FindByUserIdCardRequest{
		UserId: userID,
	}

	card, err := h.card.FindByUserIdCard(ctx, req)

	if err != nil {
		logError("FindByUserIdCard failed", err, zap.Error(err))

		return card_errors.ErrApiFailedFindByUserIdCard(c)
	}

	so := h.mapping.ToApiResponseCard(card)

	logSuccess("Success retrieve card record", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// DashboardCard godoc
// @Summary Get dashboard card data
// @Description Retrieve dashboard card data
// @Tags Card
// @Security Bearer
// @Produce json
// @Success 200 {object} response.ApiResponseDashboardCard
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/dashboard [get]
func (h *cardHandleApi) DashboardCard(c echo.Context) error {
	const method = "DashboardCard"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)
	status := "success"

	defer func() {
		end(status)
	}()

	res, err := h.card.DashboardCard(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"

		logError("Failed to retrieve dashboard card data", err, zap.Error(err))

		return card_errors.ErrApiFailedDashboardCard(c)
	}

	so := h.mapping.ToApiResponseDashboardCard(res)

	logSuccess("Success retrieve dashboard card data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// DashboardCardCardNumber godoc
// @Summary Get dashboard card data by card number
// @Description Retrieve dashboard card data for a specific card number
// @Tags Card
// @Security Bearer
// @Produce json
// @Param cardNumber path string true "Card Number"
// @Success 200 {object} response.ApiResponseDashboardCardNumber
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/dashboard/{cardNumber} [get]
func (h *cardHandleApi) DashboardCardCardNumber(c echo.Context) error {
	const method = "DashboardCardCardNumber"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)
	status := "success"

	defer func() {
		end(status)
	}()

	cardNumber := c.Param("cardNumber")

	if cardNumber == "" {
		err := errors.New("card number is empty")

		logError("Invalid card number", err, zap.Error(err))

		return card_errors.ErrApiInvalidCardNumber(c)
	}

	req := &pb.FindByCardNumberRequest{
		CardNumber: cardNumber,
	}

	res, err := h.card.DashboardCardNumber(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to retrieve dashboard card data by card number", err, zap.Error(err))

		return card_errors.ErrApiFailedDashboardCardByCardNumber(c)
	}

	so := h.mapping.ToApiResponseDashboardCardCardNumber(res)

	logSuccess("Success retrieve dashboard card data by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyBalance godoc
// @Summary Get monthly balance data
// @Description Retrieve monthly balance data for a specific year
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMonthlyBalance
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/monthly-balance [get]
func (h *cardHandleApi) FindMonthlyBalance(c echo.Context) error {
	const method = "FindMonthlyBalance"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)
	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		status = "error"

		logError("Invalid year", err, zap.Error(err))

		return card_errors.ErrApiInvalidYear(c)
	}

	req := &pb.FindYearBalance{
		Year: int32(year),
	}

	res, err := h.card.FindMonthlyBalance(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve monthly balance data", err, zap.Error(err))

		return card_errors.ErrApiFailedFindMonthlyBalance(c)
	}

	so := h.mapping.ToApiResponseMonthlyBalances(res)

	logSuccess("Success retrieve monthly balance data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyBalance godoc
// @Summary Get yearly balance data
// @Description Retrieve yearly balance data for a specific year
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseYearlyBalance
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/yearly-balance [get]
func (h *cardHandleApi) FindYearlyBalance(c echo.Context) error {
	const method = "FindYearlyBalance"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)
	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		status = "error"

		logError("Invalid year", err, zap.Error(err))

		return card_errors.ErrApiInvalidYear(c)
	}

	req := &pb.FindYearBalance{
		Year: int32(year),
	}

	res, err := h.card.FindYearlyBalance(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve yearly balance data", err, zap.Error(err))

		return card_errors.ErrApiFailedFindYearlyBalance(c)
	}

	so := h.mapping.ToApiResponseYearlyBalances(res)

	logSuccess("Success retrieve yearly balance data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTopupAmount godoc
// @Summary Get monthly topup amount data
// @Description Retrieve monthly topup amount data for a specific year
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMonthlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/monthly-topup-amount [get]
func (h *cardHandleApi) FindMonthlyTopupAmount(c echo.Context) error {
	const method = "FindMonthlyTopupAmount"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)
	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		status = "error"

		logError("Invalid year", err, zap.Error(err))

		return card_errors.ErrApiInvalidYear(c)
	}

	req := &pb.FindYearAmount{
		Year: int32(year),
	}

	res, err := h.card.FindMonthlyTopupAmount(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve monthly topup amount data", err, zap.Error(err))

		return card_errors.ErrApiFailedFindMonthlyTopupAmount(c)
	}

	so := h.mapping.ToApiResponseMonthlyAmounts(res)

	logSuccess("Success retrieve monthly topup amount data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTopupAmount godoc
// @Summary Get yearly topup amount data
// @Description Retrieve yearly topup amount data for a specific year
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/topup/yearly-topup-amount [get]
func (h *cardHandleApi) FindYearlyTopupAmount(c echo.Context) error {
	const method = "FindYearlyTopupAmount"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)
	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		status = "error"

		logError("Invalid year", err, zap.Error(err))

		return card_errors.ErrApiInvalidYear(c)
	}

	req := &pb.FindYearAmount{
		Year: int32(year),
	}

	res, err := h.card.FindYearlyTopupAmount(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve yearly topup amount data", err, zap.Error(err))

		return card_errors.ErrApiFailedFindYearlyTopupAmount(c)
	}

	so := h.mapping.ToApiResponseYearlyAmounts(res)

	logSuccess("Success retrieve yearly topup amount data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyWithdrawAmount
// godoc
// @Summary Get monthly withdraw amount data
// @Description Retrieve monthly withdraw amount data for a specific year
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMonthlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/monthly-withdraw-amount [get]
func (h *cardHandleApi) FindMonthlyWithdrawAmount(c echo.Context) error {
	const method = "FindMonthlyWithdrawAmount"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)
	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		status = "error"

		logError("Invalid year", err, zap.Error(err))

		return card_errors.ErrApiInvalidYear(c)
	}

	req := &pb.FindYearAmount{
		Year: int32(year),
	}

	res, err := h.card.FindMonthlyWithdrawAmount(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve monthly withdraw amount data", err, zap.Error(err))

		return card_errors.ErrApiFailedFindMonthlyWithdrawAmount(c)
	}

	so := h.mapping.ToApiResponseMonthlyAmounts(res)

	logSuccess("Success retrieve monthly withdraw amount data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyWithdrawAmount godoc
// @Summary Get yearly withdraw amount data
// @Description Retrieve yearly withdraw amount data for a specific year
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/yearly-withdraw-amount [get]
func (h *cardHandleApi) FindYearlyWithdrawAmount(c echo.Context) error {
	const method = "FindYearlyWithdrawAmount"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)
	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		status = "error"

		logError("Invalid year", err, zap.Error(err))

		return card_errors.ErrApiInvalidYear(c)
	}

	req := &pb.FindYearAmount{
		Year: int32(year),
	}

	res, err := h.card.FindYearlyWithdrawAmount(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to retrieve yearly withdraw amount data", err, zap.Error(err))

		return card_errors.ErrApiFailedFindYearlyWithdrawAmount(c)
	}

	so := h.mapping.ToApiResponseYearlyAmounts(res)

	logSuccess("Success retrieve yearly withdraw amount data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTransactionAmount godoc
// @Summary Get monthly transaction amount data
// @Description Retrieve monthly transaction amount data for a specific year
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMonthlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/monthly-transaction-amount [get]
func (h *cardHandleApi) FindMonthlyTransactionAmount(c echo.Context) error {
	const method = "FindMonthlyTransactionAmount"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		status = "error"

		logError("Invalid year", err, zap.Error(err))

		return card_errors.ErrApiInvalidYear(c)
	}

	req := &pb.FindYearAmount{
		Year: int32(year),
	}

	res, err := h.card.FindMonthlyTransactionAmount(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to retrieve monthly transaction amount data", err, zap.Error(err))

		return card_errors.ErrApiFailedFindMonthlyTransactionAmount(c)
	}

	so := h.mapping.ToApiResponseMonthlyAmounts(res)

	logSuccess("Success retrieve monthly transaction amount data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTransactionAmount godoc
// @Summary Get yearly transaction amount data
// @Description Retrieve yearly transaction amount data for a specific year
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/yearly-transaction-amount [get]
func (h *cardHandleApi) FindYearlyTransactionAmount(c echo.Context) error {
	const method = "FindYearlyTransactionAmount"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)
	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		status = "error"

		logError("Invalid year", err, zap.Error(err))

		return card_errors.ErrApiInvalidYear(c)
	}

	req := &pb.FindYearAmount{
		Year: int32(year),
	}

	res, err := h.card.FindYearlyTransactionAmount(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to retrieve yearly transaction amount data", err, zap.Error(err))

		return card_errors.ErrApiFailedFindYearlyTransactionAmount(c)
	}

	so := h.mapping.ToApiResponseYearlyAmounts(res)

	logSuccess("Success retrieve yearly transaction amount data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTransferSenderAmount godoc
// @Summary Get monthly transfer sender amount data
// @Description Retrieve monthly transfer sender amount data for a specific year
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMonthlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/monthly-transfer-sender-amount [get]
func (h *cardHandleApi) FindMonthlyTransferSenderAmount(c echo.Context) error {
	const method = "FindMonthlyTransferSenderAmount"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)
	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		status = "error"

		logError("Invalid year", err, zap.Error(err))

		return card_errors.ErrApiInvalidYear(c)
	}

	req := &pb.FindYearAmount{
		Year: int32(year),
	}

	res, err := h.card.FindMonthlyTransferSenderAmount(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to retrieve monthly transfer sender amount data", err, zap.Error(err))

		return card_errors.ErrApiFailedFindMonthlyTransferSenderAmount(c)
	}

	so := h.mapping.ToApiResponseMonthlyAmounts(res)

	logSuccess("Success retrieve monthly transfer sender amount data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTransferSenderAmount godoc
// @Summary Get yearly transfer sender amount data
// @Description Retrieve yearly transfer sender amount data for a specific year
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/transfer/yearly-transfer-sender-amount [get]
func (h *cardHandleApi) FindYearlyTransferSenderAmount(c echo.Context) error {
	const method = "FindYearlyTransferSenderAmount"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)
	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		status = "error"

		logError("Invalid year", err, zap.Error(err))

		return card_errors.ErrApiInvalidYear(c)
	}

	req := &pb.FindYearAmount{
		Year: int32(year),
	}

	res, err := h.card.FindYearlyTransferSenderAmount(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to retrieve yearly transfer sender amount data", err, zap.Error(err))

		return card_errors.ErrApiFailedFindYearlyTransferSenderAmount(c)
	}

	so := h.mapping.ToApiResponseYearlyAmounts(res)

	logSuccess("Success retrieve yearly transfer sender amount data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTransferReceiverAmount godoc
// @Summary Get monthly transfer receiver amount data
// @Description Retrieve monthly transfer receiver amount data for a specific year
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseMonthlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/monthly-transfer-receiver-amount [get]
func (h *cardHandleApi) FindMonthlyTransferReceiverAmount(c echo.Context) error {
	const method = "FindMonthlyTransferReceiverAmount"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		status = "error"

		logError("Invalid year", err, zap.Error(err))

		return card_errors.ErrApiInvalidYear(c)
	}

	req := &pb.FindYearAmount{
		Year: int32(year),
	}

	res, err := h.card.FindMonthlyTransferReceiverAmount(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to retrieve monthly transfer receiver amount data", err, zap.Error(err))

		return card_errors.ErrApiFailedFindMonthlyTransferReceiverAmount(c)
	}

	so := h.mapping.ToApiResponseMonthlyAmounts(res)

	logSuccess("Success retrieve monthly transfer receiver amount data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTransferReceiverAmount godoc
// @Summary Get yearly transfer receiver amount data
// @Description Retrieve yearly transfer receiver amount data for a specific year
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/yearly-transfer-receiver-amount [get]
func (h *cardHandleApi) FindYearlyTransferReceiverAmount(c echo.Context) error {
	const method = "FindYearlyTransferReceiverAmount"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		status = "error"

		logError("Invalid year", err, zap.Error(err))

		return card_errors.ErrApiInvalidYear(c)
	}

	req := &pb.FindYearAmount{
		Year: int32(year),
	}

	res, err := h.card.FindYearlyTransferReceiverAmount(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to retrieve yearly transfer receiver amount data", err, zap.Error(err))

		return card_errors.ErrApiFailedFindYearlyTransferReceiverAmount(c)
	}

	so := h.mapping.ToApiResponseYearlyAmounts(res)

	logSuccess("Success retrieve yearly transfer receiver amount data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyBalanceByCardNumber godoc
// @Summary Get monthly balance data by card number
// @Description Retrieve monthly balance data for a specific year and card number
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseMonthlyBalance
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/monthly-balance-by-card [get]
func (h *cardHandleApi) FindMonthlyBalanceByCardNumber(c echo.Context) error {
	const method = "FindMonthlyBalanceByCardNumber"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		status = "error"

		logError("Invalid year", err, zap.Error(err))

		return card_errors.ErrApiInvalidYear(c)
	}

	cardNumber := c.QueryParam("card_number")
	if cardNumber == "" {
		status = "error"
		err := errors.New("invalid card number")

		logError("Invalid card number", err, zap.Error(err))

		return card_errors.ErrApiInvalidCardNumber(c)
	}

	req := &pb.FindYearBalanceCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindMonthlyBalanceByCardNumber(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to retrieve monthly balance data by card number", err, zap.Error(err))

		return card_errors.ErrApiFailedFindMonthlyBalanceByCard(c)
	}

	so := h.mapping.ToApiResponseMonthlyBalances(res)

	logSuccess("Success retrieve monthly balance data by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyBalanceByCardNumber godoc
// @Summary Get yearly balance data by card number
// @Description Retrieve yearly balance data for a specific year and card number
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseYearlyBalance
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/yearly-balance-by-card [get]
func (h *cardHandleApi) FindYearlyBalanceByCardNumber(c echo.Context) error {
	const method = "FindYearlyBalanceByCardNumber"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		status = "error"

		logError("Invalid year", err, zap.Error(err))

		return card_errors.ErrApiInvalidYear(c)
	}

	cardNumber := c.QueryParam("card_number")
	if cardNumber == "" {
		status = "error"

		logError("Invalid card number", err, zap.Error(err))

		return card_errors.ErrApiInvalidCardNumber(c)
	}

	req := &pb.FindYearBalanceCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindYearlyBalanceByCardNumber(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to retrieve yearly balance data by card number", err, zap.Error(err))

		return card_errors.ErrApiFailedFindYearlyBalanceByCard(c)
	}

	so := h.mapping.ToApiResponseYearlyBalances(res)

	logSuccess("Success retrieve yearly balance data by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTopupAmountByCardNumber godoc
// @Summary Get monthly topup amount data by card number
// @Description Retrieve monthly topup amount data for a specific year and card number
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseMonthlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/monthly-topup-amount-by-card [get]
func (h *cardHandleApi) FindMonthlyTopupAmountByCardNumber(c echo.Context) error {
	const method = "FindMonthlyTopupAmountByCardNumber"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		status = "error"

		logError("Invalid year", err, zap.Error(err))

		return card_errors.ErrApiInvalidYear(c)
	}

	cardNumber := c.QueryParam("card_number")
	if cardNumber == "" {
		status = "error"
		err := errors.New("invalid card number")

		logError("Invalid card number", err, zap.Error(err))

		return card_errors.ErrApiInvalidCardNumber(c)
	}

	req := &pb.FindYearAmountCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindMonthlyTopupAmountByCardNumber(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to retrieve monthly topup amount data by card number", err, zap.Error(err))

		return card_errors.ErrApiFailedFindMonthlyTopupAmountByCard(c)
	}

	so := h.mapping.ToApiResponseMonthlyAmounts(res)

	logSuccess("Success retrieve monthly topup amount data by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTopupAmountByCardNumber godoc
// @Summary Get yearly topup amount data by card number
// @Description Retrieve yearly topup amount data for a specific year and card number
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/yearly-topup-amount-by-card [get]
func (h *cardHandleApi) FindYearlyTopupAmountByCardNumber(c echo.Context) error {
	const method = "FindYearlyTopupAmountByCardNumber"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		status = "error"

		logError("Invalid year", err, zap.Error(err))

		return card_errors.ErrApiInvalidYear(c)
	}

	cardNumber := c.QueryParam("card_number")

	if cardNumber == "" {
		status = "error"
		err := errors.New("invalid card number")

		logError("Invalid card number", err, zap.Error(err))

		return card_errors.ErrApiInvalidCardNumber(c)
	}

	req := &pb.FindYearAmountCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindYearlyTopupAmountByCardNumber(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to retrieve yearly topup amount data by card number", err, zap.Error(err))

		return card_errors.ErrApiFailedFindYearlyTopupAmountByCard(c)
	}

	so := h.mapping.ToApiResponseYearlyAmounts(res)

	logSuccess("Success retrieve yearly topup amount data by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyWithdrawAmountByCardNumber godoc
// @Summary Get monthly withdraw amount data by card number
// @Description Retrieve monthly withdraw amount data for a specific year and card number
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseMonthlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/monthly-withdraw-amount-by-card [get]
func (h *cardHandleApi) FindMonthlyWithdrawAmountByCardNumber(c echo.Context) error {
	const method = "FindMonthlyWithdrawAmountByCardNumber"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)

	if err != nil {
		status = "error"

		logError("Invalid year", err, zap.Error(err))

		return card_errors.ErrApiInvalidYear(c)
	}

	cardNumber := c.QueryParam("card_number")
	if cardNumber == "" {
		status = "error"
		err := errors.New("invalid card number")

		logError("Invalid card number", err, zap.Error(err))

		return card_errors.ErrApiInvalidCardNumber(c)
	}

	req := &pb.FindYearAmountCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindMonthlyWithdrawAmountByCardNumber(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve monthly withdraw amount data by card number", err, zap.Error(err))

		return card_errors.ErrApiFailedFindMonthlyWithdrawAmountByCard(c)
	}

	so := h.mapping.ToApiResponseMonthlyAmounts(res)

	logSuccess("Success retrieve monthly withdraw amount data by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyWithdrawAmountByCardNumber godoc
// @Summary Get yearly withdraw amount data by card number
// @Description Retrieve yearly withdraw amount data for a specific year and card number
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/yearly-withdraw-amount-by-card [get]
func (h *cardHandleApi) FindYearlyWithdrawAmountByCardNumber(c echo.Context) error {
	const method = "FindYearlyWithdrawAmountByCardNumber"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		status = "error"

		logError("Invalid year", err, zap.Error(err))

		return card_errors.ErrApiInvalidYear(c)
	}

	cardNumber := c.QueryParam("card_number")

	if cardNumber == "" {
		status = "error"

		logError("Invalid card number", err, zap.Error(err))

		return card_errors.ErrApiInvalidCardNumber(c)
	}

	req := &pb.FindYearAmountCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindYearlyWithdrawAmountByCardNumber(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve yearly withdraw amount data by card number", err, zap.Error(err))

		return card_errors.ErrApiFailedFindYearlyWithdrawAmountByCard(c)
	}

	so := h.mapping.ToApiResponseYearlyAmounts(res)

	logSuccess("Success retrieve yearly withdraw amount data by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTransactionAmountByCardNumber godoc
// @Summary Get monthly transaction amount data by card number
// @Description Retrieve monthly transaction amount data for a specific year and card number
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseMonthlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/monthly-transaction-amount-by-card [get]
func (h *cardHandleApi) FindMonthlyTransactionAmountByCardNumber(c echo.Context) error {
	const method = "FindMonthlyTransactionAmountByCardNumber"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		status = "error"

		logError("Invalid year", err, zap.Error(err))

		return card_errors.ErrApiInvalidYear(c)
	}

	cardNumber := c.QueryParam("card_number")

	if cardNumber == "" {
		status = "error"

		logError("Invalid card number", err, zap.Error(err))

		return card_errors.ErrApiInvalidCardNumber(c)
	}

	req := &pb.FindYearAmountCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindMonthlyTransactionAmountByCardNumber(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to retrieve monthly transaction amount data by card number", err, zap.Error(err))

		return card_errors.ErrApiFailedFindMonthlyTransactionAmountByCard(c)
	}

	so := h.mapping.ToApiResponseMonthlyAmounts(res)

	logSuccess("Success retrieve monthly transaction amount data by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTransactionAmountByCardNumber godoc
// @Summary Get yearly transaction amount data by card number
// @Description Retrieve yearly transaction amount data for a specific year and card number
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/yearly-transaction-amount-by-card [get]
func (h *cardHandleApi) FindYearlyTransactionAmountByCardNumber(c echo.Context) error {
	const method = "FindYearlyTransactionAmountByCardNumber"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		status = "error"

		logError("Invalid year", err, zap.Error(err))

		return card_errors.ErrApiInvalidYear(c)
	}

	cardNumber := c.QueryParam("card_number")

	if cardNumber == "" {
		status = "error"

		logError("Invalid card number", err, zap.Error(err))

		return card_errors.ErrApiInvalidCardNumber(c)
	}

	req := &pb.FindYearAmountCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindYearlyTransactionAmountByCardNumber(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to retrieve yearly transaction amount data by card number", err, zap.Error(err))

		return card_errors.ErrApiFailedFindYearlyTransactionAmountByCard(c)
	}

	so := h.mapping.ToApiResponseYearlyAmounts(res)

	logSuccess("Success retrieve yearly transaction amount data by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTransferSenderAmountByCardNumber godoc
// @Summary Get monthly transfer sender amount data by card number
// @Description Retrieve monthly transfer sender amount data for a specific year and card number
// @Tags Card
// @Security Bearer
// @Produce json
// @Param year query int true "Year"
// @Param card_number query string true "Card Number"
// @Success 200 {object} response.ApiResponseMonthlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/monthly-transfer-sender-amount-by-card [get]
func (h *cardHandleApi) FindMonthlyTransferSenderAmountByCardNumber(c echo.Context) error {
	const method = "FindMonthlyTransferSenderAmountByCardNumber"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		status = "error"

		logError("Invalid year", err, zap.Error(err))

		return card_errors.ErrApiInvalidYear(c)
	}

	cardNumber := c.QueryParam("card_number")
	if cardNumber == "" {
		status = "error"

		logError("Invalid card number", err, zap.Error(err))

		return card_errors.ErrApiInvalidCardNumber(c)
	}

	req := &pb.FindYearAmountCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindMonthlyTransferSenderAmountByCardNumber(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to retrieve monthly transfer sender amount data by card number", err, zap.Error(err))

		return card_errors.ErrApiFailedFindMonthlyTransferSenderAmountByCard(c)
	}

	so := h.mapping.ToApiResponseMonthlyAmounts(res)

	logSuccess("Success retrieve monthly transfer sender amount data by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTransferSenderAmountByCardNumber godoc
// @Summary Get yearly transfer sender amount by card number
// @Description Retrieve the total amount sent by a specific card number in a given year
// @Tags Card
// @Security Bearer
// @Accept json
// @Produce json
// @Param year query int true "Year for which the data is requested"
// @Param card_number query string true "Card number for which the data is requested"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/yearly-transfer-sender-amount-by-card [get]
func (h *cardHandleApi) FindYearlyTransferSenderAmountByCardNumber(c echo.Context) error {
	const method = "FindYearlyTransferSenderAmountByCardNumber"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		status = "error"

		logError("Invalid year", err, zap.Error(err))

		return card_errors.ErrApiInvalidYear(c)
	}

	cardNumber := c.QueryParam("card_number")
	if cardNumber == "" {
		status = "error"

		logError("Invalid card number", err, zap.Error(err))

		return card_errors.ErrApiInvalidCardNumber(c)
	}

	req := &pb.FindYearAmountCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindYearlyTransferSenderAmountByCardNumber(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to retrieve yearly transfer sender amount data by card number", err, zap.Error(err))

		return card_errors.ErrApiFailedFindYearlyTransferSenderAmountByCard(c)
	}

	so := h.mapping.ToApiResponseYearlyAmounts(res)

	logSuccess("Success retrieve yearly transfer sender amount data by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTransferReceiverAmountByCardNumber godoc
// @Summary Get monthly transfer receiver amount by card number
// @Description Retrieve the total amount received by a specific card number in a given year, broken down by month
// @Tags Card
// @Security Bearer
// @Accept json
// @Produce json
// @Param year query int true "Year for which the data is requested"
// @Param card_number query string true "Card number for which the data is requested"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/monthly-transfer-receiver-amount-by-card [get]
func (h *cardHandleApi) FindMonthlyTransferReceiverAmountByCardNumber(c echo.Context) error {
	const method = "FindMonthlyTransferReceiverAmountByCardNumber"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		status = "error"

		logError("Invalid year", err, zap.Error(err))

		return card_errors.ErrApiInvalidYear(c)
	}

	cardNumber := c.QueryParam("card_number")

	if cardNumber == "" {
		status = "error"

		logError("Invalid card number", err, zap.Error(err))

		return card_errors.ErrApiInvalidCardNumber(c)
	}

	req := &pb.FindYearAmountCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindMonthlyTransferReceiverAmountByCardNumber(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to retrieve monthly transfer receiver amount data by card number", err, zap.Error(err))

		return card_errors.ErrApiFailedFindMonthlyTransferReceiverAmountByCard(c)
	}

	so := h.mapping.ToApiResponseMonthlyAmounts(res)

	logSuccess("Success retrieve monthly transfer receiver amount data by card number", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearlyTransferReceiverAmountByCardNumber godoc
// @Summary Get yearly transfer receiver amount by card number
// @Description Retrieve the total amount received by a specific card number in a given year
// @Tags Card
// @Security Bearer
// @Accept json
// @Produce json
// @Param year query int true "Year for which the data is requested"
// @Param card_number query string true "Card number for which the data is requested"
// @Success 200 {object} response.ApiResponseYearlyAmount
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/card/yearly-transfer-receiver-amount-by-card [get]
func (h *cardHandleApi) FindYearlyTransferReceiverAmountByCardNumber(c echo.Context) error {
	const method = "FindYearlyTransferReceiverAmountByCardNumber"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		status = "error"

		logError("Invalid year", err, zap.Error(err))

		return card_errors.ErrApiInvalidYear(c)
	}

	cardNumber := c.QueryParam("card_number")

	if cardNumber == "" {
		status = "error"

		logError("Invalid card number", err, zap.Error(err))

		return card_errors.ErrApiInvalidCardNumber(c)
	}

	req := &pb.FindYearAmountCardNumber{
		Year:       int32(year),
		CardNumber: cardNumber,
	}

	res, err := h.card.FindYearlyTransferReceiverAmountByCardNumber(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to retrieve yearly transfer receiver amount data by card number", err, zap.Error(err))

		return card_errors.ErrApiFailedFindYearlyTransferReceiverAmountByCard(c)
	}

	logSuccess("Success retrieve yearly transfer receiver amount data by card number", zap.Bool("success", true))

	so := h.mapping.ToApiResponseYearlyAmounts(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active card by Saldo ID
// @Tags Card
// @Description Retrieve an active card associated with a Saldo ID
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponsePaginationCardDeleteAt "Card data"
// @Failure 400 {object} response.ErrorResponse "Invalid Saldo ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve card record"
// @Router /api/card/active [get]
func (h *cardHandleApi) FindByActive(c echo.Context) error {
	const method = "FindByActive"

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

	req := &pb.FindAllCardRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.card.FindByActiveCard(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve card record", err, zap.Error(err))

		return card_errors.ErrApiFailedFindByActiveCard(c)
	}

	so := h.mapping.ToApiResponsesCardDeletedAt(res)

	logSuccess("Success retrieve card record", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Retrieve trashed cards
// @Tags Card
// @Security Bearer
// @Description Retrieve a list of trashed cards
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Number of items per page (default: 10)"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponsePaginationCardDeleteAt "Card data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve card record"
// @Router /api/card/trashed [get]
func (h *cardHandleApi) FindByTrashed(c echo.Context) error {
	const method = "FindByTrashed"

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

	req := &pb.FindAllCardRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.card.FindByTrashedCard(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve card record", err, zap.Error(err))

		return card_errors.ErrApiFailedFindByTrashedCard(c)
	}

	so := h.mapping.ToApiResponsesCardDeletedAt(res)

	logSuccess("Success retrieve card record", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve card by card number
// @Tags Card
// @Description Retrieve a card by its card number
// @Accept json
// @Produce json
// @Param card_number path string true "Card number"
// @Success 200 {object} response.ApiResponseCard "Card data"
// @Failure 400 {object} response.ErrorResponse "Failed to fetch card record"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve card record"
// @Router /api/card/{card_number} [get]
func (h *cardHandleApi) FindByCardNumber(c echo.Context) error {
	const method = "FindByCardNumber"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	cardNumber := c.Param("card_number")

	if cardNumber == "" {
		status = "error"
		err := errors.New("invalid card number")

		logError("Failed to fetch card record", err, zap.Error(err))

		return card_errors.ErrApiInvalidCardNumber(c)
	}

	req := &pb.FindByCardNumberRequest{
		CardNumber: cardNumber,
	}

	res, err := h.card.FindByCardNumber(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve card record", err, zap.Error(err))

		return card_errors.ErrApiFailedFindByCardNumber(c)
	}

	so := h.mapping.ToApiResponseCard(res)

	logSuccess("Success retrieve card record", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Create a new card
// @Tags Card
// @Description Create a new card for a user
// @Accept json
// @Produce json
// @Param CreateCardRequest body requests.CreateCardRequest true "Create card request"
// @Success 200 {object} response.ApiResponseCard "Created card"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create card"
// @Router /api/card/create [post]
func (h *cardHandleApi) CreateCard(c echo.Context) error {
	const method = "CreateCard"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)
	status := "success"
	defer func() {
		end(status)
	}()

	var body requests.CreateCardRequest

	if err := c.Bind(&body); err != nil {
		status = "error"
		logError("Failed to bind CreateCard request", err, zap.Error(err))
		return card_errors.ErrApiBindCreateCard(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"
		logError("Validation failed for CreateCard", err, zap.Error(err))
		return card_errors.ErrApiValidateCreateCard(c)
	}

	req := &pb.CreateCardRequest{
		UserId:       int32(body.UserID),
		CardType:     body.CardType,
		ExpireDate:   timestamppb.New(body.ExpireDate),
		Cvv:          body.CVV,
		CardProvider: body.CardProvider,
	}

	res, err := h.card.CreateCard(ctx, req)
	if err != nil {
		status = "error"
		logError("CreateCard service failed", err)
		return card_errors.ErrApiFailedCreateCard(c)
	}

	response := h.mapping.ToApiResponseCard(res)
	logSuccess("Successfully created card", zap.Int("card_id", int(response.Data.ID)))

	return c.JSON(http.StatusOK, response)
}

// @Security Bearer
// @Summary Update a card
// @Tags Card
// @Description Update a card for a user
// @Accept json
// @Produce json
// @Param id path int true "Card ID"
// @Param UpdateCardRequest body requests.UpdateCardRequest true "Update card request"
// @Success 200 {object} response.ApiResponseCard "Updated card"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update card"
// @Router /api/card/update/{id} [post]
func (h *cardHandleApi) UpdateCard(c echo.Context) error {
	const method = "UpdateCard"

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

		logError("Invalid card id", err, zap.Error(err))

		return card_errors.ErrApiInvalidCardID(c)
	}

	var body requests.UpdateCardRequest

	if err := c.Bind(&body); err != nil {
		status = "error"

		logError("Failed to bind UpdateCard request", err, zap.Error(err))

		return card_errors.ErrApiBindUpdateCard(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"

		logError("Validation failed for UpdateCard", err, zap.Error(err))

		return card_errors.ErrApiValidateUpdateCard(c)
	}

	req := &pb.UpdateCardRequest{
		CardId:       int32(idInt),
		UserId:       int32(body.UserID),
		CardType:     body.CardType,
		ExpireDate:   timestamppb.New(body.ExpireDate),
		Cvv:          body.CVV,
		CardProvider: body.CardProvider,
	}

	res, err := h.card.UpdateCard(ctx, req)

	if err != nil {
		status = "error"

		logError("UpdateCard service failed", err, zap.Error(err))

		return card_errors.ErrApiFailedUpdateCard(c)
	}

	so := h.mapping.ToApiResponseCard(res)

	logSuccess("Successfully updated card", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Trashed a card
// @Tags Card
// @Description Trashed a card by its ID
// @Accept json
// @Produce json
// @Param id path int true "Card ID"
// @Success 200 {object} response.ApiResponseCard "Trashed card"
// @Failure 400 {object} response.ErrorResponse "Bad request or invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to trashed card"
// @Router /api/card/trashed/{id} [post]
func (h *cardHandleApi) TrashedCard(c echo.Context) error {
	const method = "TrashedCard"

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

		logError("Invalid card id", err, zap.Error(err))

		return card_errors.ErrApiInvalidCardID(c)
	}

	req := &pb.FindByIdCardRequest{
		CardId: int32(idInt),
	}

	res, err := h.card.TrashedCard(ctx, req)

	if err != nil {
		status = "error"

		logError("TrashedCard service failed", err, zap.Error(err))

		return card_errors.ErrApiFailedTrashCard(c)
	}

	so := h.mapping.ToApiResponseCard(res)

	logSuccess("Successfully trashed card", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Restore a card
// @Tags Card
// @Description Restore a card by its ID
// @Accept json
// @Produce json
// @Param id path int true "Card ID"
// @Success 200 {object} response.ApiResponseCard "Restored card"
// @Failure 400 {object} response.ErrorResponse "Bad request or invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore card"
// @Router /api/card/restore/{id} [post]
func (h *cardHandleApi) RestoreCard(c echo.Context) error {
	const method = "RestoreCard"

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

		logError("Invalid card id", err, zap.Error(err))

		return card_errors.ErrApiInvalidCardID(c)
	}

	req := &pb.FindByIdCardRequest{
		CardId: int32(idInt),
	}

	res, err := h.card.RestoreCard(ctx, req)

	if err != nil {
		status = "error"
		msg := "Failed to restore card"

		logError(msg, err, zap.Error(err))

		return card_errors.ErrApiFailedRestoreCard(c)
	}

	so := h.mapping.ToApiResponseCard(res)

	logSuccess("Successfully restored card", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Delete a card permanently
// @Tags Card
// @Description Delete a card by its ID permanently
// @Accept json
// @Produce json
// @Param id path int true "Card ID"
// @Success 200 {object} response.ApiResponseCardDelete "Deleted card"
// @Failure 400 {object} response.ErrorResponse "Bad request or invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete card"
// @Router /api/card/permanent/{id} [delete]
func (h *cardHandleApi) DeleteCardPermanent(c echo.Context) error {
	const method = "DeleteCardPermanent"

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

		logError("Invalid card id", err, zap.Error(err))

		return card_errors.ErrApiInvalidCardID(c)
	}

	req := &pb.FindByIdCardRequest{
		CardId: int32(idInt),
	}

	res, err := h.card.DeleteCardPermanent(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to delete card permanently", err, zap.Error(err))

		return card_errors.ErrApiFailedDeleteCardPermanent(c)
	}

	so := h.mapping.ToApiResponseCardDeleteAt(res)

	logSuccess("Successfully deleted card permanently", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Restore all card records
// @Tags Card
// @Description Restore all card records that were previously deleted.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseCardAll "Successfully restored all card records"
// @Failure 500 {object} response.ErrorResponse "Failed to restore all card records"
// @Router /api/card/restore/all [post]
func (h *cardHandleApi) RestoreAllCard(c echo.Context) error {
	const method = "RestoreAllCard"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	res, err := h.card.RestoreAllCard(ctx, &emptypb.Empty{})
	if err != nil {
		status = "error"

		logError("Failed to restore all cards", err, zap.Error(err))

		return card_errors.ErrApiFailedRestoreAllCard(c)
	}

	h.logger.Debug("Successfully restored all cards")

	so := h.mapping.ToApiResponseCardAll(res)

	logSuccess("Successfully restored all cards", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer.
// @Summary Permanently delete all card records
// @Tags Card
// @Description Permanently delete all card records from the database.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseCardAll "Successfully deleted all card records permanently"
// @Failure 500 {object} response.ErrorResponse "Failed to permanently delete all card records"
// @Router /api/card/permanent/all [post]
func (h *cardHandleApi) DeleteAllCardPermanent(c echo.Context) error {
	const method = "DeleteAllCardPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() {
		end(status)
	}()

	res, err := h.card.DeleteAllCardPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"

		logError("Failed to delete all cards permanently", err, zap.Error(err))

		return card_errors.ErrApiFailedDeleteAllCardPermanent(c)
	}

	so := h.mapping.ToApiResponseCardAll(res)

	logSuccess("Successfully deleted all cards permanently", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *cardHandleApi) startTracingAndLogging(
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

func (s *cardHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
