package utils

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/dotenv"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
)

var schema = `
CREATE TABLE IF NOT EXISTS users (
	user_id SERIAL PRIMARY KEY,
	firstname TEXT,
	lastname TEXT,
	email TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	verification_code TEXT,
	is_verified BOOLEAN DEFAULT false,
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS roles (
	role_id SERIAL PRIMARY KEY,
	role_name TEXT UNIQUE NOT NULL,
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS user_roles (
	user_role_id SERIAL PRIMARY KEY,
	user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
	role_id INT NOT NULL REFERENCES roles(role_id) ON DELETE CASCADE,
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
	refresh_token_id SERIAL PRIMARY KEY,
	user_id INT NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
	token VARCHAR(255) NOT NULL UNIQUE,
	expiration TIMESTAMPTZ NOT NULL,
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE "reset_tokens" (
    "id" SERIAL PRIMARY KEY,
    "user_id" BIGINT NOT NULL UNIQUE REFERENCES users(user_id) ON DELETE CASCADE,
    "token" TEXT NOT NULL UNIQUE,
    "expiry_date" TIMESTAMP NOT NULL
);

`

func SetupTestDB(t *testing.T) (*sql.DB, func()) {
	os.Setenv("APP_ENV", "test")
	setTestEnvVars()

	logger, _ := logger.NewLogger("test")
	_ = dotenv.Viper()

	conn, err := database.NewClient(logger)
	if err != nil {
		t.Fatalf("DB connection failed: %v", err)
	}

	ctx := context.Background()

	schemaName := fmt.Sprintf("schema_%s", strings.ToLower(strings.ReplaceAll(t.Name(), "/", "_")))
	_, _ = conn.ExecContext(ctx, fmt.Sprintf(`DROP SCHEMA IF EXISTS %s CASCADE`, schemaName))
	_, _ = conn.ExecContext(ctx, fmt.Sprintf(`CREATE SCHEMA %s`, schemaName))
	_, _ = conn.ExecContext(ctx, fmt.Sprintf(`SET search_path TO %s`, schemaName))

	// Apply DDL ke schema tersebut
	if _, err := conn.ExecContext(ctx, schema); err != nil {
		t.Fatalf("apply schema error: %v", err)
	}

	cleanup := func() {
		_, _ = conn.ExecContext(ctx, fmt.Sprintf(`DROP SCHEMA IF EXISTS %s CASCADE`, schemaName))
		_ = conn.Close()
		cleanupTestEnvVars()
	}

	return conn, cleanup
}

func setTestEnvVars() {
	testEnvVars := map[string]string{
		"DB_DRIVER":   "postgres",
		"DB_HOST":     "172.17.0.2",
		"DB_PORT":     "5432",
		"DB_NAME":     "example_test_payment_gateway",
		"DB_USERNAME": "postgres",
		"DB_PASSWORD": "password",
		"DB_SSLMODE":  "disable",
	}

	for key, value := range testEnvVars {
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
}

func cleanupTestEnvVars() {
	testEnvVars := []string{
		"DB_HOST", "DB_PORT", "DB_NAME", "DB_USER", "DB_PASSWORD", "DB_SSLMODE",
	}

	for _, key := range testEnvVars {
		os.Unsetenv(key)
	}
}

func CreateQueries(conn *sql.DB) *db.Queries {
	return db.New(conn)
}

func TestCtx() context.Context {
	return context.Background()
}

func InsertDummyUser(t *testing.T, ctx context.Context, queries *db.Queries) *db.User {
	user, err := queries.CreateUser(ctx, db.CreateUserParams{
		Firstname:        "Test",
		Lastname:         "User",
		Email:            "test@example.com",
		Password:         "pass",
		VerificationCode: "code123",
		IsVerified:       sql.NullBool{Bool: false, Valid: true},
	})
	if err != nil {
		t.Fatalf("insert dummy user: %v", err)
	}
	return user
}

func InsertDummyRefreshToken(t *testing.T, ctx context.Context, queries *db.Queries, userID int32) *db.RefreshToken {
	refreshToken, err := queries.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		UserID:     userID,
		Token:      "abc-token",
		Expiration: time.Now().Add(2 * time.Hour),
	})
	if err != nil {
		t.Fatalf("insert dummy refresh token: %v", err)
	}
	return refreshToken
}

func InsertDummyBulkData(t *testing.T, ctx context.Context, queries *db.Queries) (
	userIDs []int32,
	roleIDs []int32,
	refreshTokenIDs []int32,
) {

	roleNames := []string{"admin", "operator", "viewer"}
	for _, roleName := range roleNames {
		role, err := queries.CreateRole(ctx, roleName)
		if err != nil {
			t.Fatalf("failed to insert role %s: %v", roleName, err)
		}
		roleIDs = append(roleIDs, role.RoleID)
	}

	for i := 1; i <= 10; i++ {
		user, err := queries.CreateUser(ctx, db.CreateUserParams{
			Firstname:        fmt.Sprintf("First%d", i),
			Lastname:         fmt.Sprintf("Last%d", i),
			Email:            fmt.Sprintf("user%d@mail.com", i),
			Password:         fmt.Sprintf("pass%d", i),
			VerificationCode: fmt.Sprintf("verify%d", i),
			IsVerified:       sql.NullBool{Bool: true, Valid: true},
		})
		if err != nil {
			t.Fatalf("failed to insert dummy user %d: %v", i, err)
		}
		userIDs = append(userIDs, user.UserID)

		roleID := roleIDs[(i-1)%len(roleIDs)]
		_, err = queries.AssignRoleToUser(ctx, db.AssignRoleToUserParams{
			UserID: user.UserID,
			RoleID: roleID,
		})
		if err != nil {
			t.Fatalf("failed to assign role to user %d: %v", i, err)
		}

		exp := time.Now().Add(24 * time.Hour)
		rt, err := queries.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
			UserID:     user.UserID,
			Token:      fmt.Sprintf("refresh_token_%d", i),
			Expiration: exp,
		})
		if err != nil {
			t.Fatalf("failed to insert refresh token for user %d: %v", i, err)
		}
		refreshTokenIDs = append(refreshTokenIDs, rt.RefreshTokenID)
	}

	return userIDs, roleIDs, refreshTokenIDs
}
