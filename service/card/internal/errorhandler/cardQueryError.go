package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type cardQueryError struct {
	logger logger.LoggerInterface
}

func NewCardQueryError(logger logger.LoggerInterface) *cardQueryError {
	return &cardQueryError{
		logger: logger,
	}
}

func (c *cardQueryError) HandleFindAllError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponse, *int, *response.ErrorResponse) {
	return handleErrorPaginationTemplate[[]*response.CardResponse](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindAllCards, fields...)
}

func (c *cardQueryError) HandleFindByActiveError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPaginationTemplate[[]*response.CardResponseDeleteAt](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindActiveCards, fields...)
}

func (c *cardQueryError) HandleFindByTrashedError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CardResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPaginationTemplate[[]*response.CardResponseDeleteAt](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindTrashedCards, fields...)
}

func (c *cardQueryError) HandleFindByIdError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.CardResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.CardResponse](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindById, fields...)
}

func (c *cardQueryError) HandleFindByUserIdError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.CardResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.CardResponse](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindByUserID, fields...)
}

func (c *cardQueryError) HandleFindByCardNumberError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.CardResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.CardResponse](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedFindByCardNumber, fields...)
}
