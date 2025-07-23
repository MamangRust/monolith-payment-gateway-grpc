package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	merchant_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/merchant_errors/service"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// merchantStatisticByApiKeyError is a struct that implements the MerchantStatisticByApiKeyErrorHandler interface
type merchantStatisticByApiKeyError struct {
	logger logger.LoggerInterface
}

// NewMerchantStatisticByApiKeyError returns a new instance of MerchantStatisticByApiKeyError with the given logger.
func NewMerchantStatisticByApiKeyError(logger logger.LoggerInterface) MerchantStatisticByApikeyErrorHandler {
	return &merchantStatisticByApiKeyError{
		logger: logger,
	}
}

// HandleMonthlyPaymentMethodByApikeysError processes errors that occur during the retrieval of monthly payment
// methods by API key.
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
func (e *merchantStatisticByApiKeyError) HandleMonthlyPaymentMethodByApikeysError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseMonthlyPaymentMethod, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponseMonthlyPaymentMethod](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindMonthlyPaymentMethodByApikeys, fields...,
	)
}

// HandleYearlyPaymentMethodByApikeysError processes errors that occur during the retrieval of yearly payment
// methods by API key.
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
func (e *merchantStatisticByApiKeyError) HandleYearlyPaymentMethodByApikeysError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseYearlyPaymentMethod, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponseYearlyPaymentMethod](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindYearlyPaymentMethodByApikeys, fields...,
	)
}

// HandleMonthlyAmountByApikeysError processes errors that occur during the retrieval of monthly amounts by API key.
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
func (e *merchantStatisticByApiKeyError) HandleMonthlyAmountByApikeysError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseMonthlyAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponseMonthlyAmount](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindMonthlyAmountByApikeys, fields...,
	)
}

// HandleYearlyAmountByApikeysError processes errors that occur during the retrieval of yearly amounts by API key.
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
func (e *merchantStatisticByApiKeyError) HandleYearlyAmountByApikeysError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseYearlyAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponseYearlyAmount](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindYearlyAmountByApikeys, fields...,
	)
}

// HandleMonthlyTotalAmountByApikeysError processes errors that occur during the retrieval of monthly total amounts by API key.
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
func (e *merchantStatisticByApiKeyError) HandleMonthlyTotalAmountByApikeysError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseMonthlyTotalAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponseMonthlyTotalAmount](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindMonthlyTotalAmountByApikeys, fields...,
	)
}

// HandleYearlyTotalAmountByApikeysError processes errors that occur during the retrieval of yearly total amounts by API key.
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
//   - A list of MerchantResponseYearlyTotalAmount or nil if an error occurred.
//   - A standardized ErrorResponse detailing the yearly total amount retrieval failure.
func (e *merchantStatisticByApiKeyError) HandleYearlyTotalAmountByApikeysError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.MerchantResponseYearlyTotalAmount, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.MerchantResponseYearlyTotalAmount](
		e.logger, err, method, tracePrefix, span, status, merchant_errors.ErrFailedFindYearlyTotalAmountByApikeys, fields...,
	)
}
