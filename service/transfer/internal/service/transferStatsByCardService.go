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
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthTransferStatusSuccessByCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthTransferStatusSuccessByCardNumber")
	defer span.End()

	card_number := req.CardNumber
	year := req.Year
	month := req.Month

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
		attribute.String("card_number", card_number),
	)

	s.logger.Debug("Fetching monthly Transfer status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	if data := s.mencache.GetMonthTransferStatusSuccessByCard(req); data != nil {
		s.logger.Debug("Successfully fetched monthly Transfer status success from cache", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.transferStatisticByCardRepository.GetMonthTransferStatusSuccessByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleMonthTransferStatusSuccessByCardNumberError(err, "FindMonthTransferStatusSuccessByCardNumber", "transfer:by_card:month_transfer_status_success:", span, &status, zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))
	}
	so := s.mapping.ToTransferResponsesMonthStatusSuccess(records)

	s.mencache.SetMonthTransferStatusSuccessByCard(req, so)

	s.logger.Debug("Successfully fetched monthly Transfer status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	return so, nil
}

func (s *transferStatisticByCardService) FindYearlyTransferStatusSuccessByCardNumber(req *requests.YearStatusTransferCardNumber) ([]*response.TransferResponseYearStatusSuccess, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTransferStatusSuccessByCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTransferStatusSuccessByCardNumber")
	defer span.End()

	card_number := req.CardNumber
	year := req.Year

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.String("card_number", card_number),
	)

	s.logger.Debug("Fetching yearly Transfer status success", zap.Int("year", year), zap.String("card_number", card_number))

	if data := s.mencache.GetYearlyTransferStatusSuccessByCard(req); data != nil {
		s.logger.Debug("Successfully fetched yearly Transfer status success from cache", zap.Int("year", year), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.transferStatisticByCardRepository.GetYearlyTransferStatusSuccessByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleYearTransferStatusSuccessByCardNumberError(err, "FindYearlyTransferStatusSuccessByCardNumber", "FAILED_YEARLY_TRANSFER_STATUS_SUCCESS:", span, &status, zap.Int("year", year), zap.String("card_number", card_number))
	}
	so := s.mapping.ToTransferResponsesYearStatusSuccess(records)

	s.mencache.SetYearlyTransferStatusSuccessByCard(req, so)

	s.logger.Debug("Successfully fetched yearly Transfer status success", zap.Int("year", year), zap.String("card_number", card_number))

	return so, nil
}

func (s *transferStatisticByCardService) FindMonthTransferStatusFailedByCardNumber(req *requests.MonthStatusTransferCardNumber) ([]*response.TransferResponseMonthStatusFailed, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthTransferStatusFailedByCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthTransferStatusFailedByCardNumber")
	defer span.End()

	card_number := req.CardNumber
	year := req.Year
	month := req.Month

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
		attribute.String("card_number", card_number),
	)

	s.logger.Debug("Fetching monthly Transfer status Failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	if data := s.mencache.GetMonthTransferStatusFailedByCard(req); data != nil {
		s.logger.Debug("Successfully fetched monthly Transfer status Failed from cache", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.transferStatisticByCardRepository.GetMonthTransferStatusFailedByCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleMonthTransferStatusFailedByCardNumberError(err, "FindMonthTransferStatusFailedByCardNumber", "FAILED_MONTH_TRANSFER_STATUS_FAILED:", span, &status, zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))
	}
	so := s.mapping.ToTransferResponsesMonthStatusFailed(records)

	s.mencache.SetMonthTransferStatusFailedByCard(req, so)

	s.logger.Debug("Failedfully fetched monthly Transfer status Failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	return so, nil
}

func (s *transferStatisticByCardService) FindYearlyTransferStatusFailedByCardNumber(req *requests.YearStatusTransferCardNumber) ([]*response.TransferResponseYearStatusFailed, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTransferStatusFailedByCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTransferStatusFailedByCardNumber")
	defer span.End()

	card_number := req.CardNumber
	year := req.Year

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.String("card_number", card_number),
	)

	s.logger.Debug("Fetching yearly Transfer status Failed", zap.Int("year", year), zap.String("card_number", card_number))

	if data := s.mencache.GetYearlyTransferStatusFailedByCard(req); data != nil {
		s.logger.Debug("Successfully fetched yearly Transfer status Failed from cache", zap.Int("year", year), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.transferStatisticByCardRepository.GetYearlyTransferStatusFailedByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleYearTransferStatusFailedByCardNumberError(err, "FindYearlyTransferStatusFailedByCardNumber", "FAILED_YEAR_TRANSFER_STATUS_FAILED:", span, &status, zap.Int("year", year), zap.String("card_number", card_number))
	}
	so := s.mapping.ToTransferResponsesYearStatusFailed(records)

	s.mencache.SetYearlyTransferStatusFailedByCard(req, so)

	s.logger.Debug("Failedfully fetched yearly Transfer status Failed", zap.Int("year", year), zap.String("card_number", card_number))

	return so, nil
}

