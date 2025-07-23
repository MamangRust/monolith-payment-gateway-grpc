# 📦 Package `errorhandler`

**Source Path:** `./service/auth/internal/errorhandler`

## 🧩 Types

### `ErrorHandler`

ErrorHandler is a centralized error handling container that aggregates all specialized
error handlers for different domains of the application. It provides a unified interface
to handle various types of errors consistently across the system.

The struct follows the single responsibility principle by separating error handling
concerns into distinct components, each focusing on a specific error domain.

Usage:

	errorHandler := &ErrorHandler{
	    IdentityError:      NewIdentityErrorHandler(),
	    KafkaError:         NewKafkaErrorHandler(),
	    // ... initialize other handlers
	}

```go
type ErrorHandler struct {
	IdentityError IdentityErrorHandler
	KafkaError KafkaErrorHandler
	LoginError LoginErrorHandler
	MarshalError MarshalErrorHandler
	PasswordError PasswordErrorHandler
	PasswordResetError PasswordResetErrorHandler
	RandomString RandomStringErrorHandler
	RegisterError RegisterErrorHandler
	TokenError TokenErrorHandler
}
```

### `IdentityErrorHandler`

IdentityErrorHandler handles errors related to identity verification and token management.

```go
type IdentityErrorHandler interface {
	HandleInvalidTokenError func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TokenResponse, *response.ErrorResponse)
	HandleExpiredRefreshTokenError func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TokenResponse, *response.ErrorResponse)
	HandleDeleteRefreshTokenError func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TokenResponse, *response.ErrorResponse)
	HandleUpdateRefreshTokenError func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TokenResponse, *response.ErrorResponse)
	HandleValidateTokenError func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.UserResponse, *response.ErrorResponse)
	HandleGetMeError func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.UserResponse, *response.ErrorResponse)
	HandleFindByIdError func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.UserResponse, *response.ErrorResponse)
}
```

### `KafkaErrorHandler`

KafkaErrorHandler handles errors related to Kafka message processing

```go
type KafkaErrorHandler interface {
	HandleSendEmailForgotPassword func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleSendEmailRegister func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.UserResponse, *response.ErrorResponse)
	HandleSendEmailVerifyCode func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
}
```

### `LoginErrorHandler`

LoginErrorHandler handles errors during user login operations

```go
type LoginErrorHandler interface {
	HandleFindEmailError func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TokenResponse, *response.ErrorResponse)
}
```

### `MarshalErrorHandler`

MarshalErrorHandler handles errors during data marshaling operations

```go
type MarshalErrorHandler interface {
	HandleMarshalRegisterError func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.UserResponse, *response.ErrorResponse)
	HandleMarsalForgotPassword func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleMarshalVerifyCode func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
}
```

### `PasswordErrorHandler`

PasswordErrorHandler handles errors related to password operations

```go
type PasswordErrorHandler interface {
	HandlePasswordNotMatchError func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleHashPasswordError func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.UserResponse, *response.ErrorResponse)
	HandleComparePasswordError func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TokenResponse, *response.ErrorResponse)
}
```

### `PasswordResetErrorHandler`

PasswordResetErrorHandler handles errors during password reset operations

```go
type PasswordResetErrorHandler interface {
	HandleFindEmailError func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleCreateResetTokenError func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleFindTokenError func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleUpdatePasswordError func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleDeleteTokenError func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleUpdateVerifiedError func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleVerifyCodeError func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
}
```

### `RandomStringErrorHandler`

RandomStringErrorHandler handles errors during random string generation

```go
type RandomStringErrorHandler interface {
	HandleRandomStringErrorRegister func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.UserResponse, *response.ErrorResponse)
	HandleRandomStringErrorForgotPassword func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
}
```

### `RegisterErrorHandler`

RegisterErrorHandler handles errors during user registration

```go
type RegisterErrorHandler interface {
	HandleAssignRoleError func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.UserResponse, *response.ErrorResponse)
	HandleFindRoleError func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.UserResponse, *response.ErrorResponse)
	HandleFindEmailError func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.UserResponse, *response.ErrorResponse)
	HandleCreateUserError func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.UserResponse, *response.ErrorResponse)
}
```

### `TokenErrorHandler`

TokenErrorHandler handles errors during token generation

```go
type TokenErrorHandler interface {
	HandleCreateAccessTokenError func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TokenResponse, *response.ErrorResponse)
	HandleCreateRefreshTokenError func(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TokenResponse, *response.ErrorResponse)
}
```

### `identityError`

```go
type identityError struct {
	logger logger.LoggerInterface
}
```

#### Methods

##### `HandleDeleteRefreshTokenError`

