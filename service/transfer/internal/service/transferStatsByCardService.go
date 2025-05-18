package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	traceunic "github.com/MamangRust/monolith-payment-gateway-pkg/trace_unic"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transferStatisticByCardService struct {
	ctx                               context.Context
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
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &transferStatisticByCardService{
		ctx:                               ctx,
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

	records, err := s.transferStatisticByCardRepository.GetMonthTransferStatusSuccessByCardNumber(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTH_TRANSFER_SUCCESS_BY_CARD")

		s.logger.Error("Failed to fetch monthly Transfer status success", zap.Error(err), zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)

		status = "failed_find_month_transfer_success_by_card"

		return nil, transfer_errors.ErrFailedFindMonthTransferStatusSuccess
	}

	s.logger.Debug("Successfully fetched monthly Transfer status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	so := s.mapping.ToTransferResponsesMonthStatusSuccess(records)

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

	records, err := s.transferStatisticByCardRepository.GetYearlyTransferStatusSuccessByCardNumber(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_TRANSFER_SUCCESS_BY_CARD")

		s.logger.Error("Failed to fetch yearly Transfer status success", zap.Error(err), zap.Int("year", year), zap.String("card_number", card_number))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)

		status = "failed_find_yearly_transfer_success_by_card"

		return nil, transfer_errors.ErrFailedFindYearTransferStatusSuccessByCard
	}

	s.logger.Debug("Successfully fetched yearly Transfer status success", zap.Int("year", year), zap.String("card_number", card_number))

	so := s.mapping.ToTransferResponsesYearStatusSuccess(records)

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

	records, err := s.transferStatisticByCardRepository.GetMonthTransferStatusFailedByCardNumber(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTH_TRANSFER_STATUS_FAILED")

		s.logger.Error("Failed to fetch monthly Transfer status Failed", zap.Error(err), zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)

		status = "failed_find_month_transfer_status_failed"

		return nil, transfer_errors.ErrFailedFindMonthTransferStatusFailed
	}

	s.logger.Debug("Failedfully fetched monthly Transfer status Failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	so := s.mapping.ToTransferResponsesMonthStatusFailed(records)

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

	records, err := s.transferStatisticByCardRepository.GetYearlyTransferStatusFailedByCardNumber(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEAR_TRANSFER_STATUS_FAILED")

		s.logger.Error("Failed to fetch yearly Transfer status Failed", zap.Error(err), zap.Int("year", year), zap.String("card_number", card_number))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)

		status = "failed_find_year_transfer_status_failed"

		return nil, transfer_errors.ErrFailedFindYearTransferStatusFailedByCard
	}

	s.logger.Debug("Failedfully fetched yearly Transfer status Failed", zap.Int("year", year), zap.String("card_number", card_number))

	so := s.mapping.ToTransferResponsesYearStatusFailed(records)

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

	amounts, err := s.transferStatisticByCardRepository.GetMonthlyTransferAmountsBySenderCardNumber(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_TRANSFER_AMOUNTS_BY_SENDER_CARD_NUMBER")

		s.logger.Error("failed to find monthly transfer amounts by sender card number", zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)

		status = "failed_find_monthly_transfer_amounts_by_sender_card_number"

		return nil, transfer_errors.ErrFailedFindMonthlyTransferAmountsBySenderCard
	}

	responseAmounts := s.mapping.ToTransferResponsesMonthAmount(amounts)

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

	amounts, err := s.transferStatisticByCardRepository.GetMonthlyTransferAmountsByReceiverCardNumber(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_TRANSFER_AMOUNTS_BY_RECEIVER_CARD_NUMBER")

		s.logger.Error("failed to find monthly transfer amounts by receiver card number", zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)

		status = "failed_find_monthly_transfer_amounts_by_receiver_card_number"

		return nil, transfer_errors.ErrFailedFindMonthlyTransferAmountsByReceiverCard
	}

	responseAmounts := s.mapping.ToTransferResponsesMonthAmount(amounts)

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

	amounts, err := s.transferStatisticByCardRepository.GetYearlyTransferAmountsBySenderCardNumber(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_TRANSFER_AMOUNTS_BY_SENDER_CARD_NUMBER")

		s.logger.Error("failed to find yearly transfer amounts by sender card number", zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)

		status = "failed_find_yearly_transfer_amounts_by_sender_card_number"

		return nil, transfer_errors.ErrFailedFindYearlyTransferAmountsBySenderCard
	}

	responseAmounts := s.mapping.ToTransferResponsesYearAmount(amounts)

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

	amounts, err := s.transferStatisticByCardRepository.GetYearlyTransferAmountsByReceiverCardNumber(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_TRANSFER_AMOUNTS_BY_RECEIVER_CARD_NUMBER")

		s.logger.Error("failed to find yearly transfer amounts by receiver card number", zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)

		status = "failed_find_yearly_transfer_amounts_by_receiver_card_number"

		return nil, transfer_errors.ErrFailedFindYearlyTransferAmountsByReceiverCard
	}

	responseAmounts := s.mapping.ToTransferResponsesYearAmount(amounts)

	s.logger.Debug("Successfully fetched yearly transfer amounts by receiver card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseAmounts, nil
}

func (s *transferStatisticByCardService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
