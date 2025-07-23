package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/service"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// saldoCommandError represents an error handler for saldo command operations.
type saldoCommandError struct {
	logger logger.LoggerInterface
}

// NewSaldoCommandError initializes a new saldoCommandError with the provided logger.
// It returns an instance of the saldoCommandError struct.
func NewSaldoCommandError(logger logger.LoggerInterface) SaldoCommandErrorHandler {
	return &saldoCommandError{
		logger: logger,
	}
}

// HandleFindCardByNumberError processes errors during card lookup by card number
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
//   - A SaldoResponse with error details and a standardized ErrorResponse.
func (e *saldoCommandError) HandleFindCardByNumberError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.SaldoResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.SaldoResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		card_errors.ErrCardNotFoundRes,
		fields...,
	)
}

// HandleCreateSaldoError processes errors during card creation
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
//   - A SaldoResponse with error details and a standardized ErrorResponse.
func (e *saldoCommandError) HandleCreateSaldoError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.SaldoResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.SaldoResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		saldo_errors.ErrFailedCreateSaldo,
		fields...,
	)
}

// HandleUpdateSaldoError processes errors during card update
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
//   - A SaldoResponse with error details and a standardized ErrorResponse.
func (e *saldoCommandError) HandleUpdateSaldoError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.SaldoResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.SaldoResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		saldo_errors.ErrFailedUpdateSaldo,
		fields...,
	)
}

// HandleTrashSaldoError processes errors during Saldo soft deletion (trashing)
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
//   - A SaldoResponse with error details and a standardized ErrorResponse.
func (e *saldoCommandError) HandleTrashSaldoError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.SaldoResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.SaldoResponseDeleteAt](
		e.logger,
		err, method, tracePrefix, span, status,
		saldo_errors.ErrFailedTrashSaldo,
		fields...,
	)
}

// HandleRestoreSaldoError processes errors during Saldo restore (undoing trashing)
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
//   - A SaldoResponse with error details and a standardized ErrorResponse.
func (e *saldoCommandError) HandleRestoreSaldoError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.SaldoResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.SaldoResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		saldo_errors.ErrFailedRestoreSaldo,
		fields...,
	)
}

// HandleDeleteSaldoPermanentError processes errors during the permanent deletion of a Saldo.
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
func (e *saldoCommandError) HandleDeleteSaldoPermanentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		saldo_errors.ErrFailedDeleteSaldoPermanent,
		fields...,
	)
}

// HandleRestoreAllSaldoError processes errors that occur during the restoration of all Saldo.
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
func (e *saldoCommandError) HandleRestoreAllSaldoError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		saldo_errors.ErrFailedRestoreAllSaldo,
		fields...,
	)
}

// HandleDeleteAllSaldoPermanentError processes errors that occur during the permanent deletion of all Saldo.
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
func (e *saldoCommandError) HandleDeleteAllSaldoPermanentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		saldo_errors.ErrFailedDeleteAllSaldoPermanent,
		fields...,
	)
}
