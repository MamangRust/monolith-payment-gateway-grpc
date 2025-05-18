package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	traceunic "github.com/MamangRust/monolith-payment-gateway-pkg/trace_unic"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transactionStatisticByCardService struct {
	ctx                                  context.Context
	trace                                trace.Tracer
	transactionStatisticByCardRepository repository.TransactionStatisticByCardRepository
	logger                               logger.LoggerInterface
	mapping                              responseservice.TransactionResponseMapper
	requestCounter                       *prometheus.CounterVec
	requestDuration                      *prometheus.HistogramVec
}

func NewTransactionStatisticByCardService(
	ctx context.Context,
	transactionStatisticByCardRepository repository.TransactionStatisticByCardRepository,
	logger logger.LoggerInterface,
	mapping responseservice.TransactionResponseMapper,
) *transactionStatisticByCardService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transaction_statistic_by_card_service_request_total",
			Help: "Total number of requests to the TransactionStatisticByCardService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transaction_statistic_by_card_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TransactionStatisticByCardService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &transactionStatisticByCardService{
		ctx:                                  ctx,
		trace:                                otel.Tracer("transaction-statistic-by-card-service"),
		transactionStatisticByCardRepository: transactionStatisticByCardRepository,
		logger:                               logger,
		mapping:                              mapping,
		requestCounter:                       requestCounter,
		requestDuration:                      requestDuration,
	}
}

func (s *transactionStatisticByCardService) FindMonthTransactionStatusSuccessByCardNumber(req *requests.MonthStatusTransactionCardNumber) ([]*response.TransactionResponseMonthStatusSuccess, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthTransactionStatusSuccessByCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthTransactionStatusSuccessByCardNumber")
	defer span.End()

	card_number := req.CardNumber
	year := req.Year
	month := req.Month

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
		attribute.String("card_number", card_number),
	)

	s.logger.Debug("Fetching monthly Transaction status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	records, err := s.transactionStatisticByCardRepository.GetMonthTransactionStatusSuccessByCardNumber(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTH_TRANSACTION_SUCCESS_BY_CARD")

		s.logger.Error("failed to fetch monthly Transaction status success", zap.Error(err), zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch monthly Transaction status success")

		status = "failed_find_month_transaction_success_by_card"

		return nil, transaction_errors.ErrFailedFindMonthTransactionSuccessByCard
	}

	s.logger.Debug("Successfully fetched monthly Transaction status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	so := s.mapping.ToTransactionResponsesMonthStatusSuccess(records)

	return so, nil
}

func (s *transactionStatisticByCardService) FindYearlyTransactionStatusSuccessByCardNumber(req *requests.YearStatusTransactionCardNumber) ([]*response.TransactionResponseYearStatusSuccess, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTransactionStatusSuccessByCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTransactionStatusSuccessByCardNumber")
	defer span.End()

	card_number := req.CardNumber
	year := req.Year

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.String("card_number", card_number),
	)

	s.logger.Debug("Fetching yearly Transaction status success", zap.Int("year", year), zap.String("card_number", card_number))

	records, err := s.transactionStatisticByCardRepository.GetYearlyTransactionStatusSuccessByCardNumber(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEAR_TRANSACTION_SUCCESS_BY_CARD")

		s.logger.Error("failed to fetch yearly Transaction status success", zap.Error(err), zap.Int("year", year), zap.String("card_number", card_number), zap.String("trace_id", traceID))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch yearly Transaction status success")
		status = "failed_find_year_transaction_success_by_card"

		return nil, transaction_errors.ErrFailedFindYearTransactionSuccessByCard
	}

	s.logger.Debug("Successfully fetched yearly Transaction status success", zap.Int("year", year), zap.String("card_number", card_number))

	so := s.mapping.ToTransactionResponsesYearStatusSuccess(records)

	return so, nil
}

func (s *transactionStatisticByCardService) FindMonthTransactionStatusFailedByCardNumber(req *requests.MonthStatusTransactionCardNumber) ([]*response.TransactionResponseMonthStatusFailed, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthTransactionStatusFailedByCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthTransactionStatusFailedByCardNumber")
	defer span.End()

	card_number := req.CardNumber
	year := req.Year
	month := req.Month

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
		attribute.String("card_number", card_number),
	)

	s.logger.Debug("Fetching monthly Transaction status Failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	records, err := s.transactionStatisticByCardRepository.GetMonthTransactionStatusFailedByCardNumber(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTH_TRANSACTION_FAILED_BY_CARD")

		s.logger.Error("failed to fetch monthly Transaction status Failed", zap.Error(err), zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch monthly Transaction status Failed")
		status = "failed_find_month_transaction_failed_by_card"

		return nil, transaction_errors.ErrFailedFindMonthTransactionFailedByCard
	}

	s.logger.Debug("Failedfully fetched monthly Transaction status Failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	so := s.mapping.ToTransactionResponsesMonthStatusFailed(records)

	return so, nil
}

