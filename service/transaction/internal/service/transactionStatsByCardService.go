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
	"go.opentelemetry.io/otel/codes"
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
	card_number := req.CardNumber
	year := req.Year
	month := req.Month

	const method = "FindMonthTransactionStatusSuccessByCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("month", month), attribute.String("card_number", card_number))

	defer func() {
		end(status)
	}()

	if data := s.mencache.GetMonthTransactionStatusSuccessByCardCache(req); data != nil {
		logSuccess("Successfully fetched monthly Transaction status success from cache", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.transactionStatisticByCardRepository.GetMonthTransactionStatusSuccessByCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleMonthTransactionStatusSuccessByCardNumberError(err, method, "FAILED_FIND_MONTHLY_TRANSACTION_STATUS_SUCCESS_BY_CARD", span, &status, zap.Int("year", year), zap.Int("month", month), zap.Error(err))
	}
	so := s.mapping.ToTransactionResponsesMonthStatusSuccess(records)

	s.mencache.SetMonthTransactionStatusSuccessByCardCache(req, so)

	logSuccess("Successfully fetched monthly Transaction status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	return so, nil
}

func (s *transactionStatisticByCardService) FindYearlyTransactionStatusSuccessByCardNumber(req *requests.YearStatusTransactionCardNumber) ([]*response.TransactionResponseYearStatusSuccess, *response.ErrorResponse) {
	card_number := req.CardNumber
	year := req.Year

	const method = "FindYearlyTransactionStatusSuccessByCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", card_number))

	defer func() {
		end(status)
	}()

	if data := s.mencache.GetYearTransactionStatusSuccessByCardCache(req); data != nil {
		logSuccess("Successfully fetched yearly Transaction status success from cache", zap.Int("year", year), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.transactionStatisticByCardRepository.GetYearlyTransactionStatusSuccessByCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleYearlyTransactionStatusSuccessByCardNumberError(err, method, "FAILED_FIND_YEARLY_TRANSACTION_STATUS_SUCCESS_BY_CARD", span, &status, zap.Error(err))
	}
	so := s.mapping.ToTransactionResponsesYearStatusSuccess(records)

	s.mencache.SetYearTransactionStatusSuccessByCardCache(req, so)

	logSuccess("Successfully fetched yearly Transaction status success", zap.Int("year", year), zap.String("card_number", card_number))

	return so, nil
}

