package topupgrpcerrors

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/errors"
	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcTopupInvalidID    = errors.NewGrpcError("Invalid Topup ID", int(codes.InvalidArgument))
	ErrGrpcTopupInvalidMonth = errors.NewGrpcError("Invalid Topup Month", int(codes.InvalidArgument))
	ErrGrpcInvalidCardNumber = errors.NewGrpcError("Invalid card number", int(codes.InvalidArgument))
	ErrGrpcTopupInvalidYear  = errors.NewGrpcError("Invalid Topup Year", int(codes.InvalidArgument))
)
