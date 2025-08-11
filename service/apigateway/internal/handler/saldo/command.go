package saldohandler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	pb "github.com/MamangRust/monolith-payment-gateway-pb/saldo"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/api"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api/saldo"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type saldoCommandHandleApi struct {
	saldo pb.SaldoCommandServiceClient

	logger logger.LoggerInterface

	mapper apimapper.SaldoCommandResponseMapper

	trace trace.Tracer

	requestCounter *prometheus.CounterVec

	requestDuration *prometheus.HistogramVec
}

type saldoCommandHandleDeps struct {
	client pb.SaldoCommandServiceClient

	router *echo.Echo

	logger logger.LoggerInterface

	mapper apimapper.SaldoCommandResponseMapper
}

func NewSaldoCommandHandleApi(params *saldoCommandHandleDeps) *saldoCommandHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "saldo_command_handler_requests_total",
			Help: "Total number of saldo command requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "saldo_command_handler_request_duration_seconds",
			Help:    "Duration of saldo command requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	saldoHandler := &saldoCommandHandleApi{
		saldo:           params.client,
		logger:          params.logger,
		mapper:          params.mapper,
		trace:           otel.Tracer("saldo-command-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerSaldo := params.router.Group("/api/saldo-command")

	routerSaldo.POST("/create", saldoHandler.Create)
	routerSaldo.POST("/update/:id", saldoHandler.Update)
	routerSaldo.POST("/trashed/:id", saldoHandler.TrashSaldo)
	routerSaldo.POST("/restore/:id", saldoHandler.RestoreSaldo)
	routerSaldo.DELETE("/permanent/:id", saldoHandler.Delete)

	routerSaldo.POST("/restore/all", saldoHandler.RestoreAllSaldo)
	routerSaldo.POST("/permanent/all", saldoHandler.DeleteAllSaldoPermanent)

	return saldoHandler
}

