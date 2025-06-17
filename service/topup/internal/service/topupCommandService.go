package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/email"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-topup/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type topupCommandService struct {
	kafka                  *kafka.Kafka
	errorhandler           errorhandler.TopupCommandErrorHandler
	mencache               mencache.TopupCommandCache
	ctx                    context.Context
	trace                  trace.Tracer
	topupQueryRepository   repository.TopupQueryRepository
	cardRepository         repository.CardRepository
	topupCommandRepository repository.TopupCommandRepository
	saldoRepository        repository.SaldoRepository
	logger                 logger.LoggerInterface
	mapping                responseservice.TopupResponseMapper
	requestCounter         *prometheus.CounterVec
	requestDuration        *prometheus.HistogramVec
}

func NewTopupCommandService(
	kafka *kafka.Kafka,
	errorhandler errorhandler.TopupCommandErrorHandler,
	mencache mencache.TopupCommandCache,
	ctx context.Context,
	cardRepository repository.CardRepository,
	topupQueryRepository repository.TopupQueryRepository,
	topupCommandRepository repository.TopupCommandRepository,
	saldoRepository repository.SaldoRepository,
	logger logger.LoggerInterface, mapping responseservice.TopupResponseMapper) *topupCommandService {

	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "topup_command_service_request_total",
			Help: "Total number of requests to the TopupCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "topup_command_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TopupCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &topupCommandService{
		kafka:                  kafka,
		errorhandler:           errorhandler,
		mencache:               mencache,
		ctx:                    ctx,
		trace:                  otel.Tracer("topup-command-service"),
		topupQueryRepository:   topupQueryRepository,
		topupCommandRepository: topupCommandRepository,
		saldoRepository:        saldoRepository,
		cardRepository:         cardRepository,
		logger:                 logger,
		mapping:                mapping,
		requestCounter:         requestCounter,
		requestDuration:        requestDuration,
	}
}

