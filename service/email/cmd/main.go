package main

import (
	"log"
	"net/http"

	"github.com/MamangRust/monolith-payment-gateway-email/internal/config"
	"github.com/MamangRust/monolith-payment-gateway-email/internal/handler"
	"github.com/MamangRust/monolith-payment-gateway-email/internal/mailer"
	"github.com/MamangRust/monolith-payment-gateway-email/internal/metrics"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg := config.Config{
		KafkaBrokers: []string{"localhost:9092"},
		SMTPServer:   "smtp.ethereal.email",
		SMTPPort:     587,
		SMTPUser:     "julius.davis@ethereal.email",
		SMTPPass:     "4vWXpZfTMPAazhVZFU",
	}
	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}

	metrics.Register()
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":2112", nil))
	}()

	m := &mailer.Mailer{
		Server:   cfg.SMTPServer,
		Port:     cfg.SMTPPort,
		User:     cfg.SMTPUser,
		Password: cfg.SMTPPass,
	}

	h := &handler.EmailHandler{Mailer: m}

	myKafka := kafka.NewKafka(logger, cfg.KafkaBrokers)

	err = myKafka.StartConsumers([]string{
		"email-service-topic-auth-register",
		"email-service-topic-auth-forgot-password",
		"email-service-topic-saldo-create",
		"email-service-topic-topup-create",
		"email-service-topic-transfer-create",
		"email-service-topic-merchant-create",
		"email-service-topic-merchant-update-status",
		"email-service-topic-merchant-document-create",
		"email-service-topic-merchant-document-update-status",
	}, "email-service-group", h)

	if err != nil {
		log.Fatalf("Error starting consumer: %v", err)
	}
	select {}
}
