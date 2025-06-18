package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-withdraw/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type withdrawStatisticByCardService struct {
	ctx                               context.Context
	mencache                          mencache.WithdrawStasticByCardCache
	errorhandler                      errorhandler.WithdrawStatisticByCardErrorHandler
	trace                             trace.Tracer
	saldoRepository                   repository.SaldoRepository
	withdrawStatisticByCardRepository repository.WithdrawStatisticByCardRepository
	logger                            logger.LoggerInterface
	mapping                           responseservice.WithdrawResponseMapper
	requestCounter                    *prometheus.CounterVec
	requestDuration                   *prometheus.HistogramVec
}

func NewWithdrawStatisticByCardService(
	ctx context.Context,
	mencache mencache.WithdrawStasticByCardCache,
	errorhandler errorhandler.WithdrawStatisticByCardErrorHandler,
	withdrawStatisticByCardRepository repository.WithdrawStatisticByCardRepository, saldoRepository repository.SaldoRepository, logger logger.LoggerInterface, mapping responseservice.WithdrawResponseMapper) *withdrawStatisticByCardService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "withdraw_statistic_by_card_service_request_total",
			Help: "Total number of requests to the WithdrawStatisticByCardService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "withdraw_statistic_by_card_service_request_duration_seconds",
			Help:    "Histogram of request durations for the WithdrawStatisticByCardService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &withdrawStatisticByCardService{
		ctx:                               ctx,
		trace:                             otel.Tracer("withraw-statistic-by-card-service"),
		saldoRepository:                   saldoRepository,
		withdrawStatisticByCardRepository: withdrawStatisticByCardRepository,
		logger:                            logger,
		mapping:                           mapping,
		requestCounter:                    requestCounter,
		requestDuration:                   requestDuration,
		errorhandler:                      errorhandler,
		mencache:                          mencache,
	}
}

func (s *withdrawStatisticByCardService) FindMonthWithdrawStatusSuccessByCardNumber(req *requests.MonthStatusWithdrawCardNumber) ([]*response.WithdrawResponseMonthStatusSuccess, *response.ErrorResponse) {
	year := req.Year
	month := req.Month
	cardNumber := req.CardNumber

	const method = "FindMonthWithdrawStatusSuccessByCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("month", month), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthWithdrawStatusSuccessByCardNumber(req); found {
		logSuccess("Cache hit for monthly withdraw status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", cardNumber))
		return data, nil
	}

	records, err := s.withdrawStatisticByCardRepository.GetMonthWithdrawStatusSuccessByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleMonthWithdrawStatusSuccessByCardNumberError(err, method, "FAILED_GET_MONTH_WITHDRAW_STATUS_SUCCESS_BY_CARD_NUMBER", span, &status, zap.Error(err))
	}

	responseData := s.mapping.ToWithdrawResponsesMonthStatusSuccess(records)

	s.mencache.SetCachedMonthWithdrawStatusSuccessByCardNumber(req, responseData)

	logSuccess("Successfully fetched monthly withdraw status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", cardNumber))

	return responseData, nil
}

func (s *withdrawStatisticByCardService) FindYearlyWithdrawStatusSuccessByCardNumber(req *requests.YearStatusWithdrawCardNumber) ([]*response.WithdrawResponseYearStatusSuccess, *response.ErrorResponse) {

	year := req.Year
	cardNumber := req.CardNumber

	const method = "FindYearlyWithdrawStatusSuccessByCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearlyWithdrawStatusSuccessByCardNumber(req); found {
		logSuccess("Cache hit for yearly withdraw status success", zap.Int("year", year), zap.String("card_number", cardNumber))
		return data, nil
	}

	records, err := s.withdrawStatisticByCardRepository.GetYearlyWithdrawStatusSuccessByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleYearWithdrawStatusSuccessByCardNumberError(err, method, "FAILED_GET_YEARLY_WITHDRAW_STATUS_SUCCESS_BY_CARD_NUMBER", span, &status, zap.Error(err))
	}

	responseData := s.mapping.ToWithdrawResponsesYearStatusSuccess(records)

	s.mencache.SetCachedYearlyWithdrawStatusSuccessByCardNumber(req, responseData)

	logSuccess("Successfully fetched yearly withdraw status success", zap.Int("year", year), zap.String("card_number", cardNumber))

	return responseData, nil
}

