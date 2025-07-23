package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-payment-gateway-card/internal/service"
	pbhelpers "github.com/MamangRust/monolith-payment-gateway-pb"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/card"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/grpc"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/card"
	"go.uber.org/zap"
)

type cardQueryHandleGrpc struct {
	pb.UnimplementedCardQueryServiceServer

	cardQuery service.CardQueryService

	logger logger.LoggerInterface

	mapper protomapper.CardQueryProtoMapper
}

func NewCardQueryHandleGrpc(cardQuery service.CardQueryService, logger logger.LoggerInterface, mapper protomapper.CardQueryProtoMapper) CardQueryService {
	return &cardQueryHandleGrpc{
		cardQuery: cardQuery,
		logger:    logger,
		mapper:    mapper,
	}
}

// FindAllCard retrieves a paginated list of card records based on the search criteria
// specified in the request. It handles pagination and search functionality.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindAllCardRequest object containing the search parameters such as
//     search keyword, page number, and page size.
//
// Returns:
//   - An ApiResponsePaginationCard containing the paginated list of card records.
//   - An error if the operation fails.
func (s *cardQueryHandleGrpc) FindAllCard(ctx context.Context, req *pb.FindAllCardRequest) (*pb.ApiResponsePaginationCard, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCards{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	s.logger.Info("Fetching card records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	cards, totalRecords, err := s.cardQuery.FindAll(ctx, &reqService)

	if err != nil {
		s.logger.Error("FindAllCard failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbhelpers.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapper.ToProtoResponsePaginationCard(paginationMeta, "success", "Successfully fetched card records", cards)

	s.logger.Info("Successfully fetched card records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, nil
}

// FindByIdCard retrieves a card record by its ID from the database.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindByIdCardRequest object containing the card ID to be retrieved.
//
// Returns:
//   - An ApiResponseCard containing the card record fetched from the database.
//   - An error if the operation fails.
func (s *cardQueryHandleGrpc) FindByIdCard(ctx context.Context, req *pb.FindByIdCardRequest) (*pb.ApiResponseCard, error) {
	id := int(req.GetCardId())

	s.logger.Info("Fetching card record", zap.Int("card.id", id))

	if id == 0 {
		s.logger.Error("FindByIdCard failed", zap.Any("error", card_errors.ErrGrpcInvalidCardID))
		return nil, card_errors.ErrGrpcInvalidCardID
	}

	card, err := s.cardQuery.FindById(ctx, id)
	if err != nil {
		s.logger.Error("FindByIdCard failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseCard("success", "Successfully fetched card record", card)

	s.logger.Info("Successfully fetched card record", zap.Bool("success", true))

	return so, nil
}

// FindByUserIdCard retrieves a card record associated with a user ID.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindByUserIdCardRequest object containing the user ID to fetch the card record for.
//
// Returns:
//   - An ApiResponseCard containing the card record fetched from the database.
//   - An error if the operation fails, or if the provided user ID is invalid.
func (s *cardQueryHandleGrpc) FindByUserIdCard(ctx context.Context, req *pb.FindByUserIdCardRequest) (*pb.ApiResponseCard, error) {
	id := int(req.GetUserId())

	s.logger.Info("Fetching card record", zap.Int("user.id", id))

	if id == 0 {
		s.logger.Error("FindByUserIdCard failed", zap.Any("error", card_errors.ErrGrpcInvalidUserID))
		return nil, card_errors.ErrGrpcInvalidUserID
	}
	res, err := s.cardQuery.FindByUserID(ctx, id)

	if err != nil {
		s.logger.Error("FindByUserIdCard failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseCard("success", "Successfully fetched card record", res)

	s.logger.Info("Successfully fetched card record", zap.Bool("success", true))

	return so, nil
}

// FindByActiveCard retrieves a paginated list of active card records based on the search criteria
// specified in the request. It handles pagination and search functionality.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindAllCardRequest object containing the search parameters such as
//     search keyword, page number, and page size.
//
// Returns:
//   - An ApiResponsePaginationCardDeleteAt containing the paginated list of active card records.
//   - An error if the operation fails.
func (s *cardQueryHandleGrpc) FindByActiveCard(ctx context.Context, req *pb.FindAllCardRequest) (*pb.ApiResponsePaginationCardDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	s.logger.Info("Fetching card records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCards{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.cardQuery.FindByActive(ctx, &reqService)

	if err != nil {
		s.logger.Error("FindByActiveCard failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbhelpers.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapper.ToProtoResponsePaginationCardDeletedAt(paginationMeta, "success", "Successfully fetched card record", res)

	s.logger.Info("Successfully fetched card record", zap.Bool("success", true))

	return so, nil
}

// FindByTrashedCard retrieves a paginated list of trashed card records based on the search criteria
// specified in the request. It handles pagination and search functionality.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindAllCardRequest object containing the search parameters such as
//     search keyword, page number, and page size.
//
// Returns:
//   - An ApiResponsePaginationCardDeleteAt containing the paginated list of trashed card records.
//   - An error if the operation fails.
func (s *cardQueryHandleGrpc) FindByTrashedCard(ctx context.Context, req *pb.FindAllCardRequest) (*pb.ApiResponsePaginationCardDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	s.logger.Info("Fetching card records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCards{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.cardQuery.FindByTrashed(ctx, &reqService)

	if err != nil {
		s.logger.Error("FindByTrashedCard failed", zap.Any("error", err))

		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbhelpers.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapper.ToProtoResponsePaginationCardDeletedAt(paginationMeta, "success", "Successfully fetched card record", res)

	s.logger.Info("Successfully fetched card record", zap.Bool("success", true))

	return so, nil
}

// FindByCardNumber retrieves a card record associated with a given card number.
//
// Parameters:
//   - ctx: The context for request-scoped values, cancellation, and deadlines.
//   - req: A FindByCardNumberRequest object containing the card number to fetch the card record for.
//
// Returns:
//   - An ApiResponseCard containing the card record fetched from the database.
//   - An error if the operation fails, or if the provided card number is invalid.
func (s *cardQueryHandleGrpc) FindByCardNumber(ctx context.Context, req *pbhelpers.FindByCardNumberRequest) (*pb.ApiResponseCard, error) {
	card_number := req.GetCardNumber()

	s.logger.Info("Fetching card records", zap.String("card_number", card_number))

	if card_number == "" {
		s.logger.Error("FindByCardNumber failed", zap.Any("error", card_errors.ErrGrpcInvalidCardNumber))
		return nil, card_errors.ErrGrpcInvalidCardNumber
	}

	res, err := s.cardQuery.FindByCardNumber(ctx, card_number)

	if err != nil {
		s.logger.Error("FindByCardNumber failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseCard("success", "Successfully fetched card record", res)

	s.logger.Info("Successfully fetched card record", zap.Bool("success", true))

	return so, nil

}