HandleDeleteRefreshTokenError processes errors during refresh token deletion
It logs the error, records it to the trace span, and returns a standardized error response.
Parameters:
  - err: The error that occurred.
  - method: The name of the method where the error occurred.
  - tracePrefix: A prefix for generating the trace ID.
  - span: The trace span used for recording the error.
  - status: A pointer to a string that will be set with the formatted status.
  - fields: Additional fields to include in the log entry.

Returns:
  - A TokenResponse with error details and a standardized ErrorResponse.

```go
func (e *identityError) HandleDeleteRefreshTokenError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TokenResponse, *response.ErrorResponse)
```

##### `HandleExpiredRefreshTokenError`

HandleExpiredRefreshTokenError processes errors related to expired refresh tokens during identity operations.
It logs the error, records it to the trace span, and returns a standardized error response.
Parameters:
  - err: The error that occurred.
  - method: The name of the method where the error occurred.
  - tracePrefix: A prefix for generating the trace ID.
  - span: The trace span used for recording the error.
  - status: A pointer to a string that will be set with the formatted status.
  - fields: Additional fields to include in the log entry.

Returns:
  - A TokenResponse with error details and a standardized ErrorResponse.

```go
func (e *identityError) HandleExpiredRefreshTokenError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TokenResponse, *response.ErrorResponse)
```

##### `HandleFindByIdError`

HandleFindByIdError processes errors during user lookup by ID
It logs the error, records it to the trace span, and returns a standardized error response.
Parameters:
  - err: The error that occurred.
  - method: The name of the method where the error occurred.
  - tracePrefix: A prefix for generating the trace ID.
  - span: The trace span used for recording the error.
  - status: A pointer to a string that will be set with the formatted status.
  - fields: Additional fields to include in the log entry.

Returns:
  - A UserResponse with error details and a standardized ErrorResponse.

```go
func (e *identityError) HandleFindByIdError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.UserResponse, *response.ErrorResponse)
```

##### `HandleGetMeError`

HandleGetMeError processes errors during user data retrieval
It logs the error, records it to the trace span, and returns a standardized error response.
Parameters:
  - err: The error that occurred.
  - method: The name of the method where the error occurred.
  - tracePrefix: A prefix for generating the trace ID.
  - span: The trace span used for recording the error.
  - status: A pointer to a string that will be set with the formatted status.
  - fields: Additional fields to include in the log entry.

Returns:
  - A UserResponse with error details and a standardized ErrorResponse.

```go
func (e *identityError) HandleGetMeError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.UserResponse, *response.ErrorResponse)
```

##### `HandleInvalidTokenError`

HandleInvalidTokenError processes errors related to invalid tokens during identity operations.
It logs the error, records it to the trace span, and returns a standardized error response.
Parameters:
  - err: The error that occurred.
  - method: The name of the method where the error occurred.
  - tracePrefix: A prefix for generating the trace ID.
  - span: The trace span used for recording the error.
  - status: A pointer to a string that will be set with the formatted status.
  - fields: Additional fields to include in the log entry.

Returns:
  - A TokenResponse with error details and a standardized ErrorResponse.

```go
func (e *identityError) HandleInvalidTokenError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TokenResponse, *response.ErrorResponse)
```

##### `HandleUpdateRefreshTokenError`

HandleUpdateRefreshTokenError processes errors during refresh token updates.
It logs the error, records it to the trace span, and returns a standardized error response.
Parameters:
  - err: The error that occurred.
  - method: The name of the method where the error occurred.
  - tracePrefix: A prefix for generating the trace ID.
  - span: The trace span used for recording the error.
  - status: A pointer to a string that will be set with the formatted status.
  - fields: Additional fields to include in the log entry.

Returns:
  - A TokenResponse with error details and a standardized ErrorResponse.

```go
func (e *identityError) HandleUpdateRefreshTokenError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TokenResponse, *response.ErrorResponse)
```

##### `HandleValidateTokenError`

HandleValidateTokenError processes token validation errors
It logs the error, records it to the trace span, and returns a standardized error response.
Parameters:
  - err: The error that occurred.
  - method: The name of the method where the error occurred.
  - tracePrefix: A prefix for generating the trace ID.
  - span: The trace span used for recording the error.
  - status: A pointer to a string that will be set with the formatted status.
  - fields: Additional fields to include in the log entry.

Returns:
  - A UserResponse with error details and a standardized ErrorResponse.

```go
func (e *identityError) HandleValidateTokenError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.UserResponse, *response.ErrorResponse)
```

### `kafkaError`

```go
type kafkaError struct {
	logger logger.LoggerInterface
}
```

#### Methods

##### `HandleSendEmailForgotPassword`

