package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-card/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	traceunic "github.com/MamangRust/monolith-payment-gateway-pkg/trace_unic"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type cardCommandService struct {
	ctx                   context.Context
	trace                 trace.Tracer
	kafka                 kafka.Kafka
	userRepository        repository.UserRepository
	cardCommentRepository repository.CardCommandRepository
	logger                logger.LoggerInterface
	mapping               responseservice.CardResponseMapper
	requestCounter        *prometheus.CounterVec
	requestDuration       *prometheus.HistogramVec
}

func NewCardCommandService(ctx context.Context, kafka kafka.Kafka, userRepository repository.UserRepository, cardCommentRepository repository.CardCommandRepository, logger logger.LoggerInterface, mapper responseservice.CardResponseMapper) *cardCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "card_command_service_requests_total",
			Help: "Total number of requests to the CardCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "card_command_service_request_duration_seconds",
			Help:    "Histogram of request durations for the CardCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &cardCommandService{
		ctx:                   ctx,
		trace:                 otel.Tracer("card-command-service"),
		kafka:                 kafka,
		userRepository:        userRepository,
		cardCommentRepository: cardCommentRepository,
		logger:                logger,
		mapping:               mapper,
		requestCounter:        requestCounter,
		requestDuration:       requestDuration,
	}
}

func (s *cardCommandService) CreateCard(request *requests.CreateCardRequest) (*response.CardResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("CreateCard", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "CreateCard")
	defer span.End()

	span.SetAttributes(
		attribute.Int("user_id", request.UserID),
	)

	s.logger.Debug("Creating new card", zap.Any("request", request))

	_, err := s.userRepository.FindById(request.UserID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("EMAIL_NOT_FOUND")

		s.logger.Error("Failed to get user", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "User not found")
		status = "user_not_found"

		return nil, user_errors.ErrUserNotFoundRes
	}

	res, err := s.cardCommentRepository.CreateCard(request)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_CREATE_CARD")

		s.logger.Error("Failed to create card", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create card")
		status = "failed_create_card"
		return nil, card_errors.ErrFailedCreateCard
	}

	saldoPayload := map[string]any{
		"card_number":   res.CardNumber,
		"total_balance": 0,
	}

	payloadBytes, err := json.Marshal(saldoPayload)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_JSON_MARSHAL")

		s.logger.Error("Failed to marshal saldo payload", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal saldo payload")
		status = "failed_json_marshal"

		return nil, card_errors.ErrFailedCreateCard
	}

	err = s.kafka.SendMessage("saldo-service-topic-create-saldo", strconv.Itoa(res.ID), payloadBytes)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_MESSAGE_SALDO_CREATE")

		s.logger.Error("Failed to send message to saldo service", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to send message to saldo service")
		status = "failed_message_saldo_create"

		return nil, card_errors.ErrFailedCreateCard
	}

	so := s.mapping.ToCardResponse(res)

	s.logger.Debug("Successfully created new card", zap.Int("card_id", so.ID))

	return so, nil
}

func (s *cardCommandService) UpdateCard(request *requests.UpdateCardRequest) (*response.CardResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("UpdateCard", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateCard")
	defer span.End()

	span.SetAttributes(
		attribute.Int("card_id", request.CardID),
	)

	s.logger.Debug("Updating card", zap.Int("card_id", request.CardID), zap.Any("request", request))

	_, err := s.userRepository.FindById(request.UserID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("EMAIL_NOT_FOUND")

		s.logger.Error("Failed to get user", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "User not found")
		status = "user_not_found"

		return nil, user_errors.ErrUserNotFoundRes
	}

	res, err := s.cardCommentRepository.UpdateCard(request)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_CARD")

		s.logger.Error("Failed to update card", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update card")
		status = "failed_update_card"

		return nil, card_errors.ErrFailedUpdateCard
	}

	so := s.mapping.ToCardResponse(res)

	s.logger.Debug("Successfully updated card", zap.Int("cardID", so.ID))

	return so, nil
}

func (s *cardCommandService) TrashedCard(card_id int) (*response.CardResponse, *response.ErrorResponse) {
	s.logger.Debug("Trashing card", zap.Int("card_id", card_id))

	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("TrashedCard", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "TrashedCard")
	defer span.End()

	span.SetAttributes(
		attribute.Int("card_id", card_id),
	)

	res, err := s.cardCommentRepository.TrashedCard(card_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_TRASH_CARD")

		s.logger.Error("Failed to trash card", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to trash card")
		status = "failed_trash_card"
		return nil, card_errors.ErrFailedTrashCard
	}

	so := s.mapping.ToCardResponse(res)

	s.logger.Debug("Successfully trashed card", zap.Int("card_id", so.ID))

	return so, nil
}

func (s *cardCommandService) RestoreCard(card_id int) (*response.CardResponse, *response.ErrorResponse) {
	s.logger.Debug("Restoring card", zap.Int("card_id", card_id))

	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreCard", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreCard")
	defer span.End()

	span.SetAttributes(
		attribute.Int("card_id", card_id),
	)

	res, err := s.cardCommentRepository.RestoreCard(card_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_CARD")

		s.logger.Error("Failed to restore card", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore card")
		status = "failed_restore_card"

		return nil, card_errors.ErrFailedRestoreCard
	}

	so := s.mapping.ToCardResponse(res)

	s.logger.Debug("Successfully restored card", zap.Int("card_id", so.ID))

	return so, nil
}

func (s *cardCommandService) DeleteCardPermanent(card_id int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting card", zap.Int("card_id", card_id))

	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteCardPermanent", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteCardPermanent")
	defer span.End()

	span.SetAttributes(
		attribute.Int("card_id", card_id),
	)

	_, err := s.cardCommentRepository.DeleteCardPermanent(card_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_CARD")

		s.logger.Error("Failed to permanently delete card", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete card")
		status = "failed_delete_card"

		return false, card_errors.ErrFailedDeleteCard
	}

	s.logger.Debug("Successfully deleted card permanently", zap.Int("card_id", card_id))

	return true, nil
}

func (s *cardCommandService) RestoreAllCard() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all cards")

	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreAllCard", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreAllCard")
	defer span.End()

	_, err := s.cardCommentRepository.RestoreAllCard()

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ALL_CARDS")

		s.logger.Error("Failed to restore all cards", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore all cards")
		status = "failed_restore_all_cards"

		return false, card_errors.ErrFailedRestoreAllCards
	}

	s.logger.Debug("Successfully restored all cards")
	return true, nil
}

func (s *cardCommandService) DeleteAllCardPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all cards")

	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteAllCardPermanent", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteAllCardPermanent")
	defer span.End()

	_, err := s.cardCommentRepository.DeleteAllCardPermanent()

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_ALL_CARDS")

		s.logger.Error("Failed to permanently delete all cards", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete all cards")
		status = "failed_delete_all_cards"
		return false, card_errors.ErrFailedDeleteAllCards
	}

	s.logger.Debug("Successfully deleted all cards permanently")

	return true, nil
}

func (s *cardCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
