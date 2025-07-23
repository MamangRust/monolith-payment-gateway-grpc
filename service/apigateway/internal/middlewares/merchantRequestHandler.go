package middlewares

import (
	"encoding/json"

	"github.com/IBM/sarama"
)

// merchantResponseHandler handles Kafka consumer group lifecycle events for merchant-related API key validation responses.
//
// It implements the sarama.ConsumerGroupHandler interface and is primarily responsible
// for handling the setup and cleanup phases of the Kafka consumer group session.
// The actual message processing (ConsumeClaim) should be implemented to process incoming responses.
type merchantResponseHandler struct {
	// validator is the reference to ApiKeyValidator that coordinates request-response handling
	// for validating merchant API keys via Kafka.
	validator *ApiKeyValidator
}

// NewResponseHandler creates a new instance of merchantResponseHandler with the provided ApiKeyValidator.
// This handler is responsible for processing messages from a Kafka consumer group, using the validator to
// manage response channels based on correlation IDs.
func NewResponseHandler(validator *ApiKeyValidator) *merchantResponseHandler {
	return &merchantResponseHandler{
		validator: validator,
	}
}

// Setup is called when a new Kafka consumer group session begins.
//
// This method can be used to perform any necessary initialization before message consumption starts.
// In this implementation, no setup is required.
func (h *merchantResponseHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is called when a Kafka consumer group session ends.
//
// It can be used to perform cleanup tasks, such as closing resources or finalizing state.
// This implementation performs no cleanup.
func (h *merchantResponseHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim processes incoming Kafka messages from the specified consumer group claim.
// It unmarshals each message into a payload map and retrieves the correlation ID.
// If a valid correlation ID is found, it sends the message value to the corresponding
// response channel managed by the validator. Each message is marked as processed
// in the consumer group session.
func (h *merchantResponseHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
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
