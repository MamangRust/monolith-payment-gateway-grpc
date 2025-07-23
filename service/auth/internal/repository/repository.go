package repository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	refreshtokenrecord "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/refreshtoken"
	resettokenrecord "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/resettoken"
	rolemapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/role"
	userrecord "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/user"
	userrolerecord "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/userrole"
)

type Repositories struct {
	User         UserRepository
	RefreshToken RefreshTokenRepository
	UserRole     UserRoleRepository
	Role         RoleRepository
	ResetToken   ResetTokenRepository
}

type Deps struct {
	DB *db.Queries
}

func NewRepositories(deps *Deps) *Repositories {
	mapperuser := userrecord.NewUserQueryRecordMapper()
	mapperuserrole := userrolerecord.NewUserRoleRecordMapper()
	mapperrefreshtoken := refreshtokenrecord.NewRefreshTokenRecordMapper()
	mapperrole := rolemapper.NewRoleQueryRecordMapping()
	mapperresettoken := resettokenrecord.NewResetTokenRecordMapper()

	return &Repositories{
		User:         NewUserRepository(deps.DB, mapperuser),
		UserRole:     NewUserRoleRepository(deps.DB, mapperuserrole),
		RefreshToken: NewRefreshTokenRepository(deps.DB, mapperrefreshtoken),
		Role:         NewRoleRepository(deps.DB, mapperrole),
		ResetToken:   NewResetTokenRepository(deps.DB, mapperresettoken),
	}
}
