package mencache

import (
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
)

type RoleCommandCache interface {
	DeleteCachedRole(id int)
}

type RoleQueryCache interface {
	SetCachedRoles(req *requests.FindAllRoles, data []*response.RoleResponse, total *int)
	SetCachedRoleById(id int, data *response.RoleResponse)
	SetCachedRoleByUserId(userId int, data []*response.RoleResponse)
	SetCachedRoleActive(req *requests.FindAllRoles, data []*response.RoleResponseDeleteAt, total *int)
	SetCachedRoleTrashed(req *requests.FindAllRoles, data []*response.RoleResponseDeleteAt, total *int)
	GetCachedRoles(req *requests.FindAllRoles) ([]*response.RoleResponse, *int, bool)
	GetCachedRoleByUserId(userId int) ([]*response.RoleResponse, bool)
	GetCachedRoleById(id int) (*response.RoleResponse, bool)
	GetCachedRoleActive(req *requests.FindAllRoles) ([]*response.RoleResponseDeleteAt, *int, bool)
	GetCachedRoleTrashed(req *requests.FindAllRoles) ([]*response.RoleResponseDeleteAt, *int, bool)
}
