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
	start := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("FindAll", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindAll")
	defer span.End()

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Attempting to fetch withdraw list",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	if data, total, found := s.mencache.GetCachedWithdrawsCache(req); found {
		s.logger.Debug("Successfully fetched withdraws from cache",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
		)
		return data, total, nil
	}

	withdraws, totalRecords, err := s.withdrawQueryRepository.FindAll(req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, "FindAll", "FAILED_FIND_ALL_WITHDRAW", span, &status, withdraw_errors.ErrFailedFindAllWithdraws, zap.Error(err))
	}

	withdrawResponse := s.mapping.ToWithdrawsResponse(withdraws)

	s.mencache.SetCachedWithdrawsCache(req, withdrawResponse, totalRecords)

	s.logger.Debug("Withdraw list fetched successfully",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.Int("total_records", *totalRecords),
	)

	return withdrawResponse, totalRecords, nil
}

func (s *withdrawQueryService) FindAllByCardNumber(req *requests.FindAllWithdrawCardNumber) ([]*response.WithdrawResponse, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindAllByCardNumber", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindAllByCardNumber")
	defer span.End()

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search
	cardNumber := req.CardNumber

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
		attribute.String("cardNumber", cardNumber),
	)

	s.logger.Debug("Fetching withdraw by card number",
		zap.String("cardNumber", cardNumber),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if data, total, found := s.mencache.GetCachedWithdrawByCardCache(req); found {
		s.logger.Debug("Successfully fetched withdraw by card number from cache",
			zap.String("cardNumber", cardNumber),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
		return data, total, nil
	}

	withdraws, totalRecords, err := s.withdrawQueryRepository.FindAllByCardNumber(req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, "FindAllByCardNumber", "FAILED_FIND_WITHDRAW_BY_CARD", span, &status, withdraw_errors.ErrFailedFindAllWithdraws, zap.Error(err))
	}

	withdrawResponse := s.mapping.ToWithdrawsResponse(withdraws)

	s.mencache.SetCachedWithdrawByCardCache(req, withdrawResponse, totalRecords)

	return withdrawResponse, totalRecords, nil
}

func (s *withdrawQueryService) FindById(withdrawID int) (*response.WithdrawResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("FindById", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindById")
	defer span.End()

	span.SetAttributes(attribute.Int("withdraw_id", withdrawID))

	s.logger.Debug("Initiating retrieval of withdraw data by ID", zap.Int("withdraw_id", withdrawID))

	if data := s.mencache.GetCachedWithdrawCache(withdrawID); data != nil {
		s.logger.Debug("Successfully retrieved withdraw data from cache", zap.Int("withdraw_id", withdrawID))
		return data, nil
	}

	withdraw, err := s.withdrawQueryRepository.FindById(withdrawID)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "FindById", "FAILED_FIND_WITHDRAW", span, &status, withdraw_errors.ErrWithdrawNotFound, zap.Int("withdraw_id", withdrawID))
	}

	withdrawResponse := s.mapping.ToWithdrawResponse(withdraw)

	s.mencache.SetCachedWithdrawCache(withdrawResponse)

	s.logger.Debug("Withdraw data retrieval completed successfully", zap.Int("withdraw_id", withdrawID))

	return withdrawResponse, nil
}

func (s *withdrawQueryService) FindByActive(req *requests.FindAllWithdraws) ([]*response.WithdrawResponseDeleteAt, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByActive", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindByActive")
	defer span.End()

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching active withdraw",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	withdraws, totalRecords, err := s.withdrawQueryRepository.FindByActive(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, "FindByActive", "FAILED_FIND_ACTIVE_WITHDRAW", span, &status, withdraw_errors.ErrFailedFindAllWithdraws, zap.Error(err))
	}

	withdrawResponses := s.mapping.ToWithdrawsResponseDeleteAt(withdraws)

	s.logger.Debug("Successfully fetched active withdraw",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return withdrawResponses, totalRecords, nil
}

func (s *withdrawQueryService) FindByTrashed(req *requests.FindAllWithdraws) ([]*response.WithdrawResponseDeleteAt, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("FindByTrashed", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindByTrashed")
	defer span.End()

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching trashed withdraw",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	withdraws, totalRecords, err := s.withdrawQueryRepository.FindByTrashed(req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, "FindByTrashed", "FAILED_FIND_TRASHED_WITHDRAW", span, &status, withdraw_errors.ErrFailedFindAllWithdraws, zap.Error(err))
	}

	withdrawResponses := s.mapping.ToWithdrawsResponseDeleteAt(withdraws)

	s.logger.Debug("Successfully fetched trashed withdraw",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return withdrawResponses, totalRecords, nil
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
