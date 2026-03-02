package repository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
)

type Repositories struct {
	User         UserRepository
	RefreshToken RefreshTokenRepository
	UserRole     UserRoleRepository
	Role         RoleRepository
	ResetToken   ResetTokenRepository
}

func NewRepositories(db *db.Queries) *Repositories {
	return &Repositories{
		User:         NewUserRepository(db),
		UserRole:     NewUserRoleRepository(db),
		RefreshToken: NewRefreshTokenRepository(db),
		Role:         NewRoleRepository(db),
		ResetToken:   NewResetTokenRepository(db),
	}
}
