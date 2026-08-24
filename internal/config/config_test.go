package config

import (
	"net/url"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	t.Setenv("PORT", "9000")
	t.Setenv("DB_HOST", "database")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("DB_USER", "app-user")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_NAME", "app-db")
	t.Setenv("DB_SSLMODE", "require")
	t.Setenv("FIREBASE_PROJECT_ID", "test-project")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned an error: %v", err)
	}

	if cfg.ServerPort != "9000" {
		t.Errorf("ServerPort = %q, want %q", cfg.ServerPort, "9000")
	}
	if cfg.DatabaseHost != "database" {
		t.Errorf("DatabaseHost = %q, want %q", cfg.DatabaseHost, "database")
	}
	if cfg.DatabaseSSLMode != "require" {
		t.Errorf("DatabaseSSLMode = %q, want %q", cfg.DatabaseSSLMode, "require")
	}
	if cfg.FirebaseProjectID != "test-project" {
		t.Errorf(
			"FirebaseProjectID = %q, want %q",
			cfg.FirebaseProjectID,
			"test-project",
		)
	}
}

func TestLoadReportsMissingDatabaseSettings(t *testing.T) {
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("DB_NAME", "")
	t.Setenv("FIREBASE_PROJECT_ID", "")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() returned nil error, want missing environment variables error")
	}
	for _, name := range []string{"DB_HOST", "DB_USER", "DB_PASSWORD", "DB_NAME", "FIREBASE_PROJECT_ID"} {
		if !strings.Contains(err.Error(), name) {
			t.Errorf("Load() error = %q, want it to contain %q", err, name)
		}
	}
}

func TestDatabaseDSNEncodesCredentials(t *testing.T) {
	cfg := Config{
		DatabaseHost:     "db",
		DatabasePort:     "5432",
		DatabaseUser:     "app user",
		DatabasePassword: "p@ss/word",
		DatabaseName:     "app",
		DatabaseSSLMode:  "disable",
	}

	parsed, err := url.Parse(cfg.DatabaseDSN())
	if err != nil {
		t.Fatalf("DatabaseDSN() returned an invalid URL: %v", err)
	}
	password, ok := parsed.User.Password()
	if !ok {
		t.Fatal("DatabaseDSN() does not contain a password")
	}
	if parsed.User.Username() != cfg.DatabaseUser || password != cfg.DatabasePassword {
		t.Errorf("DatabaseDSN() credentials = %q:%q, want %q:%q", parsed.User.Username(), password, cfg.DatabaseUser, cfg.DatabasePassword)
	}
}
