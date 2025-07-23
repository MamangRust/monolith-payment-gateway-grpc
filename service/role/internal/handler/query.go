package handler

import (
	"context"
	"math"

	pbhelpers "github.com/MamangRust/monolith-payment-gateway-pb"
	pb "github.com/MamangRust/monolith-payment-gateway-pb/role"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-role/internal/service"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	role_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/role_errors/grpc"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/role"
	"go.uber.org/zap"
)

type roleQueryHandleGrpc struct {
	pb.UnimplementedRoleQueryServiceServer
	roleQuery service.RoleQueryService
	mapper    protomapper.RoleQueryProtoMapper
	logger    logger.LoggerInterface
}

func NewRoleQueryHandleGrpc(roleQuery service.RoleQueryService, logger logger.LoggerInterface, protomapper protomapper.RoleQueryProtoMapper) RoleQueryHandlerGrpc {
	return &roleQueryHandleGrpc{
		roleQuery: roleQuery,
		mapper:    protomapper,
		logger:    logger,
	}
}

// FindAllRole retrieves all role records with pagination and optional search.
//
// This method fetches role records from the database, applies pagination,
// and returns the results along with pagination metadata.
//
// Parameters:
//   - ctx: the context for the request, used for cancellation and timeouts.
//   - req: the request payload containing pagination and search parameters.
//
// Returns:
//   - A pointer to ApiResponsePaginationRole containing the paginated role records.
//   - An error if the operation fails, otherwise nil.
func (s *roleQueryHandleGrpc) FindAllRole(ctx context.Context, req *pb.FindAllRoleRequest) (*pb.ApiResponsePaginationRole, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	s.logger.Info("Fetching role records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := &requests.FindAllRoles{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	role, totalRecords, err := s.roleQuery.FindAll(ctx, reqService)

	if err != nil {
		s.logger.Error("FindAll failed", zap.Any("error", err))

		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbhelpers.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapper.ToProtoResponsePaginationRole(paginationMeta, "success", "Successfully fetched role records", role)

	s.logger.Info("Successfully fetched role records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, nil
}

// FindByIdRole retrieves a role record by its ID.
//
// Parameters:
//   - ctx: the context for the request, used for cancellation and timeouts.
//   - req: the request payload containing the role ID.
//
// Returns:
//   - A pointer to ApiResponseRole containing the role record.
//   - An error if the operation fails, otherwise nil.
func (s *roleQueryHandleGrpc) FindByIdRole(ctx context.Context, req *pb.FindByIdRoleRequest) (*pb.ApiResponseRole, error) {
	roleID := int(req.GetRoleId())

	s.logger.Info("Fetching role record", zap.Int("id", roleID))

	if roleID == 0 {
		s.logger.Debug("FindById failed", zap.Any("error", role_errors.ErrGrpcRoleInvalidId))
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	role, err := s.roleQuery.FindById(ctx, roleID)

	if err != nil {
		s.logger.Error("FindById failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	roleResponse := s.mapper.ToProtoResponseRole("success", "Successfully fetched role", role)

	s.logger.Info("Successfully fetched role", zap.Bool("success", true))

	return roleResponse, nil
}

// FindByUserId retrieves a role record associated with a user ID.
//
// Parameters:
//   - ctx: the context for the request, used for cancellation and timeouts.
//   - req: the request payload containing the user ID.
//
// Returns:
//   - A pointer to ApiResponseRoles containing the role records.
//   - An error if the operation fails, otherwise nil.
func (s *roleQueryHandleGrpc) FindByUserId(ctx context.Context, req *pb.FindByIdUserRoleRequest) (*pb.ApiResponsesRole, error) {
	userID := int(req.GetUserId())

	s.logger.Info("Fetching role record", zap.Int("id", userID))

	if userID == 0 {
		s.logger.Error("FindByUserId failed", zap.Any("error", role_errors.ErrGrpcRoleInvalidId))
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	role, err := s.roleQuery.FindByUserId(ctx, userID)

	if err != nil {
		s.logger.Error("FindByUserId failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	roleResponse := s.mapper.ToProtoResponsesRole("success", "Successfully fetched role by user ID", role)

	s.logger.Info("Successfully fetched role by user ID", zap.Bool("success", true))

	return roleResponse, nil
}

// FindByActive retrieves active role records with pagination and optional search.
//
// This method fetches role records from the database, applies pagination,
// and returns the results along with pagination metadata.
//
// Parameters:
//   - ctx: the context for the request, used for cancellation and timeouts.
//   - req: the request payload containing pagination and search parameters.
//
// Returns:
//   - A pointer to ApiResponsePaginationRoleDeleteAt containing the paginated role records.
//   - An error if the operation fails, otherwise nil.
func (s *roleQueryHandleGrpc) FindByActive(ctx context.Context, req *pb.FindAllRoleRequest) (*pb.ApiResponsePaginationRoleDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	s.logger.Info("Fetching active role records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := &requests.FindAllRoles{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	roles, totalRecords, err := s.roleQuery.FindByActiveRole(ctx, reqService)

	if err != nil {
		s.logger.Error("FindByActive failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbhelpers.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapper.ToProtoResponsePaginationRoleDeleteAt(paginationMeta, "success", "Successfully fetched active roles", roles)

	s.logger.Info("Successfully fetched active roles", zap.Bool("success", true))

	return so, nil
}

// FindByTrashed retrieves trashed role records with pagination and optional search.
//
// This method fetches trashed role records from the database, applies pagination,
// and returns the results along with pagination metadata.
//
// Parameters:
//   - ctx: the context for the request, used for cancellation and timeouts.
//   - req: the request payload containing pagination and search parameters.
//
// Returns:
//   - A pointer to ApiResponsePaginationRoleDeleteAt containing the paginated trashed role records.
//   - An error if the operation fails, otherwise nil.
func (s *roleQueryHandleGrpc) FindByTrashed(ctx context.Context, req *pb.FindAllRoleRequest) (*pb.ApiResponsePaginationRoleDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	s.logger.Info("Fetching trashed role records", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := &requests.FindAllRoles{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	roles, totalRecords, err := s.roleQuery.FindByTrashedRole(ctx, reqService)

	if err != nil {
		s.logger.Error("FindByTrashed failed", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbhelpers.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapper.ToProtoResponsePaginationRoleDeleteAt(paginationMeta, "success", "Successfully fetched trashed roles", roles)

	s.logger.Info("Successfully fetched trashed roles", zap.Bool("success", true))

	return so, nil
}
