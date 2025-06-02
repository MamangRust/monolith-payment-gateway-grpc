package errorhandler

import (
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type cardCommandError struct {
	logger logger.LoggerInterface
}

func NewCardCommandError(logger logger.LoggerInterface) *cardCommandError {
	return &cardCommandError{
		logger: logger,
	}
}

func (c *cardCommandError) HandleFindByIdUserError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.CardResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.CardResponse](c.logger, err, method, tracePrefix, span, status, user_errors.ErrUserNotFoundRes, fields...)
}

func (c *cardCommandError) HandleCreateCardError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.CardResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.CardResponse](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedCreateCard, fields...)
}

func (c *cardCommandError) HandleUpdateCardError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.CardResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.CardResponse](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedUpdateCard, fields...)
}

func (c *cardCommandError) HandleTrashedCardError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.CardResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.CardResponse](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedTrashCard, fields...)
}

func (c *cardCommandError) HandleRestoreCardError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (*response.CardResponse, *response.ErrorResponse) {
	return handleErrorTemplate[*response.CardResponse](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedRestoreCard, fields...)
}

func (c *cardCommandError) HandleDeleteCardPermanentError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorTemplate[bool](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedDeleteCard, fields...)
}

func (c *cardCommandError) HandleRestoreAllCardError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorTemplate[bool](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedRestoreAllCards, fields...)
}

func (c *cardCommandError) HandleDeleteAllCardPermanentError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) (bool, *response.ErrorResponse) {
	return handleErrorTemplate[bool](c.logger, err, method, tracePrefix, span, status, card_errors.ErrFailedDeleteAllCards, fields...)
}
