package repository

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	user_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/user"
)

// userQueryRepository is a struct that implements the UserQueryRepository interface
type userQueryRepository struct {
	db     *db.Queries
	mapper recordmapper.UserQueryRecordMapper
}

// NewUserQueryRepository creates a new instance of userQueryRepository.
//
// Parameters:
//   - db: A pointer to the db.Queries object for executing database queries.
//   - mapper: A UserRecordMapping that provides methods to map database rows to User domain models.
//
// Returns:
//   - A pointer to the newly created userQueryRepository instance.
func NewUserQueryRepository(db *db.Queries, mapper recordmapper.UserQueryRecordMapper) UserQueryRepository {
	return &userQueryRepository{
		db:     db,
		mapper: mapper,
	}
}

// FindAllUsers retrieves all users with optional filters and pagination.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing filter and pagination information.
//
// Returns:
//   - []*record.UserRecord: List of user records.
//   - *int: Total count of records.
//   - error: Error if retrieval fails.
func (r *userQueryRepository) FindAllUsers(ctx context.Context, req *requests.FindAllUsers) ([]*record.UserRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetUsersWithPaginationParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetUsersWithPagination(ctx, reqDb)

	if err != nil {
		return nil, nil, user_errors.ErrFindAllUsers
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	so := r.mapper.ToUsersRecordPagination(res)

	return so, &totalCount, nil
}

// FindById retrieves a user by their ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - user_id: The unique identifier of the user.
//
// Returns:
//   - *record.UserRecord: The user record if found.
//   - error: Error if retrieval fails.
func (r *userQueryRepository) FindById(ctx context.Context, user_id int) (*record.UserRecord, error) {
	res, err := r.db.GetUserByID(ctx, int32(user_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound
		}

		return nil, user_errors.ErrUserNotFound
	}

	so := r.mapper.ToUserRecord(res)

	return so, nil
}

// FindByActive retrieves all active (non-deleted) users.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing filter and pagination information.
//
// Returns:
//   - []*record.UserRecord: List of active user records.
//   - *int: Total count of records.
//   - error: Error if retrieval fails.
func (r *userQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*record.UserRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveUsersWithPaginationParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveUsersWithPagination(ctx, reqDb)

	if err != nil {
		return nil, nil, user_errors.ErrFindActiveUsers
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	so := r.mapper.ToUsersRecordActivePagination(res)

	return so, &totalCount, nil
}

// FindByTrashed retrieves all trashed (soft-deleted) users.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request containing filter and pagination information.
//
// Returns:
//   - []*record.UserRecord: List of trashed user records.
//   - *int: Total count of records.
//   - error: Error if retrieval fails.
func (r *userQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*record.UserRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedUsersWithPaginationParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedUsersWithPagination(ctx, reqDb)

	if err != nil {
		return nil, nil, user_errors.ErrFindTrashedUsers
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	so := r.mapper.ToUsersRecordTrashedPagination(res)

	return so, &totalCount, nil
}

// FindByEmail retrieves a user by their email address.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - email: The email address of the user.
//
// Returns:
//   - *record.UserRecord: The user record if found.
//   - error: Error if retrieval fails.
func (r *userQueryRepository) FindByEmail(ctx context.Context, email string) (*record.UserRecord, error) {
	res, err := r.db.GetUserByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user_errors.ErrUserNotFound
		}

		return nil, user_errors.ErrUserNotFound
	}

	so := r.mapper.ToUserRecord(res)
	return so, nil
}
