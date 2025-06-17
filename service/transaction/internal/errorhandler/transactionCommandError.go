package errorhandler

import (
	"net/http"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors"
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
	errResp := &response.ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: "Failed to parse the given time value",
		Status:  "invalid_parse_time",
	}

	return handleErrorTemplate[*response.TransactionResponse](t.logger, err, method, tracePrefix, "Invalid parse time", span, status, errResp, fields...)
}

func (t *transactionCommandError) HandleInsufficientBalanceError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	cardNumber string,
	fields ...zap.Field,
) (*response.TransactionResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransactionResponse](t.logger, err, method, tracePrefix, span, status, saldo_errors.ErrFailedInsuffientBalance, fields...)
}

func (t *transactionCommandError) HandleRepositorySingleError(err error, method, tracePrefix string, span trace.Span, status *string, errResp *response.ErrorResponse, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransactionResponse](t.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (t *transactionCommandError) HandleCreateTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransactionResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedCreateTransaction, fields...)
}

func (t *transactionCommandError) HandleUpdateTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransactionResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedUpdateTransaction, fields...)
}

func (t *transactionCommandError) HandleTrashedTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransactionResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedTrashedTransaction, fields...)
}

func (t *transactionCommandError) HandleRestoreTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransactionResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransactionResponse](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedRestoreTransaction, fields...)
}

func (t *transactionCommandError) HandleDeleteTransactionPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedDeleteTransactionPermanent, fields...)
}

func (t *transactionCommandError) HandleRestoreAllTransactionError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedRestoreAllTransactions, fields...)
}

func (t *transactionCommandError) HandleDeleteAllTransactionPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](t.logger, err, method, tracePrefix, span, status, transaction_errors.ErrFailedDeleteAllTransactionsPermanent, fields...)
}
