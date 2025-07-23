package handler

import (
	"context"
	"math"

	pbhelper "github.com/MamangRust/monolith-payment-gateway-pb"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/user"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	user_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors/grpc"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/user"
	"github.com/MamangRust/monolith-payment-gateway-user/internal/service"
	"go.uber.org/zap"
)

type userQueryHandleGrpc struct {
	pb.UnimplementedUserQueryServiceServer

	userQueryService service.UserQueryService

	logger logger.LoggerInterface

	mapper protomapper.UserQueryProtoMapper
}

func NewUserQueryHandleGrpc(query service.UserQueryService, logger logger.LoggerInterface, mapper protomapper.UserQueryProtoMapper) UserQueryHandleGrpc {
	return &userQueryHandleGrpc{
		userQueryService: query,
		logger:           logger,
		mapper:           mapper,
	}
}

// FindAll retrieves a paginated list of users based on the given request parameters.
// It supports pagination and search functionalities.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - request: A pointer to a FindAllUserRequest containing pagination and search details.
//
// Returns:
//   - A pointer to ApiResponsePaginationUser containing the list of users and pagination metadata.
//   - An error if the retrieval operation fails.
func (s *userQueryHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllUserRequest) (*pb.ApiResponsePaginationUser, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	s.logger.Debug("Fetching users", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := &requests.FindAllUsers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.userQueryService.FindAll(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to fetch users", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbhelper.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapper.ToProtoResponsePaginationUser(paginationMeta, "success", "Successfully fetched users", users)
	return so, nil
}

// FindById retrieves a user record by ID.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - request: A pointer to a FindByIdUserRequest containing the user ID.
//
// Returns:
//   - A pointer to ApiResponseUser containing the user record.
//   - An error if the retrieval operation fails.
func (s *userQueryHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUser, error) {
	id := int(request.GetId())

	s.logger.Debug("Fetching user by id", zap.Int("id", id))

	if id == 0 {
		s.logger.Error("Failed to find user by id", zap.Int("id", id))
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	user, err := s.userQueryService.FindByID(ctx, id)

	if err != nil {
		s.logger.Error("Failed to find user by id", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseUser("success", "Successfully fetched user", user)

	return so, nil
}

// FindByActive is a gRPC handler that fetches active user records according to the given request.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - request: A pointer to a FindAllUserRequest containing the pagination and search parameters.
//
// Returns:
//   - A pointer to ApiResponsePaginationUserDeleteAt containing the pagination metadata and the fetched user records.
//   - An error if the retrieval operation fails.
func (s *userQueryHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllUserRequest) (*pb.ApiResponsePaginationUserDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	s.logger.Debug("Fetching active users", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := &requests.FindAllUsers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.userQueryService.FindByActive(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to fetch active users", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbhelper.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapper.ToProtoResponsePaginationUserDeleteAt(paginationMeta, "success", "Successfully fetched active users", users)

	return so, nil
}

// FindByTrashed is a gRPC handler that fetches trashed user records according to the given request.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - request: A pointer to a FindAllUserRequest containing the pagination and search parameters.
//
// Returns:
//   - A pointer to ApiResponsePaginationUserDeleteAt containing the pagination metadata and the fetched user records.
//   - An error if the retrieval operation fails.
func (s *userQueryHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllUserRequest) (*pb.ApiResponsePaginationUserDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	s.logger.Debug("Fetching trashed users", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := &requests.FindAllUsers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.userQueryService.FindByTrashed(ctx, reqService)

	if err != nil {
		s.logger.Error("Failed to fetch trashed users", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbhelper.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapper.ToProtoResponsePaginationUserDeleteAt(paginationMeta, "success", "Successfully fetched trashed users", users)

	return so, nil
}
