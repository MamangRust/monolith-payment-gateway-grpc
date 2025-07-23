package middlewares

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

// whiteListPaths defines a list of HTTP paths that are excluded from authentication or middleware checks.
//
// These paths are typically public endpoints such as login, registration, or documentation routes.
// Middleware such as JWT authentication should skip these paths to allow anonymous access.
var whiteListPaths = []string{
	"/api/auth/login",       // Public login endpoint
	"/api/auth/register",    // Public registration endpoint
	"/api/auth/hello",       // Example open test endpoint
	"/api/auth/verify-code", // Endpoint for verifying user email/code
	"/docs/",                // Documentation path with trailing slash
	"/docs",                 // Documentation path without trailing slash
	"/swagger",              // Swagger UI endpoint
}

// WebSecurityConfig adds JWT middleware to an echo router.
//
// The middleware uses the SigningKey from the config file. It also sets the
// Skipper to skipAuth, which allows the following paths to be accessed without
// a valid JWT:
//
// - /api/auth/login
// - /api/auth/register
// - /api/auth/hello
// - /api/auth/verify-code
// - /docs/
// - /docs
// - /swagger
//
// The SuccessHandler is used to add the subject of the JWT to the context
// under the key "user_id".
//
// The ErrorHandler is used to return a 401 Unauthorized status code in case
// of a JWT error.
func WebSecurityConfig(e *echo.Echo) {
	config := echojwt.Config{
		SigningKey: []byte(viper.GetString("SECRET_KEY")),
		Skipper:    skipAuth,
		SuccessHandler: func(c echo.Context) {
			user := c.Get("user").(*jwt.Token)

			if claims, ok := user.Claims.(jwt.MapClaims); ok {
				subject := claims["sub"]
				c.Set("user_id", subject)
			}
		},
		ErrorHandler: func(c echo.Context, err error) error {
			fmt.Println("JWT Error:", err)

			return echo.ErrUnauthorized
		},
	}
	e.Use(echojwt.WithConfig(config))

}

// skipAuth is the Skipper used in the JWT middleware.
//
// It returns true for the following paths, which are skipped by the JWT middleware:
//
// - /api/auth/login
// - /api/auth/register
// - /api/auth/hello
// - /api/auth/verify-code
// - /docs/
// - /docs
// - /swagger
func skipAuth(e echo.Context) bool {
	path := e.Path()

	for _, p := range whiteListPaths {
		if path == p || strings.HasPrefix(path, "/swagger") {
			return true
		}
	}

	return false
}
