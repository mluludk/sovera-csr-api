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
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/hibiken/asynq"

	"sovera-core-api/internal/config"
	"sovera-core-api/internal/handler"
	"sovera-core-api/internal/middleware"
	"sovera-core-api/internal/repository"
	"sovera-core-api/internal/service/ai"
	"sovera-core-api/internal/service/exporter"
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

	// 4. Initialize Services & Repositories
	geminiService := ai.NewGeminiService(cfg.AIAPIKey)
	docExporter := exporter.NewDocumentExporter()

	signalRepo := repository.NewSignalRepository(dbPool)
	programRepo := repository.NewProgramRepository(dbPool)
	dealRepo := repository.NewDealRepository(dbPool)

	// 5. Create Fiber Web Application
	app := fiber.New(fiber.Config{
		AppName:      "Sovera Core API & Intelligence Engine v1.0",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	})

	// 6. Global Middlewares (CORS MUST BE FIRST)
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		AllowMethods:     "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS",
		AllowCredentials: false,
	}))
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))

	// 7. Register Handlers & Controllers
	healthHandler := handler.NewHealthHandler(dbPool)
	webhookHandler := handler.NewWebhookHandler(dbPool, asynqClient)
	signalHandler := handler.NewSignalHandler(signalRepo)
	programHandler := handler.NewProgramHandler(programRepo, geminiService)
	dealHandler := handler.NewDealHandler(dealRepo, programRepo, signalRepo, geminiService, docExporter)

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

	// Corporate Intelligence Feeds API
	apiV1.Get("/signals", signalHandler.ListSignals)
	apiV1.Get("/signals/:id/match-programs", signalHandler.MatchPrograms)

	// Institution Programs API (RLS Enforced)
	apiV1.Get("/programs", programHandler.ListPrograms)
	apiV1.Post("/programs", programHandler.CreateProgram)

	// Deal Pipeline & Proposal Studio API (RLS Enforced)
	apiV1.Get("/deals", dealHandler.ListDeals)
	apiV1.Post("/deals", dealHandler.CreateDeal)
	apiV1.Patch("/deals/:id/stage", dealHandler.UpdateStage)
	apiV1.Post("/deals/:id/generate-pitch", dealHandler.GeneratePitch)
	apiV1.Post("/deals/:id/export", dealHandler.ExportProposal)

	// 8. Graceful Shutdown Handler
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
