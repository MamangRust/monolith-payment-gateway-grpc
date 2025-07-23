package repository

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/record"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/repository"
	recordmapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/record/saldo"
)

// saldoQueryRepository is a struct that implements the SaldoQueryRepository interface
type saldoQueryRepository struct {
	db     *db.Queries
	mapper recordmapper.SaldoQueryRecordMapping
}

// NewSaldoQueryRepository initializes a new instance of saldoQueryRepository with the provided
// database queries, context, and saldo record mapper. This repository is responsible for executing
// query operations related to saldo records in the database.
//
// Parameters:
//   - db: A pointer to the db.Queries object for executing database queries.
//   - mapper: A SaldoRecordMapping that provides methods to map database rows to SaldoRecord domain models.
//
// Returns:
//   - A pointer to the newly created saldoQueryRepository instance.
func NewSaldoQueryRepository(db *db.Queries, mapper recordmapper.SaldoQueryRecordMapping) SaldoQueryRepository {
	return &saldoQueryRepository{
		db:     db,
		mapper: mapper,
	}
}

// FindAllSaldos retrieves all saldo records based on provided filters.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request object containing pagination, filtering, or search criteria.
//
// Returns:
//   - []*record.SaldoRecord: The list of saldo records.
//   - *int: The total number of records found.
//   - error: An error if the query fails.
func (r *saldoQueryRepository) FindAllSaldos(ctx context.Context, req *requests.FindAllSaldos) ([]*record.SaldoRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetSaldosParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	saldos, err := r.db.GetSaldos(ctx, reqDb)

	if err != nil {
		return nil, nil, saldo_errors.ErrFindAllSaldosFailed
	}

	var totalCount int
	if len(saldos) > 0 {
		totalCount = int(saldos[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToSaldosRecordAll(saldos), &totalCount, nil
}

// FindByActive retrieves all active saldo records (not soft-deleted).
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request object containing pagination or filtering options.
//
// Returns:
//   - []*record.SaldoRecord: The list of active saldo records.
//   - *int: The total number of active records.
//   - error: An error if the query fails.
func (r *saldoQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllSaldos) ([]*record.SaldoRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveSaldosParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveSaldos(ctx, reqDb)

	if err != nil {
		return nil, nil, saldo_errors.ErrFindActiveSaldosFailed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToSaldosRecordActive(res), &totalCount, nil

}

// FindByTrashed retrieves all trashed saldo records (soft-deleted).
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - req: The request object containing pagination or filtering options.
//
// Returns:
//   - []*record.SaldoRecord: The list of trashed saldo records.
//   - *int: The total number of trashed records.
//   - error: An error if the query fails.
func (r *saldoQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllSaldos) ([]*record.SaldoRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedSaldosParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	saldos, err := r.db.GetTrashedSaldos(ctx, reqDb)

	if err != nil {
		return nil, nil, saldo_errors.ErrFindTrashedSaldosFailed
	}

	var totalCount int
	if len(saldos) > 0 {
		totalCount = int(saldos[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapper.ToSaldosRecordTrashed(saldos), &totalCount, nil
}

// FindByCardNumber retrieves a saldo record by card number.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - card_number: The card number associated with the saldo.
//
// Returns:
//   - *record.SaldoRecord: The saldo record if found.
//   - error: An error if the record is not found or query fails.
func (r *saldoQueryRepository) FindByCardNumber(ctx context.Context, card_number string) (*record.SaldoRecord, error) {
	res, err := r.db.GetSaldoByCardNumber(ctx, card_number)

	if err != nil {
		return nil, saldo_errors.ErrFindSaldoByCardNumberFailed
	}

	return r.mapper.ToSaldoRecord(res), nil
}

// FindById retrieves a saldo record by its unique ID.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - saldo_id: The unique saldo ID to query.
//
// Returns:
//   - *record.SaldoRecord: The saldo record if found.
//   - error: An error if the record is not found or query fails.
func (r *saldoQueryRepository) FindById(ctx context.Context, saldo_id int) (*record.SaldoRecord, error) {
	res, err := r.db.GetSaldoByID(ctx, int32(saldo_id))

	if err != nil {
	}

	return r.mapper.ToSaldoRecord(res), nil
}
