package main

import (
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/apps"
)

func main() {
	server, err := apps.NewServer(&apps.Config{
		ServiceName:    "transaction-service",
		ServiceVersion: "1.0.0",
		Environment:    "production",
		OtelEndpoint:   "otel-collector:4317",
	})

	if err != nil {
		panic(err)
	}

	server.Run()
}
