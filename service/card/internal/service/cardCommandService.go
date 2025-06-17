package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-card/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-card/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-card/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors"
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
	errorhandler          errorhandler.CardCommandErrorHandler
	mencache              mencache.CardCommandCache
	trace                 trace.Tracer
	kafka                 *kafka.Kafka
	userRepository        repository.UserRepository
	cardCommentRepository repository.CardCommandRepository
	logger                logger.LoggerInterface
	mapping               responseservice.CardResponseMapper
	requestCounter        *prometheus.CounterVec
	requestDuration       *prometheus.HistogramVec
}

func NewCardCommandService(ctx context.Context, errorHandler errorhandler.CardCommandErrorHandler, mencache mencache.CardCommandCache, kafka *kafka.Kafka, userRepository repository.UserRepository, cardCommentRepository repository.CardCommandRepository, logger logger.LoggerInterface, mapper responseservice.CardResponseMapper) *cardCommandService {
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
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &cardCommandService{
		ctx:                   ctx,
		errorhandler:          errorHandler,
		mencache:              mencache,
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
	const method = "CreateCard"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	_, err := s.userRepository.FindById(request.UserID)

	if err != nil {
		return s.errorhandler.HandleFindByIdUserError(err, "CreateCard", "FAILED_USER_NOT_FOUND", span, &status, zap.Error(err))
	}

	res, err := s.cardCommentRepository.CreateCard(request)

	if err != nil {
		return s.errorhandler.HandleCreateCardError(err, "CreateCard", "FAILED_CREATE_CARD", span, &status, zap.Error(err))
	}

	saldoPayload := map[string]any{
		"card_number":   res.CardNumber,
		"total_balance": 0,
	}

	payloadBytes, err := json.Marshal(saldoPayload)

	if err != nil {
		return errorhandler.HandleMarshalError[*response.CardResponse](s.logger, err, "CreateCard", "FAILED_CREATE_CARD", span, &status, card_errors.ErrFailedCreateCard, zap.Error(err))
	}

	err = s.kafka.SendMessage("saldo-service-topic-create-saldo", strconv.Itoa(res.ID), payloadBytes)

	if err != nil {
		return errorhandler.HandleSendEmailError[*response.CardResponse](s.logger, err, "CreateCard", "FAILED_CREATE_CARD", span, &status, card_errors.ErrFailedCreateCard, zap.Error(err))
	}

	so := s.mapping.ToCardResponse(res)

	logSuccess("Successfully created card", zap.Int("card.id", so.ID), zap.String("card.card_number", so.CardNumber))

	return so, nil
}

func (s *cardCommandService) UpdateCard(request *requests.UpdateCardRequest) (*response.CardResponse, *response.ErrorResponse) {
	const method = "UpdateCard"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	_, err := s.userRepository.FindById(request.UserID)

	if err != nil {
		status = "error"

		return s.errorhandler.HandleFindByIdUserError(err, "UpdateCard", "FAILED_USER_NOT_FOUND", span, &status, zap.Error(err))
	}

	res, err := s.cardCommentRepository.UpdateCard(request)

	if err != nil {
		return s.errorhandler.HandleUpdateCardError(err, "UpdateCard", "FAILED_UPDATE_CARD", span, &status, zap.Error(err))
	}

	so := s.mapping.ToCardResponse(res)

	s.mencache.DeleteCardCommandCache(request.CardID)

	logSuccess("Successfully updated card", zap.Int("card.id", so.ID))

	return so, nil
}

func (s *cardCommandService) TrashedCard(card_id int) (*response.CardResponse, *response.ErrorResponse) {
	const method = "Login"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("card.id", card_id))

	defer func() {
		end(status)
	}()

	res, err := s.cardCommentRepository.TrashedCard(card_id)

	if err != nil {
		return s.errorhandler.HandleTrashedCardError(err, "TrashedCard", "FAILED_TO_TRASH_CARD", span, &status, zap.Error(err))
	}

	so := s.mapping.ToCardResponse(res)

	s.mencache.DeleteCardCommandCache(card_id)

	logSuccess("Successfully trashed card", zap.Int("card.id", so.ID))

	return so, nil
}

func (s *cardCommandService) RestoreCard(card_id int) (*response.CardResponse, *response.ErrorResponse) {
	const method = "RestoreCard"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("card.id", card_id))

	defer func() {
		end(status)
	}()

	res, err := s.cardCommentRepository.RestoreCard(card_id)

	if err != nil {
		return s.errorhandler.HandleRestoreCardError(err, "RestoreCard", "FAILED_TO_RESTORE_CARD", span, &status, zap.Error(err))
	}

	so := s.mapping.ToCardResponse(res)

	logSuccess("Successfully restored card", zap.Int("card.id", so.ID))

	return so, nil
}

func (s *cardCommandService) DeleteCardPermanent(card_id int) (bool, *response.ErrorResponse) {
	const method = "DeleteCardPermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("card.id", card_id))

	defer func() {
		end(status)
	}()

	_, err := s.cardCommentRepository.DeleteCardPermanent(card_id)

	if err != nil {
		return s.errorhandler.HandleDeleteCardPermanentError(err, "DeleteCardPermanent", "FAILED_TO_DELETE_CARD_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted card permanently", zap.Int("card.id", card_id))

	return true, nil
}

func (s *cardCommandService) RestoreAllCard() (bool, *response.ErrorResponse) {
	const method = "RestoreAllCard"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	_, err := s.cardCommentRepository.RestoreAllCard()

	if err != nil {
		return s.errorhandler.HandleRestoreAllCardError(err, "RestoreAllCard", "FAILED_TO_RESTORE_ALL_CARDS", span, &status, zap.Error(err))
	}

	logSuccess("Successfully restored all cards", zap.Bool("success", true))

	return true, nil
}

func (s *cardCommandService) DeleteAllCardPermanent() (bool, *response.ErrorResponse) {
	const method = "DeleteAllCardPermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	_, err := s.cardCommentRepository.DeleteAllCardPermanent()

	if err != nil {
		return s.errorhandler.HandleDeleteAllCardPermanentError(err, "DeleteAllCardPermanent", "FAILED_TO_DELETE_ALL_CARDS_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted all cards permanently", zap.Bool("success", true))

	return true, nil
}

func (s *cardCommandService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *cardCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
