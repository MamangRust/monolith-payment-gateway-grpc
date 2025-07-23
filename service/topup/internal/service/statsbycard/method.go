package topupstatsbycardservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/topup"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/errorhandler"
	cache "github.com/MamangRust/monolith-payment-gateway-topup/internal/redis/statsbycard"
	repository "github.com/MamangRust/monolith-payment-gateway-topup/internal/repository/statsbycard"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type topupStatsByCardMethodDeps struct {
	Cache cache.TopupStatsMethodByCardCache

	ErrorHandler errorhandler.TopupStatisticByCardErrorHandler

	Repository repository.TopupStatsByCardMethodRepository

	Logger logger.LoggerInterface

	Mapper responseservice.TopupStatsMethodResponseMapper
}

type topupStatsByCardMethodService struct {
	cache cache.TopupStatsMethodByCardCache

	errorHandler errorhandler.TopupStatisticByCardErrorHandler

	repository repository.TopupStatsByCardMethodRepository

	logger logger.LoggerInterface

	mapper responseservice.TopupStatsMethodResponseMapper

	observability observability.TraceLoggerObservability
}

func NewTopupStatsByCardMethodService(params *topupStatsByCardMethodDeps) TopupStatsByCardMethodService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "topup_stats_bycard_method_service_request_total",
			Help: "Total number of requests to the TopupStatsByCardMethodService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "topup_stats_bycard_method_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TopupStatsByCardMethodService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(
		otel.Tracer("topup-stats-bycard-method-service"), params.Logger, requestCounter, requestDuration)

	return &topupStatsByCardMethodService{
		cache:         params.Cache,
		errorHandler:  params.ErrorHandler,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		observability: observability,
	}
}

// FindMonthlyTopupMethodsByCardNumber retrieves monthly topup method statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing year and card number.
//
// Returns:
//   - []*response.TopupMonthMethodResponse: List of monthly topup method statistics.
//   - *response.ErrorResponse: Error details if retrieval fails.
func (s *topupStatsByCardMethodService) FindMonthlyTopupMethodsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*response.TopupMonthMethodResponse, *response.ErrorResponse) {

	year := req.Year
	cardNumber := req.CardNumber

	const method = "FindMonthlyTopupMethodsByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyTopupMethodsByCardNumberCache(ctx, req); found {
		logSuccess("Successfully fetched monthly topup methods by card number", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	records, err := s.repository.GetMonthlyTopupMethodsByCardNumber(ctx, req)
	if err != nil {
		return s.errorHandler.HandleMonthlyTopupMethodsByCardNumber(err, method, "FAILED_FIND_MONTHLY_METHODS_BY_CARD", span, &status, zap.String("card_number", cardNumber), zap.Int("year", year))
	}

	responses := s.mapper.ToTopupMonthlyMethodResponses(records)

	s.cache.SetMonthlyTopupMethodsByCardNumberCache(ctx, req, responses)

	logSuccess("Successfully fetched monthly topup methods by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

// FindYearlyTopupMethodsByCardNumber retrieves yearly topup method statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing year and card number.
//
// Returns:
//   - []*response.TopupYearlyMethodResponse: List of yearly topup method statistics.
//   - *response.ErrorResponse: Error details if retrieval fails.
func (s *topupStatsByCardMethodService) FindYearlyTopupMethodsByCardNumber(ctx context.Context, req *requests.YearMonthMethod) ([]*response.TopupYearlyMethodResponse, *response.ErrorResponse) {

	year := req.Year
	cardNumber := req.CardNumber

	const method = "FindYearlyTopupMethodsByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTopupMethodsByCardNumberCache(ctx, req); found {
		logSuccess("Successfully fetched yearly topup methods by card number", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	records, err := s.repository.GetYearlyTopupMethodsByCardNumber(ctx, req)
	if err != nil {
		return s.errorHandler.HandleYearlyTopupMethodsByCardNumber(err, method, "FAILED_FIND_YEARLY_TOPUP_METHODS_BY_CARD", span, &status, zap.Error(err))
	}

	responses := s.mapper.ToTopupYearlyMethodResponses(records)

	s.cache.SetYearlyTopupMethodsByCardNumberCache(ctx, req, responses)

	logSuccess("Successfully fetched yearly topup methods by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}
