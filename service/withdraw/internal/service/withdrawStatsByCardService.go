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
	start := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("FindMonthWithdrawStatusSuccessByCardNumber", status, start)
	}()
	_, span := s.trace.Start(s.ctx, "FindMonthWithdrawStatusSuccessByCardNumber")
	defer span.End()

	year := req.Year
	month := req.Month
	cardNumber := req.CardNumber

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
		attribute.String("card_number", cardNumber),
	)

	s.logger.Debug("Fetching monthly Withdraw status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", cardNumber))

	if data := s.mencache.GetCachedMonthWithdrawStatusSuccessByCardNumber(req); data != nil {
		s.logger.Debug("Successfully fetched monthly withdraw status success from cache", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", cardNumber))
		return data, nil
	}

	records, err := s.withdrawStatisticByCardRepository.GetMonthWithdrawStatusSuccessByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleMonthWithdrawStatusSuccessByCardNumberError(err, "FindMonthWithdrawStatusSuccessByCardNumber", "FAILED_GET_MONTH_WITHDRAW_STATUS_SUCCESS_BY_CARD_NUMBER", span, &status, zap.Error(err))
	}

	responseData := s.mapping.ToWithdrawResponsesMonthStatusSuccess(records)

	s.mencache.SetCachedMonthWithdrawStatusSuccessByCardNumber(req, responseData)

	s.logger.Debug("Successfully fetched monthly withdraw status success",
		zap.Int("year", year),
		zap.Int("month", month),
		zap.String("card_number", cardNumber),
		zap.Int("total_records", len(responseData)))

	return responseData, nil
}

func (s *withdrawStatisticByCardService) FindYearlyWithdrawStatusSuccessByCardNumber(req *requests.YearStatusWithdrawCardNumber) ([]*response.WithdrawResponseYearStatusSuccess, *response.ErrorResponse) {
	start := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("FindYearlyWithdrawStatusSuccessByCardNumber", status, start)
	}()
	_, span := s.trace.Start(s.ctx, "FindYearlyWithdrawStatusSuccessByCardNumber")
	defer span.End()

	year := req.Year
	cardNumber := req.CardNumber

	if data := s.mencache.GetCachedYearlyWithdrawStatusSuccessByCardNumber(req); data != nil {
		s.logger.Debug("Successfully fetched yearly withdraw status success from cache", zap.Int("year", year), zap.String("card_number", cardNumber))
		return data, nil
	}

	records, err := s.withdrawStatisticByCardRepository.GetYearlyWithdrawStatusSuccessByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleYearWithdrawStatusSuccessByCardNumberError(err, "FindYearlyWithdrawStatusSuccessByCardNumber", "FAILED_GET_YEARLY_WITHDRAW_STATUS_SUCCESS_BY_CARD_NUMBER", span, &status, zap.Error(err))
	}

	responseData := s.mapping.ToWithdrawResponsesYearStatusSuccess(records)

	s.mencache.SetCachedYearlyWithdrawStatusSuccessByCardNumber(req, responseData)

	s.logger.Debug("Successfully fetched yearly withdraw status success",
		zap.Int("year", year),
		zap.String("card_number", cardNumber),
		zap.Int("total_records", len(responseData)))

	return responseData, nil
}

func (s *withdrawStatisticByCardService) FindMonthWithdrawStatusFailedByCardNumber(req *requests.MonthStatusWithdrawCardNumber) ([]*response.WithdrawResponseMonthStatusFailed, *response.ErrorResponse) {
	start := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("FindMonthWithdrawStatusFailedByCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthWithdrawStatusFailedByCardNumber")
	defer span.End()

	year := req.Year
	month := req.Month
	cardNumber := req.CardNumber

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
		attribute.String("card_number", cardNumber),
	)

	s.logger.Debug("Fetching monthly Withdraw status failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", cardNumber))

	if data := s.mencache.GetCachedMonthWithdrawStatusFailedByCardNumber(req); data != nil {
		s.logger.Debug("Successfully fetched monthly withdraw status failed from cache", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", cardNumber))
		return data, nil
	}

	records, err := s.withdrawStatisticByCardRepository.GetMonthWithdrawStatusFailedByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleMonthWithdrawStatusFailedByCardNumberError(err, "FindMonthWithdrawStatusFailedByCardNumber", "FAILED_GET_MONTH_WITHDRAW_STATUS_FAILED_BY_CARD_NUMBER", span, &status, zap.Error(err))
	}

	responseData := s.mapping.ToWithdrawResponsesMonthStatusFailed(records)

	s.mencache.SetCachedMonthWithdrawStatusFailedByCardNumber(req, responseData)

	return responseData, nil
}

