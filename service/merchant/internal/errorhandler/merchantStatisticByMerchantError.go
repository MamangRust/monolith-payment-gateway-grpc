package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// merchantStatisticByMerchantError is a struct that implements the MerchantStatisticByMerchantErrorHandler interface
type merchantStatisticByMerchantError struct {
	logger logger.LoggerInterface
}

// NewMerchantStatisticByMerchantError returns a new instance of MerchantStatisticByMerchantError with the given logger.
// It returns an instance of the merchantStatisticByMerchantError struct.
func NewMerchantStatisticByMerchantError(logger logger.LoggerInterface) MerchantStatisticByMerchantErrorHandler {
	return &merchantStatisticByMerchantError{
		logger: logger,
	}
}

// HandleMonthlyPaymentMethodByMerchantsError processes errors that occur during the retrieval of monthly payment
// methods by merchants.
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
//   - A boolean indicating whether the error is fatal.
//   - A standardized ErrorResponse detailing the payment method retrieval failure.
func (e *merchantStatisticByMerchantError) HandleMonthlyPaymentMethodByMerchantsError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponseMonthlyPaymentMethod](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindMonthlyPaymentMethodByMerchants, fields...,
	)
}

// HandleYearlyPaymentMethodByMerchantsError processes errors that occur during the retrieval of yearly payment
// methods by merchants.
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
//   - A boolean indicating whether the error is fatal.
//   - A standardized ErrorResponse detailing the payment method retrieval failure.
func (e *merchantStatisticByMerchantError) HandleYearlyPaymentMethodByMerchantsError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponseYearlyPaymentMethod](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindYearlyPaymentMethodByMerchants, fields...,
	)
}

// HandleMonthlyAmountByMerchantsError processes errors that occur during the retrieval of monthly amounts by merchants.
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
//   - A boolean indicating whether the error is fatal.
//   - A standardized ErrorResponse detailing the monthly amount retrieval failure.
func (e *merchantStatisticByMerchantError) HandleMonthlyAmountByMerchantsError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponseMonthlyAmount](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindMonthlyAmountByMerchants, fields...,
	)
}

// HandleYearlyAmountByMerchantsError processes errors that occur during the retrieval of yearly amounts by merchants.
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
//   - A boolean indicating whether the error is fatal.
//   - A standardized ErrorResponse detailing the yearly amount retrieval failure.
func (e *merchantStatisticByMerchantError) HandleYearlyAmountByMerchantsError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponseYearlyAmount](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindYearlyAmountByMerchants, fields...,
	)
}

// HandleMonthlyTotalAmountByMerchantsError processes errors that occur during the retrieval of monthly total amounts by merchants.
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
//   - A boolean indicating whether the error is fatal.
//   - A standardized ErrorResponse detailing the monthly total amount retrieval failure.
func (e *merchantStatisticByMerchantError) HandleMonthlyTotalAmountByMerchantsError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponseMonthlyTotalAmount](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindMonthlyTotalAmountByMerchants, fields...,
	)
}

// HandleYearlyTotalAmountByMerchantsError processes errors that occur during the retrieval of yearly total amounts by merchants.
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
//   - A boolean indicating whether the error is fatal.
//   - A standardized ErrorResponse detailing the yearly total amount retrieval failure.
func (e *merchantStatisticByMerchantError) HandleYearlyTotalAmountByMerchantsError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponseYearlyTotalAmount](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindYearlyTotalAmountByMerchants, fields...,
	)
}
