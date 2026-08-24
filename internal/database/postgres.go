package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Connection struct {
	GORM *gorm.DB
	sql  *sql.DB
}

func Open(ctx context.Context, dsn string) (*Connection, error) {
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{DisableAutomaticPing: true})
	if err != nil {
		return nil, fmt.Errorf("open PostgreSQL connection: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("get PostgreSQL connection pool: %w", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("connect to PostgreSQL: %w", err)
	}

	return &Connection{GORM: gormDB, sql: sqlDB}, nil
}

func (c *Connection) Ping(ctx context.Context) error {
	return c.sql.PingContext(ctx)
}

func (c *Connection) Close() error {
	return c.sql.Close()
}
