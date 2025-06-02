package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/hash"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/role_errors"
	"github.com/MamangRust/monolith-payment-gateway-shared/errors/user_errors"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service"
	"github.com/MamangRust/monolith-payment-gateway-user/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-user/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-user/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type userCommandService struct {
	ctx                   context.Context
	errorhandler          errorhandler.UserCommandError
	mencache              mencache.UserCommandCache
	trace                 trace.Tracer
	userQueryRepository   repository.UserQueryRepository
	userCommandRepository repository.UserCommandRepository
	roleRepository        repository.RoleRepository
	logger                logger.LoggerInterface
	mapping               responseservice.UserResponseMapper
	hashing               hash.HashPassword
	requestCounter        *prometheus.CounterVec
	requestDuration       *prometheus.HistogramVec
}

func NewUserCommandService(
	ctx context.Context,
	errorhandler errorhandler.UserCommandError,
	mencache mencache.UserCommandCache,
	userQueryRepository repository.UserQueryRepository,
	userCommandRepository repository.UserCommandRepository,
	roleRepository repository.RoleRepository,
	logger logger.LoggerInterface,
	mapper responseservice.UserResponseMapper,
	hashing hash.HashPassword,
) *userCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_command_service_requests_total",
			Help: "Total number of requests to the UserCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "user_command_service_request_duration_seconds",
			Help:    "Histogram of request durations for the UserCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &userCommandService{
		ctx:                   ctx,
		mencache:              mencache,
		errorhandler:          errorhandler,
		trace:                 otel.Tracer("user-command-service"),
		userQueryRepository:   userQueryRepository,
		userCommandRepository: userCommandRepository,
		roleRepository:        roleRepository,
		logger:                logger,
		mapping:               mapper,
		hashing:               hashing,
		requestCounter:        requestCounter,
		requestDuration:       requestDuration,
	}
}

func (s *userCommandService) CreateUser(request *requests.CreateUserRequest) (*response.UserResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("CreateUser", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "CreateUser")
	defer span.End()

	span.SetAttributes(
		attribute.String("email", request.Email),
	)

	s.logger.Debug("Creating new user", zap.String("email", request.Email), zap.Any("request", request))

	existingUser, err := s.userQueryRepository.FindByEmail(request.Email)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "CreateUser", "FAILED_FIND_USER_BY_EMAIL", span, &status, user_errors.ErrUserEmailAlready)

	} else if existingUser != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "CreateUser", "FAILED_FIND_USER_BY_EMAIL", span, &status, user_errors.ErrUserEmailAlready)
	}

	hash, err := s.hashing.HashPassword(request.Password)
	if err != nil {
		return errorhandler.HandleErrorPasswordOperation[*response.UserResponse](s.logger, err, "CreateUser", "FAILED_HASH_PASSWORD", "hash", span, &status, user_errors.ErrUserPassword, zap.String("email", request.Email))
	}

	request.Password = hash

	const defaultRoleName = "Admin Access 1"

	role, err := s.roleRepository.FindByName(defaultRoleName)

	if err != nil || role == nil {
		return s.errorhandler.HandleRepositorySingleError(err, "CreateUser", "FAILED_FIND_ROLE", span, &status, role_errors.ErrRoleNotFoundRes)
	}

	res, err := s.userCommandRepository.CreateUser(request)

	if err != nil {
		return s.errorhandler.HandleCreateUserError(err, "CreateUser", "FAILED_CREATE_USER", span, &status, zap.String("email", request.Email))
	}

	so := s.mapping.ToUserResponse(res)

	s.logger.Debug("Successfully created new user", zap.String("email", so.Email), zap.Int("user", so.ID))

	return so, nil
}

