package kafka

import (
	"github.com/IBM/sarama"
	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/service"
)

type merchantDocumentKafkaHandler struct {
	merchantDocumentService service.MerchantDocumentCommandService
}

func NewMerchantDocumentKafkaHandler(merchantDocumentService service.MerchantDocumentCommandService) sarama.ConsumerGroupHandler {
	return &merchantDocumentKafkaHandler{
		merchantDocumentService: merchantDocumentService,
	}
}

func (s *merchantDocumentKafkaHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (s *merchantDocumentKafkaHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (s *merchantDocumentKafkaHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	return nil
}