HandleSendEmailForgotPassword processes errors that occur during the sending of forgot password emails.
It utilizes Kafka for message handling and returns a boolean indicating success or failure, along with a standardized ErrorResponse.
Parameters:
  - err: The error encountered during the email sending process.
  - method: The name of the method where the error occurred.
  - tracePrefix: A prefix for trace logging.
  - span: The tracing span for monitoring.
  - status: A pointer to a string representing the status of the operation.
  - fields: Additional logging fields for structured logging.

Returns:
  - A boolean indicating success (false) or failure (true) of the operation.
  - A pointer to a standardized ErrorResponse containing error details.

```go
func (e *kafkaError) HandleSendEmailForgotPassword(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
```

##### `HandleSendEmailRegister`

HandleSendEmailRegister processes errors that occur during the sending of registration emails.
It utilizes Kafka for message handling and returns a UserResponse with error details and a standardized ErrorResponse.
Parameters:
  - err: The error encountered during the email sending process.
  - method: The name of the method where the error occurred.
  - tracePrefix: A prefix for trace logging.
  - span: The tracing span for monitoring.
  - status: A pointer to a string representing the status of the operation.
  - fields: Additional logging fields for structured logging.

Returns:
  - A UserResponse containing user-related information if available.
  - A pointer to a standardized ErrorResponse containing error details.

```go
func (e *kafkaError) HandleSendEmailRegister(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.UserResponse, *response.ErrorResponse)
```

##### `HandleSendEmailVerifyCode`

HandleSendEmailVerifyCode processes errors that occur during the sending of verification code emails.
It utilizes Kafka for message handling and returns a boolean indicating success or failure, along with a standardized ErrorResponse.
Parameters:
  - err: The error encountered during the email sending process.
  - method: The name of the method where the error occurred.
  - tracePrefix: A prefix for trace logging.
  - span: The tracing span for monitoring.
  - status: A pointer to a string representing the status of the operation.
  - fields: Additional logging fields for structured logging.

Returns:
  - A boolean indicating success (false) or failure (true) of the operation.
  - A pointer to a standardized ErrorResponse containing error details.

```go
func (e *kafkaError) HandleSendEmailVerifyCode(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
```

### `loginError`

```go
type loginError struct {
	logger logger.LoggerInterface
}
```

#### Methods

##### `HandleFindEmailError`

HandleFindEmailError processes errors encountered during the email lookup
for login operations. It logs the error, records it to the trace span,
and returns a standardized error response.

Parameters:
- err: The error that occurred during email lookup.
- method: The name of the method where the error occurred.
- tracePrefix: A prefix for generating the trace ID.
- span: The trace span used for recording the error.
- status: A pointer to a string that will be set with the formatted status.
- fields: Additional fields to include in the log entry.

Returns:
- A TokenResponse with error details and a standardized ErrorResponse.

```go
func (e *loginError) HandleFindEmailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TokenResponse, *response.ErrorResponse)
```

### `marshalError`

```go
type marshalError struct {
	logger logger.LoggerInterface
}
```

#### Methods

##### `HandleMarsalForgotPassword`

HandleMarsalForgotPassword processes errors during forgot password data marshaling
Returns boolean status and standardized ErrorResponse

```go
func (e *marshalError) HandleMarsalForgotPassword(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
```

##### `HandleMarshalRegisterError`

HandleMarshalRegisterError processes errors during registration data marshaling
Returns UserResponse with error details and standardized ErrorResponse

```go
func (e *marshalError) HandleMarshalRegisterError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.UserResponse, *response.ErrorResponse)
```

##### `HandleMarshalVerifyCode`

HandleMarshalVerifyCode processes errors during verification code data marshaling
Returns boolean status and standardized ErrorResponse

```go
func (e *marshalError) HandleMarshalVerifyCode(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
```

### `passwordError`

```go
type passwordError struct {
	logger logger.LoggerInterface
}
```

#### Methods

##### `HandleComparePasswordError`

HandleComparePasswordError processes password comparison errors
Returns TokenResponse with error details and standardized ErrorResponse
Parameters:
  - err: The error that occurred during password comparison.
  - method: The name of the method where the error occurred.
  - tracePrefix: A prefix for generating the trace ID.
  - span: The trace span used for recording the error.
  - status: A pointer to a string that will be set with the formatted status.
  - fields: Additional fields to include in the log entry.

```go
func (e *passwordError) HandleComparePasswordError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TokenResponse, *response.ErrorResponse)
```

##### `HandleHashPasswordError`

HandleHashPasswordError processes password hashing errors
Returns UserResponse with error details and standardized ErrorResponse
Parameters:
  - err: The error that occurred during password comparison.
  - method: The name of the method where the error occurred.
  - tracePrefix: A prefix for generating the trace ID.
  - span: The trace span used for recording the error.
  - status: A pointer to a string that will be set with the formatted status.
  - fields: Additional fields to include in the log entry.