// @Summary Create a new saldo
// @Tags Saldo-Command
// @Security Bearer
// @Description Create a new saldo record with the provided card number and total balance.
// @Accept json
// @Produce json
// @Param CreateSaldoRequest body requests.CreateSaldoRequest true "Create Saldo Request"
// @Success 200 {object} response.ApiResponseSaldo "Successfully created saldo record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create saldo"
// @Router /api/saldo-command/create [post]
func (h *saldoCommandHandleApi) Create(c echo.Context) error {
	const method = "Create"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	var body requests.CreateSaldoRequest

	if err := c.Bind(&body); err != nil {
		logError("Failed to bind CreateSaldo request", err, zap.Error(err))

		return saldo_errors.ErrApiBindCreateSaldo(c)
	}

	if err := body.Validate(); err != nil {
		logError("Failed to validate CreateSaldo request", err, zap.Error(err))

		return saldo_errors.ErrApiValidateCreateSaldo(c)
	}

	res, err := h.saldo.CreateSaldo(ctx, &pb.CreateSaldoRequest{
		CardNumber:   body.CardNumber,
		TotalBalance: int32(body.TotalBalance),
	})

	if err != nil {
		logError("Failed to create saldo", err, zap.Error(err))

		return saldo_errors.ErrApiFailedCreateSaldo(c)
	}

	so := h.mapper.ToApiResponseSaldo(res)

	logSuccess("Successfully create saldo", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Update an existing saldo
// @Tags Saldo-Command
// @Security Bearer
// @Description Update an existing saldo record with the provided card number and total balance.
// @Accept json
// @Produce json
// @Param id path int true "Saldo ID"
// @Param UpdateSaldoRequest body requests.UpdateSaldoRequest true "Update Saldo Request"
// @Success 200 {object} response.ApiResponseSaldo "Successfully updated saldo record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update saldo"
// @Router /api/saldo-command/update/{id} [post]
func (h *saldoCommandHandleApi) Update(c echo.Context) error {
	const method = "Update"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	idint, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logError("Bad Request: Invalid ID", err, zap.Error(err))

		return saldo_errors.ErrApiInvalidSaldoID(c)
	}

	var body requests.UpdateSaldoRequest

	if err := c.Bind(&body); err != nil {
		logError("Failed to bind UpdateSaldo request", err, zap.Error(err))

		return saldo_errors.ErrApiBindUpdateSaldo(c)
	}

	if err := body.Validate(); err != nil {
		logError("Failed to validate UpdateSaldo request", err, zap.Error(err))

		return saldo_errors.ErrApiValidateUpdateSaldo(c)
	}

	res, err := h.saldo.UpdateSaldo(ctx, &pb.UpdateSaldoRequest{
		SaldoId:      int32(idint),
		CardNumber:   body.CardNumber,
		TotalBalance: int32(body.TotalBalance),
	})

	if err != nil {
		logError("Failed to update saldo", err, zap.Error(err))

		return saldo_errors.ErrApiFailedUpdateSaldo(c)
	}

	so := h.mapper.ToApiResponseSaldo(res)

	logSuccess("Successfully update saldo", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Soft delete a saldo
// @Tags Saldo-Command
// @Security Bearer
// @Description Soft delete an existing saldo record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Saldo ID"
// @Success 200 {object} response.ApiResponseSaldo "Successfully trashed saldo record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to trashed saldo"
// @Router /api/saldo-command/trashed/{id} [post]
func (h *saldoCommandHandleApi) TrashSaldo(c echo.Context) error {
	const method = "TrashSaldo"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		logError("Bad Request: Invalid ID", err, zap.Error(err))

		return saldo_errors.ErrApiInvalidSaldoID(c)
	}

	res, err := h.saldo.TrashedSaldo(ctx, &pb.FindByIdSaldoRequest{
		SaldoId: int32(idInt),
	})

	if err != nil {
		logError("Failed to trashed saldo", err, zap.Error(err))

		return saldo_errors.ErrApiFailedTrashSaldo(c)
	}

	so := h.mapper.ToApiResponseSaldoDeleteAt(res)

	logSuccess("Successfully trashed saldo", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Restore a trashed saldo
// @Tags Saldo-Command
// @Security Bearer
// @Description Restore an existing saldo record from the trash by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Saldo ID"
// @Success 200 {object} response.ApiResponseSaldo "Successfully restored saldo record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore saldo"
// @Router /api/saldo-command/restore/{id} [post]
func (h *saldoCommandHandleApi) RestoreSaldo(c echo.Context) error {
	const method = "RestoreSaldo"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		logError("Bad Request: Invalid ID", err, zap.Error(err))

		return saldo_errors.ErrApiInvalidSaldoID(c)
	}

	res, err := h.saldo.RestoreSaldo(ctx, &pb.FindByIdSaldoRequest{
		SaldoId: int32(idInt),
	})

	if err != nil {
		logError("Failed to restore saldo", err, zap.Error(err))

		return saldo_errors.ErrApiFailedRestoreSaldo(c)
	}

	so := h.mapper.ToApiResponseSaldo(res)

	logSuccess("Successfully restored saldo", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Permanently delete a saldo
// @Tags Saldo-Command
// @Security Bearer
// @Description Permanently delete an existing saldo record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Saldo ID"
// @Success 200 {object} response.ApiResponseSaldoDelete "Successfully deleted saldo record"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete saldo"
// @Router /api/saldo-command/permanent/{id} [delete]
func (h *saldoCommandHandleApi) Delete(c echo.Context) error {
	const method = "Delete"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		logError("Bad Request: Invalid ID", err, zap.Error(err))

		return saldo_errors.ErrApiInvalidSaldoID(c)
	}

	res, err := h.saldo.DeleteSaldoPermanent(ctx, &pb.FindByIdSaldoRequest{
		SaldoId: int32(idInt),
	})

	if err != nil {
		logError("Failed to delete saldo", err, zap.Error(err))

		return saldo_errors.ErrApiFailedDeleteSaldoPermanent(c)
	}

	so := h.mapper.ToApiResponseSaldoDelete(res)

	logSuccess("Successfully deleted saldo", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// RestoreAllSaldo restores all saldo records.
// @Summary Restore all saldo records
// @Tags Saldo-Command
// @Security Bearer
// @Description Restore all saldo records that were previously deleted.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseSaldoAll "Successfully restored all saldo records"
// @Failure 500 {object} response.ErrorResponse "Failed to restore all saldo records"
// @Router /api/saldo-command/restore/all [post]
func (h *saldoCommandHandleApi) RestoreAllSaldo(c echo.Context) error {
	const method = "RestoreAllSaldo"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	res, err := h.saldo.RestoreAllSaldo(ctx, &emptypb.Empty{})

	if err != nil {
		logError("Failed to restore all saldo", err, zap.Error(err))

		return saldo_errors.ErrApiFailedRestoreAllSaldo(c)
	}

	so := h.mapper.ToApiResponseSaldoAll(res)

	logSuccess("Successfully restored all saldo", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Summary Permanently delete all saldo records
// @Tags Saldo-Command
// @Security Bearer
// @Description Permanently delete all saldo records from the database.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseSaldoAll "Successfully deleted all saldo records permanently"
// @Failure 500 {object} response.ErrorResponse "Failed to permanently delete all saldo records"
// @Router /api/saldo-command/permanent/all [post]
func (h *saldoCommandHandleApi) DeleteAllSaldoPermanent(c echo.Context) error {
	const method = "DeleteAllSaldoPermanent"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end()
	}()

	res, err := h.saldo.DeleteAllSaldoPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		logError("Failed to delete all saldo permanently", err, zap.Error(err))

		return saldo_errors.ErrApiFailedDeleteAllSaldoPermanent(c)
	}

	so := h.mapper.ToApiResponseSaldoAll(res)

	logSuccess("Successfully deleted all saldo", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *saldoCommandHandleApi) startTracingAndLogging(
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

func (s *saldoCommandHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
