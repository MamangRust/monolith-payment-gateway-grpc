package repository

import (
	"context"
	"errors"

	apikey "github.com/MamangRust/monolith-payment-gateway-pkg/api-key"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/merchant"
)

// MerchantCommandRepository defines operations for merchant data persistence
type merchantCommandRepository struct {
	db     *db.Queries
	mapper recordmapper.MerchantCommandRecordMapper
}

// NewMerchantCommandRepository creates a new instance of merchantCommandRepository with the provided
// database queries, context, and merchant record mapper. This repository is responsible for executing
// command operations related to merchant records in the database.
//
// Parameters:
//   - db: A pointer to the db.Queries object for executing database queries.
//   - mapper: A MerchantRecordMapping that provides methods to map database rows to Merchant domain models.
//
// Returns:
//   - A pointer to the newly created merchantCommandRepository instance.
func NewMerchantCommandRepository(db *db.Queries, mapper recordmapper.MerchantCommandRecordMapper) MerchantCommandRepository {
	return &merchantCommandRepository{
		db:     db,
		mapper: mapper,
	}
}

// CreateMerchant creates a new merchant record with the provided name and user ID and returns the newly
// created record. The status of the merchant is set to "inactive" by default.
//
// Parameters:
//   - request: A CreateMerchantRequest containing the name, user ID and status of the merchant.
//
// Returns:
//   - A pointer to a MerchantRecord containing the newly created record.
//   - An error if the record could not be created.
func (r *merchantCommandRepository) CreateMerchant(ctx context.Context, request *requests.CreateMerchantRequest) (*record.MerchantRecord, error) {
	apiKey, err := apikey.GenerateApiKey()

	if err != nil {
		errApiKey := errors.New("error generate api key")

		return nil, errApiKey
	}

	req := db.CreateMerchantParams{
		Name:   request.Name,
		ApiKey: apiKey,
		UserID: int32(request.UserID),
		Status: "inactive",
	}

	res, err := r.db.CreateMerchant(ctx, req)

	if err != nil {
		return nil, merchant_errors.ErrCreateMerchantFailed
	}

	return r.mapper.ToMerchantRecord(res), nil
}

// UpdateMerchant updates the merchant record with the provided merchant ID, name, user ID, and status
// and returns the updated record.
//
// Parameters:
//   - request: A UpdateMerchantRequest containing the merchant ID, name, user ID, and status of the merchant.
//
// Returns:
//   - A pointer to a MerchantRecord containing the updated record.
//   - An error if the record could not be updated.
func (r *merchantCommandRepository) UpdateMerchant(ctx context.Context, request *requests.UpdateMerchantRequest) (*record.MerchantRecord, error) {
	req := db.UpdateMerchantParams{
		MerchantID: int32(*request.MerchantID),
		Name:       request.Name,
		UserID:     int32(request.UserID),
		Status:     request.Status,
	}

	res, err := r.db.UpdateMerchant(ctx, req)

	if err != nil {
		return nil, merchant_errors.ErrUpdateMerchantFailed
	}

	return r.mapper.ToMerchantRecord(res), nil
}

// UpdateMerchantStatus updates the status of the merchant with the provided merchant ID to the
// provided status and returns the updated record.
//
// Parameters:
//   - request: A UpdateMerchantStatusRequest containing the merchant ID and status of the merchant.
//
// Returns:
//   - A pointer to a MerchantRecord containing the updated record.
//   - An error if the record could not be updated.
func (r *merchantCommandRepository) UpdateMerchantStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*record.MerchantRecord, error) {
	req := db.UpdateMerchantStatusParams{
		MerchantID: int32(*request.MerchantID),
		Status:     request.Status,
	}

	res, err := r.db.UpdateMerchantStatus(ctx, req)

	if err != nil {
		return nil, merchant_errors.ErrUpdateMerchantStatusFailed
	}

	return r.mapper.ToMerchantRecord(res), nil
}

// TrashedMerchant trashes a merchant by its ID and returns the updated record.
//
// Parameters:
//   - merchant_id: The ID of the merchant to trash.
//
// Returns:
//   - A pointer to a MerchantRecord containing the updated record.
//   - An error if the record could not be trashed.
func (r *merchantCommandRepository) TrashedMerchant(ctx context.Context, merchant_id int) (*record.MerchantRecord, error) {
	res, err := r.db.TrashMerchant(ctx, int32(merchant_id))

	if err != nil {
		return nil, merchant_errors.ErrTrashedMerchantFailed
	}

	return r.mapper.ToMerchantRecord(res), nil
}

// RestoreMerchant restores a merchant by its ID and returns the updated record.
//
// Parameters:
//   - merchant_id: The ID of the merchant to restore.
//
// Returns:
//   - A pointer to a MerchantRecord containing the updated record.
//   - An error if the record could not be restored.
func (r *merchantCommandRepository) RestoreMerchant(ctx context.Context, merchant_id int) (*record.MerchantRecord, error) {
	res, err := r.db.RestoreMerchant(ctx, int32(merchant_id))

	if err != nil {
		return nil, merchant_errors.ErrRestoreMerchantFailed
	}

	return r.mapper.ToMerchantRecord(res), nil
}

// DeleteMerchantPermanent deletes a merchant permanently by its ID.
//
// Parameters:
//   - merchant_id: The ID of the merchant to delete.
//
// Returns:
//   - A boolean indicating whether the deletion was successful.
//   - An error if the deletion failed.
func (r *merchantCommandRepository) DeleteMerchantPermanent(ctx context.Context, merchant_id int) (bool, error) {
	err := r.db.DeleteMerchantPermanently(ctx, int32(merchant_id))

	if err != nil {
		return false, merchant_errors.ErrDeleteMerchantPermanentFailed
	}

	return true, nil
}

// RestoreAllMerchant restores all merchants that were previously trashed.
//
// Returns:
//   - A boolean indicating whether the operation was successful.
//   - An error if the operation failed.
func (r *merchantCommandRepository) RestoreAllMerchant(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllMerchants(ctx)

	if err != nil {
		return false, merchant_errors.ErrRestoreAllMerchantFailed
	}

	return true, nil
}

// DeleteAllMerchantPermanent deletes all merchants permanently.
//
// Returns:
//   - A boolean indicating whether the deletion was successful.
//   - An error if the deletion failed.
func (r *merchantCommandRepository) DeleteAllMerchantPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentMerchants(ctx)

	if err != nil {
		return false, merchant_errors.ErrDeleteAllMerchantPermanentFailed
	}

	return true, nil
}
