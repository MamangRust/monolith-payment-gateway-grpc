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
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transferCommandService struct {
	kafka                     kafka.Kafka
	ctx                       context.Context
	trace                     trace.Tracer
	cardRepository            repository.CardRepository
	saldoRepository           repository.SaldoRepository
	transferQueryRepository   repository.TransferQueryRepository
	transferCommandRepository repository.TransferCommandRepository
	logger                    logger.LoggerInterface
	mapping                   responseservice.TransferResponseMapper
	requestCounter            *prometheus.CounterVec
	requestDuration           *prometheus.HistogramVec
}

func NewTransferCommandService(
	kafka kafka.Kafka,
	ctx context.Context,
	cardRepository repository.CardRepository,
	saldoRepository repository.SaldoRepository,
	transferQueryRepository repository.TransferQueryRepository,
	transferCommandRepository repository.TransferCommandRepository,
	logger logger.LoggerInterface,
	mapping responseservice.TransferResponseMapper,
) *transferCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transfer_command_service_request_total",
			Help: "Total number of requests to the TransferCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transfer_command_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TransferCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &transferCommandService{
		kafka:                     kafka,
		ctx:                       ctx,
		trace:                     otel.Tracer("transfer-command-service"),
		cardRepository:            cardRepository,
		saldoRepository:           saldoRepository,
		transferQueryRepository:   transferQueryRepository,
		transferCommandRepository: transferCommandRepository,
		logger:                    logger,
		mapping:                   mapping,
		requestCounter:            requestCounter,
		requestDuration:           requestDuration,
	}
}

