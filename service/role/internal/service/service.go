package service

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-role/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-role/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-role/internal/repository"
	roleservicemapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/role"
	"github.com/redis/go-redis/v9"
)

// Service groups the role-related domain services, including query and command logic.
type Service struct {
	// RoleQuery handles read-only operations for roles.
	RoleQuery RoleQueryService

	// RoleCommand handles write operations (create, update, delete) for roles.
	RoleCommand RoleCommandService
}

// Deps holds the shared dependencies required to construct the role services.
type Deps struct {

	// ErrorHandler provides centralized error handling for role services.
	ErrorHandler *errorhandler.ErrorHandler

	// Mencache provides in-memory caching for role services.
	Mencache *mencache.Mencache

	// Redis is the Redis client used for distributed caching or locking.
	Redis *redis.Client

	// Repositories provides access to all repository interfaces.
	Repositories repository.Repositories

	// Logger provides structured logging capabilities.
	Logger logger.LoggerInterface
}

// NewService constructs and returns a new Service instance using the provided dependencies.
// It initializes both RoleQuery and RoleCommand services with their required dependencies.
func NewService(deps *Deps) *Service {
	roleMapper := roleservicemapper.NewRoleResponseMapper()

	return &Service{
		RoleQuery:   newRoleQueryService(deps, roleMapper.QueryMapper()),
		RoleCommand: newRoleCommandService(deps, roleMapper.CommandMapper()),
	}
}

// newRoleCommandService initializes and returns a new instance of RoleCommandService.
// It uses the provided dependencies and a role response mapper to set up
// necessary components for executing role command operations.
//
// Parameters:
// - deps: A pointer to Deps containing the shared dependencies used by the role services.
// - mapper: A pointer to RoleResponseMapper that maps domain models to API-compatible response formats.
//
// Returns:
// - A pointer to RoleCommandService, which is responsible for handling role-related command operations.
func newRoleCommandService(deps *Deps, mapper roleservicemapper.RoleCommandResponseMapper) RoleCommandService {
	return NewRoleCommandService(&roleCommandDeps{
		ErrorHandler: deps.ErrorHandler.RoleCommandError,
		Cache:        deps.Mencache.RoleCommandCache,
		Repository:   deps.Repositories,
		Logger:       deps.Logger,
		Mapper:       mapper,
	})
}

// newRoleQueryService initializes and returns a new instance of RoleQueryService.
// It uses the provided dependencies and a role response mapper to set up
// necessary components for executing role query operations.
//
// Parameters:
//   - deps: A pointer to Deps containing the required context, error handler, cache,
//     repository, and logger dependencies.
//   - mapper: A RoleResponseMapper used for mapping role domain models to API-compatible responses.
//
// Returns:
// - A RoleQueryService ready to handle role query operations.
func newRoleQueryService(deps *Deps, mapper roleservicemapper.RoleQueryResponseMapper) RoleQueryService {
	return NewRoleQueryService(&roleQueryDeps{
		ErrorHandler: deps.ErrorHandler.RoleQueryError,
		Cache:        deps.Mencache.RoleQueryCache,
		Repository:   deps.Repositories,
		Logger:       deps.Logger,
		Mapper:       mapper,
	})
}
