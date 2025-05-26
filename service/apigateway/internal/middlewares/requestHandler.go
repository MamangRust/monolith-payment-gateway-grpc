package middlewares

import (
	"encoding/json"

	"github.com/IBM/sarama"
)

type responseHandler struct {
	validator *ApiKeyValidator
}

func NewResponseHandler() *responseHandler {
	return &responseHandler{
		validator: &ApiKeyValidator{},
	}
}

func (h *responseHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *responseHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *responseHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var payload map[string]interface{}
		_ = json.Unmarshal(msg.Value, &payload)

		correlationID, _ := payload["correlation_id"].(string)
		if ch, ok := h.validator.responseChans[correlationID]; ok {
			ch <- msg.Value
		}
		session.MarkMessage(msg, "")
	}
	return nil
}
