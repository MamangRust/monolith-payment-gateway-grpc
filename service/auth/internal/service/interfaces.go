package service

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go
type RegistrationService interface {
	Register(request *requests.RegisterRequest) (*response.UserResponse, *response.ErrorResponse)
}

type LoginService interface {
	Login(request *requests.AuthRequest) (*response.TokenResponse, *response.ErrorResponse)
}

type PasswordResetService interface {
	ForgotPassword(email string) (bool, *response.ErrorResponse)
	ResetPassword(request *requests.CreateResetPasswordRequest) (bool, *response.ErrorResponse)
	VerifyCode(code string) (bool, *response.ErrorResponse)
}

type IdentifyService interface {
	RefreshToken(token string) (*response.TokenResponse, *response.ErrorResponse)
	GetMe(token string) (*response.UserResponse, *response.ErrorResponse)
}
