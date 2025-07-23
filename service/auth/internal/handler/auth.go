package handler

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-auth/internal/service"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"go.uber.org/zap"

	pb "github.com/MamangRust/monolith-payment-gateway-pb/auth"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/auth"
)

// authHandleGrpc represents the authentication service handler for the gRPC API.
type authHandleGrpc struct {
	pb.UnimplementedAuthServiceServer
	registerService      service.RegistrationService
	loginService         service.LoginService
	passwordResetService service.PasswordResetService
	identifyService      service.IdentifyService
	logger               logger.LoggerInterface
	mapper               protomapper.AuthProtoMapper
}

// NewAuthHandleGrpc creates a new instance of the authentication service handler
// for the gRPC API. It takes a pointer to the authentication service and a logger
// as arguments and returns a pointer to the handler struct.
//
// The handler wraps the authentication service and provides methods for the
// various authentication-related operations.
//
// The logger is used to log errors and other important events.
func NewAuthHandleGrpc(authService *service.Service, logger logger.LoggerInterface) pb.AuthServiceServer {
	return &authHandleGrpc{
		registerService:      authService.Register,
		loginService:         authService.Login,
		passwordResetService: authService.PasswordReset,
		identifyService:      authService.Identify,
		logger:               logger,
		mapper:               protomapper.NewAuthProtoMapper(),
	}
}

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
func (s *authHandleGrpc) VerifyCode(ctx context.Context, req *pb.VerifyCodeRequest) (*pb.ApiResponseVerifyCode, error) {
	s.logger.Info("VerifyCode called", zap.String("code", req.Code))

	_, err := s.passwordResetService.VerifyCode(ctx, req.Code)
	if err != nil {
		s.logger.Error("VerifyCode failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseVerifyCode("success", "Verification successful")

	s.logger.Info("VerifyCode success", zap.String("code", req.Code))

	return so, nil
}

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
func (s *authHandleGrpc) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.ApiResponseForgotPassword, error) {
	s.logger.Info("ForgotPassword called", zap.String("email", req.Email))

	_, err := s.passwordResetService.ForgotPassword(ctx, req.Email)
	if err != nil {
		s.logger.Error("ForgotPassword failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseForgotPassword("success", "Forgot password successful")

	s.logger.Info("ForgotPassword successful", zap.Bool("success", true))

	return so, nil
}

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
func (s *authHandleGrpc) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ApiResponseResetPassword, error) {
	s.logger.Info("ResetPassword called", zap.String("reset_token", req.ResetToken))

	_, err := s.passwordResetService.ResetPassword(ctx, &requests.CreateResetPasswordRequest{
		ResetToken:      req.ResetToken,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	})
	if err != nil {
		s.logger.Error("ResetPassword failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseResetPassword("success", "Reset password successful")

	s.logger.Info("ResetPassword successful", zap.Bool("success", true))

	return so, nil
}

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
func (s *authHandleGrpc) LoginUser(ctx context.Context, req *pb.LoginRequest) (*pb.ApiResponseLogin, error) {
	s.logger.Info("LoginUser called", zap.String("email", req.Email))

	request := &requests.AuthRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	res, err := s.loginService.Login(ctx, request)
	if err != nil {
		s.logger.Error("LoginUser failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseLogin("success", "Login successful", res)

	s.logger.Info("LoginUser successful", zap.Bool("success", true))

	return so, nil
}

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
func (s *authHandleGrpc) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.ApiResponseRefreshToken, error) {
	s.logger.Info("RefreshToken called")

	res, err := s.identifyService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		s.logger.Error("RefreshToken failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseRefreshToken("success", "Refresh token successful", res)

	s.logger.Info("RefreshToken successful", zap.Bool("success", true))

	return so, nil
}

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
func (s *authHandleGrpc) GetMe(ctx context.Context, req *pb.GetMeRequest) (*pb.ApiResponseGetMe, error) {
	s.logger.Info("GetMe called")

	res, err := s.identifyService.GetMe(ctx, req.AccessToken)
	if err != nil {
		s.logger.Error("GetMe failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseGetMe("success", "Get me successful", res)

	s.logger.Info("GetMe successful", zap.Bool("success", true))

	return so, nil
}

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
func (s *authHandleGrpc) RegisterUser(ctx context.Context, req *pb.RegisterRequest) (*pb.ApiResponseRegister, error) {
	s.logger.Info("RegisterUser called", zap.String("email", req.Email))

	request := &requests.RegisterRequest{
		FirstName:       req.Firstname,
		LastName:        req.Lastname,
		Email:           req.Email,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	}

	res, err := s.registerService.Register(ctx, request)
	if err != nil {
		s.logger.Error("RegisterUser failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseRegister("success", "Registration successful", res)

	s.logger.Info("RegisterUser successful", zap.Bool("success", true))
	return so, nil
}
