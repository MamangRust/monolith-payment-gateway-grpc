package middlewares

import (
	"encoding/json"

	"github.com/IBM/sarama"
)

type responseHandler struct {
	validator *ApiKeyValidator
}

func NewResponseHandler(validator *ApiKeyValidator) *responseHandler {
	return &responseHandler{
		validator: validator,
	}
}

func (h *responseHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *responseHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *responseHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var payload map[string]interface{}
		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			continue
		}

		correlationID, ok := payload["correlation_id"].(string)
		if !ok {
			continue
		}

		h.validator.mu.Lock()
		ch, ok := h.validator.responseChans[correlationID]
		h.validator.mu.Unlock()

		if ok {
			ch <- msg.Value
		}

		session.MarkMessage(msg, "")
	}
	return nil
}
