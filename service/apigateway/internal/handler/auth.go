package handler

import (
	"net/http"
	"strings"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/auth_errors"
	apimapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/api"
	"github.com/MamangRust/monolith-payment-gateway-shared/pb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authHandleApi struct {
	client  pb.AuthServiceClient
	logger  logger.LoggerInterface
	mapping apimapper.AuthResponseMapper
}

func NewHandlerAuth(client pb.AuthServiceClient, router *echo.Echo, logger logger.LoggerInterface, mapper apimapper.AuthResponseMapper) *authHandleApi {
	authHandler := &authHandleApi{
		client:  client,
		logger:  logger,
		mapping: mapper,
	}
	routerAuth := router.Group("/api/auth")

	routerAuth.GET("/verify-code", authHandler.VerifyCode)
	routerAuth.POST("/forgot-password", authHandler.ForgotPassword)
	routerAuth.POST("/reset-password", authHandler.ResetPassword)
	routerAuth.GET("/hello", authHandler.HandleHello)
	routerAuth.POST("/register", authHandler.Register)
	routerAuth.POST("/login", authHandler.Login)
	routerAuth.POST("/refresh-token", authHandler.RefreshToken)
	routerAuth.GET("/me", authHandler.GetMe)

	return authHandler
}

// HandleHello godoc
// @Summary Returns a "Hello" message
// @Tags Auth
// @Description Returns a simple "Hello" message for testing purposes.
// @Produce json
// @Success 200 {string} string "Hello"
// @Router /auth/hello [get]
func (h *authHandleApi) HandleHello(c echo.Context) error {
	return c.String(200, "Hello")
}

// VerifyCode godoc
// @Summary Verifies the user using a verification code
// @Tags Auth
// @Description Verifies the user's email using the verification code provided in the query parameter.
// @Produce json
// @Param verify_code query string true "Verification Code"
// @Success 200 {object} response.ApiResponseVerifyCode
// @Failure 400 {object} response.ErrorResponse
// @Router /auth/verify-code [get]
func (h *authHandleApi) VerifyCode(c echo.Context) error {
	verifyCode := c.QueryParam("verify_code")

	ctx := c.Request().Context()

	res, err := h.client.VerifyCode(ctx, &pb.VerifyCodeRequest{
		Code: verifyCode,
	})

	if err != nil {
		h.logger.Error("Failed to verify code", zap.Error(err))
		return auth_errors.ErrApiVerifyCode(c)
	}
	so := h.mapping.ToResponseVerifyCode(res)

	return c.JSON(http.StatusOK, so)
}

