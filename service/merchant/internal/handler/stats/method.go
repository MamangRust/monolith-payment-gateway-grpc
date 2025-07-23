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

type merchantStatsMethodHandleGrpc struct {
	pbmerchant.MerchantStatsMethodServiceServer

	methodstats           stats.MerchantStatsMethodService
	methodstatsbymerchant statsbymerchant.MerchantStatsByMerchantMethodService
	methodstatsbyapikey   statsbyapikey.MerchantStatsByApiKeyMethodService

	logger logger.LoggerInterface
	mapper protomapper.MerchantStatsMethodProtoMapper
}

func NewMerchantStatsMethodHandler(
	methodstats stats.MerchantStatsMethodService,
	methodstatsbymerchant statsbymerchant.MerchantStatsByMerchantMethodService,
	methodstatsbyapikey statsbyapikey.MerchantStatsByApiKeyMethodService,
	logger logger.LoggerInterface,
	mapper protomapper.MerchantStatsMethodProtoMapper,
) MerchantStatsMethodHandleGrpc {
	return &merchantStatsMethodHandleGrpc{
		methodstats:           methodstats,
		methodstatsbymerchant: methodstatsbymerchant,
		methodstatsbyapikey:   methodstatsbyapikey,
		logger:                logger,
		mapper:                mapper,
	}
}