func (s *withdrawStatisticByCardService) FindYearlyWithdrawStatusFailedByCardNumber(req *requests.YearStatusWithdrawCardNumber) ([]*response.WithdrawResponseYearStatusFailed, *response.ErrorResponse) {
	start := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("FindYearlyWithdrawStatusFailedByCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyWithdrawStatusFailedByCardNumber")
	defer span.End()

	year := req.Year
	cardNumber := req.CardNumber

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.String("card_number", cardNumber),
	)

	if data := s.mencache.GetCachedYearlyWithdrawStatusFailedByCardNumber(req); data != nil {
		s.logger.Debug("Successfully fetched yearly withdraw status failed from cache", zap.Int("year", year), zap.String("card_number", cardNumber))
		return data, nil
	}

	s.logger.Debug("Fetching yearly Withdraw status failed", zap.Int("year", year), zap.String("card_number", cardNumber))

	records, err := s.withdrawStatisticByCardRepository.GetYearlyWithdrawStatusFailedByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleYearWithdrawStatusFailedByCardNumberError(err, "FindYearlyWithdrawStatusFailedByCardNumber", "FAILED_GET_YEARLY_WITHDRAW_STATUS_FAILED_BY_CARD_NUMBER", span, &status, zap.Error(err))
	}

	responseData := s.mapping.ToWithdrawResponsesYearStatusFailed(records)

	s.mencache.SetCachedYearlyWithdrawStatusFailedByCardNumber(req, responseData)

	return responseData, nil
}

func (s *withdrawStatisticByCardService) FindMonthlyWithdrawsByCardNumber(req *requests.YearMonthCardNumber) ([]*response.WithdrawMonthlyAmountResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyWithdrawsByCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyWithdrawsByCardNumber")
	defer span.End()

	cardNumber := req.CardNumber
	year := req.Year

	span.SetAttributes(
		attribute.String("card_number", cardNumber),
		attribute.Int("year", year),
	)

	s.logger.Debug("Fetching monthly withdraws by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	if data := s.mencache.GetCachedMonthlyWithdrawsByCardNumber(req); data != nil {
		s.logger.Debug("Successfully fetched monthly withdraws by card number from cache", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	withdraws, err := s.withdrawStatisticByCardRepository.GetMonthlyWithdrawsByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleMonthlyWithdrawsAmountByCardNumberError(err, "FindMonthlyWithdrawsByCardNumber", "FAILED_GET_MONTHLY_WITHDRAW_BY_CARD_NUMBER", span, &status, zap.Error(err))
	}

	responseWithdraws := s.mapping.ToWithdrawsAmountMonthlyResponses(withdraws)

	s.mencache.SetCachedMonthlyWithdrawsByCardNumber(req, responseWithdraws)

	s.logger.Debug("Successfully fetched monthly withdraws by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseWithdraws, nil
}

func (s *withdrawStatisticByCardService) FindYearlyWithdrawsByCardNumber(req *requests.YearMonthCardNumber) ([]*response.WithdrawYearlyAmountResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyWithdrawsByCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyWithdrawsByCardNumber")
	defer span.End()

	cardNumber := req.CardNumber
	year := req.Year

	span.SetAttributes(
		attribute.String("card_number", cardNumber),
		attribute.Int("year", year),
	)

	s.logger.Debug("Fetching yearly withdraws by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	if data := s.mencache.GetCachedYearlyWithdrawsByCardNumber(req); data != nil {
		s.logger.Debug("Successfully fetched yearly withdraws by card number from cache", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	withdraws, err := s.withdrawStatisticByCardRepository.GetYearlyWithdrawsByCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleYearlyWithdrawsAmountByCardNumberError(err, "FindYearlyWithdrawsByCardNumber", "FAILED_GET_YEARLY_WITHDRAW_BY_CARD_NUMBER", span, &status, zap.Error(err))
	}

	responseWithdraws := s.mapping.ToWithdrawsAmountYearlyResponses(withdraws)

	s.mencache.SetCachedYearlyWithdrawsByCardNumber(req, responseWithdraws)

	return responseWithdraws, nil
}

func (s *withdrawStatisticByCardService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
