package kafka

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/MamangRust/monolith-payment-gateway-card/internal/service"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"go.uber.org/zap"
)

type cardKafkaHandler struct {
	logger      logger.LoggerInterface
	cardService service.CardCommandService
}

func NewCardKafkaHandler(cardService service.CardCommandService, logger logger.LoggerInterface) *cardKafkaHandler {
	return &cardKafkaHandler{
		cardService: cardService,
		logger:      logger,
	}
}

func (s *cardKafkaHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (s *cardKafkaHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (s *cardKafkaHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var payload map[string]interface{}
		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			continue
		}

		card := &requests.CreateCardRequest{
			UserID:       payload["user_id"].(int),
			CardType:     payload["card_type"].(string),
			ExpireDate:   payload["expire_date"].(time.Time),
			CVV:          payload["cvv"].(string),
			CardProvider: payload["card_provider"].(string),
		}

		_, errRes := s.cardService.CreateCard(card)

		if errRes != nil {
			s.logger.Error("card service error", zap.Any("error", errRes))

			return fmt.Errorf("card service error: %v", errRes.Message)
		}
	}

	return nil
}
