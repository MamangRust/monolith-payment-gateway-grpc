package middlewares

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ApiKeyValidator struct {
	kafka         *kafka.Kafka
	requestTopic  string
	responseTopic string
	timeout       time.Duration
	responseChans map[string]chan []byte
	mu            sync.Mutex
}

func NewApiKeyValidator(k *kafka.Kafka, requestTopic, responseTopic string, timeout time.Duration) *ApiKeyValidator {
	v := &ApiKeyValidator{
		kafka:         k,
		requestTopic:  requestTopic,
		responseTopic: responseTopic,
		timeout:       timeout,
		responseChans: make(map[string]chan []byte),
	}

	handler := &responseHandler{validator: v}
	go func() {
		err := k.StartConsumers([]string{responseTopic}, "api-gateway-group", handler)
		if err != nil {
			panic("failed to start kafka consumer: " + err.Error())
		}
	}()

	return v
}

func (v *ApiKeyValidator) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			apiKey := c.Request().Header.Get("X-Api-Key")
			if apiKey == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "API Key is required")
			}

			correlationID := uuid.NewString()
			payload := map[string]interface{}{
				"api_key":        apiKey,
				"correlation_id": correlationID,
				"reply_topic":    v.responseTopic,
			}

			data, err := json.Marshal(payload)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to encode payload")
			}

			err = v.kafka.SendMessage(v.requestTopic, correlationID, data)
			if err != nil {
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
				var response map[string]interface{}
				if err := json.Unmarshal(msg, &response); err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, "Invalid response format")
				}

				valid, ok := response["valid"].(bool)
				if !ok || !valid {
					return echo.NewHTTPError(http.StatusUnauthorized, "Invalid API Key")
				}

				merchantID, ok := response["merchant_id"]
				if !ok {
					return echo.NewHTTPError(http.StatusUnauthorized, "Merchant ID not found in response")
				}

				c.Set("merchant_id", merchantID)
				c.Set("apiKey", apiKey)
				return next(c)

			case <-time.After(v.timeout):
				return echo.NewHTTPError(http.StatusRequestTimeout, "Timeout waiting for API key validation")
			}
		}
	}
}
