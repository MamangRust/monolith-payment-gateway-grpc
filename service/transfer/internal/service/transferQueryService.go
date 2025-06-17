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
	"go.opentelemetry.io/otel/codes"
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
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindAll"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedTransfersCache(req); found {
		logSuccess("Successfully retrieved all transfer records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	transfers, totalRecords, err := s.transferQueryRepository.FindAll(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_TO_FIND_ALL_TRANSFERS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransfersResponse(transfers)

	s.mencache.SetCachedTransfersCache(req, so, totalRecords)

	logSuccess("Successfully retrieved all transfer records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *transferQueryService) FindById(transferId int) (*response.TransferResponse, *response.ErrorResponse) {
	const method = "FindById"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("transfer.id", transferId))

	defer func() {
		end(status)
	}()

	if data := s.mencache.GetCachedTransferCache(transferId); data != nil {
		logSuccess("Successfully fetched transfer from cache", zap.Int("transfer.id", transferId))
		return data, nil
	}

	transfer, err := s.transferQueryRepository.FindById(transferId)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_TO_FIND_TRANSFER_BY_ID", span, &status, transfer_errors.ErrTransferNotFound, zap.Error(err))
	}

	so := s.mapping.ToTransferResponse(transfer)

	s.mencache.SetCachedTransferCache(so)

	logSuccess("Successfully fetched transfer", zap.Int("transfer.id", transferId))

	return so, nil
}

func (s *transferQueryService) FindByActive(req *requests.FindAllTranfers) ([]*response.TransferResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindByActive"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedTransferActiveCache(req); found {
		logSuccess("Successfully retrieved active transfer records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	transfers, totalRecords, err := s.transferQueryRepository.FindByActive(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_TO_FIND_BY_ACTIVE_TRANSFERS", span, &status, transfer_errors.ErrFailedFindActiveTransfers, zap.Error(err))
	}

	so := s.mapping.ToTransfersResponseDeleteAt(transfers)

	s.mencache.SetCachedTransferActiveCache(req, so, totalRecords)

	logSuccess("Successfully retrieved active transfer records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *transferQueryService) FindByTrashed(req *requests.FindAllTranfers) ([]*response.TransferResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindByTrashed"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedTransferTrashedCache(req); found {
		logSuccess("Successfully retrieved trashed transfer records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	transfers, totalRecords, err := s.transferQueryRepository.FindByTrashed(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_TO_FIND_BY_TRASHED_TRANSFERS", span, &status, transfer_errors.ErrFailedFindTrashedTransfers, zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))
	}

	so := s.mapping.ToTransfersResponseDeleteAt(transfers)

	s.mencache.SetCachedTransferTrashedCache(req, so, totalRecords)

	logSuccess("Successfully retrieved trashed transfer records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *transferQueryService) FindTransferByTransferFrom(transfer_from string) ([]*response.TransferResponse, *response.ErrorResponse) {
	const method = "FindTransferByTransferFrom"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.String("transaction.from", transfer_from))

	defer func() {
		end(status)
	}()

	if data := s.mencache.GetCachedTransferByFrom(transfer_from); data != nil {
		logSuccess("Successfully fetched transfer from cache", zap.String("transfer_from", transfer_from))
		return data, nil
	}

	res, err := s.transferQueryRepository.FindTransferByTransferFrom(transfer_from)

	if err != nil {
		return s.errorhandler.HanldeRepositoryListError(err, method, "FAILED_TO_FIND_TRANSFER_BY_TRANSFER_FROM", span, &status, transfer_errors.ErrTransferNotFound, zap.Error(err))
	}

	so := s.mapping.ToTransfersResponse(res)

	s.mencache.SetCachedTransferByFrom(transfer_from, so)

	logSuccess("Successfully fetched transfer", zap.String("transfer_from", transfer_from))

	return so, nil
}

func (s *transferQueryService) FindTransferByTransferTo(transfer_to string) ([]*response.TransferResponse, *response.ErrorResponse) {
	const method = "FindTransferByTransferTo"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.String("transfer.to", transfer_to))

	defer func() {
		end(status)
	}()

	if data := s.mencache.GetCachedTransferByTo(transfer_to); data != nil {
		logSuccess("Successfully fetched transfer from cache", zap.String("transfer_to", transfer_to))
		return data, nil
	}

	res, err := s.transferQueryRepository.FindTransferByTransferTo(transfer_to)

	if err != nil {
		return s.errorhandler.HanldeRepositoryListError(err, method, "FAILED_TO_FIND_TRANSFER_BY_TRANSFER_TO", span, &status, transfer_errors.ErrTransferNotFound, zap.Error(err))
	}

	so := s.mapping.ToTransfersResponse(res)

	s.mencache.SetCachedTransferByTo(transfer_to, so)

	logSuccess("Successfully fetched transfer", zap.String("transfer_to", transfer_to))

	return so, nil
}

func (s *transferQueryService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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
