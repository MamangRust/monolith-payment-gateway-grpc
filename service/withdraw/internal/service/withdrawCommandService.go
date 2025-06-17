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
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type withdrawCommandService struct {
	ctx                       context.Context
	errorhandler              errorhandler.WithdrawCommandErrorHandler
	mencache                  mencache.WithdrawCommandCache
	kafka                     *kafka.Kafka
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
	mencache mencache.WithdrawCommandCache, kafka *kafka.Kafka,
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
	const method = "CreateWithdraw"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	card, err := s.cardRepository.FindUserCardByCardNumber(request.CardNumber)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_CARD_BY_CARD_NUMBER", span, &status, card_errors.ErrFailedFindByCardNumber, zap.Error(err))
	}

	saldo, err := s.saldoRepository.FindByCardNumber(request.CardNumber)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_SALDO_BY_CARD_NUMBER", span, &status, saldo_errors.ErrFailedFindSaldoByCardNumber, zap.Error(err))
	}

	if saldo.TotalBalance < request.WithdrawAmount {
		return s.errorhandler.HandleInsufficientBalanceError(err, method, "INSUFFICIENT_BALANCE", span, &status, request.CardNumber, zap.Error(err))
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
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_SALDO", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
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
			return s.errorhandler.HandleRepositorySingleError(rollbackErr, method, "FAILED_ROLLBACK_SALDO", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
		}
		if _, err := s.withdrawCommandRepository.UpdateWithdrawStatus(&requests.UpdateWithdrawStatus{
			WithdrawID: withdrawRecord.ID,
			Status:     "failed",
		}); err != nil {
			return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_WITHDRAW_STATUS", span, &status, withdraw_errors.ErrFailedUpdateWithdraw, zap.Error(err))
		}

		return s.errorhandler.HandleCreateWithdrawError(err, method, "FAILED_CREATE_WITHDRAW", span, &status, zap.Error(err))
	}
	if _, err := s.withdrawCommandRepository.UpdateWithdrawStatus(&requests.UpdateWithdrawStatus{
		WithdrawID: withdrawRecord.ID,
		Status:     "success",
	}); err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_WITHDRAW_STATUS", span, &status, withdraw_errors.ErrFailedUpdateWithdraw, zap.Error(err))
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
		return errorhandler.HandleErrorMarshal[*response.WithdrawResponse](s.logger, err, method, "FAILED_MARSHAL_EMAIL_PAYLOAD", span, &status, withdraw_errors.ErrFailedSendEmail, zap.Error(err))
	}

	err = s.kafka.SendMessage("email-service-topic-withdraw-create", strconv.Itoa(withdrawRecord.ID), payloadBytes)
	if err != nil {
		return errorhandler.HandleErrorKafkaSend[*response.WithdrawResponse](s.logger, err, method, "FAILED_SEND_EMAIL", span, &status, withdraw_errors.ErrFailedSendEmail, zap.Error(err))
	}

	so := s.mapping.ToWithdrawResponse(withdrawRecord)

	logSuccess("Successfully created withdraw", zap.Int("withdraw.id", withdrawRecord.ID))

	return so, nil
}

