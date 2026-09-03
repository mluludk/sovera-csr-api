package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// VerifyHMAC creates a Fiber middleware that validates incoming webhooks using HMAC-SHA256 signature verification.
func VerifyHMAC(secretKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Bypass verification if secret key is intentionally empty
		if secretKey == "" {
			return c.Next()
		}

		signatureHeader := c.Get("X-Hub-Signature-256")
		webhookSecretHeader := c.Get("X-Webhook-Secret")

		// 1. Direct Secret Header Match (X-Webhook-Secret)
		if webhookSecretHeader != "" && webhookSecretHeader == secretKey {
			return c.Next()
		}

		// 2. Missing Signature Header Check
		if signatureHeader == "" {
			log.Printf("[HMAC Middleware Warning] 401 Unauthorized: Missing X-Hub-Signature-256 or X-Webhook-Secret headers from %s", c.IP())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "MISSING_SIGNATURE",
				"message": "Missing X-Hub-Signature-256 header",
			})
		}

		// Expect format: sha256=<hex_hash>
		parts := strings.SplitN(signatureHeader, "=", 2)
		if len(parts) != 2 || parts[0] != "sha256" {
			log.Printf("[HMAC Middleware Warning] 401 Unauthorized: Invalid signature format '%s'", signatureHeader)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "INVALID_SIGNATURE_FORMAT",
				"message": "Signature header must be in sha256=<hex> format",
			})
		}

		expectedHashHex := parts[1]

		// Calculate HMAC-SHA256 over raw request body
		mac := hmac.New(sha256.New, []byte(secretKey))
		mac.Write(c.Body())
		computedHashHex := hex.EncodeToString(mac.Sum(nil))

		// Constant-time compare to prevent timing side-channel attacks
		if subtle.ConstantTimeCompare([]byte(computedHashHex), []byte(expectedHashHex)) != 1 {
			log.Printf("[HMAC Middleware Warning] 401 Unauthorized: Signature mismatch (computed: %s, expected: %s)", computedHashHex, expectedHashHex)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "INVALID_SIGNATURE",
				"message": "HMAC signature verification failed",
			})
		}

		return c.Next()
	}
}
