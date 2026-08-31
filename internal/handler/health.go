package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	dbPool *pgxpool.Pool
}

func NewHealthHandler(dbPool *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{dbPool: dbPool}
}

func (h *HealthHandler) HealthCheck(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	dbStatus := "UP"
	if err := h.dbPool.Ping(ctx); err != nil {
		dbStatus = "DOWN"
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "UP",
		"service":   "Sovera Core API",
		"version":   "1.0.0",
		"database":  dbStatus,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