func (s *userCommandService) UpdateUser(request *requests.UpdateUserRequest) (*response.UserResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("UpdateUser", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateUser")
	defer span.End()

	span.SetAttributes(
		attribute.Int("user_id", *request.UserID),
	)

	s.logger.Debug("Updating user", zap.Int("user_id", *request.UserID), zap.Any("request", request))

	existingUser, err := s.userQueryRepository.FindById(*request.UserID)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, "UpdateUser", "FAILED_FIND_USER", span, &status, user_errors.ErrUserNotFoundRes)
	}

	if request.Email != "" && request.Email != existingUser.Email {
		duplicateUser, _ := s.userQueryRepository.FindByEmail(request.Email)

		if duplicateUser != nil {
			return s.errorhandler.HandleRepositorySingleError(err, "UpdateUser", "FAILED_EMAIL_ALREADY", span, &status, user_errors.ErrUserEmailAlready)
		}

		existingUser.Email = request.Email
	}

	if request.Password != "" {
		hash, err := s.hashing.HashPassword(request.Password)
		if err != nil {
			return errorhandler.HandleErrorPasswordOperation[*response.UserResponse](s.logger, err, "UpdateUser", "FAILED_HASH_PASSWORD", "hash", span, &status, user_errors.ErrUserPassword, zap.Int("user_id", *request.UserID))
		}
		existingUser.Password = hash
	}

	const defaultRoleName = "Admin Access 1"

	role, err := s.roleRepository.FindByName(defaultRoleName)

	if err != nil || role == nil {
		return s.errorhandler.HandleRepositorySingleError(err, "UpdateUser", "FAILED_FIND_ROLE", span, &status, role_errors.ErrRoleNotFoundRes)
	}

	res, err := s.userCommandRepository.UpdateUser(request)

	if err != nil {
		return s.errorhandler.HandleUpdateUserError(err, "UpdateUser", "FAILED_UPDATE_USER", span, &status, zap.Int("user_id", *request.UserID))
	}

	so := s.mapping.ToUserResponse(res)

	s.mencache.DeleteUserCache(so.ID)

	s.logger.Debug("Successfully updated user", zap.Int("user_id", so.ID))

	return so, nil
}

func (s *userCommandService) TrashedUser(user_id int) (*response.UserResponseDeleteAt, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("TrashedUser", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "TrashedUser")
	defer span.End()

	span.SetAttributes(
		attribute.Int("user_id", user_id),
	)

	s.logger.Debug("Trashing user", zap.Int("user_id", user_id))

	res, err := s.userCommandRepository.TrashedUser(user_id)

	if err != nil {
		return s.errorhandler.HandleTrashedUserError(err, "TrashedUser", "FAILED_TO_TRASH_USER", span, &status, zap.Int("user_id", user_id))
	}

	so := s.mapping.ToUserResponseDeleteAt(res)

	s.mencache.DeleteUserCache(so.ID)

	s.logger.Debug("Successfully trashed user", zap.Int("user_id", user_id))

	return so, nil
}

func (s *userCommandService) RestoreUser(user_id int) (*response.UserResponseDeleteAt, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreUser", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreUser")
	defer span.End()

	span.SetAttributes(
		attribute.Int("user_id", user_id),
	)

	s.logger.Debug("Restoring user", zap.Int("user_id", user_id))

	res, err := s.userCommandRepository.RestoreUser(user_id)

	if err != nil {
		return s.errorhandler.HandleRestoreUserError(err, "RestoreUser", "FAILED_TO_RESTORE_USER", span, &status, zap.Int("user_id", user_id))
	}

	so := s.mapping.ToUserResponseDeleteAt(res)

	s.logger.Debug("Successfully restored user", zap.Int("user_id", user_id))

	return so, nil
}

func (s *userCommandService) DeleteUserPermanent(user_id int) (bool, *response.ErrorResponse) {
	start := time.Now()

	status := "success"

	defer func() {
		s.recordMetrics("DeleteUserPermanent", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteUserPermanent")
	defer span.End()

	span.SetAttributes(
		attribute.Int("user_id", user_id),
	)

	s.logger.Debug("Deleting user permanently", zap.Int("user_id", user_id))

	_, err := s.userCommandRepository.DeleteUserPermanent(user_id)

	if err != nil {
		return s.errorhandler.HandleDeleteUserError(err, "DeleteUserPermanent", "FAILED_TO_DELETE_USER", span, &status, zap.Int("user_id", user_id))
	}

	s.logger.Debug("Successfully deleted user permanently", zap.Int("user_id", user_id))

	return true, nil
}

func (s *userCommandService) RestoreAllUser() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreAllUser", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreAllUser")
	defer span.End()

	s.logger.Debug("Restoring all users")

	_, err := s.userCommandRepository.RestoreAllUser()

	if err != nil {
		return s.errorhandler.HandleRestoreAllUserError(err, "RestoreAllUser", "FAILED_RESTORE_ALL_USER", span, &status)
	}

	s.logger.Debug("Successfully restored all users")

	return true, nil
}

func (s *userCommandService) DeleteAllUserPermanent() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteAllUserPermanent", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteAllUserPermanent")
	defer span.End()

	s.logger.Debug("Permanently deleting all users")

	_, err := s.userCommandRepository.DeleteAllUserPermanent()

	if err != nil {
		return s.errorhandler.HandleDeleteAllUserError(err, "DeleteAllUserPermanent", "FAILED_DELETE_ALL_USER", span, &status)
	}

	s.logger.Debug("Successfully deleted all users permanently")

	return true, nil
}

func (s *userCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
