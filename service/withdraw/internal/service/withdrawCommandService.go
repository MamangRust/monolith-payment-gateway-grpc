package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type withdrawCommandService struct {
	ctx                       context.Context
	kafka                     kafka.Kafka
	trace                     trace.Tracer
	cardRepository            repository.CardRepository
	saldoRepository           repository.SaldoRepository
	withdrawQueryRepository   repository.WithdrawQueryRepository
	withdrawCommandRepository repository.WithdrawCommandRepository
	logger                    logger.LoggerInterface
	mapping                   responseservice.WithdrawResponseMapper
	requestCounter            *prometheus.CounterVec
	requestDuration           *prometheus.HistogramVec
}

func NewWithdrawCommandService(ctx context.Context, kafka kafka.Kafka,
	cardRepository repository.CardRepository,
	saldoRepository repository.SaldoRepository,
	withdrawCommandRepository repository.WithdrawCommandRepository,
	withdrawQueryRepository repository.WithdrawQueryRepository, logger logger.LoggerInterface, mapping responseservice.WithdrawResponseMapper) *withdrawCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "withdraw_command_service_request_total",
			Help: "Total number of requests to the WithdrawCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "withdraw_command_service_request_duration_seconds",
			Help:    "Histogram of request durations for the WithdrawCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &withdrawCommandService{
		kafka:                     kafka,
		ctx:                       ctx,
		trace:                     otel.Tracer("withdraw-command-service"),
		cardRepository:            cardRepository,
		saldoRepository:           saldoRepository,
		withdrawCommandRepository: withdrawCommandRepository,
		withdrawQueryRepository:   withdrawQueryRepository,
		logger:                    logger,
		mapping:                   mapping,
		requestCounter:            requestCounter,
		requestDuration:           requestDuration,
	}
}