```go
func (e *passwordError) HandleHashPasswordError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.UserResponse, *response.ErrorResponse)
```

##### `HandlePasswordNotMatchError`

HandlePasswordNotMatchError processes password mismatch errors
Returns boolean status and standardized ErrorResponse
Parameters:
  - err: The error that occurred during password comparison.
  - method: The name of the method where the error occurred.
  - tracePrefix: A prefix for generating the trace ID.
  - span: The trace span used for recording the error.
  - status: A pointer to a string that will be set with the formatted status.
  - fields: Additional fields to include in the log entry.

```go
func (e *passwordError) HandlePasswordNotMatchError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
```

### `passwordResetError`

```go
type passwordResetError struct {
	logger logger.LoggerInterface
}
```

#### Methods

##### `HandleCreateResetTokenError`

HandleCreateResetTokenError processes errors that occur during reset token generation.

Parameters:
  - err: The error that occurred during token creation (error)
  - method: The name of the calling method (e.g., "GenerateResetToken") (string)
  - tracePrefix: A prefix used for generating trace IDs (e.g., "GEN_RESET_TOKEN") (string)
  - span: The OpenTelemetry span for distributed tracing (trace.Span)
  - status: Pointer to a string that will be updated with the error status (*string)
  - fields: Additional context fields for structured logging (...zap.Field)

Returns:
  - bool: Always returns false indicating failure
  - *response.ErrorResponse: Standardized error response with details

```go
func (e *passwordResetError) HandleCreateResetTokenError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
```

##### `HandleDeleteTokenError`

HandleDeleteTokenError processes errors that occur when deleting used reset tokens.

Parameters:
  - err: The error that occurred during token deletion (error)
  - method: The name of the calling method (e.g., "CleanupResetToken") (string)
  - tracePrefix: A prefix used for generating trace IDs (e.g., "CLEANUP_TOKEN") (string)
  - span: The OpenTelemetry span for distributed tracing (trace.Span)
  - status: Pointer to a string that will be updated with the error status (*string)
  - fields: Additional context fields for structured logging (...zap.Field)

Returns:
  - bool: Always returns false indicating failure
  - *response.ErrorResponse: Standardized error response with details

```go
func (e *passwordResetError) HandleDeleteTokenError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
```

##### `HandleFindEmailError`

HandleFindEmailError processes errors that occur when looking up a user's email during password reset.

Parameters:
  - err: The error that occurred during email lookup (error)
  - method: The name of the calling method (e.g., "InitiatePasswordReset") (string)
  - tracePrefix: A prefix used for generating trace IDs (e.g., "INIT_PW_RESET") (string)
  - span: The OpenTelemetry span for distributed tracing (trace.Span)
  - status: Pointer to a string that will be updated with the error status (*string)
  - fields: Additional context fields for structured logging (...zap.Field)

Returns:
  - bool: Always returns false indicating failure
  - *response.ErrorResponse: Standardized error response with details

```go
func (e *passwordResetError) HandleFindEmailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
```

##### `HandleFindTokenError`

HandleFindTokenError processes errors that occur when looking up a reset token.

Parameters:
  - err: The error that occurred during token lookup (error)
  - method: The name of the calling method (e.g., "ValidateResetToken") (string)
  - tracePrefix: A prefix used for generating trace IDs (e.g., "VALIDATE_TOKEN") (string)
  - span: The OpenTelemetry span for distributed tracing (trace.Span)
  - status: Pointer to a string that will be updated with the error status (*string)
  - fields: Additional context fields for structured logging (...zap.Field)

Returns:
  - bool: Always returns false indicating failure
  - *response.ErrorResponse: Standardized error response with details

```go
func (e *passwordResetError) HandleFindTokenError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
```

##### `HandleUpdatePasswordError`

HandleUpdatePasswordError processes errors that occur during password updates.

Parameters:
  - err: The error that occurred during password update (error)
  - method: The name of the calling method (e.g., "CompletePasswordReset") (string)
  - tracePrefix: A prefix used for generating trace IDs (e.g., "COMPLETE_PW_RESET") (string)
  - span: The OpenTelemetry span for distributed tracing (trace.Span)
  - status: Pointer to a string that will be updated with the error status (*string)
  - fields: Additional context fields for structured logging (...zap.Field)

Returns:
  - bool: Always returns false indicating failure
  - *response.ErrorResponse: Standardized error response with details

```go
func (e *passwordResetError) HandleUpdatePasswordError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
```

##### `HandleUpdateVerifiedError`

HandleUpdateVerifiedError processes errors that occur when updating verification status.

