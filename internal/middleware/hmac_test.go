package middleware_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"

	"sovera-core-api/internal/middleware"
)

func TestVerifyHMAC_ValidSignature(t *testing.T) {
	secret := "test_secret_key_123"
	app := fiber.New()
	app.Post("/webhook", middleware.VerifyHMAC(secret), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	body := `{"task_id":"job_123","raw_text":"PT Maju Bersama CSR"}`
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(body))
	sigHex := hex.EncodeToString(mac.Sum(nil))

	req := httptest.NewRequest("POST", "/webhook", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Hub-Signature-256", fmt.Sprintf("sha256=%s", sigHex))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestVerifyHMAC_InvalidSignature(t *testing.T) {
	secret := "test_secret_key_123"
	app := fiber.New()
	app.Post("/webhook", middleware.VerifyHMAC(secret), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	body := `{"task_id":"job_123","raw_text":"PT Maju Bersama CSR"}`

	req := httptest.NewRequest("POST", "/webhook", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Hub-Signature-256", "sha256=invalid_hash_string")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status 401 Unauthorized, got %d", resp.StatusCode)
	}
}
