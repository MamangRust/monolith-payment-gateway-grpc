package role_test

import (
	"context"
	"net"
	"testing"

	pb "github.com/MamangRust/monolith-payment-gateway-pb/role"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-role/handler"
	"github.com/MamangRust/monolith-payment-gateway-role/repository"
	"github.com/MamangRust/monolith-payment-gateway-role/service"
	"github.com/MamangRust/monolith-payment-gateway-shared/cache"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	tests "github.com/MamangRust/monolith-payment-gateway-test"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RoleGapiTestSuite struct {
	suite.Suite
	ts            *tests.TestSuite
	grpcServer    *grpc.Server
	commandClient pb.RoleCommandServiceClient
	queryClient   pb.RoleServiceClient
	conn          *grpc.ClientConn
	roleID        int
}

func (s *RoleGapiTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)

	opts, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	redisClient := redis.NewClient(opts)

	queries := db.New(pool)
	repos := repository.NewRepositories(queries)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test", lp)
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(redisClient, log, cacheMetrics)

	roleService := service.NewService(&service.Deps{
		Repositories: repos,
		Logger:       log,
		Cache:        cacheStore,
	})

	roleHandler := handler.NewHandler(roleService)
	server := grpc.NewServer()
	pb.RegisterRoleCommandServiceServer(server, roleHandler.RoleCommand)
	pb.RegisterRoleServiceServer(server, roleHandler.RoleQuery)
	s.grpcServer = server

	lis, err := net.Listen("tcp", ":0")
	s.Require().NoError(err)

	go func() {
		_ = server.Serve(lis)
	}()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.commandClient = pb.NewRoleCommandServiceClient(conn)
	s.queryClient = pb.NewRoleServiceClient(conn)
}

func (s *RoleGapiTestSuite) TearDownSuite() {
	s.conn.Close()
	s.grpcServer.Stop()
	s.ts.Teardown()
}

func (s *RoleGapiTestSuite) Test1_CreateRole() {
	ctx := context.Background()
	req := &pb.CreateRoleRequest{
		Name: "Gapi Role",
	}

	res, err := s.commandClient.CreateRole(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal(req.Name, res.Data.Name)
	s.roleID = int(res.Data.Id)
}

func (s *RoleGapiTestSuite) Test2_FindById() {
	s.Require().NotZero(s.roleID)
	ctx := context.Background()

	req := &pb.FindByIdRoleRequest{
		RoleId: int32(s.roleID),
	}

	res, err := s.queryClient.FindByIdRole(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal(int32(s.roleID), res.Data.Id)
}

func (s *RoleGapiTestSuite) Test3_UpdateRole() {
	s.Require().NotZero(s.roleID)
	ctx := context.Background()

	req := &pb.UpdateRoleRequest{
		Id:   int32(s.roleID),
		Name: "Updated Gapi Role",
	}

	res, err := s.commandClient.UpdateRole(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal("Updated Gapi Role", res.Data.Name)
}

func (s *RoleGapiTestSuite) Test4_TrashAndRestore() {
	s.Require().NotZero(s.roleID)
	ctx := context.Background()

	req := &pb.FindByIdRoleRequest{
		RoleId: int32(s.roleID),
	}

	_, err := s.commandClient.TrashedRole(ctx, req)
	s.NoError(err)

	_, err = s.commandClient.RestoreRole(ctx, req)
	s.NoError(err)
}

func (s *RoleGapiTestSuite) Test5_DeleteRolePermanent() {
	s.Require().NotZero(s.roleID)
	ctx := context.Background()

	req := &pb.FindByIdRoleRequest{
		RoleId: int32(s.roleID),
	}

	_, err := s.commandClient.TrashedRole(ctx, req)
	s.NoError(err)

	_, err = s.commandClient.DeleteRolePermanent(ctx, req)
	s.NoError(err)
}

func TestRoleGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(RoleGapiTestSuite))
}
