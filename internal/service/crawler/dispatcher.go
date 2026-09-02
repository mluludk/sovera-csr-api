package crawler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"sovera-core-api/internal/config"
	"sovera-core-api/internal/model"
)

type Dispatcher struct {
	cfg        *config.Config
	httpClient *http.Client
}

func NewDispatcher(cfg *config.Config) *Dispatcher {
	return &Dispatcher{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// DispatchTask sends a scrape task request to Scraper Service
func (d *Dispatcher) DispatchTask(ctx context.Context, target model.CrawlingTarget, taskID string) (int, error) {
	if d.cfg.ScraperServiceURL == "" {
		return 0, fmt.Errorf("SCRAPER_SERVICE_URL is not configured")
	}

	payload := model.ScrapeTaskPayload{
		TaskID:       taskID,
		TargetID:     target.ID,
		ClientOrigin: "sovera_b2b_engine",
		SourceType:   target.SourceType,
		TargetURL:    target.TargetURL,
		CallbackURL:  d.cfg.WebhookURL,
		Config: &model.ScrapeTaskConfig{
			RenderJS:      false,
			BypassAntiBot: true,
			MaxPages:      50,
		},
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.cfg.ScraperServiceURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return 0, fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	apiKey := d.cfg.ScraperAPIKey
	if apiKey == "" {
		apiKey = d.cfg.WebhookSecretKey
	}
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to execute http post to scraper service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return resp.StatusCode, fmt.Errorf("scraper service returned non-2xx status: %d", resp.StatusCode)
	}

	return resp.StatusCode, nil
}