func (s *withdrawCommandService) Update(request *requests.UpdateWithdrawRequest) (*response.WithdrawResponse, *response.ErrorResponse) {
	const method = "UpdateWithdraw"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	_, err := s.withdrawQueryRepository.FindById(*request.WithdrawID)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_WITHDRAW", span, &status, withdraw_errors.ErrWithdrawNotFound)
	}

	saldo, err := s.saldoRepository.FindByCardNumber(request.CardNumber)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_SALDO", span, &status, saldo_errors.ErrFailedSaldoNotFound, zap.Error(err))
	}

	if saldo.TotalBalance < request.WithdrawAmount {
		return s.errorhandler.HandleInsufficientBalanceError(err, method, "FAILED_INSUFFICIENT_BALANCE", span, &status, request.CardNumber, zap.Error(err))
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
			return s.errorhandler.HandleUpdateWithdrawError(err, method, "FAILED_UPDATE_WITHDRAW_STATUS", span, &status, zap.Error(err))
		}

		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_UPDATE_SALDO", span, &status, saldo_errors.ErrFailedUpdateSaldo, zap.Error(err))
	}

	updatedWithdraw, err := s.withdrawCommandRepository.UpdateWithdraw(request)
	if err != nil {
		rollbackData := &requests.UpdateSaldoBalance{
			CardNumber:   saldo.CardNumber,
			TotalBalance: saldo.TotalBalance,
		}
		_, rollbackErr := s.saldoRepository.UpdateSaldoBalance(rollbackData)
		if rollbackErr != nil {
			return s.errorhandler.HandleUpdateWithdrawError(err, method, "FAILED_ROLLBACK_SALDO", span, &status, zap.Error(err))
		}
		if _, err := s.withdrawCommandRepository.UpdateWithdrawStatus(&requests.UpdateWithdrawStatus{
			WithdrawID: *request.WithdrawID,
			Status:     "failed",
		}); err != nil {
			return s.errorhandler.HandleUpdateWithdrawError(err, method, "FAILED_UPDATE_WITHDRAW_STATUS", span, &status, zap.Error(err))
		}

		return s.errorhandler.HandleUpdateWithdrawError(err, method, "FAILED_UPDATE_WITHDRAW", span, &status, zap.Error(err))
	}

	if _, err := s.withdrawCommandRepository.UpdateWithdrawStatus(&requests.UpdateWithdrawStatus{
		WithdrawID: updatedWithdraw.ID,
		Status:     "success",
	}); err != nil {
		return s.errorhandler.HandleUpdateWithdrawError(err, method, "FAILED_UPDATE_WITHDRAW_STATUS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToWithdrawResponse(updatedWithdraw)

	s.mencache.DeleteCachedWithdrawCache(*request.WithdrawID)

	logSuccess("Successfully updated withdraw", zap.Int("withdraw.id", updatedWithdraw.ID))

	return so, nil
}

func (s *withdrawCommandService) TrashedWithdraw(withdraw_id int) (*response.WithdrawResponse, *response.ErrorResponse) {
	const method = "TrashedWithdraw"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	res, err := s.withdrawCommandRepository.TrashedWithdraw(withdraw_id)

	if err != nil {
		return s.errorhandler.HandleTrashedWithdrawError(err, method, "FAILED_TRASHED_WITHDRAW", span, &status, zap.Error(err))
	}

	withdrawResponse := s.mapping.ToWithdrawResponse(res)

	s.mencache.DeleteCachedWithdrawCache(withdraw_id)

	logSuccess("Successfully trashed withdraw", zap.Int("withdraw.id", withdraw_id))

	return withdrawResponse, nil
}

func (s *withdrawCommandService) RestoreWithdraw(withdraw_id int) (*response.WithdrawResponse, *response.ErrorResponse) {
	const method = "RestoreWithdraw"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	res, err := s.withdrawCommandRepository.RestoreWithdraw(withdraw_id)

	if err != nil {
		return s.errorhandler.HandleRestoreWithdrawError(err, method, "FAILED_RESTORE_WITHDRAW", span, &status, zap.Error(err))
	}

	withdrawResponse := s.mapping.ToWithdrawResponse(res)

	logSuccess("Successfully restored withdraw", zap.Int("withdraw.id", withdraw_id))

	return withdrawResponse, nil
}

func (s *withdrawCommandService) DeleteWithdrawPermanent(withdraw_id int) (bool, *response.ErrorResponse) {
	const method = "DeleteWithdrawPermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	_, err := s.withdrawCommandRepository.DeleteWithdrawPermanent(withdraw_id)

	if err != nil {
		return s.errorhandler.HandleDeleteWithdrawPermanentError(err, method, "FAILED_DELETE_WITHDRAW_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted withdraw permanent", zap.Int("withdraw.id", withdraw_id))

	return true, nil
}

func (s *withdrawCommandService) RestoreAllWithdraw() (bool, *response.ErrorResponse) {
	const method = "RestoreAllWithdraw"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	_, err := s.withdrawCommandRepository.RestoreAllWithdraw()

	if err != nil {
		return s.errorhandler.HandleRestoreAllWithdrawError(err, method, "FAILED_RESTORE_ALL_WITHDRAW", span, &status, zap.Error(err))
	}

	logSuccess("Successfully restored all withdraws", zap.Bool("success", true))

	return true, nil
}

func (s *withdrawCommandService) DeleteAllWithdrawPermanent() (bool, *response.ErrorResponse) {
	const method = "DeleteAllWithdrawPermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	_, err := s.withdrawCommandRepository.DeleteAllWithdrawPermanent()

	if err != nil {
		return s.errorhandler.HandleDeleteAllWithdrawPermanentError(err, method, "FAILED_DELETE_ALL_WITHDRAW_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted all withdraw permanent", zap.Bool("success", true))

	return true, nil
}

func (s *withdrawCommandService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *withdrawCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
