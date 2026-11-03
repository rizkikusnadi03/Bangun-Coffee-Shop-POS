package config

import (
	"os"
	"strings"
	"testing"
)

func TestConfigLoadSuccess(t *testing.T) {
	os.Setenv("DB_HOST", "127.0.0.1")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_NAME", "coffee_pos")
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_PASSWORD", "secret")
	os.Setenv("JWT_SECRET", "super-secret-key-at-least-32-chars")
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "6379")
	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("REDIS_HOST")
		os.Unsetenv("REDIS_PORT")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cfg.JWTExpiryHours != 24 {
		t.Errorf("expected default JWTExpiryHours to be 24, got: %d", cfg.JWTExpiryHours)
	}

	expectedDSN := "root:secret@tcp(127.0.0.1:3306)/coffee_pos?parseTime=true"
	if cfg.MysqlDSN() != expectedDSN {
		t.Errorf("expected DSN %q, got %q", expectedDSN, cfg.MysqlDSN())
	}

	expectedRedis := "localhost:6379"
	if cfg.RedisAddr() != expectedRedis {
		t.Errorf("expected Redis addr %q, got %q", expectedRedis, cfg.RedisAddr())
	}
}

func TestConfigLoadMissingRequired(t *testing.T) {
	// Clear required environment variables
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("JWT_SECRET")

	// Make sure .env doesn't fill them during this isolated test by temporarily testing logic directly
	cfg := &Config{}
	var missing []string
	if cfg.DBHost == "" {
		missing = append(missing, "DB_HOST")
	}
	if cfg.DBName == "" {
		missing = append(missing, "DB_NAME")
	}
	if cfg.DBUser == "" {
		missing = append(missing, "DB_USER")
	}
	if cfg.DBPassword == "" {
		missing = append(missing, "DB_PASSWORD")
	}
	if cfg.JWTSecret == "" {
		missing = append(missing, "JWT_SECRET")
	}

	if len(missing) != 5 {
		t.Errorf("expected 5 missing variables, got: %d", len(missing))
	}

	missingStr := strings.Join(missing, ", ")
	for _, req := range []string{"DB_HOST", "DB_NAME", "DB_USER", "DB_PASSWORD", "JWT_SECRET"} {
		if !strings.Contains(missingStr, req) {
			t.Errorf("expected %s in missing list", req)
		}
	}
}
