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
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-withdraw/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type withdrawCommandService struct {
	ctx                       context.Context
	errorhandler              errorhandler.WithdrawCommandErrorHandler
	mencache                  mencache.WithdrawCommandCache
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

func NewWithdrawCommandService(ctx context.Context, errorhandler errorhandler.WithdrawCommandErrorHandler,
	mencache mencache.WithdrawCommandCache, kafka kafka.Kafka,
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
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &withdrawCommandService{
		kafka:                     kafka,
		ctx:                       ctx,
		errorhandler:              errorhandler,
		mencache:                  mencache,
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
		return s.errorhandler.HandleRepositorySingleError(err, "CreateWithdraw", "FAILED_FIND_CARD_BY_CARD_NUMBER", span, &status, card_errors.ErrFailedFindByCardNumber, zap.String("card_number", request.CardNumber))
	}

	saldo, err := s.saldoRepository.FindByCardNumber(request.CardNumber)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "CreateWithdraw", "FAILED_FIND_SALDO_BY_CARD_NUMBER", span, &status, saldo_errors.ErrFailedFindSaldoByCardNumber, zap.String("card_number", request.CardNumber))
	}

	if saldo.TotalBalance < request.WithdrawAmount {
		return s.errorhandler.HandleInsufficientBalanceError(err, "CreateWithdraw", "INSUFFICIENT_BALANCE", span, &status, request.CardNumber, zap.String("card_number", request.CardNumber))
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
		return s.errorhandler.HandleRepositorySingleError(err, "CreateWithdraw", "FAILED_UPDATE_SALDO", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.String("card_number", request.CardNumber))
	}
	withdrawRecord, err := s.withdrawCommandRepository.CreateWithdraw(request)
	if err != nil {
		rollbackData := &requests.UpdateSaldoWithdraw{
			CardNumber:     request.CardNumber,
			TotalBalance:   saldo.TotalBalance,
			WithdrawAmount: &request.WithdrawAmount,
			WithdrawTime:   &request.WithdrawTime,
		}
		if _, rollbackErr := s.saldoRepository.UpdateSaldoWithdraw(rollbackData); rollbackErr != nil {
			return s.errorhandler.HandleRepositorySingleError(rollbackErr, "CreateWithdraw", "FAILED_ROLLBACK_SALDO", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.String("card_number", request.CardNumber))
		}
		if _, err := s.withdrawCommandRepository.UpdateWithdrawStatus(&requests.UpdateWithdrawStatus{
			WithdrawID: withdrawRecord.ID,
			Status:     "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, "CreateWithdraw", "FAILED_UPDATE_WITHDRAW_STATUS", span, &status, withdraw_errors.ErrFailedUpdateWithdraw)
		}

		return s.errorhandler.HandleCreateWithdrawError(err, "CreateWithdraw", "FAILED_CREATE_WITHDRAW", span, &status, zap.String("card_number", request.CardNumber))
	}
	if _, err := s.withdrawCommandRepository.UpdateWithdrawStatus(&requests.UpdateWithdrawStatus{
		WithdrawID: withdrawRecord.ID,
		Status:     "success",
	}); err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "CreateWithdraw", "FAILED_UPDATE_WITHDRAW_STATUS", span, &status, withdraw_errors.ErrFailedUpdateWithdraw)
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
		return errorhandler.HandleErrorMarshal[*response.WithdrawResponse](s.logger, err, "CreateWithdraw", "FAILED_MARSHAL_EMAIL_PAYLOAD", span, &status, withdraw_errors.ErrFailedSendEmail)
	}

	err = s.kafka.SendMessage("email-service-topic-withdraw-create", strconv.Itoa(withdrawRecord.ID), payloadBytes)
	if err != nil {
		return errorhandler.HandleErrorKafkaSend[*response.WithdrawResponse](s.logger, err, "CreateWithdraw", "FAILED_SEND_EMAIL", span, &status, withdraw_errors.ErrFailedSendEmail)
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
		return s.errorhandler.HandleRepositorySingleError(err, "Update", "FAILED_FIND_WITHDRAW", span, &status, withdraw_errors.ErrWithdrawNotFound)
	}

	saldo, err := s.saldoRepository.FindByCardNumber(request.CardNumber)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "Update", "FAILED_FIND_SALDO", span, &status, saldo_errors.ErrFailedSaldoNotFound)
	}

	if saldo.TotalBalance < request.WithdrawAmount {
		return s.errorhandler.HandleInsufficientBalanceError(err, "Update", "FAILED_INSUFFICIENT_BALANCE", span, &status, request.CardNumber, zap.Int("withdraw_amount", request.WithdrawAmount))
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
		if _, err := s.withdrawCommandRepository.UpdateWithdrawStatus(&requests.UpdateWithdrawStatus{
			WithdrawID: *request.WithdrawID,
			Status:     "failed",
		}); err != nil {
			return s.errorhandler.HandleUpdateWithdrawError(err, "Update", "FAILED_UPDATE_WITHDRAW_STATUS", span, &status, zap.String("card_number", request.CardNumber))
		}

		return s.errorhandler.HandleRepositorySingleError(err, "Update", "FAILED_UPDATE_SALDO", span, &status, saldo_errors.ErrFailedUpdateSaldo)
	}

	updatedWithdraw, err := s.withdrawCommandRepository.UpdateWithdraw(request)
	if err != nil {
		rollbackData := &requests.UpdateSaldoBalance{
			CardNumber:   saldo.CardNumber,
			TotalBalance: saldo.TotalBalance,
		}
		_, rollbackErr := s.saldoRepository.UpdateSaldoBalance(rollbackData)
		if rollbackErr != nil {
			return s.errorhandler.HandleUpdateWithdrawError(err, "Update", "FAILED_ROLLBACK_SALDO", span, &status, zap.String("card_number", request.CardNumber))
		}
		if _, err := s.withdrawCommandRepository.UpdateWithdrawStatus(&requests.UpdateWithdrawStatus{
			WithdrawID: *request.WithdrawID,
			Status:     "failed",
		}); err != nil {
			return s.errorhandler.HandleUpdateWithdrawError(err, "Update", "FAILED_UPDATE_WITHDRAW_STATUS", span, &status, zap.String("card_number", request.CardNumber))
		}

		return s.errorhandler.HandleUpdateWithdrawError(err, "Update", "FAILED_UPDATE_WITHDRAW", span, &status, zap.String("card_number", request.CardNumber))
	}

	if _, err := s.withdrawCommandRepository.UpdateWithdrawStatus(&requests.UpdateWithdrawStatus{
		WithdrawID: updatedWithdraw.ID,
		Status:     "success",
	}); err != nil {
		return s.errorhandler.HandleUpdateWithdrawError(err, "Update", "FAILED_UPDATE_WITHDRAW_STATUS", span, &status, zap.String("card_number", request.CardNumber))
	}

	so := s.mapping.ToWithdrawResponse(updatedWithdraw)

	s.mencache.DeleteCachedWithdrawCache(*request.WithdrawID)

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
		return s.errorhandler.HandleTrashedWithdrawError(err, "TrashedWithdraw", "FAILED_TRASHED_WITHDRAW", span, &status, zap.Int("withdraw_id", withdraw_id))
	}

	withdrawResponse := s.mapping.ToWithdrawResponse(res)

	s.mencache.DeleteCachedWithdrawCache(withdraw_id)

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
		return s.errorhandler.HandleRestoreWithdrawError(err, "RestoreWithdraw", "FAILED_RESTORE_WITHDRAW", span, &status, zap.Int("withdraw_id", withdraw_id))
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
		return s.errorhandler.HandleDeleteWithdrawPermanentError(err, "DeleteWithdrawPermanent", "FAILED_DELETE_WITHDRAW_PERMANENT", span, &status, zap.Int("withdraw_id", withdraw_id))
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
		return s.errorhandler.HandleRestoreAllWithdrawError(err, "RestoreAllWithdraw", "FAILED_RESTORE_ALL_WITHDRAW", span, &status, zap.Error(err))
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
		return s.errorhandler.HandleDeleteAllWithdrawPermanentError(err, "DeleteAllWithdrawPermanent", "FAILED_DELETE_ALL_WITHDRAW_PERMANENT", span, &status, zap.Error(err))
	}

	s.logger.Debug("Successfully deleted all withdraws permanently")
	return true, nil
}

func (s *withdrawCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
