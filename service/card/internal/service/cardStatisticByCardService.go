package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-card/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-card/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-card/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type cardStatisticBycardService struct {
	ctx                     context.Context
	errorhandler            errorhandler.CardStatisticByNumberErrorHandler
	mencache                mencache.CardStatisticByNumberCache
	trace                   trace.Tracer
	cardStatisticRepository repository.CardStatisticByCardRepository
	logger                  logger.LoggerInterface
	mapping                 responseservice.CardResponseMapper
	requestCounter          *prometheus.CounterVec
	requestDuration         *prometheus.HistogramVec
}

func NewCardStatisticBycardService(
	ctx context.Context,
	errorhandler errorhandler.CardStatisticByNumberErrorHandler,
	mencache mencache.CardStatisticByNumberCache,
	cardStatisticRepository repository.CardStatisticByCardRepository, logger logger.LoggerInterface, mapper responseservice.CardResponseMapper) *cardStatisticBycardService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "card_statistic_bycard_request_count",
		Help: "Number of card statistic requests CardStatisticBycardService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "card_statistic_bycard_request_duration_seconds",
		Help:    "Duration of card statistic requests CardStatisticBycardService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	return &cardStatisticBycardService{
		ctx:                     ctx,
		errorhandler:            errorhandler,
		mencache:                mencache,
		trace:                   otel.Tracer("card-statistic-bycard-service"),
		cardStatisticRepository: cardStatisticRepository,
		logger:                  logger,
		mapping:                 mapper,
		requestCounter:          requestCounter,
		requestDuration:         requestDuration,
	}
}

func (s *cardStatisticBycardService) FindMonthlyBalanceByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthBalance, *response.ErrorResponse) {
	const method = "FindMonthlyBalanceByCardNumber"

	year := req.Year

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", req.CardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyBalanceCache(req); found {
		logSuccess("Cache hit for monthly balance card", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetMonthlyBalancesByCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleMonthlyBalanceByCardNumberError(err, method, "FAILED_MONTHLY_BALANCE_BY_CARD", span, &status, zap.Error(err))
	}

	so := s.mapping.ToGetMonthlyBalances(res)

	s.mencache.SetMonthlyBalanceCache(req, so)

	logSuccess("Successfully fetched monthly balance card", zap.Int("year", year))

	return so, nil
}

