package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	traceunic "github.com/MamangRust/monolith-payment-gateway-pkg/trace_unic"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type topupStatisticByCardService struct {
	ctx                            context.Context
	trace                          trace.Tracer
	topupStatisticByCardRepository repository.TopupStatisticByCardRepository
	logger                         logger.LoggerInterface
	mapping                        responseservice.TopupResponseMapper
	requestCounter                 *prometheus.CounterVec
	requestDuration                *prometheus.HistogramVec
}

func NewTopupStatisticByCardService(ctx context.Context, topupStatisticByCard repository.TopupStatisticByCardRepository, logger logger.LoggerInterface, mapping responseservice.TopupResponseMapper) *topupStatisticByCardService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "topup_statistic_by_card_service_request_total",
			Help: "Total number of requests to the TopupStatisticByCardService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "topup_statistic_by_card_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TopupStatisticByCardService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &topupStatisticByCardService{
		ctx:                            ctx,
		trace:                          otel.Tracer("topup-statistic-by-card-service"),
		topupStatisticByCardRepository: topupStatisticByCard,
		logger:                         logger,
		mapping:                        mapping,
		requestCounter:                 requestCounter,
		requestDuration:                requestDuration,
	}
}

func (s *topupStatisticByCardService) FindMonthTopupStatusSuccessByCardNumber(req *requests.MonthTopupStatusCardNumber) ([]*response.TopupResponseMonthStatusSuccess, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthTopupStatusSuccessByCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthTopupStatusSuccessByCardNumber")
	defer span.End()

	card_number := req.CardNumber
	year := req.Year
	month := req.Month

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
		attribute.String("card_number", card_number),
	)

	s.logger.Debug("Fetching monthly topup status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	records, err := s.topupStatisticByCardRepository.GetMonthTopupStatusSuccessByCardNumber(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTH_TOPUP_STATUS_SUCCESS_BY_CARD")

		s.logger.Error("Failed to fetch monthly topup status success", zap.Error(err), zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch monthly topup status success")
		status = "failed_to_fetch_monthly_topup_status_success"

		return nil, topup_errors.ErrFailedFindMonthTopupStatusSuccessByCard
	}

	s.logger.Debug("Successfully fetched monthly topup status success", zap.Int("year", year), zap.Int("month", month))

	so := s.mapping.ToTopupResponsesMonthStatusSuccess(records)

	return so, nil
}

func (s *topupStatisticByCardService) FindYearlyTopupStatusSuccessByCardNumber(req *requests.YearTopupStatusCardNumber) ([]*response.TopupResponseYearStatusSuccess, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTopupStatusSuccessByCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTopupStatusSuccessByCardNumber")
	defer span.End()

	card_number := req.CardNumber
	year := req.Year

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.String("card_number", card_number),
	)

	s.logger.Debug("Fetching yearly topup status success", zap.Int("year", year), zap.String("card_number", card_number))

	records, err := s.topupStatisticByCardRepository.GetYearlyTopupStatusSuccessByCardNumber(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_TOPUP_STATUS_SUCCESS_BY_CARD")

		s.logger.Error("Failed to fetch yearly topup status success", zap.Error(err), zap.Int("year", year), zap.String("card_number", card_number))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch yearly topup status success")
		status = "failed_to_fetch_yearly_topup_status_success"

		return nil, topup_errors.ErrFailedFindYearlyTopupStatusSuccessByCard
	}

	s.logger.Debug("Successfully fetched yearly topup status success", zap.Int("year", year))

	so := s.mapping.ToTopupResponsesYearStatusSuccess(records)

	return so, nil
}

func (s *topupStatisticByCardService) FindMonthTopupStatusFailedByCardNumber(req *requests.MonthTopupStatusCardNumber) ([]*response.TopupResponseMonthStatusFailed, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthTopupStatusFailedByCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthTopupStatusFailedByCardNumber")
	defer span.End()

	card_number := req.CardNumber
	year := req.Year
	month := req.Month

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
		attribute.String("card_number", card_number),
	)

	s.logger.Debug("Fetching monthly topup status Failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	records, err := s.topupStatisticByCardRepository.GetMonthTopupStatusFailedByCardNumber(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTH_TOPUP_STATUS_FAILED_BY_CARD")

		s.logger.Error("Failed to fetch monthly topup status Failed", zap.Error(err), zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch monthly topup status Failed")
		status = "failed_to_fetch_monthly_topup_status_failed"

		return nil, topup_errors.ErrFailedFindMonthTopupStatusFailedByCard
	}

	s.logger.Debug("Failedfully fetched monthly topup status Failed", zap.Int("year", year), zap.Int("month", month))

	so := s.mapping.ToTopupResponsesMonthStatusFailed(records)

	return so, nil
}

