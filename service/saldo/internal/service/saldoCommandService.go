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
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("CreateSaldo", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "CreateSaldo")
	defer span.End()

	span.SetAttributes(
		attribute.String("card_number", request.CardNumber),
	)

	s.logger.Debug("Creating saldo record", zap.String("card_number", request.CardNumber))

	_, err := s.cardRepository.FindCardByCardNumber(request.CardNumber)

	if err != nil {
		return s.errorhandler.HandleFindCardByNumberError(err, "CreateSaldo", "FAILED_FIND_CARD_BY_CARD_NUMBER", span, &status, zap.String("card_number", request.CardNumber))
	}

	res, err := s.saldoCommandRepository.CreateSaldo(request)

	if err != nil {
		return s.errorhandler.HandleCreateSaldoError(err, "CreateSaldo", "FAILED_CREATE_SALDO", span, &status, zap.String("card_number", request.CardNumber))
	}

	so := s.mapping.ToSaldoResponse(res)

	s.logger.Debug("Successfully created saldo record", zap.String("card_number", request.CardNumber))

	return so, nil
}

func (s *saldoCommandService) UpdateSaldo(request *requests.UpdateSaldoRequest) (*response.SaldoResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("UpdateSaldo", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateSaldo")
	defer span.End()

	span.SetAttributes(
		attribute.String("card_number", request.CardNumber),
	)

	s.logger.Debug("Updating saldo record", zap.String("card_number", request.CardNumber), zap.Float64("amount", float64(request.TotalBalance)))

	_, err := s.cardRepository.FindCardByCardNumber(request.CardNumber)

	if err != nil {
		return s.errorhandler.HandleFindCardByNumberError(err, "UpdateSaldo", "FAILED_FIND_CARD_BY_CARD_NUMBER", span, &status, zap.String("card_number", request.CardNumber))
	}

	res, err := s.saldoCommandRepository.UpdateSaldo(request)

	if err != nil {
		return s.errorhandler.HandleUpdateSaldoError(err, "UpdateSaldo", "FAILED_UPDATE_SALDO", span, &status, zap.String("card_number", request.CardNumber))
	}

	so := s.mapping.ToSaldoResponse(res)

	s.mencache.DeleteSaldoCache(res.ID)

	s.logger.Debug("Successfully updated saldo", zap.String("card_number", request.CardNumber), zap.Int("saldo_id", res.ID))

	return so, nil
}

func (s *saldoCommandService) TrashSaldo(saldo_id int) (*response.SaldoResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("TrashSaldo", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "TrashSaldo")
	defer span.End()

	span.SetAttributes(
		attribute.Int("saldo_id", saldo_id),
	)

	s.logger.Debug("Trashing saldo record", zap.Int("saldo_id", saldo_id))

	res, err := s.saldoCommandRepository.TrashedSaldo(saldo_id)

	if err != nil {
		return s.errorhandler.HandleTrashSaldoError(err, "TrashSaldo", "FAILED_TRASH_SALDO", span, &status, zap.Int("saldo_id", saldo_id))
	}
	so := s.mapping.ToSaldoResponse(res)

	s.mencache.DeleteSaldoCache(saldo_id)

	s.logger.Debug("Successfully trashed saldo", zap.Int("saldo_id", saldo_id))

	return so, nil
}

func (s *saldoCommandService) RestoreSaldo(saldo_id int) (*response.SaldoResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreSaldo", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreSaldo")
	defer span.End()

	span.SetAttributes(
		attribute.Int("saldo_id", saldo_id),
	)

	s.logger.Debug("Restoring saldo record from trash", zap.Int("saldo_id", saldo_id))

	res, err := s.saldoCommandRepository.RestoreSaldo(saldo_id)

	if err != nil {
		return s.errorhandler.HandleRestoreSaldoError(err, "RestoreSaldo", "FAILED_RESTORE_SALDO", span, &status, zap.Int("saldo_id", saldo_id))
	}

	so := s.mapping.ToSaldoResponse(res)

	s.logger.Debug("Successfully restored saldo", zap.Int("saldo_id", saldo_id))

	return so, nil
}

func (s *saldoCommandService) DeleteSaldoPermanent(saldo_id int) (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteSaldoPermanent", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteSaldoPermanent")
	defer span.End()

	span.SetAttributes(
		attribute.Int("saldo_id", saldo_id),
	)

	s.logger.Debug("Deleting saldo permanently", zap.Int("saldo_id", saldo_id))

	_, err := s.saldoCommandRepository.DeleteSaldoPermanent(saldo_id)

	if err != nil {
		return s.errorhandler.HandleDeleteSaldoPermanentError(err, "DeleteSaldoPermanent", "FAILED_DELETE_SALDO_PERMANENT", span, &status, zap.Int("saldo_id", saldo_id))
	}

	s.logger.Debug("Successfully deleted saldo permanently", zap.Int("saldo_id", saldo_id))

	return true, nil
}

func (s *saldoCommandService) RestoreAllSaldo() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreAllSaldo", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreAllSaldo")
	defer span.End()

	_, err := s.saldoCommandRepository.RestoreAllSaldo()

	if err != nil {
		return s.errorhandler.HandleRestoreAllSaldoError(err, "RestoreAllSaldo", "FAILED_RESTORE_ALL_SALDO", span, &status, zap.Error(err))
	}

	s.logger.Debug("Successfully restored all saldo")
	return true, nil
}

func (s *saldoCommandService) DeleteAllSaldoPermanent() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteAllSaldoPermanent", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteAllSaldoPermanent")
	defer span.End()

	s.logger.Debug("Permanently deleting all saldo")

	_, err := s.saldoCommandRepository.DeleteAllSaldoPermanent()

	if err != nil {
		return s.errorhandler.HandleDeleteAllSaldoPermanentError(err, "DeleteAllSaldoPermanent", "FAILED_DELETE_ALL_SALDO_PERMANENT", span, &status, zap.Error(err))
	}

	s.logger.Debug("Successfully deleted all saldo permanently")
	return true, nil
}

func (s *saldoCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
