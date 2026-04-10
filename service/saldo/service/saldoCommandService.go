package service

import (
	"context"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	mencache "github.com/MamangRust/monolith-payment-gateway-saldo/redis"
	"github.com/MamangRust/monolith-payment-gateway-saldo/repository"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errorhandler"
	saldo_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/saldo_errors/service"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type saldoCommandParams struct {
	Cache                  mencache.SaldoCommandCache
	saldoCommandRepository repository.SaldoCommandRepository
	CardRepository         repository.CardRepository
	Logger                 logger.LoggerInterface
	Observability          observability.TraceLoggerObservability
}

type saldoCommandService struct {
	ctx                    context.Context
	cache                  mencache.SaldoCommandCache
	cardRepository         repository.CardRepository
	logger                 logger.LoggerInterface
	saldoCommandRepository repository.SaldoCommandRepository
	observability          observability.TraceLoggerObservability
}

func NewSaldoCommandService(params *saldoCommandParams) SaldoCommandService {
	return &saldoCommandService{
		cache:                  params.Cache,
		saldoCommandRepository: params.saldoCommandRepository,
		cardRepository:         params.CardRepository,
		logger:                 params.Logger,
		observability:          params.Observability,
	}
}

func (s *saldoCommandService) CreateSaldo(ctx context.Context, request *requests.CreateSaldoRequest) (*db.CreateSaldoRow, error) {
	const method = "CreateSaldo"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() { end(status) }()

	_, err := s.cardRepository.FindCardByCardNumber(ctx, request.CardNumber)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateSaldoRow](s.logger, err, method, span, zap.String("card_number", request.CardNumber))
	}

	res, err := s.saldoCommandRepository.CreateSaldo(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateSaldoRow](s.logger, err, method, span, zap.String("card_number", request.CardNumber))
	}

	logSuccess("Successfully created saldo record", zap.String("card_number", request.CardNumber), zap.Float64("amount", float64(request.TotalBalance)))

	return res, nil
}

func (s *saldoCommandService) UpdateSaldo(ctx context.Context, request *requests.UpdateSaldoRequest) (*db.UpdateSaldoRow, error) {
	const method = "UpdateSaldo"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() { end(status) }()

	_, err := s.cardRepository.FindCardByCardNumber(ctx, request.CardNumber)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateSaldoRow](s.logger, err, method, span, zap.String("card_number", request.CardNumber))
	}

	res, err := s.saldoCommandRepository.UpdateSaldo(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateSaldoRow](s.logger, err, method, span, zap.String("card_number", request.CardNumber))
	}

	s.cache.DeleteSaldoCache(ctx, int(res.SaldoID))

	logSuccess("Successfully updated saldo record", zap.String("card_number", request.CardNumber), zap.Float64("amount", float64(request.TotalBalance)))

	return res, nil
}

func (s *saldoCommandService) UpdateSaldoWithdraw(ctx context.Context, request *requests.UpdateSaldoWithdraw) (*db.UpdateSaldoWithdrawRow, error) {
	const method = "UpdateSaldoWithdraw"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() { end(status) }()

	_, err := s.cardRepository.FindCardByCardNumber(ctx, request.CardNumber)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateSaldoWithdrawRow](s.logger, err, method, span, zap.String("card_number", request.CardNumber))
	}

	res, err := s.saldoCommandRepository.UpdateSaldoWithdraw(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateSaldoWithdrawRow](s.logger, err, method, span, zap.String("card_number", request.CardNumber))
	}

	s.cache.DeleteSaldoCache(ctx, int(res.SaldoID))

	logSuccess("Successfully updated saldo withdraw record", zap.String("card_number", request.CardNumber))

	return res, nil
}

func (s *saldoCommandService) TrashSaldo(ctx context.Context, saldo_id int) (*db.Saldo, error) {
	const method = "TrashSaldo"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("saldo_id", saldo_id))

	defer func() {
		end(status)
	}()

	res, err := s.saldoCommandRepository.TrashedSaldo(ctx, saldo_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Saldo](
			s.logger,
			saldo_errors.ErrFailedTrashSaldo,
			method,
			span,

			zap.Int("saldo_id", saldo_id),
		)
	}

	logSuccess("Successfully trashed saldo", zap.Int("saldo_id", saldo_id))

	return res, nil
}

func (s *saldoCommandService) RestoreSaldo(ctx context.Context, saldo_id int) (*db.Saldo, error) {
	const method = "RestoreSaldo"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("saldo_id", saldo_id))

	defer func() {
		end(status)
	}()

	res, err := s.saldoCommandRepository.RestoreSaldo(ctx, saldo_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Saldo](
			s.logger,
			saldo_errors.ErrFailedRestoreSaldo,
			method,
			span,

			zap.Int("saldo_id", saldo_id),
		)
	}

	logSuccess("Successfully restored saldo", zap.Int("saldo_id", saldo_id))

	return res, nil
}

func (s *saldoCommandService) DeleteSaldoPermanent(ctx context.Context, saldo_id int) (bool, error) {
	const method = "DeleteSaldoPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("saldo_id", saldo_id))

	defer func() {
		end(status)
	}()

	_, err := s.saldoCommandRepository.DeleteSaldoPermanent(ctx, saldo_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			saldo_errors.ErrFailedDeleteSaldoPermanent,
			method,
			span,

			zap.Int("saldo_id", saldo_id),
		)
	}

	logSuccess("Successfully deleted saldo permanently", zap.Int("saldo_id", saldo_id))

	return true, nil
}

func (s *saldoCommandService) RestoreAllSaldo(ctx context.Context) (bool, error) {
	const method = "RestoreAllSaldo"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.saldoCommandRepository.RestoreAllSaldo(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			saldo_errors.ErrFailedRestoreAllSaldo,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all saldo")
	return true, nil
}

func (s *saldoCommandService) DeleteAllSaldoPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllSaldoPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.saldoCommandRepository.DeleteAllSaldoPermanent(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			saldo_errors.ErrFailedDeleteAllSaldoPermanent,
			method,
			span,
		)
	}

	logSuccess("Successfully deleted all saldo permanently")
	return true, nil
}
