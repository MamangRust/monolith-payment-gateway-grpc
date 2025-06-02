package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-card/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-card/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-card/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type cardStatisticService struct {
	ctx                     context.Context
	errorhandler            errorhandler.CardStatisticErrorHandler
	mencache                mencache.CardStatisticCache
	trace                   trace.Tracer
	cardStatisticRepository repository.CardStatisticRepository
	logger                  logger.LoggerInterface
	mapping                 responseservice.CardResponseMapper
	requestCounter          *prometheus.CounterVec
	requestDuration         *prometheus.HistogramVec
}

func NewCardStatisticService(
	ctx context.Context,
	errorhandler errorhandler.CardStatisticErrorHandler,
	mencache mencache.CardStatisticCache,
	cardStatisticRepository repository.CardStatisticRepository, logger logger.LoggerInterface, mapper responseservice.CardResponseMapper) *cardStatisticService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "card_statistic_request_count",
		Help: "Number of card statistic requests CardStatisticService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "card_statistic_request_duration_seconds",
		Help:    "Duration of card statistic requests CardStatisticService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	return &cardStatisticService{
		ctx:                     ctx,
		errorhandler:            errorhandler,
		mencache:                mencache,
		trace:                   otel.Tracer("card-statistic-service"),
		cardStatisticRepository: cardStatisticRepository,
		logger:                  logger,
		mapping:                 mapper,
		requestCounter:          requestCounter,
		requestDuration:         requestDuration,
	}
}

