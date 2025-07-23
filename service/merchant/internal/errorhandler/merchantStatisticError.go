package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// merchantStatisticError is a struct that implements the MerchantStatisticErrorHandler interface
type merchantStatisticError struct {
	logger logger.LoggerInterface
}

// NewMerchantStatisticError returns a new instance of MerchantStatisticError with the given logger.
// It returns an instance of the merchantStatisticError struct.
func NewMerchantStatisticError(logger logger.LoggerInterface) MerchantStatisticErrorHandler {
	return &merchantStatisticError{
		logger: logger,
	}
}

// HandleMonthlyPaymentMethodsMerchantError processes errors that occur during the retrieval of monthly payment
// methods for a merchant.
// It logs the error, records it to the trace span, and returns a standardized error response.
// Args:
//   - err: The error that occurred.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of MerchantResponseMonthlyPaymentMethod containing the payment methods or nil in case of an error.
//   - A standardized ErrorResponse detailing the error encountered during the retrieval process.
func (e *merchantStatisticError) HandleMonthlyPaymentMethodsMerchantError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponseMonthlyPaymentMethod](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindMonthlyPaymentMethodsMerchant, fields...,
	)
}

// HandleYearlyPaymentMethodMerchantError processes errors that occur during the retrieval of yearly payment
// methods for a merchant.
// It logs the error, records it to the trace span, and returns a standardized error response.
// Args:
//   - err: The error that occurred.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of MerchantResponseYearlyPaymentMethod containing the payment methods or nil in case of an error.
//   - A standardized ErrorResponse detailing the error encountered during the retrieval process.
func (e *merchantStatisticError) HandleYearlyPaymentMethodMerchantError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponseYearlyPaymentMethod](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindYearlyPaymentMethodMerchant, fields...,
	)
}

// HandleMonthlyAmountMerchantError processes errors that occur during the retrieval of monthly amounts for a merchant.
// It logs the error, records it to the trace span, and returns a standardized error response.
// Args:
//   - err: The error that occurred.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - statuus: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of MerchantResponseMonthlyAmount containing the monthly amounts or nil in case of an error.
//   - A standardized ErrorResponse detailing the error encountered during the retrieval process.
func (e *merchantStatisticError) HandleMonthlyAmountMerchantError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	statuus *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponseMonthlyAmount](
		e.logger, err, method, tracePrefix, span, nil, merchant_errors.ErrFailedFindMonthlyAmountMerchant, fields...,
	)
}

// HandleYearlyAmountMerchantError processes errors that occur during the retrieval of yearly amounts for a merchant.
// It logs the error, records it to the trace span, and returns a standardized error response.
// Args:
//   - err: The error that occurred.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of MerchantResponseYearlyAmount containing the yearly amounts or nil in case of an error.
//   - A standardized ErrorResponse detailing the error encountered during the retrieval process.
func (e *merchantStatisticError) HandleYearlyAmountMerchantError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponseYearlyAmount](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindYearlyAmountMerchant, fields...,
	)
}

// HandleMonthlyTotalAmountMerchantError processes errors that occur during the retrieval of monthly total amounts for a
// merchant.
// It logs the error, records it to the trace span, and returns a standardized error response.
// Args:
//   - err: The error that occurred.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of MerchantResponseMonthlyTotalAmount containing the monthly total amounts or nil in case of an error.
//   - A standardized ErrorResponse detailing the error encountered during the retrieval process.
func (e *merchantStatisticError) HandleMonthlyTotalAmountMerchantError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponseMonthlyTotalAmount](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindMonthlyTotalAmountMerchant, fields...,
	)
}

// HandleYearlyTotalAmountMerchantError processes errors that occur during the retrieval of yearly total amounts for a
// merchant.
// It logs the error, records it to the trace span, and returns a standardized error response.
// Args:
//   - err: The error that occurred.
//   - method: The name of the method where the error occurred.
//   - tracePrefix: A prefix for generating the trace ID.
//   - span: The trace span used for recording the error.
//   - status: A pointer to a string that will be set with the formatted status.
//   - fields: Additional fields to include in the log entry.
//
// Returns:
//   - A slice of MerchantResponseYearlyTotalAmount containing the yearly total amounts or nil in case of an error.
//   - A standardized ErrorResponse detailing the error encountered during the retrieval process.
func (e *merchantStatisticError) HandleYearlyTotalAmountMerchantError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponseYearlyTotalAmount](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindYearlyTotalAmountMerchant, fields...,
	)
}
