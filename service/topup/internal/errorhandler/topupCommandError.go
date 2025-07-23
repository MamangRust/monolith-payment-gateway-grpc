package errorhandler

import (
	"net/http"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	sharederrorhandler "github.com/MamangRust/monolith-payment-gateway-shared/errorhandler"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// topupCommandError is a struct that implements the TopupCommandError interface.
type topupCommandError struct {
	logger logger.LoggerInterface
}

// NewTopupCommandError initializes a new topupCommandError with the provided logger.
// It returns an instance of the topupCommandError struct.
func NewTopupCommandError(logger logger.LoggerInterface) TopupCommandErrorHandler {
	return &topupCommandError{
		logger: logger,
	}
}

// HandleInvalidParseTimeError handles errors related to parsing time values.
// It constructs an appropriate error response for cases where the provided
// time value is invalid or cannot be parsed. The method logs the error
// and returns a TopupResponse and an ErrorResponse indicating a bad request.
//
// Parameters:
//   - err: the error encountered during time parsing.
//   - method: the name of the method where the error occurred.
//   - tracePrefix: the prefix used for tracing the error.
//   - span: the trace span for the request.
//   - status: a pointer to a string representing the status of the operation.
//   - rawTime: the raw time string that failed to parse.
//   - fields: additional context fields for logging.
//
// Returns:
//   - *response.TopupResponse: the response for the topup operation (nil in error cases).
//   - *response.ErrorResponse: the constructed error response indicating the failure.
func (e *topupCommandError) HandleInvalidParseTimeError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	rawTime string,
	fields ...zap.Field,
) (*response.TopupResponse, *response.ErrorResponse) {
	errResp := &response.ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: "Failed to parse the given time value",
		Status:  "invalid_parse_time",
	}

	return sharederrorhandler.HandleErrorTemplate[*response.TopupResponse](e.logger, err, method, tracePrefix, "Invalid parse time", span, status, errResp, fields...)

}

// HandleRepositorySingleError processes single-result errors from the repository.
// It logs the error, updates the trace span, and returns a standardized error response.
//
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
//   - A TopupResponse pointer if successful, otherwise nil.
//   - A standardized ErrorResponse describing the single-result failure.
func (e *topupCommandError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.TopupResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TopupResponse](e.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

// HandleCreateTopupError processes errors during the topup creation process.
// It logs the error, updates the trace span, and returns a standardized error response.
//
// Parameters:
//   - err: The error that occurred during the topup creation.
//   - method: The name of the method where the error originated.
//   - tracePrefix: A prefix used for generating the trace ID.
//   - span: The tracing span used for recording error details.
//   - status: A pointer to a string that will be updated with the formatted status.
//   - fields: Additional contextual fields for logging.
//
// Returns:
//   - A TopupResponse pointer if the operation is successful, otherwise nil.
//   - A standardized ErrorResponse describing the creation failure.
func (e *topupCommandError) HandleCreateTopupError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TopupResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TopupResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		topup_errors.ErrFailedCreateTopup,
		fields...,
	)
}

// HandleUpdateTopupError processes errors during the topup update process.
// It logs the error, updates the trace span, and returns a standardized error response.
//
// Parameters:
//   - err: The error that occurred during the topup update.
//   - method: The name of the method where the error originated.
//   - tracePrefix: A prefix used for generating the trace ID.
//   - span: The tracing span used for recording error details.
//   - status: A pointer to a string that will be updated with the formatted status.
//   - fields: Additional contextual fields for logging.
//
// Returns:
//   - A TopupResponse pointer if the operation is successful, otherwise nil.
//   - A standardized ErrorResponse describing the update failure.
func (e *topupCommandError) HandleUpdateTopupError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TopupResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TopupResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		topup_errors.ErrFailedUpdateTopup,
		fields...,
	)
}

// HandleTrashedTopupError processes errors during the topup trash process.
// It logs the error, updates the trace span, and returns a standardized error response.
//
// Parameters:
//   - err: The error that occurred during the topup trash.
//   - method: The name of the method where the error originated.
//   - tracePrefix: A prefix used for generating the trace ID.
//   - span: The tracing span used for recording error details.
//   - status: A pointer to a string that will be updated with the formatted status.
//   - fields: Additional contextual fields for logging.
//
// Returns:
//   - A TopupResponseDeleteAt pointer if the operation is successful, otherwise nil.
//   - A standardized ErrorResponse describing the trash failure.
func (e *topupCommandError) HandleTrashedTopupError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TopupResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.TopupResponseDeleteAt](
		e.logger,
		err, method, tracePrefix, span, status,
		topup_errors.ErrFailedTrashTopup,
		fields...,
	)
}

// HandleRestoreTopupError processes errors during the topup restore process.
// It logs the error, updates the trace span, and returns a standardized error response.
//
// Parameters:
//   - err: The error that occurred during the topup restore.
//   - method: The name of the method where the error originated.
//   - tracePrefix: A prefix used for generating the trace ID.
//   - span: The tracing span used for recording error details.
//   - status: A pointer to a string that will be updated with the formatted status.
//   - fields: Additional contextual fields for logging.
//
// Returns:
//   - A TopupResponseDeleteAt pointer if the operation is successful, otherwise nil.
//   - A standardized ErrorResponse describing the restore failure.
func (e *topupCommandError) HandleRestoreTopupError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (*response.TopupResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.TopupResponse](
		e.logger,
		err, method, tracePrefix, span, status,
		topup_errors.ErrFailedRestoreTopup,
		fields...,
	)
}

// HandleDeleteTopupPermanentError processes errors that occur during the permanent deletion of a Topup.
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
func (e *topupCommandError) HandleDeleteTopupPermanentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		topup_errors.ErrFailedDeleteTopup,
		fields...,
	)
}

// HandleRestoreAllTopupError processes errors that occur during the restoration of all Topups.
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
//   - A boolean indicating whether the restoration was successful.
//   - A standardized ErrorResponse detailing the restoration error.
func (e *topupCommandError) HandleRestoreAllTopupError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		topup_errors.ErrFailedRestoreAllTopups,
		fields...,
	)
}

// HandleDeleteAllTopupPermanentError processes errors that occur during the permanent deletion of all Topups.
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
func (e *topupCommandError) HandleDeleteAllTopupPermanentError(
	err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](
		e.logger,
		err, method, tracePrefix, span, status,
		topup_errors.ErrFailedDeleteAllTopups,
		fields...,
	)
}