// ForgotPassword godoc
// @Summary Sends a reset token to the user's email
// @Tags Auth
// @Description Initiates password reset by sending a reset token to the provided email.
// @Accept json
// @Produce json
// @Param request body requests.ForgotPasswordRequest true "Forgot Password Request"
// @Success 200 {object} response.ApiResponseForgotPassword
// @Failure 400 {object} response.ErrorResponse
// @Router /auth/forgot-password [post]
func (h *authHandleApi) ForgotPassword(c echo.Context) error {
	var body requests.ForgotPasswordRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Invalid request format", zap.Error(err))
		return auth_errors.ErrBindForgotPassword(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation failed", zap.Error(err))
		return auth_errors.ErrValidateForgotPassword(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.ForgotPassword(ctx, &pb.ForgotPasswordRequest{
		Email: body.Email,
	})

	if err != nil {
		h.logger.Error("Failed to forgot password", zap.Error(err))
		return auth_errors.ErrApiForgotPassword(c)
	}

	so := h.mapping.ToResponseForgotPassword(res)

	return c.JSON(http.StatusOK, so)
}

// ResetPassword godoc
// @Summary Resets the user's password using a reset token
// @Tags Auth
// @Description Allows user to reset their password using a valid reset token.
// @Accept json
// @Produce json
// @Param request body requests.CreateResetPasswordRequest true "Reset Password Request"
// @Success 200 {object} response.ApiResponseResetPassword
// @Failure 400 {object} response.ErrorResponse
// @Router /auth/reset-password [post]
func (h *authHandleApi) ResetPassword(c echo.Context) error {
	var body requests.CreateResetPasswordRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Invalid request format", zap.Error(err))
		return auth_errors.ErrBindResetPassword(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation failed", zap.Error(err))
		return auth_errors.ErrValidateResetPassword(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.ResetPassword(ctx, &pb.ResetPasswordRequest{
		ResetToken:      body.ResetToken,
		Password:        body.Password,
		ConfirmPassword: body.ConfirmPassword,
	})

	if err != nil {
		h.logger.Error("Failed to reset password", zap.Error(err))
		return auth_errors.ErrApiResetPassword(c)
	}

	so := h.mapping.ToResponseResetPassword(res)

	return c.JSON(http.StatusOK, so)
}

// Register godoc
// @Summary Register a new user
// @Tags Auth
// @Description Registers a new user with the provided details.
// @Accept json
// @Produce json
// @Param request body requests.CreateUserRequest true "User registration data"
// @Success 200 {object} response.ApiResponseRegister "Success"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/auth/register [post]
func (h *authHandleApi) Register(c echo.Context) error {
	var body requests.CreateUserRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Invalid request format", zap.Error(err))
		return auth_errors.ErrBindRegister(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation failed", zap.Error(err))
		return auth_errors.ErrValidateRegister(c)
	}

	data := &pb.RegisterRequest{
		Firstname:       body.FirstName,
		Lastname:        body.LastName,
		Email:           body.Email,
		Password:        body.Password,
		ConfirmPassword: body.ConfirmPassword,
	}

	ctx := c.Request().Context()

	res, err := h.client.RegisterUser(ctx, data)

	if err != nil {
		h.logger.Error("Registration failed", zap.Error(err))
		return auth_errors.ErrApiRegister(c)
	}

	so := h.mapping.ToResponseRegister(res)

	return c.JSON(http.StatusOK, so)
}

// Login godoc
// @Summary Authenticate a user
// @Tags Auth
// @Description Authenticates a user using the provided email and password.
// @Accept json
// @Produce json
// @Param request body requests.AuthRequest true "User login credentials"
// @Success 200 {object} response.ApiResponseLogin "Success"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/auth/login [post]
func (h *authHandleApi) Login(c echo.Context) error {
	var body requests.AuthRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Invalid request format", zap.Error(err))
		return auth_errors.ErrBindLogin(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation failed", zap.Error(err))
		return auth_errors.ErrValidateRegister(c)
	}

	data := &pb.LoginRequest{
		Email:    body.Email,
		Password: body.Password,
	}

	ctx := c.Request().Context()

	res, err := h.client.LoginUser(ctx, data)

	if err != nil {
		if status.Code(err) == codes.Unauthenticated {
			h.logger.Debug("Invalid login attempt", zap.String("email", body.Email))
			return auth_errors.ErrInvalidLogin(c)
		}

		h.logger.Error("Login failed", zap.Error(err))

		if status.Code(err) == codes.Internal && strings.Contains(err.Error(), "empty token") {
			return auth_errors.ErrInvalidAccessToken(c)
		}
	}

	mappedResponse := h.mapping.ToResponseLogin(res)

	if mappedResponse.Data == nil || mappedResponse.Data.AccessToken == "" {
		h.logger.Error("Empty token in final response", zap.Any("response", mappedResponse))
		return auth_errors.ErrApiLogin(c)
	}

	return c.JSON(http.StatusOK, mappedResponse)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Tags Auth
// @Security Bearer
// @Description Refreshes the access token using a valid refresh token.
// @Accept json
// @Produce json
// @Param request body requests.RefreshTokenRequest true "Refresh token data"
// @Success 200 {object} response.ApiResponseRefreshToken "Success"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/auth/refresh-token [post]
func (h *authHandleApi) RefreshToken(c echo.Context) error {
	var body requests.RefreshTokenRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Invalid request format", zap.Error(err))
		return auth_errors.ErrBindRefreshToken(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation failed", zap.Error(err))
		return auth_errors.ErrValidateRefreshToken(c)
	}

	res, err := h.client.RefreshToken(c.Request().Context(), &pb.RefreshTokenRequest{
		RefreshToken: body.RefreshToken,
	})

	if err != nil {
		h.logger.Error("Token refresh failed", zap.Error(err))
		return auth_errors.ErrApiRefreshToken(c)
	}

	so := h.mapping.ToResponseRefreshToken(res)

	return c.JSON(http.StatusOK, so)
}

// GetMe godoc
// @Summary Get current user information
// @Tags Auth
// @Security Bearer
// @Description Retrieves the current user's information using a valid access token from the Authorization header.
// @Produce json
// @Security BearerToken
// @Success 200 {object} response.ApiResponseGetMe "Success"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/auth/me [get]
func (h *authHandleApi) GetMe(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")

	h.logger.Debug("Authorization header: ", zap.String("authHeader", authHeader))

	if !strings.HasPrefix(authHeader, "Bearer ") {
		h.logger.Debug("Authorization header is missing or invalid format")
		return auth_errors.ErrInvalidAccessToken(c)
	}

	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

	res, err := h.client.GetMe(c.Request().Context(), &pb.GetMeRequest{
		AccessToken: accessToken,
	})

	if err != nil {
		h.logger.Error("Failed to get user information", zap.Error(err))
		return auth_errors.ErrApiGetMe(c)
	}

	so := h.mapping.ToResponseGetMe(res)

	return c.JSON(http.StatusOK, so)
}
