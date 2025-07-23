package handler

import (
	"context"

	pb "github.com/MamangRust/monolith-payment-gateway-pb/transaction"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	transaction_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/transaction_errors/grpc"
	protomapper "github.com/MamangRust/monolith-payment-gateway-shared/mapper/proto/transaction"
	"github.com/MamangRust/monolith-payment-gateway-transaction/internal/service"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type transactionCommandHandleGrpc struct {
	pb.UnimplementedTransactionCommandServiceServer

	service service.TransactionCommandService
	logger  logger.LoggerInterface
	mapper  protomapper.TransactionCommandProtoMapper
}

func NewTransactionCommandHandleGrpc(service service.TransactionCommandService, logger logger.LoggerInterface, mapper protomapper.TransactionCommandProtoMapper) TransactionCommandHandleGrpc {
	return &transactionCommandHandleGrpc{
		service: service,
		logger:  logger,
		mapper:  mapper,
	}
}

// CreateTransaction implements the gRPC service for creating a new transaction.
// It validates the incoming request data, then calls the transaction command service to create the transaction.
// On success, it returns a gRPC response containing the created transaction data or an error if the creation operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - request: A pointer to a CreateTransactionRequest containing the transaction details.
//
// Returns:
//   - A pointer to ApiResponseTransaction containing the created transaction data on success.
//   - An error if the creation operation fails.
func (t *transactionCommandHandleGrpc) CreateTransaction(ctx context.Context, request *pb.CreateTransactionRequest) (*pb.ApiResponseTransaction, error) {
	transactionTime := request.GetTransactionTime().AsTime()
	merchantID := int(request.GetMerchantId())

	req := &requests.CreateTransactionRequest{
		CardNumber:      request.GetCardNumber(),
		Amount:          int(request.GetAmount()),
		PaymentMethod:   request.GetPaymentMethod(),
		MerchantID:      &merchantID,
		TransactionTime: transactionTime,
	}

	t.logger.Info("Creating transaction",
		zap.String("card_number", req.CardNumber),
		zap.Int("amount", req.Amount),
		zap.String("payment_method", req.PaymentMethod),
		zap.Int("merchant_id", *req.MerchantID),
		zap.Time("transaction_time", req.TransactionTime),
	)

	if err := req.Validate(); err != nil {
		t.logger.Error("failed to create transaction", zap.Any("error", err))
		return nil, transaction_errors.ErrGrpcValidateCreateTransactionRequest
	}

	res, err := t.service.Create(ctx, request.ApiKey, req)
	if err != nil {
		t.logger.Error("failed to create transaction", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapper.ToProtoResponseTransaction("success", "Successfully created transaction", res)

	t.logger.Info("Successfully created transaction",
		zap.String("card_number", req.CardNumber),
		zap.Int("amount", req.Amount),
		zap.String("payment_method", req.PaymentMethod),
		zap.Int("merchant_id", *req.MerchantID),
		zap.Time("transaction_time", req.TransactionTime),
	)

	return so, nil
}

// UpdateTransaction implements the gRPC service for updating an existing transaction by ID.
// It validates the incoming request data, then calls the transaction command service to update the transaction.
// On success, it returns a gRPC response containing the updated transaction data or an error if the update operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - request: A pointer to an UpdateTransactionRequest containing the updated transaction details.
//
// Returns:
//   - A pointer to ApiResponseTransaction containing the updated transaction data on success.
//   - An error if the update operation fails.
func (t *transactionCommandHandleGrpc) UpdateTransaction(ctx context.Context, request *pb.UpdateTransactionRequest) (*pb.ApiResponseTransaction, error) {
	id := int(request.GetTransactionId())

	t.logger.Info("Updating transaction",
		zap.Int("id", id))

	if id == 0 {
		t.logger.Error("failed to update transaction", zap.Any("error", transaction_errors.ErrGrpcTransactionInvalidID))
		return nil, transaction_errors.ErrGrpcTransactionInvalidID
	}

	transactionTime := request.GetTransactionTime().AsTime()
	merchantID := int(request.GetMerchantId())

	req := &requests.UpdateTransactionRequest{
		TransactionID:   &id,
		CardNumber:      request.GetCardNumber(),
		Amount:          int(request.GetAmount()),
		PaymentMethod:   request.GetPaymentMethod(),
		MerchantID:      &merchantID,
		TransactionTime: transactionTime,
	}

	if err := req.Validate(); err != nil {
		t.logger.Error("failed to update transaction", zap.Any("error", err))
		return nil, transaction_errors.ErrGrpcValidateUpdateTransactionRequest
	}

	res, err := t.service.Update(ctx, request.ApiKey, req)
	if err != nil {
		t.logger.Error("failed to update transaction", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapper.ToProtoResponseTransaction("success", "Successfully updated transaction", res)

	t.logger.Info("Successfully updated transaction",
		zap.Int("id", id),
		zap.String("card_number", req.CardNumber),
		zap.Int("amount", req.Amount),
		zap.String("payment_method", req.PaymentMethod),
		zap.Int("merchant_id", *req.MerchantID),
		zap.Time("transaction_time", req.TransactionTime),
	)

	return so, nil
}

// TrashedTransaction implements the gRPC service for trashing a transaction by its ID.
// It validates the incoming request data, then calls the transaction command service to trash the transaction.
// On success, it returns a gRPC response containing the trashed transaction data or an error if the trashing operation fails.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - request: A pointer to a FindByIdTransactionRequest containing the transaction ID to trash.
//
// Returns:
//   - A pointer to ApiResponseTransaction containing the trashed transaction data on success.
//   - An error if the trashing operation fails.
func (t *transactionCommandHandleGrpc) TrashedTransaction(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransactionDeleteAt, error) {
	id := int(request.GetTransactionId())

	t.logger.Info("Trashing transaction",
		zap.Int("id", id))

	if id == 0 {
		t.logger.Error("failed to trashed transaction", zap.Any("error", transaction_errors.ErrGrpcTransactionInvalidID))
		return nil, transaction_errors.ErrGrpcTransactionInvalidID
	}

	res, err := t.service.TrashedTransaction(ctx, id)

	if err != nil {
		t.logger.Error("failed to trashed transaction", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapper.ToProtoResponseTransactionDeleteAt("success", "Successfully trashed transaction", res)

	t.logger.Info("Successfully trashed transaction",
		zap.Int("id", id),
	)

	return so, nil
}

// RestoreTransaction implements the gRPC service for restoring a transaction from trashed.
//
// It logs the process of restoring the transaction and validates the input id, ensuring that it is a positive
// integer. If the id is invalid, it logs the error and returns an appropriate error response. The function queries
// the transaction command service to restore the transaction and maps the result to a protocol buffer response
// object.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - request: A pointer to a FindByIdTransactionRequest containing the transaction id to restore.
//
// Returns:
//   - A pointer to ApiResponseTransaction containing the restored transaction data on success.
//   - An error if the restoration operation fails, or if the provided id is invalid.
func (t *transactionCommandHandleGrpc) RestoreTransaction(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransaction, error) {
	id := int(request.GetTransactionId())

	t.logger.Info("Restoring transaction",
		zap.Int("id", id))

	if id == 0 {
		t.logger.Error("failed to restore transaction", zap.Any("error", transaction_errors.ErrGrpcTransactionInvalidID))
		return nil, transaction_errors.ErrGrpcTransactionInvalidID
	}

	res, err := t.service.RestoreTransaction(ctx, id)

	if err != nil {
		t.logger.Error("failed to restore transaction", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapper.ToProtoResponseTransaction("success", "Successfully restored transaction", res)

	t.logger.Info("Successfully restored transaction",
		zap.Int("id", id),
	)

	return so, nil
}

// DeleteTransaction implements the gRPC service for permanently deleting a transaction by ID.
// It logs the deletion process and validates the transaction ID, ensuring it is a positive integer.
// If the ID is invalid, it logs the error and returns an appropriate error response. The function
// calls the transaction command service to delete the transaction permanently and maps the result
// to a protocol buffer response object.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - request: A pointer to a FindByIdTransactionRequest containing the transaction ID to delete.
//
// Returns:
//   - A pointer to ApiResponseTransactionDelete containing the deletion confirmation on success.
//   - An error if the deletion operation fails, or if the provided ID is invalid.
func (t *transactionCommandHandleGrpc) DeleteTransaction(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransactionDelete, error) {
	id := int(request.GetTransactionId())

	t.logger.Info("Deleting transaction",
		zap.Int("id", id))

	if id == 0 {
		t.logger.Error("failed to delete transaction", zap.Any("error", transaction_errors.ErrGrpcTransactionInvalidID))
		return nil, transaction_errors.ErrGrpcTransactionInvalidID
	}

	_, err := t.service.DeleteTransactionPermanent(ctx, id)

	if err != nil {
		t.logger.Error("failed to delete transaction", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapper.ToProtoResponseTransactionDelete("success", "Successfully deleted transaction")

	t.logger.Info("Successfully deleted transaction",
		zap.Int("id", id),
	)

	return so, nil

}

// RestoreAllTransaction implements the gRPC service for restoring all trashed transactions.
//
// It logs the process of restoring all transactions and queries the transaction command service to restore all
// transactions. The function maps the result to a protocol buffer response object.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - _ *emptypb.Empty: Unused parameter.
//
// Returns:
//   - A pointer to ApiResponseTransactionAll containing the restored transactions on success.
//   - An error if the restoration operation fails.
func (t *transactionCommandHandleGrpc) RestoreAllTransaction(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTransactionAll, error) {
	t.logger.Info("Restoring all transaction")

	_, err := t.service.RestoreAllTransaction(ctx)

	if err != nil {
		t.logger.Error("failed to restore all transaction", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapper.ToProtoResponseTransactionAll("success", "Successfully restore all transaction")

	t.logger.Info("Successfully restore all transaction")

	return so, nil
}

// DeleteAllTransactionPermanent implements the gRPC service for permanently deleting all trashed transactions.
//
// It logs the process of deleting all transactions and queries the transaction command service to delete all
// transactions. The function maps the result to a protocol buffer response object.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation signals, and deadlines.
//   - _ *emptypb.Empty: Unused parameter.
//
// Returns:
//   - A pointer to ApiResponseTransactionAll containing the deleted transactions on success.
//   - An error if the deletion operation fails.
func (t *transactionCommandHandleGrpc) DeleteAllTransactionPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTransactionAll, error) {
	t.logger.Info("Deleting all transaction permanent")

	_, err := t.service.DeleteAllTransactionPermanent(ctx)

	if err != nil {
		t.logger.Error("failed to delete all transaction permanent", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := t.mapper.ToProtoResponseTransactionAll("success", "Successfully delete transaction permanent")

	t.logger.Info("Successfully delete all transaction permanent")

	return so, nil
}
