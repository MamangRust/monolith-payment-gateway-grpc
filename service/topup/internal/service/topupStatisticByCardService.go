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
	"go.opentelemetry.io/otel/codes"
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
	card_number := req.CardNumber
	year := req.Year
	month := req.Month

	const method = "FindMonthTopupStatusSuccessByCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("month", month), attribute.String("card_number", card_number))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthTopupStatusSuccessByCardNumberCache(req); found {
		logSuccess("Successfully fetched monthly topup status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.topupStatisticByCardRepository.GetMonthTopupStatusSuccessByCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleMonthTopupStatusSuccessByCardNumber(err, method, "FAILED_FIND_MONTHLY_METHODS_BY_CARD", span, &status, zap.Error(err))
	}
	so := s.mapping.ToTopupResponsesMonthStatusSuccess(records)

	s.mencache.SetMonthTopupStatusSuccessByCardNumberCache(req, so)

	logSuccess("Successfully fetched monthly topup status success", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	return so, nil
}

func (s *topupStatisticByCardService) FindYearlyTopupStatusSuccessByCardNumber(req *requests.YearTopupStatusCardNumber) ([]*response.TopupResponseYearStatusSuccess, *response.ErrorResponse) {
	card_number := req.CardNumber
	year := req.Year

	const method = "FindYearlyTopupStatusSuccessByCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", card_number))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyTopupStatusSuccessByCardNumberCache(req); found {
		logSuccess("Successfully fetched yearly topup status success", zap.Int("year", year), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.topupStatisticByCardRepository.GetYearlyTopupStatusSuccessByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleYearlyTopupStatusSuccessByCardNumber(err, method, "FAILED_FIND_YEARLY_METHODS_BY_CARD", span, &status, zap.Error(err))
	}
	so := s.mapping.ToTopupResponsesYearStatusSuccess(records)

	s.mencache.SetYearlyTopupStatusSuccessByCardNumberCache(req, so)

	logSuccess("Successfully fetched yearly topup status success", zap.Int("year", year), zap.String("card_number", card_number))

	return so, nil
}

func (s *topupStatisticByCardService) FindMonthTopupStatusFailedByCardNumber(req *requests.MonthTopupStatusCardNumber) ([]*response.TopupResponseMonthStatusFailed, *response.ErrorResponse) {
	card_number := req.CardNumber
	year := req.Year
	month := req.Month

	const method = "FindMonthTopupStatusFailedByCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("month", month), attribute.String("card_number", card_number))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthTopupStatusFailedByCardNumberCache(req); found {
		logSuccess("Successfully fetched monthly topup status Failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.topupStatisticByCardRepository.GetMonthTopupStatusFailedByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleMonthTopupStatusFailedByCardNumber(err, method, "FAILED_FIND_MONTHLY_METHODS_BY_CARD", span, &status, zap.Error(err))
	}
	so := s.mapping.ToTopupResponsesMonthStatusFailed(records)

	s.mencache.SetMonthTopupStatusFailedByCardNumberCache(req, so)

	logSuccess("Successfully fetched monthly topup status Failed", zap.Int("year", year), zap.Int("month", month), zap.String("card_number", card_number))

	return so, nil
}

func (s *topupStatisticByCardService) FindYearlyTopupStatusFailedByCardNumber(req *requests.YearTopupStatusCardNumber) ([]*response.TopupResponseYearStatusFailed, *response.ErrorResponse) {

	card_number := req.CardNumber
	year := req.Year

	const method = "FindYearlyTopupStatusFailedByCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", card_number))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyTopupStatusFailedByCardNumberCache(req); found {
		logSuccess("Successfully fetched yearly topup status Failed", zap.Int("year", year), zap.String("card_number", card_number))
		return data, nil
	}

	records, err := s.topupStatisticByCardRepository.GetYearlyTopupStatusFailedByCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleYearlyTopupStatusFailedByCardNumber(err, method, "FAILED_FIND_YEARLY_AMOUNTS_BY_CARD", span, &status, zap.Int("year", year), zap.String("card_number", card_number))
	}
	so := s.mapping.ToTopupResponsesYearStatusFailed(records)

	s.mencache.SetYearlyTopupStatusFailedByCardNumberCache(req, so)

	logSuccess("Successfully fetched yearly topup status Failed", zap.Int("year", year), zap.String("card_number", card_number))

	return so, nil
}

