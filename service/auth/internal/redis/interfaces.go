package mencache

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

// IdentityCache defines the interface for identity-related caching operations.
// It provides methods to manage refresh tokens and user information in cache.
//
//go:generate mockgen -source=interfaces.go -destination=mocks/cache.go
type IdentityCache interface {
	// SetRefreshToken stores a refresh token in the cache with expiration.
	//
	// Parameters:
	//   - ctx: the context for the caching operation
	//   - token: the refresh token string
	//   - expiration: the duration until the token expires
	SetRefreshToken(ctx context.Context, token string, expiration time.Duration)

	// GetRefreshToken retrieves a refresh token from the cache.
	//
	// Parameters:
	//   - ctx: the context for the caching operation
	//   - token: the refresh token string to retrieve
	//
	// Returns:
	//   - The stored token string and a boolean indicating whether it was found.
	GetRefreshToken(ctx context.Context, token string) (string, bool)

	// DeleteRefreshToken removes a refresh token from the cache.
	//
	// Parameters:
	//   - ctx: the context for the caching operation
	//   - token: the refresh token to delete
	DeleteRefreshToken(ctx context.Context, token string)

	// SetCachedUserInfo stores user information in the cache with expiration.
	//
	// Parameters:
	//   - ctx: the context for the caching operation
	//   - user: the user data to cache
	//   - expiration: the duration until the cache entry expires
	SetCachedUserInfo(ctx context.Context, user *response.UserResponse, expiration time.Duration)

	// GetCachedUserInfo retrieves cached user information by user ID.
	//
	// Parameters:
	//   - ctx: the context for the caching operation
	//   - userId: the unique identifier of the user
	//
	// Returns:
	//   - The cached UserResponse and a boolean indicating whether it was found.
	GetCachedUserInfo(ctx context.Context, userId string) (*response.UserResponse, bool)

	// DeleteCachedUserInfo removes cached user information from the cache.
	//
	// Parameters:
	//   - ctx: the context for the caching operation
	//   - userId: the unique identifier of the user to remove from cache
	DeleteCachedUserInfo(ctx context.Context, userId string)
}

// LoginCache provides caching mechanisms for login-related data
// such as cached token responses by email.
type LoginCache interface {
	// SetCachedLogin stores login token response in cache with expiration.
	//
	// Parameters:
	//   - ctx: the context for the caching operation
	//   - email: the email address used for login
	//   - data: the token response to cache
	//   - expiration: the duration until the cache entry expires
	SetCachedLogin(ctx context.Context, email string, data *response.TokenResponse, expiration time.Duration)

	// GetCachedLogin retrieves a cached token response by email.
	//
	// Parameters:
	//   - ctx: the context for the caching operation
	//   - email: the email address for which to retrieve the cached login data
	//
	// Returns:
	//   - The cached TokenResponse and a boolean indicating whether it was found.
	GetCachedLogin(ctx context.Context, email string) (*response.TokenResponse, bool)
}

// PasswordResetCache provides caching mechanisms for password reset operations,
// such as storing and retrieving reset tokens and verification codes.
type PasswordResetCache interface {
	// SetResetTokenCache stores a password reset token associated with a user ID.
	//
	// Parameters:
	//   - ctx: the context for the caching operation
	//   - token: the reset token string
	//   - userID: the ID of the user who requested the reset
	//   - expiration: the duration until the token expires
	SetResetTokenCache(ctx context.Context, token string, userID int, expiration time.Duration)

	// GetResetTokenCache retrieves a user ID associated with a given reset token.
	//
	// Parameters:
	//   - ctx: the context for the caching operation
	//   - token: the reset token to look up
	//
	// Returns:
	//   - The associated user ID and a boolean indicating if the token was found.
	GetResetTokenCache(ctx context.Context, token string) (int, bool)

	// DeleteResetTokenCache removes a reset token from the cache.
	//
	// Parameters:
	//   - ctx: the context for the caching operation
	//   - token: the reset token to delete
	DeleteResetTokenCache(ctx context.Context, token string)

	// DeleteVerificationCodeCache deletes the verification code associated with an email.
	//
	// Parameters:
	//   - ctx: the context for the caching operation
	//   - email: the email address whose verification code should be deleted
	DeleteVerificationCodeCache(ctx context.Context, email string)
}

// RegisterCache provides caching mechanisms for registration processes,
// such as storing verification codes.
type RegisterCache interface {
	// SetVerificationCodeCache stores a verification code for an email with expiration.
	//
	// Parameters:
	//   - ctx: the context for the caching operation
	//   - email: the email address associated with the code
	//   - code: the verification code to cache
	//   - expiration: the duration until the code expires
	SetVerificationCodeCache(ctx context.Context, email string, code string, expiration time.Duration)
}
