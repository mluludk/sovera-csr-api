package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"

	"sovera-core-api/internal/repository"
	"sovera-core-api/internal/service/ai"
)

type ExtractionWorker struct {
	geminiService *ai.GeminiService
	signalRepo    *repository.SignalRepository
}

func NewExtractionWorker(geminiService *ai.GeminiService, signalRepo *repository.SignalRepository) *ExtractionWorker {
	return &ExtractionWorker{
		geminiService: geminiService,
		signalRepo:    signalRepo,
	}
}

func (w *ExtractionWorker) ProcessExtractionTask(ctx context.Context, task *asynq.Task) error {
	var payload LLMExtractionPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal LLM extraction task payload: %w", err)
	}

	log.Printf("[Asynq Worker] Starting LLM Entity Extraction for Job [%s] Content Hash [%s]", payload.JobID, payload.ContentHash)

	textToExtract := payload.RawText
	if textToExtract == "" {
		textToExtract = payload.MarkdownContent
	}

	// 1. LLM Structured Entity Extraction via Gemini
	extractedSignal, err := w.geminiService.ExtractCorporateSignal(ctx, textToExtract)
	if err != nil {
		return fmt.Errorf("LLM signal extraction failed: %w", err)
	}

	// 2. Vector Embedding Generation (1536 dim)
	textToEmbed := fmt.Sprintf("%s %s %s %s", extractedSignal.CompanyName, extractedSignal.IndustrySector, extractedSignal.CSRPillarFocus, extractedSignal.Summary)
	embedding, err := w.geminiService.GenerateEmbedding(ctx, textToEmbed)
	if err != nil {
		return fmt.Errorf("vector embedding generation failed: %w", err)
	}

	// 3. Persist Extracted Signal & Vector to Database
	signalID, err := w.signalRepo.SaveSignal(ctx, extractedSignal, payload.SourceType, payload.SourceURL, payload.ContentHash, embedding)
	if err != nil {
		return fmt.Errorf("database signal persistence failed: %w", err)
	}

	log.Printf("[Asynq Worker] Successfully saved Signal [%s] for company [%s] with intent score [%d]", signalID, extractedSignal.CompanyName, extractedSignal.IntentScore)
	return nil
}
