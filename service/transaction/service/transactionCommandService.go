package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/email"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errorhandler"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/service"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"

	mencache "github.com/MamangRust/monolith-payment-gateway-transaction/redis"
	"github.com/MamangRust/monolith-payment-gateway-transaction/repository"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// transactionCommandServiceDeps groups dependencies for transaction command service.
type transactionCommandServiceDeps struct {
	Kafka                        *kafka.Kafka
	Mencache                     mencache.TransactionCommandCache
	Tracer                       trace.Tracer
	MerchantRepository           repository.MerchantRepository
	CardRepository               repository.CardRepository
	SaldoRepository              repository.SaldoRepository
	TransactionQueryRepository   repository.TransactionQueryRepository
	TransactionCommandRepository repository.TransactionCommandRepository
	Logger                       logger.LoggerInterface
	Observability                observability.TraceLoggerObservability
}

// transactionCommandService handles transaction write operations.
type transactionCommandService struct {
	kafka                        *kafka.Kafka
	cache                        mencache.TransactionCommandCache
	merchantRepository           repository.MerchantRepository
	cardRepository               repository.CardRepository
	saldoRepository              repository.SaldoRepository
	transactionQueryRepository   repository.TransactionQueryRepository
	transactionCommandRepository repository.TransactionCommandRepository
	logger                       logger.LoggerInterface
	observability                observability.TraceLoggerObservability
}

func NewTransactionCommandService(
	params *transactionCommandServiceDeps,
) TransactionCommandService {
	return &transactionCommandService{
		kafka:                        params.Kafka,
		cache:                        params.Mencache,
		merchantRepository:           params.MerchantRepository,
		cardRepository:               params.CardRepository,
		saldoRepository:              params.SaldoRepository,
		transactionCommandRepository: params.TransactionCommandRepository,
		transactionQueryRepository:   params.TransactionQueryRepository,
		logger:                       params.Logger,
		observability:                params.Observability,
	}
}

func (s *transactionCommandService) Create(ctx context.Context, apiKey string, request *requests.CreateTransactionRequest) (*db.UpdateTransactionStatusRow, error) {
	const method = "CreateTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("apikey", apiKey))
	defer func() { end(status) }()

	merchant, err := s.merchantRepository.FindByApiKey(ctx, apiKey)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionStatusRow](s.logger, err, method, span, zap.String("api_key", apiKey))
	}

	card, err := s.cardRepository.FindUserCardByCardNumber(ctx, request.CardNumber)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionStatusRow](s.logger, err, method, span, zap.String("card_number", request.CardNumber))
	}

	saldo, err := s.saldoRepository.FindByCardNumber(ctx, card.CardNumber)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionStatusRow](s.logger, err, method, span, zap.String("card_number", card.CardNumber))
	}

	if int(saldo.TotalBalance) < request.Amount {
		status = "error"
		err := errors.New("insufficient balance")
		return errorhandler.HandleError[*db.UpdateTransactionStatusRow](s.logger, err, method, span, zap.Float64("current_balance", float64(saldo.TotalBalance)), zap.Float64("requested_amount", float64(request.Amount)))
	}

	newUserBalance := int(saldo.TotalBalance) - request.Amount
	if _, err := s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
		CardNumber:   card.CardNumber,
		TotalBalance: newUserBalance,
	}); err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionStatusRow](s.logger, err, method, span, zap.String("card_number", card.CardNumber))
	}
	merchantId := int(merchant.MerchantID)

	request.MerchantID = &merchantId
	transaction, err := s.transactionCommandRepository.CreateTransaction(ctx, request)
	if err != nil {
		status = "error"
		if _, rollbackErr := s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
			CardNumber:   card.CardNumber,
			TotalBalance: int(saldo.TotalBalance),
		}); rollbackErr != nil {
			return errorhandler.HandleError[*db.UpdateTransactionStatusRow](s.logger, rollbackErr, method, span, zap.String("card_number", card.CardNumber))
		}
		s.markTransactionAsFailed(ctx, int(transaction.TransactionID), method, span)
		return errorhandler.HandleError[*db.UpdateTransactionStatusRow](s.logger, err, method, span, zap.Int("transaction_id", int(transaction.TransactionID)))
	}

	updatedTransaction, err := s.transactionCommandRepository.UpdateTransactionStatus(ctx, &requests.UpdateTransactionStatus{
		TransactionID: int(transaction.TransactionID),
		Status:        "success",
	})
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionStatusRow](s.logger, err, method, span, zap.Int("transaction_id", int(transaction.TransactionID)))
	}

	merchantCard, err := s.cardRepository.FindCardByUserId(ctx, int(merchant.UserID))
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionStatusRow](s.logger, err, method, span, zap.Int("merchant_id", int(merchant.MerchantID)))
	}

	merchantSaldo, err := s.saldoRepository.FindByCardNumber(ctx, merchantCard.CardNumber)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionStatusRow](s.logger, err, method, span, zap.String("merchant_card_number", merchantCard.CardNumber))
	}

	newMerchantBalance := int(merchantSaldo.TotalBalance) + request.Amount
	if _, err := s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
		CardNumber:   merchantCard.CardNumber,
		TotalBalance: newMerchantBalance,
	}); err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionStatusRow](s.logger, err, method, span, zap.String("merchant_card_number", merchantCard.CardNumber))
	}

	go func() {
		htmlBody := email.GenerateEmailHTML(map[string]string{
			"Title":   "Transaction Successful",
			"Message": fmt.Sprintf("Your transaction of %d has been processed successfully.", request.Amount),
			"Button":  "View History",
			"Link":    "https://sanedge.example.com/transaction/history",
		})

		emailPayload := map[string]any{
			"email":   card.Email,
			"subject": "Transaction Successful - SanEdge",
			"body":    htmlBody,
		}

		payloadBytes, err := json.Marshal(emailPayload)
		if err != nil {
			s.logger.Error("failed to marshal email payload for transaction", zap.Error(err), zap.Int("transaction_id", int(updatedTransaction.TransactionID)))
			return
		}

		if s.kafka != nil {
			err = s.kafka.SendMessage("email-service-topic-transaction-create", strconv.Itoa(int(updatedTransaction.TransactionID)), payloadBytes)
			if err != nil {
				s.logger.Error("failed to send transaction email via kafka", zap.Error(err), zap.Int("transaction_id", int(updatedTransaction.TransactionID)))
			}
		} else {
			s.logger.Warn("Kafka is nil, skipping transaction email", zap.Int("transaction_id", int(updatedTransaction.TransactionID)))
		}
	}()

	logSuccess("Successfully created transaction", zap.Int("transaction.id", int(updatedTransaction.TransactionID)))

	return updatedTransaction, nil
}

