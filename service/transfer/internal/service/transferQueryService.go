package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-transfer/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-transfer/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transferQueryService struct {
	ctx                     context.Context
	errorhandler            errorhandler.TransferQueryErrorHandler
	mencache                mencache.TransferQueryCache
	trace                   trace.Tracer
	transferQueryRepository repository.TransferQueryRepository
	logger                  logger.LoggerInterface
	mapping                 responseservice.TransferResponseMapper
	requestCounter          *prometheus.CounterVec
	requestDuration         *prometheus.HistogramVec
}

func NewTransferQueryService(ctx context.Context, errorhandler errorhandler.TransferQueryErrorHandler,
	mencache mencache.TransferQueryCache, transferQueryRepository repository.TransferQueryRepository, logger logger.LoggerInterface, mapping responseservice.TransferResponseMapper) *transferQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transfer_query_service_request_total",
			Help: "Total number of requests to the TransferQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transfer_query_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TransferQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &transferQueryService{
		ctx:                     ctx,
		errorhandler:            errorhandler,
		mencache:                mencache,
		trace:                   otel.Tracer("transfer-query-service"),
		transferQueryRepository: transferQueryRepository,
		logger:                  logger,
		mapping:                 mapping,
		requestCounter:          requestCounter,
		requestDuration:         requestDuration,
	}
}

func (s *transferQueryService) FindAll(req *requests.FindAllTranfers) ([]*response.TransferResponse, *int, *response.ErrorResponse) {
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

	s.logger.Debug("Fetching transfer",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if data, total, found := s.mencache.GetCachedTransfersCache(req); found {
		s.logger.Debug("Successfully fetched transfers from cache",
			zap.Int("totalRecords", *total))
		return data, total, nil
	}

	transfers, totalRecords, err := s.transferQueryRepository.FindAll(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, "FindAll", "FAILED_TO_FIND_ALL_TRANSFERS", span, &status)
	}

	so := s.mapping.ToTransfersResponse(transfers)

	s.mencache.SetCachedTransfersCache(req, so, totalRecords)

	s.logger.Debug("Successfully fetched transfer",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *transferQueryService) FindById(transferId int) (*response.TransferResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindById", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindById")
	defer span.End()

	span.SetAttributes(
		attribute.Int("transfer_id", transferId),
	)

	s.logger.Debug("Fetching transfer by ID", zap.Int("transfer_id", transferId))

	if data := s.mencache.GetCachedTransferCache(transferId); data != nil {
		s.logger.Debug("Successfully fetched transfer from cache",
			zap.Int("transfer_id", transferId))
		return data, nil
	}

	transfer, err := s.transferQueryRepository.FindById(transferId)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "FindById", "FAILED_TO_FIND_TRANSFER_BY_ID", span, &status, transfer_errors.ErrTransferNotFound, zap.Int("transfer_id", transferId))
	}

	so := s.mapping.ToTransferResponse(transfer)

	s.mencache.SetCachedTransferCache(so)

	s.logger.Debug("Successfully fetched transfer", zap.Int("transfer_id", transferId))

	return so, nil
}

func (s *transferQueryService) FindByActive(req *requests.FindAllTranfers) ([]*response.TransferResponseDeleteAt, *int, *response.ErrorResponse) {
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

	s.logger.Debug("Fetching active transfer",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if data, total, found := s.mencache.GetCachedTransferActiveCache(req); found {
		s.logger.Debug("Successfully fetched active transfers from cache",
			zap.Int("totalRecords", *total))
		return data, total, nil
	}

	transfers, totalRecords, err := s.transferQueryRepository.FindByActive(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, "FindByActive", "FAILED_TO_FIND_BY_ACTIVE_TRANSFERS", span, &status, transfer_errors.ErrFailedFindActiveTransfers, zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))
	}

	so := s.mapping.ToTransfersResponseDeleteAt(transfers)

	s.mencache.SetCachedTransferActiveCache(req, so, totalRecords)

	s.logger.Debug("Successfully fetched active transfer",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *transferQueryService) FindByTrashed(req *requests.FindAllTranfers) ([]*response.TransferResponseDeleteAt, *int, *response.ErrorResponse) {
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

	s.logger.Debug("Fetching trashed transfer",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if data, total, found := s.mencache.GetCachedTransferTrashedCache(req); found {
		s.logger.Debug("Successfully fetched trashed transfers from cache",
			zap.Int("totalRecords", *total))
		return data, total, nil
	}

	transfers, totalRecords, err := s.transferQueryRepository.FindByTrashed(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, "FindByTrashed", "FAILED_TO_FIND_BY_TRASHED_TRANSFERS", span, &status, transfer_errors.ErrFailedFindTrashedTransfers, zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))
	}

	so := s.mapping.ToTransfersResponseDeleteAt(transfers)

	s.mencache.SetCachedTransferTrashedCache(req, so, totalRecords)

	s.logger.Debug("Successfully fetched trashed transfer",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *transferQueryService) FindTransferByTransferFrom(transfer_from string) ([]*response.TransferResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindTransferByTransferFrom", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindTransferByTransferFrom")
	defer span.End()

	span.SetAttributes(
		attribute.String("transfer_from", transfer_from),
	)

	s.logger.Debug("Starting fetch transfer by transfer_from",
		zap.String("transfer_from", transfer_from),
	)

	if data := s.mencache.GetCachedTransferByFrom(transfer_from); data != nil {
		s.logger.Debug("Successfully fetched transfer by transfer_from from cache", zap.String("transfer_from", transfer_from))
		return data, nil
	}

	res, err := s.transferQueryRepository.FindTransferByTransferFrom(transfer_from)

	if err != nil {
		return s.errorhandler.HanldeRepositoryListError(err, "FindTransferByTransferFrom", "FAILED_TO_FIND_TRANSFER_BY_TRANSFER_FROM", span, &status, transfer_errors.ErrTransferNotFound, zap.String("transfer_from", transfer_from))
	}

	so := s.mapping.ToTransfersResponse(res)

	s.mencache.SetCachedTransferByFrom(transfer_from, so)

	s.logger.Debug("Successfully fetched transfer record by transfer_from",
		zap.String("transfer_from", transfer_from),
	)

	return so, nil
}

func (s *transferQueryService) FindTransferByTransferTo(transfer_to string) ([]*response.TransferResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindTransferByTransferTo", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindTransferByTransferTo")
	defer span.End()

	span.SetAttributes(
		attribute.String("transfer_to", transfer_to),
	)

	s.logger.Debug("Starting fetch transfer by transfer_to",
		zap.String("transfer_to", transfer_to),
	)

	if data := s.mencache.GetCachedTransferByTo(transfer_to); data != nil {
		s.logger.Debug("Successfully fetched transfer by transfer_to from cache", zap.String("transfer_to", transfer_to))
		return data, nil
	}

	res, err := s.transferQueryRepository.FindTransferByTransferTo(transfer_to)

	if err != nil {
		return s.errorhandler.HanldeRepositoryListError(err, "FindTransferByTransferTo", "FAILED_TO_FIND_TRANSFER_BY_TRANSFER_TO", span, &status, transfer_errors.ErrTransferNotFound, zap.String("transfer_to", transfer_to))
	}

	so := s.mapping.ToTransfersResponse(res)

	s.mencache.SetCachedTransferByTo(transfer_to, so)

	s.logger.Debug("Successfully fetched transfer record by transfer_to",
		zap.String("transfer_to", transfer_to),
	)

	return so, nil
}

func (s *transferQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func (s *transferQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
