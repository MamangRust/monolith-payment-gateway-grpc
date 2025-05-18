package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-card/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	traceunic "github.com/MamangRust/monolith-payment-gateway-pkg/trace_unic"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors"
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
	trace                   trace.Tracer
	cardDashboardRepository repository.CardDashboardRepository
	logger                  logger.LoggerInterface
	mapping                 responseservice.CardResponseMapper
	requestCounter          *prometheus.CounterVec
	requestDuration         *prometheus.HistogramVec
}

func NewCardDashboardService(
	ctx context.Context,
	cardDashboardRepository repository.CardDashboardRepository,
	logger logger.LoggerInterface,
	mapper responseservice.CardResponseMapper) *cardDashboardService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "card_dashboard_request_count",
		Help: "Number of card dashboard requests CardDashboardService",
	}, []string{"status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "card_dashboard_request_duration_seconds",
		Help:    "Duration of card dashboard requests CardDashboardService",
		Buckets: prometheus.DefBuckets,
	}, []string{"status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	return &cardDashboardService{
		ctx:                     ctx,
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

	_, span := s.trace.Start(s.ctx, "DashboardCard")

	s.logger.Debug("Starting DashboardCard service")

	totalBalance, err := s.cardDashboardRepository.GetTotalBalances()
	if err != nil {
		traceID := traceunic.GenerateTraceID("TOTAL_BALANCE_NOT_FOUND")

		s.logger.Error("Failed to get total balance", zap.String("trace_id", traceID), zap.Error(err))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Total balance not found")
		status = "total_balance_not_found"

		return nil, card_errors.ErrFailedFindTotalBalances
	}

	totalTopup, err := s.cardDashboardRepository.GetTotalTopAmount()
	if err != nil {
		traceID := traceunic.GenerateTraceID("TOTAL_TOP_AMOUNT_NOT_FOUND")

		s.logger.Error("Failed to get total top amount", zap.String("trace_id", traceID), zap.Error(err))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Total top amount not found")
		status = "total_top_amount_not_found"

		return nil, card_errors.ErrFailedFindTotalTopAmount
	}

	totalWithdraw, err := s.cardDashboardRepository.GetTotalWithdrawAmount()
	if err != nil {
		traceID := traceunic.GenerateTraceID("TOTAL_WITHDRAW_AMOUNT_NOT_FOUND")

		s.logger.Error("Failed to get total withdraw amount", zap.String("trace_id", traceID), zap.Error(err))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Total withdraw amount not found")
		status = "total_withdraw_amount_not_found"

		return nil, card_errors.ErrFailedFindTotalWithdrawAmount
	}

	totalTransaction, err := s.cardDashboardRepository.GetTotalTransactionAmount()
	if err != nil {
		traceID := traceunic.GenerateTraceID("TOTAL_TRANSACTION_AMOUNT_NOT_FOUND")

		s.logger.Error("Failed to get total transaction amount", zap.String("trace_id", traceID), zap.Error(err))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Total transaction amount not found")
		status = "total_transaction_amount_not_found"

		return nil, card_errors.ErrFailedFindTotalTransactionAmount
	}

	totalTransfer, err := s.cardDashboardRepository.GetTotalTransferAmount()
	if err != nil {
		traceID := traceunic.GenerateTraceID("TOTAL_TRANSFER_AMOUNT_NOT_FOUND")

		s.logger.Error("Failed to get total transfer amount", zap.String("trace_id", traceID), zap.Error(err))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Total transfer amount not found")
		status = "total_transfer_amount_not_found"

		return nil, card_errors.ErrFailedFindTotalTransferAmount
	}

	s.logger.Debug("Completed DashboardCard service",
		zap.Int("total_balance", int(*totalBalance)),
		zap.Int("total_topup", int(*totalTopup)),
		zap.Int("total_withdraw", int(*totalWithdraw)),
		zap.Int("total_transaction", int(*totalTransaction)),
		zap.Int("total_transfer", int(*totalTransfer)),
	)

	return &response.DashboardCard{
		TotalBalance:     totalBalance,
		TotalTopup:       totalTopup,
		TotalWithdraw:    totalWithdraw,
		TotalTransaction: totalTransaction,
		TotalTransfer:    totalTransfer,
	}, nil
}

