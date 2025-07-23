package service

import (
	"github.com/MamangRust/monolith-payment-gateway-auth/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-auth/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-auth/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-pkg/auth"
	"github.com/MamangRust/monolith-payment-gateway-pkg/hash"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/user"
)

// Service provides core user authentication and identity operations.
type Service struct {
	// Login handles user authentication and session issuance.
	Login LoginService

	// Register handles user account creation and role assignment.
	Register RegistrationService

	// PasswordReset handles the password recovery flow, including token verification.
	PasswordReset PasswordResetService

	// Identify resolves user identity from tokens and manages token sessions.
	Identify IdentifyService
}

// Deps holds all external dependencies required to initialize the Service.
type Deps struct {

	// ErrorHandler is the centralized error handler for all sub-services.
	ErrorHandler *errorhandler.ErrorHandler

	// Mencache provides in-memory caching for authentication-related data.
	Mencache *mencache.Mencache

	// Repositories holds all required repositories (User, Role, Token, etc.).
	Repositories *repository.Repositories

	// Token manages JWT access and refresh tokens.
	Token auth.TokenManager

	// Hash is the password hashing and verification utility.
	Hash hash.HashPassword

	// Logger is the logging interface used across services.
	Logger logger.LoggerInterface

	// Kafka is the event publisher for sending user-related events.
	Kafka *kafka.Kafka
}

// NewService initializes and returns the core authentication service bundle.
func NewService(deps *Deps) *Service {
	tokenService := NewTokenService(&tokenServiceDeps{
		Token:        deps.Token,
		RefreshToken: deps.Repositories.RefreshToken,
		Logger:       deps.Logger,
	})

	mapper := responseservice.NewUserQueryResponseMapper()

	return &Service{
		Login:         newLogin(deps, tokenService),
		Register:      newRegister(deps, mapper),
		PasswordReset: newPasswordReset(deps),
		Identify:      newIdentity(deps, tokenService),
	}
}

// newLogin initializes and returns the LoginService.
func newLogin(deps *Deps, tokenService *tokenService) LoginService {
	return NewLoginService(&LoginServiceDeps{
		ErrorPassword:  deps.ErrorHandler.PasswordError,
		ErrorToken:     deps.ErrorHandler.TokenError,
		ErrorHandler:   deps.ErrorHandler.LoginError,
		Cache:          deps.Mencache.LoginCache,
		Logger:         deps.Logger,
		Hash:           deps.Hash,
		UserRepository: deps.Repositories.User,
		RefreshToken:   deps.Repositories.RefreshToken,
		Token:          deps.Token,
		TokenService:   tokenService,
	})
}

// newRegister initializes and returns the RegistrationService.
func newRegister(deps *Deps, mapper responseservice.UserQueryResponseMapper) RegistrationService {
	return NewRegisterService(&RegisterServiceDeps{
		ErrorHandler:  deps.ErrorHandler.RegisterError,
		ErrorPassword: deps.ErrorHandler.PasswordError,
		ErrorMarshal:  deps.ErrorHandler.MarshalError,
		ErrorKafka:    deps.ErrorHandler.KafkaError,
		Cache:         deps.Mencache.RegisterCache,
		User:          deps.Repositories.User,
		Role:          deps.Repositories.Role,
		UserRole:      deps.Repositories.UserRole,
		Hash:          deps.Hash,
		Kafka:         deps.Kafka,
		Logger:        deps.Logger,
		Mapper:        mapper,
	})
}

// newPasswordReset initializes the reset forgot password, reset password and verify code services.
func newPasswordReset(deps *Deps) PasswordResetService {
	return NewPasswordResetService(&PasswordResetServiceDeps{
		ErrorHandler:      deps.ErrorHandler.PasswordResetError,
		ErrorRandomString: deps.ErrorHandler.RandomString,
		ErrorMarshal:      deps.ErrorHandler.MarshalError,
		ErrorPassword:     deps.ErrorHandler.PasswordError,
		ErrorKafka:        deps.ErrorHandler.KafkaError,
		Cache:             deps.Mencache.PasswordResetCache,
		Kafka:             deps.Kafka,
		Logger:            deps.Logger,
		User:              deps.Repositories.User,
		ResetToken:        deps.Repositories.ResetToken,
	})
}

// newIdentity initializes the IdentifyService for identity verification and token refresh.
func newIdentity(deps *Deps, tokenService *tokenService) IdentifyService {
	return NewIdentityService(&IdentityServiceDeps{
		ErrorHandler: deps.ErrorHandler.IdentityError,
		ErrorToken:   deps.ErrorHandler.TokenError,
		Cache:        deps.Mencache.IdentityCache,
		Token:        deps.Token,
		RefreshToken: deps.Repositories.RefreshToken,
		User:         deps.Repositories.User,
		Logger:       deps.Logger,
		TokenService: tokenService,
	})
}