func (s *transactionCommandService) Update(ctx context.Context, apiKey string, request *requests.UpdateTransactionRequest) (*db.UpdateTransactionRow, error) {
	const method = "UpdateTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() { end(status) }()

	transaction, err := s.transactionQueryRepository.FindById(ctx, *request.TransactionID)
	if err != nil {
		status = "error"
		s.markTransactionAsFailed(ctx, *request.TransactionID, method, span)
		return errorhandler.HandleError[*db.UpdateTransactionRow](s.logger, err, method, span, zap.Int("transaction_id", *request.TransactionID))
	}

	merchant, err := s.merchantRepository.FindByApiKey(ctx, apiKey)
	if err != nil || transaction.MerchantID != merchant.MerchantID {
		status = "error"
		s.markTransactionAsFailed(ctx, *request.TransactionID, method, span)
		return errorhandler.HandleError[*db.UpdateTransactionRow](s.logger, err, method, span, zap.String("api_key", apiKey))
	}

	card, err := s.cardRepository.FindCardByCardNumber(ctx, transaction.CardNumber)
	if err != nil {
		status = "error"
		s.markTransactionAsFailed(ctx, *request.TransactionID, method, span)
		return errorhandler.HandleError[*db.UpdateTransactionRow](s.logger, err, method, span, zap.String("card_number", transaction.CardNumber))
	}

	saldo, err := s.saldoRepository.FindByCardNumber(ctx, card.CardNumber)
	if err != nil {
		status = "error"
		s.markTransactionAsFailed(ctx, *request.TransactionID, method, span)
		return errorhandler.HandleError[*db.UpdateTransactionRow](s.logger, err, method, span, zap.String("card_number", card.CardNumber))
	}

	rollbackBalance := saldo.TotalBalance + transaction.Amount
	if _, err := s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
		CardNumber:   card.CardNumber,
		TotalBalance: int(rollbackBalance),
	}); err != nil {
		status = "error"
		s.markTransactionAsFailed(ctx, *request.TransactionID, method, span)
		return errorhandler.HandleError[*db.UpdateTransactionRow](s.logger, err, method, span, zap.String("card_number", card.CardNumber))
	}

	if int(rollbackBalance) < request.Amount {
		status = "error"
		err := errors.New("insufficient balance after rollback")
		s.markTransactionAsFailed(ctx, *request.TransactionID, method, span)
		return errorhandler.HandleError[*db.UpdateTransactionRow](s.logger, err, method, span, zap.Float64("current_balance", float64(rollbackBalance)), zap.Float64("requested_amount", float64(request.Amount)))
	}

	newBalance := int(rollbackBalance) - request.Amount
	if _, err := s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
		CardNumber:   card.CardNumber,
		TotalBalance: newBalance,
	}); err != nil {
		status = "error"
		s.markTransactionAsFailed(ctx, *request.TransactionID, method, span)
		return errorhandler.HandleError[*db.UpdateTransactionRow](s.logger, err, method, span, zap.String("card_number", card.CardNumber))
	}

	parsedTime := transaction.TransactionTime

	merchantId := int(transaction.MerchantID)
	transactionId := int(transaction.TransactionID)

	updateReq := &requests.UpdateTransactionRequest{
		TransactionID:   &transactionId,
		CardNumber:      transaction.CardNumber,
		Amount:          request.Amount,
		PaymentMethod:   request.PaymentMethod,
		MerchantID:      &merchantId,
		TransactionTime: parsedTime,
	}
	res, err := s.transactionCommandRepository.UpdateTransaction(ctx, updateReq)
	if err != nil {
		status = "error"
		s.markTransactionAsFailed(ctx, *request.TransactionID, method, span)
		return errorhandler.HandleError[*db.UpdateTransactionRow](s.logger, err, method, span, zap.Int("transaction_id", *request.TransactionID))
	}

	if _, err := s.transactionCommandRepository.UpdateTransactionStatus(ctx, &requests.UpdateTransactionStatus{
		TransactionID: int(transaction.TransactionID),
		Status:        "success",
	}); err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](s.logger, err, method, span, zap.Int("transaction_id", int(transaction.TransactionID)))
	}

	logSuccess("Successfully updated transaction", zap.Int("transaction.id", int(res.TransactionID)))

	return res, nil
}

