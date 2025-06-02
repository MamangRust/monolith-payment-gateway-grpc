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
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DashboardCard", status, startTime)
	}()

	if data, found := s.mencache.GetDashboardCardCache(); found {
		s.logger.Debug("DashboardCard cache hit")
		return data, nil
	}
	s.logger.Debug("DashboardCard cache miss")

	_, span := s.trace.Start(s.ctx, "DashboardCard")

	s.logger.Debug("Starting DashboardCard service")

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

	s.logger.Debug("Completed DashboardCard service",
		zap.Int("total_balance", int(*totalBalance)),
		zap.Int("total_topup", int(*totalTopup)),
		zap.Int("total_withdraw", int(*totalWithdraw)),
		zap.Int("total_transaction", int(*totalTransaction)),
		zap.Int("total_transfer", int(*totalTransfer)),
	)

	return result, nil
}

func (s *cardDashboardService) DashboardCardCardNumber(cardNumber string) (*response.DashboardCardCardNumber, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DashboardCardCardNumber", status, startTime)
	}()

	if data, found := s.mencache.GetDashboardCardCardNumberCache(cardNumber); found {
		s.logger.Debug("DashboardCardCardNumber cache hit", zap.String("card_number", cardNumber))
		return data, nil
	}
	s.logger.Debug("DashboardCardCardNumber cache miss", zap.String("card_number", cardNumber))

	_, span := s.trace.Start(s.ctx, "DashboardCardCardNumber")
	span.SetAttributes(attribute.String("card_number", cardNumber))

	s.logger.Debug("Starting DashboardCardCardNumber service",
		zap.String("card_number", cardNumber),
	)

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

	s.logger.Debug("Completed DashboardCardCardNumber service",
		zap.String("card_number", cardNumber),
		zap.Int("total_balance", int(*totalBalance)),
		zap.Int("total_topup", int(*totalTopup)),
		zap.Int("total_withdraw", int(*totalWithdraw)),
		zap.Int("total_transaction", int(*totalTransaction)),
		zap.Int("total_transfer_sent", int(*totalTransferSent)),
		zap.Int("total_transfer_received", int(*totalTransferReceived)),
	)

	return result, nil
}

func (s *cardDashboardService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
