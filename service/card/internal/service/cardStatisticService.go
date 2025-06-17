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
	"go.opentelemetry.io/otel/codes"
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
	const method = "FindMonthlyBalance"
	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyBalanceCache(year); found {
		logSuccess("Monthly balance cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetMonthlyBalance(year)
	if err != nil {
		return s.errorhandler.HandleMonthlyBalanceError(err, method, "FAILED_MONTHLY_BALANCE", span, &status, zap.Error(err))
	}

	so := s.mapping.ToGetMonthlyBalances(res)

	s.mencache.SetMonthlyBalanceCache(year, so)

	logSuccess("Monthly balance retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticService) FindYearlyBalance(year int) ([]*response.CardResponseYearlyBalance, *response.ErrorResponse) {
	const method = "FindYearlyBalance"
	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyBalanceCache(year); found {
		logSuccess("Yearly balance cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetYearlyBalance(year)
	if err != nil {
		return s.errorhandler.HandleYearlyBalanceError(err, method, "FAILED_YEARLY_BALANCE", span, &status, zap.Error(err))
	}

	so := s.mapping.ToGetYearlyBalances(res)

	s.mencache.SetYearlyBalanceCache(year, so)

	logSuccess("Yearly balance retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticService) FindMonthlyTopupAmount(year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	const method = "FindMonthlyTopupAmount"
	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	s.logger.Debug("FindMonthlyTopupAmount called", zap.Int("year", year))

	if data, found := s.mencache.GetMonthlyTopupAmountCache(year); found {
		logSuccess("Monthly topup amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetMonthlyTopupAmount(year)

	if err != nil {
		return s.errorhandler.HandleMonthlyTopupAmountError(err, method, "FAILED_MONTHLY_TOPUP_AMOUNT", span, &status, zap.Error(err))
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

	s.mencache.SetMonthlyTopupAmountCache(year, so)

	logSuccess("Monthly topup amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticService) FindYearlyTopupAmount(year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	const method = "FindYearlyTopupAmount"
	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyTopupAmountCache(year); found {
		logSuccess("Yearly topup amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetYearlyTopupAmount(year)

	if err != nil {
		return s.errorhandler.HandleYearlyTopupAmountError(err, method, "FAILED_YEARLY_TOPUP_AMOUNT", span, &status, zap.Error(err))
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.mencache.SetYearlyTopupAmountCache(year, so)

	logSuccess("Yearly topup amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticService) FindMonthlyWithdrawAmount(year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	const method = "FindMonthlyWithdrawAmount"
	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyWithdrawAmountCache(year); found {
		logSuccess("Monthly withdraw amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetMonthlyWithdrawAmount(year)

	if err != nil {
		return s.errorhandler.HandleMonthlyWithdrawAmountError(err, method, "FAILED_MONTHLY_WITHDRAW_AMOUNT", span, &status, zap.Error(err))
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

	s.mencache.SetMonthlyWithdrawAmountCache(year, so)

	logSuccess("Monthly withdraw amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticService) FindYearlyWithdrawAmount(year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	const method = "FindYearlyWithdrawAmount"
	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	res, err := s.cardStatisticRepository.GetYearlyWithdrawAmount(year)

	if err != nil {
		return s.errorhandler.HandleYearlyWithdrawAmountError(err, method, "FAILED_YEARLY_WITHDRAW_AMOUNT", span, &status, zap.Error(err))
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.mencache.SetYearlyWithdrawAmountCache(year, so)

	logSuccess("Yearly withdraw amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticService) FindMonthlyTransactionAmount(year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	const method = "FindMonthlyTransactionAmount"
	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyTransactionAmountCache(year); found {
		logSuccess("Monthly transaction amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetMonthlyTransactionAmount(year)

	if err != nil {
		return s.errorhandler.HandleMonthlyTransactionAmountError(err, method, "FAILED_MONTHLY_TRANSACTION_AMOUNT", span, &status, zap.Error(err))
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

	s.mencache.SetMonthlyTransactionAmountCache(year, so)

	logSuccess("Monthly transaction amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticService) FindYearlyTransactionAmount(year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	const method = "FindYearlyTransactionAmount"
	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyTransactionAmountCache(year); found {
		logSuccess("Yearly transaction amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetYearlyTransactionAmount(year)

	if err != nil {
		return s.errorhandler.HandleYearlyTransactionAmountError(err, method, "FAILED_YEARLY_TRANSACTION_AMOUNT", span, &status, zap.Error(err))
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.mencache.SetYearlyTransactionAmountCache(year, so)

	logSuccess("Yearly transaction amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticService) FindMonthlyTransferAmountSender(year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	const method = "FindMonthlyTransferAmountSender"
	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyTransferAmountSenderCache(year); found {
		logSuccess("Monthly transfer sender amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetMonthlyTransferAmountSender(year)
	if err != nil {
		return s.errorhandler.HandleMonthlyTransferAmountSenderError(err, method, "FAILED_MONTHLY_TRANSFER_AMOUNT_SENDER", span, &status, zap.Error(err))
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

	s.mencache.SetMonthlyTransferAmountSenderCache(year, so)

	logSuccess("Monthly transfer sender amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticService) FindYearlyTransferAmountSender(year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	const method = "FindYearlyTransferAmountSender"
	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyTransferAmountSenderCache(year); found {
		logSuccess("Yearly transfer sender amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetYearlyTransferAmountSender(year)

	if err != nil {
		return s.errorhandler.HandleYearlyTransferAmountSenderError(err, method, "FAILED_YEARLY_TRANSFER_AMOUNT_SENDER", span, &status, zap.Error(err))
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.mencache.SetYearlyTransferAmountSenderCache(year, so)

	logSuccess("Yearly transfer sender amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticService) FindMonthlyTransferAmountReceiver(year int) ([]*response.CardResponseMonthAmount, *response.ErrorResponse) {
	const method = "FindMonthlyTransferAmountReceiver"
	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyTransferAmountReceiverCache(year); found {
		logSuccess("Monthly transfer receiver amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetMonthlyTransferAmountReceiver(year)

	if err != nil {
		return s.errorhandler.HandleMonthlyTransferAmountReceiverError(err, method, "FAILED_MONTHLY_TRANSFER_AMOUNT_RECEIVER", span, &status, zap.Error(err))
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

	s.mencache.SetMonthlyTransferAmountReceiverCache(year, so)

	logSuccess("Monthly transfer receiver amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticService) FindYearlyTransferAmountReceiver(year int) ([]*response.CardResponseYearAmount, *response.ErrorResponse) {
	const method = "FindYearlyTransferAmountReceiver"
	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyTransferAmountReceiverCache(year); found {
		logSuccess("Yearly transfer receiver amount cache hit", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cardStatisticRepository.GetYearlyTransferAmountReceiver(year)

	if err != nil {
		return s.errorhandler.HandleYearlyTransferAmountReceiverError(err, method, "FAILED_YEARLY_TRANSFER_AMOUNT_RECEIVER", span, &status, zap.Error(err))
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.mencache.SetYearlyTransferAmountReceiverCache(year, so)

	logSuccess("Yearly transfer receiver amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *cardStatisticService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