func (s *transferCommandService) CreateTransaction(request *requests.CreateTransferRequest) (*response.TransferResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("CreateTransaction", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "CreateTransaction")
	defer span.End()

	span.SetAttributes(
		attribute.String("transfer_from", request.TransferFrom),
		attribute.String("transfer_to", request.TransferTo),
		attribute.Float64("transfer_amount", float64(request.TransferAmount)),
	)

	s.logger.Debug("Starting create transaction process",
		zap.Any("request", request),
	)

	card, err := s.cardRepository.FindUserCardByCardNumber(request.TransferFrom)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_CARD_BY_CARD_NUMBER")

		s.logger.Error("Card not found for card number",
			zap.String("card_number", request.TransferFrom),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Card not found for card number")
		status = "card_not_found_for_card_number"

		return nil, card_errors.ErrCardNotFoundRes
	}

	_, err = s.cardRepository.FindCardByCardNumber(request.TransferTo)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_CARD_BY_CARD_NUMBER")

		s.logger.Error("Card not found for card number",
			zap.String("card_number", request.TransferTo),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Card not found for card number")
		status = "card_not_found_for_card_number"

		return nil, card_errors.ErrCardNotFoundRes
	}

	senderSaldo, err := s.saldoRepository.FindByCardNumber(request.TransferFrom)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_SALDO_BY_CARD_NUMBER")

		s.logger.Error("Failed to find sender's saldo by card number",
			zap.String("card_number", request.TransferFrom),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find sender's saldo by card number")
		status = "failed_to_find_sender_saldo_by_card_number"

		return nil, saldo_errors.ErrFailedSaldoNotFound
	}

	receiverSaldo, err := s.saldoRepository.FindByCardNumber(request.TransferTo)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_SALDO_BY_CARD_NUMBER")

		s.logger.Error("Failed to find receiver's saldo by card number",
			zap.String("card_number", request.TransferTo),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find receiver's saldo by card number")
		status = "failed_to_find_receiver_saldo_by_card_number"

		return nil, saldo_errors.ErrFailedSaldoNotFound
	}

	if senderSaldo.TotalBalance < request.TransferAmount {
		traceID := traceunic.GenerateTraceID("INSUFFICIENT_BALANCE_FOR_SENDER")

		s.logger.Error("Insufficient balance for sender",
			zap.String("card_number", request.TransferFrom),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.SetStatus(codes.Error, "Insufficient balance for sender")
		status = "insufficient_balance_for_sender"

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Insufficient balance for sender",
			Code:    http.StatusBadRequest,
		}
	}

	senderSaldo.TotalBalance -= request.TransferAmount
	_, err = s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   senderSaldo.CardNumber,
		TotalBalance: senderSaldo.TotalBalance,
	})
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_SALDO_BALANCE")

		s.logger.Error("Failed to update sender saldo",
			zap.String("card_number", request.TransferFrom),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update sender saldo")
		status = "failed_to_update_sender_saldo"

		return nil, saldo_errors.ErrFailedUpdateSaldo
	}

	receiverSaldo.TotalBalance += request.TransferAmount
	_, err = s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   receiverSaldo.CardNumber,
		TotalBalance: receiverSaldo.TotalBalance,
	})
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_SALDO_BALANCE")

		s.logger.Error("Failed to update receiver saldo",
			zap.String("card_number", request.TransferTo),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update receiver saldo")
		status = "failed_to_update_receiver_saldo"

		senderSaldo.TotalBalance += request.TransferAmount
		_, rollbackErr := s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
			CardNumber:   senderSaldo.CardNumber,
			TotalBalance: senderSaldo.TotalBalance,
		})
		if rollbackErr != nil {
			traceID := traceunic.GenerateTraceID("FAILED_ROLLBACK_SALDO_BALANCE")

			s.logger.Error("Failed to rollback sender saldo",
				zap.String("card_number", request.TransferFrom),
				zap.Error(rollbackErr))
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(rollbackErr)
			span.SetStatus(codes.Error, "Failed to rollback sender saldo")
			status = "failed_to_rollback_sender_saldo"
		}

		return nil, saldo_errors.ErrFailedUpdateSaldo
	}

	transfer, err := s.transferCommandRepository.CreateTransfer(request)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_CREATE_TRANSFER")

		s.logger.Error("Failed to create transfer",
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create transfer")
		status = "failed_to_create_transfer"

		senderSaldo.TotalBalance += request.TransferAmount
		receiverSaldo.TotalBalance -= request.TransferAmount

		_, rollbackErr := s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
			CardNumber:   senderSaldo.CardNumber,
			TotalBalance: senderSaldo.TotalBalance,
		})
		if rollbackErr != nil {
			traceID := traceunic.GenerateTraceID("FAILED_ROLLBACK_SALDO_BALANCE")

			s.logger.Error("Failed to rollback sender saldo",
				zap.String("card_number", request.TransferFrom),
				zap.Error(rollbackErr))
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(rollbackErr)
			span.SetStatus(codes.Error, "Failed to rollback sender saldo")
			status = "failed_to_rollback_sender_saldo"

			return nil, saldo_errors.ErrFailedUpdateSaldo
		}

		_, rollbackErr = s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
			CardNumber:   receiverSaldo.CardNumber,
			TotalBalance: receiverSaldo.TotalBalance,
		})
		if rollbackErr != nil {
			traceID := traceunic.GenerateTraceID("FAILED_ROLLBACK_SALDO_BALANCE")

			s.logger.Error("Failed to rollback receiver saldo",
				zap.String("card_number", request.TransferTo),
				zap.Error(rollbackErr))
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(rollbackErr)
			span.SetStatus(codes.Error, "Failed to rollback receiver saldo")
			status = "failed_to_rollback_receiver_saldo"

			return nil, saldo_errors.ErrFailedUpdateSaldo
		}

		if _, err := s.transferCommandRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
			TransferID: transfer.ID,
			Status:     "failed",
		}); err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_UPDATE_TRANSFER_STATUS")

			s.logger.Error("Failed to update transfer status",
				zap.Error(err))
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to update transfer status")
			status = "failed_to_update_transfer_status"

			return nil, transfer_errors.ErrFailedUpdateTransfer
		}

		return nil, transfer_errors.ErrFailedCreateTransfer
	}

	res, err := s.transferCommandRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
		TransferID: transfer.ID,
		Status:     "success",
	})

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_TRANSFER_STATUS")

		s.logger.Error("Failed to update transfer status",
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update transfer status")
		status = "failed_to_update_transfer_status"

		return nil, transfer_errors.ErrFailedUpdateTransfer
	}

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   "Transfer Successful",
		"Message": fmt.Sprintf("Your Transfer of %d has been processed successfully.", request.TransferAmount),
		"Button":  "View History",
		"Link":    "https://sanedge.example.com/withdraw/history",
	})

	emailPayload := map[string]any{
		"email":   card.Email,
		"subject": "Transfer Successful - SanEdge",
		"body":    htmlBody,
	}

	payloadBytes, err := json.Marshal(emailPayload)
	if err != nil {
		traceID := traceunic.GenerateTraceID("TRANSFER_ERR")
		s.logger.Error("Failed to marshal transfer email payload", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal transfer email payload")
		return nil, withdraw_errors.ErrFailedSendEmail
	}

	err = s.kafka.SendMessage("email-service-topic-transfer-create", strconv.Itoa(res.ID), payloadBytes)
	if err != nil {
		traceID := traceunic.GenerateTraceID("TrANSFER_ERR")
		s.logger.Error("Failed to send transfer email message", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to send transfer email")
		return nil, withdraw_errors.ErrFailedSendEmail
	}

	so := s.mapping.ToTransferResponse(transfer)

	s.logger.Debug("successfully create transaction",
		zap.Int("transfer_id", transfer.ID),
	)

	return so, nil
}

func (s *transferCommandService) UpdateTransaction(request *requests.UpdateTransferRequest) (*response.TransferResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("UpdateTransaction", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateTransaction")
	defer span.End()

	span.SetAttributes(
		attribute.Int("transfer_id", *request.TransferID),
	)

	s.logger.Debug("Starting update transaction process",
		zap.Int("transfer_id", *request.TransferID),
	)

	transfer, err := s.transferQueryRepository.FindById(*request.TransferID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TRANSFER_BY_ID")

		s.logger.Error("Failed to find transfer by ID",
			zap.Int("transfer_id", *request.TransferID),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find transfer by ID")
		status = "failed_to_find_transfer_by_id"

		return nil, transfer_errors.ErrTransferNotFound
	}

	amountDifference := request.TransferAmount - transfer.TransferAmount

	senderSaldo, err := s.saldoRepository.FindByCardNumber(transfer.TransferFrom)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_SALDO_BY_CARD_NUMBER")

		s.logger.Error("Failed to find saldo by card number",
			zap.String("card_number", transfer.TransferFrom),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find saldo by card number")
		status = "failed_to_find_saldo_by_card_number"

		if _, err := s.transferCommandRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
			TransferID: *request.TransferID,
			Status:     "failed",
		}); err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_UPDATE_TRANSFER_STATUS")

			s.logger.Error("Failed to update transfer status",
				zap.Error(err))
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to update transfer status")
			status = "failed_to_update_transfer_status"
		}

		return nil, saldo_errors.ErrFailedSaldoNotFound
	}

	newSenderBalance := senderSaldo.TotalBalance - amountDifference
	if newSenderBalance < 0 {
		traceID := traceunic.GenerateTraceID("INSUFFICIENT_BALANCE")

		s.logger.Error("Insufficient balance for sender",
			zap.String("card_number", senderSaldo.CardNumber),
			zap.Float64("amount_difference", float64(amountDifference)),
			zap.Float64("new_sender_balance", float64(newSenderBalance)),
			zap.Error(err))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Insufficient balance for sender")
		status = "insufficient_balance"

		if _, err := s.transferCommandRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
			TransferID: *request.TransferID,
			Status:     "failed",
		}); err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_UPDATE_TRANSFER_STATUS")

			s.logger.Error("Failed to update transfer status",
				zap.Error(err))
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to update transfer status")
			status = "failed_to_update_transfer_status"

			return nil, transfer_errors.ErrFailedUpdateTransfer
		}

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Insufficient balance for sender",
			Code:    http.StatusBadRequest,
		}
	}

	senderSaldo.TotalBalance = newSenderBalance
	_, err = s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   senderSaldo.CardNumber,
		TotalBalance: senderSaldo.TotalBalance,
	})
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_SALDO_BALANCE")

		s.logger.Error("Failed to update saldo balance",
			zap.String("card_number", senderSaldo.CardNumber),
			zap.Float64("amount_difference", float64(amountDifference)),
			zap.Float64("new_sender_balance", float64(newSenderBalance)),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update saldo balance")
		status = "failed_to_update_saldo_balance"

		if _, err := s.transferCommandRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
			TransferID: *request.TransferID,
			Status:     "failed",
		}); err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_UPDATE_TRANSFER_STATUS")

			s.logger.Error("Failed to update transfer status",
				zap.Error(err))
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to update transfer status")
			status = "failed_to_update_transfer_status"

			return nil, transfer_errors.ErrFailedUpdateTransfer
		}

		return nil, saldo_errors.ErrFailedUpdateSaldo
	}

	receiverSaldo, err := s.saldoRepository.FindByCardNumber(transfer.TransferTo)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_SALDO_BY_CARD_NUMBER")

		s.logger.Error("Failed to find saldo by card number",
			zap.String("card_number", transfer.TransferTo),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find saldo by card number")
		status = "failed_to_find_saldo_by_card_number"

		rollbackSenderBalance := &requests.UpdateSaldoBalance{
			CardNumber:   transfer.TransferFrom,
			TotalBalance: senderSaldo.TotalBalance + amountDifference,
		}
		_, rollbackErr := s.saldoRepository.UpdateSaldoBalance(rollbackSenderBalance)
		if rollbackErr != nil {
			traceID := traceunic.GenerateTraceID("FAILED_UPDATE_SALDO_BALANCE")

			s.logger.Error("Failed to update saldo balance",
				zap.String("card_number", senderSaldo.CardNumber),
				zap.Float64("amount_difference", float64(amountDifference)),
				zap.Float64("new_sender_balance", float64(newSenderBalance)),
				zap.Error(err))
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to update saldo balance")
			status = "failed_to_update_saldo_balance"

			return nil, saldo_errors.ErrFailedUpdateSaldo
		}

		if _, err := s.transferCommandRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
			TransferID: *request.TransferID,
			Status:     "failed",
		}); err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_UPDATE_TRANSFER_STATUS")

			s.logger.Error("Failed to update transfer status",
				zap.Error(err))
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to update transfer status")
			status = "failed_to_update_transfer_status"

			return nil, transfer_errors.ErrFailedUpdateTransfer
		}

		return nil, saldo_errors.ErrFailedSaldoNotFound
	}

	newReceiverBalance := receiverSaldo.TotalBalance + amountDifference
	receiverSaldo.TotalBalance = newReceiverBalance
	_, err = s.saldoRepository.UpdateSaldoBalance(&requests.UpdateSaldoBalance{
		CardNumber:   receiverSaldo.CardNumber,
		TotalBalance: receiverSaldo.TotalBalance,
	})
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_SALDO_BALANCE")

		s.logger.Error("Failed to update saldo balance",
			zap.String("card_number", receiverSaldo.CardNumber),
			zap.Float64("amount_difference", float64(amountDifference)),
			zap.Float64("new_receiver_balance", float64(newReceiverBalance)),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update saldo balance")
		status = "failed_to_update_saldo_balance"

		rollbackSenderBalance := &requests.UpdateSaldoBalance{
			CardNumber:   transfer.TransferFrom,
			TotalBalance: senderSaldo.TotalBalance + amountDifference,
		}
		rollbackReceiverBalance := &requests.UpdateSaldoBalance{
			CardNumber:   transfer.TransferTo,
			TotalBalance: receiverSaldo.TotalBalance - amountDifference,
		}

		if _, err := s.saldoRepository.UpdateSaldoBalance(rollbackSenderBalance); err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_UPDATE_SALDO_BALANCE")

			s.logger.Error("Failed to rollback sender's saldo after receiver update failure",
				zap.Error(err))
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to rollback sender's saldo after receiver update failure")
			status = "failed_to_rollback_sender_saldo_after_receiver_update_failure"

			return nil, saldo_errors.ErrFailedUpdateSaldo
		}
		if _, err := s.saldoRepository.UpdateSaldoBalance(rollbackReceiverBalance); err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_UPDATE_SALDO_BALANCE")

			s.logger.Error("Failed to rollback receiver's saldo after sender update failure",
				zap.Error(err))
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to rollback receiver's saldo after sender update failure")
			status = "failed_to_rollback_receiver_saldo_after_sender_update_failure"

			return nil, saldo_errors.ErrFailedUpdateSaldo
		}

		if _, err := s.transferCommandRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
			TransferID: *request.TransferID,
			Status:     "failed",
		}); err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_UPDATE_TRANSFER_STATUS")

			s.logger.Error("Failed to update transfer status",
				zap.Error(err))
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to update transfer status")
			status = "failed_to_update_transfer_status"

			return nil, transfer_errors.ErrFailedUpdateTransfer
		}

		return nil, saldo_errors.ErrFailedUpdateSaldo
	}

	updatedTransfer, err := s.transferCommandRepository.UpdateTransfer(request)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_TRANSFER")

		s.logger.Error("Failed to update transfer",
			zap.Int("transfer_id", *request.TransferID),
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update transfer")
		status = "failed_to_update_transfer"

		rollbackSenderBalance := &requests.UpdateSaldoBalance{
			CardNumber:   transfer.TransferFrom,
			TotalBalance: senderSaldo.TotalBalance + amountDifference,
		}
		rollbackReceiverBalance := &requests.UpdateSaldoBalance{
			CardNumber:   transfer.TransferTo,
			TotalBalance: receiverSaldo.TotalBalance - amountDifference,
		}

		if _, err := s.saldoRepository.UpdateSaldoBalance(rollbackSenderBalance); err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_UPDATE_SALDO_BALANCE")

			s.logger.Error("Failed to rollback sender's saldo after receiver update failure",
				zap.Error(err))
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to rollback sender's saldo after receiver update failure")
			status = "failed_to_rollback_sender_saldo_after_receiver_update_failure"

			return nil, saldo_errors.ErrFailedUpdateSaldo
		}
		if _, err := s.saldoRepository.UpdateSaldoBalance(rollbackReceiverBalance); err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_UPDATE_SALDO_BALANCE")

			s.logger.Error("Failed to rollback receiver's saldo after sender update failure",
				zap.Error(err))
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to rollback receiver's saldo after sender update failure")
			status = "failed_to_rollback_receiver_saldo_after_sender_update_failure"

			return nil, saldo_errors.ErrFailedUpdateSaldo
		}

		if _, err := s.transferCommandRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
			TransferID: *request.TransferID,
			Status:     "failed",
		}); err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_UPDATE_TRANSFER_STATUS")

			s.logger.Error("Failed to update transfer status",
				zap.Error(err))
			span.SetAttributes(attribute.String("trace.id", traceID))
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to update transfer status")
			status = "failed_to_update_transfer_status"

			return nil, transfer_errors.ErrFailedUpdateTransfer
		}

		return nil, transfer_errors.ErrFailedUpdateTransfer
	}

	if _, err := s.transferCommandRepository.UpdateTransferStatus(&requests.UpdateTransferStatus{
		TransferID: *request.TransferID,
		Status:     "success",
	}); err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_TRANSFER_STATUS")

		s.logger.Error("Failed to update transfer status",
			zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update transfer status")
		status = "failed_to_update_transfer_status"

		return nil, transfer_errors.ErrFailedUpdateTransfer
	}

	so := s.mapping.ToTransferResponse(updatedTransfer)

	s.logger.Debug("successfully update transaction",
		zap.Int("transfer_id", *request.TransferID),
	)

	return so, nil
}

