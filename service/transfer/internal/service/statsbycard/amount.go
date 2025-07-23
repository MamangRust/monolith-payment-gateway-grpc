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

type transferStatsByCardAmountDeps struct {
	ErrorHandler errorhandler.TransferStatisticByCardErrorHandler

	Cache mencache.TransferStatsByCardAmountCache

	Sender repository.TransferStatsByCardAmountSenderRepository

	Receiver repository.TransferStatsByCardAmountReceiverRepository

	Logger logger.LoggerInterface

	Mapper responseservice.TransferAmountResponseMapper
}

type transferStatsByCardAmountService struct {
	errorHandler errorhandler.TransferStatisticByCardErrorHandler

	cache mencache.TransferStatsByCardAmountCache

	sender repository.TransferStatsByCardAmountSenderRepository

	receiver repository.TransferStatsByCardAmountReceiverRepository

	logger logger.LoggerInterface

	mapper responseservice.TransferAmountResponseMapper

	observability observability.TraceLoggerObservability
}

func NewTransferStatsByCardAmountService(params *transferStatsByCardAmountDeps) TransferStatsByCardAmountService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transfer_stats_by_card_amount_service_request_total",
			Help: "Total number of requests to the TransferStatsAmountService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transfer_stats_by_card_amount_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TransferStatsAmountService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(
		otel.Tracer("transfer-stats-bycard-amount-service"), params.Logger, requestCounter, requestDuration)

	return &transferStatsByCardAmountService{
		errorHandler:  params.ErrorHandler,
		cache:         params.Cache,
		sender:        params.Sender,
		receiver:      params.Receiver,
		logger:        params.Logger,
		mapper:        params.Mapper,
		observability: observability,
	}
}

// FindMonthlyTransferAmountsBySenderCardNumber retrieves monthly transfer amounts by sender card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request with card number, month, and year.
//
// Returns:
//   - []*response.TransferMonthAmountResponse: Monthly transfer amount stats (sent).
//   - *response.ErrorResponse: Error response if any.
func (s *transferStatsByCardAmountService) FindMonthlyTransferAmountsBySenderCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	const method = "FindMonthlyTransferAmountsBySenderCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyTransferAmountsBySenderCard(ctx, req); found {
		logSuccess("Successfully fetched monthly transfer amounts by sender card number from cache", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	amounts, err := s.sender.GetMonthlyTransferAmountsBySenderCardNumber(ctx, req)
	if err != nil {
		return s.errorHandler.HandleMonthlyTransferAmountsBySenderError(err, method, "FAILED_MONTH_TRANSFER_AMOUNTS_BY_SENDER:", span, &status, zap.Error(err))
	}

	responseAmounts := s.mapper.ToTransferResponsesMonthAmount(amounts)

	s.cache.SetMonthlyTransferAmountsBySenderCard(ctx, req, responseAmounts)

	logSuccess("Successfully fetched monthly transfer amounts by sender card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseAmounts, nil
}

// FindMonthlyTransferAmountsByReceiverCardNumber retrieves monthly transfer amounts by receiver card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request with card number, month, and year.
//
// Returns:
//   - []*response.TransferMonthAmountResponse: Monthly transfer amount stats (received).
//   - *response.ErrorResponse: Error response if any.
//   - An ErrorResponse if the retrieval or mapper fails.
func (s *transferStatsByCardAmountService) FindMonthlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse) {

	cardNumber := req.CardNumber
	year := req.Year

	const method = "FindMonthlyTransferAmountsByReceiverCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetMonthlyTransferAmountsByReceiverCard(ctx, req); found {
		logSuccess("Successfully fetched monthly transfer amounts by receiver card number from cache", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	amounts, err := s.receiver.GetMonthlyTransferAmountsByReceiverCardNumber(ctx, req)
	if err != nil {
		return s.errorHandler.HandleMonthlyTransferAmountsByReceiverError(err, method, "FAILED_MONTH_TRANSFER_AMOUNTS_BY_RECEIVER:", span, &status, zap.Error(err))
	}

	responseAmounts := s.mapper.ToTransferResponsesMonthAmount(amounts)

	s.cache.SetMonthlyTransferAmountsByReceiverCard(ctx, req, responseAmounts)

	logSuccess("Successfully fetched monthly transfer amounts by receiver card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseAmounts, nil
}

// FindYearlyTransferAmountsBySenderCardNumber retrieves yearly transfer amounts by sender card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request with card number and year.
//
// Returns:
//   - []*response.TransferYearAmountResponse: Yearly transfer amount stats (sent).
//   - *response.ErrorResponse: Error response if any.
func (s *transferStatsByCardAmountService) FindYearlyTransferAmountsBySenderCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*response.TransferYearAmountResponse, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	const method = "FindYearlyTransferAmountsBySenderCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTransferAmountsBySenderCard(ctx, req); found {
		logSuccess("Successfully fetched yearly transfer amounts by sender card number from cache", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	amounts, err := s.sender.GetYearlyTransferAmountsBySenderCardNumber(ctx, req)
	if err != nil {
		return s.errorHandler.HandleYearlyTransferAmountsBySenderError(err, method, "FAILED_YEAR_TRANSFER_AMOUNTS_BY_SENDER", span, &status, zap.Error(err))
	}

	responseAmounts := s.mapper.ToTransferResponsesYearAmount(amounts)

	s.cache.SetYearlyTransferAmountsBySenderCard(ctx, req, responseAmounts)

	logSuccess("Successfully fetched yearly transfer amounts by sender card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseAmounts, nil
}

// FindYearlyTransferAmountsByReceiverCardNumber retrieves yearly transfer amounts by receiver card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request with card number and year.
//
// Returns:
//   - []*response.TransferYearAmountResponse: Yearly transfer amount stats (received).
//   - *response.ErrorResponse: Error response if any.
func (s *transferStatsByCardAmountService) FindYearlyTransferAmountsByReceiverCardNumber(ctx context.Context, req *requests.MonthYearCardNumber) ([]*response.TransferYearAmountResponse, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	const method = "FindYearlyTransferAmountsByReceiverCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetYearlyTransferAmountsByReceiverCard(ctx, req); found {
		logSuccess("Successfully fetched yearly transfer amounts by receiver card number from cache", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	amounts, err := s.receiver.GetYearlyTransferAmountsByReceiverCardNumber(ctx, req)

	if err != nil {
		return s.errorHandler.HandleYearlyTransferAmountsByReceiverError(err, method, "FAILED_YEAR_TRANSFER_AMOUNTS_BY_RECEIVER", span, &status, zap.Error(err))
	}

	responseAmounts := s.mapper.ToTransferResponsesYearAmount(amounts)

	s.cache.SetYearlyTransferAmountsByReceiverCard(ctx, req, responseAmounts)

	logSuccess("Successfully fetched yearly transfer amounts by receiver card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseAmounts, nil
}
