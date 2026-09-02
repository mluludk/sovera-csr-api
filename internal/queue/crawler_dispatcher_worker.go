package queue

import (
	"context"
	"fmt"
	"log"
	"time"

	"sovera-core-api/internal/model"
	"sovera-core-api/internal/repository"
	"sovera-core-api/internal/service/crawler"

	"github.com/hibiken/asynq"
)

type CrawlerDispatcherHandler struct {
	crawlerRepo *repository.CrawlerRepository
	dispatcher  *crawler.Dispatcher
}

func NewCrawlerDispatcherHandler(crawlerRepo *repository.CrawlerRepository, dispatcher *crawler.Dispatcher) *CrawlerDispatcherHandler {
	return &CrawlerDispatcherHandler{
		crawlerRepo: crawlerRepo,
		dispatcher:  dispatcher,
	}
}

func (h *CrawlerDispatcherHandler) HandleDispatchCrawlingTask(ctx context.Context, t *asynq.Task) error {
	log.Println("Starting execution of periodic crawling target dispatching worker...")

	dueTargets, err := h.crawlerRepo.GetDueTargets(ctx, 20)
	if err != nil {
		return fmt.Errorf("failed to fetch due crawling targets: %w", err)
	}

	if len(dueTargets) == 0 {
		log.Println("No due crawling targets found. Skipping dispatching cycle.")
		return nil
	}

	log.Printf("Found %d due crawling targets. Dispatching to Scraper Service...", len(dueTargets))

	for _, target := range dueTargets {
		taskID := fmt.Sprintf("task_%s_%d", target.ID[:8], time.Now().Unix())

		// Create dispatch log in database
		initialLog := model.CrawlingLog{
			TargetID: &target.ID,
			TaskID:   taskID,
			Status:   "DISPATCHED",
		}
		if err := h.crawlerRepo.CreateLog(ctx, initialLog); err != nil {
			log.Printf("Error creating crawling log for target %s: %v", target.ID, err)
		}

		// Dispatch scrape task via HTTP
		statusCode, err := h.dispatcher.DispatchTask(ctx, target, taskID)
		if err != nil {
			log.Printf("Failed to dispatch target %s (TaskID: %s): %v", target.ID, taskID, err)
			errMsg := err.Error()
			_ = h.crawlerRepo.UpdateLogStatus(ctx, taskID, "FAILED", nil, nil, &errMsg)
			continue
		}

		// Update target next_run_at in database
		if err := h.crawlerRepo.UpdateTargetNextRun(ctx, target.ID, target.CheckIntervalHours); err != nil {
			log.Printf("Failed to update target next run time for %s: %v", target.ID, err)
		}

		log.Printf("Successfully dispatched target '%s' (TaskID: %s, HTTP %d)", target.SourceName, taskID, statusCode)
	}

	log.Println("Crawling dispatching cycle completed successfully.")
	return nil
}
