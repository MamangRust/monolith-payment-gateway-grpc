package myhandlerkafka

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/service"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"go.uber.org/zap"
)

type merchantKafkaHandler struct {
	logger          logger.LoggerInterface
	merchantService service.MerchantQueryService
	kafka           *kafka.Kafka
}

func NewMerchantKafkaHandler(merchantService service.MerchantQueryService, kafka *kafka.Kafka, logger logger.LoggerInterface) sarama.ConsumerGroupHandler {
	return &merchantKafkaHandler{
		merchantService: merchantService,
		kafka:           kafka,
		logger:          logger,
	}
}

func (m *merchantKafkaHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (m *merchantKafkaHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (m *merchantKafkaHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var payload map[string]interface{}
		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			m.logger.Error("Failed to unmarshal message", zap.Error(err))
			continue
		}

		apiKey, _ := payload["api_key"].(string)
		correlationID, _ := payload["correlation_id"].(string)
		replyTopic, _ := payload["reply_topic"].(string)

		resp := map[string]interface{}{
			"correlation_id": correlationID,
			"valid":          false,
		}

		merchant, err := m.merchantService.FindByApiKey(apiKey)
		if err == nil {
			resp["valid"] = true
			resp["merchant_id"] = merchant.ID
		}

		respBytes, _ := json.Marshal(resp)

		sendErr := m.kafka.SendMessage(replyTopic, correlationID, respBytes)
		if sendErr != nil {
			m.logger.Error("Failed to send Kafka response", zap.Error(sendErr))
		}
		session.MarkMessage(msg, "")
	}
	return nil
}
