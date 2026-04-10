package apps

import (
	pb "github.com/MamangRust/monolith-payment-gateway-pb/saldo"
	pbstats "github.com/MamangRust/monolith-payment-gateway-pb/saldo/stats"
	"github.com/MamangRust/monolith-payment-gateway-saldo/handler"
	"github.com/MamangRust/monolith-payment-gateway-saldo/repository"
	"github.com/MamangRust/monolith-payment-gateway-saldo/service"
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
		pb.RegisterSaldoQueryServiceServer(gs, h)
		pb.RegisterSaldoCommandServiceServer(gs, h)
		pbstats.RegisterSaldoStatsBalanceServiceServer(gs, h)
		pbstats.RegisterSaldoStatsTotalBalanceServer(gs, h)
	}

	return srv, nil
}
