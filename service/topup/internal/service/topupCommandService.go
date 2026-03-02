package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/email"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/errorhandler"
	topup_errors "github.com/MamangRust/monolith-payment-gateway-shared/errors/topup_errors/service"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"

	mencache "github.com/MamangRust/monolith-payment-gateway-topup/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-topup/internal/repository"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// topupCommandDeps groups dependencies for top-up command service.
type topupCommandDeps struct {
	Kafka                  *kafka.Kafka
	Cache                  mencache.TopupCommandCache
	CardRepository         repository.CardRepository
	TopupQueryRepository   repository.TopupQueryRepository
	TopupCommandRepository repository.TopupCommandRepository
	SaldoRepository        repository.SaldoRepository
	Logger                 logger.LoggerInterface
	Observability          observability.TraceLoggerObservability
}

// topupCommandService handles top-up command operations.
type topupCommandService struct {
	kafka                  *kafka.Kafka
	cache                  mencache.TopupCommandCache
	topupQueryRepository   repository.TopupQueryRepository
	cardRepository         repository.CardRepository
	topupCommandRepository repository.TopupCommandRepository
	saldoRepository        repository.SaldoRepository
	logger                 logger.LoggerInterface
	observability          observability.TraceLoggerObservability
}

func NewTopupCommandService(
	params *topupCommandDeps,
) TopupCommandService {
	return &topupCommandService{
		kafka:                  params.Kafka,
		cache:                  params.Cache,
		topupQueryRepository:   params.TopupQueryRepository,
		topupCommandRepository: params.TopupCommandRepository,
		saldoRepository:        params.SaldoRepository,
		cardRepository:         params.CardRepository,
		logger:                 params.Logger,
		observability:          params.Observability,
	}
}

func (s *topupCommandService) CreateTopup(ctx context.Context, request *requests.CreateTopupRequest) (*db.UpdateTopupStatusRow, error) {
	const method = "CreateTopup"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() { end(status) }()

	card, err := s.cardRepository.FindUserCardByCardNumber(ctx, request.CardNumber)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTopupStatusRow](s.logger, err, method, span, zap.String("card_number", request.CardNumber))
	}

	topup, err := s.topupCommandRepository.CreateTopup(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTopupStatusRow](s.logger, err, method, span, zap.String("card_number", request.CardNumber))
	}

	saldo, err := s.saldoRepository.FindByCardNumber(ctx, request.CardNumber)
	if err != nil {
		status = "error"
		s.markTopupAsFailed(ctx, int(topup.TopupID), method, span)
		return errorhandler.HandleError[*db.UpdateTopupStatusRow](s.logger, err, method, span, zap.String("card_number", request.CardNumber))
	}

	newBalance := int(saldo.TotalBalance) + request.TopupAmount
	_, err = s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
		CardNumber:   request.CardNumber,
		TotalBalance: newBalance,
	})
	if err != nil {
		status = "error"
		s.markTopupAsFailed(ctx, int(topup.TopupID), method, span)
		return errorhandler.HandleError[*db.UpdateTopupStatusRow](s.logger, err, method, span, zap.String("card_number", request.CardNumber))
	}

	expireDate, err := time.Parse("2006-01-02", card.ExpireDate.Time.String())
	if err != nil {
		status = "error"
		s.markTopupAsFailed(ctx, int(topup.TopupID), method, span)
		return errorhandler.HandleError[*db.UpdateTopupStatusRow](s.logger, err, method, span, zap.String("card_number", request.CardNumber))
	}

	_, err = s.cardRepository.UpdateCard(ctx, &requests.UpdateCardRequest{
		CardID:       int(card.CardID),
		UserID:       int(card.UserID),
		CardType:     card.CardType,
		ExpireDate:   expireDate,
		CVV:          card.Cvv,
		CardProvider: card.CardProvider,
	})
	if err != nil {
		status = "error"
		s.markTopupAsFailed(ctx, int(topup.TopupID), method, span)
		return errorhandler.HandleError[*db.UpdateTopupStatusRow](s.logger, err, method, span, zap.String("card_number", request.CardNumber))
	}

	updatedTopup, err := s.topupCommandRepository.UpdateTopupStatus(ctx, &requests.UpdateTopupStatus{
		TopupID: int(topup.TopupID),
		Status:  "success",
	})
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTopupStatusRow](s.logger, err, method, span, zap.Int("topup_id", int(topup.TopupID)))
	}

	go func() {
		htmlBody := email.GenerateEmailHTML(map[string]string{
			"Title":   "Topup Successful",
			"Message": fmt.Sprintf("Your topup of %d has been processed successfully.", request.TopupAmount),
			"Button":  "View History",
			"Link":    "https://sanedge.example.com/topup/history",
		})

		emailPayload := map[string]any{
			"email":   card.Email,
			"subject": "Topup Successful - SanEdge",
			"body":    htmlBody,
		}

		payloadBytes, err := json.Marshal(emailPayload)
		if err != nil {
			s.logger.Error("failed to marshal email payload for topup", zap.Error(err), zap.Int("topup_id", int(topup.TopupID)))
			return
		}

		err = s.kafka.SendMessage("email-service-topic-topup-create", strconv.Itoa(int(updatedTopup.TopupID)), payloadBytes)
		if err != nil {
			s.logger.Error("failed to send topup email via kafka", zap.Error(err), zap.Int("topup_id", int(topup.TopupID)))
		}
	}()

	logSuccess("Topup created successfully", zap.String("cardNumber", request.CardNumber), zap.Int("topupID", int(topup.TopupID)), zap.Float64("topupAmount", float64(request.TopupAmount)))

	return updatedTopup, nil
}

