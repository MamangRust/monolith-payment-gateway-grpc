package withdrawstatsbycardservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/withdraw"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-withdraw/internal/redis/statsbycard"
	repository "github.com/MamangRust/monolith-payment-gateway-withdraw/internal/repository/statsbycard"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type withdrawStatsByCardAmountDeps struct {
	ErrorHandler errorhandler.WithdrawStatisticByCardErrorHandler

	Cache mencache.WithdrawStatsByCardAmountCache

	Repository repository.WithdrawStatsByCardAmountRepository

	Logger logger.LoggerInterface

	Mapper responseservice.WithdrawStatsAmountResponseMapper
}

type withdrawStatsByCardAmountService struct {
	errorhandler errorhandler.WithdrawStatisticByCardErrorHandler

	cache mencache.WithdrawStatsByCardAmountCache

	repository repository.WithdrawStatsByCardAmountRepository

	logger logger.LoggerInterface

	mapper responseservice.WithdrawStatsAmountResponseMapper

	observability observability.TraceLoggerObservability
}

func NewWithdrawStatsByCardAmountService(deps *withdrawStatsByCardAmountDeps) WithdrawStatsByCardAmountService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "withdraw_stats_bycard_amount_service_request_total",
			Help: "Total number of requests to the WithdrawStatsByCardAmountService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "withdraw_stats_bycard_amount_service_request_duration_seconds",
			Help:    "Histogram of request durations for the WithdrawStatsByCardAmountService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(
		otel.Tracer("withdraw-stats-bycard-amount-service"), deps.Logger, requestCounter, requestDuration)

	return &withdrawStatsByCardAmountService{
		errorhandler:  deps.ErrorHandler,
		cache:         deps.Cache,
		repository:    deps.Repository,
		logger:        deps.Logger,
		mapper:        deps.Mapper,
		observability: observability,
	}
}

// FindMonthlyWithdrawsByCardNumber retrieves total monthly withdraw amounts by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing the card number, month, and year.
//
// Returns:
//   - []*response.WithdrawMonthlyAmountResponse: List of monthly withdraw amounts for the given card number.
//   - *response.ErrorResponse: Error information if any occurred.
func (s *withdrawStatsByCardAmountService) FindMonthlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*response.WithdrawMonthlyAmountResponse, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	const method = "FindMonthlyWithdrawsByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMonthlyWithdrawsByCardNumber(ctx, req); found {
		logSuccess("Cache hit for monthly withdraws by card number", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	withdraws, err := s.repository.GetMonthlyWithdrawsByCardNumber(ctx, req)
	if err != nil {
		return s.errorhandler.HandleMonthlyWithdrawsAmountByCardNumberError(err, method, "FAILED_GET_MONTHLY_WITHDRAW_BY_CARD_NUMBER", span, &status, zap.Error(err))
	}

	responseWithdraws := s.mapper.ToWithdrawsAmountMonthlyResponses(withdraws)

	s.cache.SetCachedMonthlyWithdrawsByCardNumber(ctx, req, responseWithdraws)

	logSuccess("Successfully fetched monthly withdraws by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseWithdraws, nil
}

// FindYearlyWithdrawsByCardNumber retrieves total yearly withdraw amounts by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing the card number and year.
//
// Returns:
//   - []*response.WithdrawYearlyAmountResponse: List of yearly withdraw amounts for the given card number.
//   - *response.ErrorResponse: Error information if any occurred.
func (s *withdrawStatsByCardAmountService) FindYearlyWithdrawsByCardNumber(ctx context.Context, req *requests.YearMonthCardNumber) ([]*response.WithdrawYearlyAmountResponse, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	const method = "FindYearlyWithdrawsByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedYearlyWithdrawsByCardNumber(ctx, req); found {
		logSuccess("Cache hit for yearly withdraws by card number", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	withdraws, err := s.repository.GetYearlyWithdrawsByCardNumber(ctx, req)

	if err != nil {
		return s.errorhandler.HandleYearlyWithdrawsAmountByCardNumberError(err, method, "FAILED_GET_YEARLY_WITHDRAW_BY_CARD_NUMBER", span, &status, zap.Error(err))
	}

	responseWithdraws := s.mapper.ToWithdrawsAmountYearlyResponses(withdraws)

	s.cache.SetCachedYearlyWithdrawsByCardNumber(ctx, req, responseWithdraws)

	logSuccess("Successfully fetched yearly withdraws by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseWithdraws, nil
}
