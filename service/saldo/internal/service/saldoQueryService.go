package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-saldo/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/service"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/saldo"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// saldoQueryParams contains the dependencies required to construct a saldoQueryService.
type saldoQueryParams struct {
	// Ctx is the context used throughout the service.
	Ctx context.Context

	// ErrorHandler handles domain-specific errors for saldo queries.
	ErrorHandler errorhandler.SaldoQueryErrorHandler

	// Cache provides in-memory caching for saldo query operations.
	Cache mencache.SaldoQueryCache

	// Repository provides access to the saldo query data layer.
	Repository repository.SaldoQueryRepository

	// Logger is the structured logger used by the service.
	Logger logger.LoggerInterface

	// Mapper maps domain models to response DTOs.
	Mapper responseservice.SaldoQueryResponseMapper
}

// saldoQueryService handles read-only operations for saldo data.
type saldoQueryService struct {
	// ctx is the base context shared across service methods.
	ctx context.Context

	// errorhandler handles domain-specific errors for saldo queries.
	errorhandler errorhandler.SaldoQueryErrorHandler

	// mencache provides caching functionality for saldo queries.
	mencache mencache.SaldoQueryCache

	// saldoQueryRepository provides data access for saldo-related queries.
	saldoQueryRepository repository.SaldoQueryRepository

	// logger is the structured logger used in the service.
	logger logger.LoggerInterface

	// mapper maps internal saldo entities to response DTOs.
	mapper responseservice.SaldoQueryResponseMapper

	observability observability.TraceLoggerObservability
}

// NewSaldoQueryService initializes a new saldoQueryService with the provided parameters.
//
// It sets up the prometheus metrics for counting and measuring the duration of saldo query requests.
//
// Parameters:
// - params: A pointer to a saldoQueryParams containing the necessary dependencies.
//
// Returns:
// - A pointer to a newly created saldoQueryService.
func NewSaldoQueryService(
	params *saldoQueryParams,
) SaldoQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "saldo_query_service_request_total",
			Help: "Total number of requests to the SaldoQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "saldo_query_service_request_duration_seconds",
			Help:    "Histogram of request durations for the SaldoQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("saldo-query-service"), params.Logger, requestCounter, requestDuration)

	return &saldoQueryService{
		ctx:                  params.Ctx,
		errorhandler:         params.ErrorHandler,
		mencache:             params.Cache,
		saldoQueryRepository: params.Repository,
		logger:               params.Logger,
		mapper:               params.Mapper,
		observability:        observability,
	}
}

// FindAll retrieves all saldo records with optional pagination and filtering.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing filters such as pagination, status, etc.
//
// Returns:
//   - []*response.SaldoResponse: The list of saldo responses.
//   - *int: The total number of records found.
//   - *response.ErrorResponse: An error response if the operation fails.
func (s *saldoQueryService) FindAll(ctx context.Context, req *requests.FindAllSaldos) ([]*response.SaldoResponse, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedSaldos(ctx, req); found {
		logSuccess("Successfully retrieved all saldo records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	res, totalRecords, err := s.saldoQueryRepository.FindAllSaldos(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_SALDOS", span, &status, zap.Error(err))
	}

	so := s.mapper.ToSaldoResponses(res)

	logSuccess("Successfully retrieved all saldo records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

// FindByActive retrieves all active saldo records.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing filter options such as page and page size.
//
// Returns:
//   - []*response.SaldoResponseDeleteAt: The list of active saldo records.
//   - *int: The total number of active records.
//   - *response.ErrorResponse: An error response if the operation fails.
func (s *saldoQueryService) FindByActive(ctx context.Context, req *requests.FindAllSaldos) ([]*response.SaldoResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindByActive"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedSaldoByActive(ctx, req); found {
		logSuccess("Successfully fetched active saldo from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	res, totalRecords, err := s.saldoQueryRepository.FindByActive(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_ACTIVE_SALDOS", span, &status, saldo_errors.ErrFailedFindActiveSaldos, zap.Error(err))
	}

	so := s.mapper.ToSaldoResponsesDeleteAt(res)

	s.mencache.SetCachedSaldoByActive(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched active saldo", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

// FindByTrashed retrieves all trashed (soft-deleted) saldo records.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing filter options such as page and page size.
//
// Returns:
//   - []*response.SaldoResponseDeleteAt: The list of trashed saldo records.
//   - *int: The total number of trashed records.
//   - *response.ErrorResponse: An error response if the operation fails.
func (s *saldoQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllSaldos) ([]*response.SaldoResponseDeleteAt, *int, *response.ErrorResponse) {
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	const method = "FindByTrashed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedSaldoByTrashed(ctx, req); found {
		logSuccess("Successfully fetched trashed saldo from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	res, totalRecords, err := s.saldoQueryRepository.FindByTrashed(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_TRASHED_SALDOS", span, &status, saldo_errors.ErrFailedFindTrashedSaldos, zap.Error(err))
	}
	so := s.mapper.ToSaldoResponsesDeleteAt(res)

	s.mencache.SetCachedSaldoByTrashed(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched trashed saldo", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

// FindById retrieves a saldo by its unique ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - saldo_id: The ID of the saldo to retrieve.
//
// Returns:
//   - *response.SaldoResponse: The saldo response if found.
//   - *response.ErrorResponse: An error response if the saldo is not found or an error occurs.
func (s *saldoQueryService) FindById(ctx context.Context, saldo_id int) (*response.SaldoResponse, *response.ErrorResponse) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("saldo.id", saldo_id))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedSaldoById(ctx, saldo_id); found {
		logSuccess("Successfully fetched saldo from cache", zap.Int("saldo.id", saldo_id))
		return data, nil
	}

	res, err := s.saldoQueryRepository.FindById(ctx, saldo_id)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_SALDO", span, &status, saldo_errors.ErrFailedSaldoNotFound, zap.Error(err))
	}

	so := s.mapper.ToSaldoResponse(res)

	s.mencache.SetCachedSaldoById(ctx, saldo_id, so)

	logSuccess("Successfully fetched saldo", zap.Int("saldo.id", saldo_id))

	return so, nil
}

// FindByCardNumber retrieves a saldo by its associated card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - card_number: The card number associated with the saldo.
//
// Returns:
//   - *response.SaldoResponse: The saldo response if found.
//   - *response.ErrorResponse: An error response if the saldo is not found or an error occurs.
func (s *saldoQueryService) FindByCardNumber(ctx context.Context, card_number string) (*response.SaldoResponse, *response.ErrorResponse) {
	const method = "FindByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("card_number", card_number))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedSaldoByCardNumber(ctx, card_number); found {
		logSuccess("Successfully fetched saldo by card number from cache", zap.String("card_number", card_number))
		return data, nil
	}

	res, err := s.saldoQueryRepository.FindByCardNumber(ctx, card_number)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_SALDO_BY_CARD_NUMBER", span, &status, saldo_errors.ErrFailedSaldoNotFound, zap.Error(err))
	}

	so := s.mapper.ToSaldoResponse(res)

	s.mencache.SetCachedSaldoByCardNumber(ctx, card_number, so)

	logSuccess("Successfully fetched saldo by card number", zap.String("card_number", card_number))

	return so, nil
}

// normalizePagination normalizes pagination page and pageSize arguments.
//
// If page or pageSize is less than or equal to 0, it is set to the default value of 1 or 10, respectively.
//
// Parameters:
//   - page: The input page number.
//   - pageSize: The input page size.
//
// Returns:
//   - The normalized page number.
//   - The normalized page size.
func (s *saldoQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}
