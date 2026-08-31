package main

import (
	"context"
	"log"

	"github.com/hibiken/asynq"

	"sovera-core-api/internal/config"
)

func main() {
	cfg := config.LoadConfig()
	log.Printf("Starting Sovera Background Worker connected to Redis [%s]...", cfg.RedisURL)

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.RedisURL},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"raw-ingestion-queue":       10,
				"llm-extraction-queue":      5,
				"proposal-generation-queue": 3,
			},
		},
	)

	mux := asynq.NewServeMux()

	// Handler stubs for Asynq task types
	mux.HandleFunc("task:raw_ingestion", func(ctx context.Context, task *asynq.Task) error {
		log.Printf("Processing raw ingestion task: %s", task.Type())
		return nil
	})

	log.Println("Asynq Worker server listening for queue jobs...")
	if err := srv.Run(mux); err != nil {
		log.Fatalf("Could not run Asynq worker server: %v", err)
	}
}
