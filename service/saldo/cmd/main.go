package main

import (
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/apps"
)

func main() {
	server, err := apps.NewServer(&apps.Config{
		ServiceName:    "saldo-service",
		ServiceVersion: "1.0.0",
		Environment:    "production",
		OtelEndpoint:   "otel-collector:4317",
	})

	if err != nil {
		panic(err)
	}

	server.Run()
}
