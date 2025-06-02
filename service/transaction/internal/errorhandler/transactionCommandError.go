package errorhandler

import (
	"net/http"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	traceunic "github.com/MamangRust/monolith-payment-gateway-pkg/trace_unic"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transactionCommandError struct {
	logger logger.LoggerInterface
}

func NewTransactionCommandError(logger logger.LoggerInterface) *transactionCommandError {
	return &transactionCommandError{logger: logger}
}

func (t *transactionCommandError) HandleInvalidParseTimeError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	rawTime string,
	fields ...zap.Field,
) (*response.TransactionResponse, *response.ErrorResponse) {

	traceID := traceunic.GenerateTraceID("INVALID_PARSE_TIME")

	t.logger.Error("Invalid time parse error",
		append(fields,
			zap.String("trace.id", traceID),
			zap.String("raw_time", rawTime),
			zap.String("method", method),
			zap.String("trace_prefix", tracePrefix),
			zap.Error(err),
		)...,
	)

	span.SetAttributes(attribute.String("trace.id", traceID))
	span.RecordError(err)
	span.SetStatus(codes.Error, "Invalid parse time")

	if status != nil {
		*status = "invalid_parse_time"
	}

	return nil, &response.ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: "Failed to parse the given time value",
		Status:  "invalid_parse_time",
	}
}

func (t *transactionCommandError) HandleInsufficientBalanceError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	cardNumber string,
	fields ...zap.Field,
) (*response.TransactionResponse, *response.ErrorResponse) {

	traceID := traceunic.GenerateTraceID("INSUFFICIENT_BALANCE")

	t.logger.Error("Insufficient balance",
		append(fields,
			zap.String("trace.id", traceID),
			zap.String("card_number", cardNumber),
			zap.String("method", method),
			zap.String("trace_prefix", tracePrefix),
			zap.Error(err),
		)...,
	)

	span.SetAttributes(attribute.String("trace.id", traceID))
	span.RecordError(err)
	span.SetStatus(codes.Error, "Insufficient balance")

	if status != nil {
		*status = "insufficient_balance"
	}

	return nil, saldo_errors.ErrFailedInsuffientBalance
}

func (t *transactionCommandError) HandleRepositorySingleError(err error, method, tracePrefix string, span trace.Span, status *string, errResp *response.ErrorResponse, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.TransactionResponse](t.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (t *transactionCommandError) HandleCreateTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.TransactionResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedCreateTransaction, fields...)
}

func (t *transactionCommandError) HandleUpdateTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.TransactionResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedUpdateTransaction, fields...)
}

func (t *transactionCommandError) HandleTrashedTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.TransactionResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedTrashedTransaction, fields...)
}

func (t *transactionCommandError) HandleRestoreTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.TransactionResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedRestoreTransaction, fields...)
}

func (t *transactionCommandError) HandleDeleteTransactionPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorTemplate[bool](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedDeleteTransactionPermanent, fields...)
}

func (t *transactionCommandError) HandleRestoreAllTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorTemplate[bool](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedRestoreAllTransactions, fields...)
}

func (t *transactionCommandError) HandleDeleteAllTransactionPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorTemplate[bool](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedDeleteAllTransactionsPermanent, fields...)
}