func (s *withdrawStatisticByCardService) FindMonthWithdrawStatusFailedByCardNumber(req *requests.MonthStatusWithdrawCardNumber) ([]*response.WithdrawResponseMonthStatusFailed, *response.ErrorResponse) {
	year := req.Year
	month := req.Month
	cardNumber := req.CardNumber

	const method = "FindMonthWithdrawStatusFailedByCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("month", month), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthWithdrawStatusFailedByCardNumber(req); found {
		logSuccess("Cache hit for monthly withdraw status failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", cardNumber))
		return data, nil
	}

	records, err := s.withdrawStatisticByCardRepository.GetMonthWithdrawStatusFailedByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleMonthWithdrawStatusFailedByCardNumberError(err, method, "FAILED_GET_MONTH_WITHDRAW_STATUS_FAILED_BY_CARD_NUMBER", span, &status, zap.Error(err))
	}

	responseData := s.mapping.ToWithdrawResponsesMonthStatusFailed(records)

	s.mencache.SetCachedMonthWithdrawStatusFailedByCardNumber(req, responseData)

	logSuccess("Successfully fetched monthly withdraw status failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", cardNumber))

	return responseData, nil
}

func (s *withdrawStatisticByCardService) FindYearlyWithdrawStatusFailedByCardNumber(req *requests.YearStatusWithdrawCardNumber) ([]*response.WithdrawResponseYearStatusFailed, *response.ErrorResponse) {
	year := req.Year
	cardNumber := req.CardNumber

	const method = "FindYearlyWithdrawStatusFailedByCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearlyWithdrawStatusFailedByCardNumber(req); found {
		logSuccess("Cache hit for yearly withdraw status failed", zap.Int("year", year), zap.String("card_number", cardNumber))
		return data, nil
	}

	records, err := s.withdrawStatisticByCardRepository.GetYearlyWithdrawStatusFailedByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleYearWithdrawStatusFailedByCardNumberError(err, method, "FAILED_GET_YEARLY_WITHDRAW_STATUS_FAILED_BY_CARD_NUMBER", span, &status, zap.Error(err))
	}

	responseData := s.mapping.ToWithdrawResponsesYearStatusFailed(records)

	s.mencache.SetCachedYearlyWithdrawStatusFailedByCardNumber(req, responseData)

	logSuccess("Successfully fetched yearly withdraw status failed", zap.Int("year", year), zap.String("card_number", cardNumber))

	return responseData, nil
}

func (s *withdrawStatisticByCardService) FindMonthlyWithdrawsByCardNumber(req *requests.YearMonthCardNumber) ([]*response.WithdrawMonthlyAmountResponse, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	const method = "FindMonthlyWithdrawsByCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthlyWithdrawsByCardNumber(req); found {
		logSuccess("Cache hit for monthly withdraws by card number", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	withdraws, err := s.withdrawStatisticByCardRepository.GetMonthlyWithdrawsByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleMonthlyWithdrawsAmountByCardNumberError(err, method, "FAILED_GET_MONTHLY_WITHDRAW_BY_CARD_NUMBER", span, &status, zap.Error(err))
	}

	responseWithdraws := s.mapping.ToWithdrawsAmountMonthlyResponses(withdraws)

	s.mencache.SetCachedMonthlyWithdrawsByCardNumber(req, responseWithdraws)

	logSuccess("Successfully fetched monthly withdraws by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseWithdraws, nil
}

func (s *withdrawStatisticByCardService) FindYearlyWithdrawsByCardNumber(req *requests.YearMonthCardNumber) ([]*response.WithdrawYearlyAmountResponse, *response.ErrorResponse) {
	cardNumber := req.CardNumber
	year := req.Year

	const method = "FindYearlyWithdrawsByCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearlyWithdrawsByCardNumber(req); found {
		logSuccess("Cache hit for yearly withdraws by card number", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	withdraws, err := s.withdrawStatisticByCardRepository.GetYearlyWithdrawsByCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleYearlyWithdrawsAmountByCardNumberError(err, method, "FAILED_GET_YEARLY_WITHDRAW_BY_CARD_NUMBER", span, &status, zap.Error(err))
	}

	responseWithdraws := s.mapping.ToWithdrawsAmountYearlyResponses(withdraws)

	s.mencache.SetCachedYearlyWithdrawsByCardNumber(req, responseWithdraws)

	logSuccess("Successfully fetched yearly withdraws by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseWithdraws, nil
}

func (s *withdrawStatisticByCardService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *withdrawStatisticByCardService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