func (s *topupCommandService) UpdateTopup(ctx context.Context, request *requests.UpdateTopupRequest) (*db.UpdateTopupStatusRow, error) {
	const method = "UpdateTopup"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() { end(status) }()

	_, err := s.cardRepository.FindCardByCardNumber(ctx, request.CardNumber)
	if err != nil {
		status = "error"
		s.markTopupAsFailed(ctx, *request.TopupID, method, span)
		return errorhandler.HandleError[*db.UpdateTopupStatusRow](s.logger, err, method, span, zap.String("card_number", request.CardNumber))
	}

	existingTopup, err := s.topupQueryRepository.FindById(ctx, *request.TopupID)
	if err != nil {
		status = "error"
		s.markTopupAsFailed(ctx, *request.TopupID, method, span)
		return errorhandler.HandleError[*db.UpdateTopupStatusRow](s.logger, err, method, span, zap.Int("topup_id", *request.TopupID))
	}

	_, err = s.topupCommandRepository.UpdateTopup(ctx, request)
	if err != nil {
		status = "error"
		s.markTopupAsFailed(ctx, *request.TopupID, method, span)
		return errorhandler.HandleError[*db.UpdateTopupStatusRow](s.logger, err, method, span, zap.Int("topup_id", *request.TopupID))
	}

	currentSaldo, err := s.saldoRepository.FindByCardNumber(ctx, request.CardNumber)
	if err != nil {
		status = "error"
		s.markTopupAsFailed(ctx, *request.TopupID, method, span)
		return errorhandler.HandleError[*db.UpdateTopupStatusRow](s.logger, err, method, span, zap.String("card_number", request.CardNumber))
	}

	topupDifference := request.TopupAmount - int(existingTopup.TopupAmount)

	newBalance := int(currentSaldo.TotalBalance) + topupDifference
	_, err = s.saldoRepository.UpdateSaldoBalance(ctx, &requests.UpdateSaldoBalance{
		CardNumber:   request.CardNumber,
		TotalBalance: newBalance,
	})
	if err != nil {
		status = "error"
		// 6. Jalankan logika rollback
		_, rollbackErr := s.topupCommandRepository.UpdateTopupAmount(ctx, &requests.UpdateTopupAmount{
			TopupID:     *request.TopupID,
			TopupAmount: int(existingTopup.TopupAmount),
		})
		if rollbackErr != nil {
			return errorhandler.HandleError[*db.UpdateTopupStatusRow](s.logger, rollbackErr, method, span, zap.Int("topup_id", *request.TopupID))
		}
		s.markTopupAsFailed(ctx, *request.TopupID, method, span)
		return errorhandler.HandleError[*db.UpdateTopupStatusRow](s.logger, err, method, span, zap.String("card_number", request.CardNumber))
	}

	updatedTopup, err := s.topupCommandRepository.UpdateTopupStatus(ctx, &requests.UpdateTopupStatus{
		TopupID: *request.TopupID,
		Status:  "success",
	})
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTopupStatusRow](s.logger, err, method, span, zap.Int("topup_id", *request.TopupID))
	}

	logSuccess("UpdateTopup process completed", zap.Int("topup_id", *request.TopupID))

	return updatedTopup, nil
}

