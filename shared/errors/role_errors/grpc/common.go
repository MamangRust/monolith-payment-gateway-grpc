package rolegrpcerrors

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/errors"
	"google.golang.org/grpc/codes"
)

var ErrGrpcRoleInvalidId = errors.NewGrpcError("Invalid Role ID", int(codes.NotFound))