func (s *transactionStatisticByCardService) FindYearlyTransactionStatusFailedByCardNumber(req *requests.YearStatusTransactionCardNumber) ([]*response.TransactionResponseYearStatusFailed, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTransactionStatusFailedByCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTransactionStatusFailedByCardNumber")
	defer span.End()

	card_number := req.CardNumber
	year := req.Year

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.String("card_number", card_number),
	)

	s.logger.Debug("Fetching yearly Transaction status Failed", zap.Int("year", year), zap.String("card_number", card_number))

	records, err := s.transactionStatisticByCardRepository.GetYearlyTransactionStatusFailedByCardNumber(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEAR_TRANSACTION_FAILED_BY_CARD")

		s.logger.Error("failed to fetch yearly Transaction status Failed", zap.Error(err), zap.Int("year", year), zap.String("card_number", card_number))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch yearly Transaction status Failed")
		status = "failed_find_year_transaction_failed_by_card"

		return nil, transaction_errors.ErrFailedFindYearTransactionFailedByCard
	}

	s.logger.Debug("Failedfully fetched yearly Transaction status Failed", zap.Int("year", year))

	so := s.mapping.ToTransactionResponsesYearStatusFailed(records)

	return so, nil
}

func (s *transactionStatisticByCardService) FindMonthlyPaymentMethodsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*response.TransactionMonthMethodResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyPaymentMethodsByCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyPaymentMethodsByCardNumber")
	defer span.End()

	cardNumber := req.CardNumber
	year := req.Year

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.String("card_number", cardNumber),
	)

	s.logger.Debug("Fetching monthly payment methods by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	records, err := s.transactionStatisticByCardRepository.GetMonthlyPaymentMethodsByCardNumber(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_PAYMENT_METHODS_BY_CARD")

		s.logger.Error("Failed to fetch monthly payment methods by card number", zap.Error(err), zap.String("card_number", cardNumber), zap.Int("year", year))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch monthly payment methods by card number")
		status = "failed_to_fetch_monthly_payment_methods_by_card"

		return nil, transaction_errors.ErrFailedFindMonthlyPaymentMethodsByCard
	}

	responses := s.mapping.ToTransactionMonthlyMethodResponses(records)

	s.logger.Debug("Successfully fetched monthly payment methods by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

func (s *transactionStatisticByCardService) FindYearlyPaymentMethodsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*response.TransactionYearMethodResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyPaymentMethodsByCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyPaymentMethodsByCardNumber")
	defer span.End()

	cardNumber := req.CardNumber
	year := req.Year

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.String("card_number", cardNumber),
	)

	s.logger.Debug("Fetching yearly payment methods by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	records, err := s.transactionStatisticByCardRepository.GetYearlyPaymentMethodsByCardNumber(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_PAYMENT_METHODS_BY_CARD")

		s.logger.Error("Failed to fetch yearly payment methods by card number", zap.Error(err), zap.String("card_number", cardNumber), zap.Int("year", year))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch yearly payment methods by card number")
		status = "failed_to_fetch_yearly_payment_methods_by_card"

		return nil, transaction_errors.ErrFailedFindYearlyPaymentMethodsByCard
	}

	responses := s.mapping.ToTransactionYearlyMethodResponses(records)

	s.logger.Debug("Successfully fetched yearly payment methods by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

func (s *transactionStatisticByCardService) FindMonthlyAmountsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*response.TransactionMonthAmountResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyAmountsByCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyAmountsByCardNumber")
	defer span.End()

	cardNumber := req.CardNumber
	year := req.Year

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.String("card_number", cardNumber),
	)

	s.logger.Debug("Fetching monthly amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	records, err := s.transactionStatisticByCardRepository.GetMonthlyAmountsByCardNumber(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_AMOUNTS_BY_CARD")

		s.logger.Error("Failed to fetch monthly amounts by card number", zap.Error(err), zap.String("card_number", cardNumber), zap.Int("year", year))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch monthly amounts by card number")
		status = "failed_to_fetch_monthly_amounts_by_card"

		return nil, transaction_errors.ErrFailedFindMonthlyAmountsByCard
	}

	responses := s.mapping.ToTransactionMonthlyAmountResponses(records)

	s.logger.Debug("Successfully fetched monthly amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

func (s *transactionStatisticByCardService) FindYearlyAmountsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*response.TransactionYearlyAmountResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyAmountsByCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyAmountsByCardNumber")
	defer span.End()

	cardNumber := req.CardNumber
	year := req.Year

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.String("card_number", cardNumber),
	)

	s.logger.Debug("Fetching yearly amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	records, err := s.transactionStatisticByCardRepository.GetYearlyAmountsByCardNumber(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_AMOUNTS_BY_CARD")

		s.logger.Error("Failed to fetch yearly amounts by card number", zap.Error(err), zap.String("card_number", cardNumber), zap.Int("year", year))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch yearly amounts by card number")
		status = "failed_to_fetch_yearly_amounts_by_card"

		return nil, transaction_errors.ErrFailedFindYearlyAmountsByCard
	}

	responses := s.mapping.ToTransactionYearlyAmountResponses(records)

	s.logger.Debug("Successfully fetched yearly amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

func (s *transactionStatisticByCardService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
