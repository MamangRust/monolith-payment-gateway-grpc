package apps

import (
	"fmt"

	"github.com/MamangRust/monolith-payment-gateway-auth/handler"
	"github.com/MamangRust/monolith-payment-gateway-auth/repository"
	"github.com/MamangRust/monolith-payment-gateway-auth/service"

	pb "github.com/MamangRust/monolith-payment-gateway-pb"
	"github.com/MamangRust/monolith-payment-gateway-pkg/auth"
	"github.com/MamangRust/monolith-payment-gateway-pkg/hash"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/server"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	tokenManager, err := auth.NewManager(viper.GetString("SECRET_KEY"))
	if err != nil {
		return nil, fmt.Errorf("failed to create token manager: %w", err)
	}

	hasher := hash.NewHashingPassword()
	repositories := repository.NewRepositories(srv.DB)
	myKafka := kafka.NewKafka(srv.Logger, []string{viper.GetString("KAFKA_BROKERS")})

	services := service.NewService(&service.Deps{
		Cache:        srv.CacheStore,
		Repositories: repositories,
		Token:        tokenManager,
		Hash:         hasher,
		Logger:       srv.Logger,
		Kafka:        myKafka,
	})

	handlers := handler.NewHandler(&handler.Deps{Service: services, Logger: srv.Logger})

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterAuthServiceServer(gs, handlers.Auth)
	}

	return srv, nil
}
