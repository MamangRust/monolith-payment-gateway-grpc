package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-role/internal/service"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/role_errors"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto"
	"github.com/MamangRust/monolith-payment-gateway-shared/pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type roleHandleGrpc struct {
	pb.UnimplementedRoleServiceServer
	roleQuery   service.RoleQueryService
	roleCommand service.RoleCommandService
	mapping     protomapper.RoleProtoMapper
	logger      logger.LoggerInterface
}

func NewRoleHandleGrpc(service *service.Service, logger logger.LoggerInterface) *roleHandleGrpc {
	return &roleHandleGrpc{
		roleQuery:   service.RoleQuery,
		roleCommand: service.RoleCommand,
		mapping:     protomapper.NewRoleProtoMapper(),
		logger:      logger,
	}
}

func (s *roleHandleGrpc) FindAllRole(ctx context.Context, req *pb.FindAllRoleRequest) (*pb.ApiResponsePaginationRole, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	s.logger.Debug("Fetching role records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllRoles{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	role, totalRecords, err := s.roleQuery.FindAll(&reqService)

	if err != nil {
		s.logger.Debug("FindAll failed", zap.Any("error", err))

		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationRole(paginationMeta, "success", "Successfully fetched role records", role)

	return so, nil
}

func (s *roleHandleGrpc) FindByIdRole(ctx context.Context, req *pb.FindByIdRoleRequest) (*pb.ApiResponseRole, error) {
	roleID := int(req.GetRoleId())

	s.logger.Debug("Fetching role record", zap.Int("id", roleID))

	if roleID == 0 {
		s.logger.Debug("FindById failed", zap.Any("error", role_errors.ErrGrpcRoleInvalidId))
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	role, err := s.roleQuery.FindById(roleID)

	if err != nil {
		s.logger.Debug("FindById failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	roleResponse := s.mapping.ToProtoResponseRole("success", "Successfully fetched role", role)

	return roleResponse, nil
}

func (s *roleHandleGrpc) FindByUserId(ctx context.Context, req *pb.FindByIdUserRoleRequest) (*pb.ApiResponsesRole, error) {
	userID := int(req.GetUserId())

	s.logger.Debug("Fetching role record", zap.Int("id", userID))

	if userID == 0 {
		s.logger.Debug("FindByUserId failed", zap.Any("error", role_errors.ErrGrpcRoleInvalidId))
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	role, err := s.roleQuery.FindByUserId(userID)

	if err != nil {
		s.logger.Debug("FindByUserId failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	roleResponse := s.mapping.ToProtoResponsesRole("success", "Successfully fetched role by user ID", role)

	return roleResponse, nil
}

func (s *roleHandleGrpc) FindByActive(ctx context.Context, req *pb.FindAllRoleRequest) (*pb.ApiResponsePaginationRoleDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	s.logger.Debug("Fetching active role records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllRoles{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	roles, totalRecords, err := s.roleQuery.FindByActiveRole(&reqService)

	if err != nil {
		s.logger.Debug("FindByActive failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationRoleDeleteAt(paginationMeta, "success", "Successfully fetched active roles", roles)

	return so, nil
}

func (s *roleHandleGrpc) FindByTrashed(ctx context.Context, req *pb.FindAllRoleRequest) (*pb.ApiResponsePaginationRoleDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	s.logger.Debug("Fetching trashed role records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllRoles{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	roles, totalRecords, err := s.roleQuery.FindByTrashedRole(&reqService)

	if err != nil {
		s.logger.Debug("FindByTrashed failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapping.ToProtoResponsePaginationRoleDeleteAt(paginationMeta, "success", "Successfully fetched trashed roles", roles)

	return so, nil
}

func (s *roleHandleGrpc) CreateRole(ctx context.Context, reqPb *pb.CreateRoleRequest) (*pb.ApiResponseRole, error) {
	req := &requests.CreateRoleRequest{
		Name: reqPb.Name,
	}

	s.logger.Debug("Creating role", zap.Any("request", req))

	if err := req.Validate(); err != nil {
		s.logger.Debug("CreateRole failed", zap.Any("error", err))
		return nil, role_errors.ErrGrpcFailedCreateRole
	}

	role, err := s.roleCommand.CreateRole(req)

	if err != nil {
		s.logger.Debug("CreateRole failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseRole("success", "Successfully created role", role)

	return so, nil
}

func (s *roleHandleGrpc) UpdateRole(ctx context.Context, reqPb *pb.UpdateRoleRequest) (*pb.ApiResponseRole, error) {
	roleID := int(reqPb.GetId())

	s.logger.Debug("Updating role", zap.Int("id", roleID))

	if roleID == 0 {
		s.logger.Debug("UpdateRole failed", zap.Any("error", role_errors.ErrGrpcRoleInvalidId))
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	name := reqPb.GetName()

	req := &requests.UpdateRoleRequest{
		ID:   &roleID,
		Name: name,
	}

	if err := req.Validate(); err != nil {
		s.logger.Debug("UpdateRole failed", zap.Any("error", err))
		return nil, role_errors.ErrGrpcValidateUpdateRole
	}

	role, err := s.roleCommand.UpdateRole(req)

	if err != nil {
		s.logger.Debug("UpdateRole failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseRole("success", "Successfully updated role", role)

	return so, nil
}

func (s *roleHandleGrpc) TrashedRole(ctx context.Context, req *pb.FindByIdRoleRequest) (*pb.ApiResponseRole, error) {
	roleID := int(req.GetRoleId())

	s.logger.Debug("Trashing role", zap.Int("id", roleID))

	if roleID == 0 {
		s.logger.Debug("TrashedRole failed", zap.Any("error", role_errors.ErrGrpcRoleInvalidId))
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	role, err := s.roleCommand.TrashedRole(roleID)

	if err != nil {
		s.logger.Debug("TrashedRole failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseRole("success", "Successfully trashed role", role)

	return so, nil
}

func (s *roleHandleGrpc) RestoreRole(ctx context.Context, req *pb.FindByIdRoleRequest) (*pb.ApiResponseRole, error) {
	roleID := int(req.GetRoleId())

	s.logger.Debug("Restoring role", zap.Int("id", roleID))

	if roleID == 0 {
		s.logger.Debug("RestoreRole failed", zap.Any("error", role_errors.ErrGrpcRoleInvalidId))
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	role, err := s.roleCommand.RestoreRole(roleID)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseRole("success", "Successfully restored role", role)

	return so, nil
}

func (s *roleHandleGrpc) DeleteRolePermanent(ctx context.Context, req *pb.FindByIdRoleRequest) (*pb.ApiResponseRoleDelete, error) {
	id := int(req.GetRoleId())

	s.logger.Debug("Deleting role permanently", zap.Int("id", id))

	if id == 0 {
		s.logger.Debug("DeleteRolePermanent failed", zap.Any("error", role_errors.ErrGrpcRoleInvalidId))
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	_, err := s.roleCommand.DeleteRolePermanent(id)

	if err != nil {
		s.logger.Debug("DeleteRolePermanent failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseRoleDelete("success", "Successfully deleted role permanently")

	return so, nil
}

func (s *roleHandleGrpc) RestoreAllRole(ctx context.Context, req *emptypb.Empty) (*pb.ApiResponseRoleAll, error) {
	s.logger.Debug("Restoring all roles")

	_, err := s.roleCommand.RestoreAllRole()

	if err != nil {
		s.logger.Debug("RestoreAllRole failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseRoleAll("success", "Successfully restored all roles")

	return so, nil
}

func (s *roleHandleGrpc) DeleteAllRolePermanent(ctx context.Context, req *emptypb.Empty) (*pb.ApiResponseRoleAll, error) {
	s.logger.Debug("Deleting all roles permanently")

	_, err := s.roleCommand.DeleteAllRolePermanent()

	if err != nil {
		s.logger.Debug("DeleteAllRolePermanent failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseRoleAll("success", "Successfully deleted all roles")

	return so, nil
}
