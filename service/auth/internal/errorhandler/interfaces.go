package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// ErrorHandler interfaces define contracts for handling specific types of errors across the authentication service.
// These interfaces standardize error handling with consistent tracing, logging and response formats.

// IdentityErrorHandler handles errors related to identity verification and token management.
//
//go:generate mockgen -source=interfaces.go -destination=mocks/cache.go
type IdentityErrorHandler interface {
	// HandleInvalidTokenError processes invalid token errors during identity operations
	// Returns TokenResponse with error details and standardized ErrorResponse
	HandleInvalidTokenError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.TokenResponse, *response.ErrorResponse)

	// HandleExpiredRefreshTokenError processes expired refresh token errors
	// Returns TokenResponse with error details and standardized ErrorResponse
	HandleExpiredRefreshTokenError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.TokenResponse, *response.ErrorResponse)

	// HandleDeleteRefreshTokenError processes errors during refresh token deletion
	// Returns TokenResponse with error details and standardized ErrorResponse
	HandleDeleteRefreshTokenError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.TokenResponse, *response.ErrorResponse)

	// HandleUpdateRefreshTokenError processes errors during refresh token updates
	// Returns TokenResponse with error details and standardized ErrorResponse
	HandleUpdateRefreshTokenError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.TokenResponse, *response.ErrorResponse)

	// HandleValidateTokenError processes token validation errors
	// Returns UserResponse with error details and standardized ErrorResponse
	HandleValidateTokenError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.UserResponse, *response.ErrorResponse)

	// HandleGetMeError processes errors during user data retrieval
	// Returns UserResponse with error details and standardized ErrorResponse
	HandleGetMeError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.UserResponse, *response.ErrorResponse)

	// HandleFindByIdError processes errors during user lookup by ID
	// Returns UserResponse with error details and standardized ErrorResponse
	HandleFindByIdError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.UserResponse, *response.ErrorResponse)
}

// KafkaErrorHandler handles errors related to Kafka message processing
type KafkaErrorHandler interface {
	// HandleSendEmailForgotPassword processes errors during forgot password email sending
	// Returns boolean status and standardized ErrorResponse
	HandleSendEmailForgotPassword(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)

	// HandleSendEmailRegister processes errors during registration email sending
	// Returns UserResponse with error details and standardized ErrorResponse
	HandleSendEmailRegister(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.UserResponse, *response.ErrorResponse)

	// HandleSendEmailVerifyCode processes errors during verification code email sending
	// Returns boolean status and standardized ErrorResponse
	HandleSendEmailVerifyCode(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
}

// LoginErrorHandler handles errors during user login operations
type LoginErrorHandler interface {
	// HandleFindEmailError processes errors during email lookup for login
	// Returns TokenResponse with error details and standardized ErrorResponse
	HandleFindEmailError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.TokenResponse, *response.ErrorResponse)
}

// MarshalErrorHandler handles errors during data marshaling operations
type MarshalErrorHandler interface {
	// HandleMarshalRegisterError processes errors during registration data marshaling
	// Returns UserResponse with error details and standardized ErrorResponse
	HandleMarshalRegisterError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.UserResponse, *response.ErrorResponse)

	// HandleMarsalForgotPassword processes errors during forgot password data marshaling
	// Returns boolean status and standardized ErrorResponse
	HandleMarsalForgotPassword(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)

	// HandleMarshalVerifyCode processes errors during verification code data marshaling
	// Returns boolean status and standardized ErrorResponse
	HandleMarshalVerifyCode(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
}

// PasswordErrorHandler handles errors related to password operations
type PasswordErrorHandler interface {
	// HandlePasswordNotMatchError processes password mismatch errors
	// Returns boolean status and standardized ErrorResponse
	HandlePasswordNotMatchError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)

	// HandleHashPasswordError processes password hashing errors
	// Returns UserResponse with error details and standardized ErrorResponse
	HandleHashPasswordError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.UserResponse, *response.ErrorResponse)

	// HandleComparePasswordError processes password comparison errors
	// Returns TokenResponse with error details and standardized ErrorResponse
	HandleComparePasswordError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.TokenResponse, *response.ErrorResponse)
}

// PasswordResetErrorHandler handles errors during password reset operations
type PasswordResetErrorHandler interface {
	// HandleFindEmailError processes errors during email lookup for password reset
	// Returns boolean status and standardized ErrorResponse
	HandleFindEmailError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)

	// HandleCreateResetTokenError processes errors during reset token creation
	// Returns boolean status and standardized ErrorResponse
	HandleCreateResetTokenError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)

	// HandleFindTokenError processes errors during reset token lookup
	// Returns boolean status and standardized ErrorResponse
	HandleFindTokenError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)

	// HandleUpdatePasswordError processes errors during password updates
	// Returns boolean status and standardized ErrorResponse
	HandleUpdatePasswordError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)

	// HandleDeleteTokenError processes errors during token deletion
	// Returns boolean status and standardized ErrorResponse
	HandleDeleteTokenError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)

	// HandleUpdateVerifiedError processes errors during verification status updates
	// Returns boolean status and standardized ErrorResponse
	HandleUpdateVerifiedError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)

	// HandleVerifyCodeError processes errors during verification code validation
	// Returns boolean status and standardized ErrorResponse
	HandleVerifyCodeError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
}

// RandomStringErrorHandler handles errors during random string generation
type RandomStringErrorHandler interface {
	// HandleRandomStringErrorRegister processes errors during random string generation for registration
	// Returns UserResponse with error details and standardized ErrorResponse
	HandleRandomStringErrorRegister(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.UserResponse, *response.ErrorResponse)

	// HandleRandomStringErrorForgotPassword processes errors during random string generation for password reset
	// Returns boolean status and standardized ErrorResponse
	HandleRandomStringErrorForgotPassword(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (bool, *response.ErrorResponse)
}

// RegisterErrorHandler handles errors during user registration
type RegisterErrorHandler interface {
	// HandleAssignRoleError processes errors during role assignment
	// Returns UserResponse with error details and standardized ErrorResponse
	HandleAssignRoleError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.UserResponse, *response.ErrorResponse)

	// HandleFindRoleError processes errors during role lookup
	// Returns UserResponse with error details and standardized ErrorResponse
	HandleFindRoleError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.UserResponse, *response.ErrorResponse)

	// HandleFindEmailError processes errors during email lookup for registration
	// Returns UserResponse with error details and standardized ErrorResponse
	HandleFindEmailError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.UserResponse, *response.ErrorResponse)

	// HandleCreateUserError processes errors during user creation
	// Returns UserResponse with error details and standardized ErrorResponse
	HandleCreateUserError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.UserResponse, *response.ErrorResponse)
}

// TokenErrorHandler handles errors during token generation
type TokenErrorHandler interface {
	// HandleCreateAccessTokenError processes errors during access token creation
	// Returns TokenResponse with error details and standardized ErrorResponse
	HandleCreateAccessTokenError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.TokenResponse, *response.ErrorResponse)

	// HandleCreateRefreshTokenError processes errors during refresh token creation
	// Returns TokenResponse with error details and standardized ErrorResponse
	HandleCreateRefreshTokenError(
		err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field,
	) (*response.TokenResponse, *response.ErrorResponse)
}