func (s *transactionStatisticByCardService) FindMonthTransactionStatusFailedByCardNumber(req *requests.MonthStatusTransactionCardNumber) ([]*response.TransactionResponseMonthStatusFailed, *response.ErrorResponse) {
	card_number := req.CardNumber
	year := req.Year
	month := req.Month

	const method = "FindMonthTransactionStatusFailedByCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("month", month), attribute.String("card_number", card_number))

	defer func() {
		end(status)
	}()

	if data := s.mencache.GetMonthTransactionStatusFailedByCardCache(req); data != nil {
		logSuccess("Successfully fetched monthly Transaction status Failed from cache", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.transactionStatisticByCardRepository.GetMonthTransactionStatusFailedByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleMonthTransactionStatusFailedByCardNumberError(err, method, "FAILED_MONTHLY_TRANSACTION_STATUS_FAILED_BY_CARD", span, &status, zap.Error(err))
	}
	so := s.mapping.ToTransactionResponsesMonthStatusFailed(records)

	s.mencache.SetMonthTransactionStatusFailedByCardCache(req, so)

	logSuccess("Successfully fetched monthly Transaction status Failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	return so, nil
}

func (s *transactionStatisticByCardService) FindYearlyTransactionStatusFailedByCardNumber(req *requests.YearStatusTransactionCardNumber) ([]*response.TransactionResponseYearStatusFailed, *response.ErrorResponse) {
	card_number := req.CardNumber
	year := req.Year

	const method = "FindYearlyTransactionStatusFailedByCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", card_number))

	defer func() {
		end(status)
	}()

	if data := s.mencache.GetYearTransactionStatusFailedByCardCache(req); data != nil {
		logSuccess("Successfully fetched yearly Transaction status Failed from cache", zap.Int("year", year), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.transactionStatisticByCardRepository.GetYearlyTransactionStatusFailedByCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleYearlyTransactionStatusFailedByCardNumberError(err, method, "FAILED_FIND_YEARLY_TRANSACTION_STATUS_FAILED_BY_CARD", span, &status, zap.Error(err))
	}
	so := s.mapping.ToTransactionResponsesYearStatusFailed(records)

	s.mencache.SetYearTransactionStatusFailedByCardCache(req, so)

	logSuccess("Successfully fetched yearly Transaction status Failed", zap.Int("year", year), zap.String("card_number", card_number))

	return so, nil
}

func (s *transactionStatisticByCardService) FindMonthlyPaymentMethodsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*response.TransactionMonthMethodResponse, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	const method = "FindMonthlyPaymentMethodsByCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data := s.mencache.GetMonthlyPaymentMethodsByCardCache(req); data != nil {
		logSuccess("Successfully fetched monthly payment methods by card number from cache", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	records, err := s.transactionStatisticByCardRepository.GetMonthlyPaymentMethodsByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleMonthlyPaymentMethodsByCardNumberError(err, method, "FAILED_FIND_MONTHLY_PAYMENT_METHODS_BY_CARD", span, &status, zap.Error(err))
	}

	responses := s.mapping.ToTransactionMonthlyMethodResponses(records)

	s.mencache.SetMonthlyPaymentMethodsByCardCache(req, responses)

	logSuccess("Successfully fetched monthly payment methods by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

func (s *transactionStatisticByCardService) FindYearlyPaymentMethodsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*response.TransactionYearMethodResponse, *response.ErrorResponse) {

	cardNumber := req.CardNumber
	year := req.Year

	const method = "FindYearlyPaymentMethodsByCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data := s.mencache.GetYearlyPaymentMethodsByCardCache(req); data != nil {
		logSuccess("Successfully fetched yearly payment methods by card number from cache", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	records, err := s.transactionStatisticByCardRepository.GetYearlyPaymentMethodsByCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleYearlyPaymentMethodsByCardNumberError(err, method, "FAILED_FIND_YEARLY_PAYMENT_METHODS_BY_CARD", span, &status, zap.Error(err))
	}

	responses := s.mapping.ToTransactionYearlyMethodResponses(records)

	s.mencache.SetYearlyPaymentMethodsByCardCache(req, responses)

	logSuccess("Successfully fetched yearly payment methods by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

func (s *transactionStatisticByCardService) FindMonthlyAmountsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*response.TransactionMonthAmountResponse, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	const method = "FindMonthlyAmountsByCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data := s.mencache.GetMonthlyAmountsByCardCache(req); data != nil {
		logSuccess("Successfully fetched monthly amounts by card number from cache", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	records, err := s.transactionStatisticByCardRepository.GetMonthlyAmountsByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleMonthlyAmountsByCardNumberError(err, method, "FAILED_FIND_MONTHLY_AMOUNTS_BY_CARD", span, &status, zap.Error(err))
	}

	responses := s.mapping.ToTransactionMonthlyAmountResponses(records)

	s.mencache.SetMonthlyAmountsByCardCache(req, responses)

	logSuccess("Successfully fetched monthly amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

func (s *transactionStatisticByCardService) FindYearlyAmountsByCardNumber(req *requests.MonthYearPaymentMethod) ([]*response.TransactionYearlyAmountResponse, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	const method = "FindYearlyAmountsByCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data := s.mencache.GetYearlyAmountsByCardCache(req); data != nil {
		logSuccess("Successfully fetched yearly amounts by card number from cache", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	records, err := s.transactionStatisticByCardRepository.GetYearlyAmountsByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleYearlyAmountsByCardNumberError(err, method, "FAILED_FIND_YEARLY_AMOUNTS_BY_CARD", span, &status, zap.Error(err))
	}

	responses := s.mapping.ToTransactionYearlyAmountResponses(records)

	s.mencache.SetYearlyAmountsByCardCache(req, responses)

	logSuccess("Successfully fetched yearly amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

func (s *transactionStatisticByCardService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *transactionStatisticByCardService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
