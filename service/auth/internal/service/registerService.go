package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-auth/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-payment-gateway-auth/internal/redis"
	"github.com/MamangRust/monolith-payment-gateway-auth/internal/repository"

	"github.com/MamangRust/monolith-payment-gateway-pkg/email"
	"github.com/MamangRust/monolith-payment-gateway-pkg/hash"
	"github.com/MamangRust/monolith-payment-gateway-pkg/kafka"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	randomstring "github.com/MamangRust/monolith-payment-gateway-pkg/random_string"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/requests"
	"github.com/MamangRust/monolith-payment-gateway-shared/domain/response"
	responseservice "github.com/MamangRust/monolith-payment-gateway-shared/mapper/response/service/user"
	"github.com/MamangRust/monolith-payment-gateway-shared/observability"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// RegisterServiceDeps holds dependencies for constructing a registerService.
type RegisterServiceDeps struct {
	// ErrorHandler handles registration-related errors.
	ErrorHandler errorhandler.RegisterErrorHandler

	// ErrorPassword handles password hashing or validation errors.
	ErrorPassword errorhandler.PasswordErrorHandler

	// ErrorRandomString handles random string generation errors.
	ErrorRandomString errorhandler.RandomStringErrorHandler

	// ErrorMarshal handles marshalling (e.g., JSON) errors.
	ErrorMarshal errorhandler.MarshalErrorHandler

	// ErrorKafka handles Kafka publishing errors.
	ErrorKafka errorhandler.KafkaErrorHandler

	// Cache provides caching for registration flow.
	Cache mencache.RegisterCache

	// User is the repository for managing user data.
	User repository.UserRepository

	// Role is the repository for managing role definitions.
	Role repository.RoleRepository

	// UserRole is the repository for assigning roles to users.
	UserRole repository.UserRoleRepository

	// Hash handles password hashing logic.
	Hash hash.HashPassword

	// Kafka is the Kafka client for publishing registration events.
	Kafka *kafka.Kafka

	// Logger logs service operations and errors.
	Logger logger.LoggerInterface

	// Mapping maps user entities to response models.
	Mapper responseservice.UserQueryResponseMapper
}

// registerService handles user registration logic, including role assignment,
// password hashing, caching, and event publishing to Kafka.
type registerService struct {
	// Handles registration-specific errors.
	errohandler errorhandler.RegisterErrorHandler

	// Handles password validation and hashing.
	errorPassword errorhandler.PasswordErrorHandler

	// Handles errors during random string generation.
	errorRandomString errorhandler.RandomStringErrorHandler

	// Handles JSON or struct marshalling errors.
	errorMarshal errorhandler.MarshalErrorHandler

	// Handles Kafka-related errors.
	errorKafka errorhandler.KafkaErrorHandler

	// Caching layer for temporary registration data.
	mencache mencache.RegisterCache

	// User data repository.
	user repository.UserRepository

	// Role repository for looking up roles.
	role repository.RoleRepository

	// Repository for assigning roles to users.
	userRole repository.UserRoleRepository

	// Utility for password hashing.
	hash hash.HashPassword

	// Kafka publisher for registration events.
	kafka *kafka.Kafka

	// Logger for events and diagnostics.
	logger logger.LoggerInterface

	// Mapper to convert user data into response objects.
	mapper responseservice.UserQueryResponseMapper

	observability observability.TraceLoggerObservability
}

// NewRegisterService initializes and returns the RegisterService with the given parameters.
// It sets up the prometheus metrics for request counters and durations, and registers them.
// The function takes RegisterServiceDeps which includes context, error handlers, cache, logger,
// user repository, role repository, user role repository, hash, kafka, and token service.
// Returns a pointer to the initialized registerService.
func NewRegisterService(params *RegisterServiceDeps) *registerService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "register_service_requests_total",
			Help: "Total number of requests to the RegisterService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "register_service_request_duration_seconds",
			Help:    "Histogram of request durations for the RegisterService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	observability := observability.NewTraceLoggerObservability(otel.Tracer("register-service"), params.Logger, requestCounter, requestDuration)

	return &registerService{
		errorPassword:     params.ErrorPassword,
		errohandler:       params.ErrorHandler,
		errorRandomString: params.ErrorRandomString,
		errorMarshal:      params.ErrorMarshal,
		errorKafka:        params.ErrorKafka,
		mencache:          params.Cache,
		user:              params.User,
		role:              params.Role,
		userRole:          params.UserRole,
		hash:              params.Hash,
		kafka:             params.Kafka,
		logger:            params.Logger,
		mapper:            params.Mapper,
		observability:     observability,
	}
}