func (s *topupCommandService) TrashedTopup(ctx context.Context, topup_id int) (*db.Topup, error) {
	const method = "TrashedTopup"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("topup_id", topup_id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Starting TrashedTopup process", zap.Int("topup_id", topup_id))

	res, err := s.topupCommandRepository.TrashedTopup(ctx, topup_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Topup](
			s.logger,
			topup_errors.ErrFailedTrashTopup,
			method,
			span,

			zap.Int("topup_id", topup_id),
		)
	}

	logSuccess("TrashedTopup process completed", zap.Int("topup_id", topup_id))

	return res, nil
}

func (s *topupCommandService) RestoreTopup(ctx context.Context, topup_id int) (*db.Topup, error) {
	const method = "RestoreTopup"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("topup_id", topup_id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Starting RestoreTopup process", zap.Int("topup_id", topup_id))

	res, err := s.topupCommandRepository.RestoreTopup(ctx, topup_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Topup](
			s.logger,
			topup_errors.ErrFailedRestoreTopup,
			method,
			span,

			zap.Int("topup_id", topup_id),
		)
	}

	logSuccess("RestoreTopup process completed", zap.Int("topup_id", topup_id))

	return res, nil
}

func (s *topupCommandService) DeleteTopupPermanent(ctx context.Context, topup_id int) (bool, error) {
	const method = "DeleteTopupPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("topup_id", topup_id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Starting DeleteTopupPermanent process", zap.Int("topup_id", topup_id))

	_, err := s.topupCommandRepository.DeleteTopupPermanent(ctx, topup_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			topup_errors.ErrFailedDeleteTopup,
			method,
			span,

			zap.Int("topup_id", topup_id),
		)
	}

	logSuccess("DeleteTopupPermanent process completed", zap.Int("topup_id", topup_id))

	return true, nil
}

func (s *topupCommandService) RestoreAllTopup(ctx context.Context) (bool, error) {
	const method = "RestoreAllTopup"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	s.logger.Debug("Restoring all topups")

	_, err := s.topupCommandRepository.RestoreAllTopup(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			topup_errors.ErrFailedRestoreAllTopups,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all topups")
	return true, nil
}

func (s *topupCommandService) DeleteAllTopupPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllTopupPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	s.logger.Debug("Permanently deleting all topups")

	_, err := s.topupCommandRepository.DeleteAllTopupPermanent(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			topup_errors.ErrFailedDeleteAllTopups,
			method,
			span,
		)
	}

	logSuccess("Successfully deleted all topups permanently")
	return true, nil
}

func (s *topupCommandService) markTopupAsFailed(ctx context.Context, topupID int, method string, span trace.Span) {
	req := requests.UpdateTopupStatus{
		TopupID: topupID,
		Status:  "failed",
	}
	go func() {
		if _, err := s.topupCommandRepository.UpdateTopupStatus(ctx, &req); err != nil {
			s.logger.Error("compensation: failed to mark topup as failed", zap.Error(err), zap.Int("topup_id", topupID), zap.String("method", method))
		}
	}()
}
