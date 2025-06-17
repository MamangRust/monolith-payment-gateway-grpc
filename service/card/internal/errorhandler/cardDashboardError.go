package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type cardDasboardError struct {
	logger logger.LoggerInterface
}

func NewCardDashboardError(logger logger.LoggerInterface) *cardDasboardError {
	return &cardDasboardError{
		logger: logger,
	}
}

func (c *cardDasboardError) HandleTotalBalanceError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.DashboardCard, *response.ErrorResponse) {
	return handleErrorRepository[*response.DashboardCard](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindTotalBalances, fields...)
}

func (c *cardDasboardError) HandleTotalTopupAmountError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.DashboardCard, *response.ErrorResponse) {
	return handleErrorRepository[*response.DashboardCard](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindTotalTopAmount, fields...)
}

func (c *cardDasboardError) HandleTotalWithdrawAmountError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.DashboardCard, *response.ErrorResponse) {
	return handleErrorRepository[*response.DashboardCard](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindTotalWithdrawAmount, fields...)
}

func (c *cardDasboardError) HandleTotalTransactionAmountError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.DashboardCard, *response.ErrorResponse) {
	return handleErrorRepository[*response.DashboardCard](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindTotalTransactionAmount, fields...)
}

func (c *cardDasboardError) HandleTotalTransferAmountError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.DashboardCard, *response.ErrorResponse) {
	return handleErrorRepository[*response.DashboardCard](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindTotalTransferAmount, fields...)
}

func (c *cardDasboardError) HandleTotalBalanceCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.DashboardCardCardNumber, *response.ErrorResponse) {
	return handleErrorRepository[*response.DashboardCardCardNumber](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindTotalBalanceByCard, fields...)
}

func (c *cardDasboardError) HandleTotalTopupAmountCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.DashboardCardCardNumber, *response.ErrorResponse) {
	return handleErrorRepository[*response.DashboardCardCardNumber](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindTotalTopupAmountByCard, fields...)
}

func (c *cardDasboardError) HandleTotalWithdrawAmountCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.DashboardCardCardNumber, *response.ErrorResponse) {
	return handleErrorRepository[*response.DashboardCardCardNumber](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindTotalWithdrawAmountByCard, fields...)
}

func (c *cardDasboardError) HandleTotalTransactionAmountCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.DashboardCardCardNumber, *response.ErrorResponse) {
	return handleErrorRepository[*response.DashboardCardCardNumber](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindTotalTransactionAmountByCard, fields...)
}

func (c *cardDasboardError) HandleTotalTransferAmountBySender(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.DashboardCardCardNumber, *response.ErrorResponse) {
	return handleErrorRepository[*response.DashboardCardCardNumber](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindTotalTransferAmountBySender, fields...)
}

func (c *cardDasboardError) HandleTotalTransferAmountByReceiver(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.DashboardCardCardNumber, *response.ErrorResponse) {
	return handleErrorRepository[*response.DashboardCardCardNumber](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindTotalTransferAmountByReceiver, fields...)
}
