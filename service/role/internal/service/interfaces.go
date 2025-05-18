package service

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type RoleQueryService interface {
	FindAll(req *requests.FindAllRoles) ([]*response.RoleResponse, *int, *response.ErrorResponse)
	FindByActiveRole(req *requests.FindAllRoles) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashedRole(req *requests.FindAllRoles) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(role_id int) (*response.RoleResponse, *response.ErrorResponse)
	FindByUserId(id int) ([]*response.RoleResponse, *response.ErrorResponse)
}

type RoleCommandService interface {
	CreateRole(request *requests.CreateRoleRequest) (*response.RoleResponse, *response.ErrorResponse)
	UpdateRole(request *requests.UpdateRoleRequest) (*response.RoleResponse, *response.ErrorResponse)
	TrashedRole(role_id int) (*response.RoleResponse, *response.ErrorResponse)
	RestoreRole(role_id int) (*response.RoleResponse, *response.ErrorResponse)
	DeleteRolePermanent(role_id int) (bool, *response.ErrorResponse)

	RestoreAllRole() (bool, *response.ErrorResponse)
	DeleteAllRolePermanent() (bool, *response.ErrorResponse)
}
