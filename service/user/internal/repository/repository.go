package repository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	rolemapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/role"
	usermapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/user"
)

type Repositories interface {
	UserQuery() UserQueryRepository
	UserCommand() UserCommandRepository
	Role() RoleRepository
}

type repositories struct {
	userQuery   UserQueryRepository
	userCommand UserCommandRepository
	role        RoleRepository
}

func NewRepositories(db *db.Queries) Repositories {
	usermapper := usermapper.NewUserRecordMapper()
	rolemapper := rolemapper.NewRoleQueryRecordMapping()

	return &repositories{
		userCommand: NewUserCommandRepository(db, usermapper.CommandMapper()),
		userQuery:   NewUserQueryRepository(db, usermapper.QueryMapper()),
		role:        NewRoleRepository(db, rolemapper),
	}
}

func (r *repositories) UserQuery() UserQueryRepository {
	return r.userQuery
}
func (r *repositories) UserCommand() UserCommandRepository {
	return r.userCommand
}
func (r *repositories) Role() RoleRepository {
	return r.role
}