func (s *cardStatisticBycardService) FindYearlyBalanceByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearlyBalance, *response.ErrorResponse) {
	const method = "FindYearlyBalanceByCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", req.Year), attribute.String("card_number", req.CardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyBalanceCache(req); found {
		logSuccess("Cache hit for yearly balance card", zap.Int("year", req.Year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetYearlyBalanceByCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleYearlyBalanceByCardNumberError(err, method, "FAILED_YEARLY_BALANCE_BY_CARD", span, &status, zap.Error(err))
	}

	so := s.mapping.ToGetYearlyBalances(res)

	s.mencache.SetYearlyBalanceCache(req, so)

	logSuccess("Successfully fetched yearly balance card", zap.Int("year", req.Year))

	return so, nil
}

func (s *cardStatisticBycardService) FindMonthlyTopupAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	const method = "FindMonthlyTopupAmountByCardNumber"

	cardNumber := req.CardNumber
	year := req.Year

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyTopupAmountCache(req); found {
		logSuccess("Cache hit for monthly topup amount card", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetMonthlyTopupAmountByCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleMonthlyTopupAmountByCardNumberError(err, "FindMonthlyTopupAmount", "FAILED_MONTHLY_TOPUP_AMOUNT_BY_CARD", span, &status, zap.Error(err))
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

	s.mencache.SetMonthlyTopupAmountCache(req, so)

	s.logger.Debug("Monthly topup amount by card number retrieved successfully",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticBycardService) FindYearlyTopupAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	const method = "FindYearlyTopupAmountByCardNumber"

	cardNumber := req.CardNumber
	year := req.Year

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyTopupAmountCache(req); found {
		logSuccess("Cache hit for yearly topup amount card", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetYearlyTopupAmountByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleYearlyTopupAmountByCardNumberError(err, method, "FAILED_YEARLY_TOPUP_AMOUNT_BY_CARD", span, &status, zap.Error(err))
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.mencache.SetYearlyTopupAmountCache(req, so)

	logSuccess("Successfully fetched yearly topup amount card", zap.String("card_number", cardNumber), zap.Int("year", year))

	return so, nil
}

func (s *cardStatisticBycardService) FindMonthlyWithdrawAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	const method = "RefreshToken"

	cardNumber := req.CardNumber
	year := req.Year

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyWithdrawAmountCache(req); found {
		logSuccess("Cache hit for monthly withdraw amount card", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetMonthlyWithdrawAmountByCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleMonthlyWithdrawAmountByCardNumberError(err, method, "FAILED_MONTHLY_WITHDRAW_AMOUNT_BY_CARD", span, &status, zap.Error(err))
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

	s.mencache.SetMonthlyWithdrawAmountCache(req, so)

	logSuccess("Successfully fetched monthly withdraw amount card", zap.String("card_number", cardNumber), zap.Int("year", year))

	return so, nil
}

func (s *cardStatisticBycardService) FindYearlyWithdrawAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	const method = "FindYearlyWithdrawAmountByCardNumber"

	cardNumber := req.CardNumber
	year := req.Year

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyWithdrawAmountCache(req); found {
		logSuccess("Cache hit for yearly withdraw amount card", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetYearlyWithdrawAmountByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleYearlyWithdrawAmountByCardNumberError(err, method, "FAILED_YEARLY_WITHDRAW_AMOUNT_BY_CARD", span, &status, zap.Error(err))
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.mencache.SetYearlyWithdrawAmountCache(req, so)

	logSuccess("Successfully fetched yearly withdraw amount card", zap.String("card_number", cardNumber), zap.Int("year", year))

	return so, nil
}

func (s *cardStatisticBycardService) FindMonthlyTransactionAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	const method = "FindMonthlyTransactionAmountByCardNumber"

	cardNumber := req.CardNumber
	year := req.Year

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyTransactionAmountCache(req); found {
		logSuccess("Cache hit for monthly transaction amount card", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetMonthlyTransactionAmountByCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleMonthlyTransactionAmountByCardNumberError(err, method, "FAILED_MONTHLY_TRANSACTION_AMOUNT_BY_CARD", span, &status, zap.Error(err))
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

	s.mencache.SetMonthlyTransactionAmountCache(req, so)

	logSuccess("Successfully fetched monthly transaction amount card", zap.String("card_number", cardNumber), zap.Int("year", year))

	return so, nil
}

func (s *cardStatisticBycardService) FindYearlyTransactionAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	const method = "FindYearlyTransactionAmountByCardNumber"

	cardNumber := req.CardNumber
	year := req.Year

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyTransactionAmountCache(req); found {
		logSuccess("Cache hit for yearly transaction amount card", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetYearlyTransactionAmountByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleYearlyTransactionAmountByCardNumberError(err, method, "FAILED_YEARLY_TRANSACTION_AMOUNT_BY_CARD", span, &status, zap.Error(err))
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.mencache.SetYearlyTransactionAmountCache(req, so)

	s.logger.Debug("Yearly transaction amount by card number retrieved successfully",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticBycardService) FindMonthlyTransferAmountBySender(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	const method = "FindMonthlyTransferAmountBySender"

	cardNumber := req.CardNumber
	year := req.Year

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyTransferBySenderCache(req); found {
		logSuccess("Cache hit for monthly transfer sender amount card", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetMonthlyTransferAmountBySender(req)

	if err != nil {
		return s.errorhandler.HandleMonthlyTransferAmountBySenderError(err, method, "FAILED_MONTHLY_TRANSFER_AMOUNT_BY_SENDER", span, &status, zap.Error(err))
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

	s.mencache.SetMonthlyTransferBySenderCache(req, so)

	logSuccess("Successfully fetched monthly transfer sender amount card", zap.String("card_number", cardNumber), zap.Int("year", year))

	return so, nil
}

func (s *cardStatisticBycardService) FindYearlyTransferAmountBySender(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	const method = "FindYearlyTransferAmountBySender"

	cardNumber := req.CardNumber
	year := req.Year

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyTransferBySenderCache(req); found {
		logSuccess("Cache hit for yearly transfer sender amount card", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetYearlyTransferAmountBySender(req)
	if err != nil {
		return s.errorhandler.HandleYearlyTransferAmountBySenderError(err, method, "FAILED_YEARLY_TRANSFER_AMOUNT_BY_SENDER", span, &status, zap.Error(err))
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.mencache.SetYearlyTransferBySenderCache(req, so)

	logSuccess("Successfully fetched yearly transfer sender amount card", zap.String("card_number", cardNumber), zap.Int("year", year))

	return so, nil
}

func (s *cardStatisticBycardService) FindMonthlyTransferAmountByReceiver(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	const method = "FindMonthlyTransferAmountByReceiver"

	cardNumber := req.CardNumber
	year := req.Year

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyTransferByReceiverCache(req); found {
		logSuccess("Cache hit for monthly transfer receiver amount card", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetMonthlyTransferAmountByReceiver(req)

	if err != nil {
		return s.errorhandler.HandleMonthlyTransferAmountByReceiverError(err, method, "FAILED_MONTHLY_TRANSFER_AMOUNT_BY_RECEIVER", span, &status, zap.Error(err))
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

	s.mencache.SetMonthlyTransferByReceiverCache(req, so)

	logSuccess("Successfully fetched monthly transfer receiver amount card", zap.String("card_number", cardNumber), zap.Int("year", year))

	return so, nil
}

func (s *cardStatisticBycardService) FindYearlyTransferAmountByReceiver(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	const method = "FindYearlyTransferAmountByReceiver"

	cardNumber := req.CardNumber
	year := req.Year

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyTransferByReceiverCache(req); found {
		logSuccess("Cache hit for yearly transfer receiver amount card", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetYearlyTransferAmountByReceiver(req)
	if err != nil {
		return s.errorhandler.HandleYearlyTransferAmountByReceiverError(err, method, "FAILED_YEARLY_TRANSFER_AMOUNT_BY_RECEIVER", span, &status, zap.Error(err))
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.mencache.SetYearlyTransferByReceiverCache(req, so)

	logSuccess("Successfully fetched yearly transfer receiver amount card", zap.String("card_number", cardNumber), zap.Int("year", year))

	return so, nil
}

func (s *cardStatisticBycardService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *cardStatisticBycardService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
