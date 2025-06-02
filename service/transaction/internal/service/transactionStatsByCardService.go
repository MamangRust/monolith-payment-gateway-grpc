package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-transaction/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transactionStatisticByCardService struct {
	ctx                                  context.Context
	errorhandler                         errorhandler.TransactionStatisticByCardErrorHandler
	mencache                             mencache.TransactionStatisticByCardCache
	trace                                trace.Tracer
	transactionStatisticByCardRepository repository.TransactionStatisticByCardRepository
	logger                               logger.LoggerInterface
	mapping                              responseservice.TransactionResponseMapper
	requestCounter                       *prometheus.CounterVec
	requestDuration                      *prometheus.HistogramVec
}

func NewTransactionStatisticByCardService(
	ctx context.Context,
	errorhandler errorhandler.TransactionStatisticByCardErrorHandler,
	mencache mencache.TransactionStatisticByCardCache,
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
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &transactionStatisticByCardService{
		ctx:                                  ctx,
		errorhandler:                         errorhandler,
		mencache:                             mencache,
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

	if data := s.mencache.GetMonthTransactionStatusSuccessByCardCache(req); data != nil {
		s.logger.Debug("Successfully fetched monthly Transaction status success from cache", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.transactionStatisticByCardRepository.GetMonthTransactionStatusSuccessByCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleMonthTransactionStatusSuccessByCardNumberError(err, "FindMonthTransactionStatusSuccessByCardNumber", "FAILED_FIND_MONTHLY_TRANSACTION_STATUS_SUCCESS_BY_CARD", span, &status, zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))
	}
	so := s.mapping.ToTransactionResponsesMonthStatusSuccess(records)

	s.mencache.SetMonthTransactionStatusSuccessByCardCache(req, so)

	s.logger.Debug("Successfully fetched monthly Transaction status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

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

	if data := s.mencache.GetYearTransactionStatusSuccessByCardCache(req); data != nil {
		s.logger.Debug("Successfully fetched yearly Transaction status success from cache", zap.Int("year", year), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.transactionStatisticByCardRepository.GetYearlyTransactionStatusSuccessByCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleYearlyTransactionStatusSuccessByCardNumberError(err, "FindYearlyTransactionStatusSuccessByCardNumber", "FAILED_FIND_YEARLY_TRANSACTION_STATUS_SUCCESS_BY_CARD", span, &status, zap.Int("year", year), zap.String("card_number", card_number))
	}
	so := s.mapping.ToTransactionResponsesYearStatusSuccess(records)

	s.mencache.SetYearTransactionStatusSuccessByCardCache(req, so)

	s.logger.Debug("Successfully fetched yearly Transaction status success", zap.Int("year", year), zap.String("card_number", card_number))

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

	if data := s.mencache.GetMonthTransactionStatusFailedByCardCache(req); data != nil {
		s.logger.Debug("Successfully fetched monthly Transaction status Failed from cache", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.transactionStatisticByCardRepository.GetMonthTransactionStatusFailedByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleMonthTransactionStatusFailedByCardNumberError(err, "FindMonthTransactionStatusFailedByCardNumber", "FAILED_MONTHLY_TRANSACTION_STATUS_FAILED_BY_CARD", span, &status, zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))
	}
	so := s.mapping.ToTransactionResponsesMonthStatusFailed(records)

	s.mencache.SetMonthTransactionStatusFailedByCardCache(req, so)

	s.logger.Debug("Failedfully fetched monthly Transaction status Failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

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

	if data := s.mencache.GetYearTransactionStatusFailedByCardCache(req); data != nil {
		return data, nil
	}

	records, err := s.transactionStatisticByCardRepository.GetYearlyTransactionStatusFailedByCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleYearlyTransactionStatusFailedByCardNumberError(err, "FindYearlyTransactionStatusFailedByCardNumber", "FAILED_FIND_YEARLY_TRANSACTION_STATUS_FAILED_BY_CARD", span, &status, zap.Int("year", year), zap.String("card_number", card_number))
	}
	so := s.mapping.ToTransactionResponsesYearStatusFailed(records)

	s.mencache.SetYearTransactionStatusFailedByCardCache(req, so)

	s.logger.Debug("Failedfully fetched yearly Transaction status Failed", zap.Int("year", year))

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

	if data := s.mencache.GetMonthlyPaymentMethodsByCardCache(req); data != nil {
		s.logger.Debug("Successfully fetched monthly payment methods by card number from cache", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	records, err := s.transactionStatisticByCardRepository.GetMonthlyPaymentMethodsByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleMonthlyPaymentMethodsByCardNumberError(err, "FindMonthlyPaymentMethodsByCardNumber", "FAILED_FIND_MONTHLY_PAYMENT_METHODS_BY_CARD", span, &status, zap.String("card_number", cardNumber), zap.Int("year", year))
	}

	responses := s.mapping.ToTransactionMonthlyMethodResponses(records)

	s.mencache.SetMonthlyPaymentMethodsByCardCache(req, responses)

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

	if data := s.mencache.GetYearlyPaymentMethodsByCardCache(req); data != nil {
		s.logger.Debug("Successfully fetched yearly payment methods by card number from cache", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	records, err := s.transactionStatisticByCardRepository.GetYearlyPaymentMethodsByCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleYearlyPaymentMethodsByCardNumberError(err, "FindYearlyPaymentMethodsByCardNumber", "FAILED_FIND_YEARLY_PAYMENT_METHODS_BY_CARD", span, &status, zap.String("card_number", cardNumber), zap.Int("year", year))
	}

	responses := s.mapping.ToTransactionYearlyMethodResponses(records)

	s.mencache.SetYearlyPaymentMethodsByCardCache(req, responses)

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

	if data := s.mencache.GetMonthlyAmountsByCardCache(req); data != nil {
		s.logger.Debug("Successfully fetched monthly amounts by card number from cache", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	records, err := s.transactionStatisticByCardRepository.GetMonthlyAmountsByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleMonthlyAmountsByCardNumberError(err, "FindMonthlyAmountsByCardNumber", "FAILED_FIND_MONTHLY_AMOUNTS_BY_CARD", span, &status, zap.String("card_number", cardNumber), zap.Int("year", year))
	}

	responses := s.mapping.ToTransactionMonthlyAmountResponses(records)

	s.mencache.SetMonthlyAmountsByCardCache(req, responses)

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

	if data := s.mencache.GetYearlyAmountsByCardCache(req); data != nil {
		s.logger.Debug("Successfully fetched yearly amounts by card number from cache", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	records, err := s.transactionStatisticByCardRepository.GetYearlyAmountsByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleYearlyAmountsByCardNumberError(err, "FindYearlyAmountsByCardNumber", "FAILED_FIND_YEARLY_AMOUNTS_BY_CARD", span, &status, zap.String("card_number", cardNumber), zap.Int("year", year))
	}

	responses := s.mapping.ToTransactionYearlyAmountResponses(records)

	s.mencache.SetYearlyAmountsByCardCache(req, responses)

	s.logger.Debug("Successfully fetched yearly amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

func (s *transactionStatisticByCardService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
