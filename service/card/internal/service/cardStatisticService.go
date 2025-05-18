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

type cardStatisticService struct {
	ctx                     context.Context
	trace                   trace.Tracer
	cardStatisticRepository repository.CardStatisticRepository
	logger                  logger.LoggerInterface
	mapping                 responseservice.CardResponseMapper
	requestCounter          *prometheus.CounterVec
	requestDuration         *prometheus.HistogramVec
}

func NewCardStatisticService(
	ctx context.Context,
	cardStatisticRepository repository.CardStatisticRepository, logger logger.LoggerInterface, mapper responseservice.CardResponseMapper) *cardStatisticService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "card_statistic_request_count",
		Help: "Number of card statistic requests CardStatisticService",
	}, []string{"status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "card_statistic_request_duration_seconds",
		Help:    "Duration of card statistic requests CardStatisticService",
		Buckets: prometheus.DefBuckets,
	}, []string{"status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	return &cardStatisticService{
		ctx:                     ctx,
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

	res, err := s.cardStatisticRepository.GetMonthlyBalance(year)
	if err != nil {
		traceID := traceunic.GenerateTraceID("MONTHLY_BALANCE_NOT_FOUND")

		s.logger.Error("Failed to retrieve monthly balance", zap.String("trace_id", traceID), zap.Int("year", year), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Monthly balance not found")
		status = "monthly_balance_not_found"

		return nil, card_errors.ErrFailedFindMonthlyBalance
	}

	so := s.mapping.ToGetMonthlyBalances(res)

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

	res, err := s.cardStatisticRepository.GetYearlyBalance(year)
	if err != nil {
		traceID := traceunic.GenerateTraceID("YEARLY_BALANCE_NOT_FOUND")

		s.logger.Error("Failed to retrieve yearly balance", zap.String("trace_id", traceID), zap.Int("year", year), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Yearly balance not found")
		status = "yearly_balance_not_found"

		return nil, card_errors.ErrFailedFindYearlyBalance
	}

	so := s.mapping.ToGetYearlyBalances(res)

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

	res, err := s.cardStatisticRepository.GetMonthlyTopupAmount(year)

	if err != nil {
		traceID := traceunic.GenerateTraceID("MONTHLY_TOPUP_AMOUNT_NOT_FOUND")

		s.logger.Error("Failed to retrieve monthly topup amount", zap.String("trace_id", traceID), zap.Int("year", year), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Monthly topup amount not found")
		status = "monthly_topup_amount_not_found"

		return nil, card_errors.ErrFailedFindMonthlyTopupAmount
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

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

	res, err := s.cardStatisticRepository.GetYearlyTopupAmount(year)

	if err != nil {
		traceID := traceunic.GenerateTraceID("YEARLY_TOPUP_AMOUNT_NOT_FOUND")

		s.logger.Error("Failed to retrieve yearly topup amount", zap.String("trace_id", traceID), zap.Int("year", year), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Yearly topup amount not found")
		status = "yearly_topup_amount_not_found"

		return nil, card_errors.ErrFailedFindYearlyTopupAmount
	}

	so := s.mapping.ToGetYearlyAmounts(res)

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

	res, err := s.cardStatisticRepository.GetMonthlyWithdrawAmount(year)

	if err != nil {
		traceID := traceunic.GenerateTraceID("MONTHLY_WITHDRAW_AMOUNT_NOT_FOUND")

		s.logger.Error("Failed to retrieve monthly withdraw amount", zap.String("trace_id", traceID), zap.Int("year", year), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Monthly withdraw amount not found")
		status = "monthly_withdraw_amount_not_found"

		return nil, card_errors.ErrFailedFindMonthlyWithdrawAmount
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

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
		traceID := traceunic.GenerateTraceID("YEARLY_WITHDRAW_AMOUNT_NOT_FOUND")

		s.logger.Error("Failed to retrieve yearly withdraw amount", zap.String("trace_id", traceID), zap.Int("year", year), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Yearly withdraw amount not found")
		status = "yearly_withdraw_amount_not_found"

		return nil, card_errors.ErrFailedFindYearlyWithdrawAmount
	}

	so := s.mapping.ToGetYearlyAmounts(res)

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

	res, err := s.cardStatisticRepository.GetMonthlyTransactionAmount(year)

	if err != nil {
		traceID := traceunic.GenerateTraceID("MONTHLY_TRANSACTION_AMOUNT_NOT_FOUND")

		s.logger.Error("Failed to retrieve monthly transaction amount", zap.String("trace_id", traceID), zap.Int("year", year), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Monthly transaction amount not found")
		status = "monthly_transaction_amount_not_found"

		return nil, card_errors.ErrFailedFindMonthlyTransactionAmount
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

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

	res, err := s.cardStatisticRepository.GetYearlyTransactionAmount(year)

	if err != nil {
		traceID := traceunic.GenerateTraceID("YEARLY_TRANSACTION_AMOUNT_NOT_FOUND")

		s.logger.Error("Failed to retrieve yearly transaction amount", zap.String("trace_id", traceID), zap.Int("year", year), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Yearly transaction amount not found")
		status = "yearly_transaction_amount_not_found"

		return nil, card_errors.ErrFailedFindYearlyTransactionAmount
	}

	so := s.mapping.ToGetYearlyAmounts(res)

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

	res, err := s.cardStatisticRepository.GetMonthlyTransferAmountSender(year)
	if err != nil {
		traceID := traceunic.GenerateTraceID("MONTHLY_TRANSFER_AMOUNT_SENDER_NOT_FOUND")

		s.logger.Error("Failed to retrieve monthly transfer sender amount", zap.String("trace_id", traceID), zap.Int("year", year), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Monthly transfer sender amount not found")
		status = "monthly_transfer_amount_sender_not_found"

		return nil, card_errors.ErrFailedFindMonthlyTransferAmountSender
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

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

	res, err := s.cardStatisticRepository.GetYearlyTransferAmountSender(year)

	if err != nil {
		traceID := traceunic.GenerateTraceID("YEARLY_TRANSFER_AMOUNT_SENDER_NOT_FOUND")

		s.logger.Error("Failed to retrieve yearly transfer sender amount", zap.String("trace_id", traceID), zap.Int("year", year), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Yearly transfer sender amount not found")
		status = "yearly_transfer_amount_sender_not_found"

		return nil, card_errors.ErrFailedFindYearlyTransferAmountSender
	}

	so := s.mapping.ToGetYearlyAmounts(res)

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

	res, err := s.cardStatisticRepository.GetMonthlyTransferAmountReceiver(year)

	if err != nil {
		traceID := traceunic.GenerateTraceID("MONTHLY_TRANSFER_AMOUNT_RECEIVER_NOT_FOUND")

		s.logger.Error("Failed to retrieve monthly transfer receiver amount", zap.String("trace_id", traceID), zap.Int("year", year), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Monthly transfer receiver amount not found")
		status = "monthly_transfer_amount_receiver_not_found"

		return nil, card_errors.ErrFailedFindMonthlyTransferAmountReceiver
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

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

	res, err := s.cardStatisticRepository.GetYearlyTransferAmountReceiver(year)

	if err != nil {
		traceID := traceunic.GenerateTraceID("YEARLY_TRANSFER_AMOUNT_RECEIVER_NOT_FOUND")

		s.logger.Error("Failed to retrieve yearly transfer receiver amount", zap.String("trace_id", traceID), zap.Int("year", year), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Yearly transfer receiver amount not found")
		status = "yearly_transfer_amount_receiver_not_found"

		return nil, card_errors.ErrFailedFindYearlyTransferAmountReceiver
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.logger.Debug("Yearly transfer receiver amount retrieved successfully",
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