Parameters:
  - err: The error that occurred during status update (error)
  - method: The name of the calling method (e.g., "MarkAsVerified") (string)
  - tracePrefix: A prefix used for generating trace IDs (e.g., "MARK_VERIFIED") (string)
  - span: The OpenTelemetry span for distributed tracing (trace.Span)
  - status: Pointer to a string that will be updated with the error status (*string)
  - fields: Additional context fields for structured logging (...zap.Field)

Returns:
  - bool: Always returns false indicating failure
  - *response.ErrorResponse: Standardized error response with details

```go
func (e *passwordResetError) HandleUpdateVerifiedError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
```

##### `HandleVerifyCodeError`

HandleVerifyCodeError processes errors that occur during verification code validation.

Parameters:
  - err: The error that occurred during code validation (error)
  - method: The name of the calling method (e.g., "CheckVerificationCode") (string)
  - tracePrefix: A prefix used for generating trace IDs (e.g., "CHECK_VERIFY_CODE") (string)
  - span: The OpenTelemetry span for distributed tracing (trace.Span)
  - status: Pointer to a string that will be updated with the error status (*string)
  - fields: Additional context fields for structured logging (...zap.Field)

Returns:
  - bool: Always returns false indicating failure
  - *response.ErrorResponse: Standardized error response with details

```go
func (e *passwordResetError) HandleVerifyCodeError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
```

### `randomStringError`

```go
type randomStringError struct {
	logger logger.LoggerInterface
}
```

#### Methods

##### `HandleRandomStringErrorForgotPassword`

HandleRandomStringErrorForgotPassword processes errors that occur during random string generation
for forgot password operations. It leverages handleErrorGenerateRandomString to log the error and
update the trace span with relevant error details.

Parameters:
  - err: The error that occurred during random string generation.
  - method: The name of the method where the error occurred.
  - tracePrefix: A prefix for generating the trace ID.
  - span: The OpenTelemetry span for distributed tracing.
  - status: A pointer to a string that will be updated with the error status.
  - fields: Additional context fields for structured logging.

Returns:
  - A boolean indicating whether the operation was successful (false) or not (true).
  - A standardized ErrorResponse containing error details.

```go
func (h *randomStringError) HandleRandomStringErrorForgotPassword(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
```

##### `HandleRandomStringErrorRegister`

HandleRandomStringErrorRegister processes errors that occur during random string generation
for user registration. It leverages handleErrorGenerateRandomString to log the error and update
the trace span with relevant error details.

Parameters:
  - err: The error that occurred during random string generation.
  - method: The name of the method where the error occurred.
  - tracePrefix: A prefix for generating the trace ID.
  - span: The OpenTelemetry span for distributed tracing.
  - status: A pointer to a string that will be updated with the error status.
  - fields: Additional context fields for structured logging.

Returns:
  - A UserResponse with user-related information if available.
  - A standardized ErrorResponse containing error details.

```go
func (r randomStringError) HandleRandomStringErrorRegister(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.UserResponse, *response.ErrorResponse)
```

### `registerError`

```go
type registerError struct {
	logger logger.LoggerInterface
}
```

#### Methods

##### `HandleAssignRoleError`

HandleAssignRoleError processes errors that occur during role assignment to a new user.
This typically happens after successful user creation but before completing registration.

Parameters:
  - err: The error that occurred during role assignment (error)
  - method: The name of the calling method (e.g., "CompleteRegistration") (string)
  - tracePrefix: A prefix used for generating trace IDs (e.g., "COMPLETE_REG") (string)
  - span: The OpenTelemetry span for distributed tracing (trace.Span)
  - status: Pointer to a string that will be updated with error status (e.g., "complete_reg_error_assign_role") (*string)
  - fields: Additional context fields for structured logging (...zap.Field)

Returns:
  - *response.UserResponse: Nil user response since operation failed
  - *response.ErrorResponse: Standardized error response with user_not_found error details

```go
func (e *registerError) HandleAssignRoleError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.UserResponse, *response.ErrorResponse)
```

##### `HandleCreateUserError`

HandleCreateUserError processes errors that occur during core user creation in the registration flow.
This handles failures in the primary user record creation operation.

Parameters:
  - err: The error that occurred during user creation (error)
  - method: The name of the calling method (e.g., "CreateUserRecord") (string)
  - tracePrefix: A prefix used for generating trace IDs (e.g., "CREATE_USER") (string)
  - span: The OpenTelemetry span for distributed tracing (trace.Span)
  - status: Pointer to a string that will be updated with error status (*string)
  - fields: Additional context fields for structured logging (...zap.Field)

Returns:
  - *response.UserResponse: Nil user response since operation failed
  - *response.ErrorResponse: Standardized error response with user_not_found error details

```go
func (e *registerError) HandleCreateUserError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.UserResponse, *response.ErrorResponse)
```

