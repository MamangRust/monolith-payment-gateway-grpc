package transactiongrpcerrors

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/errors"
	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcValidateCreateTransactionRequest = errors.NewGrpcError("Invalid input for create card", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateTransactionRequest = errors.NewGrpcError("Invalid input for update card", int(codes.InvalidArgument))
)
