package errorhandler

import "github.com/MamangRust/monolith-payment-gateway-pkg/logger"

// ErrorHandler is a centralized error handling container that aggregates all specialized
// error handlers for different domains of the application. It provides a unified interface
// to handle various types of errors consistently across the system.
//
// The struct follows the single responsibility principle by separating error handling
// concerns into distinct components, each focusing on a specific error domain.
type ErrorHandler struct {
	// IdentityError handles errors related to user identity verification,
	// token validation, and user information retrieval
	IdentityError IdentityErrorHandler

	// KafkaError handles errors related to Kafka message processing,
	// including message publishing failures and serialization errors
	KafkaError KafkaErrorHandler

	// LoginError handles authentication failures and user login errors,
	// including invalid credentials and account lockouts
	LoginError LoginErrorHandler

	// MarshalError handles data serialization/deserialization errors,
	// including JSON/XML marshaling and unmarshaling failures
	MarshalError MarshalErrorHandler

	// PasswordError handles password-related operation failures,
	// including hashing, validation, and comparison errors
	PasswordError PasswordErrorHandler

	// PasswordResetError handles errors in the password reset flow,
	// including token generation, validation, and password update failures
	PasswordResetError PasswordResetErrorHandler

	// RandomString handles errors during cryptographically secure
	// random string generation, used for tokens and verification codes
	RandomString RandomStringErrorHandler

	// RegisterError handles user registration failures,
	// including validation errors, duplicate emails, and role assignment issues
	RegisterError RegisterErrorHandler

	// TokenError handles JWT token operation failures,
	// including token generation, signing, and verification errors
	TokenError TokenErrorHandler
}

// NewErrorHandler initializes a new ErrorHandler with all the other error handlers.
// It takes a logger as input and returns a pointer to the ErrorHandler struct.
func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		IdentityError:      NewIdentityError(logger),
		KafkaError:         NewKafkaError(logger),
		LoginError:         NewLoginError(logger),
		MarshalError:       NewMarshalError(logger),
		PasswordError:      NewPasswordError(logger),
		PasswordResetError: NewPasswordResetError(logger),
		RandomString:       NewRandomStringError(logger),
		RegisterError:      NewRegisterError(logger),
		TokenError:         NewTokenError(logger),
	}
}