// FindMonthlyPaymentMethodsMerchant retrieves monthly payment methods for a merchant by year.
// It handles invalid years and returns a gRPC response containing the monthly payment methods
// or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancelation signals, and deadlines.
//   - req: A pointer to a FindYearMerchant containing the year.
//
// Returns:
//   - A pointer to ApiResponseMerchantMonthlyPaymentMethod containing the monthly payment methods on success.
//   - An error if the retrieval operation fails.
func (s *merchantStatsMethodHandleGrpc) FindMonthlyPaymentMethodsMerchant(ctx context.Context, req *pbmerchant.FindYearMerchant) (*pbmerchant.ApiResponseMerchantMonthlyPaymentMethod, error) {
	year := int(req.GetYear())

	if year <= 0 {
		s.logger.Error("invalid year failed", zap.Any("error", merchant_errors.ErrGrpcMerchantInvalidYear))
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	res, err := s.methodstats.FindMonthlyPaymentMethodsMerchant(ctx, year)

	if err != nil {
		s.logger.Error("FindMonthlyPaymentMethodsMerchant failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMonthlyPaymentMethods("success", "Successfully fetched monthly payment methods for merchant", res)

	s.logger.Info("Successfully fetched monthly payment methods for merchant", zap.Bool("success", true))

	return so, nil
}

// FindYearlyPaymentMethodMerchant retrieves yearly payment methods for a merchant by year.
// It handles invalid years and returns a gRPC response containing the yearly payment methods
// or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to a FindYearMerchant containing the year.
//
// Returns:
//   - A pointer to ApiResponseMerchantYearlyPaymentMethod containing the yearly payment methods on success.
//   - An error if the retrieval operation fails.
func (s *merchantStatsMethodHandleGrpc) FindYearlyPaymentMethodMerchant(ctx context.Context, req *pbmerchant.FindYearMerchant) (*pbmerchant.ApiResponseMerchantYearlyPaymentMethod, error) {
	year := int(req.GetYear())

	if year <= 0 {
		s.logger.Error("invalid year failed", zap.Any("error", merchant_errors.ErrGrpcMerchantInvalidYear))
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	res, err := s.methodstats.FindYearlyPaymentMethodMerchant(ctx, year)

	if err != nil {
		s.logger.Error("FindYearlyPaymentMethodMerchant failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseYearlyPaymentMethods("success", "Successfully fetched yearly payment methods for merchant", res)

	s.logger.Info("Successfully fetched yearly payment methods for merchant", zap.Bool("success", true))

	return so, nil
}

// FindMonthlyPaymentMethodByMerchants retrieves monthly payment methods for a specific merchant by year.
// It validates the merchant ID and year, and returns a gRPC response containing the monthly payment methods
// or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to a FindYearMerchantById containing the merchant ID and year.
//
// Returns:
//   - A pointer to ApiResponseMerchantMonthlyPaymentMethod containing the monthly payment methods on success.
//   - An error if the retrieval operation fails.
func (s *merchantStatsMethodHandleGrpc) FindMonthlyPaymentMethodByMerchants(ctx context.Context, req *pbmerchant.FindYearMerchantById) (*pbmerchant.ApiResponseMerchantMonthlyPaymentMethod, error) {
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

	reqService := requests.MonthYearPaymentMethodMerchant{
		MerchantID: int(req.MerchantId),
		Year:       int(year),
	}

	res, err := s.methodstatsbymerchant.FindMonthlyPaymentMethodByMerchants(ctx, &reqService)

	if err != nil {
		s.logger.Error("FindMonthlyPaymentMethodByMerchants failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMonthlyPaymentMethods("success", "Successfully fetched monthly payment methods by merchant", res)

	s.logger.Info("Successfully fetched monthly payment methods by merchant", zap.Bool("success", true))

	return so, nil
}

// FindYearlyPaymentMethodByMerchants retrieves yearly payment methods for a specific merchant by year.
// It validates the merchant ID and year, and returns a gRPC response containing the yearly payment methods
// or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to a FindYearMerchantById containing the merchant ID and year.
//
// Returns:
//   - A pointer to ApiResponseMerchantYearlyPaymentMethod containing the yearly payment methods on success.
//   - An error if the retrieval operation fails.
func (s *merchantStatsMethodHandleGrpc) FindYearlyPaymentMethodByMerchants(ctx context.Context, req *pbmerchant.FindYearMerchantById) (*pbmerchant.ApiResponseMerchantYearlyPaymentMethod, error) {
	merchantId := req.GetMerchantId()
	year := req.GetYear()

	if year <= 0 {
		s.logger.Error("Invalid year failed", zap.Any("error", merchant_errors.ErrGrpcMerchantInvalidYear))
		return nil, merchant_errors.ErrGrpcMerchantInvalidYear
	}

	if merchantId <= 0 {
		s.logger.Error("Invalid id failed", zap.Any("error", merchant_errors.ErrGrpcMerchantInvalidID))
		return nil, merchant_errors.ErrGrpcMerchantInvalidID
	}

	reqService := requests.MonthYearPaymentMethodMerchant{
		MerchantID: int(req.MerchantId),
		Year:       int(year),
	}

	res, err := s.methodstatsbymerchant.FindYearlyPaymentMethodByMerchants(ctx, &reqService)

	if err != nil {
		s.logger.Error("FindYearlyPaymentMethodByMerchants failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseYearlyPaymentMethods("success", "Successfully fetched yearly payment methods by merchant", res)

	s.logger.Info("Successfully fetched yearly payment methods by merchant", zap.Bool("success", true))

	return so, nil
}

// FindMonthlyPaymentMethodByApikey retrieves a merchant's monthly payment methods by API key and year.
// It validates the API key and year, and returns a gRPC response containing the monthly payment methods
// on success or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to a FindYearMerchantByApikey containing the API key and year.
//
// Returns:
//   - A pointer to ApiResponseMerchantMonthlyPaymentMethod containing the monthly payment methods
//     on success.
//   - An error if the retrieval operation fails.
func (s *merchantStatsMethodHandleGrpc) FindMonthlyPaymentMethodByApikey(ctx context.Context, req *pbmerchant.FindYearMerchantByApikey) (*pbmerchant.ApiResponseMerchantMonthlyPaymentMethod, error) {
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

	reqService := requests.MonthYearPaymentMethodApiKey{
		Year:   int(year),
		Apikey: api_key,
	}

	res, err := s.methodstatsbyapikey.FindMonthlyPaymentMethodByApikeys(ctx, &reqService)

	if err != nil {
		s.logger.Error("FindMonthlyPaymentMethodByApikey failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseMonthlyPaymentMethods("success", "Successfully fetched monthly payment methods by merchant", res)

	s.logger.Info("Successfully fetched monthly payment methods by merchant", zap.Bool("success", true))

	return so, nil
}

// FindYearlyPaymentMethodByApikey retrieves a merchant's yearly payment methods by API key and year.
// It validates the API key and year, and returns a gRPC response containing the yearly payment methods
// on success or an error if the operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to a FindYearMerchantByApikey containing the API key and year.
//
// Returns:
//   - A pointer to ApiResponseMerchantYearlyPaymentMethod containing the yearly payment methods
//     on success.
//   - An error if the retrieval operation fails.
func (s *merchantStatsMethodHandleGrpc) FindYearlyPaymentMethodByApikey(ctx context.Context, req *pbmerchant.FindYearMerchantByApikey) (*pbmerchant.ApiResponseMerchantYearlyPaymentMethod, error) {
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

	reqService := requests.MonthYearPaymentMethodApiKey{
		Year:   int(year),
		Apikey: api_key,
	}

	res, err := s.methodstatsbyapikey.FindYearlyPaymentMethodByApikeys(ctx, &reqService)

	if err != nil {
		s.logger.Error("FindYearlyPaymentMethodByApikey failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseYearlyPaymentMethods("success", "Successfully fetched yearly payment methods by merchant", res)

	s.logger.Info("Successfully fetched yearly payment methods by merchant", zap.Bool("success", true))

	return so, nil
}
