package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	role_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/role_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// roleQueryError is a struct that implements the RoleQueryError interface
type roleQueryError struct {
	logger logger.LoggerInterface
}

// NewRoleQueryError returns a new instance of roleQueryError with the given logger.
// It returns an instance of the roleQueryError struct.
func NewRoleQueryError(logger logger.LoggerInterface) RoleQueryErrorHandler {
	return &roleQueryError{
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
//   - A slice of RoleResponse pointers if successful, otherwise nil.
//   - A pointer to an integer representing additional pagination details, otherwise nil.
//   - A standardized ErrorResponse describing the pagination failure.
func (e *roleQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.RoleResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.RoleResponse](
		e.logger, err, method, tracePrefix, span, status, role_errors.ErrFailedFindAll, fields...,
	)
}

// HandleRepositoryPaginationDeletedError processes pagination errors from the repository
// when retrieving deleted role documents. It logs the error, updates the trace span,
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
//   - A slice of RoleResponseDeleteAt pointers if successful, otherwise nil.
//   - A pointer to an integer representing additional pagination details, otherwise nil.
//   - A standardized ErrorResponse describing the pagination failure.
func (e *roleQueryError) HandleRepositoryPaginationDeletedError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.RoleResponseDeleteAt](
		e.logger, err, method, tracePrefix, span, status, errResp, fields...,
	)
}

// HandleRepositoryListError processes list errors from the repository.
// It logs the error, updates the trace span, and returns a standardized error response.
//
// Parameters:
//   - err: The error that occurred during the list operation.
//   - method: The name of the method where the error originated.
//   - tracePrefix: A prefix used for generating the trace ID.
//   - span: The tracing span used for recording error details.
//   - status: A pointer to a string that will be updated with the formatted status.
//   - fields: Additional contextual fields for logging.
//
// Returns:
//   - A slice of RoleResponse pointers if successful, otherwise nil.
//   - A standardized ErrorResponse describing the list failure.
func (e *roleQueryError) HandleRepositoryListError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.RoleResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.RoleResponse](e.logger, err, method, tracePrefix, span, status, role_errors.ErrFailedFindAll, fields...)
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
//   - defaultErr: A pointer to an ErrorResponse that will be updated with the error details.
//   - fields: Additional contextual fields for logging.
//
// Returns:
//   - A RoleResponse pointer if successful, otherwise nil.
//   - A standardized ErrorResponse describing the single-result failure.
func (e *roleQueryError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (*response.RoleResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.RoleResponse](e.logger, err, method, tracePrefix, span, status, defaultErr, fields...)
}