func (s *topupCommandService) CreateTopup(request *requests.CreateTopupRequest) (*response.TopupResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("CreateTopup", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "CreateTopup")
	defer span.End()

	span.SetAttributes(
		attribute.String("card_number", request.CardNumber),
		attribute.Float64("topup_amount", float64(request.TopupAmount)),
	)

	s.logger.Debug("Starting CreateTopup process",
		zap.String("cardNumber", request.CardNumber),
		zap.Float64("topupAmount", float64(request.TopupAmount)),
	)

	card, err := s.cardRepository.FindUserCardByCardNumber(request.CardNumber)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "CreateTopup", "FAILED_FIND_CARD_BY_CARD_NUMBER", span, &status, card_errors.ErrCardNotFoundRes)
	}

	topup, err := s.topupCommandRepository.CreateTopup(request)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "CreateTopup", "FAILED_CREATE_TOPUP", span, &status, topup_errors.ErrFailedCreateTopup)
	}

	saldo, err := s.saldoRepository.FindByCardNumber(request.CardNumber)
	if err != nil {
		req := requests.UpdateTopupStatus{
			TopupID: topup.ID,
			Status:  "failed",
		}
		s.topupCommandRepository.UpdateTopupStatus(&req)

		return s.errorhandler.HandleRepositorySingleError(err, "CreateTopup", "FAILED_FIND_SALDO_BY_CARD_NUMBER", span, &status, saldo_errors.ErrFailedSaldoNotFound)
	}

	newBalance := saldo.TotalBalance + request.TopupAmount
	_, err = s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   request.CardNumber,
		TotalBalance: newBalance,
	})
	if err != nil {

		req := requests.UpdateTopupStatus{
			TopupID: topup.ID,
			Status:  "failed",
		}
		s.topupCommandRepository.UpdateTopupStatus(&req)

		return s.errorhandler.HandleRepositorySingleError(err, "CreateTopup", "FAILED_UPDATE_SALDO_BALANCE", span, &status, topup_errors.ErrFailedUpdateTopup)
	}

	expireDate, err := time.Parse("2006-01-02", card.ExpireDate)
	if err != nil {
		req := requests.UpdateTopupStatus{
			TopupID: topup.ID,
			Status:  "failed",
		}
		s.topupCommandRepository.UpdateTopupStatus(&req)

		return s.errorhandler.HandleInvalidParseTimeError(err, "CreateTopup", "FAILED_PARSE_EXPIRE_DATE", span, &status, card.ExpireDate, zap.String("expire_date", card.ExpireDate))
	}

	_, err = s.cardRepository.UpdateCard(&requests.UpdateCardRequest{
		CardID:       card.ID,
		UserID:       card.UserID,
		CardType:     card.CardType,
		ExpireDate:   expireDate,
		CVV:          card.CVV,
		CardProvider: card.CardProvider,
	})
	if err != nil {
		req := requests.UpdateTopupStatus{
			TopupID: topup.ID,
			Status:  "failed",
		}
		s.topupCommandRepository.UpdateTopupStatus(&req)

		return s.errorhandler.HandleRepositorySingleError(err, "CreateTopup", "FAILED_UPDATE_CARD", span, &status, topup_errors.ErrFailedUpdateTopup)
	}

	req := requests.UpdateTopupStatus{
		TopupID: topup.ID,
		Status:  "success",
	}

	res, err := s.topupCommandRepository.UpdateTopupStatus(&req)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "CreateTopup", "FAILED_UPDATE_TOPUP_STATUS", span, &status, topup_errors.ErrFailedUpdateTopup)
	}

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   "Topup Successful",
		"Message": fmt.Sprintf("Your topup of %d has been processed successfully.", request.TopupAmount),
		"Button":  "View History",
		"Link":    "https://sanedge.example.com/topup/history",
	})

	emailPayload := map[string]any{
		"email":   card.Email,
		"subject": "Topup Successful - SanEdge",
		"body":    htmlBody,
	}

	payloadBytes, err := json.Marshal(emailPayload)
	if err != nil {
		return errorhandler.HandleErrorJSONMarshal[*response.TopupResponse](s.logger, err, "CreateTopup", "FAILED_JSON_MARSHAL", span, &status, topup_errors.ErrFailedUpdateTopup)
	}

	err = s.kafka.SendMessage("email-service-topic-topup-create", strconv.Itoa(res.ID), payloadBytes)
	if err != nil {
		return errorhandler.HandleErrorKafkaSend[*response.TopupResponse](s.logger, err, "CreateTopup", "FAILED_KAFKA_SEND", span, &status, topup_errors.ErrFailedUpdateTopup)
	}

	so := s.mapping.ToTopupResponse(topup)

	s.logger.Debug("CreateTopup process completed",
		zap.String("cardNumber", request.CardNumber),
		zap.Float64("topupAmount", float64(request.TopupAmount)),
		zap.Float64("newBalance", float64(newBalance)),
	)

	return so, nil
}

