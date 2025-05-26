//go:generate mockgen -source=interfaces.go -destination=mocks/mock_auth_handle.go -package=mocks
package handler

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-shared/pb"
)

type AuthHandleGrpc interface {
	pb.AuthServiceServer
	LoginUser(ctx context.Context, req *pb.LoginRequest) (*pb.ApiResponseLogin, error)
	RegisterUser(ctx context.Context, req *pb.RegisterRequest) (*pb.ApiResponseRegister, error)
}
