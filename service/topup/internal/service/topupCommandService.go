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
	traceunic "github.com/MamangRust/monolith-payment-gateway-pkg/trace_unic"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type topupCommandService struct {
	kafka                  kafka.Kafka
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
	kafka kafka.Kafka,
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
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &topupCommandService{
		kafka:                  kafka,
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
		traceID := traceunic.GenerateTraceID("FAILED_FIND_CARD_BY_CARD_NUMBER")

		s.logger.Error("Card not found for card number",
			zap.String("card_number", request.CardNumber),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Card not found for card number")
		status = "card_not_found_for_card_number"

		s.logger.Error("failed to find card by number", zap.Error(err))
		return nil, card_errors.ErrCardNotFoundRes
	}

	topup, err := s.topupCommandRepository.CreateTopup(request)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_CREATE_TOPUP")

		s.logger.Error("failed to create topup", zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create topup")
		status = "failed_to_create_topup"

		return nil, topup_errors.ErrFailedCreateTopup
	}

	saldo, err := s.saldoRepository.FindByCardNumber(request.CardNumber)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_SALDO_BY_CARD_NUMBER")

		s.logger.Error("Failed to retrieve saldo details",
			zap.String("card_number", request.CardNumber),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve saldo details")
		status = "failed_to_retrieve_saldo_details"

		req := requests.UpdateTopupStatus{
			TopupID: topup.ID,
			Status:  "failed",
		}
		s.topupCommandRepository.UpdateTopupStatus(&req)
		return nil, saldo_errors.ErrFailedSaldoNotFound
	}

	newBalance := saldo.TotalBalance + request.TopupAmount
	_, err = s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   request.CardNumber,
		TotalBalance: newBalance,
	})
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_SALDO_BALANCE")

		s.logger.Error("Failed to update saldo balance",
			zap.String("card_number", request.CardNumber),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update saldo balance")
		status = "failed_to_update_saldo_balance"

		req := requests.UpdateTopupStatus{
			TopupID: topup.ID,
			Status:  "failed",
		}
		s.topupCommandRepository.UpdateTopupStatus(&req)
		return nil, topup_errors.ErrFailedUpdateTopup
	}

	expireDate, err := time.Parse("2006-01-02", card.ExpireDate)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_PARSE_EXPIRE_DATE")

		s.logger.Error("Failed to parse expire date",
			zap.String("expire_date", card.ExpireDate),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to parse expire date")
		status = "failed_to_parse_expire_date"

		req := requests.UpdateTopupStatus{
			TopupID: topup.ID,
			Status:  "failed",
		}
		s.topupCommandRepository.UpdateTopupStatus(&req)
		return nil, topup_errors.ErrFailedUpdateTopup
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
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_CARD")

		s.logger.Error("Failed to update card",
			zap.String("card_number", request.CardNumber),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update card")
		status = "failed_to_update_card"

		req := requests.UpdateTopupStatus{
			TopupID: topup.ID,
			Status:  "failed",
		}
		s.topupCommandRepository.UpdateTopupStatus(&req)
		return nil, card_errors.ErrFailedUpdateCard
	}

	req := requests.UpdateTopupStatus{
		TopupID: topup.ID,
		Status:  "success",
	}

	res, err := s.topupCommandRepository.UpdateTopupStatus(&req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_TOPUP_STATUS")

		s.logger.Error("Failed to update topup status",
			zap.Int("topup_id", topup.ID),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update topup status")
		status = "failed_to_update_topup_status"

		return nil, topup_errors.ErrFailedUpdateTopup
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
		traceID := traceunic.GenerateTraceID("TOPUP_ERR")
		s.logger.Error("Failed to marshal topup email payload", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal topup email payload")
		return nil, withdraw_errors.ErrFailedSendEmail
	}

	err = s.kafka.SendMessage("email-service-topic-topup-create", strconv.Itoa(res.ID), payloadBytes)
	if err != nil {
		traceID := traceunic.GenerateTraceID("TOPUP_ERR")
		s.logger.Error("Failed to send topup email message", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to send topup email")
		return nil, withdraw_errors.ErrFailedSendEmail
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
		traceID := traceunic.GenerateTraceID("FAILED_FIND_CARD_BY_CARD_NUMBER")

		s.logger.Error("Card not found for card number",
			zap.String("card_number", request.CardNumber),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Card not found for card number")
		status = "card_not_found_for_card_number"

		req := requests.UpdateTopupStatus{
			TopupID: *request.TopupID,
			Status:  "failed",
		}

		s.topupCommandRepository.UpdateTopupStatus(&req)

		return nil, card_errors.ErrCardNotFoundRes
	}

	existingTopup, err := s.topupQueryRepository.FindById(*request.TopupID)
	if err != nil || existingTopup == nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TOPUP_BY_ID")

		s.logger.Error("Topup not found for topup id",
			zap.Int("topup_id", *request.TopupID),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Topup not found for topup id")
		status = "topup_not_found_for_topup_id"

		req := requests.UpdateTopupStatus{
			TopupID: *request.TopupID,
			Status:  "failed",
		}

		s.topupCommandRepository.UpdateTopupStatus(&req)
		return nil, topup_errors.ErrTopupNotFoundRes
	}

	topupDifference := request.TopupAmount - existingTopup.TopupAmount

	_, err = s.topupCommandRepository.UpdateTopup(request)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_TOPUP")

		s.logger.Error("Failed to update topup",
			zap.Int("topup_id", *request.TopupID),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update topup")
		status = "failed_to_update_topup"

		req := requests.UpdateTopupStatus{
			TopupID: *request.TopupID,
			Status:  "failed",
		}

		s.topupCommandRepository.UpdateTopupStatus(&req)
		return nil, topup_errors.ErrFailedUpdateTopup
	}

	currentSaldo, err := s.saldoRepository.FindByCardNumber(request.CardNumber)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_SALDO_BY_CARD_NUMBER")

		s.logger.Error("Failed to retrieve saldo details",
			zap.String("card_number", request.CardNumber),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve saldo details")
		status = "failed_to_retrieve_saldo_details"

		req := requests.UpdateTopupStatus{
			TopupID: *request.TopupID,
			Status:  "failed",
		}

		s.topupCommandRepository.UpdateTopupStatus(&req)

		return nil, saldo_errors.ErrFailedSaldoNotFound
	}

	if currentSaldo == nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_SALDO_BY_CARD_NUMBER")

		s.logger.Error("Failed to retrieve saldo details",
			zap.String("card_number", request.CardNumber),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve saldo details")
		status = "failed_to_retrieve_saldo_details"

		req := requests.UpdateTopupStatus{
			TopupID: *request.TopupID,
			Status:  "failed",
		}

		s.topupCommandRepository.UpdateTopupStatus(&req)

		return nil, card_errors.ErrCardNotFoundRes
	}

	newBalance := currentSaldo.TotalBalance + topupDifference
	_, err = s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   request.CardNumber,
		TotalBalance: newBalance,
	})
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_SALDO_BALANCE")

		s.logger.Error("Failed to update saldo balance",
			zap.String("card_number", request.CardNumber),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update saldo balance")
		status = "failed_to_update_saldo_balance"

		_, rollbackErr := s.topupCommandRepository.UpdateTopupAmount(&requests.UpdateTopupAmount{
			TopupID:     *request.TopupID,
			TopupAmount: existingTopup.TopupAmount,
		})
		if rollbackErr != nil {
			traceID := traceunic.GenerateTraceID("FAILED_ROLLBACK_TOPUP_AMOUNT")

			s.logger.Error("Failed to rollback topup amount",
				zap.Int("topup_id", *request.TopupID),
				zap.Error(rollbackErr))
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(rollbackErr)
			span.SetStatus(codes.Error, "Failed to rollback topup amount")
			status = "failed_to_rollback_topup_amount"
		}

		req := requests.UpdateTopupStatus{
			TopupID: *request.TopupID,
			Status:  "failed",
		}

		s.topupCommandRepository.UpdateTopupStatus(&req)
		return nil, saldo_errors.ErrFailedUpdateSaldo
	}

	updatedTopup, err := s.topupQueryRepository.FindById(*request.TopupID)
	if err != nil || updatedTopup == nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TOPUP_BY_ID")

		s.logger.Error("Failed to fetch topup by ID",
			zap.Int("topup_id", *request.TopupID),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch topup by ID")
		status = "failed_to_retrieve_topup_by_id"

		req := requests.UpdateTopupStatus{
			TopupID: *request.TopupID,
			Status:  "failed",
		}

		s.topupCommandRepository.UpdateTopupStatus(&req)
		return nil, topup_errors.ErrTopupNotFoundRes
	}

	req := requests.UpdateTopupStatus{
		TopupID: *request.TopupID,
		Status:  "success",
	}

	_, err = s.topupCommandRepository.UpdateTopupStatus(&req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_TOPUP_STATUS")

		s.logger.Error("Failed to update topup status",
			zap.Int("topup_id", *request.TopupID),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update topup status")
		status = "failed_to_update_topup_status"

		return nil, topup_errors.ErrFailedUpdateTopup
	}

	so := s.mapping.ToTopupResponse(updatedTopup)

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
		traceID := traceunic.GenerateTraceID("FAILED_TRASH_TOPUP")

		s.logger.Error("Failed to trash topup",
			zap.Int("topup_id", topup_id),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to trash topup")
		status = "failed_to_trash_topup"

		return nil, topup_errors.ErrFailedTrashTopup
	}

	so := s.mapping.ToTopupResponseDeleteAt(res)

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
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_TOPUP")

		s.logger.Error("Failed to restore topup",
			zap.Int("topup_id", topup_id),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore topup")
		status = "failed_to_restore_topup"

		return nil, topup_errors.ErrFailedRestoreTopup
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
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_TOPUP_PERMANENT")

		s.logger.Error("Failed to permanently delete topup",
			zap.Int("topup_id", topup_id),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete topup")
		status = "failed_to_permanently_delete_topup"

		return false, topup_errors.ErrFailedDeleteTopup
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
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ALL_TOPUPS")

		s.logger.Error("Failed to restore all topups", zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore all topups")
		status = "failed_to_restore_all_topups"
		return false, topup_errors.ErrFailedRestoreAllTopups
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
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_ALL_TOPUPS")

		s.logger.Error("Failed to permanently delete all topups", zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete all topups")
		status = "failed_to_permanently_delete_all_topups"

		return false, topup_errors.ErrFailedDeleteAllTopups
	}

	s.logger.Debug("Successfully deleted all topups permanently")
	return true, nil
}

func (s *topupCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