func (s *topupCommandService) UpdateTopup(request *requests.UpdateTopupRequest) (*response.TopupResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("update_topup", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateTopup")
	defer span.End()

	s.logger.Debug("Starting UpdateTopup process",
		zap.String("cardNumber", request.CardNumber),
		zap.Int("topupID", *request.TopupID),
		zap.Float64("newTopupAmount", float64(request.TopupAmount)),
	)

	_, err := s.cardRepository.FindCardByCardNumber(request.CardNumber)
	if err != nil {
		req := requests.UpdateTopupStatus{
			TopupID: *request.TopupID,
			Status:  "failed",
		}

		s.topupCommandRepository.UpdateTopupStatus(&req)

		return s.errorhandler.HandleRepositorySingleError(err, "UpdateTopup", "FAILED_FIND_CARD_BY_CARD_NUMBER", span, &status, topup_errors.ErrFailedUpdateTopup)
	}

	existingTopup, err := s.topupQueryRepository.FindById(*request.TopupID)
	if err != nil || existingTopup == nil {
		req := requests.UpdateTopupStatus{
			TopupID: *request.TopupID,
			Status:  "failed",
		}

		s.topupCommandRepository.UpdateTopupStatus(&req)

		return s.errorhandler.HandleRepositorySingleError(err, "UpdateTopup", "FAILED_FIND_TOPUP_BY_ID", span, &status, topup_errors.ErrFailedUpdateTopup)
	}

	topupDifference := request.TopupAmount - existingTopup.TopupAmount

	_, err = s.topupCommandRepository.UpdateTopup(request)
	if err != nil {
		req := requests.UpdateTopupStatus{
			TopupID: *request.TopupID,
			Status:  "failed",
		}

		s.topupCommandRepository.UpdateTopupStatus(&req)

		return s.errorhandler.HandleRepositorySingleError(err, "UpdateTopup", "FAILED_UPDATE_TOPUP", span, &status, topup_errors.ErrFailedUpdateTopup)
	}

	currentSaldo, err := s.saldoRepository.FindByCardNumber(request.CardNumber)
	if err != nil {
		req := requests.UpdateTopupStatus{
			TopupID: *request.TopupID,
			Status:  "failed",
		}

		s.topupCommandRepository.UpdateTopupStatus(&req)

		return s.errorhandler.HandleRepositorySingleError(err, "UpdateTopup", "FAILED_FIND_SALDO_BY_CARD_NUMBER", span, &status, topup_errors.ErrFailedUpdateTopup)
	}

	if currentSaldo == nil {
		req := requests.UpdateTopupStatus{
			TopupID: *request.TopupID,
			Status:  "failed",
		}

		s.topupCommandRepository.UpdateTopupStatus(&req)

		return s.errorhandler.HandleRepositorySingleError(err, "UpdateTopup", "FAILED_FIND_SALDO_BY_CARD_NUMBER", span, &status, topup_errors.ErrFailedUpdateTopup)
	}

	newBalance := currentSaldo.TotalBalance + topupDifference
	_, err = s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   request.CardNumber,
		TotalBalance: newBalance,
	})
	if err != nil {
		_, rollbackErr := s.topupCommandRepository.UpdateTopupAmount(&requests.UpdateTopupAmount{
			TopupID:     *request.TopupID,
			TopupAmount: existingTopup.TopupAmount,
		})
		if rollbackErr != nil {
			return s.errorhandler.HandleRepositorySingleError(rollbackErr, "UpdateTopup", "FAILED_ROLLBACK_TOPUP_AMOUNT", span, &status, topup_errors.ErrFailedUpdateTopup)
		}

		req := requests.UpdateTopupStatus{
			TopupID: *request.TopupID,
			Status:  "failed",
		}

		s.topupCommandRepository.UpdateTopupStatus(&req)

		return s.errorhandler.HandleRepositorySingleError(err, "UpdateTopup", "FAILED_UPDATE_SALDO_BALANCE", span, &status, topup_errors.ErrFailedUpdateTopup)
	}

	updatedTopup, err := s.topupQueryRepository.FindById(*request.TopupID)
	if err != nil || updatedTopup == nil {
		req := requests.UpdateTopupStatus{
			TopupID: *request.TopupID,
			Status:  "failed",
		}

		s.topupCommandRepository.UpdateTopupStatus(&req)

		return s.errorhandler.HandleRepositorySingleError(err, "UpdateTopup", "FAILED_FIND_TOPUP_BY_ID", span, &status, topup_errors.ErrFailedUpdateTopup)
	}

	req := requests.UpdateTopupStatus{
		TopupID: *request.TopupID,
		Status:  "success",
	}

	_, err = s.topupCommandRepository.UpdateTopupStatus(&req)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "UpdateTopup", "FAILED_UPDATE_TOPUP_STATUS", span, &status, topup_errors.ErrFailedUpdateTopup)
	}

	so := s.mapping.ToTopupResponse(updatedTopup)

	s.mencache.DeleteCachedTopupCache(*request.TopupID)

	s.logger.Debug("UpdateTopup process completed",
		zap.String("cardNumber", request.CardNumber),
		zap.Int("topupID", *request.TopupID),
		zap.Float64("newTopupAmount", float64(request.TopupAmount)),
		zap.Float64("newBalance", float64(newBalance)),
	)

	return so, nil
}

