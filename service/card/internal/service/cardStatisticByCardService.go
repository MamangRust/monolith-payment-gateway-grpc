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
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyBalance", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyBalance")
	defer span.End()

	year := req.Year

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("FindMonthlyBalance called", zap.Int("year", year))

	if data, found := s.mencache.GetMonthlyBalanceCache(req); found {
		s.logger.Debug("Cache hit for monthly balance card", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetMonthlyBalancesByCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleMonthlyBalanceByCardNumberError(err, "FindMonthlyBalance", "FAILED_MONTHLY_BALANCE_BY_CARD", span, &status)
	}

	so := s.mapping.ToGetMonthlyBalances(res)

	s.mencache.SetMonthlyBalanceCache(req, so)

	s.logger.Debug("Monthly balance retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticBycardService) FindYearlyBalanceByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearlyBalance, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyBalance", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyBalance")
	defer span.End()

	year := req.Year

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("FindYearlyBalance called", zap.Int("year", year))

	if data, found := s.mencache.GetYearlyBalanceCache(req); found {
		s.logger.Debug("Cache hit for yearly balance card", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetYearlyBalanceByCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleYearlyBalanceByCardNumberError(err, "FindYearlyBalance", "FAILED_YEARLY_BALANCE_BY_CARD", span, &status)
	}

	so := s.mapping.ToGetYearlyBalances(res)

	s.mencache.SetYearlyBalanceCache(req, so)

	s.logger.Debug("Yearly balance retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticBycardService) FindMonthlyTopupAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyTopupAmount", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyTopupAmount")
	defer span.End()

	cardNumber := req.CardNumber
	year := req.Year

	span.SetAttributes(
		attribute.String("card_number", cardNumber),
		attribute.Int("year", year),
	)

	s.logger.Debug("FindMonthlyTopupAmountByCardNumber called",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
	)

	if data, found := s.mencache.GetMonthlyTopupAmountCache(req); found {
		s.logger.Debug("Cache hit for monthly topup amount card", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetMonthlyTopupAmountByCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleMonthlyTopupAmountByCardNumberError(err, "FindMonthlyTopupAmount", "FAILED_MONTHLY_TOPUP_AMOUNT_BY_CARD", span, &status)
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
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTopupAmount", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTopupAmount")
	defer span.End()

	cardNumber := req.CardNumber
	year := req.Year

	span.SetAttributes(
		attribute.String("card_number", cardNumber),
		attribute.Int("year", year),
	)

	s.logger.Debug("FindYearlyTopupAmountByCardNumber called",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
	)

	if data, found := s.mencache.GetYearlyTopupAmountCache(req); found {
		s.logger.Debug("Cache hit for yearly topup amount card", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetYearlyTopupAmountByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleYearlyTopupAmountByCardNumberError(err, "FindYearlyTopupAmount", "FAILED_YEARLY_TOPUP_AMOUNT_BY_CARD", span, &status)
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.mencache.SetYearlyTopupAmountCache(req, so)

	s.logger.Debug("Yearly topup amount by card number retrieved successfully",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticBycardService) FindMonthlyWithdrawAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyWithdrawAmount", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyWithdrawAmount")
	defer span.End()

	cardNumber := req.CardNumber
	year := req.Year

	span.SetAttributes(
		attribute.String("card_number", cardNumber),
		attribute.Int("year", year),
	)

	s.logger.Debug("FindMonthlyWithdrawAmountByCardNumber called",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
	)

	if data, found := s.mencache.GetMonthlyWithdrawAmountCache(req); found {
		s.logger.Debug("Cache hit for monthly withdraw amount card", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetMonthlyWithdrawAmountByCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleMonthlyWithdrawAmountByCardNumberError(err, "FindMonthlyWithdrawAmount", "FAILED_MONTHLY_WITHDRAW_AMOUNT_BY_CARD", span, &status)
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

	s.mencache.SetMonthlyWithdrawAmountCache(req, so)

	s.logger.Debug("Monthly withdraw amount by card number retrieved successfully",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticBycardService) FindYearlyWithdrawAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyWithdrawAmount", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyWithdrawAmount")
	defer span.End()

	span.SetAttributes(
		attribute.String("card_number", req.CardNumber),
		attribute.Int("year", req.Year),
	)

	cardNumber := req.CardNumber
	year := req.Year

	s.logger.Debug("FindYearlyWithdrawAmountByCardNumber called",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
	)

	if data, found := s.mencache.GetYearlyWithdrawAmountCache(req); found {
		s.logger.Debug("Cache hit for yearly withdraw amount card", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetYearlyWithdrawAmountByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleYearlyWithdrawAmountByCardNumberError(err, "FindYearlyWithdrawAmount", "FAILED_YEARLY_WITHDRAW_AMOUNT_BY_CARD", span, &status)
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.mencache.SetYearlyWithdrawAmountCache(req, so)

	s.logger.Debug("Yearly withdraw amount by card number retrieved successfully",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticBycardService) FindMonthlyTransactionAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyTransactionAmount", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyTransactionAmount")
	defer span.End()

	cardNumber := req.CardNumber
	year := req.Year

	span.SetAttributes(
		attribute.String("card_number", cardNumber),
		attribute.Int("year", year),
	)

	s.logger.Debug("FindMonthlyTransactionAmountByCardNumber called",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
	)

	if data, found := s.mencache.GetMonthlyTransactionAmountCache(req); found {
		s.logger.Debug("Cache hit for monthly transaction amount card", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetMonthlyTransactionAmountByCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleMonthlyTransactionAmountByCardNumberError(err, "FindMonthlyTransactionAmount", "FAILED_MONTHLY_TRANSACTION_AMOUNT_BY_CARD", span, &status)
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

	s.mencache.SetMonthlyTransactionAmountCache(req, so)

	s.logger.Debug("Monthly transaction amount by card number retrieved successfully",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticBycardService) FindYearlyTransactionAmountByCardNumber(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTransactionAmount", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTransactionAmount")
	defer span.End()

	cardNumber := req.CardNumber
	year := req.Year

	span.SetAttributes(
		attribute.String("card_number", cardNumber),
		attribute.Int("year", year),
	)

	s.logger.Debug("FindYearlyTransactionAmountByCardNumber called",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
	)

	if data, found := s.mencache.GetYearlyTransactionAmountCache(req); found {
		s.logger.Debug("Cache hit for yearly transaction amount card", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetYearlyTransactionAmountByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleYearlyTransactionAmountByCardNumberError(err, "FindYearlyTransactionAmount", "FAILED_YEARLY_TRANSACTION_AMOUNT_BY_CARD", span, &status)
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
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyTransferAmountBySender", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyTransferAmountBySender")
	defer span.End()

	cardNumber := req.CardNumber
	year := req.Year

	span.SetAttributes(
		attribute.String("card_number", cardNumber),
		attribute.Int("year", year),
	)

	s.logger.Debug("FindMonthlyTransferAmountBySender called",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
	)

	if data, found := s.mencache.GetMonthlyTransferBySenderCache(req); found {
		s.logger.Debug("Cache hit for monthly transfer amount card", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetMonthlyTransferAmountBySender(req)

	if err != nil {
		return s.errorhandler.HandleMonthlyTransferAmountBySenderError(err, "FindMonthlyTransferAmountBySender", "FAILED_MONTHLY_TRANSFER_AMOUNT_BY_SENDER", span, &status)
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

	s.mencache.SetMonthlyTransferBySenderCache(req, so)

	s.logger.Debug("Monthly transfer sender amount by card number retrieved successfully",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticBycardService) FindYearlyTransferAmountBySender(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTransferAmountBySender", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTransferAmountBySender")
	defer span.End()

	cardNumber := req.CardNumber
	year := req.Year

	span.SetAttributes(
		attribute.String("card_number", cardNumber),
		attribute.Int("year", year),
	)

	s.logger.Debug("FindYearlyTransferAmountBySender called",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
	)

	if data, found := s.mencache.GetYearlyTransferBySenderCache(req); found {
		s.logger.Debug("Cache hit for yearly transfer amount card", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetYearlyTransferAmountBySender(req)
	if err != nil {
		return s.errorhandler.HandleYearlyTransferAmountBySenderError(err, "FindYearlyTransferAmountBySender", "FAILED_YEARLY_TRANSFER_AMOUNT_BY_SENDER", span, &status)
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.mencache.SetYearlyTransferBySenderCache(req, so)

	s.logger.Debug("Yearly transfer sender amount by card number retrieved successfully",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticBycardService) FindMonthlyTransferAmountByReceiver(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyTransferAmountByReceiver", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyTransferAmountByReceiver")
	defer span.End()

	cardNumber := req.CardNumber
	year := req.Year

	span.SetAttributes(
		attribute.String("card_number", cardNumber),
		attribute.Int("year", year),
	)

	s.logger.Debug("FindMonthlyTransferAmountByReceiver called",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
	)

	if data, found := s.mencache.GetMonthlyTransferByReceiverCache(req); found {
		s.logger.Debug("Cache hit for monthly transfer amount card", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetMonthlyTransferAmountByReceiver(req)

	if err != nil {
		return s.errorhandler.HandleMonthlyTransferAmountByReceiverError(err, "FindMonthlyTransferAmountByReceiver", "FAILED_MONTHLY_TRANSFER_AMOUNT_BY_RECEIVER", span, &status)
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

	s.mencache.SetMonthlyTransferByReceiverCache(req, so)

	s.logger.Debug("Monthly transfer receiver amount by card number retrieved successfully",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticBycardService) FindYearlyTransferAmountByReceiver(req *requests.MonthYearCardNumberCard) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTransferAmountByReceiver", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTransferAmountByReceiver")
	defer span.End()

	cardNumber := req.CardNumber
	year := req.Year

	span.SetAttributes(
		attribute.String("card_number", cardNumber),
		attribute.Int("year", year),
	)

	s.logger.Debug("FindYearlyTransferAmountByReceiver called",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
	)

	if data, found := s.mencache.GetYearlyTransferByReceiverCache(req); found {
		s.logger.Debug("Cache hit for yearly transfer amount card", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetYearlyTransferAmountByReceiver(req)
	if err != nil {
		return s.errorhandler.HandleYearlyTransferAmountByReceiverError(err, "FindYearlyTransferAmountByReceiver", "FAILED_YEARLY_TRANSFER_AMOUNT_BY_RECEIVER", span, &status)
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.mencache.SetYearlyTransferByReceiverCache(req, so)

	s.logger.Debug("Yearly transfer receiver amount by card number retrieved successfully",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticBycardService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
