package handler

import (
	"context"

	pb "github.com/MamangRust/monolith-payment-gateway-pb/role"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-role/internal/service"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	role_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/role_errors/grpc"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/role"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type roleCommandHandleGrpc struct {
	pb.UnimplementedRoleCommandServiceServer
	roleCommand service.RoleCommandService
	mapper      protomapper.RoleCommandProtoMapper
	logger      logger.LoggerInterface
}

func NewRoleCommandHandleGrpc(roleCommand service.RoleCommandService, logger logger.LoggerInterface, protomapper protomapper.RoleCommandProtoMapper) RoleCommandHandlerGrpc {
	return &roleCommandHandleGrpc{
		roleCommand: roleCommand,
		mapper:      protomapper,
		logger:      logger,
	}
}

// CreateRole creates a new role with the provided details.
//
// This method validates the request payload and calls the command service to
// create the role. If the request is invalid or the operation fails,
// it returns an error.
//
// Parameters:
//   - ctx: the context for the request, used for cancellation and timeouts.
//   - reqPb: the request payload containing the role details.
//
// Returns:
//   - A pointer to ApiResponseRole containing the created role.
//   - An error if the operation fails, otherwise nil.
func (s *roleCommandHandleGrpc) CreateRole(ctx context.Context, reqPb *pb.CreateRoleRequest) (*pb.ApiResponseRole, error) {
	req := &requests.CreateRoleRequest{
		Name: reqPb.Name,
	}

	s.logger.Info("Creating role", zap.Any("request", req))

	if err := req.Validate(); err != nil {
		s.logger.Error("CreateRole failed", zap.Any("error", err))
		return nil, role_errors.ErrGrpcValidateCreateRole
	}

	role, err := s.roleCommand.CreateRole(ctx, req)

	if err != nil {
		s.logger.Error("CreateRole failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseRole("success", "Successfully created role", role)

	s.logger.Info("Successfully created role", zap.Bool("success", true))

	return so, nil
}

// UpdateRole updates an existing role with the provided details.
//
// This method validates the request payload and calls the command service to
// update the role. If the request is invalid or the operation fails,
// it returns an error.
//
// Parameters:
//   - ctx: the context for the request, used for cancellation and timeouts.
//   - reqPb: the request payload containing the role details.
//
// Returns:
//   - A pointer to ApiResponseRole containing the updated role.
//   - An error if the operation fails, otherwise nil.
func (s *roleCommandHandleGrpc) UpdateRole(ctx context.Context, reqPb *pb.UpdateRoleRequest) (*pb.ApiResponseRole, error) {
	roleID := int(reqPb.GetId())

	s.logger.Info("Updating role", zap.Int("id", roleID))

	if roleID == 0 {
		s.logger.Error("UpdateRole failed", zap.Any("error", role_errors.ErrGrpcRoleInvalidId))
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	name := reqPb.GetName()

	req := &requests.UpdateRoleRequest{
		ID:   &roleID,
		Name: name,
	}

	if err := req.Validate(); err != nil {
		s.logger.Info("UpdateRole failed", zap.Any("error", err))
		return nil, role_errors.ErrGrpcValidateUpdateRole
	}

	role, err := s.roleCommand.UpdateRole(ctx, req)

	if err != nil {
		s.logger.Info("UpdateRole failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseRole("success", "Successfully updated role", role)

	s.logger.Info("Successfully updated role", zap.Bool("success", true))

	return so, nil
}

// TrashedRole trashes a role with the given ID.
//
// Parameters:
//   - ctx: the context for the request, used for cancellation and timeouts.
//   - req: the request payload containing the role ID.
//
// Returns:
//   - A pointer to ApiResponseRole containing the trashed role.
//   - An error if the operation fails, otherwise nil.
func (s *roleCommandHandleGrpc) TrashedRole(ctx context.Context, req *pb.FindByIdRoleRequest) (*pb.ApiResponseRoleDeleteAt, error) {
	roleID := int(req.GetRoleId())

	s.logger.Info("Trashing role", zap.Int("id", roleID))

	if roleID == 0 {
		s.logger.Error("TrashedRole failed", zap.Any("error", role_errors.ErrGrpcRoleInvalidId))
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	role, err := s.roleCommand.TrashedRole(ctx, roleID)

	if err != nil {
		s.logger.Error("TrashedRole failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseRoleDeleteAt("success", "Successfully trashed role", role)

	s.logger.Info("Successfully trashed role", zap.Bool("success", true))

	return so, nil
}

// RestoreRole restores a trashed role with the given ID.
//
// Parameters:
//   - ctx: the context for the request, used for cancellation and timeouts.
//   - req: the request payload containing the role ID.
//
// Returns:
//   - A pointer to ApiResponseRole containing the restored role.
//   - An error if the operation fails, otherwise nil.
func (s *roleCommandHandleGrpc) RestoreRole(ctx context.Context, req *pb.FindByIdRoleRequest) (*pb.ApiResponseRole, error) {
	roleID := int(req.GetRoleId())

	s.logger.Info("Restoring role", zap.Int("id", roleID))

	if roleID == 0 {
		s.logger.Error("RestoreRole failed", zap.Any("error", role_errors.ErrGrpcRoleInvalidId))
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	role, err := s.roleCommand.RestoreRole(ctx, roleID)

	if err != nil {
		s.logger.Error("RestoreRole failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseRole("success", "Successfully restored role", role)

	s.logger.Info("Successfully restored role", zap.Bool("success", true))

	return so, nil
}

// DeleteRolePermanent deletes a role permanently by its ID.
//
// Parameters:
//   - ctx: the context for the request, used for cancellation and timeouts.
//   - req: the request payload containing the role ID.
//
// Returns:
//   - A pointer to ApiResponseRoleDelete containing the deleted role.
//   - An error if the operation fails, otherwise nil.
func (s *roleCommandHandleGrpc) DeleteRolePermanent(ctx context.Context, req *pb.FindByIdRoleRequest) (*pb.ApiResponseRoleDelete, error) {
	id := int(req.GetRoleId())

	s.logger.Info("Deleting role permanently", zap.Int("id", id))

	if id == 0 {
		s.logger.Error("DeleteRolePermanent failed", zap.Any("error", role_errors.ErrGrpcRoleInvalidId))
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	_, err := s.roleCommand.DeleteRolePermanent(ctx, id)

	if err != nil {
		s.logger.Error("DeleteRolePermanent failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseRoleDelete("success", "Successfully deleted role permanently")

	s.logger.Info("Successfully deleted role permanently", zap.Bool("success", true))

	return so, nil
}

// RestoreAllRole restores all trashed roles.
//
// Parameters:
//   - ctx: the context for the request, used for cancellation and timeouts.
//   - req: an empty request object.
//
// Returns:
//   - An ApiResponseRoleAll containing the restored roles.
//   - An error if the operation fails.
func (s *roleCommandHandleGrpc) RestoreAllRole(ctx context.Context, req *emptypb.Empty) (*pb.ApiResponseRoleAll, error) {
	s.logger.Info("Restoring all roles")

	_, err := s.roleCommand.RestoreAllRole(ctx)

	if err != nil {
		s.logger.Error("RestoreAllRole failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseRoleAll("success", "Successfully restored all roles")

	s.logger.Info("Successfully restored all roles", zap.Bool("success", true))

	return so, nil
}

// DeleteAllRolePermanent permanently deletes all roles.
//
// Parameters:
//   - ctx: the context for the request, used for cancellation and timeouts.
//   - req: an empty request object.
//
// Returns:
//   - An ApiResponseRoleAll containing the result of the deletion.
//   - An error if the operation fails.
func (s *roleCommandHandleGrpc) DeleteAllRolePermanent(ctx context.Context, req *emptypb.Empty) (*pb.ApiResponseRoleAll, error) {
	s.logger.Info("Deleting all roles permanently")

	_, err := s.roleCommand.DeleteAllRolePermanent(ctx)

	if err != nil {
		s.logger.Error("DeleteAllRolePermanent failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapper.ToProtoResponseRoleAll("success", "Successfully deleted all roles")

	s.logger.Info("Successfully deleted all roles", zap.Bool("success", true))

	return so, nil
}
