package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-transfer/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transferStatisticByCardService struct {
	ctx                               context.Context
	mencache                          mencache.TransferStatisticByCardCache
	errorhandler                      errorhandler.TransferStatisticByCardErrorHandler
	trace                             trace.Tracer
	cardRepository                    repository.CardRepository
	saldoRepository                   repository.SaldoRepository
	transferStatisticByCardRepository repository.TransferStatisticByCardRepository
	logger                            logger.LoggerInterface
	mapping                           responseservice.TransferResponseMapper
	requestCounter                    *prometheus.CounterVec
	requestDuration                   *prometheus.HistogramVec
}

func NewTransferStatisticByCardService(
	ctx context.Context,
	mencache mencache.TransferStatisticByCardCache,
	errorhandler errorhandler.TransferStatisticByCardErrorHandler,
	cardRepository repository.CardRepository,
	transferStatisticByCardRepository repository.TransferStatisticByCardRepository,
	saldoRepository repository.SaldoRepository, logger logger.LoggerInterface, mapping responseservice.TransferResponseMapper) *transferStatisticByCardService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transfer_statistic_by_card_service_request_total",
			Help: "Total number of requests to the TransferStatisticByCardService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transfer_statistic_by_card_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TransferStatisticByCardService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &transferStatisticByCardService{
		ctx:                               ctx,
		errorhandler:                      errorhandler,
		mencache:                          mencache,
		trace:                             otel.Tracer("transfer-statistic-by-card-service"),
		cardRepository:                    cardRepository,
		transferStatisticByCardRepository: transferStatisticByCardRepository,
		saldoRepository:                   saldoRepository,
		logger:                            logger,
		mapping:                           mapping,
		requestCounter:                    requestCounter,
		requestDuration:                   requestDuration,
	}
}

func (s *transferStatisticByCardService) FindMonthTransferStatusSuccessByCardNumber(req *requests.MonthStatusTransferCardNumber) ([]*response.TransferResponseMonthStatusSuccess, *response.ErrorResponse) {
	card_number := req.CardNumber
	year := req.Year
	month := req.Month

	const method = "FindAll"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("month", month), attribute.String("card_number", card_number))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthTransferStatusSuccessByCard(req); found {
		logSuccess("Successfully fetched monthly Transfer status success from cache", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.transferStatisticByCardRepository.GetMonthTransferStatusSuccessByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleMonthTransferStatusSuccessByCardNumberError(err, method, "FAILED_MONTH_TRANSFER_STATUS_SUCCESS:", span, &status, zap.Error(err))
	}
	so := s.mapping.ToTransferResponsesMonthStatusSuccess(records)

	s.mencache.SetMonthTransferStatusSuccessByCard(req, so)

	logSuccess("Successfully fetched monthly Transfer status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	return so, nil
}

func (s *transferStatisticByCardService) FindYearlyTransferStatusSuccessByCardNumber(req *requests.YearStatusTransferCardNumber) ([]*response.TransferResponseYearStatusSuccess, *response.ErrorResponse) {

	card_number := req.CardNumber
	year := req.Year

	const method = "FindYearlyTransferStatusSuccessByCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", card_number))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyTransferStatusSuccessByCard(req); found {
		logSuccess("Successfully fetched yearly Transfer status success from cache", zap.Int("year", year), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.transferStatisticByCardRepository.GetYearlyTransferStatusSuccessByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleYearTransferStatusSuccessByCardNumberError(err, method, "FAILED_YEARLY_TRANSFER_STATUS_SUCCESS:", span, &status, zap.Error(err))
	}
	so := s.mapping.ToTransferResponsesYearStatusSuccess(records)

	s.mencache.SetYearlyTransferStatusSuccessByCard(req, so)

	logSuccess("Successfully fetched yearly Transfer status success", zap.Int("year", year), zap.String("card_number", card_number))

	return so, nil
}

func (s *transferStatisticByCardService) FindMonthTransferStatusFailedByCardNumber(req *requests.MonthStatusTransferCardNumber) ([]*response.TransferResponseMonthStatusFailed, *response.ErrorResponse) {
	card_number := req.CardNumber
	year := req.Year
	month := req.Month

	const method = "FindMonthTransferStatusFailedByCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("month", month), attribute.String("card_number", card_number))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthTransferStatusFailedByCard(req); found {
		logSuccess("Successfully fetched monthly Transfer status Failed from cache", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.transferStatisticByCardRepository.GetMonthTransferStatusFailedByCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleMonthTransferStatusFailedByCardNumberError(err, method, "FAILED_MONTH_TRANSFER_STATUS_FAILED:", span, &status, zap.Error(err))
	}
	so := s.mapping.ToTransferResponsesMonthStatusFailed(records)

	s.mencache.SetMonthTransferStatusFailedByCard(req, so)

	logSuccess("Successfully fetched monthly Transfer status Failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	return so, nil
}

func (s *transferStatisticByCardService) FindYearlyTransferStatusFailedByCardNumber(req *requests.YearStatusTransferCardNumber) ([]*response.TransferResponseYearStatusFailed, *response.ErrorResponse) {
	card_number := req.CardNumber
	year := req.Year

	const method = "FindYearlyTransferStatusFailedByCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", card_number))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyTransferStatusFailedByCard(req); found {
		logSuccess("Successfully fetched yearly Transfer status Failed from cache", zap.Int("year", year), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.transferStatisticByCardRepository.GetYearlyTransferStatusFailedByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleYearTransferStatusFailedByCardNumberError(err, method, "FAILED_YEAR_TRANSFER_STATUS_FAILED:", span, &status, zap.Error(err))
	}
	so := s.mapping.ToTransferResponsesYearStatusFailed(records)

	s.mencache.SetYearlyTransferStatusFailedByCard(req, so)

	logSuccess("Successfully fetched yearly Transfer status Failed", zap.Int("year", year), zap.String("card_number", card_number))

	return so, nil
}

