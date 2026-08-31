package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"

	"sovera-core-api/internal/queue"
)

type WebhookHandler struct {
	dbPool      *pgxpool.Pool
	asynqClient *asynq.Client
}

func NewWebhookHandler(dbPool *pgxpool.Pool, asynqClient *asynq.Client) *WebhookHandler {
	return &WebhookHandler{
		dbPool:      dbPool,
		asynqClient: asynqClient,
	}
}

type CrawlerPayload struct {
	TaskID          string `json:"task_id"`
	SourceType      string `json:"source_type"`
	SourceURL       string `json:"source_url"`
	AuthorOrAccount string `json:"author_or_account"`
	PublishedDate   string `json:"published_date"`
	RawText         string `json:"raw_text"`
	MarkdownContent string `json:"markdown_content"`
}

func (h *WebhookHandler) HandleCrawlerWebhook(c *fiber.Ctx) error {
	var payload CrawlerPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "INVALID_PAYLOAD",
			"message": "Failed to parse JSON body",
		})
	}

	if payload.RawText == "" && payload.MarkdownContent == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"error":   "EMPTY_CONTENT",
			"message": "Payload raw_text or markdown_content cannot be empty",
		})
	}

	// Calculate SHA-256 content_hash for deduplication
	contentToHash := payload.RawText
	if contentToHash == "" {
		contentToHash = payload.MarkdownContent
	}
	hash := sha256.Sum256([]byte(contentToHash))
	contentHash := hex.EncodeToString(hash[:])

	jobID := "job_ingest_" + uuid.New().String()[:8]

	// Enqueue Asynq task into Redis
	taskPayload := queue.LLMExtractionPayload{
		JobID:           jobID,
		TaskID:          payload.TaskID,
		SourceType:      payload.SourceType,
		SourceURL:       payload.SourceURL,
		AuthorOrAccount: payload.AuthorOrAccount,
		PublishedDate:   payload.PublishedDate,
		RawText:         payload.RawText,
		MarkdownContent: payload.MarkdownContent,
		ContentHash:     contentHash,
	}

	if h.asynqClient != nil {
		task, err := queue.NewLLMExtractionTask(taskPayload)
		if err == nil {
			info, enqueueErr := h.asynqClient.Enqueue(task)
			if enqueueErr != nil {
				log.Printf("Failed to enqueue extraction task to Asynq: %v", enqueueErr)
			} else {
				log.Printf("Successfully enqueued task [%s] to queue [%s]", info.ID, info.Queue)
			}
		}
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"success":      true,
		"message":      "Payload queued for background ingestion & LLM extraction",
		"job_id":       jobID,
		"content_hash": contentHash,
		"queued_at":    time.Now().UTC().Format(time.RFC3339),
	})
}
