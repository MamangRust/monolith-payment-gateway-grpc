package topupstatsbycardservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/topup"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-topup/internal/redis/statsbycard"
	repository "github.com/MamangRust/monolith-payment-gateway-topup/internal/repository/statsbycard"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type topupStatsByCardAmountDeps struct {
	Cache mencache.TopupStatsAmountByCardCache

	ErrorHandler errorhandler.TopupStatisticByCardErrorHandler

	Repository repository.TopupStatsByCardAmountRepository

	Logger logger.LoggerInterface

	Mapper responseservice.TopupStatsAmountResponseMapper
}

type topupStatsByCardAmountService struct {
	cache mencache.TopupStatsAmountByCardCache

	errorHandler errorhandler.TopupStatisticByCardErrorHandler

	repository repository.TopupStatsByCardAmountRepository

	logger logger.LoggerInterface

	mapper responseservice.TopupStatsAmountResponseMapper

	observability observability.TraceLoggerObservability
}

func NewTopupStatsByCardAmountService(params *topupStatsByCardAmountDeps) TopupStatsByCardAmountService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "topup_stats_bycard_amount_service_request_total",
			Help: "Total number of requests to the TopupStatsByCardAmountService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "topup_stats_bycard_amount_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TopupStatsByCardAmountService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(
		otel.Tracer("topup-stats-bycard-amount-service"), params.Logger, requestCounter, requestDuration)

	return &topupStatsByCardAmountService{
		cache:         params.Cache,
		errorHandler:  params.ErrorHandler,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		observability: observability,
	}
}

// FindMonthlyTopupAmountsByCardNumber retrieves monthly topup amount statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing year and card number.
//
// Returns:
//   - []*response.TopupMonthAmountResponse: List of monthly topup amount statistics.
//   - *response.ErrorResponse: Error details if retrieval fails.
func (s *topupStatsByCardAmountService) FindMonthlyTopupAmountsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*response.TopupMonthAmountResponse, *response.ErrorResponse) {
	year := req.Year
	cardNumber := req.CardNumber

	const method = "FindMonthlyTopupAmountsByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyTopupAmountsByCardNumberCache(ctx, req); found {
		logSuccess("Successfully fetched monthly topup amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	records, err := s.repository.GetMonthlyTopupAmountsByCardNumber(ctx, req)
	if err != nil {
		return s.errorHandler.HandleMonthlyTopupAmountsByCardNumber(err, method, "FAILED_FIND_MONTHLY_AMOUNTS_BY_CARD", span, &status, zap.Error(err))
	}

	responses := s.mapper.ToTopupMonthlyAmountResponses(records)

	s.cache.SetMonthlyTopupAmountsByCardNumberCache(ctx, req, responses)

	logSuccess("Successfully fetched monthly topup amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

// FindYearlyTopupAmountsByCardNumber retrieves yearly topup amount statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing year and card number.
//
// Returns:
//   - []*response.TopupYearlyAmountResponse: List of yearly topup amount statistics.
//   - *response.ErrorResponse: Error details if retrieval fails.
func (s *topupStatsByCardAmountService) FindYearlyTopupAmountsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*response.TopupYearlyAmountResponse, *response.ErrorResponse) {
	year := req.Year
	cardNumber := req.CardNumber

	const method = "FindMonthTopupStatusSuccessByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTopupAmountsByCardNumberCache(ctx, req); found {
		logSuccess("Successfully fetched yearly topup amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	records, err := s.repository.GetYearlyTopupAmountsByCardNumber(ctx, req)
	if err != nil {
		return s.errorHandler.HandleYearlyTopupAmountsByCardNumber(err, method, "FAILED_FIND_YEARLY_TOPUP_AMOUNTS_BY_CARD", span, &status, zap.Error(err))
	}

	responses := s.mapper.ToTopupYearlyAmountResponses(records)

	s.cache.SetYearlyTopupAmountsByCardNumberCache(ctx, req, responses)

	logSuccess("Successfully fetched yearly topup amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}
