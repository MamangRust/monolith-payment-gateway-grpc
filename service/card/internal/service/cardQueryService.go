package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-card/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-card/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-card/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/card"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// cardQueryServiceDeps holds the dependencies required to initialize the cardQueryService.
// This struct is used during service construction and supports dependency injection.
type cardQueryServiceDeps struct {
	// ErrorHandler handles domain-specific errors related to card queries.
	ErrorHandler errorhandler.CardQueryErrorHandler

	// Cache provides in-memory or Redis-based caching for card query data.
	Cache mencache.CardQueryCache

	// CardQueryRepository provides access to the card-related data from the data store.
	CardQueryRepository repository.CardQueryRepository

	// Logger is used for structured logging of operations and errors.
	Logger logger.LoggerInterface

	// Mapper maps internal domain models to response DTOs.
	Mapper responseservice.CardQueryResponseMapper
}

// cardQueryService implements the CardQueryService interface.
// It handles business logic for querying card-related data,
// with support for caching, tracing, metrics, and logging.
type cardQueryService struct {
	// errorhandler processes errors specific to card query logic.
	errorhandler errorhandler.CardQueryErrorHandler

	// mencache provides caching to improve response time and reduce database load.
	mencache mencache.CardQueryCache

	// cardQueryRepository provides methods to retrieve card data from the data source.
	cardQueryRepository repository.CardQueryRepository

	// logger logs information, errors, and operational metrics.
	logger logger.LoggerInterface

	// mapper transforms internal models into external-facing response formats.
	mapper responseservice.CardQueryResponseMapper

	observability observability.TraceLoggerObservability
}

// NewCardQueryService initializes a new instance of cardQueryService with the provided parameters.
//
// It sets up Prometheus metrics for counting and measuring the duration of card query requests.
//
// Parameters:
// - params: A pointer to cardQueryServiceDeps containing the necessary dependencies.
//
// Returns:
// - A pointer to a newly created cardQueryService.
func NewCardQueryService(
	params *cardQueryServiceDeps,
) CardQueryService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "card_query_request_count",
		Help: "Number of card query requests CardQueryService",
	}, []string{"method", "status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "card_query_request_duration_seconds",
		Help:    "Duration of card query requests CardQueryService",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("card-query-service"), params.Logger, requestCounter, requestDuration)

	return &cardQueryService{
		errorhandler:        params.ErrorHandler,
		cardQueryRepository: params.CardQueryRepository,
		logger:              params.Logger,
		mapper:              params.Mapper,
		observability:       observability,
	}
}