func (s *topupStatisticByCardService) FindMonthlyTopupMethodsByCardNumber(req *requests.YearMonthMethod) ([]*response.TopupMonthMethodResponse, *response.ErrorResponse) {

	year := req.Year
	cardNumber := req.CardNumber

	const method = "FindMonthlyTopupMethodsByCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyTopupMethodsByCardNumberCache(req); found {
		logSuccess("Successfully fetched monthly topup methods by card number", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	records, err := s.topupStatisticByCardRepository.GetMonthlyTopupMethodsByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleMonthlyTopupMethodsByCardNumber(err, method, "FAILED_FIND_MONTHLY_METHODS_BY_CARD", span, &status, zap.String("card_number", cardNumber), zap.Int("year", year))
	}

	responses := s.mapping.ToTopupMonthlyMethodResponses(records)

	s.mencache.SetMonthlyTopupMethodsByCardNumberCache(req, responses)

	logSuccess("Successfully fetched monthly topup methods by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

func (s *topupStatisticByCardService) FindYearlyTopupMethodsByCardNumber(req *requests.YearMonthMethod) ([]*response.TopupYearlyMethodResponse, *response.ErrorResponse) {

	year := req.Year
	cardNumber := req.CardNumber

	const method = "FindYearlyTopupMethodsByCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyTopupMethodsByCardNumberCache(req); found {
		logSuccess("Successfully fetched yearly topup methods by card number", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	records, err := s.topupStatisticByCardRepository.GetYearlyTopupMethodsByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleYearlyTopupMethodsByCardNumber(err, method, "FAILED_FIND_YEARLY_TOPUP_METHODS_BY_CARD", span, &status, zap.Error(err))
	}

	responses := s.mapping.ToTopupYearlyMethodResponses(records)

	s.mencache.SetYearlyTopupMethodsByCardNumberCache(req, responses)

	logSuccess("Successfully fetched yearly topup methods by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

func (s *topupStatisticByCardService) FindMonthlyTopupAmountsByCardNumber(req *requests.YearMonthMethod) ([]*response.TopupMonthAmountResponse, *response.ErrorResponse) {
	year := req.Year
	cardNumber := req.CardNumber

	const method = "FindMonthlyTopupAmountsByCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyTopupAmountsByCardNumberCache(req); found {
		logSuccess("Successfully fetched monthly topup amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	records, err := s.topupStatisticByCardRepository.GetMonthlyTopupAmountsByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleMonthlyTopupAmountsByCardNumber(err, method, "FAILED_FIND_MONTHLY_AMOUNTS_BY_CARD", span, &status, zap.Error(err))
	}

	responses := s.mapping.ToTopupMonthlyAmountResponses(records)

	s.mencache.SetMonthlyTopupAmountsByCardNumberCache(req, responses)

	logSuccess("Successfully fetched monthly topup amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

func (s *topupStatisticByCardService) FindYearlyTopupAmountsByCardNumber(req *requests.YearMonthMethod) ([]*response.TopupYearlyAmountResponse, *response.ErrorResponse) {
	year := req.Year
	cardNumber := req.CardNumber

	const method = "FindMonthTopupStatusSuccessByCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.String("card_number", cardNumber))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyTopupAmountsByCardNumberCache(req); found {
		logSuccess("Successfully fetched yearly topup amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))
		return data, nil
	}

	records, err := s.topupStatisticByCardRepository.GetYearlyTopupAmountsByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleYearlyTopupAmountsByCardNumber(err, method, "FAILED_FIND_YEARLY_TOPUP_AMOUNTS_BY_CARD", span, &status, zap.Error(err))
	}

	responses := s.mapping.ToTopupYearlyAmountResponses(records)

	s.mencache.SetYearlyTopupAmountsByCardNumberCache(req, responses)

	logSuccess("Successfully fetched yearly topup amounts by card number", zap.String("card_number", cardNumber), zap.Int("year", year))

	return responses, nil
}

func (s *topupStatisticByCardService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *topupStatisticByCardService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
