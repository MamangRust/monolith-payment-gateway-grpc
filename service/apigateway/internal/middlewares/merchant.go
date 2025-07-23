package middlewares

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	mencache "github.com/MamangRust/monolith-payment-gateway-apigateway/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// ApiKeyValidator is responsible for validating merchant API keys via a Kafka request-response pattern.
//
// It sends API key validation requests to a Kafka topic and listens for responses on a separate topic.
// The validator tracks pending responses using correlation IDs and response channels, ensuring thread-safe access
// with a mutex and supports timeouts for unresponsive requests.
type ApiKeyValidator struct {
	// kafka is the Kafka client used to publish API key validation requests and subscribe to responses.
	kafka *kafka.Kafka

	// logger is the structured logger used to record events, warnings, and errors related to API key validation.
	logger logger.LoggerInterface

	// requestTopic is the Kafka topic where validation requests are published.
	requestTopic string

	// responseTopic is the Kafka topic where validation responses are expected.
	responseTopic string

	// timeout specifies the maximum duration to wait for a validation response before returning an error.
	timeout time.Duration

	// responseChans is a map of correlation IDs to channels that receive the corresponding Kafka response payloads.
	// It is used to track in-flight validation requests and pair them with their responses.
	responseChans map[string]chan []byte

	// mu is a mutex used to protect concurrent access to responseChans.
	mu sync.Mutex

	cache mencache.MerchantCache
}

// NewApiKeyValidator returns a new ApiKeyValidator instance. It starts a Kafka consumer
// that listens to the response topic and waits for a response to the validation
// request published to the request topic. The validator is used as a middleware
// to validate the API Key in the request header and store the merchant ID in the
// Echo context.
func NewApiKeyValidator(k *kafka.Kafka, requestTopic, responseTopic string, timeout time.Duration, logger logger.LoggerInterface, cache mencache.MerchantCache) *ApiKeyValidator {
	v := &ApiKeyValidator{
		kafka:         k,
		requestTopic:  requestTopic,
		responseTopic: responseTopic,
		timeout:       timeout,
		responseChans: make(map[string]chan []byte),
		logger:        logger,
		cache:         cache,
	}

	handler := &merchantResponseHandler{validator: v}
	go func() {
		err := k.StartConsumers([]string{responseTopic}, "api-gateway-group", handler)
		if err != nil {
			v.logger.Fatal("Failed to start kafka consumer", zap.Error(err))
			panic("failed to start kafka consumer: " + err.Error())
		}
	}()

	return v
}

// Middleware returns an Echo middleware that validates the API Key in the request
// header by publishing a message to a Kafka topic and waiting for a response.
// If the validation is successful, the merchant ID is stored in the Echo context
// and the request is passed to the next handler. If the validation fails or
// times out, an HTTP error is returned.
func (v *ApiKeyValidator) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			apiKey := c.Request().Header.Get("X-Api-Key")
			if apiKey == "" {
				v.logger.Error("Missing API Key in request header")
				return echo.NewHTTPError(http.StatusUnauthorized, "API Key is required")
			}

			if merchantID, found := v.cache.GetMerchantCache(c.Request().Context(), apiKey); found {
				v.logger.Info("Merchant ID found in cache", zap.String("apiKey", apiKey), zap.String("merchant_id", merchantID))
				c.Set("merchant_id", merchantID)
				c.Set("apiKey", apiKey)
				return next(c)
			}

			correlationID := uuid.NewString()
			v.logger.Info("Cache miss, sending Kafka request", zap.String("apiKey", apiKey), zap.String("correlationID", correlationID))

			payload := requests.MerchantRequestPayload{
				ApiKey:        apiKey,
				CorrelationID: correlationID,
				ReplyTopic:    v.responseTopic,
			}

			data, err := json.Marshal(payload)
			if err != nil {
				v.logger.Error("Failed to encode payload", zap.Error(err), zap.String("correlationID", correlationID))
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to encode payload")
			}

			err = v.kafka.SendMessage(v.requestTopic, correlationID, data)
			if err != nil {
				v.logger.Error("Failed to send Kafka message", zap.Error(err), zap.String("correlationID", correlationID))
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to send Kafka message")
			}

			respChan := make(chan []byte, 1)

			v.mu.Lock()
			v.responseChans[correlationID] = respChan
			v.mu.Unlock()

			defer func() {
				v.mu.Lock()
				delete(v.responseChans, correlationID)
				v.mu.Unlock()
			}()

			select {
			case msg := <-respChan:
				var response response.MerchantResponsePayload
				if err := json.Unmarshal(msg, &response); err != nil {
					v.logger.Error("Failed to decode Kafka response", zap.Error(err), zap.String("correlationID", correlationID))
					return echo.NewHTTPError(http.StatusInternalServerError, "Invalid response format")
				}

				if !response.Valid || response.MerchantID == 0 {
					v.logger.Error("Invalid API Key validation result", zap.String("correlationID", correlationID), zap.String("apiKey", apiKey))
					return echo.NewHTTPError(http.StatusUnauthorized, "Invalid API Key")
				}

				merchantIDStr := strconv.Itoa(int(response.MerchantID))

				v.cache.SetMerchantCache(c.Request().Context(), merchantIDStr, apiKey)

				v.logger.Info("API Key validated successfully", zap.String("merchant_id", merchantIDStr), zap.String("correlationID", correlationID))

				c.Set("merchant_id", merchantIDStr)
				c.Set("apiKey", apiKey)
				return next(c)

			case <-time.After(v.timeout):
				v.logger.Error("Timeout waiting for Kafka response", zap.String("correlationID", correlationID), zap.Duration("timeout", v.timeout))
				return echo.NewHTTPError(http.StatusRequestTimeout, "Timeout waiting for API key validation")
			}
		}
	}
}
