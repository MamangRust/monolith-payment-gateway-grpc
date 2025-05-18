package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	traceunic "github.com/MamangRust/monolith-payment-gateway-pkg/trace_unic"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors"
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
	trace                  trace.Tracer
	cardRepository         repository.CardRepository
	logger                 logger.LoggerInterface
	mapping                responseservice.SaldoResponseMapper
	saldoCommandRepository repository.SaldoCommandRepository
	requestCounter         *prometheus.CounterVec
	requestDuration        *prometheus.HistogramVec
}

func NewSaldoCommandService(ctx context.Context, saldo repository.SaldoCommandRepository, card repository.CardRepository, logger logger.LoggerInterface, mapping responseservice.SaldoResponseMapper) *saldoCommandService {
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
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &saldoCommandService{
		ctx:                    ctx,
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
		traceID := traceunic.GenerateTraceID("FAILED_FIND_CARD_BY_CARD_NUMBER")

		s.logger.Error("Card not found for card number",
			zap.String("card_number", request.CardNumber),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Card not found for card number")
		status = "card_not_found_for_card_number"

		return nil, card_errors.ErrCardNotFoundRes
	}

	res, err := s.saldoCommandRepository.CreateSaldo(request)

	if err != nil {
		s.logger.Error("Failed to create saldo record",
			zap.Error(err))

		return nil, saldo_errors.ErrFailedCreateSaldo
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
		traceID := traceunic.GenerateTraceID("FAILED_FIND_CARD_BY_CARD_NUMBER")

		s.logger.Error("Card not found for card number",
			zap.String("card_number", request.CardNumber),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Card not found for card number")
		status = "card_not_found_for_card_number"

		return nil, card_errors.ErrCardNotFoundRes
	}

	res, err := s.saldoCommandRepository.UpdateSaldo(request)

	if err != nil {
		s.logger.Error("Failed to update saldo", zap.Error(err), zap.String("card_number", request.CardNumber))
		return nil, saldo_errors.ErrFailedUpdateSaldo
	}

	so := s.mapping.ToSaldoResponse(res)

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
		traceID := traceunic.GenerateTraceID("FAILED_TRASH_SALDO")

		s.logger.Error("Failed to trash saldo",
			zap.Int("saldo_id", saldo_id),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to trash saldo")
		status = "failed_to_trash_saldo"

		return nil, saldo_errors.ErrFailedTrashSaldo
	}
	so := s.mapping.ToSaldoResponse(res)

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
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_SALDO")

		s.logger.Error("Failed to restore saldo",
			zap.Int("saldo_id", saldo_id),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore saldo")
		status = "failed_to_restore_saldo"

		return nil, saldo_errors.ErrFailedRestoreSaldo
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
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_SALDO_PERMANENT")

		s.logger.Error("Failed to permanently delete saldo",
			zap.Int("saldo_id", saldo_id),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete saldo")
		status = "failed_to_permanently_delete_saldo"

		return false, saldo_errors.ErrFailedDeleteSaldoPermanent
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
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ALL_SALDO")

		s.logger.Error("Failed to restore all saldo", zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore all saldo")
		status = "failed_to_restore_all_saldo"

		return false, saldo_errors.ErrFailedRestoreAllSaldo
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
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_ALL_SALDO_PERMANENT")

		s.logger.Error("Failed to permanently delete all saldo", zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete all saldo")
		status = "failed_to_permanently_delete_all_saldo"

		return false, saldo_errors.ErrFailedDeleteAllSaldoPermanent
	}

	s.logger.Debug("Successfully deleted all saldo permanently")
	return true, nil
}

func (s *saldoCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