##### `HandleFindEmailError`

HandleFindEmailError processes errors that occur when checking for email existence during registration.
Used to prevent duplicate email registrations.

Parameters:
  - err: The error that occurred during email lookup (error)
  - method: The name of the calling method (e.g., "CheckEmailAvailability") (string)
  - tracePrefix: A prefix used for generating trace IDs (e.g., "CHECK_EMAIL") (string)
  - span: The OpenTelemetry span for distributed tracing (trace.Span)
  - status: Pointer to a string that will be updated with error status (*string)
  - fields: Additional context fields for structured logging (...zap.Field)

Returns:
  - *response.UserResponse: Nil user response since operation failed
  - *response.ErrorResponse: Standardized error response with user_not_found error details

```go
func (e *registerError) HandleFindEmailError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.UserResponse, *response.ErrorResponse)
```

##### `HandleFindRoleError`

HandleFindRoleError processes errors that occur when looking up role information during registration.
Used when assigning default or requested roles to new users.

Parameters:
  - err: The error that occurred during role lookup (error)
  - method: The name of the calling method (e.g., "AssignDefaultRole") (string)
  - tracePrefix: A prefix used for generating trace IDs (e.g., "ASSIGN_ROLE") (string)
  - span: The OpenTelemetry span for distributed tracing (trace.Span)
  - status: Pointer to a string that will be updated with error status (*string)
  - fields: Additional context fields for structured logging (...zap.Field)

Returns:
  - *response.UserResponse: Nil user response since operation failed
  - *response.ErrorResponse: Standardized error response with role_not_found error details

```go
func (e *registerError) HandleFindRoleError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.UserResponse, *response.ErrorResponse)
```

### `tokenError`

```go
type tokenError struct {
	logger logger.LoggerInterface
}
```

#### Methods

##### `HandleCreateAccessTokenError`

HandleCreateAccessTokenError processes errors that occur during access token generation.
This includes JWT signing errors, claims validation failures, and token encoding issues.

Parameters:
  - err: The error that occurred during access token creation (error)
  - method: The name of the calling method (e.g., "GenerateAccessToken") (string)
  - tracePrefix: A prefix used for generating trace IDs (e.g., "GEN_ACCESS_TOKEN") (string)
  - span: The OpenTelemetry span for distributed tracing (trace.Span)
  - status: Pointer to a string that will be updated with error status (e.g., "gen_access_token_error_create_token") (*string)
  - fields: Additional context fields for structured logging (...zap.Field)

Returns:
  - *response.TokenResponse: Nil token response since operation failed
  - *response.ErrorResponse: Standardized error response with failed_create_access error details,
    typically containing error code 500 (Internal Server Error)

```go
func (e *tokenError) HandleCreateAccessTokenError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TokenResponse, *response.ErrorResponse)
```

##### `HandleCreateRefreshTokenError`

HandleCreateRefreshTokenError processes errors that occur during refresh token generation.
This handles failures in token persistence, cryptographic operations, and storage errors.

Parameters:
  - err: The error that occurred during refresh token creation (error)
  - method: The name of the calling method (e.g., "GenerateRefreshToken") (string)
  - tracePrefix: A prefix used for generating trace IDs (e.g., "GEN_REFRESH_TOKEN") (string)
  - span: The OpenTelemetry span for distributed tracing (trace.Span)
  - status: Pointer to a string that will be updated with error status (*string)
  - fields: Additional context fields for structured logging (...zap.Field)

Returns:
  - *response.TokenResponse: Nil token response since operation failed
  - *response.ErrorResponse: Standardized error response with failed_create_refresh error details,
    typically containing error code 500 (Internal Server Error)

```go
func (e *tokenError) HandleCreateRefreshTokenError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.TokenResponse, *response.ErrorResponse)
```

## 🚀 Functions

### `HandleInvalidFormatUserIDError`

HandleInvalidFormatUserIDError handles errors due to invalid user ID formats.
It logs the error, updates the trace span, and returns a zero value of the specified type T
along with a standardized ErrorResponse.

Parameters:
  - logger: The logger instance for error logging (logger.LoggerInterface)
  - err: The error encountered due to invalid user ID format (error)
  - method: The name of the method where the error occurred (string)
  - tracePrefix: Prefix used for generating trace IDs (string)
  - span: The OpenTelemetry span for distributed tracing (trace.Span)
  - status: Pointer to a string to be updated with the error status (*string)
  - fields: Additional context fields for structured logging (...zap.Field)

Returns:
  - A zero value of the specified type T
  - A pointer to the error response (*response.ErrorResponse)

```go
func HandleInvalidFormatUserIDError[T any](logger logger.LoggerInterface, err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (T, *response.ErrorResponse)
```

