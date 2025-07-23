package merchantstatshandler

import (
	"context"

	stats "github.com/MamangRust/monolith-payment-gateway-merchant/internal/service/stats"
	statsbyapikey "github.com/MamangRust/monolith-payment-gateway-merchant/internal/service/statsbyapikey"
	statsbymerchant "github.com/MamangRust/monolith-payment-gateway-merchant/internal/service/statsbymerchant"
	pbmerchant "github.com/MamangRust/monolith-payment-gateway-pb/merchant"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/grpc"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/merchant"

	"go.uber.org/zap"
)

type merchantStatsTotalAmountHandleGrpc struct {
	pbmerchant.MerchantStatsTotalAmountServiceServer

	amountstats           stats.MerchantStatsTotalAmountService
	amountstatsbymerchant statsbymerchant.MerchantStatsByMerchantTotalAmountService
	amountstatsbyapikey   statsbyapikey.MerchantStatsByApiKeyTotalAmountService

	logger logger.LoggerInterface
	mapper protomapper.MerchantStatsTotalAmountProtoMapper
}

func NewMerchantStatsTotalAmountHandler(
	amountstats stats.MerchantStatsTotalAmountService,
	amountstatsbymerchant statsbymerchant.MerchantStatsByMerchantTotalAmountService,
	amountstatsbyapikey statsbyapikey.MerchantStatsByApiKeyTotalAmountService,
	logger logger.LoggerInterface,
	mapper protomapper.MerchantStatsTotalAmountProtoMapper,
) MerchantStatsTotalAmountHandleGrpc {
	return &merchantStatsTotalAmountHandleGrpc{
		amountstats:           amountstats,
		amountstatsbymerchant: amountstatsbymerchant,
		amountstatsbyapikey:   amountstatsbyapikey,
		logger:                logger,
		mapper:                mapper,
	}
}

