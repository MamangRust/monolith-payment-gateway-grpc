package apps

import (
	pb "github.com/MamangRust/monolith-payment-gateway-pb/role"
	"github.com/MamangRust/monolith-payment-gateway-pkg/server"
	"github.com/MamangRust/monolith-payment-gateway-role/handler"
	"github.com/MamangRust/monolith-payment-gateway-role/repository"
	"github.com/MamangRust/monolith-payment-gateway-role/service"
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
		pb.RegisterRoleServiceServer(gs, h.RoleQuery)
		pb.RegisterRoleCommandServiceServer(gs, h.RoleCommand)
	}

	return srv, nil
}
