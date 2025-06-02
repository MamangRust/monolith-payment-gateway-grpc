package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-topup/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type topupStatisticByCardService struct {
	ctx                            context.Context
	mencache                       mencache.TopupStatisticByCardCache
	errorhandler                   errorhandler.TopupStatisticByCardErrorHandler
	trace                          trace.Tracer
	topupStatisticByCardRepository repository.TopupStatisticByCardRepository
	logger                         logger.LoggerInterface
	mapping                        responseservice.TopupResponseMapper
	requestCounter                 *prometheus.CounterVec
	requestDuration                *prometheus.HistogramVec
}

func NewTopupStatisticByCardService(ctx context.Context, mencache mencache.TopupStatisticByCardCache,
	errorhandler errorhandler.TopupStatisticByCardErrorHandler, topupStatisticByCard repository.TopupStatisticByCardRepository, logger logger.LoggerInterface, mapping responseservice.TopupResponseMapper) *topupStatisticByCardService {
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
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &topupStatisticByCardService{
		ctx:                            ctx,
		mencache:                       mencache,
		errorhandler:                   errorhandler,
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

	if data := s.mencache.GetMonthTopupStatusSuccessByCardNumberCache(req); data != nil {
		s.logger.Debug("Successfully fetched monthly topup status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.topupStatisticByCardRepository.GetMonthTopupStatusSuccessByCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleMonthTopupStatusSuccessByCardNumber(err, "FindMonthTopupStatusSuccessByCardNumber", "FAILED_FIND_MONTHLY_METHODS_BY_CARD", span, &status, zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))
	}
	so := s.mapping.ToTopupResponsesMonthStatusSuccess(records)

	s.mencache.SetMonthTopupStatusSuccessByCardNumberCache(req, so)

	s.logger.Debug("Successfully fetched monthly topup status success", zap.Int("year", year), zap.Int("month", month))

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

	if data := s.mencache.GetYearlyTopupStatusSuccessByCardNumberCache(req); data != nil {
		s.logger.Debug("Successfully fetched yearly topup status success", zap.Int("year", year), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.topupStatisticByCardRepository.GetYearlyTopupStatusSuccessByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleYearlyTopupStatusSuccessByCardNumber(err, "FindYearlyTopupStatusSuccessByCardNumber", "FAILED_FIND_YEARLY_METHODS_BY_CARD", span, &status, zap.Int("year", year), zap.String("card_number", card_number))
	}
	so := s.mapping.ToTopupResponsesYearStatusSuccess(records)

	s.mencache.SetYearlyTopupStatusSuccessByCardNumberCache(req, so)

	s.logger.Debug("Successfully fetched yearly topup status success", zap.Int("year", year))

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

	if data := s.mencache.GetMonthTopupStatusFailedByCardNumberCache(req); data != nil {
		s.logger.Debug("Successfully fetched monthly topup status Failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.topupStatisticByCardRepository.GetMonthTopupStatusFailedByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleMonthTopupStatusFailedByCardNumber(err, "FindMonthTopupStatusFailedByCardNumber", "FAILED_FIND_MONTHLY_METHODS_BY_CARD", span, &status, zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))
	}
	so := s.mapping.ToTopupResponsesMonthStatusFailed(records)

	s.mencache.SetMonthTopupStatusFailedByCardNumberCache(req, so)

	s.logger.Debug("Failedfully fetched monthly topup status Failed", zap.Int("year", year), zap.Int("month", month))

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

	if data := s.mencache.GetYearlyTopupStatusFailedByCardNumberCache(req); data != nil {
		s.logger.Debug("Successfully fetched yearly topup status Failed", zap.Int("year", year), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.topupStatisticByCardRepository.GetYearlyTopupStatusFailedByCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleYearlyTopupStatusFailedByCardNumber(err, "FindYearlyTopupStatusFailedByCardNumber", "FAILED_FIND_YEARLY_AMOUNTS_BY_CARD", span, &status, zap.Int("year", year), zap.String("card_number", card_number))
	}
	so := s.mapping.ToTopupResponsesYearStatusFailed(records)

	s.mencache.SetYearlyTopupStatusFailedByCardNumberCache(req, so)

	s.logger.Debug("Failedfully fetched yearly topup status Failed", zap.Int("year", year))

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

	if data := s.mencache.GetMonthlyTopupMethodsByCardNumberCache(req); data != nil {
		s.logger.Debug("Successfully fetched monthly topup methods by card number", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	records, err := s.topupStatisticByCardRepository.GetMonthlyTopupMethodsByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleMonthlyTopupMethodsByCardNumber(err, "FindMonthlyTopupMethodsByCardNumber", "FAILED_FIND_MONTHLY_METHODS_BY_CARD", span, &status, zap.String("card_number", cardNumber), zap.Int("year", year))
	}

	responses := s.mapping.ToTopupMonthlyMethodResponses(records)

	s.mencache.SetMonthlyTopupMethodsByCardNumberCache(req, responses)

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

	if data := s.mencache.GetYearlyTopupMethodsByCardNumberCache(req); data != nil {
		s.logger.Debug("Successfully fetched yearly topup methods by card number", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	records, err := s.topupStatisticByCardRepository.GetYearlyTopupMethodsByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleYearlyTopupMethodsByCardNumber(err, "FindYearlyTopupMethodsByCardNumber", "FAILED_FIND_YEARLY_TOPUP_METHODS_BY_CARD", span, &status, zap.String("card_number", cardNumber), zap.Int("year", year))
	}

	responses := s.mapping.ToTopupYearlyMethodResponses(records)

	s.mencache.SetYearlyTopupMethodsByCardNumberCache(req, responses)

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

	if data := s.mencache.GetMonthlyTopupAmountsByCardNumberCache(req); data != nil {
		s.logger.Debug("Successfully fetched monthly topup amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	records, err := s.topupStatisticByCardRepository.GetMonthlyTopupAmountsByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleMonthlyTopupAmountsByCardNumber(err, "FindMonthlyTopupAmountsByCardNumber", "FAILED_FIND_MONTHLY_AMOUNTS_BY_CARD", span, &status, zap.String("card_number", cardNumber), zap.Int("year", year))
	}

	responses := s.mapping.ToTopupMonthlyAmountResponses(records)

	s.mencache.SetMonthlyTopupAmountsByCardNumberCache(req, responses)

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

	if data := s.mencache.GetYearlyTopupAmountsByCardNumberCache(req); data != nil {
		s.logger.Debug("Successfully fetched yearly topup amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	records, err := s.topupStatisticByCardRepository.GetYearlyTopupAmountsByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleYearlyTopupAmountsByCardNumber(err, "FindYearlyTopupAmountsByCardNumber", "FAILED_FIND_YEARLY_TOPUP_AMOUNTS_BY_CARD", span, &status, zap.String("card_number", cardNumber), zap.Int("year", year))
	}

	responses := s.mapping.ToTopupYearlyAmountResponses(records)

	s.mencache.SetYearlyTopupAmountsByCardNumberCache(req, responses)

	s.logger.Debug("Successfully fetched yearly topup amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

func (s *topupStatisticByCardService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