// FindMonthlyTotalAmountMerchant retrieves monthly transaction amounts for a merchant by year.
// It validates the year and returns a gRPC response containing the monthly transaction amounts
// or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to a FindYearMerchant containing the year.
//
// Returns:
//   - A pointer to ApiResponseMerchantMonthlyTotalAmount containing the monthly transaction amounts on success.
//   - An error if the retrieval operation fails.
func (s *merchantStatsTotalAmountHandleGrpc) FindMonthlyTotalAmountMerchant(ctx context.Context, req *pbmerchant.FindYearMerchant) (*pbmerchant.ApiResponseMerchantMonthlyTotalAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		s.logger.Error("invalid year failed", zap.Any("error", merchant_errors.ErrGrpcMerchantInvalidYear))
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	res, err := s.amountstats.FindMonthlyTotalAmountMerchant(ctx, year)

	if err != nil {
		s.logger.Error("FindMonthlyTotalAmountMerchant failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMonthlyTotalAmounts("success", "Successfully fetched monthly amount for merchant", res)

	s.logger.Info("Successfully fetched monthly amount for merchant", zap.Bool("success", true))

	return so, nil
}

// FindYearlyTotalAmountMerchant retrieves yearly transaction amounts for a merchant by year.
// It validates the year and returns a gRPC response containing the yearly transaction amounts
// or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to a FindYearMerchant containing the year.
//
// Returns:
//   - A pointer to ApiResponseMerchantYearlyTotalAmount containing the yearly transaction amounts on success.
//   - An error if the retrieval operation fails.
func (s *merchantStatsTotalAmountHandleGrpc) FindYearlyTotalAmountMerchant(ctx context.Context, req *pbmerchant.FindYearMerchant) (*pbmerchant.ApiResponseMerchantYearlyTotalAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		s.logger.Error("invalid year failed", zap.Any("error", merchant_errors.ErrGrpcMerchantInvalidYear))
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	res, err := s.amountstats.FindYearlyTotalAmountMerchant(ctx, year)

	if err != nil {
		s.logger.Error("FindYearlyTotalAmountMerchant failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseYearlyTotalAmounts("success", "Successfully fetched yearly amount for merchant", res)

	s.logger.Info("Successfully fetched yearly amount for merchant", zap.Bool("success", true))

	return so, nil
}

// FindMonthlyTotalAmountByMerchants retrieves monthly transaction amounts for a merchant by year.
// It validates the year and returns a gRPC response containing the monthly transaction amounts
// or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to a FindYearMerchantById containing the year.
//
// Returns:
//   - A pointer to ApiResponseMerchantMonthlyTotalAmount containing the monthly transaction amounts on success.
//   - An error if the retrieval operation fails.
func (s *merchantStatsTotalAmountHandleGrpc) FindMonthlyTotalAmountByMerchants(ctx context.Context, req *pbmerchant.FindYearMerchantById) (*pbmerchant.ApiResponseMerchantMonthlyTotalAmount, error) {
	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		s.logger.Error("invalid year failed", zap.Any("error", merchant_errors.ErrGrpcMerchantInvalidYear))
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	if id <= 0 {
		s.logger.Error("invalid merchant id failed", zap.Any("error", merchant_errors.ErrGrpcMerchantInvalidID))
		return nil, merchant_errors.ErrGrpcMerchantInvalidID
	}

	res, err := s.amountstatsbymerchant.FindMonthlyTotalAmountByMerchants(ctx, &requests.MonthYearTotalAmountMerchant{
		MerchantID: id,
		Year:       year,
	})

	if err != nil {
		s.logger.Error("FindMonthlyTotalAmountMerchant failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMonthlyTotalAmounts("success", "Successfully fetched monthly amount for merchant", res)

	s.logger.Info("Successfully fetched monthly amount for merchant", zap.Bool("success", true))

	return so, nil
}

// FindYearlyTotalAmountByMerchants retrieves yearly transaction amounts for a merchant by year.
// It validates the year and returns a gRPC response containing the yearly transaction amounts
// or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to a FindYearMerchantById containing the year.
//
// Returns:
//   - A pointer to ApiResponseMerchantYearlyTotalAmount containing the yearly transaction amounts on success.
//   - An error if the retrieval operation fails.
func (s *merchantStatsTotalAmountHandleGrpc) FindYearlyTotalAmountByMerchants(ctx context.Context, req *pbmerchant.FindYearMerchantById) (*pbmerchant.ApiResponseMerchantYearlyTotalAmount, error) {
	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		s.logger.Error("invalid year failed", zap.Any("error", merchant_errors.ErrGrpcMerchantInvalidYear))
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	if id <= 0 {
		s.logger.Error("invalid merchant id failed", zap.Any("error", merchant_errors.ErrGrpcMerchantInvalidID))
		return nil, merchant_errors.ErrGrpcMerchantInvalidID
	}

	res, err := s.amountstatsbymerchant.FindYearlyTotalAmountByMerchants(ctx, &requests.MonthYearTotalAmountMerchant{
		MerchantID: id,
		Year:       year,
	})

	if err != nil {
		s.logger.Error("FindYearlyTotalAmountMerchant failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseYearlyTotalAmounts("success", "Successfully fetched yearly amount for merchant", res)

	s.logger.Info("Successfully fetched yearly amount for merchant", zap.Bool("success", true))

	return so, nil
}

// FindMonthlyTotalAmountByMerchants retrieves monthly transaction amounts for a merchant by year.
// It validates the year and returns a gRPC response containing the monthly transaction amounts
// or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to a FindYearMerchantById containing the year.
//
// Returns:
//   - A pointer to ApiResponseMerchantMonthlyTotalAmount containing the monthly transaction amounts on success.
//   - An error if the retrieval operation fails.
func (s *merchantStatsTotalAmountHandleGrpc) FindMonthlyTotalAmountByApikey(ctx context.Context, req *pbmerchant.FindYearMerchantByApikey) (*pbmerchant.ApiResponseMerchantMonthlyTotalAmount, error) {
	year := int(req.GetYear())
	apikey := req.GetApiKey()

	if year <= 0 {
		s.logger.Error("invalid year failed", zap.Any("error", merchant_errors.ErrGrpcMerchantInvalidYear))
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	res, err := s.amountstatsbyapikey.FindMonthlyTotalAmountByApikeys(ctx, &requests.MonthYearTotalAmountApiKey{
		Year:   year,
		Apikey: apikey,
	})

	if err != nil {
		s.logger.Error("FindMonthlyTotalAmountMerchant failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMonthlyTotalAmounts("success", "Successfully fetched monthly amount for merchant", res)

	s.logger.Info("Successfully fetched monthly amount for merchant", zap.Bool("success", true))

	return so, nil
}

// FindYearlyTotalAmountByApikey retrieves yearly transaction amounts for a merchant by year.
// It validates the year and returns a gRPC response containing the yearly transaction amounts
// or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to a FindYearMerchantByApikey containing the year.
//
// Returns:
//   - A pointer to ApiResponseMerchantYearlyTotalAmount containing the yearly transaction amounts on success.
//   - An error if the retrieval operation fails.
func (s *merchantStatsTotalAmountHandleGrpc) FindYearlyTotalAmountByApikey(ctx context.Context, req *pbmerchant.FindYearMerchantByApikey) (*pbmerchant.ApiResponseMerchantYearlyTotalAmount, error) {
	year := int(req.GetYear())
	apikey := req.GetApiKey()

	if year <= 0 {
		s.logger.Error("invalid year failed", zap.Any("error", merchant_errors.ErrGrpcMerchantInvalidYear))
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	res, err := s.amountstatsbyapikey.FindYearlyTotalAmountByApikeys(ctx, &requests.MonthYearTotalAmountApiKey{
		Apikey: apikey,
		Year:   year,
	})

	if err != nil {
		s.logger.Error("FindYearlyTotalAmountMerchant failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseYearlyTotalAmounts("success", "Successfully fetched yearly amount for merchant", res)

	s.logger.Info("Successfully fetched yearly amount for merchant", zap.Bool("success", true))

	return so, nil
}
