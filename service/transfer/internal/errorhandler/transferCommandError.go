package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transferCommandError struct {
	logger logger.LoggerInterface
}

func NewTransferCommandError(logger logger.LoggerInterface) *transferCommandError {
	return &transferCommandError{
		logger: logger,
	}
}

func (t *transferCommandError) HandleSenderInsufficientBalanceError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	senderCardNumber string,
	fields ...zap.Field,
) (*response.TransferResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.TransferResponse](t.logger, err, method, tracePrefix, "InsufficientBalance", span, status, saldo_errors.ErrFailedInsuffientBalance, fields...)
}

func (t *transferCommandError) HandleReceiverInsufficientBalanceError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	receiverCardNumber string,
	fields ...zap.Field,
) (*response.TransferResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.TransferResponse](t.logger, err, method, tracePrefix, "InsufficientBalance", span, status, saldo_errors.ErrFailedInsuffientBalance, fields...)
}

func (t *transferCommandError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.TransferResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransferResponse](t.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (t *transferCommandError) HandleCreateTransferError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransferResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransferResponse](t.logger, err, method, tracePrefix, span, status, transfer_errors.ErrFailedCreateTransfer, fields...)
}

func (t *transferCommandError) HandleUpdateTransferError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransferResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransferResponse](t.logger, err, method, tracePrefix, span, status, transfer_errors.ErrFailedUpdateTransfer, fields...)
}

func (t *transferCommandError) HandleTrashedTransferError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransferResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransferResponse](t.logger, err, method, tracePrefix, span, status, transfer_errors.ErrFailedTrashedTransfer, fields...)
}

func (t *transferCommandError) HandleRestoreTransferError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TransferResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TransferResponse](t.logger, err, method, tracePrefix, span, status, transfer_errors.ErrFailedRestoreTransfer, fields...)
}

func (t *transferCommandError) HandleDeleteTransferPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](t.logger, err, method, tracePrefix, span, status, transfer_errors.ErrFailedDeleteTransferPermanent, fields...)
}

func (t *transferCommandError) HandleRestoreAllTransferError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](t.logger, err, method, tracePrefix, span, status, transfer_errors.ErrFailedRestoreAllTransfers, fields...)
}

func (t *transferCommandError) HandleDeleteAllTransferPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](t.logger, err, method, tracePrefix, span, status, transfer_errors.ErrFailedDeleteAllTransfersPermanent, fields...)
}
