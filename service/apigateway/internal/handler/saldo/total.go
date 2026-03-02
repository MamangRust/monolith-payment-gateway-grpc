package saldohandler

import (
	"net/http"
	"strconv"

	saldo_cache "github.com/MamangRust/monolith-payment-gateway-apigateway/internal/redis/api/saldo"
	pbsaldo "github.com/MamangRust/monolith-payment-gateway-pb/saldo"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/saldo/stats"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/saldo"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type saldoTotalBalanceHandleApi struct {
	saldo pb.SaldoStatsTotalBalanceClient

	logger logger.LoggerInterface

	mapper apimapper.SaldoStatsTotalResponseMapper

	cache saldo_cache.SaldoMencache

	apiHandler errors.ApiHandler
}

type saldoTotalBalanceHandleDeps struct {
	client pb.SaldoStatsTotalBalanceClient

	router *echo.Echo

	logger logger.LoggerInterface

	mapper apimapper.SaldoStatsTotalResponseMapper

	cache saldo_cache.SaldoMencache

	apiHandler errors.ApiHandler
}

func NewSaldoTotalBalanceHandleApi(params *saldoTotalBalanceHandleDeps) *saldoTotalBalanceHandleApi {
	saldoHandler := &saldoTotalBalanceHandleApi{
		saldo:      params.client,
		logger:     params.logger,
		mapper:     params.mapper,
		cache:      params.cache,
		apiHandler: params.apiHandler,
	}

	routerSaldo := params.router.Group("/api/saldo-stats-total-balance")

	routerSaldo.GET("/monthly-total-balance", params.apiHandler.Handle("find-monthly-total-saldo-balance", saldoHandler.FindMonthlyTotalSaldoBalance))
	routerSaldo.GET("/yearly-total-balance", params.apiHandler.Handle("find-yearly-total-saldo-balance", saldoHandler.FindYearTotalSaldoBalance))

	return saldoHandler
}

// FindMonthlyTotalSaldoBalance retrieves the total saldo balance for a specific month and year.
// @Summary Get monthly total saldo balance
// @Tags Saldo
// @Security Bearer
// @Description Retrieve the total saldo balance for a specific month and year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseMonthTotalSaldo "Monthly total saldo balance"
// @Failure 400 {object} response.ErrorResponse "Invalid year or month parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve monthly total saldo balance"
// @Router /api/saldo-stats-total-balances/monthly-total-balance [get]
func (h *saldoTotalBalanceHandleApi) FindMonthlyTotalSaldoBalance(c echo.Context) error {
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month <= 0 || month > 12 {
		return errors.NewBadRequestError("invalid month parameter")
	}

	ctx := c.Request().Context()

	reqCache := &requests.MonthTotalSaldoBalance{
		Year:  year,
		Month: month,
	}

	cachedData, found := h.cache.GetMonthlyTotalSaldoBalanceCache(ctx, reqCache)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.saldo.FindMonthlyTotalSaldoBalance(ctx, &pbsaldo.FindMonthlySaldoTotalBalance{
		Year:  int32(year),
		Month: int32(month),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly total saldo balance", zap.Error(err))
		return h.handleGrpcError(err, "FindMonthlyTotalSaldoBalance")
	}

	apiResponse := h.mapper.ToApiResponseMonthTotalSaldo(res)
	h.cache.SetMonthlyTotalSaldoCache(ctx, reqCache, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// FindYearTotalSaldoBalance retrieves the total saldo balance for a specific year.
// @Summary Get yearly total saldo balance
// @Tags Saldo
// @Security Bearer
// @Description Retrieve the total saldo balance for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Success 200 {object} response.ApiResponseYearTotalSaldo "Yearly total saldo balance"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly total saldo balance"
// @Router /api/saldo-stats-total-balance/yearly-total-balance [get]
func (h *saldoTotalBalanceHandleApi) FindYearTotalSaldoBalance(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year <= 0 {
		return errors.NewBadRequestError("invalid year parameter")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetYearTotalSaldoBalanceCache(ctx, year)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.saldo.FindYearTotalSaldoBalance(ctx, &pbsaldo.FindYearlySaldo{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve year total saldo balance", zap.Error(err))
		return h.handleGrpcError(err, "FindYearTotalSaldoBalance")
	}

	apiResponse := h.mapper.ToApiResponseYearTotalSaldo(res)
	h.cache.SetYearTotalSaldoBalanceCache(ctx, year, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

func (h *saldoTotalBalanceHandleApi) handleGrpcError(err error, operation string) *errors.AppError {
	st, ok := status.FromError(err)
	if !ok {
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}

	switch st.Code() {
	case codes.NotFound:
		return errors.NewNotFoundError("Saldo").WithInternal(err)

	case codes.AlreadyExists:
		return errors.NewConflictError("Saldo already exists").WithInternal(err)

	case codes.InvalidArgument:
		return errors.NewBadRequestError(st.Message()).WithInternal(err)

	case codes.PermissionDenied:
		return errors.ErrForbidden.WithInternal(err)

	case codes.Unauthenticated:
		return errors.ErrUnauthorized.WithInternal(err)

	case codes.ResourceExhausted:
		return errors.ErrTooManyRequests.WithInternal(err)

	case codes.Unavailable:
		return errors.NewServiceUnavailableError("Saldo service").WithInternal(err)

	case codes.DeadlineExceeded:
		return errors.ErrTimeout.WithInternal(err)

	default:
		return errors.NewInternalError(err).WithMessage("Failed to " + operation)
	}
}
