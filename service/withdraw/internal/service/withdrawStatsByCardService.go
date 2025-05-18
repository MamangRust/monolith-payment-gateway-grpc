package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	traceunic "github.com/MamangRust/monolith-payment-gateway-pkg/trace_unic"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type withdrawStatisticByCardService struct {
	ctx                               context.Context
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
		[]string{"method"},
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
	card_number := req.CardNumber
	month := req.Month

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
		attribute.String("card_number", card_number),
	)

	s.logger.Debug("Fetching monthly Withdraw status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	records, err := s.withdrawStatisticByCardRepository.GetMonthWithdrawStatusSuccessByCardNumber(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTH_WITHDRAW_SUCCESS")

		s.logger.Error("failed to fetch monthly Withdraw status success", zap.Error(err), zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch monthly Withdraw status success")
		status = "failed_find_month_withdraw_status_success"

		return nil, withdraw_errors.ErrFailedFindMonthWithdrawStatusSuccess
	}

	s.logger.Debug("Successfully fetched monthly Withdraw status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	so := s.mapping.ToWithdrawResponsesMonthStatusSuccess(records)

	return so, nil
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
	card_number := req.CardNumber

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.String("card_number", card_number),
	)

	s.logger.Debug("Fetching yearly Withdraw status success", zap.Int("year", year), zap.String("card_number", card_number))

	records, err := s.withdrawStatisticByCardRepository.GetYearlyWithdrawStatusSuccessByCardNumber(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEAR_WITHDRAW_SUCCESS")

		s.logger.Error("failed to fetch yearly Withdraw status success", zap.Error(err), zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch yearly Withdraw status success")
		status = "failed_find_year_withdraw_status_success"

		return nil, withdraw_errors.ErrFailedFindYearWithdrawStatusSuccess
	}

	s.logger.Debug("Successfully fetched yearly Withdraw status success", zap.Int("year", year), zap.String("card_number", card_number))

	so := s.mapping.ToWithdrawResponsesYearStatusSuccess(records)

	return so, nil
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
	card_number := req.CardNumber
	month := req.Month

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
		attribute.String("card_number", card_number),
	)

	s.logger.Debug("Fetching monthly Withdraw status Failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	records, err := s.withdrawStatisticByCardRepository.GetMonthWithdrawStatusFailedByCardNumber(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTH_WITHDRAW_FAILED")

		s.logger.Error("failed to fetch monthly Withdraw status Failed", zap.Error(err), zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch monthly Withdraw status Failed")
		status = "failed_find_month_withdraw_status_failed"

		return nil, withdraw_errors.ErrFailedFindMonthWithdrawStatusFailed
	}

	s.logger.Debug("Failedfully fetched monthly Withdraw status Failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	so := s.mapping.ToWithdrawResponsesMonthStatusFailed(records)

	return so, nil
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
	card_number := req.CardNumber

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.String("card_number", card_number),
	)

	s.logger.Debug("Fetching yearly Withdraw status Failed", zap.Int("year", year), zap.String("card_number", card_number))

	records, err := s.withdrawStatisticByCardRepository.GetYearlyWithdrawStatusFailedByCardNumber(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEAR_WITHDRAW_FAILED")

		s.logger.Error("failed to fetch yearly Withdraw status Failed", zap.Error(err), zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch yearly Withdraw status Failed")
		status = "failed_find_year_withdraw_status_failed"

		return nil, withdraw_errors.ErrFailedFindYearWithdrawStatusFailed
	}

	s.logger.Debug("Failedfully fetched yearly Withdraw status Failed", zap.Int("year", year), zap.String("card_number", card_number))

	so := s.mapping.ToWithdrawResponsesYearStatusFailed(records)

	return so, nil
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

	withdraws, err := s.withdrawStatisticByCardRepository.GetMonthlyWithdrawsByCardNumber(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTH_WITHDRAW")

		s.logger.Error("failed to find monthly withdraws by card number", zap.Error(err), zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find monthly withdraws by card number")
		status = "failed_find_month_withdraw"

		return nil, withdraw_errors.ErrFailedFindMonthlyWithdraws
	}

	responseWithdraws := s.mapping.ToWithdrawsAmountMonthlyResponses(withdraws)

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

	withdraws, err := s.withdrawStatisticByCardRepository.GetYearlyWithdrawsByCardNumber(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEAR_WITHDRAW")

		s.logger.Error("failed to find yearly withdraws by card number", zap.Error(err), zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find yearly withdraws by card number")
		status = "failed_find_year_withdraw"

		return nil, withdraw_errors.ErrFailedFindYearlyWithdraws
	}

	responseWithdraws := s.mapping.ToWithdrawsAmountYearlyResponses(withdraws)

	s.logger.Debug("Successfully fetched yearly withdraws by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responseWithdraws, nil
}

func (s *withdrawStatisticByCardService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
