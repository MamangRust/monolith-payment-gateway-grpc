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

type merchantStatsAmountHandleGrpc struct {
	pbmerchant.MerchantStatsAmountServiceServer

	amountstats           stats.MerchantStatsAmountService
	amountstatsbymerchant statsbymerchant.MerchantStatsByMerchantAmountService
	amountstatsbyapikey   statsbyapikey.MerchantStatsByApiKeyAmountService

	logger logger.LoggerInterface
	mapper protomapper.MerchantStatsAmountProtoMapper
}

func NewMerchantStatsAmountHandler(
	amountstats stats.MerchantStatsAmountService,
	amountstatsbymerchant statsbymerchant.MerchantStatsByMerchantAmountService,
	amountstatsbyapikey statsbyapikey.MerchantStatsByApiKeyAmountService,
	logger logger.LoggerInterface,
	mapper protomapper.MerchantStatsAmountProtoMapper,
) MerchantStatsAmountHandleGrpc {
	return &merchantStatsAmountHandleGrpc{
		amountstats:           amountstats,
		amountstatsbymerchant: amountstatsbymerchant,
		amountstatsbyapikey:   amountstatsbyapikey,
		logger:                logger,
		mapper:                mapper,
	}
}

// FindMonthlyAmountMerchant retrieves monthly transaction amounts for a merchant by year.
// It handles invalid years and returns a gRPC response containing the monthly transaction amounts
// or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to a FindYearMerchant containing the year.
//
// Returns:
//   - A pointer to ApiResponseMerchantMonthlyAmount containing the monthly transaction amounts on success.
//   - An error if the retrieval operation fails.
func (s *merchantStatsAmountHandleGrpc) FindMonthlyAmountMerchant(ctx context.Context, req *pbmerchant.FindYearMerchant) (*pbmerchant.ApiResponseMerchantMonthlyAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		s.logger.Error("invalid year failed", zap.Any("error", merchant_errors.ErrGrpcMerchantInvalidYear))
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	res, err := s.amountstats.FindMonthlyAmountMerchant(ctx, year)

	if err != nil {
		s.logger.Error("FindMonthlyAmountMerchant failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMonthlyAmounts("success", "Successfully fetched monthly amount for merchant", res)

	s.logger.Info("Successfully fetched monthly amount for merchant", zap.Bool("success", true))

	return so, nil
}

// FindYearlyAmountMerchant retrieves yearly transaction amounts for a merchant by year.
// It validates the year and returns a gRPC response containing the yearly transaction amounts
// or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to a FindYearMerchant containing the year.
//
// Returns:
//   - A pointer to ApiResponseMerchantYearlyAmount containing the yearly transaction amounts on success.
//   - An error if the retrieval operation fails.
func (s *merchantStatsAmountHandleGrpc) FindYearlyAmountMerchant(ctx context.Context, req *pbmerchant.FindYearMerchant) (*pbmerchant.ApiResponseMerchantYearlyAmount, error) {
	year := int(req.GetYear())

	if year <= 0 {
		s.logger.Error("invalid year failed", zap.Any("error", merchant_errors.ErrGrpcMerchantInvalidYear))
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	res, err := s.amountstats.FindYearlyAmountMerchant(ctx, year)

	if err != nil {
		s.logger.Error("FindYearlyAmountMerchant failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseYearlyAmounts("success", "Successfully fetched yearly amount for merchant", res)

	return so, nil
}

// FindMonthlyAmountByMerchants retrieves monthly transaction amounts for a specific merchant by year.
// It validates the merchant ID and year, and returns a gRPC response containing the monthly transaction amounts
// or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to a FindYearMerchantById containing the merchant ID and year.
//
// Returns:
//   - A pointer to ApiResponseMerchantMonthlyAmount containing the monthly transaction amounts on success.
//   - An error if the retrieval operation fails.
func (s *merchantStatsAmountHandleGrpc) FindMonthlyAmountByMerchants(ctx context.Context, req *pbmerchant.FindYearMerchantById) (*pbmerchant.ApiResponseMerchantMonthlyAmount, error) {
	merchantId := req.GetMerchantId()
	year := req.GetYear()

	if year <= 0 {
		s.logger.Error("invalid year failed", zap.Any("error", merchant_errors.ErrGrpcMerchantInvalidYear))
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	if merchantId <= 0 {
		s.logger.Error("invalid id failed", zap.Any("error", merchant_errors.ErrGrpcMerchantInvalidID))
		return nil, merchant_errors.ErrGrpcMerchantInvalidID
	}

	reqService := requests.MonthYearAmountMerchant{
		MerchantID: int(req.MerchantId),
		Year:       int(year),
	}

	res, err := s.amountstatsbymerchant.FindMonthlyAmountByMerchants(ctx, &reqService)

	if err != nil {
		s.logger.Error("FindMonthlyAmountByMerchants failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMonthlyAmounts("success", "Successfully fetched monthly amount by merchant", res)

	s.logger.Info("Successfully fetched monthly amount by merchant", zap.Bool("success", true))

	return so, nil
}

// FindYearlyAmountByMerchants retrieves yearly transaction amounts for a specific merchant by year.
// It validates the merchant ID and year, and returns a gRPC response containing the yearly transaction amounts
// or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to a FindYearMerchantById containing the merchant ID and year.
//
// Returns:
//   - A pointer to ApiResponseMerchantYearlyAmount containing the yearly transaction amounts on success.
//   - An error if the retrieval operation fails.
func (s *merchantStatsAmountHandleGrpc) FindYearlyAmountByMerchants(ctx context.Context, req *pbmerchant.FindYearMerchantById) (*pbmerchant.ApiResponseMerchantYearlyAmount, error) {
	merchantId := req.GetMerchantId()
	year := req.GetYear()

	if year <= 0 {
		s.logger.Error("invalid year failed", zap.Any("error", merchant_errors.ErrGrpcMerchantInvalidYear))
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	if merchantId <= 0 {
		s.logger.Error("invalid id failed", zap.Any("error", merchant_errors.ErrGrpcMerchantInvalidID))
		return nil, merchant_errors.ErrGrpcMerchantInvalidID
	}

	reqService := requests.MonthYearAmountMerchant{
		MerchantID: int(req.MerchantId),
		Year:       int(year),
	}
	res, err := s.amountstatsbymerchant.FindYearlyAmountByMerchants(ctx, &reqService)

	if err != nil {
		s.logger.Error("FindYearlyAmountByMerchants failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseYearlyAmounts("success", "Successfully fetched yearly amount by merchant", res)

	s.logger.Info("Successfully fetched yearly amount by merchant", zap.Bool("success", true))

	return so, nil
}

// FindMonthlyAmountByApikey retrieves a merchant's monthly amount by API key and year.
// It validates the API key and year, and returns a gRPC response containing the monthly amount
// on success or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to a FindYearMerchantByApikey containing the API key and year.
//
// Returns:
//   - A pointer to ApiResponseMerchantMonthlyAmount containing the monthly amount
//     on success.
//   - An error if the retrieval operation fails.
func (s *merchantStatsAmountHandleGrpc) FindMonthlyAmountByApikey(ctx context.Context, req *pbmerchant.FindYearMerchantByApikey) (*pbmerchant.ApiResponseMerchantMonthlyAmount, error) {
	api_key := req.GetApiKey()
	year := req.GetYear()

	if year <= 0 {
		s.logger.Error("invalid year failed", zap.Any("error", merchant_errors.ErrGrpcMerchantInvalidYear))
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	if api_key == "" {
		s.logger.Error("invalid api key failed", zap.Any("error", merchant_errors.ErrGrpcMerchantInvalidApiKey))
		return nil, merchant_errors.ErrGrpcMerchantInvalidApiKey
	}

	reqService := requests.MonthYearAmountApiKey{
		Apikey: api_key,
		Year:   int(year),
	}

	res, err := s.amountstatsbyapikey.FindMonthlyAmountByApikeys(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthlyAmountByApikey failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMonthlyAmounts("success", "Successfully fetched monthly amount by merchant", res)

	s.logger.Info("Successfully fetched monthly amount by merchant", zap.Bool("success", true))

	return so, nil
}

// FindYearlyAmountByApikey retrieves a merchant's yearly amount by API key and year.
// It validates the API key and year, and returns a gRPC response containing the yearly amount
// on success or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to a FindYearMerchantByApikey containing the API key and year.
//
// Returns:
//   - A pointer to ApiResponseMerchantYearlyAmount containing the yearly amount
//     on success.
//   - An error if the retrieval operation fails.
func (s *merchantStatsAmountHandleGrpc) FindYearlyAmountByApikey(ctx context.Context, req *pbmerchant.FindYearMerchantByApikey) (*pbmerchant.ApiResponseMerchantYearlyAmount, error) {
	api_key := req.GetApiKey()
	year := req.GetYear()

	if year <= 0 {
		s.logger.Error("invalid year failed", zap.Any("error", merchant_errors.ErrGrpcMerchantInvalidYear))
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	if api_key == "" {
		s.logger.Error("invalid api key failed", zap.Any("error", merchant_errors.ErrGrpcMerchantInvalidApiKey))
		return nil, merchant_errors.ErrGrpcMerchantInvalidApiKey
	}

	reqService := requests.MonthYearAmountApiKey{
		Apikey: api_key,
		Year:   int(year),
	}

	res, err := s.amountstatsbyapikey.FindYearlyAmountByApikeys(ctx, &reqService)

	if err != nil {
		s.logger.Error("FindYearlyAmountByApikey failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseYearlyAmounts("success", "Successfully fetched yearly amount by merchant", res)

	s.logger.Info("Successfully fetched yearly amount by merchant", zap.Bool("success", true))

	return so, nil
}
