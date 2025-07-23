package repository

import (
	"context"
	"database/sql"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	merchantdocument_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_document_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/merchantdocument"
)

// merchantDocumentCommandRepository implements the MerchantDocumentCommandRepository interface
type merchantDocumentCommandRepository struct {
	db     *db.Queries
	mapper recordmapper.MerchantDocumentCommandMapper
}

// NewMerchantDocumentCommandRepository initializes a new instance of
// merchantDocumentCommandRepository with the provided database queries,
// context, and card record mapper. This repository is responsible for
// executing command operations related to merchant document records in the
// database.
//
// Parameters:
//   - db: A pointer to the db.Queries object for executing database queries.
//   - ctx: The context to be used for database operations, allowing for
//     cancellation and timeout.
//   - mapper: A MerchantDocumentMapping that provides methods to map
//     database rows to MerchantDocumentRecord domain models.
//
// Returns:
//   - A pointer to the newly created merchantDocumentCommandRepository
//     instance.
func NewMerchantDocumentCommandRepository(db *db.Queries, mapper recordmapper.MerchantDocumentCommandMapper) MerchantDocumentCommandRepository {
	return &merchantDocumentCommandRepository{
		db:     db,
		mapper: mapper,
	}
}

// CreateMerchantDocument creates a new merchant document with the given request
// parameters and returns a domain model object for the newly created document
// record. The request must contain a valid merchant ID, document type, and
// document URL. The document status is set to "pending" by default. If the
// operation fails to create a new document, an error is returned.
//
// Parameters:
//   - request: A pointer to the CreateMerchantDocumentRequest object
//     containing the merchant ID, document type, document URL, and any
//     additional information to be stored for the document.
//
// Returns:
//   - A pointer to the newly created merchant document record domain model
//     (MerchantDocumentRecord).
//   - An error if the operation fails.
func (r *merchantDocumentCommandRepository) CreateMerchantDocument(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*record.MerchantDocumentRecord, error) {
	req := db.CreateMerchantDocumentParams{
		MerchantID:   int32(request.MerchantID),
		DocumentType: request.DocumentType,
		DocumentUrl:  request.DocumentUrl,
		Status:       "pending",
		Note:         sql.NullString{String: "", Valid: true},
	}

	res, err := r.db.CreateMerchantDocument(ctx, req)
	if err != nil {
		return nil, merchantdocument_errors.ErrCreateMerchantDocumentFailed
	}

	return r.mapper.ToGetMerchantDocument(res), nil
}

// UpdateMerchantDocument updates an existing merchant document with the provided
// request parameters. It modifies the document type, URL, status, and note based
// on the request. If the update operation fails, an error is returned.
//
// Parameters:
//   - request: A pointer to the UpdateMerchantDocumentRequest object containing
//     the document ID, merchant ID, document type, document URL, status, and note.
//
// Returns:
//   - A pointer to the updated merchant document record domain model
//     (MerchantDocumentRecord).
//   - An error if the operation fails.
func (r *merchantDocumentCommandRepository) UpdateMerchantDocument(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*record.MerchantDocumentRecord, error) {
	req := db.UpdateMerchantDocumentParams{
		DocumentID:   int32(request.MerchantID),
		DocumentType: request.DocumentType,
		DocumentUrl:  request.DocumentUrl,
		Status:       request.Status,
		Note:         sql.NullString{String: request.Note, Valid: true},
	}

	res, err := r.db.UpdateMerchantDocument(ctx, req)
	if err != nil {
		return nil, merchantdocument_errors.ErrUpdateMerchantDocumentFailed
	}

	return r.mapper.ToGetMerchantDocument(res), nil
}

// UpdateMerchantDocumentStatus updates the status and note of an existing
// merchant document based on the provided request parameters. It modifies
// the status and note fields of the document. If the update operation fails,
// an error is returned.
//
// Parameters:
//   - request: A pointer to the UpdateMerchantDocumentStatusRequest object
//     containing the document ID, status, and note.
//
// Returns:
//   - A pointer to the updated merchant document record domain model
//     (MerchantDocumentRecord).
//   - An error if the operation fails.
func (r *merchantDocumentCommandRepository) UpdateMerchantDocumentStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*record.MerchantDocumentRecord, error) {
	req := db.UpdateMerchantDocumentStatusParams{
		DocumentID: int32(request.MerchantID),
		Status:     request.Status,
		Note:       sql.NullString{String: request.Note, Valid: true},
	}

	res, err := r.db.UpdateMerchantDocumentStatus(ctx, req)
	if err != nil {
		return nil, merchantdocument_errors.ErrUpdateMerchantDocumentStatusFailed
	}

	return r.mapper.ToGetMerchantDocument(res), nil
}

// TrashedMerchantDocument trashes a merchant document by its ID and returns the updated record.
//
// Parameters:
//   - documentID: The ID of the merchant document to trash.
//
// Returns:
//   - A pointer to the updated merchant document record domain model
//     (MerchantDocumentRecord).
//   - An error if the operation fails.
func (r *merchantDocumentCommandRepository) TrashedMerchantDocument(ctx context.Context, documentID int) (*record.MerchantDocumentRecord, error) {
	res, err := r.db.TrashMerchantDocument(ctx, int32(documentID))
	if err != nil {
		return nil, merchantdocument_errors.ErrTrashedMerchantDocumentFailed
	}

	return r.mapper.ToGetMerchantDocument(res), nil
}

// RestoreMerchantDocument restores a merchant by its ID and returns the updated record.
//
// Parameters:
//   - documentID: The ID of the merchant to restore.
//
// Returns:
//   - A pointer to a MerchantDocumentRecord containing the updated record.
//   - An error if the record could not be restored.
func (r *merchantDocumentCommandRepository) RestoreMerchantDocument(ctx context.Context, documentID int) (*record.MerchantDocumentRecord, error) {
	res, err := r.db.RestoreMerchantDocument(ctx, int32(documentID))
	if err != nil {
		return nil, merchantdocument_errors.ErrRestoreMerchantDocumentFailed
	}

	return r.mapper.ToGetMerchantDocument(res), nil
}

// DeleteMerchantDocumentPermanent permanently deletes a merchant document by its ID.
//
// Parameters:
//   - documentID: The ID of the document to be deleted.
//
// Returns:
//   - A boolean indicating whether the deletion was successful.
//   - An error if the deletion failed.
func (r *merchantDocumentCommandRepository) DeleteMerchantDocumentPermanent(ctx context.Context, documentID int) (bool, error) {
	err := r.db.DeleteMerchantDocumentPermanently(ctx, int32(documentID))
	if err != nil {
		return false, merchantdocument_errors.ErrDeleteMerchantDocumentPermanentFailed
	}

	return true, nil
}

// RestoreAllMerchantDocument restores all merchant documents.
//
// Returns:
//   - A boolean indicating whether the restoration was successful.
//   - An error if the restoration failed.
func (r *merchantDocumentCommandRepository) RestoreAllMerchantDocument(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllMerchantDocuments(ctx)
	if err != nil {
		return false, merchantdocument_errors.ErrRestoreAllMerchantDocumentsFailed
	}

	return true, nil
}

// DeleteAllMerchantDocumentPermanent deletes all merchant documents permanently.
//
// Returns:
//   - A boolean indicating whether the deletion was successful.
//   - An error if the deletion failed.
func (r *merchantDocumentCommandRepository) DeleteAllMerchantDocumentPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentMerchantDocuments(ctx)
	if err != nil {
		return false, merchantdocument_errors.ErrDeleteAllMerchantDocumentsPermanentFailed
	}

	return true, nil
}
