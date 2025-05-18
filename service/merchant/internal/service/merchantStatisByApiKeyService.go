package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-merchant/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantStatisByApiKeyService struct {
	ctx                              context.Context
	trace                            trace.Tracer
	merchantStatisByApiKeyRepository repository.MerchantStatisticByApiKeyRepository
	logger                           logger.LoggerInterface
	mapping                          responseservice.MerchantResponseMapper
	requestCounter                   *prometheus.CounterVec
	requestDuration                  *prometheus.HistogramVec
}

func NewMerchantStatisByApiKeyService(
	ctx context.Context,
	merchantStatisByApiKeyRepository repository.MerchantStatisticByApiKeyRepository,
	logger logger.LoggerInterface,
	mapping responseservice.MerchantResponseMapper,
) *merchantStatisByApiKeyService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_statis_by_apikey_service_requests_total",
			Help: "Total number of requests to the MerchantStatisByApiKeyService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_statis_by_apikey_service_request_duration_seconds",
			Help:    "Histogram of request durations for the MerchantStatisByApiKeyService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &merchantStatisByApiKeyService{
		ctx:                              ctx,
		trace:                            otel.Tracer("merchant-statis-by-apikey-service"),
		merchantStatisByApiKeyRepository: merchantStatisByApiKeyRepository,
		logger:                           logger,
		mapping:                          mapping,
		requestCounter:                   requestCounter,
		requestDuration:                  requestDuration,
	}
}

func (s *merchantStatisByApiKeyService) FindMonthlyPaymentMethodByApikeys(req *requests.MonthYearPaymentMethodApiKey) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse) {
	api_key := req.Apikey
	year := req.Year

	s.logger.Debug("Finding monthly payment methods by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	res, err := s.merchantStatisByApiKeyRepository.GetMonthlyPaymentMethodByApikey(req)

	if err != nil {
		s.logger.Error("Failed to find monthly payment methods by merchant", zap.Error(err), zap.String("api_key", api_key), zap.Int("year", year))

		return nil, merchant_errors.ErrFailedFindMonthlyPaymentMethodByApikeys
	}

	so := s.mapping.ToMerchantMonthlyPaymentMethods(res)

	s.logger.Debug("Successfully found monthly payment methods by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisByApiKeyService) FindYearlyPaymentMethodByApikeys(req *requests.MonthYearPaymentMethodApiKey) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse) {
	api_key := req.Apikey
	year := req.Year

	s.logger.Debug("Finding yearly payment methods by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	res, err := s.merchantStatisByApiKeyRepository.GetYearlyPaymentMethodByApikey(req)

	if err != nil {
		s.logger.Error("Failed to find yearly payment methods by merchant", zap.Error(err), zap.String("api_key", api_key), zap.Int("year", year))

		return nil, merchant_errors.ErrFailedFindYearlyPaymentMethodByApikeys
	}

	so := s.mapping.ToMerchantYearlyPaymentMethods(res)

	s.logger.Debug("Successfully found yearly payment methods by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisByApiKeyService) FindMonthlyAmountByApikeys(req *requests.MonthYearAmountApiKey) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse) {
	api_key := req.Apikey
	year := req.Year

	s.logger.Debug("Finding monthly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	res, err := s.merchantStatisByApiKeyRepository.GetMonthlyAmountByApikey(req)

	if err != nil {
		s.logger.Error("Failed to find monthly amount by merchant", zap.Error(err), zap.String("api_key", api_key), zap.Int("year", year))

		return nil, merchant_errors.ErrFailedFindMonthlyAmountByApikeys
	}

	so := s.mapping.ToMerchantMonthlyAmounts(res)

	s.logger.Debug("Successfully found monthly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisByApiKeyService) FindYearlyAmountByApikeys(req *requests.MonthYearAmountApiKey) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse) {
	api_key := req.Apikey
	year := req.Year

	s.logger.Debug("Finding yearly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	res, err := s.merchantStatisByApiKeyRepository.GetYearlyAmountByApikey(req)

	if err != nil {
		s.logger.Error("Failed to find yearly amount by merchant", zap.Error(err), zap.String("api_key", api_key), zap.Int("year", year))

		return nil, merchant_errors.ErrFailedFindYearlyAmountByApikeys
	}

	so := s.mapping.ToMerchantYearlyAmounts(res)

	s.logger.Debug("Successfully found yearly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisByApiKeyService) FindMonthlyTotalAmountByApikeys(req *requests.MonthYearTotalAmountApiKey) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse) {
	api_key := req.Apikey
	year := req.Year

	s.logger.Debug("Finding monthly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	res, err := s.merchantStatisByApiKeyRepository.GetMonthlyTotalAmountByApikey(req)

	if err != nil {
		s.logger.Error("Failed to find monthly amount by merchant", zap.Error(err), zap.String("api_key", api_key), zap.Int("year", year))

		return nil, merchant_errors.ErrFailedFindMonthlyTotalAmountByApikeys
	}

	so := s.mapping.ToMerchantMonthlyTotalAmounts(res)

	s.logger.Debug("Successfully found monthly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	return so, nil
}

func (s *merchantStatisByApiKeyService) FindYearlyTotalAmountByApikeys(req *requests.MonthYearTotalAmountApiKey) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse) {
	api_key := req.Apikey
	year := req.Year

	s.logger.Debug("Finding yearly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	res, err := s.merchantStatisByApiKeyRepository.GetYearlyTotalAmountByApikey(req)

	if err != nil {
		s.logger.Error("Failed to find yearly amount by merchant", zap.Error(err), zap.String("api_key", api_key), zap.Int("year", year))

		return nil, merchant_errors.ErrFailedFindYearlyTotalAmountByApikeys
	}

	so := s.mapping.ToMerchantYearlyTotalAmounts(res)

	s.logger.Debug("Successfully found yearly amount by merchant", zap.String("api_key", api_key), zap.Int("year", year))

	return so, nil
}
