package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// cardQueryError is a struct that implements the CardQueryError interface
type cardQueryError struct {
	logger logger.LoggerInterface
}

// NewCardQueryError returns a new instance of cardQueryError
func NewCardQueryError(logger logger.LoggerInterface) CardQueryErrorHandler {
	return &cardQueryError{
		logger: logger,
	}
}

// HandleFindAllError processes errors when fetching all cards.
// It logs the error, records it to the trace span, and returns a paginated CardResponse with error details.
// Parameters:
//   - err: The error that occurred.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of CardResponse pointers with paginated card details.
//   - A pointer to an integer indicating the total count of cards.
//   - A standardized ErrorResponse detailing the error.
func (c *cardQueryError) HandleFindAllError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.CardResponse](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindAllCards, fields...)
}

// HandleFindByActiveError processes errors when fetching all active cards.
// It logs the error, records it to the trace span, and returns a paginated CardResponse with error details.
// Parameters:
//   - err: The error that occurred.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of CardResponse pointers with paginated card details.
//   - A pointer to an integer indicating the total count of cards.
//   - A standardized ErrorResponse detailing the error.
func (c *cardQueryError) HandleFindByActiveError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.CardResponseDeleteAt](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindActiveCards, fields...)
}

// HandleFindByTrashedError processes errors when fetching all trashed cards.
// It logs the error, records it to the trace span, and returns a paginated CardResponse with error details.
// Parameters:
//   - err: The error that occurred.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of CardResponse pointers with paginated card details.
//   - A pointer to an integer indicating the total count of cards.
//   - A standardized ErrorResponse detailing the error.
func (c *cardQueryError) HandleFindByTrashedError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.CardResponseDeleteAt](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindTrashedCards, fields...)
}

// HandleFindByIdError processes errors during card lookup by ID
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
//   - A CardResponse with error details and a standardized ErrorResponse.
func (c *cardQueryError) HandleFindByIdError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.CardResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.CardResponse](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindById, fields...)
}

// HandleFindByUserIdError processes errors during card lookup by user ID
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
//   - A CardResponse with error details and a standardized ErrorResponse.
func (c *cardQueryError) HandleFindByUserIdError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.CardResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.CardResponse](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindByUserID, fields...)
}

// HandleFindByCardNumberError processes errors during card lookup by card number
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
//   - A CardResponse with error details and a standardized ErrorResponse.
func (c *cardQueryError) HandleFindByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.CardResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.CardResponse](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindByCardNumber, fields...)
}
