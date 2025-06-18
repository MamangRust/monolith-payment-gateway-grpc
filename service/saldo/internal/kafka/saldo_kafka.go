package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/service"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"go.uber.org/zap"
)

type saldoKafkaHandler struct {
	logger       logger.LoggerInterface
	saldoService service.SaldoCommandService
}

func NewSaldoKafkaHandler(saldoService service.SaldoCommandService, logger logger.LoggerInterface) sarama.ConsumerGroupHandler {
	return &saldoKafkaHandler{
		saldoService: saldoService,
		logger:       logger,
	}
}

func (s *saldoKafkaHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (s *saldoKafkaHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (s *saldoKafkaHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var payload map[string]interface{}

		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			return err
		}

		_, errRes := s.saldoService.CreateSaldo(&requests.CreateSaldoRequest{
			CardNumber:   payload["card_number"].(string),
			TotalBalance: int(payload["total_balance"].(float64)),
		})

		if errRes != nil {
			s.logger.Error("card service error", zap.Any("error", errRes))

			return fmt.Errorf("card service error: %v", errRes.Message)
		}
	}

	return nil
}
