package service

import (
	"context"
	"encoding/json"
	"strconv"

	mencache "github.com/MamangRust/monolith-payment-gateway-card/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-card/internal/repository"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-payment-gateway-shared/errorhandler"
	card_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/service"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// cardCommandServiceDeps defines dependencies for cardCommandService.
type cardCommandServiceDeps struct {
	Cache                 mencache.CardCommandCache
	Kafka                 *kafka.Kafka
	UserRepository        repository.UserRepository
	CardCommandRepository repository.CardCommandRepository
	Logger                logger.LoggerInterface
	Observability         observability.TraceLoggerObservability
}

// cardCommandService implements CardCommandService.
type cardCommandService struct {
	cache                 mencache.CardCommandCache
	kafka                 *kafka.Kafka
	userRepository        repository.UserRepository
	cardCommandRepository repository.CardCommandRepository
	logger                logger.LoggerInterface
	observability         observability.TraceLoggerObservability
}

func NewCardCommandService(params *cardCommandServiceDeps) CardCommandService {

	return &cardCommandService{
		cache:                 params.Cache,
		kafka:                 params.Kafka,
		userRepository:        params.UserRepository,
		cardCommandRepository: params.CardCommandRepository,
		logger:                params.Logger,
		observability:         params.Observability,
	}
}

func (s *cardCommandService) CreateCard(ctx context.Context, request *requests.CreateCardRequest) (*db.CreateCardRow, error) {
	const method = "CreateCard"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.userRepository.FindById(ctx, request.UserID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.CreateCardRow](s.logger, err, method, span, zap.Int("user_id", request.UserID))
	}

	res, err := s.cardCommandRepository.CreateCard(ctx, request)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.CreateCardRow](s.logger, err, method, span, zap.Int("user_id", request.UserID))
	}

	go func() {
		saldoPayload := map[string]any{
			"card_number":   res.CardNumber,
			"total_balance": 0,
		}

		payloadBytes, err := json.Marshal(saldoPayload)
		if err != nil {
			s.logger.Error("failed to marshal saldo payload for new card", zap.Error(err), zap.Int("card_id", int(res.CardID)))
			return
		}

		err = s.kafka.SendMessage("saldo-service-topic-create-saldo", strconv.Itoa(int(res.CardID)), payloadBytes)
		if err != nil {
			s.logger.Error("failed to send create saldo message to kafka", zap.Error(err), zap.Int("card_id", int(res.CardID)))
		}
	}()

	logSuccess("Successfully created card", zap.Int("card.id", int(res.CardID)), zap.String("card.card_number", res.CardNumber))

	return res, nil
}

func (s *cardCommandService) UpdateCard(ctx context.Context, request *requests.UpdateCardRequest) (*db.UpdateCardRow, error) {
	const method = "UpdateCard"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.userRepository.FindById(ctx, request.UserID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.UpdateCardRow](s.logger, err, method, span, zap.Int("user_id", request.UserID))
	}

	res, err := s.cardCommandRepository.UpdateCard(ctx, request)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.UpdateCardRow](s.logger, err, method, span, zap.Int("card_id", request.CardID))
	}

	s.cache.DeleteCardCommandCache(ctx, request.CardID)

	logSuccess("Successfully updated card", zap.Int("card.id", int(res.CardID)))

	return res, nil
}

func (s *cardCommandService) TrashedCard(ctx context.Context, card_id int) (*db.Card, error) {
	const method = "TrashedCard"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("card_id", card_id))

	defer func() {
		end(status)
	}()

	res, err := s.cardCommandRepository.TrashedCard(ctx, card_id)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Card](
			s.logger,
			card_errors.ErrFailedTrashCard,
			method,
			span,

			zap.Int("card_id", card_id),
		)
	}

	s.cache.DeleteCardCommandCache(ctx, card_id)

	logSuccess("Successfully trashed card", zap.Int("card_id", int(res.CardID)))

	return res, nil
}

func (s *cardCommandService) RestoreCard(ctx context.Context, card_id int) (*db.Card, error) {
	const method = "RestoreCard"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("card_id", card_id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Restoring card", zap.Int("card_id", card_id))

	res, err := s.cardCommandRepository.RestoreCard(ctx, card_id)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Card](
			s.logger,
			card_errors.ErrFailedRestoreCard,
			method,
			span,

			zap.Int("card_id", card_id),
		)
	}

	s.cache.DeleteCardCommandCache(ctx, card_id)

	logSuccess("Successfully restored card", zap.Int("card_id", int(res.CardID)))

	return res, nil
}

func (s *cardCommandService) DeleteCardPermanent(ctx context.Context, card_id int) (bool, error) {
	const method = "DeleteCardPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("card_id", card_id))

	defer func() {
		end(status)
	}()

	_, err := s.cardCommandRepository.DeleteCardPermanent(ctx, card_id)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			card_errors.ErrFailedDeleteCard,
			method,
			span,

			zap.Int("card_id", card_id),
		)
	}

	s.cache.DeleteCardCommandCache(ctx, card_id)

	logSuccess("Successfully deleted card permanently", zap.Int("card_id", card_id))

	return true, nil
}

func (s *cardCommandService) RestoreAllCard(ctx context.Context) (bool, error) {
	const method = "RestoreAllCard"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.cardCommandRepository.RestoreAllCard(ctx)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			card_errors.ErrFailedRestoreAllCards,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all cards")

	return true, nil
}

func (s *cardCommandService) DeleteAllCardPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllCardPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.cardCommandRepository.DeleteAllCardPermanent(ctx)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			card_errors.ErrFailedDeleteAllCards,
			method,
			span,
		)
	}

	logSuccess("Successfully deleted all cards permanently")

	return true, nil
}