func (s *cardDashboardService) DashboardCardCardNumber(cardNumber string) (*response.DashboardCardCardNumber, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DashboardCardCardNumber", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "DashboardCardCardNumber")
	span.SetAttributes(attribute.String("card_number", cardNumber))

	s.logger.Debug("Starting DashboardCardCardNumber service",
		zap.String("card_number", cardNumber),
	)

	totalBalance, err := s.cardDashboardRepository.GetTotalBalanceByCardNumber(cardNumber)
	if err != nil {
		traceID := traceunic.GenerateTraceID("TOTAL_BALANCE_NOT_FOUND")

		s.logger.Error("Failed to get total balance", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Total balance not found")
		status = "total_balance_not_found"

		return nil, card_errors.ErrFailedFindTotalBalanceByCard
	}

	totalTopup, err := s.cardDashboardRepository.GetTotalTopupAmountByCardNumber(cardNumber)
	if err != nil {
		traceID := traceunic.GenerateTraceID("TOTAL_TOP_AMOUNT_NOT_FOUND")

		s.logger.Error("Failed to get total top amount", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Total top amount not found")
		status = "total_top_amount_not_found"

		return nil, card_errors.ErrFailedFindTotalTopupAmountByCard
	}

	totalWithdraw, err := s.cardDashboardRepository.GetTotalWithdrawAmountByCardNumber(cardNumber)
	if err != nil {
		traceID := traceunic.GenerateTraceID("TOTAL_WITHDRAW_AMOUNT_NOT_FOUND")

		s.logger.Error("Failed to get total withdraw amount", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Total withdraw amount not found")
		status = "total_withdraw_amount_not_found"

		return nil, card_errors.ErrFailedFindTotalWithdrawAmountByCard
	}

	totalTransaction, err := s.cardDashboardRepository.GetTotalTransactionAmountByCardNumber(cardNumber)
	if err != nil {
		traceID := traceunic.GenerateTraceID("TOTAL_TRANSACTION_AMOUNT_NOT_FOUND")

		s.logger.Error("Failed to get total transaction amount", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Total transaction amount not found")
		status = "total_transaction_amount_not_found"

		return nil, card_errors.ErrFailedFindTotalTransactionAmountByCard
	}

	totalTransferSent, err := s.cardDashboardRepository.GetTotalTransferAmountBySender(cardNumber)
	if err != nil {
		traceID := traceunic.GenerateTraceID("TOTAL_TRANSFER_AMOUNT_NOT_FOUND")

		s.logger.Error("Failed to get total transfer amount sent by card",
			zap.String("card_number", cardNumber),
			zap.Error(err),
		)
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Total transfer amount not found")
		status = "total_transfer_amount_not_found"

		return nil, card_errors.ErrFailedFindTotalTransferAmountBySender
	}

	totalTransferReceived, err := s.cardDashboardRepository.GetTotalTransferAmountByReceiver(cardNumber)
	if err != nil {
		traceID := traceunic.GenerateTraceID("TOTAL_TRANSFER_AMOUNT_NOT_FOUND")

		s.logger.Error("Failed to get total transfer amount received by card",
			zap.String("card_number", cardNumber),
			zap.Error(err),
		)
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Total transfer amount not found")
		status = "total_transfer_amount_not_found"

		return nil, card_errors.ErrFailedFindTotalTransferAmountByReceiver
	}

	s.logger.Debug("Completed DashboardCardCardNumber service",
		zap.String("card_number", cardNumber),
		zap.Int("total_balance", int(*totalBalance)),
		zap.Int("total_topup", int(*totalTopup)),
		zap.Int("total_withdraw", int(*totalWithdraw)),
		zap.Int("total_transaction", int(*totalTransaction)),
		zap.Int("total_transfer_sent", int(*totalTransferSent)),
		zap.Int("total_transfer_received", int(*totalTransferReceived)),
	)

	return &response.DashboardCardCardNumber{
		TotalBalance:          totalBalance,
		TotalTopup:            totalTopup,
		TotalWithdraw:         totalWithdraw,
		TotalTransaction:      totalTransaction,
		TotalTransferSend:     totalTransferSent,
		TotalTransferReceiver: totalTransferReceived,
	}, nil
}

func (s *cardDashboardService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
