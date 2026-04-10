package apps

import (
	pb "github.com/MamangRust/monolith-payment-gateway-pb/card"
	"github.com/MamangRust/monolith-payment-gateway-card/handler"
	"github.com/MamangRust/monolith-payment-gateway-card/repository"
	"github.com/MamangRust/monolith-payment-gateway-card/service"
	"github.com/MamangRust/monolith-payment-gateway-pkg/server"
	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	repos := repository.NewRepositories(srv.DB)
	svc := service.NewService(&service.Deps{
		Cache:        srv.CacheStore,
		Logger:       srv.Logger,
		Repositories: repos,
	})
	h := handler.NewHandler(svc)

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterCardQueryServiceServer(gs, h)
		pb.RegisterCardCommandServiceServer(gs, h)
	}

	return srv, nil
}
