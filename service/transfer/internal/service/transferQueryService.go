package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	traceunic "github.com/MamangRust/monolith-payment-gateway-pkg/trace_unic"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/transfer_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
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
	trace                   trace.Tracer
	transferQueryRepository repository.TransferQueryRepository
	logger                  logger.LoggerInterface
	mapping                 responseservice.TransferResponseMapper
	requestCounter          *prometheus.CounterVec
	requestDuration         *prometheus.HistogramVec
}

func NewTransferQueryService(ctx context.Context, transferQueryRepository repository.TransferQueryRepository, logger logger.LoggerInterface, mapping responseservice.TransferResponseMapper) *transferQueryService {
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
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &transferQueryService{
		ctx:                     ctx,
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

	pageSize := req.PageSize
	page := req.Page
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

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	transfers, totalRecords, err := s.transferQueryRepository.FindAll(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_TRANSFERS")

		s.logger.Error("failed to find all transfer",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.String("trace_id", traceID),
		)

		span.SetAttributes(attribute.String("trace_id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find all transfer")

		status = "failed_to_find_all_transfer"

		return nil, nil, transfer_errors.ErrFailedFindAllTransfers
	}

	so := s.mapping.ToTransfersResponse(transfers)

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

	transfer, err := s.transferQueryRepository.FindById(transferId)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TRANSFER_BY_ID")

		s.logger.Error("Failed to find transfer by ID", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find transfer by ID")
		status = "failed_find_transfer_by_id"

		return nil, transfer_errors.ErrTransferNotFound
	}

	so := s.mapping.ToTransferResponse(transfer)

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

	page := req.Page
	pageSize := req.PageSize
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

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	transfers, totalRecords, err := s.transferQueryRepository.FindByActive(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ACTIVE_TRANSFERS")

		s.logger.Error("Failed to find active transfer",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.String("trace_id", traceID),
		)

		span.SetAttributes(attribute.String("trace_id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find active transfer")

		status = "failed_find_active_transfer"

		return nil, nil, transfer_errors.ErrFailedFindActiveTransfers
	}

	so := s.mapping.ToTransfersResponseDeleteAt(transfers)

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

	page := req.Page
	pageSize := req.PageSize
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

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	transfers, totalRecords, err := s.transferQueryRepository.FindByTrashed(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TRASHED_TRANSFERS")

		s.logger.Error("Failed to find trashed transfer",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.String("trace_id", traceID),
		)

		span.SetAttributes(attribute.String("trace_id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find trashed transfer")

		status = "failed_find_trashed_transfer"

		return nil, nil, transfer_errors.ErrFailedFindTrashedTransfers
	}

	so := s.mapping.ToTransfersResponseDeleteAt(transfers)

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

	res, err := s.transferQueryRepository.FindTransferByTransferFrom(transfer_from)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TRANSFER_BY_TRANSFER_FROM")

		s.logger.Error("Failed to fetch transfers by transfer_from",
			zap.Error(err),
			zap.String("transfer_from", transfer_from),
			zap.String("trace_id", traceID),
		)

		span.SetAttributes(attribute.String("trace_id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch transfers by transfer_from")

		status = "failed_find_transfer_by_transfer_from"

		return nil, transfer_errors.ErrTransferNotFound
	}

	so := s.mapping.ToTransfersResponse(res)

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

	res, err := s.transferQueryRepository.FindTransferByTransferTo(transfer_to)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TRANSFER_BY_TRANSFER_TO")

		s.logger.Error("Failed to fetch transfers by transfer_to",
			zap.Error(err),
			zap.String("transfer_to", transfer_to),
			zap.String("trace_id", traceID),
		)

		span.SetAttributes(attribute.String("trace_id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch transfers by transfer_to")

		status = "failed_find_transfer_by_transfer_to"
		return nil, transfer_errors.ErrTransferNotFound
	}

	so := s.mapping.ToTransfersResponse(res)

	s.logger.Debug("Successfully fetched transfer record by transfer_to",
		zap.String("transfer_to", transfer_to),
	)

	return so, nil
}

func (s *transferQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
