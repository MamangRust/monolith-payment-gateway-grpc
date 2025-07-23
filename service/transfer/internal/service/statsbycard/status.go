package transferstatsbycardservice

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/transfer"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-transfer/internal/redis/statsbycard"
	repository "github.com/MamangRust/monolith-payment-gateway-transfer/internal/repository/statsbycard"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type transferStatsByCardStatusDeps struct {
	ErrorHandler errorhandler.TransferStatisticByCardErrorHandler

	Cache mencache.TransferStatsByCardStatusCache

	Repository repository.TransferStatsByCardStatusRepository

	Logger logger.LoggerInterface

	Mapper responseservice.TransferStatsStatusResponseMapper
}

type transferStatsByCardStatusService struct {
	errorHandler errorhandler.TransferStatisticByCardErrorHandler

	cache mencache.TransferStatsByCardStatusCache

	repository repository.TransferStatsByCardStatusRepository

	logger logger.LoggerInterface

	mapper responseservice.TransferStatsStatusResponseMapper

	observability observability.TraceLoggerObservability
}

func NewTransferStatsByCardStatusService(params *transferStatsByCardStatusDeps) TransferStatsByCardStatusService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transfer_stats_bycard_status_service_request_total",
			Help: "Total number of requests to the TransferStatsStatusService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transfer_stats_bycard_status_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TransferStatsStatusService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(
		otel.Tracer("transfer-stats-bycard-status-service"), params.Logger, requestCounter, requestDuration)

	return &transferStatsByCardStatusService{
		errorHandler:  params.ErrorHandler,
		cache:         params.Cache,
		repository:    params.Repository,
		logger:        params.Logger,
		mapper:        params.Mapper,
		observability: observability,
	}
}

// FindMonthTransferStatusSuccessByCardNumber retrieves monthly successful transfer stats by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing card number, month, and year.
//
// Returns:
//   - []*response.TransferResponseMonthStatusSuccess: Monthly success transfer stats.
//   - *response.ErrorResponse: Error response if any.
func (s *transferStatsByCardStatusService) FindMonthTransferStatusSuccessByCardNumber(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*response.TransferResponseMonthStatusSuccess, *response.ErrorResponse) {
	card_number := req.CardNumber
	year := req.Year
	month := req.Month

	const method = "FindAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("month", month), attribute.String("card_number", card_number))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthTransferStatusSuccessByCard(ctx, req); found {
		logSuccess("Successfully fetched monthly Transfer status success from cache", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.repository.GetMonthTransferStatusSuccessByCardNumber(ctx, req)
	if err != nil {
		return s.errorHandler.HandleMonthTransferStatusSuccessByCardNumberError(err, method, "FAILED_MONTH_TRANSFER_STATUS_SUCCESS:", span, &status, zap.Error(err))
	}
	so := s.mapper.ToTransferResponsesMonthStatusSuccess(records)

	s.cache.SetMonthTransferStatusSuccessByCard(ctx, req, so)

	logSuccess("Successfully fetched monthly Transfer status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	return so, nil
}

// FindYearlyTransferStatusSuccessByCardNumber retrieves yearly successful transfer stats by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing card number and year.
//
// Returns:
//   - []*response.TransferResponseYearStatusSuccess: Yearly success transfer stats.
//   - *response.ErrorResponse: Error response if any.
func (s *transferStatsByCardStatusService) FindYearlyTransferStatusSuccessByCardNumber(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*response.TransferResponseYearStatusSuccess, *response.ErrorResponse) {

	card_number := req.CardNumber
	year := req.Year

	const method = "FindYearlyTransferStatusSuccessByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", card_number))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTransferStatusSuccessByCard(ctx, req); found {
		logSuccess("Successfully fetched yearly Transfer status success from cache", zap.Int("year", year), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.repository.GetYearlyTransferStatusSuccessByCardNumber(ctx, req)
	if err != nil {
		return s.errorHandler.HandleYearTransferStatusSuccessByCardNumberError(err, method, "FAILED_YEARLY_TRANSFER_STATUS_SUCCESS:", span, &status, zap.Error(err))
	}
	so := s.mapper.ToTransferResponsesYearStatusSuccess(records)

	s.cache.SetYearlyTransferStatusSuccessByCard(ctx, req, so)

	logSuccess("Successfully fetched yearly Transfer status success", zap.Int("year", year), zap.String("card_number", card_number))

	return so, nil
}

// FindMonthTransferStatusFailedByCardNumber retrieves monthly failed transfer stats by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing card number, month, and year.
//
// Returns:
//   - []*response.TransferResponseMonthStatusFailed: Monthly failed transfer stats.
//   - *response.ErrorResponse: Error response if any.
func (s *transferStatsByCardStatusService) FindMonthTransferStatusFailedByCardNumber(ctx context.Context, req *requests.MonthStatusTransferCardNumber) ([]*response.TransferResponseMonthStatusFailed, *response.ErrorResponse) {
	card_number := req.CardNumber
	year := req.Year
	month := req.Month

	const method = "FindMonthTransferStatusFailedByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("month", month), attribute.String("card_number", card_number))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthTransferStatusFailedByCard(ctx, req); found {
		logSuccess("Successfully fetched monthly Transfer status Failed from cache", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.repository.GetMonthTransferStatusFailedByCardNumber(ctx, req)

	if err != nil {
		return s.errorHandler.HandleMonthTransferStatusFailedByCardNumberError(err, method, "FAILED_MONTH_TRANSFER_STATUS_FAILED:", span, &status, zap.Error(err))
	}
	so := s.mapper.ToTransferResponsesMonthStatusFailed(records)

	s.cache.SetMonthTransferStatusFailedByCard(ctx, req, so)

	logSuccess("Successfully fetched monthly Transfer status Failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	return so, nil
}

// FindYearlyTransferStatusFailedByCardNumber retrieves yearly failed transfer stats by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing card number and year.
//
// Returns:
//   - []*response.TransferResponseYearStatusFailed: Yearly failed transfer stats.
//   - *response.ErrorResponse: Error response if any.
func (s *transferStatsByCardStatusService) FindYearlyTransferStatusFailedByCardNumber(ctx context.Context, req *requests.YearStatusTransferCardNumber) ([]*response.TransferResponseYearStatusFailed, *response.ErrorResponse) {
	card_number := req.CardNumber
	year := req.Year

	const method = "FindYearlyTransferStatusFailedByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", card_number))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTransferStatusFailedByCard(ctx, req); found {
		logSuccess("Successfully fetched yearly Transfer status Failed from cache", zap.Int("year", year), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.repository.GetYearlyTransferStatusFailedByCardNumber(ctx, req)
	if err != nil {
		return s.errorHandler.HandleYearTransferStatusFailedByCardNumberError(err, method, "FAILED_YEAR_TRANSFER_STATUS_FAILED:", span, &status, zap.Error(err))
	}
	so := s.mapper.ToTransferResponsesYearStatusFailed(records)

	s.cache.SetYearlyTransferStatusFailedByCard(ctx, req, so)

	logSuccess("Successfully fetched yearly Transfer status Failed", zap.Int("year", year), zap.String("card_number", card_number))

	return so, nil
}
