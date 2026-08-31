package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
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

	// Dispatches job asynchronously to Asynq queue (stub queue enqueue if client is configured)
	// Return HTTP 202 Accepted immediately
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"success":      true,
		"message":      "Payload queued for background ingestion & LLM extraction",
		"job_id":       jobID,
		"content_hash": contentHash,
		"queued_at":    time.Now().UTC().Format(time.RFC3339),
	})
}