func (s *withdrawCommandService) Create(request *requests.CreateWithdrawRequest) (*response.WithdrawResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("CreateWithdraw", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "CreateWithdraw")

	defer span.End()

	s.logger.Debug("Creating new withdraw", zap.Any("request", request))

	span.SetAttributes(
		attribute.String("card_number", request.CardNumber),
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
		status = "failed_find_card_by_card_number"

		return nil, card_errors.ErrFailedFindByCardNumber
	}

	saldo, err := s.saldoRepository.FindByCardNumber(request.CardNumber)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_SALDO_BY_CARD_NUMBER")

		s.logger.Error("Failed to find saldo by card number", zap.Error(err), zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find saldo by card number")
		status = "failed_find_saldo_by_card_number"

		return nil, saldo_errors.ErrFailedSaldoNotFound
	}

	if saldo == nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_SALDO_BY_CARD_NUMBER")

		s.logger.Error("Saldo not found for the specified user ID", zap.String("cardNumber", request.CardNumber), zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Saldo not found for the specified user ID")
		status = "failed_find_saldo_by_card_number"

		return nil, saldo_errors.ErrFailedSaldoNotFound
	}
	if saldo.TotalBalance < request.WithdrawAmount {
		traceID := traceunic.GenerateTraceID("INSUFFICIENT_BALANCE_FOR_WITHDRAWAL")

		s.logger.Error("Insufficient balance for withdrawal", zap.String("cardNumber", request.CardNumber), zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Insufficient balance for withdrawal")
		status = "insufficient_balance_for_withdrawal"

		return nil, saldo_errors.ErrFailedInsuffientBalance
	}
	newTotalBalance := saldo.TotalBalance - request.WithdrawAmount
	updateData := &requests.UpdateSaldoWithdraw{
		CardNumber:     request.CardNumber,
		TotalBalance:   newTotalBalance,
		WithdrawAmount: &request.WithdrawAmount,
		WithdrawTime:   &request.WithdrawTime,
	}
	_, err = s.saldoRepository.UpdateSaldoWithdraw(updateData)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_SALDO")

		s.logger.Error("Failed to update saldo", zap.Error(err), zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update saldo")
		status = "failed_update_saldo"

		return nil, saldo_errors.ErrFailedUpdateSaldo
	}
	withdrawRecord, err := s.withdrawCommandRepository.CreateWithdraw(request)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_CREATE_WITHDRAW")

		s.logger.Error("Failed to create withdraw", zap.Error(err), zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create withdraw")
		status = "failed_create_withdraw"

		rollbackData := &requests.UpdateSaldoWithdraw{
			CardNumber:     request.CardNumber,
			TotalBalance:   saldo.TotalBalance,
			WithdrawAmount: &request.WithdrawAmount,
			WithdrawTime:   &request.WithdrawTime,
		}
		if _, rollbackErr := s.saldoRepository.UpdateSaldoWithdraw(rollbackData); rollbackErr != nil {
			traceID := traceunic.GenerateTraceID("FAILED_ROLLBACK_SALDO")

			s.logger.Error("Failed to rollback saldo", zap.Error(rollbackErr), zap.String("traceID", traceID))

			span.SetAttributes(
				attribute.String("traceID", traceID),
			)

			span.RecordError(rollbackErr)
			span.SetStatus(codes.Error, "Failed to rollback saldo")
			status = "failed_rollback_saldo"

			return nil, saldo_errors.ErrFailedUpdateSaldo
		}
		if _, err := s.withdrawCommandRepository.UpdateWithdrawStatus(&requests.UpdateWithdrawStatus{
			WithdrawID: withdrawRecord.ID,
			Status:     "failed",
		}); err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_UPDATE_WITHDRAW_STATUS")

			s.logger.Error("Failed to update withdraw status", zap.Error(err), zap.String("traceID", traceID))

			span.SetAttributes(
				attribute.String("traceID", traceID),
			)

			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to update withdraw status")
			status = "failed_update_withdraw_status"

			return nil, withdraw_errors.ErrFailedUpdateWithdraw
		}
		return nil, withdraw_errors.ErrFailedCreateWithdraw
	}
	if _, err := s.withdrawCommandRepository.UpdateWithdrawStatus(&requests.UpdateWithdrawStatus{
		WithdrawID: withdrawRecord.ID,
		Status:     "success",
	}); err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_WITHDRAW_STATUS")

		s.logger.Error("Failed to update withdraw status", zap.Error(err), zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update withdraw status")
		status = "failed_update_withdraw_status"

		return nil, withdraw_errors.ErrFailedUpdateWithdraw
	}

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   "Withdraw Successful",
		"Message": fmt.Sprintf("Your withdrawal of %d has been processed successfully.", request.WithdrawAmount),
		"Button":  "View History",
		"Link":    "https://sanedge.example.com/withdraw/history",
	})

	emailPayload := map[string]any{
		"email":   card.Email,
		"subject": "Withdraw Successful - SanEdge",
		"body":    htmlBody,
	}

	payloadBytes, err := json.Marshal(emailPayload)
	if err != nil {
		traceID := traceunic.GenerateTraceID("WITHDRAW_ERR")
		s.logger.Error("Failed to marshal withdraw email payload", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal withdraw email payload")
		return nil, withdraw_errors.ErrFailedSendEmail
	}

	err = s.kafka.SendMessage("email-service-topic-withdraw-create", strconv.Itoa(withdrawRecord.ID), payloadBytes)
	if err != nil {
		traceID := traceunic.GenerateTraceID("WITHDRAW_ERR")
		s.logger.Error("Failed to send withdraw email message", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to send withdraw email")
		return nil, withdraw_errors.ErrFailedSendEmail
	}

	so := s.mapping.ToWithdrawResponse(withdrawRecord)
	s.logger.Debug("Successfully created withdraw", zap.Int("withdraw_id", withdrawRecord.ID))
	return so, nil
}

func (s *withdrawCommandService) Update(request *requests.UpdateWithdrawRequest) (*response.WithdrawResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("Update", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "Update")
	defer span.End()

	s.logger.Debug("Updating withdraw", zap.Int("withdraw_id", *request.WithdrawID), zap.Any("request", request))

	_, err := s.withdrawQueryRepository.FindById(*request.WithdrawID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("WITHDRAW_NOT_FOUND")

		s.logger.Error("Withdraw not found", zap.Error(err), zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Withdraw not found")
		status = "withdraw_not_found"

		return nil, withdraw_errors.ErrWithdrawNotFound
	}
	saldo, err := s.saldoRepository.FindByCardNumber(request.CardNumber)
	if err != nil {
		traceID := traceunic.GenerateTraceID("SALDO_NOT_FOUND")

		s.logger.Error("Saldo not found", zap.Error(err), zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Saldo not found")
		status = "saldo_not_found"

		return nil, saldo_errors.ErrFailedSaldoNotFound
	}
	if saldo.TotalBalance < request.WithdrawAmount {
		traceID := traceunic.GenerateTraceID("INSUFFICIENT_BALANCE")

		s.logger.Error("Insufficient balance for withdrawal update", zap.Error(err), zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Insufficient balance for withdrawal update")
		status = "insufficient_balance"

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Insufficient balance for withdrawal update.",
			Code:    http.StatusBadRequest,
		}
	}
	newTotalBalance := saldo.TotalBalance - request.WithdrawAmount
	updateSaldoData := &requests.UpdateSaldoWithdraw{
		CardNumber:     saldo.CardNumber,
		TotalBalance:   newTotalBalance,
		WithdrawAmount: &request.WithdrawAmount,
		WithdrawTime:   &request.WithdrawTime,
	}
	_, err = s.saldoRepository.UpdateSaldoWithdraw(updateSaldoData)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_SALDO")

		s.logger.Error("Failed to update saldo", zap.Error(err), zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update saldo")
		status = "failed_update_saldo"

		if _, err := s.withdrawCommandRepository.UpdateWithdrawStatus(&requests.UpdateWithdrawStatus{
			WithdrawID: *request.WithdrawID,
			Status:     "failed",
		}); err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_UPDATE_WITHDRAW_STATUS")

			s.logger.Error("Failed to update withdraw status", zap.Error(err), zap.String("traceID", traceID))

			span.SetAttributes(
				attribute.String("traceID", traceID),
			)

			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to update withdraw status")
			status = "failed_update_withdraw_status"

			return nil, withdraw_errors.ErrFailedUpdateWithdraw
		}
		return nil, withdraw_errors.ErrFailedUpdateWithdraw
	}
	updatedWithdraw, err := s.withdrawCommandRepository.UpdateWithdraw(request)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_WITHDRAW")

		s.logger.Error("Failed to update withdraw", zap.Error(err), zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update withdraw")
		status = "failed_update_withdraw"

		rollbackData := &requests.UpdateSaldoBalance{
			CardNumber:   saldo.CardNumber,
			TotalBalance: saldo.TotalBalance,
		}
		_, rollbackErr := s.saldoRepository.UpdateSaldoBalance(rollbackData)
		if rollbackErr != nil {
			traceID := traceunic.GenerateTraceID("FAILED_ROLLBACK_SALDO")

			s.logger.Error("Failed to rollback saldo", zap.Error(err), zap.String("traceID", traceID))

			span.SetAttributes(
				attribute.String("traceID", traceID),
			)

			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to rollback saldo")
			status = "failed_rollback_saldo"

			return nil, saldo_errors.ErrFailedUpdateSaldo
		}
		if _, err := s.withdrawCommandRepository.UpdateWithdrawStatus(&requests.UpdateWithdrawStatus{
			WithdrawID: *request.WithdrawID,
			Status:     "failed",
		}); err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_UPDATE_WITHDRAW_STATUS")

			s.logger.Error("Failed to update withdraw status", zap.Error(err), zap.String("traceID", traceID))

			span.SetAttributes(
				attribute.String("traceID", traceID),
			)

			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to update withdraw status")
			status = "failed_update_withdraw_status"

			return nil, withdraw_errors.ErrFailedUpdateWithdraw
		}
		return nil, withdraw_errors.ErrFailedUpdateWithdraw
	}
	if _, err := s.withdrawCommandRepository.UpdateWithdrawStatus(&requests.UpdateWithdrawStatus{
		WithdrawID: updatedWithdraw.ID,
		Status:     "success",
	}); err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_WITHDRAW_STATUS")

		s.logger.Error("Failed to update withdraw status", zap.Error(err), zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update withdraw status")
		status = "failed_update_withdraw_status"

		return nil, withdraw_errors.ErrFailedUpdateWithdraw
	}
	so := s.mapping.ToWithdrawResponse(updatedWithdraw)
	s.logger.Debug("Successfully updated withdraw", zap.Int("withdraw_id", so.ID))
	return so, nil
}

func (s *withdrawCommandService) TrashedWithdraw(withdraw_id int) (*response.WithdrawResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("TrashedWithdraw", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "TrashedWithdraw")
	defer span.End()

	span.SetAttributes(
		attribute.Int("withdraw_id", withdraw_id),
	)

	s.logger.Debug("Trashing withdraw", zap.Int("withdraw_id", withdraw_id))

	res, err := s.withdrawCommandRepository.TrashedWithdraw(withdraw_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_TRASHED_WITHDRAW")

		s.logger.Error("Failed to trashed withdraw", zap.Error(err), zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to trashed withdraw")
		status = "failed_trashed_withdraw"

		return nil, withdraw_errors.ErrFailedTrashedWithdraw
	}

	withdrawResponse := s.mapping.ToWithdrawResponse(res)

	s.logger.Debug("Successfully trashed withdraw", zap.Int("withdraw_id", withdraw_id))

	return withdrawResponse, nil
}

func (s *withdrawCommandService) RestoreWithdraw(withdraw_id int) (*response.WithdrawResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreWithdraw", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreWithdraw")
	defer span.End()

	span.SetAttributes(
		attribute.Int("withdraw_id", withdraw_id),
	)

	s.logger.Debug("Restoring withdraw", zap.Int("withdraw_id", withdraw_id))

	res, err := s.withdrawCommandRepository.RestoreWithdraw(withdraw_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_WITHDRAW")

		s.logger.Error("Failed to restore withdraw", zap.Error(err), zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore withdraw")
		status = "failed_restore_withdraw"

		return nil, withdraw_errors.ErrFailedRestoreWithdraw
	}

	withdrawResponse := s.mapping.ToWithdrawResponse(res)

	s.logger.Debug("Successfully restored withdraw", zap.Int("withdraw_id", withdraw_id))

	return withdrawResponse, nil
}

func (s *withdrawCommandService) DeleteWithdrawPermanent(withdraw_id int) (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteWithdrawPermanent", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteWithdrawPermanent")
	defer span.End()

	span.SetAttributes(
		attribute.Int("withdraw_id", withdraw_id),
	)

	s.logger.Debug("Deleting withdraw permanently", zap.Int("withdraw_id", withdraw_id))

	_, err := s.withdrawCommandRepository.DeleteWithdrawPermanent(withdraw_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_WITHDRAW_PERMANENT")

		s.logger.Error("Failed to delete withdraw permanently", zap.Error(err), zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to delete withdraw permanently")
		status = "failed_delete_withdraw_permanent"

		return false, withdraw_errors.ErrFailedDeleteWithdrawPermanent
	}

	s.logger.Debug("Successfully deleted withdraw permanently", zap.Int("withdraw_id", withdraw_id))

	return true, nil
}

func (s *withdrawCommandService) RestoreAllWithdraw() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreAllWithdraw", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreAllWithdraw")
	defer span.End()

	s.logger.Debug("Restoring all withdraws")

	_, err := s.withdrawCommandRepository.RestoreAllWithdraw()

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ALL_WITHDRAW")

		s.logger.Error("Failed to restore all withdraws", zap.Error(err), zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore all withdraws")
		status = "failed_restore_all_withdraw"

		return false, withdraw_errors.ErrFailedRestoreAllWithdraw
	}

	s.logger.Debug("Successfully restored all withdraws")
	return true, nil
}

func (s *withdrawCommandService) DeleteAllWithdrawPermanent() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteAllWithdrawPermanent", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteAllWithdrawPermanent")
	defer span.End()

	s.logger.Debug("Permanently deleting all withdraws")

	_, err := s.withdrawCommandRepository.DeleteAllWithdrawPermanent()

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_ALL_WITHDRAW_PERMANENT")

		s.logger.Error("Failed to delete all withdraws permanently", zap.Error(err), zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to delete all withdraws permanently")
		status = "failed_delete_all_withdraw_permanent"

		return false, withdraw_errors.ErrFailedDeleteAllWithdrawPermanent
	}

	s.logger.Debug("Successfully deleted all withdraws permanently")
	return true, nil
}

func (s *withdrawCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
