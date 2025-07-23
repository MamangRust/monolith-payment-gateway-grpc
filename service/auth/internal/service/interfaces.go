package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// RegistrationService defines the service layer for handling user registration.
//
//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go
type RegistrationService interface {
	// Register creates a new user account with the given registration request.
	//
	// Parameters:
	//   - ctx: the context for the operation (e.g., timeout, logging, tracing)
	//   - request: the registration request payload
	//
	// Returns:
	//   - A UserResponse if registration is successful, or an ErrorResponse if it fails.
	Register(ctx context.Context, request *requests.RegisterRequest) (*response.UserResponse, *response.ErrorResponse)
}

// LoginService defines the service layer for user authentication.
type LoginService interface {
	// Login authenticates a user using their credentials and returns a token upon success.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - request: the authentication request payload containing email and password
	//
	// Returns:
	//   - A TokenResponse if authentication is successful, or an ErrorResponse if it fails.
	Login(ctx context.Context, request *requests.AuthRequest) (*response.TokenResponse, *response.ErrorResponse)
}

// PasswordResetService handles user password recovery and reset operations.
type PasswordResetService interface {
	// ForgotPassword initiates the password reset process by sending a verification code to the user's email.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - email: the user's email address
	//
	// Returns:
	//   - true if the code was sent successfully, or an ErrorResponse if it fails.
	ForgotPassword(ctx context.Context, email string) (bool, *response.ErrorResponse)

	// ResetPassword sets a new password for the user using the provided reset token and new password.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - request: the payload containing reset token and new password
	//
	// Returns:
	//   - true if the password reset is successful, or an ErrorResponse if it fails.
	ResetPassword(ctx context.Context, request *requests.CreateResetPasswordRequest) (bool, *response.ErrorResponse)

	// VerifyCode validates the verification code sent to the user's email.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - code: the verification code to validate
	//
	// Returns:
	//   - true if the code is valid, or an ErrorResponse if the code is invalid or expired.
	VerifyCode(ctx context.Context, code string) (bool, *response.ErrorResponse)
}

// IdentifyService handles operations related to identifying and refreshing authenticated users.
type IdentifyService interface {
	// RefreshToken generates a new access token using a valid refresh token.
	//
	// Parameters:
	//   - ctx: the context for the operation (used for timeout, tracing, etc.)
	//   - token: the refresh token string
	//
	// Returns:
	//   - A new TokenResponse if the token is valid, or an ErrorResponse if the refresh fails.
	RefreshToken(ctx context.Context, token string) (*response.TokenResponse, *response.ErrorResponse)

	// GetMe retrieves the current user's profile information based on the access token.
	//
	// Parameters:
	//   - ctx: the context for the operation
	//   - token: the access token string
	//
	// Returns:
	//   - A UserResponse representing the authenticated user, or an ErrorResponse if unauthorized or failed.
	GetMe(ctx context.Context, token string) (*response.UserResponse, *response.ErrorResponse)
}
