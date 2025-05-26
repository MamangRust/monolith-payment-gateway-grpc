package middlewares

import (
	"encoding/json"
	"net/http"
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
		err := k.StartConsumers([]string{responseTopic}, "merchant-transaction", handler)
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

			data, _ := json.Marshal(payload)
			err := v.kafka.SendMessage(v.requestTopic, correlationID, data)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Kafka send error")
			}

			respChan := make(chan []byte, 1)
			v.responseChans[correlationID] = respChan
			defer delete(v.responseChans, correlationID)

			select {
			case msg := <-respChan:
				var response map[string]interface{}
				_ = json.Unmarshal(msg, &response)
				if response["valid"].(bool) {
					c.Set("merchant_id", response["merchant_id"])
					return next(c)
				}
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid API Key")
			case <-time.After(v.timeout):
				return echo.NewHTTPError(http.StatusRequestTimeout, "Timeout waiting for Kafka response")
			}
		}
	}
}
