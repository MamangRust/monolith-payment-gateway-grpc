package withdrawgrpcerrors

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/errors"
	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcValidateCreateWithdrawRequest = errors.NewGrpcError("Invalid input for create withdraw", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateWithdrawRequest = errors.NewGrpcError("Invalid input for update withdraw", int(codes.InvalidArgument))
)
