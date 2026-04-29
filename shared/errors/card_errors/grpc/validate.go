package cardgrpcerrors

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/errors"
	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcValidateCreateCardRequest = errors.NewGrpcError("Invalid input for create card", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateCardRequest = errors.NewGrpcError("Invalid input for update card", int(codes.InvalidArgument))
)
