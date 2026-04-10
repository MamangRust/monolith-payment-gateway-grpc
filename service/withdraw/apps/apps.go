package apps

import (
	pbcard "github.com/MamangRust/monolith-payment-gateway-pb/card"
	pbsaldo "github.com/MamangRust/monolith-payment-gateway-pb/saldo"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/withdraw"
	pbstats "github.com/MamangRust/monolith-payment-gateway-pb/withdraw/stats"
	"github.com/MamangRust/monolith-payment-gateway-pkg/adapter"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/server"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/handler"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/repository"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/service"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	// gRPC Clients for cross-service communication
	connSaldo, _ := grpc.NewClient(viper.GetString("GRPC_SALDO_ADDR"))
	connCard, _ := grpc.NewClient(viper.GetString("GRPC_CARD_ADDR"))

	saldoClientQuery := pbsaldo.NewSaldoQueryServiceClient(connSaldo)
	saldoClientCmd := pbsaldo.NewSaldoCommandServiceClient(connSaldo)
	cardClientQuery := pbcard.NewCardQueryServiceClient(connCard)
	cardClientCmd := pbcard.NewCardCommandServiceClient(connCard)

	saldoAdapter := adapter.NewSaldoAdapter(saldoClientQuery, saldoClientCmd)
	cardAdapter := adapter.NewCardAdapter(cardClientQuery, cardClientCmd)

	repos := repository.NewRepositories(srv.DB, cardAdapter, saldoAdapter)
	myKafka := kafka.NewKafka(srv.Logger, []string{viper.GetString("KAFKA_BROKERS")})

	svc := service.NewService(&service.Deps{
		Kafka:        myKafka,
		Repositories: repos,
		Logger:       srv.Logger,
		Cache:        srv.CacheStore,
	})
	h := handler.NewHandler(svc)

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterWithdrawQueryServiceServer(gs, h)
		pb.RegisterWithdrawCommandServiceServer(gs, h)
		pbstats.RegisterWithdrawStatsAmountServiceServer(gs, h)
		pbstats.RegisterWithdrawStatsStatusServiceServer(gs, h)
	}

	return srv, nil
}
