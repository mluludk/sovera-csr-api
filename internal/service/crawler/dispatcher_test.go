package crawler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"sovera-core-api/internal/config"
	"sovera-core-api/internal/model"
)

func TestDispatcher_DispatchTask_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"status":"ACCEPTED"}`))
	}))
	defer ts.Close()

	cfg := &config.Config{
		ScraperServiceURL: ts.URL,
		WebhookURL:        "http://localhost:4000/api/v1/webhooks/crawler",
		WebhookSecretKey:  "secret_key_123",
	}

	dispatcher := NewDispatcher(cfg)
	target := model.CrawlingTarget{
		ID:         "target_12345678",
		SourceType: "PDF_DOCUMENT",
		TargetURL:  "https://example.com/report.pdf",
	}

	statusCode, err := dispatcher.DispatchTask(context.Background(), target, "task_123")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if statusCode != http.StatusAccepted {
		t.Errorf("expected status code 202, got %d", statusCode)
	}
}