func (s *topupCommandService) TrashedTopup(topup_id int) (*response.TopupResponseDeleteAt, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("TrashedTopup", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "TrashedTopup")
	defer span.End()

	span.SetAttributes(
		attribute.Int("topup_id", topup_id),
	)

	s.logger.Debug("Starting TrashedTopup process",
		zap.Int("topup_id", topup_id),
	)

	res, err := s.topupCommandRepository.TrashedTopup(topup_id)

	if err != nil {
		return s.errorhandler.HandleTrashedTopupError(err, "TrashedTopup", "FAILED_TRASH_TOPUP", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTopupResponseDeleteAt(res)

	s.mencache.DeleteCachedTopupCache(topup_id)

	s.logger.Debug("TrashedTopup process completed",
		zap.Int("topup_id", topup_id),
	)

	return so, nil
}

func (s *topupCommandService) RestoreTopup(topup_id int) (*response.TopupResponseDeleteAt, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreTopup", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreTopup")
	defer span.End()

	span.SetAttributes(
		attribute.Int("topup_id", topup_id),
	)

	s.logger.Debug("Starting RestoreTopup process",
		zap.Int("topupID", topup_id),
	)

	res, err := s.topupCommandRepository.RestoreTopup(topup_id)

	if err != nil {
		return s.errorhandler.HandleRestoreTopupError(err, "RestoreTopup", "FAILED_RESTORE_TOPUP", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTopupResponseDeleteAt(res)

	s.logger.Debug("RestoreTopup process completed",
		zap.Int("topupID", topup_id),
	)

	return so, nil
}

func (s *topupCommandService) DeleteTopupPermanent(topup_id int) (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteTopupPermanent", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteTopupPermanent")
	defer span.End()

	span.SetAttributes(
		attribute.Int("topup_id", topup_id),
	)

	s.logger.Debug("Starting DeleteTopupPermanent process",
		zap.Int("topupID", topup_id),
	)

	_, err := s.topupCommandRepository.DeleteTopupPermanent(topup_id)

	if err != nil {
		return s.errorhandler.HandleDeleteTopupPermanentError(err, "DeleteTopupPermanent", "FAILED_DELETE_TOPUP_PERMANENT", span, &status, zap.Error(err))
	}

	s.logger.Debug("DeleteTopupPermanent process completed",
		zap.Int("topupID", topup_id),
	)

	return true, nil
}

func (s *topupCommandService) RestoreAllTopup() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreAllTopup", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreAllTopup")
	defer span.End()

	s.logger.Debug("Restoring all topups")

	_, err := s.topupCommandRepository.RestoreAllTopup()

	if err != nil {
		return s.errorhandler.HandleRestoreAllTopupError(err, "RestoreAllTopup", "FAILED_RESTORE_ALL_TOPUP", span, &status, zap.Error(err))
	}

	s.logger.Debug("Successfully restored all topups")
	return true, nil
}

func (s *topupCommandService) DeleteAllTopupPermanent() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteAllTopupPermanent", status, start)
	}()
	_, span := s.trace.Start(s.ctx, "DeleteAllTopupPermanent")

	defer span.End()

	s.logger.Debug("Permanently deleting all topups")

	_, err := s.topupCommandRepository.DeleteAllTopupPermanent()

	if err != nil {
		return s.errorhandler.HandleDeleteAllTopupPermanentError(err, "DeleteAllTopupPermanent", "FAILED_DELETE_ALL_TOPUP_PERMANENT", span, &status, zap.Error(err))
	}

	s.logger.Debug("Successfully deleted all topups permanently")
	return true, nil
}

func (s *topupCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
