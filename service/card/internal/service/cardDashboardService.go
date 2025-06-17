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

type cardDashboardService struct {
	ctx                     context.Context
	errorhandler            errorhandler.CardDashboardErrorHandler
	mencache                mencache.CardDashboardCache
	trace                   trace.Tracer
	cardDashboardRepository repository.CardDashboardRepository
	logger                  logger.LoggerInterface
	mapping                 responseservice.CardResponseMapper
	requestCounter          *prometheus.CounterVec
	requestDuration         *prometheus.HistogramVec
}

func NewCardDashboardService(
	ctx context.Context,
	errorhandler errorhandler.CardDashboardErrorHandler,
	mencache mencache.CardDashboardCache,
	cardDashboardRepository repository.CardDashboardRepository,
	logger logger.LoggerInterface,
	mapper responseservice.CardResponseMapper) *cardDashboardService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "card_dashboard_request_count",
		Help: "Number of card dashboard requests CardDashboardService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "card_dashboard_request_duration_seconds",
		Help:    "Duration of card dashboard requests CardDashboardService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	return &cardDashboardService{
		ctx:                     ctx,
		errorhandler:            errorhandler,
		mencache:                mencache,
		trace:                   otel.Tracer("card-dashboard-service"),
		cardDashboardRepository: cardDashboardRepository,
		logger:                  logger,
		mapping:                 mapper,
		requestCounter:          requestCounter,
		requestDuration:         requestDuration,
	}
}

func (s *cardDashboardService) DashboardCard() (*response.DashboardCard, *response.ErrorResponse) {
	const method = "DashboardCard"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetDashboardCardCache(); found {
		s.logger.Debug("DashboardCard cache hit")
		return data, nil
	}

	totalBalance, err := s.cardDashboardRepository.GetTotalBalances()
	if err != nil {
		return s.errorhandler.HandleTotalBalanceError(err, "DashboardCard", "FAILED_FIND_TOTAL_BALANCE", span, &status, zap.Error(err))
	}

	totalTopup, err := s.cardDashboardRepository.GetTotalTopAmount()
	if err != nil {
		return s.errorhandler.HandleTotalTopupAmountError(err, "DashboardCard", "FAILED_FIND_TOTAL_TOPUP", span, &status, zap.Error(err))
	}

	totalWithdraw, err := s.cardDashboardRepository.GetTotalWithdrawAmount()
	if err != nil {
		return s.errorhandler.HandleTotalWithdrawAmountError(err, "DashboardCard", "FAILED_FIND_TOTAL_WITHDRAW", span, &status, zap.Error(err))
	}

	totalTransaction, err := s.cardDashboardRepository.GetTotalTransactionAmount()
	if err != nil {
		return s.errorhandler.HandleTotalTransactionAmountError(err, "DashboardCard", "FAILED_FIND_TOTAL_TRANSACTION", span, &status, zap.Error(err))
	}

	totalTransfer, err := s.cardDashboardRepository.GetTotalTransferAmount()
	if err != nil {
		return s.errorhandler.HandleTotalTransferAmountError(err, "DashboardCard", "FAILED_FIND_TOTAL_TRANSFER", span, &status, zap.Error(err))
	}

	result := &response.DashboardCard{
		TotalBalance:     totalBalance,
		TotalTopup:       totalTopup,
		TotalWithdraw:    totalWithdraw,
		TotalTransaction: totalTransaction,
		TotalTransfer:    totalTransfer,
	}

	s.mencache.SetDashboardCardCache(result)

	logSuccess("Success find dashboard card", zap.Bool("success", true))

	return result, nil
}

func (s *cardDashboardService) DashboardCardCardNumber(cardNumber string) (*response.DashboardCardCardNumber, *response.ErrorResponse) {
	const method = "DashboardCardCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetDashboardCardCardNumberCache(cardNumber); found {
		s.logger.Debug("DashboardCardCardNumber cache hit", zap.String("card_number", cardNumber))
		return data, nil
	}
	s.logger.Debug("DashboardCardCardNumber cache miss", zap.String("card_number", cardNumber))

	totalBalance, err := s.cardDashboardRepository.GetTotalBalanceByCardNumber(cardNumber)
	if err != nil {
		return s.errorhandler.HandleTotalBalanceCardNumberError(err, "DashboardCardCardNumber", "FAILED_FIND_TOTAL_BALANCE_BY_CARD", span, &status, zap.Error(err))
	}

	totalTopup, err := s.cardDashboardRepository.GetTotalTopupAmountByCardNumber(cardNumber)
	if err != nil {
		return s.errorhandler.HandleTotalTopupAmountCardNumberError(err, "DashboardCardCardNumber", "FAILED_FIND_TOTAL_TOPUP_BY_CARD", span, &status, zap.Error(err))
	}

	totalWithdraw, err := s.cardDashboardRepository.GetTotalWithdrawAmountByCardNumber(cardNumber)
	if err != nil {
		return s.errorhandler.HandleTotalWithdrawAmountCardNumberError(err, "DashboardCardCardNumber", "FAILED_FIND_TOTAL_WITHDRAW_BY_CARD", span, &status, zap.Error(err))
	}

	totalTransaction, err := s.cardDashboardRepository.GetTotalTransactionAmountByCardNumber(cardNumber)
	if err != nil {
		return s.errorhandler.HandleTotalTransactionAmountCardNumberError(err, "DashboardCardCardNumber", "FAILED_FIND_TOTAL_TRANSACTION_BY_CARD", span, &status, zap.Error(err))
	}

	totalTransferSent, err := s.cardDashboardRepository.GetTotalTransferAmountBySender(cardNumber)
	if err != nil {
		return s.errorhandler.HandleTotalTransferAmountBySender(err, "DashboardCardCardNumber", "FAILED_FIND_TOTAL_TRANSFER_BY_CARD", span, &status, zap.Error(err))
	}

	totalTransferReceived, err := s.cardDashboardRepository.GetTotalTransferAmountByReceiver(cardNumber)
	if err != nil {
		return s.errorhandler.HandleTotalTransferAmountByReceiver(err, "DashboardCardCardNumber", "FAILED_FIND_TOTAL_TRANSFER_BY_CARD", span, &status, zap.Error(err))
	}

	result := &response.DashboardCardCardNumber{
		TotalBalance:          totalBalance,
		TotalTopup:            totalTopup,
		TotalWithdraw:         totalWithdraw,
		TotalTransaction:      totalTransaction,
		TotalTransferSend:     totalTransferSent,
		TotalTransferReceiver: totalTransferReceived,
	}

	s.mencache.SetDashboardCardCardNumberCache(cardNumber, result)

	logSuccess("Success find dashboard card card number", zap.Bool("success", true))

	return result, nil
}

func (s *cardDashboardService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

	s.logger.Debug("Start: " + method)

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

func (s *cardDashboardService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
