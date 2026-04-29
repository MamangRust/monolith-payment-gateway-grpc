package merchantdocumentgrpcerrors

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/errors"
	"google.golang.org/grpc/codes"
)

var ErrGrpcMerchantInvalidID = errors.NewGrpcError("Invalid merchant id", int(codes.InvalidArgument))
