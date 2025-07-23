package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/service"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"go.uber.org/zap"
)

// saldoKafkaHandler is a struct that implements the sarama.ConsumerGroupHandler interface
type saldoKafkaHandler struct {
	logger       logger.LoggerInterface
	saldoService service.SaldoCommandService
}

// NewSaldoKafkaHandler creates a new Kafka consumer group handler for processing saldo-related Kafka messages.
//
// It takes a saldo command service and a logger as parameters.
// The handler is responsible for consuming messages from Kafka topics related to saldo operations.
// It implements the sarama.ConsumerGroupHandler interface to manage consumer group lifecycle events.
func NewSaldoKafkaHandler(saldoService service.SaldoCommandService, logger logger.LoggerInterface) sarama.ConsumerGroupHandler {
	return &saldoKafkaHandler{
		saldoService: saldoService,
		logger:       logger,
	}
}

// Setup is called when the consumer group is first initialized.
func (s *saldoKafkaHandler) Setup(session sarama.ConsumerGroupSession) error {
	s.logger.Info("saldo kafka handler setup")
	return nil
}

// Cleanup is called when the consumer group is closed.
func (s *saldoKafkaHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	s.logger.Info("saldo kafka handler cleanup")
	return nil
}

// ConsumeClaim processes incoming Kafka messages from the specified consumer group claim.
// It unmarshals each message into a payload map and creates a new saldo request.
// If a valid saldo request is found, it calls the CreateSaldo method of the saldoService
// with the created request. If the CreateSaldo method returns an error, it logs the
// error and returns an error with the message "card service error: <error message>".
// Each message is marked as processed in the consumer group session.
func (s *saldoKafkaHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.logger.Info("saldo kafka handler consume claim")

	for msg := range claim.Messages() {
		var payload map[string]interface{}

		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			return err
		}

		_, errRes := s.saldoService.CreateSaldo(ctx, &requests.CreateSaldoRequest{
			CardNumber:   payload["card_number"].(string),
			TotalBalance: int(payload["total_balance"].(float64)),
		})

		if errRes != nil {
			s.logger.Error("card service error", zap.Any("error", errRes))

			return fmt.Errorf("card service error: %v", errRes.Message)
		}
	}

	s.logger.Info("saldo kafka handler consume claim success", zap.Bool("success", true))

	return nil
}
