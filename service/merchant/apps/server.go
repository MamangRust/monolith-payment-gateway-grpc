package apps

import (
	pb "github.com/MamangRust/monolith-payment-gateway-pb/merchant"
	"github.com/MamangRust/monolith-payment-gateway-merchant/handler"
	"github.com/MamangRust/monolith-payment-gateway-merchant/repository"
	"github.com/MamangRust/monolith-payment-gateway-merchant/service"
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
		pb.RegisterMerchantQueryServiceServer(gs, h)
		pb.RegisterMerchantCommandServiceServer(gs, h)
	}

	return srv, nil
}