func (s *topupStatisticByCardService) FindYearlyTopupStatusFailedByCardNumber(req *requests.YearTopupStatusCardNumber) ([]*response.TopupResponseYearStatusFailed, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTopupStatusFailedByCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTopupStatusFailedByCardNumber")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", req.Year),
		attribute.String("card_number", req.CardNumber),
	)

	card_number := req.CardNumber
	year := req.Year

	s.logger.Debug("Fetching yearly topup status Failed", zap.Int("year", year), zap.String("card_number", card_number))

	records, err := s.topupStatisticByCardRepository.GetYearlyTopupStatusFailedByCardNumber(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_TOPUP_STATUS_FAILED_BY_CARD")

		s.logger.Error("Failed to fetch yearly topup status Failed", zap.Error(err), zap.Int("year", year), zap.String("card_number", card_number))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch yearly topup status Failed")
		status = "failed_to_fetch_yearly_topup_status_failed"
		return nil, topup_errors.ErrFailedFindYearlyTopupStatusFailedByCard
	}

	s.logger.Debug("Failedfully fetched yearly topup status Failed", zap.Int("year", year))

	so := s.mapping.ToTopupResponsesYearStatusFailed(records)

	return so, nil
}

func (s *topupStatisticByCardService) FindMonthlyTopupMethodsByCardNumber(req *requests.YearMonthMethod) ([]*response.TopupMonthMethodResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyTopupMethodsByCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyTopupMethodsByCardNumber")
	defer span.End()

	year := req.Year
	cardNumber := req.CardNumber

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.String("card_number", cardNumber),
	)

	s.logger.Debug("Fetching monthly topup methods by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	records, err := s.topupStatisticByCardRepository.GetMonthlyTopupMethodsByCardNumber(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTH_TOPUP_METHODS_BY_CARD")

		s.logger.Error("Failed to fetch monthly topup methods by card number", zap.Error(err), zap.String("card_number", cardNumber), zap.Int("year", year))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch monthly topup methods by card number")
		status = "failed_to_fetch_monthly_topup_methods_by_card"

		return nil, topup_errors.ErrFailedFindMonthlyTopupMethodsByCard
	}

	responses := s.mapping.ToTopupMonthlyMethodResponses(records)

	s.logger.Debug("Successfully fetched monthly topup methods by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

func (s *topupStatisticByCardService) FindYearlyTopupMethodsByCardNumber(req *requests.YearMonthMethod) ([]*response.TopupYearlyMethodResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTopupMethodsByCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTopupMethodsByCardNumber")
	defer span.End()

	year := req.Year
	cardNumber := req.CardNumber

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.String("card_number", cardNumber),
	)

	s.logger.Debug("Fetching yearly topup methods by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	records, err := s.topupStatisticByCardRepository.GetYearlyTopupMethodsByCardNumber(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_TOPUP_METHODS_BY_CARD")

		s.logger.Error("Failed to fetch yearly topup methods by card number", zap.Error(err), zap.String("card_number", cardNumber), zap.Int("year", year))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch yearly topup methods by card number")
		status = "failed_to_fetch_yearly_topup_methods_by_card"

		return nil, topup_errors.ErrFailedFindYearlyTopupMethodsByCard
	}

	responses := s.mapping.ToTopupYearlyMethodResponses(records)

	s.logger.Debug("Successfully fetched yearly topup methods by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

func (s *topupStatisticByCardService) FindMonthlyTopupAmountsByCardNumber(req *requests.YearMonthMethod) ([]*response.TopupMonthAmountResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyTopupAmountsByCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyTopupAmountsByCardNumber")
	defer span.End()

	year := req.Year
	cardNumber := req.CardNumber

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.String("card_number", cardNumber),
	)
	s.logger.Debug("Fetching monthly topup amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	records, err := s.topupStatisticByCardRepository.GetMonthlyTopupAmountsByCardNumber(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTH_TOPUP_AMOUNTS_BY_CARD")

		s.logger.Error("Failed to fetch monthly topup amounts by card number", zap.Error(err), zap.String("card_number", cardNumber), zap.Int("year", year))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch monthly topup amounts by card number")
		status = "failed_to_fetch_monthly_topup_amounts_by_card"

		return nil, topup_errors.ErrFailedFindMonthlyTopupAmountsByCard
	}

	responses := s.mapping.ToTopupMonthlyAmountResponses(records)

	s.logger.Debug("Successfully fetched monthly topup amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

func (s *topupStatisticByCardService) FindYearlyTopupAmountsByCardNumber(req *requests.YearMonthMethod) ([]*response.TopupYearlyAmountResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTopupAmountsByCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTopupAmountsByCardNumber")
	defer span.End()

	year := req.Year
	cardNumber := req.CardNumber

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.String("card_number", cardNumber),
	)

	s.logger.Debug("Fetching yearly topup amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	records, err := s.topupStatisticByCardRepository.GetYearlyTopupAmountsByCardNumber(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_TOPUP_AMOUNTS_BY_CARD")

		s.logger.Error("Failed to fetch yearly topup amounts by card number", zap.Error(err), zap.String("card_number", cardNumber), zap.Int("year", year))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch yearly topup amounts by card number")
		status = "failed_to_fetch_yearly_topup_amounts_by_card"

		return nil, topup_errors.ErrFailedFindYearlyTopupAmountsByCard
	}

	responses := s.mapping.ToTopupYearlyAmountResponses(records)

	s.logger.Debug("Successfully fetched yearly topup amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

func (s *topupStatisticByCardService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