func (s *cardStatisticService) FindMonthlyBalance(year int) ([]*response.CardResponseMonthBalance, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyBalance", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyBalance")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("FindMonthlyBalance called", zap.Int("year", year))

	if data, found := s.mencache.GetMonthlyBalanceCache(year); found {
		s.logger.Debug("Monthly balance cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetMonthlyBalance(year)
	if err != nil {
		return s.errorhandler.HandleMonthlyBalanceError(err, "FindMonthlyBalance", "FAILED_MONTHLY_BALANCE", span, &status)
	}

	so := s.mapping.ToGetMonthlyBalances(res)

	s.mencache.SetMonthlyBalanceCache(year, so)

	s.logger.Debug("Monthly balance retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticService) FindYearlyBalance(year int) ([]*response.CardResponseYearlyBalance, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyBalance", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyBalance")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("FindYearlyBalance called", zap.Int("year", year))

	if data, found := s.mencache.GetYearlyBalanceCache(year); found {
		s.logger.Debug("Yearly balance cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetYearlyBalance(year)
	if err != nil {
		return s.errorhandler.HandleYearlyBalanceError(err, "FindYearlyBalance", "FAILED_YEARLY_BALANCE", span, &status)
	}

	so := s.mapping.ToGetYearlyBalances(res)

	s.mencache.SetYearlyBalanceCache(year, so)

	s.logger.Debug("Yearly balance retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticService) FindMonthlyTopupAmount(year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyTopupAmount", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyTopupAmount")
	defer span.End()

	s.logger.Debug("FindMonthlyTopupAmount called", zap.Int("year", year))

	if data, found := s.mencache.GetMonthlyTopupAmountCache(year); found {
		s.logger.Debug("Monthly topup amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetMonthlyTopupAmount(year)

	if err != nil {
		return s.errorhandler.HandleMonthlyTopupAmountError(err, "FindMonthlyTopupAmount", "FAILED_MONTHLY_TOPUP_AMOUNT", span, &status)
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

	s.mencache.SetMonthlyTopupAmountCache(year, so)

	s.logger.Debug("Monthly topup amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticService) FindYearlyTopupAmount(year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTopupAmount", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTopupAmount")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("FindYearlyTopupAmount called", zap.Int("year", year))

	if data, found := s.mencache.GetYearlyTopupAmountCache(year); found {
		s.logger.Debug("Yearly topup amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetYearlyTopupAmount(year)

	if err != nil {
		return s.errorhandler.HandleYearlyTopupAmountError(err, "FindYearlyTopupAmount", "FAILED_YEARLY_TOPUP_AMOUNT", span, &status)
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.mencache.SetYearlyTopupAmountCache(year, so)

	s.logger.Debug("Yearly topup amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticService) FindMonthlyWithdrawAmount(year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyWithdrawAmount", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyWithdrawAmount")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("FindMonthlyWithdrawAmount called", zap.Int("year", year))

	if data, found := s.mencache.GetMonthlyWithdrawAmountCache(year); found {
		s.logger.Debug("Monthly withdraw amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetMonthlyWithdrawAmount(year)

	if err != nil {
		return s.errorhandler.HandleMonthlyWithdrawAmountError(err, "FindMonthlyWithdrawAmount", "FAILED_MONTHLY_WITHDRAW_AMOUNT", span, &status)
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

	s.mencache.SetMonthlyWithdrawAmountCache(year, so)

	s.logger.Debug("Monthly withdraw amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticService) FindYearlyWithdrawAmount(year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyWithdrawAmount", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyWithdrawAmount")

	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("FindYearlyWithdrawAmount called", zap.Int("year", year))

	res, err := s.cardStatisticRepository.GetYearlyWithdrawAmount(year)

	if err != nil {
		return s.errorhandler.HandleYearlyWithdrawAmountError(err, "FindYearlyWithdrawAmount", "FAILED_YEARLY_WITHDRAW_AMOUNT", span, &status)
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.mencache.SetYearlyWithdrawAmountCache(year, so)

	s.logger.Debug("Yearly withdraw amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticService) FindMonthlyTransactionAmount(year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyTransactionAmount", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyTransactionAmount")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("FindMonthlyTransactionAmount called", zap.Int("year", year))

	if data, found := s.mencache.GetMonthlyTransactionAmountCache(year); found {
		s.logger.Debug("Monthly transaction amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetMonthlyTransactionAmount(year)

	if err != nil {
		return s.errorhandler.HandleMonthlyTransactionAmountError(err, "FindMonthlyTransactionAmount", "FAILED_MONTHLY_TRANSACTION_AMOUNT", span, &status)
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

	s.mencache.SetMonthlyTransactionAmountCache(year, so)

	s.logger.Debug("Monthly transaction amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticService) FindYearlyTransactionAmount(year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	s.logger.Debug("FindYearlyTransactionAmount called", zap.Int("year", year))

	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTransactionAmount", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTransactionAmount")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	if data, found := s.mencache.GetYearlyTransactionAmountCache(year); found {
		s.logger.Debug("Yearly transaction amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetYearlyTransactionAmount(year)

	if err != nil {
		return s.errorhandler.HandleYearlyTransactionAmountError(err, "FindYearlyTransactionAmount", "FAILED_YEARLY_TRANSACTION_AMOUNT", span, &status)
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.mencache.SetYearlyTransactionAmountCache(year, so)

	s.logger.Debug("Yearly transaction amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticService) FindMonthlyTransferAmountSender(year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyTransferAmountSender", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyTransferAmountSender")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("FindMonthlyTransferAmountSender called", zap.Int("year", year))

	if data, found := s.mencache.GetMonthlyTransferAmountSenderCache(year); found {
		s.logger.Debug("Monthly transfer sender amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetMonthlyTransferAmountSender(year)
	if err != nil {
		return s.errorhandler.HandleMonthlyTransferAmountSenderError(err, "FindMonthlyTransferAmountSender", "FAILED_MONTHLY_TRANSFER_AMOUNT_SENDER", span, &status)
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

	s.mencache.SetMonthlyTransferAmountSenderCache(year, so)

	s.logger.Debug("Monthly transfer sender amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticService) FindYearlyTransferAmountSender(year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTransferAmountSender", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTransferAmountSender")
	defer span.End()

	s.logger.Debug("FindYearlyTransferAmountSender called", zap.Int("year", year))

	if data, found := s.mencache.GetYearlyTransferAmountSenderCache(year); found {
		s.logger.Debug("Yearly transfer sender amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetYearlyTransferAmountSender(year)

	if err != nil {
		return s.errorhandler.HandleYearlyTransferAmountSenderError(err, "FindYearlyTransferAmountSender", "FAILED_YEARLY_TRANSFER_AMOUNT_SENDER", span, &status)
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.mencache.SetYearlyTransferAmountSenderCache(year, so)

	s.logger.Debug("Yearly transfer sender amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticService) FindMonthlyTransferAmountReceiver(year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyTransferAmountReceiver", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyTransferAmountReceiver")
	defer span.End()

	s.logger.Debug("FindMonthlyTransferAmountReceiver called", zap.Int("year", year))

	if data, found := s.mencache.GetMonthlyTransferAmountReceiverCache(year); found {
		s.logger.Debug("Monthly transfer receiver amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetMonthlyTransferAmountReceiver(year)

	if err != nil {
		return s.errorhandler.HandleMonthlyTransferAmountReceiverError(err, "FindMonthlyTransferAmountReceiver", "FAILED_MONTHLY_TRANSFER_AMOUNT_RECEIVER", span, &status)
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

	s.mencache.SetMonthlyTransferAmountReceiverCache(year, so)

	s.logger.Debug("Monthly transfer receiver amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticService) FindYearlyTransferAmountReceiver(year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTransferAmountReceiver", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTransferAmountReceiver")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("FindYearlyTransferAmountReceiver called", zap.Int("year", year))

	if data, found := s.mencache.GetYearlyTransferAmountReceiverCache(year); found {
		s.logger.Debug("Yearly transfer receiver amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetYearlyTransferAmountReceiver(year)

	if err != nil {
		return s.errorhandler.HandleYearlyTransferAmountReceiverError(err, "FindYearlyTransferAmountReceiver", "FAILED_YEARLY_TRANSFER_AMOUNT_RECEIVER", span, &status)
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.mencache.SetYearlyTransferAmountReceiverCache(year, so)

	s.logger.Debug("Yearly transfer receiver amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
