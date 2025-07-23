package main

import (
	"github.com/MamangRust/monolith-payment-gateway-user/internal/apps"
	"go.uber.org/zap"
)

// main starts the gRPC server for the User Service.
//
// It sets up all required dependencies and handles graceful shutdown
// when the service is interrupted.
func main() {
	server, shutdown, err := apps.NewServer()

	if err != nil {
		server.Logger.Fatal("Failed to create server", zap.Error(err))
		panic(err)
	}

	defer func() {
		if err := shutdown(server.Ctx); err != nil {
			server.Logger.Error("Failed to shutdown tracer", zap.Error(err))
		}
	}()

	server.Run()
}
