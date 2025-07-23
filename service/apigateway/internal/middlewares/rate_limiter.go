package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

// RateLimiter wraps a token bucket rate limiter to control the rate of operations,
// such as API requests or background jobs.
//
// It helps to prevent abuse and ensure fair usage of system resources by limiting
// how frequently certain actions can be performed.
type RateLimiter struct {
	// limiter is the underlying token bucket rate limiter.
	// It manages the rate at which tokens are added and how many tokens are allowed per burst.
	limiter *rate.Limiter
}

// NewRateLimiter creates a new RateLimiter with the given rate and burst parameters.
// The rate is specified in requests per second, and the burst parameter
// specifies the maximum number of requests allowed in a single burst of
// traffic.
func NewRateLimiter(rps int, burst int) *RateLimiter {
	limiter := rate.NewLimiter(rate.Limit(rps), burst)
	return &RateLimiter{
		limiter: limiter,
	}
}

// Limit limits the number of requests to the given handler to the specified
// rate. If the rate is exceeded, it returns a 429 (Too Many Requests) error.
// The rate is specified in requests per second, and the burst parameter
// specifies how many requests can be made before the rate limiting kicks in.
func (rl *RateLimiter) Limit(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !rl.limiter.Allow() {
			return c.JSON(http.StatusTooManyRequests, map[string]string{
				"error": "Too many requests, please try again later",
			})
		}
		return next(c)
	}
}
