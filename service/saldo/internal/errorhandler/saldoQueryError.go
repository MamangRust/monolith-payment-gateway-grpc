package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// saldoQueryError is a struct that implements the SaldoQueryError interface
type saldoQueryError struct {
	logger logger.LoggerInterface
}

// NewSaldoQueryError returns a new instance of SaldoQueryError with the given logger.
//
// It is used to create a new SaldoQueryError handler with the given logger.
func NewSaldoQueryError(logger logger.LoggerInterface) SaldoQueryErrorHandler {
	return &saldoQueryError{
		logger: logger,
	}
}

// HandleRepositoryPaginationError processes pagination errors from the repository.
// It logs the error, updates the trace span, and returns a standardized error response.
//
// Parameters:
//   - err: The error that occurred during the pagination operation.
//   - method: The name of the method where the error originated.
//   - tracePrefix: A prefix used for generating the trace ID.
//   - span: The tracing span used for recording error details.
//   - status: A pointer to a string that will be updated with the formatted status.
//   - fields: Additional contextual fields for logging.
//
// Returns:
//   - A slice of SaldoResponse pointers if successful, otherwise nil.
//   - A pointer to an integer representing additional pagination details, otherwise nil.
//   - A standardized ErrorResponse describing the pagination failure.
func (e *saldoQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.SaldoResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.SaldoResponse](e.logger, err, method, tracePrefix, span, status, saldo_errors.ErrFailedFindAllSaldos, fields...)
}

// HandleRepositoryPaginationDeleteAtError processes pagination errors from the repository
// when retrieving deleted saldo documents. It logs the error, updates the trace span,
// and returns a standardized error response.
//
// Parameters:
//   - err: The error that occurred during the pagination operation.
//   - method: The name of the method where the error originated.
//   - tracePrefix: A prefix used for generating the trace ID.
//   - span: The tracing span used for recording error details.
//   - status: A pointer to a string that will be updated with the formatted status.
//   - errResp: A pointer to an ErrorResponse that will be updated with the error details.
//   - fields: Additional contextual fields for logging.
//
// Returns:
//   - A slice of SaldoResponseDeleteAt pointers if successful, otherwise nil.
//   - A pointer to an integer representing additional pagination details, otherwise nil.
//   - A standardized ErrorResponse describing the pagination failure.
func (e *saldoQueryError) HandleRepositoryPaginationDeleteAtError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.SaldoResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.SaldoResponseDeleteAt](e.logger, err, method, tracePrefix, span, status, errResp, fields...)
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
//   - A SaldoResponse pointer if successful, otherwise nil.
//   - A standardized ErrorResponse describing the single-result failure.
func (e *saldoQueryError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.SaldoResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.SaldoResponse](e.logger, err, method, tracePrefix, span, status, errResp, fields...)
}
