package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	user_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// userQueryError handles error logging for user query operations.
type userQueryError struct {
	logger logger.LoggerInterface
}

// NewUserQueryError initializes a new userQueryError with the provided logger.
// It returns an instance of the userQueryError struct.
func NewUserQueryError(logger logger.LoggerInterface) UserQueryError {
	return &userQueryError{logger: logger}
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
//   - A slice of UserResponse pointers if successful, otherwise nil.
//   - A pointer to an integer representing additional pagination details, otherwise nil.
//   - A standardized ErrorResponse describing the pagination failure.
func (u *userQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.UserResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.UserResponse](u.logger, err, method, tracePrefix, span, status, user_errors.ErrUserNotFoundRes, fields...)
}

// HandleRepositoryPaginationDeleteAtError processes pagination errors from the repository
// when retrieving deleted user documents. It logs the error, updates the trace span,
// and returns a standardized error response.
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
//   - A slice of UserResponseDeleteAt pointers if successful, otherwise nil.
//   - A pointer to an integer representing additional pagination details, otherwise nil.
//   - A standardized ErrorResponse describing the pagination failure.
func (u *userQueryError) HandleRepositoryPaginationDeleteAtError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.UserResponseDeleteAt](u.logger, err, method, tracePrefix, span, status, user_errors.ErrUserNotFoundRes, fields...)
}

// HandleRepositorySingleError processes single-result errors from the user repository.
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
//   - A UserResponse pointer if successful, otherwise nil.
//   - A standardized ErrorResponse describing the single-result failure.
func (u *userQueryError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.UserResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.UserResponse](u.logger, err, method, tracePrefix, span, status, errResp, fields...)
}