### `HandleRepositorySingleError`

HandleRepositorySingleError is the public interface for single-result repository errors.
It provides standardized handling of database operation failures.

Parameters:
  - logger: Structured logger
  - err: Database error
  - method: Calling method
  - tracePrefix: Trace prefix
  - span: Tracing span
  - status: Status reference
  - defaultErr: Default error
  - fields: Context fields

Returns:
  - Zero value and error response

```go
func HandleRepositorySingleError[T any](logger logger.LoggerInterface, err error, method, tracePrefix string, span trace.Span, status *string, defaultErr *response.ErrorResponse, fields ...zap.Field) (T, *response.ErrorResponse)
```

### `HandleTokenError`

HandleTokenError is a helper function used to process errors related to token operations.
It logs the error using the provided logger, records the error to the trace span,
and sets the status to a standardized format. It returns a zero value of the specified type T
and a pointer to a standardized ErrorResponse.

Parameters:

	logger - LoggerInterface used for logging the error.
	err - The error that occurred.
	method - The name of the method where the error occurred.
	tracePrefix - A prefix for generating the trace ID.
	span - The trace span used for recording the error.
	status - A pointer to a string that will be set with the formatted status.
	defaultErr - The default error response to return.
	fields - Additional fields to include in the log entry.

Returns:

	A zero value of type T and a pointer to an ErrorResponse.

```go
func HandleTokenError[T any](logger logger.LoggerInterface, err error, method, tracePrefix string, span trace.Span, status *string, defaultErr *response.ErrorResponse, fields ...zap.Field) (T, *response.ErrorResponse)
```

### `handleErrorGenerateRandomString`

handleErrorGenerateRandomString handles errors that occur during random string generation operations.
It wraps the error using the standard error template with a predefined "generate random string" error message.

Parameters:
  - logger: The logger instance for recording error details (logger.LoggerInterface)
  - err: The error that occurred during random string generation (error)
  - method: The name of the method where the error occurred (e.g., "CreateVerificationCode") (string)
  - tracePrefix: The prefix for generating trace IDs (e.g., "CREATE_VERIFICATION_CODE") (string)
  - span: The OpenTelemetry span for distributed tracing (trace.Span)
  - status: Pointer to a string that will be updated with the error status (e.g., "create_verification_code_error_generate_random_string") (*string)
  - defaultErr: The predefined error response to return (*response.ErrorResponse)
  - fields: Additional context fields for structured logging (...zap.Field)

Returns:
  - A zero value of the specified type T
  - A pointer to the error response (*response.ErrorResponse)

```go
func handleErrorGenerateRandomString[T any](logger logger.LoggerInterface, err error, method, tracePrefix string, span trace.Span, status *string, defaultErr *response.ErrorResponse, fields ...zap.Field) (T, *response.ErrorResponse)
```

### `handleErrorInvalidID`

handleErrorInvalidID handles errors related to invalid ID formats or values.
It wraps the error using the standard error template with a predefined "invalid id" error message.

Parameters:
  - logger: The logger instance for recording error details (logger.LoggerInterface)
  - err: The error that occurred due to invalid ID (error)
  - method: The name of the method where the error occurred (e.g., "GetUserByID") (string)
  - tracePrefix: The prefix for generating trace IDs (e.g., "GET_USER_BY_ID") (string)
  - span: The OpenTelemetry span for distributed tracing (trace.Span)
  - status: Pointer to a string that will be updated with the error status (e.g., "get_user_by_id_error_invalid_id") (*string)
  - defaultErr: The predefined error response to return (*response.ErrorResponse)
  - fields: Additional context fields for structured logging (...zap.Field)

Returns:
  - A zero value of the specified type T
  - A pointer to the error response (*response.ErrorResponse)

```go
func handleErrorInvalidID[T any](logger logger.LoggerInterface, err error, method, tracePrefix string, span trace.Span, status *string, defaultErr *response.ErrorResponse, fields ...zap.Field) (T, *response.ErrorResponse)
```

### `handleErrorJSONMarshal`

handleErrorJSONMarshal specializes error handling for JSON marshaling failures.

Parameters:
  - logger: Logger instance
  - err: Marshaling error
  - method: Calling method name
  - tracePrefix: Trace prefix
  - span: Tracing span
  - status: Status reference
  - defaultErr: Default error
  - fields: Log fields

Returns:
  - Zero value and error response

```go
func handleErrorJSONMarshal[T any](logger logger.LoggerInterface, err error, method, tracePrefix string, span trace.Span, status *string, defaultErr *response.ErrorResponse, fields ...zap.Field) (T, *response.ErrorResponse)
```

### `handleErrorKafkaSend`

handleErrorKafkaSend specializes error handling for Kafka producer failures.

