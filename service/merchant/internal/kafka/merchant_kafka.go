package myhandlerkafka

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/service"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"go.uber.org/zap"
)

type merchantKafkaHandler struct {
	logger          logger.LoggerInterface
	merchantService service.MerchantQueryService
}

func NewMerchantKafkaHandler(merchantService service.MerchantQueryService, logger logger.LoggerInterface) sarama.ConsumerGroupHandler {
	return &merchantKafkaHandler{
		merchantService: merchantService,
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
			m.logger.Debug("Failed to unmarshal message: %v", zap.Error(err))
			continue
		}

		apiKeyRaw, ok := payload["api_key"]
		if !ok {
			m.logger.Debug("invalid api_key format")
			continue
		}

		apiKey, ok := apiKeyRaw.(string)
		if !ok {
			m.logger.Debug("invalid api_key format")
			continue
		}

		_, err := m.merchantService.FindByApiKey(apiKey)
		if err != nil {
			m.logger.Debug("invalid api_key")
			continue
		}

		m.logger.Debug("valid api_key, processing merchant payload...")

		m.logger.Debug("Received Merchant Data: %v", zap.Any("payload", payload))

		session.MarkMessage(msg, "")
	}
	return nil
}
