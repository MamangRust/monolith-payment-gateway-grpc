package main

import "github.com/MamangRust/monolith-payment-gateway-saldo/internal/apps"

func main() {
	server, err := apps.NewServer()

	if err != nil {
		panic(err)
	}

	server.Run()
}
