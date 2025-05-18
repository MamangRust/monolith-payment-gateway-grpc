package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-card/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	traceunic "github.com/MamangRust/monolith-payment-gateway-pkg/trace_unic"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
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

type cardStatisticBycardService struct {
	ctx                     context.Context
	trace                   trace.Tracer
	cardStatisticRepository repository.CardStatisticByCardRepository
	logger                  logger.LoggerInterface
	mapping                 responseservice.CardResponseMapper
	requestCounter          *prometheus.CounterVec
	requestDuration         *prometheus.HistogramVec
}

func NewCardStatisticBycardService(
	ctx context.Context,
	cardStatisticRepository repository.CardStatisticByCardRepository, logger logger.LoggerInterface, mapper responseservice.CardResponseMapper) *cardStatisticBycardService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "card_statistic_bycard_request_count",
		Help: "Number of card statistic requests CardStatisticBycardService",
	}, []string{"status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "card_statistic_bycard_request_duration_seconds",
		Help:    "Duration of card statistic requests CardStatisticBycardService",
		Buckets: prometheus.DefBuckets,
	}, []string{"status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	return &cardStatisticBycardService{
		ctx:                     ctx,
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

	res, err := s.cardStatisticRepository.GetMonthlyBalancesByCardNumber(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("MONTHLY_BALANCE_NOT_FOUND")

		s.logger.Error("Failed to retrieve monthly balance", zap.String("trace_id", traceID), zap.Int("year", year), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Monthly balance not found")
		status = "monthly_balance_not_found"

		return nil, card_errors.ErrFailedFindMonthlyBalanceByCard
	}

	so := s.mapping.ToGetMonthlyBalances(res)

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

	res, err := s.cardStatisticRepository.GetYearlyBalanceByCardNumber(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("YEARLY_BALANCE_NOT_FOUND")

		s.logger.Error("Failed to retrieve yearly balance", zap.String("trace_id", traceID), zap.Int("year", year), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Yearly balance not found")
		status = "yearly_balance_not_found"

		return nil, card_errors.ErrFailedFindYearlyBalanceByCard
	}

	so := s.mapping.ToGetYearlyBalances(res)

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

	res, err := s.cardStatisticRepository.GetMonthlyTopupAmountByCardNumber(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("MONTHLY_TOPUP_AMOUNT_NOT_FOUND")

		s.logger.Error("Failed to retrieve monthly topup amount by card number", zap.String("trace_id", traceID), zap.String("card_number", cardNumber), zap.Int("year", year), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Monthly topup amount not found")
		status = "monthly_topup_amount_not_found"

		return nil, card_errors.ErrFailedFindMonthlyTopupAmountByCard
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

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

	res, err := s.cardStatisticRepository.GetYearlyTopupAmountByCardNumber(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("YEARLY_TOPUP_AMOUNT_NOT_FOUND")

		s.logger.Error("Failed to retrieve yearly topup amount by card number", zap.String("trace_id", traceID), zap.String("card_number", cardNumber), zap.Int("year", year), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Yearly topup amount not found")
		status = "yearly_topup_amount_not_found"

		return nil, card_errors.ErrFailedFindYearlyTopupAmountByCard
	}

	so := s.mapping.ToGetYearlyAmounts(res)

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

	res, err := s.cardStatisticRepository.GetMonthlyWithdrawAmountByCardNumber(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("MONTHLY_WITHDRAW_AMOUNT_NOT_FOUND")

		s.logger.Error("Failed to retrieve monthly withdraw amount by card number", zap.String("trace_id", traceID), zap.String("card_number", cardNumber), zap.Int("year", year), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Monthly withdraw amount not found")
		status = "monthly_withdraw_amount_not_found"

		return nil, card_errors.ErrFailedFindMonthlyWithdrawAmountByCard
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

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

	res, err := s.cardStatisticRepository.GetYearlyWithdrawAmountByCardNumber(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("YEARLY_WITHDRAW_AMOUNT_NOT_FOUND")

		s.logger.Error("Failed to retrieve yearly withdraw amount by card number", zap.String("trace_id", traceID), zap.String("card_number", cardNumber), zap.Int("year", year), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Yearly withdraw amount not found")
		status = "yearly_withdraw_amount_not_found"

		return nil, card_errors.ErrFailedFindYearlyWithdrawAmountByCard
	}

	so := s.mapping.ToGetYearlyAmounts(res)

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

	res, err := s.cardStatisticRepository.GetMonthlyTransactionAmountByCardNumber(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("MONTHLY_TRANSACTION_AMOUNT_NOT_FOUND")

		s.logger.Error("Failed to retrieve monthly transaction amount by card number", zap.String("trace_id", traceID), zap.String("card_number", cardNumber), zap.Int("year", year), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Monthly transaction amount not found")
		status = "monthly_transaction_amount_not_found"

		return nil, card_errors.ErrFailedFindMonthlyTransactionAmountByCard
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

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

	res, err := s.cardStatisticRepository.GetYearlyTransactionAmountByCardNumber(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("YEARLY_TRANSACTION_AMOUNT_NOT_FOUND")

		s.logger.Error("Failed to retrieve yearly transaction amount by card number", zap.String("trace_id", traceID), zap.String("card_number", cardNumber), zap.Int("year", year), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Yearly transaction amount not found")
		status = "yearly_transaction_amount_not_found"

		return nil, card_errors.ErrFailedFindYearlyTransactionAmountByCard
	}

	so := s.mapping.ToGetYearlyAmounts(res)

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

	res, err := s.cardStatisticRepository.GetMonthlyTransferAmountBySender(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("MONTHLY_TRANSFER_AMOUNT_NOT_FOUND")

		s.logger.Error("Failed to retrieve monthly transfer sender amount by card number", zap.String("trace_id", traceID), zap.String("card_number", cardNumber), zap.Int("year", year), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Monthly transfer amount not found")
		status = "monthly_transfer_amount_not_found"

		return nil, card_errors.ErrFailedFindMonthlyTransferAmountBySender
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

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

	res, err := s.cardStatisticRepository.GetYearlyTransferAmountBySender(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("YEARLY_TRANSFER_AMOUNT_NOT_FOUND")

		s.logger.Error("Failed to retrieve yearly transfer sender amount by card number", zap.String("trace_id", traceID), zap.String("card_number", cardNumber), zap.Int("year", year), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Yearly transfer amount not found")
		status = "yearly_transfer_amount_not_found"

		return nil, card_errors.ErrFailedFindYearlyTransferAmountBySender
	}

	so := s.mapping.ToGetYearlyAmounts(res)

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

	res, err := s.cardStatisticRepository.GetMonthlyTransferAmountByReceiver(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("MONTHLY_TRANSFER_AMOUNT_NOT_FOUND")

		s.logger.Error("Failed to retrieve monthly transfer receiver amount by card number", zap.String("trace_id", traceID), zap.String("card_number", cardNumber), zap.Int("year", year), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Monthly transfer amount not found")
		status = "monthly_transfer_amount_not_found"

		return nil, card_errors.ErrFailedFindMonthlyTransferAmountByReceiver
	}

	so := s.mapping.ToGetMonthlyAmounts(res)

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

	res, err := s.cardStatisticRepository.GetYearlyTransferAmountByReceiver(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("YEARLY_TRANSFER_AMOUNT_NOT_FOUND")

		s.logger.Error("Failed to retrieve yearly transfer receiver amount by card number", zap.String("trace_id", traceID), zap.String("card_number", cardNumber), zap.Int("year", year), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Yearly transfer amount not found")
		status = "yearly_transfer_amount_not_found"

		return nil, card_errors.ErrFailedFindYearlyTransferAmountByReceiver
	}

	so := s.mapping.ToGetYearlyAmounts(res)

	s.logger.Debug("Yearly transfer receiver amount by card number retrieved successfully",
		zap.String("card_number", cardNumber),
		zap.Int("year", year),
		zap.Int("result_count", len(so)),
	)

	return so, nil
}

func (s *cardStatisticBycardService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
