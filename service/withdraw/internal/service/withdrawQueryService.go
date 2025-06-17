package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/withdraw_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-withdraw/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-withdraw/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type withdrawQueryService struct {
	ctx                     context.Context
	errorhandler            errorhandler.WithdrawQueryErrorHandler
	mencache                mencache.WithdrawQueryCache
	trace                   trace.Tracer
	withdrawQueryRepository repository.WithdrawQueryRepository
	logger                  logger.LoggerInterface
	mapping                 responseservice.WithdrawResponseMapper
	requestCounter          *prometheus.CounterVec
	requestDuration         *prometheus.HistogramVec
}

func NewWithdrawQueryService(ctx context.Context, errorhandler errorhandler.WithdrawQueryErrorHandler,
	mencache mencache.WithdrawQueryCache, withdrawQueryRepository repository.WithdrawQueryRepository, logger logger.LoggerInterface, mapping responseservice.WithdrawResponseMapper) *withdrawQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "withdraw_query_service_request_total",
			Help: "Total number of requests to the WithdrawQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "withdraw_query_service_request_duration_seconds",
			Help:    "Histogram of request durations for the WithdrawQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &withdrawQueryService{
		ctx:                     ctx,
		mencache:                mencache,
		errorhandler:            errorhandler,
		trace:                   otel.Tracer("withdraw-query-service"),
		withdrawQueryRepository: withdrawQueryRepository,
		logger:                  logger,
		mapping:                 mapping,
		requestCounter:          requestCounter,
		requestDuration:         requestDuration,
	}
}

func (s *withdrawQueryService) FindAll(req *requests.FindAllWithdraws) ([]*response.WithdrawResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindAll"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedWithdrawsCache(req); found {
		logSuccess("Successfully retrieved all withdraw records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	withdraws, totalRecords, err := s.withdrawQueryRepository.FindAll(req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_WITHDRAW", span, &status, withdraw_errors.ErrFailedFindAllWithdraws, zap.Error(err))
	}

	withdrawResponse := s.mapping.ToWithdrawsResponse(withdraws)

	s.mencache.SetCachedWithdrawsCache(req, withdrawResponse, totalRecords)

	logSuccess("Successfully retrieved all withdraw records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return withdrawResponse, totalRecords, nil
}

func (s *withdrawQueryService) FindAllByCardNumber(req *requests.FindAllWithdrawCardNumber) ([]*response.WithdrawResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindAllByCardNumber"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedWithdrawByCardCache(req); found {
		logSuccess("Successfully retrieved all withdraw records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}
	withdraws, totalRecords, err := s.withdrawQueryRepository.FindAllByCardNumber(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_WITHDRAW_BY_CARD", span, &status, withdraw_errors.ErrFailedFindAllWithdraws, zap.Error(err))
	}

	withdrawResponse := s.mapping.ToWithdrawsResponse(withdraws)

	s.mencache.SetCachedWithdrawByCardCache(req, withdrawResponse, totalRecords)

	logSuccess("Successfully retrieved all withdraw records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return withdrawResponse, totalRecords, nil
}

func (s *withdrawQueryService) FindById(withdrawID int) (*response.WithdrawResponse, *response.ErrorResponse) {
	const method = "FindById"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("withdraw.id", withdrawID))

	defer func() {
		end(status)
	}()

	if data := s.mencache.GetCachedWithdrawCache(withdrawID); data != nil {
		logSuccess("Successfully retrieved withdraw from cache", zap.Int("withdraw_id", withdrawID))
		return data, nil
	}

	withdraw, err := s.withdrawQueryRepository.FindById(withdrawID)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_WITHDRAW", span, &status, withdraw_errors.ErrWithdrawNotFound, zap.Error(err))
	}

	withdrawResponse := s.mapping.ToWithdrawResponse(withdraw)

	s.mencache.SetCachedWithdrawCache(withdrawResponse)

	logSuccess("Successfully retrieved withdraw", zap.Int("withdraw.id", withdrawID))

	return withdrawResponse, nil
}

func (s *withdrawQueryService) FindByActive(req *requests.FindAllWithdraws) ([]*response.WithdrawResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindByActive"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	withdraws, totalRecords, err := s.withdrawQueryRepository.FindByActive(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_ACTIVE_WITHDRAW", span, &status, withdraw_errors.ErrFailedFindAllWithdraws, zap.Error(err))
	}

	withdrawResponses := s.mapping.ToWithdrawsResponseDeleteAt(withdraws)

	logSuccess("Successfully retrieved all withdraw records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return withdrawResponses, totalRecords, nil
}

func (s *withdrawQueryService) FindByTrashed(req *requests.FindAllWithdraws) ([]*response.WithdrawResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindByTrashed"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	withdraws, totalRecords, err := s.withdrawQueryRepository.FindByTrashed(req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_TRASHED_WITHDRAW", span, &status, withdraw_errors.ErrFailedFindAllWithdraws, zap.Error(err))
	}

	withdrawResponses := s.mapping.ToWithdrawsResponseDeleteAt(withdraws)

	logSuccess("Successfully retrieved all withdraw records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return withdrawResponses, totalRecords, nil
}

func (s *withdrawQueryService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *withdrawQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func (s *withdrawQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
