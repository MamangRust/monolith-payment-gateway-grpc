package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/MamangRust/monolith-payment-gateway-apigateway/internal/app"
)

func main() {
	client, shutdown, err := app.RunClient()
	if err != nil {
		log.Fatalf("failed to run client: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	client.Logger.Info("Gracefully shutting down...")
	shutdown()
}
