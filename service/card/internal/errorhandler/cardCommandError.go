package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/service"
	user_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// cardCommandError is a struct that implements the CardCommandError interface
type cardCommandError struct {
	logger logger.LoggerInterface
}

// NewCardCommandError initializes a new cardCommandError with the provided logger.
// It returns an instance of the cardCommandError struct.
func NewCardCommandError(logger logger.LoggerInterface) CardCommandErrorHandler {
	return &cardCommandError{
		logger: logger,
	}
}

// HandleFindByIdUserError processes errors during user lookup by ID
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
//   - A CardResponse with error details and a standardized ErrorResponse.
func (c *cardCommandError) HandleFindByIdUserError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.CardResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.CardResponse](c.logger, err, method, tracePrefix, span, status, user_errors.ErrUserNotFoundRes, fields...)
}

// HandleCreateCardError processes errors that occur during card creation.
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
//   - A CardResponse, which is nil since operation failed.
//   - A standardized ErrorResponse detailing the card creation failure.
func (c *cardCommandError) HandleCreateCardError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.CardResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.CardResponse](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedCreateCard, fields...)
}

// HandleUpdateCardError processes errors that occur during card updates.
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
//   - A CardResponse, which is nil since operation failed.
//   - A standardized ErrorResponse detailing the card update failure.
func (c *cardCommandError) HandleUpdateCardError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.CardResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.CardResponse](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedUpdateCard, fields...)
}

// HandleTrashedCardError processes errors that occur during card trashing.
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
//   - A CardResponse, which is nil since operation failed.
//   - A standardized ErrorResponse detailing the card trashing failure.
func (c *cardCommandError) HandleTrashedCardError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.CardResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.CardResponseDeleteAt](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedTrashCard, fields...)
}

// HandleRestoreCardError processes errors that occur during card restoration.
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
//   - A CardResponse, which is nil since operation failed.
//   - A standardized ErrorResponse detailing the card restoration failure.
func (c *cardCommandError) HandleRestoreCardError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.CardResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.CardResponse](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedRestoreCard, fields...)
}

// HandleDeleteCardPermanentError processes errors that occur during card deletion.
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
//   - A boolean indicating whether the error is fatal.
//   - A standardized ErrorResponse detailing the card deletion failure.
func (c *cardCommandError) HandleDeleteCardPermanentError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedDeleteCard, fields...)
}

// HandleRestoreAllCardError processes errors that occur during card restoration.
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
//   - A boolean indicating whether the error is fatal.
//   - A standardized ErrorResponse detailing the card restoration failure.
func (c *cardCommandError) HandleRestoreAllCardError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedRestoreAllCards, fields...)
}

// HandleDeleteAllCardPermanentError processes errors that occur during the permanent deletion of all cards.
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
//   - A boolean indicating whether the error is fatal.
//   - A standardized ErrorResponse detailing the card deletion failure.
func (c *cardCommandError) HandleDeleteAllCardPermanentError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedDeleteAllCards, fields...)
}
