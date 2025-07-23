package repository

import (
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	rolerecordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/role"
)

// Repositories is a struct containing role command and query repositories.
type Repositories interface {
	RoleQueryRepository
	RoleCommandRepository
}

type repositories struct {
	RoleQueryRepository
	RoleCommandRepository
}

// NewRepositories creates a new instance of Repositories with the provided database
// queries, context, and role record mapper. This repository is responsible for
// executing command and query operations related to role records in the database.
//
// Parameters:
//   - deps: A pointer to Deps containing the required dependencies.
//
// Returns:
//   - A pointer to the newly created Repositories instance.
func NewRepositories(db *db.Queries) Repositories {
	roleMapper := rolerecordmapper.NewRoleRecordMapper()

	return &repositories{
		RoleQueryRepository:   NewRoleQueryRepository(db, roleMapper.QueryMapper()),
		RoleCommandRepository: NewRoleCommandRepository(db, roleMapper.CommandMapper()),
	}
}