func (s *transferCommandService) TrashedTransfer(transfer_id int) (*response.TransferResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("TrashedTransfer", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "TrashedTransfer")
	defer span.End()

	span.SetAttributes(
		attribute.Int("transfer_id", transfer_id),
	)

	s.logger.Debug("Starting trashed transfer process",
		zap.Int("transfer_id", transfer_id),
	)

	res, err := s.transferCommandRepository.TrashedTransfer(transfer_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_TRASHED_TRANSFER")

		s.logger.Error("Failed to trashed transfer", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to trashed transfer")
		status = "failed_trashed_transfer"

		return nil, transfer_errors.ErrFailedTrashedTransfer
	}

	so := s.mapping.ToTransferResponse(res)

	s.logger.Debug("successfully trashed transfer",
		zap.Int("transfer_id", transfer_id),
	)

	return so, nil
}

func (s *transferCommandService) RestoreTransfer(transfer_id int) (*response.TransferResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreTransfer", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreTransfer")
	defer span.End()

	span.SetAttributes(
		attribute.Int("transfer_id", transfer_id),
	)

	s.logger.Debug("Starting restore transfer process",
		zap.Int("transfer_id", transfer_id),
	)

	res, err := s.transferCommandRepository.RestoreTransfer(transfer_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_TRANSFER")

		s.logger.Error("Failed to restore transfer", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore transfer")
		status = "failed_restore_transfer"

		return nil, transfer_errors.ErrFailedRestoreTransfer
	}

	so := s.mapping.ToTransferResponse(res)

	s.logger.Debug("successfully restore transfer",
		zap.Int("transfer_id", transfer_id),
	)

	return so, nil
}

func (s *transferCommandService) DeleteTransferPermanent(transfer_id int) (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteTransferPermanent", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteTransferPermanent")
	defer span.End()

	span.SetAttributes(
		attribute.Int("transfer_id", transfer_id),
	)

	s.logger.Debug("Starting delete transfer permanent process",
		zap.Int("transfer_id", transfer_id),
	)

	_, err := s.transferCommandRepository.DeleteTransferPermanent(transfer_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_TRANSFER_PERMANENT")

		s.logger.Error("Failed to delete permanent transfer", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to delete permanent transfer")
		status = "failed_delete_transfer_permanent"

		return false, transfer_errors.ErrFailedDeleteTransferPermanent
	}

	s.logger.Debug("successfully delete permanent transfer",
		zap.Int("transfer_id", transfer_id),
	)

	return true, nil
}

func (s *transferCommandService) RestoreAllTransfer() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreAllTransfer", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreAllTransfer")
	defer span.End()

	s.logger.Debug("Restoring all transfers")

	_, err := s.transferCommandRepository.RestoreAllTransfer()

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ALL_TRANSFERS")

		s.logger.Error("Failed to restore all transfers", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore all transfers")
		status = "failed_to_restore_all_transfers"

		return false, transfer_errors.ErrFailedRestoreAllTransfers
	}

	s.logger.Debug("Successfully restored all transfers")

	return true, nil
}

func (s *transferCommandService) DeleteAllTransferPermanent() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteAllTransferPermanent", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteAllTransferPermanent")
	defer span.End()

	s.logger.Debug("Permanently deleting all transfers")

	_, err := s.transferCommandRepository.DeleteAllTransferPermanent()

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_ALL_TRANSFERS_PERMANENT")

		s.logger.Error("Failed to delete all transfers permanently", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to delete all transfers permanently")
		status = "failed_to_delete_all_transfers_permanent"

		return false, transfer_errors.ErrFailedDeleteAllTransfersPermanent
	}

	s.logger.Debug("Successfully deleted all transfers permanently")
	return true, nil
}

func (s *transferCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