func (s *transactionCommandService) TrashedTransaction(ctx context.Context, transaction_id int) (*db.Transaction, error) {
	const method = "TrashedTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transaction_id", transaction_id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Starting TrashedTransaction process", zap.Int("transaction_id", transaction_id))

	res, err := s.transactionCommandRepository.TrashedTransaction(ctx, transaction_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Transaction](
			s.logger,
			transaction_errors.ErrFailedTrashedTransaction,
			method,
			span,

			zap.Int("transaction_id", transaction_id),
		)
	}

	logSuccess("Successfully trashed transaction", zap.Int("transaction_id", transaction_id))

	return res, nil
}

func (s *transactionCommandService) RestoreTransaction(ctx context.Context, transaction_id int) (*db.Transaction, error) {
	const method = "RestoreTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transaction_id", transaction_id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Starting RestoreTransaction process", zap.Int("transaction_id", transaction_id))

	res, err := s.transactionCommandRepository.RestoreTransaction(ctx, transaction_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Transaction](
			s.logger,
			transaction_errors.ErrFailedRestoreTransaction,
			method,
			span,

			zap.Int("transaction_id", transaction_id),
		)
	}

	logSuccess("Successfully restored transaction", zap.Int("transaction_id", transaction_id))

	return res, nil
}

func (s *transactionCommandService) DeleteTransactionPermanent(ctx context.Context, transaction_id int) (bool, error) {
	const method = "DeleteTransactionPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transaction_id", transaction_id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Starting DeleteTransactionPermanent process", zap.Int("transaction_id", transaction_id))

	_, err := s.transactionCommandRepository.DeleteTransactionPermanent(ctx, transaction_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			transaction_errors.ErrFailedDeleteTransactionPermanent,
			method,
			span,

			zap.Int("transaction_id", transaction_id),
		)
	}

	logSuccess("Successfully permanently deleted transaction", zap.Int("transaction_id", transaction_id))

	return true, nil
}

func (s *transactionCommandService) RestoreAllTransaction(ctx context.Context) (bool, error) {
	const method = "RestoreAllTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	s.logger.Debug("Restoring all transactions")

	_, err := s.transactionCommandRepository.RestoreAllTransaction(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			transaction_errors.ErrFailedRestoreAllTransactions,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all transactions")
	return true, nil
}

func (s *transactionCommandService) DeleteAllTransactionPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllTransactionPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	s.logger.Debug("Permanently deleting all transactions")

	_, err := s.transactionCommandRepository.DeleteAllTransactionPermanent(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			transaction_errors.ErrFailedDeleteAllTransactionsPermanent,
			method,
			span,
		)
	}

	logSuccess("Successfully deleted all transactions permanently")
	return true, nil
}

func (s *transactionCommandService) markTransactionAsFailed(ctx context.Context, transactionID int, method string, span trace.Span) {
	req := &requests.UpdateTransactionStatus{
		TransactionID: transactionID,
		Status:        "failed",
	}
	go func() {
		if _, err := s.transactionCommandRepository.UpdateTransactionStatus(ctx, req); err != nil {
			s.logger.Error("compensation: failed to mark transaction as failed", zap.Error(err), zap.Int("transaction_id", transactionID), zap.String("method", method))
		}
	}()
}