Parameters:
  - logger: Logger instance
  - err: Kafka send error
  - method: Calling method
  - tracePrefix: Trace prefix
  - span: Tracing span
  - status: Status reference
  - defaultErr: Default error
  - fields: Log fields

Returns:
  - Zero value and error response

```go
func handleErrorKafkaSend[T any](logger logger.LoggerInterface, err error, method, tracePrefix string, span trace.Span, status *string, defaultErr *response.ErrorResponse, fields ...zap.Field) (T, *response.ErrorResponse)
```

### `handleErrorPasswordOperation`

handleErrorPasswordOperation handles errors related to password operations (hashing, validation, etc.).
It allows specifying a custom operation name in the error message for more context.

Parameters:
  - logger: The logger instance for recording error details (logger.LoggerInterface)
  - err: The error that occurred during password operation (error)
  - method: The name of the method where the error occurred (e.g., "ChangePassword") (string)
  - tracePrefix: The prefix for generating trace IDs (e.g., "CHANGE_PASSWORD") (string)
  - operation: The specific password operation that failed (e.g., "hashing", "validation") (string)
  - span: The OpenTelemetry span for distributed tracing (trace.Span)
  - status: Pointer to a string that will be updated with the error status (e.g., "change_password_error_hashing") (*string)
  - defaultErr: The predefined error response to return (*response.ErrorResponse)
  - fields: Additional context fields for structured logging (...zap.Field)

Returns:
  - A zero value of the specified type T
  - A pointer to the error response (*response.ErrorResponse)

```go
func handleErrorPasswordOperation[T any](logger logger.LoggerInterface, err error, method, tracePrefix, operation string, span trace.Span, status *string, defaultErr *response.ErrorResponse, fields ...zap.Field) (T, *response.ErrorResponse)
```

### `handleErrorRepository`

handleErrorRepository specializes handleErrorTemplate for repository layer errors.
It automatically sets the error message to "repository error" and follows the
same standardized error handling pattern.

Parameters:
  - logger: LoggerInterface instance for structured logging
  - err: The error from repository operation
  - method: Name of the calling method
  - tracePrefix: Prefix for trace ID generation
  - span: OpenTelemetry span for tracing
  - status: Pointer to status string to be updated
  - errorResp: Predefined error response
  - fields: Additional contextual log fields

Returns:
  - Zero value of type T
  - Pointer to response.ErrorResponse

```go
func handleErrorRepository[T any](logger logger.LoggerInterface, err error, method, tracePrefix string, span trace.Span, status *string, errorResp *response.ErrorResponse, fields ...zap.Field) (T, *response.ErrorResponse)
```

### `handleErrorTemplate`

Package errorhandler provides standardized error handling utilities for the application.
It includes templates for handling various types of errors with consistent logging,
tracing, and response formatting across the codebase.
handleErrorTemplate is a generic function for handling errors with predefined error response templates.
Parameters:
  - logger: LoggerInterface instance for structured logging
  - err: The error that occurred
  - method: Name of the method where error occurred (e.g., "GetUser")
  - tracePrefix: Prefix for trace ID generation (e.g., "GET_USER")
  - errorMessage: Descriptive error message for logging
  - span: OpenTelemetry span for distributed tracing
  - status: Pointer to status string that will be updated
  - errorResp: Predefined error response template
  - fields: Additional zap fields for contextual logging

Returns:
  - Zero value of type T
  - Pointer to response.ErrorResponse

```go
func handleErrorTemplate[T any](logger logger.LoggerInterface, err error, method, tracePrefix, errorMessage string, span trace.Span, status *string, errorResp *response.ErrorResponse, fields ...zap.Field) (T, *response.ErrorResponse)
```

### `handleErrorTokenTemplate`

handleErrorTokenTemplate specializes handleErrorTemplate for token-related errors.
It automatically sets the error message to "token error" and follows the
standardized error handling pattern for authentication/authorization failures.

Parameters:
  - logger: LoggerInterface instance
  - err: The token-related error
  - method: Name of the calling method
  - tracePrefix: Trace ID prefix
  - span: OpenTelemetry span
  - status: Pointer to status string
  - defaultErr: Default error response
  - fields: Additional log fields

Returns:
  - Zero value of type T
  - Pointer to response.ErrorResponse

```go
func handleErrorTokenTemplate[T any](logger logger.LoggerInterface, err error, method, tracePrefix string, span trace.Span, status *string, defaultErr *response.ErrorResponse, fields ...zap.Field) (T, *response.ErrorResponse)
```

### `toSnakeCase`

toSnakeCase converts a camelCase string to a snake_case string.

Parameters:

  - s: CamelCase string

Returns:

  - Snake case equivalent of the input string.

```go
func toSnakeCase(s string) string
```

