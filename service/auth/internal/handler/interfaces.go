package handler

import (
	"context"

	pb "github.com/MamangRust/monolith-payment-gateway-pb/auth"
)

// AuthHandleGrpc defines the gRPC service interface for authentication operations.
// It combines the generated gRPC server interface with concrete method signatures
// for all authentication-related operations.
//
// This interface serves as the contract for the gRPC authentication service handler,
// providing methods for user authentication, registration, password management,
// and token operations.
//
//go:generate mockgen -source=interfaces.go -destination=mocks/handler.go
type AuthHandleGrpc interface {
	pb.AuthServiceServer // Embeds the generated gRPC service interface

	// LoginUser authenticates a user and returns access tokens.
	// It logs the operation's start and end, and utilizes the loginService
	// to perform the authentication. If the service indicates a failure, an error
	// is logged and returned. On success, a success message is logged, and a
	// protobuf response is returned indicating successful authentication and containing
	// tokens.
	//
	// Parameters:
	//   - ctx: The context for managing request-scoped values, cancelation signals, and deadlines.
	//   - req: A pointer to LoginRequest containing the user email and password.
	//
	// Returns:
	//   - A pointer to ApiResponseLogin containing the authentication result on success.
	//   - An error if the authentication process fails.
	LoginUser(ctx context.Context, req *pb.LoginRequest) (*pb.ApiResponseLogin, error)

	// RegisterUser creates a new user account.
	// It logs the operation's start and end, and utilizes the registerService
	// to perform the registration. If the service indicates a failure, an error
	// is logged and returned. On success, a success message is logged, and a
	// protobuf response is returned indicating successful registration.
	//
	// Parameters:
	//   - ctx: The context for managing request-scoped values, cancelation signals, and deadlines.
	//   - req: A pointer to RegisterRequest containing the user registration details.
	//
	// Returns:
	//   - A pointer to ApiResponseRegister containing the registration result on success.
	//   - An error if the registration process fails.
	RegisterUser(ctx context.Context, req *pb.RegisterRequest) (*pb.ApiResponseRegister, error)

	// ForgotPassword initiates the password reset process.
	// It logs the operation's start and end, and utilizes the passwordResetService
	// to perform the password reset. If the service indicates a failure, an error
	// is logged and returned. On success, a success message is logged, and a
	// protobuf response is returned indicating successful password reset initiation.
	//
	// Parameters:
	//   - ctx: The context for managing request-scoped values, cancelation signals, and deadlines.
	//   - req: A pointer to ForgotPasswordRequest containing the email address of the user to reset.
	//
	// Returns:
	//   - A pointer to ApiResponseForgotPassword containing the password reset result on success.
	//   - An error if the password reset process fails.
	ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.ApiResponseForgotPassword, error)

	// VerifyCode validates a password reset verification code.
	// It logs the operation's start and end, and utilizes the passwordResetService
	// to perform the verification. If the service indicates a failure, an error
	// is logged and returned. On success, a success message is logged, and a
	// protobuf response is returned indicating successful verification.
	//
	// Parameters:
	//   - ctx: The context for managing request-scoped values, cancelation signals, and deadlines.
	//   - req: A pointer to VerifyCodeRequest containing the verification code to be validated.
	//
	// Returns:
	//   - A pointer to ApiResponseVerifyCode containing the verification result on success.
	//   - An error if the verification process fails.
	VerifyCode(ctx context.Context, req *pb.VerifyCodeRequest) (*pb.ApiResponseVerifyCode, error)

	// ResetPassword completes the password reset process.
	// It logs the operation's start and end, and utilizes the passwordResetService
	// to perform the password reset. If the service indicates a failure, an error
	// is logged and returned. On success, a success message is logged, and a
	// protobuf response is returned indicating successful password reset completion.
	//
	// Parameters:
	//   - ctx: The context for managing request-scoped values, cancelation signals, and deadlines.
	//   - req: A pointer to ResetPasswordRequest containing the new password and reset token.
	//
	// Returns:
	//   - A pointer to ApiResponseResetPassword containing the password reset result on success.
	//   - An error if the password reset process fails.
	ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ApiResponseResetPassword, error)

	// GetMe retrieves the current authenticated user's details.
	// It logs the operation's start and end, and utilizes the identifyService
	// to perform the user retrieval. If the service indicates a failure, an error
	// is logged and returned. On success, a success message is logged, and a
	// protobuf response is returned indicating successful user retrieval and containing
	// the user's information.
	//
	// Parameters:
	//   - ctx: The context for managing request-scoped values, cancelation signals, and deadlines.
	//   - req: A pointer to GetMeRequest containing the access token.
	//
	// Returns:
	//   - A pointer to ApiResponseGetMe containing the user retrieval result on success.
	//   - An error if the user retrieval process fails.
	GetMe(ctx context.Context, req *pb.GetMeRequest) (*pb.ApiResponseGetMe, error)

	// RefreshToken generates new access tokens using a refresh token.
	// It logs the operation's start and end, and utilizes the identifyService
	// to perform the token refresh. If the service indicates a failure, an error
	// is logged and returned. On success, a success message is logged, and a
	// protobuf response is returned indicating successful token refresh and containing
	// the new tokens.
	//
	// Parameters:
	//   - ctx: The context for managing request-scoped values, cancelation signals, and deadlines.
	//   - req: A pointer to RefreshTokenRequest containing the refresh token.
	//
	// Returns:
	//   - A pointer to ApiResponseRefreshToken containing the token refresh result on success.
	//   - An error if the token refresh process fails.
	RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.ApiResponseRefreshToken, error)
}
