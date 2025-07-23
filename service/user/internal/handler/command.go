package handler

import (
	"context"

	pb "github.com/MamangRust/monolith-payment-gateway-pb/user"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	user_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors/grpc"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/user"
	"github.com/MamangRust/monolith-payment-gateway-user/internal/service"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userCommandHandleGrpc struct {
	pb.UnimplementedUserCommandServiceServer

	userCommandService service.UserCommandService

	logger logger.LoggerInterface

	mapper protomapper.UserCommandProtoMapper
}

func NewUserCommandHandleGrpc(query service.UserCommandService, logger logger.LoggerInterface, mapper protomapper.UserCommandProtoMapper) UserCommandHandleGrpc {
	return &userCommandHandleGrpc{
		userCommandService: query,
		logger:             logger,
		mapper:             mapper,
	}
}

// Create is a gRPC handler that creates a new user according to the given request.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - request: A pointer to a CreateUserRequest containing the user details.
//
// Returns:
//   - A pointer to ApiResponseUser containing the created user record.
//   - An error if the creation operation fails.
func (s *userCommandHandleGrpc) Create(ctx context.Context, request *pb.CreateUserRequest) (*pb.ApiResponseUser, error) {
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

	user, err := s.userCommandService.CreateUser(ctx, req)

	if err != nil {
		s.logger.Error("Failed to create user", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseUser("success", "Successfully created user", user)

	return so, nil
}

// Update is a gRPC handler that updates an existing user according to the given request.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - request: A pointer to an UpdateUserRequest containing the user ID and updated details.
//
// Returns:
//   - A pointer to ApiResponseUser containing the updated user record.
//   - An error if the update operation fails.
func (s *userCommandHandleGrpc) Update(ctx context.Context, request *pb.UpdateUserRequest) (*pb.ApiResponseUser, error) {
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

	user, err := s.userCommandService.UpdateUser(ctx, req)

	if err != nil {
		s.logger.Error("Failed to update user", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseUser("success", "Successfully updated user", user)

	return so, nil
}

// TrashedUser trashes a user account identified by the provided ID.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - request: A pointer to a FindByIdUserRequest containing the user ID to trash.
//
// Returns:
//   - A pointer to ApiResponseUserDeleteAt containing the trashed user data on success.
//   - An error if the operation fails, or if the provided ID is invalid.
func (s *userCommandHandleGrpc) TrashedUser(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUserDeleteAt, error) {
	id := int(request.GetId())

	s.logger.Debug("Trashing user", zap.Int("id", id))

	if id == 0 {
		s.logger.Error("Failed to trashed user", zap.Int("id", id))
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	user, err := s.userCommandService.TrashedUser(ctx, id)

	if err != nil {
		s.logger.Error("Failed to trashed user", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseUserDeleteAt("success", "Successfully trashed user", user)

	return so, nil
}

// RestoreUser restores a user account from the trash by its ID.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - request: A pointer to a FindByIdUserRequest containing the user ID to restore.
//
// Returns:
//   - A pointer to ApiResponseUserDeleteAt containing the restored user data on success.
//   - An error if the operation fails, or if the provided ID is invalid.
func (s *userCommandHandleGrpc) RestoreUser(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUser, error) {
	id := int(request.GetId())

	s.logger.Debug("Restoring user", zap.Int("id", id))

	if id == 0 {
		s.logger.Error("Failed to restore user", zap.Int("id", id))
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	user, err := s.userCommandService.RestoreUser(ctx, id)

	if err != nil {
		s.logger.Error("Failed to restore user", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseUser("success", "Successfully restored user", user)

	return so, nil
}

// DeleteUserPermanent deletes a user permanently by its ID.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - req: A pointer to a FindByIdUserRequest containing the user ID to delete.
//
// Returns:
//   - A pointer to ApiResponseUserDelete containing the deleted user data on success.
//   - An error if the operation fails or if the provided ID is invalid.
func (s *userCommandHandleGrpc) DeleteUserPermanent(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUserDelete, error) {
	id := int(request.GetId())

	s.logger.Debug("Deleting user permanently", zap.Int("id", id))

	if id == 0 {
		s.logger.Error("Failed to delete user permanently", zap.Int("id", id))
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	_, err := s.userCommandService.DeleteUserPermanent(ctx, id)

	if err != nil {
		s.logger.Error("Failed to delete user permanently", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseUserDelete("success", "Successfully deleted user permanently")

	return so, nil
}

// RestoreAllUser restores all trashed user records.
//
// Parameters:
//   - ctx: the context for the request, used for cancellation and timeouts.
//   - _: an empty request object.
//
// Returns:
//   - An ApiResponseUserAll containing the restored users.
//   - An error if the operation fails.
func (s *userCommandHandleGrpc) RestoreAllUser(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseUserAll, error) {
	s.logger.Info("Restoring all users")

	_, err := s.userCommandService.RestoreAllUser(ctx)

	if err != nil {
		s.logger.Error("Failed to restore all users", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseUserAll("success", "Successfully restored all users")

	s.logger.Info("Successfully restored all users", zap.Bool("success", true))

	return so, nil
}

// DeleteAllUserPermanent permanently deletes all user records that were previously trashed.
//
// Parameters:
//   - ctx: the context for the request, used for cancellation and timeouts.
//   - _: an empty request object.
//
// Returns:
//   - An ApiResponseUserAll containing the deleted users.
//   - An error if the operation fails.
func (s *userCommandHandleGrpc) DeleteAllUserPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseUserAll, error) {
	s.logger.Debug("Deleting all user permanently")

	_, err := s.userCommandService.DeleteAllUserPermanent(ctx)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseUserAll("success", "Successfully deleted all user permanently")

	return so, nil
}
