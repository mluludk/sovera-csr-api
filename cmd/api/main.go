package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/hibiken/asynq"

	"sovera-core-api/internal/config"
	"sovera-core-api/internal/handler"
	"sovera-core-api/internal/middleware"
	"sovera-core-api/internal/repository"
)

func main() {
	// 1. Load Environment Configuration
	cfg := config.LoadConfig()
	log.Printf("Starting Sovera Core API Server in [%s] mode...", cfg.Environment)

	// 2. Initialize PostgreSQL DB Pool with pgx/v5
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbPool, err := repository.InitDBPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Printf("Warning: Database pool initialization failed: %v", err)
	} else {
		log.Println("PostgreSQL connection pool initialized successfully.")
		defer dbPool.Close()

		// Run SQL DDL Migrations
		if err := repository.RunMigrations(context.Background(), dbPool, "db/migrations"); err != nil {
			log.Printf("Warning: Automatic migration run failed: %v", err)
		} else {
			log.Println("Database DDL migrations executed successfully.")
		}
	}

	// 3. Initialize Asynq Redis Client
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.RedisURL})
	defer asynqClient.Close()

	// 4. Create Fiber Web Application
	app := fiber.New(fiber.Config{
		AppName:      "Sovera Core API & Intelligence Engine v1.0",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	})

	// 5. Global Middlewares
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))

	// 6. Register Handlers & Controllers
	healthHandler := handler.NewHealthHandler(dbPool)
	webhookHandler := handler.NewWebhookHandler(dbPool, asynqClient)

	// Root & Health check routes
	app.Get("/health", healthHandler.HealthCheck)

	apiV1 := app.Group("/api/v1")
	apiV1.Get("/health", healthHandler.HealthCheck)

	// Webhook Ingestion Endpoint (Protected by HMAC Verification)
	apiV1.Post(
		"/webhooks/crawler",
		middleware.VerifyHMAC(cfg.WebhookSecretKey),
		webhookHandler.HandleCrawlerWebhook,
	)

	// 7. Graceful Shutdown Handler
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		addr := fmt.Sprintf(":%s", cfg.Port)
		log.Printf("Sovera API Listening on http://localhost%s/api/v1", addr)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Fiber server error: %v", err)
		}
	}()

	<-shutdownChan
	log.Println("Gracefully shutting down Sovera API server...")

	_ = app.Shutdown()
	log.Println("Server stopped cleanly.")
}
