package myhandlerkafka

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/service"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
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
		var payload requests.MerchantRequestPayload
		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			m.logger.Error("Failed to unmarshal merchant request", zap.Error(err))
			continue
		}

		resp := response.MerchantResponsePayload{
			CorrelationID: payload.CorrelationID,
			Valid:         false,
		}

		merchant, err := m.merchantService.FindByApiKey(payload.ApiKey)
		if err == nil && merchant != nil {
			resp.Valid = true
			resp.MerchantID = int64(merchant.ID)
		}

		respBytes, _ := json.Marshal(resp)
		sendErr := m.kafka.SendMessage(payload.ReplyTopic, payload.CorrelationID, respBytes)
		if sendErr != nil {
			m.logger.Error("Failed to send Kafka response", zap.Error(sendErr))
		}

		session.MarkMessage(msg, "")
	}
	return nil
}
