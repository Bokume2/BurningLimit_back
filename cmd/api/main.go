package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Bokume2/FirstHackathon2026Summer_back/internal/config"
	"github.com/Bokume2/FirstHackathon2026Summer_back/internal/database"
	"github.com/labstack/echo/v5"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	connectCtx, cancelConnect := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelConnect()

	db, err := database.Open(connectCtx, cfg.DatabaseDSN())
	if err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("close database: %v", err)
		}
	}()

	e := echo.New()
	e.GET("/health", func(c *echo.Context) error {
		pingCtx, cancelPing := context.WithTimeout(c.Request().Context(), 2*time.Second)
		defer cancelPing()

		if err := db.Ping(pingCtx); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "unhealthy"})
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	server := echo.StartConfig{
		Address:         ":" + cfg.ServerPort,
		GracefulTimeout: 10 * time.Second,
	}
	if err := server.Start(ctx, e); err != nil {
		return fmt.Errorf("start HTTP server: %w", err)
	}

	return nil
}
