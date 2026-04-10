package apps

import (
	pbcard "github.com/MamangRust/monolith-payment-gateway-pb/card"
	pbmerchant "github.com/MamangRust/monolith-payment-gateway-pb/merchant"
	pbsaldo "github.com/MamangRust/monolith-payment-gateway-pb/saldo"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/transaction"
	pbstats "github.com/MamangRust/monolith-payment-gateway-pb/transaction/stats"
	"github.com/MamangRust/monolith-payment-gateway-pkg/adapter"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/server"
	"github.com/MamangRust/monolith-payment-gateway-transaction/handler"
	"github.com/MamangRust/monolith-payment-gateway-transaction/repository"
	"github.com/MamangRust/monolith-payment-gateway-transaction/service"
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
	connMerchant, _ := grpc.NewClient(viper.GetString("GRPC_MERCHANT_ADDR"))

	saldoClientQuery := pbsaldo.NewSaldoQueryServiceClient(connSaldo)
	saldoClientCmd := pbsaldo.NewSaldoCommandServiceClient(connSaldo)
	cardClientQuery := pbcard.NewCardQueryServiceClient(connCard)
	cardClientCmd := pbcard.NewCardCommandServiceClient(connCard)
	merchantClientQuery := pbmerchant.NewMerchantQueryServiceClient(connMerchant)

	saldoAdapter := adapter.NewSaldoAdapter(saldoClientQuery, saldoClientCmd)
	cardAdapter := adapter.NewCardAdapter(cardClientQuery, cardClientCmd)
	merchantAdapter := adapter.NewMerchantAdapter(merchantClientQuery)

	repos := repository.NewRepositories(srv.DB, saldoAdapter, cardAdapter, merchantAdapter)
	myKafka := kafka.NewKafka(srv.Logger, []string{viper.GetString("KAFKA_BROKERS")})
	svc := service.NewService(&service.Deps{
		Kafka:        myKafka,
		Repositories: repos,
		Logger:       srv.Logger,
		Cache:        srv.CacheStore,
	})
	h := handler.NewHandler(svc)

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterTransactionQueryServiceServer(gs, h)
		pb.RegisterTransactionCommandServiceServer(gs, h)
		pbstats.RegisterTransactionStatsAmountServiceServer(gs, h)
		pbstats.RegisterTransactionStatsMethodServiceServer(gs, h)
		pbstats.RegisterTransactionStatsStatusServiceServer(gs, h)
	}

	return srv, nil
}
