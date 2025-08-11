package service

import (
	"context"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-saldo/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-saldo/internal/repository"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/saldo"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

// saldoCommandParams holds the dependencies required to construct a saldoCommandService.
type saldoCommandParams struct {
	// Ctx is the base context for the service.
	Ctx context.Context

	// ErrorHandler handles domain-specific errors for saldo command operations.
	ErrorHandler errorhandler.SaldoCommandErrorHandler

	// Cache provides in-memory caching for saldo command operations.
	Cache mencache.SaldoCommandCache

	// SaldoRepository provides access to persistent storage for saldo commands.
	SaldoRepository repository.SaldoCommandRepository

	// CardRepository provides access to card data related to saldo operations.
	CardRepository repository.CardRepository

	// Logger is used for structured logging.
	Logger logger.LoggerInterface

	// Mapper maps internal saldo entities to response DTOs.
	Mapper responseservice.SaldoCommandResponseMapper
}

// saldoCommandService handles write operations for saldo, such as top-up and adjustment.
type saldoCommandService struct {
	// ctx is the base context shared across the service.
	ctx context.Context

	// errorhandler handles domain-specific errors for saldo commands.
	errorhandler errorhandler.SaldoCommandErrorHandler

	// mencache provides in-memory caching for saldo data.
	mencache mencache.SaldoCommandCache

	// cardRepository accesses card information related to saldo operations.
	cardRepository repository.CardRepository

	// logger provides structured logging capability.
	logger logger.LoggerInterface

	// mapper converts internal saldo models to response formats.
	mapper responseservice.SaldoCommandResponseMapper

	// saldoCommandRepository provides persistence operations for saldo commands.
	saldoCommandRepository repository.SaldoCommandRepository

	observability observability.TraceLoggerObservability
}

// NewSaldoCommandService initializes a new instance of saldoCommandService with the provided parameters.
// It sets up the prometheus metrics for counting and measuring the duration of saldo command requests.
//
// Parameters:
// - params: A pointer to a saldoCommandParams containing the necessary dependencies.
//
// Returns:
// - A pointer to a newly created saldoCommandService.
func NewSaldoCommandService(params *saldoCommandParams) SaldoCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "saldo_command_service_request_total",
			Help: "Total number of requests to the SaldoCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "saldo_command_service_request_duration_seconds",
			Help:    "Histogram of request durations for the SaldoCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("saldo-command-service"), params.Logger, requestCounter, requestDuration)

	return &saldoCommandService{
		ctx:                    params.Ctx,
		errorhandler:           params.ErrorHandler,
		mencache:               params.Cache,
		saldoCommandRepository: params.SaldoRepository,
		cardRepository:         params.CardRepository,
		logger:                 params.Logger,
		mapper:                 params.Mapper,
		observability:          observability,
	}
}

// CreateSaldo creates a new saldo record in the system.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The request payload containing saldo creation data.
//
// Returns:
//   - *response.SaldoResponse: The created saldo response.
//   - *response.ErrorResponse: An error response if creation fails.
func (s *saldoCommandService) CreateSaldo(ctx context.Context, request *requests.CreateSaldoRequest) (*response.SaldoResponse, *response.ErrorResponse) {
	const method = "CreateSaldo"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	rescard, err := s.cardRepository.FindCardByCardNumber(ctx, request.CardNumber)

	if err != nil {
		return s.errorhandler.HandleFindCardByNumberError(err, method, "FAILED_FIND_CARD_BY_CARD_NUMBER", span, &status, zap.Error(err))
	}

	s.logger.Info("response card", zap.Any("card", rescard), zap.Any("err", err))

	res, err := s.saldoCommandRepository.CreateSaldo(ctx, request)

	if err != nil {
		return s.errorhandler.HandleCreateSaldoError(err, method, "FAILED_CREATE_SALDO", span, &status, zap.Error(err))
	}

	so := s.mapper.ToSaldoResponse(res)

	logSuccess("Successfully created saldo record", zap.String("card_number", request.CardNumber), zap.Float64("amount", float64(request.TotalBalance)))

	return so, nil
}

