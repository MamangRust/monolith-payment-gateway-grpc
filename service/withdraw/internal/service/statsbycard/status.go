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

type withdrawStatsByCardStatusDeps struct {
	ErrorHandler errorhandler.WithdrawStatisticByCardErrorHandler

	Cache mencache.WithdrawStatsByCardStatusCache

	Repository repository.WithdrawStatsByCardStatusRepository

	Logger logger.LoggerInterface

	Mapper responseservice.WithdrawStatsStatusResponseMapper
}

type withdrawStatsByCardStatusService struct {
	errorhandler errorhandler.WithdrawStatisticByCardErrorHandler

	cache mencache.WithdrawStatsByCardStatusCache

	repository repository.WithdrawStatsByCardStatusRepository

	logger logger.LoggerInterface

	mapper responseservice.WithdrawStatsStatusResponseMapper

	observability observability.TraceLoggerObservability
}

func NewWithdrawStatsByCardStatusService(deps *withdrawStatsByCardStatusDeps) WithdrawStatsByCardStatusService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "withdraw_stats_bycard_status_service_request_total",
			Help: "Total number of requests to the WithdrawStatsByCardStatusService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "withdraw_stats_bycard_status_service_request_duration_seconds",
			Help:    "Histogram of request durations for the WithdrawStatsByCardStatusService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(
		otel.Tracer("withdraw-stats-bycard-status-service"), deps.Logger, requestCounter, requestDuration)

	return &withdrawStatsByCardStatusService{
		errorhandler:  deps.ErrorHandler,
		cache:         deps.Cache,
		repository:    deps.Repository,
		logger:        deps.Logger,
		mapper:        deps.Mapper,
		observability: observability,
	}
}

// FindMonthWithdrawStatusSuccessByCardNumber retrieves monthly successful withdraw statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing the card number, month, and year.
//
// Returns:
//   - []*response.WithdrawResponseMonthStatusSuccess: List of successful monthly withdraw statistics for the given card number.
//   - *response.ErrorResponse: Error information if any occurred.
func (s *withdrawStatsByCardStatusService) FindMonthWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*response.WithdrawResponseMonthStatusSuccess, *response.ErrorResponse) {
	year := req.Year
	month := req.Month
	cardNumber := req.CardNumber

	const method = "FindMonthWithdrawStatusSuccessByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("month", month), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMonthWithdrawStatusSuccessByCardNumber(ctx, req); found {
		logSuccess("Cache hit for monthly withdraw status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", cardNumber))
		return data, nil
	}

	records, err := s.repository.GetMonthWithdrawStatusSuccessByCardNumber(ctx, req)
	if err != nil {
		return s.errorhandler.HandleMonthWithdrawStatusSuccessByCardNumberError(err, method, "FAILED_GET_MONTH_WITHDRAW_STATUS_SUCCESS_BY_CARD_NUMBER", span, &status, zap.Error(err))
	}

	responseData := s.mapper.ToWithdrawResponsesMonthStatusSuccess(records)

	s.cache.SetCachedMonthWithdrawStatusSuccessByCardNumber(ctx, req, responseData)

	logSuccess("Successfully fetched monthly withdraw status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", cardNumber))

	return responseData, nil
}