func (s *transferStatisticByCardService) FindMonthlyTransferAmountsBySenderCardNumber(req *requests.MonthYearCardNumber) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyTransferAmountsBySenderCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyTransferAmountsBySenderCardNumber")
	defer span.End()

	cardNumber := req.CardNumber
	year := req.Year

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.String("card_number", cardNumber),
	)

	s.logger.Debug("Fetching monthly transfer amounts by sender card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	if data := s.mencache.GetMonthlyTransferAmountsBySenderCard(req); data != nil {
		s.logger.Debug("Successfully fetched monthly transfer amounts by sender card number from cache", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	amounts, err := s.transferStatisticByCardRepository.GetMonthlyTransferAmountsBySenderCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleMonthlyTransferAmountsBySenderError(err, "FindMonthlyTransferAmountsBySenderCardNumber", "FAILED_MONTH_TRANSFER_AMOUNTS_BY_SENDER:", span, &status, zap.String("card_number", cardNumber), zap.Int("year", year))
	}

	responseAmounts := s.mapping.ToTransferResponsesMonthAmount(amounts)

	s.mencache.SetMonthlyTransferAmountsBySenderCard(req, responseAmounts)

	s.logger.Debug("Successfully fetched monthly transfer amounts by sender card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseAmounts, nil
}

func (s *transferStatisticByCardService) FindMonthlyTransferAmountsByReceiverCardNumber(req *requests.MonthYearCardNumber) ([]*response.TransferMonthAmountResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyTransferAmountsByReceiverCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyTransferAmountsByReceiverCardNumber")
	defer span.End()

	cardNumber := req.CardNumber
	year := req.Year

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.String("card_number", cardNumber),
	)

	s.logger.Debug("Fetching monthly transfer amounts by receiver card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	if data := s.mencache.GetMonthlyTransferAmountsByReceiverCard(req); data != nil {
		s.logger.Debug("Successfully fetched monthly transfer amounts by receiver card number", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	amounts, err := s.transferStatisticByCardRepository.GetMonthlyTransferAmountsByReceiverCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleMonthlyTransferAmountsByReceiverError(err, "FindMonthlyTransferAmountsByReceiverCardNumber", "FAILED_MONTH_TRANSFER_AMOUNTS_BY_RECEIVER:", span, &status, zap.String("card_number", cardNumber), zap.Int("year", year))
	}

	responseAmounts := s.mapping.ToTransferResponsesMonthAmount(amounts)

	s.mencache.SetMonthlyTransferAmountsByReceiverCard(req, responseAmounts)

	s.logger.Debug("Successfully fetched monthly transfer amounts by receiver card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseAmounts, nil
}

func (s *transferStatisticByCardService) FindYearlyTransferAmountsBySenderCardNumber(req *requests.MonthYearCardNumber) ([]*response.TransferYearAmountResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTransferAmountsBySenderCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTransferAmountsBySenderCardNumber")
	defer span.End()

	cardNumber := req.CardNumber
	year := req.Year

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.String("card_number", cardNumber),
	)

	s.logger.Debug("Fetching yearly transfer amounts by sender card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	if data := s.mencache.GetYearlyTransferAmountsBySenderCard(req); data != nil {
		s.logger.Debug("Successfully fetched yearly transfer amounts by sender card number from cache", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	amounts, err := s.transferStatisticByCardRepository.GetYearlyTransferAmountsBySenderCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleYearlyTransferAmountsBySenderError(err, "FindYearlyTransferAmountsBySenderCardNumber", "FAILED_YEAR_TRANSFER_AMOUNTS_BY_SENDER:", span, &status, zap.String("card_number", cardNumber), zap.Int("year", year))
	}

	responseAmounts := s.mapping.ToTransferResponsesYearAmount(amounts)

	s.mencache.SetYearlyTransferAmountsBySenderCard(req, responseAmounts)

	s.logger.Debug("Successfully fetched yearly transfer amounts by sender card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseAmounts, nil
}

func (s *transferStatisticByCardService) FindYearlyTransferAmountsByReceiverCardNumber(req *requests.MonthYearCardNumber) ([]*response.TransferYearAmountResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTransferAmountsByReceiverCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTransferAmountsByReceiverCardNumber")
	defer span.End()

	cardNumber := req.CardNumber
	year := req.Year

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.String("card_number", cardNumber),
	)

	s.logger.Debug("Fetching yearly transfer amounts by receiver card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	if data := s.mencache.GetYearlyTransferAmountsByReceiverCard(req); data != nil {
		s.logger.Debug("Successfully fetched yearly transfer amounts by receiver card number", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	amounts, err := s.transferStatisticByCardRepository.GetYearlyTransferAmountsByReceiverCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleYearlyTransferAmountsByReceiverError(err, "FindYearlyTransferAmountsByReceiverCardNumber", "FAILED_YEAR_TRANSFER_AMOUNTS_BY_RECEIVER", span, &status, zap.String("card_number", cardNumber), zap.Int("year", year))
	}

	responseAmounts := s.mapping.ToTransferResponsesYearAmount(amounts)

	s.mencache.SetYearlyTransferAmountsByReceiverCard(req, responseAmounts)

	s.logger.Debug("Successfully fetched yearly transfer amounts by receiver card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseAmounts, nil
}

func (s *transferStatisticByCardService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