// FindAll retrieves a paginated list of card records based on the search criteria
// specified in the request. It queries the database and returns a slice of CardResponse,
// the total count of records, and an error if any occurred.
//
// Parameters:
//   - req: A FindAllCards request object containing the search parameters
//     such as search keyword, page number, and page size.
//
// Returns:
//   - A slice of CardResponse representing the card records fetched from the database.
//   - A pointer to an int representing the total number of records matching the search criteria.
//   - An ErrorResponse if the operation fails, nil otherwise.
func (s *cardQueryService) FindAll(ctx context.Context, req *requests.FindAllCards) ([]*response.CardResponse, *int, *response.ErrorResponse) {
	const method = "FindAll"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetFindAllCache(ctx, req); found {
		logSuccess("Successfully fetched card records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	cards, totalRecords, err := s.cardQueryRepository.FindAllCards(ctx, req)

	if err != nil {
		return s.errorhandler.HandleFindAllError(err, method, "FAILED_FIND_ALL_CARD", span, &status, zap.Error(err))
	}

	responseData := s.mapper.ToCardsResponse(cards)

	s.mencache.SetFindAllCache(ctx, req, responseData, totalRecords)

	logSuccess("Successfully fetched card records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return responseData, totalRecords, nil
}

// FindByActive retrieves a paginated list of active card records based on the search criteria
// specified in the request. It queries the database and returns a slice of CardResponseDeleteAt,
// the total count of records, and an error if any occurred.
//
// Parameters:
//   - req: A FindAllCards request object containing the search parameters
//     such as search keyword, page number, and page size.
//
// Returns:
//   - A slice of CardResponseDeleteAt representing the card records fetched from the database.
//   - A pointer to an int representing the total number of records matching the search criteria.
//   - An ErrorResponse if the operation fails, nil otherwise.
func (s *cardQueryService) FindByActive(ctx context.Context, req *requests.FindAllCards) ([]*response.CardResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByActive"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetByActiveCache(ctx, req); found {
		logSuccess("Successfully fetched active card records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	res, totalRecords, err := s.cardQueryRepository.FindByActive(ctx, req)

	if err != nil {
		return s.errorhandler.HandleFindByActiveError(err, method, "FAILED_FIND_ACTIVE_CARD", span, &status, zap.Error(err))
	}

	responseData := s.mapper.ToCardsResponseDeleteAt(res)

	s.mencache.SetByActiveCache(ctx, req, responseData, totalRecords)

	logSuccess("Successfully fetched active card records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return responseData, totalRecords, nil
}

// FindByTrashed retrieves a paginated list of trashed card records based on the search
// criteria specified in the request. It queries the database and returns a slice of
// CardResponseDeleteAt, the total count of records, and an error if any occurred.
//
// Parameters:
//   - req: A FindAllCards request object containing the search parameters
//     such as search keyword, page number, and page size.
//
// Returns:
//   - A slice of CardResponseDeleteAt representing the trashed card records fetched from the database.
//   - A pointer to an int representing the total number of records matching the search criteria.
//   - An ErrorResponse if the operation fails, nil otherwise.
func (s *cardQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllCards) ([]*response.CardResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByTrashed"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetByTrashedCache(ctx, req); found {
		logSuccess("Successfully fetched trashed card records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	res, totalRecords, err := s.cardQueryRepository.FindByTrashed(ctx, req)
	if err != nil {
		return s.errorhandler.HandleFindByTrashedError(err, method, "FAILED_FIND_TRASHED_CARD", span, &status, zap.Error(err))
	}

	responseData := s.mapper.ToCardsResponseDeleteAt(res)

	s.mencache.SetByTrashedCache(ctx, req, responseData, totalRecords)

	logSuccess("Successfully fetched trashed card records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return responseData, totalRecords, nil
}

// FindById retrieves a card record by its ID from the database.
//
// Parameters:
//   - card_id: The ID of the card to be retrieved.
//
// Returns:
//   - A pointer to a CardResponse representing the card record fetched from the database.
//   - An ErrorResponse if the operation fails, nil otherwise.
func (s *cardQueryService) FindById(ctx context.Context, card_id int) (*response.CardResponse, *response.ErrorResponse) {
	const method = "FindByActive"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("card.id", card_id))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetByIdCache(ctx, card_id); found {
		logSuccess("Successfully fetched card from cache", zap.Int("card.id", card_id))
		return data, nil
	}

	res, err := s.cardQueryRepository.FindById(ctx, card_id)

	if err != nil {
		return s.errorhandler.HandleFindByIdError(err, method, "FAILED_TO_FIND_CARD", span, &status, zap.Error(err))
	}

	so := s.mapper.ToCardResponse(res)

	s.mencache.SetByIdCache(ctx, card_id, so)

	logSuccess("Successfully fetched card", zap.Int("card.id", so.ID))

	return so, nil
}

// FindByUserID retrieves a card record associated with a user ID from the database.
//
// Parameters:
//   - user_id: The ID of the user who owns the card.
//
// Returns:
//   - A pointer to a CardResponse representing the card record fetched from the database.
//   - An ErrorResponse if the operation fails, nil otherwise.
func (s *cardQueryService) FindByUserID(ctx context.Context, user_id int) (*response.CardResponse, *response.ErrorResponse) {
	const method = "FindByUserId"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("user.id", user_id))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetByUserIDCache(ctx, user_id); found {
		logSuccess("Successfully fetched card records by user ID from cache", zap.Int("user.id", user_id))
		return data, nil
	}

	res, err := s.cardQueryRepository.FindCardByUserId(ctx, user_id)

	if err != nil {
		return s.errorhandler.HandleFindByUserIdError(err, method, "FAILED_FIND_CARD_BY_USER_ID", span, &status, zap.Error(err))
	}

	so := s.mapper.ToCardResponse(res)

	s.mencache.SetByUserIDCache(ctx, user_id, so)

	logSuccess("Successfully fetched card records by user ID", zap.Int("user.id", user_id))

	return so, nil
}

// FindByCardNumber retrieves a card record associated with a card number from the database.
//
// Parameters:
//   - card_number: The card number of the card to be retrieved.
//
// Returns:
//   - A pointer to a CardResponse representing the card record fetched from the database.
//   - An ErrorResponse if the operation fails, nil otherwise.
func (s *cardQueryService) FindByCardNumber(ctx context.Context, card_number string) (*response.CardResponse, *response.ErrorResponse) {
	const method = "FindByCardNumber"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("card.card_number", card_number))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetByCardNumberCache(ctx, card_number); found {
		logSuccess("Successfully fetched card record by card number from cache", zap.String("card_number", card_number))
		return data, nil
	}

	res, err := s.cardQueryRepository.FindCardByCardNumber(ctx, card_number)

	if err != nil {
		return s.errorhandler.HandleFindByCardNumberError(err, method, "FAILED_FIND_CARD_BY_CARD_NUMBER", span, &status, zap.Error(err))
	}

	so := s.mapper.ToCardResponse(res)

	s.mencache.SetByCardNumberCache(ctx, card_number, so)

	logSuccess("Successfully fetched card record by card number", zap.String("card_number", card_number))

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
func (s *cardQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}