// FindYearlyWithdrawStatusSuccessByCardNumber retrieves yearly successful withdraw statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing the card number and year.
//
// Returns:
//   - []*response.WithdrawResponseYearStatusSuccess: List of successful yearly withdraw statistics for the given card number.
//   - *response.ErrorResponse: Error information if any occurred.
func (s *withdrawStatsByCardStatusService) FindYearlyWithdrawStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*response.WithdrawResponseYearStatusSuccess, *response.ErrorResponse) {

	year := req.Year
	cardNumber := req.CardNumber

	const method = "FindYearlyWithdrawStatusSuccessByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedYearlyWithdrawStatusSuccessByCardNumber(ctx, req); found {
		logSuccess("Cache hit for yearly withdraw status success", zap.Int("year", year), zap.String("card_number", cardNumber))
		return data, nil
	}

	records, err := s.repository.GetYearlyWithdrawStatusSuccessByCardNumber(ctx, req)
	if err != nil {
		return s.errorhandler.HandleYearWithdrawStatusSuccessByCardNumberError(err, method, "FAILED_GET_YEARLY_WITHDRAW_STATUS_SUCCESS_BY_CARD_NUMBER", span, &status, zap.Error(err))
	}

	responseData := s.mapper.ToWithdrawResponsesYearStatusSuccess(records)

	s.cache.SetCachedYearlyWithdrawStatusSuccessByCardNumber(ctx, req, responseData)

	logSuccess("Successfully fetched yearly withdraw status success", zap.Int("year", year), zap.String("card_number", cardNumber))

	return responseData, nil
}

// FindMonthWithdrawStatusFailedByCardNumber retrieves monthly failed withdraw statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing the card number, month, and year.
//
// Returns:
//   - []*response.WithdrawResponseMonthStatusFailed: List of failed monthly withdraw statistics for the given card number.
//   - *response.ErrorResponse: Error information if any occurred.
func (s *withdrawStatsByCardStatusService) FindMonthWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusWithdrawCardNumber) ([]*response.WithdrawResponseMonthStatusFailed, *response.ErrorResponse) {
	year := req.Year
	month := req.Month
	cardNumber := req.CardNumber

	const method = "FindMonthWithdrawStatusFailedByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("month", month), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMonthWithdrawStatusFailedByCardNumber(ctx, req); found {
		logSuccess("Cache hit for monthly withdraw status failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", cardNumber))
		return data, nil
	}

	records, err := s.repository.GetMonthWithdrawStatusFailedByCardNumber(ctx, req)
	if err != nil {
		return s.errorhandler.HandleMonthWithdrawStatusFailedByCardNumberError(err, method, "FAILED_GET_MONTH_WITHDRAW_STATUS_FAILED_BY_CARD_NUMBER", span, &status, zap.Error(err))
	}

	responseData := s.mapper.ToWithdrawResponsesMonthStatusFailed(records)

	s.cache.SetCachedMonthWithdrawStatusFailedByCardNumber(ctx, req, responseData)

	logSuccess("Successfully fetched monthly withdraw status failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", cardNumber))

	return responseData, nil
}

// FindYearlyWithdrawStatusFailedByCardNumber retrieves yearly failed withdraw statistics by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing the card number and year.
//
// Returns:
//   - []*response.WithdrawResponseYearStatusFailed: List of failed yearly withdraw statistics for the given card number.
//   - *response.ErrorResponse: Error information if any occurred.
func (s *withdrawStatsByCardStatusService) FindYearlyWithdrawStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusWithdrawCardNumber) ([]*response.WithdrawResponseYearStatusFailed, *response.ErrorResponse) {
	year := req.Year
	cardNumber := req.CardNumber

	const method = "FindYearlyWithdrawStatusFailedByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedYearlyWithdrawStatusFailedByCardNumber(ctx, req); found {
		logSuccess("Cache hit for yearly withdraw status failed", zap.Int("year", year), zap.String("card_number", cardNumber))
		return data, nil
	}

	records, err := s.repository.GetYearlyWithdrawStatusFailedByCardNumber(ctx, req)
	if err != nil {
		return s.errorhandler.HandleYearWithdrawStatusFailedByCardNumberError(err, method, "FAILED_GET_YEARLY_WITHDRAW_STATUS_FAILED_BY_CARD_NUMBER", span, &status, zap.Error(err))
	}

	responseData := s.mapper.ToWithdrawResponsesYearStatusFailed(records)

	s.cache.SetCachedYearlyWithdrawStatusFailedByCardNumber(ctx, req, responseData)

	logSuccess("Successfully fetched yearly withdraw status failed", zap.Int("year", year), zap.String("card_number", cardNumber))

	return responseData, nil
}