// Register creates a new user account with the given registration request.
//
// Parameters:
//   - ctx: the context for the operation (e.g., timeout, logging, tracing)
//   - request: the registration request payload
//
// Returns:
//   - A UserResponse if registration is successful, or an ErrorResponse if it fails.
func (s *registerService) Register(ctx context.Context, request *requests.RegisterRequest) (*response.UserResponse, *response.ErrorResponse) {
	const method = "Register"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("email", request.Email))

	defer func() {
		end(status)
	}()

	existingUser, err := s.user.FindByEmail(ctx, request.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			existingUser = nil
		} else {
			return s.errohandler.HandleFindEmailError(err, "Register", "REGISTER_ERR", span, &status,
				zap.String("email", request.Email), zap.Error(err))
		}
	}

	if existingUser != nil {
		err := errors.New("user already exists")
		return s.errohandler.HandleFindEmailError(err, "Register", "REGISTER_ERR", span, &status,
			zap.String("email", request.Email), zap.Error(err))
	}

	passwordHash, err := s.hash.HashPassword(request.Password)
	if err != nil {
		return s.errorPassword.HandleHashPasswordError(err, "Register", "REGISTER_ERR", span, &status, zap.Error(err))
	}
	request.Password = passwordHash

	const defaultRoleName = "Admin_Role_10"

	role, err := s.role.FindByName(ctx, defaultRoleName)

	if err != nil || role == nil {
		return s.errohandler.HandleFindRoleError(err, "Register", "REGISTER_ERR", span, &status,
			zap.String("role_name", defaultRoleName), zap.Error(err))
	}

	random, err := randomstring.GenerateRandomString(10)
	if err != nil {
		return s.errorRandomString.HandleRandomStringErrorRegister(err, "Register", "REGISTER_ERR", span, &status, zap.Error(err))
	}

	request.VerifiedCode = random
	request.IsVerified = false

	newUser, err := s.user.CreateUser(ctx, request)
	if err != nil {
		return s.errohandler.HandleCreateUserError(err, "Register", "REGISTER_ERR", span, &status, zap.Error(err))
	}

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   "Welcome to SanEdge",
		"Message": "Your account has been successfully created.",
		"Button":  "Verify Now",
		"Link":    "https://sanedge.example.com/login?verify_code=" + request.VerifiedCode,
	})

	emailPayload := map[string]any{
		"email":   request.Email,
		"subject": "Welcome to SanEdge",
		"body":    htmlBody,
	}

	payloadBytes, err := json.Marshal(emailPayload)
	if err != nil {
		return s.errorMarshal.HandleMarshalRegisterError(err, "Register", "MARSHAL_ERR", span, &status, zap.Error(err))
	}

	err = s.kafka.SendMessage("email-service-topic-auth-register", strconv.Itoa(newUser.ID), payloadBytes)
	if err != nil {
		return s.errorKafka.HandleSendEmailRegister(err, "Register", "SEND_EMAIL_ERR", span, &status, zap.Error(err))
	}

	_, err = s.userRole.AssignRoleToUser(ctx, &requests.CreateUserRoleRequest{
		UserId: newUser.ID,
		RoleId: role.ID,
	})
	if err != nil {
		return s.errohandler.HandleAssignRoleError(err, "Register", "ASSIGN_ROLE_ERR", span, &status, zap.Error(err))
	}

	s.mencache.SetVerificationCodeCache(ctx, request.Email, random, 15*time.Minute)

	userResponse := s.mapper.ToUserResponse(newUser)

	logSuccess("User registered successfully",
		zap.String("email", request.Email),
		zap.String("first_name", request.FirstName),
		zap.String("last_name", request.LastName),
	)

	return userResponse, nil
}
