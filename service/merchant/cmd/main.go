package main

import (
	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/apps"
)

func main() {
	server, err := apps.NewServer(&apps.Config{
		ServiceName:    "merchant-service",
		ServiceVersion: "1.0.0",
		Environment:    "production",
		OtelEndpoint:   "otel-collector:4317",
	})

	if err != nil {
		panic(err)
	}

	server.Run()
}
