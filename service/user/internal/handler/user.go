package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto"
	"github.com/MamangRust/monolith-payment-gateway-shared/pb"
	"github.com/MamangRust/monolith-payment-gateway-user/internal/service"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userHandleGrpc struct {
	pb.UnimplementedUserServiceServer
	userQueryService   service.UserQueryService
	userCommandService service.UserCommandService
	logger             logger.LoggerInterface
	mapping            protomapper.UserProtoMapper
}

func NewUserHandleGrpc(user service.Service, logger logger.LoggerInterface) *userHandleGrpc {
	return &userHandleGrpc{
		userQueryService:   user.UserQuery,
		userCommandService: user.UserCommand,
		logger:             logger,
		mapping:            protomapper.NewUserProtoMapper(),
	}
}

func (s *userHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllUserRequest) (*pb.ApiResponsePaginationUser, error) {
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

	reqService := requests.FindAllUsers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.userQueryService.FindAll(&reqService)

	if err != nil {
		s.logger.Error("Failed to fetch users", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationUser(paginationMeta, "success", "Successfully fetched users", users)
	return so, nil
}

func (s *userHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUser, error) {
	id := int(request.GetId())

	s.logger.Debug("Fetching user by id", zap.Int("id", id))

	if id == 0 {
		s.logger.Error("Failed to find user by id", zap.Int("id", id))
		return nil, user_errors.ErrGrpcUserNotFound
	}

	user, err := s.userQueryService.FindByID(id)

	if err != nil {
		s.logger.Error("Failed to find user by id", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseUser("success", "Successfully fetched user", user)

	return so, nil
}

func (s *userHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllUserRequest) (*pb.ApiResponsePaginationUserDeleteAt, error) {
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

	reqService := requests.FindAllUsers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.userQueryService.FindByActive(&reqService)

	if err != nil {
		s.logger.Error("Failed to fetch active users", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapping.ToProtoResponsePaginationUserDeleteAt(paginationMeta, "success", "Successfully fetched active users", users)

	return so, nil
}

func (s *userHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllUserRequest) (*pb.ApiResponsePaginationUserDeleteAt, error) {
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

	reqService := requests.FindAllUsers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.userQueryService.FindByTrashed(&reqService)

	if err != nil {
		s.logger.Error("Failed to fetch trashed users", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationUserDeleteAt(paginationMeta, "success", "Successfully fetched trashed users", users)

	return so, nil
}

func (s *userHandleGrpc) Create(ctx context.Context, request *pb.CreateUserRequest) (*pb.ApiResponseUser, error) {
	req := &requests.CreateUserRequest{
		FirstName:       request.GetFirstname(),
		LastName:        request.GetLastname(),
		Email:           request.GetEmail(),
		Password:        request.GetPassword(),
		ConfirmPassword: request.GetConfirmPassword(),
	}

	s.logger.Debug("Creating user", zap.Any("request", req))

	if err := req.Validate(); err != nil {
		s.logger.Error("Failed to create user", zap.Any("error", err))
		return nil, user_errors.ErrGrpcValidateCreateUser
	}

	user, err := s.userCommandService.CreateUser(req)

	if err != nil {
		s.logger.Error("Failed to create user", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseUser("success", "Successfully created user", user)

	return so, nil
}

func (s *userHandleGrpc) Update(ctx context.Context, request *pb.UpdateUserRequest) (*pb.ApiResponseUser, error) {
	id := int(request.GetId())

	s.logger.Debug("Updating user", zap.Int("id", id))

	if id == 0 {
		s.logger.Error("Failed to update user", zap.Int("id", id))
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	req := &requests.UpdateUserRequest{
		UserID:          &id,
		FirstName:       request.GetFirstname(),
		LastName:        request.GetLastname(),
		Email:           request.GetEmail(),
		Password:        request.GetPassword(),
		ConfirmPassword: request.GetConfirmPassword(),
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Failed to update user", zap.Any("error", err))
		return nil, user_errors.ErrGrpcValidateCreateUser
	}

	user, err := s.userCommandService.UpdateUser(req)

	if err != nil {
		s.logger.Error("Failed to update user", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseUser("success", "Successfully updated user", user)

	return so, nil
}

func (s *userHandleGrpc) TrashedUser(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUserDeleteAt, error) {
	id := int(request.GetId())

	s.logger.Debug("Trashing user", zap.Int("id", id))

	if id == 0 {
		s.logger.Error("Failed to trashed user", zap.Int("id", id))
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	user, err := s.userCommandService.TrashedUser(id)

	if err != nil {
		s.logger.Error("Failed to trashed user", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseUserDeleteAt("success", "Successfully trashed user", user)

	return so, nil
}

func (s *userHandleGrpc) RestoreUser(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUserDeleteAt, error) {
	id := int(request.GetId())

	s.logger.Debug("Restoring user", zap.Int("id", id))

	if id == 0 {
		s.logger.Error("Failed to restore user", zap.Int("id", id))
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	user, err := s.userCommandService.RestoreUser(id)

	if err != nil {
		s.logger.Error("Failed to restore user", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseUserDeleteAt("success", "Successfully restored user", user)

	return so, nil
}

func (s *userHandleGrpc) DeleteUserPermanent(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUserDelete, error) {
	id := int(request.GetId())

	s.logger.Debug("Deleting user permanently", zap.Int("id", id))

	if id == 0 {
		s.logger.Error("Failed to delete user permanently", zap.Int("id", id))
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	_, err := s.userCommandService.DeleteUserPermanent(id)

	if err != nil {
		s.logger.Error("Failed to delete user permanently", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseUserDelete("success", "Successfully deleted user permanently")

	return so, nil
}

func (s *userHandleGrpc) RestoreAllUser(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseUserAll, error) {
	s.logger.Debug("Restoring all user")

	_, err := s.userCommandService.RestoreAllUser()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseUserAll("success", "Successfully restore all user")

	return so, nil
}

func (s *userHandleGrpc) DeleteAllUserPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseUserAll, error) {
	s.logger.Debug("Deleting all user permanently")

	_, err := s.userCommandService.DeleteAllUserPermanent()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseUserAll("success", "Successfully delete user permanen")

	return so, nil
}
