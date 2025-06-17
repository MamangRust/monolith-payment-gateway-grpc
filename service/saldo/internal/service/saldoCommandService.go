package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-saldo/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type saldoCommandService struct {
	ctx                    context.Context
	errorhandler           errorhandler.SaldoCommandErrorHandler
	mencache               mencache.SaldoCommandCache
	trace                  trace.Tracer
	cardRepository         repository.CardRepository
	logger                 logger.LoggerInterface
	mapping                responseservice.SaldoResponseMapper
	saldoCommandRepository repository.SaldoCommandRepository
	requestCounter         *prometheus.CounterVec
	requestDuration        *prometheus.HistogramVec
}

func NewSaldoCommandService(ctx context.Context, errorhandler errorhandler.SaldoCommandErrorHandler,
	mencache mencache.SaldoCommandCache, saldo repository.SaldoCommandRepository, card repository.CardRepository, logger logger.LoggerInterface, mapping responseservice.SaldoResponseMapper) *saldoCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "saldo_command_service_request_total",
			Help: "Total number of requests to the SaldoCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "saldo_command_service_request_duration_seconds",
			Help:    "Histogram of request durations for the SaldoCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &saldoCommandService{
		ctx:                    ctx,
		errorhandler:           errorhandler,
		mencache:               mencache,
		trace:                  otel.Tracer("saldo-command-service"),
		saldoCommandRepository: saldo,
		cardRepository:         card,
		logger:                 logger,
		mapping:                mapping,
		requestCounter:         requestCounter,
		requestDuration:        requestDuration,
	}
}

func (s *saldoCommandService) CreateSaldo(request *requests.CreateSaldoRequest) (*response.SaldoResponse, *response.ErrorResponse) {
	const method = "CreateSaldo"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	_, err := s.cardRepository.FindCardByCardNumber(request.CardNumber)

	if err != nil {
		return s.errorhandler.HandleFindCardByNumberError(err, method, "FAILED_FIND_CARD_BY_CARD_NUMBER", span, &status, zap.Error(err))
	}

	res, err := s.saldoCommandRepository.CreateSaldo(request)

	if err != nil {
		return s.errorhandler.HandleCreateSaldoError(err, method, "FAILED_CREATE_SALDO", span, &status, zap.Error(err))
	}

	so := s.mapping.ToSaldoResponse(res)

	logSuccess("Successfully created saldo record", zap.String("card_number", request.CardNumber), zap.Float64("amount", float64(request.TotalBalance)))

	return so, nil
}

func (s *saldoCommandService) UpdateSaldo(request *requests.UpdateSaldoRequest) (*response.SaldoResponse, *response.ErrorResponse) {
	const method = "UpdateSaldo"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	_, err := s.cardRepository.FindCardByCardNumber(request.CardNumber)

	if err != nil {
		return s.errorhandler.HandleFindCardByNumberError(err, method, "FAILED_FIND_CARD_BY_CARD_NUMBER", span, &status, zap.Error(err))
	}

	res, err := s.saldoCommandRepository.UpdateSaldo(request)

	if err != nil {
		return s.errorhandler.HandleUpdateSaldoError(err, "UpdateSaldo", "FAILED_UPDATE_SALDO", span, &status, zap.Error(err))
	}

	so := s.mapping.ToSaldoResponse(res)

	s.mencache.DeleteSaldoCache(res.ID)

	logSuccess("Successfully updated saldo record", zap.String("card_number", request.CardNumber), zap.Float64("amount", float64(request.TotalBalance)))

	return so, nil
}

func (s *saldoCommandService) TrashSaldo(saldo_id int) (*response.SaldoResponse, *response.ErrorResponse) {
	const method = "TrashSaldo"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	res, err := s.saldoCommandRepository.TrashedSaldo(saldo_id)

	if err != nil {
		return s.errorhandler.HandleTrashSaldoError(err, method, "FAILED_TRASH_SALDO", span, &status, zap.Error(err))
	}
	so := s.mapping.ToSaldoResponse(res)

	s.mencache.DeleteSaldoCache(saldo_id)

	logSuccess("Successfully trashed saldo record", zap.Int("saldo.id", saldo_id))

	return so, nil
}

func (s *saldoCommandService) RestoreSaldo(saldo_id int) (*response.SaldoResponse, *response.ErrorResponse) {
	const method = "RestoreSaldo"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	res, err := s.saldoCommandRepository.RestoreSaldo(saldo_id)

	if err != nil {
		return s.errorhandler.HandleRestoreSaldoError(err, method, "FAILED_RESTORE_SALDO", span, &status, zap.Error(err))
	}

	so := s.mapping.ToSaldoResponse(res)

	logSuccess("Successfully restored saldo record", zap.Int("saldo.id", saldo_id))

	return so, nil
}

func (s *saldoCommandService) DeleteSaldoPermanent(saldo_id int) (bool, *response.ErrorResponse) {
	const method = "DeleteSaldoPermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	_, err := s.saldoCommandRepository.DeleteSaldoPermanent(saldo_id)

	if err != nil {
		return s.errorhandler.HandleDeleteSaldoPermanentError(err, method, "FAILED_DELETE_SALDO_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted permanent saldo record", zap.Int("saldo.id", saldo_id))

	return true, nil
}

func (s *saldoCommandService) RestoreAllSaldo() (bool, *response.ErrorResponse) {
	const method = "RestoreAllSaldo"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	_, err := s.saldoCommandRepository.RestoreAllSaldo()

	if err != nil {
		return s.errorhandler.HandleRestoreAllSaldoError(err, method, "FAILED_RESTORE_ALL_SALDO", span, &status, zap.Error(err))
	}

	logSuccess("Successfully restored all saldo", zap.Bool("success", true))

	return true, nil
}

func (s *saldoCommandService) DeleteAllSaldoPermanent() (bool, *response.ErrorResponse) {
	const method = "DeleteAllSaldoPermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	_, err := s.saldoCommandRepository.DeleteAllSaldoPermanent()

	if err != nil {
		return s.errorhandler.HandleDeleteAllSaldoPermanentError(err, method, "FAILED_DELETE_ALL_SALDO_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted all permanent saldo", zap.Bool("success", true))

	return true, nil
}

func (s *saldoCommandService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
	trace.Span,
	func(string),
	string,
	func(string, ...zap.Field),
) {
	start := time.Now()
	status := "success"

	_, span := s.trace.Start(s.ctx, method)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.AddEvent("Start: " + method)

	s.logger.Info("Start: " + method)

	end := func(status string) {
		s.recordMetrics(method, status, start)
		code := codes.Ok
		if status != "success" {
			code = codes.Error
		}
		span.SetStatus(code, status)
		span.End()
	}

	logSuccess := func(msg string, fields ...zap.Field) {
		span.AddEvent(msg)
		s.logger.Info(msg, fields...)
	}

	return span, end, status, logSuccess
}

func (s *saldoCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
