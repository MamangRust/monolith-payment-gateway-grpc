package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/MamangRust/monolith-payment-gateway-email/internal/config"
	"github.com/MamangRust/monolith-payment-gateway-email/internal/handler"
	"github.com/MamangRust/monolith-payment-gateway-email/internal/mailer"
	"github.com/MamangRust/monolith-payment-gateway-email/internal/metrics"
	"github.com/MamangRust/monolith-payment-gateway-pkg/dotenv"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	otel_pkg "github.com/MamangRust/monolith-payment-gateway-pkg/otel"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// main initializes and starts the email service. It sets up the logger, loads
// environment configurations, initializes the tracer provider, registers the
// metrics endpoint, and starts the Kafka consumers for various topics related
// to email notifications. The function will block indefinitely after starting
// the service.
func main() {
	logger, err := logger.NewLogger("email-service")
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}

	if err := dotenv.Viper(); err != nil {
		logger.Fatal("Failed to load .env file", zap.Error(err))
	}

	ctx := context.Background()

	cfg := config.Config{
		KafkaBrokers: []string{viper.GetString("KAFKA_BROKERS")},
		SMTPServer:   viper.GetString("SMTP_SERVER"),
		SMTPPort:     viper.GetInt("SMTP_PORT"),
		SMTPUser:     viper.GetString("SMTP_USER"),
		SMTPPass:     viper.GetString("SMTP_PASS"),
	}

	metricsAddr := fmt.Sprintf(":%s", viper.GetString("METRIC_EMAIL_ADDR"))

	metrics.Register()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(metricsAddr, nil))
	}()

	shutdownTracerProvider, err := otel_pkg.InitTracerProvider("email-service", ctx)

	if err != nil {
		logger.Fatal("Failed to initialize tracer provider", zap.Error(err))
	}

	defer func() {
		if err := shutdownTracerProvider(ctx); err != nil {
			logger.Fatal("Failed to shutdown tracer provider", zap.Error(err))
		}
	}()

	m := mailer.NewMailer(
		ctx,
		cfg.SMTPServer,
		cfg.SMTPPort,
		cfg.SMTPUser,
		cfg.SMTPPass,
		logger,
	)

	h := handler.NewEmailHandler(ctx, logger, m)

	myKafka := kafka.NewKafka(logger, cfg.KafkaBrokers)

	err = myKafka.StartConsumers([]string{
		"email-service-topic-auth-register",
		"email-service-topic-auth-forgot-password",
		"email-service-topic-auth-verify-code-success",
		"email-service-topic-saldo-create",
		"email-service-topic-topup-create",
		"email-service-topic-transaction-create",
		"email-service-topic-transfer-create",
		"email-service-topic-merchant-create",
		"email-service-topic-merchant-update-status",
		"email-service-topic-merchant-document-create",
		"email-service-topic-merchant-document-update-status",
	}, "email-service-group", h)

	if err != nil {
		log.Fatalf("Error starting consumer: %v", err)
	}

	logger.Info("Email service started", zap.String("service", "email-service"))

	select {}

}