// UpdateSaldo updates an existing saldo record.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - request: The request payload containing updated saldo data.
//
// Returns:
//   - *response.SaldoResponse: The updated saldo response.
//   - *response.ErrorResponse: An error response if update fails.
func (s *saldoCommandService) UpdateSaldo(ctx context.Context, request *requests.UpdateSaldoRequest) (*response.SaldoResponse, *response.ErrorResponse) {
	const method = "UpdateSaldo"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.cardRepository.FindCardByCardNumber(ctx, request.CardNumber)

	if err != nil {
		return s.errorhandler.HandleFindCardByNumberError(err, method, "FAILED_FIND_CARD_BY_CARD_NUMBER", span, &status, zap.Error(err))
	}

	res, err := s.saldoCommandRepository.UpdateSaldo(ctx, request)

	if err != nil {
		return s.errorhandler.HandleUpdateSaldoError(err, "UpdateSaldo", "FAILED_UPDATE_SALDO", span, &status, zap.Error(err))
	}

	so := s.mapper.ToSaldoResponse(res)

	s.mencache.DeleteSaldoCache(ctx, res.ID)

	logSuccess("Successfully updated saldo record", zap.String("card_number", request.CardNumber), zap.Float64("amount", float64(request.TotalBalance)))

	return so, nil
}

// TrashSaldo moves a saldo to the trash (soft delete).
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - saldo_id: The ID of the saldo to trash.
//
// Returns:
//   - *response.SaldoResponse: The trashed saldo response.
//   - *response.ErrorResponse: An error response if trashing fails.
func (s *saldoCommandService) TrashSaldo(ctx context.Context, saldo_id int) (*response.SaldoResponseDeleteAt, *response.ErrorResponse) {
	const method = "TrashSaldo"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.saldoCommandRepository.TrashedSaldo(ctx, saldo_id)

	if err != nil {
		return s.errorhandler.HandleTrashSaldoError(err, method, "FAILED_TRASH_SALDO", span, &status, zap.Error(err))
	}
	so := s.mapper.ToSaldoResponseDeleteAt(res)

	s.mencache.DeleteSaldoCache(ctx, saldo_id)

	logSuccess("Successfully trashed saldo record", zap.Int("saldo.id", saldo_id))

	return so, nil
}

// RestoreSaldo restores a previously trashed saldo.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - saldo_id: The ID of the saldo to restore.
//
// Returns:
//   - *response.SaldoResponse: The restored saldo response.
//   - *response.ErrorResponse: An error response if restoring fails.
func (s *saldoCommandService) RestoreSaldo(ctx context.Context, saldo_id int) (*response.SaldoResponse, *response.ErrorResponse) {
	const method = "RestoreSaldo"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.saldoCommandRepository.RestoreSaldo(ctx, saldo_id)

	if err != nil {
		return s.errorhandler.HandleRestoreSaldoError(err, method, "FAILED_RESTORE_SALDO", span, &status, zap.Error(err))
	}

	so := s.mapper.ToSaldoResponse(res)

	logSuccess("Successfully restored saldo record", zap.Int("saldo.id", saldo_id))

	return so, nil
}

// DeleteSaldoPermanent permanently deletes a saldo record.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//   - saldo_id: The ID of the saldo to delete permanently.
//
// Returns:
//   - bool: True if the deletion is successful.
//   - *response.ErrorResponse: An error response if deletion fails.
func (s *saldoCommandService) DeleteSaldoPermanent(ctx context.Context, saldo_id int) (bool, *response.ErrorResponse) {
	const method = "DeleteSaldoPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.saldoCommandRepository.DeleteSaldoPermanent(ctx, saldo_id)

	if err != nil {
		return s.errorhandler.HandleDeleteSaldoPermanentError(err, method, "FAILED_DELETE_SALDO_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted permanent saldo record", zap.Int("saldo.id", saldo_id))

	return true, nil
}

// RestoreAllSaldo restores all trashed saldo records.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//
// Returns:
//   - bool: True if restoration is successful.
//   - *response.ErrorResponse: An error response if operation fails.
func (s *saldoCommandService) RestoreAllSaldo(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "RestoreAllSaldo"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.saldoCommandRepository.RestoreAllSaldo(ctx)

	if err != nil {
		return s.errorhandler.HandleRestoreAllSaldoError(err, method, "FAILED_RESTORE_ALL_SALDO", span, &status, zap.Error(err))
	}

	logSuccess("Successfully restored all saldo", zap.Bool("success", true))

	return true, nil
}

// DeleteAllSaldoPermanent permanently deletes all trashed saldo records.
//
// Parameters:
//   - ctx: The context for timeout and cancellation.
//
// Returns:
//   - bool: True if deletion is successful.
//   - *response.ErrorResponse: An error response if operation fails.
func (s *saldoCommandService) DeleteAllSaldoPermanent(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "DeleteAllSaldoPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	_, err := s.saldoCommandRepository.DeleteAllSaldoPermanent(ctx)

	if err != nil {
		return s.errorhandler.HandleDeleteAllSaldoPermanentError(err, method, "FAILED_DELETE_ALL_SALDO_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted all permanent saldo", zap.Bool("success", true))

	return true, nil
}