func (s *transferStatisticByCardService) FindMonthlyTransferAmountsBySenderCardNumber(req *requests.MonthYearCardNumber) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	const method = "FindMonthlyTransferAmountsBySenderCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyTransferAmountsBySenderCard(req); found {
		logSuccess("Successfully fetched monthly transfer amounts by sender card number from cache", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	amounts, err := s.transferStatisticByCardRepository.GetMonthlyTransferAmountsBySenderCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleMonthlyTransferAmountsBySenderError(err, method, "FAILED_MONTH_TRANSFER_AMOUNTS_BY_SENDER:", span, &status, zap.Error(err))
	}

	responseAmounts := s.mapping.ToTransferResponsesMonthAmount(amounts)

	s.mencache.SetMonthlyTransferAmountsBySenderCard(req, responseAmounts)

	logSuccess("Successfully fetched monthly transfer amounts by sender card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseAmounts, nil
}

func (s *transferStatisticByCardService) FindMonthlyTransferAmountsByReceiverCardNumber(req *requests.MonthYearCardNumber) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse) {

	cardNumber := req.CardNumber
	year := req.Year

	const method = "FindMonthlyTransferAmountsByReceiverCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyTransferAmountsByReceiverCard(req); found {
		logSuccess("Successfully fetched monthly transfer amounts by receiver card number from cache", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	amounts, err := s.transferStatisticByCardRepository.GetMonthlyTransferAmountsByReceiverCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleMonthlyTransferAmountsByReceiverError(err, method, "FAILED_MONTH_TRANSFER_AMOUNTS_BY_RECEIVER:", span, &status, zap.Error(err))
	}

	responseAmounts := s.mapping.ToTransferResponsesMonthAmount(amounts)

	s.mencache.SetMonthlyTransferAmountsByReceiverCard(req, responseAmounts)

	logSuccess("Successfully fetched monthly transfer amounts by receiver card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseAmounts, nil
}

func (s *transferStatisticByCardService) FindYearlyTransferAmountsBySenderCardNumber(req *requests.MonthYearCardNumber) ([]*response.TransferYearAmountResponse, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	const method = "FindYearlyTransferAmountsBySenderCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyTransferAmountsBySenderCard(req); found {
		logSuccess("Successfully fetched yearly transfer amounts by sender card number from cache", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	amounts, err := s.transferStatisticByCardRepository.GetYearlyTransferAmountsBySenderCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleYearlyTransferAmountsBySenderError(err, method, "FAILED_YEAR_TRANSFER_AMOUNTS_BY_SENDER", span, &status, zap.Error(err))
	}

	responseAmounts := s.mapping.ToTransferResponsesYearAmount(amounts)

	s.mencache.SetYearlyTransferAmountsBySenderCard(req, responseAmounts)

	logSuccess("Successfully fetched yearly transfer amounts by sender card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseAmounts, nil
}

func (s *transferStatisticByCardService) FindYearlyTransferAmountsByReceiverCardNumber(req *requests.MonthYearCardNumber) ([]*response.TransferYearAmountResponse, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	const method = "FindYearlyTransferAmountsByReceiverCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyTransferAmountsByReceiverCard(req); found {
		logSuccess("Successfully fetched yearly transfer amounts by receiver card number from cache", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	amounts, err := s.transferStatisticByCardRepository.GetYearlyTransferAmountsByReceiverCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleYearlyTransferAmountsByReceiverError(err, method, "FAILED_YEAR_TRANSFER_AMOUNTS_BY_RECEIVER", span, &status, zap.Error(err))
	}

	responseAmounts := s.mapping.ToTransferResponsesYearAmount(amounts)

	s.mencache.SetYearlyTransferAmountsByReceiverCard(req, responseAmounts)

	logSuccess("Successfully fetched yearly transfer amounts by receiver card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseAmounts, nil
}

func (s *transferStatisticByCardService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
	trace.Span,
	func(string),
	string,
	func(string, ...zap.Field),
) {
	start := time.Now()
	status := "success"

	_, span := s.trace.Start(s.ctx, method)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.AddEvent("Start: " + method)

	s.logger.Info("Start: " + method)

	end := func(status string) {
		s.recordMetrics(method, status, start)
		code := codes.Ok
		if status != "success" {
			code = codes.Error
		}
		span.SetStatus(code, status)
		span.End()
	}

	logSuccess := func(msg string, fields ...zap.Field) {
		span.AddEvent(msg)
		s.logger.Info(msg, fields...)
	}

	return span, end, status, logSuccess
}

func (s *transferStatisticByCardService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
