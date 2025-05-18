package middlewares

import (
	"net/http"

	"github.com/MamangRust/payment-gateway-monolith-grpc/shared/pb"
	"github.com/labstack/echo/v4"
)

// Todo: ini bakal pakai nats aja untuk ambil api-keynya
func ApiKeyMiddleware(merchantService pb.MerchantServiceClient) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			apiKey := c.Request().Header.Get("X-Api-Key")
			if apiKey == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "API Key is required")
			}

			_, err := merchantService.FindByApiKey(c.Request().Context(), &pb.FindByApiKeyRequest{
				ApiKey: apiKey,
			})
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid API Key")
			}

			c.Set("apiKey", apiKey)

			return next(c)
		}
	}
}
