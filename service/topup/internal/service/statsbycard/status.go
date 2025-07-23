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

type topupStatsByCardStatusDeps struct {
	Cache cache.TopupStatsStatusByCardCache

	ErrorHandler errorhandler.TopupStatisticByCardErrorHandler

	Repository repository.TopupStatsByCardStatusRepository

	Logger logger.LoggerInterface

	Mapper responseservice.TopupStatsStatusResponseMapper
}

type topupStatsByCardStatusService struct {
	cache cache.TopupStatsStatusByCardCache

	errorHandler errorhandler.TopupStatisticByCardErrorHandler

	repository repository.TopupStatsByCardStatusRepository

	logger logger.LoggerInterface

	mapper responseservice.TopupStatsStatusResponseMapper

	observability observability.TraceLoggerObservability
}

func NewTopupStatsByCardStatusService(params *topupStatsByCardStatusDeps) TopupStatsByCardStatusService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "topup_stats_bycard_status_service_request_total",
			Help: "Total number of requests to the TopupStatsByCardStatusService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "topup_stats_bycard_status_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TopupStatsByCardStatusService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(
		otel.Tracer("topup-stats-bycard-status-service"), params.Logger, requestCounter, requestDuration)

	return &topupStatsByCardStatusService{
		cache:         params.Cache,
		errorHandler:  params.ErrorHandler,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		observability: observability,
	}
}

// FindMonthTopupStatusSuccessByCardNumber retrieves monthly successful topup statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing month, year, and card number.
//
// Returns:
//   - []*response.TopupResponseMonthStatusSuccess: List of monthly successful topup statistics.
//   - *response.ErrorResponse: Error details if retrieval fails.
func (s *topupStatsByCardStatusService) FindMonthTopupStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*response.TopupResponseMonthStatusSuccess, *response.ErrorResponse) {
	card_number := req.CardNumber
	year := req.Year
	month := req.Month

	const method = "FindMonthTopupStatusSuccessByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("month", month), attribute.String("card_number", card_number))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthTopupStatusSuccessByCardNumberCache(ctx, req); found {
		logSuccess("Successfully fetched monthly topup status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.repository.GetMonthTopupStatusSuccessByCardNumber(ctx, req)

	if err != nil {
		return s.errorHandler.HandleMonthTopupStatusSuccessByCardNumber(err, method, "FAILED_FIND_MONTHLY_METHODS_BY_CARD", span, &status, zap.Error(err))
	}
	so := s.mapper.ToTopupResponsesMonthStatusSuccess(records)

	s.cache.SetMonthTopupStatusSuccessByCardNumberCache(ctx, req, so)

	logSuccess("Successfully fetched monthly topup status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	return so, nil
}

// FindYearlyTopupStatusSuccessByCardNumber retrieves yearly successful topup statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing year and card number.
//
// Returns:
//   - []*response.TopupResponseYearStatusSuccess: List of yearly successful topup statistics.
//   - *response.ErrorResponse: Error details if retrieval fails.
func (s *topupStatsByCardStatusService) FindYearlyTopupStatusSuccessByCardNumber(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*response.TopupResponseYearStatusSuccess, *response.ErrorResponse) {
	card_number := req.CardNumber
	year := req.Year

	const method = "FindYearlyTopupStatusSuccessByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", card_number))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTopupStatusSuccessByCardNumberCache(ctx, req); found {
		logSuccess("Successfully fetched yearly topup status success", zap.Int("year", year), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.repository.GetYearlyTopupStatusSuccessByCardNumber(ctx, req)
	if err != nil {
		return s.errorHandler.HandleYearlyTopupStatusSuccessByCardNumber(err, method, "FAILED_FIND_YEARLY_METHODS_BY_CARD", span, &status, zap.Error(err))
	}
	so := s.mapper.ToTopupResponsesYearStatusSuccess(records)

	s.cache.SetYearlyTopupStatusSuccessByCardNumberCache(ctx, req, so)

	logSuccess("Successfully fetched yearly topup status success", zap.Int("year", year), zap.String("card_number", card_number))

	return so, nil
}

// FindMonthTopupStatusFailedByCardNumber retrieves monthly failed topup statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing month, year, and card number.
//
// Returns:
//   - []*response.TopupResponseMonthStatusFailed: List of monthly failed topup statistics.
//   - *response.ErrorResponse: Error details if retrieval fails.
func (s *topupStatsByCardStatusService) FindMonthTopupStatusFailedByCardNumber(ctx context.Context, req *requests.MonthTopupStatusCardNumber) ([]*response.TopupResponseMonthStatusFailed, *response.ErrorResponse) {
	card_number := req.CardNumber
	year := req.Year
	month := req.Month

	const method = "FindMonthTopupStatusFailedByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("month", month), attribute.String("card_number", card_number))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthTopupStatusFailedByCardNumberCache(ctx, req); found {
		logSuccess("Successfully fetched monthly topup status Failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.repository.GetMonthTopupStatusFailedByCardNumber(ctx, req)
	if err != nil {
		return s.errorHandler.HandleMonthTopupStatusFailedByCardNumber(err, method, "FAILED_FIND_MONTHLY_METHODS_BY_CARD", span, &status, zap.Error(err))
	}
	so := s.mapper.ToTopupResponsesMonthStatusFailed(records)

	s.cache.SetMonthTopupStatusFailedByCardNumberCache(ctx, req, so)

	logSuccess("Successfully fetched monthly topup status Failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	return so, nil
}

// FindYearlyTopupStatusFailedByCardNumber retrieves yearly failed topup statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing year and card number.
//
// Returns:
//   - []*response.TopupResponseYearStatusFailed: List of yearly failed topup statistics.
//   - *response.ErrorResponse: Error details if retrieval fails.
func (s *topupStatsByCardStatusService) FindYearlyTopupStatusFailedByCardNumber(ctx context.Context, req *requests.YearTopupStatusCardNumber) ([]*response.TopupResponseYearStatusFailed, *response.ErrorResponse) {

	card_number := req.CardNumber
	year := req.Year

	const method = "FindYearlyTopupStatusFailedByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", card_number))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTopupStatusFailedByCardNumberCache(ctx, req); found {
		logSuccess("Successfully fetched yearly topup status Failed", zap.Int("year", year), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.repository.GetYearlyTopupStatusFailedByCardNumber(ctx, req)

	if err != nil {
		return s.errorHandler.HandleYearlyTopupStatusFailedByCardNumber(err, method, "FAILED_FIND_YEARLY_AMOUNTS_BY_CARD", span, &status, zap.Int("year", year), zap.String("card_number", card_number))
	}
	so := s.mapper.ToTopupResponsesYearStatusFailed(records)

	s.cache.SetYearlyTopupStatusFailedByCardNumberCache(ctx, req, so)

	logSuccess("Successfully fetched yearly topup status Failed", zap.Int("year", year), zap.String("card_number", card_number))

	return so, nil
}
