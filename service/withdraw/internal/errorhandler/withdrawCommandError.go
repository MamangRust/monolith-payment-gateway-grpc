package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	sharederrorhandler "github.com/MamangRust/monolith-payment-gateway-shared/errorhandler"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/service"
	withdraw_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// withdrawCommandError handles error logging for withdraw command operations.
type withdrawCommandError struct {
	logger logger.LoggerInterface
}

// NewWithdrawCommandError initializes a new withdrawCommandError with the provided logger.
// It returns an instance of the withdrawCommandError struct.
func NewWithdrawCommandError(logger logger.LoggerInterface) WithdrawCommandErrorHandler {
	return &withdrawCommandError{
		logger: logger,
	}
}

// HandleInsufficientBalanceError handles errors related to insufficient balance during a withdraw operation.
// It logs the error, records it in the trace span, and returns a standardized error response.
// Parameters:
//   - err: The error that occurred.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - cardNumber: The card number associated with the insufficient balance.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A WithdrawResponse with error details and a standardized ErrorResponse.
func (t *withdrawCommandError) HandleInsufficientBalanceError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	cardNumber string,
	fields ...zap.Field,
) (*response.WithdrawResponse, *response.ErrorResponse) {
	return sharederrorhandler.HandleErrorTemplate[*response.WithdrawResponse](t.logger, err, method, tracePrefix, "Insufficient Balance", span, status, saldo_errors.ErrFailedInsuffientBalance, fields...)
}

// HandleRepositorySingleError processes single-result errors from the repository.
// It logs the error, records it to the trace span, and returns a standardized error response.
// Parameters:
//   - err: The error that occurred during the single-result operation.
//   - method: The name of the method where the error originated.
//   - tracePrefix: A prefix used for generating the trace ID.
//   - span: The tracing span used for recording error details.
//   - status: A pointer to a string that will be updated with the formatted status.
//   - errResp: A pointer to an ErrorResponse that will be updated with the error details.
//   - fields: Additional contextual fields for logging.
//
// Returns:
//   - A WithdrawResponse pointer if successful, otherwise nil.
//   - A standardized ErrorResponse describing the single-result failure.
func (w *withdrawCommandError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.WithdrawResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.WithdrawResponse](w.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

// HandleCreateWithdrawError processes errors during creation of a new withdraw.
// It logs the error, records it to the trace span, and returns a standardized error response.
// Args:
//   - err: The error that occurred.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A WithdrawResponse with error details and a standardized ErrorResponse.
func (w *withdrawCommandError) HandleCreateWithdrawError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.WithdrawResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.WithdrawResponse](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedCreateWithdraw, fields...)
}

// HandleUpdateWithdrawError processes errors during updating a withdraw.
// It logs the error, records it to the trace span, and returns a standardized error response.
// Args:
//   - err: The error that occurred.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A WithdrawResponse with error details and a standardized ErrorResponse.
func (w *withdrawCommandError) HandleUpdateWithdrawError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.WithdrawResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.WithdrawResponse](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedUpdateWithdraw, fields...)
}

// HandleTrashedWithdrawError processes errors during trashing a withdraw.
// It logs the error, records it to the trace span, and returns a standardized error response.
// Args:
//   - err: The error that occurred.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A WithdrawResponse with error details and a standardized ErrorResponse.
func (w *withdrawCommandError) HandleTrashedWithdrawError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.WithdrawResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.WithdrawResponseDeleteAt](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedTrashedWithdraw, fields...)
}

// HandleRestoreWithdrawError processes errors during restore of a withdraw.
// It logs the error, records it to the trace span, and returns a standardized error response.
// Args:
//   - err: The error that occurred.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A WithdrawResponse with error details and a standardized ErrorResponse.
func (w *withdrawCommandError) HandleRestoreWithdrawError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.WithdrawResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.WithdrawResponse](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedRestoreWithdraw, fields...)
}

// HandleDeleteWithdrawPermanentError processes errors during the permanent deletion of a withdraw.
// It logs the error, records it to the trace span, and returns a standardized error response.
// Parameters:
//   - err: The error that occurred.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A boolean indicating whether the deletion was successful.
//   - A standardized ErrorResponse detailing the deletion error.
func (w *withdrawCommandError) HandleDeleteWithdrawPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedDeleteWithdrawPermanent, fields...)
}

// HandleRestoreAllWithdrawError processes errors during the restoration of all withdraw.
// It logs the error, records it to the trace span, and returns a standardized error response.
// Parameters:
//   - err: The error that occurred.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A boolean indicating whether the restoration was successful.
//   - A standardized ErrorResponse detailing the restoration error.
func (w *withdrawCommandError) HandleRestoreAllWithdrawError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedRestoreAllWithdraw, fields...)
}

// HandleDeleteAllWithdrawPermanentError handles errors that occur during the permanent deletion of all withdraws.
// It logs the error, records it to the trace span, and returns a standardized error response.
//
// Parameters:
//   - err: The error that occurred.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A boolean indicating whether the deletion was successful.
//   - A standardized ErrorResponse detailing the deletion error.
func (w *withdrawCommandError) HandleDeleteAllWithdrawPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](w.logger, err, method, tracePrefix, span, status, withdraw_errors.ErrFailedDeleteAllWithdrawPermanent, fields...)
}
