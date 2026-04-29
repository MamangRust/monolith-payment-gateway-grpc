package saldogrpcerrors

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/errors"
	"google.golang.org/grpc/codes"
)

var (
	ErrGrpcValidateCreateSaldo         = errors.NewGrpcError("Invalid input for create saldo", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateSaldo         = errors.NewGrpcError("Invalid input for update saldo", int(codes.InvalidArgument))
	ErrGrpcValidateUpdateSaldoWithdraw = errors.NewGrpcError("Invalid input for update saldo withdraw", int(codes.InvalidArgument))
)
