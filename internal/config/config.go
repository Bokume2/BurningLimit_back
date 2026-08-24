package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort        string
	DatabaseHost      string
	DatabasePort      string
	DatabaseUser      string
	DatabasePassword  string
	DatabaseName      string
	DatabaseSSLMode   string
	FirebaseProjectID string
}

func Load() (Config, error) {
	cfg := Config{
		ServerPort:        valueOrDefault("PORT", "8080"),
		DatabaseHost:      strings.TrimSpace(os.Getenv("DB_HOST")),
		DatabasePort:      valueOrDefault("DB_PORT", "5432"),
		DatabaseUser:      strings.TrimSpace(os.Getenv("DB_USER")),
		DatabasePassword:  os.Getenv("DB_PASSWORD"),
		DatabaseName:      strings.TrimSpace(os.Getenv("DB_NAME")),
		DatabaseSSLMode:   valueOrDefault("DB_SSLMODE", "disable"),
		FirebaseProjectID: strings.TrimSpace(os.Getenv("FIREBASE_PROJECT_ID")),
	}

	missing := make([]string, 0, 5)
	for name, value := range map[string]string{
		"DB_HOST":             cfg.DatabaseHost,
		"DB_USER":             cfg.DatabaseUser,
		"DB_PASSWORD":         cfg.DatabasePassword,
		"DB_NAME":             cfg.DatabaseName,
		"FIREBASE_PROJECT_ID": cfg.FirebaseProjectID,
	} {
		if value == "" {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return Config{}, fmt.Errorf("required environment variables are not set: %s", strings.Join(missing, ", "))
	}

	if err := validatePort("PORT", cfg.ServerPort); err != nil {
		return Config{}, err
	}
	if err := validatePort("DB_PORT", cfg.DatabasePort); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) DatabaseDSN() string {
	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.DatabaseUser, c.DatabasePassword),
		Host:   net.JoinHostPort(c.DatabaseHost, c.DatabasePort),
		Path:   c.DatabaseName,
	}
	query := dsn.Query()
	query.Set("sslmode", c.DatabaseSSLMode)
	dsn.RawQuery = query.Encode()

	return dsn.String()
}

func valueOrDefault(name, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value
	}
	return fallback
}

func validatePort(name, value string) error {
	port, err := strconv.Atoi(value)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("%s must be a number between 1 and 65535", name)
	}
	return nil
}
