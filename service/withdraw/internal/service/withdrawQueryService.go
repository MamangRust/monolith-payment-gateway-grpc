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

type withdrawQueryService struct {
	ctx                     context.Context
	trace                   trace.Tracer
	withdrawQueryRepository repository.WithdrawQueryRepository
	logger                  logger.LoggerInterface
	mapping                 responseservice.WithdrawResponseMapper
	requestCounter          *prometheus.CounterVec
	requestDuration         *prometheus.HistogramVec
}

func NewWithdrawQueryService(ctx context.Context, withdrawQueryRepository repository.WithdrawQueryRepository, logger logger.LoggerInterface, mapping responseservice.WithdrawResponseMapper) *withdrawQueryService {
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
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &withdrawQueryService{
		ctx:                     ctx,
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

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching withdraw",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	withdraws, totalRecords, err := s.withdrawQueryRepository.FindAll(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_WITHDRAW")

		s.logger.Error("Failed to fetch withdraw",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch withdraw")
		status = "failed_find_all_withdraw"

		return nil, nil, withdraw_errors.ErrFailedFindAllWithdraws
	}

	withdrawResponse := s.mapping.ToWithdrawsResponse(withdraws)

	s.logger.Debug("Successfully fetched withdraw",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

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

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching withdraw",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	withdraws, totalRecords, err := s.withdrawQueryRepository.FindAllByCardNumber(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_WITHDRAW")

		s.logger.Error("Failed to fetch withdraw",
			zap.Error(err),
			zap.String("traceID", traceID),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch withdraw")
		status = "failed_find_all_withdraw"

		return nil, nil, withdraw_errors.ErrFailedFindAllWithdrawsByCard
	}

	withdrawResponse := s.mapping.ToWithdrawsResponse(withdraws)

	s.logger.Debug("Successfully fetched withdraw",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

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

	span.SetAttributes(
		attribute.Int("withdraw_id", withdrawID),
	)

	s.logger.Debug("Fetching withdraw by ID", zap.Int("withdraw_id", withdrawID))

	withdraw, err := s.withdrawQueryRepository.FindById(withdrawID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_WITHDRAW_BY_ID")

		s.logger.Error("Failed to fetch withdraw by ID",
			zap.Error(err),
			zap.String("traceID", traceID),
			zap.Int("withdraw_id", withdrawID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch withdraw by ID")
		status = "failed_find_withdraw_by_id"

		return nil, withdraw_errors.ErrWithdrawNotFound
	}
	so := s.mapping.ToWithdrawResponse(withdraw)

	s.logger.Debug("Successfully fetched withdraw", zap.Int("withdraw_id", withdrawID))

	return so, nil
}

func (s *withdrawQueryService) FindByActive(req *requests.FindAllWithdraws) ([]*response.WithdrawResponseDeleteAt, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByActive", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindByActive")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
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

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	withdraws, totalRecords, err := s.withdrawQueryRepository.FindByActive(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ACTIVE_WITHDRAW")

		s.logger.Error("Failed to fetch active withdraw",
			zap.Error(err),
			zap.String("traceID", traceID),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch active withdraw")
		status = "failed_find_active_withdraw"

		return nil, nil, withdraw_errors.ErrFailedFindActiveWithdraws
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

	page := req.Page
	pageSize := req.PageSize
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

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	withdraws, totalRecords, err := s.withdrawQueryRepository.FindByTrashed(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TRASHED_WITHDRAW")

		s.logger.Error("Failed to fetch trashed withdraw",
			zap.Error(err),
			zap.String("traceID", traceID),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch trashed withdraw")
		status = "failed_find_trashed_withdraw"

		return nil, nil, withdraw_errors.ErrFailedFindTrashedWithdraws
	}

	withdrawResponses := s.mapping.ToWithdrawsResponseDeleteAt(withdraws)

	s.logger.Debug("Successfully fetched trashed withdraw",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return withdrawResponses, totalRecords, nil
}

func (s *withdrawQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
