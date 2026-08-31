package main

import (
	"context"
	"log"
	"time"

	"github.com/hibiken/asynq"

	"sovera-core-api/internal/config"
	"sovera-core-api/internal/queue"
	"sovera-core-api/internal/repository"
	"sovera-core-api/internal/service/ai"
)

func main() {
	cfg := config.LoadConfig()
	log.Printf("Starting Sovera Background Worker connected to Redis [%s]...", cfg.RedisURL)

	// 1. Database Pool Connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbPool, err := repository.InitDBPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Printf("Warning: Worker running without live DB pool: %v", err)
	} else {
		defer dbPool.Close()
	}

	// 2. Services & Worker Dependencies
	geminiService := ai.NewGeminiService(cfg.AIAPIKey)
	signalRepo := repository.NewSignalRepository(dbPool)
	extractionWorker := queue.NewExtractionWorker(geminiService, signalRepo)

	// 3. Asynq Server Setup
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.RedisURL},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				queue.QueueRawIngestion:       10,
				queue.QueueLLMExtraction:      5,
				queue.QueueProposalGeneration: 3,
			},
		},
	)

	mux := asynq.NewServeMux()

	// Register Asynq task handler for LLM extraction
	mux.HandleFunc(queue.TypeLLMExtraction, extractionWorker.ProcessExtractionTask)

	log.Println("Asynq Worker server listening for queue jobs...")
	if err := srv.Run(mux); err != nil {
		log.Fatalf("Could not run Asynq worker server: %v", err)
	}
}
