package queue

import (
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

const (
	TypeDispatchCrawling   = "task:dispatch_crawling"
	TypeRawIngestion       = "task:raw_ingestion"
	TypeLLMExtraction      = "task:llm_extraction"
	TypeProposalGeneration = "task:proposal_generation"

	QueueDispatchCrawling   = "dispatch-crawling-queue"
	QueueRawIngestion       = "raw-ingestion-queue"
	QueueLLMExtraction      = "llm-extraction-queue"
	QueueProposalGeneration = "proposal-generation-queue"
)

type LLMExtractionPayload struct {
	JobID           string `json:"job_id"`
	TaskID          string `json:"task_id"`
	SourceType      string `json:"source_type"`
	SourceURL       string `json:"source_url"`
	AuthorOrAccount string `json:"author_or_account"`
	PublishedDate   string `json:"published_date"`
	RawText         string `json:"raw_text"`
	MarkdownContent string `json:"markdown_content"`
	ContentHash     string `json:"content_hash"`
}

func NewDispatchCrawlingTask() (*asynq.Task, error) {
	return asynq.NewTask(TypeDispatchCrawling, []byte("{}"), asynq.Queue(QueueDispatchCrawling), asynq.MaxRetry(3)), nil
}

func NewLLMExtractionTask(payload LLMExtractionPayload) (*asynq.Task, error) {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal extraction task payload: %w", err)
	}

	return asynq.NewTask(TypeLLMExtraction, bytes, asynq.Queue(QueueLLMExtraction), asynq.MaxRetry(5)), nil
}
