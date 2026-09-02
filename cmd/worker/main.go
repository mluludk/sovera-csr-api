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
	"sovera-core-api/internal/service/crawler"
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
	crawlerRepo := repository.NewCrawlerRepository(dbPool)
	dispatcher := crawler.NewDispatcher(cfg)

	extractionWorker := queue.NewExtractionWorker(geminiService, signalRepo)
	dispatcherWorker := queue.NewCrawlerDispatcherHandler(crawlerRepo, dispatcher)

	// 3. Asynq Scheduler for Periodic Tasks
	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{Addr: cfg.RedisURL},
		&asynq.SchedulerOpts{},
	)

	// Schedule task:dispatch_crawling
	dispatchTask, err := queue.NewDispatchCrawlingTask()
	if err == nil {
		// Instant startup sweep: Enqueue 1x dispatch task immediately on worker launch/restart
		startupClient := asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.RedisURL})
		if info, enqueueErr := startupClient.Enqueue(dispatchTask); enqueueErr != nil {
			log.Printf("Notice: Could not enqueue instant startup crawling task: %v", enqueueErr)
		} else {
			log.Printf("Enqueued instant startup crawling target check to Redis queue (TaskID: %s)", info.ID)
		}
		startupClient.Close()

		// Register periodic cron to run every 1 hour (cron: "0 * * * *")
		if entryID, err := scheduler.Register("0 * * * *", dispatchTask); err != nil {
			log.Printf("Warning: Could not register periodic crawling dispatch cron: %v", err)
		} else {
			log.Printf("Registered periodic crawling dispatch cron with entry ID: %s", entryID)
		}
	}

	go func() {
		if err := scheduler.Run(); err != nil {
			log.Printf("Scheduler error: %v", err)
		}
	}()

	// 4. Asynq Server Setup
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.RedisURL},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				queue.QueueDispatchCrawling:   10,
				queue.QueueRawIngestion:       10,
				queue.QueueLLMExtraction:      5,
				queue.QueueProposalGeneration: 3,
			},
		},
	)

	mux := asynq.NewServeMux()

	// Register Asynq task handlers
	mux.HandleFunc(queue.TypeDispatchCrawling, dispatcherWorker.HandleDispatchCrawlingTask)
	mux.HandleFunc(queue.TypeLLMExtraction, extractionWorker.ProcessExtractionTask)

	log.Println("Asynq Worker server listening for queue jobs...")
	if err := srv.Run(mux); err != nil {
		log.Fatalf("Could not run Asynq worker server: %v", err)
	}
}
