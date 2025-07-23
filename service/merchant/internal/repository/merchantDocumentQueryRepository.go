package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	merchantdocument_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_document_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/merchantdocument"
)

// merchantDocumentQueryRepository provides methods to query merchant documents from the database.
type merchantDocumentQueryRepository struct {
	db     *db.Queries
	mapper recordmapper.MerchantDocumentQueryMapper
}

// NewMerchantDocumentQueryRepository creates a new instance of merchantDocumentQueryRepository with the provided
// database queries, context, and merchant document record mapper. This repository is responsible for querying
// merchant documents from the database.
//
// Parameters:
//   - db: A pointer to the db.Queries object for executing database queries.
//   - ctx: The context to be used for database operations, allowing for cancellation and timeout.
//   - mapper: A MerchantDocumentMapping that provides methods to map database rows to MerchantDocumentRecord domain models.
//
// Returns:
//   - A pointer to the newly created merchantDocumentQueryRepository instance.
func NewMerchantDocumentQueryRepository(db *db.Queries, mapper recordmapper.MerchantDocumentQueryMapper) MerchantDocumentQueryRepository {
	return &merchantDocumentQueryRepository{
		db:     db,
		mapper: mapper,
	}
}

// FindAllDocuments queries the database for all merchant documents matching the search query and returns the records paginated by the given page and page size.
// The returned records are mapped to MerchantDocumentRecord domain models using the provided mapper.
//
// Parameters:
//   - req: A pointer to the FindAllMerchantDocuments request object containing the search query, page, and page size.
//
// Returns:
//   - A slice of MerchantDocumentRecord objects representing the matching merchant documents.
//   - A pointer to an integer representing the total number of matching records.
//   - An error if the query fails, which is wrapped in ErrFindAllMerchantDocumentsFailed.
func (r *merchantDocumentQueryRepository) FindAllDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*record.MerchantDocumentRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	params := db.GetMerchantDocumentsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	docs, err := r.db.GetMerchantDocuments(ctx, params)
	if err != nil {
		return nil, nil, merchantdocument_errors.ErrFindAllMerchantDocumentsFailed
	}

	var totalCount int
	if len(docs) > 0 {
		totalCount = int(docs[0].TotalCount)
	}

	return r.mapper.ToMerchantDocumentsRecord(docs), &totalCount, nil
}

// FindByActive queries the database for all active merchant documents matching the search query and returns the records paginated by the given page and page size.
// The returned records are mapped to MerchantDocumentRecord domain models using the provided mapper.
//
// Parameters:
//   - req: A pointer to the FindAllMerchantDocuments request object containing the search query, page, and page size.
//
// Returns:
//   - A slice of MerchantDocumentRecord objects representing the matching active merchant documents.
//   - A pointer to an integer representing the total number of matching records.
//   - An error if the query fails, which is wrapped in ErrFindActiveMerchantDocumentsFailed.
func (r *merchantDocumentQueryRepository) FindByActiveDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*record.MerchantDocumentRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	params := db.GetActiveMerchantDocumentsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	docs, err := r.db.GetActiveMerchantDocuments(ctx, params)
	if err != nil {
		return nil, nil, merchantdocument_errors.ErrFindActiveMerchantDocumentsFailed
	}

	var totalCount int
	if len(docs) > 0 {
		totalCount = int(docs[0].TotalCount)
	}

	return r.mapper.ToMerchantDocumentsActiveRecord(docs), &totalCount, nil
}

// FindByTrashed queries the database for all trashed merchant documents matching the search query and returns the records paginated by the given page and page size.
// The returned records are mapped to MerchantDocumentRecord domain models using the provided mapper.
//
// Parameters:
//   - req: A pointer to the FindAllMerchantDocuments request object containing the search query, page, and page size.
//
// Returns:
//   - A slice of MerchantDocumentRecord objects representing the matching trashed merchant documents.
//   - A pointer to an integer representing the total number of matching records.
//   - An error if the query fails, which is wrapped in ErrFindTrashedMerchantDocumentsFailed.
func (r *merchantDocumentQueryRepository) FindByTrashedDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*record.MerchantDocumentRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	params := db.GetTrashedMerchantDocumentsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	docs, err := r.db.GetTrashedMerchantDocuments(ctx, params)
	if err != nil {
		return nil, nil, merchantdocument_errors.ErrFindTrashedMerchantDocumentsFailed
	}

	var totalCount int
	if len(docs) > 0 {
		totalCount = int(docs[0].TotalCount)
	}

	return r.mapper.ToMerchantDocumentsTrashedRecord(docs), &totalCount, nil
}

// FindById queries the database for a merchant document by ID and returns the record.
// The returned record is mapped to a MerchantDocumentRecord domain model using the provided mapper.
//
// Parameters:
//   - id: The ID of the merchant document to retrieve.
//
// Returns:
//   - A pointer to a MerchantDocumentRecord object representing the merchant document.
//   - An error if the query fails, which is wrapped in ErrFindMerchantDocumentByIdFailed.
func (r *merchantDocumentQueryRepository) FindByIdDocument(ctx context.Context, id int) (*record.MerchantDocumentRecord, error) {
	doc, err := r.db.GetMerchantDocument(ctx, int32(id))
	if err != nil {
		return nil, merchantdocument_errors.ErrFindMerchantDocumentByIdFailed
	}
	return r.mapper.ToGetMerchantDocument(doc), nil
}
